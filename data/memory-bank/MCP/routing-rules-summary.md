# Cipher Routing Rules — Summary (extracted)

Created: 2025-11-13T03:33:43.142Z

Source: [.kilocode/rules/cipher-routing-rules.md](.kilocode/rules/cipher-routing-rules.md:1)

Concise summary of critical routing rules and operational constraints to store in memory:

- DOMAIN-SPECIFIC FIRST (priority)
  - Always prefer domain-specific MCP tools over generic fallbacks.
  - Examples: GitHub → `github`; Web scraping → `firecrawl`; Code analysis → `code-index`; API testing → `schemathesis`; Test execution → `pytest`; File ops → `filesystem`; Web search → `brave-search`; HTTP → `httpie`.

- TASK CATEGORIZATION & ROUTING
  - Development Tasks: code-index → filesystem → github → memory-bank
  - Web Research: brave-search → firecrawl → memory-bank
  - API Testing: schemathesis → httpie → pytest
  - File Management: filesystem → file-batch → memory-bank
  - Documentation: context7 → filesystem → memory-bank

- PERFORMANCE CONSTRAINTS
  - Serial execution only (maxParallelCalls: 1)
  - Minimize tool calls per task (max 8)
  - Use batch operations where available
  - Cache results in memory-bank for reuse

- MEMORY INTEGRATION (required)
  - Store important decisions in the memory-bank.
  - Search memory-bank first for existing solutions before routing.
  - Search routing patterns before making decisions: `memory_bank_search("routing patterns")`.
  - Workspace memory auto-capture session summaries automatically.

- evalLlm (routing intelligence)
  - Cipher uses evalLlm (gpt-5-mini) for routing decisions; follow its routing recommendations.
  - Store successful routing patterns in memory-bank at `/home/jrede/dev/MCP/data/memory-bank/routing-patterns/`.

- SESSION MANAGEMENT
  - Generate unique session IDs: `session-<timestamp>-<random>`.
  - Use persistent sessions for related operations and clean up when complete.

- REFERENCE LOCATIONS (in-repo)
  - Routing rules doc: `.kilocode/rules/cipher-routing-rules.md` [`.kilocode/rules/cipher-routing-rules.md:1`]
  - Cipher config (source-of-truth): `cipher.yml` (not found during scan; recommend adding or exposing when available)
  - Routing patterns memory folder: `data/memory-bank/routing-patterns/`

Notes & Recommended next steps:
- I could parse `cipher.yml` and extract concrete server/timeouts/routing entries and store a structured summary; `cipher.yml` was not present in the repo root when I scanned — confirm if it exists elsewhere or if you want me to proceed when you add it.
- This file captures the enforced operational rules (domain-first, memory-first, serial constraints) so agents/operators can reference them quickly.

Filename: routing-rules-summary.md
