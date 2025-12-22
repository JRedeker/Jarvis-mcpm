# Jarvis Tools Specification

## Research Validation Summary

**Validated:** December 22, 2025
**Sources:** Context7 (mcp-go v0.43.2, Docker Docs), Go Effective docs

| Decision | Recommendation | Rationale |
|:---------|:---------------|:----------|
| Docker approach | CLI via `exec.Command` | Matches existing patterns, no new deps |
| Log retrieval | `docker compose logs` | Captures stdout+stderr (supervisorctl only stderr) |
| Interface changes | Extend `DockerRunner` | Additive, backward compatible |
| Error handling | `mcp.NewToolResultError()` | Follows mcp-go patterns |

---

## Interface Requirements

### DockerRunner Interface Extension

The `DockerRunner` interface SHALL be extended with the following methods:

```go
type DockerRunner interface {
    // Existing methods (unchanged)
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

#### Scenario: Interface backward compatibility
- **WHEN** existing code calls `DockerRunner` methods
- **THEN** behavior SHALL be identical to previous version
- **AND** no breaking changes SHALL occur

---

## ADDED Requirements

### Requirement: Docker Container Lifecycle Management

The `jarvis_system` tool SHALL provide actions to manage Docker container lifecycle without requiring direct shell access.

**Implementation:** Uses `docker compose` CLI commands via `exec.CommandContext()`.

#### Scenario: Rebuild containers after code changes
- **WHEN** agent calls `jarvis_system(action="rebuild")`
- **THEN** the handler SHALL execute `docker compose build --no-cache`
- **AND** then execute `docker compose up -d`
- **AND** the output SHALL indicate success or failure with details

#### Scenario: Rebuild specific service
- **WHEN** agent calls `jarvis_system(action="rebuild", service="mcpm-daemon")`
- **THEN** only the specified service SHALL be rebuilt
- **AND** the command SHALL be `docker compose build --no-cache mcpm-daemon`
- **AND** other containers SHALL remain unchanged

#### Scenario: Stop all containers
- **WHEN** agent calls `jarvis_system(action="stop")`
- **THEN** the handler SHALL execute `docker compose stop`
- **AND** the output SHALL list stopped services

#### Scenario: Stop specific service
- **WHEN** agent calls `jarvis_system(action="stop", service="mcp-qdrant")`
- **THEN** the handler SHALL execute `docker compose stop mcp-qdrant`
- **AND** only the specified service SHALL be stopped

#### Scenario: Start containers
- **WHEN** agent calls `jarvis_system(action="start")`
- **THEN** the handler SHALL execute `docker compose start`
- **AND** the output SHALL indicate service status

#### Scenario: Start specific service
- **WHEN** agent calls `jarvis_system(action="start", service="mcpm-daemon")`
- **THEN** the handler SHALL execute `docker compose start mcpm-daemon`
- **AND** only the specified service SHALL be started

#### Scenario: View container logs
- **WHEN** agent calls `jarvis_system(action="docker_logs", service="mcpm-daemon", lines=100)`
- **THEN** the handler SHALL execute `docker compose logs --tail 100 mcpm-daemon`
- **AND** the output SHALL be formatted for readability
- **AND** lines parameter SHALL default to 50 if not provided
- **AND** lines SHALL be capped at 500 maximum

#### Scenario: View container status
- **WHEN** agent calls `jarvis_system(action="docker_status")`
- **THEN** the handler SHALL execute `docker compose ps --format json`
- **AND** the JSON output SHALL be parsed into structured format
- **AND** the output SHALL include container name, status, health, and ports

#### Scenario: Docker not available
- **WHEN** agent calls any docker action and Docker is not running
- **THEN** the handler SHALL return `mcp.NewToolResultError()` with clear message
- **AND** the message SHALL suggest checking Docker installation

---

### Requirement: Enhanced Log Retrieval

The `jarvis_diagnose` tool SHALL provide enhanced log retrieval using `docker compose logs` instead of supervisorctl.

**Implementation Change (from research):** Use `docker compose logs` for all log retrieval as it captures both stdout and stderr, whereas supervisorctl only captures stderr.

#### Scenario: Retrieve stderr logs (default - backward compatible)
- **WHEN** agent calls `jarvis_diagnose(action="logs", profile="research")`
- **THEN** the handler SHALL execute `docker compose logs --tail 50 mcpm-research`
- **AND** stderr content SHALL be returned (default behavior preserved)

#### Scenario: Retrieve stdout logs
- **WHEN** agent calls `jarvis_diagnose(action="logs", profile="research", log_type="stdout")`
- **THEN** the handler SHALL execute `docker compose logs --tail 50 mcpm-research`
- **AND** stdout content SHALL be filtered and returned

#### Scenario: Retrieve all logs
- **WHEN** agent calls `jarvis_diagnose(action="logs", profile="research", log_type="all")`
- **THEN** the handler SHALL execute `docker compose logs --tail 50 mcpm-research`
- **AND** both stdout and stderr SHALL be returned
- **AND** logs SHALL be clearly labeled by type

#### Scenario: Aggregate logs across profiles
- **WHEN** agent calls `jarvis_diagnose(action="logs", aggregate=true)`
- **THEN** logs SHALL be retrieved for profiles: essentials, memory, dev-core, data, research
- **AND** each profile's logs SHALL be prefixed with profile name in markdown headers
- **AND** profiles that fail to retrieve logs SHALL show error message

#### Scenario: Aggregate with log type filter
- **WHEN** agent calls `jarvis_diagnose(action="logs", aggregate=true, log_type="stderr")`
- **THEN** stderr logs from all profiles SHALL be returned

#### Scenario: Custom line count
- **WHEN** agent calls `jarvis_diagnose(action="logs", profile="research", lines=200)`
- **THEN** the handler SHALL retrieve last 200 lines
- **AND** lines SHALL be capped at 500 maximum

---

### Requirement: Test Suite Execution

The `jarvis_project` tool SHALL provide a `test` action to run project test suites.

**Implementation:** Uses `CommandRunner.RunInDir()` with appropriate working directories.

#### Scenario: Run Go tests
- **WHEN** agent calls `jarvis_project(action="test", suite="go")`
- **THEN** the handler SHALL execute `go test -v ./...` in `Jarvis/` directory
- **AND** timeout SHALL be 5 minutes
- **AND** the output SHALL include PASS/FAIL summary

#### Scenario: Run Go tests with verbose output
- **WHEN** agent calls `jarvis_project(action="test", suite="go", verbose=true)`
- **THEN** the handler SHALL execute `go test -v ./...`
- **AND** full verbose test output SHALL be included

#### Scenario: Run Node tests
- **WHEN** agent calls `jarvis_project(action="test", suite="node")`
- **THEN** the handler SHALL execute `npm test` in `MCPM/` directory
- **AND** timeout SHALL be 5 minutes
- **AND** the output SHALL include test results

#### Scenario: Run integration tests
- **WHEN** agent calls `jarvis_project(action="test", suite="integration")`
- **THEN** the handler SHALL execute `./scripts/tests/run-bats.sh` in project root
- **AND** timeout SHALL be 10 minutes
- **AND** the output SHALL include test results

#### Scenario: Run all tests
- **WHEN** agent calls `jarvis_project(action="test", suite="all")`
- **THEN** Go, Node, and integration tests SHALL be executed sequentially
- **AND** total timeout SHALL be 20 minutes
- **AND** the output SHALL include combined results with per-suite summary
- **AND** format SHALL be markdown with suite headers

#### Scenario: Test suite not found
- **WHEN** agent calls `jarvis_project(action="test", suite="invalid")`
- **THEN** `mcp.NewToolResultError()` SHALL be returned
- **AND** error message SHALL list valid options: go, node, integration, all

#### Scenario: Test failure handling
- **WHEN** tests fail (non-zero exit code)
- **THEN** the handler SHALL return text result (NOT error result)
- **AND** output SHALL clearly indicate FAIL status
- **AND** test output SHALL be included for debugging

---

### Requirement: Selective Component Builds

The `jarvis_system` tool SHALL provide a `build` action to build specific components.

#### Scenario: Build Jarvis binary
- **WHEN** agent calls `jarvis_system(action="build", component="jarvis")`
- **THEN** the handler SHALL execute `go build -o jarvis .` in `Jarvis/` directory
- **AND** the output SHALL indicate success or compilation errors

#### Scenario: Build MCPM dependencies
- **WHEN** agent calls `jarvis_system(action="build", component="mcpm")`
- **THEN** the handler SHALL execute `npm install` in `MCPM/` directory
- **AND** the output SHALL indicate installed packages count

#### Scenario: Build daemon container
- **WHEN** agent calls `jarvis_system(action="build", component="daemon")`
- **THEN** the handler SHALL execute `docker compose build mcpm-daemon`
- **AND** the output SHALL indicate build status

#### Scenario: Build all components
- **WHEN** agent calls `jarvis_system(action="build", component="all")`
- **THEN** Jarvis, MCPM, and daemon SHALL be built sequentially
- **AND** the output SHALL include status for each component
- **AND** format SHALL be markdown with component headers

#### Scenario: Invalid component
- **WHEN** agent calls `jarvis_system(action="build", component="invalid")`
- **THEN** `mcp.NewToolResultError()` SHALL be returned
- **AND** error message SHALL list valid options: jarvis, mcpm, daemon, all

---

### Requirement: Configuration Export and Import

The `jarvis_config` tool SHALL provide `export` and `import` actions for configuration backup and restore.

#### Scenario: Export configuration
- **WHEN** agent calls `jarvis_config(action="export")`
- **THEN** a JSON blob SHALL be returned containing:
  - `servers`: contents of `~/.config/mcpm/servers.json`
  - `profiles`: contents of `~/.config/mcpm/profiles.json`
  - `clients`: contents of client configs if they exist
- **AND** the export SHALL include `version` and `exported_at` metadata
- **AND** a warning about sensitive data (API keys) SHALL be included

#### Scenario: Import configuration from data
- **WHEN** agent calls `jarvis_config(action="import", data="<json>")`
- **THEN** the JSON SHALL be parsed and validated
- **AND** configuration files SHALL be updated atomically
- **AND** the output SHALL indicate which files were updated

#### Scenario: Import invalid JSON
- **WHEN** agent calls `jarvis_config(action="import", data="invalid")`
- **THEN** `mcp.NewToolResultError()` SHALL be returned
- **AND** error message SHALL include JSON parsing error details

#### Scenario: Import partial configuration
- **WHEN** agent calls `jarvis_config(action="import", data="<json with only servers>")`
- **THEN** only the provided configuration sections SHALL be updated
- **AND** missing sections SHALL be skipped with a warning in output

#### Scenario: Import with backup
- **WHEN** agent calls `jarvis_config(action="import", data="<json>")`
- **THEN** existing configs SHALL NOT be backed up (rely on git for history)
- **AND** import SHALL overwrite existing values

---

## MODIFIED Requirements

### Requirement: System Tool Actions

The `jarvis_system` tool action enum SHALL be updated:

```go
mcp.WithEnum("bootstrap", "restart", "restart_infra", "rebuild", "stop", "start", "docker_logs", "docker_status", "build")
```

#### Scenario: List valid actions on error
- **WHEN** agent calls `jarvis_system` with invalid action
- **THEN** error message SHALL list all valid actions
- **AND** format SHALL be: `Invalid action 'X'. Valid actions: bootstrap, restart, restart_infra, rebuild, stop, start, docker_logs, docker_status, build`

---

### Requirement: Diagnose Tool Parameters

The `jarvis_diagnose` tool `logs` action SHALL support new parameters:

```go
mcp.WithString("log_type",
    mcp.Description("Log type: stderr (default), stdout, all"),
    mcp.Enum("stderr", "stdout", "all"),
),
mcp.WithBoolean("aggregate",
    mcp.Description("Aggregate logs from all profiles"),
),
```

#### Scenario: Default behavior preserved
- **WHEN** agent calls `jarvis_diagnose(action="logs", profile="x")` without new parameters
- **THEN** behavior SHALL be identical to previous version
- **AND** `log_type` SHALL default to `stderr`
- **AND** `aggregate` SHALL default to `false`

---

### Requirement: Project Tool Actions

The `jarvis_project` tool action enum SHALL be updated:

```go
mcp.WithEnum("analyze", "diff", "devops", "test")
```

#### Scenario: List valid actions on error
- **WHEN** agent calls `jarvis_project` with invalid action
- **THEN** error message SHALL list: analyze, diff, devops, test

---

### Requirement: Config Tool Actions

The `jarvis_config` tool action enum SHALL be updated:

```go
mcp.WithEnum("get", "set", "list", "migrate", "export", "import")
```

#### Scenario: List valid actions on error
- **WHEN** agent calls `jarvis_config` with invalid action
- **THEN** error message SHALL list: get, set, list, migrate, export, import

---

## Error Handling Patterns

Based on mcp-go best practices (from research):

### User Errors (validation, invalid input)
```go
return mcp.NewToolResultError("Invalid service name: must not contain spaces"), nil
```

### System Errors (Docker down, file not found)
```go
return nil, fmt.Errorf("docker compose failed: %w", err)
```

### Partial Failures (some operations succeed, some fail)
```go
return &mcp.CallToolResult{
    Content: []mcp.Content{
        mcp.NewTextContent(fmt.Sprintf("Partial success:\n%s\n\nFailures:\n%s", successes, failures)),
    },
}, nil
```
