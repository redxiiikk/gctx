package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config mirrors mgit.yaml.
type Config struct {
	SSHPrivateKey string `yaml:"ssh_private_key"`
	GitUsername   string `yaml:"git_username"`
	GitEmail      string `yaml:"git_email"`
}

func loadConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir, err := filepath.Abs(filepath.Clean(wd))
	if err != nil {
		return nil, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	home = filepath.Clean(home)

	for {
		cfgPath := filepath.Join(dir, "mgit.yaml")
		data, err := os.ReadFile(cfgPath)
		if err == nil {
			var c Config
			if err := yaml.Unmarshal(data, &c); err != nil {
				return nil, err
			}
			return &c, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
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
	return nil, nil
}

func expandPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return home
	}
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
