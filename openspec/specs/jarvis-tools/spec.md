# jarvis-tools Specification

## Purpose
TBD - created by archiving change reduce-context-tokens. Update Purpose after archive.
## Requirements
### Requirement: Consolidated Tool Architecture
The system SHALL expose exactly 8 MCP tools using action-based consolidation to minimize context token consumption.

#### Scenario: Tool listing returns consolidated tools
- **WHEN** a client requests tools/list
- **THEN** the response SHALL contain exactly 8 tools
- **AND** all tool names SHALL be prefixed with `jarvis_`

#### Scenario: Payload size reduction
- **WHEN** the tool listing is serialized to JSON
- **THEN** the compact JSON payload SHALL be less than 4,000 bytes

---

### Requirement: Status Tool
The system SHALL provide a `jarvis_check_status` tool for system health checks.

#### Scenario: Health check execution
- **WHEN** `jarvis_check_status` is called with no parameters
- **THEN** the system SHALL return MCPM doctor output, API health, and daemon status

---

### Requirement: Server Management Tool
The system SHALL provide a `jarvis_server` tool for managing MCP servers with action-based routing.

#### Scenario: List servers
- **WHEN** `jarvis_server` is called with `action=list`
- **THEN** the system SHALL return all installed MCP servers with status

#### Scenario: Get server info
- **WHEN** `jarvis_server` is called with `action=info` and `name=<server>`
- **THEN** the system SHALL return detailed server configuration

#### Scenario: Install server
- **WHEN** `jarvis_server` is called with `action=install` and `name=<server>`
- **THEN** the system SHALL install the server from registry

#### Scenario: Uninstall server
- **WHEN** `jarvis_server` is called with `action=uninstall` and `name=<server>`
- **THEN** the system SHALL remove the server and clean up configuration

#### Scenario: Search servers
- **WHEN** `jarvis_server` is called with `action=search` and `query=<term>`
- **THEN** the system SHALL search the registry and return matches

#### Scenario: Edit server
- **WHEN** `jarvis_server` is called with `action=edit`, `name=<server>`, and optional config params
- **THEN** the system SHALL modify the server configuration

#### Scenario: Create server
- **WHEN** `jarvis_server` is called with `action=create`, `name=<server>`, and `type=<stdio|streamable-http>`
- **THEN** the system SHALL register a new custom server

#### Scenario: Usage statistics
- **WHEN** `jarvis_server` is called with `action=usage`
- **THEN** the system SHALL return tool usage statistics

#### Scenario: Invalid action error
- **WHEN** `jarvis_server` is called with an invalid action
- **THEN** the system SHALL return an error listing valid actions

---

### Requirement: Profile Management Tool
The system SHALL provide a `jarvis_profile` tool for managing MCPM profiles with action-based routing.

#### Scenario: List profiles
- **WHEN** `jarvis_profile` is called with `action=list`
- **THEN** the system SHALL return all configured profiles

#### Scenario: Create profile
- **WHEN** `jarvis_profile` is called with `action=create` and `name=<profile>`
- **THEN** the system SHALL create a new profile

#### Scenario: Edit profile
- **WHEN** `jarvis_profile` is called with `action=edit`, `name=<profile>`, and server modifications
- **THEN** the system SHALL modify the profile configuration

#### Scenario: Delete profile
- **WHEN** `jarvis_profile` is called with `action=delete` and `name=<profile>`
- **THEN** the system SHALL remove the profile

#### Scenario: Suggest profile
- **WHEN** `jarvis_profile` is called with `action=suggest`
- **THEN** the system SHALL analyze the environment and recommend profiles

#### Scenario: Restart profiles
- **WHEN** `jarvis_profile` is called with `action=restart` and optional `name=<profile>`
- **THEN** the system SHALL restart the specified or all profiles

---

### Requirement: Client Management Tool
The system SHALL provide a `jarvis_client` tool for configuring MCP client applications.

#### Scenario: List clients
- **WHEN** `jarvis_client` is called with `action=list`
- **THEN** the system SHALL return all known MCP clients with detection status

#### Scenario: Edit client
- **WHEN** `jarvis_client` is called with `action=edit`, `name=<client>`, and profile/server modifications
- **THEN** the system SHALL update the client configuration

#### Scenario: Import client config
- **WHEN** `jarvis_client` is called with `action=import` and `name=<client>`
- **THEN** the system SHALL import a template configuration

#### Scenario: Configure client path
- **WHEN** `jarvis_client` is called with `action=config`, `name=<client>`, and optional `path=<config_path>`
- **THEN** the system SHALL get or set the client config path

---

### Requirement: Configuration Tool
The system SHALL provide a `jarvis_config` tool for MCPM configuration management.

#### Scenario: Get config value
- **WHEN** `jarvis_config` is called with `action=get` and `key=<key>`
- **THEN** the system SHALL return the configuration value

#### Scenario: Set config value
- **WHEN** `jarvis_config` is called with `action=set`, `key=<key>`, and `value=<value>`
- **THEN** the system SHALL update the configuration

#### Scenario: List config
- **WHEN** `jarvis_config` is called with `action=list`
- **THEN** the system SHALL return all configuration values

#### Scenario: Migrate config
- **WHEN** `jarvis_config` is called with `action=migrate`
- **THEN** the system SHALL upgrade configuration to latest format with backup

---

### Requirement: Project Analysis Tool
The system SHALL provide a `jarvis_project` tool for project detection and DevOps scaffolding.

#### Scenario: Analyze project
- **WHEN** `jarvis_project` is called with `action=analyze`
- **THEN** the system SHALL detect languages, frameworks, and existing configs

#### Scenario: Fetch diff context
- **WHEN** `jarvis_project` is called with `action=diff` and optional `staged=<bool>`
- **THEN** the system SHALL return git status and diff for review

#### Scenario: Apply DevOps stack
- **WHEN** `jarvis_project` is called with `action=devops` and optional type/force params
- **THEN** the system SHALL scaffold linting, pre-commit hooks, and CI workflows

---

### Requirement: System Management Tool
The system SHALL provide a `jarvis_system` tool for bootstrap and restart operations.

#### Scenario: Bootstrap system
- **WHEN** `jarvis_system` is called with `action=bootstrap`
- **THEN** the system SHALL install MCPM, default servers, and start infrastructure

#### Scenario: Restart service
- **WHEN** `jarvis_system` is called with `action=restart`
- **THEN** the system SHALL gracefully restart Jarvis

#### Scenario: Restart infrastructure
- **WHEN** `jarvis_system` is called with `action=restart_infra`
- **THEN** the system SHALL reboot Docker infrastructure with health checks

---

### Requirement: Server Sharing Tool
The system SHALL provide a `jarvis_share` tool for remote server access.

#### Scenario: Start sharing
- **WHEN** `jarvis_share` is called with `action=start`, `name=<server>`, and optional port/auth params
- **THEN** the system SHALL expose the server via secure tunnel

#### Scenario: Stop sharing
- **WHEN** `jarvis_share` is called with `action=stop` and `name=<server>`
- **THEN** the system SHALL revoke tunnel access

#### Scenario: List shared servers
- **WHEN** `jarvis_share` is called with `action=list`
- **THEN** the system SHALL return active shares with URLs and status

---

### Requirement: Concise Tool Descriptions
The system SHALL use concise descriptions optimized for AI comprehension.

#### Scenario: Description length constraint
- **WHEN** tools are registered
- **THEN** each tool description SHALL be 60 characters or less

#### Scenario: Action enumeration in description
- **WHEN** a tool supports multiple actions
- **THEN** the description SHALL list available actions

---

### Requirement: Proper Tool Annotations
The system SHALL set meaningful annotation values for each tool.

#### Scenario: Read-only tools marked correctly
- **WHEN** a tool only reads data (status, list, info, search, suggest, diff, analyze)
- **THEN** the tool SHALL have `readOnlyHint: true` and `destructiveHint: false`

#### Scenario: Modifying tools marked correctly
- **WHEN** a tool modifies state (install, uninstall, edit, create, delete, restart, bootstrap, devops)
- **THEN** the tool SHALL have `readOnlyHint: false`

#### Scenario: Idempotent tools marked correctly
- **WHEN** a tool is idempotent (list, info, status, suggest, analyze)
- **THEN** the tool SHALL have `idempotentHint: true`
