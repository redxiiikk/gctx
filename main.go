package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Fprintln(os.Stdin, "hello world")

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gitArgs := os.Args[1:]
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	sshCmd := gitSSHCommand(cfg)
	if sshCmd != "" && needsSSHAuth(gitArgs) {
		env := environWithout(map[string]struct{}{
			"GIT_SSH_COMMAND": {},
			"GIT_SSH":         {},
		})
		cmd.Env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	}

	err = cmd.Run()
	if err == nil {
		return
	}
	if ee, ok := errors.AsType[*exec.ExitError](err); ok {
		os.Exit(ee.ExitCode())
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
