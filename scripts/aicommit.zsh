# CommitGen helper for zsh
# Source this file from your ~/.zshrc to add the `aicommit` command.
# Requires `commitgen` on PATH (or export COMMITGEN_CLI=/path/to/commitgen).

aicommit() {
  emulate -L zsh
  setopt errexit pipefail nounset

  local commitgen_cli=${COMMITGEN_CLI:-commitgen}
  if ! command -v "$commitgen_cli" >/dev/null 2>&1; then
    print -u2 "aicommit: $commitgen_cli not found on PATH"
    return 1
  fi

  local msg
  if ! msg="$("$commitgen_cli" suggest --ai --plain 2> >(cat >&2))"; then
    print -u2 "aicommit: $commitgen_cli suggest failed; aborting git commit"
    return 1
  fi

  msg="${msg//$'\n'/ }"
  msg="${msg## }"
  msg="${msg%% }"
  if [[ -z "$msg" ]]; then
    print -u2 "aicommit: empty message from $commitgen_cli; aborting git commit"
    return 1
  fi

  git commit -m "$msg" "$@"
}

# Reuse git commit completion so `aicommit <TAB>` behaves identically.
if whence compdef >/dev/null 2>&1; then
  autoload -Uz compinit && compinit
  compdef _git aicommit=git-commit
fi
