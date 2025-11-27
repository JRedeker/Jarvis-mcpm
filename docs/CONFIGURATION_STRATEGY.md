# MCPM & Jarvis Configuration Strategy: The 3-Layer Stack

**Version:** 2.0 (November 2025)
**Status:** Active Standard

## Overview

To ensure scalability, conflict avoidance, and "agentic" intelligence, we have adopted a **3-Layer Profile Stacking Architecture**. This approach replaces monolithic configuration with modular, composable layers that are assembled dynamically by Jarvis at runtime.

## The 3 Layers

### Layer 1: Environment (The Base)
*   **Naming Convention:** `project-<name>` (e.g., `project-pokeedge`, `project-new`)
*   **Purpose:** Defines the *workspace context*. It contains the tools necessary for the specific project's domain (e.g., Databases, specific APIs).
*   **Exclusivity:** An agent can only be in one Environment at a time.
*   **Default:** `project-new` (Standard scaffolding tools: `context7`, `fetch`, `search`).

### Layer 2: Client Adapter (The Glue)
*   **Naming Convention:** `client-<name>` (e.g., `client-codex`, `client-gemini`)
*   **Purpose:** Adds tools that are specific to the AI Client being used (e.g., `morph-fast-apply`, specialized rendering tools).
*   **Constraint:** MUST NOT contain general tools (like `time` or `memory`) to avoid duplication.

### Layer 3: Global Capabilities (The Augment)
*   **Naming Convention:** `memory`, `testing-all-tools`
*   **Purpose:** Cross-cutting concerns that are always active or toggled on demand.
*   **Standard:** `memory` (Persistent storage: `basic-memory`, `mem0`, `qdrant`).

## The Jarvis Logic ("The Brain")

Jarvis (the MCP Server) is configured **Directly** in all clients (via binary path). It exposes a `suggest_profile` tool that implements the stacking logic:

```go
// Pseudo-logic for suggest_profile(client_name, testing_mode)
profiles = []

// 1. Determine Environment
if cwd contains "pokeedge" -> profiles.add("project-pokeedge")
else -> profiles.add("project-new")

// 2. Add Client Adapter
if client_name == "codex" -> profiles.add("client-codex")
if client_name == "gemini" -> profiles.add("client-gemini")

// 3. Add Globals
profiles.add("memory")
if testing_mode -> profiles.add("testing-all-tools")

return profiles // e.g., ["project-new", "client-codex", "memory"]
```

## Configuration Rules

1.  **Jarvis Must Be Direct:** Never put `jarvis` inside an MCPM profile. Configure it directly in the client's config file pointing to the binary. This ensures management capabilities are available even if profiles break.
2.  **No Overlap:** Tools should exist in exactly one layer.
    *   `time`, `fetch` -> Layer 1 (Project)
    *   `morph-fast-apply` -> Layer 2 (Client)
    *   `mem0` -> Layer 3 (Global)
3.  **Dynamic Loading:** Clients should be configured to either:
    *   Ask Jarvis for the profile list on startup (Ideal).
    *   Or hardcode the layers if dynamic loading isn't supported (e.g., `mcpm profile run client-codex` + `mcpm profile run project-pokeedge`).

## Current Profile Map

| Profile Name | Layer | Contents |
| :--- | :--- | :--- |
| `project-new` | 1 (Env) | `context7`, `fetch`, `search`, `time` |
| `project-pokeedge` | 1 (Env) | `project-new` + `schemathesis`, `httpie`, `pytest` |
| `client-codex` | 2 (Client) | `morph-fast-apply` |
| `client-gemini` | 2 (Client) | `morph-fast-apply` |
| `memory` | 3 (Global) | `basic-memory`, `mem0`, `qdrant` |
| `testing-all-tools` | 3 (Global) | **ALL** tools (for CI/CD) |
