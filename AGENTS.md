# Jarvis & MCPM Agent Instructions

**Current Date:** November 2025
**Version:** 2.1

## 🚨 Core Mandate: Use Jarvis Tools, Not Shell

You are an advanced AI agent. You must **NOT** use `run_shell_command` to execute `mcpm` or `jarvis` binaries directly unless explicitly instructed or debugging a tool failure.

**ALWAYS** use the provided MCP tools:
*   `manage_client(...)` -> NOT `mcpm client ...`
*   `manage_profile(...)` -> NOT `mcpm profile ...`
*   `install_server(...)` -> NOT `mcpm install ...`
*   `search_servers(...)` -> NOT `mcpm search ...`

These tools ensure proper state management, logging, and error handling within the Jarvis ecosystem.

## 📁 Configuration Locations

We use a 3-Layer Architecture. Jarvis needs to know where to write configurations.

### Standard Clients
Most clients (Claude Desktop, Cursor, Windsurf) use standard locations which Jarvis detects automatically.

### Custom/Specialized Clients (Kilo Code, Roo, Cline)
These are VS Code extensions and have unique storage paths. You must often **register** these paths with Jarvis.

**Kilo Code / Cline / Roo Paths (Linux/WSL):**
*   **Kilo Code:** `~/.vscode-server/data/User/globalStorage/kilocode.kilo-code/settings/mcp_settings.json`
*   **Cline:** `~/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`
*   **Roo Code:** `~/.vscode-server/data/User/globalStorage/rooveterinaryinc.roo-cline/settings/mcp_settings.json`

*(Note: On macOS, these are usually in `~/Library/Application Support/Code/User/globalStorage/...`)*

## 🛠️ Workflows

### 1. Setting Up a Client (The Right Way)

If a client (like Kilo Code) isn't detected or uses a custom path, register it first:

```javascript
// Tell Jarvis where Kilo Code lives
use_tool("manage_client", {
  "action": "config",
  "client_name": "cline", // Kilo Code is often compatible with Cline manager
  "config_path": "/home/user/.vscode-server/.../mcp_settings.json"
})
```

### 2. Applying Profiles (The 3-Layer Stack)

Don't add individual servers manually. Use profiles.

```javascript
// 1. Create/Update Environment Profile
use_tool("manage_profile", {
  "action": "create",
  "name": "project-pokeedge",
  "add_servers": "context7,fetch-mcp,time,brave-search"
})

// 2. Apply to Client
use_tool("manage_client", {
  "action": "edit",
  "client_name": "cline",
  "add_profile": "project-pokeedge,memory"
})
```

### 3. Debugging

If a tool fails or a client cannot connect, start with `check_status`:

```javascript
// Run comprehensive system diagnostics
use_tool("check_status", {})
```

If you must use the shell to debug, remember:
*   `export MCPM_NON_INTERACTIVE=true`
*   `export MCPM_FORCE=true`

## 🏗️ Architecture Reference

**Layer 1: Environment** (`project-xyz`)
*   Domain-specific tools (APIs, DBs, Search).
*   Exclusivity: One environment active at a time.

**Layer 2: Client Adapter** (`client-xyz`)
*   Client-specific tools (e.g., `morph-fast-apply` for Codex).

**Layer 3: Global** (`memory`)
*   Always-on capabilities (Memory, core utilities).

---
*Refer to `docs/CONFIGURATION_STRATEGY.md` for deep architectural details.*
