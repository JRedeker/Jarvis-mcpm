# Jarvis Tools Specification

## MODIFIED Requirements

### Requirement: Server Creation
The `create_server` tool MUST support `streamable-http` as a transport type and use it as the default for network-based servers.

#### Scenario: Create HTTP Server
- **WHEN** user calls `create_server` with `type="streamable-http"` and `url="http://localhost:3000/mcp"`
- **THEN** the server SHALL be registered successfully
- **AND** `check_status` SHALL be able to validate its health

#### Scenario: Default Network Transport
- **WHEN** user creates a server with `url` but no `type` specified
- **THEN** the tool SHALL default `type` to `streamable-http`
- **AND** it SHALL NOT default to `sse`

### Requirement: Health Checks
The `check_status` tool MUST validate Streamable HTTP endpoints.

#### Scenario: Check HTTP Endpoint
- **WHEN** `check_status` runs against a Streamable HTTP server
- **THEN** it SHALL perform an MCP health probe
- **AND** it SHALL report the server as "healthy" if reachable

#### Scenario: Unreachable HTTP Endpoint
- **WHEN** `check_status` runs against an unreachable HTTP endpoint
- **THEN** it SHALL fail gracefully with a descriptive error message
