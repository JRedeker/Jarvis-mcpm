# Ticket: Performance logging architecture for Cipher (local logs and OpenRouter requests)

Summary
- Goal: Improve logging performance and consistency for Cipher’s local logs in ./logs/ and OpenRouter request logs, minimizing hot-path I/O and standardizing structured JSON output with rotation and optional shipping to external backends.
- Scope: Replace synchronous on-path file writes with a queue-based logger, adopt structured JSON logs, configure rotation/compression, add sampling/toggles, and optionally integrate OpenTelemetry or a log shipper in later phases.

Current state and bottlenecks (observed)
- LLM inference server initializes global logging and writes a separate JSONL cost file:
  - Global config via [logging.basicConfig()](servers/llm-inference-mcp.py:35)
  - Module logger [logger](servers/llm-inference-mcp.py:36)
  - Cost logging implemented as manual file append in [log_cost()](servers/llm-inference-mcp.py:51) using [open()](servers/llm-inference-mcp.py:80) to write to [COST_LOG_FILE](servers/llm-inference-mcp.py:39) (JSONL)
  - OpenRouter HTTP client usage [httpx.AsyncClient()](servers/llm-inference-mcp.py:456) with request/response handling on the hot path
- Routing metadata server uses basic file + console handlers:
  - File handler and stderr stream via [logging.FileHandler()](servers/routing-metadata-mcp.py:41) and [logging.StreamHandler()](servers/routing-metadata-mcp.py:43)
- Health server example shows rotating logs but not JSON formatting:
  - Rotation pattern via [logging.handlers.RotatingFileHandler()](tests-and-notes/health-server.py:382)
- Bottlenecks:
  - Synchronous file writes on hot paths (e.g., manual open/write in cost logging) can block request handling.
  - Mixed formats (plain text vs JSONL) and inconsistent field schema reduce downstream usability.
  - No compression on rotated logs and inconsistent rotation policies.
  - Sensitive headers/fields may be emitted without guards in debug scenarios.
  - No explicit sampling controls for verbose subsystems (e.g., httpx).

Non-goals
- Centralized log aggregation mandated in Phase 1 (remains optional)
- Replacing metrics/Prometheus; this ticket focuses on logs

Objectives
- Non-blocking logging on hot paths using a queue + background writer.
- Structured JSON logs with stable schema across servers.
- File rotation and compression policies to bound disk usage and improve locality.
- Configurable verbosity, sampling, and sensitive-field filtering.
- Optional Phase 2: ship logs to OTLP collector (OpenTelemetry) or external log stores.

Options (with tradeoffs)

A) Stdlib Queue-based JSON logging (recommended for Phase 1)
- Pattern: producers attach a QueueHandler that enqueues records; a single QueueListener runs a JSON FileHandler with rotation.
- Serialization: prefer orjson (fast) with fallback to json for portability; emit JSONL per record for easy tailing and shipping.
- Rotation: size-based (RotatingFileHandler) or time-based (TimedRotatingFileHandler) with post-rotate gzip.
- Pros: Minimal dependency footprint, very low overhead on hot path, fully local/offline compatible, simple ops in ./logs/.
- Cons: No centralized view by default; multi-process safety for rotation can require care (see concurrent-log-handler alternative).

B) OpenTelemetry logs exporter (optional Phase 2)
- Pattern: keep local queue logger, add OTLP LogRecord exporter to a collector; unify with traces/metrics.
- Pros: Centralized aggregation, correlation with spans/metrics, flexible backends (Loki, Elastic, OpenSearch).
- Cons: Additional moving parts, configuration cost, network overhead, tighter dependency on collector availability.

C) Local queue + external shipper (Fluent Bit or Vector) watching ./logs/*.jsonl
- Pattern: app writes fast locally (Option A), shipper tails files and forwards to remote store.
- Pros: Keeps application concerns minimal, robust shipping/retry/transform handled by sidecar/daemon.
- Cons: Operational complexity of managing a shipper across environments.

Recommendation
- Phase 1: Implement Option A (stdlib queue + JSON logs + rotation + compression) across MCP servers (llm-inference, routing-metadata), including OpenRouter cost logs. Keep changes local, additive, and behind env toggles.
- Phase 2: Add optional OTLP logs exporter to a collector, guarded by env flags; keep JSONL file sink as the default fallback.
- Phase 3: If centralization without OTLP is preferred, adopt Fluent Bit or Vector to tail ./logs/*.jsonl and send to Loki/Elastic.

Proposed design details

1) Non-blocking queue-based logging
- Introduce a startup initializer per server to:
  - Create one multiprocessing-safe queue (reasonable maxsize to apply backpressure).
  - Attach QueueHandler to the root or per-module logger.
  - Start a QueueListener with a JSON FileHandler (rotating), optionally also console.
- Replace synchronous open() writes for cost logs in [log_cost()](servers/llm-inference-mcp.py:51) with logger.emit of a “cost” event (JSON), so the file I/O occurs only in the background listener.
- Retain explicit file targets, but unify under logs/ with consistent filenames:
  - logs/llm-inference.jsonl
  - logs/openrouter-costs.jsonl
  - logs/routing-metadata.jsonl

2) Structured JSON schema (consistent keys)
- Minimal baseline keys in all records:
  - ts (ISO8601), level, logger, pid, thread, session_id (if available), request_id (if available), tool/server name, event, fields (object)
- LLM/OpenRouter-specific fields (emit from llm-inference):
  - tier_id, model, tokens: {prompt, completion, total}, cost: {input, output, total}, latency_ms, http_status, provider: "openrouter"
- Routing metadata-specific fields:
  - selected_tool, recommended_tool, is_compliant, detected_domain, status, session_id

3) Rotation and compression
- Use RotatingFileHandler size-based policy for JSONL outputs (e.g., 10–50 MB per segment) with backupCount N.
- Add gzip compression on rollover (custom rotator and namer) to keep older segments compact.
- Use TimedRotatingFileHandler for very chatty logs if we prefer chronological cutovers.

4) Sensitive data filtering and sampling
- Ensure headers and payloads are filtered; reuse patterns similar to OpenAI SDK’s filter:
  - For debug-level http logs, never emit Authorization or cookies; redact tokens.
- Sampling controls (env toggles) for verbose subsystems (e.g., only 1 of N debug http events).

5) Configuration toggles (cipher.yml)
- Add environment flags to servers for tuning without code changes:
  - LOG_LEVEL (default INFO)
  - LOG_JSON (default true)
  - LOG_ROTATE_SIZE (e.g., bytes: 10_485_760)
  - LOG_BACKUPS (e.g., 5)
  - LOG_GZIP (true/false)
  - LOG_QUEUE_MAX (queue size)
  - LOG_SAMPLE_HTTP (integer; e.g., sample 1 in N http logs at DEBUG)
- Apply for:
  - LLM Inference server at [cipher.yml](cipher.yml:319)
  - Routing Metadata server at [cipher.yml](cipher.yml:330)

Targeted code touch points

- llm-inference
  - Replace [logging.basicConfig()](servers/llm-inference-mcp.py:35) with queue-based initializer and JSON file sink (logs/llm-inference.jsonl).
  - In [log_cost()](servers/llm-inference-mcp.py:51), emit a structured log record (event: "openrouter_cost") instead of manual [open()](servers/llm-inference-mcp.py:80) writes. Retain JSONL filename but route via the queue listener’s rotating handler for logs/openrouter-costs.jsonl.
  - Around [httpx.AsyncClient()](servers/llm-inference-mcp.py:456), capture:
    - latency_ms, http_status, request_id header (if present), and the tier_id/model mapped in the same structured event.
  - Avoid logging request/response bodies by default; only status/metadata at INFO; enable DEBUG via LOG_LEVEL for deeper diagnostics with strict redaction.

- routing-metadata
  - Replace plain text [logging.FileHandler()](servers/routing-metadata-mcp.py:41) with JSON rotating file handler attached to the queue listener (logs/routing-metadata.jsonl).
  - Emit structured events for validation and tracking calls with fields selected_tool, recommended_tool, detected_domain, is_compliant, status.

- health-server (reference only)
  - Example of rotation exists at [logging.handlers.RotatingFileHandler()](tests-and-notes/health-server.py:382). Adopt JSON formatting and gzip rollover for consistency if this service is part of the deployable set.

Option A: Pros/Cons in Cipher context
- Pros
  - Lowest risk to implement; relies on stdlib components.
  - Eliminates hot-path blocking I/O; background listener performs file writes.
  - Unifies format and schema, simplifying downstream processing and shipping.
  - Works offline and in developer environments without extra infra.
- Cons
  - No centralized aggregation; cross-host queries require shipping solution.
  - Multi-process rotation safety can require either careful configuration or a 3rd-party rotating handler.

Option B: OTLP exporter (Phase 2)
- Pros
  - Central view and correlation with traces/metrics; can route to Loki/Tempo/Elastic via collector.
- Cons
  - Additional dependencies and infra; ensure graceful fallback when collector is unavailable.

Option C: External shipper (Fluent Bit/Vector)
- Pros
  - Application remains simple; shipper handles backpressure/retries.
- Cons
  - Operational complexity, packaging, and cross-platform installs.

Acceptance criteria (Phase 1)
- All server logs in ./logs/ are JSONL with a stable schema and include ts, level, logger, event, and fields.
- Hot-path synchronous file writes removed from LLM inference (cost and request logs are emitted via queue).
- Size-based rotation with N backups and gzip compression in place for llm-inference.jsonl, routing-metadata.jsonl, and openrouter-costs.jsonl.
- Default LOG_LEVEL=INFO, with DEBUG enabling additional non-sensitive diagnostics; sensitive data redacted at all levels.
- Measured overhead of logging on tool-call critical path ≤ 1 ms p95 at INFO level with queue logging enabled.
- Documented env toggles in cipher.yml server blocks for llm-inference and routing-metadata.

Rollout plan
1) Implement queue-based JSON logging (per server initializer), introduce JSON formatter (prefer orjson with fallback), enable rotation+gzip.
2) Replace manual cost file writes in [log_cost()](servers/llm-inference-mcp.py:51) with queue-emitted structured logs; keep artifact name logs/openrouter-costs.jsonl.
3) Add env toggles in [cipher.yml](cipher.yml:319) and [cipher.yml](cipher.yml:330); set sane defaults.
4) Benchmarks:
   - Micro: issue 1,000 synthetic tool calls to llm-inference and measure per-call added latency at INFO vs DEBUG, with and without queue.
   - Sustained: 10 minutes at 5–10 rps; verify no dropped records (or measure drops if queue saturation policy enabled), assess file rotations/compression.
5) QA: Validate schema fields in all logs; verify redaction; tail -f confirms continuous JSONL output; simulate full disk to confirm graceful degradation.
6) Optional: Add OTLP exporter feature flag and validate collector upload.

Out-of-scope notes
- Centralized log storage (e.g., Loki/Elastic) is deferred to Phase 2/3.
- Metrics remain under Prometheus MCP; logs and metrics correlation becomes available if/when OTLP is adopted.

References (touch points)
- LLM Inference logging init and cost logging:
  - [logging.basicConfig()](servers/llm-inference-mcp.py:35)
  - [logger](servers/llm-inference-mcp.py:36)
  - [log_cost()](servers/llm-inference-mcp.py:51)
  - [open()](servers/llm-inference-mcp.py:80)
  - [httpx.AsyncClient()](servers/llm-inference-mcp.py:456)
- Routing Metadata logging:
  - [logging.FileHandler()](servers/routing-metadata-mcp.py:41)
  - [logging.StreamHandler()](servers/routing-metadata-mcp.py:43)
- Health server rotation example:
  - [logging.handlers.RotatingFileHandler()](tests-and-notes/health-server.py:382)
- Cipher configuration for server env toggles:
  - [cipher.yml](cipher.yml:319)
  - [cipher.yml](cipher.yml:330)

Deliverables
- Phase 1 PR implementing queue-based JSON logging with rotation/gzip and schema standardization across llm-inference and routing-metadata servers, plus cipher.yml env toggles and a short benchmarking report.
- Phase 2 PR (optional) adding OTLP exporter integration behind env flags, with deployment notes.