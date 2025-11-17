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
	return `# commitgen zsh snippet (native ghost text)
typeset -g _CG_BIN_PATH=""
typeset -g _CG_PREVIEW_INIT=0

_cg_find_bin() {
  if [[ -n ${COMMITGEN_BIN-} && -x ${COMMITGEN_BIN} ]]; then
    _CG_BIN_PATH="${COMMITGEN_BIN}"
    return 0
  fi

  if [[ -n "$_CG_BIN_PATH" && -x "$_CG_BIN_PATH" ]]; then
    return 0
  fi

  if _cg_bin=$(command -v commitgen 2>/dev/null); then
    _CG_BIN_PATH="$_cg_bin"
    return 0
  fi

  local repo_root
  repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || repo_root=""
  if [[ -n "$repo_root" && -x "$repo_root/bin/commitgen" ]]; then
    _CG_BIN_PATH="$repo_root/bin/commitgen"
    return 0
  fi

  return 1
}

_cg_run_commitgen() {
  _cg_find_bin || return 1
  "$_CG_BIN_PATH" "$@"
}

_cg_fetch_suggestion() {
  local suggestion
  suggestion=$(_cg_run_commitgen cached --plain 2>/dev/null) || true
  if [[ -z "$suggestion" ]]; then
    suggestion=$(_cg_run_commitgen suggest --ai --plain 2>/dev/null) || true
  fi
  [[ -n "$suggestion" ]] && echo "$suggestion"
}

_cg_match_prefix() {
  local left="$LBUFFER"
  case "$left" in
    'git commit -m "'*)
      printf 'git commit -m "'
      return 0
      ;;
    'gc "'*)
      printf 'gc "'
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

_cg_update_preview_widget() {
  local prefix
  prefix=$(_cg_match_prefix) || { POSTDISPLAY=; return; }

  local after=${LBUFFER#${prefix}}
  if [[ "$after" == *\"* ]]; then
    POSTDISPLAY=
    return
  fi

  local suggestion
  suggestion=$(_cg_fetch_suggestion)
  if [[ -z "$suggestion" ]]; then
    POSTDISPLAY=
    return
  fi

  local typed="${after}"
  if [[ "$typed" == "$suggestion" ]]; then
    POSTDISPLAY=
    return
  fi

  POSTDISPLAY=${suggestion}\"
}

_cg_accept_preview_widget() {
  if [[ -z "$POSTDISPLAY" ]]; then
    _cg_update_preview_widget
  fi

  if [[ -n "$POSTDISPLAY" ]]; then
    LBUFFER+="$POSTDISPLAY"
    POSTDISPLAY=
  fi
  zle -R
}

if [[ $_CG_PREVIEW_INIT -eq 0 ]]; then
  typeset -ga zle_highlight
  if (( ${zle_highlight[(I)special:*]} == 0 )); then
    zle_highlight+=(special:fg=240)
  fi

  zle -N zle-line-pre-redraw _cg_update_preview_widget
  zle -N cg-accept-preview _cg_accept_preview_widget
  bindkey '^F' cg-accept-preview
  bindkey '^[[C' cg-accept-preview
  bindkey '^[[F' cg-accept-preview
  typeset -g _CG_PREVIEW_INIT=1
fi
`
}
