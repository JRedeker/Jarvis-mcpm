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
- `tools.go` - Tool definitions and handlers for all Jarvis capabilities
- `go.mod` - Go dependencies (requires Go 1.24+)

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

## Important Jarvis Tools

Jarvis exposes these key tools to AI agents (defined in `tools.go`):

- `bootstrap_system()` - Self-initialization (installs MCPM, starts Docker)
- `apply_devops_stack(project_type, force, enable_ai_review)` - Scaffolds projects with linting, pre-commit hooks, CI/CD
- `analyze_project()` - Returns JSON report of detected languages and configs
- `install_server(name)` - Installs MCP servers via `mcpm install`
- `search_servers(query)` - Searches available MCP servers
- `manage_client(action, client_name, ...)` - Manages client configurations with path persistence
- `manage_profile(action, name, ...)` - Manages MCPM profiles
- `restart_infrastructure()` - Reboots Docker containers
- `check_status()` - System diagnostics
- `fetch_diff_context(staged)` - Returns git status and diff for self-review
- `suggest_profile(client_name, testing)` - Determines which profiles to activate

## Configuration Strategy

### Never Put Jarvis in MCPM Profiles
Jarvis must be configured directly in client config files pointing to the binary. This ensures management capabilities remain available even if profiles break.

### Example Client Configuration
```json
{
  "mcpServers": {
    "mcpm_jarvis": {
      "command": "/path/to/Jarvis/jarvis",
      "args": []
    },
    "mcpm_profile_project-pokeedge": {
      "command": "mcpm",
      "args": ["profile", "run", "project-pokeedge"]
    },
    "mcpm_profile_memory": {
      "command": "mcpm",
      "args": ["profile", "run", "memory"]
    }
  }
}
```

### Standard Config Locations (Linux)
- Claude CLI: `~/.claude.json`
- Claude Desktop/VSCode: `~/.config/Claude/claude_desktop_config.json`

## Development Guidelines

### Go (Jarvis)
- Single-file architecture: `main.go` contains server setup, `tools.go` contains all tool logic
- Uses `github.com/mark3labs/mcp-go` for MCP protocol
- Logging to `logs/jarvis.log` (auto-created)
- Process management for shared servers uses `sync.Mutex` for thread safety

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
1. Run `check_status()` for diagnostics
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
