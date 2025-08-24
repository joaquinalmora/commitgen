package hook

// Installs git hooks that call binary to prefill messages

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallHook() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get current directory:", err)
		return
	}

	hookPath := filepath.Join(cwd, ".git", "hooks", "prepare-commit-msg")

	if _, err := os.Stat(hookPath); err == nil {
		fmt.Println("Error: prepare-commit-msg hook already exists")
		fmt.Println("Please remove it before installing this hook.")
		return
	} else if !os.IsNotExist(err) {
		fmt.Println("failed to check if hook exists:", err)
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
		fmt.Println("failed to write hook file:", err)
		return
	}

	if err := os.Chmod(hookPath, 0o755); err != nil {
		fmt.Println("failed to make hook executable:", err)
		return
	}
	fmt.Println("prepare-commit-msg hook installed successfully at", hookPath)
}
