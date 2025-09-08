package errors

import (
	"testing"
)

func TestUserError(t *testing.T) {
	err := UserError{
		Message: "Test error",
		Help:    "Test help message",
		Code:    1,
	}

	expected := "Test error\n\nHelp: Test help message"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestUserErrorWithoutHelp(t *testing.T) {
	err := UserError{
		Message: "Test error",
		Code:    1,
	}

	expected := "Test error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestNoGitRepo(t *testing.T) {
	err := NoGitRepo()
	if err.Code != 1 {
		t.Errorf("Expected code 1, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestNoStagedChanges(t *testing.T) {
	err := NoStagedChanges()
	if err.Code != 1 {
		t.Errorf("Expected code 1, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestInvalidAPIKey(t *testing.T) {
	err := InvalidAPIKey("openai")
	if err.Code != 2 {
		t.Errorf("Expected code 2, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestAIProviderError(t *testing.T) {
	testErr := UserError{Message: "test error"}
	err := AIProviderError("openai", testErr)
	if err.Code != 3 {
		t.Errorf("Expected code 3, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestNetworkError(t *testing.T) {
	testErr := UserError{Message: "connection failed"}
	err := NetworkError(testErr)
	if err.Code != 4 {
		t.Errorf("Expected code 4, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestConfigError(t *testing.T) {
	err := ConfigError("api_key", "invalid")
	if err.Code != 5 {
		t.Errorf("Expected code 5, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}

func TestGitError(t *testing.T) {
	testErr := UserError{Message: "git failed"}
	err := GitError("diff", testErr)
	if err.Code != 6 {
		t.Errorf("Expected code 6, got %d", err.Code)
	}
	if err.Message == "" {
		t.Error("Expected non-empty message")
	}
	if err.Help == "" {
		t.Error("Expected non-empty help")
	}
}
