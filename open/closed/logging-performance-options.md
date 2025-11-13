# Logging Performance Options for Cipher

## Description

Create a more performant way to create and store logs for Cipher. This includes:

- Local logs in ./logs/
- Logs for OpenRouter requests and other external API calls
- Ideally, use an existing tool within Cipher or add a new component to handle this efficiently

## Status

Open

## Priority

High

## Related Components

- cipher.yml (configuration)
- servers/ (potential new logging server or integration with existing ones)
- logs/ (local log directory)
- OpenRouter integration (external API logging)

## Options to Consider

1. **Use Existing Tools in Cipher**
   - Integrate with existing MCP servers (e.g., logs-mcp if available)
   - Leverage OpenTelemetry for structured logging (OTLP exporter to a backend)
   - Use filesystem MCP for log storage with compression/batching for performance

2. **Add New Logging Component**
   - Create a dedicated logging server (e.g., servers/logs-performance.py)
   - Implement batching, compression, and rotation for high-volume logs
   - Support for external services (e.g., OpenRouter) with request/response tracing

3. **Performance Goals**
   - Minimize I/O overhead for local logs
   - Structured logging for external requests (JSONL or OTLP)
   - Scalable storage (e.g., rotate files, use database for queries)
   - Low latency for log emission during high-load operations

## Acceptance Criteria

- [ ] Logs are emitted efficiently without blocking main processes
- [ ] Local logs in ./logs/ are performant (e.g., async writes, buffering)
- [ ] External API logs (OpenRouter) are captured with metadata (request ID, status, latency)
- [ ] Configuration in cipher.yml for log levels, backends, and rotation
- [ ] Verification script to test log emission and retrieval
- [ ] Documentation updated in docs/ for logging setup

## Next Steps

1. Evaluate current logging in cipher.py and servers/
2. Research integration with OpenTelemetry or similar
3. Prototype a logging middleware or server
4. Test with sample loads to measure performance gains

## Related Tickets

- [ ] #logs-mcp-restore (if integrating with existing logs server)
- [ ] #docker-deferred (if deployment affects logging)