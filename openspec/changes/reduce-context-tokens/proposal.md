# Change: Reduce Context Tokens for Jarvis MCP Connections

## Why

When an AI client (Claude, OpenCode, etc.) connects to Jarvis via MCP, it receives a full tool listing that consumes significant context tokens. Currently, Jarvis exposes **24 individual tools** that generate approximately **11 KB of JSON metadata** per connection. This metadata includes:

- Tool names and descriptions (~3 KB)
- Input schemas with parameter definitions (~4 KB)
- Default annotations that are always included (~3 KB)
- Redundant/verbose JSON structures

For AI agents with limited context windows (especially in multi-tool scenarios where Jarvis is one of many MCP servers), this overhead directly reduces the tokens available for actual reasoning and conversation.

### Quantified Impact

| Metric | Current | Problem |
|--------|---------|---------|
| Tools | 24 | High cognitive load for AI to understand |
| JSON Payload | 11,065 bytes | ~2,750 tokens consumed just for tool discovery |
| Description Text | 2,937 bytes | Verbose, repetitive phrasing |
| Annotations Overhead | ~3,240 bytes | Always sent even when using defaults |

## What Changes

### 1. Tool Consolidation (Primary)
Reduce from 24 tools to **8 consolidated tools** using action-based patterns:

| Current Tools | Consolidated Tool |
|---------------|-------------------|
| `check_status` | `jarvis_check_status` (unchanged) |
| `list_servers`, `server_info`, `install_server`, `uninstall_server`, `search_servers`, `edit_server`, `create_server` | `jarvis_server` (with `action` param) |
| `manage_profile`, `suggest_profile`, `restart_profiles` | `jarvis_profile` (with `action` param) |
| `manage_client` | `jarvis_client` (with `action` param) |
| `manage_config`, `migrate_config` | `jarvis_config` (with `action` param) |
| `analyze_project`, `fetch_diff_context`, `apply_devops_stack` | `jarvis_project` (with `action` param) |
| `bootstrap_system`, `restart_service`, `restart_infrastructure` | `jarvis_system` (with `action` param) |
| `share_server`, `stop_sharing_server`, `list_shared_servers` | `jarvis_share` (with `action` param) |
| `usage_stats` | Merged into `jarvis_server action=usage` |

### 2. Description Optimization (Secondary)
- Shorten descriptions to essential information only
- Remove redundant phrases like "Use this to..." and "Returns..."
- Average description reduced from 122 chars to ~60 chars

### 3. Annotation Optimization (SDK Limitation)
- Mark read-only tools appropriately (`readOnlyHint: true`)
- The SDK always sends annotations; we'll use meaningful values
- Future: Consider SDK PR to support omitempty for annotations

### 4. Tool Naming Convention
- Add `jarvis_` prefix to all tools for clear namespace identification
- Helps AI agents understand tool provenance in multi-server setups

## Impact

### Affected Code
- `Jarvis/handlers/server.go` - Tool definitions
- `Jarvis/handlers/handlers.go` - Handler implementations
- `Jarvis/main.go` - Server configuration
- `AGENTS.md` - Tool documentation for AI assistants

### Breaking Changes
- **BREAKING**: All tool names change (24 tools become 8)
- **BREAKING**: Parameters restructured to use action-based pattern
- AI clients using old tool names will need to adapt

### Expected Results

| Metric | Current | After | Improvement |
|--------|---------|-------|-------------|
| Tools | 24 | 8 | 67% reduction |
| JSON Payload | ~11 KB | ~3 KB | 73% reduction |
| Est. Tokens | ~2,750 | ~750 | 73% reduction |

### Risks
- **Migration**: Existing AI workflows using specific tool names will break
- **Discoverability**: Action-based tools require AI to understand parameters
- **Complexity**: Handlers become more complex with action routing

### Mitigations
- Clear documentation in AGENTS.md
- Helpful error messages when invalid actions used
- Maintain internal handler logic separation for testability
