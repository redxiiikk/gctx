package gitcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/redxiiikk/mgit/internal/config"
)

// NeedsSSHAuth reports whether this git invocation may spawn ssh to a remote.
func NeedsSSHAuth(args []string) bool {
	sub, rest := FirstGitSubcommand(args)
	if sub == "" {
		return false
	}
	switch sub {
	case "clone", "fetch", "pull", "push", "ls-remote":
		return true
	case "archive":
		return archiveNeedsSSH(rest)
	case "remote":
		return remoteNeedsSSH(rest)
	case "submodule":
		return submoduleNeedsSSH(rest)
	case "send-pack", "receive-pack", "upload-pack", "upload-archive":
		return true
	default:
		return false
	}
}

func remoteNeedsSSH(rest []string) bool {
	if len(rest) == 0 {
		return false
	}
	switch rest[0] {
	case "show", "update", "prune":
		return true
	default:
		return false
	}
}

func archiveNeedsSSH(rest []string) bool {
	for _, a := range rest {
		if a == "--remote" || strings.HasPrefix(a, "--remote=") {
			return true
		}
	}
	return false
}

func submoduleNeedsSSH(rest []string) bool {
	if len(rest) == 0 {
		return false
	}
	switch rest[0] {
	case "add", "update", "sync", "init":
		return true
	default:
		return false
	}
}

// SSHEnvVars returns the GIT_SSH_COMMAND env var built from cfg, or nil if no
// SSH key is configured. It returns an error if the configured key file does
// not exist on disk.
func SSHEnvVars(cfg *config.Config) ([]EnvVar, error) {
	if cfg == nil {
		return nil, nil
	}
	key := config.ExpandPath(cfg.SSHPrivateKey)
	if key == "" {
		return nil, nil
	}
	if _, err := os.Stat(key); err != nil {
		return nil, fmt.Errorf("SSH private key not found: %s", key)
	}
	cmd := "ssh -i " + shellQuoteSingle(key) + " -o IdentitiesOnly=yes"
	return []EnvVar{{Key: "GIT_SSH_COMMAND", Value: cmd}}, nil
}

func shellQuoteSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}
