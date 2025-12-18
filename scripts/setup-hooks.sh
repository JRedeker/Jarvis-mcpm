#!/bin/bash

# Setup Git Hooks for MCP Project
# This script configures local git hooks to ensure quality checks run before pushing.
# Matches CI workflow to prevent push/CI mismatches.

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "Setting up Git hooks..."

# 1. Install pre-commit (if available)
if command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit hooks..."
    cd "$PROJECT_ROOT"
    pre-commit install
else
    echo "Warning: 'pre-commit' not found. Skipping pre-commit installation."
    echo "   Please install it with: pip install pre-commit"
fi

# 2. Create pre-push hook
PRE_PUSH="$HOOKS_DIR/pre-push"
echo "Creating pre-push hook at $PRE_PUSH..."

cat > "$PRE_PUSH" << 'EOF'
#!/bin/bash
# MCP Pre-push Hook
# Runs ALL CI checks locally before allowing a push.
# This prevents CI failures after push.

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running pre-push quality checks (mirrors CI)...${NC}"

# 1. Linting (Pre-commit on all files)
echo -e "${YELLOW}Running Linters (pre-commit)...${NC}"
if command -v pre-commit &> /dev/null; then
    pre-commit run --all-files
else
    echo -e "${RED}Warning: pre-commit not found. Skipping linting.${NC}"
fi

# 2. Go Tests (Jarvis) - matches CI go-tests job
echo -e "${YELLOW}Running Jarvis Tests (Go)...${NC}"
if [ -d "Jarvis" ]; then
    cd Jarvis
    if ! go test -v ./...; then
        echo -e "${RED}Go tests failed.${NC}"
        exit 1
    fi
    cd ..
fi

# 3. Shell Tests (bats) - matches CI shell-tests job
echo -e "${YELLOW}Running Shell Tests (bats)...${NC}"
if command -v bats &> /dev/null; then
    if [ -d "scripts/tests" ]; then
        if ! bats scripts/tests/*.bats; then
            echo -e "${RED}Bats tests failed.${NC}"
            exit 1
        fi
    fi
else
    echo -e "${RED}Warning: bats not found. Skipping shell tests.${NC}"
    echo -e "${RED}Install with: sudo apt install bats (Linux) or brew install bats-core (macOS)${NC}"
fi

# 4. Docker Build Check (optional, skipped if Docker unavailable)
echo -e "${YELLOW}Checking Docker build...${NC}"
if command -v docker &> /dev/null; then
    if [ -f "mcpm-daemon/Dockerfile" ]; then
        if ! docker build -q -t mcpm-daemon-test ./mcpm-daemon > /dev/null 2>&1; then
            echo -e "${RED}Warning: Docker build failed (non-blocking).${NC}"
        else
            echo -e "${GREEN}Docker build OK${NC}"
            docker rmi mcpm-daemon-test > /dev/null 2>&1 || true
        fi
    fi
else
    echo -e "${YELLOW}Docker not available. Skipping Docker build check.${NC}"
fi

echo -e "${GREEN}All checks passed. Pushing allowed.${NC}"
exit 0
EOF

chmod +x "$PRE_PUSH"

echo "Git hooks configured successfully."
echo ""
echo "The pre-push hook now runs:"
echo "  1. pre-commit (linting)"
echo "  2. Go tests"
echo "  3. Bats shell tests"
echo "  4. Docker build check (optional)"
echo ""
echo "This mirrors the CI workflow to prevent failures after push."
