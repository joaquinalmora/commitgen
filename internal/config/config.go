package config

import (
	"os"
	"strconv"
)

type Config struct {
	AI struct {
		Enabled  bool
		Provider string
		APIKey   string
		Model    string
		BaseURL  string
	}
	MaxFiles      int
	PatchBytes    int
	UseAIFallback bool
}

func Load() Config {
	var cfg Config

	cfg.AI.Enabled = getEnvBool("COMMITGEN_AI", false)
	cfg.AI.Provider = getEnv("COMMITGEN_AI_PROVIDER", "openai")
	cfg.AI.APIKey = getEnv("COMMITGEN_AI_API_KEY", "")
	cfg.AI.Model = getEnv("COMMITGEN_AI_MODEL", "gpt-4o-mini")
	cfg.AI.BaseURL = getEnv("COMMITGEN_AI_BASE_URL", "")

	cfg.MaxFiles = getEnvInt("COMMITGEN_MAX_FILES", 10)
	cfg.PatchBytes = getEnvInt("COMMITGEN_PATCH_BYTES", 100*1024)
	cfg.UseAIFallback = getEnvBool("COMMITGEN_AI_FALLBACK", true)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
