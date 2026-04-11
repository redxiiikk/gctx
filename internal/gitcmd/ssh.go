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

// SSHCommand returns the value for GIT_SSH_COMMAND built from cfg, or an empty
// string if no SSH key is configured. It returns an error if the configured
// key file does not exist on disk.
func SSHCommand(cfg *config.Config) (string, error) {
	if cfg == nil {
		return "", nil
	}
	key := config.ExpandPath(cfg.SSHPrivateKey)
	if key == "" {
		return "", nil
	}
	if _, err := os.Stat(key); err != nil {
		return "", fmt.Errorf("SSH private key not found: %s", key)
	}
	return "ssh -i " + shellQuoteSingle(key) + " -o IdentitiesOnly=yes", nil
}

// EnvironWithout returns os.Environ() with the named keys removed.
func EnvironWithout(keys map[string]struct{}) []string {
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

func shellQuoteSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}
