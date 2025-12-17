## REMOVED Requirements

### Requirement: SSE Transport Type Support
**Reason**: SSE is deprecated in MCP spec 2025-03-26. Complete removal simplifies the CLI and prevents user confusion.
**Migration**: Users must use `--type streamable-http` instead of `--type sse`.

#### Scenario: SSE type accepted
- **WHEN** user ran `mcpm new server --type sse --url <url>`
- **THEN** the server was registered with SSE transport
- **STATUS**: Removed - `sse` is no longer a valid type

## MODIFIED Requirements

### Requirement: Valid Transport Types
The `mcpm new` command SHALL only accept `stdio`, `http`, and `streamable-http` as valid transport types.

#### Scenario: Valid types enumeration
- **WHEN** the CLI validates the `--type` parameter
- **THEN** it SHALL accept: `stdio`, `http`, `streamable-http`
- **AND** it SHALL reject `sse` with an error

#### Scenario: SSE type rejection with helpful error
- **WHEN** user runs `mcpm new server --type sse --url <url>`
- **THEN** the CLI SHALL return an error
- **AND** the error message SHALL state that SSE has been removed
- **AND** the error message SHALL suggest using `streamable-http` instead

### Requirement: Profile Run Flags
The `mcpm profile run` and `mcpm run` commands SHALL only support `--http` flag for network transport.

#### Scenario: HTTP flag only
- **WHEN** running a profile or server in network mode
- **THEN** the `--http` flag SHALL be the only transport option
- **AND** the `--sse` flag SHALL NOT be recognized

#### Scenario: SSE flag rejection
- **WHEN** user runs `mcpm profile run --sse myprofile`
- **THEN** the CLI SHALL return an error indicating `--sse` is no longer supported
