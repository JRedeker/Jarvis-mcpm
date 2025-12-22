# Design: Enhance Jarvis Operations

## Context

Jarvis provides 9 MCP tools with 31 actions for AI agents to manage infrastructure. However, a gap analysis revealed that agents must still use direct shell commands for:
- Docker container lifecycle (rebuild, stop, start, logs)
- Running tests (Go, Node, integration)
- Enhanced log viewing (stdout, aggregation)
- Selective component builds
- Configuration backup/restore

This design addresses these gaps while maintaining Jarvis's architectural principles:
- Dependency injection for testability
- Interface-based abstractions
- AI-friendly output formatting

## Research Validation

**Date:** December 22, 2025
**Sources:** Context7 (mcp-go, Docker Docs), Docker SDK documentation, Go Effective docs

### Validated Patterns

| Pattern | Status | Evidence |
|:--------|:-------|:---------|
| Action-based routing | **Valid** | mcp-go supports `WithEnum()` for action parameters |
| Error handling | **Valid** | `mcp.NewToolResultError()` for user errors, Go `error` for system errors |
| CLI over SDK | **Recommended** | Docker SDK adds ~50MB deps; CLI via `exec.Command` matches existing patterns |
| Interface extension | **Valid** | Additive changes to `DockerRunner` maintain backward compatibility |

### Key Findings

1. **Docker CLI Approach**: Use `docker compose` CLI commands via `exec.Command` rather than Docker SDK
   - Matches existing `RealDockerRunner` implementation
   - No new dependencies required
   - JSON output available via `--format json` flag

2. **Log Retrieval**: Use `docker compose logs` instead of extending supervisorctl
   - Supervisord only captures stderr by default
   - Container-level logging captures both stdout and stderr
   - Better cross-platform compatibility

3. **Test Execution**: Use `CommandRunner.RunInDir()` for different working directories
   - Go tests in `Jarvis/`
   - Node tests in `MCPM/`
   - Integration tests in project root

## Goals

1. **Eliminate shell fallbacks** - All common operations available via Jarvis tools
2. **Maintain testability** - All new code uses dependency injection
3. **Preserve tool consolidation** - Extend existing tools rather than adding new ones
4. **Keep descriptions concise** - Stay under 150 character limit per tool

## Non-Goals

1. **Git operations** - Intentionally delegated to host agent (Claude, OpenCode)
2. **Real-time streaming** - MCP doesn't support streaming; use polling patterns
3. **Interactive operations** - All operations must be non-interactive
4. **Docker SDK integration** - Use CLI for simplicity and reduced dependencies

## Decisions

### Decision 1: Extend Existing Tools vs Add New Tools

**Choice**: Extend existing tools with new actions

**Rationale**:
- Maintains the 9-tool consolidation from v3.0
- Avoids payload size increases
- Follows established patterns (jarvis_system already has bootstrap, restart, restart_infra)

**Alternatives Considered**:
- Add `jarvis_docker` tool - Rejected: increases tool count, duplicates some jarvis_system functionality
- Add `jarvis_test` tool - Rejected: testing is a project concern, fits `jarvis_project`

### Decision 2: Docker Interface Extension

**Choice**: Extend `DockerRunner` interface with new methods using CLI approach

```go
type DockerRunner interface {
    // Existing methods
    ComposeUp(ctx context.Context, services ...string) error
    ComposeDown(ctx context.Context) error
    ComposeRestart(ctx context.Context, services ...string) error
    ComposePs(ctx context.Context) ([]ContainerStatus, error)
    ExecSupervisorctl(ctx context.Context, action, target string) (string, error)

    // NEW - Phase 1: Docker Operations
    ComposeBuild(ctx context.Context, noCache bool, services ...string) error
    ComposeStop(ctx context.Context, services ...string) error
    ComposeStart(ctx context.Context, services ...string) error
    ComposeLogs(ctx context.Context, service string, lines int) (string, error)
}
```

**Implementation Pattern** (validated against existing code):
```go
func (r *RealDockerRunner) ComposeBuild(ctx context.Context, noCache bool, services ...string) error {
    args := []string{"compose", "build"}
    if noCache {
        args = append(args, "--no-cache")
    }
    args = append(args, services...)
    cmd := exec.CommandContext(ctx, "docker", args...)
    cmd.Dir = r.projectRoot
    return cmd.Run()
}

func (r *RealDockerRunner) ComposeLogs(ctx context.Context, service string, lines int) (string, error) {
    args := []string{"compose", "logs", "--tail", strconv.Itoa(lines)}
    if service != "" {
        args = append(args, service)
    }
    cmd := exec.CommandContext(ctx, "docker", args...)
    cmd.Dir = r.projectRoot
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

**Rationale**:
- Consistent with existing patterns in `handlers.go`
- Enables mocking for tests via `MockDockerClient`
- CLI approach validated as simpler than Docker SDK

### Decision 3: Log Retrieval Strategy (UPDATED)

**Choice**: Use `docker compose logs` for all log types instead of supervisorctl

**Rationale** (from research):
- Supervisord only captures stderr by default
- `docker compose logs` captures both stdout and stderr at container level
- Better formatting and timestamp support
- Cross-platform compatibility

**Implementation Change**:
```go
// OLD approach (supervisorctl - limited to stderr)
output, err := h.Docker.ExecSupervisorctl(ctx, "tail", fmt.Sprintf("mcpm-%s", profile))

// NEW approach (docker compose logs - both streams)
output, err := h.Docker.ComposeLogs(ctx, fmt.Sprintf("mcpm-%s", profile), lines)
```

**Log Type Handling**:
| log_type | Implementation |
|:---------|:---------------|
| `stderr` | Filter container logs (backward compatible default) |
| `stdout` | Filter container logs |
| `all` | Return full container logs |

### Decision 4: Test Runner Implementation

**Choice**: Use `CommandRunner` interface to execute test commands

**Rationale**:
- Tests are just shell commands (`go test`, `npm test`)
- CommandRunner already exists and is mockable
- Captures output for AI-friendly formatting

**Test Suite Mapping** (validated):
| Suite | Command | Working Directory | Timeout |
|:------|:--------|:------------------|:--------|
| `go` | `go test -v ./...` | `Jarvis/` | 5 min |
| `node` | `npm test` | `MCPM/` | 5 min |
| `integration` | `./scripts/tests/run-bats.sh` | Project root | 10 min |
| `all` | Run all above sequentially | - | 20 min |

**Output Formatting**:
```go
func formatTestResult(suite string, output string, err error) string {
    status := "PASS"
    if err != nil {
        status = "FAIL"
    }
    return fmt.Sprintf("## %s Tests: %s\n\n```\n%s\n```", suite, status, output)
}
```

### Decision 5: Log Aggregation Strategy

**Choice**: Server-side aggregation in Jarvis handlers

**Rationale**:
- MCPM API would need new endpoints
- Jarvis can aggregate by calling `ComposeLogs` multiple times
- Reduces API surface area

**Implementation**:
```go
func (h *Handler) DiagnoseLogsAggregated(ctx context.Context, lines int) (string, error) {
    profiles := []string{"essentials", "memory", "dev-core", "data", "research"}
    var result strings.Builder

    for _, profile := range profiles {
        logs, err := h.Docker.ComposeLogs(ctx, fmt.Sprintf("mcpm-%s", profile), lines)
        if err != nil {
            result.WriteString(fmt.Sprintf("## %s: ERROR\n%v\n\n", profile, err))
            continue
        }
        result.WriteString(fmt.Sprintf("## %s\n```\n%s\n```\n\n", profile, logs))
    }
    return result.String(), nil
}
```

### Decision 6: Config Export Format

**Choice**: JSON export matching existing config structure

**Rationale**:
- `servers.json` and `profiles.json` are already JSON
- Easy to restore via import
- Human-readable for debugging

**Export Structure**:
```json
{
  "version": "1.0",
  "exported_at": "2025-12-22T15:00:00Z",
  "warning": "This export may contain sensitive data (API keys). Handle securely.",
  "servers": { /* contents of servers.json */ },
  "profiles": { /* contents of profiles.json */ },
  "clients": { /* contents of clients.json if exists */ }
}
```

## Risks / Trade-offs

### Risk 1: Increased Tool Complexity
- **Risk**: More actions per tool increases cognitive load
- **Mitigation**: Group related actions logically; provide clear documentation
- **Trade-off**: Accept some complexity for operational completeness

### Risk 2: Test Execution Time
- **Risk**: Running full test suite (`all`) could timeout MCP calls
- **Mitigation**: Use generous timeout (20 minutes for all); provide per-suite options
- **Trade-off**: Long-running tests may need to be run separately

### Risk 3: Docker Command Variations
- **Risk**: Docker commands differ across environments (docker vs podman)
- **Mitigation**: Use `docker compose` which handles abstraction
- **Trade-off**: Require Docker Compose; don't support raw Docker

### Risk 4: Log Streaming Performance (NEW)
- **Risk**: Large log files could consume excessive memory
- **Mitigation**: Limit to `--tail N` lines (default 50, max 500)
- **Trade-off**: Cannot retrieve full historical logs in one call

## Migration Plan

1. **Phase 1** (Docker ops): No migration needed - additive changes
2. **Phase 2** (Logging): Backward compatible - new optional parameters
3. **Phase 3** (Testing): No migration needed - new action
4. **Phase 4** (Builds): No migration needed - new action
5. **Phase 5** (Config): No migration needed - new actions

All changes are additive and backward compatible.

## Mock Updates Required

The `MockDockerClient` in `testing/mocks/docker_mock.go` needs these additions:

```go
type MockDockerClient struct {
    // ... existing fields ...

    // NEW fields for Phase 1
    ComposeBuildError  error
    ComposeStopError   error
    ComposeStartError  error
    ComposeLogsOutput  map[string]string
    ComposeLogsError   error
}

func (m *MockDockerClient) ComposeBuild(ctx context.Context, noCache bool, services ...string) error {
    m.recordCall("ComposeBuild", noCache, services)
    return m.ComposeBuildError
}

func (m *MockDockerClient) ComposeStop(ctx context.Context, services ...string) error {
    m.recordCall("ComposeStop", services)
    return m.ComposeStopError
}

func (m *MockDockerClient) ComposeStart(ctx context.Context, services ...string) error {
    m.recordCall("ComposeStart", services)
    return m.ComposeStartError
}

func (m *MockDockerClient) ComposeLogs(ctx context.Context, service string, lines int) (string, error) {
    m.recordCall("ComposeLogs", service, lines)
    if m.ComposeLogsError != nil {
        return "", m.ComposeLogsError
    }
    if output, ok := m.ComposeLogsOutput[service]; ok {
        return output, nil
    }
    return fmt.Sprintf("Logs for %s (last %d lines)", service, lines), nil
}
```

## Open Questions (Resolved)

1. **Q**: Should `docker_logs` support real-time streaming?
   **A**: No - MCP doesn't support streaming. Use `lines` parameter with polling.

2. **Q**: Should test failures return error results or text results?
   **A**: Text results with clear PASS/FAIL indicators. Errors reserved for tool failures.

3. **Q**: Should config export include sensitive env vars (API keys)?
   **A**: Yes, but with a warning in output. User responsible for secure handling.

4. **Q**: Should we use Docker SDK or CLI? (NEW)
   **A**: CLI via `exec.Command`. Validated as simpler, no new deps, matches existing patterns.

5. **Q**: Should logs use supervisorctl or docker compose logs? (NEW)
   **A**: `docker compose logs`. Supervisorctl only captures stderr; container logs capture both streams.
