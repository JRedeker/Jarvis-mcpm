# Technical Design

## Architecture Changes

### 1. Daemon Endpoints
The `mcpm-daemon` (using `supervisord` or direct execution) currently exposes SSE endpoints. We will switch the listening mode to Streamable HTTP.

*   **Current:** `http://localhost:<port>/sse`
*   **New:** `http://localhost:<port>/mcp` (or as defined by the spec, often just the root or a dedicated path like `/message`).

### 2. Jarvis Updates
*   **Server Registration:** `create_server` and `edit_server` tools need to accept `streamable-http` as a valid type and default to it.
*   **Health Checks:** `check_status` (and `doctor`) must validate Streamable HTTP endpoints instead of looking for SSE headers.
*   **Process Management:** Ensure the `restart_profiles` logic works with the new transport.

### 3. Configuration & Scripts
*   **`update_configs.py`:** Update logic to rewrite `transport: sse` to `transport: streamable-http` and update URLs.
*   **Templates:** Create/Update `claude-code-http.json` (replacing `claude-code-sse.json`).
*   **Environment Variables:** Review `.env` or Docker env vars if they specify transport explicitly.

## Data Model

No major schema changes in the database, but `servers.json` and `profiles.json` (managed by `mcpm`) will have updated fields:
```json
{
  "url": "http://localhost:6276/mcp",
  "transport": "streamable-http" // instead of "sse"
}
```

## Migration Strategy

1.  **Dual Support (Transitional):** If possible, have `mcpm-daemon` support both for a short period. If not, a hard cutover is acceptable given this is a dev tool.
2.  **Automated Config Update:** The `update_configs.py` script will be the primary migration tool for users.
