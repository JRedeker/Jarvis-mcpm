# Change: Add Granular Client Editing for AI Agents

## Why

AI agents currently lack fine-grained control over MCP client configurations. The existing `jarvis_client` tool can only add/remove entire profiles, but cannot:

1. **Enable/disable individual servers** without removing them
2. **Set environment variables** on servers (e.g., API keys)
3. **Modify server URLs** for custom endpoints
4. **Add custom headers** for authenticated remotes

This forces agents to either:
- Manually edit JSON files (error-prone, loses context)
- Remove and re-add servers to change settings (destructive)
- Ask users to manually configure sensitive settings

## What Changes

Extend `jarvis_client` with new parameters for granular server management:

- `set_enabled`: Toggle server enabled state (`"server_name=true"` or `"server_name=false"`)
- `set_env`: Set environment variables (`"server_name:API_KEY=value"`)
- `set_url`: Modify server URL (`"server_name=http://..."`)
- `set_header`: Add/modify headers (`"server_name:Authorization=Bearer xxx"`)
- `remove_env`: Remove environment variable (`"server_name:VAR_NAME"`)
- `remove_header`: Remove header (`"server_name:Header-Name"`)

## Impact

- Affected specs: `jarvis-tools` (Client Management Tool requirement)
- Affected code: `Jarvis/handlers/handlers.go`, `Jarvis/handlers/opencode.go`
- **NOT breaking**: All existing parameters remain, new parameters are additive
- **Security consideration**: Env vars may contain secrets - logging must redact values
