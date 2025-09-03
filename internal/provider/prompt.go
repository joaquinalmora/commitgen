package provider

import (
	"fmt"
	"os"
	"strings"
)

// buildSharedPrompt creates the core prompt content used by both providers
func buildSharedPrompt(files []string, patch string) string {
	var prompt strings.Builder

	prompt.WriteString("Analyze these code changes and generate a professional commit message.\n\n")

	// Provide more context about the files
	prompt.WriteString("Files modified:\n")
	for i, file := range files {
		if i >= 5 {
			prompt.WriteString(fmt.Sprintf("... and %d more files\n", len(files)-5))
			break
		}
		prompt.WriteString("- " + file + "\n")
	}

	// Give guidance for multi-file commits
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

// GetBuiltinConventions returns the embedded conventions.md content
func GetBuiltinConventions() (string, error) {
	return loadConventions()
}

// LoadConventionsWithSource loads conventions and returns the source
func LoadConventionsWithSource() (content string, source string, err error) {
	// Try custom conventions first
	customPaths := []string{"./conventions.md", "./internal/provider/conventions.md"}

	for _, path := range customPaths {
		if contentBytes, err := os.ReadFile(path); err == nil {
			return string(contentBytes), path, nil
		}
	}

	// Fall back to built-in
	content, err = loadConventions()
	if err != nil {
		return "", "", err
	}

	return content, "built-in", nil
}
