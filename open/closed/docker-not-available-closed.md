# Blocker: Docker not available — OTLP logging backend cannot be deployed

Status: Open
Severity: High (Infra)
Created: 2025-11-12
Owner: Platform Infra

Summary
- The OpenTelemetry OTLP Collector + Loki + Grafana stack cannot be started because Docker is not available on this host.
- This blocks the recommended centralized logging path. The JSONL rotating fallback under ./logs remains the active sink.
- OTLP exporter will attempt to use http://localhost:4317 but fail fast until the collector is running; failures are non-fatal when fallback is enabled.

Impact
- No centralized, queryable logs in Loki; Logs MCP queries will not work against a nonexistent backend.
- Only local JSONL logs are produced under ./logs with rotation+gzip.
- Affected implementation and docs:
  - [cipher.yml](cipher.yml)
  - [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py)
  - [servers/routing-metadata-mcp.py](servers/routing-metadata-mcp.py)
  - [servers/logs-mcp.py](servers/logs-mcp.py)
  - [docs/otel-collector-loki-config.yaml](docs/otel-collector-loki-config.yaml)
  - [scripts/deploy-otlp-stack.sh](scripts/deploy-otlp-stack.sh)
  - [scripts/test-otlp-logs.py](scripts/test-otlp-logs.py)
  - Related ticket: [tickets/logging-performance-options.md](tickets/logging-performance-options.md)

Evidence
- Reproduction command and output:
  - Command:
    - docker ps 2>/dev/null || echo "Docker not available - will use JSONL fallback only"
  - Output:
    - Docker not available - will use JSONL fallback only

Workaround (temporary)
- Continue using JSONL rotating logs under ./logs for forensics and development.
- Leave OTEL_* environment variables configured in [cipher.yml](cipher.yml); exporter failures will be non-fatal while the collector is down.
- Use application-level metrics/log sampling to keep I/O minimal until the backend is available.

Resolution paths
- Option A (Windows + WSL):
  - Install/launch Docker Desktop for Windows.
  - Enable WSL integration for the target Linux distribution (Docker Desktop Settings → Resources → WSL Integration).
  - From the WSL shell, verify:
    - docker ps
  - Start the telemetry stack:
    - ./scripts/deploy-otlp-stack.sh start
- Option B (Linux native daemon):
  - Install Docker Engine or distro package (example for Ubuntu):
    - sudo apt-get update
    - sudo apt-get install -y docker.io docker-compose-plugin
    - sudo usermod -aG docker "$USER"
    - newgrp docker
    - sudo systemctl enable --now docker
  - Verify:
    - docker ps
  - Start the telemetry stack:
    - ./scripts/deploy-otlp-stack.sh start
- Option C (no-Docker alternative):
  - Run OpenTelemetry Collector and Loki as native services or binaries.
  - Configure endpoints to match the provided pipeline in [docs/otel-collector-loki-config.yaml](docs/otel-collector-loki-config.yaml) (Collector OTLP gRPC on 4317, Loki on 3100).
  - Update [cipher.yml](cipher.yml) only if endpoints differ.

Validation checklist
- docker ps returns without error and lists running containers when the stack is up.
- ./scripts/deploy-otlp-stack.sh start succeeds; containers include otel-collector, loki, (optional) grafana.
- Ports reachable:
  - Collector gRPC: 4317
  - Loki HTTP API: 3100
  - Grafana UI (if enabled): 3000
- Run verifier:
  - ./scripts/test-otlp-logs.py
  - Expect synthetic logs to appear in Loki; Logs MCP queries from [servers/logs-mcp.py](servers/logs-mcp.py) return results.

Acceptance criteria
- Docker daemon is installed and running (or native services equivalent in place).
- Telemetry stack started successfully; Loki and Collector are reachable.
- OTLP logs are received in Loki with expected fields (service.name, event, latency_ms, http_status, tier_id, model, token usage, cost).
- JSONL fallback remains enabled for resilience; no secrets present in logs per redaction rules.

Notes
- Until this blocker is resolved, centralized log querying (Loki/Grafana, Logs MCP) is unavailable; local JSONL remains the authoritative source.
- If Docker cannot be installed in this environment, proceed with Option C and document the chosen native service method alongside any endpoint changes.

Links
- Plan/spec: [tickets/logging-performance-options.md](tickets/logging-performance-options.md)
- Config: [cipher.yml](cipher.yml)
- Collector config: [docs/otel-collector-loki-config.yaml](docs/otel-collector-loki-config.yaml)
- Deployment: [scripts/deploy-otlp-stack.sh](scripts/deploy-otlp-stack.sh)
- Verifier: [scripts/test-otlp-logs.py](scripts/test-otlp-logs.py)
- Logs MCP: [servers/logs-mcp.py](servers/logs-mcp.py)