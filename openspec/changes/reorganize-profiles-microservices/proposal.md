# Profile Architecture Reorganization: From Monolith to Micro-Profiles

## Summary
This proposal transitions the default Jarvis/MCPM configuration strategy from a monolithic "toolbox" profile to a set of composable "micro-profiles." This change addresses stability issues caused by mixing lightweight, essential tools with heavy, network-dependent, or Docker-based tools in a single profile.

## Problem
Currently, the default `toolbox` profile aggregates a wide variety of tools:
- **Essential:** `time`, `fetch-mcp` (fast, robust).
- **Heavy/Network:** `brave-search`, `firecrawl` (Docker-based, network-dependent, prone to timeouts).
- **Complex:** `context7`, `morph-fast-apply` (specialized logic).

When a heavy tool (like `brave-search` or `firecrawl`) fails to start or times out (common with Docker/Network issues), the entire `toolbox` profile becomes unresponsive. This denies the agent access to *all* tools, including critical ones like `time` or memory, effectively paralyzing the workflow.

## Solution
We will decompose the monolithic `toolbox` into domain-specific micro-profiles, each with its own failure domain and Supervisor process.

### New Profile Taxonomy
1.  **`essentials`**: Fast, local utilities that must always be available (e.g., `time`, `fetch-mcp`).
2.  **`research`**: High-latency, network-bound tools (e.g., `brave-search`, `firecrawl`). Failures here should not crash the agent.
3.  **`dev-core`**: Coding and logic tools (e.g., `context7`, `morph-fast-apply`).
4.  **`memory`**: State and persistence (e.g., `basic-memory`, `mem0-mcp`).
5.  **`data`**: Heavy database services (e.g., `qdrant`).

## Impact
- **Stability:** A timeout in `research` will no longer block `essentials` or `dev-core`.
- **Performance:** Parallel startup of profiles.
- **Debuggability:** `jarvis_diagnose` can target specific domains (e.g., "Check research tools").
- **Scalability:** Easier to add new categories without destabilizing existing ones.
