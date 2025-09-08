#!/bin/bash
# commitgen setup script - makes AI setup simple for users

set -e

echo "🚀 Setting up commitgen with AI..."

# Check if commitgen is installed
if ! command -v commitgen &> /dev/null; then
    echo "❌ commitgen not found. Please install first:"
    echo "   brew install joaquinalmora/tap/commitgen"
    exit 1
fi

# Copy environment template
if [ ! -f ~/.env ]; then
    echo "📝 Creating environment template..."
    cp .env.example ~/.env
    echo "✅ Created ~/.env file"
else
    echo "📝 ~/.env already exists"
fi

# Check for API key
if grep -q "your-openai-api-key-here" ~/.env 2>/dev/null; then
    echo ""
    echo "🔑 NEXT STEP: Add your OpenAI API key"
    echo "   1. Get API key: https://platform.openai.com/api-keys"
    echo "   2. Edit: nano ~/.env"
    echo "   3. Replace: your-openai-api-key-here"
    echo ""
    read -p "Press Enter when you've added your API key..."
fi

# Install shell integration
echo "🔧 Installing shell integration..."
commitgen install-shell

# Test AI
echo "🤖 Testing AI connection..."
if commitgen suggest --ai >/dev/null 2>&1; then
    echo "✅ AI is working!"
else
    echo "⚠️  AI test failed. Check your API key in ~/.env"
fi

echo ""
echo "🎉 Setup complete! Try typing: git commit -m \""
echo "   AI suggestions will appear automatically!"
