#!/bin/bash

# Setup script for git hooks
# Run this once after cloning the repository

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🔧 Setting up git hooks...${NC}"

# Get the root directory
ROOT_DIR=$(git rev-parse --show-toplevel)

# Configure git to use .githooks directory
git config core.hooksPath "$ROOT_DIR/.githooks"

echo -e "${GREEN}✅ Git hooks configured!${NC}"
echo ""

# Check if gitleaks is installed
if command -v gitleaks &> /dev/null; then
    echo -e "${GREEN}✅ gitleaks is installed ($(gitleaks version))${NC}"
else
    echo -e "${YELLOW}⚠️  gitleaks is not installed${NC}"
    echo ""
    echo "Install it to enable security scanning:"
    echo "  macOS:   brew install gitleaks"
    echo "  Linux:   sudo apt install gitleaks"
    echo "  Go:      go install github.com/gitleaks/gitleaks/v8@latest"
    echo ""
fi

echo ""
echo "Setup complete! The pre-commit hook will now scan for:"
echo "  • API keys, tokens, and passwords"
echo "  • Absolute paths with usernames"
echo "  • Database URLs with embedded credentials"
echo "  • Private keys and certificates"
