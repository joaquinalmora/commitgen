# Technical Reference

Developer-facing documentation for extending and maintaining commitgen.

## Architecture Overview

commitgen follows a modular architecture with clear separation of concerns:

```text
CLI → Config → Provider → AI API → Response
             ↓
        Cache System ← Git Hooks
             ↓  
     Shell Integration
```

## Repository Layout

**Core Components:**
- `cmd/commitgen/main.go` — CLI entrypoint and command routing
- `internal/config/config.go` — Environment configuration management
- `internal/diff/diff.go` — Git interactions and patch generation
- `internal/prompt/prompt.go` — Heuristic fallback message generation

**AI Integration:**
- `internal/provider/provider.go` — Provider interface and factory
- `internal/provider/openai.go` — OpenAI API integration  
- `internal/provider/ollama.go` — Ollama local AI integration
- `internal/provider/conventions.md` — Commit message standards

**Performance & Integration:**
- `internal/cache/cache.go` — High-performance caching system
- `internal/hook/hook.go` — Git hook installer and manager
- `internal/shell/shell.go` — Shell integration (zsh autosuggestions)
- `internal/doctor/doctor.go` — System diagnostics

## AI Provider System

### Interface

```go
type Provider interface {
    GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error)
}
```

### Implementation

```go
// Factory pattern for provider selection
func GetProvider(cfg *Config) (Provider, error) {
    switch cfg.Provider {
    case "openai":
        return NewOpenAIProvider(cfg.OpenAI), nil
    case "ollama":
        return NewOllamaProvider(cfg.Ollama), nil
    default:
        return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
    }
}
```

### Configuration

Environment variables control AI behavior:

```bash
OPENAI_API_KEY=sk-xxx          # OpenAI API authentication
COMMITGEN_PROVIDER=ollama      # Provider selection
COMMITGEN_MODEL=llama3.2:3b    # Model override
OLLAMA_HOST=localhost:11434    # Ollama server endpoint
```

## Cache System

### Performance Goals

- **Target**: Sub-100ms commit message retrieval
- **Method**: Background pre-generation via git hooks
- **Storage**: Content-addressed JSON files
- **Invalidation**: SHA256 hash-based change detection

### Implementation

```go
type CacheEntry struct {
    Message   string    `json:"message"`
    Files     []string  `json:"files"`
    DiffHash  string    `json:"diff_hash"`
    Timestamp time.Time `json:"timestamp"`
    Provider  string    `json:"provider"`
}
```

### Workflow

1. `git add` → `post-index-change` hook → background cache generation
2. `git commit` → `prepare-commit-msg` hook → instant cache retrieval
3. Cache miss → fallback to real-time AI generation

## Git Hook System

### Hooks Installed

**prepare-commit-msg**: Injects cached AI messages into commit editor
```bash
#!/bin/sh
# .git/hooks/prepare-commit-msg
commitgen cached > "$1" 2>/dev/null || true
```

**post-index-change**: Triggers background cache generation  
```bash  
#!/bin/sh
# .git/hooks/post-index-change  
commitgen cache >/dev/null 2>&1 &
```

### Hook Management

```go
func InstallHooks(repoPath string) error {
    // Backup existing hooks
    // Install new hooks with proper permissions
    // Validate installation
}

func UninstallHooks(repoPath string) error {
    // Restore backed up hooks
    // Clean up commitgen hooks
}
```

## Shell Integration

### zsh Autosuggestions

The shell integration provides ghost text suggestions:

```bash
git commit -m "   # ← AI suggestion appears here
```

### Implementation Strategy

1. **Plugin Detection**: Check for `zsh-autosuggestions`
2. **Strategy Injection**: Add `commitgen` to suggestion strategies
3. **Context Matching**: Trigger only on `git commit -m "` patterns
4. **API Integration**: Call `commitgen suggest --ai --plain`

## Development Guidelines

### Adding New AI Providers

1. Implement the `Provider` interface
2. Add configuration options to `internal/config/`
3. Register in provider factory
4. Add tests for API integration
5. Update documentation

### Testing Strategy

**Unit Tests**: Mock provider responses, test error handling
```bash
go test ./internal/provider/...
```

**Integration Tests**: Test with real APIs (optional)
```bash
OPENAI_API_KEY=sk-xxx go test ./e2e/...
```

**Manual Testing**: End-to-end workflow validation
```bash
git add . && commitgen suggest --ai --verbose
```

### Error Handling

commitgen follows graceful degradation:

1. **AI Provider Failure** → Fallback to heuristics  
2. **Network Issues** → Use cached message if available
3. **Cache Miss** → Real-time generation
4. **Total Failure** → Empty message (let user write manually)

### Performance Considerations

**Cache Strategy**: Content-addressed storage prevents duplicate work
**Background Processing**: Non-blocking cache generation via git hooks  
**Timeout Handling**: 30-second timeout for AI API calls
**Memory Usage**: Streaming diff processing for large repositories

## Troubleshooting

### Common Development Issues

**Provider Not Found**: Check factory registration in `provider.go`
**Cache Not Working**: Verify git hook installation and permissions
**Shell Integration Issues**: Check zsh plugin compatibility
**API Timeouts**: Increase timeout or implement retry logic

### Debugging Tools

```bash
commitgen doctor              # System health check
commitgen suggest --verbose   # Debug AI generation  
commitgen cache --debug       # Debug cache system
ls -la .git/hooks/            # Check hook installation
```

## Contributing

### Code Style

- Follow Go conventions (`gofmt`, `golint`)
- Add tests for new functionality  
- Update documentation for API changes
- Use conventional commits for development

### Release Process

1. Update version in `cmd/commitgen/main.go`
2. Create git tag: `git tag v1.x.x`
3. Push tag: `git push origin v1.x.x`
4. GitHub Actions handles release automation

This architecture supports the project's goals of high performance, reliability, and extensibility while maintaining a clean separation between AI providers, caching, and user interfaces.
