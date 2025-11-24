# 🤖 Jarvis - MCP Automation Gateway

> **"Just ask, and it shall be configured."**

Jarvis is a specialized **Model Context Protocol (MCP) Server** designed as the autonomous companion to the [MCPM (Model Context Protocol Manager)](https://github.com/pathintegral-institute/mcpm.sh) ecosystem.

While **MCPM** provides the CLI and global registry for managing servers, **Jarvis** acts as the intelligent "Hands" that allow AI agents (like Claude, Cursor, or generic LLMs) to autonomously manage their own tools, infrastructure, and configuration.

By connecting Jarvis to your AI client, you give the AI the power to install new capabilities, fix its own environment, and manage local servers without you needing to open a terminal.

---

## 🌟 Extension of MCPM

Jarvis is built to work seamlessly with MCPM v2.0. It does not replace MCPM but extends it into the agentic realm.

*   **Registry Integration:** Jarvis uses the official [MCP Registry](https://mcpm.sh/registry) (via the MCPM CLI) to discover and verify tools.
*   **Unified Config:** It respects the global configuration model of MCPM, meaning changes made by Jarvis are reflected in your CLI `mcpm ls` and vice-versa.
*   **Vibes:** Built with the same philosophy: *Open Source. Agent Friendly. Simple.*

---

## 🚀 Quick Start

### 1. Build Jarvis
Jarvis is a single binary. You only need Go installed to build it.

```bash
cd Jarvis
go build -o jarvis .
```

### 2. Connect Your Agent
Add Jarvis to your MCP client configuration.

**For Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):**
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/absolute/path/to/MCP/Jarvis/jarvis",
      "args": []
    }
  }
}
```

### 3. Bootstrap the System
Once connected, simply open your AI client and type:

> **"Please bootstrap the MCP system."**

Jarvis will use its `bootstrap_system` tool to:
1.  Install the **MCPM** (Package Manager) core dependencies.
2.  Link the `mcpm` CLI to your system path.
3.  Start the **Infrastructure** (PostgreSQL & Qdrant) via Docker.

---

## 🛠️ Capabilities

Jarvis exposes the following tools to your AI agent:

### 📦 Package Management
*   **`install_server(name)`**: Installs a new MCP server from the registry (e.g., "brave", "sqlite").
*   **`uninstall_server(name)`**: Removes a server and its configuration.
*   **`list_servers()`**: Shows all currently installed and active servers.
*   **`search_servers(query)`**: Finds new tools in the `technologies.toml` registry.
*   **`server_info(name)`**: Displays detailed metadata about a specific tool.

### ⚙️ System Configuration
*   **`bootstrap_system()`**: The one-click setup tool. Initializes the entire environment.
*   **`check_status()`**: Runs a system doctor (checks Node, Docker, Config health).
*   **`manage_config(action, key, value)`**: Modifies core MCPM settings.
*   **`manage_profile(action, name)`**: Creates or switches between tool profiles (e.g., "dev", "research").

### 🌐 Sharing & Tunnels
*   **`share_server(name)`**: Exposes a local MCP server via a secure public tunnel (great for demos or remote agents).
*   **`stop_sharing_server(name)`**: Tears down the tunnel.
*   **`list_shared_servers()`**: Shows active tunnels.

---

## 🏗️ Architecture

Jarvis acts as a **Translation Layer**:

```mermaid
graph LR
    Agent[AI Agent] <-->|MCP Protocol| Jarvis[Jarvis Server]
    Jarvis <-->|Exec| CLI[MCPM Node.js CLI]
    CLI <-->|Manage| Registry[technologies.toml]
    CLI <-->|Control| Docker[Docker Infrastructure]
```

*   **Jarvis (Go):** Handles the high-performance MCP connection, threading, and process lifecycle.
*   **MCPM (Node.js):** Handles the logic of package resolution, dependency management (`npm`), and file generation.

## 🐛 Troubleshooting

**"Command not found"**
If Jarvis reports errors running `mcpm`, ensure you have run the `bootstrap_system` tool at least once, or manually run `npm link` inside the `MCPM/` directory.

**Docker issues**
Jarvis requires Docker to be running for the database components. If `bootstrap_system` fails, check if Docker Desktop/Engine is active.

## 📜 License
MIT License - Compatible with the [MCPM ecosystem](https://github.com/pathintegral-institute/mcpm.sh).
