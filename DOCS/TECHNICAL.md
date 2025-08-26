# commitgen — Technical Reference

This document explains the code layout, runtime behavior, and developer-facing details useful when extending or packaging the project.

## Repository Layout (Key Files)

- `cmd/commitgen/main.go` — CLI entrypoint and command wiring with AI integration
- `internal/diff/diff.go` — Git interactions to read staged files and produce trimmed patch
- `internal/prompt/prompt.go` — Heuristics for fallback commit message generation
- `internal/provider/` — AI provider implementations (OpenAI, etc.)
- `internal/cache/cache.go` — High-performance caching system for AI responses
- `internal/hook/hook.go` — Dual git hook installer (prepare-commit-msg + post-index-change)
- `internal/shell/shell.go` — zsh plugin-first snippet writer and `.zshrc` block manager
- `internal/doctor/doctor.go` — Environment & installation diagnostics

## Commands (What They Do)

### Core Commands

- `commitgen suggest [--ai] [--plain] [--verbose]` — Generate commit message suggestion
  - `--ai`: Use AI provider (OpenAI) for professional commit messages
  - `--plain`: Output exactly one trimmed subject line (shell/script friendly)
  - `--verbose`: Print diagnostic output to stderr
  - Without `--ai`: Use heuristic-based fallback generation

### Cache Management

- `commitgen cache` — Pre-generate and cache AI commit message for staged changes
- `commitgen cached` — Retrieve most recent cached commit message (instant)

### Installation Management

- `commitgen install-hook` — Install dual git hook system for auto-cache workflow
- `commitgen uninstall-hook` — Remove git hooks with backup restoration
- `commitgen install-shell` — Install zsh shell integration snippet
- `commitgen uninstall-shell` — Remove zsh integration
- `commitgen doctor` — Run comprehensive environment diagnostics

## AI Integration Architecture

### Provider System

```text
suggest() → config.Load() → provider.GetProvider() → AI API Call
                                     ↓
                              Fallback to heuristics
```

**Current Providers:**

- `internal/provider/openai.go` — OpenAI API integration with conventional commits
- `internal/provider/provider.go` — Abstract interface for extensibility

### Configuration

Environment variables control AI behavior:

- `OPENAI_API_KEY` — OpenAI API key for AI suggestions
- `COMMITGEN_AI=1` — Force AI mode (equivalent to `--ai` flag)  
- `COMMITGEN_DEBUG=1` — Enable verbose debug output

## Auto-Cache System

### Architecture Overview

The cache system provides 50x performance improvement by pre-generating AI responses:

```text
git add → post-index-change hook → background cache generation
git commit → prepare-commit-msg hook → instant cached message insertion
```

### Cache Storage

```text
~/.cache/commitgen/
├── messages/
│   ├── a1b2c3d4.json    # Cached message with metadata  
│   └── e5f6g7h8.json    # Content hash → message mapping
└── latest.json          # Pointer to most recent cache entry
```

### Cache Lifecycle

1. **Content Hashing**: SHA256 of staged diff for deduplication
2. **Generation**: AI API call with commit conventions  
3. **Storage**: JSON with message, timestamp, hash metadata
4. **Expiry**: 24-hour automatic cleanup
5. **Retrieval**: Hash-based lookup for instant access

## Dual Git Hook System

### Hook Implementation

**prepare-commit-msg Hook:**

- Triggers on `git commit` before editor opens
- Tries cached message first (instant: ~0.05s)
- Falls back to real-time AI generation if no cache
- Inserts message into commit file automatically
- Preserves existing commit messages (non-destructive)

**post-index-change Hook:**

- Triggers on `git add` when index changes
- Runs cache generation in background (non-blocking)
- No impact on `git add` performance
- Ensures fresh cache available for next commit

### Installation Safety

- Backs up existing hooks automatically
- Restores backups on uninstall
- Detects commitgen-created hooks to prevent conflicts
- Cross-platform shell compatibility

## Performance Benchmarks

### Real-World Measurements

```bash
# Real-time AI (no cache)
time ./bin/commitgen suggest --ai
# Result: 3.142s user time

# Cached retrieval  
time ./bin/commitgen cached
# Result: 0.060s user time (52x faster)

# Complete workflow with auto-cache
time git add . && time git commit --no-edit  
# Result: 0.030s + 0.051s = 0.081s total
```

### Cache Hit Scenarios

- **Identical content**: 100% cache hit (SHA256 match)
- **Modified content**: New generation triggered  
- **24h expiry**: Automatic cleanup prevents stale data
- **Storage efficiency**: Deduplication via content hashing

## Suggest Flow (High-Level)

### AI Mode (Default with `--ai`)

1. `suggest --ai` calls `internal/diff.StagedChanges(limit)` to get `files []string` and trimmed `patch string`
2. `config.Load()` determines AI provider settings and API configuration
3. `provider.GetProvider()` creates configured AI provider (OpenAI, etc.)
4. `provider.GenerateCommitMessage(files, patch)` sends request to AI API
5. AI response parsed and formatted according to conventional commits
6. Output formatted by flags (`--plain` vs normal; debug to stderr with `--verbose`)
7. Fallback to heuristics if AI fails

### Heuristic Mode (Fallback)

1. `suggest` (no `--ai`) calls `internal/diff.StagedChanges(limit)`
2. `prompt.MakePrompt(files, patch)` applies bucket heuristics:
   - `isTestsOnly` / `isDocsOnly` / `isConfigOnly` / `isRenameOnly`
   - Default: `feat: update <file1>, <file2> (+N more)`
3. Output formatted according to flags

### Cache Flow

1. `cache` command generates message via AI mode and stores result
2. Content hashed with SHA256 for deduplication  
3. Message stored in `~/.cache/commitgen/messages/<hash>.json`
4. `latest.json` updated to point to newest cache entry
5. `cached` command retrieves via hash lookup (instant)

## Plugin-First zsh Design

This project uses a plugin-first approach for zsh inline suggestions with complete AI integration.

### Goal

Surface AI-generated, professional commit suggestions as inline ghost text in interactive zsh shells. Provide safe, conservative fallback for users without `zsh-autosuggestions`.

### User-Facing Contract

- **Input**: Local repository with staged changes + OpenAI API access
- **Output**: Professional conventional commit message (AI-generated)
- **Fallback**: Intelligent heuristics when AI unavailable
- **Privacy**: Only git diff sent to AI (no personal data)
- **Performance**: Auto-cache provides instant responses

### Plugin-First Behavior (zsh-autosuggestions)

- Detects `zsh-autosuggestions` using environment vars and plugin paths
- Provides strategy function `_zsh_autosuggest_strategy_commitgen`
- Calls `commitgen cached` first (instant), falls back to `commitgen suggest --ai`
- Prepends strategy for prioritization over history suggestions

**Rationale**: Delivers true inline ghost text with native accept/decline behavior (right-arrow to accept). Non-invasive: only reads staged files, doesn't modify commit files.

### Native Fallback (POSTDISPLAY)

- If plugin absent, shows dim preview using `zle -M`
- Provides widget `cg-accept-preview` bound to Ctrl-F
- Conservative approach: explicit user acceptance required

## Development

### Run Tests

```bash
go test ./...
```

### Run Diagnostics

```bash
./bin/commitgen doctor
```

### Cache Development Tools

```bash
# View cache contents
./bin/commitgen cached

# Force regenerate cache
./bin/commitgen cache

# Cache location
ls -la ~/.cache/commitgen/
```

### Implementation Notes

- Shell snippet: `internal/shell/commitgen.zsh` → `~/.config/commitgen.zsh`
- Cache suggestion is synchronous and fast (0.06s via auto-cache)
- Store cache in `XDG_CACHE_HOME` with user-only permissions
- Hooks are idempotent and tolerant of re-installation

## Quality Assurance

### Testing Strategy

- **Unit Tests**: `internal/prompt`, `internal/cache`, `internal/provider`
- **Integration Tests**: End-to-end workflow with git operations
- **Manual QA**: Interactive testing in zsh with/without plugins
- **Performance Tests**: Cache hit rates and response times

### Edge Cases

- **Slow AI responses**: Auto-cache eliminates delays  
- **Plugin detection**: Supports Oh My Zsh, antigen, zplug
- **Shell customizations**: Snippet is idempotent and re-source safe
- **Network failures**: Graceful fallback to heuristics

## Extension Points

### Adding New AI Providers

1. Implement `internal/provider/Provider` interface
2. Add provider detection in `provider.GetProvider()`
3. Update configuration documentation
4. Add comprehensive error handling

### Cache System Extensions

- **Custom TTL**: Configurable expiry times
- **Compression**: Reduce storage footprint  
- **Distributed Cache**: Share between team members
- **Metrics**: Cache hit rates and performance tracking

**Note**: Use `commitgen doctor` to validate installations and diagnose issues. Keep this document synchronized with code changes.

