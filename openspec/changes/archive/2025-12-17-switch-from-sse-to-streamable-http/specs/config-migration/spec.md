# Configuration Migration Specification

## ADDED Requirements

### Requirement: Automated Config Update
The `update_configs.py` script MUST automatically migrate existing SSE configurations to Streamable HTTP.

#### Scenario: Migrate SSE Config
- **WHEN** `update_configs.py` is run against a config with `transport: "sse"`
- **THEN** it SHALL update the transport to `streamable-http`
- **AND** it SHALL rewrite `/sse` URLs to `/mcp`
- **AND** it SHALL create a backup before modifying

#### Scenario: Identify SSE configurations
- **WHEN** `update_configs.py` scans configuration files
- **THEN** it SHALL identify all configs with `transport: "sse"`
- **AND** it SHALL report which files will be migrated

## MODIFIED Requirements

### Requirement: Configuration Templates
New configuration templates MUST use Streamable HTTP defaults.

#### Scenario: Use HTTP Template
- **WHEN** a user applies the `claude-code-http.json` template
- **THEN** it SHALL contain `transport: "streamable-http"` settings
- **AND** it SHALL use `/mcp` endpoint paths

### Requirement: Project Defaults
The project configuration MUST default to Streamable HTTP for all new daemon and profile setups.

#### Scenario: New Profile Default Transport
- **WHEN** a new MCPM profile is created
- **THEN** it SHALL use Streamable HTTP transport by default
- **AND** fresh daemon installations SHALL default to HTTP mode
