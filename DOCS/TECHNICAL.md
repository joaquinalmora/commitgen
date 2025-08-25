# commitgen — Technical Reference

This document explains the code layout, the runtime behavior of the CLI and hook, and developer-facing details useful when extending or packaging the project.

## Repository layout (key files)
- `cmd/commitgen/main.go` — CLI entrypoint. Registers commands and dispatches. Central flags: `--plain`, `--verbose` for the `suggest` command.
- `internal/diff/diff.go` — Git interactions: reads staged file list and staged patch. Uses `git diff --cached --name-only` and `git diff --cached --unified=0`. The staged patch is truncated to a file-size limit passed to `StagedChanges`.
- `internal/prompt/prompt.go` — Heuristics that build the suggested commit subject from the list of staged files and the trimmed patch. Buckets include: docs-only, tests-only, config-only, rename-only, and a default summarizer (first two files + `(+N more)`).
- `internal/hook/hook.go` — Hook installer. Writes a `prepare-commit-msg` script that invokes the stable `bin/commitgen` binary to populate the commit message file with the suggestion.
- `internal/shell/shell.go` — Shell installer/uninstaller. Writes `~/.config/commitgen.zsh` (plugin-first snippet) and manages a guarded block in `~/.zshrc` for idempotent installation.
- `internal/provider/provider.go` — Placeholder for future AI provider integrations.

## Commands (what they do)
- `commitgen suggest [--plain] [--verbose]` — Generate a suggested commit subject. Behavior:
  - `--plain`: prints exactly one trimmed subject line to stdout (suitable for shells/hooks). If no staged changes, prints nothing and exits 0.
  - `--verbose`: prints diagnostic output (patch preview, sizes, and the message) to stderr; stdout remains the message for non-plain mode.
  - Without flags, the tool prints a human-facing message to stdout. `--verbose` directs debug to stderr.
- `commitgen install-hook` — Writes a `prepare-commit-msg` hook into the repository's `.git/hooks/prepare-commit-msg` file (current implementation writes the file in the current working directory's `.git` folder).
- `commitgen init-shell` — Writes `~/.config/commitgen.zsh` and appends a guarded `source` block to `~/.zshrc` (idempotent).
- `commitgen uninstall-shell` — Removes the guarded block and deletes `~/.config/commitgen.zsh` (idempotent).

## Suggest flow (high-level)
1. `suggest` calls `internal/diff.StagedChanges(limit)` to get `files []string` and a trimmed `patch string`.
2. `prompt.MakePrompt(files, patch)` applies bucket heuristics:
   - `isTestsOnly` / `isDocsOnly` / `isConfigOnly` / `isRenameOnly`
   - Default: `Update <file1>, <file2> (+N more)`.
3. Output is formatted according to flags (`--plain` vs normal; debug to stderr when `--verbose`).

## Hook script details
- The hook written by `InstallHook()` is a POSIX `sh` script that:
  - Skips when the message file already has content, or when the hook is called from merge/squash/rebase flows.
  - Runs `bin/commitgen suggest` and writes the single-line suggestion into the commit message file (if non-empty).

## Shell integration (zsh)
- Installer writes a **plugin-first** snippet. Behavior:
  - If `zsh-autosuggestions` is available, a custom strategy `_zsh_autosuggest_strategy_commitgen` is defined to call `commitgen suggest --plain` and surface the result as ghost text (dimmed) that the plugin accepts with →.
  - Otherwise a native fallback provides a quick POSTDISPLAY-style preview and a `cg-accept-preview` widget; this fallback is intentionally conservative and is a prototype that requires refinement (debounce, accurate inside-quote detection, and robust accept behavior).

## Tests and verification
- Unit tests: `internal/prompt/prompt_test.go` checks that `MakePrompt` produces single-line, non-empty prompts for representative inputs.
- Run tests: `go test ./...` (or target `internal/prompt` for prompt tests).
- Build binary: `go build -o bin/commitgen ./cmd/commitgen`.

## Developer notes & best practices
- Keep `--plain` output free of non-message stdout; use stderr for diagnostics so `commitgen suggest --plain` is safe to embed in shell commands and hooks.
- Keep git interactions local to `diff` package; if running from subdirectories, prefer `git rev-parse --show-toplevel` to locate the repo root when installing hooks.
- Shell installer is conservative: it appends a guarded block and will not remove manual edits made outside the guarded block. `uninstall-shell` removes only the guarded block created by `init-shell`.

## Next technical tasks (short)
- Harden native zsh fallback (POSTDISPLAY-only, debounce, correct insertion on accept).
- Implement `uninstall-hook` with backup/restore semantics.
- Add `commitgen doctor` to validate environment, snippet sourcing, hook presence, and binary availability.
- Add CI (unit tests + build on Linux/macOS) and packaging steps (releases and Homebrew tap).

## File map quick reference
- `MakePrompt(files, patch)` — prompt heuristics: `internal/prompt/prompt.go`.
- `StagedChanges(limit)` — git patch retrieval and truncation: `internal/diff/diff.go`.
- Hook writer — `internal/hook/hook.go`.
- Shell installer — `internal/shell/shell.go`.

---
This document is intended as an internal developer reference — keep it up to date as the codebase changes. Use it as the source when composing a user-facing README or technical blog posts.
