# IDE Configuration Guide for MCPM

**Version:** 2.0
**Date:** 2025-11-22
**Status:** Active

---

## Overview

This guide explains how to configure your IDEs (Cline, Kilo Code, Claude Desktop) to work with **MCPM (Model Context Protocol Manager)**. MCPM simplifies the management of MCP servers, allowing you to easily install, configure, and connect tools to your development environment.

## Supported IDEs

- **Cline** (VS Code extension)
- **Kilo Code** (Cursor-based IDE)
- **Claude Desktop** (standalone application)

---

## 1. Prerequisites

Ensure you have `mcpm` installed and available in your system path.

```bash
npm install -g mcpm
```

Verify the installation:

```bash
mcpm --version
```

## 2. General Configuration Strategy

Instead of manually editing JSON configuration files with complex server details, you can use `mcpm` to generate or manage these configurations.

### The `mcpm` Approach

1.  **Install Servers:** Use `mcpm install <package>` to add servers to your local registry.
2.  **Generate Config:** Use `mcpm config` to output the JSON configuration needed for your IDE.
3.  **Connect:** Paste the configuration into your IDE's settings.

---

## 3. Cline (VS Code) Configuration

### Automatic Configuration (Recommended)

If you are using the Cline extension in VS Code, you can often point it to your MCPM configuration or copy the generated config directly.

1.  **Generate Config:**
    ```bash
    mcpm config --format cline
    ```
2.  **Copy Output:** Copy the JSON output from the command.
3.  **Update Settings:**
    *   Open VS Code Settings (`Ctrl/Cmd + ,`).
    *   Search for "cline mcp".
    *   Edit `settings.json` and paste the configuration under `cline.mcp`.

### Example `settings.json`

```json
{
  "cline.mcp": {
    "mcpServers": {
      "brave-search": {
        "command": "npx",
        "args": ["-y", "@brave/brave-search-mcp-server"],
        "env": {
          "BRAVE_API_KEY": "your-key-here"
        }
      },
      "filesystem": {
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
      }
    }
  }
}
```

*Note: The above is just an example. Run `mcpm config` to get the exact configuration for your installed servers.*

---

## 4. Kilo Code (Cursor) Configuration

Kilo Code (and Cursor) uses a similar configuration structure to VS Code.

1.  **Generate Config:**
    ```bash
    mcpm config --format cursor
    ```
2.  **Locate Config File:**
    *   Typically found at `~/.cursor/mcp.json` or accessed via the Cursor settings UI.
3.  **Update Config:**
    *   Paste the generated JSON into the configuration file.

---

## 5. Claude Desktop Configuration

Claude Desktop requires a specific configuration file location.

1.  **Generate Config:**
    ```bash
    mcpm config --format claude
    ```
2.  **Update Config File:**
    *   **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
    *   **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

    Open the file and replace/merge the `mcpServers` section with the output from `mcpm`.

---

## 6. Managing Servers with MCPM

### Installing a New Server

To add a new tool to your IDE, simply install it via `mcpm`:

```bash
mcpm install @modelcontextprotocol/server-github
```

After installation, re-run the config generation command for your IDE and update the settings to include the new server.

### Removing a Server

```bash
mcpm uninstall @modelcontextprotocol/server-github
```

### Updating Servers

```bash
mcpm update all
```

---

## 7. Troubleshooting

### Tools Not Showing Up?

1.  **Check Installation:** Run `mcpm list` to ensure the server is installed and recognized.
2.  **Verify Config:** Double-check that the JSON in your IDE settings matches the output of `mcpm config`.
3.  **Restart IDE:** Sometimes a restart is required for the IDE to pick up new MCP server processes.
4.  **Check Logs:**
    *   **Cline:** Check the "MCP" output channel in VS Code.
    *   **Claude Desktop:** Check the application logs for connection errors.

### Environment Variables

Some servers require API keys (e.g., Brave Search, GitHub). Ensure these are set in your environment or explicitly added to the generated configuration.

`mcpm` allows you to set environment variables for specific servers:

```bash
mcpm env set brave-search BRAVE_API_KEY=your_key_here
```

This ensures that when you generate the config, the keys are correctly populated (or referenced).
