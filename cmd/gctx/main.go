package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/redxiiikk/gctx/internal/config"
	"github.com/redxiiikk/gctx/internal/runner"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "gctx" {
		os.Exit(runGctxCmd(os.Args[2:]))
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

func runGctxCmd(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gctx gctx <command>")
		fmt.Fprintln(os.Stderr, "commands: version, init, completion")
		return 1
	}
	switch args[0] {
	case "version":
		return cmdVersion()
	case "init":
		return cmdInit()
	case "completion":
		return cmdCompletion(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown gctx command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "commands: version, init, completion")
		return 1
	}
}
