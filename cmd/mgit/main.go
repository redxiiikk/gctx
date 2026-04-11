package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/redxiiikk/mgit/internal/config"
	"github.com/redxiiikk/mgit/internal/runner"
)

func main() {
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
