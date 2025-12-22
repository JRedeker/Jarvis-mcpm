# Change: Enhance Jarvis Operations to Eliminate Shell Command Gaps

**Status:** APPROVED (Research Validated)
**Validated:** December 22, 2025
**Estimated Effort:** 16-20 hours

## Why

AI agents currently must fall back to direct shell commands (`docker exec`, `docker compose`, manual test runs) for critical operations that Jarvis tools don't cover. This breaks the "use Jarvis tools, not shell" mandate and reduces agent autonomy. The gaps include:

1. **Docker operations** - Agents cannot rebuild containers, view logs, or recover from crashes
2. **Testing** - No way to run Go/Node/integration tests to verify changes
3. **Enhanced logging** - Only stderr logs available via supervisorctl, no stdout or aggregated views
4. **Selective builds** - Bootstrap is all-or-nothing; cannot rebuild individual components
5. **Config backup** - No export/import for disaster recovery

## Research Validation

| Decision | Recommendation | Status |
|:---------|:---------------|:-------|
| Docker approach | CLI via `exec.Command` | Validated |
| Log retrieval | `docker compose logs` (not supervisorctl) | Changed |
| Interface extension | Extend `DockerRunner` interface | Validated |
| Error handling | `mcp.NewToolResultError()` for user errors | Validated |

**Key Finding:** Changed Phase 2 to use `docker compose logs` instead of supervisorctl because supervisorctl only captures stderr. Container-level logging captures both stdout and stderr.

## What Changes

### Phase 1: Docker Operations (HIGH PRIORITY)
- Add `rebuild`, `stop`, `start`, `docker_logs`, `docker_status` actions to `jarvis_system`
- Extend `DockerRunner` interface with `ComposeBuild`, `ComposeStop`, `ComposeStart`, `ComposeLogs`
- Enable agents to manage container lifecycle without shell access

### Phase 2: Enhanced Logging (MEDIUM PRIORITY)
- Extend `jarvis_diagnose logs` with `log_type` parameter (stderr|stdout|all)
- Add `aggregate` parameter to combine logs across profiles
- **Changed:** Use `docker compose logs` instead of supervisorctl for log retrieval

### Phase 3: Test Runner (MEDIUM PRIORITY)
- Add `test` action to `jarvis_project`
- Support Go, Node, and integration test suites with configurable timeouts

### Phase 4: Selective Builds (MEDIUM PRIORITY)
- Add `build` action to `jarvis_system`
- Target specific components (jarvis|mcpm|daemon|all)

### Phase 5: Config Backup (LOW PRIORITY)
- Add `export` and `import` actions to `jarvis_config`
- Enable full configuration backup and restore with versioning

## Impact

- **Affected specs**: jarvis-tools (new capability)
- **Affected code**:
  - `Jarvis/handlers/handlers.go` - New handler methods, DockerRunner interface extension
  - `Jarvis/handlers/consolidated.go` - Action routing for new actions
  - `Jarvis/handlers/server.go` - Tool definitions with updated enums
  - `Jarvis/handlers/handlers_test.go` - Test coverage for new handlers
  - `Jarvis/testing/mocks/docker_mock.go` - Mock implementations for new methods
  - `AGENTS.md` - Documentation updates
  - `docs/CONFIGURATION_STRATEGY.md` - Documentation updates

## Success Criteria

1. AI agents can perform ALL common operations via Jarvis tools
2. Zero required shell fallbacks for standard workflows
3. Full test coverage for new actions (93 tasks defined)
4. Documentation updated with new capabilities
5. Payload size remains under 7000 bytes
6. Tool descriptions under 150 characters

## Risks and Mitigations

| Risk | Mitigation |
|:-----|:-----------|
| Test timeout | Use generous timeouts (up to 20 min for full suite) |
| Docker unavailable | Return clear error message with troubleshooting steps |
| Log size | Cap at 500 lines maximum |
| Breaking changes | All changes are additive; backward compatible |
