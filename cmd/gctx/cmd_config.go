package main

import (
	"fmt"

	"github.com/redxiiikk/gctx/internal/config"
	"gopkg.in/yaml.v3"
)

func cmdConfig() int {
	c, err := config.Load()
	if err != nil {
		return 1
	}

	if c.IsEmpty() {
		fmt.Println("No configuration found.")
		return 0
	}

	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return 1
	}

	fmt.Println(string(yamlBytes))

	return 0
}
