package shell

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	snippetRelPath = ".config/commitgen.zsh"
	guardStart     = "# >>> commitgen >>> (managed)"
	guardEnd       = "# <<< commitgen <<<"
)

func InstallShell() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(home, snippetRelPath)
	cfgDir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return err
	}

	snippet := pluginFirstSnippet()

	if err := os.WriteFile(cfgPath, []byte(snippet), 0o644); err != nil {
		return err
	}

	zshrcPath := filepath.Join(home, ".zshrc")
	zshrcBytes, _ := os.ReadFile(zshrcPath)
	zshrc := string(zshrcBytes)

	if containsGuard(zshrc) {
		return nil
	}

	block := fmt.Sprintf("%s\n[[ -f \"%s\" ]] && source \"%s\"\n%s\n", guardStart, cfgPath, cfgPath, guardEnd)

	f, err := os.OpenFile(zshrcPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + block); err != nil {
		return err
	}

	return nil
}

func UninstallShell() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(home, snippetRelPath)
	zshrcPath := filepath.Join(home, ".zshrc")

	zshrcBytes, err := os.ReadFile(zshrcPath)
	if err != nil {
		_ = os.Remove(cfgPath)
		return nil
	}
	zshrc := string(zshrcBytes)

	new, changed := removeGuardedBlock(zshrc)
	if changed {
		if err := os.WriteFile(zshrcPath, []byte(new), 0o644); err != nil {
			return err
		}
	}

	if err := os.Remove(cfgPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func containsGuard(s string) bool {
	return filepathHas(s, guardStart) && filepathHas(s, guardEnd)
}

func filepathHas(s, sub string) bool {
	return len(s) > 0 && (stringIndex(s, sub) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func removeGuardedBlock(s string) (string, bool) {
	start := stringIndex(s, guardStart)
	if start < 0 {
		return s, false
	}
	end := stringIndex(s[start:], guardEnd)
	if end < 0 {
		return s, false
	}
	end += start + len(guardEnd)
	new := s[:start]
	if end < len(s) {
		new += s[end:]
	}
	return new, true
}

func pluginFirstSnippet() string {
	return `# commitgen zsh snippet (plugin-first)
# plugin strategy for zsh-autosuggestions
_cg_autosuggest_available() {
  # detect autosuggestions plugin
  if [[ -n ${ZSH_AUTOSUGGEST_DIR-} ]]; then
    return 0
  fi
  if [[ -f ${ZSH:-$HOME/.oh-my-zsh}/custom/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh ]]; then
    return 0
  fi
  return 1
}

if _cg_autosuggest_available; then
  _zsh_autosuggest_strategy_commitgen() {
    # only run for exact patterns: git commit -m " or gc "
    [[ $BUFFER == "git commit -m \""* ]] || [[ $BUFFER == "gc \""* ]] || return 1
    commitgen suggest --ai --plain 2>/dev/null
  }
  # ensure strategy order: commitgen first, then history
  if [[ -n "${ZSH_AUTOSUGGEST_STRATEGY-}" ]]; then
    export ZSH_AUTOSUGGEST_STRATEGY=(commitgen ${ZSH_AUTOSUGGEST_STRATEGY[@]})
  else
    export ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)
  fi
else
  # native POSTDISPLAY preview fallback
  _cg_update_preview() {
    [[ $BUFFER == "git commit -m \""* ]] || [[ $BUFFER == "gc \""* ]] || return 1
    # extract inside-quotes content
    local inside=${BUFFER#*\"}
    inside=${inside%%\"*}
    local sug
    sug=$(commitgen suggest --ai --plain 2>/dev/null || true)
    [[ -z $sug ]] && return 1
    if [[ $inside == "$sug" ]]; then
      return 1
    fi
    zle -M "$sug"
  }
  function cg-accept-preview() {
    # insert suggestion into the quoted message
    local before=${BUFFER%%\"*}\"
    local after=
    # replace or append
    BUFFER=${BUFFER%%\"*}\"$(commitgen suggest --ai --plain 2>/dev/null)\"${BUFFER#*\"}
    zle reset-prompt
  }
  zle -N cg-accept-preview
  autoload -Uz add-zsh-hook
  add-zsh-hook -Uz preexec _cg_update_preview
  bindkey -M emacs '^F' cg-accept-preview
  bindkey -M viins '^F' cg-accept-preview
  bindkey '^[[C' cg-accept-preview
fi
`
}
