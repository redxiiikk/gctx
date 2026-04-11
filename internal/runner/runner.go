package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/redxiiikk/mgit/internal/config"
	"github.com/redxiiikk/mgit/internal/gitcmd"
)

// Run executes git with gitArgs, injecting SSH and author identity environment
// variables as dictated by cfg. It returns the git process exit code, or a
// non-nil error if the process could not be started or the SSH key is missing.
// A nil cfg is valid and disables all injection.
func Run(cfg *config.Config, gitArgs []string) (int, error) {
	needSSH := gitcmd.NeedsSSHAuth(gitArgs)
	var sshCmd string
	if needSSH {
		var err error
		sshCmd, err = gitcmd.SSHCommand(cfg)
		if err != nil {
			return 1, err
		}
		// SSHCommand returns "" when no key is configured; skip injection.
		needSSH = sshCmd != ""
	}

	pairs := gitcmd.AuthorEnvPairs(cfg)
	needAuthor := len(pairs) > 0 && gitcmd.NeedsAuthorIdentity(gitArgs)

	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if needSSH || needAuthor {
		cmd.Env = buildEnv(sshCmd, pairs, needSSH, needAuthor)
	}

	if err := cmd.Run(); err == nil {
		return 0, nil
	} else {
		if ee, ok := errors.AsType[*exec.ExitError](err); ok {
			return ee.ExitCode(), nil
		}
		return 1, fmt.Errorf("running git: %w", err)
	}
}

// buildEnv constructs the subprocess environment, stripping any pre-existing
// SSH or author identity variables before appending the new values.
func buildEnv(sshCmd string, authorPairs []string, needSSH, needAuthor bool) []string {
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
	env := gitcmd.EnvironWithout(drops)
	if needSSH {
		env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	}
	if needAuthor {
		env = append(env, authorPairs...)
	}
	return env
}
