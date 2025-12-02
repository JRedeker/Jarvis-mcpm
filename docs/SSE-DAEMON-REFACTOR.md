# MCPM SSE Daemon Refactor - Status Summary

**Date:** 2025-11-28
**Status:** Completed (Released v4.0)

## Problem Statement

MCP servers were running in **stdio mode**, meaning each Claude Code session spawned its own duplicate processes. This caused:
- 6+ instances of basic-memory and mem0 running simultaneously
- No hot-reload capability (config changes required client restart)
- Resource waste and no centralized management
- The original Firecrawl 404 bug was caused by stale processes not picking up new API key config

## Solution: Single Docker Container with SSE

Created `mcpm-daemon` - a single Docker container running all MCP profiles as SSE services via supervisor.

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    mcp-daemon Container                          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │ p-pokeedge   │  │ memory       │  │ morph        │           │
│  │ --sse :6276  │  │ --sse :6277  │  │ --sse :6278  │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
│                                                                  │
│  Managed by: supervisord                                         │
└─────────────────────────────────────────────────────────────────┘
          │                    │                  │
          ▼                    ▼                  ▼
     Claude Code          Kilo Code           Cursor
     (connects via SSE URLs instead of spawning processes)
```

### Files Created

| File | Purpose |
|------|---------|
| `mcpm-daemon/Dockerfile` | Python 3.11 + MCPM + supervisor |
| `mcpm-daemon/entrypoint.sh` | Generates supervisor config, starts profiles |
| `docker-compose.yml` | Updated with mcpm-daemon service |
| `.env` | API keys for all MCP servers |
| `config-templates/claude-code-sse.json` | Template for SSE client config |
| `Jarvis/main.go` | Added `restart_profiles` tool (not yet built) |
| `Jarvis/tools.go` | Added `handleRestartProfiles` handler |

### Port Assignments

| Profile | Port | Servers |
|---------|------|---------|
| p-pokeedge | 6276 | context7, fetch-mcp, kagimcp, time, firecrawl |
| memory | 6277 | basic-memory, mem0-mcp |
| morph | 6278 | morph-fast-apply |

## Current State

### Working
- ✅ mcpm-daemon container builds and runs
- ✅ All 3 profiles start successfully via supervisor
- ✅ SSE endpoints accessible at `http://localhost:627X/sse/`
- ✅ Kilo Code config created at `~/.config/kilo/mcp.json`

### Issues Found During Testing
- ❌ `fetch-mcp` had placeholder path `{ABSOLUTE PATH TO FILE HERE}`
  - **Fix needed:** Update to `npx -y mcp-fetch-server`
  - Must use Jarvis (currently disconnected from session)

### Pending
- Jarvis MCP connection lost after `restart_service` call
- Need to rebuild Jarvis with new `restart_profiles` tool
- Full testing in Kilo Code incomplete

## Next Steps

1. **Reconnect to Jarvis** (restart Claude Code or new session)

2. **Fix fetch-mcp config via Jarvis:**
   ```
   edit_server(name="fetch-mcp", command="npx", args="-y mcp-fetch-server")
   ```

3. **Restart mcpm-daemon** to pick up config change:
   ```bash
   docker compose restart mcpm-daemon
   ```

4. **Rebuild Jarvis** with the new `restart_profiles` tool:
   ```bash
   cd Jarvis && go build -o jarvis .
   ```

5. **Test in Kilo Code** - verify all tools work via SSE

6. **Apply to Claude Code** - update `~/.claude.json` with SSE URLs

7. **Update documentation** - CLAUDE.md, CONFIGURATION_STRATEGY.md

## Key Commands

```bash
# Check daemon status
docker logs mcp-daemon

# Restart all profiles
docker compose restart mcpm-daemon

# Restart specific profile (after Jarvis rebuild)
docker exec mcp-daemon supervisorctl restart mcpm-p-pokeedge

# View profile logs
docker exec mcp-daemon cat /var/log/mcpm/p-pokeedge.out.log
```

## Config Locations

- MCPM servers: `~/.config/mcpm/servers.json`
- MCPM profiles: `~/.config/mcpm/profiles.json`
- Claude Code: `~/.claude.json`
- Kilo Code: `~/.config/kilo/mcp.json`
- API Keys: `/home/jrede/dev/MCP/.env`
