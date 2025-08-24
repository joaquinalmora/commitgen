package main

//Entry point for CLI
import (
	"fmt"
	"os"
	"sort"

	"github.com/joaquinalmora/commitgen/internal/diff"
	"github.com/joaquinalmora/commitgen/internal/prompt"
)

type Command struct {
	Description string
	Run         func(args []string)
}

var commands = map[string]Command{
	"suggest": {
		Description: "Suggest a commit message based on staged changes",
		Run: func(args []string) {
			suggest()
		},
	},
	"install-hook": {
		Description: "Install a git commit hook to auto-suggest commit messages",
		Run: func(args []string) {
			installHook()
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
		fmt.Println("Unknown command: %s\n\n", name)
		printUsage(commands)
		return
	}

	cmd.Run(os.Args[2:])
}

func suggest() {
	files, patch, err := diff.StagedChanges(100 * 1024) //100KB limit

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(len(patch), "bytes of staged changes")
	fmt.Println(patch[:min(100, len(patch))])
	fmt.Println(prompt.MakePrompt(files, patch))
}
