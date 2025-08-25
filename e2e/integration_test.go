package main_test

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Builds the commitgen binary into the temp dir, creates a small git repo,
// stages a file and asserts `commitgen suggest --plain` prints a single line.
func TestSuggestPlainIntegration(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmp := t.TempDir()
	// find repo root by walking up to where go.mod lives
	repoRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(repoRoot)
		if parent == repoRoot {
			t.Fatalf("could not find repo root from %s", wd)
		}
		repoRoot = parent
	}

	// build binary into tmp
	binPath := filepath.Join(tmp, "commitgen")
	build := exec.Command("go", "build", "-o", binPath, "./cmd/commitgen")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}

	// init git repo
	repo := filepath.Join(tmp, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v\n%s", err, string(out))
	}
	// set local git config so commits would work
	cfg1 := exec.Command("git", "config", "user.email", "e2e@example.com")
	cfg1.Dir = repo
	cfg1.Run()
	cfg2 := exec.Command("git", "config", "user.name", "e2e")
	cfg2.Dir = repo
	cfg2.Run()

	// create and stage a file
	fpath := filepath.Join(repo, "demo.txt")
	if err := os.WriteFile(fpath, []byte("hello e2e"), 0o644); err != nil {
		t.Fatal(err)
	}
	add := exec.Command("git", "add", "demo.txt")
	add.Dir = repo
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v\n%s", err, string(out))
	}

	// run commitgen suggest --plain (use built binary)
	suggest := exec.Command(binPath, "suggest", "--plain")
	suggest.Dir = repo
	out, err := suggest.Output()
	if err != nil {
		// capture stderr if available
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("suggest failed: %v\n%s", err, string(ee.Stderr))
		}
		t.Fatalf("suggest failed: %v", err)
	}

	// check output is a single non-empty line
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	lines := 0
	var last string
	for scanner.Scan() {
		last = scanner.Text()
		lines++
	}
	if lines != 1 {
		t.Fatalf("expected 1 line, got %d, output=%q", lines, string(out))
	}
	if strings.TrimSpace(last) == "" {
		t.Fatalf("empty suggestion")
	}
}
