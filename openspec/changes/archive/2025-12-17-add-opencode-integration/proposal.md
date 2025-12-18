# Change: Add OpenCode Client Integration

## Why
OpenCode is a modern open-source AI coding agent (CLI/TUI/IDE) that supports the Model Context Protocol (MCP). To expand Jarvis's reach beyond Claude Desktop and VS Code, we need to integrate OpenCode as a supported client. This enables users of OpenCode to leverage Jarvis's infrastructure management tools and the composable profile system.

## What Changes
- Add OpenCode as a recognized client in the MCPM client registry
- Create configuration templates for OpenCode's `opencode.json` format
- Document OpenCode-specific MCP server configuration (local stdio + remote HTTP transports)
- Update Jarvis `manage_client` tool to support OpenCode config path detection and editing
- Provide example configurations for common profile stacks (jarvis, memory, project profiles)

## Impact
- **Affected specs**: New capability `opencode-client` (client integration)
- **Affected code**:
  - `Jarvis/handlers/handlers.go` - client detection logic
  - `MCPM/index.js` - client registry
  - `config-templates/` - new OpenCode template
  - `docs/` - documentation updates
- **No breaking changes**: Additive feature only
