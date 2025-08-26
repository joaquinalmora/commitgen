# Quick Start - commitgen

This file shows the minimal steps a new user should run to try commitgen with full AI capabilities and auto-cache system.

## 1. Build and Setup

```bash
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
go build -o bin/commitgen ./cmd/commitgen

# Set up OpenAI API key for AI-powered suggestions
export OPENAI_API_KEY="sk-your-api-key-here"
```

## 2. Install Auto-Cache System (Recommended)

```bash
# Install the complete dual-hook system for ultimate performance
./bin/commitgen install-hook
```

This enables:

- **Background AI generation** after `git add` (no delays)
- **Instant commit messages** during `git commit` (0.05s)
- **50x performance improvement** over real-time AI
- **Zero workflow changes** - use normal git commands

## 3. Optional: Install Shell Integration

```bash
# Install zsh integration for ghost text suggestions
./bin/commitgen install-shell
```

## 4. Test the Complete System

```bash
# Run diagnostics to verify everything is working
./bin/commitgen doctor

# Test AI-powered suggestions
echo "# Test feature" >> test.md
git add test.md
./bin/commitgen suggest --ai

# Test the auto-cache workflow
git commit --no-edit  # Uses cached AI message instantly
```

## Performance Comparison

```bash
# Method 1: Real-time AI (slow)
time ./bin/commitgen suggest --ai
# Result: ~3.0 seconds

# Method 2: Cached retrieval (instant)  
time ./bin/commitgen cached
# Result: ~0.06 seconds (50x faster!)

# Method 3: Auto-cache workflow (seamless)
time git add file.md && time git commit --no-edit
# Result: 0.03s + 0.05s = 0.08s total
```

## Command Examples

### AI-Powered Suggestions

```bash
# Generate AI commit message (recommended)
./bin/commitgen suggest --ai

# AI with verbose debug output
./bin/commitgen suggest --ai --verbose

# Force AI mode via environment variable
COMMITGEN_AI=1 ./bin/commitgen suggest
```

### Cache Management

```bash
# Pre-generate cache for current staged changes
./bin/commitgen cache

# Get the most recent cached message
./bin/commitgen cached

# Check if cache exists for current changes
./bin/commitgen cached >/dev/null && echo "Cache available" || echo "No cache"
```

### Fallback Heuristics

```bash
# Use heuristic-based suggestions (no AI)
./bin/commitgen suggest

# Plain output for scripting
./bin/commitgen suggest --plain

# Verbose heuristic analysis
./bin/commitgen suggest --verbose
```

## Workflow Examples

### Standard Development (with auto-cache)

```bash
# 1. Make changes
echo "export function newFeature() {}" >> src/app.js

# 2. Stage changes (triggers background AI cache)
git add src/app.js

# 3. Commit with instant AI message
git commit --no-edit
# Result: "feat(app): add new feature implementation function"
```

### Manual Cache Management

```bash
# Generate cache manually for better control
git add .
./bin/commitgen cache

# Use cached message in commit
git commit -m "$(./bin/commitgen cached)"
```

### AI + Custom Message

```bash
# Get AI suggestion as starting point
git add .
AI_MSG=$(./bin/commitgen suggest --ai)
echo "AI suggests: $AI_MSG"

# Customize and commit
git commit -m "feat: implement user authentication system

- Add JWT token validation
- Implement password hashing
- Add user session management

Based on AI suggestion: $AI_MSG"
```
