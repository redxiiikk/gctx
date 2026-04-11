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
	path := filepath.Join(wd, "mgit.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
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
