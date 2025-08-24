package diff

import (
	//"fmt"
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

func StagedChanges(filesLimitBytes int) (files []string, patch string, err error) {
	files, err = StagedFiles()
	if err != nil {
		return nil, "", err
	}

	if len(files) == 0 {
		return files, "", nil
	}

	stagedChangesBytes, err := exec.Command("git", "diff", "--cached", "--unified=0").Output()
	if err != nil {
		return nil, "", err
	}

	patch = string(stagedChangesBytes)
	if len(patch) > filesLimitBytes {
		patch = patch[:filesLimitBytes]
	}

	return files, patch, nil

}
