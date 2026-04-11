package gitcmd

import (
	"reflect"
	"testing"
)

func TestFirstGitSubcommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantSub  string
		wantRest []string
	}{
		{
			name:     "empty args",
			args:     nil,
			wantSub:  "",
			wantRest: nil,
		},
		{
			name:     "subcommand only",
			args:     []string{"push"},
			wantSub:  "push",
			wantRest: []string{},
		},
		{
			name:     "subcommand with args",
			args:     []string{"push", "origin", "main"},
			wantSub:  "push",
			wantRest: []string{"origin", "main"},
		},
		{
			name:     "-c flag before subcommand",
			args:     []string{"-c", "user.name=test", "commit", "-m", "msg"},
			wantSub:  "commit",
			wantRest: []string{"-m", "msg"},
		},
		{
			name:     "-C flag before subcommand",
			args:     []string{"-C", "/tmp/repo", "push"},
			wantSub:  "push",
			wantRest: []string{},
		},
		{
			name:     "--git-dir space-separated before subcommand",
			args:     []string{"--git-dir", "/tmp/.git", "fetch"},
			wantSub:  "fetch",
			wantRest: []string{},
		},
		{
			name:     "--git-dir=value before subcommand",
			args:     []string{"--git-dir=/tmp/.git", "fetch"},
			wantSub:  "fetch",
			wantRest: []string{},
		},
		{
			name:     "--work-tree=value before subcommand",
			args:     []string{"--work-tree=/tmp", "status"},
			wantSub:  "status",
			wantRest: []string{},
		},
		{
			name:     "--namespace=value before subcommand",
			args:     []string{"--namespace=ns1", "push"},
			wantSub:  "push",
			wantRest: []string{},
		},
		{
			name:     "bare single-dash flags before subcommand",
			args:     []string{"-p", "log"},
			wantSub:  "log",
			wantRest: []string{},
		},
		{
			name:     "multiple global flags before subcommand",
			args:     []string{"-c", "core.pager=cat", "-C", "/tmp", "log", "--oneline"},
			wantSub:  "log",
			wantRest: []string{"--oneline"},
		},
		{
			name:     "bare dash is not a flag",
			args:     []string{"-"},
			wantSub:  "-",
			wantRest: []string{},
		},
		{
			name:     "only global flags, no subcommand",
			args:     []string{"-c", "core.editor=vim"},
			wantSub:  "",
			wantRest: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSub, gotRest := FirstGitSubcommand(tt.args)
			if gotSub != tt.wantSub {
				t.Errorf("sub = %q, want %q", gotSub, tt.wantSub)
			}
			if !reflect.DeepEqual(gotRest, tt.wantRest) {
				t.Errorf("rest = %v, want %v", gotRest, tt.wantRest)
			}
		})
	}
}
