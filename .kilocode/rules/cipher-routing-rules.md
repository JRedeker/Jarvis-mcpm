[cipher-routing]

## Cipher-Aggregator Routing Rules

When using MCP tools through cipher-aggregator (http://localhost:3020/sse), you MUST follow these routing rules from cipher.yml systemPrompt:

### 1. DOMAIN-SPECIFIC FIRST (CRITICAL)
Always prioritize domain-specific tools over generic ones:
- GitHub operations → `github` MCP (NEVER use fetch/curl)
- Web scraping → `firecrawl` MCP (NEVER use fetch)
- Code analysis → `code-index` MCP (NEVER manually scan files)
- API testing → `schemathesis` MCP
- Test execution → `pytest` MCP
- File operations → `filesystem` MCP (not manual bash)
- Web search → `brave-search` MCP (not fetch/curl)
- HTTP requests → `httpie` MCP (for testing APIs)

### 2. TASK CATEGORIZATION & ROUTING
**Development Tasks**: code-index → filesystem → github → memory-bank
**Web Research**: brave-search → firecrawl → memory-bank
**API Testing**: schemathesis → httpie → pytest
**File Management**: filesystem → file-batch → memory-bank
**Documentation**: context7 → filesystem → memory-bank

### 3. PERFORMANCE CONSTRAINTS
- Serial execution only (maxParallelCalls: 1)
- Minimize tool calls per task (max 8)
- Use batch operations when available
- Cache results in memory-bank for reuse

### 4. MEMORY INTEGRATION
You MUST use cipher's memory systems as documented:
- **Store important decisions** via memory-bank
- **Search memory-bank FIRST** for known solutions before routing
- **Search routing patterns** before making decisions: `memory_bank_search("routing patterns")`
- Workspace memory autoCapture stores session summaries automatically

### 5. evalLlm AWARENESS
Cipher uses evalLlm (gpt-5-mini) for intelligent routing:
- Tool routing is already optimized by cipher's evalLlm
- Follow the routing decisions made by cipher
- Store successful routing patterns in memory-bank
- Routing patterns location: `/home/jrede/dev/MCP/data/memory-bank/routing-patterns/`

### 6. SESSION MANAGEMENT
- Generate unique session IDs: 'session-' + timestamp + random
- Use persistent sessions for related operations
- Clean up sessions when task completes

## Reference Documentation
- Full tool guide: `/home/jrede/dev/MCP/AGENTS.md`
- Routing guide: `/home/jrede/dev/MCP/tool-routing-guide.md`
- evalLlm docs: `/home/jrede/dev/MCP/docs/eval-llm-configuration.md`
- Cipher config: `/home/jrede/dev/MCP/cipher.yml`