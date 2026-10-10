// Package cli implements the S1 command skeleton. No command performs CI or Git writes.
package cli

import (
	"fmt"
	"io"
)

const (
	ExitSuccess     = 0 // Help/version only in S1.
	ExitInternal    = 1
	ExitUsage       = 2
	ExitUnavailable = 3
	ExitBlocked     = 4 // Reserved for future Git gate dispatch.
)

// BuildInfo is injected by main so command behavior can be tested in-process.
type BuildInfo struct {
	Version string
	Commit  string
	Go      string
}

const usage = "Usage: ci-see [--help|--version|init|status]\n" +
	"  ci-see         Run local validation (not implemented in S1)\n" +
	"  ci-see init    Configure repository hooks (not implemented in S1)\n" +
	"  ci-see status  Read validation result (not implemented in S1)\n"

// Run returns an explicit exit code. It never exits the process or trusts caller input
// as validation evidence or developer approval.
func Run(args []string, stdout, stderr io.Writer, build BuildInfo) int {
	if stdout == nil || stderr == nil {
		return ExitInternal
	}
	if len(args) == 0 {
		return unavailable(stderr, "validation", "no workflows were executed")
	}

	switch args[0] {
	case "--help", "-h":
		if len(args) == 1 {
			return printTo(stdout, usage, ExitSuccess)
		}
	case "help":
		if len(args) == 1 {
			return printTo(stdout, usage, ExitSuccess)
		}
		if len(args) == 2 && (args[1] == "init" || args[1] == "status") {
			return printTo(stdout, "Usage: ci-see "+args[1]+" [--help]\nNot implemented in S1.\n", ExitSuccess)
		}
	case "--version", "version":
		if len(args) == 1 {
			if build.Version == "" {
				build.Version = "dev"
			}
			if build.Commit == "" {
				build.Commit = "unknown"
			}
			if build.Go == "" {
				build.Go = "unknown"
			}
			return printTo(stdout, fmt.Sprintf("ci-see %s (commit %s; %s)\n", build.Version, build.Commit, build.Go), ExitSuccess)
		}
	case "init", "status":
		if len(args) == 2 && (args[1] == "--help" || args[1] == "-h") {
			return printTo(stdout, "Usage: ci-see "+args[0]+" [--help]\nNot implemented in S1.\n", ExitSuccess)
		}
		if len(args) == 1 {
			if args[0] == "init" {
				return unavailable(stderr, "initialization", "no hooks or configuration were written")
			}
			return unavailable(stderr, "result lookup", "no validation result is available")
		}
	}

	return printTo(stderr, "ci-see: invalid command or arguments; use 'ci-see --help'.\n"+usage, ExitUsage)
}

func unavailable(stderr io.Writer, operation, detail string) int {
	return printTo(stderr, fmt.Sprintf("ci-see: %s is not implemented in S1; %s.\n", operation, detail), ExitUnavailable)
}

func printTo(w io.Writer, message string, successCode int) int {
	if _, err := io.WriteString(w, message); err != nil {
		return ExitInternal
	}
	return successCode
}
