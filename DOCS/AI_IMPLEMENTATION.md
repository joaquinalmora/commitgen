# AI Implementation Guide

This document outlines the AI implementation strategy for commitgen, building on the solid foundation already in place.

## Current Status ✅

### Foundation Complete
- ✅ Provider interface defined in `internal/provider/provider.go`
- ✅ Configuration system in `internal/config/config.go` 
- ✅ CLI integration with `--ai` flag in `cmd/commitgen/main.go`
- ✅ Fallback mechanism to heuristics when AI fails
- ✅ Environment variable support (`COMMITGEN_AI=1`, `COMMITGEN_AI_PROVIDER`, etc.)

### Ready for Implementation
The AI integration follows a clean architecture:

```
suggest() -> AI Provider -> LLM API -> Response
    ↓ (on error)
   Heuristics (fallback)
```

## Environment Variables

```bash
# Enable AI mode
export COMMITGEN_AI=1

# Provider configuration  
export COMMITGEN_AI_PROVIDER=openai        # or "ollama"
export COMMITGEN_AI_API_KEY=sk-xxxxx
export COMMITGEN_AI_MODEL=gpt-4o-mini      # or "llama3.2:3b"
export COMMITGEN_AI_BASE_URL=               # custom endpoint (optional)

# Advanced settings
export COMMITGEN_MAX_FILES=10               # max files to analyze
export COMMITGEN_PATCH_BYTES=102400         # max patch size (100KB)
export COMMITGEN_AI_FALLBACK=true           # fallback to heuristics on error
```

## Implementation Roadmap

### Phase 1: OpenAI Provider (Recommended First)
**File:** `internal/provider/openai.go`

```go
type OpenAIProvider struct {
    apiKey  string
    model   string
    baseURL string
    client  *http.Client
}

func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, files []string, patch string) (string, error) {
    // 1. Build prompt from files + patch
    // 2. Call OpenAI API
    // 3. Parse response and return clean commit message
}
```

**Key considerations:**
- Use `gpt-4o-mini` as default (cost-effective, fast)
- Prompt engineering for concise, conventional commits
- Proper error handling and rate limiting
- Token limit awareness (truncate large patches)

### Phase 2: Ollama Provider (Local LLM Support)
**File:** `internal/provider/ollama.go`

```go
type OllamaProvider struct {
    baseURL string
    model   string
    client  *http.Client
}
```

**Benefits:**
- Privacy-first (everything local)
- No API costs
- Support for models like `llama3.2:3b`, `qwen2.5-coder`

### Phase 3: Prompt Engineering
**File:** `internal/provider/prompt.go`

Create consistent prompts that work across providers:

```
You are a git commit message generator. Generate a concise, conventional commit message.

Files changed: file1.go, file2.md
Diff preview:
[truncated patch]

Rules:
- Use conventional commit format (feat:, fix:, docs:, etc.)
- Keep under 50 characters for subject line
- Be specific but concise
- Focus on what changed, not how

Output only the commit message, nothing else.
```

## Integration Points

### CLI Usage
```bash
# Use AI with OpenAI
export COMMITGEN_AI_API_KEY=sk-xxxxx
./bin/commitgen suggest --ai

# Use AI with Ollama (local)
export COMMITGEN_AI_PROVIDER=ollama
export COMMITGEN_AI_MODEL=llama3.2:3b
./bin/commitgen suggest --ai

# Environment-based (no flags needed)
export COMMITGEN_AI=1
./bin/commitgen suggest
```

### Shell Integration
The existing zsh integration will automatically use AI when configured:

```bash
git commit -m "  # Shows AI-generated suggestion as ghost text
```

### Git Hook Integration  
The prepare-commit-msg hook will automatically use AI when enabled.

## Testing Strategy

### Unit Tests
- Mock provider responses
- Test fallback mechanisms  
- Validate prompt construction

### Integration Tests
- Test with real API keys (optional, CI)
- Test Ollama local server
- Test error scenarios

### Manual Testing
```bash
# Test different scenarios
echo "test" > test.txt && git add test.txt
./bin/commitgen suggest --ai --verbose

# Test fallback
COMMITGEN_AI_API_KEY=invalid ./bin/commitgen suggest --ai --verbose
```

## Error Handling

The system gracefully handles:
- Missing API keys → fallback to heuristics
- Network errors → fallback to heuristics  
- Rate limiting → fallback to heuristics
- Invalid responses → fallback to heuristics

Verbose mode shows what's happening:
```
$ ./bin/commitgen suggest --ai --verbose
Using AI provider: openai
AI generation error: rate limit exceeded
Falling back to heuristics
Update test.txt
```

## Next Steps

1. **Implement OpenAI provider** - Start with basic API integration
2. **Add prompt engineering** - Optimize for quality commit messages
3. **Implement Ollama provider** - Support local LLMs
4. **Add comprehensive tests** - Ensure reliability
5. **Update documentation** - User guide for AI features

## Technical Notes

- The `context.Context` supports cancellation for long-running requests
- Provider interface allows easy addition of new LLM services
- Configuration is environment-driven for security (no API keys in files)
- Fallback ensures the tool always works, even when AI fails

The foundation is solid and ready for AI implementation. Start with OpenAI provider as it's the most straightforward to implement and test.
