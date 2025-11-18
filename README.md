<h1 align="center">
  :deciduous_tree: MCPJungle (jarvis) :deciduous_tree:
</h1>
<p align="center">
  Self-hosted MCP Gateway for your private AI agents
</p>
<p align="center">
  <a href="https://discord.gg/CapV4Z3krk" style="text-decoration: none;">
    <img src="https://img.shields.io/badge/Discord-MCPJungle-5865F2?style=flat-square&logo=discord&logoColor=white" alt="Discord" style="max-width: 100%;">
  </a>
</p>

**MCPJungle (jarvis)** is a single source-of-truth registry for all [Model Context Protocol](https://modelcontextprotocol.io/introduction) Servers running in your Organisation. This implementation uses the simplified jarvis architecture - no cipher aggregator complexity.

🧑‍💻 Developers use it to register & manage MCP servers and the tools they provide from a central place.

🤖 MCP Clients use it to discover and consume all these tools from a single "Gateway" MCP Server.

![diagram](./assets/mcpjungle-diagram/mcpjungle-diagram.png)

<p align="center">MCPJungle jarvis is the only MCP Server your AI agents need to connect to!</p>

## 🚀 Quick Start - jarvis Architecture

This quickstart guide will show you how to:
1. Start the jarvis server locally
2. Register MCP servers in jarvis
3. Connect your Claude/Cursor to jarvis to access all MCP tools

### Prerequisites
- All required API keys in `.env` file:
  - `BRAVE_API_KEY`, `FIRECRAWL_API_KEY`, `MORPH_API_KEY`, `TAVILY_API_KEY`, `OPENAI_API_KEY`

### Start the jarvis server
```bash
# Download and start MCPJungle binary
wget https://github.com/mcpjungle/MCPJungle/releases/download/0.2.16/mcpjungle_Linux_x86_64.tar.gz
tar -xzf mcpjungle_Linux_x86_64.tar.gz
chmod +x mcpjungle

# Start jarvis server
./mcpjungle start --port 8080

# Verify it's running
curl http://localhost:8080/health
```

### Register all MCP servers at once
```bash
# Register all 6 pre-configured servers
./mcpjungle register -c config/jarvis/servers/context7.json
./mcpjungle register -c config/jarvis/servers/brave-search.json
./mcpjungle register -c config/jarvis/servers/filesystem.json
./mcpjungle register -c config/jarvis/servers/firecrawl.json
./mcpjungle register -c config/jarvis/servers/morph-fast-apply.json
./mcpjungle register -c config/jarvis/servers/gpt-researcher.json

# Verify registration
./mcpjungle list servers
./mcpjungle list tools
```

### Connect to jarvis

Use the following configuration for your Claude/Cursor MCP servers config:
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:8080/mcp",
        "--allow-http"
      ]
    }
  }
}
```

Once jarvis is added as an MCP to your Claude, try asking it:
```text
Use context7 to get the documentation for `/lodash/lodash`
```

Claude will then call the `context7__get-library-docs` tool via jarvis.

## 📋 Available Tools (34 Total)

### Documentation & Research
- **context7**: Library documentation lookup (2 tools)
- **gpt-researcher**: Deep web research with AI analysis (5 tools)

### Web & Search
- **brave-search**: Web search, news, images, videos (6 tools)
- **firecrawl**: Web scraping and content extraction (6 tools)

### Development
- **filesystem**: File operations, reading, writing (14 tools)
- **morph-fast-apply**: AI-powered code editing (1 tool)

## 🏗️ Architecture - jarvis Implementation

This repository implements the **simplified jarvis architecture**:

```
[IDEs] → [jarvis:8080] → [MCP Servers]
```

**Key Features:**
- ✅ **Single aggregation layer** - No cipher complexity
- ✅ **34 tools available** - All major MCP servers integrated
- ✅ **Health monitoring** - Built-in health checks
- ✅ **Tool groups** - Organize tools by use case
- ✅ **Access control** - Enterprise security features

**Previous cipher aggregator files have been archived** - see `archive/cipher-aggregator/`

## 📊 Phase Status

| Phase | Status | Progress | Target Date |
|-------|--------|----------|-------------|
| **Phase 0.5: Documentation** | ✅ COMPLETED | 100% | 2025-11-18 |
| **Phase 1: Core Setup** | ✅ COMPLETED | 100% | 2025-11-18 |
| **Phase 2: Memory Research** | ⏸️ Ready | 0% | 2025-11-25 |
| **Phase 3: Memory Implementation** | ⏸️ Not Started | 0% | 2025-12-02 |
| **Phase 4: IDE Migration** | ⏸️ Not Started | 0% | 2025-12-06 |

## 🔧 Installation & Setup

### Option 1: Direct Binary (Recommended)
```bash
# Download latest release
wget https://github.com/mcpjungle/MCPJungle/releases/download/0.2.16/mcpjungle_Linux_x86_64.tar.gz
tar -xzf mcpjungle_Linux_x86_64.tar.gz
chmod +x mcpjungle

# Start jarvis
./mcpjungle start --port 8080
```

### Option 2: Docker Compose
```bash
# Use provided docker-compose.yml
docker compose up -d
```

### Option 3: Homebrew (macOS)
```bash
brew install mcpjungle/mcpjungle/mcpjungle
mcpjungle start
```

## 🛠️ Usage

### Basic Operations
```bash
# Check health
curl http://localhost:8080/health

# List registered servers
./mcpjungle list servers

# List available tools
./mcpjungle list tools

# Test a tool
./mcpjungle invoke context7__get-library-docs --input '{"context7CompatibleLibraryID": "/lodash/lodash"}'
```

### Server Management
```bash
# Register a new server
./mcpjungle register -c config.json

# Remove a server
./mcpjungle deregister <server-name>

# Enable/disable tools
./mcpjungle enable <server-name>__<tool-name>
./mcpjungle disable <server-name>__<tool-name>
```

### Tool Groups
```bash
# Create tool group
./mcpjungle create group -c group-config.json

# List groups
./mcpjungle list groups

# Use group-specific endpoint
http://localhost:8080/v0/groups/<group-name>/mcp
```

## 📚 Documentation

### Core Documents
- **[MCP-MASTER.md](MCP-MASTER.md)** - Master implementation plan
- **[docs/phase1-completion-report.md](docs/phase1-completion-report.md)** - Phase 1 completion details
- **[docs/config/actual-configurations.md](docs/config/actual-configurations.md)** - Actual server configs used

### Guides
- **[docs/guides/install-mcpjungle.md](docs/guides/install-mcpjungle.md)** - Installation procedures
- **[docs/guides/server-registration.md](docs/guides/server-registration.md)** - Server registration guide
- **[docs/guides/ide-configuration.md](docs/guides/ide-configuration.md)** - IDE setup instructions

### Architecture
- **[docs/architecture/simplified-architecture.md](docs/architecture/simplified-architecture.md)** - Technical architecture
- **[docs/config/port-allocation.md](docs/config/port-allocation.md)** - Port allocation matrix

## 🔮 Next Steps

### Phase 2: Memory Research (Starting Soon)
Research and implement memory solutions:
1. Test memory-bank MCP server
2. Evaluate Cipher default mode
3. Create comparison matrix
4. Implement chosen solution

### Future Enhancements
- Memory persistence across sessions
- Advanced tool grouping
- Analytics and monitoring
- Custom memory solutions

## 🤝 Contributing

We welcome contributions! See our documentation for:
- Development setup
- Architecture decisions
- Contribution guidelines

**Current Focus**: Memory solution research and implementation

---

**Status**: ✅ **Phase 1 Complete** - 6/6 servers registered, 34 tools available
**Next**: Phase 2 memory research starting 2025-11-25
