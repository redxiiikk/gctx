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
	var inject []gitcmd.EnvVar

	if gitcmd.NeedsSSHAuth(gitArgs) {
		sshEnv, err := gitcmd.SSHEnvVars(cfg)
		if err != nil {
			return 1, err
		}
		inject = append(inject, sshEnv...)
	}

	if gitcmd.NeedsAuthorIdentity(gitArgs) {
		inject = append(inject, gitcmd.AuthorEnvVars(cfg)...)
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(inject) > 0 {
		printInjection(cfg, inject)
		cmd.Env = buildEnv(inject)
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

// printInjection prints what mgit will inject so the user is aware before git
// runs. SSH key path is read from cfg for human-readable display; name and
// email are read from inject directly since they are already trimmed.
func printInjection(cfg *config.Config, inject []gitcmd.EnvVar) {
	if cfg == nil {
		return
	}
	if findEnvVar(inject, "GIT_SSH_COMMAND") != "" {
		fmt.Fprintf(os.Stderr, "[mgit] using SSH key: %s\n", config.ExpandPath(cfg.SSHPrivateKey))
	}
	if name := findEnvVar(inject, "GIT_AUTHOR_NAME"); name != "" {
		fmt.Fprintf(os.Stderr, "[mgit] using author name: %s\n", name)
	}
	if email := findEnvVar(inject, "GIT_AUTHOR_EMAIL"); email != "" {
		fmt.Fprintf(os.Stderr, "[mgit] using author email: %s\n", email)
	}
}

// buildEnv constructs the subprocess environment, dropping only the variables
// that will be reinjected — so unconfigured fields are left untouched.
// GIT_SSH is also dropped when GIT_SSH_COMMAND is present to prevent the
// legacy override from interfering.
func buildEnv(inject []gitcmd.EnvVar) []string {
	drops := map[string]struct{}{}
	for _, ev := range inject {
		drops[ev.Key] = struct{}{}
		if ev.Key == "GIT_SSH_COMMAND" {
			drops["GIT_SSH"] = struct{}{}
		}
	}
	env := osEnvironExcluding(drops)
	for _, ev := range inject {
		env = append(env, ev.Key+"="+ev.Value)
	}
	return env
}

// findEnvVar returns the Value of the first EnvVar whose Key matches key,
// or an empty string if not found.
func findEnvVar(vars []gitcmd.EnvVar, key string) string {
	for _, ev := range vars {
		if ev.Key == key {
			return ev.Value
		}
	}
	return ""
}

// osEnvironExcluding returns os.Environ() with the specified keys removed.
func osEnvironExcluding(keys map[string]struct{}) []string {
	var out []string
	for _, e := range os.Environ() {
		name := strings.SplitN(e, "=", 2)[0]
		if _, drop := keys[name]; drop {
			continue
		}
		out = append(out, e)
	}
	return out
}
