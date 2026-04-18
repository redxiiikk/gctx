package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrConfigNotFound is returned by Load when no gctx.yaml is found between the
// working directory and the user's home directory.
var ErrConfigNotFound = errors.New("gctx.yaml not found")

// configFileName is the on-disk file name used for gctx configuration.
const configFileName = "gctx.yaml"

// Config mirrors the fields in gctx.yaml.
type Config struct {
	SSHPrivateKey string `yaml:"ssh_private_key,omitempty"`
	GitUsername   string `yaml:"git_username,omitempty"`
	GitEmail      string `yaml:"git_email,omitempty"`

	// Path is the absolute path the config was loaded from (or should be
	// written to). It is not persisted in gctx.yaml (yaml:"-").
	Path string `yaml:"-"`
}

// Load searches for gctx.yaml starting from the current working directory and
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
		cfgPath := filepath.Join(dir, configFileName)
		data, err := os.ReadFile(cfgPath)
		if err == nil {
			var c Config
			if err := yaml.Unmarshal(data, &c); err != nil {
				return nil, fmt.Errorf("parsing %s: %w", cfgPath, err)
			}
			c.Path = cfgPath
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

// Save serializes the config to its associated path, omitting fields that are
// empty. It creates the file (or truncates an existing one).
func (c *Config) Save() error {
	if c == nil {
		return errors.New("config: Save called on nil Config")
	}
	if c.Path == "" {
		return errors.New("config: Path is not set; use Load or assign Path")
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	if err := os.WriteFile(c.Path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", c.Path, err)
	}
	return nil
}

// Keys returns the list of valid configuration keys recognized by Set/Get.
func Keys() []string {
	return []string{"ssh_private_key", "git_username", "git_email"}
}

// Set updates the field identified by key. It returns an error when the key is
// not recognized.
func (c *Config) Set(key, value string) error {
	switch key {
	case "ssh_private_key":
		c.SSHPrivateKey = value
	case "git_username":
		c.GitUsername = value
	case "git_email":
		c.GitEmail = value
	default:
		return fmt.Errorf("unknown config key %q (valid keys: %s)", key, strings.Join(Keys(), ", "))
	}
	return nil
}

// Get returns the current value for the given key, or an error when the key is
// not recognized.
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "ssh_private_key":
		return c.SSHPrivateKey, nil
	case "git_username":
		return c.GitUsername, nil
	case "git_email":
		return c.GitEmail, nil
	default:
		return "", fmt.Errorf("unknown config key %q (valid keys: %s)", key, strings.Join(Keys(), ", "))
	}
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

func (c *Config) IsEmpty() bool {
	return c.SSHPrivateKey == "" && c.GitUsername == "" && c.GitEmail == ""
}
