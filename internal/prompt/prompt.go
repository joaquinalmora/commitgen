package prompt

import (
	"fmt"
	"strings"
)

// Builds the message, start of with heuristic and then a prompt text for AI

func MakePrompt(files []string) string {
	n := 2

	if len(files) == 0 {
		return "No staged files"
	}

	if len(files) < n {
		n = len(files)
	}
	head := files[:n]
	rest := len(files) - n

	base := strings.Join(head, ", ")

	if rest > 0 {
		return "Update " + base + fmt.Sprintf(" (+%d more)", rest)
	}
	return "Update " + base

}
