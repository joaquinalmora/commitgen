# commitgen — Technical Reference

This document explains the code layout, the runtime behavior of the CLI and hook, and developer-facing details useful when extending or packaging the project.

## Repository layout (key files)

- `cmd/commitgen/main.go` — CLI entrypoint and command wiring.
- `internal/diff/diff.go` — Git interactions to read staged files and produce a trimmed patch.
- `internal/prompt/prompt.go` — Heuristics to produce a one-line suggested subject.
- `internal/hook/hook.go` — Hook installer writing `prepare-commit-msg` scripts.
- `internal/shell/shell.go` — zsh plugin-first snippet writer and guarded `.zshrc` block manager.

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

## Hook script details

- The hook (`prepare-commit-msg`) is a POSIX `sh` script written by `internal/hook.InstallHook()`.
- The script skips when the message file already has content or when the hook is triggered by merge/squash/rebase flows. It invokes `commitgen suggest --plain` and writes the suggestion to the commit message file if non-empty.

## Shell integration (zsh)

- Installer writes a plugin-first snippet at `~/.config/commitgen.zsh` and adds a guarded block to `~/.zshrc`.
- If `zsh-autosuggestions` is available, a custom strategy surfaces the suggestion as ghost text. Otherwise a conservative POSTDISPLAY fallback provides a preview and an accept widget. The native fallback is a prototype and needs debounce/robustness work.

## Tests and verification

- Unit tests: package-level `_test.go` files (e.g. `internal/prompt/prompt_test.go`, `internal/shell/shell_test.go`). Run `go test ./...`.
- Integration/e2e: `e2e/` contains end-to-end tests that build the binary and exercise the CLI in a temp git repo. `e2e/` is intentionally ignored by default to avoid shipping test artifacts.
- CI: a minimal workflow at `.github/workflows/ci.yml` runs a gofmt check and `go test ./...` on Ubuntu for push/PR. macOS and e2e jobs are exposed as manual workflows (`workflow_dispatch`).

## Developer notes & best practices

- Keep `--plain` output free of non-message stdout; diagnostics should use stderr so hooks can safely capture stdout.
- Target Unix-like shells (macOS + Linux) for shell integration. Native Windows support is out of scope for now.
- When adding dependencies run `go mod tidy` locally and commit `go.sum` so CI caching works correctly.

## Next technical tasks (short)

- Harden the native zsh fallback (POSTDISPLAY-only, debounce, accurate inside-quote detection).
- Implement `uninstall-hook` with backup/restore semantics.
- Add a `commitgen doctor` command to validate environment, snippet sourcing, hook presence, and binary availability.

## File map quick reference

- `MakePrompt(files, patch)` — prompt heuristics: `internal/prompt/prompt.go`.
- `StagedChanges(limit)` — git patch retrieval and truncation: `internal/diff/diff.go`.
- Hook writer — `internal/hook/hook.go`.
- Shell installer — `internal/shell/shell.go`.

This document is intended as an internal developer reference — keep it up to date as the codebase changes. Use it as the source when composing a user-facing README or technical blog posts.

Note: repository `.gitignore` contains `*_test.go` and `e2e/` so test files and the end-to-end folder remain local-only and are not pushed to remotes.

CI: A minimal GitHub Actions workflow lives at `.github/workflows/ci.yml`. It runs a gofmt check and unit tests on Ubuntu for push/PR. macOS and the e2e integration job are exposed as manual workflows (`workflow_dispatch`) to avoid slowing down routine PRs.
# commitgen — Technical Reference

This document explains the code layout, the runtime behavior of the CLI and hook, and developer-facing details useful when extending or packaging the project.

## Repository layout (key files)

## Commands (what they do)
  - `--plain`: prints exactly one trimmed subject line to stdout (suitable for shells/hooks). If no staged changes, prints nothing and exits 0.
  - `--verbose`: prints diagnostic output (patch preview, sizes, and the message) to stderr; stdout remains the message for non-plain mode.
  - Without flags, the tool prints a human-facing message to stdout. `--verbose` directs debug to stderr.
  - Note: non-plain CLI runs now exit non-zero when there are no staged files; `--plain` remains silent (exit 0) so hooks may call it safely.
  - All diagnostic and installer messages are written to stderr so `--plain` output remains clean on stdout.

## Suggest flow (high-level)
1. `suggest` calls `internal/diff.StagedChanges(limit)` to get `files []string` and a trimmed `patch string`.
2. `prompt.MakePrompt(files, patch)` applies bucket heuristics:
   - `isTestsOnly` / `isDocsOnly` / `isConfigOnly` / `isRenameOnly`
   - Default: `Update <file1>, <file2> (+N more)`.
3. Output is formatted according to flags (`--plain` vs normal; debug to stderr when `--verbose`).

## Hook script details
  - Skips when the message file already has content, or when the hook is called from merge/squash/rebase flows.
  - Runs `bin/commitgen suggest` and writes the single-line suggestion into the commit message file (if non-empty).

## Shell integration (zsh)
  - If `zsh-autosuggestions` is available, a custom strategy `_zsh_autosuggest_strategy_commitgen` is defined to call `commitgen suggest --plain` and surface the result as ghost text (dimmed) that the plugin accepts with →.
  - Otherwise a native fallback provides a quick POSTDISPLAY-style preview and a `cg-accept-preview` widget; this fallback is intentionally conservative and is a prototype that requires refinement (debounce, accurate inside-quote detection, and robust accept behavior).

## Tests and verification

## Developer notes & best practices

## Next technical tasks (short)

## File map quick reference

This document is intended as an internal developer reference — keep it up to date as the codebase changes. Use it as the source when composing a user-facing README or technical blog posts.
Note: repository `.gitignore` now contains `*_test.go` and `e2e/` so test files and the end-to-end folder remain local-only and are not pushed to remotes.

CI: A minimal GitHub Actions workflow lives at `.github/workflows/ci.yml`. It runs a gofmt check and unit tests on Ubuntu for push/PR. macOS and the e2e integration job are exposed as manual workflows (workflow_dispatch) to avoid slowing down routine PRs.
