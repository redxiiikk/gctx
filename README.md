# mgit

`mgit` is a thin Git wrapper that automatically injects per-project SSH keys
and author identity from a local configuration file. It is a transparent
pass-through: every argument you pass to `mgit` is forwarded to `git` unchanged.

## Why

When you work across multiple organisations (each with its own SSH key and
author identity), `mgit` removes the need to juggle `GIT_SSH_COMMAND`,
`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, or `includeIf` stanzas in
`~/.gitconfig`. Drop a `mgit.yaml` in a project (or workspace root) and every
`mgit` invocation in that tree picks it up automatically.

## Installation

```bash
go install github.com/redxiiikk/mgit/cmd/mgit@latest
```

Or build from source using [Task](https://taskfile.dev):

```bash
git clone https://github.com/redxiiikk/mgit.git
cd mgit
task build    # builds to dist/mgit
task install  # installs to $GOPATH/bin and sets up shell completion
```

`task install` auto-detects your current shell from `$SHELL` and installs the
appropriate completion script. You can also run it standalone at any time:

```bash
task completion
```

## Usage

Use `mgit` everywhere you would use `git`:

```bash
mgit clone git@github.com:org/repo.git
mgit commit -m "fix: something"
mgit push origin main
```

### Built-in commands

`mgit` also exposes its own subcommands under the `mgit mgit` namespace:

| Command | Description |
|---------|-------------|
| `mgit mgit version` | Show version and build date |
| `mgit mgit init` | Interactively create a `mgit.yaml` in the current directory |
| `mgit mgit completion <shell>` | Print a shell completion script (`bash`, `zsh`, or `fish`) |

### Shell tab completion

#### Option 1 — eval (recommended)

Add a single line to your shell's startup file. The completion script is
generated fresh on every shell start, so it stays in sync automatically
whenever `mgit` is updated.

**Zsh** — add to `~/.zshrc`:

```zsh
eval "$(mgit mgit completion zsh)"
```

**Bash** — add to `~/.bashrc`:

```bash
eval "$(mgit mgit completion bash)"
```

**Fish** — add to `~/.config/fish/config.fish`:

```fish
mgit mgit completion fish | source
```

#### Option 2 — static file

Write the script to disk once and source it from your shell config.
You will need to re-run the command after updating `mgit`.

**Zsh**

```zsh
mkdir -p ~/.zfunc
mgit mgit completion zsh > ~/.zfunc/_mgit
# Ensure ~/.zshrc contains (before compinit):
#   fpath=(~/.zfunc $fpath)
#   autoload -Uz compinit && compinit
```

**Bash**

```bash
mkdir -p ~/.bash_completion.d
mgit mgit completion bash > ~/.bash_completion.d/mgit
# Ensure ~/.bashrc contains:
#   source ~/.bash_completion.d/mgit
```

**Fish**

```fish
mgit mgit completion fish > ~/.config/fish/completions/mgit.fish
```

---

Once active, `mgit mgit <tab>` completes built-in subcommands, and all other
`mgit <tab>` invocations delegate to `git`'s native completion.

## Configuration

Create a file named `mgit.yaml` in your project directory (or any ancestor
directory up to your home directory). `mgit` searches upward from the current
working directory and uses the first file it finds.

```yaml
# ~/.ssh/id_ed25519_work will be passed to ssh via GIT_SSH_COMMAND.
# Supports ~ expansion.
ssh_private_key: ~/.ssh/id_ed25519_work

# Injected as GIT_AUTHOR_NAME / GIT_COMMITTER_NAME for object-creating commands.
git_username: Your Name

# Injected as GIT_AUTHOR_EMAIL / GIT_COMMITTER_EMAIL.
git_email: you@example.com
```

All fields are optional. See `mgit.example.yaml` for an annotated template.

### Config file search order

```
/home/user/work/acme/project/   ← starts here (cwd)
/home/user/work/acme/           ← walks upward
/home/user/work/
/home/user/                     ← stops at home directory
```

The first `mgit.yaml` found wins. If no file is found, `mgit` runs `git`
with no modifications to the environment.

### SSH key injection

When a git subcommand that contacts a remote is detected (`clone`, `fetch`,
`pull`, `push`, `ls-remote`, `remote show/update/prune`, `submodule add/update`,
etc.), `GIT_SSH_COMMAND` is set to:

```
ssh -i '<key_path>' -o IdentitiesOnly=yes
```

`IdentitiesOnly=yes` ensures the ssh-agent does not offer unrelated keys.

### Author identity injection

For commands that create Git objects (`commit`, `merge`, `rebase`,
`cherry-pick`, `revert`, `tag`, `pull`, `am`, `stash push/save`, `notes add`,
etc.), the following variables are set:

- `GIT_AUTHOR_NAME` / `GIT_COMMITTER_NAME`
- `GIT_AUTHOR_EMAIL` / `GIT_COMMITTER_EMAIL`

Any pre-existing values for these variables in the environment are replaced.

## IDE Integration

> **Note:** `mgit` works as a drop-in Git replacement for IDE operations that
> invoke an external `git` process. IDEs that use embedded Git libraries
> (e.g. JGit, libgit2) internally may bypass the configured executable, so
> full compatibility cannot be guaranteed.

### VSCode

Add to `.vscode/settings.json`:

```json
{
  "git.path": "/usr/local/bin/mgit"
}
```

### JetBrains IDEs

Open **Settings / Preferences** → **Version Control** → **Git**, set
**Path to Git executable** to the absolute path of `mgit`, then click **Test**.

## Project structure

```
mgit/
├── cmd/mgit/
│   ├── main.go               — entry point and mgit subcommand dispatcher
│   ├── cmd_version.go        — `mgit mgit version`
│   ├── cmd_init.go           — `mgit mgit init`
│   └── cmd_completion.go     — `mgit mgit completion`
├── internal/
│   ├── config/               — config loading and path expansion
│   ├── gitcmd/               — subcommand detection, SSH and author helpers
│   └── runner/               — subprocess execution
├── mgit.example.yaml
└── README.md
```

## License

MIT
