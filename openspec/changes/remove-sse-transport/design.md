# Technical Design: Remove SSE Transport

## Context

The MCP ecosystem has fully transitioned to Streamable HTTP as the standard transport protocol. SSE support was retained during the initial migration for backward compatibility, particularly for the `morph` profile which had issues with HTTP roots/session handling. This design documents the complete removal of SSE.

## Goals

- **Simplicity:** Single transport protocol reduces code complexity and documentation
- **Spec Compliance:** Full alignment with MCP 2025-03-26 specification
- **Clean Codebase:** Remove dead code paths and conditional transport logic

## Non-Goals

- Maintaining any backward compatibility with SSE
- Supporting mixed SSE/HTTP deployments
- Creating automated migration tooling beyond documentation

## Decisions

### Decision 1: Hard Removal vs Deprecation Warning

**Decision:** Hard removal with no deprecation period.

**Rationale:**
- SSE was already deprecated in previous change
- This is a development tool, not production infrastructure
- Users can manually update configs with documented migration path
- Clean break is simpler than maintaining deprecation warnings

**Alternatives Considered:**
- Emit warning when SSE is used → Adds complexity, delays cleanup
- Auto-migrate on startup → Modifying user configs automatically is risky

### Decision 2: Morph Profile Handling

**Decision:** Force `morph` to use HTTP transport; if issues persist, document workaround.

**Rationale:**
- The `morph-fast-apply` server should now support HTTP
- If it doesn't, the profile can be temporarily disabled
- One server's compatibility shouldn't block platform modernization

**Alternatives Considered:**
- Keep SSE just for morph → Defeats purpose of full removal
- Remove morph from default profiles → Acceptable if HTTP doesn't work

### Decision 3: Error Behavior for Invalid Type

**Decision:** `mcpm new --type sse` SHALL return an error with helpful message.

**Rationale:**
- Silent conversion could confuse users
- Clear error tells users what to do
- Consistent with "fail fast" principle

**Error Message:**
```
Error: Invalid type 'sse'. SSE transport has been removed.
Use 'streamable-http' instead: mcpm new myserver --type streamable-http --url <url>
```

## Architecture Changes

### Daemon Entrypoint (`mcpm-daemon/entrypoint.sh`)

**Before:**
```bash
declare -A SSE_PROFILES=(
    ["morph"]=1
)

if [ -n "${SSE_PROFILES[$profile]}" ]; then
    transport="--sse"
else
    transport="--http"
fi
```

**After:**
```bash
transport="--http"
# SSE support removed - all profiles use HTTP
```

### MCPM CLI (`MCPM/index.js`)

**Before:**
```javascript
const validTypes = ['stdio', 'sse', 'http', 'streamable-http'];
```

**After:**
```javascript
const validTypes = ['stdio', 'http', 'streamable-http'];
// Note: 'http' is alias for 'streamable-http'
```

### Jarvis Tool Descriptions (`Jarvis/handlers/server.go`)

**Before:**
```go
mcp.Description("Transport type: 'stdio', 'sse', or 'streamable-http'")
mcp.Description("URL (for sse/http types)")
```

**After:**
```go
mcp.Description("Transport type: 'stdio' or 'streamable-http'")
mcp.Description("URL (for streamable-http type)")
```

## Files to Modify

| File | Change Type | Description |
|------|-------------|-------------|
| `mcpm-daemon/entrypoint.sh` | Code | Remove SSE_PROFILES, use --http only |
| `MCPM/index.js` | Code | Remove 'sse' from validTypes, add error |
| `Jarvis/handlers/server.go` | Code | Update tool descriptions |
| `temp_claude.json` | Config | Migrate /sse to /mcp |
| `README.md` | Doc | Update example configs |
| `Jarvis/README.md` | Doc | Remove SSE references |
| `docs/MCPM-documentation.md` | Doc | Remove --sse flag docs |
| `docs/SSE-DAEMON-REFACTOR.md` | Doc | Archive (move to archive/) |

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Morph profile breaks with HTTP | Medium | Test before release; document workaround |
| Users with custom SSE configs | Low | Clear migration docs; error messages guide users |
| External scripts using `--sse` | Low | Scripts should be updated; error message helps |

## Migration Plan

1. **Phase 1:** Update daemon entrypoint (removes SSE fallback)
2. **Phase 2:** Update MCPM CLI (removes SSE type)
3. **Phase 3:** Update Jarvis tools (removes SSE from descriptions)
4. **Phase 4:** Update configs and documentation
5. **Phase 5:** Archive historical SSE documentation

## Open Questions

None - this is a straightforward removal based on the completed transition to HTTP.
