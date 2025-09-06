package provider

import (
	"fmt"
	"os"
	"strings"
)

func buildSharedPrompt(files []string, patch string) string {
	var prompt strings.Builder

	prompt.WriteString("Analyze these code changes and generate a professional commit message.\n\n")

	prompt.WriteString("Files modified:\n")
	for i, file := range files {
		if i >= 5 {
			prompt.WriteString(fmt.Sprintf("... and %d more files\n", len(files)-5))
			break
		}
		prompt.WriteString("- " + file + "\n")
	}

	if len(files) > 3 {
		prompt.WriteString("\nThis is a multi-file change. Focus on the main purpose and scope.\n")
		prompt.WriteString("Look for the common theme across all changes.\n")
		prompt.WriteString("If changes are mixed (e.g., docs + code + tests), prioritize the most significant functional change.\n\n")
	}

	prompt.WriteString("Code changes (git diff):\n")
	if len(patch) > 2000 {
		prompt.WriteString(patch[:2000] + "...\n")
	} else {
		prompt.WriteString(patch + "\n")
	}

	prompt.WriteString("\nInstructions:\n")
	prompt.WriteString("- Use conventional commit format: type(scope): description\n")
	prompt.WriteString("- Keep description under 50 characters\n")
	prompt.WriteString("- Be specific about what changed, not just file types\n")
	prompt.WriteString("- For multi-file changes, describe the main functional change\n")
	prompt.WriteString("- Common types: feat, fix, docs, style, refactor, test, chore\n")
	prompt.WriteString("- If primarily adding new functionality, use 'feat'\n")
	prompt.WriteString("- If primarily fixing issues, use 'fix'\n")
	prompt.WriteString("- If mixed changes, choose based on the most significant change\n\n")

	prompt.WriteString("Generate only the commit message text. No explanations or formatting.")

	return prompt.String()
}

func GetBuiltinConventions() (string, error) {
	return loadConventions()
}

func LoadConventionsWithSource() (content string, source string, err error) {
	customPaths := []string{"./conventions.md", "./internal/provider/conventions.md"}

	for _, path := range customPaths {
		if contentBytes, err := os.ReadFile(path); err == nil {
			return string(contentBytes), path, nil
		}
	}

	content, err = loadConventions()
	if err != nil {
		return "", "", err
	}

	return content, "built-in", nil
}
