package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/redxiiikk/mgit/internal/config"
	"github.com/redxiiikk/mgit/internal/runner"
)

// Set via -ldflags at build time.
var (
	version   = "dev-build-version"
	buildDate = "N/A"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "mgit" {
		os.Exit(runMgitCmd(os.Args[2:]))
	}

	cfg, err := config.Load()
	if err != nil && !errors.Is(err, config.ErrConfigNotFound) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code, err := runner.Run(cfg, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(code)
}

func runMgitCmd(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: mgit mgit <command>")
		fmt.Fprintln(os.Stderr, "commands: version, init")
		return 1
	}
	switch args[0] {
	case "version":
		return cmdVersion()
	case "init":
		return cmdInit()
	default:
		fmt.Fprintf(os.Stderr, "unknown mgit command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "commands: version, init")
		return 1
	}
}
