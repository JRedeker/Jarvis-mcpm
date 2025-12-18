# Tasks: Reduce Context Tokens

## 1. Core Implementation

### 1.1 Tool Definition Refactoring
- [x] 1.1.1 Create new consolidated tool definitions in `handlers/server.go`
  - Define `jarvis_check_status` tool (no action param)
  - Define `jarvis_server` tool with action enum
  - Define `jarvis_profile` tool with action enum
  - Define `jarvis_client` tool with action enum
  - Define `jarvis_config` tool with action enum
  - Define `jarvis_project` tool with action enum
  - Define `jarvis_system` tool with action enum
  - Define `jarvis_share` tool with action enum

- [x] 1.1.2 Set appropriate annotations for each tool
  - Mark read-only tools: `jarvis_check_status`, list/info/search/suggest actions
  - Mark idempotent tools: list/info/status/analyze actions
  - Ensure all tools have meaningful annotation values

- [x] 1.1.3 Write concise descriptions (max 60 chars each)
  - Include action list in description for multi-action tools
  - Remove verbose explanatory text

### 1.2 Handler Implementation
- [x] 1.2.1 Create action router in `handlers/consolidated.go`
  - Implemented consolidated handlers with switch-based routing
  - Return helpful error messages for invalid actions

- [x] 1.2.2 Implement `Server()` consolidated handler
  - Route to existing `listServers`, `serverInfo`, etc. based on action
  - Validate action parameter, return helpful error if invalid

- [x] 1.2.3 Implement `Profile()` consolidated handler
  - Route list/create/edit/delete to `ManageProfile()` with mapped action
  - Route suggest to `SuggestProfile()`
  - Route restart to `RestartProfiles()`

- [x] 1.2.4 Implement `Client()` consolidated handler
  - Route all actions to `ManageClient()` with mapped action

- [x] 1.2.5 Implement `Config()` consolidated handler
  - Route get/set/list to `ManageConfig()` with mapped action
  - Route migrate to `MigrateConfig()`

- [x] 1.2.6 Implement `Project()` consolidated handler
  - Route analyze to `AnalyzeProject()`
  - Route diff to `FetchDiffContext()`
  - Route devops to `ApplyDevOpsStack()`

- [x] 1.2.7 Implement `System()` consolidated handler
  - Route bootstrap to `BootstrapSystem()`
  - Route restart to `RestartService()`
  - Route restart_infra to `RestartInfrastructure()`

- [x] 1.2.8 Implement `Share()` consolidated handler
  - Route start to `ShareServer()`
  - Route stop to `StopSharingServer()`
  - Route list to `ListSharedServers()`

### 1.3 Registration Update
- [x] 1.3.1 Update `GetToolDefinitions()` to return only 8 tools
- [x] 1.3.2 Remove old tool definitions
- [x] 1.3.3 Update `RegisterToolsWithMCPServer()` if needed

## 2. Testing

### 2.1 Unit Tests
- [x] 2.1.1 Add tests for action routing logic
  - Test valid actions route correctly
  - Test invalid actions return helpful errors
  - Test missing required params handled

- [x] 2.1.2 Update existing handler tests
  - Existing tests still pass (they test the internal handlers)
  - Added new tests for consolidated handlers

- [x] 2.1.3 Add payload size verification test
  - Test that ListToolsResult JSON is < 6KB (achieved 5.3KB)
  - Test description lengths are <= 80 chars
  - Test exactly 8 tools returned

### 2.2 Integration Tests
- [x] 2.2.1 Test full tool listing response
  - Verify exactly 8 tools returned
  - Verify all tool names start with `jarvis_`

- [x] 2.2.2 Test each consolidated tool end-to-end
  - Verify each action works through the new interface

## 3. Documentation

### 3.1 Update AGENTS.md
- [x] 3.1.1 Document new tool names and actions
  - Add table mapping old tools to new action-based calls
  - Update usage examples

- [x] 3.1.2 Add migration notes
  - Document breaking changes
  - Provide before/after examples

### 3.2 Update README
- [ ] 3.2.1 Update tool listing in README if present
- [ ] 3.2.2 Add context token reduction as a feature

## 4. Validation

### 4.1 Metrics Verification
- [x] 4.1.1 Measure actual payload size reduction
  - Before: ~11KB (~2,750 tokens)
  - After: ~5.3KB (~1,325 tokens)
  - **Achieved 52% reduction** (saves ~1,400 tokens per connection)

- [x] 4.1.2 Verify functionality preservation
  - Run full test suite: All 110 tests pass
  - Action routing verified for all 8 consolidated tools

## Dependencies

- Tasks 1.1.x can be done in parallel
- Tasks 1.2.x depend on 1.1.x completion
- Tasks 2.x depend on 1.x completion
- Tasks 3.x can start after 1.x is stable
- Task 4.x is final validation

## Estimated Effort

| Section | Estimate |
|---------|----------|
| 1.1 Tool Definitions | 1-2 hours |
| 1.2 Handler Implementation | 2-3 hours |
| 1.3 Registration | 30 min |
| 2.1 Unit Tests | 2-3 hours |
| 2.2 Integration Tests | 1-2 hours |
| 3.x Documentation | 1 hour |
| 4.x Validation | 30 min |
| **Total** | **8-12 hours** |
