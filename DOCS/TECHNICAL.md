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
