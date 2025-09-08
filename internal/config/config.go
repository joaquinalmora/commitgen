package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AI struct {
		Enabled  bool   `yaml:"enabled"`
		Provider string `yaml:"provider"`
		APIKey   string `yaml:"api_key"`
		Model    string `yaml:"model"`
		BaseURL  string `yaml:"base_url"`
	} `yaml:"ai"`

	Performance struct {
		PatchBytes int    `yaml:"patch_bytes"`
		CacheTTL   string `yaml:"cache_ttl"`
		MaxFiles   int    `yaml:"max_files"`
	} `yaml:"performance"`

	Git struct {
		AutoInstallHook bool   `yaml:"auto_install_hook"`
		CommitTemplate  string `yaml:"commit_template"`
	} `yaml:"git"`

	Output struct {
		Verbose bool `yaml:"verbose"`
		Plain   bool `yaml:"plain"`
		Colors  bool `yaml:"colors"`
	} `yaml:"output"`

	Advanced struct {
		ConventionsFile string `yaml:"conventions_file"`
		FallbackEnabled bool   `yaml:"fallback_enabled"`
		Debug           bool   `yaml:"debug"`
	} `yaml:"advanced"`

	MaxFiles      int
	PatchBytes    int
	UseAIFallback bool
}

func Load() Config {
	loadEnvFiles()

	var cfg Config

	cfg = loadFromYAML(cfg)

	cfg.AI.Enabled = getEnvBool("COMMITGEN_AI", cfg.AI.Enabled)
	cfg.AI.Provider = "openai"

	// Always use OpenAI API key
	if apiKey := getEnv("OPENAI_API_KEY", ""); apiKey != "" {
		cfg.AI.APIKey = apiKey
	}

	cfg.AI.Model = getEnv("COMMITGEN_MODEL", cfg.AI.Model)
	if cfg.AI.Model == "" {
		cfg.AI.Model = "gpt-4o-mini"
	}

	cfg.AI.BaseURL = getEnv("COMMITGEN_BASE_URL", cfg.AI.BaseURL)
	if cfg.AI.BaseURL == "" {
		cfg.AI.BaseURL = "https://api.openai.com/v1"
	}

	cfg.Performance.MaxFiles = getEnvInt("COMMITGEN_MAX_FILES", cfg.Performance.MaxFiles)
	if cfg.Performance.MaxFiles == 0 {
		cfg.Performance.MaxFiles = 10
	}

	cfg.Performance.PatchBytes = getEnvInt("COMMITGEN_PATCH_BYTES", cfg.Performance.PatchBytes)
	if cfg.Performance.PatchBytes == 0 {
		cfg.Performance.PatchBytes = 100 * 1024
	}

	cfg.MaxFiles = cfg.Performance.MaxFiles
	cfg.PatchBytes = cfg.Performance.PatchBytes
	cfg.UseAIFallback = getEnvBool("COMMITGEN_AI_FALLBACK", true)

	return cfg
}

func loadFromYAML(cfg Config) Config {
	cfg.AI.Provider = "openai"
	cfg.AI.Model = "gpt-4o-mini"
	cfg.Performance.PatchBytes = 100 * 1024
	cfg.Performance.MaxFiles = 10
	cfg.Performance.CacheTTL = "24h"
	cfg.Output.Colors = true
	cfg.Advanced.FallbackEnabled = true

	configPaths := []string{
		"./commitgen.yaml",
		"./commitgen.yml",
		"~/.commitgen.yaml",
		"~/.commitgen.yml",
	}

	for _, path := range configPaths {
		if path[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				continue
			}
			path = filepath.Join(homeDir, path[1:])
		}

		if data, err := os.ReadFile(path); err == nil {
			var yamlCfg Config
			if err := yaml.Unmarshal(data, &yamlCfg); err == nil {
				if yamlCfg.AI.Provider != "" {
					cfg.AI.Provider = yamlCfg.AI.Provider
				}
				if yamlCfg.AI.Model != "" {
					cfg.AI.Model = yamlCfg.AI.Model
				}
				if yamlCfg.AI.APIKey != "" {
					cfg.AI.APIKey = yamlCfg.AI.APIKey
				}
				if yamlCfg.AI.BaseURL != "" {
					cfg.AI.BaseURL = yamlCfg.AI.BaseURL
				}
				cfg.AI.Enabled = yamlCfg.AI.Enabled

				if yamlCfg.Performance.PatchBytes > 0 {
					cfg.Performance.PatchBytes = yamlCfg.Performance.PatchBytes
				}
				if yamlCfg.Performance.MaxFiles > 0 {
					cfg.Performance.MaxFiles = yamlCfg.Performance.MaxFiles
				}
				if yamlCfg.Performance.CacheTTL != "" {
					cfg.Performance.CacheTTL = yamlCfg.Performance.CacheTTL
				}

				cfg.Git = yamlCfg.Git
				cfg.Output = yamlCfg.Output
				cfg.Advanced = yamlCfg.Advanced
				break
			}
		}
	}

	return cfg
}

func loadEnvFiles() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	envFiles := []string{
		filepath.Join(home, ".env"), // ~/.env (global)
		".env",                      // ./.env (project-local)
		".env.local",                // ./.env.local (local overrides)
	}

	for _, envFile := range envFiles {
		_ = godotenv.Load(envFile)
	}
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
