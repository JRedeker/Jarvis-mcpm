Title: VS Code LM → Cipher → OpenRouter routing plan

Goal
Route all VS Code LM API requests through a dedicated provider that forwards prompts to Cipher, which then smart-routes to OpenRouter.

Why
- Centralize auth/policy/observability in Cipher
- Keep model/vendor agility via Cipher evalLlm tiering
- Reuse Cipher’s routing, cost controls, telemetry, and guardrails

High-level Flow
1) VS Code LM API (user’s extension call)
2) Cipher LM provider (VS Code LM provider extension)
3) Cipher Aggregator endpoint (SSE) → llm_inference_auto tool
4) evalLlm picks best tier/model
5) OpenRouter executes the model call
6) Streamed tokens returned up the chain to VS Code

Provider Responsibilities
- Convert LM API chat/completion requests to cipher llm_inference_auto payload:
  - task_description: short purpose string
  - messages: normalized to OpenAI-style chat format already supported by Cipher
  - optional override_tier/model hints passed through as metadata
- Streaming: bridge Cipher SSE streaming to VS Code LM streaming tokens
- Error mapping: translate Cipher/OpenRouter errors to LM provider errors (with retry hints)
- Model surfacing: expose a minimal list (e.g., “auto (Cipher)”, plus curated options) and treat “auto” as Cipher smart routing

Kilo Code Integration
- In Kilo Code, select the VS Code LM provider; if the default provider is Cipher LM, all Kilo Code LM calls route through Cipher automatically
- No API keys stored in Kilo Code; Cipher holds credentials and executes routing

Considerations
- Rate limits & timeouts: enforce centrally in Cipher; surface friendly messages in VS Code UI
- Privacy: ensure workspace snippets sent to Cipher comply with user/org policy
- Tool calling: phase 1 = pure chat/completion; phase 2 can map model tool-calls to MCP tools if/when needed

Minimal Implementation Plan
- Build Cipher LM provider extension for VS Code LM API
- Provider settings: Cipher URL, project/session, optional default tier/model hint
- Map request/response, add streaming, basic error translation
- Validate with Kilo Code end-to-end
