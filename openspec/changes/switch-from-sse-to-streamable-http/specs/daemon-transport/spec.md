# Daemon Transport Specification

## MODIFIED Requirements

### Requirement: Transport Protocol
The `mcpm-daemon` and managed profiles MUST expose Streamable HTTP endpoints instead of SSE.

#### Acceptance Criteria
- [ ] **AC-1:** The daemon listens on ports 6276-6279 and accepts standard HTTP POST requests at `/mcp` (or `/message`).
- [ ] **AC-2:** The daemon properly handles the MCP initialization handshake over HTTP.
- [ ] **AC-3:** Legacy SSE endpoints are either removed or redirect/inform the user to use the new endpoint.
- [ ] **AC-4:** `curl -X POST http://localhost:6276/mcp` returns a valid MCP response (e.g., JSON-RPC error or success) without hanging.

#### Scenario: Start Daemon
Given the `mcpm-daemon` is started
When I check the listening ports (6276-6279)
Then they should accept HTTP POST requests at `/mcp` (or configured path)
And they should NOT require an `Accept: text/event-stream` header for initial connection handshake
