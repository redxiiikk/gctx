package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// writeConfig creates a gctx.yaml in dir with the given content.
func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "gctx.yaml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLoad_NotFound(t *testing.T) {
	// Use a temp dir that has no gctx.yaml anywhere up to its root.
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	_, err := Load()
	if !errors.Is(err, ErrConfigNotFound) {
		t.Errorf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestLoad_FoundInCwd(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, "git_username: Alice\ngit_email: alice@example.com\n")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GitUsername != "Alice" {
		t.Errorf("GitUsername = %q, want %q", cfg.GitUsername, "Alice")
	}
	if cfg.GitEmail != "alice@example.com" {
		t.Errorf("GitEmail = %q, want %q", cfg.GitEmail, "alice@example.com")
	}
}

func TestLoad_FoundInParent(t *testing.T) {
	parent := t.TempDir()
	writeConfig(t, parent, "git_username: Parent\n")

	child := filepath.Join(parent, "subdir", "project")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GitUsername != "Parent" {
		t.Errorf("GitUsername = %q, want %q", cfg.GitUsername, "Parent")
	}
}

func TestLoad_CwdTakesPrecedenceOverParent(t *testing.T) {
	parent := t.TempDir()
	writeConfig(t, parent, "git_username: Parent\n")

	child := filepath.Join(parent, "child")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, child, "git_username: Child\n")
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GitUsername != "Child" {
		t.Errorf("GitUsername = %q, want %q", cfg.GitUsername, "Child")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir, ":\tinvalid: yaml: [\n")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}

	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"   ", ""},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
		{"~", home},
		{"~/foo/bar", filepath.Join(home, "foo/bar")},
		{"~/", home},
	}

	for _, tt := range tests {
		got := ExpandPath(tt.input)
		if got != tt.want {
			t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
