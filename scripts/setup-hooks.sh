#!/bin/bash

# Setup Git Hooks for MCP Project
# This script configures local git hooks to ensure quality checks run before pushing.

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "🔧 Setting up Git hooks..."

# 1. Install pre-commit (if available)
if command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit hooks..."
    cd "$PROJECT_ROOT"
    pre-commit install
else
    echo "⚠️ 'pre-commit' not found. Skipping pre-commit installation."
    echo "   Please install it with: pip install pre-commit"
fi

# 2. Create pre-push hook
PRE_PUSH="$HOOKS_DIR/pre-push"
echo "Creating pre-push hook at $PRE_PUSH..."

cat > "$PRE_PUSH" << 'EOF'
#!/bin/bash
# MCP Pre-push Hook
# Runs Linters and Tests before allowing a push.

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Running pre-push quality checks...${NC}"

# 1. Linting (Pre-commit on all files)
echo -e "${YELLOW}🔍 Running Linters (pre-commit)...${NC}"
if command -v pre-commit &> /dev/null; then
    pre-commit run --all-files
else
    echo -e "${RED}⚠️ pre-commit not found. Skipping linting.${NC}"
fi

# 2. Go Tests (Jarvis)
echo -e "${YELLOW}🧪 Running Jarvis Tests (Go)...${NC}"
if [ -d "Jarvis" ]; then
    cd Jarvis
    if ! go test -v ./...; then
        echo -e "${RED}❌ Go tests failed.${NC}"
        exit 1
    fi
    cd ..
fi

# 3. Python Tests (MCPM)
echo -e "${YELLOW}🧪 Running MCPM Tests (Python)...${NC}"
if [ -d "mcpm_source" ]; then
    cd mcpm_source
    if command -v uv &> /dev/null; then
        # Use uv if available
        if ! uv run pytest; then
            echo -e "${RED}❌ Python tests failed.${NC}"
            exit 1
        fi
    else
        # Fallback
        if ! pytest; then
             echo -e "${RED}❌ Python tests failed.${NC}"
             exit 1
        fi
    fi
    cd ..
fi

echo -e "${GREEN}✅ All checks passed. Pushing allowed.${NC}"
exit 0
EOF

chmod +x "$PRE_PUSH"

echo "✅ Git hooks configured successfully."
