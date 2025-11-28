# Jarvis & MCPM Agent Instructions

**Current Date:** November 28, 2025
**Version:** 2.2 (The "Smart Stack" Edition)

## 🚨 Core Mandate: Use Jarvis Tools, Not Shell

You are an advanced AI agent. You must **NOT** use `run_shell_command` to execute `mcpm` or `jarvis` binaries directly unless explicitly instructed or debugging a tool failure.

**ALWAYS** use the provided MCP tools:
*   `manage_client(...)` -> NOT `mcpm client ...`
*   `manage_profile(...)` -> NOT `mcpm profile ...`
*   `install_server(...)` -> NOT `mcpm install ...`
*   `search_servers(...)` -> NOT `mcpm search ...`

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
use_tool("manage_client", {
  "action": "config",
  "client_name": "cline",
  "config_path": "/home/user/.vscode-server/.../mcp_settings.json"
})

// 2. Apply the Profile Stack
use_tool("manage_client", {
  "action": "edit",
  "client_name": "cline",
  "add_profile": "project-pokeedge,memory" // Stack Layer 1 + Layer 3
})
```

### B. Intelligent Refactoring (`morph-fast-apply`)
The `morph-fast-apply` server is now a standard part of the `project-pokeedge` profile (Layer 1).
This tool allows you to make semantic edits without worrying about exact line numbers.

**Usage Pattern:**
1.  Check if available: `list_servers()`
2.  Use it: When the user asks for a complex refactor, prefer `morph-fast-apply` tools over raw file overwrites if safe to do so.

### C. Handling Output (The Presentation Layer)
Jarvis now returns formatted Markdown with emojis (✅/❌) and code blocks.
*   **Do not parse raw JSON manually** if the tool returns text.
*   **Present the output** to the user clearly.

## 📂 Reference Paths (Linux)

*   **Claude CLI:** `~/.claude.json`
*   **Claude Desktop/VSCode:** `~/.config/Claude/claude_desktop_config.json`

## 🚑 Debugging

If tools fail:
1.  Run `check_status()` for diagnostics.
2.  If you must use shell: `export MCPM_NON_INTERACTIVE=true` and `export MCPM_FORCE=true`.
