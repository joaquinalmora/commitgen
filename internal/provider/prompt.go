package provider

import (
	"os"
)

func GetBuiltinConventions() (string, error) {
	return loadConventions()
}

func LoadConventionsWithSource() (content string, source string, err error) {
	customPaths := []string{"./conventions.md", "./internal/provider/conventions.md"}

	for _, path := range customPaths {
		if contentBytes, err := os.ReadFile(path); err == nil {
			return string(contentBytes), path, nil
		}
	}

	content, err = loadConventions()
	if err != nil {
		return "", "", err
	}

	return content, "built-in", nil
}
