# Jarvis API Reference (v3.1)

> Updated for consolidated tool architecture - 24 tools → 9 tools

This document provides a complete reference for all Jarvis MCP tools and the MCPM REST API.

---

## Overview

Jarvis v3.1 consolidates 24 tools into 9 action-based tools for **52% context token reduction** (~1,400 tokens saved per connection).

### Tool Summary

| Tool | Actions | Description |
|:-----|:--------|:------------|
| `jarvis_check_status` | - | System health diagnostics |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage | Server management |
| `jarvis_profile` | list, create, edit, delete, suggest, restart | Profile management |
| `jarvis_client` | list, edit, import, config | AI client configuration |
| `jarvis_config` | get, set, list, migrate | MCPM configuration |
| `jarvis_project` | analyze, diff, devops | Project analysis & DevOps |
| `jarvis_system` | bootstrap, restart, restart_infra | System operations |
| `jarvis_share` | start, stop, list | Server sharing |
| `jarvis_diagnose` | profile_health, test_endpoint, logs, full | MCP profile debugging |

---

## MCPM REST API

The MCPM API Server provides a REST interface for programmatic access to MCPM functionality.

**Base URL:** `http://localhost:6275/api/v1`

### Health Check

| Method | Endpoint | Description |
|:-------|:---------|:------------|
| GET | `/health` | System health check |

### Server Management

| Method | Endpoint | Description |
|:-------|:---------|:------------|
| GET | `/servers` | List all servers |
| GET | `/servers/:name` | Get server info |
| POST | `/servers/:name/install` | Install server |
| DELETE | `/servers/:name` | Uninstall server |
| POST | `/servers` | Create custom server |
| PUT | `/servers/:name` | Edit server config |
| GET | `/servers/search?q=query` | Search servers |

### Profile Management

| Method | Endpoint | Description |
|:-------|:---------|:------------|
| GET | `/profiles` | List all profiles |
| GET | `/profiles/:name` | Get profile details |
| POST | `/profiles` | Create profile |
| PUT | `/profiles/:name` | Edit profile |
| DELETE | `/profiles/:name` | Delete profile |

### Client Management

| Method | Endpoint | Description |
|:-------|:---------|:------------|
| GET | `/clients` | List all clients |
| GET | `/clients/:name` | Get client info |
| PUT | `/clients/:name` | Edit client config |

### System Operations

| Method | Endpoint | Description |
|:-------|:---------|:------------|
| GET | `/search?q=query` | Search servers (alias) |
| GET | `/usage` | Usage statistics |
| POST | `/migrate` | Migrate configuration |

### Response Format

All responses follow a consistent JSON structure:

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

Error responses:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message"
  }
}
```

### Configuration

| Environment Variable | Default | Description |
|:---------------------|:--------|:------------|
| `MCPM_API_PORT` | `6275` | API server port |
| `MCPM_API_HOST` | `0.0.0.0` | API server bind address |
| `JARVIS_MCPM_TRANSPORT` | `http` | Transport: `http` or `cli` |
| `MCPM_API_URL` | `http://localhost:6275` | API server URL (for clients) |

---

## MCP Tools Reference

### `jarvis_check_status`

Comprehensive system health check for MCPM, Docker, and all services. Validates Node.js, Python, dependencies, running containers, and HTTP endpoints. Returns actionable fix suggestions.

**Parameters:** None

**Example:**
```javascript
jarvis_check_status()
```

**Returns:**
```markdown
## System Status Report
✅ Node.js: v20.10.0
✅ MCPM CLI: installed
✅ Docker: running

### Infrastructure
✅ PostgreSQL: healthy (port 5432)
✅ Qdrant: healthy (port 6333)
✅ MCPM Daemon: healthy (3 profiles active)

🚀 **ALL SYSTEMS GO!** Jarvis is ready to assist.
```

---

### `jarvis_server`

Manage MCP servers: list, info, install, uninstall, search, edit, create, usage.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `list`, `info`, `install`, `uninstall`, `search`, `edit`, `create`, `usage` |
| `name` | string | For most | Server name |
| `query` | string | For search | Search query |
| `type` | string | For create | Transport: `stdio` or `streamable-http` |
| `command` | string | No | Command to run server |
| `args` | string | No | Arguments (space-separated) |
| `env` | string | No | Environment variables (KEY=val, comma-separated) |
| `url` | string | No | URL for HTTP servers |
| `headers` | string | No | HTTP headers (KEY=val, comma-separated) |

**Examples:**

```javascript
// List all installed servers
jarvis_server({ action: "list" })

// Get server details
jarvis_server({ action: "info", name: "context7" })

// Install a server
jarvis_server({ action: "install", name: "brave-search" })

// Uninstall a server
jarvis_server({ action: "uninstall", name: "old-server" })

// Search registry
jarvis_server({ action: "search", query: "documentation" })

// Edit server config
jarvis_server({ action: "edit", name: "my-server", env: "API_KEY=xxx" })

// Create custom server (stdio)
jarvis_server({
  action: "create",
  name: "my-custom-server",
  type: "stdio",
  command: "python",
  args: "-m my_server"
})

// Create custom server (HTTP)
jarvis_server({
  action: "create",
  name: "remote-server",
  type: "streamable-http",
  url: "http://localhost:8080/mcp"
})

// View usage statistics
jarvis_server({ action: "usage" })
```

---

### `jarvis_profile`

Manage MCPM profiles: list, create, edit, delete, suggest, restart.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `list`, `create`, `edit`, `delete`, `suggest`, `restart` |
| `name` | string | For most | Profile name |
| `new_name` | string | No | New name when renaming (for edit) |
| `add_servers` | string | No | Servers to add, comma-separated (for edit) |
| `remove_servers` | string | No | Servers to remove, comma-separated (for edit) |
| `profile` | string | No | Specific profile to restart (for restart) |
| `testing` | boolean | No | Include testing profile (for suggest) |

**Examples:**

```javascript
// List all profiles
jarvis_profile({ action: "list" })

// Create a new profile
jarvis_profile({ action: "create", name: "project-frontend" })

// Add servers to profile
jarvis_profile({
  action: "edit",
  name: "project-frontend",
  add_servers: "brave-search,context7"
})

// Remove servers from profile
jarvis_profile({
  action: "edit",
  name: "project-frontend",
  remove_servers: "old-server"
})

// Rename a profile
jarvis_profile({
  action: "edit",
  name: "old-name",
  new_name: "new-name"
})

// Delete a profile
jarvis_profile({ action: "delete", name: "unused-profile" })

// Get profile suggestions for current directory
jarvis_profile({ action: "suggest" })

// Get suggestions including testing tools
jarvis_profile({ action: "suggest", testing: true })

// Restart all profiles
jarvis_profile({ action: "restart" })

// Restart specific profile
jarvis_profile({ action: "restart", profile: "toolbox" })
```

---

### `jarvis_client`

Configure MCP clients: list, edit, import, config.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `list`, `edit`, `import`, `config` |
| `client_name` | string | For edit/import/config | Client: `opencode`, `claude-code`, `claude-desktop` |
| `add_server` | string | No | Server to add (for edit) |
| `remove_server` | string | No | Server to remove (for edit) |
| `add_profile` | string | No | Profile to add (for edit) |
| `remove_profile` | string | No | Profile to remove (for edit) |
| `config_path` | string | No | Config file path (for config) |

**Examples:**

```javascript
// List all detected clients
jarvis_client({ action: "list" })

// Add profiles to OpenCode
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "jarvis,memory,toolbox"
})

// Remove a profile from client
jarvis_client({
  action: "edit",
  client_name: "opencode",
  remove_profile: "old-profile"
})

// Import template configuration
jarvis_client({ action: "import", client_name: "opencode" })

// Get/set config path
jarvis_client({ action: "config", client_name: "opencode" })
jarvis_client({
  action: "config",
  client_name: "cline",
  config_path: "/home/user/.vscode-server/.../mcp_settings.json"
})
```

---

### `jarvis_config`

Manage MCPM config: get, set, list, migrate.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `get`, `set`, `list`, `migrate` |
| `key` | string | For get/set | Config key |
| `value` | string | For set | Value to set |

**Examples:**

```javascript
// List all config values
jarvis_config({ action: "list" })

// Get a specific config value
jarvis_config({ action: "get", key: "default_profile" })

// Set a config value
jarvis_config({ action: "set", key: "default_profile", value: "memory" })

// Migrate config to latest format
jarvis_config({ action: "migrate" })
```

---

### `jarvis_project`

Project tools: analyze, diff, devops.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `analyze`, `diff`, `devops` |
| `staged` | boolean | No | Show only staged changes (for diff) |
| `project_type` | string | No | Override: `python`, `go`, `node`, `typescript` (for devops) |
| `force` | boolean | No | Overwrite existing configs (for devops) |
| `enable_ai_review` | boolean | No | Add PR Agent workflow (for devops) |

**Examples:**

```javascript
// Analyze current project
jarvis_project({ action: "analyze" })

// Returns:
{
  "path": "/home/user/my-project",
  "languages": ["go", "javascript"],
  "frameworks": [],
  "configs": {
    "has_git": true,
    "has_pre_commit": false,
    "has_github_workflows": true,
    "has_pr_agent": false,
    "has_dependabot": true,
    "has_gitleaks": false
  },
  "key_files": ["go.mod", "package.json"]
}

// Get git diff for review
jarvis_project({ action: "diff" })

// Get only staged changes
jarvis_project({ action: "diff", staged: true })

// Apply DevOps stack
jarvis_project({ action: "devops" })

// Apply with specific project type
jarvis_project({ action: "devops", project_type: "python" })

// Force overwrite existing configs
jarvis_project({ action: "devops", force: true })

// Include AI review workflow
jarvis_project({ action: "devops", enable_ai_review: true })
```

---

### `jarvis_system`

System ops: bootstrap, restart, restart_infra.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `bootstrap`, `restart`, `restart_infra` |

**Examples:**

```javascript
// Complete system initialization
jarvis_system({ action: "bootstrap" })
// Installs MCPM, starts Docker, installs default servers

// Restart Jarvis service
jarvis_system({ action: "restart" })

// Restart Docker infrastructure (PostgreSQL, Qdrant)
jarvis_system({ action: "restart_infra" })
```

---

### `jarvis_share`

Share MCP servers: start, stop, list.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `start`, `stop`, `list` |
| `name` | string | For start/stop | Server name |
| `port` | string | No | Port for shared server (for start) |
| `no_auth` | boolean | No | Disable authentication (for start) |

**Examples:**

```javascript
// List active shares
jarvis_share({ action: "list" })

// Start sharing a server
jarvis_share({ action: "start", name: "context7" })

// Share on specific port without auth
jarvis_share({ action: "start", name: "my-server", port: "8080", no_auth: true })

// Stop sharing
jarvis_share({ action: "stop", name: "context7" })
```

---

### `jarvis_diagnose`

**NEW in v3.1:** Debug MCP profile issues: profile_health, test_endpoint, logs, full.

**Parameters:**

| Name | Type | Required | Description |
|:-----|:-----|:---------|:------------|
| `action` | string | Yes | Operation: `profile_health`, `test_endpoint`, `logs`, `full` |
| `profile` | string | For logs | Profile name to get logs for |
| `endpoint` | string | For test_endpoint | MCP endpoint URL to test |

**Examples:**

```javascript
// Check overall profile health (supervisor status)
jarvis_diagnose({ action: "profile_health" })
// Returns: Status of all supervised profiles (running/stopped/failed)

// Test if a specific MCP endpoint is responding
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6279/mcp" })
// Returns: MCP protocol test results including available tools

// Get stderr logs from a profile's subprocess
jarvis_diagnose({ action: "logs", profile: "qdrant" })
// Returns: Recent stderr output for debugging startup failures

// Comprehensive diagnostic report
jarvis_diagnose({ action: "full" })
// Returns: Combined profile health, endpoint tests, and recommendations
```

**Use Cases:**

| Symptom | Action | What It Reveals |
|---------|--------|-----------------|
| "qdrant mcp failed to get tools" | `profile_health` | Shows if supervisor has the profile running |
| "Connection refused on port 6279" | `test_endpoint` | Tests MCP protocol handshake |
| "Profile starts then crashes" | `logs` | Shows stderr from failed subprocess |
| "Everything seems broken" | `full` | Complete diagnostic for support requests |

**Common Workflow:**

```javascript
// Step 1: Check if profiles are running
jarvis_diagnose({ action: "profile_health" })

// Step 2: If a profile is failing, get its logs
jarvis_diagnose({ action: "logs", profile: "qdrant" })

// Step 3: Test if the endpoint responds correctly
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6279/mcp" })

// Step 4: Get comprehensive report for debugging
jarvis_diagnose({ action: "full" })
```

---

## Migration Guide (v2.x → v3.0)

### Old → New Tool Mapping

| Old Tool (v2.x) | New Tool (v3.0) |
|:----------------|:----------------|
| `check_status` | `jarvis_check_status` |
| `list_servers` | `jarvis_server(action="list")` |
| `server_info` | `jarvis_server(action="info")` |
| `install_server` | `jarvis_server(action="install")` |
| `uninstall_server` | `jarvis_server(action="uninstall")` |
| `search_servers` | `jarvis_server(action="search")` |
| `edit_server` | `jarvis_server(action="edit")` |
| `create_server` | `jarvis_server(action="create")` |
| `usage_stats` | `jarvis_server(action="usage")` |
| `manage_profile(action="ls")` | `jarvis_profile(action="list")` |
| `manage_profile(action="create")` | `jarvis_profile(action="create")` |
| `manage_profile(action="edit")` | `jarvis_profile(action="edit")` |
| `manage_profile(action="delete")` | `jarvis_profile(action="delete")` |
| `suggest_profile` | `jarvis_profile(action="suggest")` |
| `restart_profiles` | `jarvis_profile(action="restart")` |
| `manage_client(action="ls")` | `jarvis_client(action="list")` |
| `manage_client(action="edit")` | `jarvis_client(action="edit")` |
| `manage_client(action="import")` | `jarvis_client(action="import")` |
| `manage_client(action="config")` | `jarvis_client(action="config")` |
| `manage_config(action="get")` | `jarvis_config(action="get")` |
| `manage_config(action="set")` | `jarvis_config(action="set")` |
| `manage_config(action="ls")` | `jarvis_config(action="list")` |
| `migrate_config` | `jarvis_config(action="migrate")` |
| `analyze_project` | `jarvis_project(action="analyze")` |
| `fetch_diff_context` | `jarvis_project(action="diff")` |
| `apply_devops_stack` | `jarvis_project(action="devops")` |
| `bootstrap_system` | `jarvis_system(action="bootstrap")` |
| `restart_service` | `jarvis_system(action="restart")` |
| `restart_infrastructure` | `jarvis_system(action="restart_infra")` |
| `share_server` | `jarvis_share(action="start")` |
| `stop_sharing_server` | `jarvis_share(action="stop")` |
| `list_shared_servers` | `jarvis_share(action="list")` |

---

## Tool Categories Summary

| Category | Tool | Actions |
|:---------|:-----|:--------|
| System Health | `jarvis_check_status` | - |
| Server Management | `jarvis_server` | list, info, install, uninstall, search, edit, create, usage |
| Profile Management | `jarvis_profile` | list, create, edit, delete, suggest, restart |
| Client Configuration | `jarvis_client` | list, edit, import, config |
| MCPM Configuration | `jarvis_config` | get, set, list, migrate |
| Project Tools | `jarvis_project` | analyze, diff, devops |
| System Operations | `jarvis_system` | bootstrap, restart, restart_infra |
| Server Sharing | `jarvis_share` | start, stop, list |
| Diagnostics | `jarvis_diagnose` | profile_health, test_endpoint, logs, full |

---

*Updated for Jarvis v3.1 - Diagnostics Edition*
