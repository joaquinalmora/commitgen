package prompt

// Builds the message, start of with heuristic and then a prompt text for AI

func MakePrompt(files []string) string {
	prompt := ""
	if len(files) == 0 {
		prompt = "No staged changes"
	} else if len(files) == 1 {
		prompt = "Update: " + files[0]
	} else {
		prompt = "Update multiple files"
	}
	return prompt
}
