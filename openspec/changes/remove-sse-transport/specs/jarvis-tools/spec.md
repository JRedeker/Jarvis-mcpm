## MODIFIED Requirements

### Requirement: Server Transport Type Parameter
The `create_server` and `edit_server` tools SHALL only accept `stdio` and `streamable-http` as valid transport types.

#### Scenario: Valid transport types
- **WHEN** a user calls `create_server` with `type` parameter
- **THEN** the tool SHALL accept `stdio` or `streamable-http`
- **AND** the tool SHALL NOT accept `sse` as a valid type

#### Scenario: Tool description accuracy
- **WHEN** the tool schema is displayed to an AI agent
- **THEN** the description SHALL state "Transport type: 'stdio' or 'streamable-http'"
- **AND** the description SHALL NOT mention `sse`

#### Scenario: URL parameter description
- **WHEN** the `url` parameter description is displayed
- **THEN** it SHALL state "URL (for streamable-http type)"
- **AND** it SHALL NOT reference `sse`

### Requirement: Health Check Endpoints
The `check_status` tool SHALL verify Streamable HTTP endpoints exclusively.

#### Scenario: Endpoint health verification
- **WHEN** `check_status` checks daemon connectivity
- **THEN** it SHALL test `/mcp` endpoints
- **AND** it SHALL NOT test `/sse` endpoints
