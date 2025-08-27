# Setup Guide

Comprehensive installation and configuration guide for commitgen with AI integration.

## Quick Install

### Option 1: Automated Setup (Recommended)

```bash
# Install commitgen
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest

# Run automated setup script
curl -sSL https://raw.githubusercontent.com/joaquinalmora/commitgen/main/scripts/setup.sh | bash
```

The script handles:
- Environment file creation
- API key configuration prompts
- Shell integration installation
- AI connection testing

### Option 2: Manual Installation

```bash
# Clone and build
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
go build -o bin/commitgen ./cmd/commitgen

# Add to PATH (optional)
sudo cp bin/commitgen /usr/local/bin/
```

## AI Provider Setup

### OpenAI Setup (Cloud AI)

1. **Get API Key**: Visit [OpenAI API Keys](https://platform.openai.com/api-keys)

2. **Configure Environment**:
   ```bash
   # Copy environment template
   cp .env.example ~/.env
   
   # Edit and add your API key
   nano ~/.env
   ```

3. **Add to ~/.env**:
   ```bash
   OPENAI_API_KEY=sk-your-actual-api-key-here
   COMMITGEN_PROVIDER=openai
   COMMITGEN_MODEL=gpt-4o-mini
   ```

4. **Test Connection**:
   ```bash
   source ~/.zshrc
   commitgen suggest --ai --verbose
   ```

### Ollama Setup (Local AI)

1. **Install Ollama**:
   ```bash
   # macOS
   brew install ollama
   
   # Linux
   curl -fsSL https://ollama.ai/install.sh | sh
   ```

2. **Start Ollama Service**:
   ```bash
   ollama serve
   ```

3. **Pull a Model**:
   ```bash
   # Recommended: lightweight and fast
   ollama pull llama3.2:3b
   
   # Alternative: better for coding
   ollama pull qwen2.5-coder:7b
   ```

4. **Configure Environment**:
   ```bash
   # Add to ~/.env
   COMMITGEN_PROVIDER=ollama
   COMMITGEN_MODEL=llama3.2:3b
   OLLAMA_HOST=http://localhost:11434
   ```

## Integration Setup

### Git Hooks (Auto-Cache)

Enable background AI cache generation for instant commit messages:

```bash
# Install hooks
commitgen install-hook

# Test the workflow
echo "test" > test.txt
git add test.txt    # ← Cache generation starts in background
git commit          # ← Uses cached message instantly
```

**What gets installed:**
- `prepare-commit-msg`: Uses cached messages during commit
- `post-index-change`: Generates cache when files are staged

### Shell Integration (Ghost Text)

Enable AI suggestions as you type commit commands:

```bash
# Install shell integration
commitgen install-shell

# Reload shell
source ~/.zshrc
```

**How it works:**
- Type `git commit -m "` → AI suggestion appears as ghost text
- Type `gc "` → Works with git aliases too
- Press `→` or `Ctrl+F` to accept suggestion

## Environment Configuration

### Security Best Practices

✅ **DO:**
- Use `.env` files for API keys
- Keep `.env` files local (never commit them)
- Use the provided `.env.example` as a template

❌ **DON'T:**
- Put API keys directly in shell config files
- Commit API keys to version control
- Share API keys in screenshots or logs

### Environment Variables

```bash
# AI Provider Configuration
OPENAI_API_KEY=sk-your-key              # OpenAI API key
COMMITGEN_PROVIDER=openai               # Provider: 'openai' or 'ollama'
COMMITGEN_MODEL=gpt-4o-mini             # Model name
OLLAMA_HOST=http://localhost:11434      # Ollama server URL

# Performance Settings
COMMITGEN_CACHE_TTL=24h                 # Cache lifetime
COMMITGEN_MAX_FILES=10                  # Max files to analyze
COMMITGEN_PATCH_BYTES=102400            # Max patch size (100KB)

# Advanced Settings
COMMITGEN_AI_FALLBACK=true              # Fallback to heuristics on AI failure
COMMITGEN_VERBOSE=false                 # Enable verbose logging
```

## Troubleshooting

### Common Issues

**Q: AI not working**
```bash
# Check system status
commitgen doctor

# Test AI with verbose output
commitgen suggest --ai --verbose
```

**Q: Shell integration not working**
```bash
# Reinstall shell integration
commitgen uninstall-shell
commitgen install-shell
source ~/.zshrc
```

**Q: Git hooks not working**
```bash
# Check hook installation
ls -la .git/hooks/prepare-commit-msg
ls -la .git/hooks/post-index-change

# Reinstall hooks
commitgen uninstall-hook
commitgen install-hook
```

**Q: Environment variables not loading**
```bash
# Check if .env exists
ls -la ~/.env

# Check shell configuration
tail -5 ~/.zshrc

# Manually load environment
source ~/.env
```

### Advanced Configuration

#### Custom API Endpoints

For custom OpenAI-compatible endpoints:

```bash
OPENAI_API_KEY=your-key
OPENAI_BASE_URL=https://your-custom-endpoint.com/v1
```

#### Multiple Git Repositories

Each repository can have its own hook configuration:

```bash
# Install hooks per repository
cd /path/to/repo1
commitgen install-hook

cd /path/to/repo2  
commitgen install-hook
```

#### Shell Compatibility

Supported shells:
- ✅ zsh (recommended)
- ✅ zsh + oh-my-zsh
- ✅ zsh + zsh-autosuggestions

For other shells, use git hooks without shell integration.

## Verification

### Test Your Setup

1. **Test AI Connection**:
   ```bash
   commitgen suggest --ai --verbose
   ```

2. **Test Shell Integration**:
   ```bash
   # Type this and look for ghost text:
   git commit -m "
   ```

3. **Test Git Hooks**:
   ```bash
   echo "test" > test.txt
   git add test.txt
   git commit  # Should show AI-generated message
   ```

4. **Check System Health**:
   ```bash
   commitgen doctor
   ```

### Expected Output

When everything is working:

```bash
$ commitgen doctor
Git repo: ok
commitgen binary on PATH: /usr/local/bin/commitgen
prepare-commit-msg hook: .git/hooks/prepare-commit-msg
zsh snippet installed: ~/.config/commitgen.zsh
.zshrc contains commitgen guarded block
staged files: 1
zsh-autosuggestions: detected
```

## Uninstallation

To remove commitgen completely:

```bash
# Remove shell integration
commitgen uninstall-shell

# Remove git hooks (per repository)
commitgen uninstall-hook

# Remove binary
sudo rm /usr/local/bin/commitgen

# Remove environment file
rm ~/.env

# Remove configuration
rm -rf ~/.config/commitgen.zsh
rm -rf ~/.cache/commitgen
```
