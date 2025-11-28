# MCPM & Jarvis Configuration Strategy: The 3-Layer Stack

**Version:** 3.0 (November 2025)
**Status:** Active Standard

## Overview

To ensure scalability, conflict avoidance, and "agentic" intelligence, we have adopted a **3-Layer Profile Stacking Architecture**. This approach replaces monolithic configuration with modular, composable layers that balance auto-detection with manual opt-in capabilities.

## The 3 Layers

### Layer 1: PROJECT (Auto-Detected)
*   **Naming Convention:** `p-<name>` (e.g., `p-pokeedge`, `p-new`)
*   **Purpose:** Workspace-specific tools based on the project you're working in
*   **Detection:** Auto-detected by Jarvis based on current working directory
*   **Exclusivity:** An agent can only be in one PROJECT profile at a time
*   **Examples:**
    *   `p-pokeedge` → `brave-search`, `context7`, `firecrawl`, `time`
    *   `p-new` → `context7`, `firecrawl`, `time` (default for new projects)

### Layer 2: CAPABILITY (Manual Opt-In)
*   **Naming Convention:** Tool-specific names (e.g., `morph`, `qdrant`)
*   **Purpose:** Opt-in tool sets based on client capabilities or specific needs
*   **Detection:** NOT auto-detected - users manually add to their client configs
*   **Use Cases:**
    *   `morph` - AI refactoring for clients **without** built-in morph (e.g., Claude Desktop)
    *   `qdrant` - Vector search when needed for specific tasks
*   **Anti-pattern:** Do NOT create client-specific profiles like `c-codex` or `c-gemini`. Different client **instances** have different needs regardless of their type.

### Layer 3: ENVIRONMENT (Auto-Applied Global)
*   **Naming Convention:** `memory`, `testing-all-tools`
*   **Purpose:** Cross-cutting concerns that are always active or toggled on demand
*   **Detection:** Auto-applied by Jarvis
*   **Standard:** `memory` (Persistent storage: `basic-memory`, `mem0-mcp`)
*   **Conditional:** `testing-all-tools` (enabled when testing_mode=true)

## The Jarvis Logic ("The Brain")

Jarvis (the MCP Server) is configured **directly** in all clients (via binary path). It exposes a `suggest_profile` tool that implements the stacking logic:

```go
// Simplified logic for suggest_profile(testing_mode)
profiles = []

// 1. LAYER 1: PROJECT (Auto-detected from cwd)
if cwd contains "pokeedge" -> profiles.add("p-pokeedge")
else if cwd contains "codex" -> profiles.add("p-codex")
else -> profiles.add("p-new")  // Default

// 2. LAYER 2: CAPABILITY (Not auto-suggested - manual only)
// Users add these directly to their client configs as needed

// 3. LAYER 3: ENVIRONMENT (Auto-applied globals)
profiles.add("memory")  // Always active
if testing_mode -> profiles.add("testing-all-tools")

return profiles  // e.g., ["p-new", "memory"]
```

## Configuration Rules

### 1. Jarvis Must Be Direct
Never put `jarvis` inside an MCPM profile. Configure it directly in the client's config file pointing to the binary. This ensures management capabilities are available even if profiles break.

### 2. No Overlap
Tools should exist in exactly one layer to avoid duplication and version conflicts.

### 3. Shortened Server Names (Critical for OpenAI API)
**IMPORTANT:** To comply with OpenAI's 64-character limit on tool names, server names in client configs should be SHORT and NOT include `mcpm_profile_` prefix.

**Correct:**
```json
{
  "mcpServers": {
    "p-pokeedge": {
      "command": "mcpm",
      "args": ["profile", "run", "p-pokeedge"]
    }
  }
}
```

**Incorrect (causes tool name overflow):**
```json
{
  "mcpServers": {
    "mcpm_profile_p-pokeedge": {  // ❌ Too long!
      "command": "mcpm",
      "args": ["profile", "run", "p-pokeedge"]
    }
  }
}
```

Tool names are generated as: `mcp__{server_name}__{server}_{tool}`
Example: `mcp__p-pokeedge__brave-search_brave_web_search` (46 chars ✓)
Bad example: `mcp__mcpm_profile_p-pokeedge__brave-search_brave_web_search` (59 chars, but longer tools exceed 64)

## Current Profile Map

| Profile Name | Layer | Contents | Auto-Suggested? |
| :--- | :--- | :--- | :--- |
| `p-new` | 1 (Project) | `context7`, `firecrawl`, `time` | ✓ (default) |
| `p-pokeedge` | 1 (Project) | `brave-search`, `context7`, `firecrawl`, `time` | ✓ (if in pokeedge dir) |
| `morph` | 2 (Capability) | `morph-fast-apply` | ✗ (manual) |
| `qdrant` | 2 (Capability) | `mcp-server-qdrant` | ✗ (manual) |
| `memory` | 3 (Environment) | `basic-memory`, `mem0-mcp` | ✓ (always) |
| `testing-all-tools` | 3 (Environment) | **ALL** tools (for CI/CD) | ✓ (if testing=true) |

## Example Client Configurations

### Kilocode (Has Built-In Morph)
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/dev/MCP/Jarvis/jarvis",
      "args": []
    },
    "p-pokeedge": {
      "command": "mcpm",
      "args": ["profile", "run", "p-pokeedge"]
    },
    "memory": {
      "command": "mcpm",
      "args": ["profile", "run", "memory"]
    }
  }
}
```
**Note:** No `morph` profile - Kilocode already has morph built-in!

### Claude Desktop (Needs Morph)
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/dev/MCP/Jarvis/jarvis",
      "args": []
    },
    "p-pokeedge": {
      "command": "mcpm",
      "args": ["profile", "run", "p-pokeedge"]
    },
    "morph": {
      "command": "mcpm",
      "args": ["profile", "run", "morph"]
    },
    "memory": {
      "command": "mcpm",
      "args": ["profile", "run", "memory"]
    }
  }
}
```
**Note:** Added `morph` profile manually because Claude Desktop lacks native refactoring.

### Specialized Project with Vector Search
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/dev/MCP/Jarvis/jarvis",
      "args": []
    },
    "p-research": {
      "command": "mcpm",
      "args": ["profile", "run", "p-research"]
    },
    "qdrant": {
      "command": "mcpm",
      "args": ["profile", "run", "qdrant"]
    },
    "memory": {
      "command": "mcpm",
      "args": ["profile", "run", "memory"]
    }
  }
}
```
**Note:** Added `qdrant` profile for vector search capabilities.

## Client Configuration Locations

Standard configuration paths for wiring Jarvis into clients on Linux:
*   **Claude CLI:** `~/.claude.json`
*   **Claude VSCode / Desktop:** `~/.config/Claude/claude_desktop_config.json`
*   **Kilocode:** Check Kilocode documentation for config path

## Migration from v2.0

### Changes in v3.0 (November 2025)

1. **Removed Client Profiles**
   - Deleted: `c-codex`, `c-gemini`
   - Reason: Client identity ≠ capability requirements. What you need depends on the task, not the client type.

2. **Renamed PROJECT Profiles**
   - `project-pokeedge` → `p-pokeedge`
   - `project-new` → `p-new`
   - Reason: Shorter names to comply with OpenAI's 64-char tool name limit

3. **Created CAPABILITY Profiles**
   - Added: `morph` (AI refactoring)
   - Added: `qdrant` (vector search)
   - Reason: Opt-in based on needs, not client type

4. **Shortened Server Names**
   - Remove `mcpm_profile_` prefix from server names in client configs
   - Use just the profile name: `p-pokeedge`, `morph`, `memory`
   - Reason: Saves 13 characters per tool name, prevents OpenAI API errors

5. **Updated suggest_profile Logic**
   - Removed Layer 2 auto-detection (client adapters)
   - Now returns: `[project, memory]` or `[project, memory, testing]`
   - CAPABILITY profiles are manually configured, never auto-suggested

### Migration Steps

1. **Update profile names:**
   ```bash
   mcpm profile edit project-pokeedge --name p-pokeedge
   mcpm profile edit project-new --name p-new
   ```

2. **Create capability profiles:**
   ```bash
   mcpm profile create morph
   mcpm profile edit morph --add-server morph-fast-apply

   mcpm profile create qdrant
   mcpm profile edit qdrant --add-server mcp-server-qdrant
   ```

3. **Remove old client profiles:**
   ```bash
   mcpm profile rm c-codex --force
   mcpm profile rm c-gemini --force
   ```

4. **Update client configs** to use shortened server names (remove `mcpm_profile_` prefix)

5. **Manually add capability profiles** to client configs as needed (e.g., add `morph` to Claude Desktop but not Kilocode)

## Design Principles

*   **Auto-Detect What You Can:** PROJECT profiles are environment-aware
*   **Opt-In What You Should:** CAPABILITY profiles respect client differences
*   **Always Provide Essentials:** ENVIRONMENT profiles ensure core functionality
*   **Keep Names Short:** Comply with API limitations (64-char tool names)
*   **Separate Concerns:** Tools in one layer, configured once, no duplication
