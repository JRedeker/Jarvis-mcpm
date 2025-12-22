## MODIFIED Requirements

### Requirement: Client Management Tool
The system SHALL provide a `jarvis_client` tool for configuring MCP client applications with granular server-level control.

#### Scenario: List clients
- **WHEN** `jarvis_client` is called with `action=list`
- **THEN** the system SHALL return all known MCP clients with detection status

#### Scenario: Edit client
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and profile/server modifications
- **THEN** the system SHALL update the client configuration

#### Scenario: Import client config
- **WHEN** `jarvis_client` is called with `action=import` and `client_name=<client>`
- **THEN** the system SHALL import a template configuration

#### Scenario: Configure client path
- **WHEN** `jarvis_client` is called with `action=config`, `client_name=<client>`, and optional `config_path=<path>`
- **THEN** the system SHALL get or set the client config path

#### Scenario: Enable or disable server
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `set_enabled=<server>=<true|false>`
- **THEN** the system SHALL set the server's `enabled` field without removing the server
- **AND** the system SHALL preserve all other server configuration

#### Scenario: Set server environment variable
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `set_env=<server>:<KEY>=<value>`
- **THEN** the system SHALL add or update the environment variable on the server
- **AND** the system SHALL NOT log the value (security)

#### Scenario: Set server URL
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `set_url=<server>=<url>`
- **THEN** the system SHALL update the server's URL
- **AND** the system SHALL validate the URL format before writing

#### Scenario: Set server header
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `set_header=<server>:<Header-Name>=<value>`
- **THEN** the system SHALL add or update the HTTP header on the server
- **AND** the system SHALL NOT log the value (security)

#### Scenario: Remove server environment variable
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `remove_env=<server>:<KEY>`
- **THEN** the system SHALL remove the specified environment variable from the server

#### Scenario: Remove server header
- **WHEN** `jarvis_client` is called with `action=edit`, `client_name=<client>`, and `remove_header=<server>:<Header-Name>`
- **THEN** the system SHALL remove the specified HTTP header from the server

#### Scenario: Server not found error
- **WHEN** a granular edit operation targets a server that does not exist in the client config
- **THEN** the system SHALL return an error with the server name and list of available servers

#### Scenario: Invalid URL format error
- **WHEN** `set_url` is called with a malformed URL
- **THEN** the system SHALL return an error explaining the expected URL format
