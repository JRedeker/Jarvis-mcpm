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

**Current Date:** December 19, 2025
**Version:** 3.1 (The "Diagnostic" Edition)

## 🚨 Core Mandate: Use Jarvis Tools, Not Shell

You are an advanced AI agent. You must **NOT** use `run_shell_command` to execute `mcpm` or `jarvis` binaries directly unless explicitly instructed or debugging a tool failure.

**ALWAYS** use the provided MCP tools (consolidated in v3.0):

| Tool | Actions | Example |
|:-----|:--------|:--------|
| `jarvis_check_status` | (none) | System health check |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage | `jarvis_server(action="install", name="context7")` |
| `jarvis_profile` | list, create, edit, delete, suggest, restart | `jarvis_profile(action="list")` |
| `jarvis_client` | list, edit, import, config | `jarvis_client(action="edit", client_name="opencode", add_profile="memory")` |
| `jarvis_config` | get, set, list, migrate | `jarvis_config(action="list")` |
| `jarvis_project` | analyze, diff, devops | `jarvis_project(action="analyze")` |
| `jarvis_system` | bootstrap, restart, restart_infra | `jarvis_system(action="bootstrap")` |
| `jarvis_share` | start, stop, list | `jarvis_share(action="list")` |
| `jarvis_diagnose` | profile_health, test_endpoint, logs, full | `jarvis_diagnose(action="profile_health")` |

## 🧠 The 3-Layer Stack Philosophy

We do not manage monolithic configurations. We manage **Composable Profiles**.
See `docs/CONFIGURATION_STRATEGY.md` for the full architectural standard.

### 1. Layer 1: Environment (`project-<name>`)
*   **What:** The workspace context.
*   **Tools:** Domain-specific (Database, API, Search, Fetch).
*   **Example:** `project-pokeedge` contains `context7`, `fetch`, `morph-fast-apply`.

### 2. Layer 2: Client Adapter (`client-<name>`)
*   **What:** Capabilities specific to the AI Client (VS Code vs. Terminal).
*   **Tools:** Rendering aids, specific diff applicators (if not in Layer 1).

### 3. Layer 3: Global (`memory`, `testing`)
*   **What:** Always-on capabilities.
*   **Tools:** `basic-memory`, `mem0`, `qdrant`.

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
  add_profile: "project-pokeedge,memory" // Stack Layer 1 + Layer 3
})
```

### B. Intelligent Refactoring (`morph-fast-apply`)
The `morph-fast-apply` server is now a standard part of the `project-pokeedge` profile (Layer 1).
This tool allows you to make semantic edits without worrying about exact line numbers.

**Usage Pattern:**
1.  Check if available: `jarvis_server(action="list")`
2.  Use it: When the user asks for a complex refactor, prefer `morph-fast-apply` tools over raw file overwrites if safe to do so.

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
  add_profile: "jarvis,memory,p-pokeedge"
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
| 6276 | p-pokeedge profile |
| 6277 | memory profile |
| 6278 | morph profile |
| 6279 | qdrant profile |
| 6280 | p-new profile |

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
jarvis_diagnose({ action: "logs", profile: "qdrant" })

// Step 3: Test if MCP endpoint is responding correctly
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6279/mcp" })

// Step 4: Get comprehensive diagnostic report
jarvis_diagnose({ action: "full" })
```

### Legacy Debugging (if jarvis_diagnose unavailable)
1.  Run `jarvis_check_status()` for diagnostics (now includes API server health).
2.  Check API server: `curl http://localhost:6275/api/v1/health`
3.  If you must use shell: `export MCPM_NON_INTERACTIVE=true` and `export MCPM_FORCE=true`.

## 📋 Migration Guide (v2.x → v3.0)

**Breaking Change:** Tool names have been consolidated for context token efficiency.

### Old → New Tool Mapping

| Old Tool | New Tool + Action |
|:---------|:------------------|
| `check_status` | `jarvis_check_status` |
| `list_servers` | `jarvis_server(action="list")` |
| `server_info` | `jarvis_server(action="info")` |
| `install_server` | `jarvis_server(action="install")` |
| `uninstall_server` | `jarvis_server(action="uninstall")` |
| `search_servers` | `jarvis_server(action="search")` |
| `edit_server` | `jarvis_server(action="edit")` |
| `create_server` | `jarvis_server(action="create")` |
| `usage_stats` | `jarvis_server(action="usage")` |
| `manage_profile` | `jarvis_profile(action="list\|create\|edit\|delete")` |
| `suggest_profile` | `jarvis_profile(action="suggest")` |
| `restart_profiles` | `jarvis_profile(action="restart")` |
| `manage_client` | `jarvis_client(action="list\|edit\|import\|config")` |
| `manage_config` | `jarvis_config(action="get\|set\|list")` |
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

### Benefits
- **52% payload reduction**: ~5.3KB vs ~11KB (saves ~1,400 tokens per connection)
- **Cleaner namespace**: All tools prefixed with `jarvis_`
- **Action-based routing**: Easier to discover related operations

## 🆕 New in v3.1: Diagnostic Tools

The `jarvis_diagnose` tool enables AI agents to self-debug when MCP tools fail:

| Action | Purpose |
|:-------|:--------|
| `profile_health` | Check supervisor status for all MCP profiles |
| `test_endpoint` | Test MCP protocol on a specific endpoint |
| `logs` | Retrieve subprocess stderr logs for a profile |
| `full` | Comprehensive diagnostic report |

### Common Diagnostic Workflows

**Profile won't start:**
```javascript
jarvis_diagnose({ action: "logs", profile: "qdrant", lines: 100 })
// Look for Python errors, connection refused, config issues
```

**Tools missing from profile:**
```javascript
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6276/mcp" })
// Shows tool count and any errors from tools/list
```

**After config changes:**
```javascript
jarvis_profile({ action: "restart", profile: "p-pokeedge" })
jarvis_diagnose({ action: "profile_health" })
// Verify profile restarted successfully
```
