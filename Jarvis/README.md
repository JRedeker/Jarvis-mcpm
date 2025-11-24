# Jarvis MCP Server

**Version:** 1.0.0
**Language:** Go

## Overview

Jarvis is a specialized MCP server that acts as a management gateway for the MCPM (Model Context Protocol Manager) ecosystem. It exposes `mcpm` CLI functionality as MCP tools, allowing AI agents to programmatically install, configure, and manage other MCP servers.

## Role & Responsibilities

*   **Gateway:** Serves as the primary interface for AI agents to interact with the MCPM system.
*   **Management:** Provides tools to list, install, uninstall, and configure MCP servers.
*   **Execution:** Wraps the `mcpm` CLI commands and exposes them via the Model Context Protocol.

## Interaction with MCPM

Jarvis does not implement package management logic itself. Instead, it relies on the `mcpm` CLI tool being installed and available in the system PATH.

*   **Command Execution:** Jarvis executes `mcpm` commands (e.g., `mcpm install`, `mcpm list`) as subprocesses.
*   **Environment:** It assumes a configured environment where `mcpm` is accessible.

## Available Tools & Usage

### Server Operations

*   **`list_servers`**: List all installed MCP servers.
    ```json
    { "name": "list_servers" }
    ```

*   **`search_servers`**: Search for available servers in the registry.
    ```json
    {
      "name": "search_servers",
      "arguments": { "query": "memory" }
    }
    ```

*   **`install_server`**: Install a new MCP server.
    ```json
    {
      "name": "install_server",
      "arguments": { "name": "sqlite" }
    }
    ```

*   **`uninstall_server`**: Remove an installed server.
    ```json
    {
      "name": "uninstall_server",
      "arguments": { "name": "sqlite" }
    }
    ```

*   **`server_info`**: Get detailed information about a server.
    ```json
    {
      "name": "server_info",
      "arguments": { "name": "sqlite" }
    }
    ```

*   **`create_server`**: Create a new server configuration.
    ```json
    {
      "name": "create_server",
      "arguments": {
        "name": "custom-tool",
        "type": "stdio",
        "command": "python",
        "args": "server.py"
      }
    }
    ```

*   **`edit_server`**: Edit a server configuration.
    ```json
    {
      "name": "edit_server",
      "arguments": { "name": "custom-tool" }
    }
    ```

### System Management

*   **`check_status`**: Check the health of the MCPM system.
    ```json
    { "name": "check_status" }
    ```

*   **`usage_stats`**: Display comprehensive analytics and usage data.
    ```json
    { "name": "usage_stats" }
    ```

### Configuration & Profiles

*   **`manage_config`**: Manage MCPM configuration settings.
    ```json
    {
      "name": "manage_config",
      "arguments": {
        "action": "set",
        "key": "theme",
        "value": "dark"
      }
    }
    ```

*   **`manage_profile`**: Manage MCPM profiles.
    ```json
    {
      "name": "manage_profile",
      "arguments": {
        "action": "create",
        "name": "development"
      }
    }
    ```

*   **`manage_client`**: Manage MCP client configurations.
    ```json
    {
      "name": "manage_client",
      "arguments": { "action": "ls" }
    }
    ```

*   **`migrate_config`**: Migrate v1 configuration to v2.
    ```json
    { "name": "migrate_config" }
    ```

## Development

### Prerequisites

*   Go 1.24+
*   `mcpm` installed globally (or in PATH)

### Building

```bash
go build -o jarvis .
```

### Running

```bash
./jarvis
```