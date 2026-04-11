package gitcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/redxiiikk/mgit/internal/config"
)

func TestNeedsSSHAuth(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{[]string{"push"}, true},
		{[]string{"push", "origin", "main"}, true},
		{[]string{"fetch"}, true},
		{[]string{"pull"}, true},
		{[]string{"clone", "git@github.com:foo/bar"}, true},
		{[]string{"ls-remote"}, true},
		{[]string{"send-pack"}, true},
		{[]string{"receive-pack"}, true},
		{[]string{"upload-pack"}, true},
		{[]string{"upload-archive"}, true},

		// archive: needs --remote flag
		{[]string{"archive", "--remote=origin", "HEAD"}, true},
		{[]string{"archive", "--remote", "origin", "HEAD"}, true},
		{[]string{"archive", "HEAD"}, false},

		// remote sub-commands
		{[]string{"remote", "show", "origin"}, true},
		{[]string{"remote", "update"}, true},
		{[]string{"remote", "prune", "origin"}, true},
		{[]string{"remote", "add", "origin", "url"}, false},
		{[]string{"remote"}, false},

		// submodule sub-commands
		{[]string{"submodule", "add", "url"}, true},
		{[]string{"submodule", "update"}, true},
		{[]string{"submodule", "sync"}, true},
		{[]string{"submodule", "init"}, true},
		{[]string{"submodule", "status"}, false},
		{[]string{"submodule"}, false},

		// local-only commands
		{[]string{"status"}, false},
		{[]string{"log"}, false},
		{[]string{"commit", "-m", "msg"}, false},
		{[]string{}, false},

		// global flags before subcommand
		{[]string{"-C", "/tmp", "push"}, true},
		{[]string{"--git-dir=/tmp/.git", "fetch"}, true},
	}

	for _, tt := range tests {
		got := NeedsSSHAuth(tt.args)
		if got != tt.want {
			t.Errorf("NeedsSSHAuth(%v) = %v, want %v", tt.args, got, tt.want)
		}
	}
}

func TestSSHCommand_NoConfig(t *testing.T) {
	cmd, err := SSHCommand(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "" {
		t.Errorf("expected empty command, got %q", cmd)
	}
}

func TestSSHCommand_NoKey(t *testing.T) {
	cfg := &config.Config{}
	cmd, err := SSHCommand(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "" {
		t.Errorf("expected empty command, got %q", cmd)
	}
}

func TestSSHCommand_KeyExists(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg := &config.Config{SSHPrivateKey: f.Name()}
	cmd, err := SSHCommand(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "ssh -i") {
		t.Errorf("expected 'ssh -i ...' in command, got %q", cmd)
	}
	if !strings.Contains(cmd, "IdentitiesOnly=yes") {
		t.Errorf("expected 'IdentitiesOnly=yes' in command, got %q", cmd)
	}
}

func TestSSHCommand_KeyMissing(t *testing.T) {
	cfg := &config.Config{SSHPrivateKey: "/nonexistent/path/id_rsa"}
	_, err := SSHCommand(cfg)
	if err == nil {
		t.Error("expected error for missing key file, got nil")
	}
}

func TestSSHCommand_TildeExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	// Create a temp key file inside home so tilde expansion can be verified.
	f, err := os.CreateTemp(home, "mgit_test_key_*")
	if err != nil {
		t.Skip("cannot create temp file in home dir")
	}
	f.Close()
	defer os.Remove(f.Name())

	rel := "~/" + filepath.Base(f.Name())
	cfg := &config.Config{SSHPrivateKey: rel}
	cmd, err := SSHCommand(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmd, "ssh -i") {
		t.Errorf("expected 'ssh -i ...' in command, got %q", cmd)
	}
}

func TestShellQuoteSingle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/home/user/.ssh/id_rsa", `'/home/user/.ssh/id_rsa'`},
		{"path with spaces", `'path with spaces'`},
		{"it's a key", `'it'\''s a key'`},
		{"", `''`},
	}
	for _, tt := range tests {
		got := shellQuoteSingle(tt.input)
		if got != tt.want {
			t.Errorf("shellQuoteSingle(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
