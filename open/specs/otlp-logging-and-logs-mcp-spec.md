# Spec: OTLP logging for Cipher + Logs MCP for efficient querying

Purpose
- Implement high-performance logging across Cipher MCP servers using OpenTelemetry (OTLP).
- Primary transport: OTLP/Protobuf over gRPC to an OpenTelemetry Collector; fallback to local JSONL rotating logs.
- Provide a small Logs MCP to query backend logs (e.g., Loki) efficiently for developer/agent workflows.

Out of scope
- Full tracing/metrics deployment (logs only). Future work can unify logs, traces, metrics via OTel.

Architecture overview
- Data flow:
  - MCP servers (LLM inference, routing metadata) emit logs via OTel SDK → BatchLogRecordProcessor → OTLP exporter (gRPC, protobuf) → OTel Collector → Backend (Grafana Loki or Elastic/OpenSearch).
  - Local fallback: JSONL rotating logs with gzip in ./logs when enabled (for offline/dev and incident forensics).
- Control plane:
  - Env-based configuration in [cipher.yml](cipher.yml) per server for transport/batching/compression.
  - Optional logs MCP provides read/query tools for the log backend (domain-specific tool, aggregator-first).

Components affected
- Server initialization for OTel Logs (per-service):
  - LLM inference server (init around [logging.basicConfig()](servers/llm-inference-mcp.py:35); module logger [logger](servers/llm-inference-mcp.py:36))
  - Routing metadata server (logger setup at [logging.basicConfig()](servers/routing-metadata-mcp.py:38))
- Logging emission touch points:
  - Replace manual cost file writes in [log_cost()](servers/llm-inference-mcp.py:51) and [open()](servers/llm-inference-mcp.py:80) with OTel LogRecord emission.
  - Instrument HTTP call around [httpx.AsyncClient()](servers/llm-inference-mcp.py:456) to capture latency/http_status/request_id.
- Configuration:
  - Add OTEL_* env vars under server blocks in [cipher.yml](cipher.yml:319) and [cipher.yml](cipher.yml:330).
- Fallback reference:
  - Local rotation baseline at [logging.handlers.RotatingFileHandler()](tests-and-notes/health-server.py:382) and setup in [setup_logging()](tests-and-notes/health-server.py:379).

Transport, encoding, batching and compression (why and how)
- Primary: OTLP/Protobuf over gRPC
  - Fastest in CPU/bytes with persistent HTTP/2 multiplexed connection.
  - Env:
    - OTEL_LOGS_EXPORTER=otlp
    - OTEL_EXPORTER_OTLP_PROTOCOL=grpc
    - OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
    - OTEL_EXPORTER_OTLP_COMPRESSION=gzip (recommended off-host)
- Fallback: OTLP/HTTP with Protobuf (when gRPC blocked)
  - Env:
    - OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
    - OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
- Batching (reduce overhead; negligible delivery latency):
  - OTEL_BSP_SCHEDULE_DELAY=2000
  - OTEL_BSP_MAX_QUEUE_SIZE=2048
  - OTEL_BSP_MAX_EXPORT_BATCH_SIZE=512
- Avoid OTLP/JSON for steady state; use only for debugging.

Data model and redaction
- Resource attributes (per service):
  - service.name: “llm-inference-mcp” or “routing-metadata-mcp”
  - service.version: git SHA or release tag
  - service.instance.id: hostname/PID or UUID
- Common attributes:
  - session_id, tool/server, event, latency_ms, http_status, is_error
- LLM-specific:
  - tier_id, model, tokens.{prompt,completion,total}, cost.{input,output,total}
- Redaction and limits:
  - Never log Authorization/cookies; filter values in debug flows.
  - OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT=4096 (guardrail)

Fallback JSONL sink (when enabled)
- Use rotating JSONL (gzip on rollover) for dev/offline and resilience:
  - logs/llm-inference.jsonl
  - logs/routing-metadata.jsonl
  - logs/openrouter-costs.jsonl
- Reference rotation: [tests-and-notes/health-server.py](tests-and-notes/health-server.py:379)

Detailed implementation plan

1) Dependencies (Python)
- Add OTel logging exporter stack:
  - opentelemetry-api
  - opentelemetry-sdk (and/or opentelemetry-sdk-logs based on version)
  - opentelemetry-exporter-otlp-proto-grpc
  - opentelemetry-semantic-conventions

2) OTel setup per service
- Create a LoggerProvider with Resource attributes.
- Add BatchLogRecordProcessor with OTLPLogExporter (gRPC).
- Bridge stdlib logging so existing logger.info/error emit OTel logs.
- Maintain local JSONL file sink behind env toggles for fallback (dual-sink acceptable).

3) LLM inference server changes
- Replace manual file append in [log_cost()](servers/llm-inference-mcp.py:51) and [open()](servers/llm-inference-mcp.py:80):
  - Emit structured log with event="openrouter_cost" and attributes: tier_id, model, tokens.*, cost.*, task summary (truncated).
- Surround OpenRouter call at [httpx.AsyncClient()](servers/llm-inference-mcp.py:456):
  - Measure latency_ms (start-end).
  - Record http_status and request id header if available (e.g., X-Request-ID).
  - Emit event="openrouter_request" at INFO by default with attributes; DEBUG path can include headers minus sensitive fields.
- Ensure sensitive field redaction (Authorization, cookies) on any debug/trace records.

4) Routing metadata server changes
- Initialize OTel logger provider and exporter at startup (around [logging.basicConfig()](servers/routing-metadata-mcp.py:38)).
- Emit structured events for:
  - validate_tool_selection result (selected_tool, recommended_tool, detected_domain, is_compliant, status).
  - track_tool_execution (execution_time_ms, success, error_message).

5) cipher.yml configuration
- Under LLM inference [cipher.yml](cipher.yml:319) and routing metadata [cipher.yml](cipher.yml:330), add env:
  - OTEL_LOGS_EXPORTER=otlp
  - OTEL_EXPORTER_OTLP_PROTOCOL=grpc
  - OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
  - OTEL_EXPORTER_OTLP_COMPRESSION=gzip
  - OTEL_BSP_SCHEDULE_DELAY=2000
  - OTEL_BSP_MAX_QUEUE_SIZE=2048
  - OTEL_BSP_MAX_EXPORT_BATCH_SIZE=512
  - OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT=4096
  - Retain LOG_JSON/rotation flags for fallback.

6) Collector configuration (reference)
- Collector receives OTLP logs and forwards to backend.
- Example pipelines (conceptual):
  - receivers: otlp (grpc at :4317, http at :4318)
  - processors: batch, attributes (redaction if needed), filter
  - exporters: loki (or elastic/opensearch)
- Validate retention and label mapping to avoid high-cardinality blowups in Loki/Elastic.

7) Backend mapping (Loki example)
- Loki labels:
  - service_name, level, is_error, session_id (bounded), model/tier_id (bounded)
- Place verbose strings (task text) in log line/body; keep labels concise to control cardinality.

8) Performance targets and tests
- Overhead budget: ≤ 1 ms p95 added latency at INFO versus baseline for typical tool calls.
- Micro benchmark: 1,000 calls; compare baseline vs OTLP gRPC + batch + gzip.
- Load test: 10 min @ 5–10 rps; confirm zero exporter drops, acceptable Collector/Backend ingestion.
- Failure drills: stop Collector; exporter should retry/backoff; ensure hot path is not blocked and fallback JSONL continues.

9) Failure handling
- Export retries with backoff.
- If queue pressure persists, prefer dropping oldest/low-importance logs rather than blocking requests.
- Ensure JSONL fallback captures critical events.

10) Rollout plan
- Phase 1 (pilot): Enable OTel exporter + batch in routing metadata only; verify pipeline.
- Phase 2: Enable in LLM inference; turn on compression for non-local endpoints.
- Phase 3: Tune batch sizes and labels; enable dashboards in Grafana/Kibana; optionally disable JSONL fallback in production if unnecessary.
- Backout: Disable OTEL_LOGS_EXPORTER (or set to none) to fall back to JSONL only.

Small Logs MCP (spec)

Goals
- Provide a minimal, efficient domain-specific MCP to query logs backend (Loki first) for developer/agent workflows.
- Tools: range query, recent tail, counts, and labels listing. Bounded outputs and rate-limited.

Server
- Path: [servers/logs-mcp.py](servers/logs-mcp.py)
- Type: stdio (Python3)
- Env:
  - LOGS_BACKEND=loki
  - LOKI_BASE_URL (e.g., http://localhost:3100)
  - LOKI_TENANT (optional)
  - LOKI_TOKEN (optional bearer)
  - LOGS_DEFAULT_LIMIT=500
  - LOGS_TIMEOUT_MS=10000
  - LOGS_MAX_BYTES=200000  (cap response size)
- Dependencies: httpx

Tools and input schemas
- logs_query_range
  - Arguments:
    - query: required (LogQL string or simplified filters we translate into LogQL)
    - start, end: RFC3339 or epoch (ms)
    - limit: integer (default from env)
    - direction: “backward” | “forward” (default backward)
    - labels: object of key:value for label matchers (optional)
  - Returns: text content with a summarized table (timestamp, level, service, msg excerpt) and a counts footer; truncation note if hit size cap.
- logs_tail
  - Arguments:
    - query: required
    - since: RFC3339/epoch; default now-5m
    - limit: integer
    - labels: optional
  - Returns: last N entries with timestamps and key attributes; note on truncation if any.
- logs_count
  - Arguments:
    - query: required
    - start, end: required
    - labels: optional
  - Returns: simple JSON text summarizing total count and rate (per minute).
- logs_labels
  - Arguments:
    - prefix: optional string to filter label keys
  - Returns: available label keys and sample values (capped).

Rate limiting and safety
- Enforce per-call timeout (LOGS_TIMEOUT_MS).
- Reject calls with excessive time windows or query sizes; require narrowing if exceeds caps.
- Truncate response to LOGS_MAX_BYTES with a “[truncated]” note.

cipher.yml integration (example)
- Add server entry:
  - Name: logs
  - Command: python3
  - Args: /home/jrede/dev/MCP/servers/logs-mcp.py
  - Env: LOKI_BASE_URL, LOGS_DEFAULT_LIMIT, LOGS_TIMEOUT_MS, LOGS_MAX_BYTES, LOKI_TENANT/LOKI_TOKEN if applicable.

Example use flows
- Investigate a routing incident:
  - Call logs_query_range with labels {service_name: “routing-metadata-mcp”, session_id: “…”, is_error: true} and a 10m window.
- Check OpenRouter errors:
  - Query event="openrouter_request" and http_status>=500 in last 15m; show summaries and counts.
- Tail recent LLM cost logs:
  - logs_tail with event="openrouter_cost" and model="openai/gpt-5" for last 5m.

Security and privacy
- All auth tokens via env; do not log them.
- No PII in labels; short-lived session IDs OK if bounded.
- Respect tenant isolation headers for multi-tenant backends.

Acceptance criteria
- OTel logs exported via OTLP/Protobuf over gRPC to Collector, visible in backend.
- JSONL fallback present and working when enabled.
- Measured p95 overhead ≤ 1 ms at INFO; batch + compression functional.
- Logs MCP exposes the four tools and returns bounded, useful summaries for common queries.
- cipher.yml updated with OTEL_* flags per service and logs MCP server block.

Tasks checklist (to PRs)
- Add OTel dependencies, initialize in [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:35) and [servers/routing-metadata-mcp.py](servers/routing-metadata-mcp.py:38).
- Refactor [log_cost()](servers/llm-inference-mcp.py:51) to emit OTel logs; instrument HTTP call at [httpx.AsyncClient()](servers/llm-inference-mcp.py:456).
- Add OTEL_* env vars to [cipher.yml](cipher.yml:319) and [cipher.yml](cipher.yml:330); keep JSONL flags.
- Provide Collector config sample and backend mapping notes in this spec.
- Implement [servers/logs-mcp.py](servers/logs-mcp.py) with tools logs_query_range, logs_tail, logs_count, logs_labels.
- Add logs MCP block to [cipher.yml](cipher.yml) with env and timeouts.
- Benchmarks + validation documented; update ticket status and link to this spec.

Open questions
- Preferred backend (Loki vs Elastic/OpenSearch) for production? This affects label strategy and dashboards.
- Retention and storage quotas for logs at chosen backend.
- Should we gate DEBUG logs by sampling ratio env (LOG_SAMPLE_DEBUG_N) globally?