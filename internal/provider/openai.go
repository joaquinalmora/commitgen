package provider

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joaquinalmora/commitgen/internal/errors"
)

//go:embed conventions.md
var conventionsFS embed.FS

type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

type openAIRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []choice     `json:"choices"`
	Error   *openAIError `json:"error,omitempty"`
}

type choice struct {
	Message message `json:"message"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewOpenAIProvider(config Config) (Provider, error) {
	if config.APIKey == "" {
		return nil, errors.InvalidAPIKey("OpenAI")
	}

	if !strings.HasPrefix(config.APIKey, "sk-") {
		return nil, errors.InvalidAPIKey("OpenAI")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := config.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIProvider{
		apiKey:  config.APIKey,
		model:   model,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) IsConfigured() bool {
	return p.apiKey != ""
}

func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error) {
	prompt := buildPrompt(files, patch)

	conventions, err := loadConventions()
	if err != nil {
		conventions = "Use conventional commit format: type: description (under 50 chars)"
	}

	reqBody := openAIRequest{
		Model: p.model,
		Messages: []message{
			{
				Role:    "system",
				Content: conventions,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   100,
		Temperature: 0.1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", errors.NetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return "", errors.InvalidAPIKey("OpenAI")
		case http.StatusTooManyRequests:
			return "", errors.UserError{
				Message: "OpenAI API rate limit exceeded",
				Help:    "Wait a moment and try again, or upgrade your OpenAI plan",
				Code:    7,
			}
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			return "", errors.UserError{
				Message: "OpenAI service temporarily unavailable",
				Help:    "Try again in a few moments or use '--cached' for a previous message",
				Code:    8,
			}
		default:
			return "", fmt.Errorf("OpenAI API error (HTTP %d): %s", resp.StatusCode, string(body))
		}
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", errors.AIProviderError("OpenAI", fmt.Errorf("%s", openAIResp.Error.Message))
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	message := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	message = strings.Trim(message, `"'`)

	if strings.HasPrefix(message, "```") {
		lines := strings.Split(message, "\n")
		if len(lines) > 1 {
			message = strings.Join(lines[1:], "\n")
		}
		message = strings.TrimSuffix(message, "```")
		message = strings.TrimSpace(message)
	}

	if len(message) > 72 {
		lines := strings.Split(message, "\n")
		firstLine := lines[0]

		if len(firstLine) > 72 {
			words := strings.Fields(firstLine)
			var result []string
			length := 0

			for _, word := range words {
				if length+len(word)+1 > 72 {
					break
				}
				result = append(result, word)
				length += len(word) + 1
			}

			if len(result) > 0 {
				truncated := strings.Join(result, " ")
				// Check if the last word is a connector word that suggests incomplete thought
				lastWord := strings.ToLower(result[len(result)-1])
				if lastWord == "and" || lastWord == "or" || lastWord == "but" || lastWord == "with" || lastWord == "for" || lastWord == "to" {
					// Remove the connector word to avoid incomplete sentences
					if len(result) > 1 {
						truncated = strings.Join(result[:len(result)-1], " ")
					} else {
						// If only connector word, fall back to character truncation
						truncated = firstLine[:69] + "..."
					}
				}
				message = truncated
			} else {
				message = firstLine[:69] + "..."
			}
		} else {
			message = firstLine
		}
	}

	message = strings.TrimSpace(message)
	message = trimTrailingConnector(message)
	if message == "" {
		message = "chore: update files"
	}

	return message, nil
}

func buildPrompt(files []string, patch string) string {
	var prompt strings.Builder

	prompt.WriteString("Analyze these code changes and generate a professional commit message:\n\n")

	prompt.WriteString("Files modified:\n")
	for i, file := range files {
		if i >= 5 {
			prompt.WriteString(fmt.Sprintf("... and %d more files\n", len(files)-5))
			break
		}
		prompt.WriteString("- " + file + "\n")
	}

	prompt.WriteString("\nCode changes (git diff):\n")
	if len(patch) > 2000 {
		prompt.WriteString(patch[:2000] + "...\n")
	} else {
		prompt.WriteString(patch + "\n")
	}

	prompt.WriteString("\nGenerate only the commit message text. Do not include any markdown formatting, code blocks, or explanations. Return only the raw commit message.")

	return prompt.String()
}

func loadConventions() (string, error) {
	customPath := os.Getenv("COMMITGEN_CONVENTIONS_FILE")
	if customPath != "" {
		content, err := os.ReadFile(customPath)
		if err != nil {
			return "", fmt.Errorf("failed to read custom conventions file: %w", err)
		}
		return string(content), nil
	}

	content, err := conventionsFS.ReadFile("conventions.md")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func trimTrailingConnector(message string) string {
	if message == "" {
		return message
	}

	connectors := []string{" and", " or", " but", " with", " for", " to"}
	lower := strings.ToLower(message)

	for _, connector := range connectors {
		if strings.HasSuffix(lower, connector) {
			message = strings.TrimSpace(message[:len(message)-len(connector)])
			break
		}
	}

	return strings.TrimSpace(message)
}
