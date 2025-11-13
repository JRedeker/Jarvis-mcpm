# VS Code LM Provider (Cipher-first) — Implementation Specification

Status: Proposed
Owners: Platform/Agents
Related: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py), [docs/llm-inference-setup.md](docs/llm-inference-setup.md), [cipher.yml](cipher.yml), [data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md](data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md)

1) Objective

Provide a VS Code “LM API provider” that routes all IDE language model calls to Cipher first, which then invokes the llm-inference MCP tool to reach OpenRouter. This centralizes auth, routing, policy, and observability, while preserving OpenRouter prompt caching and fallbacks.

2) Goals / Non‑Goals

Goals
- Route VS Code LM requests via Cipher → llm-inference → OpenRouter
- Maintain OpenRouter caching and reliability
- Support per-request programmatic model control via tier override
- Provide an “Auto (Cipher)” model option and curated tier shortcuts
- Enable phased streaming (Phase 1 non-streaming, Phase 2 token streaming)
- Centralize telemetry and cost tracking

Non‑Goals
- Provider must not hold OpenRouter API keys
- Provider does not directly call OpenRouter
- No bespoke prompt guards beyond Cipher policies

3) User Stories

- As a developer, selecting “Auto (Cipher)” should choose a suitable model/tier automatically with no extra config.
- As an operator, I need a quick way to force a tier (e.g., m3) for a single request.
- As a platform owner, I need centralized logging, cost accounting, and safe headers applied once, not per extension.

4) High-Level Architecture

Sequence (non-streaming, Phase 1)
1. VS Code LM API → Provider receives request
2. Provider builds payload for Cipher (task_description + OpenAI-style messages)
3. Provider calls Cipher Aggregator JSON-RPC over SSE/HTTP (localhost)
4. Cipher calls llm-inference MCP → llm_inference_auto
5. llm-inference MCP:
   - Builds OpenRouter request with headers via [python.get_openrouter_headers()](servers/llm-inference-mcp.py:94)
   - POSTs to OpenRouter [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:456)
   - Returns final text to Cipher
6. Provider returns text to VS Code LM API

Two-level routing confirmation
- Level 1 (Cipher): route to MCP server [docs/llm-inference-setup.md](docs/llm-inference-setup.md:305), [docs/llm-inference-setup.md](docs/llm-inference-setup.md:312)
- Level 2 (LLM MCP): auto-select or override tier → model, temps, tokens [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:500)

Rationale
- Cached prompt path: “All requests go through OpenRouter for caching benefits.” [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:378)
- Recommended VS Code LM flow: [data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md](data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md:11), [data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md](data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md:26)

5) Components

- VS Code LM Provider (TypeScript)
  - Registers as an LM provider surfaced as “Cipher (Auto)”
  - Bridges VS Code messages → Cipher JSON-RPC calls
  - Optional quick-pick to choose override tier per request
- Cipher Aggregator client
  - JSON-RPC over HTTP(S)/SSE to http://localhost:3020/sse with sessionId
- llm-inference MCP server (Python)
  - Tool: llm_inference_auto (auto routing + override_tier)
  - Tool: list_available_models (operator UI in IDE)
  - Calls OpenRouter with modeled headers [python.get_openrouter_headers()](servers/llm-inference-mcp.py:94)

6) Provider UX / Settings

Surfaced models
- Auto (Cipher) — default
- Tier L0 (Speed)
- Tier M1 (Light)
- Tier M2 (Balanced)
- Tier M3 (Deep Reasoning)
- Tier M4 (Max Reasoning)

Settings (extension > settings.json)
- cipherAggregatorUrl (default http://localhost:3020/sse)
- sessionId (default generated per window)
- defaultModel (“auto-cipher”)
- allowTierOverride (true)
- telemetryEnabled (true)

7) Data Contracts

Provider → Cipher Aggregator (JSON-RPC)
- Method: tools/call
- name: "llm_inference_auto"
- arguments:
  {
    "task_description": "Short intent or request title",
    "messages": [
      {"role":"system","content":"..."},
      {"role":"user","content":"..."}
    ],
    "override_tier": "m3" // optional
  }

Notes
- override_tier schema: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:402)
- Invocation site: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:422)
- POST to OpenRouter: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:456)
- Tier inventory tool: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:500)

Response (Phase 1)
- Final text content (non-streaming)
- Optionally includes tier/model and token usage in body [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:480)

“Model info” command (IDE)
- Calls list_available_models and renders returned text [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:500)

8) Tier Mapping (LLM MCP)

- Tier selection is automatic unless override_tier is provided; the mapping of tier→model/temp/tokens is maintained in llm-inference MCP, not the provider.
- Central defaults and eval LLM config live in [cipher.yml](cipher.yml:61), [cipher.yml](cipher.yml:349)

9) Streaming Plan

Phase 1 (initial release)
- Non-streaming: provider waits for full result text from llm_inference_auto
- Minimal complexity; immediate compatibility with current Python tool

Phase 2 (token streaming)
Server changes (llm-inference MCP)
- Add OpenRouter streaming (payload: "stream": true) and consume chunks
- Emit MCP incremental outputs (or chunked envelopes) back to Cipher
- Preserve headers via [python.get_openrouter_headers()](servers/llm-inference-mcp.py:94)

Provider changes
- Use EventSource or fetch reader to receive streamed tokens from Cipher Aggregator
- Forward tokens to VS Code LM streaming API
- Finalize completion on “done” signal

10) Security and Privacy

- Provider never stores or transmits OpenRouter API keys (lives only in Cipher)
- Transport to Cipher is local HTTP/SSE; consider localhost-only by default
- Redact secrets in prompt telemetry; log minimal routing metadata client-side
- Honor workspace trust; disable provider if workspace isn’t trusted (configurable)

11) Error Handling and Fallbacks

Error classes
- Cipher Aggregator unavailable: show actionable message (status, retry option)
- llm-inference tool error: display tool error body; expose “retry different tier” CTA
- OpenRouter HTTP error: message includes status/body [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:491)
- Unknown tool: fail-fast [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:512)

Fallback strategy
- Offer retry with a lighter tier (l0/m1) for timeouts
- Offer retry with “auto” if a forced override fails

12) Performance

- Local hop (provider → Cipher) is sub-10ms typical
- Remote LLM latency dominates; no extra remote LLM network round-trips
- Streaming adds near-zero overhead vs. direct OpenRouter streaming
- Budget: p95 provider overhead < 25ms; total RTT ~ OpenRouter p95 ± few ms

13) Telemetry & Cost Tracking

- Correlate requests with sessionId
- Server-side cost logging already present (openrouter-costs.jsonl); ensure enabled [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:470)
- Provider telemetry events (success, duration, tier used, error classes) without prompt content
- Optional anonymized metrics toggle in settings

13a) Response Metadata Display (Critical Feature)

**Requirement**: Display routing metadata at the end of every Cline response

**Format**:
```
---
🤖 [Tier: {tier} | Model: {model_name} | Cost: ${cost} | Tokens: {prompt}→{completion} ({total})]
```

**Implementation**:

Server Response Format (Already Implemented):
- llm_inference_auto returns metadata in response footer
- Machine-readable HTML comment: `<!-- METADATA: tier=m2|model=minimax-01|cost=0.001827|tokens=731,1009,1740 -->`
- Human-readable cost line included in text

Provider Extraction Logic (To Be Implemented):
1. Parse llm_inference_auto response for metadata comment
2. Extract tier, model, cost, tokens from pipe-delimited format
3. Append formatted metadata line to final Cline response
4. Display must appear at very end of response

**Parsing Logic**:
```typescript
// Extract metadata from response
const metadataMatch = response.match(/<!-- METADATA: tier=([^|]+)\|model=([^|]+)\|cost=([^|]+)\|tokens=([^,]+),([^,]+),([^ ]+) -->/);

if (metadataMatch) {
  const [_, tier, model, cost, promptTokens, completionTokens, totalTokens] = metadataMatch;

  // Format display line
  const metadataLine = `\n---\n🤖 [Tier: ${tier} | Model: ${model} | Cost: $${cost} | Tokens: ${promptTokens}→${completionTokens} (${totalTokens})]`;

  // Append to response
  return originalResponse + metadataLine;
}
```

**Why This Matters**:
- Primary use case for VS Code LM Provider
- Transparency: Users see which model/tier was used
- Cost tracking: Real-time cost visibility
- Auditability: Track routing decisions

**See Also**:
- Original ticket: [vscode-lm-provider-cipher/tickets/implement-response-metadata-enrichment.md](vscode-lm-provider-cipher/tickets/implement-response-metadata-enrichment.md)
- Architecture issue: [vscode-lm-provider-cipher/tickets/metadata-display-architecture-issue.md](vscode-lm-provider-cipher/tickets/metadata-display-architecture-issue.md)

14) Testing Strategy

Unit (Provider)
- Request shaping: messages → OpenAI format; task_description generation
- Settings parsing; override tier selection

Integration
- With live Cipher Aggregator (localhost), exercise:
  - Auto selection (no override) [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:433)
  - Manual override_tier path [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:402)
  - list_available_models UI

E2E
- VS Code extension host tests: invoke LM completions; verify outputs and telemetry
- Fault injection: aggregator down, tool errors, OpenRouter 429/5xx

Streaming (Phase 2)
- Token arrival order, flow control, finalization markers

15) Rollout Plan

- Phase 0: Non-streaming provider behind feature flag (“cipherLmProvider.enabled”)
- Phase 1: Default to “Auto (Cipher)” for internal users; gather latency/cost metrics
- Phase 2: Enable streaming, roll out broadly
- Migration: document steps to select the provider in IDE; add command palette actions

16) Acceptance Criteria

- Functional
  - Auto and override tiers produce valid completions via OpenRouter path
  - “Model info” lists tiers from MCP tool
- Performance
  - p95 provider overhead < 25ms (non-streaming)
- Reliability
  - Successful completion rate (no provider-side failures) ≥ 99.5%
- Observability
  - Requests logged with tier and duration; costs captured server-side
- Security
  - Provider holds no OpenRouter key; traffic limited to Cipher Aggregator

17) Open Questions

- Do we need per-workspace default tier hints (metadata) passed to llm-inference?
- Should provider expose “explain routing decision” by calling a routing-metadata tool?
- For streaming, should Cipher aggregate partials or pass through raw OpenRouter SSE?

18) Future Work

- “Explain my tier” function using routing metadata service
- Inline token cost estimates to IDE
- Guardrail feedback loops (e.g., detect PII, request downgrade to safer tier)
- Cross-IDE provider parity (JetBrains)

19) Implementation Notes and References

- Use existing server and docs:
  - llm-inference MCP: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py)
  - OpenRouter headers: [python.get_openrouter_headers()](servers/llm-inference-mcp.py:94)
  - POST to OpenRouter: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:456)
  - list_available_models tool: [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:500)
  - Two-level routing flow: [docs/llm-inference-setup.md](docs/llm-inference-setup.md:305), [docs/llm-inference-setup.md](docs/llm-inference-setup.md:312)
  - Cipher defaults: [cipher.yml](cipher.yml:61), [cipher.yml](cipher.yml:349)
  - VS Code → Cipher LM recommendation: [data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md](data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md:11), [data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md](data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md:26)

20) Milestones

- M0: Spec approval (this document)
- M1: Provider skeleton, non-streaming happy path
- M2: Override tier UI + Model info command
- M3: Telemetry + performance validation, internal default enablement
- M4: Streaming end-to-end
- M5: General availability

Appendix A — Example JSON-RPC Calls

List available tools
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 1,
  "params": { "sessionId": "session-<ts>-<rand>" }
}

Call llm_inference_auto
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "id": 2,
  "params": {
    "sessionId": "session-<ts>-<rand>",
    "name": "llm_inference_auto",
    "arguments": {
      "task_description": "Summarize these changes",
      "messages": [
        {"role":"system","content":"You are a helpful assistant."},
        {"role":"user","content":"Summarize the diff in 5 bullets."}
      ],
      "override_tier": "m2"
    }
  }
}

Appendix B — VS Code Provider Behaviors

- Default model label: “Auto (Cipher)”
- Curated menu for tier overrides per request
- “Model info (Cipher)” command that renders [servers/llm-inference-mcp.py](servers/llm-inference-mcp.py:500) output
- Respect workspace trust and show diagnostics on connection failures

File Location
- This spec lives under [specs/vscode-lm-provider-cipher-spec.md](specs/vscode-lm-provider-cipher-spec.md)
