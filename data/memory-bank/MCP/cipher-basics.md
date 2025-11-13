# Cipher — Critical Basics (concise reference)

Created: 2025-11-13T03:32:36.534Z

Summary (key items to store for agents/operators):

- Canonical name & location
  - "cipher" → `cipher-aggregator` (single aggregator for this workspace).
  - Repo root: /home/jrede/dev/MCP (use this as project root).

- Primary runtime & control
  - SSE port: 3020 (SSE server at http://127.0.0.1:3020/sse).
  - PID file: `cipher-aggregator.pid` in repo root.
  - Logs: `logs/` (cipher-aggregator-*.log); check tail -f logs/cipher-aggregator.log for startup/connectivity.
  - Management script: `mcp-manager.sh` — provides start/stop/restart/status/health commands.
  - Typical start command uses cipher binary with: `--mode mcp --mcp-transport-type sse --mcp-host 127.0.0.1 --mcp-port 3020 -a /home/jrede/dev/MCP/cipher.yml`.

- Configuration
  - Main config: `cipher.yml` at repo root — source of truth for routing, server blocks, env toggles and tool timeouts.
  - .env must contain OPENAI_API_KEY for embedding/semantic tools; .env Unix line endings required by some scripts.
  - Health helper config: `health-config.env` and optional health server at port 3021.

- Agent / routing rules (behavioral)
  - Agents only see the aggregator (never individual servers). Always connect to `cipher-aggregator` and call `tools/list` to discover available tools.
  - Routing and domain rules live in `cipher.yml` and helper doc `.kilocode/rules/cipher-routing-rules.md` (DOMAIN-SPECIFIC-FIRST, etc.).
  - Use the alias: when user/agents say "cipher", treat as `cipher-aggregator`.

- Memory & workspace tools available (important for storing/searching knowledge)
  - Memory Bank unified endpoints (aggregated names): `memory_bank_write`, `memory_bank_read`, `memory_bank_update`, `memory_bank_list_projects`.
  - Built-in cipher semantic tools: `cipher_memory_search`, `cipher_extract_and_operate_memory`, `cipher_workspace_search`, `cipher_workspace_store`, `cipher_store_reasoning_memory`, `cipher_search_reasoning_patterns`, `cipher_add_node`/`cipher_add_edge`/graph tools.
  - Use `cipher_extract_and_operate_memory` for structured extraction+store and `cipher_memory_search` for semantic retrieval.
  - Workspace memory path: `/home/jrede/dev/MCP/data/workspace-memory` and project memory bank: `/home/jrede/dev/MCP/data/memory-bank`.

- Useful operational checks & patterns
  - Validate cipher.yml syntax: `python3 -c "import yaml; yaml.safe_load(open('cipher.yml'))"`.
  - Check aggregator health: `curl http://localhost:3020/health` (or use `mcp-manager.sh health`).
  - If tools missing (embeddings/graph) check OPENAI_API_KEY, KNOWLEDGE_GRAPH_ENABLED env toggles in cipher.yml.
  - For startup issues inspect `mcp-manager.sh` logs and `logs/cipher-aggregator-*.log` (look for "Successfully connected" lines).

- Key commands & reminders for agents/operators
  - Start: `./mcp-manager.sh start`
  - Restart (apply config changes): `./mcp-manager.sh restart` then tail logs.
  - List tools: request aggregator `tools/list` (or use `mcp-manager.sh list_tools` wrapper).
  - Treat aggregator as single-entrypoint: agents must not attempt to call individual MCP servers directly.

Notes / Recommendations stored here:
- Use `cipher_memory_search` and `cipher_extract_and_operate_memory` as first-class tools for knowledge storage and retrieval.
- Always prefer domain-specific tools as per routing rules (see `.kilocode/rules/cipher-routing-rules.md`).
- When in doubt, check `cipher.yml` and aggregator logs; store any discovered routing patterns in memory for reuse.

Filename: cipher-basics.md
