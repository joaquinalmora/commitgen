# commitgen

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

AI-powered commit message generation for Git. Generate professional, contextual commit messages from your staged changes using OpenAI or local Ollama models with intelligent caching and seamless git workflow integration.

## ⚡ Quick Start

### Installation

```bash
# Option 1: Homebrew (Recommended)
brew tap joaquinalmora/tap
brew install commitgen

# Option 2: Go Install
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest

# Option 3: From Source
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen && make build && make install
```

### Setup

```bash
# Interactive configuration (recommended)
commitgen init

# Add your OpenAI API key when prompted
# Then install git hooks and shell integration
commitgen install-hook
commitgen install-shell
```

### Basic Usage

```bash
# Stage your changes
git add .

# Generate AI-powered commit message
commitgen suggest --ai

# Or just commit and get auto-suggestions
git commit -m "  # AI suggestions appear here
```

## ✨ Features

- **🤖 AI-Powered**: OpenAI and Ollama support for professional commit messages
- **⚡ Auto-Cache**: Intelligent caching with 50x performance boost
- **🔄 Git Integration**: Automatic hooks for seamless workflow
- **👻 Shell Integration**: Ghost text suggestions as you type `git commit -m "`
- **🛡️ Smart Fallback**: Heuristic generation when AI is unavailable
- **🚀 Zero Config**: One-command setup with intelligent defaults
- **🎯 Context-Aware**: Analyzes actual code changes, not just file names

## 🔧 Configuration

### Environment Variables

```bash
# Core configuration
OPENAI_API_KEY=sk-your-key-here
COMMITGEN_PROVIDER=openai               # 'openai' or 'ollama'  
COMMITGEN_MODEL=gpt-4o-mini             # Model selection

# Performance tuning
COMMITGEN_CACHE_TTL=24h                 # Cache lifetime
COMMITGEN_MAX_FILES=10                  # Max files analyzed
COMMITGEN_PATCH_BYTES=102400            # Max patch size
```

### YAML Configuration

Create `commitgen.yaml`:

```yaml
ai:
  enabled: true
  provider: "openai"
  model: "gpt-4o-mini"
  api_key: ""  # Use environment variable

performance:
  patch_bytes: 4000
  cache_ttl: "24h"
  max_files: 10

output:
  verbose: false
  colors: true
```

## 🚀 Usage Examples

### Basic Commands

```bash
commitgen suggest                       # Generate commit message
commitgen suggest --ai                  # Force AI generation
commitgen suggest --verbose             # Show detailed logs
commitgen cache                         # Pre-generate cache
commitgen cache --clear                 # Clear cache
commitgen doctor                        # System health check
```

### Git Integration

```bash
# Install git hooks for auto-suggestions
commitgen install-hook

# Now git commit will auto-suggest messages
git commit
# Message: "feat: add user authentication with JWT tokens"

# Uninstall hooks
commitgen uninstall-hook
```

### Shell Integration

```bash
# Install shell integration (zsh)
commitgen install-shell

# Now typing 'git commit -m "' will show ghost text suggestions
git commit -m "feat: add user auth and↩ # <-- AI suggestion appears
```

## 🤖 AI Providers

| Provider | Best For | Setup |
|----------|----------|-------|
| **OpenAI** | Production quality, best results | `export OPENAI_API_KEY=sk-...` |
| **Ollama** | Local/private, no API costs | Install Ollama + model |

### OpenAI Setup

1. Get API key from [OpenAI Platform](https://platform.openai.com/api-keys)
2. Set environment variable: `export OPENAI_API_KEY=sk-your-key`
3. Optional: Choose model: `export COMMITGEN_MODEL=gpt-4o-mini`

### Ollama Setup

1. Install [Ollama](https://ollama.ai)
2. Pull a model: `ollama pull llama3.2:3b`
3. Configure: `export COMMITGEN_PROVIDER=ollama COMMITGEN_MODEL=llama3.2:3b`

## 🛠️ Troubleshooting

### Common Issues

**"No API key" error:**

- Run `commitgen init` for interactive setup
- Set `OPENAI_API_KEY` environment variable
- Verify API key at OpenAI platform

**"Not a git repository" error:**

- Ensure you're in a git repository (`git init` if needed)
- Check current directory with `pwd`

**"No staged changes" error:**

- Stage files first: `git add <files>`
- Check status: `git status`

**Shell integration not working:**

- Run `commitgen doctor` for diagnostics
- Reinstall: `commitgen uninstall-shell && commitgen install-shell`
- Source shell: `source ~/.zshrc`

### Debug Mode

```bash
commitgen suggest --verbose             # Show detailed logs
commitgen doctor                        # System diagnostics
```

## 🎯 How It Works

1. **Analyze Changes**: Reads git staged changes and file modifications
2. **Generate Context**: Creates intelligent prompts from code diffs
3. **AI Processing**: Sends context to AI provider (OpenAI/Ollama)
4. **Smart Caching**: Caches results for identical changes
5. **Integration**: Provides suggestions via git hooks or shell integration

## 📚 Documentation

- **[Setup Guide](docs/SETUP.md)** - Complete installation and configuration
- **[Technical Reference](docs/TECHNICAL.md)** - Architecture and development
- **[Contributing Guide](docs/CONTRIBUTING.md)** - Development guidelines
- **[Changelog](docs/CHANGELOG.md)** - Release history
- **[AI Commit Workflow](DOCS/AI_COMMIT_WORKFLOW.md)** - zsh `aicommit` function + sample git hook

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make changes with tests
4. Submit a pull request

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for detailed guidelines.

## � Uninstall

Complete removal:

```bash
# Remove integrations
commitgen uninstall-shell
commitgen uninstall-hook  # Run in each repository

# Remove binary (if installed manually)
sudo rm /usr/local/bin/commitgen

# Or via Homebrew
brew uninstall commitgen

# Clean up configuration
rm ~/.env ~/.config/commitgen.zsh
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.
