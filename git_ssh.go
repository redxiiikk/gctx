package main

import (
	"os"
	"strings"
)

// needsSSHAuth reports whether this git invocation may spawn ssh to a remote.
// It parses past common global options (-c, -C, --git-dir, …) to find the subcommand.
func needsSSHAuth(args []string) bool {
	sub, rest := firstGitSubcommand(args)
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

func firstGitSubcommand(args []string) (sub string, restAfterSub []string) {
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-c" && i+1 < len(args):
			i += 2
		case a == "-C" || a == "--git-dir" || a == "--work-tree" || a == "--namespace":
			if i+1 < len(args) {
				i += 2
			} else {
				i++
			}
		case strings.HasPrefix(a, "--git-dir=") || strings.HasPrefix(a, "--work-tree=") ||
			strings.HasPrefix(a, "--namespace="):
			i++
		case strings.HasPrefix(a, "-") && a != "-":
			i++
		default:
			return a, args[i+1:]
		}
	}
	return "", nil
}

func shellQuoteSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}

// gitSSHCommand returns the value for GIT_SSH_COMMAND, or empty if nothing to set.
func gitSSHCommand(c *Config) string {
	if c == nil {
		return ""
	}
	key := expandPath(c.SSHPrivateKey)
	if key == "" {
		return ""
	}
	return "ssh -i " + shellQuoteSingle(key) + " -o IdentitiesOnly=yes"
}

func environWithout(keys map[string]struct{}) []string {
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
