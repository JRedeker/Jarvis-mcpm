# MCP Ecosystem: Technical Reference

**System Context:** Model Context Protocol (MCP) Management & Execution Environment
**Primary Function:** Orchestration of MCP servers via CLI management and agentic gateways.

## 1. System Architecture

The repository implements a layered architecture separating infrastructure, package management, and execution gateways.

```mermaid
graph TD
    subgraph "Layer 1: Infrastructure (Docker)"
        Daemon[MCPM Daemon :6276-6280] -->|Hosts| Profiles[Profile Runners]
        Profiles -->|Store| DB[(PostgreSQL :5432)]
        Profiles -->|Index| Vector[(Qdrant :6333)]
    end

    subgraph "Layer 2: Package Management (MCPM)"
        MCPM_CLI[MCPM CLI] -->|Reads| Registry[Technology Registry]
        MCPM_CLI -->|Builds| Daemon_Image[Daemon Image]
    end

    subgraph "Layer 3: Execution Gateway (Jarvis)"
        Agent[AI Agent] -->|Stdio| Jarvis[Jarvis Server]
        Agent -.->|HTTP| Daemon
        Jarvis -->|Manage| Daemon
        Jarvis -->|Install| MCPM_CLI
    end

    MCPM_CLI -.->|Configures| Daemon
```

## 2. Component Specifications

### 2.1. Jarvis Gateway (`./Jarvis/`)
*   **Type:** MCP Server (Go)
*   **Role:** Agentic Interface & Presentation Layer.
*   **Mechanism:** Wraps `mcpm` CLI commands into executable MCP tools.
*   **Enhancement:** Acts as a "Presentation Layer" by capturing raw CLI output, stripping ANSI color codes, and wrapping the result in clean Markdown with status emojis (✅/❌) for optimal LLM consumption.
*   **Key Tools:**
    *   `apply_devops_stack(project_type)`: Scaffolds new projects or upgrades existing ones with safe-mode logic.
    *   `analyze_project()`: Returns JSON structure of the current directory (languages, configs).
    *   `restart_infrastructure()`: Reboots the Docker stack via management script.
    *   `install_server(name)`: Invokes `mcpm install`.
    *   `check_status()`: Diagnostics via `mcpm doctor`.
*   **Dependency:** Requires `mcpm` binary in system PATH.

### 2.2. MCPM Core (`./MCPM/`)
*   **Type:** Node.js CLI Application
*   **Role:** Package Manager & Configuration Generator.
*   **Data Source:** `config/technologies.toml` (Registry of supported servers).
*   **Outputs:**
    *   `cline_mcp_settings.json`: Configuration for Cline IDE.
    *   `cursor_mcp_settings.json`: Configuration for Cursor IDE.

### 2.3. Infrastructure (`./`)
*   **Management:** `scripts/manage-mcp.sh` (Unified start/stop/logs/test controller).
*   **Container Runtime:** Docker Compose (`docker-compose.yml`).
*   **Services:**
    *   **MCPM Daemon:** `mcp-mcpm-daemon` (Hosts Streamable HTTP Profiles on ports 6276+).
    *   **PostgreSQL:** `postgres:15` on port `5432`.
    *   **Qdrant:** `qdrant/qdrant:latest` on port `6333` (Vector Store).

## 3. Operational Workflows

### 3.1. Server Installation Path
1.  **Trigger:** Agent calls `jarvis.install_server("brave")`.
2.  **Execution:** Jarvis spawns `mcpm install brave`.
3.  **Resolution:** MCPM resolves "brave" from `technologies.toml`.
4.  **Action:** MCPM installs npm package `@modelcontextprotocol/server-brave-search`.
5.  **Config:** MCPM updates local registry and regenerates IDE config files.

### 3.2. Semantic Search Path (Future)
1.  **Trigger:** Agent calls `jarvis.search_servers("web search")`.
2.  **Query:** Jarvis queries Qdrant vector store (via `memory` server or direct client).
3.  **Result:** Returns relevant server packages based on description embeddings.

## 4. Development Environment Setup

### 4.1. Prerequisites
*   **Runtime:** Go 1.24+, Node.js 18+, Docker Engine.
*   **Environment:** Linux/macOS (Windows via WSL2).

### 4.2. Initialization Sequence

The system is designed to be bootstrapped by the agent itself.

1.  **Build Jarvis:**
    ```bash
    cd Jarvis
    go build -o jarvis .
    ```

2.  **Configure Agent:**
    Use the **3-Layer Profile Stack**. Do not just add the binary.
    *   **Reference:** See `AGENTS.md` or `docs/CONFIGURATION_STRATEGY.md`.
    *   **Pattern:** Wire `jarvis` (direct) AND `mcpm_profile_<name>` (via mcpm) into your client config.

3.  **Bootstrap via Agent:**
    Start your agent and give the instruction:
    > "Please bootstrap the system."

    Jarvis will automatically:
    - Install the MCPM CLI dependencies.
    - Link the `mcpm` command to your system.
    - Start the Docker infrastructure (Postgres & Qdrant).

4.  **Verify:**
    Ask the agent: *"Check system status"* to confirm everything is running.

## 5. Documentation Index

*   **Architecture:** [`ARCHITECTURE-MAP.md`](./ARCHITECTURE-MAP.md)
*   **MCPM Spec:** [`docs/MCPM-documentation.md`](./docs/MCPM-documentation.md)
*   **Server Registry:** [`MCPM/config/technologies.toml`](./MCPM/config/technologies.toml)

## 6. Reference Material

*   **`mcpm_source/`**: This directory contains source code for reference purposes only. It is not part of the active system or build pipeline.
