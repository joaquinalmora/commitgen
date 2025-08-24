package prompt

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Builds the message, start of with heuristic and then a prompt text for AI

func MakePrompt(files []string, patch string) string {
	n := 2

	if len(files) == 0 {
		return "No staged files"
	}

	if isTestsOnly(files) {
		return "Update tests"
	}
	if isDocsOnly(files) {
		return "Update documentation"
	}

	if isConfigOnly(files) {
		return "Update configuration"
	}

	if isRenameOnly(patch) {
		return "Rename files"
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

func matchesAnySuffix(file string, suffixes []string) bool {
	file = strings.ToLower(file)
	for _, s := range suffixes {
		if strings.HasSuffix(file, s) {
			return true
		}
	}
	return false
}

func matchesAnyPrefix(file string, prefixes []string) bool {
	file = strings.ToLower(file)
	for _, p := range prefixes {
		if strings.HasPrefix(file, p) {
			return true
		}
	}
	return false
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func isTestsOnly(files []string) bool {
	testSuffixes := []string{".test.js", ".spec.js", ".test.ts", ".spec.ts", "_test.go", "_mock.go", "_test.py", "test_", "_test"}
	testPrefixes := []string{"test/", "tests/", "src/test/"}

	for _, f := range files {
		if !(matchesAnySuffix(f, testSuffixes) || matchesAnyPrefix(f, testPrefixes)) {
			return false
		}
	}
	return true
}

func isDocsOnly(files []string) bool {
	docSuffixes := []string{".md", ".rst", ".adoc"}
	docPrefixes := []string{"docs/", "doc/"}
	docSpecials := []string{"readme", "changelog", "contributing"}

	for _, f := range files {
		lf := strings.ToLower(f)
		base := strings.TrimSuffix(filepath.Base(lf), filepath.Ext(lf))
		if !(matchesAnySuffix(lf, docSuffixes) || matchesAnyPrefix(lf, docPrefixes) || contains(docSpecials, base)) {
			return false
		}
	}
	return true
}

func isConfigOnly(files []string) bool {
	configSuffixes := []string{".json", ".yaml", ".yml", ".toml", ".ini", ".cfg", ".config", ".xml", ".env", ".properties", ".conf"}
	configSpecials := []string{".gitignore", ".gitattributes", "dockerfile", "makefile", ".editorconfig", "go.mod", "go.sum", "package.json", "package-lock", "pnpm-lock", "requirements", "pipfile", "pyproject"}
	configPrefixes := []string{"config/", "configs/", "settings/", ".github/", ".vscode/"}

	for _, f := range files {
		lf := strings.ToLower(f)
		base := strings.TrimSuffix(filepath.Base(lf), filepath.Ext(lf))
		if !(matchesAnySuffix(lf, configSuffixes) || contains(configSpecials, base) || matchesAnyPrefix(lf, configPrefixes) || strings.HasPrefix(base, ".env")) {
			return false
		}
	}
	return true
}

func isRenameOnly(patch string) bool {
	lf := strings.ToLower(patch)

	if !(strings.Contains(lf, "rename from)") && strings.Contains(lf, "rename to")) {
		return false
	}

	if strings.Contains(lf, "@@") {
		return false
	}

	return true
}
