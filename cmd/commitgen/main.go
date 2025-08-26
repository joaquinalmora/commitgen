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
	"github.com/joaquinalmora/commitgen/internal/hook"
	"github.com/joaquinalmora/commitgen/internal/prompt"
	"github.com/joaquinalmora/commitgen/internal/provider"
	"github.com/joaquinalmora/commitgen/internal/shell"
)

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
		fmt.Fprintln(os.Stderr, "Error: not a git repository (no .git directory found)")
		os.Exit(1)
	}

	plain := hasFlag(args, "--plain")
	verbose := hasFlag(args, "--verbose")
	useAI := hasFlag(args, "--ai")
	useCache := hasFlag(args, "--cached")

	cfg := config.Load()
	if cfg.AI.Enabled {
		useAI = true
	}

	files, patch, err := diff.StagedChanges(cfg.PatchBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if len(patch) == 0 {
		if plain {
			return
		}
		fmt.Fprintln(os.Stderr, "No staged files.")
		os.Exit(1)
	}

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
		if verbose {
			fmt.Fprintln(os.Stderr, "Using cached message for these changes")
		}
		if plain {
			fmt.Println(cached.Message)
		} else {
			fmt.Println(cached.Message)
		}
		return
	}

	var msg string

	if useAI && cfg.AI.APIKey != "" {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using AI provider:", cfg.AI.Provider)
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
				fmt.Fprintln(os.Stderr, "Falling back to heuristics")
			}
			msg = prompt.MakePrompt(files, patch)
		} else {
			ctx := context.Background()
			aiMsg, err := aiProvider.GenerateCommitMessage(ctx, files, patch)
			if err != nil {
				if verbose {
					fmt.Fprintln(os.Stderr, "AI generation error:", err)
					fmt.Fprintln(os.Stderr, "Falling back to heuristics")
				}
				msg = prompt.MakePrompt(files, patch)
			} else {
				msg = aiMsg
				if verbose {
					fmt.Fprintln(os.Stderr, "Generated using AI")
				}
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
