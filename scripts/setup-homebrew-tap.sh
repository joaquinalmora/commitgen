#!/bin/bash
# Setup script for creating Homebrew tap repository
# Run this once to set up automated Homebrew distribution

set -e

echo "🍺 Setting up Homebrew tap for commitgen..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI (gh) is required but not installed.${NC}"
    echo -e "${YELLOW}Install with: brew install gh${NC}"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${RED}❌ Please authenticate with GitHub CLI first:${NC}"
    echo -e "${YELLOW}Run: gh auth login${NC}"
    exit 1
fi

echo -e "${GREEN}✅ GitHub CLI is installed and authenticated${NC}"

# Create the homebrew-tap repository
echo "📦 Creating joaquinalmora/homebrew-tap repository..."

if gh repo view joaquinalmora/homebrew-tap &> /dev/null; then
    echo -e "${YELLOW}⚠️  Repository joaquinalmora/homebrew-tap already exists${NC}"
else
    gh repo create joaquinalmora/homebrew-tap \
        --public \
        --description "Homebrew formulae for joaquinalmora's tools" \
        --clone=false
    echo -e "${GREEN}✅ Created joaquinalmora/homebrew-tap repository${NC}"
fi

# Create a GitHub token for Homebrew tap updates
echo ""
echo -e "${YELLOW}🔑 Next steps:${NC}"
echo "1. Create a GitHub Personal Access Token with 'repo' permissions:"
echo "   https://github.com/settings/tokens/new"
echo ""
echo "2. Add it as a repository secret named 'HOMEBREW_TAP_GITHUB_TOKEN':"
echo "   https://github.com/joaquinalmora/commitgen/settings/secrets/actions"
echo ""
echo "3. The token should have these permissions:"
echo "   - Contents: Write (to update the tap repository)"
echo "   - Metadata: Read"
echo ""
echo -e "${GREEN}✅ Homebrew tap setup script completed!${NC}"
echo ""
echo "Once you've added the token, you can create your first release with:"
echo -e "${YELLOW}git tag v0.1.0 && git push origin v0.1.0${NC}"
