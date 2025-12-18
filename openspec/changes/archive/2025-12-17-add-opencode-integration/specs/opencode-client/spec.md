## ADDED Requirements

### Requirement: OpenCode Client Detection
The system SHALL detect OpenCode installations by checking standard configuration file locations in priority order: environment variable override, project-local config, then global config.

#### Scenario: Detect OpenCode via environment variable
- **WHEN** `OPENCODE_CONFIG` environment variable is set
- **THEN** the system SHALL use that path as the OpenCode configuration file

#### Scenario: Detect OpenCode via project-local config
- **WHEN** `OPENCODE_CONFIG` is not set AND `./opencode.json` exists in the current directory
- **THEN** the system SHALL use that file as the OpenCode configuration

#### Scenario: Detect OpenCode via global config
- **WHEN** `OPENCODE_CONFIG` is not set AND no project-local config exists AND `~/.config/opencode/opencode.json` exists
- **THEN** the system SHALL use the global config file

#### Scenario: OpenCode not detected
- **WHEN** no OpenCode configuration file is found in any standard location
- **THEN** the system SHALL report that OpenCode is not configured

### Requirement: OpenCode Configuration Format
The system SHALL generate OpenCode-compatible MCP server configurations using the `mcp` object schema with `type`, `command`, `url`, `enabled`, `environment`, and `headers` properties.

#### Scenario: Configure local stdio server
- **WHEN** adding a stdio-based MCP server (e.g., Jarvis) to OpenCode
- **THEN** the system SHALL generate a config entry with `type: "local"` and `command` array

#### Scenario: Configure remote HTTP server
- **WHEN** adding an HTTP-based MCP server (e.g., memory profile) to OpenCode
- **THEN** the system SHALL generate a config entry with `type: "remote"` and `url` property

#### Scenario: Preserve existing configuration
- **WHEN** modifying an OpenCode config that contains other settings (theme, model, etc.)
- **THEN** the system SHALL preserve all non-MCP configuration properties

### Requirement: OpenCode Client Registry
The system SHALL include OpenCode in the client registry with standard configuration paths and format identifier.

#### Scenario: List available clients includes OpenCode
- **WHEN** a user queries available MCP clients
- **THEN** OpenCode SHALL appear in the list with its display name and detected status

#### Scenario: Configure client for OpenCode
- **WHEN** a user runs `manage_client` with `client_name: "opencode"`
- **THEN** the system SHALL use OpenCode-specific config paths and format

### Requirement: OpenCode Configuration Template
The system SHALL provide a configuration template for OpenCode that includes Jarvis (local) and common HTTP profiles.

#### Scenario: Generate starter configuration
- **WHEN** a user requests initial OpenCode MCP setup
- **THEN** the system SHALL provide a template with Jarvis as a local server and memory/project profiles as remote servers

#### Scenario: Template uses correct endpoints
- **WHEN** generating OpenCode configuration for HTTP profiles
- **THEN** the system SHALL use `/mcp` endpoint paths (not deprecated `/sse` paths)
