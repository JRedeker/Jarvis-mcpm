# Cipher-Aggregator Only Routing Guidance

Use ONLY the cipher-aggregator for all MCP tool routing and interactions.

- Aggregator endpoint: http://localhost:3020/sse
- Source of truth for routing policy:
  - [cipher.yml](cipher.yml)
  - [.kilocode/rules/cipher-routing-rules.md](.kilocode/rules/cipher-routing-rules.md)

## Domain-Specific First (authoritative mapping)
Always route to the domain-specific MCP server, never to generic fallbacks:
- GitHub operations → `github` MCP
- Web scraping → `firecrawl` MCP
- Code analysis → `code-index` MCP
- API testing → `schemathesis` MCP
- Test execution → `pytest` MCP
- File operations → `filesystem` MCP
- Web search → `brave-search` MCP
- HTTP requests → `httpie` MCP

See detailed policy and execution constraints in [cipher.yml](cipher.yml) and [.kilocode/rules/cipher-routing-rules.md](.kilocode/rules/cipher-routing-rules.md).

## Prohibited
- Do not reference or use any vendor-specific non-aggregator tools, packages, servers, or docs.
- Do not create, load, or reference any vendor-specific rules files (for example: vendor-rules.md or vendor-rules.mdc).
- Do not route via fetch/curl for domains covered by MCP servers listed above.

## Performance & Safety Constraints
- Serial execution only (maxParallelCalls: 1)
- Minimize tool calls per task (max 8)
- Prefer batch operations when available
- Use workspace memory for caching significant decisions and patterns

## Memory Integration (required)
- Store important routing decisions and successful patterns via the memory-bank MCP as documented in [cipher.yml](cipher.yml).
- Search memory-bank before choosing tools when applicable.

## Session Management
- Generate unique session IDs per task: “session-<timestamp>-<nonce>”
- Reuse a session ID for related operations and clean up when complete

Adhere strictly to [cipher.yml](cipher.yml) and [.kilocode/rules/cipher-routing-rules.md](.kilocode/rules/cipher-routing-rules.md). Any vendor-specific references that conflict with this policy are disallowed.
