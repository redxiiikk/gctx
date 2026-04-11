package gitcmd

import (
	"strings"

	"github.com/redxiiikk/mgit/internal/config"
)

// NeedsAuthorIdentity reports whether the git invocation may create new objects
// that reference user.name / user.email.
func NeedsAuthorIdentity(args []string) bool {
	sub, rest := FirstGitSubcommand(args)
	switch sub {
	case "commit", "merge", "rebase", "cherry-pick", "revert", "pull", "am":
		return true
	case "tag":
		return true
	case "stash":
		return stashNeedsAuthor(rest)
	case "notes":
		return notesNeedsAuthor(rest)
	default:
		return false
	}
}

func stashNeedsAuthor(rest []string) bool {
	if len(rest) == 0 {
		return true
	}
	switch rest[0] {
	case "push", "save", "branch", "store", "create":
		return true
	default:
		return false
	}
}

func notesNeedsAuthor(rest []string) bool {
	if len(rest) == 0 {
		return false
	}
	switch rest[0] {
	case "add", "append", "merge", "edit":
		return true
	default:
		return false
	}
}

// AuthorEnvVars returns the GIT_AUTHOR_* / GIT_COMMITTER_* env vars derived
// from cfg. Only fields that are non-empty in the config produce entries.
func AuthorEnvVars(cfg *config.Config) []EnvVar {
	if cfg == nil {
		return nil
	}
	name := strings.TrimSpace(cfg.GitUsername)
	email := strings.TrimSpace(cfg.GitEmail)
	if name == "" && email == "" {
		return nil
	}
	var vars []EnvVar
	if name != "" {
		vars = append(vars,
			EnvVar{"GIT_AUTHOR_NAME", name},
			EnvVar{"GIT_COMMITTER_NAME", name},
		)
	}
	if email != "" {
		vars = append(vars,
			EnvVar{"GIT_AUTHOR_EMAIL", email},
			EnvVar{"GIT_COMMITTER_EMAIL", email},
		)
	}
	return vars
}
