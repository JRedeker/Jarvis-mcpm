#!/usr/bin/env bash
# Run Loki and OpenTelemetry Collector (contrib) locally without Docker.
# Usage:
#   scripts/run-telemetry-local.sh start|stop|restart|status
#
# Starts:
#   - Loki on http://localhost:3100 (config: docs/loki-local-config.yaml)
#   - OTEL Collector (gRPC OTLP) on 0.0.0.0:4317 (config: docs/otel-collector-local-config.yaml)
#
# Data and logs:
#   - Binaries: ./bin
#   - Loki data: ./data/loki
#   - Process logs: ./logs/loki-local.log, ./logs/otelcol-local.log
#   - PIDs: ./.pids/{loki.pid,otelcol.pid}

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

BIN="$REPO_ROOT/bin"
LOGS="$REPO_ROOT/logs"
PIDS="$REPO_ROOT/.pids"
DOCS="$REPO_ROOT/docs"
DATA="$REPO_ROOT/data"

mkdir -p "$BIN" "$LOGS" "$PIDS" "$DOCS" "$DATA" "$DATA/loki" "$DATA/loki/chunks" "$DATA/loki/index" "$DATA/loki/compactor" "$DATA/loki/boltdb-cache" "$DATA/loki/rules"

LOKI_BIN="$BIN/loki-linux-amd64"
OTEL_BIN="$BIN/otelcol-contrib"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || { echo "Missing required command: $1"; exit 1; }
}

download_loki() {
  if [ ! -x "$LOKI_BIN" ]; then
    echo "[bootstrap] Downloading Loki (linux-amd64, zip)..."
    need_cmd curl
    need_cmd unzip
    # Clean up any previous incorrect artifact
    rm -f "$BIN/loki-linux-amd64.gz" || true
    local zip_path="$BIN/loki-linux-amd64.zip"
    curl -sSL "https://github.com/grafana/loki/releases/latest/download/loki-linux-amd64.zip" -o "$zip_path"
    unzip -o "$zip_path" -d "$BIN" >/dev/null
    rm -f "$zip_path"
    chmod +x "$LOKI_BIN"
    echo "[bootstrap] Loki ready at $LOKI_BIN"
  fi
}

download_otel() {
  if [ ! -x "$OTEL_BIN" ]; then
    echo "[bootstrap] Downloading OpenTelemetry Collector (contrib, linux-amd64)..."
    need_cmd curl
    need_cmd tar
    need_cmd grep
    need_cmd sed
    # Try to resolve latest linux_amd64 asset name via GitHub API
    local api_url="https://api.github.com/repos/open-telemetry/opentelemetry-collector-releases/releases/latest"
    local asset_url
    asset_url="$(curl -sSL "$api_url" | grep -Eo 'https://[^"]*otelcol-contrib_[^"]*_linux_amd64\.tar\.gz' | head -n1 || true)"
    if [ -z "$asset_url" ]; then
      echo "[bootstrap] Could not resolve latest asset via API (rate limit?). Trying a direct name (may fail)..."
      asset_url="https://github.com/open-telemetry/opentelemetry-collector-releases/releases/latest/download/otelcol-contrib_linux_amd64.tar.gz"
    fi
    local tar_path="$BIN/otelcol-contrib_linux_amd64.tar.gz"
    curl -sSL "$asset_url" -o "$tar_path"
    tar -xz -C "$BIN" -f "$tar_path"
    rm -f "$tar_path"
    chmod +x "$OTEL_BIN"
    echo "[bootstrap] otelcol-contrib ready at $OTEL_BIN"
  fi
}

ensure_loki_config() {
  if [ ! -f "$DOCS/loki-local-config.yaml" ]; then
    echo "[bootstrap] Writing default Loki config to $DOCS/loki-local-config.yaml"
    cat > "$DOCS/loki-local-config.yaml" <<'YML'
# Minimal local Loki config for dev (no Docker)
# Stores data under ./data/loki; ensure you run from repo root.
# Matches Collector exporter endpoint: http://localhost:3100/loki/api/v1/push

server:
  http_listen_port: 3100
  grpc_listen_port: 0

common:
  path_prefix: ./data/loki
  storage:
    filesystem:
      chunks_directory: ./data/loki/chunks
      rules_directory: ./data/loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: ./data/loki/index
    cache_location: ./data/loki/boltdb-cache
    shared_store: filesystem
  filesystem:
    directory: ./data/loki/chunks

compactor:
  working_directory: ./data/loki/compactor
  shared_store: filesystem

limits_config:
  allow_structured_metadata: true
  ingestion_rate_mb: 8
  ingestion_burst_size_mb: 16

table_manager:
  retention_deletes_enabled: true
  retention_period: 168h  # 7 days for dev

ruler:
  storage:
    type: local
    local:
      directory: ./data/loki/rules
  rule_path: ./data/loki/rules-temp
  enable_api: true
YML
  fi
}

ensure_collector_config() {
  if [ ! -f "$DOCS/otel-collector-local-config.yaml" ]; then
    echo "[bootstrap] Writing default OTEL Collector config to $DOCS/otel-collector-local-config.yaml"
    cat > "$DOCS/otel-collector-local-config.yaml" <<'YML'
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch: {}

exporters:
  loki:
    endpoint: http://localhost:3100/loki/api/v1/push

extensions:
  health_check: {}

service:
  extensions: [health_check]
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [loki]
YML
  fi
}

start_loki() {
  if pgrep -f "$LOKI_BIN" >/dev/null 2>&1; then
    echo "[loki] Already running"
    return 0
  fi
  echo "[loki] Starting..."
  nohup "$LOKI_BIN" -config.file "$DOCS/loki-local-config.yaml" > "$LOGS/loki-local.log" 2>&1 &
  echo $! > "$PIDS/loki.pid"
  echo "[loki] PID $(cat "$PIDS/loki.pid")"
}

start_otel() {
  if pgrep -f "$OTEL_BIN --config $DOCS/otel-collector-local-config.yaml" >/dev/null 2>&1; then
    echo "[otelcol] Already running"
    return 0
  fi
  echo "[otelcol] Starting..."
  nohup "$OTEL_BIN" --config "$DOCS/otel-collector-local-config.yaml" > "$LOGS/otelcol-local.log" 2>&1 &
  echo $! > "$PIDS/otelcol.pid"
  echo "[otelcol] PID $(cat "$PIDS/otelcol.pid")"
}

wait_http() {
  local url="$1"
  local timeout="${2:-30}"
  local i=0
  until curl -s -o /dev/null -f "$url"; do
    sleep 1
    i=$((i+1))
    if [ "$i" -ge "$timeout" ]; then
      echo "[wait] Timeout waiting for $url"
      return 1
    fi
  done
  return 0
}

start_all() {
  download_loki
  download_otel
  ensure_loki_config
  ensure_collector_config

  start_loki
  wait_http "http://localhost:3100/ready" 30 || { echo "[loki] Not ready in time; see $LOGS/loki-local.log"; exit 1; }

  start_otel
  # Health check extension listens on 13133 when enabled
  wait_http "http://localhost:13133/health" 30 || echo "[otelcol] Health endpoint not ready yet; continuing..."

  echo
  echo "[ok] Started Loki (http://localhost:3100) and OTLP Collector (grpc://localhost:4317)"
  echo "[info] Logs: $LOGS/loki-local.log, $LOGS/otelcol-local.log"
}

stop_proc() {
  local name="$1"
  local pidfile="$PIDS/$name.pid"
  if [ -f "$pidfile" ]; then
    local pid
    pid="$(cat "$pidfile")"
    if kill -0 "$pid" >/dev/null 2>&1; then
      echo "Stopping $name (PID $pid)..."
      kill "$pid" || true
      sleep 1
      if kill -0 "$pid" >/dev/null 2>&1; then
        echo "Force killing $name (PID $pid)..."
        kill -9 "$pid" || true
      fi
    fi
    rm -f "$pidfile"
  else
    echo "$name not running (no pid file)"
  fi
}

stop_all() {
  stop_proc "otelcol"
  stop_proc "loki"
  echo "[ok] Stopped"
}

status() {
  echo "== Status =="
  if [ -f "$PIDS/loki.pid" ]; then
    pid=$(cat "$PIDS/loki.pid"); ps -p "$pid" -o pid,cmd= || echo "Loki pid file present but process missing"
  else
    echo "Loki: no pid file"
  fi
  if [ -f "$PIDS/otelcol.pid" ]; then
    pid=$(cat "$PIDS/otelcol.pid"); ps -p "$pid" -o pid,cmd= || echo "otelcol pid file present but process missing"
  else
    echo "otelcol: no pid file"
  fi
  echo "-- HTTP checks --"
  curl -s -o /dev/null -w "Loki /ready: %{http_code}\n" http://localhost:3100/ready || true
  curl -s -o /dev/null -w "OTel /health: %{http_code}\n" http://localhost:13133/health || true
}

case "${1:-start}" in
  start) start_all ;;
  stop) stop_all ;;
  restart) stop_all; start_all ;;
  status) status ;;
  *) echo "Usage: $0 {start|stop|restart|status}"; exit 1 ;;
esac