package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrConfigNotFound is returned by Load when no mgit.yaml is found between the
// working directory and the user's home directory.
var ErrConfigNotFound = errors.New("mgit.yaml not found")

// Config mirrors the fields in mgit.yaml.
type Config struct {
	SSHPrivateKey string `yaml:"ssh_private_key"`
	GitUsername   string `yaml:"git_username"`
	GitEmail      string `yaml:"git_email"`
}

// Load searches for mgit.yaml starting from the current working directory and
// walking upward until the user's home directory. It returns ErrConfigNotFound
// if no file is located.
func Load() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	dir, err := filepath.Abs(wd)
	if err != nil {
		return nil, fmt.Errorf("resolving working directory: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}
	home = filepath.Clean(home)

	for {
		cfgPath := filepath.Join(dir, "mgit.yaml")
		data, err := os.ReadFile(cfgPath)
		if err == nil {
			var c Config
			if err := yaml.Unmarshal(data, &c); err != nil {
				return nil, fmt.Errorf("parsing %s: %w", cfgPath, err)
			}
			return &c, nil
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading %s: %w", cfgPath, err)
		}

		if dir == home {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, ErrConfigNotFound
}

// ExpandPath expands a leading ~ or ~/ in p to the user's home directory.
func ExpandPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		if p == "~" {
			return home
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
