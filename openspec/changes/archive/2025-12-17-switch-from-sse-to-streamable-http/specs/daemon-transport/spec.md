# Daemon Transport Specification

## MODIFIED Requirements

### Requirement: Transport Protocol
The `mcpm-daemon` and managed profiles MUST expose Streamable HTTP endpoints instead of SSE.

#### Scenario: Daemon startup with HTTP transport
- **WHEN** the `mcpm-daemon` is started
- **THEN** it SHALL listen on ports 6276-6279
- **AND** it SHALL accept HTTP POST requests at `/mcp`
- **AND** it SHALL NOT require `Accept: text/event-stream` header

#### Scenario: MCP handshake over HTTP
- **WHEN** a client sends an initialization request to `/mcp`
- **THEN** the daemon SHALL handle the MCP handshake over HTTP
- **AND** `curl -X POST http://localhost:6276/mcp` SHALL return a valid JSON-RPC response
