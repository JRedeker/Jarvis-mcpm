# MCPM & Jarvis Configuration Strategy: The Streamable HTTP Daemon Architecture

**Version:** 4.1 (March 2025)
**Status:** Active Standard
**Core Change:** Shift from "SSE" to "Streamable HTTP" (MCP 2025-03-26 Spec).

## Overview

We have moved to a **Single Daemon Architecture**. Instead of every AI client spawning its own duplicate `mcpm` processes (wasteful, hard to manage), Jarvis orchestrates a single Docker container (`mcpm-daemon`) that hosts all tool profiles as **Streamable HTTP endpoints**.

Clients simply connect to these running HTTP endpoints.

## The 3-Layer Port Mapping

We still use the 3-Layer logic, but now layers map to **Ports**, not just profile names.

### Layer 1: PROJECT (Port 6276)
*   **Profile Name:** `p-<name>` (e.g., `p-pokeedge`)
*   **Port:** `6276` (Default Project Port)
*   **URL:** `http://localhost:6276/mcp`
*   **Purpose:** Workspace-specific tools (Search, Dev Tools).
*   **Switching:** To switch projects, you change which profile is running on port 6276 (or Jarvis creates a new port mapping). Currently, `p-pokeedge` is the primary active project.

### Layer 2: CAPABILITY (Ports 6278, 6279, ...)
*   **Profile Names:** `morph`, `qdrant`
*   **Ports:**
    *   `morph`: `6278` (`http://localhost:6278/mcp`)
    *   `qdrant`: `6279` (`http://localhost:6279/mcp`)
*   **Purpose:** Specialized, heavy capabilities (AI Refactoring, Vector DB).
*   **Opt-In:** Clients configure these URLs if they need the capability.

### Layer 3: ENVIRONMENT (Port 6277)
*   **Profile Name:** `memory`
*   **Port:** `6277` (`http://localhost:6277/mcp`)
*   **Purpose:** Persistent Memory (User preferences, historical context).
*   **Global:** Always active, always on this port.

## Configuration Rules

### 1. Use Streamable HTTP, Not Stdio or SSE
Clients **MUST** be configured to use `transport: streamable-http` (or rely on default HTTP detection) and the `url` field. Do not use `command: mcpm` for profiles anymore.

### 2. Jarvis is Stdio
Jarvis itself (the gateway) remains a local **stdio** command because it needs to manage the local Docker daemon and file system directly.

### 3. Short Names
Continue using short profile names (`p-pokeedge` not `project-pokeedge`) to keep tool names within API limits.

### 4. Timeouts for Remote MCP Servers
Remote MCP servers (HTTP endpoints) should have explicit timeouts configured to handle slow startup or network issues:

```json
{
  "p-pokeedge": {
    "type": "remote",
    "url": "http://localhost:6276/mcp",
    "timeout": 30000,
    "enabled": true
  }
}
```

**Recommended timeouts:**
| Profile Type | Timeout | Reason |
|--------------|---------|--------|
| Simple tools (fetch, time) | 10000ms | Fast response expected |
| AI-powered (qdrant, memory) | 30000ms | May need model warm-up |
| Heavy processing (morph) | 60000ms | Complex code operations |

### 5. Debugging Failed Profiles
Use `jarvis_diagnose` to debug profiles that fail to load:

```javascript
// Check if supervisor is running the profile
jarvis_diagnose({ action: "profile_health" })

// Get stderr logs from the subprocess
jarvis_diagnose({ action: "logs", profile: "qdrant" })

// Test MCP endpoint connectivity
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6279/mcp" })
```

## Current Port Map

| Profile | Port | URL | Servers |
| :--- | :--- | :--- | :--- |
| `p-pokeedge` | 6276 | `http://localhost:6276/mcp` | `brave-search`, `context7`, `firecrawl`, `time`, `fetch-mcp` |
| `memory` | 6277 | `http://localhost:6277/mcp` | `basic-memory`, `mem0-mcp` |
| `morph` | 6278 | `http://localhost:6278/mcp` | `morph-fast-apply` |
| `qdrant` | 6279 | `http://localhost:6279/mcp` | `mcp-server-qdrant` |

## Example Client Configurations

### OpenCode
OpenCode uses a different configuration format with native Jarvis support. See [docs/tech/opencode-integration.md](tech/opencode-integration.md) for full details.

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "jarvis": {
      "type": "local",
      "command": ["/home/user/dev/MCP/Jarvis/jarvis"],
      "enabled": true
    },
    "p-pokeedge": {
      "type": "remote",
      "url": "http://localhost:6276/mcp",
      "enabled": true
    },
    "memory": {
      "type": "remote",
      "url": "http://localhost:6277/mcp",
      "enabled": true
    },
    "morph": {
      "type": "remote",
      "url": "http://localhost:6278/mcp",
      "enabled": true
    }
  }
}
```

**Quick Setup:**
```javascript
// Use Jarvis to import starter config (v3.0 syntax)
jarvis_client({
  action: "import",
  client_name: "opencode"
})
```

### Claude Desktop / Claude Code / Gemini CLI
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/dev/MCP/Jarvis/jarvis",
      "args": []
    },
    "p-pokeedge": {
      "url": "http://localhost:6276/mcp",
      "transport": "streamable-http"
    },
    "memory": {
      "url": "http://localhost:6277/mcp",
      "transport": "streamable-http"
    },
    "morph": {
      "url": "http://localhost:6278/mcp",
      "transport": "streamable-http"
    },
    "qdrant": {
      "url": "http://localhost:6279/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### Kilo Code
```json
{
  "mcpServers": {
    "p-pokeedge": {
      "url": "http://localhost:6276/mcp",
      "transport": "streamable-http"
    },
    "memory": {
      "url": "http://localhost:6277/mcp",
      "transport": "streamable-http"
    },
    "morph": {
      "url": "http://localhost:6278/mcp",
      "transport": "streamable-http"
    }
  }
}
```
**Note:** Kilo Code doesn't strictly need Jarvis if it just wants tools, but Jarvis provides the system management capabilities.

## Migration Guide

**Note:** SSE transport has been fully removed as of v4.2 (MCP 2025-03-26 spec compliance). All profiles now use Streamable HTTP transport exclusively.

### If upgrading from SSE-based configuration:
1.  **Update Docker:** Ensure `mcpm-daemon` is running in `docker-compose.yml` (latest image).
2.  **Update Configs:** Replace `transport: sse` with `transport: streamable-http` and change `/sse` to `/mcp` in all endpoint URLs.
3.  **Restart:** Restart clients.
4.  **Verify:** Use `jarvis_check_status()` to confirm the daemon is healthy.
