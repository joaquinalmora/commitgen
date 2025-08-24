package diff

import (
	"os/exec"
	"strings"
)

// Reads staged file list and patch from git.
// Returns filenames and a trimmed diff string

func StagedFiles() ([]string, error) {
	changedFilesBytes, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return nil, err
	}

	rawFiles := strings.Split(strings.TrimSpace(string(changedFilesBytes)), "\n")

	var files []string
	for _, f := range rawFiles {
		if f != "" {
			files = append(files, f)
		}
	}

	return files, nil
}
