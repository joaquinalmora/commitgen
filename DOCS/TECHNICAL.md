# Technical Reference

Architecture and development guide for commitgen.

## Architecture

```text
CLI → Config → Provider → AI API → Response
             ↓
        Cache System ← Git Hooks
             ↓  
     Shell Integration
```

## Key Components

| Component | Purpose | Location |
|-----------|---------|----------|
| **CLI** | Command routing & flags | `cmd/commitgen/main.go` |
| **Config** | Environment management | `internal/config/config.go` |
| **Providers** | AI integrations | `internal/provider/*.go` |
| **Cache** | Performance optimization | `internal/cache/cache.go` |
| **Hooks** | Git integration | `internal/hook/hook.go` |
| **Shell** | Terminal integration | `internal/shell/shell.go` |
| **Doctor** | System diagnostics | `internal/doctor/doctor.go` |

## Provider System

### Interface
```go
type Provider interface {
    GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error)
}
```

### Adding New Providers
1. Implement `Provider` interface
2. Add to factory in `provider.go`
3. Add configuration options
4. Update documentation

## Cache Architecture

**Goal**: Sub-100ms commit message retrieval  
**Method**: Background pre-generation via git hooks  
**Storage**: Content-addressed JSON files  
**Invalidation**: SHA256 hash-based change detection  

### Workflow
1. `git add` → `post-index-change` hook → background cache generation
2. `git commit` → `prepare-commit-msg` hook → instant cache retrieval  
3. Cache miss → fallback to real-time AI generation

## Performance Design

- **Content Addressing**: Prevents duplicate AI calls for same changes
- **Background Processing**: Non-blocking cache generation
- **Graceful Degradation**: AI → Cache → Heuristics → User input
- **Timeout Handling**: 30s AI timeout with fallback

## Development

### Testing
```bash
go test ./...                              # Unit tests
OPENAI_API_KEY=sk-xxx go test ./e2e/...    # Integration tests  
commitgen doctor                           # System validation
```

### Debugging
```bash
commitgen suggest --verbose    # Debug AI generation
commitgen cache --debug        # Debug cache system
ls -la .git/hooks/             # Check hook installation
```

### Code Style
- Follow Go conventions (`gofmt`, `golint`)
- Add tests for new functionality
- Update docs for API changes
- Use conventional commits

This architecture prioritizes performance, reliability, and extensibility while maintaining clean separation of concerns.
