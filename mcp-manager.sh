#!/usr/bin/env bash

# =============================================================================
# MCP Server Management Script - Production-Grade
# Purpose: Manage cipher-aggregator with robust safety, monitoring, and conflict resolution
# Author: Auto-generated and reviewed
# Version: 2.0
# =============================================================================

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MCP_DIR="${SCRIPT_DIR}"
CIPHER_CONFIG="$MCP_DIR/cipher.yml"
LOG_DIR="$MCP_DIR/logs"
PID_FILE="$MCP_DIR/cipher-aggregator.pid"
SSE_PORT=3020
SSE_HOST="127.0.0.1"
MAX_STARTUP_WAIT=30
LOG_MAX_SIZE_MB=100

# Enhanced startup and health check configuration
STARTUP_GRACE_PERIOD=5
SSE_CHECK_TIMEOUT=3
SSE_CHECK_RETRIES=5
SSE_CONNECTION_POOL=false
HEALTH_PORT=3021
HEALTH_ENABLED=true
MONITOR_CHECK_INTERVAL=30

# Logging configuration
MCP_LOG_FORMAT="${MCP_LOG_FORMAT:-text}"
MCP_LOG_LEVEL="${MCP_LOG_LEVEL:-INFO}"
DEBUG="${DEBUG:-false}"
QUIET="${QUIET:-false}"

# Critical aggregator mode settings - hardcoded for reliability
export MCP_SERVER_MODE=aggregator
export AGGREGATOR_CONFLICT_RESOLUTION=prefix

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[$(date +'%H:%M:%S') INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%H:%M:%S') SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%H:%M:%S') WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[$(date +'%H:%M:%S') ERROR]${NC} $1"
}

log_debug() {
    echo -e "${PURPLE}[$(date +'%H:%M:%S') DEBUG]${NC} $1" >&2
}

# Structured JSON logging functions
log_json() {
    local level="$1"
    shift
    local message="$*"

    if [[ "$MCP_LOG_FORMAT" == "json" ]]; then
        # Check if jq is available for JSON formatting
        if command -v jq >/dev/null 2>&1; then
            echo '{"timestamp":"'$(date -u +%Y-%m-%dT%H:%M:%SZ)'","level":"'$level'","message":"'$message'","component":"mcp-manager","version":"2.0"}' | jq -r .
        else
            # Fallback without jq formatting
            echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $level: $message"
        fi
    fi
}

log_structured() {
    local event_type="$1"
    shift
    local data="$*"

    if [[ "$MCP_LOG_FORMAT" == "json" ]]; then
        if command -v jq >/dev/null 2>&1; then
            echo '{"timestamp":"'$(date -u +%Y-%m-%dT%H:%M:%SZ)'","event":"'$event_type'","data":'$data',"component":"mcp-manager"}' | jq -r .
        else
            echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] EVENT: $event_type $data"
        fi
    fi
}

# Create log with timestamp
setup_logging() {
    mkdir -p "$LOG_DIR"
    LOG_FILE="$LOG_DIR/mcp-manager-$(date +%Y%m%d-%H%M%S).log"

    # Redirect both stdout and stderr to log file
    if [[ "$MCP_LOG_FORMAT" == "json" ]]; then
        log_structured "script_started" '{"version":"2.0","log_format":"json"}'
        exec > >(tee -a "$LOG_FILE" | while IFS= read -r line; do log_json "info" "$line"; done)
        exec 2>&1 | while IFS= read -r line; do log_json "error" "$line"; done
    else
        exec > >(tee -a "$LOG_FILE")
        exec 2>&1
    fi
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--debug)
                DEBUG=true
                shift
                ;;
            -q|--quiet)
                QUIET=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                COMMAND="$1"
                shift
                ;;
        esac
    done
    COMMAND="${COMMAND:-help}"
}

# Ensure required directories exist
ensure_directories() {
    mkdir -p "$LOG_DIR" "$MCP_DIR/servers" 2>/dev/null || true
    [[ -w "$LOG_DIR" ]] || {
        echo "Cannot write to log directory: $LOG_DIR" >&2
        exit 1
    }
}

# Detect cipher binary path
find_cipher_binary() {
    local candidate_paths=(
        "$(which cipher 2>/dev/null || true)"
        "$HOME/.npm-global/bin/cipher"
        "$HOME/.local/bin/cipher"
        "/usr/local/bin/cipher"
    )

    for path in "${candidate_paths[@]}"; do
        if [[ -n "$path" && -x "$path" ]]; then
            echo "$path"
            return 0
        fi
    done

    # Try npm global as fallback
    local npm_path="$HOME/.npm-global/bin/cipher"
    if [[ -x "$npm_path" ]]; then
        echo "$npm_path"
        return 0
    fi

    echo "Cipher binary not found in PATH" >&2
    return 1
}

# Check if port is available
is_port_available() {
    local port="$1"
    if command -v lsof >/dev/null 2>&1; then
        ! lsof -i :$port >/dev/null 2>&1
    elif command -v netstat >/dev/null 2>&1; then
        ! netstat -ln 2>/dev/null | grep -q ":$port "
    else
        log_warning "Cannot check port availability (lsof/netstat not available)"
        return 0  # Assume available if can't check
    fi
}

# Wait for port to be available
wait_for_port() {
    local port="$1"
    local max_wait="${2:-10}"
    local count=0

    while ! is_port_available "$port" && [[ $count -lt $max_wait ]]; do
        log_debug "Port $port in use, waiting... ($count/$max_wait)"
        sleep 1
        ((count++))
    done

    if [[ $count -eq $max_wait ]]; then
        return 1
    fi
    return 0
}

# Enhanced SSE server connectivity test with exponential backoff
test_sse_server() {
    local max_attempts="${SSE_CHECK_RETRIES:-5}"
    local timeout="${SSE_CHECK_TIMEOUT:-3}"
    local attempt=1
    local backoff=1

    while [[ $attempt -le $max_attempts ]]; do
        log_debug "Testing SSE server (attempt $attempt/$max_attempts)..."

        # Test basic connectivity first (faster)
        if timeout $timeout bash -c "exec 3<>/dev/tcp/$SSE_HOST/$SSE_PORT" 2>/dev/null; then
            log_debug "SSE port is reachable"
            return 0  # Port reachable = server is responding
            fi
        ((attempt++))
        if [[ $attempt -le $max_attempts ]]; then
            log_debug "SSE test failed, waiting ${backoff}s before retry..."
            sleep $backoff
            # Exponential backoff: 1s, 2s, 4s, 8s, 8s...
            backoff=$((backoff < 8 ? backoff * 2 : 8))
        fi
    done

    log_warning "SSE server not responding after $max_attempts attempts"
    return 1
}

# Check if cipher-aggregator is running
is_cipher_running() {
    if [[ ! -f "$PID_FILE" ]]; then
        return 1
    fi

    local pid
    pid=$(cat "$PID_FILE")

    if ! ps -p "$pid" >/dev/null 2>&1; then
        rm -f "$PID_FILE"
        return 1
    fi

    # Verify it's actually the cipher process
    if ! ps -p "$pid" -o cmd= 2>/dev/null | grep -q "cipher.*mcp"; then
        log_warning "PID $pid exists but is not cipher-aggregator"
        rm -f "$PID_FILE"
        return 1
    fi

    return 0
}

# Get cipher-aggregator PID
get_cipher_pid() {
    if [[ -f "$PID_FILE" ]]; then
        cat "$PID_FILE"
    fi
}

# Smart conflict detection - target only known MCP servers
detect_conflicts() {
    local conflicts=()

    # Known MCP server patterns to watch for
    local mcp_patterns=(
        "firecrawl-mcp"
        "@morph-llm/morph"
        "@allpepper/memory-bank"
        "@upstash/context7"
        "@sveltejs/mcp"
        "mcp-server-filesystem"
        "@executeautomation/playwright"
        "mcp-server-memory-bank"
    )

    # Get all Node processes
    while IFS= read -r line; do
        local pid=$(echo "$line" | awk '{print $2}')
        local cmd=$(echo "$line" | awk '{for(i=11;i<=NF;i++) printf "%s ", $i; print ""}')

        for pattern in "${mcp_patterns[@]}"; do
            if [[ "$cmd" == *"$pattern"* ]]; then
                conflicts+=("$pid")
                log_debug "Found conflict: PID $pid ($pattern)"
                break
            fi
        done
    done < <(ps aux | grep -E "node.*npm|node.*mcp" | grep -v grep | grep -v "cipher" || true)

    printf '%s\n' "${conflicts[@]}"
}

# Safe conflict cleanup
kill_conflicts() {
    local conflicts=($(detect_conflicts))

    if [[ ${#conflicts[@]} -eq 0 ]]; then
        log_success "No conflicting MCP servers detected"
        return 0
    fi

    log_warning "Found ${#conflicts[@]} conflicting MCP servers: ${conflicts[*]}"
    log_info "Attempting graceful termination..."

    # Graceful termination
    for pid in "${conflicts[@]}"; do
        kill -TERM "$pid" 2>/dev/null || true
    done

    # Wait and check
    sleep 3
    local remaining=($(detect_conflicts))

    if [[ ${#remaining[@]} -gt 0 ]]; then
        log_warning "Forcing remaining conflicts..."
        for pid in "${remaining[@]}"; do
            kill -KILL "$pid" 2>/dev/null || true
        done
        sleep 1
    fi

    # Final verification
    local final_check=($(detect_conflicts))
    if [[ ${#final_check[@]} -eq 0 ]]; then
        log_success "All conflicting MCP servers eliminated"
    else
        log_error "Failed to eliminate all conflicts: ${final_check[*]}"
        return 1
    fi
}

# Clean up old log files
cleanup_logs() {
    if [[ ! -d "$LOG_DIR" ]]; then
        return 0
    fi

    local log_size_mb
    log_size_mb=$(du -sm "$LOG_DIR" 2>/dev/null | cut -f1 || echo "0")

    if [[ $log_size_mb -gt $LOG_MAX_SIZE_MB ]]; then
        log_warning "Log directory size ($log_size_mb MB) exceeds limit ($LOG_MAX_SIZE_MB MB)"
        log_info "Cleaning up old log files..."

        # Keep only recent logs (last 7 days)
        find "$LOG_DIR" -name "*.log" -mtime +7 -delete 2>/dev/null || true
        find "$LOG_DIR" -name "*.log" -size +50M -delete 2>/dev/null || true

        local new_size_mb
        new_size_mb=$(du -sm "$LOG_DIR" 2>/dev/null | cut -f1 || echo "0")
        log_success "Log cleanup completed. New size: ${new_size_mb} MB"
    fi
}

# Get cipher version
get_cipher_version() {
    local cipher_path
    cipher_path=$(find_cipher_binary) || return 1

    if command -v "$cipher_path" >/dev/null 2>&1; then
        "$cipher_path" --version 2>/dev/null | head -1 || echo "unknown"
    else
        echo "binary not found"
    fi
}

# Start cipher-aggregator
start_cipher() {
    log_info "Starting cipher-aggregator..."

    # Pre-flight checks
    if [[ ! -f "$CIPHER_CONFIG" ]]; then
        log_error "Cipher configuration not found: $CIPHER_CONFIG"
        return 1
    fi

    # Check if already running
    if is_cipher_running; then
        local pid
        pid=$(get_cipher_pid)
        log_warning "Cipher-aggregator is already running (PID: $pid)"
        return 0
    fi

    # Check port availability
    if ! is_port_available "$SSE_PORT"; then
        log_error "Port $SSE_PORT is already in use"
        return 1
    fi

    # Clear conflicts
    kill_conflicts || true

    # Ensure we're in the right directory
    cd "$MCP_DIR"

    # Find cipher binary
    local cipher_path
    cipher_path=$(find_cipher_binary) || {
        log_error "Cannot find cipher binary"
        return 1
    }

    log_info "Using cipher binary: $cipher_path"
    log_info "Cipher version: $(get_cipher_version)"

    # Load environment variables from .env file
    if [[ -f "$MCP_DIR/.env" ]]; then
        log_info "Loading environment variables from .env..."
        set -a  # Automatically export all variables
        source "$MCP_DIR/.env"
        set +a
        # Use safe parameter expansion to avoid unbound variable errors under 'set -u'
        log_success "Loaded .env file (OPENAI_API_KEY present: $([ -n \"${OPENAI_API_KEY:-}\" ] && echo 'yes' || echo 'no'))"
    else
        log_warning ".env file not found at $MCP_DIR/.env"
    fi

    # Ensure critical environment variables are set
    export MCP_SERVER_MODE=aggregator
    export AGGREGATOR_CONFLICT_RESOLUTION=prefix
    log_info "Aggregator mode: MCP_SERVER_MODE=$MCP_SERVER_MODE, AGGREGATOR_CONFLICT_RESOLUTION=$AGGREGATOR_CONFLICT_RESOLUTION"

    # Start cipher-aggregator in background
    log_info "Starting cipher-aggregator on $SSE_HOST:$SSE_PORT..."

    # Start with proper backgrounding
    nohup "$cipher_path" \
        --mode mcp \
        --mcp-transport-type sse \
        --mcp-host "$SSE_HOST" \
        --mcp-port "$SSE_PORT" \
        -a "$CIPHER_CONFIG" \
        > "$LOG_DIR/cipher-aggregator-$(date +%Y%m%d-%H%M%S).log" 2>&1 &

    local cipher_pid=$!
    echo "$cipher_pid" > "$PID_FILE"

    log_info "Started cipher-aggregator (PID: $cipher_pid), waiting for startup..."

    # Give cipher time to initialize before starting health checks
    local grace_period="${STARTUP_GRACE_PERIOD:-5}"
    log_info "Waiting ${grace_period}s for cipher initialization..."
    sleep $grace_period

    # Wait for startup with enhanced timeout and retry logic
    local count=0
    local startup_success=false

    while [[ $count -lt $MAX_STARTUP_WAIT ]]; do
        sleep 2
        ((count+=2))

        # Check if process is still running
        if ! ps -p "$cipher_pid" >/dev/null 2>&1; then
            log_error "Cipher-aggregator process died during startup"
            rm -f "$PID_FILE"
            cat "$LOG_DIR"/cipher-aggregator-*.log 2>/dev/null | tail -20
            return 1
        fi

        # Check if SSE server is responding with enhanced test
        if test_sse_server; then
            startup_success=true
            break
        fi

        log_debug "Waiting for startup... ($count/${MAX_STARTUP_WAIT}s)"
    done

    if [[ "$startup_success" == "true" ]]; then
        log_success "Cipher-aggregator started successfully (PID: $cipher_pid)"
        log_success "SSE server ready at http://$SSE_HOST:$SSE_PORT/sse"
        return 0
    else
        log_error "Cipher-aggregator failed to start within $MAX_STARTUP_WAIT seconds"
        log_error "Check logs: $LOG_DIR/cipher-aggregator-*.log"
        cat "$LOG_DIR"/cipher-aggregator-*.log 2>/dev/null | tail -20
        stop_cipher  # Clean up on failure
        return 1
    fi
}

# Stop cipher-aggregator
stop_cipher() {
    log_info "Stopping cipher-aggregator..."

    if ! is_cipher_running; then
        log_warning "Cipher-aggregator is not running"
        rm -f "$PID_FILE"  # Clean up stale PID file
        return 0
    fi

    local pid
    pid=$(get_cipher_pid)

    # Graceful shutdown
    log_info "Sending TERM signal to PID $pid..."
    kill -TERM "$pid" 2>/dev/null || true

    # Wait for graceful shutdown
    local count=0
    while is_cipher_running && [[ $count -lt 10 ]]; do
        sleep 1
        ((count++))
    done

    # Force kill if still running
    if is_cipher_running; then
        log_warning "Graceful shutdown failed, forcing..."
        kill -KILL "$pid" 2>/dev/null || true
        sleep 2
    fi

    rm -f "$PID_FILE"

    # Verify it's really stopped
    if ! is_cipher_running; then
        log_success "Cipher-aggregator stopped successfully"
    else
        log_error "Failed to stop cipher-aggregator (PID $pid may be stuck)"
        return 1
    fi
}

# Restart cipher-aggregator
restart_cipher() {
    log_info "Restarting cipher-aggregator..."
    stop_cipher
    sleep 2
    start_cipher
}

# Show detailed status
status_cipher() {
    if is_cipher_running; then
        local pid
        pid=$(get_cipher_pid)

        # Get detailed process info
        local cmd
        cmd=$(ps -p "$pid" -o cmd= 2>/dev/null || echo "unknown")

        local start_time
        start_time=$(ps -p "$pid" -o lstart= 2>/dev/null || echo "unknown")

        local memory_mb
        memory_mb=$(ps -p "$pid" -o %mem= 2>/dev/null | head -1 || echo "unknown")

        log_success "✅ Cipher-aggregator is running (PID: $pid)"
        echo -e "   ${BLUE}Command:${NC} $cmd"
        echo -e "   ${BLUE}Started:${NC} $start_time"
        echo -e "   ${BLUE}Memory:${NC} ${memory_mb}%"
        echo -e "   ${BLUE}Endpoint:${NC} http://$SSE_HOST:$SSE_PORT/sse"

        # Test SSE server
        if test_sse_server; then
            log_success "✅ SSE server responding"
        else
            log_warning "⚠️  SSE server not responding"
        fi

        # Check connections
        if command -v lsof >/dev/null 2>&1; then
            local connections
            connections=$(lsof -i :$SSE_PORT 2>/dev/null | grep -c ESTABLISHED || echo "0")
            echo -e "   ${BLUE}Active Connections:${NC} $connections"
        fi

        return 0
    else
        log_error "❌ Cipher-aggregator is not running"

        # Show recent logs if available
        if ls "$LOG_DIR"/cipher-aggregator-*.log >/dev/null 2>&1; then
            echo -e "\n${YELLOW}Recent errors from logs:${NC}"
            grep -E "(ERROR|error|Error)" "$LOG_DIR"/cipher-aggregator-*.log 2>/dev/null | tail -3 || echo "No recent errors found"
        fi

        return 1
    fi
}

# Comprehensive cleanup
cleanup() {
    log_info "Performing comprehensive MCP cleanup..."

    # Stop cipher-aggregator
    stop_cipher

    # Kill conflicts
    kill_conflicts

    # Clean up logs
    cleanup_logs

    log_success "Cleanup completed"
}

# Monitor mode - watch for conflicts and auto-resolve
monitor() {
    log_info "Starting monitor mode (Ctrl+C to stop)..."
    log_info "Will check every 30 seconds and auto-recover issues"

    local check_count=0

    while true; do
        ((check_count++))

        if is_cipher_running; then
            # Check SSE server health
            if ! test_sse_server; then
                log_warning "$(date): SSE server not responding, restarting..."
                restart_cipher
            fi

            # Check for conflicts
            local conflicts=($(detect_conflicts))
            if [[ ${#conflicts[@]} -gt 0 ]]; then
                log_warning "$(date): Found ${#conflicts[@]} conflicts, cleaning up..."
                kill_conflicts
            fi

        else
            log_warning "$(date): Cipher-aggregator not running, starting..."
            start_cipher
        fi

        # Show heartbeat every 5 minutes
        if [[ $((check_count % 10)) -eq 0 ]]; then
            log_success "Monitor heartbeat: $check_count checks completed"
        fi

        sleep 30  # Check every 30 seconds
    done
}

# List available MCP tools (placeholder)
list_tools() {
    if ! is_cipher_running; then
        log_error "Cipher-aggregator is not running"
        return 1
    fi

    log_info "Available MCP servers via cipher-aggregator:"
    if [[ -f "$CIPHER_CONFIG" ]]; then
        # Parse cipher.yml for server names
        grep -E "^  [a-z-]+:" "$CIPHER_CONFIG" 2>/dev/null | \
            sed 's/^  //; s/:$//' | \
            sort | \
            while read -r server; do
                echo "  - $server"
            done
    fi
}

# Show detailed configuration
show_config() {
    log_info "MCP Management Configuration:"

    echo -e "\n${BLUE}Paths:${NC}"
    echo "  Script Directory: $SCRIPT_DIR"
    echo "  MCP Directory: $MCP_DIR"
    echo "  Config File: $CIPHER_CONFIG"
    echo "  Log Directory: $LOG_DIR"
    echo "  PID File: $PID_FILE"

    echo -e "\n${BLUE}Network:${NC}"
    echo "  SSE Host: $SSE_HOST"
    echo "  SSE Port: $SSE_PORT"
    echo "  SSE URL: http://$SSE_HOST:$SSE_PORT/sse"

    echo -e "\n${BLUE}System:${NC}"
    echo "  Cipher Binary: $(find_cipher_binary 2>/dev/null || echo 'Not found')"
    echo "  Cipher Version: $(get_cipher_version)"
    echo "  Max Startup Wait: ${MAX_STARTUP_WAIT}s"
    echo "  Log Size Limit: ${LOG_MAX_SIZE_MB}MB"

    if [[ -f "$CIPHER_CONFIG" ]]; then
        local server_count
        server_count=$(grep -E "^  [a-z-]+:" "$CIPHER_CONFIG" 2>/dev/null | wc -l || echo "0")
        echo "  Configured Servers: $server_count"
    fi

    echo -e "\n${BLUE}Status:${NC}"
    if is_cipher_running; then
        echo "  Cipher-Aggregator: Running"
    else
        echo "  Cipher-Aggregator: Stopped"
    fi

    local conflicts_count
    conflicts_count=$((${#$(detect_conflicts)} - 0))
    echo "  Conflicts Detected: $conflicts_count"
}

# Health check function
health_check() {
    local health_score=0
    local max_score=100

    echo -e "${BLUE}=== MCP System Health Check ===${NC}"

    # Check cipher running
    if is_cipher_running; then
        echo -e "✅ Cipher-Aggregator: Running"
        ((health_score+=25))
    else
        echo -e "❌ Cipher-Aggregator: Not running"
    fi

    # Check SSE server
    if test_sse_server; then
        echo -e "✅ SSE Server: Responding"
        ((health_score+=25))
    else
        echo -e "❌ SSE Server: Not responding"
    fi

    # Check port availability
    if is_port_available "$SSE_PORT"; then
        if ! is_cipher_running; then
            echo -e "✅ Port $SSE_PORT: Available"
            ((health_score+=25))
        else
            echo -e "✅ Port $SSE_PORT: In use (by cipher)"
            ((health_score+=25))
        fi
    else
        echo -e "❌ Port $SSE_PORT: Occupied by other process"
    fi

    # Check conflicts
    local conflicts
    conflicts=($(detect_conflicts))
    local conflicts_count=${#conflicts[@]}
    if [[ $conflicts_count -eq 0 ]]; then
        echo -e "✅ No Conflicts: Clean environment"
        ((health_score+=25))
    else
        echo -e "❌ Conflicts Found: $conflicts_count MCP servers"
    fi

    echo -e "\n${BLUE}Health Score: $health_score/$max_score${NC}"

    if [[ $health_score -eq $max_score ]]; then
        echo -e "${GREEN}✅ System is healthy!${NC}"
    elif [[ $health_score -ge 75 ]]; then
        echo -e "${YELLOW}⚠️  System is mostly healthy${NC}"
    else
        echo -e "${RED}❌ System has issues that need attention${NC}"
    fi
}

# Show help
show_help() {
    cat << EOF
MCP Server Management Script v2.0
Usage: $0 <command> [options]

Commands:
  start         Start cipher-aggregator
  stop          Stop cipher-aggregator
  restart       Restart cipher-aggregator
  status        Show detailed status
  health        Run system health check
  cleanup       Kill conflicts and clean logs
  monitor       Auto-monitor and recover
  list          List available MCP tools
  config        Show configuration

Options:
  -d, --debug   Enable debug output
  -q, --quiet   Suppress non-error output
  -h, --help    Show this help

Examples:
  $0 start                    # Start cipher-aggregator
  $0 health                   # Check system health
  $0 monitor                  # Auto-recovery mode
  $0 status                   # Detailed status
  $0 cleanup                  # Clean conflicts

For more information, see $MCP_DIR/README.md
EOF
}

# Main script logic
main() {
    # Setup logging first
    setup_logging

    # Parse command line arguments
    parse_args "$@"

    # Ensure directories exist
    ensure_directories

    # Execute command
    case "$COMMAND" in
        start)
            start_cipher
            ;;
        stop)
            stop_cipher
            ;;
        restart)
            restart_cipher
            ;;
        status)
            status_cipher
            ;;
        health)
            health_check
            ;;
        cleanup)
            cleanup
            ;;
        monitor)
            monitor
            ;;
        list)
            list_tools
            ;;
        config)
            show_config
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
