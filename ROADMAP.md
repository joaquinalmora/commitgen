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

**Inline suggestions status (today)**
- `suggest --plain` implemented and returning single-line subject.
- Prototype zsh integration present: custom strategy function for `zsh-autosuggestions` and native zsh POSTDISPLAY widget attempted.
- Current behavior observed: in native zsh path the suggestion is being inserted directly instead of shown as ghost-text; needs refinement so it stays as preview until accepted (key-binding).

**CLI plumbing**
- Minimal command registry and `--help`/usage printer.

---

## ⚠️ Gaps / known issues

- Inline suggestions prototype exists, but native zsh path currently inserts text immediately instead of previewing; needs POSTDISPLAY-only preview with explicit accept binding.
- **Hook uninstall** is missing (`uninstall-hook`).
- **Doctor** command missing (env checks).
- `install-hook` assumes **repo root**; doesn’t search parent dirs for `.git`.
- `suggest` must print **only** the message (no debug lines) for hook usage.
- No **unit tests** for prompt heuristics.
- No **AI integration** (LLM) yet.
- No **packaging** (Homebrew tap, release binaries) or CI.
- No **config flags/env** (e.g., `--max-files`, `--patch-bytes`, enable/disable buckets).
- Shell setup automation missing: need safe, idempotent installer that writes a separate snippet (e.g., ~/.config/commitgen.zsh) and sources it from ~/.zshrc, with uninstall.
- Cross-shell coverage not implemented: bash fallback helpers and fish integration TBD.

---

## 🎯 Goals (what “done” looks like)

1) **Inline suggestions in terminal**
   - **With oh-my-zsh + zsh-autosuggestions:** custom strategy that supplies `commitgen suggest` text when typing `git commit -m "`. Accept with →.
   - **Without plugins (plain zsh):** a `zle-line-pre-redraw` widget that paints a ghost suggestion via `POSTDISPLAY` and accepts with → by binding a custom widget.
   - preview via POSTDISPLAY, accept with a bound key (Right Arrow or Ctrl+F), never auto-insert.
   - Bash/Fish: provide simple one-liners and optional shell snippets (bash completions/fish functions) to surface suggestions or integrate with their UX.

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

6) **Guided install & uninstall**
   - `commitgen init shell` writes `~/.config/commitgen.zsh` and appends a guarded `source` line to `~/.zshrc` if missing; supports `--ohmyzsh` and `--native`.
   - `commitgen uninstall shell` removes the guarded block and the snippet.
   - `commitgen doctor` verifies repo, hook, binary on PATH, zsh snippet sourced, and staged-change path.

---

## 📌 Immediate next steps (bite-sized)

1. **Finalize `--plain` output**: ensure zero extra prints in plain mode (already working; add test).
2. **Fix native zsh preview**: adjust widget so it only sets `POSTDISPLAY` (no insertion), and bind accept to a key (Right Arrow or Ctrl+F). Add minimal debounce.
3. **Harden autosuggestions strategy**: when buffer matches commit prefix, return empty suggestion with success to block history; ensure strategy runs first.
4. **Add `commitgen init shell`**: generate `~/.config/commitgen.zsh` with either autosuggestions strategy or native widget; append a single guarded `# >>> commitgen >>>` block to `~/.zshrc` to source it. Make idempotent.
5. **Add `commitgen uninstall shell`**: remove the guarded block and snippet. Idempotent.
6. **Add `commitgen doctor`**: checks (in git repo, binary on PATH, hook executable, snippet sourced, staged changes). Non‑zero exit on failure.
7. **Add `uninstall-hook`** and improve `install-hook` (backup/restore; walk to .git).
8. **README updates**: quickstart, inline setup (oh-my-zsh + native zsh), troubleshooting, uninstall.

---

- [ ] Confirm zsh-autosuggestions plugin is installed in ~/.oh-my-zsh/custom/plugins
- [ ] Re-enable in plugins=(git zsh-autosuggestions zsh-syntax-highlighting)
- [ ] Verify commitgen strategy shows ghost text in dim style
- [ ] If plugin missing, run git clone https://github.com/zsh-users/zsh-autosuggestions

## 🧭 Stretch tasks (when the core is solid)

- **AI integration** with a provider adapter (OpenAI / local Ollama).
- **Config file** (e.g., `.commitgen.yaml`) to tweak buckets, limits.
- **Multiplatform hook script** considerations (currently POSIX `sh` is fine for macOS/Linux).
- Bash helpers: function/alias to insert `$(commitgen suggest --plain)` and optional readline bindings.
- Fish integration: a simple fish function to preview/accept suggestions.
- CI: cross‑shell smoke tests using docker images for zsh/bash/fish.

---

## ✅ Definition of done

- Hook installs/uninstalls cleanly; `doctor` passes.
- `suggest` is deterministic and single-line.
- Inline suggestions work:
  - With **oh-my-zsh** (zsh-autosuggestions).
  - Native zsh shows preview (POSTDISPLAY), requires explicit accept, and does not auto-insert.
- Install/uninstall and doctor commands succeed end‑to‑end on a fresh machine.
- README covers install, usage, inline setup, troubleshooting.
- Tests pass in CI; release binaries available.

---

## 💤 Parking note (today’s status)
- CLI: `commitgen suggest --plain` works and is the source of the ghost text (verified).
- Native zsh inline preview: functional; accepts with → and Ctrl+F; preview color not reliably dim in this terminal/theme.
- zsh-autosuggestions: NOT enabled by default; when enabled, visuals look great using custom `commitgen` strategy.
- Current shell set to **native only** (no plugin).

### When I come back
1) **Decide primary UX**
   - Option A: Make **autosuggestions the default** (prettier ghost text); provide one-command install (`commitgen init shell`) that clones plugin if missing and wires custom strategy.
   - Option B: Keep **native** as default; continue trying to enforce `zle_highlight` styling (standout/bold,fg=242) or accept plain preview.
2) Implement `commitgen init shell` + `uninstall shell` + `doctor`.
3) (Optional) Add AI provider after UX is solid.

### Quick test
```bash
echo "test $(date)" >> demo.txt && git add demo.txt
git commit -m "
# expect: preview from commitgen, accept with → or Ctrl+F, no re-suggest after insert

