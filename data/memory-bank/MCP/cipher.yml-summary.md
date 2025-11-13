# cipher.yml — Structured Summary

Created: 2025-11-13T03:50:17.022Z

Source files referenced: [`cipher.yml`](cipher.yml:1), routing rules [`.kilocode/rules/cipher-routing-rules.md`](.kilocode/rules/cipher-routing-rules.md:1), manager script [`mcp-manager.sh`](mcp-manager.sh:15).

Summary (extracted from [`cipher.yml`](cipher.yml:1)):

- Global settings
  - toolExecution.callTimeout: 45000 ms (45s)
  - maxParallelCalls: 1
  - maxToolCallsPerTask: 8
  - workspaceMemory.path: `./data/workspace-memory` (autoCapture: true)
  - memoryBank.path: `./data/memory-bank`

- Servers (name → enabled, transport, command placeholder)
  - routing-metadata: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/routing-metadata-mcp.py`
  - llm-inference: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/eval_llm_venv/bin/python /home/jrede/dev/MCP/servers/llm-inference-mcp.py`
  - httpie: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/httpie-mcp.py`
  - schemathesis: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/schemathesis-mcp.py`
  - firecrawl: enabled: true, transport: stdio, command: `npx -y @brave/firecrawl-mcp-server`
  - code-index: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/code-index-mcp.py`
  - memory-bank: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/memory-bank-mcp.py`
  - logs: enabled: true, transport: stdio, command: `/home/jrede/dev/MCP/servers/logs-mcp.py`
  - magic-mcp: enabled: false (disabled in template / ticketed)

- Env & integration notes
  - Default env keys referenced: [`OPENAI_API_KEY`](.env:1), `KNOWLEDGE_GRAPH_ENABLED`, `LOKI_BASE_URL`, OTEL config.
  - Secrets and API keys should be held in [`.env`](.env:1) (do not commit secrets directly into `cipher.yml`).

- Operational hooks & scripts
  - Management script references: [`mcp-manager.sh`](mcp-manager.sh:15) uses `-a /home/jrede/dev/MCP/cipher.yml` to start the aggregator.
  - Health checks: `curl http://localhost:3020/health` and optional health server on port 3021 per [`mcp-manager.sh`](mcp-manager.sh:28).

- Routing / policy (enforced by aggregator)
  - Domain-specific-first routing per rules in [`.kilocode/rules/cipher-routing-rules.md`](.kilocode/rules/cipher-routing-rules.md:1).
  - Memory-first pattern: store decisions in memory bank and search `routing-patterns` before making new routing choices.
  - evalLlm (gpt-5-mini) provides routing recommendations; store successful patterns at `data/memory-bank/routing-patterns/`.

Next steps I can perform (select):
- Parse the real `cipher.yml` if you provide/restore it; I'll update memory with exact server commands, timeouts, and env flags.
- Save this structured summary (already done) and additionally push a more detailed per-server checklist into memory (start/verify steps for each server).

Filename: [`data/memory-bank/MCP/cipher.yml-summary.md`](data/memory-bank/MCP/cipher.yml-summary.md:1)
