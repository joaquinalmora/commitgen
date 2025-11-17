#!/bin/bash
# commitgen automated setup script

set -e

echo "Setting up commitgen with AI..."

# Check if commitgen is installed
if ! command -v commitgen &> /dev/null; then
    echo "commitgen not found. Installing via go install..."
    if ! command -v go &> /dev/null; then
        echo "Go not found. Please install Go first or use Homebrew:"
        echo "   brew tap joaquinalmora/tap && brew install commitgen"
        exit 1
    fi
    go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest
    echo "commitgen installed"
fi

# Ask user preference for configuration method
echo ""
echo "Choose configuration method:"
echo "1. Interactive setup (recommended) - creates commitgen.yaml"
echo "2. Environment variables - creates .env file"
read -p "Enter choice [1-2] (default: 1): " config_choice

if [[ "$config_choice" == "2" ]]; then
    # Environment variables setup
    echo "Setting up environment variables..."
    commitgen env-example
    if [[ -f .env.example ]]; then
        cp .env.example ~/.env
        echo "Created ~/.env file"
        echo ""
        echo "NEXT STEP: Add your OpenAI API key"
        echo "   1. Get API key: https://platform.openai.com/api-keys"
        echo "   2. Edit: nano ~/.env"
        echo "   3. Replace: your-openai-api-key-here"
        echo ""
        read -p "Press Enter when you've added your API key..."
    fi
else
    # Interactive setup (default)
    echo "Starting interactive configuration..."
    commitgen init
fi

# Install integrations
echo "Installing git hooks..."
if git rev-parse --git-dir >/dev/null 2>&1; then
    commitgen install-hook
    echo "Git hooks installed"
else
    echo "Not in a git repository. Git hooks will be installed per-repo later."
fi

echo "Installing shell integration..."
commitgen install-shell
echo "Shell integration installed"

# Test the setup
echo "Testing configuration..."
if commitgen doctor >/dev/null 2>&1; then
    echo "All checks passed!"
else
    echo "Some checks failed. Run 'commitgen doctor' for details."
fi

echo ""
echo "Setup complete!"
echo ""
echo "Next steps:"
echo "1. Restart your terminal or run: source ~/.zshrc"
echo "2. In a git repo, try: git add . && git commit -m \""
echo "3. AI suggestions will appear automatically!"
echo ""
echo "Troubleshooting: run 'commitgen doctor' for diagnostics"
