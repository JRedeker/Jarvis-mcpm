# MCPM Architecture Map

This document provides a comprehensive map of the MCPM (Model Context Protocol Manager) system architecture, file structure, and component relationships within this repository.

## 1. System Overview

**MCPJungle (jarvis)** is a centralized gateway and registry for Model Context Protocol (MCP) servers. It simplifies the connection between AI clients (IDEs like Cline, Cursor) and various MCP tools by providing a single aggregation point.

### Core Components

*   **Jarvis (MCP Gateway):** The central server (running on port 8080) that aggregates tools from multiple registered MCP servers.
*   **MCP Servers:** Independent services providing specific capabilities (e.g., file system access, web search, GitHub integration).
*   **PostgreSQL:** The backend database for storing server registrations, tool configurations, and metadata.
*   **Docker Compose:** Orchestrates the deployment of the Jarvis gateway, database, and other containerized services.
*   **MCPM (CLI):** A command-line interface (currently transitioning/integrated) for managing packages and server registrations.

## 2. Architecture Diagram

```mermaid
graph TD
    subgraph Clients
        IDE1[Cline IDE]
        IDE2[Kilo Code IDE]
    end

    subgraph "MCPJungle (Jarvis) Host"
        Gateway[Jarvis Gateway :8080]
        DB[(PostgreSQL :5432)]
        
        subgraph "Managed MCP Servers"
            S1[context7 (HTTP)]
            S2[brave-search (stdio)]
            S3[filesystem (stdio)]
            S4[firecrawl (stdio)]
            S5[morph-fast-apply (stdio)]
            S6[gpt-researcher (stdio)]
            S7[fetch (stdio)]
            S8[github (stdio)]
            S9[memory (stdio)]
            S10[playwright (stdio)]
            S11[sqlite (stdio)]
        end
    end

    IDE1 -->|MCP Protocol| Gateway
    IDE2 -->|MCP Protocol| Gateway
    
    Gateway -->|Query/Store| DB
    Gateway -->|HTTP| S1
    Gateway -->|stdio| S2
    Gateway -->|stdio| S3
    Gateway -->|stdio| S4
    Gateway -->|stdio| S5
    Gateway -->|stdio| S6
    Gateway -->|stdio| S7
    Gateway -->|stdio| S8
    Gateway -->|stdio| S9
    Gateway -->|stdio| S10
    Gateway -->|stdio| S11
```

## 3. File Structure & Key Files

### Root Directory
*   `README.md`: Project overview, quick start guide, and available tools list.
*   `docker-compose.yml`: Defines the container services (PostgreSQL, Qdrant).
*   `package.json`: Node.js dependencies, including `mcpm` and server packages.
*   `config/`: Configuration directory.
    *   `technologies.toml`: Central registry of all supported MCP technologies and their metadata. **(Critical)**

### Documentation (`docs/`)
*   `architecture.md`: High-level architecture documentation.
*   `INFORMATION-ARCHITECTURE.md`: Structure of information within the system.
*   `architecture/`: Detailed architectural documents.
    *   `simplified-architecture.md`: Explains the "Jarvis" single-gateway model.
    *   `optimization-roadmap.md`: Plans for system optimization and the role of `mcpm`.
*   `config/`: Configuration documentation.
    *   `port-allocation.md`: Matrix of ports used by services.
    *   `tool-groups.md`: Definitions for tool grouping (Universal, Backend, Frontend).
    *   `actual-configurations.md`: Reference for deployed configurations.
*   `guides/`: User and setup guides.
    *   `server-registration.md`: How to register new servers.
    *   `ide-configuration.md`: Connecting IDEs to Jarvis.
*   `runbooks/`: Operational procedures and status reports.
    *   `phase1b-server-registration-final.md`: Status of server registrations.
*   `tech/`: Technical documentation for individual MCP servers (e.g., `brave-search-mcp.md`, `filesystem-mcp.md`).

## 4. Component Relationships

### Jarvis Gateway <-> MCP Servers
*   **Relationship:** Aggregation / Proxy
*   **Mechanism:** Jarvis connects to servers via `stdio` (spawning processes) or `HTTP` (remote streams).
*   **Configuration:** Defined in JSON files (referenced in docs) and stored in PostgreSQL.

### Jarvis Gateway <-> PostgreSQL
*   **Relationship:** Persistence
*   **Mechanism:** TCP connection to port 5432.
*   **Data:** Stores server registry, tool definitions, and potentially configuration state.

### MCPM (CLI) <-> Jarvis
*   **Relationship:** Management
*   **Mechanism:** The CLI interacts with the Jarvis API or configuration files to register/deregister servers.
*   **Status:** `mcpm` package is a dependency; usage is evolving towards direct integration or helper scripts.

### Tool Groups
*   **Concept:** Logical grouping of tools for specific workflows (Universal, Backend, Frontend).
*   **Implementation:** Configured within Jarvis to expose specific subsets of tools via dedicated endpoints (e.g., `/v0/groups/universal/mcp`).

## 6. Key Data Flows

1.  **Registration:** Admin uses CLI/API -> Jarvis records server details -> PostgreSQL.
2.  **Discovery:** Client (IDE) connects to Jarvis -> Jarvis queries active servers/DB -> Returns list of tools.
3.  **Invocation:** Client requests tool execution -> Jarvis routes request to appropriate MCP server -> Server executes -> Result returned to Client.

## 7. Technology Stack

*   **Runtime:** Node.js (for most MCP servers and tooling).
*   **Database:** PostgreSQL (metadata), Qdrant (vector search/memory - planned).
*   **Containerization:** Docker & Docker Compose.
*   **Protocol:** Model Context Protocol (MCP).
*   **Languages:** TypeScript/JavaScript (primary), Python (some servers like `gpt-researcher`).