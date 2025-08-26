# Quick Start - commitgen

This file shows the minimal steps a new user should run to try commitgen with full AI capabilities (OpenAI and Ollama) and auto-cache system.

## 1. Build and Setup

```bash
git clone https://github.com/joaquinalmora/commitgen.git
cd commitgen
go build -o bin/commitgen ./cmd/commitgen
```

### Option A: OpenAI Setup (Cloud AI)

```bash
# Set up OpenAI API key for cloud AI-powered suggestions
export OPENAI_API_KEY="sk-your-api-key-here"
export COMMITGEN_AI_PROVIDER="openai"  # Default
```

### Option B: Ollama Setup (Local AI)

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model (choose one)
ollama pull llama3.2:3b              # Recommended: Fast, lightweight
ollama pull qwen2.5-coder:7b         # Best for code
ollama pull codellama:7b             # Code-specialized
ollama pull deepseek-coder:6.7b      # High performance

# Configure commitgen
export COMMITGEN_AI_PROVIDER="ollama"
export COMMITGEN_AI_MODEL="llama3.2:3b"  # Or your chosen model
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
# Generate AI commit message (OpenAI or Ollama)
./bin/commitgen suggest --ai

# AI with verbose debug output
./bin/commitgen suggest --ai --verbose

# Force AI mode via environment variable
COMMITGEN_AI=1 ./bin/commitgen suggest

# Test different providers
COMMITGEN_AI_PROVIDER=openai ./bin/commitgen suggest --ai
COMMITGEN_AI_PROVIDER=ollama ./bin/commitgen suggest --ai
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

### Provider Comparison Workflow

```bash
# Test OpenAI vs Ollama on the same changes
git add .

# OpenAI (cloud)
COMMITGEN_AI_PROVIDER=openai ./bin/commitgen suggest --ai
# Example: "feat(auth): implement JWT token validation middleware"

# Ollama (local)  
COMMITGEN_AI_PROVIDER=ollama ./bin/commitgen suggest --ai
# Example: "feat(auth): add JWT authentication middleware"

# Choose your preferred provider for commit
git commit -m "$(COMMITGEN_AI_PROVIDER=ollama ./bin/commitgen suggest --ai)"
```

### Privacy-First Local AI

```bash
# Complete privacy with Ollama - no data leaves your machine
export COMMITGEN_AI_PROVIDER="ollama"
export COMMITGEN_AI_MODEL="qwen2.5-coder:7b"  # Best for code

# Normal workflow with full privacy
echo "export function validateUser() {}" >> src/auth.js
git add src/auth.js
git commit --no-edit  # Uses local AI, instant via cache
```
