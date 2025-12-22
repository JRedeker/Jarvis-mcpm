# Tasks: Add Granular Client Editing

## 1. Tool Definition Updates
- [ ] 1.1 Add `set_enabled` parameter to `jarvis_client` tool schema
- [ ] 1.2 Add `set_env` parameter to `jarvis_client` tool schema
- [ ] 1.3 Add `set_url` parameter to `jarvis_client` tool schema
- [ ] 1.4 Add `set_header` parameter to `jarvis_client` tool schema
- [ ] 1.5 Add `remove_env` parameter to `jarvis_client` tool schema
- [ ] 1.6 Add `remove_header` parameter to `jarvis_client` tool schema

## 2. OpenCode Handler Implementation
- [ ] 2.1 Implement `setServerEnabled()` in opencode.go
- [ ] 2.2 Implement `setServerEnv()` in opencode.go
- [ ] 2.3 Implement `setServerUrl()` in opencode.go
- [ ] 2.4 Implement `setServerHeader()` in opencode.go
- [ ] 2.5 Implement `removeServerEnv()` in opencode.go
- [ ] 2.6 Implement `removeServerHeader()` in opencode.go
- [ ] 2.7 Integrate new functions into `openCodeEdit()` handler

## 3. CLI Fallback for Non-OpenCode Clients
- [ ] 3.1 Map new parameters to MCPM CLI arguments for claude-code
- [ ] 3.2 Map new parameters to MCPM CLI arguments for claude-desktop
- [ ] 3.3 Handle clients without CLI support gracefully

## 4. Security & Validation
- [ ] 4.1 Implement value redaction in logs for env vars and headers
- [ ] 4.2 Validate server name exists before modification
- [ ] 4.3 Validate URL format for `set_url`
- [ ] 4.4 Add input sanitization for env var names

## 5. Testing
- [ ] 5.1 Unit tests for `setServerEnabled()`
- [ ] 5.2 Unit tests for `setServerEnv()`
- [ ] 5.3 Unit tests for `setServerUrl()`
- [ ] 5.4 Unit tests for `setServerHeader()`
- [ ] 5.5 Unit tests for `removeServerEnv()`
- [ ] 5.6 Unit tests for `removeServerHeader()`
- [ ] 5.7 Integration test: full edit workflow
- [ ] 5.8 Test error cases (server not found, invalid URL, etc.)

## 6. Documentation
- [ ] 6.1 Update `docs/API_REFERENCE.md` with new parameters
- [ ] 6.2 Update `docs/tech/opencode-integration.md` with examples
- [ ] 6.3 Update `AGENTS.md` with granular editing examples
