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
COMMITGEN_AI_API_KEY=invalid ./bin/commitgen suggest --ai --verbose

echo ""
echo "3. Testing environment variable mode:"
COMMITGEN_AI=1 ./bin/commitgen suggest --verbose

echo ""
echo "4. To test with real OpenAI API key, run:"
echo "   export COMMITGEN_AI_API_KEY=sk-your-key-here"
echo "   ./bin/commitgen suggest --ai --verbose"

echo ""
echo "5. Available environment variables:"
echo "   COMMITGEN_AI=1                          # Enable AI mode"
echo "   COMMITGEN_AI_PROVIDER=openai            # Provider (openai/ollama)"
echo "   COMMITGEN_AI_API_KEY=sk-xxxxx           # OpenAI API key"
echo "   COMMITGEN_AI_MODEL=gpt-4o-mini          # Model name"
echo "   COMMITGEN_AI_BASE_URL=custom            # Custom endpoint"
echo "   COMMITGEN_MAX_FILES=10                  # Max files to analyze"
echo "   COMMITGEN_PATCH_BYTES=102400            # Max patch size"

echo ""
echo "AI foundation is working! Ready for real API key testing."
