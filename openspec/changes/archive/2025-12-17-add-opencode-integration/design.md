# Design: OpenCode Client Integration

## Context
OpenCode is an open-source AI coding agent developed by SST that supports MCP servers via two transport types:
1. **Local (stdio)**: Spawns a binary/command and communicates via stdin/stdout
2. **Remote (Streamable HTTP)**: Connects to HTTP endpoints

OpenCode uses a JSON configuration file (`opencode.json` or `opencode.jsonc`) that can be placed:
- Globally: `~/.config/opencode/opencode.json`
- Per-project: `./opencode.json` in project root

This differs from Claude Desktop's simpler `mcpServers` structure.

## Goals
- Enable Jarvis to detect and configure OpenCode as a managed client
- Support both global and per-project OpenCode configurations
- Provide ready-to-use templates for common MCP setups
- Maintain consistency with existing client management patterns (Claude, Cline, etc.)

## Non-Goals
- Modifying OpenCode's source code
- Supporting deprecated SSE transport (OpenCode uses Streamable HTTP for remote)
- Implementing OpenCode-specific features beyond MCP configuration

## Decisions

### 1. Configuration Format Mapping
**Decision**: Map Jarvis profiles to OpenCode's `mcp` config section.

OpenCode format:
```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "server-name": {
      "type": "local" | "remote",
      "command": ["path", "args"],  // for local
      "url": "http://...",          // for remote
      "enabled": true,
      "environment": {},
      "headers": {}
    }
  }
}
```

Jarvis profile mapping:
- `jarvis` (stdio) -> `type: "local"`, `command: ["/path/to/jarvis"]`
- `toolbox`, `memory`, etc. (HTTP) -> `type: "remote"`, `url: "http://localhost:PORT/mcp"`

**Rationale**: Direct mapping preserves OpenCode's native format while enabling Jarvis to manage configurations programmatically.

### 2. Config Path Detection
**Decision**: Use standard OpenCode config locations with environment variable override.

Detection order:
1. `$OPENCODE_CONFIG` (if set)
2. `./opencode.json` (project-local)
3. `~/.config/opencode/opencode.json` (global)

**Rationale**: Matches OpenCode's own resolution order, ensuring Jarvis edits the correct file.

### 3. Client Registry Entry
**Decision**: Add `opencode` to MCPM client registry with detected paths.

```javascript
// MCPM client registry
opencode: {
  name: "OpenCode",
  configPaths: [
    "{env:OPENCODE_CONFIG}",
    "./opencode.json",
    "~/.config/opencode/opencode.json"
  ],
  format: "opencode-mcp"
}
```

**Rationale**: Enables `manage_client` tool to automatically detect and configure OpenCode.

### 4. Template Structure
**Decision**: Create `config-templates/opencode.json` with common profile stack.

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "jarvis": {
      "type": "local",
      "command": ["${JARVIS_PATH}"],
      "enabled": true
    },
    "memory": {
      "type": "remote",
      "url": "http://localhost:6277/mcp",
      "enabled": true
    }
  }
}
```

**Rationale**: Provides users with a working starting point that follows the 3-layer stack philosophy.

## Alternatives Considered

### A. Symlink to Claude Config
**Rejected**: OpenCode has a different schema; symlinking would break parsing.

### B. OpenCode Plugin/Extension
**Rejected**: Out of scope; requires OpenCode ecosystem knowledge and adds maintenance burden.

### C. HTTP-Only Configuration
**Rejected**: Jarvis itself uses stdio transport for direct integration; forcing HTTP would add unnecessary latency and complexity.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| OpenCode config schema changes | Medium | Pin to documented schema version, monitor OpenCode releases |
| Path detection fails on Windows | Low | Linux/Unix focus per project constraints; document limitation |
| User has conflicting global/local configs | Low | Document precedence; prefer project-local edits |

## Migration Plan
1. Add client registry entry (no breaking changes)
2. Create config template
3. Update `manage_client` handler to support OpenCode format
4. Update documentation
5. Test with real OpenCode installation

## Open Questions
- Should we support OpenCode's per-agent MCP configuration? (Currently out of scope)
- Should templates include all profiles or minimal starter set? (Recommend minimal)
