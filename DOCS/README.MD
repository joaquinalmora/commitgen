# commitgen

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

AI-powered commit message generation for Git. Generate professional, contextual commit messages from your staged changes using OpenAI or local Ollama models with intelligent caching and seamless git workflow integration.

## ⚡ Quick Setup (30 seconds)

```bash
# 1. Install
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest

# 2. Setup AI (automated script)
curl -sSL https://raw.githubusercontent.com/joaquinalmora/commitgen/main/scripts/setup.sh | bash

# 3. Start using AI commit messages!
git add . && git commit  # AI suggestions appear automatically
```

**That's it!** 🎉 The setup script handles everything: environment, API keys, and shell integration.

## Features
as
- **🤖 AI-Powered**: OpenAI and Ollama integration for professional commit messages
- **🏠 Local & Cloud**: Choose between OpenAI API or local Ollama models  
- **⚡ Auto-Cache System**: Background pre-generation with 50x performance boost (3s → 0.06s)
- **🔄 Seamless Integration**: Zero-friction git workflow with dual hook system
- **📝 Smart Conventions**: Comprehensive commit message standards with company-specific rules
- **🛡️ Intelligent Fallback**: Smart heuristics when AI is unavailable
- **🚀 Shell Integration**: Ghost text suggestions in zsh (with `zsh-autosuggestions`)
- **🔧 Easy Setup**: One-command installation with automatic hook management

## Manual Setup

For advanced users or custom installations, see [SETUP.md](SETUP.md).

## Usage Examples

### Basic Usage

```bash
# Stage changes and generate AI commit message
git add .
commitgen suggest --ai
```

### Shell Integration (Automatic)

```bash
# Type either command for AI suggestions:
git commit -m "    # ← AI suggestions appear here
gc "              # ← Works with your git aliases too
```

### Git Hook Integration (Automatic)

```bash
git add .
git commit        # Opens editor with AI-generated message
```

## AI Providers

### OpenAI (Recommended)

- **Models**: gpt-4o, gpt-4o-mini, gpt-3.5-turbo
- **Setup**: Add `OPENAI_API_KEY` to your environment
- **Best for**: Production use, highest quality

### Ollama (Local, Free)

- **Models**: llama3.2:3b, qwen2.5-coder:7b, codellama:7b, deepseek-coder:6.7b
- **Setup**: Install Ollama locally
- **Best for**: Privacy, no API costs, offline use

## Performance

### Auto-Cache System

- **Background Generation**: AI suggestions pre-generated on `git add`
- **Instant Retrieval**: 0.06s response time (50x faster than realtime)
- **Smart Invalidation**: Cache updates automatically on new changes

### Benchmarks

- **Without Cache**: 3.2s (network + AI processing)
- **With Cache**: 0.06s (instant retrieval)
- **Cache Hit Rate**: 95%+ in typical workflows

## Commands

```bash
commitgen suggest [--ai] [--plain] [--verbose]  # Generate commit message
commitgen cache                                 # Pre-generate and cache message
commitgen cached                                # Show cached message
commitgen install-hook                          # Install git hooks
commitgen install-shell                         # Install shell integration  
commitgen doctor                                # System diagnostics
```

## Configuration

Environment variables for customization:

```bash
# AI Provider
export OPENAI_API_KEY="sk-your-key"           # OpenAI API key
export COMMITGEN_PROVIDER="ollama"            # Use Ollama instead
export COMMITGEN_MODEL="llama3.2:3b"          # Custom model

# Cache Settings  
export COMMITGEN_CACHE_TTL="24h"              # Cache lifetime
export COMMITGEN_MAX_FILES="10"               # Max files to analyze

# Advanced
export COMMITGEN_PATCH_BYTES="102400"         # Max patch size (100KB)
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

For technical details, see [DOCS/TECHNICAL.md](DOCS/TECHNICAL.md).

## License

MIT License - see [LICENSE](LICENSE) file for details.
