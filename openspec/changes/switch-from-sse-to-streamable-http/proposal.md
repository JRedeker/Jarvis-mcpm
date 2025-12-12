# Switch from SSE to Streamable HTTP

## Background
The Model Context Protocol (MCP) has introduced **Streamable HTTP** as the preferred transport mechanism, deprecating Server-Sent Events (SSE). Our current architecture relies heavily on a "Single SSE Daemon" (`mcpm-daemon`) that exposes tools via SSE endpoints (e.g., `http://localhost:6276/sse`). To stay aligned with the spec and ensure future compatibility, we must migrate our transport layer to Streamable HTTP.

## Problem
*   **Deprecation:** SSE is deprecated in the latest MCP spec (2025-03-26).
*   **Compatibility:** Future clients and tools may drop support for SSE.
*   **Performance/Features:** Streamable HTTP offers better handling of connection lifecycles and is more standard for modern APIs.

## Goals
1.  Migrate `mcpm-daemon` and all profiles to use Streamable HTTP.
2.  **Establish Streamable HTTP as the default transport** for all new projects, server registrations, and templates.
3.  Update `Jarvis` to manage and interact with Streamable HTTP endpoints.
4.  Update configuration templates and scripts to reflect the change (e.g., `/sse` -> `/mcp`).
5.  Ensure backward compatibility where possible or provide a clean migration path.

## Non-Goals
*   Changing the underlying `mcpm` logic unrelated to transport.
*   Refactoring the entire `Jarvis` codebase beyond transport handling.

## Risks
*   **Client Compatibility:** Older clients might not support Streamable HTTP immediately (though most MCP clients should be updating).
*   **Downtime:** The migration requires restarting the daemon and potentially updating client configs (Claude Desktop, etc.).

## User Stories
*   **As an AI Agent Developer,** I want my tools to communicate using the latest MCP standard (Streamable HTTP) so that my infrastructure remains compatible with future client updates.
*   **As a Platform User,** I want my existing configuration to be automatically migrated from SSE to Streamable HTTP so that I don't have to manually edit JSON files or debug connection issues.
*   **As a System Administrator,** I want `check_status` to correctly report the health of HTTP endpoints so that I can reliably monitor my tool server's uptime.
*   **As a new User,** I want new server configurations to default to Streamable HTTP so that I am using best practices from day one without extra configuration.
