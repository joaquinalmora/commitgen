package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OllamaProvider struct {
	model   string
	baseURL string
	client  *http.Client
}

type ollamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func NewOllamaProvider(config Config) (Provider, error) {
	if config.Model == "" {
		config.Model = "llama3.2:3b"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}

	return &OllamaProvider{
		model:   config.Model,
		baseURL: config.BaseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) IsConfigured() bool {
	return p.model != "" && p.baseURL != ""
}

func (p *OllamaProvider) GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error) {
	conventions, err := loadConventions()
	if err != nil {
		return "", fmt.Errorf("failed to load conventions: %w", err)
	}

	prompt := buildOllamaPrompt(files, patch, conventions)

	reqBody := ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"num_predict": 100,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("Ollama API error: %s", ollamaResp.Error)
	}

	message := strings.TrimSpace(ollamaResp.Response)
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
				message = strings.Join(result, " ")
			} else {
				message = firstLine[:69] + "..."
			}
		} else {
			message = firstLine
		}
	}

	return message, nil
}

func buildOllamaPrompt(files []string, patch string, conventions string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a professional software developer writing commit messages. ")
	prompt.WriteString("Analyze the following code changes and generate a single, concise commit message.\n\n")

	prompt.WriteString("Follow these commit message conventions:\n")
	prompt.WriteString(conventions)
	prompt.WriteString("\n\n")

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

	prompt.WriteString("\nGenerate ONLY the commit message. Do not include explanations, markdown formatting, or code blocks. ")
	prompt.WriteString("Respond with just the raw commit message text that follows conventional commit format.")

	return prompt.String()
}
