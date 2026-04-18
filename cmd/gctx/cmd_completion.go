package main

import (
	"fmt"
	"os"
)

const bashCompletion = `# gctx bash completion
# Recommended — add to ~/.bashrc so completions stay up-to-date automatically:
#   eval "$(gctx gctx completion bash)"
#
# Alternative — write a static file:
#   gctx gctx completion bash > ~/.bash_completion.d/gctx

_gctx_complete() {
    local cur prev
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Handle 'gctx gctx <subcommand>' completion
    if [[ "${COMP_WORDS[1]}" == "gctx" && "$COMP_CWORD" -eq 2 ]]; then
        COMPREPLY=($(compgen -W "version init completion config" -- "$cur"))
        return 0
    fi

    # 'gctx gctx config <key>' — valid config keys
    if [[ "${COMP_WORDS[1]}" == "gctx" && "${COMP_WORDS[2]}" == "config" && "$COMP_CWORD" -eq 3 ]]; then
        COMPREPLY=($(compgen -W "ssh_private_key git_username git_email" -- "$cur"))
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

complete -o bashdefault -o default -o nospace -F _gctx_complete gctx
`

const zshCompletion = `#compdef gctx
# gctx zsh completion
# Recommended — add to ~/.zshrc so completions stay up-to-date automatically:
#   eval "$(gctx gctx completion zsh)"
#
# Alternative — write a static file (requires fpath setup, see README):
#   gctx gctx completion zsh > ~/.zfunc/_gctx

_gctx() {
    # 'gctx gctx <subcommand>' — complete built-in subcommands
    if [[ "${words[2]}" == "gctx" ]]; then
        if (( CURRENT == 3 )); then
            local -a subcmds=(
                'version:Show gctx version information'
                'init:Initialize a new gctx.yaml configuration file'
                'completion:Generate shell completion scripts (bash/zsh/fish)'
                'config:Show, get, or set gctx.yaml options'
            )
            _describe 'gctx subcommand' subcmds
        elif (( CURRENT == 4 )) && [[ "${words[3]}" == "config" ]]; then
            _values 'config key' ssh_private_key git_username git_email
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
compdef _gctx gctx
`

const fishCompletion = `# gctx fish completion
# Recommended — add to ~/.config/fish/config.fish so completions stay up-to-date:
#   gctx gctx completion fish | source
#
# Alternative — write a static file:
#   gctx gctx completion fish > ~/.config/fish/completions/gctx.fish

# Disable default file completions
complete -c gctx -f

# 'gctx gctx <subcommand>' — complete built-in subcommands
complete -c gctx -n '__fish_seen_subcommand_from gctx' \
    -a 'version' -d 'Show gctx version information'
complete -c gctx -n '__fish_seen_subcommand_from gctx' \
    -a 'init' -d 'Initialize a new gctx.yaml configuration file'
complete -c gctx -n '__fish_seen_subcommand_from gctx' \
    -a 'completion' -d 'Generate shell completion scripts'
complete -c gctx -n '__fish_seen_subcommand_from gctx' \
    -a 'config' -d 'Show, get, or set gctx.yaml options'

# 'gctx gctx config <key>' — valid config keys
complete -c gctx -n '__fish_seen_subcommand_from gctx; and __fish_prev_arg_in config' \
    -a 'ssh_private_key' -d 'Path to SSH private key'
complete -c gctx -n '__fish_seen_subcommand_from gctx; and __fish_prev_arg_in config' \
    -a 'git_username' -d 'Git author name'
complete -c gctx -n '__fish_seen_subcommand_from gctx; and __fish_prev_arg_in config' \
    -a 'git_email' -d 'Git author email'

# For everything else, wrap git completions
complete -c gctx -n 'not __fish_seen_subcommand_from gctx' -w git
`

func cmdCompletion(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gctx gctx completion <shell>")
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
