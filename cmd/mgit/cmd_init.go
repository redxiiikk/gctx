package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/redxiiikk/mgit/internal/config"
)

func cmdInit() int {
	reader := bufio.NewReader(os.Stdin)

	sshKey := prompt(reader, "SSH private key (leave empty to skip): ")
	username := prompt(reader, "Git username (leave empty to skip): ")
	email := prompt(reader, "Git email (leave empty to skip): ")

	cfg := &config.Config{
		SSHPrivateKey: sshKey,
		GitUsername:   username,
		GitEmail:      email,
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting working directory: %v\n", err)
		return 1
	}

	if err := config.Write(cfg, wd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Println("Created mgit.yaml")
	return 0
}

func prompt(r *bufio.Reader, label string) string {
	fmt.Print(label)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}
