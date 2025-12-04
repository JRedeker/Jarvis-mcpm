# Jarvis API Reference

> Auto-generated from tool definitions in `handlers/server.go`

This document provides a complete reference for all Jarvis MCP tools.

---

## System Management

### `check_status`

Comprehensive system health check for MCPM, Docker, and all services. Validates Node.js, Python, dependencies, running containers, and SSE endpoints. Returns actionable fix suggestions for any issues found.

---


## Server Management

### `list_servers`

Shows all installed MCP servers with their status, transport type, and profile associations. Use this to inventory available tools before making changes.

---

### `server_info`

Detailed information about a specific server including command, args, environment variables, installation source, and usage statistics. Essential before editing or troubleshooting a server.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |

---

### `install_server`

Installs an MCP server from the registry with automatic dependency resolution. Validates the server exists before installing and suggests alternatives for typos.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |

---

### `uninstall_server`

Removes an installed MCP server and cleans up its configuration. Warns about profile impact before removal.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |

---

### `search_servers`

Search the MCP server registry by keyword, category, or capability. Returns matching servers with descriptions and installation commands.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `query` | string |

---

### `edit_server`

Modify an installed server's configuration including command, arguments, environment variables, and URL.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |
| `command` | string |
| `args` | string |
| `env` | string |
| `url` | string |
| `headers` | string |

---

### `create_server`

Register a new custom MCP server that isn't in the registry.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |
| `type` | string |
| `command` | string |
| `args` | string |
| `env` | string |
| `url` | string |
| `headers` | string |

---

### `usage_stats`

Shows tool usage statistics across all servers and profiles.

---


## Profile Management

### `manage_profile`

Create, edit, delete, or list MCPM profiles. Profiles group servers for specific use cases.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `action` | string |
| `name` | string |
| `new_name` | string |
| `add_servers` | string |
| `remove_servers` | string |

---

### `suggest_profile`

Intelligently determines optimal MCPM profile stack by analyzing working directory, client type, and mode.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `testing` | boolean |

---

### `restart_profiles`

Restarts the MCPM daemon container to reload all MCP profiles with updated configurations.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `profile` | string |

---


## Client Management

### `manage_client`

Configure MCP client applications (Claude Code, Claude Desktop, etc.) with profiles and servers.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `action` | string |
| `client_name` | string |
| `add_server` | string |
| `remove_server` | string |
| `add_profile` | string |
| `remove_profile` | string |
| `config_path` | string |

---


## Configuration

### `manage_config`

Get, set, or list MCPM configuration values.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `action` | string |
| `key` | string |
| `value` | string |

---

### `migrate_config`

Upgrades MCPM configuration to the latest format with automatic backup.

---


## Project Analysis

### `analyze_project`

Analyzes the current project to detect languages, frameworks, and existing DevOps configurations. Returns JSON report.

---

### `fetch_diff_context`

Returns git status and diff for self-review before commits. Helps catch issues before pushing.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `staged` | boolean |

---

### `apply_devops_stack`

Scaffolds projects with linting, pre-commit hooks, and CI/CD workflows based on detected project type.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `project_type` | string |
| `force` | boolean |
| `enable_ai_review` | boolean |

---


## System Bootstrap

### `bootstrap_system`

Complete system initialization: installs MCPM, sets up default servers (context7, brave-search, github), and starts Docker infrastructure (PostgreSQL, Qdrant). One command to get fully operational.

---

### `restart_service`

Gracefully restarts Jarvis to apply configuration changes or resolve stuck states. Automatically saves state and reconnects active sessions.

---

### `restart_infrastructure`

Safely reboots Docker infrastructure (PostgreSQL, Qdrant) with health checks and automatic reconnection. Resolves database connection issues, clears stale locks, and ensures all services are healthy.

---


## Server Sharing

### `share_server`

Exposes local MCP servers via secure tunnels with optional authentication. Enables remote teams to access your tools without VPN or port forwarding.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |
| `port` | string |
| `no_auth` | boolean |

---

### `stop_sharing_server`

Revokes tunnel access and terminates shared server sessions. Immediately disconnects all remote clients.

**Parameters:**

| Name | Type |
|:-----|:-----|
| `name` | string |

---

### `list_shared_servers`

Shows all active server shares with tunnel URLs, authentication status, and connected clients.

---


## Tool Categories Summary

| Category | Description |
|:---------|:------------|
| System Management | Health checks, bootstrapping, service management |
| Server Management | Install, uninstall, search, configure MCP servers |
| Profile Management | Create and manage composable tool profiles |
| Client Management | Configure AI clients (Claude, Cursor, etc.) |
| Configuration | Manage global settings and migrations |
| Project Analysis | Analyze projects and apply DevOps stacks |
| System Bootstrap | System initialization and infrastructure |
| Server Sharing | Share servers via secure tunnels |

---

*Generated by `scripts/generate-api-docs.sh`*
