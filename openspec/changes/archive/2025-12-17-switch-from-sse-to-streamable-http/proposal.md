# Change: Switch from SSE to Streamable HTTP

## Why

The MCP specification (2025-03-26) deprecated SSE in favor of Streamable HTTP. Our current `mcpm-daemon` architecture uses SSE endpoints (`/sse`) which will become incompatible with future clients. Migrating now ensures:

1. **Spec compliance:** Align with the official MCP transport standard
2. **Future compatibility:** Clients are dropping SSE support
3. **Better connection handling:** Streamable HTTP has improved lifecycle management

## What Changes

- **mcpm-daemon:** Replace `/sse` endpoints with `/mcp` using Streamable HTTP transport
- **Jarvis tools:** Update `create_server` to accept `streamable-http` type and default to it for URL-based servers
- **Jarvis health checks:** Update `check_status` to verify HTTP endpoints instead of SSE
- **Configuration migration:** Create `update_configs.py` script to migrate existing SSE configs automatically
- **Templates:** Update `config-templates/` to use `streamable-http` transport by default

## Impact

- **Affected specs:** daemon-transport, jarvis-tools, config-migration
- **Affected code:**
  - `mcpm-daemon/entrypoint.sh` (transport flags)
  - `Jarvis/handlers/server.go` (type validation)
  - `config-templates/*.json` (transport settings)
- **User impact:** Existing SSE configurations MUST be migrated (automated script provided)

## Migration Path

1. Run `update_configs.py` to migrate existing configs
2. Restart `mcpm-daemon` to apply new transport
3. Verify endpoints respond at `/mcp` instead of `/sse`

## Non-Goals

- Changing underlying `mcpm` logic unrelated to transport
- Refactoring Jarvis beyond transport handling

## Risks

- **Client compatibility:** Older clients may not support Streamable HTTP immediately
- **Downtime:** Migration requires daemon restart and config updates
