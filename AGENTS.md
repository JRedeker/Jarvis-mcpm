<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Jarvis & MCPM Agent Instructions

**Current Date:** December 22, 2025
**Version:** 5.1 (Enhanced Operations Edition)

## 🚨 Core Mandate: Use Jarvis Tools, Not Shell

You are an advanced AI agent. You must **NOT** use `run_shell_command` to execute `mcpm` or `jarvis` binaries directly unless explicitly instructed or debugging a tool failure.

**ALWAYS** use the provided MCP tools (consolidated in v3.0, enhanced in v5.1):

| Tool | Actions | Example |
|:-----|:--------|:--------|
| `jarvis_check_status` | (none) | System health check |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage | `jarvis_server(action="install", name="context7")` |
| `jarvis_profile` | list, create, edit, delete, suggest, restart | `jarvis_profile(action="list")` |
| `jarvis_client` | list, edit, import, config | `jarvis_client(action="edit", client_name="opencode", add_profile="memory")` |
| `jarvis_config` | get, set, list, migrate, **export**, **import** | `jarvis_config(action="export", path="backup.json")` |
| `jarvis_project` | analyze, diff, devops, **test** | `jarvis_project(action="test", verbose=true)` |
| `jarvis_system` | bootstrap, restart, restart_infra, **rebuild**, **stop**, **start**, **docker_logs**, **docker_status**, **build** | `jarvis_system(action="rebuild", no_cache=true)` |
| `jarvis_share` | start, stop, list | `jarvis_share(action="list")` |
| `jarvis_diagnose` | profile_health, test_endpoint, logs, full, config_sync | `jarvis_diagnose(action="profile_health")` |

## 🧠 The Micro-Profile Stack Philosophy

We do not manage monolithic configurations. We manage **Composable Micro-Profiles**.
See `docs/CONFIGURATION_STRATEGY.md` for the full architectural standard.

Instead of a single `toolbox` profile, we have domain-specific profiles that run in separate failure domains.

### 1. Essentials (`essentials`) - Port 6276
*   **What:** Fast, local utilities. Always on.
*   **Tools:** `time`, `fetch-mcp`.

### 2. Memory (`memory`) - Port 6277
*   **What:** Persistent storage.
*   **Tools:** `basic-memory`, `mem0-mcp`.

### 3. Developer Core (`dev-core`) - Port 6278
*   **What:** Coding intelligence.
*   **Tools:** `context7`.

### 4. Data (`data`) - Port 6279
*   **What:** Heavy databases.
*   **Tools:** `mcp-server-qdrant`, `postgres`.

### 5. Research (`research`) - Port 6281
*   **What:** High-risk network tools (Docker/Web).
*   **Tools:** `kagimcp`, `firecrawl`, `arxiv-mcp`.

## 🛠️ Key Operational Workflows

### A. Setting Up a New Client
If a client (like Kilo Code/Cline) isn't detected, register it:

```javascript
// 1. Tell Jarvis where the config file lives
jarvis_client({
  action: "config",
  client_name: "cline",
  config_path: "/home/user/.vscode-server/.../mcp_settings.json"
})

// 2. Apply the Profile Stack
jarvis_client({
  action: "edit",
  client_name: "cline",
  add_profile: "essentials,memory,dev-core,research,data"
})
```

### B. Using Context7 for Documentation
The `context7` server in `dev-core` provides library documentation lookup.

**Usage Pattern:**
1.  Check if available: `jarvis_server(action="list")`
2.  Use it: When working with libraries, use context7 to fetch documentation.

### C. Handling Output (The Presentation Layer)
Jarvis now returns formatted Markdown with emojis (✅/❌) and code blocks.
*   **Do not parse raw JSON manually** if the tool returns text.
*   **Present the output** to the user clearly.

## 📂 Reference Paths (Linux)

*   **OpenCode:** `~/.config/opencode/opencode.json` (global) or `./opencode.json` (project)
*   **Claude CLI:** `~/.claude.json`
*   **Claude Desktop/VSCode:** `~/.config/Claude/claude_desktop_config.json`

### OpenCode Quick Setup

OpenCode has native Jarvis support. To configure:

```javascript
// Import a starter configuration
jarvis_client({
  action: "import",
  client_name: "opencode"
})

// Or add profiles individually
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "essentials,memory,dev-core,research,data"
})
```

## 🌐 MCPM API Server (NEW)

Jarvis now supports two transports for communicating with MCPM:

### Transport Selection

| Transport | Environment Variable | Description |
|:----------|:--------------------|:------------|
| HTTP (default) | `JARVIS_MCPM_TRANSPORT=http` | Calls MCPM REST API |
| CLI (fallback) | `JARVIS_MCPM_TRANSPORT=cli` | Spawns `mcpm` subprocess |

The HTTP transport is preferred as it:
- Returns structured JSON responses (no text parsing)
- Faster (single long-running process vs subprocess per call)
- Type-safe and testable
- Falls back to CLI automatically if API server is unreachable

### API Server Configuration

| Variable | Default | Description |
|:---------|:--------|:------------|
| `MCPM_API_URL` | `http://localhost:6275` | API server URL |
| `MCPM_API_PORT` | `6275` | Port (for running locally) |

### Port Allocation

| Port | Service |
|:-----|:--------|
| 6275 | MCPM API Server |
| 6276 | essentials |
| 6277 | memory |
| 6278 | dev-core |
| 6279 | data |
| 6280 | p-new |
| 6281 | research |

### Starting the API Server

```bash
# Via CLI
mcpm serve --port 6275

# Via Docker (automatic when using docker-compose)
docker compose up mcpm-daemon
```

## 🚑 Debugging

If tools fail, use the `jarvis_diagnose` tool (NEW in v3.1):

```javascript
// Step 1: Check overall profile health
jarvis_diagnose({ action: "profile_health" })

// Step 2: If a specific profile is failing, get its logs
jarvis_diagnose({ action: "logs", profile: "research" })

// Step 3: Test if MCP endpoint is responding correctly
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6281/mcp" })

// Step 4: Get comprehensive diagnostic report
jarvis_diagnose({ action: "full" })
```

### Legacy Debugging (if jarvis_diagnose unavailable)
1.  Run `jarvis_check_status()` for diagnostics (now includes API server health).
2.  Check API server: `curl http://localhost:6275/api/v1/health`
3.  If you must use shell: `export MCPM_NON_INTERACTIVE=true` and `export MCPM_FORCE=true`.

## 🆕 New in v3.1: Diagnostic Tools

The `jarvis_diagnose` tool enables AI agents to self-debug when MCP tools fail:

| Action | Purpose |
|:-------|:--------|
| `profile_health` | Check supervisor status for all MCP profiles |
| `test_endpoint` | Test MCP protocol on a specific endpoint |
| `logs` | Retrieve subprocess stderr logs for a profile |
| `full` | Comprehensive diagnostic report |
| `config_sync` | Audit/fix mismatches between servers.json and profiles.json |

### Common Diagnostic Workflows

**Profile won't start:**
```javascript
jarvis_diagnose({ action: "logs", profile: "research", lines: 100 })
// Look for Python errors, connection refused, config issues
```

**Tools missing from profile:**
```javascript
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6276/mcp" })
// Shows tool count and any errors from tools/list
```

**After config changes:**
```javascript
jarvis_profile({ action: "restart", profile: "research" })
jarvis_diagnose({ action: "profile_health" })
// Verify profile restarted successfully
```

**Config out of sync (profile_tags vs profiles.json):**
```javascript
// Audit for mismatches
jarvis_diagnose({ action: "config_sync" })

// Auto-fix mismatches
jarvis_diagnose({ action: "config_sync", auto_fix: true })
```

## 🆕 New in v5.1: Enhanced Operations

### Docker Control (jarvis_system)

New actions for granular Docker management:

| Action | Purpose | Example |
|:-------|:--------|:--------|
| `rebuild` | Build & restart services | `jarvis_system(action="rebuild", no_cache=true)` |
| `stop` | Stop without removing | `jarvis_system(action="stop", service="mcp-daemon")` |
| `start` | Start stopped services | `jarvis_system(action="start")` |
| `docker_logs` | Get container logs | `jarvis_system(action="docker_logs", service="mcp-daemon", lines=100)` |
| `docker_status` | Detailed container status | `jarvis_system(action="docker_status")` |
| `build` | Build components selectively | `jarvis_system(action="build", component="jarvis")` |

### Test Runner (jarvis_project)

Run tests directly from Jarvis with auto-detection:

```javascript
// Auto-detect project type and run tests
jarvis_project({ action: "test" })

// With verbose output
jarvis_project({ action: "test", verbose: true })

// Specific package
jarvis_project({ action: "test", package: "./handlers/..." })

// Override project type
jarvis_project({ action: "test", project_type: "go" })
```

Supported project types: go, python, node, typescript

### Config Backup (jarvis_config)

Export and import MCPM configurations:

```javascript
// Export (scrubs secrets by default)
jarvis_config({ action: "export", path: "backup.json" })

// Export with secrets
jarvis_config({ action: "export", path: "backup.json", include_secrets: true })

// Import (creates backup of existing config)
jarvis_config({ action: "import", path: "backup.json" })
```

### Enhanced Logging (jarvis_diagnose)

Logs now use `docker compose logs` instead of supervisorctl for better output:

```javascript
// Get logs from mcp-daemon
jarvis_diagnose({ action: "logs", lines: 100 })

// Filter logs for a specific profile
jarvis_diagnose({ action: "logs", profile: "research", lines: 50 })
```
