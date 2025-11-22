# MCPM Server Registration Guide

**Version:** 2.0
**Date:** 2025-11-22
**Status:** Active

---

## Overview

This guide explains how to register and manage MCP servers using **MCPM (Model Context Protocol Manager)**. MCPM automates the installation, configuration, and registration process, eliminating the need for manual JSON file editing.

## 1. Installing Servers

The primary way to register a server is to install it via the `mcpm` CLI.

### Basic Installation

To install a server from the official registry or npm:

```bash
mcpm install @modelcontextprotocol/server-github
```

This command:
1.  Downloads the package.
2.  Registers it in the local MCPM registry.
3.  Makes it available for IDE configuration generation.

### Installing with Custom Configuration

Some servers require specific arguments or environment variables. You can provide these during installation or configure them afterwards.

```bash
# Install with environment variables
mcpm install @brave/brave-search-mcp-server --env BRAVE_API_KEY=your_key
```

---

## 2. Managing Configuration

Once a server is installed, you can modify its configuration using the `mcpm` CLI.

### Setting Environment Variables

If a server requires API keys or other secrets:

```bash
mcpm env set brave-search BRAVE_API_KEY=your_actual_api_key
```

### Updating Arguments

To change the arguments passed to the server executable:

```bash
mcpm config set filesystem --args "/path/to/allowed/directory"
```

### Viewing Configuration

To see the current configuration for a specific server:

```bash
mcpm inspect filesystem
```

---

## 3. Server Management

### Listing Installed Servers

To see all registered servers:

```bash
mcpm list
```

### Updating Servers

To update a specific server to the latest version:

```bash
mcpm update @modelcontextprotocol/server-github
```

To update all installed servers:

```bash
mcpm update all
```

### Uninstalling Servers

To remove a server and its configuration:

```bash
mcpm uninstall @modelcontextprotocol/server-github
```

---

## 4. Advanced Registration

### Registering Local/Custom Servers

If you are developing your own MCP server or have a local script, you can register it manually:

```bash
mcpm register my-local-server \
  --command "python3" \
  --args "/path/to/my/server.py" \
  --env MY_VAR=value
```

### Registering HTTP Servers

For remote MCP servers accessible via HTTP:

```bash
mcpm register remote-server \
  --url "https://api.example.com/mcp" \
  --transport "sse"
```

---

## 5. Troubleshooting

### Installation Failures

*   **Network Issues:** Ensure you have internet access and can reach the npm registry.
*   **Permissions:** If you encounter permission errors, try running with appropriate user privileges (avoid `sudo` if possible, or fix npm permissions).

### Server Not Starting

*   **Check Logs:** Use `mcpm logs <server-name>` (if available) or check your IDE's output.
*   **Verify Command:** Run `mcpm inspect <server-name>` to ensure the command and arguments are correct.
*   **Test Manually:** Try running the command string manually in your terminal to see if it errors out.

### Environment Variables Not Picked Up

*   Ensure you used `mcpm env set` and not just exported them in your shell (unless the server is configured to inherit shell envs).
*   Regenerate your IDE configuration after changing environment variables.
