# OpenCode Integration Guide

This guide explains how to integrate Jarvis and the MCP profile stack with OpenCode, the open-source AI coding agent.

## Overview

OpenCode is a modern AI coding agent developed by SST that supports the Model Context Protocol (MCP). Jarvis provides native support for OpenCode, allowing you to:

- Automatically detect OpenCode configuration files
- Add/remove MCP profiles with a single command
- Import a starter configuration with all common profiles

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
use_tool("manage_client", {
  "action": "import",
  "client_name": "opencode"
})
```

This creates a configuration with:
- `jarvis` (local stdio) - Infrastructure management
- `p-pokeedge` (remote HTTP) - Project tools
- `memory` (remote HTTP) - Persistent memory
- `morph` (remote HTTP) - AI refactoring

### Option 2: Manual Configuration

Add profiles one at a time:

```javascript
// First, add Jarvis
use_tool("manage_client", {
  "action": "edit",
  "client_name": "opencode",
  "add_profile": "jarvis"
})

// Then add HTTP profiles
use_tool("manage_client", {
  "action": "edit",
  "client_name": "opencode",
  "add_profile": "memory,p-pokeedge"
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

## Port Reference

| Profile | Port | URL | Description |
|---------|------|-----|-------------|
| `p-pokeedge` | 6276 | `http://localhost:6276/mcp` | Project tools (search, fetch, context) |
| `memory` | 6277 | `http://localhost:6277/mcp` | Persistent memory |
| `morph` | 6278 | `http://localhost:6278/mcp` | AI code refactoring |
| `qdrant` | 6279 | `http://localhost:6279/mcp` | Vector database |
| `p-new` | 6280 | `http://localhost:6280/mcp` | New project template |

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
   use_tool("manage_client", {
     "action": "config",
     "client_name": "opencode"
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
use_tool("manage_client", {
  "action": "edit",
  "client_name": "opencode",
  "config_path": "/custom/path/opencode.json",
  "add_profile": "memory"
})
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
- [CONFIGURATION_STRATEGY.md](../CONFIGURATION_STRATEGY.md) - 3-Layer Stack Philosophy
