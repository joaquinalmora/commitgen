package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func InstallHook() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get current directory:", err)
		return
	}

	binPath := resolveBinaryPath(cwd)

	installPrepareCommitMsgHook(cwd, binPath)
	installPostIndexChangeHook(cwd, binPath)
}

func installPrepareCommitMsgHook(cwd string, binPath string) {
	hookPath := filepath.Join(cwd, ".git", "hooks", "prepare-commit-msg")

	if _, err := os.Stat(hookPath); err == nil {
		fmt.Fprintln(os.Stderr, "Warning: prepare-commit-msg hook already exists, backing up...")
		backupPath := hookPath + ".backup"
		if err := os.Rename(hookPath, backupPath); err != nil {
			fmt.Fprintln(os.Stderr, "failed to backup existing hook:", err)
			return
		}
	}

	script := fmt.Sprintf(`#!/bin/sh
# commitgen prepare-commit-msg hook
MSG_FILE="$1"
SOURCE="$2"

# Don't override existing messages
if [ -s "$MSG_FILE" ]; then
	exit 0
fi

# Skip for special commits
case "$SOURCE" in
merge|squash|rebase)
	exit 0
	;;
esac

# Try cached message first (instant)
CACHED_MSG=$(%s cached 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$CACHED_MSG" ]; then
	printf '%%s\n' "$CACHED_MSG" > "$MSG_FILE"
	exit 0
fi

# Fallback to real-time generation
SUGGEST_MSG=$(%s suggest 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$SUGGEST_MSG" ] && [ "$SUGGEST_MSG" != "No staged files" ]; then
	printf '%%s\n' "$SUGGEST_MSG" > "$MSG_FILE"
fi
`, binPath, binPath)

	if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write prepare-commit-msg hook:", err)
		return
	}

	fmt.Fprintln(os.Stderr, "prepare-commit-msg hook installed successfully")
}

func installPostIndexChangeHook(cwd string, binPath string) {
	hookPath := filepath.Join(cwd, ".git", "hooks", "post-index-change")

	if _, err := os.Stat(hookPath); err == nil {
		fmt.Fprintln(os.Stderr, "Warning: post-index-change hook already exists, backing up...")
		backupPath := hookPath + ".backup"
		if err := os.Rename(hookPath, backupPath); err != nil {
			fmt.Fprintln(os.Stderr, "failed to backup existing hook:", err)
			return
		}
	}

	script := fmt.Sprintf(`#!/bin/sh
# commitgen post-index-change hook (auto-cache)
# This runs automatically when git add is executed

# Only run if there are staged changes
if git diff --cached --quiet; then
	exit 0
fi

# Generate cache in background (don't slow down git add)
%s cache >/dev/null 2>&1 &
`, binPath)

	if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write post-index-change hook:", err)
		return
	}

	fmt.Fprintln(os.Stderr, "post-index-change hook installed successfully")
	fmt.Fprintln(os.Stderr, "Auto-cache enabled: commit messages will be pre-generated on git add")
}

func UninstallHook() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get current directory:", err)
		return
	}

	uninstallPrepareCommitMsgHook(cwd)
	uninstallPostIndexChangeHook(cwd)
}

func resolveBinaryPath(cwd string) string {
	localBin := filepath.Join(cwd, "bin", "commitgen")
	if info, err := os.Stat(localBin); err == nil {
		if info.Mode()&0o111 != 0 {
			return strconv.Quote(localBin)
		}
	}

	if globalBin, err := exec.LookPath("commitgen"); err == nil {
		return strconv.Quote(globalBin)
	}

	return "commitgen"
}

func uninstallPrepareCommitMsgHook(cwd string) {
	hookPath := filepath.Join(cwd, ".git", "hooks", "prepare-commit-msg")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No prepare-commit-msg hook found")
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "failed to check prepare-commit-msg hook:", err)
		return
	}

	content, err := os.ReadFile(hookPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read hook file:", err)
		return
	}

	hookContent := string(content)
	if !containsCommitgenSignature(hookContent) {
		fmt.Fprintln(os.Stderr, "prepare-commit-msg hook exists but doesn't appear to be created by commitgen")
		return
	}

	if err := os.Remove(hookPath); err != nil {
		fmt.Fprintln(os.Stderr, "failed to remove prepare-commit-msg hook:", err)
		return
	}

	backupPath := hookPath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Rename(backupPath, hookPath); err != nil {
			fmt.Fprintln(os.Stderr, "failed to restore backup hook:", err)
		} else {
			fmt.Fprintln(os.Stderr, "prepare-commit-msg hook removed, backup restored")
			return
		}
	}

	fmt.Fprintln(os.Stderr, "prepare-commit-msg hook removed successfully")
}

func uninstallPostIndexChangeHook(cwd string) {
	hookPath := filepath.Join(cwd, ".git", "hooks", "post-index-change")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "failed to check post-index-change hook:", err)
		return
	}

	content, err := os.ReadFile(hookPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read post-index-change hook file:", err)
		return
	}

	hookContent := string(content)
	if !stringContains(hookContent, "commitgen cache") {
		fmt.Fprintln(os.Stderr, "post-index-change hook exists but doesn't appear to be created by commitgen")
		return
	}

	if err := os.Remove(hookPath); err != nil {
		fmt.Fprintln(os.Stderr, "failed to remove post-index-change hook:", err)
		return
	}

	backupPath := hookPath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Rename(backupPath, hookPath); err != nil {
			fmt.Fprintln(os.Stderr, "failed to restore post-index-change backup hook:", err)
		} else {
			fmt.Fprintln(os.Stderr, "post-index-change hook removed, backup restored")
			return
		}
	}

	fmt.Fprintln(os.Stderr, "post-index-change hook removed successfully")
}

func containsCommitgenSignature(content string) bool {
	return filepath.Base(content) != content &&
		(filepath.Dir(content) != "." || len(content) > 50) &&
		len(content) > 10 &&
		stringContains(content, "commitgen suggest")
}

func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
