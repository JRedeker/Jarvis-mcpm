#!/bin/bash

# Jarvis Setup Script
# Builds Jarvis and outputs the configuration needed for MCP clients.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
JARVIS_DIR="$PROJECT_ROOT/Jarvis"
BINARY_PATH="$JARVIS_DIR/jarvis"

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${CYAN}🤖 Setting up Jarvis...${NC}"

# 1. Build Jarvis
if [ ! -d "$JARVIS_DIR" ]; then
    echo "Error: Jarvis directory not found at $JARVIS_DIR"
    exit 1
fi

echo "Building binary..."
cd "$JARVIS_DIR"
go build -o jarvis .
echo -e "${GREEN}✅ Build successful.${NC}"

# 2. Verify Path
ABS_BINARY_PATH=$(realpath "$BINARY_PATH")

# 3. Output Configuration
echo ""
echo -e "${YELLOW}🎉 Jarvis is ready.${NC}"
echo "To connect your Agent, add this to your MCP configuration file:"
echo "(e.g., claude_desktop_config.json or Kilo Code settings)"
echo ""
echo -e "${GREEN}"
cat <<EOF
{
  "mcpServers": {
    "jarvis": {
      "command": "$ABS_BINARY_PATH",
      "args": []
    }
  }
}
EOF
echo -e "${NC}"
echo ""
echo "👉 Tip: You can copy the block above and paste it to your Agent."
echo "   Or tell your Agent: 'Configure yourself using the jarvis binary at $ABS_BINARY_PATH'"
