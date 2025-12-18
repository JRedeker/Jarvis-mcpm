## REMOVED Requirements

### Requirement: SSE Transport Fallback
**Reason**: SSE is deprecated in MCP spec 2025-03-26. Streamable HTTP is now the only supported remote transport.
**Migration**: All profiles use `--http` flag. Users with SSE configurations must update to use `/mcp` endpoints with `transport: "streamable-http"`.

#### Scenario: SSE fallback for problematic profiles
- **WHEN** a profile was configured in SSE_PROFILES array
- **THEN** it used `--sse` transport flag
- **STATUS**: Removed - all profiles now use `--http`

## MODIFIED Requirements

### Requirement: Daemon Transport Configuration
The mcpm-daemon SHALL use Streamable HTTP as the exclusive transport protocol for all MCP profiles.

#### Scenario: All profiles use HTTP transport
- **WHEN** the daemon starts any profile
- **THEN** it SHALL use the `--http` flag
- **AND** it SHALL NOT support `--sse` flag

#### Scenario: Morph profile uses HTTP
- **WHEN** the `morph` profile is started
- **THEN** it SHALL use HTTP transport like all other profiles
- **AND** the SSE_PROFILES exception array SHALL be removed

#### Scenario: Endpoint paths
- **WHEN** a client connects to a profile
- **THEN** the endpoint path SHALL be `/mcp`
- **AND** the `/sse` path SHALL NOT be available
