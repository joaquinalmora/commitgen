# Installation & Setup Guide

Complete installation and configuration guide for commitgen.

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew tap joaquinalmora/tap
brew install commitgen
commitgen init  # Interactive configuration
```

### Option 2: Go Install

```bash
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest
commitgen init
```

### Option 3: From Source

```bash
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
make build
make install  # Installs to /usr/local/bin
```

## Configuration

### Interactive Setup (Recommended)

```bash
commitgen init
```

This will guide you through:

- Configuration file location (local vs global)
- OpenAI API key setup
- AI model selection
- Default preferences

### Manual Configuration

#### Environment Variables (Alternative Setup)

Create a `.env` file in your home directory:

```bash
# Generate template file
commitgen env-example
cp .env.example ~/.env

# Edit with your actual API key
nano ~/.env
```

#### YAML Configuration (Advanced)

Create `commitgen.yaml`:

```yaml
ai:
  enabled: true
  provider: "openai"
  model: "gpt-4o-mini"
  api_key: ""  # Use environment variable instead

performance:
  patch_bytes: 4000
  cache_ttl: "24h"
  max_files: 10

output:
  verbose: false
  colors: true
```

## OpenAI API Setup

1. Visit [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Create a new API key
3. Add it to your configuration:
   - Via `commitgen init`
   - In `.env` file: `OPENAI_API_KEY=sk-...`
   - In `commitgen.yaml`: `api_key: "sk-..."`

## Git Integration

### Basic Usage

```bash
git add .
commitgen suggest          # Generate commit message
commitgen suggest --ai     # Force AI generation
commitgen suggest --verbose # Show detailed logs
```

### Git Hooks (Auto-suggestions)

```bash
commitgen install-hook     # Install prepare-commit-msg hook
git commit                 # Will auto-suggest messages
commitgen uninstall-hook   # Remove hook
```

> The installer writes two hooks: `.git/hooks/prepare-commit-msg` (injects the suggestion when the commit message is empty, preferring `commitgen cached`) and `.git/hooks/post-index-change` (auto-runs `commitgen cache` after every `git add`). `post-index-change` was added in Git 2.44, so uninstall the hooks if you're on an older Git release or if your hosting provider rejects unfamiliar hooks.

## Shell Integration

### Automatic Setup

```bash
commitgen install-shell
```

This command assumes **zsh** and appends a guarded block to `~/.zshrc` that sources `~/.config/commitgen.zsh`. The snippet implements the preview widget and key bindings, and you can remove it later with `commitgen uninstall-shell`.

### Manual zsh Integration

For zsh with oh-my-zsh and zsh-autosuggestions:

1. **Ensure zsh-autosuggestions is enabled**:
   ```bash
   # In ~/.zshrc plugins list
   plugins=(... zsh-autosuggestions)
   ```

2. **Configure autosuggestions strategy** (add before sourcing oh-my-zsh):
   ```bash
   ZSH_AUTOSUGGEST_USE_ASYNC=1
   ZSH_AUTOSUGGEST_STRATEGY=(commitgen history)
   ```

3. **Add commitgen function**:
   ```bash
   # Add to ~/.zshrc
   commitgen() {
     if [[ "$1" == "suggest" ]] && [[ -z "$2" ]]; then
       command commitgen suggest --plain 2>/dev/null
     else
       command commitgen "$@"
     fi
   }
   ```

4. **Reload shell**:
   ```bash
   source ~/.zshrc
   ```

### Verification

```bash
cd your-git-repo
git add some-file.txt
# Type 'git commit -m "' and press TAB
# Should show AI-generated suggestion
```

## Troubleshooting

### Common Issues

**"No API key" error:**
- Run `commitgen init` to set up configuration
- Check that `OPENAI_API_KEY` is set correctly
- Verify API key is valid at OpenAI platform

**"Not a git repository" error:**
- Ensure you're in a git repository
- Run `git init` if needed

**"No staged changes" error:**
- Stage files with `git add <files>`
- Check `git status` to see staged changes

**Shell integration not working:**
- Run `commitgen doctor` for diagnostics
- Ensure zsh-autosuggestions plugin is enabled
- Check that commitgen is in PATH

### Debug Mode

```bash
commitgen suggest --verbose  # Show detailed logs
commitgen doctor             # System diagnostics
```

### Getting Help

```bash
commitgen --help            # Show available commands
commitgen suggest --help    # Command-specific help
```

### Example `.env` entries

```bash
echo 'OPENAI_API_KEY=sk-your-key-here' >> ~/.env
echo 'COMMITGEN_MODEL=gpt-4o' >> ~/.env  # Default: gpt-4o-mini
```

## Integration Setup

### Git Hooks
```bash
commitgen install-hook  # Enable auto-cache and commit integration
```

### Shell Integration

Basic setup (recommended):
```bash
commitgen install-shell  # Automated setup
source ~/.zshrc
```

For manual setup or advanced shell configurations, see the Shell Integration section above for detailed instructions.

## Environment Configuration

Core variables:

| Variable | Purpose | Default |
|----------|---------|---------|
| `OPENAI_API_KEY` | Required API key for OpenAI | _none_ |
| `COMMITGEN_AI` | Enable AI mode automatically | `false` |
| `COMMITGEN_MODEL` | OpenAI model override | `gpt-4o-mini` |
| `COMMITGEN_BASE_URL` | Custom OpenAI-compatible endpoint | `https://api.openai.com/v1` |
| `COMMITGEN_MAX_FILES` | Max files included in prompts | `10` |
| `COMMITGEN_PATCH_BYTES` | Max diff bytes sent to AI | `102400` |
| `COMMITGEN_AI_FALLBACK` | Toggle heuristic fallback | `true` |
| `COMMITGEN_CONVENTIONS_FILE` | Custom markdown conventions | _unset_ |

## Troubleshooting

### Common Issues

#### AI not working
```bash
commitgen doctor                    # Check system status
commitgen suggest --ai --verbose    # Test with debug output
```

#### Shell integration not working  
```bash
commitgen uninstall-shell && commitgen install-shell
source ~/.zshrc
```

#### Git hooks not working
```bash
ls -la .git/hooks/prepare-commit-msg  # Check installation
commitgen uninstall-hook && commitgen install-hook
```

#### Environment not loading
```bash
ls -la ~/.env                       # Check file exists  
source ~/.env && env | grep COMMIT  # Test loading
```

## Advanced Configuration

### Custom Conventions
```bash
# Export and customize commit standards
cp internal/provider/conventions.md ~/.commitgen-conventions.md
echo "COMMITGEN_CONVENTIONS_FILE=$HOME/.commitgen-conventions.md" >> ~/.env
```

### Multiple Repositories
Each repo can have independent hook configuration:
```bash
cd /path/to/repo && commitgen install-hook
```

### Custom API Endpoints
```bash
COMMITGEN_BASE_URL=https://your-endpoint.com/v1
OPENAI_API_KEY=your-key
```

## Verification

Quick health check:
```bash
commitgen doctor  # Should show all green
```

Test workflow:
```bash
echo "test" > test.txt && git add test.txt
git commit -m ""  # Look for ghost text suggestions
```

## Uninstall

Complete removal:
```bash
commitgen uninstall-shell
commitgen uninstall-hook  # Run in each repo
sudo rm /usr/local/bin/commitgen
rm ~/.env ~/.config/commitgen.zsh
```
