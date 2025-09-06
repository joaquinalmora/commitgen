package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joaquinalmora/commitgen/internal/cache"
	"github.com/joaquinalmora/commitgen/internal/config"
	"github.com/joaquinalmora/commitgen/internal/diff"
	"github.com/joaquinalmora/commitgen/internal/doctor"
	"github.com/joaquinalmora/commitgen/internal/errors"
	"github.com/joaquinalmora/commitgen/internal/hook"
	"github.com/joaquinalmora/commitgen/internal/logger"
	"github.com/joaquinalmora/commitgen/internal/prompt"
	"github.com/joaquinalmora/commitgen/internal/provider"
	"github.com/joaquinalmora/commitgen/internal/shell"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// handleError prints user-friendly error messages and exits with appropriate code
func handleError(err error) {
	if userErr, ok := err.(errors.UserError); ok {
		fmt.Fprintln(os.Stderr, "Error:", userErr.Message)
		if userErr.Help != "" {
			fmt.Fprintln(os.Stderr, "\nHelp:", userErr.Help)
		}
		os.Exit(userErr.Code)
	} else {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type Command struct {
	Description string
	Run         func(args []string)
}

var commands = map[string]Command{
	"suggest": {
		Description: "Suggest a commit message based on staged changes [--ai] [--plain] [--verbose]",
		Run: func(args []string) {
			suggest(args)
		},
	},
	"install-hook": {
		Description: "Install a git commit hook to auto-suggest commit messages",
		Run: func(args []string) {
			hook.InstallHook()
		},
	},
	"uninstall-hook": {
		Description: "Remove the git commit hook installed by commitgen",
		Run: func(args []string) {
			hook.UninstallHook()
		},
	},
	"install-shell": {
		Description: "(alias) Install shell snippet and guarded .zshrc block",
		Run: func(args []string) {
			if err := shell.InstallShell(); err != nil {
				fmt.Fprintln(os.Stderr, "install shell failed:", err)
			}
		},
	},
	"cache": {
		Description: "Generate and cache commit message for current staged changes",
		Run: func(args []string) {
			generateCache(args)
		},
	},
	"cached": {
		Description: "Get the most recent cached commit message",
		Run: func(args []string) {
			getCached(args)
		},
	},
	"doctor": {
		Description: "Run environment checks and print a diagnostic report",
		Run: func(args []string) {
			if err := doctor.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "doctor checks failed:", err)
				os.Exit(1)
			}
		},
	},
	"version": {
		Description: "Show version information",
		Run: func(args []string) {
			fmt.Printf("commitgen %s\n", version)
			if len(args) > 0 && args[0] == "--verbose" {
				fmt.Printf("Commit: %s\n", commit)
				fmt.Printf("Built: %s\n", date)
			}
		},
	},
	"init": {
		Description: "Create a configuration file interactively",
		Run: func(args []string) {
			initConfig(args)
		},
	},
	"uninstall-shell": {
		Description: "Remove shell snippet and guarded .zshrc block",
		Run: func(args []string) {
			if err := shell.UninstallShell(); err != nil {
				fmt.Fprintln(os.Stderr, "uninstall shell failed:", err)
			}
		},
	},
}

func printUsage(commands map[string]Command) {
	fmt.Println("Usage: commitgen <command> [options]")
	fmt.Println("Available commands:")

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %-13s %s\n", k, commands[k].Description)
	}

}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage(commands)
		return
	}

	name := os.Args[1]

	cmd, ok := commands[name]

	if !ok {
		fmt.Printf("Unknown command: %s\n\n", name)
		printUsage(commands)
		return
	}

	cmd.Run(os.Args[2:])
}

func inGitRepo() bool {
	cwd, _ := os.Getwd()
	_, err := os.Stat(filepath.Join(cwd, ".git"))
	return err == nil
}

func suggest(args []string) {
	if !inGitRepo() {
		handleError(errors.NoGitRepo())
	}

	plain := hasFlag(args, "--plain")
	verbose := hasFlag(args, "--verbose")
	useAI := hasFlag(args, "--ai")
	useCache := hasFlag(args, "--cached")

	logger.SetVerbose(verbose)

	cfg := config.Load()
	if cfg.AI.Enabled {
		useAI = true
	}

	logger.Debug("Configuration loaded: AI enabled=%v, provider=%s", cfg.AI.Enabled, cfg.AI.Provider)

	files, patch, err := diff.StagedChanges(cfg.PatchBytes)
	if err != nil {
		handleError(errors.GitError("reading staged changes", err))
	}

	if len(patch) == 0 {
		if plain {
			return
		}
		handleError(errors.NoStagedChanges())
	}

	logger.Debug("Found %d changed files, patch size: %d bytes", len(files), len(patch))

	c := cache.New()

	if useCache {
		cached, err := c.GetLatest()
		if err == nil {
			if verbose {
				fmt.Fprintln(os.Stderr, "Using cached message from", cached.Timestamp.Format("15:04:05"))
			}
			if plain {
				fmt.Println(cached.Message)
			} else {
				fmt.Println(cached.Message)
			}
			return
		}
		if verbose {
			fmt.Fprintln(os.Stderr, "No valid cache found, generating new message")
		}
	}

	cached, err := c.Get(files, patch)
	if err == nil && !useCache {
		logger.Debug("Using cached message for these changes")
		if plain {
			fmt.Println(cached.Message)
		} else {
			fmt.Println(cached.Message)
		}
		return
	}

	var msg string

	if useAI && cfg.AI.APIKey != "" {
		logger.Info("Using AI provider: %s", cfg.AI.Provider)

		providerConfig := provider.Config{
			Provider: cfg.AI.Provider,
			APIKey:   cfg.AI.APIKey,
			Model:    cfg.AI.Model,
			BaseURL:  cfg.AI.BaseURL,
		}

		aiProvider, err := provider.GetProvider(providerConfig)
		if err != nil {
			logger.Warn("AI provider initialization failed: %v", err)
			logger.Info("Falling back to heuristic message generation")
			msg = prompt.MakePrompt(files, patch)
		} else {
			ctx := context.Background()
			logger.Debug("Sending request to AI provider...")
			aiMsg, err := aiProvider.GenerateCommitMessage(ctx, files, patch)
			if err != nil {
				logger.Warn("AI generation failed: %v", err)
				logger.Info("Falling back to heuristic message generation")
				msg = prompt.MakePrompt(files, patch)
			} else {
				msg = aiMsg
				logger.Debug("Successfully generated commit message using AI")
				c.Set(files, patch, msg, cfg.AI.Provider)
			}
		}
	} else {
		if useAI && verbose {
			fmt.Fprintln(os.Stderr, "AI requested but no API key configured, using heuristics")
		}
		msg = prompt.MakePrompt(files, patch)
		if useAI {
			c.Set(files, patch, msg, "heuristics")
		}
	}

	if plain {
		s := strings.TrimSpace(msg)
		if s != "" {
			fmt.Println(s)
		}
		return
	}

	if verbose {
		fmt.Fprintln(os.Stderr, len(patch), "bytes of staged changes")
		fmt.Fprintln(os.Stderr, patch[:min(100, len(patch))])
		fmt.Fprintln(os.Stderr, msg)
		return
	}

	fmt.Println(msg)
}

func generateCache(args []string) {
	if !inGitRepo() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository (no .git directory found)")
		os.Exit(1)
	}

	verbose := hasFlag(args, "--verbose")
	cfg := config.Load()

	files, patch, err := diff.StagedChanges(cfg.PatchBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if len(patch) == 0 {
		if verbose {
			fmt.Fprintln(os.Stderr, "No staged files to cache.")
		}
		return
	}

	c := cache.New()

	var msg string
	var providerName string

	if cfg.AI.APIKey != "" {
		if verbose {
			fmt.Fprintln(os.Stderr, "Generating AI cache for", len(files), "files")
		}

		providerConfig := provider.Config{
			Provider: cfg.AI.Provider,
			APIKey:   cfg.AI.APIKey,
			Model:    cfg.AI.Model,
			BaseURL:  cfg.AI.BaseURL,
		}

		aiProvider, err := provider.GetProvider(providerConfig)
		if err != nil {
			if verbose {
				fmt.Fprintln(os.Stderr, "AI provider error:", err)
			}
			msg = prompt.MakePrompt(files, patch)
			providerName = "heuristics"
		} else {
			ctx := context.Background()
			aiMsg, err := aiProvider.GenerateCommitMessage(ctx, files, patch)
			if err != nil {
				if verbose {
					fmt.Fprintln(os.Stderr, "AI generation error:", err)
				}
				msg = prompt.MakePrompt(files, patch)
				providerName = "heuristics"
			} else {
				msg = aiMsg
				providerName = cfg.AI.Provider
			}
		}
	} else {
		msg = prompt.MakePrompt(files, patch)
		providerName = "heuristics"
	}

	err = c.Set(files, patch, msg, providerName)
	if err != nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Cache save error:", err)
		}
		return
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "Cached message using", providerName+":", msg)
	}
}

func getCached(args []string) {
	plain := hasFlag(args, "--plain")
	verbose := hasFlag(args, "--verbose")

	c := cache.New()
	cached, err := c.GetLatest()
	if err != nil {
		if !plain {
			fmt.Fprintln(os.Stderr, "No cached messages found")
		}
		os.Exit(1)
	}

	if verbose && !plain {
		fmt.Fprintln(os.Stderr, "Cached at:", cached.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Fprintln(os.Stderr, "Provider:", cached.Provider)
		fmt.Fprintln(os.Stderr, "Files:", strings.Join(cached.Files, ", "))
	}

	fmt.Println(cached.Message)
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func initConfig(args []string) {
	fmt.Println("🔧 CommitGen Configuration Setup")
	fmt.Println("This will create a commitgen.yaml configuration file.")
	fmt.Println()

	global := hasFlag(args, "--global")
	if !global {
		fmt.Print("Create config file globally (~/.commitgen.yaml) or locally (./commitgen.yaml)? [global/local] (default: local): ")
		var choice string
		fmt.Scanln(&choice)
		global = choice == "global" || choice == "g"
	}

	configPath := "./commitgen.yaml"
	if global {
		if homeDir, err := os.UserHomeDir(); err == nil {
			configPath = filepath.Join(homeDir, ".commitgen.yaml")
		}
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠️  Configuration file already exists at %s\n", configPath)
		fmt.Print("Overwrite? [y/N]: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" && confirm != "yes" {
			fmt.Println("Configuration setup cancelled.")
			return
		}
	}

	fmt.Print("Enter your OpenAI API key (or press Enter to configure later): ")
	var apiKey string
	fmt.Scanln(&apiKey)

	fmt.Print("Choose AI model [gpt-4o/gpt-4o-mini/gpt-3.5-turbo] (default: gpt-4o-mini): ")
	var model string
	fmt.Scanln(&model)
	if model == "" {
		model = "gpt-4o-mini"
	}

	fmt.Print("Enable AI by default? [y/N]: ")
	var aiEnabledStr string
	fmt.Scanln(&aiEnabledStr)
	aiEnabled := aiEnabledStr == "y" || aiEnabledStr == "Y" || aiEnabledStr == "yes"

	configContent := fmt.Sprintf(`# CommitGen Configuration
# Generated by 'commitgen init'

ai:
  enabled: %t
  provider: "openai"
  model: "%s"
  api_key: "%s"
  base_url: ""

performance:
  patch_bytes: 4000
  cache_ttl: "24h"
  max_files: 10

git:
  auto_install_hook: false
  commit_template: ""

output:
  verbose: false
  plain: false
  colors: true

advanced:
  conventions_file: ""
  fallback_enabled: true
  debug: false
`, aiEnabled, model, apiKey)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("❌ Failed to create config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuration file created at %s\n", configPath)
	fmt.Println()
	fmt.Println("Next steps:")
	if apiKey == "" {
		fmt.Println("1. Add your OpenAI API key to the config file or set OPENAI_API_KEY environment variable")
	}
	fmt.Println("2. Customize the configuration as needed")
	fmt.Println("3. Run 'commitgen suggest' to test your setup")
}
