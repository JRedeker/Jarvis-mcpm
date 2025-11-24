# MCPM: Node.js CLI Bundle

**Version:** 1.0.0 (Fork)
**Language:** Node.js

## Overview

This directory contains the **specialized Node.js implementation** of the MCPM (Model Context Protocol Manager) CLI.

It is a streamlined fork designed specifically for the Jarvis ecosystem, replacing the original Python-based CLI to ensure stability and provide agent-centric features.

## 🚀 Features

*   **Native Node.js:** No complex Python virtual environments or Pydantic version conflicts.
*   **Agent-First Output:** Commands like `install` output clean JSON configuration snippets that Agents can read and use immediately.
*   **Docker Preference:** Automatically detects `docker` fields in the registry and prefers containerized installation for security.

## 📂 Directory Structure

*   **`index.js`**: The main CLI entry point (executable). Uses `commander` for argument parsing.
*   **`config/technologies.toml`**: The "Registry". A TOML file defining all supported tools, their repositories, and installation methods.
*   **`FORK_DETAILS.md`**: A deep dive into *why* this fork exists and how it differs from the upstream Python version.

## 🛠️ Development

### Prerequisites
*   Node.js 18+
*   NPM

### Installation (Manual)
If you are working on the CLI code directly:

```bash
npm install
npm link
```

This makes the `mcpm` command available globally on your system, pointing to this directory.

### Adding New Tools
To add a new tool to the ecosystem, simply edit `config/technologies.toml`.

```toml
[technologies.backend]
my_new_tool = {
    repo = "https://github.com/user/my-tool",
    docker = "user/my-tool:latest",  # Optional: For Docker support
    description = "A cool new tool"
}
```

## 🔗 Reference
See [FORK_DETAILS.md](./FORK_DETAILS.md) for a technical comparison with the original architecture.
