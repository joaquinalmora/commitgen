package main

//Entry point for CLI
import (
	"fmt"
	"os"

	"github.com/joaquinalmora/commitgen/internal/diff"
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
	fmt.Println("suggest message")
	files, err := diff.StagedFiles()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(files)
}
