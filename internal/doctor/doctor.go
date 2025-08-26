package doctor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Result holds checks and messages
type Result struct {
	OK      bool
	Message string
}

// Run performs a set of local environment checks and prints a compact report to stdout.
// It returns a non-nil error when a fatal check fails (useful for CI or scripts).
func Run() error {
	var out bytes.Buffer
	ok := true

	// 1) are we in a git repo?
	if inGitRepo() {
		fmt.Fprintln(&out, "✔ Git repo: ok")
	} else {
		fmt.Fprintln(&out, "✖ Git repo: .git not found (not inside a git repo)")
		ok = false
	}

	// 2) binary on PATH?
	if p, err := exec.LookPath("commitgen"); err == nil {
		fmt.Fprintf(&out, "✔ commitgen binary on PATH: %s\n", p)
	} else {
		fmt.Fprintln(&out, "✖ commitgen binary not on PATH (build with: go build -o bin/commitgen ./cmd/commitgen)")
		ok = false
	}

	// 3) prepare-commit-msg hook present?
	if inGitRepo() {
		root, _ := gitRoot()
		hookPath := filepath.Join(root, ".git", "hooks", "prepare-commit-msg")
		if _, err := os.Stat(hookPath); err == nil {
			fmt.Fprintf(&out, "✔ prepare-commit-msg hook: %s\n", hookPath)
		} else {
			fmt.Fprintln(&out, "✖ prepare-commit-msg hook: not found (ok if not installed)")
		}
	}

	// 4) zsh snippet installed?
	if home, err := os.UserHomeDir(); err == nil {
		snippet := filepath.Join(home, ".config", "commitgen.zsh")
		if _, err := os.Stat(snippet); err == nil {
			fmt.Fprintf(&out, "✔ zsh snippet installed: %s\n", snippet)
		} else {
			fmt.Fprintln(&out, "✖ zsh snippet not found at ~/.config/commitgen.zsh (ok if you haven't installed shell integration)")
		}

		// zshrc guarded block
		zshrc := filepath.Join(home, ".zshrc")
		if b, err := os.ReadFile(zshrc); err == nil {
			if bytes.Contains(b, []byte("# >>> commitgen >>> (managed)")) {
				fmt.Fprintln(&out, "✔ .zshrc contains commitgen guarded block")
			} else {
				fmt.Fprintln(&out, "✖ .zshrc does not contain commitgen guarded block")
			}
		}
	}

	// 5) staged changes?
	if inGitRepo() {
		if list, _ := gitStagedList(); len(list) > 0 {
			fmt.Fprintf(&out, "✔ staged files: %d\n", len(list))
		} else {
			fmt.Fprintln(&out, "✖ no staged files (commitgen suggest requires staged changes to produce suggestions)")
		}
	}

	// 6) zsh-autosuggestions presence (best-effort)
	if autosuggestAvailable() {
		fmt.Fprintln(&out, "✔ zsh-autosuggestions: detected")
	} else {
		fmt.Fprintln(&out, "✖ zsh-autosuggestions: not detected (plugin-first UX will not be active)")
	}

	// print summary
	fmt.Println(out.String())

	if ok {
		return nil
	}
	return fmt.Errorf("doctor detected issues; see output")
}

func inGitRepo() bool {
	if _, err := os.Stat(".git"); err == nil {
		return true
	}
	// fallback: git rev-parse
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	return err == nil
}

func gitRoot() (string, error) {
	b, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

func gitStagedList() ([]string, error) {
	b, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return nil, err
	}
	out := bytes.TrimSpace(b)
	if len(out) == 0 {
		return []string{}, nil
	}
	parts := bytes.Split(out, []byte("\n"))
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		res = append(res, string(p))
	}
	return res, nil
}

func autosuggestAvailable() bool {
	// check common env var
	if _, ok := os.LookupEnv("ZSH_AUTOSUGGEST_DIR"); ok {
		return true
	}
	// check common oh-my-zsh path
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	p := filepath.Join(home, ".oh-my-zsh", "custom", "plugins", "zsh-autosuggestions", "zsh-autosuggestions.zsh")
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
