package main

//Entry point for CLI
import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		cliArg := os.Args[1]

		if cliArg == "suggest" {
			suggest()
		} else {
			fmt.Println("Unknown command")
		}
	} else {
		fmt.Println("Usage: commitgen <command>\nAvailable commands:\n  suggest")
	}
}

func suggest() {
	fmt.Println("suggest message")
}
