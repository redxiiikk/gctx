package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Fprintln(os.Stderr, "hello world")

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
	needSSH := sshCmd != "" && needsSSHAuth(gitArgs)
	pairs := authorEnvPairs(cfg)
	needAuthor := len(pairs) > 0 && needsAuthorIdentity(gitArgs)

	if needSSH || needAuthor {
		drops := map[string]struct{}{
			"GIT_SSH_COMMAND": {},
			"GIT_SSH":         {},
		}
		if needAuthor {
			drops["GIT_AUTHOR_NAME"] = struct{}{}
			drops["GIT_AUTHOR_EMAIL"] = struct{}{}
			drops["GIT_COMMITTER_NAME"] = struct{}{}
			drops["GIT_COMMITTER_EMAIL"] = struct{}{}
		}
		env := environWithout(drops)
		if needSSH {
			env = append(env, "GIT_SSH_COMMAND="+sshCmd)
		}
		if needAuthor {
			env = append(env, pairs...)
		}
		cmd.Env = env
	}

	err = cmd.Run()
	if err == nil {
		return
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
