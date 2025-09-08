package errors

import (
	"fmt"
	"strings"
)

// UserError represents an error that should be displayed to the user
type UserError struct {
	Message string
	Help    string
	Code    int
}

func (e UserError) Error() string {
	if e.Help != "" {
		return fmt.Sprintf("%s\n\nHelp: %s", e.Message, e.Help)
	}
	return e.Message
}

func NoGitRepo() UserError {
	return UserError{
		Message: "Not a git repository",
		Help:    "Navigate to a git repository directory or run 'git init' to create one",
		Code:    1,
	}
}

func NoStagedChanges() UserError {
	return UserError{
		Message: "No staged changes found",
		Help:    "Stage your changes with 'git add <files>' before generating a commit message",
		Code:    1,
	}
}

func InvalidAPIKey(provider string) UserError {
	var helpMsg string
	switch strings.ToLower(provider) {
	case "openai":
		helpMsg = "Get your API key from https://platform.openai.com/api-keys and set it with:\n" +
			"  export OPENAI_API_KEY=your-key-here\n" +
			"  Or add it to your ~/.env file"
	default:
		helpMsg = fmt.Sprintf("Check your %s API key configuration", provider)
	}

	return UserError{
		Message: fmt.Sprintf("Invalid or missing %s API key", provider),
		Help:    helpMsg,
		Code:    2,
	}
}

func AIProviderError(provider string, err error) UserError {
	return UserError{
		Message: fmt.Sprintf("AI provider (%s) failed: %v", provider, err),
		Help: "Check your internet connection and API key. " +
			"Commitgen will fall back to heuristic message generation.",
		Code: 3,
	}
}

func NetworkError(err error) UserError {
	return UserError{
		Message: fmt.Sprintf("Network error: %v", err),
		Help: "Check your internet connection and try again. " +
			"You can use '--cached' to use a previously generated message or " +
			"use commitgen without AI for heuristic suggestions.",
		Code: 4,
	}
}

func ConfigError(field string, value string) UserError {
	return UserError{
		Message: fmt.Sprintf("Invalid configuration: %s = %s", field, value),
		Help:    "Check your configuration file or environment variables",
		Code:    5,
	}
}

func GitError(operation string, err error) UserError {
	return UserError{
		Message: fmt.Sprintf("Git operation failed (%s): %v", operation, err),
		Help:    "Ensure you're in a valid git repository and have the necessary permissions",
		Code:    6,
	}
}
