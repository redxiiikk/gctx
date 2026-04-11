package main

import (
	"fmt"
	"os"
)

const bashCompletion = `# mgit bash completion
# Recommended — add to ~/.bashrc so completions stay up-to-date automatically:
#   eval "$(mgit mgit completion bash)"
#
# Alternative — write a static file:
#   mgit mgit completion bash > ~/.bash_completion.d/mgit

_mgit_complete() {
    local cur prev
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Handle 'mgit mgit <subcommand>' completion
    if [[ "${COMP_WORDS[1]}" == "mgit" && "$COMP_CWORD" -eq 2 ]]; then
        COMPREPLY=($(compgen -W "version init completion" -- "$cur"))
        return 0
    fi

    # Delegate to git completion for all other cases
    if declare -f __git_wrap__git_main > /dev/null 2>&1; then
        local old_words=("${COMP_WORDS[@]}")
        local old_cword=$COMP_CWORD
        COMP_WORDS=("git" "${COMP_WORDS[@]:1}")
        __git_wrap__git_main
        COMP_WORDS=("${old_words[@]}")
        COMP_CWORD=$old_cword
    elif declare -f _git > /dev/null 2>&1; then
        local old_words=("${COMP_WORDS[@]}")
        COMP_WORDS=("git" "${COMP_WORDS[@]:1}")
        _git
        COMP_WORDS=("${old_words[@]}")
    fi
}

complete -o bashdefault -o default -o nospace -F _mgit_complete mgit
`

const zshCompletion = `#compdef mgit
# mgit zsh completion
# Recommended — add to ~/.zshrc so completions stay up-to-date automatically:
#   eval "$(mgit mgit completion zsh)"
#
# Alternative — write a static file (requires fpath setup, see README):
#   mgit mgit completion zsh > ~/.zfunc/_mgit

_mgit() {
    # 'mgit mgit <subcommand>' — complete built-in subcommands
    if [[ "${words[2]}" == "mgit" ]]; then
        if (( CURRENT == 3 )); then
            local -a subcmds=(
                'version:Show mgit version information'
                'init:Initialize a new mgit.yaml configuration file'
                'completion:Generate shell completion scripts (bash/zsh/fish)'
            )
            _describe 'mgit subcommand' subcmds
        fi
        return
    fi

    # All other cases: delegate to git completion
    words[1]=(git)
    service=git
    _git
}

# #compdef is only processed when installed via fpath.
# compdef explicitly registers the function for both eval and fpath installs.
compdef _mgit mgit
`

const fishCompletion = `# mgit fish completion
# Recommended — add to ~/.config/fish/config.fish so completions stay up-to-date:
#   mgit mgit completion fish | source
#
# Alternative — write a static file:
#   mgit mgit completion fish > ~/.config/fish/completions/mgit.fish

# Disable default file completions
complete -c mgit -f

# 'mgit mgit <subcommand>' — complete built-in subcommands
complete -c mgit -n '__fish_seen_subcommand_from mgit' \
    -a 'version' -d 'Show mgit version information'
complete -c mgit -n '__fish_seen_subcommand_from mgit' \
    -a 'init' -d 'Initialize a new mgit.yaml configuration file'
complete -c mgit -n '__fish_seen_subcommand_from mgit' \
    -a 'completion' -d 'Generate shell completion scripts'

# For everything else, wrap git completions
complete -c mgit -n 'not __fish_seen_subcommand_from mgit' -w git
`

func cmdCompletion(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: mgit mgit completion <shell>")
		fmt.Fprintln(os.Stderr, "supported shells: bash, zsh, fish")
		return 1
	}

	switch args[0] {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "supported shells: bash, zsh, fish")
		return 1
	}

	return 0
}
