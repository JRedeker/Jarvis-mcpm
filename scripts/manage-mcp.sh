#!/bin/bash

# MCP Environment Management Script
# Handles the full lifecycle of the MCP infrastructure (Docker) and provides logging utilities.
# Updated for Streamable HTTP transport (MCP 2025-03-26 spec)

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

mkdir -p "$LOG_DIR"

function log() {
    echo "[$TIMESTAMP] $1" | tee -a "$LOG_DIR/management.log"
}

function start() {
    log "Starting MCP infrastructure..."
    cd "$PROJECT_ROOT"
    docker compose up -d
    log "Infrastructure started. Checking status..."
    docker compose ps
    echo ""
    echo -e "${CYAN}📦 Streamable HTTP endpoints available:${NC}"
    echo "  toolbox:  http://localhost:6276/mcp"
    echo "  memory:      http://localhost:6277/mcp"
    echo "  morph:       http://localhost:6278/mcp"
}

function stop() {
    log "Stopping MCP infrastructure..."
    cd "$PROJECT_ROOT"
    docker compose down
    log "Infrastructure stopped."
}

function restart() {
    log "Restarting MCP infrastructure..."
    stop
    sleep 2
    start
}

function rebuild() {
    log "Rebuilding mcpm-daemon container..."
    cd "$PROJECT_ROOT"
    docker compose build mcpm-daemon
    log "Restarting with new image..."
    docker compose up -d mcpm-daemon
    log "Rebuild complete."
    health
}

function status() {
    log "Checking MCP infrastructure status..."
    cd "$PROJECT_ROOT"
    docker compose ps
}

function logs() {
    local service="${2:-}"
    cd "$PROJECT_ROOT"
    if [ -n "$service" ]; then
        docker compose logs -f "$service"
    else
        docker compose logs -f
    fi
}

function health() {
    echo -e "${CYAN}🔍 Checking MCP service health...${NC}"
    echo ""

    # Check Docker containers
    echo -e "${CYAN}Docker Containers:${NC}"
    cd "$PROJECT_ROOT"
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
    echo ""

    # Check profile endpoints
    echo -e "${CYAN}Profile Endpoints (Streamable HTTP):${NC}"

    declare -A PROFILES
    PROFILES["toolbox"]=6276
    PROFILES["memory"]=6277
    PROFILES["morph"]=6278

    for profile in "${!PROFILES[@]}"; do
        port="${PROFILES[$profile]}"
        # Try to connect to the HTTP endpoint
        if curl -sf "http://localhost:$port/mcp" -o /dev/null -m 2 2>/dev/null || \
           curl -sf "http://localhost:$port/health" -o /dev/null -m 2 2>/dev/null; then
            echo -e "  ${GREEN}✓${NC} $profile (port $port): healthy"
        else
            echo -e "  ${RED}✗${NC} $profile (port $port): not responding"
        fi
    done
    echo ""

    # Check Jarvis Binary
    echo -e "${CYAN}Jarvis Binary:${NC}"
    if [ -f "$PROJECT_ROOT/Jarvis/jarvis" ]; then
        echo -e "  ${GREEN}✓${NC} Binary exists: $PROJECT_ROOT/Jarvis/jarvis"
        echo -e "  Version: $($PROJECT_ROOT/Jarvis/jarvis -help 2>&1 | head -1 || echo 'unknown')"
    else
        echo -e "  ${YELLOW}○${NC} Binary not built. Run: ./scripts/setup-jarvis.sh"
    fi
    echo ""

    # Check Supervisor Status (MCPM Daemon)
    echo -e "${CYAN}Daemon Process Status (Supervisor):${NC}"
    if docker ps --format '{{.Names}}' | grep -q "^mcpm-daemon$"; then
        # Capture status, filter out informational "No token" warnings
        status_output=$(docker exec mcpm-daemon supervisorctl status 2>&1 | grep -v "No token data found" || true)

        if [ -n "$status_output" ]; then
            echo "$status_output" | while read -r line; do
                if echo "$line" | grep -q "RUNNING"; then
                    echo -e "  ${GREEN}✓${NC} $line"
                else
                    echo -e "  ${RED}✗${NC} $line"
                fi
            done
        else
            echo -e "  ${YELLOW}⚠ No status output or supervisor not reachable${NC}"
        fi
    else
         echo -e "  ${YELLOW}⚠ mcpm-daemon container is not running${NC}"
    fi
}

function test() {
    log "Running full test suite..."

    # Go Tests
    log "Testing Jarvis (Go)..."
    cd "$PROJECT_ROOT/Jarvis"
    go test -v ./... || { log "❌ Go tests failed"; exit 1; }

    # Python Tests (if available)
    if [ -d "$PROJECT_ROOT/mcpm_source" ] && [ -f "$PROJECT_ROOT/mcpm_source/pyproject.toml" ]; then
        log "Testing MCPM (Python)..."
        cd "$PROJECT_ROOT/mcpm_source"
        uv run pytest || { log "❌ Python tests failed"; exit 1; }
    else
        log "Skipping Python tests (mcpm_source not configured)"
    fi

    log "✅ All tests passed."
}

function help() {
    echo "MCP Infrastructure Management Script"
    echo ""
    echo "Usage: $0 COMMAND [options]"
    echo ""
    echo "Commands:"
    echo "  start      Start all Docker containers"
    echo "  stop       Stop all Docker containers"
    echo "  restart    Restart all Docker containers"
    echo "  rebuild    Rebuild mcpm-daemon and restart"
    echo "  status     Show container status"
    echo "  health     Check health of all services"
    echo "  logs       Follow container logs (optionally: logs <service>)"
    echo "  test       Run full test suite"
    echo ""
    echo "Examples:"
    echo "  $0 start           # Start infrastructure"
    echo "  $0 health          # Check all services"
    echo "  $0 logs mcpm-daemon # Follow daemon logs"
    echo "  $0 rebuild         # Rebuild after config changes"
}

case "${1:-}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    rebuild)
        rebuild
        ;;
    status)
        status
        ;;
    health)
        health
        ;;
    logs)
        logs "$@"
        ;;
    test)
        test
        ;;
    ""|help|--help|-h)
        help
        ;;
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        echo ""
        help
        exit 1
        ;;
esac
