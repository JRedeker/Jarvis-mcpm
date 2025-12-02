# MCPM & Jarvis Configuration Strategy: The SSE Daemon Architecture

**Version:** 4.0 (November 2025)
**Status:** Active Standard
**Core Change:** Shift from "Stdio Process Per Client" to "Single SSE Daemon".

## Overview

We have moved to a **Single Daemon Architecture**. Instead of every AI client spawning its own duplicate `mcpm` processes (wasteful, hard to manage), Jarvis orchestrates a single Docker container (`mcpm-daemon`) that hosts all tool profiles as **SSE (Server-Sent Events) endpoints**.

Clients simply connect to these running HTTP endpoints.

## The 3-Layer Port Mapping

We still use the 3-Layer logic, but now layers map to **Ports**, not just profile names.

### Layer 1: PROJECT (Port 6276)
*   **Profile Name:** `p-<name>` (e.g., `p-pokeedge`)
*   **Port:** `6276` (Default Project Port)
*   **URL:** `http://localhost:6276/sse`
*   **Purpose:** Workspace-specific tools (Search, Dev Tools).
*   **Switching:** To switch projects, you change which profile is running on port 6276 (or Jarvis creates a new port mapping). Currently, `p-pokeedge` is the primary active project.

### Layer 2: CAPABILITY (Ports 6278, 6279, ...)
*   **Profile Names:** `morph`, `qdrant`
*   **Ports:**
    *   `morph`: `6278` (`http://localhost:6278/sse`)
    *   `qdrant`: `6279` (`http://localhost:6279/sse`)
*   **Purpose:** Specialized, heavy capabilities (AI Refactoring, Vector DB).
*   **Opt-In:** Clients configure these URLs if they need the capability.

### Layer 3: ENVIRONMENT (Port 6277)
*   **Profile Name:** `memory`
*   **Port:** `6277` (`http://localhost:6277/sse`)
*   **Purpose:** Persistent Memory (User preferences, historical context).
*   **Global:** Always active, always on this port.

## Configuration Rules

### 1. Use SSE, Not Stdio
Clients **MUST** be configured to use `transport: sse` and the `url` field. Do not use `command: mcpm` for profiles anymore.

### 2. Jarvis is Stdio
Jarvis itself (the gateway) remains a local **stdio** command because it needs to manage the local Docker daemon and file system directly.

### 3. Short Names
Continue using short profile names (`p-pokeedge` not `project-pokeedge`) to keep tool names within API limits.

## Current Port Map

| Profile | Port | URL | Servers |
| :--- | :--- | :--- | :--- |
| `p-pokeedge` | 6276 | `http://localhost:6276/sse` | `brave-search`, `context7`, `firecrawl`, `time`, `fetch-mcp` |
| `memory` | 6277 | `http://localhost:6277/sse` | `basic-memory`, `mem0-mcp` |
| `morph` | 6278 | `http://localhost:6278/sse` | `morph-fast-apply` |
| `qdrant` | 6279 | `http://localhost:6279/sse` | `mcp-server-qdrant` |

## Example Client Configurations

### Claude Desktop / Claude Code / Gemini CLI
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/dev/MCP/Jarvis/jarvis",
      "args": []
    },
    "p-pokeedge": {
      "url": "http://localhost:6276/sse",
      "transport": "sse"
    },
    "memory": {
      "url": "http://localhost:6277/sse",
      "transport": "sse"
    },
    "morph": {
      "url": "http://localhost:6278/sse",
      "transport": "sse"
    },
    "qdrant": {
      "url": "http://localhost:6279/sse",
      "transport": "sse"
    }
  }
}
```

### Kilo Code
```json
{
  "mcpServers": {
    "p-pokeedge": {
      "url": "http://localhost:6276/sse",
      "transport": "sse"
    },
    "memory": {
      "url": "http://localhost:6277/sse",
      "transport": "sse"
    },
    "morph": {
      "url": "http://localhost:6278/sse",
      "transport": "sse"
    }
  }
}
```
**Note:** Kilo Code doesn't strictly need Jarvis if it just wants tools, but Jarvis provides the system management capabilities.

## Migration Guide (v3.0 Stdio -> v4.0 SSE)

1.  **Update Docker:** Ensure `mcpm-daemon` is running in `docker-compose.yml`.
2.  **Update Configs:** Replace all `command: mcpm` entries in client configs with `url: http://...` entries.
3.  **Restart:** Restart clients. They will connect instantly.
4.  **Verify:** Use `jarvis.check_status()` to confirm the daemon is healthy.
