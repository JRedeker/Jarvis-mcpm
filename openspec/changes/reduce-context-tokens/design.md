# Design: Reduce Context Tokens for Jarvis MCP

## Context

Jarvis is an MCP server that exposes tools for managing other MCP servers, profiles, and infrastructure. When AI clients connect, they receive metadata about all available tools. This metadata consumes context tokens that could otherwise be used for reasoning.

### Stakeholders
- **AI Clients**: Claude, OpenCode, Cline - need efficient tool discovery
- **Developers**: Use Jarvis for infrastructure automation
- **Jarvis Maintainers**: Need to balance usability vs. efficiency

### Constraints
- MCP protocol requires tool metadata on connection
- SDK (mcp-go v0.43.2) always includes annotation fields
- Must maintain backward compatibility with MCPM CLI patterns
- Cannot change how AI clients consume MCP responses

## Goals / Non-Goals

### Goals
1. Reduce tool listing payload by >60%
2. Maintain all existing functionality
3. Improve AI comprehension through clearer tool organization
4. Establish naming conventions for multi-server clarity

### Non-Goals
- Modifying the MCP SDK (future work)
- Changing MCPM CLI interface
- Adding new functionality
- Per-session tool filtering (possible future enhancement)

## Research Findings

### Current State Analysis

```
Tool Categories (24 total):
├── System Management (1): check_status
├── Server Management (8): list_servers, server_info, install_server,
│                          uninstall_server, search_servers, edit_server,
│                          create_server, usage_stats
├── Profile Management (3): manage_profile, suggest_profile, restart_profiles
├── Client Management (1): manage_client
├── Configuration (2): manage_config, migrate_config
├── Project Analysis (3): analyze_project, fetch_diff_context, apply_devops_stack
├── Infrastructure (3): bootstrap_system, restart_service, restart_infrastructure
└── Server Sharing (3): share_server, stop_sharing_server, list_shared_servers
```

### Token Analysis

| Component | Bytes | Est. Tokens |
|-----------|-------|-------------|
| Tool names | 360 | 90 |
| Descriptions | 2,937 | 734 |
| Input schemas | 4,528 | 1,132 |
| Annotations | 3,240 | 810 |
| JSON structure | ~1,000 | 250 |
| **Total** | **~11,065** | **~2,750** |

### SDK Behavior (mcp-go v0.43.2)

The `NewTool()` function sets default annotations:
```go
Annotations: ToolAnnotation{
    Title:           "",
    ReadOnlyHint:    ToBoolPtr(false),
    DestructiveHint: ToBoolPtr(true),
    IdempotentHint:  ToBoolPtr(false),
    OpenWorldHint:   ToBoolPtr(true),
}
```

These are ALWAYS serialized via custom `MarshalJSON`:
```go
m["annotations"] = t.Annotations  // No conditional check
```

**Limitation**: Cannot omit annotations without forking the SDK.

### Consolidation Patterns in Other Systems

1. **GitHub CLI**: Uses `gh <noun> <verb>` pattern
2. **kubectl**: Uses `kubectl <verb> <resource>` pattern
3. **Docker CLI**: Uses `docker <command> [OPTIONS]` pattern

Common pattern: **Noun-based tools with action parameters**

## Decisions

### Decision 1: Action-Based Consolidation

**Choice**: Consolidate 24 tools into 8 using action parameters

**Rationale**:
- Reduces tool count by 67%
- Follows established CLI patterns (gh, kubectl)
- Maintains logical grouping
- Allows future action additions without new tools

**Alternatives Considered**:
1. **Keep all tools, shorten descriptions only**: Only ~20% reduction
2. **Use tool prefixes without consolidation**: No structural improvement
3. **Dynamic tool filtering per session**: Complex, SDK support limited

### Decision 2: Naming Convention

**Choice**: Use `jarvis_<noun>` pattern

**Rationale**:
- Clear namespace in multi-server environments
- AI can identify tool source easily
- Consistent with other MCP servers (e.g., `memory_`, `context7_`)

**Examples**:
- `jarvis_server` not `server_manager`
- `jarvis_profile` not `profile_management`

### Decision 3: Action Parameter Design

**Choice**: Required `action` parameter with proper enum constraint

Per mcp-go SDK documentation, use `mcp.Enum()` for type-safe action validation:

```go
mcp.WithString("action",
    mcp.Required(),
    mcp.Enum("list", "info", "install", "uninstall", "search", "edit", "create", "usage"),
    mcp.Description("Operation to perform"),
)
```

**Rationale**:
- SDK-native enum validation (documented pattern in mcp-go)
- Schema includes allowed values in `enum` field
- AI clients can validate before calling
- Consistent with MCP specification for inputSchema

**SDK Reference**: The `mcp.Enum()` function is the recommended pattern per mcp-go documentation for constrained string values

### Decision 4: Description Strategy

**Choice**: Concise, action-focused descriptions

**Before** (122 chars avg):
```
"Installs an MCP server from the registry with automatic dependency
resolution. Validates the server exists before installing and suggests
alternatives for typos."
```

**After** (~50 chars):
```
"Manage MCP servers: install, configure, search."
```

**Rationale**:
- AI doesn't need verbose explanations
- Action list provides sufficient context
- Error messages provide detailed guidance when needed

## Consolidated Tool Architecture

### Tool Mapping

```
jarvis_check_status
├── action: (none - single purpose)
└── Maps to: check_status handler

jarvis_server
├── action: list → list_servers
├── action: info → server_info
├── action: install → install_server
├── action: uninstall → uninstall_server
├── action: search → search_servers
├── action: edit → edit_server
├── action: create → create_server
└── action: usage → usage_stats

jarvis_profile
├── action: list → manage_profile(action=ls)
├── action: create → manage_profile(action=create)
├── action: edit → manage_profile(action=edit)
├── action: delete → manage_profile(action=delete)
├── action: suggest → suggest_profile
└── action: restart → restart_profiles

jarvis_client
├── action: list → manage_client(action=ls)
├── action: edit → manage_client(action=edit)
├── action: import → manage_client(action=import)
└── action: config → manage_client(action=config)

jarvis_config
├── action: get → manage_config(action=get)
├── action: set → manage_config(action=set)
├── action: list → manage_config(action=ls)
└── action: migrate → migrate_config

jarvis_project
├── action: analyze → analyze_project
├── action: diff → fetch_diff_context
└── action: devops → apply_devops_stack

jarvis_system
├── action: bootstrap → bootstrap_system
├── action: restart → restart_service
└── action: restart_infra → restart_infrastructure

jarvis_share
├── action: start → share_server
├── action: stop → stop_sharing_server
└── action: list → list_shared_servers
```

### Handler Architecture

Maintain internal handler separation for testability. Use SDK's type-safe parameter extraction:

```go
// Public consolidated tool handler
func (h *Handler) Server(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Use SDK's RequireString for required parameters (per mcp-go docs)
    action, err := req.RequireString("action")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    // Go switch pattern per Effective Go
    switch action {
    case "list":
        return h.listServers(ctx, req)
    case "info":
        return h.serverInfo(ctx, req)
    case "install":
        return h.installServer(ctx, req)
    case "uninstall":
        return h.uninstallServer(ctx, req)
    case "search":
        return h.searchServers(ctx, req)
    case "edit":
        return h.editServer(ctx, req)
    case "create":
        return h.createServer(ctx, req)
    case "usage":
        return h.usageStats(ctx, req)
    default:
        return mcp.NewToolResultError("Invalid action. Valid: list|info|install|uninstall|search|edit|create|usage"), nil
    }
}

// Private handlers remain unchanged for testing
func (h *Handler) listServers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Existing implementation
}
```

**SDK Pattern Notes**:
- `req.RequireString()` returns error if parameter missing (mcp-go documented pattern)
- `req.GetString(key, default)` for optional parameters
- Switch statement follows Go's comma-separated case pattern when needed

## Projected Results

### Payload Comparison

**Current (24 tools)**:
```json
{
  "tools": [
    {"name": "check_status", "description": "Comprehensive...", "inputSchema": {...}, "annotations": {...}},
    {"name": "list_servers", "description": "Shows all...", "inputSchema": {...}, "annotations": {...}},
    // ... 22 more tools
  ]
}
// ~11,065 bytes
```

**After (8 tools)**:
```json
{
  "tools": [
    {"name": "jarvis_check_status", "description": "System health check.", "inputSchema": {...}, "annotations": {...}},
    {"name": "jarvis_server", "description": "Manage MCP servers.", "inputSchema": {"action": {...}, "name": {...}}, "annotations": {...}},
    // ... 6 more tools
  ]
}
// ~3,000 bytes (estimated)
```

### Token Savings

| Metric | Before | After | Reduction |
|--------|--------|-------|-----------|
| Tools | 24 | 8 | 67% |
| Payload | 11 KB | 3 KB | 73% |
| Tokens | ~2,750 | ~750 | 73% |

## Risks / Trade-offs

### Risk 1: Breaking Change
- **Impact**: Existing AI workflows break
- **Mitigation**: Clear migration docs, helpful error messages
- **Severity**: High but one-time

### Risk 2: Action Discoverability
- **Impact**: AI may not know valid actions
- **Mitigation**: Actions listed in description, errors suggest valid actions
- **Severity**: Medium

### Risk 3: Handler Complexity
- **Impact**: Consolidated handlers have more branching
- **Mitigation**: Keep internal handlers separate, test individually
- **Severity**: Low

## Migration Plan

### Phase 1: Implement Consolidated Tools
1. Create new `jarvis_*` tool definitions
2. Implement action routing in handlers
3. Update tests for new structure

### Phase 2: Documentation Update
1. Update AGENTS.md with new tool names
2. Add migration guide for breaking changes
3. Update examples in docs

### Phase 3: Deprecation (Optional Future)
1. Keep old tool names as aliases (if needed)
2. Log deprecation warnings
3. Remove in next major version

## Open Questions

1. **Alias Support**: Should we keep old tool names as deprecated aliases?
   - Recommendation: No, clean break is simpler

2. **Per-Session Filtering**: Should different clients see different tools?
   - Recommendation: Future enhancement, not in this change

3. **SDK Enhancement**: Should we contribute omitempty support upstream?
   - Recommendation: Yes, as separate effort after this change
