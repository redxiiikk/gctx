package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/redxiiikk/gctx/internal/config"
	"gopkg.in/yaml.v3"
)

func cmdConfig(args []string) int {
	switch len(args) {
	case 0:
		return cmdConfigShow()
	case 1:
		return cmdConfigGet(args[0])
	case 2:
		return cmdConfigSet(args[0], args[1])
	default:
		fmt.Fprintln(os.Stderr, "usage: gctx gctx config [<key> [<value>]]")
		fmt.Fprintf(os.Stderr, "valid keys: %s\n", strings.Join(config.Keys(), ", "))
		return 1
	}
}

func cmdConfigShow() int {
	c, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrConfigNotFound) {
			fmt.Println("No configuration found.")
			return 0
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if c.IsEmpty() {
		fmt.Println("No configuration found.")
		return 0
	}

	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Print(string(yamlBytes))
	return 0
}

func cmdConfigGet(key string) int {
	c, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrConfigNotFound) {
			fmt.Fprintln(os.Stderr, "No configuration found. Run 'gctx gctx init' first.")
			return 1
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	value, err := c.Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println(value)
	return 0
}

func cmdConfigSet(key, value string) int {
	c, err := config.Load()
	if err != nil {
		if !errors.Is(err, config.ErrConfigNotFound) {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		wd, werr := os.Getwd()
		if werr != nil {
			fmt.Fprintf(os.Stderr, "getting working directory: %v\n", werr)
			return 1
		}
		c = &config.Config{Path: filepath.Join(wd, "gctx.yaml")}
	}

	if err := c.Set(key, value); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if err := c.Save(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Printf("Updated %s in %s\n", key, c.Path)
	return 0
}
