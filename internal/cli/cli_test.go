package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type badWriter struct{}

func (badWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestCommands(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		code    int
		outPart string
		errPart string
	}{
		{"root", nil, ExitUnavailable, "", "validation is not implemented"},
		{"init", []string{"init"}, ExitUnavailable, "", "no hooks or configuration were written"},
		{"status", []string{"status"}, ExitUnavailable, "", "no validation result is available"},
		{"help", []string{"--help"}, ExitSuccess, "Usage: ci-see", ""},
		{"short help", []string{"-h"}, ExitSuccess, "Usage: ci-see", ""},
		{"help command", []string{"help"}, ExitSuccess, "Usage: ci-see", ""},
		{"help init", []string{"help", "init"}, ExitSuccess, "Usage: ci-see init", ""},
		{"help status", []string{"help", "status"}, ExitSuccess, "Usage: ci-see status", ""},
		{"init help", []string{"init", "--help"}, ExitSuccess, "Usage: ci-see init", ""},
		{"status help", []string{"status", "-h"}, ExitSuccess, "Usage: ci-see status", ""},
		{"version", []string{"--version"}, ExitSuccess, "ci-see v-test (commit abc; go-test)", ""},
		{"version command", []string{"version"}, ExitSuccess, "ci-see v-test", ""},
		{"invalid subcommand", []string{"magic"}, ExitUsage, "", "invalid command"},
		{"invalid flag", []string{"--force"}, ExitUsage, "", "invalid command"},
		{"bypass", []string{"--bypass"}, ExitUsage, "", "invalid command"},
		{"forged proof", []string{"--verified=true"}, ExitUsage, "", "invalid command"},
		{"fingerprint", []string{"--fingerprint=deadbeef"}, ExitUsage, "", "invalid command"},
		{"init bypass", []string{"init", "--force"}, ExitUsage, "", "invalid command"},
		{"version invalid", []string{"init", "--version"}, ExitUsage, "", "invalid command"},
		{"extra help args", []string{"--help", "init"}, ExitUsage, "", "invalid command"},
		{"extra root args", []string{"status", "junk"}, ExitUsage, "", "invalid command"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := Run(test.args, &out, &errOut, BuildInfo{Version: "v-test", Commit: "abc", Go: "go-test"})
			if code != test.code {
				t.Errorf("code = %d, want %d", code, test.code)
			}
			if !strings.Contains(out.String(), test.outPart) || !strings.Contains(errOut.String(), test.errPart) {
				t.Errorf("stdout=%q stderr=%q", out.String(), errOut.String())
			}
			if test.code == ExitSuccess && errOut.Len() != 0 {
				t.Errorf("success wrote stderr: %q", errOut.String())
			}
			if test.code != ExitSuccess && out.Len() != 0 {
				t.Errorf("failed command wrote stdout: %q", out.String())
			}
		})
	}
}

func TestClaimedApprovalEnvironmentHasNoEffect(t *testing.T) {
	t.Setenv("CI_SEE_FORCE", "true")
	t.Setenv("CI_SEE_VERIFIED", "true")
	t.Setenv("CI_SEE_APPROVED", "true")
	for _, args := range [][]string{nil, {"init"}, {"status"}, {"--approved"}} {
		var out, errOut bytes.Buffer
		if got := Run(args, &out, &errOut, BuildInfo{}); got == ExitSuccess {
			t.Fatalf("env or args authorized %v", args)
		}
	}
}

func TestWriteFailuresFailClosed(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"--version"}} {
		if got := Run(args, badWriter{}, &bytes.Buffer{}, BuildInfo{}); got != ExitInternal {
			t.Errorf("write failure = %d", got)
		}
	}
	for _, args := range [][]string{nil, {"init"}, {"bogus"}} {
		if got := Run(args, &bytes.Buffer{}, badWriter{}, BuildInfo{}); got != ExitInternal {
			t.Errorf("write failure = %d", got)
		}
	}
	if got := Run(nil, nil, &bytes.Buffer{}, BuildInfo{}); got != ExitInternal {
		t.Errorf("nil stdout = %d", got)
	}
	if got := Run(nil, &bytes.Buffer{}, nil, BuildInfo{}); got != ExitInternal {
		t.Errorf("nil stderr = %d", got)
	}
}

func TestVersionDefaults(t *testing.T) {
	var out, errOut bytes.Buffer
	if got := Run([]string{"--version"}, &out, &errOut, BuildInfo{}); got != ExitSuccess {
		t.Fatalf("version = %d", got)
	}
	if !strings.Contains(out.String(), "ci-see dev (commit unknown; unknown)") {
		t.Fatalf("unexpected default version: %q", out.String())
	}
}

func TestInitDoesNotModifyTemporaryRepository(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0700); err != nil {
		t.Fatal(err)
	}
	hook := filepath.Join(hooks, "pre-commit")
	if err := os.WriteFile(hook, []byte("user hook\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if got := Run([]string{"init"}, &out, &errOut, BuildInfo{}); got != ExitUnavailable {
		t.Fatalf("init status = %d", got)
	}
	data, err := os.ReadFile(hook)
	if err != nil || string(data) != "user hook\n" {
		t.Fatalf("user hook modified: data=%q err=%v", data, err)
	}
	entries, err := os.ReadDir(hooks)
	if err != nil || len(entries) != 1 {
		t.Fatalf("unexpected hooks: %v %v", entries, err)
	}
}
