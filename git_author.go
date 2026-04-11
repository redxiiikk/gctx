package main

import "strings"

// needsAuthorIdentity reports whether git may use user.name / user.email for new objects.
func needsAuthorIdentity(args []string) bool {
	sub, rest := firstGitSubcommand(args)
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

// authorEnvPairs returns GIT_AUTHOR_* / GIT_COMMITTER_* entries for exec.Command.Env.
func authorEnvPairs(cfg *Config) []string {
	if cfg == nil {
		return nil
	}
	name := strings.TrimSpace(cfg.GitUsername)
	email := strings.TrimSpace(cfg.GitEmail)
	if name == "" && email == "" {
		return nil
	}
	var p []string
	if name != "" {
		p = append(p, "GIT_AUTHOR_NAME="+name, "GIT_COMMITTER_NAME="+name)
	}
	if email != "" {
		p = append(p, "GIT_AUTHOR_EMAIL="+email, "GIT_COMMITTER_EMAIL="+email)
	}
	return p
}
