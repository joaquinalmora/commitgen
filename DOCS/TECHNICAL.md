# commitgen — Technical Reference

This document explains the code layout, the runtime behavior of the CLI and hook, and developer-facing details useful when extending or packaging the project.

## Repository layout (key files)

- `cmd/commitgen/main.go` — CLI entrypoint and command wiring.
- `internal/diff/diff.go` — Git interactions to read staged files and produce a trimmed patch.
- `internal/prompt/prompt.go` — Heuristics to produce a one-line suggested subject.
- `internal/hook/hook.go` — Hook installer writing `prepare-commit-msg` scripts.
- `internal/shell/shell.go` — zsh plugin-first snippet writer and guarded `.zshrc` block manager.
- `internal/doctor/doctor.go` — environment & install diagnostics (`commitgen doctor`).

## Commands (what they do)

- `commitgen suggest [--plain] [--verbose]` — Generate a suggested commit subject.
  - `--plain`: prints exactly one trimmed subject line to stdout (suitable for shells/hooks). If no staged changes, prints nothing and exits 0.
  - `--verbose`: prints diagnostic output (patch preview, sizes, and the message) to stderr.
  - Non-plain mode prints a human-facing message; non-plain runs now exit non-zero when there are no staged files (so CI and CLI callers can detect failure).

## Suggest flow (high-level)

1. `suggest` calls `internal/diff.StagedChanges(limit)` to get `files []string` and a trimmed `patch string`.
2. `prompt.MakePrompt(files, patch)` applies bucket heuristics:
   - `isTestsOnly` / `isDocsOnly` / `isConfigOnly` / `isRenameOnly`.
   - Default: `Update <file1>, <file2> (+N more)`.
3. Output is formatted according to flags (`--plain` vs normal; debug to stderr when `--verbose`).

## Plugin-first zsh design

This project uses a plugin-first approach for zsh inline suggestions.

### Goal

Surface a deterministic, single-line commit suggestion as inline ghost text in interactive zsh shells. Provide a safe, conservative fallback for users who don't have `zsh-autosuggestions` installed.

### User-facing contract

- Input: local repository and staged changes (no network calls).
- Output: a single-line suggestion string (suitable for `git commit -m "<suggestion>"`).
- Error modes: empty workspace or no staged files → no suggestion (graceful no-op).
- Privacy: all work happens locally; nothing is transmitted.

### Plugin-first behavior (zsh-autosuggestions)

- Detect presence of `zsh-autosuggestions` using common env vars and default plugin paths.
- Provide a minimal strategy function `_zsh_autosuggest_strategy_commitgen` that calls `commitgen suggest --plain` when the user is typing a `git commit -m "..."` (or the configured alias).
- Prepend the strategy so `commitgen` is tried before history suggestions.

Rationale: This delivers true inline ghost text with the plugin's native accept/decline behavior (right-arrow to accept). It's non-invasive: it only reads staged files and does not change commit files.

### Native fallback (POSTDISPLAY)

- If the plugin isn't present, show a dim preview using `zle -M` when the user types `git commit -m "`.
- Provide a widget `cg-accept-preview` bound to Ctrl-F (and optionally Right/End keys) that inserts the suggestion into the command buffer.
- Keep fallback conservative: no automatic edits to `COMMIT_EDITMSG`; users must accept the preview explicitly.

### Implementation notes

- The snippet lives under `internal/shell/commitgen.zsh` and the installer writes it to `~/.config/commitgen.zsh` and adds a guarded source block in `~/.zshrc`.
- Suggestion generation is synchronous but should be fast; prefer caching or asynchronous invocation if real-world delays appear.
- Store any short-lived cache in `XDG_CACHE_HOME` with user-only permissions.

### Edge cases and mitigations

- Slow suggestion generation: run asynchronously or add debounce to avoid prompt stalls.
- Plugin detection: support common plugin managers (Oh My Zsh, antigen, zplug) via environment checks and common paths.
- Shell customizations: the snippet must be idempotent and tolerant of re-sourcing.

### Tests and QA

- Unit tests: `internal/prompt` and `internal/shell` should have deterministic unit tests for suggestion output and installer idempotency.
- Manual QA: interactive tests in zsh with and without `zsh-autosuggestions` to verify ghost behavior and fallback accept widget.

### Next work items

- Harden the fallback (debounce + accurate inside-quote detection).
- Add optional `prepare-commit-msg` prefill mode (opt-in) as an alternative UX.

Note: `commitgen doctor` has been implemented to help validate local installs and environment (see `internal/doctor/doctor.go`).

Keep this document synchronized with code changes. Use it as the developer reference when working on the shell integration and related tests.

