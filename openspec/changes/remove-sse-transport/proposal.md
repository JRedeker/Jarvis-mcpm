# Change: Remove SSE Transport Support Completely

## Why

SSE (Server-Sent Events) was deprecated in the MCP specification (2025-03-26) and replaced by Streamable HTTP. Our previous change (`switch-from-sse-to-streamable-http`) established HTTP as the default but retained SSE as a fallback. This creates:

1. **Maintenance burden:** Two transport code paths to maintain
2. **Documentation confusion:** Users unsure which transport to use
3. **Technical debt:** SSE fallback logic (e.g., `SSE_PROFILES` in daemon) complicates the codebase
4. **Mixed signals:** Tool descriptions still mention SSE as an option

It's time to complete the migration by removing SSE entirely.

## What Changes

### **BREAKING** Changes
- **Daemon:** Remove `--sse` flag and `SSE_PROFILES` fallback logic from `mcpm-daemon/entrypoint.sh`
- **MCPM CLI:** Remove `sse` from valid transport types in `mcpm new --type` command
- **Jarvis Tools:** Remove `sse` from `create_server` and `edit_server` type descriptions
- **Configuration:** Migrate all remaining `/sse` URLs to `/mcp`

### Documentation Updates
- **README.md:** Update all example configs to use `streamable-http` transport
- **Jarvis/README.md:** Remove SSE endpoint references
- **docs/MCPM-documentation.md:** Remove `--sse` flag documentation, mark as removed
- **docs/SSE-DAEMON-REFACTOR.md:** Archive or remove (historical document)
- **AGENTS.md:** Remove any SSE references

### Configuration File Updates
- **temp_claude.json:** Migrate from SSE to HTTP transport
- **config-templates/:** Ensure no SSE references remain

## Impact

- **Affected specs:** daemon-transport, jarvis-tools, mcpm-cli
- **Affected code:**
  - `mcpm-daemon/entrypoint.sh` (remove SSE_PROFILES)
  - `MCPM/index.js` (remove 'sse' from validTypes)
  - `Jarvis/handlers/server.go` (update tool descriptions)
- **Affected docs:** README.md, Jarvis/README.md, docs/MCPM-documentation.md
- **User impact:** Users with SSE configurations must migrate to HTTP (migration path documented)

## Migration Path for Users

1. Replace `transport: "sse"` with `transport: "streamable-http"` in client configs
2. Replace `/sse` with `/mcp` in endpoint URLs
3. Restart the mcpm-daemon to apply changes
