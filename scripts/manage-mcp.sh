#!/bin/bash

# MCP Environment Management Script
# Handles the full lifecycle of the MCP infrastructure (Docker) and provides logging utilities.

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

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

function status() {
    log "Checking MCP infrastructure status..."
    cd "$PROJECT_ROOT"
    docker compose ps
}

function logs() {
    cd "$PROJECT_ROOT"
    docker compose logs -f
}

function test() {
    log "Running full test suite..."

    # Go Tests
    log "Testing Jarvis (Go)..."
    cd "$PROJECT_ROOT/Jarvis"
    go test -v ./... || { log "❌ Go tests failed"; exit 1; }

    # Python Tests
    log "Testing MCPM (Python)..."
    cd "$PROJECT_ROOT/mcpm_source"
    uv run pytest || { log "❌ Python tests failed"; exit 1; }

    log "✅ All tests passed."
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs
        ;;
    test)
        test
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs|test}"
        exit 1
        ;;
esac
