package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joaquinalmora/commitgen/internal/diff"
	"github.com/joaquinalmora/commitgen/internal/hook"
	"github.com/joaquinalmora/commitgen/internal/prompt"
	"github.com/joaquinalmora/commitgen/internal/shell"
)

type Command struct {
	Description string
	Run         func(args []string)
}

var commands = map[string]Command{
	"suggest": {
		Description: "Suggest a commit message based on staged changes",
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
	"init-shell": {
		Description: "Install shell snippet and guarded .zshrc block",
		Run: func(args []string) {
			if err := shell.InstallShell(); err != nil {
				fmt.Fprintln(os.Stderr, "install shell failed:", err)
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
	// ensure we're in a repo
	if !inGitRepo() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository (no .git directory found)")
		os.Exit(1)
	}
	plain := hasFlag(args, "--plain")
	verbose := hasFlag(args, "--verbose")

	files, patch, err := diff.StagedChanges(100 * 1024)
	if err != nil {
		// always report the error to stderr and exit non-zero
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if len(patch) == 0 {
		// hooks expect --plain to silently do nothing; CLI non-plain should fail
		if plain {
			return
		}
		fmt.Fprintln(os.Stderr, "No staged files.")
		os.Exit(1)
	}

	msg := prompt.MakePrompt(files, patch)

	if plain {
		s := strings.TrimSpace(msg)
		if s != "" {
			fmt.Println(s)
		}
		return
	}

	if verbose {
		// diagnostics to stderr
		fmt.Fprintln(os.Stderr, len(patch), "bytes of staged changes")
		fmt.Fprintln(os.Stderr, patch[:min(100, len(patch))])
		fmt.Fprintln(os.Stderr, msg)
		return
	}

	// default: human message to stdout
	fmt.Println(msg)
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
