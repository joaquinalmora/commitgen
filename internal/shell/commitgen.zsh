#!/usr/bin/env zsh
# commitgen zsh snippet (plugin-first)

# Detect zsh-autosuggestions by a common env var or path
_cg_autosuggest_available() {
  [[ -n ${ZSH_AUTOSUGGEST_DIR-} ]] && return 0
  [[ -f ${ZSH:-$HOME/.oh-my-zsh}/custom/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh ]] && return 0
  return 1
}

if _cg_autosuggest_available; then
  # Provide a simple strategy function for zsh-autosuggestions
  _zsh_autosuggest_strategy_commitgen() {
    # run when typing a git commit with -m
    [[ $BUFFER == *git\ commit* && $BUFFER == *-m* ]] || return 1
    commitgen suggest --plain 2>/dev/null
  }

  # Prepend commitgen to the strategy list so it appears first
  if [[ -n ${ZSH_AUTOSUGGEST_STRATEGY-} ]]; then
    ZSH_AUTOSUGGEST_STRATEGY=(commitgen ${ZSH_AUTOSUGGEST_STRATEGY[@]})
  else
    ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)
  fi

else
  # Fallback: show a dim preview message (zle -M) and bind an accept key
  _cg_update_preview() {
    [[ $BUFFER == git\ commit* && $BUFFER == *-m\ \"* ]] || return 1
    # Extract inside-quotes content
    local inside=${BUFFER#*\"}
    inside=${inside%%\"*}
    local sug
    sug=$(commitgen suggest --plain 2>/dev/null || true)
    [[ -z $sug ]] && return 1
    [[ $inside == "$sug" ]] && return 1
    zle -M "$sug"
  }

  cg-accept-preview() {
    local sug
    sug=$(commitgen suggest --plain 2>/dev/null || true)
    [[ -z $sug ]] && return 1
    
    if [[ $BUFFER == *\"*\"* ]]; then
      # Buffer already has quotes, replace content between them
      local prefix=${BUFFER%%\"*}
      local rest=${BUFFER#*\"}
      local after=${rest#*\"}
      BUFFER="${prefix}\"${sug}\"${after}"
    elif [[ $BUFFER == *\"* ]]; then
      # Buffer has opening quote, close it with suggestion
      local prefix=${BUFFER%%\"*}
      BUFFER="${prefix}\"${sug}\""
    elif [[ $BUFFER == *-m* ]]; then
      # Has -m flag but no quotes yet
      local prefix=${BUFFER%%-m*}
      BUFFER="${prefix}-m \"${sug}\""
    else
      # No -m flag, append it with suggestion
      BUFFER="${BUFFER} -m \"${sug}\""
    fi
    zle reset-prompt
  }

  zle -N cg-accept-preview
  autoload -Uz add-zsh-hook
  add-zsh-hook -Uz preexec _cg_update_preview
  bindkey -M emacs '^F' cg-accept-preview
  bindkey -M viins '^F' cg-accept-preview
  bindkey '^[[C' cg-accept-preview
fi