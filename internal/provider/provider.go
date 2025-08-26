package provider

import (
	"context"
	"fmt"
)

type Provider interface {
	GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error)
	Name() string
	IsConfigured() bool
}

type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}

type ProviderError struct {
	Provider string
	Err      error
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider %s: %v", e.Provider, e.Err)
}

func GetProvider(config Config) (Provider, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIProvider(config)
	case "ollama":
		return NewOllamaProvider(config)
	default:
		return nil, fmt.Errorf("unknown provider: %s", config.Provider)
	}
}

func NewOllamaProvider(config Config) (Provider, error) {
	return nil, fmt.Errorf("Ollama provider not implemented yet")
}
