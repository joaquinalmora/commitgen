package prompt

import (
	"fmt"
	"path/filepath"
	"strings"
)

func MakePrompt(files []string, patch string) string {
	n := 2

	if len(files) == 0 {
		return "chore: add missing files"
	}

	if isTestsOnly(files) {
		return analyzeTestChanges(patch)
	}
	if isDocsOnly(files) {
		return analyzeDocChanges(patch)
	}

	if isConfigOnly(files) {
		return analyzeConfigChanges(patch)
	}

	if isRenameOnly(patch) {
		return "refactor: rename files for clarity"
	}

	commitType := analyzeCommitType(patch)

	if len(files) < n {
		n = len(files)
	}
	head := files[:n]
	rest := len(files) - n

	base := strings.Join(head, ", ")

	if rest > 0 {
		return fmt.Sprintf("%s: update %s and %d more files", commitType, base, rest)
	}
	return fmt.Sprintf("%s: update %s", commitType, base)
}

func analyzeCommitType(patch string) string {
	lowerPatch := strings.ToLower(patch)

	if strings.Contains(lowerPatch, "fix") ||
		strings.Contains(lowerPatch, "bug") ||
		strings.Contains(lowerPatch, "error") ||
		strings.Contains(lowerPatch, "issue") {
		return "fix"
	}

	if strings.Contains(lowerPatch, "performance") ||
		strings.Contains(lowerPatch, "optimize") ||
		strings.Contains(lowerPatch, "speed") ||
		strings.Contains(lowerPatch, "memory") {
		return "perf"
	}

	if strings.Contains(lowerPatch, "security") ||
		strings.Contains(lowerPatch, "auth") ||
		strings.Contains(lowerPatch, "permission") {
		return "security"
	}

	if strings.Contains(lowerPatch, "refactor") ||
		strings.Contains(lowerPatch, "cleanup") ||
		strings.Contains(lowerPatch, "restructure") {
		return "refactor"
	}

	if strings.Contains(lowerPatch, "style") ||
		strings.Contains(lowerPatch, "format") ||
		strings.Contains(lowerPatch, "lint") {
		return "style"
	}

	addedLines := strings.Count(patch, "\n+")
	removedLines := strings.Count(patch, "\n-")

	if addedLines > removedLines*2 {
		return "feat"
	}

	return "chore"
}

func analyzeTestChanges(patch string) string {
	lowerPatch := strings.ToLower(patch)

	if strings.Contains(lowerPatch, "fix") {
		return "test: fix failing tests"
	}

	addedLines := strings.Count(patch, "\n+")
	removedLines := strings.Count(patch, "\n-")

	if addedLines > removedLines {
		return "test: add test coverage"
	}

	return "test: update test cases"
}

func analyzeDocChanges(patch string) string {
	lowerPatch := strings.ToLower(patch)

	if strings.Contains(lowerPatch, "readme") {
		return "docs: update README"
	}

	if strings.Contains(lowerPatch, "api") {
		return "docs: update API documentation"
	}

	if strings.Contains(lowerPatch, "fix") || strings.Contains(lowerPatch, "typo") {
		return "docs: fix documentation errors"
	}

	return "docs: update documentation"
}

func analyzeConfigChanges(patch string) string {
	lowerPatch := strings.ToLower(patch)

	if strings.Contains(lowerPatch, "dependency") ||
		strings.Contains(lowerPatch, "package") ||
		strings.Contains(lowerPatch, "version") {
		return "chore: update dependencies"
	}

	if strings.Contains(lowerPatch, "ci") ||
		strings.Contains(lowerPatch, "workflow") ||
		strings.Contains(lowerPatch, "pipeline") {
		return "ci: update CI configuration"
	}

	if strings.Contains(lowerPatch, "build") ||
		strings.Contains(lowerPatch, "webpack") ||
		strings.Contains(lowerPatch, "gulp") {
		return "build: update build configuration"
	}

	return "chore: update configuration"
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

	if !(strings.Contains(lf, "rename from") && strings.Contains(lf, "rename to")) {
		return false
	}

	if strings.Contains(lf, "@@") {
		return false
	}

	return true
}
