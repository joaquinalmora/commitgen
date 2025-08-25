# Inline Terminal Suggestions (oh-my-zsh **and** plain zsh)

You can get **ghost text** suggestions for commit messages while typing
```
git commit -m "<cursor>"
```
and accept with **Right Arrow (→)**. Two setups are supported:

- ✅ **oh-my-zsh + zsh-autosuggestions** (recommended if you already use it)
- ✅ **plain zsh (no plugins)** using ZLE widgets

> The Git hook continues to work either way. Inline suggestions are an extra UX layer.

---

## 0) Put `commitgen` on PATH

```bash
# from your repo
go build -o ~/bin/commitgen ./cmd/commitgen
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Ensure it prints exactly one line:
```bash
commitgen suggest
```

---

## 1) oh-my-zsh + zsh-autosuggestions

### Check the plugin is enabled
```zsh
grep -n '^plugins=' ~/.zshrc
# Expect: plugins=(git zsh-autosuggestions ...)
```

### Add config (before sourcing oh-my-zsh)
Edit `~/.zshrc`, above `source $ZSH/oh-my-zsh.sh`:
```zsh
# commitgen autosuggestion: enable async and set our strategy first
ZSH_AUTOSUGGEST_USE_ASYNC=1
ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)
```

### Define the strategy (after sourcing oh-my-zsh)
Place this after `source $ZSH/oh-my-zsh.sh`:
```zsh
# Suggest commit message inline when typing: git commit -m "<cursor>"
_zsh_autosuggest_strategy_commitgen() {
  # We check the left buffer (what you've typed so far) for an opened quote
  if [[ $LBUFFER == 'git commit -m "'* || $LBUFFER == "git commit -m '"* ]]; then
    local sug
    # Use 'commitgen' if it's on PATH; otherwise put the absolute path here
    sug=$(commitgen suggest 2>/dev/null)
    [[ -n $sug && $sug != "No staged files" ]] || return 1
    suggestion=$sug
    return 0
  fi
  return 1
}
```

Reload:
```zsh
zsh -n ~/.zshrc   # syntax check; no output = good
source ~/.zshrc
```

**Use it:** stage a change, type `git commit -m "` and a faded inline suggestion appears; press **→** to accept.

### Troubleshooting
- `~/.zshrc:NN: unmatched ' or "` → use the snippet above (no regex).  
- Strategy not applied? Ensure `ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)` is **before** sourcing oh-my-zsh, and the function is **after**.
- No suggestion? Ensure staged files exist and `commitgen suggest` prints a single line.

---

## 2) Plain zsh (no plugins)

We’ll use a ZLE hook to paint ghost text with `POSTDISPLAY` and a custom Right‑Arrow key that accepts the suggestion.

Add this to `~/.zshrc`:

```zsh
# --- commitgen inline (plain zsh) ---
# Compute suggestion only when typing: git commit -m "<cursor>"
_commitgen_inline_compute() {
  POSTDISPLAY=""
  if [[ $LBUFFER == 'git commit -m "'* || $LBUFFER == "git commit -m '"* ]]; then
    local sug
    # Use absolute path if commitgen isn't on PATH
    sug=$(commitgen suggest 2>/dev/null)
    [[ -n $sug && $sug != "No staged files" ]] || return
    # Show as ghost text
    POSTDISPLAY="$sug"
  fi
}

# Accept suggestion on Right Arrow: insert POSTDISPLAY into buffer if present
_commitgen_accept_or_right() {
  if [[ -n $POSTDISPLAY ]]; then
    LBUFFER+="$POSTDISPLAY"
    POSTDISPLAY=""
    zle redisplay
  else
    zle forward-char   # normal Right Arrow
  fi
}
zle -N _commitgen_inline_compute
zle -N _commitgen_accept_or_right

# Run our compute hook before each redraw
autoload -Uz add-zle-hook-widget
add-zle-hook-widget line-pre-redraw _commitgen_inline_compute

# Bind Right Arrow to accept-or-move
bindkey -M emacs '^[[C' _commitgen_accept_or_right   # Right Arrow
bindkey -M vicmd '^[[C' _commitgen_accept_or_right
bindkey -M viins '^[[C' _commitgen_accept_or_right
# --- end commitgen inline ---
```

Reload:
```zsh
zsh -n ~/.zshrc   # syntax check
source ~/.zshrc
```

**Use it:** stage a change, type `git commit -m "`, see ghost text; press **→** to accept.

### Notes
- This doesn’t require oh-my-zsh or plugins.
- It computes suggestions on each redraw; if you want to throttle, add a tiny cache keyed by repo state (see below).

---

## Optional: micro-cache for performance (both setups)

```zsh
typeset -gA _COMMITGEN_CACHE

_commitgen_key() {
  local root head tree
  root=$(git rev-parse --show-toplevel 2>/dev/null) || return 1
  head=$(git rev-parse --short HEAD 2>/dev/null)
  tree=$(git write-tree 2>/dev/null)
  print -- "$root|$head|$tree"
}
```

**oh-my-zsh strategy**: compute `key=$(_commitgen_key)` and reuse `$_COMMITGEN_CACHE[$key]`.  
**plain zsh**: same inside `_commitgen_inline_compute`.

---

## Uninstall (either setup)

- **oh-my-zsh**: remove the `ZSH_AUTOSUGGEST_*` lines and the `_zsh_autosuggest_strategy_commitgen` function; then `source ~/.zshrc`.
- **plain zsh**: remove the block, or run:
  ```zsh
  add-zle-hook-widget -d line-pre-redraw _commitgen_inline_compute
  bindkey -M emacs '^[[C' forward-char
  bindkey -M vicmd '^[[C' forward-char
  bindkey -M viins '^[[C' forward-char
  ```

---

## Coexistence

- The Git **hook** keeps working (pre-fills editor or `-m ""`).
- Inline suggestions trigger only for `git commit -m "` / `'`, everything else remains as-is.
- With oh-my-zsh, we set `(commitgen history)` so history suggestions continue to work for all other commands.
