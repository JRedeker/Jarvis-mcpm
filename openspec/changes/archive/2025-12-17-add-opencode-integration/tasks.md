# Tasks: Add OpenCode Client Integration

## 1. Client Registry
- [x] 1.1 Add OpenCode entry to MCPM client registry in `MCPM/index.js`
- [x] 1.2 Define config paths: `$OPENCODE_CONFIG`, `./opencode.json`, `~/.config/opencode/opencode.json`
- [x] 1.3 Add format identifier `opencode-mcp` for config generation

## 2. Configuration Template
- [x] 2.1 Create `config-templates/opencode.json` with Jarvis (local) + memory profile (remote)
- [x] 2.2 Include schema reference `https://opencode.ai/config.json`
- [x] 2.3 Use placeholder `${JARVIS_PATH}` for Jarvis binary path

## 3. Jarvis Handler Updates
- [x] 3.1 Update `handlers.go` client detection to recognize OpenCode config paths
- [x] 3.2 Implement OpenCode config format writer (mcp object structure)
- [x] 3.3 Add OpenCode to `manage_client` action `config` path detection
- [x] 3.4 Add OpenCode to `manage_client` action `edit` for adding/removing profiles

## 4. Testing
- [x] 4.1 Add unit tests for OpenCode config detection in `handlers_test.go`
- [x] 4.2 Add unit tests for OpenCode config format generation
- [x] 4.3 Add integration test with mock OpenCode config file

## 5. Documentation
- [x] 5.1 Update `AGENTS.md` with OpenCode client reference paths
- [x] 5.2 Create `docs/tech/opencode-integration.md` with setup instructions
- [x] 5.3 Update `docs/CONFIGURATION_STRATEGY.md` to mention OpenCode support

## 6. Validation
- [x] 6.1 Manual test: Create OpenCode config using Jarvis tools
- [x] 6.2 Manual test: Verify Jarvis stdio connection works from OpenCode
- [x] 6.3 Manual test: Verify HTTP profile connections work from OpenCode
- [x] 6.4 Run full test suite: `go test -v ./...` and `bats scripts/tests/`
