# commitgen — Project Roadmap

This document captures **what’s built today**, **what’s missing**, and a **step‑by‑step plan** to get to the “polished, resume‑ready” version with inline terminal suggestions and clean installation.

---

## ✅ Current state (what you have)

**Language & layout**
- Go module with standard layout:
  - `cmd/commitgen/` (CLI entrypoint)
  - `internal/diff/` (reads staged files & diff)
  - `internal/prompt/` (heuristics for message text)
  - `internal/hook/` (Git hook installer)

**Commands**
- `commitgen suggest`
  - Collects staged files + trimmed patch (`git diff --cached --unified=0`).
  - Heuristics:
    - **docs-only** → “Update documentation”
    - **tests-only** → “Update tests”
    - **config-only** → “Update configuration”
    - **rename-only** (diff similarity 100%, no hunks) → “Rename files”
    - Else → “Update <file1>, <file2> (+N more)”
- `commitgen install-hook`
  - Writes a **prepare-commit-msg** hook to `.git/hooks/prepare-commit-msg`
  - Calls a **stable binary** at `bin/commitgen` (`binPath`), not `go run` temp path
  - Skips if message already has content; skips merge/squash/rebase
  - Writes suggestion into the message file

**CLI plumbing**
- Minimal command registry and `--help`/usage printer.

---

## ⚠️ Gaps / known issues

- **Inline suggestions (ghost text) not implemented** yet.
- **Hook uninstall** is missing (`uninstall-hook`).
- **Doctor** command missing (env checks).
- `install-hook` assumes **repo root**; doesn’t search parent dirs for `.git`.
- `suggest` must print **only** the message (no debug lines) for hook usage.
- No **unit tests** for prompt heuristics.
- No **AI integration** (LLM) yet.
- No **packaging** (Homebrew tap, release binaries) or CI.
- No **config flags/env** (e.g., `--max-files`, `--patch-bytes`, enable/disable buckets).

---

## 🎯 Goals (what “done” looks like)

1) **Inline suggestions in terminal**
   - **With oh-my-zsh + zsh-autosuggestions:** custom strategy that supplies `commitgen suggest` text when typing `git commit -m "`. Accept with →.
   - **Without plugins (plain zsh):** a `zle-line-pre-redraw` widget that paints a ghost suggestion via `POSTDISPLAY` and accepts with → by binding a custom widget.

2) **Robust Git hook lifecycle**
   - `install-hook`: backup existing hook or abort with message; write executable script.
   - `uninstall-hook`: restore backup or remove hook.
   - `doctor`: verify repo, hook presence/permissions, binary path, staged files present.

3) **Quality-of-life**
   - Clean user-facing README including:
     - What it does, quickstart, examples
     - Inline suggestions setup (oh-my-zsh + native zsh)
     - Troubleshooting
   - Configurable behavior (env/flags).

4) **Optionally: AI mode**
   - Provider interface (`internal/ai`), prompt built from `files + trimmed patch`.
   - Toggle via `--ai` or `COMMITGEN_AI=1` with fallback to heuristics.

5) **Packaging**
   - Makefile, GitHub Actions CI (build/test), release artifacts.
   - Homebrew tap (`brew install joaquinalmora/tap/commitgen`).

---

## 📌 Immediate next steps (bite-sized)

1) **Finish inline suggestions**
   - (A) **oh-my-zsh path**: add autosuggestion strategy (see `INLINE_SUGGESTIONS.md`).
   - (B) **native zsh path**: add `~/.zshrc` widget using `POSTDISPLAY` + custom Right-Arrow handler.

2) **Add uninstall & doctor**
   - `uninstall-hook`: check `.git/hooks/prepare-commit-msg`, remove (or restore `.bak`). Print result.
   - `doctor`: print a checklist (in repo? hook executable? `bin/commitgen` exists? staged files?). Exit non‑zero on failure.

3) **Polish suggest output**
   - Ensure `suggest` prints **exactly one line**. No lengths, no previews.

4) **Walk-up .git**
   - In `InstallHook()`, optionally walk parent dirs until you find `.git` so users can run from subfolders.

5) **Tests**
   - Unit tests for `isDocsOnly`, `isTestsOnly`, `isConfigOnly`, `isRenameOnly`.
   - Golden tests for the summary line (e.g., 1 file, N files).

---

## 🧭 Stretch tasks (when the core is solid)

- **AI integration** with a provider adapter (OpenAI / local Ollama).
- **Config file** (e.g., `.commitgen.yaml`) to tweak buckets, limits.
- **Multiplatform hook script** considerations (currently POSIX `sh` is fine for macOS/Linux).

---

## ✅ Definition of done

- Hook installs/uninstalls cleanly; `doctor` passes.
- `suggest` is deterministic and single-line.
- Inline suggestions work:
  - With **oh-my-zsh** (zsh-autosuggestions).
  - Without plugins (native zsh widget).
- README covers install, usage, inline setup, troubleshooting.
- Tests pass in CI; release binaries available.
