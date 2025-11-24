# MCP Ecosystem: Developer Context

## Project Overview

This repository hosts the **Model Context Protocol (MCP) Management & Execution Environment**. It is a modular system designed to orchestrate MCP servers, providing both a CLI package manager (`mcpm`) and an agentic gateway (`jarvis`) for AI-driven control.

### Core Components

1.  **Infrastructure (Root):** Shared services (PostgreSQL, Qdrant) defined in `docker-compose.yml`.
2.  **MCPM (`./MCPM/`):** A Node.js-based Package Manager CLI. It manages the registry of MCP servers, handles installation, and generates IDE configurations (Cline, Cursor).
3.  **Jarvis (`./Jarvis/`):** A Go-based MCP Server. It acts as a gateway, exposing `mcpm` capabilities as executable tools to AI agents.

## Architecture

The system follows a layered architecture:

*   **Layer 1: Infrastructure:** Docker containers for state persistence (Postgres) and vector search (Qdrant).
*   **Layer 2: Logic (MCPM):** The "brains" of the package management, driven by `config/technologies.toml`.
*   **Layer 3: Interface (Jarvis):** The "hands" that execute commands on behalf of an agent.

See `ARCHITECTURE-MAP.md` for a detailed visual graph.

## Building and Running

### Prerequisites
*   **Runtime:** Go 1.24+, Node.js 18+
*   **Container:** Docker Engine & Docker Compose
*   **OS:** Linux/macOS (or WSL2 on Windows)

### Initialization Sequence

1.  **Start Infrastructure:**
    ```bash
    docker-compose up -d
    ```

2.  **Setup MCPM (CLI):**
    ```bash
    cd MCPM
    npm install
    npm link  # Makes 'mcpm' available globally
    cd ..
    ```

3.  **Build Jarvis (Gateway):**
    ```bash
    cd Jarvis
    go mod download
    go build -o jarvis .
    cd ..
    ```

### Running the System
*   **Manual CLI:** Use `mcpm <command>` directly in the terminal (e.g., `mcpm list`, `mcpm install brave`).
*   **Agent Gateway:** Run `./Jarvis/jarvis` to start the MCP server. This is typically done by configuring your AI client (like Claude Desktop or an IDE) to run this executable.

## Key Files & Directories

*   **`docker-compose.yml`**: Defines the Postgres (port 5432) and Qdrant (port 6333) services.
*   **`docs/`**: Comprehensive documentation.
    *   `docs/MCPM-documentation.md`: Detailed spec for the package manager.
    *   `docs/tech/`: Guides for specific MCP servers (Brave, Filesystem, etc.).
*   **`MCPM/config/technologies.toml`**: The central registry of supported MCP servers. **Modify this to add new tools.**
*   **`Jarvis/main.go`**: Entry point for the Go MCP server.

## Development Conventions

*   **Registry:** New MCP servers are added to `MCPM/config/technologies.toml` first.
*   **Language Standards:**
    *   **Go (Jarvis):** Follow standard Go idioms.
    *   **Node.js (MCPM):** Uses standard NPM structure.
*   **Architecture:** The system is designed to be "agent-first". Tools should be exposed in a way that is semantically clear to an LLM.
