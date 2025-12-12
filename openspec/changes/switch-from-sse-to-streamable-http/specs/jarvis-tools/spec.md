# Jarvis Tools Specification

## MODIFIED Requirements

### Requirement: Server Creation
The `create_server` tool MUST support `streamable-http` as a transport type and use it as the default for network-based servers.

#### Acceptance Criteria
- [ ] **AC-1:** `create_server` accepts `type="streamable-http"`.
- [ ] **AC-2:** When `type` is omitted but `url` is provided, `create_server` defaults `type` to `streamable-http`.
- [ ] **AC-3:** The generated server configuration in `servers.json` correctly reflects these values.

#### Scenario: Create HTTP Server
Given I use `create_server`
When I specify `type="streamable-http"` and `url="http://localhost:3000/mcp"`
Then the server should be registered successfully
And `check_status` should be able to validate its health

#### Scenario: Default Network Transport
Given I create a new remote server
When I do not specify a transport type but provide a URL
Then it MUST default to `streamable-http` instead of `sse`

### Requirement: Health Checks
The `check_status` tool MUST validate Streamable HTTP endpoints.

#### Acceptance Criteria
- [ ] **AC-1:** `check_status` successfully pings a running Streamable HTTP server.
- [ ] **AC-2:** `check_status` fails gracefully with a descriptive error if the HTTP endpoint is unreachable or returns 5xx.
- [ ] **AC-3:** `doctor` command output includes status for HTTP-based servers.

#### Scenario: Check HTTP Endpoint
Given a running Streamable HTTP server
When `check_status` runs
Then it should report the server as "healthy" by performing a standard MCP health probe (e.g., JSON-RPC ping or capabilities check)
