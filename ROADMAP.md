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

**Commands & recent progress**
- `commitgen suggest`
  - Collects staged files + trimmed patch (`git diff --cached --unified=0`).
  - Heuristics:
    - **docs-only** → “Update documentation”
    - **tests-only** → “Update tests”
    - **config-only** → “Update configuration”
    - **rename-only** (diff similarity 100%, no hunks) → “Rename files”
    - Else → “Update <file1>, <file2> (+N more)”
  - New: `--plain` mode prints exactly one trimmed subject line to stdout (no debug output).
  - New: `--verbose` flag prints debug info to stderr (keeps stdout clean for hook usage).
- `commitgen install-hook`
  - Writes a **prepare-commit-msg** hook to `.git/hooks/prepare-commit-msg` (installer improved: exported `InstallHook` in `internal/hook`).
  - Hook calls the stable binary; design assumes `bin/commitgen` path but installer will be hardened next.

**Inline suggestions status (today)**
- `suggest --plain` implemented and has a unit test (`internal/prompt/prompt_test.go`) asserting single-line output.
- CLI changes made in `cmd/commitgen/main.go`: plain/verbose logic, gated debug output; `prompt.MakePrompt(files, patch)` remains the source of the subject.
- `internal/diff/StagedChanges` already returns a trimmed patch (`--unified=0`) and truncates to a configurable byte limit.
- Prototype zsh integration for both `zsh-autosuggestions` strategy and native `POSTDISPLAY` exists, but native preview needs refinement (should only preview, not insert).

**Installer test (recent)**
- Added `commitgen init-shell` / `commitgen uninstall-shell` commands and `internal/shell` package.
- Actions taken:
  - Built `bin/commitgen` locally.
  - Ran `commitgen init-shell`: it created `~/.config/commitgen.zsh` and appended a guarded block sourcing it to `~/.zshrc`.
  - Ran `commitgen uninstall-shell`: it removed the snippet file (`~/.config/commitgen.zsh`).
- Observations & caveats:
  - `uninstall-shell` removed the `~/.config/commitgen.zsh` file but did not (and should not) remove unrelated, manually-added commitgen lines that were already present in `~/.zshrc` prior to the install run. Review `~/.zshrc` for any legacy commitgen snippets and remove them manually if desired.
  - Unit tests passed (`internal/prompt`), and the code was built successfully.

**CLI plumbing**
- Minimal command registry and `--help`/usage printer. Recent edits include clearer stdout/stderr separation to keep `--plain` safe for shell consumption.

---

## ⚠️ Gaps / known issues

- Native zsh preview currently inserts text in some setups; it must be changed to a POSTDISPLAY-only preview and require explicit accept (Right Arrow / Ctrl+F).
- `uninstall-hook` is still missing and `install-hook` should be hardened to search repo root (`git rev-parse --show-toplevel`) and backup existing hooks.
- `commitgen doctor` command not implemented yet (env and install verification).
- Shell setup automation (idempotent `commitgen init shell` / `uninstall shell`) is missing.
- No AI/LLM provider integration yet (`internal/provider` is a placeholder).
- Packaging & CI (build/test workflows, release artifacts) not implemented.
- Additional config flags (e.g. `--max-files`, `--patch-bytes`) not yet exposed.

---

## 🎯 Goals (what “done” looks like)

1) **Inline suggestions in terminal**
  - **Plugin-first (preferred):** `zsh-autosuggestions` strategy that returns commit suggestions as ghost text (dimmed) and accepts with →.
  - **Native fallback:** POSTDISPLAY-based preview that never auto-inserts; requires explicit accept (Right Arrow / Ctrl+F). Native should only show preview when appropriate (inside quotes, or when empty if configured).
  - Ensure `commitgen suggest --plain` remains the canonical single-line source for suggestions used by both strategies.

2) **Robust Git hook lifecycle**
  - `install-hook`: backup existing hook or abort with message; write executable script. (`internal/hook` updated to export `InstallHook`.)
  - `uninstall-hook`: implement restore/remove behavior.
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
  - `commitgen init shell` should write `~/.config/commitgen.zsh` and append a guarded `source` block to `~/.zshrc` (idempotent). Support `--ohmyzsh` and `--native` options.
  - `commitgen uninstall shell` should remove the guarded block and snippet.
  - `commitgen doctor` should verify repo, hook, binary on PATH, zsh snippet sourced, and staged-change path.

---

## 📌 Immediate next steps (bite-sized)

1. Confirm `--plain` behavior (done):
  - `commitgen suggest --plain` prints exactly one trimmed line to stdout. Unit test added in `internal/prompt/prompt_test.go`.
2. Fix native zsh preview (high priority):
  - Ensure the native widget only sets `POSTDISPLAY` and does not insert text.
  - Add explicit accept widget (`cg-accept-preview`) bound to → / Ctrl+F.
3. Clean up legacy .zshrc entries (short):
  - Because some users may have earlier manual snippets, add a small audit step to `commitgen doctor` or improve `uninstall-shell` to detect and remove only the guarded block created by `init-shell` (avoid touching other manual content).
3. Implement `commitgen init shell` (idempotent installer):
  - Create `~/.config/commitgen.zsh` with plugin strategy and native fallback. Append guarded block to `~/.zshrc`.
4. Implement `commitgen uninstall shell` and `uninstall-hook`.
5. Implement `commitgen doctor` to validate installs and environment.
6. Add README docs for inline setup and uninstall.
7. Add CI and packaging (next sprint).

---

 - [ ] Confirm zsh-autosuggestions plugin is installed in `~/.oh-my-zsh/custom/plugins` (optional if using plugin-first UX)
 - [ ] Re-enable in `plugins=(...)` or let `commitgen init shell --ohmyzsh` handle it
 - [ ] Verify `commitgen` strategy shows ghost text in dim style when plugin present

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
# expect: preview from commitgen; accept with → or Ctrl+F; no re-suggest after insert
```

---

## 🔭 Decision: default UX policy
- **Default to plugin visuals when available**: if `zsh-autosuggestions` is present, use the custom `commitgen` strategy for dim “ghost” text and familiar → acceptance.
- **Native fallback otherwise**: ship a `POSTDISPLAY`-based preview that works everywhere (even if some terminals don’t dim the ghost). Accept with → / Ctrl+F.

## 🧰 Installer (`commitgen init shell`) — specification
- Detect whether `zsh-autosuggestions` is installed (`$ZSH/custom/plugins/zsh-autosuggestions/...`).
- Generate `~/.config/commitgen.zsh`:
  - **Plugin path**: define `_zsh_autosuggest_strategy_commitgen`, set `ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)`; do **not** load history for other commands.
  - **Native path**: define `_cg_update_preview` (no auto-insert; compares inside-quotes to avoid re-suggest), `cg-accept-preview` (computes on demand), keybinds for `^[[C` (→), `^F` (Ctrl+F), `^[[F` (End). Include optional `COMMITGEN_PREVIEW_ONLY_EMPTY` toggle.
- Append a **guarded** block to `~/.zshrc` (idempotent):
  ```
  # >>> commitgen >>> (managed)
  [[ -f ~/.config/commitgen.zsh ]] &amp;&amp; source ~/.config/commitgen.zsh
  # <<< commitgen <<<
  ```
- Options:
  - `--native` (force native even if plugin exists)
  - `--ohmyzsh` (force plugin path; clone if missing and user agrees)
- Safety: do not duplicate blocks; create a timestamped backup of `.zshrc` before first edit.

## ♻️ Uninstall &amp; verification
- **`commitgen uninstall shell`**: remove the guarded block from `~/.zshrc` and delete `~/.config/commitgen.zsh`. Idempotent.
- **`commitgen doctor`** (non-zero on failure):
  - In a git repo (`git rev-parse --is-inside-work-tree`).
  - Binary on PATH (`command -v commitgen`).
  - Hook present/executable (if installed).
  - Shell snippet sourced (detect via an env marker the snippet sets).
  - Staged changes are visible (`git diff --cached --name-only` not empty).
  - If plugin path: verify `_zsh_autosuggest_start` exists and strategy includes `commitgen`.

## 🧪 Native preview constraints &amp; toggles
- `POSTDISPLAY` styling is theme/terminal-dependent; some setups won’t render dim grey. Functionality remains correct.
- Prevent repeats by comparing what’s **inside quotes** with the suggestion; if equal, do not set `POSTDISPLAY`.
- Optional calmer mode:
  ```zsh
  # Only show when the message field is empty
  export COMMITGEN_PREVIEW_ONLY_EMPTY=1
  ```

## 🎹 Key bindings (standardize)
- Bind in emacs and vi keymaps:
  ```zsh
  bindkey -M emacs '^F'  cg-accept-preview
  bindkey -M viins  '^F'  cg-accept-preview
  bindkey -M vicmd  '^F'  cg-accept-preview
  bindkey '^[[C'    cg-accept-preview   # →
  bindkey '^[[F'    cg-accept-preview   # End
  bindkey '^[OC'    cg-accept-preview   # → (application cursor mode)
  ```

## 🛠️ Troubleshooting playbook (quick commands)
- **Am I in a git repo / do I have staged files?**
  ```zsh
  git rev-parse --is-inside-work-tree
  git diff --cached --name-only
  ```
- **Is commitgen working?**
  ```zsh
  commitgen suggest --plain | hexdump -C
  ```
- **Which path is active right now?**
  ```zsh
  echo "autosuggest_loaded=${+functions[_zsh_autosuggest_start]}  pre_redraw=${+widgets[zle-line-pre-redraw]}  accept=${+widgets[cg-accept-preview]}"
  ```
- **Nuke half-loaded autosuggestions in this shell (errors like \`_zsh_autosuggest_fetch\`):**
  ```zsh
  unfunction -m '_zsh_autosuggest_*' 2>/dev/null; for w in zle-line-pre-redraw zle-keymap-select zle-line-init zle-line-finish; do zle -D $w 2>/dev/null; done; bindkey -e
  ```

## ⏱️ First 30 minutes when you’re back
1) Decide default path: plugin-first (recommended) or native-first.
2) Implement `commitgen init shell` (idempotent) + `uninstall shell`.
3) Implement `commitgen doctor` with the checks above.
4) Add README “Inline setup” with both paths, plus uninstall and doctor.
5) Add a small test: `suggest --plain` must emit exactly 1 line; unit test prompt buckets.

## ✅ Parking sanity checklist
- Native preview inserts and does not re-suggest after accept.
- Plugin visuals (if enabled) show dim ghost, accept with →.
- `commitgen suggest --plain` prints a single line for staged changes.
- `install-hook`/`uninstall-hook` lifecycle covered.
- `doctor` reports green on a fresh shell.
