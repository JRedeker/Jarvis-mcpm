# Ticket: Defer Docker-based telemetry deployment

Status: Deferred
Severity: Medium
Created: 2025-11-12
Owner: Platform Infra

Summary
- We attempted to bring up an OpenTelemetry Collector + Loki stack locally to enable OTLP logging ingestion, but the host environment cannot reliably run Docker and is missing several required CLI utilities (curl/unzip/tar/gzip).
- Rather than spending time installing Docker or fixing the environment now, we are deferring the Docker-based local telemetry deployment and proceeding with the JSONL fallback and Logs MCP for the short term.

Why we deferred
- On inspection the environment is WSL on Ubuntu 24.04, Docker CLI is present but the Docker daemon/socket is not available on the host, and essential extraction tools were missing. Attempts to bootstrap a local telemetry stack using the repository helper script failed.
- Installing Docker Desktop / enabling WSL integration or installing system packages may be acceptable later, but is outside the scope of this CI/QA step and blocks progress now.

What we did already
- Implemented OpenTelemetry SDK integration in-app and installed OTel Python packages (local .venv).
- Added a local bootstrap script: [`scripts/run-telemetry-local.sh`](scripts/run-telemetry-local.sh:1) to run binaries without Docker (attempted, but environment lacks unzip/tar/gzip).
- Created a blocker: [`open/blockers/docker-not-available.md`](open/blockers/docker-not-available.md:1).
- Created guidance & config files: [`docs/loki-local-config.yaml`](docs/loki-local-config.yaml:1) and collector examples [`docs/otel-collector-loki-config.yaml`](docs/otel-collector-loki-config.yaml:1).
- Implemented Logs MCP to query Loki when available: [`servers/logs-mcp.py`](servers/logs-mcp.py:1).
- Kept JSONL rotating fallback active in `./logs` for forensics and offline use (apps still write structured JSONL).

Impact / Risks
- Centralized queryable logs (Loki/Grafana) are not available locally; Logs MCP queries will return empty until a backend is reachable.
- Debugging across services will rely on local JSONL files and in-app structured logs.
- No immediate data loss: JSONL fallback is active; redaction and sampling settings still apply.

Next steps / Options to unblock (pick one)
1) Remotely-hosted Loki (recommended short-term)
   - Provision a managed Loki/observability endpoint (Grafana Cloud, self-hosted accessible endpoint).
   - Update [`cipher.yml`](cipher.yml:1) with LOKI_BASE_URL and set the Logs MCP backend to that URL so Logs MCP queries succeed immediately.
   - Acceptance: Logs from apps appear in the remote Loki and Logs MCP queries return results.

2) Enable Docker (preferred local path)
   - On the developer/machine: install Docker Desktop and enable WSL integration, or install Docker Engine on WSL and start the daemon.
   - Re-run local deploy script: [`scripts/deploy-otlp-stack.sh`](scripts/deploy-otlp-stack.sh:1) or local bootstrap: [`scripts/run-telemetry-local.sh`](scripts/run-telemetry-local.sh:1).
   - Acceptance: docker ps shows containers including otel-collector and loki; test script [`scripts/test-otlp-logs.py`](scripts/test-otlp-logs.py:1) verifies ingestion.

3) Native binaries (no Docker)
   - Manually install required tools (curl/unzip/tar/gzip) and run the bootstrap helper:
     - `./scripts/run-telemetry-local.sh start` (will download binaries and start loki + otelcol).
   - Acceptance: Collector gRPC at 4317 and Loki HTTP at 3100 are reachable and ingestion verified.

Action items (deferred)
- [ ] Decide which unblock path to follow (Remote Loki / Docker / Native binaries).
- [ ] If choosing Remote Loki: provision endpoint and update `cipher.yml` with LOKI_BASE_URL.
- [ ] If choosing Docker or Native: allocate time to perform host installs or enable WSL Docker integration.
- [ ] Re-run verification: [`scripts/test-otlp-logs.py`](scripts/test-otlp-logs.py:1) and validate Logs MCP queries.

Acceptance criteria to close this ticket
- Either:
  - Loki is reachable (local or remote) and OTLP logs show the expected structured fields (service.name, event, latency_ms, http_status, tier_id, model, token usage, cost), and Logs MCP returns results; OR
  - A documented, approved alternative (remote Loki) is in place and validated.
- JSONL fallback remains active and contains equivalent structured payloads for forensics.
- The previously-open blocker [`open/blockers/docker-not-available.md`](open/blockers/docker-not-available.md:1) is closed or updated to reflect the chosen unblock path.

References
- Decision/spec: [`tickets/logging-performance-options.md`](tickets/logging-performance-options.md:1)
- Blocker details: [`open/blockers/docker-not-available.md`](open/blockers/docker-not-available.md:1)
- Bootstrap script: [`scripts/run-telemetry-local.sh`](scripts/run-telemetry-local.sh:1)
- Local Loki config: [`docs/loki-local-config.yaml`](docs/loki-local-config.yaml:1)
- Logs MCP server: [`servers/logs-mcp.py`](servers/logs-mcp.py:1)