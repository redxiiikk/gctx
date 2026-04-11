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

Or build from source:

```bash
git clone https://github.com/redxiiikk/mgit.git
cd mgit
go build -o mgit ./cmd/mgit
```

## Usage

Use `mgit` everywhere you would use `git`:

```bash
mgit clone git@github.com:org/repo.git
mgit commit -m "fix: something"
mgit push origin main
```

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

## Project structure

```
mgit/
├── cmd/mgit/main.go          — entry point
├── internal/
│   ├── config/               — config loading and path expansion
│   ├── gitcmd/               — subcommand detection, SSH and author helpers
│   └── runner/               — subprocess execution
├── mgit.example.yaml
└── README.md
```

## License

MIT
