package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		printInjection(cfg, needSSH, needAuthor)
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

// printInjection prints the mgit configuration that will be injected so the
// user is aware of what mgit is doing before git runs.
func printInjection(cfg *config.Config, needSSH, needAuthor bool) {
	if cfg == nil {
		return
	}
	if needSSH {
		fmt.Fprintf(os.Stderr, "[mgit] using SSH key: %s\n", config.ExpandPath(cfg.SSHPrivateKey))
	}
	if needAuthor {
		name := strings.TrimSpace(cfg.GitUsername)
		email := strings.TrimSpace(cfg.GitEmail)
		if name != "" {
			fmt.Fprintf(os.Stderr, "[mgit] using author name: %s\n", name)
		}
		if email != "" {
			fmt.Fprintf(os.Stderr, "[mgit] using author email: %s\n", email)
		}
	}
}

// buildEnv constructs the subprocess environment, stripping only the variables
// that will be reinjected — so unconfigured fields are left untouched.
func buildEnv(sshCmd string, authorPairs []string, needSSH, needAuthor bool) []string {
	drops := map[string]struct{}{}
	if needSSH {
		drops["GIT_SSH_COMMAND"] = struct{}{}
		drops["GIT_SSH"] = struct{}{}
	}
	if needAuthor {
		// Only drop the vars that are actually present in authorPairs so that
		// unconfigured fields (e.g. no name, only email) keep their original
		// values from the environment.
		for _, pair := range authorPairs {
			key := strings.SplitN(pair, "=", 2)[0]
			drops[key] = struct{}{}
		}
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
