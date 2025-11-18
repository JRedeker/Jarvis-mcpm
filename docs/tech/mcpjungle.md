# MCPJungle Technology Documentation

**GitHub**: https://github.com/mcpjungle/MCPJungle
**Stars**: ~694
**Status**: ✅ Selected as Primary MCP Hub
**Purpose**: Self-hosted MCP Gateway and Registry for AI agents
**Instance Name**: `jarvis` (our MCPJungle hub)
**Latest Release**: 0.2.16

---

## Overview

MCPJungle is a self-hosted MCP Gateway that serves as a single source-of-truth registry for all Model Context Protocol servers in an organization.

**Key Features:**
- 📦 **Unified Gateway**: Single MCP endpoint for all registered servers
- 🔧 **Tool Groups**: Organize and expose subsets of tools to different clients
- 🔐 **Enterprise Mode**: Access control, authentication, OpenTelemetry
- 🚀 **CLI Management**: Register, deregister, enable/disable tools via CLI
- ⚡ **Streamable HTTP**: Primary transport, also supports STDIO
- 🏗️ **Go-based**: Written in Go (99%), production-ready performance

**Architecture:**
```
┌─────────────────────────────────────────────┐
│       MCP Clients (Claude, Cursor, etc)      │
│         ws://localhost:8080/mcp              │
└────────────────────┬────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────┐
│         MCPJungle Gateway (jarvis)          │
│  • Tool registry & discovery                │
│  • Tool groups & access control             │
│  • Authentication & authorization           │
│  • OpenTelemetry metrics                    │
└─────────┬───────────────────┬───────────────┘
          │                   │
    ┌─────▼──────┐     ┌─────▼──────┐
    │  context7   │     │  filesystem │
    │ (HTTP)      │     │  (STDIO)    │
    └─────────────┘     └─────────────┘
```

---

## Installation

### Homebrew (Recommended for macOS)
```bash
brew install mcpjungle/mcpjungle/mcpjungle

# Verify installation
mcpjungle version
```

### Direct Binary Download
Download from [GitHub Releases](https://github.com/mcpjungle/MCPJungle/releases)

> **Note**: On macOS, Homebrew is required because binaries are not yet Notarized.

### Docker
```bash
docker pull mcpjungle/mcpjungle

# Standard image (minimal, production-ready)
docker pull mcpjungle/mcpjungle:latest

# STDIO image (includes npx, uvx for local STDIO servers)
docker pull mcpjungle/mcpjungle:latest-stdio
```

---

## Quick Start

### 1. Start Server (Docker Compose)
```bash
# Download docker-compose.yaml
curl -O https://raw.githubusercontent.com/mcpjungle/MCPJungle/refs/heads/main/docker-compose.yaml

# Start server + PostgreSQL
docker compose up -d

# Verify server
curl http://localhost:8080/health
```

### 2. Install CLI
```bash
brew install mcpjungle/mcpjungle/mcpjungle
```

### 3. Register MCP Server
```bash
# Register remote HTTP server (e.g., context7)
mcpjungle register --name context7 --url https://mcp.context7.com/mcp

# List registered tools
mcpjungle list tools

# Test tool invocation
mcpjungle invoke context7__get-library-docs --input '{"library": "lodash/lodash"}'
```

### 4. Connect IDE Client

**Claude Desktop:**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "command": "npx",
      "args": ["mcp-remote", "http://localhost:8080/mcp", "--allow-http"]
    }
  }
}
```

**Cursor:**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

---

## Configuration

### Server Modes

**Development Mode** (default for local use):
```bash
# Starts with SQLite, no auth required
mcpjungle start
```

**Enterprise Mode** (production deployments):
```bash
# Enforces auth, access control, OpenTelemetry
mcpjungle start --enterprise

# Or via environment variable
export SERVER_MODE=enterprise
mcpjungle start

# Or via docker-compose.prod.yaml
curl -O https://raw.githubusercontent.com/mcpjungle/MCPJungle/refs/heads/main/docker-compose.prod.yaml
docker compose -f docker-compose.prod.yaml up -d
```

### Database Configuration

**SQLite** (default for local dev):
```bash
# Creates mcpjungle.db in current directory
mcpjungle start
```

**PostgreSQL** (recommended for production):
```bash
# Via DSN
export DATABASE_URL=postgres://admin:root@localhost:5432/mcpjungle_db
mcpjungle start

# Or via individual env vars
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=admin
export POSTGRES_PASSWORD=secret
export POSTGRES_DB=mcpjungle_db
mcpjungle start
```

---

## Registering MCP Servers

### Streamable HTTP Servers

**Command-line registration:**
```bash
mcpjungle register \
  --name context7 \
  --description "Documentation lookup via llms.txt" \
  --url https://mcp.context7.com/mcp
```

**Configuration file:**
```json
{
  "name": "context7",
  "transport": "streamable_http",
  "description": "Documentation lookup",
  "url": "https://mcp.context7.com/mcp",
  "bearer_token": "optional-api-token"
}
```

```bash
mcpjungle register -c ./context7.json
```

### STDIO Servers

**Configuration file example** (`filesystem.json`):
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "Filesystem MCP server",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "."],
  "env": {
    "NODE_ENV": "production"
  }
}
```

```bash
mcpjungle register -c ./filesystem.json
```

**Docker filesystem access:**
When running MCPJungle in Docker, mount host directories as volumes:
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/host"]
}
```

> The `docker-compose.yaml` mounts `$PWD` as `/host` in the container.

### Deregistering Servers
```bash
mcpjungle deregister context7
mcpjungle deregister filesystem
```

---

## Tool Groups

Tool groups allow exposing only specific subsets of tools to different clients.

### Creating Tool Groups

**Example 1: Cherry-pick specific tools**
```json
{
  "name": "claude-tools",
  "description": "Handpicked tools for Claude Desktop",
  "included_tools": [
    "filesystem__read_file",
    "context7__get-library-docs",
    "time__get_current_time"
  ]
}
```

**Example 2: Include entire servers with exclusions**
```json
{
  "name": "research-tools",
  "description": "All tools from time and context7 except convert_time",
  "included_servers": ["time", "context7"],
  "excluded_tools": ["time__convert_time"]
}
```

**Example 3: Mixed approach**
```json
{
  "name": "comprehensive",
  "description": "Mix of manual tools, servers, and exclusions",
  "included_tools": ["filesystem__read_file"],
  "included_servers": ["time"],
  "excluded_tools": ["time__convert_time"]
}
```

### Creating and Using Groups
```bash
# Create group
mcpjungle create group -c ./claude-tools.json

# Output: http://127.0.0.1:8080/v0/groups/claude-tools/mcp

# List all groups
mcpjungle list groups

# View group details
mcpjungle get group claude-tools

# List tools in group
mcpjungle list tools --group claude-tools

# Invoke tool in group context
mcpjungle invoke filesystem__read_file \
  --group claude-tools \
  --input '{"path": "README.md"}'

# Delete group
mcpjungle delete group claude-tools
```

### Group Limitations
- Cannot update existing groups (must delete and recreate)
- Prompts not yet supported in groups
- In enterprise mode, only admins can create groups (standard users coming soon)

---

## Tool Management

### Tool Naming Convention
Tools follow the canonical name pattern: `<server-name>__<tool-name>`

Examples:
- `context7__get-library-docs`
- `filesystem__read_file`
- `time__get_current_time`

### Enabling/Disabling Tools
```bash
# Disable specific tool
mcpjungle disable tool context7__get-library-docs

# Re-enable tool
mcpjungle enable tool context7__get-library-docs

# Disable all tools in a server
mcpjungle disable tool context7

# Disable entire server (all tools + prompts)
mcpjungle disable server context7

# Re-enable server
mcpjungle enable server context7

# Disable a prompt
mcpjungle disable prompt "huggingface__Model Details"
```

> **Note**: New servers are registered with all tools/prompts **enabled** by default.

---

## Enterprise Features

### Access Control

In enterprise mode, create MCP clients with specific access permissions:

```bash
# Initialize server (creates admin user)
mcpjungle init-server

# Create MCP client with access to specific servers
mcpjungle create mcp-client cursor-local --allow "calculator,github"

# Output:
# MCP client 'cursor-local' created successfully!
# Access token: 1YHf2LwE1LXtp5lW_vM-gmdYHlPHdqwnILitBhXE4Aw
```

**Cursor configuration with auth:**
```json
{
  "mcpServers": {
    "mcpjungle": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer 1YHf2LwE1LXtp5lW_vM-gmdYHlPHdqwnILitBhXE4Aw"
      }
    }
  }
}
```

### Authentication

**Static Bearer Tokens** (for SaaS MCP servers):
```bash
# Register with bearer token
mcpjungle register \
  --name huggingface \
  --url https://huggingface.co/mcp \
  --bearer-token <your-hf-api-token>

# Or via config file
{
  "name": "huggingface",
  "transport": "streamable_http",
  "url": "https://huggingface.co/mcp",
  "bearer_token": "<your-hf-api-token>"
}
```

> OAuth flow support coming soon!

### OpenTelemetry

**Enterprise mode** (enabled by default):
```bash
mcpjungle start --enterprise
# Metrics available at http://localhost:8080/metrics
```

**Development mode** (opt-in):
```bash
export OTEL_ENABLED=true
export OTEL_RESOURCE_ATTRIBUTES=deployment.environment.name=dev
mcpjungle start
```

---

## Prompts

MCPJungle supports MCP [Prompts](https://modelcontextprotocol.io/specification/2025-06-18/server/prompts):

```bash
# List all prompts from huggingface server
mcpjungle list prompts --server huggingface

# Get prompt with custom arguments
mcpjungle get prompt "huggingface__Model Details" \
  --arg model_id="openai/gpt-oss-120b"

# Disable/enable prompts
mcpjungle disable prompt "huggingface__Model Details"
mcpjungle enable prompt "huggingface__Model Details"
```

---

## Known Limitations

### 1. Stateless Connections
- MCPJungle creates a new connection for every tool call
- For STDIO servers, a new subprocess is started per tool call
- No support for stateful/long-running connections
- Performance overhead but prevents memory leaks

### 2. Authentication
- Static bearer tokens only (OAuth coming soon)

### 3. Tool Groups
- Cannot update existing groups (must delete/recreate)
- Prompts not yet supported in groups

---

## Production Deployment

### Docker Production Setup
```bash
# Use production docker-compose
curl -O https://raw.githubusercontent.com/mcpjungle/MCPJungle/refs/heads/main/docker-compose.prod.yaml

# Use standard image for remote servers only
docker compose -f docker-compose.prod.yaml up -d

# Or use stdio image if you need local STDIO servers
MCPJUNGLE_IMAGE_TAG=latest-stdio docker compose -f docker-compose.prod.yaml up -d
```

### Recommendations
- Use external PostgreSQL cluster (not Docker container)
- Run in `enterprise` mode
- Enable OpenTelemetry for monitoring
- Use tool groups for per-client access control
- Mount persistent volumes for logs and data

---

## Troubleshooting

### Server won't start
```bash
# Check port availability
lsof -i :8080

# Check logs (Docker)
docker compose logs mcpjungle

# Check database connectivity
mcpjungle start --verbose
```

### STDIO servers fail
- Check mcpjungle server logs for stderr output
- Ensure using `latest-stdio` image if relying on npx/uvx
- Verify command and args are correct
- Check file paths and permissions

### Tools not appearing
```bash
# List all registered servers
mcpjungle list servers

# List all tools
mcpjungle list tools

# Check if tool is disabled
mcpjungle get tool <tool-name>

# Re-enable if needed
mcpjungle enable tool <tool-name>
```

---

## Command Reference

### Server Management
```bash
mcpjungle start [--enterprise]     # Start server
mcpjungle init-server               # Initialize enterprise server
mcpjungle version                   # Show version
```

### Server Registration
```bash
mcpjungle register -c <config.json>              # Register from file
mcpjungle register --name <name> --url <url>     # Register HTTP server
mcpjungle deregister <server-name>               # Remove server
mcpjungle list servers                           # List all servers
```

### Tool Management
```bash
mcpjungle list tools [--server <name>] [--group <name>]
mcpjungle invoke <tool-name> --input '<json>'
mcpjungle usage <tool-name>
mcpjungle enable tool <name>
mcpjungle disable tool <name>
```

### Tool Groups
```bash
mcpjungle create group -c <config.json>
mcpjungle list groups
mcpjungle get group <name>
mcpjungle delete group <name>
```

### Prompts
```bash
mcpjungle list prompts [--server <name>]
mcpjungle get prompt <name> [--arg key=value]
mcpjungle enable prompt <name>
mcpjungle disable prompt <name>
```

### Access Control (Enterprise)
```bash
mcpjungle create mcp-client <name> --allow "server1,server2"
mcpjungle list mcp-clients
mcpjungle delete mcp-client <name>
```

---

## Related Documentation

- [MCP Protocol Specification](https://modelcontextprotocol.io/introduction)
- [MCPJungle GitHub](https://github.com/mcpjungle/MCPJungle)
- [MCPJungle Docker Hub](https://hub.docker.com/r/mcpjungle/mcpjungle)
- [MCPJungle Discord](https://discord.gg/CapV4Z3krk)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 0)
**Next Steps**: Create jarvis instance configuration and integrate with Cipher
