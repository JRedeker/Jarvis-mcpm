# Jarvis: Server Component

**Version:** 3.0.0 (Context Efficiency Edition)
**Language:** Go 1.24+
**MCP SDK:** [mcp-go v0.43.2](https://github.com/mark3labs/mcp-go)
**Last Updated:** December 2025

## Overview

This directory contains the source code for **Jarvis**, the Go-based MCP server that acts as the gateway between your AI Agent and the underlying system infrastructure.

While the [Root README](../README.md) covers high-level usage and tools, this document focuses on the internal development and structure of the Go application.

## 📂 Directory Structure

*   **`main.go`**: Entry point for the MCP server using `github.com/mark3labs/mcp-go`.
*   **`handlers/`**: Tool handlers with dependency injection for testability
    *   `server.go` - Consolidated tool definitions (8 tools in v3.0)
    *   `consolidated.go` - Action-based routing handlers
    *   `handlers.go` - Core business logic implementations
    *   `registry.go` - Handler registration
*   **`testing/`**: Test utilities (mocks, helpers, fixtures)
*   **`Dockerfile`**: Container image for deployed environments.

## 🛠️ Building from Source

### Prerequisites
*   Go 1.24 or higher
*   Make (optional, standard Go commands work fine)

### Compilation
Jarvis is built as a static binary.

```bash
go mod download
go build -o jarvis .
```

### Running Locally
You can run the binary directly to test the Stdio transport (though it will expect MCP JSON-RPC messages on stdin).

```bash
./jarvis
```

## 🧩 Key Implementation Details

### Consolidated Tool Architecture (v3.0)

In v3.0, Jarvis consolidated 24 tools into 8 action-based tools for **52% context token reduction**:

| Tool | Actions | Description |
|:-----|:--------|:------------|
| `jarvis_check_status` | - | System health check |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage | Server management |
| `jarvis_profile` | list, create, edit, delete, suggest, restart | Profile management |
| `jarvis_client` | list, edit, import, config | Client configuration |
| `jarvis_config` | get, set, list, migrate | MCPM configuration |
| `jarvis_project` | analyze, diff, devops | Project analysis & DevOps |
| `jarvis_system` | bootstrap, restart, restart_infra | System operations |
| `jarvis_share` | start, stop, list | Server sharing |

### `jarvis_system` (action="bootstrap")
Self-initialization tool:
1.  Locates the project root by looking for the `MCPM/` directory.
2.  Runs `npm install` and `npm link` inside `MCPM/`.
3.  Runs `docker-compose up -d` in the root.

### `jarvis_project` (action="devops")
Transforms Jarvis into a Project Architect.
*   **Input:** `project_type` (optional), `enable_ai_review` (bool), `force` (bool).
*   **Safe Mode:** Checks for existing configs first. Use `force=true` to overwrite.
*   **Actions:**
    1.  **Git:** Runs `git init`.
    2.  **Pre-Commit:** Generates `.pre-commit-config.yaml` tailored to the language.
    3.  **Hooks:** Runs `pip install pre-commit && pre-commit install`.
    4.  **AI Review:** Generates `.github/workflows/pr_agent.yml`.
    5.  **Ignore:** Creates `.gitignore`.

### `jarvis_project` (action="analyze")
Enables intelligent decision making for agents.
*   **Output:** JSON report of detected languages (`go`, `python`, `node`) and existing configurations (Git, Pre-commit, Workflows).
*   **Use Case:** Agents call this *before* applying the stack to determine the correct strategy.

### `jarvis_system` (action="restart_infra")
Self-healing capability for the environment.
*   **Action:** Executes `scripts/manage-mcp.sh restart`.
*   **Result:** Safely reboots the Postgres and Qdrant containers and logs the output.

### `jarvis_project` (action="diff")
Enables the "Local Review Loop" for agents.
*   **Input:** `staged` (bool).
*   **Output:** A formatted Markdown report containing:
    *   Current Working Directory.
    *   `git status` output.
    *   `git diff` content (staged or HEAD).
*   **Use Case:** Allows the Agent to "see" its own changes before committing, enabling self-correction.

### `jarvis_profile` (action="restart")
Manages the `mcpm-daemon` container, allowing hot-reloads of MCP servers.
*   **Input:** `profile` (optional string).
*   **Action:**
    *   If `profile` provided: Runs `docker exec mcp-daemon supervisorctl restart mcpm-<profile>`.
    *   If empty: Runs `docker compose restart mcpm-daemon`.
*   **Benefit:** Updates API keys or server configs without disconnecting the AI client.

### `jarvis_profile` (action="suggest")
Jarvis implements a "3-Layer Stacking" logic to determine the active toolset dynamically.
*   **Input:** `testing` (bool).
*   **Logic:**
    1.  **Base:** Detects project context (e.g., `p-pokeedge`) from CWD. Defaults to `p-new`.
    2.  **Global:** Appends `memory`. Appends `testing-all-tools` if `testing=true`.
*   **Output:** JSON array of profile names (which correspond to HTTP ports).

### `jarvis_client` (Client Configuration Management)
Jarvis exposes advanced configuration for MCP clients.
*   **Actions:** `list`, `edit`, `import`, `config`.
*   **HTTP Support:** Can configure clients to point to local HTTP endpoints (`http://localhost:XXXX/mcp`) instead of spawning stdio processes.

### `jarvis_share` (Server Tunneling)
Jarvis manages a map of running background processes for server sharing.
*   **Concurrency:** Uses `sync.Mutex` to safely manage the state of shared tunnels.
*   **Process Management:** Captures `stdout` to detect the public URL generated by `mcpm share`.

## 🧪 Testing

Jarvis includes a comprehensive test suite covering all critical functionality:

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests with race detection
go test -v -race ./...
```

### Test Coverage

**6 test functions covering 23+ test cases:**

1. **TestSetupLogging** - Verifies logging initialization and file creation
2. **TestBuildManageClientArgs** - Tests argument building for client management (5 scenarios)
3. **TestHandlersBasicExecution** - Tests tools without arguments (4 tools)
   - `handleListServers`
   - `handleCheckStatus`
   - `handleUsageStats`
   - `handleListSharedServers`
4. **TestHandlersWithArguments** - Tests tools with parameters (3 tools)
   - `handleInstallServer`
   - `handleSearchServers`
   - `handleServerInfo`
5. **TestMonitorShareProcess** - Tests URL extraction from share output (3 scenarios)

### Pre-Commit Hooks

Jarvis enforces quality standards via pre-commit hooks:
- **Go formatting** (`go fmt`)
- **Secret detection** (`gitleaks`)
- **Trailing whitespace** removal
- **End of file** fixes
- **Large file** detection
- **Merge conflict** markers

### CI/CD

GitHub Actions workflow runs on every push:
- Linting with `golangci-lint`
- Full test suite with race detection
- Build verification

## 🎯 Version 3.0: Context Efficiency Edition

The v3.0 release consolidates 24 tools into 8 action-based tools for **52% context token reduction** (~1,400 tokens saved per connection).

### Key Changes

- **Tool Consolidation:** 24 → 8 tools using action-based routing
- **Namespace:** All tools prefixed with `jarvis_`
- **Payload Reduction:** ~11KB → ~5.3KB
- **Breaking Change:** Old tool names replaced (see migration guide in AGENTS.md)

### Smart Error Handling

All consolidated tools include **intelligent validation and helpful suggestions**:

**jarvis_server (action="install"):**
- ❌ Validates server name format (rejects spaces/slashes)
- 💡 Suggests `jarvis_server(action="search")` on 404 errors
- ✅ Provides next steps after installation
- Detects "already installed" and suggests profile management

**jarvis_server (action="search"):**
- ❌ Validates non-empty query
- 💡 Provides tips when no results found
- 💡 Adds next step guidance to use `jarvis_server(action="info")`

**jarvis_client:**
- ❌ Validates action against allowed values (list, edit, import, config)
- 💡 Provides examples for each action type
- 💡 Lists common client names when required
- ✅ Adds next step guidance after successful edits

### Documentation Updates

- **AGENTS.md** includes full migration guide (old → new tool mapping)
- **CLAUDE.md** updated with v3.0 consolidated tool reference
- Clear guidance on when to use Jarvis vs direct MCPM CLI

## 📦 Dependencies

| Package | Version | Description |
|:--------|:--------|:------------|
| [mcp-go](https://github.com/mark3labs/mcp-go) | v0.43.2 | MCP server SDK |
| [testify](https://github.com/stretchr/testify) | v1.11.1 | Testing framework |
| [go-cmp](https://github.com/google/go-cmp) | v0.7.0 | Value comparison |

## 🐳 Docker

Build and run Jarvis as a container:

```bash
# Build image
make docker

# Or manually
docker build -t jarvis:3.0.0 .

# Run (stdio mode)
docker run --rm -it jarvis:latest
```

## 📄 License

MIT License - See [LICENSE](../LICENSE) for details.
