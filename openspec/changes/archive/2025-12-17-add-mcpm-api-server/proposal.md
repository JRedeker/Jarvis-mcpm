# Change: Add MCPM API Server

## Why

Jarvis (Go) currently interacts with MCPM by spawning subprocess calls to the CLI. This approach has several drawbacks:

1. **Subprocess overhead** - Each call spawns a new Node.js process
2. **Text parsing** - CLI output must be parsed as strings, losing type information
3. **Brittle** - If CLI output format changes, Jarvis breaks
4. **No structured errors** - Error handling relies on exit codes and stderr text

A structured HTTP API provides type-safe JSON responses, faster execution (single process), and better testability.

## What Changes

- **MCPM API module:** New `MCPM/api/` directory with Express HTTP server
- **CLI command:** New `mcpm serve` command to start the API server
- **Jarvis HTTP transport:** New `HTTPMcpmRunner` implementation replacing subprocess calls
- **Daemon integration:** API server runs in mcpm-daemon container on port 6275
- **API endpoints:** REST endpoints for all MCPM operations (servers, profiles, clients, system)

## Proposed Solution
Add a thin HTTP API server to MCPM that wraps the existing Python APIs (`GlobalConfigManager`, `ProfileConfigManager`, `RepositoryManager`). Jarvis will call this API instead of the CLI.

### Architecture Change
```
Before:
  Jarvis (Go) → [subprocess] → mcpm CLI → MCPM Python APIs

After:
  Jarvis (Go) → [HTTP/JSON] → MCPM API Server → MCPM Python APIs
```

### Key Benefits
1. **Structured responses** - JSON with typed fields instead of text parsing
2. **Faster** - Single long-running process vs subprocess per call
3. **Type-safe** - Both ends can validate JSON schemas
4. **Testable** - API server can be tested independently
5. **Reusable** - Other tools could use the same API

## Impact

- **Affected specs:** New capability `mcpm-api-server`
- **Affected code:**
  - `MCPM/api/` - New API server module
  - `MCPM/index.js` - Add `serve` command
  - `Jarvis/handlers/http_mcpm.go` - HTTP transport implementation
  - `mcpm-daemon/entrypoint.sh` - Start API server
  - `docker-compose.yml` - Expose port 6275
- **No breaking changes:** CLI remains for human use, HTTP is additive

## Scope

### In Scope
- New HTTP API server module in MCPM (`MCPM/api/`)
- New `mcpm serve` CLI command to start the API server
- Update Jarvis `McpmRunner` interface to support HTTP transport
- New `HTTPMcpmRunner` implementation in Go
- Run API server in mcpm-daemon container
- API endpoints for all operations Jarvis currently uses

### Out of Scope
- Removing the existing CLI (it remains for human use)
- Changing Jarvis's MCP tool interface (tools stay the same)
- Authentication (runs on localhost only)
- WebSocket/streaming (simple request/response is sufficient)

## Risk Assessment
- **Low risk**: Additive change, existing CLI unaffected
- **Rollback**: Can revert Jarvis to CLI calls if issues arise

## Success Criteria
1. All existing Jarvis tests pass with HTTP transport
2. API server starts successfully in mcpm-daemon
3. No subprocess calls to `mcpm` CLI from Jarvis
4. Response times comparable or better than CLI approach
