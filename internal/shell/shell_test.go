package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallUninstallShell(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmp); err != nil {
		t.Fatalf("setenv HOME: %v", err)
	}
	defer os.Setenv("HOME", oldHome)

	cfgPath := filepath.Join(tmp, ".config", "commitgen.zsh")
	zshrcPath := filepath.Join(tmp, ".zshrc")

	if err := InstallShell(); err != nil {
		t.Fatalf("InstallShell failed: %v", err)
	}
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatalf("expected snippet at %s", cfgPath)
	}
	zb, err := os.ReadFile(zshrcPath)
	if err != nil {
		t.Fatalf("reading .zshrc: %v", err)
	}
	if !strings.Contains(string(zb), guardStart) || !strings.Contains(string(zb), guardEnd) {
		t.Fatalf("guarded block not found in .zshrc")
	}

	if err := InstallShell(); err != nil {
		t.Fatalf("second InstallShell failed: %v", err)
	}

	if err := UninstallShell(); err != nil {
		t.Fatalf("UninstallShell failed: %v", err)
	}
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Fatalf("expected snippet removed, still exists")
	}
	zb2, err := os.ReadFile(zshrcPath)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("reading .zshrc after uninstall: %v", err)
		}
	} else {
		if strings.Contains(string(zb2), guardStart) || strings.Contains(string(zb2), guardEnd) {
			t.Fatalf("guarded block still present in .zshrc after uninstall")
		}
	}

	if err := UninstallShell(); err != nil {
		t.Fatalf("UninstallShell on clean state failed: %v", err)
	}
}
