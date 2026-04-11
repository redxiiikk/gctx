package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Fprintln(os.Stderr, "hello world")

	cmd := exec.Command("git", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return
	}
	if ee, ok := err.(*exec.ExitError); ok {
		os.Exit(ee.ExitCode())
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
