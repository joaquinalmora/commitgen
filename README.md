# commitgen

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

AI-powered commit message generation for Git. Generate professional, contextual commit messages from your staged changes using OpenAI or local Ollama models with intelligent caching and seamless git workflow integration.

## ⚡ Quick Start

```bash
# Install and setup (30 seconds)
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest
curl -sSL https://raw.githubusercontent.com/joaquinalmora/commitgen/main/scripts/setup.sh | bash

# Start using AI commit messages
git add . && git commit  # AI suggestions appear automatically
```

## ✨ Features

- **🤖 AI-Powered**: OpenAI and Ollama support for professional commit messages
- **⚡ Auto-Cache**: 50x performance boost with background pre-generation
- **🔄 Git Integration**: Automatic hooks for seamless workflow
- **� Shell Integration**: Ghost text suggestions as you type
- **🛡️ Smart Fallback**: Heuristics when AI is unavailable
- **� Zero Config**: One-command setup with intelligent defaults

## 📖 Documentation

- **[Complete Documentation](docs/README.md)** - Full documentation index
- **[Setup Guide](docs/SETUP.md)** - Installation, configuration, and shell integration
- **[Technical Reference](docs/TECHNICAL.md)** - Architecture and development guide
- **[Contributing](docs/CONTRIBUTING.md)** - Development guidelines and contribution process

## 🚀 Usage

### Basic Commands

```bash
commitgen suggest --ai              # Generate AI commit message
commitgen cache                     # Pre-generate cache
commitgen install-hook              # Setup git hooks
commitgen install-shell             # Setup shell integration
commitgen doctor                    # Check system health
```

### AI Providers

| Provider | Best For | Setup |
|----------|----------|-------|
| **OpenAI** | Production quality | Add `OPENAI_API_KEY` |
| **Ollama** | Local/private | Install Ollama + model |

### Configuration

```bash
# Quick config via environment
export OPENAI_API_KEY="sk-your-key"
export COMMITGEN_PROVIDER="openai"     # or "ollama"
export COMMITGEN_MODEL="gpt-4o-mini"   # or "llama3.2:3b"
```

## 🎯 How It Works

1. **Stage changes**: `git add .`
2. **AI generates**: Background cache or real-time
3. **Shell suggests**: Ghost text appears as you type
4. **Accept/commit**: Press → or just commit normally

## 🤝 Contributing

1. Fork and create feature branch
2. Make changes with tests
3. Submit pull request

See [DOCS/TECHNICAL.md](DOCS/TECHNICAL.md) for architecture details.

## 📄 License

MIT License - see LICENSE file for details.

MIT License - see [LICENSE](LICENSE) file for details.
