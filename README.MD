# commitgen

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

AI-powered commit message generation for Git. Generate professional, contextual commit messages from your staged changes using OpenAI or local Ollama models with intelligent caching and seamless git workflow integration.

## Features

- **🤖 AI-Powered**: OpenAI and Ollama integration for professional commit messages
- **🏠 Local & Cloud**: Choose between OpenAI API or local Ollama models
- **⚡ Auto-Cache System**: Background pre-generation with 50x performance boost (3s → 0.06s)
- **🔄 Seamless Integration**: Zero-friction git workflow with dual hook system
- **📝 Smart Conventions**: Comprehensive commit message standards with company-specific rules
- **🛡️ Intelligent Fallback**: Smart heuristics when AI is unavailable
- **🚀 Shell Integration**: Ghost text suggestions in zsh (with `zsh-autosuggestions`)
- **🔧 Easy Setup**: One-command installation with automatic hook management

## Quick Start

### 1. Build and Setup

```bash
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
go build -o bin/commitgen ./cmd/commitgen

# Option A: OpenAI Setup
export OPENAI_API_KEY="your-api-key-here"

# Option B: Ollama Setup (local AI)
# Install Ollama first: https://ollama.ai
ollama pull llama3.2:3b
export COMMITGEN_AI_PROVIDER="ollama"

# Test AI-powered suggestions
./bin/commitgen suggest --ai
```

### 2. Install Auto-Cache System (Recommended)

Experience the ultimate git workflow with instant, AI-powered commit messages:

```bash
./bin/commitgen install-hook
```

**What this does:**

- Installs dual git hooks for seamless automation
- `git add` → triggers background AI cache generation
- `git commit` → instantly uses cached message (0.05s)
- **50x performance improvement** over real-time generation

### 3. Install Shell Integration (Optional)

```bash
./bin/commitgen install-shell
```

Now when you type `git commit -m "` the suggestion appears as ghost text. Accept with → or Ctrl+F.

## Commands

```bash
# AI-powered suggestions (recommended)
commitgen suggest --ai [--plain] [--verbose]  # Generate AI commit message
commitgen cache                               # Pre-generate and cache message
commitgen cached                              # Get most recent cached message

# Fallback heuristics
commitgen suggest [--plain] [--verbose]       # Generate heuristic-based message

# System management
commitgen doctor                              # Verify installation and environment
commitgen install-shell                       # Install zsh shell integration
commitgen uninstall-shell                    # Remove zsh shell integration  
commitgen install-hook                        # Install auto-cache git hooks
commitgen uninstall-hook                     # Remove git hooks
```

## How it works

### AI Mode (Default with `--ai`)

commitgen integrates with multiple AI providers to generate professional commit messages:

- **OpenAI Integration**: Cloud-based models (gpt-4o, gpt-4o-mini, gpt-3.5-turbo)
- **Ollama Support**: Local models (llama3.2, qwen2.5-coder, codellama, deepseek-coder)
- **Analyzes staged changes** → understands context and impact
- **Applies company conventions** → follows established standards
- **Generates professional messages** → proper scope, type, and description
- **Intelligent caching** → 24-hour cache with hash-based deduplication
- **Smart fallback** → uses heuristics if AI is unavailable

### Auto-Cache System

The dual git hook system provides seamless automation:

1. **`git add`** → triggers `post-index-change` hook → background cache generation
2. **`git commit`** → triggers `prepare-commit-msg` hook → instant cached message insertion
3. **Performance**: 3000ms → 60ms (50x improvement)

### Heuristic Mode (Fallback)

When AI is unavailable, intelligent heuristics analyze staged changes:

- **Documentation-only changes** → "docs: update documentation"
- **Test-only changes** → "test: update tests"  
- **Configuration-only changes** → "config: update configuration"
- **Rename-only changes** → "refactor: rename files"
- **Mixed changes** → "feat: update file1, file2 (+N more)"

## Configuration

### AI Provider Setup

#### OpenAI (Cloud)

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="sk-your-api-key-here"
export COMMITGEN_AI_PROVIDER="openai"          # Default
export COMMITGEN_AI_MODEL="gpt-4o-mini"        # Default model
```

#### Ollama (Local)

```bash
# Install and setup Ollama
curl -fsSL https://ollama.ai/install.sh | sh
ollama pull llama3.2:3b

# Configure commitgen
export COMMITGEN_AI_PROVIDER="ollama" 
export COMMITGEN_AI_MODEL="llama3.2:3b"        # Default
export COMMITGEN_AI_BASE_URL="http://localhost:11434"  # Default
```

#### Alternative Models

```bash
# OpenAI models
export COMMITGEN_AI_MODEL="gpt-4o"
export COMMITGEN_AI_MODEL="gpt-4o-mini"
export COMMITGEN_AI_MODEL="gpt-3.5-turbo"

# Ollama models  
ollama pull qwen2.5-coder:7b && export COMMITGEN_AI_MODEL="qwen2.5-coder:7b"
ollama pull codellama:7b && export COMMITGEN_AI_MODEL="codellama:7b"
ollama pull deepseek-coder:6.7b && export COMMITGEN_AI_MODEL="deepseek-coder:6.7b"
```

### Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key for cloud AI suggestions
- `COMMITGEN_AI_PROVIDER`: AI provider (`openai` or `ollama`, default: `openai`)
- `COMMITGEN_AI_MODEL`: Model to use (provider-specific defaults)
- `COMMITGEN_AI_BASE_URL`: Custom API endpoint (provider-specific defaults)
- `COMMITGEN_AI=1`: Force AI mode (equivalent to `--ai` flag)
- `COMMITGEN_DEBUG=1`: Enable verbose output for debugging

## Development

### Run tests

```bash
go test ./...
```

### Run diagnostics

```bash
./bin/commitgen doctor
```

### Cache Management

```bash
# View cache contents
./bin/commitgen cached

# Force regenerate cache
./bin/commitgen cache

# Cache location: ~/.cache/commitgen/
```

## Privacy & Security

- **OpenAI**: Only sends git diff output to OpenAI (no personal data)
- **Ollama**: Complete privacy with local AI models (no data sent anywhere)
- **Local Processing**: All analysis and caching happens locally
- **Secure Storage**: Cache uses SHA256 hashing for content identification
- **No Persistence**: No commit history or personal information stored
- **Optional Integration**: All AI features are opt-in

## Performance

| Operation | Time | Notes |
|-----------|------|-------|
| `git add` | ~0.03s | Triggers background cache (no delay) |
| `git commit` | ~0.05s | Uses cached message (instant) |
| Real-time AI | ~3.0s | Fallback when cache unavailable |
| **Improvement** | **50x faster** | Background caching vs real-time |

## What's Next

✅ **AI integration** is complete and fully operational with multiple providers!

- **Auto-cache system** provides instant commit messages
- **OpenAI & Ollama integration** for cloud and local AI models
- **Dual git hooks** enable seamless developer workflow
- **Intelligent fallback** ensures reliability across all scenarios

See `DOCS/USAGE.md` for detailed examples and `DOCS/TECHNICAL.md` for implementation details.
