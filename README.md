# gctx

`gctx` is a thin Git wrapper that automatically injects per-project SSH keys
and author identity from a local configuration file. It is a transparent
pass-through: every argument you pass to `gctx` is forwarded to `git` unchanged.

## Why

When you work across multiple organisations (each with its own SSH key and
author identity), `gctx` removes the need to juggle `GIT_SSH_COMMAND`,
`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, or `includeIf` stanzas in
`~/.gitconfig`. Drop a `gctx.yaml` in a project (or workspace root) and every
`gctx` invocation in that tree picks it up automatically.

## Installation

```bash
go install github.com/redxiiikk/gctx/cmd/gctx@latest
```

Or build from source with [mise](https://mise.jdx.dev/) (tasks are defined in `mise.toml`):

> Before building, install [goreleaser](https://goreleaser.com/) using
> their [installation guide](https://goreleaser.com/getting-started/install).

```bash
git clone https://github.com/redxiiikk/gctx.git
cd gctx
mise install    # optional: install the Go version from mise.toml
mise build      # builds to dist/gctx_darwin_{.Arch}/gctx
```

### Development

From the repository root:

| Command                     | Description                                                                                                  |
|-----------------------------|--------------------------------------------------------------------------------------------------------------|
| `mise test`                 | Run `go test ./...`                                                                                          |
| `mise lint`                 | Run `gofmt -w -s .` on the tree                                                                              |
| `mise run run <subcommand>` | Run the CLI from source (`go run ./cmd/gctx gctx …`). Subcommands: `init`, `version`, `config`, `completion` |

Examples:

```bash
mise run run version
mise run run config
```

The `completion` subcommand also needs a shell (`bash`, `zsh`, or `fish`). Because of that, call it directly:

```bash
go run ./cmd/gctx gctx completion zsh
```

## Usage

Use `gctx` everywhere you would use `git`:

```bash
gctx clone git@github.com:org/repo.git
gctx commit -m "fix: something"
gctx push origin main
```

### Built-in commands

`gctx` also exposes its own subcommands under the `gctx gctx` namespace:

| Command                        | Description                                                    |
|--------------------------------|----------------------------------------------------------------|
| `gctx gctx version`            | Show version and build metadata                                |
| `gctx gctx init`               | Interactively create a `gctx.yaml` in the current directory    |
| `gctx gctx config`             | Print the resolved effective config as YAML (from search path) |
| `gctx gctx completion <shell>` | Print a shell completion script (`bash`, `zsh`, or `fish`)     |

### Shell tab completion

#### Option 1 — eval (recommended)

Add a single line to your shell's startup file. The completion script is
generated fresh on every shell start, so it stays in sync automatically
whenever `gctx` is updated.

**Zsh** — add to `~/.zshrc`:

```zsh
eval "$(gctx gctx completion zsh)"
```

**Bash** — add to `~/.bashrc`:

```bash
eval "$(gctx gctx completion bash)"
```

**Fish** — add to `~/.config/fish/config.fish`:

```fish
gctx gctx completion fish | source
```

#### Option 2 — static file

Write the script to disk once and source it from your shell config.
You will need to re-run the command after updating `gctx`.

**Zsh**

```zsh
mkdir -p ~/.zfunc
gctx gctx completion zsh > ~/.zfunc/_gctx
# Ensure ~/.zshrc contains (before compinit):
#   fpath=(~/.zfunc $fpath)
#   autoload -Uz compinit && compinit
```

**Bash**

```bash
mkdir -p ~/.bash_completion.d
gctx gctx completion bash > ~/.bash_completion.d/gctx
# Ensure ~/.bashrc contains:
#   source ~/.bash_completion.d/gctx
```

**Fish**

```fish
gctx gctx completion fish > ~/.config/fish/completions/gctx.fish
```

---

Once active, `gctx gctx <tab>` completes built-in subcommands, and all other
`gctx <tab>` invocations delegate to `git`'s native completion.

## Configuration

Create a file named `gctx.yaml` in your project directory (or any ancestor
directory up to your home directory). `gctx` searches upward from the current
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

All fields are optional. See `gctx.example.yaml` for an annotated template.

### Config file search order

```
/home/user/work/acme/project/   ← starts here (cwd)
/home/user/work/acme/           ← walks upward
/home/user/work/
/home/user/                     ← stops at home directory
```

The first `gctx.yaml` found wins. If no file is found, `gctx` runs `git`
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

> **Note:** `gctx` works as a drop-in Git replacement for IDE operations that
> invoke an external `git` process. IDEs that use embedded Git libraries
> (e.g. JGit, libgit2) internally may bypass the configured executable, so
> full compatibility cannot be guaranteed.

### VSCode

Add to `.vscode/settings.json`:

```json
{
  "git.path": "/usr/local/bin/gctx"
}
```

### JetBrains IDEs

Open **Settings / Preferences** → **Version Control** → **Git**, set
**Path to Git executable** to the absolute path of `gctx`, then click **Test**.

## License

MIT
