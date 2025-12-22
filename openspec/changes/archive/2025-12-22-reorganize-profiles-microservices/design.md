# Design: Composable Micro-Profiles

## Architectural Shift
The architecture moves from a "Layered" approach where Layer 1 was a single `toolbox` to a "Composable" approach where Layer 1 is a *stack* of independent profiles.

### The New Stack
The Client (e.g., OpenCode) will configure its `mcpServers` by composing these profiles together.

```json
// Concept
"mcpServers": {
  ...loadProfile("essentials"),
  ...loadProfile("memory"),
  ...loadProfile("dev-core"),
  ...loadProfile("research"),
  ...loadProfile("data")
}
```

## Profile Definitions

### 1. `essentials` (Port: 6276 - Reused)
*   **Characteristics:** Fast startup (<1s), no external dependencies, low resource usage.
*   **Contents:**
    *   `time`: Local time.
    *   `fetch-mcp`: Simple HTTP fetch (if lightweight).
    *   `filesystem`: (Native/Built-in, but conceptually here).

### 2. `memory` (Port: 6277 - Existing)
*   **Characteristics:** Persistence, state management.
*   **Contents:**
    *   `basic-memory`: Key-value storage.
    *   `mem0-mcp`: Long-term memory.

### 3. `dev-core` (Port: 6278 - Reused/New)
*   **Characteristics:** Code analysis, AST manipulation, specific logic.
*   **Contents:**
    *   `context7`: Code context analysis.
    *   `morph-fast-apply`: Intelligent code editing.

### 4. `research` (Port: 6281 - NEW)
*   **Characteristics:** High latency, Docker containers, API keys, Network I/O. **High Failure Probability.**
*   **Contents:**
    *   `brave-search`: Web search (Docker).
    *   `firecrawl`: Web scraping (Docker/API).
    *   `arxiv-mcp`: Paper search.

### 5. `data` (Port: 6279 - Existing)
*   **Characteristics:** Heavy databases.
*   **Contents:**
    *   `mcp-server-qdrant`: Vector DB.
    *   `postgres`: Relational DB.

## Port Allocation Strategy
We need to assign a unique port to the new `research` profile to allow parallel execution.

| Profile | Port | Status |
|:---|:---|:---|
| `essentials` | 6276 | Was `toolbox`. Keeps default port. |
| `memory` | 6277 | Unchanged. |
| `dev-core` | 6278 | Was `morph` (standalone). Consolidated. |
| `data` | 6279 | Was `qdrant`. Renamed/Expanded. |
| `p-new` | 6280 | Experimental/Testing. |
| `research` | 6281 | **NEW**. |

## Migration Path
1.  **Stop** existing `toolbox`.
2.  **Create** new profiles.
3.  **Update** Client Config (OpenCode) to reference the new list.
4.  **Restart** Infrastructure.

## Documentation Updates
- `AGENTS.md`: Update instructions to use the new profile names.
- `CONFIGURATION_STRATEGY.md`: Document the Micro-Profile pattern.
