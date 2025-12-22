# OpenCode Integration Guide

> Updated for Jarvis v5.1 and Micro-Profile Architecture

This guide explains how to integrate Jarvis and the MCP micro-profile stack with OpenCode, the open-source AI coding agent.

## Overview

OpenCode is a modern AI coding agent developed by SST that supports the Model Context Protocol (MCP). Jarvis provides native support for OpenCode, allowing you to:

- Automatically detect OpenCode configuration files
- Add/remove MCP profiles with a single command
- Import a starter configuration with the full micro-profile stack
- Manage individual server configurations within profiles

## Configuration Format

OpenCode uses a different configuration format than Claude Desktop. Here's the structure:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "server-name": {
      "type": "local" | "remote",
      "command": ["path", "args"],  // for local (stdio)
      "url": "http://...",          // for remote (HTTP)
      "enabled": true,
      "environment": {},
      "headers": {}
    }
  }
}
```

### Transport Types

| Type | Use Case | Example |
|------|----------|---------|
| `local` | Direct binary execution (stdio) | Jarvis gateway |
| `remote` | HTTP endpoints | Profile servers on mcpm-daemon |

## Configuration Locations

OpenCode looks for configuration in this order:

1. `$OPENCODE_CONFIG` environment variable
2. `./opencode.json` in current project directory
3. `~/.config/opencode/opencode.json` (global)

## Quick Start

### Option 1: Import Starter Configuration

The fastest way to get started:

```javascript
jarvis_client({
  action: "import",
  client_name: "opencode"
})
```

This creates a configuration with the full micro-profile stack:
- `jarvis` (local stdio) - Infrastructure management gateway
- `essentials` (remote HTTP) - Fast, local utilities (time, fetch)
- `memory` (remote HTTP) - Persistent memory (basic-memory, mem0)
- `dev-core` (remote HTTP) - Coding intelligence (context7)
- `research` (remote HTTP) - Web research tools (kagi, firecrawl)
- `data` (remote HTTP) - Databases (qdrant, postgres)

### Option 2: Manual Configuration

Add profiles one at a time:

```javascript
// First, add Jarvis
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "jarvis"
})

// Then add the micro-profile stack
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "essentials,memory,dev-core,research,data"
})
```

### Option 3: Direct File Edit

Copy the template from `config-templates/opencode.json`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "jarvis": {
      "type": "local",
      "command": ["/path/to/MCP/Jarvis/jarvis"],
      "enabled": true
    },
    "essentials": {
      "type": "remote",
      "url": "http://localhost:6276/mcp",
      "enabled": true
    },
    "memory": {
      "type": "remote",
      "url": "http://localhost:6277/mcp",
      "enabled": true
    },
    "dev-core": {
      "type": "remote",
      "url": "http://localhost:6278/mcp",
      "enabled": true
    },
    "data": {
      "type": "remote",
      "url": "http://localhost:6279/mcp",
      "enabled": true
    },
    "research": {
      "type": "remote",
      "url": "http://localhost:6281/mcp",
      "enabled": true
    }
  }
}
```

## Port Reference (Micro-Profile Stack)

| Profile | Port | URL | Tools |
|---------|------|-----|-------|
| `essentials` | 6276 | `http://localhost:6276/mcp` | time, fetch-mcp |
| `memory` | 6277 | `http://localhost:6277/mcp` | basic-memory, mem0-mcp |
| `dev-core` | 6278 | `http://localhost:6278/mcp` | context7 |
| `data` | 6279 | `http://localhost:6279/mcp` | mcp-server-qdrant, postgres |
| `p-new` | 6280 | `http://localhost:6280/mcp` | Reserved for new profiles |
| `research` | 6281 | `http://localhost:6281/mcp` | kagimcp, firecrawl, arxiv-mcp |

### Why Micro-Profiles?

Instead of a monolithic `toolbox` profile, the micro-profile architecture provides:

1. **Failure Isolation**: A crash in `research` doesn't take down `memory`
2. **Selective Loading**: Enable only what you need (e.g., disable `data` if not using databases)
3. **Resource Control**: Heavy tools (qdrant) run in separate containers
4. **Easier Debugging**: Issues are isolated to specific domains

See [CONFIGURATION_STRATEGY.md](../CONFIGURATION_STRATEGY.md) for the full architectural philosophy.

## Prerequisites

Before using the MCP servers:

1. **Build Jarvis:**
   ```bash
   cd /path/to/MCP/Jarvis
   go build -o jarvis .
   ```

2. **Start Docker Infrastructure:**
   ```bash
   cd /path/to/MCP
   ./scripts/manage-mcp.sh start
   ```

3. **Verify Status:**
   ```bash
   ./scripts/manage-mcp.sh health
   ```

## Troubleshooting

### OpenCode not detecting Jarvis tools

1. Check if config exists:
   ```javascript
   jarvis_client({
     action: "config",
     client_name: "opencode"
   })
   ```

2. Verify Jarvis is built:
   ```bash
   ls -la /path/to/MCP/Jarvis/jarvis
   ```

3. Check Docker is running:
   ```bash
   docker ps | grep mcp
   ```

### HTTP profiles not connecting

1. Verify mcpm-daemon is healthy:
   ```bash
   docker logs mcp-daemon
   ```

2. Test endpoint directly:
   ```bash
   curl http://localhost:6276/mcp
   ```

3. Restart infrastructure:
   ```bash
   ./scripts/manage-mcp.sh restart
   ```

### Config path issues

Use explicit path:
```javascript
jarvis_client({
  action: "edit",
  client_name: "opencode",
  config_path: "/custom/path/opencode.json",
  add_profile: "memory"
})
```

### Profile tools not loading

Use the diagnostic tool:
```javascript
// Check profile health
jarvis_diagnose({ action: "profile_health" })

// Get logs from failing profile
jarvis_diagnose({ action: "logs", profile: "memory" })

// Test endpoint connectivity
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6277/mcp" })

// Check if configs are in sync
jarvis_diagnose({ action: "config_sync" })
```

## Comparison with Claude Desktop

| Feature | OpenCode | Claude Desktop |
|---------|----------|----------------|
| Config format | `mcp` object | `mcpServers` object |
| Transport key | `type` | `transport` |
| Local command | `["cmd", "arg"]` array | `command` string + `args` array |
| HTTP URL | `url` | `url` |
| Schema | `https://opencode.ai/config.json` | None |

## Related Documentation

- [OpenCode Docs: MCP Servers](https://opencode.ai/docs/mcp-servers/)
- [OpenCode Docs: Config](https://opencode.ai/docs/config/)
- [CONFIGURATION_STRATEGY.md](../CONFIGURATION_STRATEGY.md) - Micro-Profile Stack Philosophy
- [API_REFERENCE.md](../API_REFERENCE.md) - Complete Jarvis tool reference

---

*Updated for Jarvis v5.1 - December 2025*
