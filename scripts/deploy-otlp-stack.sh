#!/bin/bash
# OTLP Logging Stack Deployment for Cipher MCP
#
# This script deploys a local observability stack for OTLP logging:
# - OpenTelemetry Collector (receives OTLP logs via gRPC)
# - Grafana Loki (log storage and querying)
# - Grafana (visualization)
#
# Usage:
#   ./scripts/deploy-otlp-stack.sh [start|stop|restart|status|logs]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
STACK_NAME="cipher-otlp"

# Default command
COMMAND="${1:-start}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! docker ps &> /dev/null; then
        log_error "Docker daemon is not running. Please start Docker."
        exit 1
    fi

    log_info "Docker is available"
}

# Create Docker Compose file
create_docker_compose() {
    cat > "${PROJECT_ROOT}/docker-compose.otlp.yml" << 'EOF'
version: '3.8'

services:
  # OpenTelemetry Collector
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    container_name: cipher-otel-collector
    command: ["--config=/etc/otel-collector-config.yml"]
    volumes:
      - ./config/otel-collector-config.yml:/etc/otel-collector-config.yml
    ports:
      - "4317:4317"   # OTLP gRPC receiver
      - "4318:4318"   # OTLP HTTP receiver
      - "8888:8888"   # Prometheus metrics
      - "13133:13133" # Health check
    networks:
      - cipher-logs
    restart: unless-stopped

  # Grafana Loki
  loki:
    image: grafana/loki:latest
    container_name: cipher-loki
    command: -config.file=/etc/loki/local-config.yaml
    ports:
      - "3100:3100"
    volumes:
      - ./config/loki-config.yml:/etc/loki/local-config.yaml
      - loki-data:/loki
    networks:
      - cipher-logs
    restart: unless-stopped

  # Grafana for visualization
  grafana:
    image: grafana/grafana:latest
    container_name: cipher-grafana
    ports:
      - "3000:3000"
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
      - ./config/grafana-datasources.yml:/etc/grafana/provisioning/datasources/datasources.yml
    networks:
      - cipher-logs
    restart: unless-stopped
    depends_on:
      - loki

networks:
  cipher-logs:
    driver: bridge

volumes:
  loki-data:
  grafana-data:
EOF
    log_info "Created docker-compose.otlp.yml"
}

# Create OTel Collector config
create_otel_config() {
    mkdir -p "${PROJECT_ROOT}/config"

    cat > "${PROJECT_ROOT}/config/otel-collector-config.yml" << 'EOF'
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 2s
    send_batch_size: 512
    send_batch_max_size: 1024

  attributes:
    actions:
      - key: sensitive
        action: delete
      - key: authorization
        action: delete

  resource:
    attributes:
      - key: deployment.environment
        value: local
        action: upsert

exporters:
  loki:
    endpoint: http://loki:3100/loki/api/v1/push
    labels:
      resource:
        service.name: "service_name"
        service.instance.id: "service_instance_id"
      attributes:
        level: "level"
        session_id: "session_id"
        tier_id: "tier_id"
        event: "event"

  logging:
    loglevel: info

  prometheus:
    endpoint: "0.0.0.0:8888"

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch, attributes, resource]
      exporters: [loki, logging]

  telemetry:
    logs:
      level: info
    metrics:
      level: detailed
      address: 0.0.0.0:8888
EOF
    log_info "Created OTel Collector config"
}

# Create Loki config
create_loki_config() {
    cat > "${PROJECT_ROOT}/config/loki-config.yml" << 'EOF'
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
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
      schema: v11
      index:
        prefix: index_
        period: 24h

ruler:
  alertmanager_url: http://localhost:9093

limits_config:
  retention_period: 168h  # 7 days
  max_query_length: 365d
  reject_old_samples: false
  allow_structured_metadata: true
EOF
    log_info "Created Loki config"
}

# Create Grafana datasource config
create_grafana_config() {
    cat > "${PROJECT_ROOT}/config/grafana-datasources.yml" << 'EOF'
apiVersion: 1

datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    isDefault: true
    jsonData:
      maxLines: 1000
    version: 1
    editable: true
EOF
    log_info "Created Grafana datasource config"
}

# Start the stack
start_stack() {
    log_info "Starting OTLP logging stack..."

    check_docker
    create_docker_compose
    create_otel_config
    create_loki_config
    create_grafana_config

    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.otlp.yml up -d

    log_info "Waiting for services to be ready..."
    sleep 5

    # Check health
    if curl -s http://localhost:13133/health > /dev/null 2>&1; then
        log_info "✓ OTel Collector is healthy"
    else
        log_warn "✗ OTel Collector may not be ready yet"
    fi

    if curl -s http://localhost:3100/ready > /dev/null 2>&1; then
        log_info "✓ Loki is healthy"
    else
        log_warn "✗ Loki may not be ready yet"
    fi

    log_info ""
    log_info "OTLP Stack is running!"
    log_info "  - OTel Collector gRPC: http://localhost:4317"
    log_info "  - OTel Collector HTTP: http://localhost:4318"
    log_info "  - Loki API: http://localhost:3100"
    log_info "  - Grafana UI: http://localhost:3000 (admin/admin)"
    log_info ""
    log_info "View logs: ./scripts/deploy-otlp-stack.sh logs"
}

# Stop the stack
stop_stack() {
    log_info "Stopping OTLP logging stack..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.otlp.yml down
    log_info "Stack stopped"
}

# Restart the stack
restart_stack() {
    stop_stack
    sleep 2
    start_stack
}

# Show stack status
show_status() {
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.otlp.yml ps
}

# Show logs
show_logs() {
    SERVICE="${2:-}"
    cd "$PROJECT_ROOT"

    if [ -z "$SERVICE" ]; then
        docker compose -f docker-compose.otlp.yml logs -f
    else
        docker compose -f docker-compose.otlp.yml logs -f "$SERVICE"
    fi
}

# Main command handler
case "$COMMAND" in
    start)
        start_stack
        ;;
    stop)
        stop_stack
        ;;
    restart)
        restart_stack
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs "$@"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs [service]}"
        echo ""
        echo "Services: otel-collector, loki, grafana"
        exit 1
        ;;
esac