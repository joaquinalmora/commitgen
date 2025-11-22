# commitgen

<div align="center">

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)  
<em>Built with the tools and technologies:</em>  

![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=default&logo=Go&logoColor=white)
![GNU Make](https://img.shields.io/badge/GNU%20Make-000000.svg?style=default&logo=GNU&logoColor=white)
![Shell](https://img.shields.io/badge/Shell-4EAA25.svg?style=default&logo=GNU%20Bash&logoColor=white)
![GitHub](https://img.shields.io/badge/GitHub-181717.svg?style=default&logo=GitHub&logoColor=white)

</div>


## Quick Start

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
git commit -m ""  # AI suggestions appear here
```

## Features

- **AI-Powered**: OpenAI GPT-4o-mini (and other OpenAI models) for professional commit messages
- **Auto-Cache**: Intelligent caching with 50x performance boost
- **Git Integration**: Automatic hooks for seamless workflow
- **Shell Integration**: Ghost text suggestions as you type `git commit -m "`
- **Smart Fallback**: Heuristic generation when AI is unavailable
- **Zero Config**: One-command setup with intelligent defaults
- **Context-Aware**: Analyzes actual code changes, not just file names

## Configuration

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `OPENAI_API_KEY` | API key used by the OpenAI provider | _required for AI_ |
| `COMMITGEN_AI` | Enable AI automatically (otherwise pass `--ai`) | `false` |
| `COMMITGEN_MODEL` | OpenAI model name | `gpt-4o-mini` |
| `COMMITGEN_BASE_URL` | Override the OpenAI API URL for proxies/self-hosting | `https://api.openai.com/v1` |
| `COMMITGEN_MAX_FILES` | Max staged files included in the prompt | `10` |
| `COMMITGEN_PATCH_BYTES` | Max bytes of diff sent to the AI | `102400` |
| `COMMITGEN_AI_FALLBACK` | Disable (`false`) or enable (`true`) heuristic fallback | `true` |
| `COMMITGEN_CONVENTIONS_FILE` | Path to custom commit-style markdown | _unset_ |

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

## Usage Examples

### Basic Commands

```bash
commitgen suggest                       # Generate commit message
commitgen suggest --ai                  # Force AI generation
commitgen suggest --cached              # Reuse the last cached AI result
commitgen suggest --verbose             # Show detailed logs
commitgen cached --plain                # Print cached message without formatting
commitgen cache                         # Pre-generate cache
commitgen cache --clear                 # Clear cache
commitgen init                          # Interactive config (local)
commitgen init --global                 # Interactive config in ~/.commitgen.yaml
commitgen env-example                   # Write .env.example
commitgen doctor                        # System health check
commitgen version --verbose             # Include git commit + build date
```

### CLI Reference

| Command | What it does | Helpful flags |
|---------|--------------|---------------|
| `commitgen suggest` | Generates commit text from staged changes | `--ai`, `--cached`, `--plain`, `--verbose` |
| `commitgen cache` | Performs AI/heuristic generation and stores the result | `--clear`, `--verbose` |
| `commitgen cached` | Prints the most recent cached commit message (used by hooks/shell) | `--plain`, `--verbose` |
| `commitgen install-hook` / `uninstall-hook` | Manage `.git/hooks/prepare-commit-msg` and `.git/hooks/post-index-change` | _n/a_ |
| `commitgen install-shell` / `uninstall-shell` | Manage the guarded `~/.zshrc` block + `~/.config/commitgen.zsh` snippet | _n/a_ |
| `commitgen init` | Interactive YAML config generator (supports `--global`) | `--global` |
| `commitgen env-example` | Writes `.env.example` with the current defaults | _n/a_ |
| `commitgen doctor` | Runs environment checks | `--verbose` (via `COMMITGEN_AI=1` etc.) |
| `commitgen version` | Prints version/build metadata | `--verbose` |

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

> `commitgen install-hook` writes both `.git/hooks/prepare-commit-msg` (inserts the suggestion when the message is empty) and `.git/hooks/post-index-change` (warms the cache every time you run `git add`). The cache-first behavior depends on `commitgen cached`, so keep the binary accessible to your repo. `post-index-change` is new in Git 2.44, so skip the auto-cache hook (or remove it via `commitgen uninstall-hook`) if you are on an older Git release or a hosting platform that disallows it.

### Shell Integration

Shell ghost text currently targets **zsh** only. The installer writes a guarded block to `~/.zshrc` that sources `~/.config/commitgen.zsh`, which contains the preview widget and `^F` binding.

```bash
# Install shell integration (zsh)
commitgen install-shell

# Now typing 'git commit -m "' will show ghost text suggestions
git commit -m "feat: add user auth and↩ # <-- AI suggestion appears
```

The installer leaves any existing `~/.zshrc` content alone, and you can remove the snippet at any time with `commitgen uninstall-shell` (which deletes both the guard block and `~/.config/commitgen.zsh`).

## AI Providers

OpenAI is the only wired-up provider today. Local/Ollama support is still being designed, so the CLI will ignore `COMMITGEN_PROVIDER=ollama` (or similar) until the provider package grows that implementation.

| Provider | Status | Setup |
|----------|--------|-------|
| **OpenAI** | ✅ Supported | `export OPENAI_API_KEY=sk-...` then choose a model with `COMMITGEN_MODEL` |
| **Local/Ollama** | 🔜 Planned | Not implemented yet (follow project updates) |

### OpenAI Setup

1. Get API key from [OpenAI Platform](https://platform.openai.com/api-keys)
2. Set environment variable: `export OPENAI_API_KEY=sk-your-key`
3. Optional: Choose model: `export COMMITGEN_MODEL=gpt-4o-mini`

### Local/Ollama (Roadmap)

Ollama/local-hosted models are high on the roadmap, but the provider code does not talk to Ollama yet. Track the repo releases for updates if you need on-device inference.

## Troubleshooting

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

## How It Works

1. **Analyze Changes**: Reads git staged changes and file modifications
2. **Generate Context**: Creates intelligent prompts from code diffs
3. **AI Processing**: Sends context to the OpenAI provider
4. **Smart Caching**: Caches results for identical changes
5. **Integration**: Provides suggestions via git hooks or shell integration

## Documentation

- **[Setup Guide](docs/SETUP.md)** - Complete installation and configuration
- **[AI Commit Workflow](DOCS/AI_COMMIT_WORKFLOW.md)** - zsh `aicommit` function + sample git hook

Technical, contributing, and changelog docs are on the roadmap—follow the repo releases if you need those references.

## Helper Scripts

- `scripts/test-ai.sh` — sanity-checks the OpenAI flow using the same env vars as the app (`OPENAI_API_KEY`, `COMMITGEN_MODEL`, etc.).
- `scripts/export-conventions.sh` — placeholder for a future conventions exporter (currently a no-op script).
- `scripts/setup-company-conventions.sh` — placeholder for company-specific conventions bootstrap (currently a no-op script).

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make changes with tests
4. Submit a pull request

Detailed contribution guidelines are on the roadmap—feel free to open a PR or issue if you need additional context.

## Uninstall

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

## License

MIT License - see [LICENSE](LICENSE) file for details.
