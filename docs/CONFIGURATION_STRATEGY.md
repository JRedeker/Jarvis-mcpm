# Configuration Strategy

**Version:** 5.1 (December 2025)
**Status:** Active Standard
**Core Change:** Composable Micro-Profiles with Enhanced Operations.

## Overview

Jarvis manages a set of "Composable Micro-Profiles" to provide a stable, scalable, and fault-tolerant environment for AI agents. Instead of every AI client spawning its own duplicate processes, Jarvis orchestrates a single Docker container (`mcpm-daemon`) that hosts all tool profiles as **Streamable HTTP endpoints**.

## Architectural Philosophy

Instead of loading a single monolithic "toolbox" profile, agents load a *stack* of specialized profiles. Each profile runs in its own process/container failure domain. If one profile crashes (e.g., due to a network timeout or Docker issue), the others remain available.

### The Standard Stack

| Profile | Port | URL | Purpose | Contents |
|:---|:---|:---|:---|:---|
| **`essentials`** | 6276 | `http://localhost:6276/mcp` | High-availability local tools. | `time`, `fetch-mcp` |
| **`memory`** | 6277 | `http://localhost:6277/mcp` | State persistence. | `basic-memory`, `mem0-mcp` |
| **`dev-core`** | 6278 | `http://localhost:6278/mcp` | Coding intelligence. | `context7` |
| **`data`** | 6279 | `http://localhost:6279/mcp` | Heavy storage/Vector DB. | `mcp-server-qdrant` |
| **`research`** | 6281 | `http://localhost:6281/mcp` | **High Risk** network/web tools. | `kagimcp`, `firecrawl`, `arxiv-mcp` |

*Note: `p-new` (Port 6280) is reserved for experimental/newly installed tools.*

## Configuring Clients

To configure an AI client (like OpenCode) to use this stack, add all profiles to its configuration.

**Via Jarvis:**
```javascript
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "essentials,memory,dev-core,research,data"
})
```

## Adding New Tools

When installing a new tool, assign it to the appropriate profile based on its characteristics:

1.  **Fast & Local?** -> `essentials`
2.  **Web/Network/Docker?** -> `research`
3.  **Coding Logic?** -> `dev-core`
4.  **Database?** -> `data`

```javascript
// Example: Adding a new web tool
jarvis_server({ action: "install", name: "google-maps" })
jarvis_profile({ action: "edit", name: "research", add_servers: "google-maps" })
```

## Troubleshooting

### Debugging Failed Profiles
Use `jarvis_diagnose` to debug profiles that fail to load:

```javascript
// Check if supervisor is running the profile
jarvis_diagnose({ action: "profile_health" })

// Get stderr logs from the subprocess
jarvis_diagnose({ action: "logs", profile: "research" })

// Test MCP endpoint connectivity
jarvis_diagnose({ action: "test_endpoint", endpoint: "http://localhost:6281/mcp" })
```

### Timeouts
Remote MCP servers (HTTP endpoints) should have explicit timeouts configured to handle slow startup or network issues:

- **Essentials/Dev-Core:** 10s
- **Memory/Data:** 30s
- **Research:** 60s (Heavy network usage)

## Configuration Synchronization

Jarvis maintains two configuration files that must stay in sync:

| File | Purpose | Used By |
|:-----|:--------|:--------|
| `~/.config/mcpm/servers.json` | Server definitions with `profile_tags` | Metadata/display |
| `~/.config/mcpm/profiles.json` | Profile → Server mappings | **Daemon (source of truth)** |

### Automatic Synchronization

As of v3.2, when you edit a profile via `jarvis_profile(action="edit")`, the server's `profile_tags` are automatically updated to stay in sync.

### Manual Audit & Fix

If configurations become out of sync (e.g., after manual edits), use the config_sync diagnostic:

```javascript
// Audit for mismatches
jarvis_diagnose({ action: "config_sync" })

// Auto-fix mismatches (updates servers.json to match profiles.json)
jarvis_diagnose({ action: "config_sync", auto_fix: true })
```

### API Endpoints

The MCPM API also exposes these endpoints:

- `GET /api/v1/audit` - Returns mismatch report
- `POST /api/v1/audit/fix` - Auto-fixes mismatches

## Configuration Backup & Restore (v5.1)

Jarvis supports exporting and importing MCPM configurations for backup and disaster recovery.

### Exporting Configuration

```javascript
// Export with secrets scrubbed (default)
jarvis_config({ action: "export", path: "mcpm-backup.json" })

// Export with secrets included (use with caution)
jarvis_config({ action: "export", path: "mcpm-backup.json", include_secrets: true })
```

The export includes:
- `servers.json` - All server definitions
- `profiles.json` - Profile to server mappings
- Version and timestamp metadata

**Security:** By default, environment variables containing sensitive patterns (API_KEY, SECRET, PASSWORD, TOKEN, CREDENTIALS) are replaced with `[SCRUBBED]`.

### Importing Configuration

```javascript
jarvis_config({ action: "import", path: "mcpm-backup.json" })
```

Import behavior:
1. Creates `.bak` backups of existing configurations
2. Warns if the backup contains scrubbed secrets
3. Writes new configurations to `~/.config/mcpm/`

**Note:** After importing, restart profiles to apply changes:
```javascript
jarvis_profile({ action: "restart" })
```

## Docker Operations (v5.1)

AI agents can now manage Docker containers without shell access:

```javascript
// Rebuild and restart services
jarvis_system({ action: "rebuild", no_cache: true })

// Stop/start without removing
jarvis_system({ action: "stop", service: "mcp-daemon" })
jarvis_system({ action: "start" })

// View container logs
jarvis_system({ action: "docker_logs", service: "mcp-daemon", lines: 100 })

// Check container status
jarvis_system({ action: "docker_status" })

// Build specific components
jarvis_system({ action: "build", component: "jarvis" })  // Go binary
jarvis_system({ action: "build", component: "mcpm-daemon" })  // Docker image
```
