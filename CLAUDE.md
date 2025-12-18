<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **Jarvis**, an intelligent infrastructure layer for AI agents built on the Model Context Protocol (MCP). It transforms AI agents from passive chatbots into full-stack DevOps engineers by providing tool management, guardrails, memory, and research capabilities.

The repository contains three main components:
- **Jarvis** (Go): MCP server gateway that exposes system management tools to AI agents
- **MCPM** (Node.js): Package manager for installing and configuring MCP servers
- **Shared Infrastructure** (Docker): PostgreSQL and Qdrant vector database

## Core Architecture

### 3-Layer Profile Stack
The system uses a composable configuration architecture with three layers:

1. **Layer 1 (Environment)**: `project-<name>` - Workspace-specific tools (e.g., `project-pokeedge`, `project-new`)
2. **Layer 2 (Client Adapter)**: `client-<name>` - AI client-specific tools (e.g., `client-codex`)
3. **Layer 3 (Global)**: `memory`, `testing-all-tools` - Cross-cutting capabilities

Jarvis determines which profiles to activate based on context (working directory, client name, testing mode).

### Component Relationships
```
AI Agent → Jarvis (MCP Server) → MCPM CLI → Docker Infrastructure
                                    ↓
                            Server Registry (technologies.toml)
```

Jarvis acts as a **Presentation Layer**, capturing raw CLI output, stripping ANSI codes, and formatting responses in clean Markdown with status emojis (✅/❌).

## Common Commands

### Building and Running

**Build Jarvis:**
```bash
cd Jarvis
go build -o jarvis .
```

**Setup Jarvis (builds and shows config):**
```bash
./scripts/setup-jarvis.sh
```

**Manage Docker infrastructure:**
```bash
./scripts/manage-mcp.sh start    # Start PostgreSQL and Qdrant
./scripts/manage-mcp.sh stop     # Stop containers
./scripts/manage-mcp.sh restart  # Restart containers
./scripts/manage-mcp.sh status   # Check container status
./scripts/manage-mcp.sh logs     # View container logs
```

**Install MCPM CLI:**
```bash
cd MCPM
npm install
npm link  # Makes 'mcpm' command available globally
```

### Testing

**Run all tests:**
```bash
./scripts/manage-mcp.sh test  # Runs both Go and Python tests
```

**Run Go tests only:**
```bash
cd Jarvis
go test -v ./...
```

**Run Python tests (mcpm_source):**
```bash
cd mcpm_source
uv run pytest
```

### Development

**Format Go code:**
```bash
cd Jarvis
gofmt -w .
```

**Format Python code:**
```bash
cd mcpm_source
uv run ruff format .
```

**Lint Go code:**
```bash
cd Jarvis
golangci-lint run
```

## Key Files and Directories

### Jarvis/ (Go MCP Server)
- `main.go` - Entry point, MCP server setup, logging configuration
- `handlers/` - Tool handlers with dependency injection for testability
  - `handlers.go` - Core handler implementations
  - `server.go` - MCP tool definitions
  - `registry.go` - Handler registration
- `testing/` - Test utilities (mocks, helpers, fixtures)
- `smoketests/` - Integration smoke tests
- `go.mod` - Go dependencies (requires Go 1.23+)

### MCPM/ (Node.js CLI)
- `index.js` - CLI entry point using Commander
- `package.json` - Node dependencies (commander, toml, chalk, etc.)
- Note: This is a streamlined fork optimized for the Jarvis ecosystem

### mcpm_source/ (Python Reference)
- Contains original Python MCPM source for reference only
- Not part of active build pipeline
- Uses `uv` for dependency management

### Infrastructure
- `docker-compose.yml` - Defines PostgreSQL and Qdrant services
- `scripts/` - Management utilities for setup and operations
- `.env.template` - Environment variable template

### Documentation
- `ARCHITECTURE-MAP.md` - System hierarchy and component roles
- `AGENTS.md` - Instructions for AI agents using Jarvis
- `docs/CONFIGURATION_STRATEGY.md` - 3-Layer Profile Stack details
- `docs/TECHNICAL_ARCHITECTURE.md` - Component specifications

## Using Jarvis (Primary Interface - ALWAYS PREFER THIS)

**Jarvis is the intelligent gateway to MCPM** designed specifically for AI agents. It provides:

✅ **Clean, Structured Output** - No ANSI codes, stripped warnings, formatted responses
✅ **Smart Error Handling** - Actionable error messages with suggestions for fixes
✅ **Validation & Safety** - Prevents common mistakes before they happen
✅ **Context Awareness** - Understands your project and suggests appropriate actions
✅ **Batch Operations** - Complex workflows that would require multiple CLI commands

### When to Use Jarvis vs Direct MCPM

**✓ ALWAYS use Jarvis tools:**
- Installing servers → `jarvis_server(action="install", name="...")` not `mcpm install`
- Managing profiles → `jarvis_profile(action="edit", ...)` not `mcpm profile`
- Checking system health → `jarvis_check_status()` not `mcpm doctor`
- Configuring clients → `jarvis_client(action="edit", ...)` not `mcpm client`
- Searching servers → `jarvis_server(action="search", query="...")` not `mcpm search`

**⚠️ Only use direct MCPM CLI when:**
- Debugging Jarvis itself
- Running scripts outside of an AI agent context
- Following documentation that predates Jarvis

### Quick Reference

| Task | Use Jarvis Tool | Not Direct CLI |
|------|----------------|----------------|
| Install MCP server | `jarvis_server(action="install", name="brave-search")` | ~~`mcpm install brave-search`~~ |
| Check system health | `jarvis_check_status()` | ~~`mcpm doctor`~~ |
| Add to profile | `jarvis_profile(action="edit", name="p-pokeedge", add_servers="...")` | ~~`mcpm profile edit`~~ |
| Search available tools | `jarvis_server(action="search", query="documentation")` | ~~`mcpm search documentation`~~ |
| Bootstrap environment | `jarvis_system(action="bootstrap")` | ~~Multiple manual steps~~ |

## Jarvis Tool Reference

Jarvis exposes **8 consolidated MCP tools** (v3.0) for context token efficiency. All handlers are defined in `handlers/server.go` with implementations in `handlers/consolidated.go`.

### Complete Tool Reference Table

| Tool | Actions | Key Parameters |
|:-----|:--------|:---------------|
| `jarvis_check_status` | - | System health check for MCPM, Docker, services |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage | `action` (required), `name`, `query`, `type`, `command`, `args`, `env`, `url`, `headers` |
| `jarvis_profile` | list, create, edit, delete, suggest, restart | `action` (required), `name`, `new_name`, `add_servers`, `remove_servers`, `profile`, `testing` |
| `jarvis_client` | list, edit, import, config | `action` (required), `client_name`, `add_server`, `remove_server`, `add_profile`, `remove_profile`, `config_path` |
| `jarvis_config` | get, set, list, migrate | `action` (required), `key`, `value` |
| `jarvis_project` | analyze, diff, devops | `action` (required), `staged`, `project_type`, `force`, `enable_ai_review` |
| `jarvis_system` | bootstrap, restart, restart_infra | `action` (required) |
| `jarvis_share` | start, stop, list | `action` (required), `name`, `port`, `no_auth` |

### Usage Examples

```javascript
// Install a server
jarvis_server({ action: "install", name: "brave-search" })

// List all profiles
jarvis_profile({ action: "list" })

// Configure OpenCode client
jarvis_client({ action: "edit", client_name: "opencode", add_profile: "memory" })

// Bootstrap entire system
jarvis_system({ action: "bootstrap" })

// Analyze current project
jarvis_project({ action: "analyze" })
```

> See full API documentation: `docs/API_REFERENCE.md`

## Configuration Strategy

### Never Put Jarvis in MCPM Profiles
Jarvis must be configured directly in client config files pointing to the binary. This ensures management capabilities remain available even if profiles break.

### Example Client Configuration
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/path/to/Jarvis/jarvis",
      "args": []
    },
    "p-pokeedge": {
      "url": "http://localhost:6276/mcp",
      "transport": "streamable-http"
    },
    "memory": {
      "url": "http://localhost:6277/mcp",
      "transport": "streamable-http"
    }
  }
}
```

### Standard Config Locations (Linux)
- Claude CLI: `~/.claude.json`
- Claude Desktop/VSCode: `~/.config/Claude/claude_desktop_config.json`

## Development Guidelines

### Go (Jarvis)
- **Package architecture**: `main.go` for server setup, `handlers/` package for tool logic
- **Dependency Injection**: All handlers accept interfaces for testability
  - `McpmRunner` - MCPM CLI operations
  - `DockerRunner` - Docker compose operations
  - `GitRunner` - Git operations
  - `FileSystem` - File system operations
  - `CommandRunner` - Shell command execution
  - `ProcessManager` - Shared server process management
- Uses `github.com/mark3labs/mcp-go` for MCP protocol
- Logging to `logs/jarvis.log` (auto-created)
- **Testing**: Mocks in `testing/mocks/`, helpers in `testing/helpers/`
- Run tests: `go test -v ./...`
- Run with coverage: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

### Node.js (MCPM)
- Commander-based CLI with subcommands
- Reads registry from `config/technologies.toml`
- Outputs IDE-specific config files (Cursor, Cline)
- Prefers Docker installation when available

### Python (mcpm_source)
- Reference implementation only
- Uses `uv` for dependency management (NOT pip/venv)
- Format with `ruff`, follows semantic commit messages
- NEVER commit without explicit user request

## Git Workflow

This project follows **semantic versioning** with conventional commits:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks

**BREAKING CHANGE:** in commit footer triggers major version bump (rare - requires explicit user approval).

## Infrastructure Services

### PostgreSQL
- Container: `mcp-postgres`
- Port: 5432
- User: `mcp` / Password: `mcp_password` / DB: `mcp_db`
- Health check: `pg_isready`

### Qdrant (Vector Store)
- Container: `mcp-qdrant`
- Ports: 6333 (HTTP), 6334 (gRPC)
- Used for semantic search and memory

## Debugging

If Jarvis tools fail:
1. Run `jarvis_check_status()` for diagnostics
2. Check logs: `logs/jarvis.log`, `logs/management.log`
3. For MCPM issues, set environment variables:
   - `export MCPM_NON_INTERACTIVE=true`
   - `export MCPM_FORCE=true`
4. Verify containers: `docker compose ps`
5. View container logs: `./scripts/manage-mcp.sh logs`

## Key Design Principles

- **Agent-First**: Jarvis transforms agents into self-sufficient DevOps engineers
- **Security**: Pre-commit hooks block secrets, enforce formatting standards
- **Composability**: 3-Layer Stack prevents config duplication
- **Self-Healing**: Agents can diagnose and repair infrastructure issues
- **Dynamic Expansion**: Agents can install new tools on-demand via MCPM registry
