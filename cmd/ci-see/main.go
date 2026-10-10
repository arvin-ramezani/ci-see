package main

import (
	"os"
	"runtime"

	"github.com/arvin-ramezani/ci-see/internal/cli"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Go:      runtime.Version(),
	}))
}
