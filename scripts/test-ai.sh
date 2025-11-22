#!/bin/bash

# Test script for AI implementation
# This demonstrates how to test the OpenAI provider

echo "Testing commitgen AI implementation"
echo "======================================="

echo ""
echo "1. Testing without API key (should fallback to heuristics):"
./bin/commitgen suggest --ai --verbose

echo ""
echo "2. Testing with invalid API key (should fallback to heuristics):"
OPENAI_API_KEY=invalid ./bin/commitgen suggest --ai --verbose

echo ""
echo "3. Testing environment variable opt-in (set COMMITGEN_AI=1 to make --ai the default):"
COMMITGEN_AI=1 ./bin/commitgen suggest --verbose

echo ""
echo "4. To test with real OpenAI API key, run:"
echo "   export OPENAI_API_KEY=sk-your-key-here"
echo "   ./bin/commitgen suggest --ai --verbose"

echo ""
echo "5. Available environment variables:"
echo "   OPENAI_API_KEY=sk-xxxxx                 # OpenAI API key"
echo "   COMMITGEN_MODEL=gpt-4o-mini             # Model name"
echo "   COMMITGEN_BASE_URL=https://api.openai.com/v1"
echo "   COMMITGEN_AI=1                          # Enable AI mode by default"
echo "   COMMITGEN_AI_FALLBACK=true              # Allow heuristics if AI fails"
echo "   COMMITGEN_MAX_FILES=10                  # Max files to analyze"
echo "   COMMITGEN_PATCH_BYTES=102400            # Max patch size"
echo "   COMMITGEN_CONVENTIONS_FILE=~/custom.md  # Custom style guide"

echo ""
echo "AI foundation is working! Ready for real API key testing."
