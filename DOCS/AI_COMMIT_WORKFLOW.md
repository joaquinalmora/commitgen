# AI Commit Workflow (zsh + git)

This document shows how to integrate CommitGen (the `commitgen` CLI) into your shell without interfering with existing `git` aliases or completions.

## Requirements

- `commitgen` CLI on `PATH` (built or installed) capable of `commitgen suggest --ai --plain`.
- zsh shell (the snippets rely on zsh widgets).

## 1. Shell Function (`aicommit`)

Source `scripts/aicommit.zsh` from your `~/.zshrc`:

```zsh
source /path/to/commitgen/scripts/aicommit.zsh
```

This provides:

- `aicommit [git-commit-args]` – runs `commitgen suggest --ai --plain` once, validates the output, then executes `git commit -m "<message>" ...`.
- Git-style completion via `compdef _git aicommit=git-commit`, so tab-completion behaves exactly like `git commit`.
- Graceful failure: if `commitgen` errors or returns an empty string, `git commit` is not executed.

### Usage Examples

```bash
aicommit                   # generate message, run git commit
aicommit -a                # pass -a through to git
aicommit path/to/file      # commit staged subset
```

## 2. Optional `prepare-commit-msg` Hook

Copy `hooks/prepare-commit-msg.ai` to `.git/hooks/prepare-commit-msg` and make it executable:

```bash
cp hooks/prepare-commit-msg.ai .git/hooks/prepare-commit-msg
chmod +x .git/hooks/prepare-commit-msg
```

Behavior:

- Runs only when the commit message file is empty (no `-m`, no template, not a merge/squash).
- Calls `commitgen suggest --ai --plain` exactly once and writes the result to the message file.
- Skips generation when `aicommit` is used (because the message is already supplied) or when `commitgen` is unavailable.

This ensures there is never a double generation: either you invoke `aicommit` manually, or Git triggers the hook—but not both for the same commit.
