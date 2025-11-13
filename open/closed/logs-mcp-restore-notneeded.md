# Logs MCP Restore

## Description

Restore the logs-mcp server to full operational status, including configuration, dependencies, and integration with the cipher-aggregator.

## Status

Open

## Priority

High

## Related Tickets

- [ ] #docker-deferred (if related to deployment)
- [ ] #routing-metadata-connection-issue (if related to integration)

## Steps

1. Verify the server script at servers/logs-mcp.py
2. Check dependencies in eval_llm_venv
3. Test the tool integration with cipher.yml configuration
4. Run verification tests for log querying

## Acceptance Criteria

- [ ] The logs-mcp server starts without errors
- [ ] Tool calls to logs-mcp succeed
- [ ] Integration with other MCP servers works
- [ ] No errors in logs/cipher-aggregator.log

## Environment Notes

- Current configuration in cipher.yml under mcpServers.logs
- Dependencies: Python 3, httpx, opentelemetry-api, opentelemetry-sdk, opentelemetry-exporter-otlp-proto-grpc
- Test script: scripts/test-otlp-logs.py
