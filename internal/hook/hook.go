package hook

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallHook() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get current directory:", err)
		return
	}

	hookPath := filepath.Join(cwd, ".git", "hooks", "prepare-commit-msg")

	if _, err := os.Stat(hookPath); err == nil {
		fmt.Fprintln(os.Stderr, "Error: prepare-commit-msg hook already exists")
		fmt.Fprintln(os.Stderr, "Please remove it before installing this hook.")
		return
	} else if !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "failed to check if hook exists:", err)
		return
	}

	binPath := filepath.Join(cwd, "bin", "commitgen")

	script := fmt.Sprintf(`#!/bin/sh
		MSG_FILE="$1"
		SOURCE="$2"

		if [ -s "$MSG_FILE" ]; then
		exit 0
		fi

		case "$SOURCE" in
		merge|squash|rebase)
			exit 0
			;;
		esac

		SUGGEST="%s suggest"

		# Run and capture suggestion
		OUTPUT=$($SUGGEST)

		if [ -z "$OUTPUT" ] || [ "$OUTPUT" = "No staged files" ]; then
		exit 0
		fi

		printf '%%s\n' "$OUTPUT" > "$MSG_FILE"
		`, binPath)

	if err := os.WriteFile(hookPath, []byte(script), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write hook file:", err)
		return
	}

	if err := os.Chmod(hookPath, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to make hook executable:", err)
		return
	}
	fmt.Fprintln(os.Stderr, "prepare-commit-msg hook installed successfully at", hookPath)
}

func UninstallHook() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get current directory:", err)
		return
	}

	hookPath := filepath.Join(cwd, ".git", "hooks", "prepare-commit-msg")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "No prepare-commit-msg hook found")
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "failed to check hook:", err)
		return
	}

	content, err := os.ReadFile(hookPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read hook file:", err)
		return
	}

	hookContent := string(content)
	if !containsCommitgenSignature(hookContent) {
		fmt.Fprintln(os.Stderr, "Hook exists but doesn't appear to be created by commitgen")
		fmt.Fprintln(os.Stderr, "Please remove manually if needed")
		return
	}

	if err := os.Remove(hookPath); err != nil {
		fmt.Fprintln(os.Stderr, "failed to remove hook:", err)
		return
	}

	fmt.Fprintln(os.Stderr, "prepare-commit-msg hook removed successfully")
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
