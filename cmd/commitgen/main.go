package main

//Entry point for CLI
import (
	"fmt"
	"os"

	"github.com/joaquinalmora/commitgen/internal/diff"
	"github.com/joaquinalmora/commitgen/internal/prompt"
)

func main() {
	if len(os.Args) > 1 {
		cliArg := os.Args[1]

		if cliArg == "suggest" {
			suggest()
		} else {
			fmt.Println("Unknown command: " + cliArg)
		}
	} else {
		fmt.Println("Usage: commitgen <command>\nAvailable commands:\n  suggest")
	}
}

func suggest() {
	files, patch, err := diff.StagedChanges(100 * 1024) //100KB limit

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(len(patch), "bytes of staged changes")
	fmt.Println(patch[:min(100, len(patch))])
	fmt.Println(prompt.MakePrompt(files))
}
