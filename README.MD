# commitgen

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

Intelligent commit message suggestions for Git. Generate contextual commit messages from your staged changes using smart heuristics, with AI enhancement coming soon.

## Features

- **Smart heuristics**: Detects documentation, tests, configuration, and rename-only changes
- **Shell integration**: Ghost text suggestions in zsh (with `zsh-autosuggestions`)
- **Git hooks**: Optional prepare-commit-msg hook for automatic suggestions
- **Clean CLI**: Simple commands with `--plain` mode for scripting
- **Privacy-first**: All processing happens locally, no data sent anywhere

## Quick Start

### 1. Build and test

```bash
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
go build -o bin/commitgen ./cmd/commitgen

# Test with your staged changes
./bin/commitgen suggest
```

### 2. Install shell integration (optional)

```bash
./bin/commitgen install-shell
```

Now when you type `git commit -m "` the suggestion appears as ghost text. Accept with → or Ctrl+F.

### 3. Install git hook (optional)

```bash
./bin/commitgen install-hook
```

Suggestions will automatically populate in your commit editor.

## Commands

```bash
commitgen suggest [--plain] [--verbose]  # Generate commit message suggestion
commitgen doctor                         # Verify installation and environment
commitgen install-shell                  # Install zsh shell integration
commitgen uninstall-shell               # Remove zsh shell integration  
commitgen install-hook                   # Install git prepare-commit-msg hook
commitgen uninstall-hook                # Remove git hook
```

## How it works

commitgen analyzes your staged changes and applies intelligent heuristics:

- **Documentation-only changes** → "Update documentation"
- **Test-only changes** → "Update tests"
- **Configuration-only changes** → "Update configuration"
- **Rename-only changes** → "Rename files"
- **Mixed changes** → "Update file1, file2 (+N more)"

## Development

### Run tests

```bash
go test ./...
```

### Run diagnostics

```bash
./bin/commitgen doctor
```

## Privacy & Safety

- All analysis happens locally on your machine
- No data is transmitted to external services  
- Only reads staged files and git diff output
- Optional hooks do not modify existing commit messages

## What's Next

🚀 **AI integration** is coming next! The foundation is complete and ready for LLM provider integration to complement the existing heuristics.

See `DOCS/ROADMAP.md` for development status and `DOCS/USAGE.md` for detailed usage examples.
