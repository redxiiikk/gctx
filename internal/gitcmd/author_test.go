package gitcmd

import (
	"reflect"
	"testing"

	"github.com/redxiiikk/mgit/internal/config"
)

func TestNeedsAuthorIdentity(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		// Commands that always need identity
		{[]string{"commit", "-m", "msg"}, true},
		{[]string{"merge", "feature"}, true},
		{[]string{"rebase", "main"}, true},
		{[]string{"cherry-pick", "abc123"}, true},
		{[]string{"revert", "HEAD"}, true},
		{[]string{"pull"}, true},
		{[]string{"am"}, true},
		{[]string{"tag", "v1.0"}, true},

		// stash sub-commands
		{[]string{"stash"}, true},
		{[]string{"stash", "push"}, true},
		{[]string{"stash", "save"}, true},
		{[]string{"stash", "branch", "new-branch"}, true},
		{[]string{"stash", "store"}, true},
		{[]string{"stash", "create"}, true},
		{[]string{"stash", "pop"}, false},
		{[]string{"stash", "drop"}, false},
		{[]string{"stash", "list"}, false},
		{[]string{"stash", "show"}, false},

		// notes sub-commands
		{[]string{"notes", "add"}, true},
		{[]string{"notes", "append"}, true},
		{[]string{"notes", "merge"}, true},
		{[]string{"notes", "edit"}, true},
		{[]string{"notes", "remove"}, false},
		{[]string{"notes", "list"}, false},
		{[]string{"notes"}, false},

		// Commands that never need identity
		{[]string{"push"}, false},
		{[]string{"fetch"}, false},
		{[]string{"status"}, false},
		{[]string{"log"}, false},
		{[]string{"diff"}, false},
		{[]string{}, false},

		// Global flags before subcommand
		{[]string{"-C", "/tmp", "commit", "-m", "x"}, true},
		{[]string{"--git-dir=/tmp/.git", "push"}, false},
	}

	for _, tt := range tests {
		got := NeedsAuthorIdentity(tt.args)
		if got != tt.want {
			t.Errorf("NeedsAuthorIdentity(%v) = %v, want %v", tt.args, got, tt.want)
		}
	}
}

func TestAuthorEnvPairs(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		if got := AuthorEnvPairs(nil); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("empty config", func(t *testing.T) {
		cfg := &config.Config{}
		if got := AuthorEnvPairs(cfg); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("name and email set", func(t *testing.T) {
		cfg := &config.Config{GitUsername: "Alice", GitEmail: "alice@example.com"}
		got := AuthorEnvPairs(cfg)
		want := []string{
			"GIT_AUTHOR_NAME=Alice",
			"GIT_COMMITTER_NAME=Alice",
			"GIT_AUTHOR_EMAIL=alice@example.com",
			"GIT_COMMITTER_EMAIL=alice@example.com",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("name only", func(t *testing.T) {
		cfg := &config.Config{GitUsername: "Bob"}
		got := AuthorEnvPairs(cfg)
		want := []string{"GIT_AUTHOR_NAME=Bob", "GIT_COMMITTER_NAME=Bob"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("email only", func(t *testing.T) {
		cfg := &config.Config{GitEmail: "bob@example.com"}
		got := AuthorEnvPairs(cfg)
		want := []string{"GIT_AUTHOR_EMAIL=bob@example.com", "GIT_COMMITTER_EMAIL=bob@example.com"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("whitespace-only fields treated as empty", func(t *testing.T) {
		cfg := &config.Config{GitUsername: "  ", GitEmail: "\t"}
		if got := AuthorEnvPairs(cfg); got != nil {
			t.Errorf("expected nil for whitespace-only fields, got %v", got)
		}
	})
}
