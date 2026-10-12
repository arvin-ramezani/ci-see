// Command d3-preflight validates signed D3 run manifests offline. It cannot redeem or execute them.
package main

import (
	"crypto/ed25519"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arvin-ramezani/ci-see/internal/d3preflight"
)

const (
	exitOK         = 0
	exitInternal   = 1
	exitUsage      = 2
	exitIncomplete = 3
	exitBlocked    = 4
)

func main() {
	// No owner-approved trust pin or trusted clock is configured in this build.
	// Tests inject disposable fixture keys and a deterministic clock directly.
	os.Exit(runDefault(os.Args[1:], os.Stdout, os.Stderr))
}

func runDefault(args []string, stdout, stderr io.Writer) int {
	keys := map[string]ed25519.PublicKey{}
	return run(args, stdout, stderr, keys, d3preflight.Clock{})
}

func run(args []string, stdout, stderr io.Writer, keys map[string]ed25519.PublicKey, clock d3preflight.Clock) int {
	if stdout == nil || stderr == nil {
		return exitInternal
	}
	if hasDuplicateFlags(args) {
		return printResult(stderr, "INCOMPLETE", "duplicate_argument", exitUsage)
	}
	flags := flag.NewFlagSet("d3-preflight", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	manifestPath := flags.String("manifest", "", "manifest JSON path")
	approvalPath := flags.String("approval", "", "detached approval JSON path")
	repository := flags.String("repository", "", "expected owner/repository")
	pullRequest := flags.Int("pr", 0, "expected pull request number")
	headSHA := flags.String("head-sha", "", "expected commit SHA")
	treeSHA := flags.String("tree-sha", "", "expected tree SHA")
	runID := flags.String("run-id", "", "expected run UUID")
	hostID := flags.String("host-id", "", "expected host identity")
	vmID := flags.String("vm-id", "", "expected VM identity")
	snapshotID := flags.String("snapshot-id", "", "expected snapshot identity")
	guestID := flags.String("guest-id", "", "expected guest identity")
	leg := flags.String("leg", "", "expected execution leg (windows-native or wsl2)")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return printResult(stderr, "INCOMPLETE", "arguments_invalid", exitUsage)
	}
	if *manifestPath == "" || *approvalPath == "" || *repository == "" || *pullRequest < 1 || *headSHA == "" || *treeSHA == "" ||
		*runID == "" || *hostID == "" || *vmID == "" || *snapshotID == "" || *guestID == "" || *leg == "" {
		return printResult(stderr, string(d3preflight.StatusIncomplete), "expected_scope_or_input_missing", exitIncomplete)
	}
	manifestBytes, ok := readBounded(*manifestPath, d3preflight.MaxManifest)
	if !ok {
		return printResult(stderr, string(d3preflight.StatusIncomplete), "manifest_unavailable", exitIncomplete)
	}
	approvalBytes, ok := readBounded(*approvalPath, d3preflight.MaxApproval)
	if !ok {
		return printResult(stderr, string(d3preflight.StatusIncomplete), "approval_unavailable", exitIncomplete)
	}
	result := d3preflight.Validate(d3preflight.Input{
		ManifestBytes: manifestBytes,
		ApprovalBytes: approvalBytes,
		TrustedKeys:   keys,
		Expected: d3preflight.Scope{
			Repository: *repository, PullRequest: *pullRequest, HeadSHA: *headSHA, TreeSHA: *treeSHA, RunID: *runID,
			HostID: *hostID, VMID: *vmID, SnapshotID: *snapshotID, GuestID: *guestID, ExecutionLeg: *leg,
		},
		Clock: clock,
	})
	if result.Status != d3preflight.StatusConfigValid {
		code := exitBlocked
		if result.Status == d3preflight.StatusIncomplete {
			code = exitIncomplete
		}
		return printResult(stderr, string(result.Status), result.Code, code)
	}
	if _, err := fmt.Fprintln(stdout, "SIGNATURE_VALID"); err != nil {
		return exitInternal
	}
	if _, err := fmt.Fprintln(stdout, "CONFIG_VALID"); err != nil {
		return exitInternal
	}
	if _, err := fmt.Fprintln(stdout, "execution_authorized=false approval_redeemed=false"); err != nil {
		return exitInternal
	}
	return exitOK
}

func hasDuplicateFlags(args []string) bool {
	seen := make(map[string]struct{}, len(args))
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		name, _, _ := strings.Cut(strings.TrimPrefix(arg, "--"), "=")
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			return true
		}
		seen[name] = struct{}{}
	}
	return false
}

func readBounded(path string, limit int) ([]byte, bool) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, int64(limit)+1))
	if err != nil || len(data) > limit {
		return nil, false
	}
	return data, true
}

func printResult(writer io.Writer, status, code string, exit int) int {
	if _, err := fmt.Fprintf(writer, "d3-preflight: %s code=%s execution_authorized=false approval_redeemed=false\n", status, code); err != nil {
		return exitInternal
	}
	return exit
}

