# Configuration Migration Specification

## ADDED Requirements

### Requirement: Automated Config Update
The `update_configs.py` script MUST automatically migrate existing SSE configurations to Streamable HTTP.

#### Acceptance Criteria
- [ ] **AC-1:** `update_configs.py` identifies configs with `transport: "sse"`.
- [ ] **AC-2:** Detected SSE configs are rewritten to `transport: "streamable-http"`.
- [ ] **AC-3:** URLs ending in `/sse` are rewritten to `/mcp`.
- [ ] **AC-4:** The script creates a backup of the original configuration before modifying.

#### Scenario: Migrate SSE Config
Given a `claude.json` or `profiles.json` with `transport: sse` and `/sse` URLs
When `update_configs.py` is run
Then the transport should be updated to `streamable-http`
And the URL path should be updated from `/sse` to `/mcp`

## MODIFIED Requirements

### Requirement: Configuration Templates
New configuration templates MUST use Streamable HTTP defaults.

#### Acceptance Criteria
- [ ] **AC-1:** `config-templates/claude-code-sse.json` is replaced by or updated to `config-templates/claude-code-http.json`.
- [ ] **AC-2:** The new template uses `streamable-http` and `/mcp` endpoints.

#### Scenario: Use Template
Given I want to configure a new client
When I use the `claude-code-http.json` template
Then it should contain `streamable-http` transport settings by default

### Requirement: Project Defaults
The project configuration MUST default to Streamable HTTP for all new daemon and profile setups.

#### Acceptance Criteria
- [ ] **AC-1:** Fresh installations of the daemon default to HTTP mode.
- [ ] **AC-2:** New profiles created via `mcpm` use HTTP transport unless specified otherwise.

#### Scenario: New Profile
Given I create a new MCPM profile
When it is served by the daemon
Then it MUST use Streamable HTTP transport by default
