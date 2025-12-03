#!/bin/bash
set -e

# MCPM Daemon Entrypoint
# Starts all configured profiles as Streamable HTTP services using supervisor
# Updated from SSE to HTTP transport (MCP 2025-03-26 spec)

echo "=== MCPM Daemon Starting ==="

# Export all API keys to environment (they're passed via docker-compose)
# These will be inherited by supervisor-managed processes
export FIRECRAWL_API_KEY="${FIRECRAWL_API_KEY:-}"
export KAGI_API_KEY="${KAGI_API_KEY:-}"
export BRAVE_API_KEY="${BRAVE_API_KEY:-}"
export QDRANT_API_KEY="${QDRANT_API_KEY:-}"
export OPENROUTER_API_KEY="${OPENROUTER_API_KEY:-}"
export MORPH_API_KEY="${MORPH_API_KEY:-}"
export OPENAI_API_KEY="${OPENAI_API_KEY:-}"
export QDRANT_URL="${QDRANT_URL:-http://localhost:6333}"

echo "Environment configured with API keys"

# Profile to port mapping
declare -A PROFILE_PORTS=(
    ["p-pokeedge"]=6276
    ["memory"]=6277
    ["morph"]=6278
    ["qdrant"]=6279
    ["p-new"]=6280
)

# Generate supervisor config for each enabled profile
generate_supervisor_config() {
    cat > /etc/supervisor/conf.d/mcpm-profiles.conf << 'SUPERVISOR_HEADER'
[supervisord]
nodaemon=true
logfile=/var/log/mcpm/supervisord.log
pidfile=/var/run/supervisord.pid
loglevel=info

[unix_http_server]
file=/var/run/supervisor.sock
chmod=0700

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///var/run/supervisor.sock

SUPERVISOR_HEADER

    # Check which profiles to run (from MCPM_PROFILES env var, or default to all)
    if [ -n "$MCPM_PROFILES" ]; then
        IFS=',' read -ra PROFILES <<< "$MCPM_PROFILES"
    else
        PROFILES=("p-pokeedge" "memory" "morph")
    fi

    for profile in "${PROFILES[@]}"; do
        port="${PROFILE_PORTS[$profile]}"
        if [ -n "$port" ]; then
            echo "Configuring profile: $profile on port $port"
            cat >> /etc/supervisor/conf.d/mcpm-profiles.conf << EOF

[program:mcpm-$profile]
command=mcpm profile run --http --host 0.0.0.0 --port $port $profile
autostart=true
autorestart=true
stderr_logfile=/var/log/mcpm/$profile.err.log
stdout_logfile=/var/log/mcpm/$profile.out.log
environment=HOME="/root",PATH="/root/.local/bin:/usr/local/bin:/usr/bin:/bin"
EOF
        fi
    done
}

# Wait for MCPM config to be available (mounted volume)
wait_for_config() {
    local max_wait=30
    local waited=0

    while [ ! -f /root/.config/mcpm/servers.json ]; do
        if [ $waited -ge $max_wait ]; then
            echo "ERROR: MCPM config not found after ${max_wait}s"
            echo "Make sure to mount ~/.config/mcpm to /root/.config/mcpm"
            exit 1
        fi
        echo "Waiting for MCPM config... (${waited}s)"
        sleep 1
        ((waited++))
    done

    echo "MCPM config found!"
}

# Main
wait_for_config
generate_supervisor_config

echo "=== Starting Supervisor ==="
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/mcpm-profiles.conf
