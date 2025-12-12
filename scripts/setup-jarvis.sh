#!/bin/bash

# Jarvis Setup Script
# Builds Jarvis and outputs the configuration needed for MCP clients.
# Supports both stdio (default) and Streamable HTTP transport modes.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
JARVIS_DIR="$PROJECT_ROOT/Jarvis"
BINARY_PATH="$JARVIS_DIR/jarvis"

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
HTTP_MODE=false
HTTP_PORT="6275"
HTTP_HOST="127.0.0.1"
AUTO_CONFIG=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --http)
            HTTP_MODE=true
            shift
            ;;
        --port)
            HTTP_PORT="$2"
            shift 2
            ;;
        --host)
            HTTP_HOST="$2"
            shift 2
            ;;
        --auto-config)
            AUTO_CONFIG=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --http         Configure for Streamable HTTP mode instead of stdio"
            echo "  --port PORT    HTTP port (default: 6275)"
            echo "  --host HOST    HTTP host (default: 127.0.0.1)"
            echo "  --auto-config  Automatically update detected client configs"
            echo "  --help         Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${CYAN}🤖 Setting up Jarvis...${NC}"

# 1. Build Jarvis
if [ ! -d "$JARVIS_DIR" ]; then
    echo -e "${RED}Error: Jarvis directory not found at $JARVIS_DIR${NC}"
    exit 1
fi

echo "Building binary..."
cd "$JARVIS_DIR"
go build -o jarvis .
echo -e "${GREEN}✅ Build successful.${NC}"

# 2. Verify Path
ABS_BINARY_PATH=$(realpath "$BINARY_PATH")

# 3. Detect client config locations
declare -A CLIENT_CONFIGS
CLIENT_CONFIGS["Claude Code"]="$HOME/.claude.json"
CLIENT_CONFIGS["Claude Desktop"]="$HOME/.config/Claude/claude_desktop_config.json"
CLIENT_CONFIGS["VSCode Claude"]="$HOME/.config/Code/User/globalStorage/anthropic.claude-code/settings.json"

echo ""
echo -e "${CYAN}📍 Detected client config locations:${NC}"
for client in "${!CLIENT_CONFIGS[@]}"; do
    config_path="${CLIENT_CONFIGS[$client]}"
    if [ -f "$config_path" ]; then
        echo -e "  ${GREEN}✓${NC} $client: $config_path"
    else
        echo -e "  ${YELLOW}○${NC} $client: $config_path (not found)"
    fi
done

# 4. Generate configuration based on mode
echo ""
echo -e "${YELLOW}🎉 Jarvis is ready.${NC}"

if [ "$HTTP_MODE" = true ]; then
    HTTP_URL="http://${HTTP_HOST}:${HTTP_PORT}/mcp"
    echo ""
    echo -e "${CYAN}Transport: Streamable HTTP (MCP 2025-03-26 spec)${NC}"
    echo ""
    echo "To start Jarvis in HTTP mode:"
    echo -e "${GREEN}  $ABS_BINARY_PATH --http --host $HTTP_HOST --port $HTTP_PORT${NC}"
    echo ""
    echo "Add this to your MCP client configuration:"
    echo ""
    echo -e "${GREEN}"
    cat <<EOF
{
  "mcpServers": {
    "jarvis": {
      "url": "$HTTP_URL"
    }
  }
}
EOF
    echo -e "${NC}"
else
    echo ""
    echo -e "${CYAN}Transport: stdio (default for direct client connections)${NC}"
    echo ""
    echo "Add this to your MCP client configuration:"
    echo "(e.g., claude_desktop_config.json or ~/.claude.json)"
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
fi

# 5. Also show profile configurations for daemon
echo ""
echo -e "${CYAN}📦 Profile endpoints (from mcpm-daemon):${NC}"
echo "  p-pokeedge:  http://localhost:6276/mcp"
echo "  memory:      http://localhost:6277/mcp"
echo "  morph:       http://localhost:6278/mcp"
echo "  qdrant:      http://localhost:6279/mcp"
echo "  p-new:       http://localhost:6280/mcp"
echo ""

# 6. Auto-config if requested
if [ "$AUTO_CONFIG" = true ]; then
    echo -e "${CYAN}🔧 Auto-configuring clients...${NC}"

    for client in "${!CLIENT_CONFIGS[@]}"; do
        config_path="${CLIENT_CONFIGS[$client]}"
        if [ -f "$config_path" ]; then
            echo -e "  Updating $client..."
            # Create backup
            cp "$config_path" "${config_path}.backup.$(date +%Y%m%d_%H%M%S)"

            # Use jq to update if available
            if command -v jq &> /dev/null; then
                if [ "$HTTP_MODE" = true ]; then
                    jq --arg url "$HTTP_URL" \
                       '.mcpServers.jarvis = {"url": $url}' \
                       "$config_path" > "${config_path}.tmp" && mv "${config_path}.tmp" "$config_path"
                else
                    jq --arg cmd "$ABS_BINARY_PATH" \
                       '.mcpServers.jarvis = {"command": $cmd, "args": []}' \
                       "$config_path" > "${config_path}.tmp" && mv "${config_path}.tmp" "$config_path"
                fi
                echo -e "    ${GREEN}✓ Updated (using jq)${NC}"

            # Fallback to Python if jq is missing
            elif command -v python3 &> /dev/null; then
                export HTTP_MODE
                export HTTP_URL
                export ABS_BINARY_PATH
                export CONFIG_PATH="$config_path"

                python3 -c '
import sys, json, os

config_path = os.environ["CONFIG_PATH"]
http_mode = os.environ.get("HTTP_MODE") == "true"
http_url = os.environ.get("HTTP_URL")
binary_path = os.environ.get("ABS_BINARY_PATH")

try:
    with open(config_path, "r") as f:
        data = json.load(f)

    if "mcpServers" not in data:
        data["mcpServers"] = {}

    if http_mode:
        data["mcpServers"]["jarvis"] = {"url": http_url}
    else:
        data["mcpServers"]["jarvis"] = {"command": binary_path, "args": []}

    with open(config_path, "w") as f:
        json.dump(data, f, indent=2)
except Exception as e:
    sys.exit(1)
'
                if [ $? -eq 0 ]; then
                    echo -e "    ${GREEN}✓ Updated (using python3)${NC}"
                else
                    echo -e "    ${RED}✗ Update failed (python3 error)${NC}"
                fi

            else
                echo -e "    ${YELLOW}⚠ jq and python3 not installed - manual update required${NC}"
            fi
        fi
    done
fi

echo ""
echo "👉 Tip: You can copy the config block above and paste it to your Agent."
echo "   Or tell your Agent: 'Configure yourself using the jarvis binary at $ABS_BINARY_PATH'"
echo ""
echo "📚 For HTTP mode: $0 --http"
echo "📚 For auto-config: $0 --auto-config"
