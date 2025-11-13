# ADR-0001 — Deprecate local LM routing via Cipher; route Cline directly to OpenRouter

Status
- Accepted

Date
- 2025-11-13

Context
- The prior architecture planned a VS Code LM Provider that routed IDE LM calls to Cipher Aggregator, which in turn invoked the llm-inference MCP to reach OpenRouter. This centralized auth, routing policy, and cost telemetry in Cipher while maintaining OpenRouter caching and reliability.
- Product direction has changed: we now want Cline to call OpenRouter directly for LM calls and keep Cipher strictly as the MCP server and tool aggregator.

Decision
- Deprecate and archive the VS Code Cipher-first LM Provider.
- Remove any IDE LM routing through Cipher; no LM traffic should be sent to Cipher.
- Cline calls OpenRouter directly for all LM requests.
- Cipher remains the MCP server and aggregator for domain-specific tools only.

Scope
- In scope:
  - Provider deprecation and archival under clients/vscode-lm-provider
  - Documentation updates to reflect OpenRouter-direct flow
  - Removal of provider-specific CI and integration harness
  - Migration guidance for teams currently targeting Cipher for LM
- Out of scope:
  - Any changes to Cipher’s role as an MCP aggregator for tools
  - Any changes to llm-inference MCP beyond marking its usage as optional/legacy for IDE LM flows

Rationale
- Simplifies the LM call path and removes local routing complexity in the IDE layer.
- Leverages OpenRouter’s direct capabilities and caching without an intermediary hop.
- Reduces maintenance burden around IDE-specific provider packaging and transport logic.

Consequences
- VS Code LM Provider code becomes legacy and must be archived to avoid confusion.
- Documentation must be updated to prevent developers from configuring IDE → Cipher LM routing.
- Teams must configure OpenRouter credentials directly in Cline.

Implementation Plan
1) Mark spec as Deprecated and point to this ADR
   - Update: open/specs/vscode-lm-provider-cipher/spec.md (Status: Deprecated; link to ADR-0001)
2) Archive client package
   - clients/vscode-lm-provider: add ARCHIVED.md, set package.json private, exclude from CI, keep code for historical reference
3) Remove provider-specific CI
   - Delete or disable .github/workflows/vscode-lm-provider-integration.yml if present; avoid adding new workflows
4) Update top-level docs
   - README: configure OpenRouter directly in Cline; remove mentions of IDE → Cipher LM routing
   - data/memory-bank/routing-patterns/*: revise routing pattern to reflect Cline → OpenRouter direct for LM
   - docs/llm-inference-setup.md: clarify llm-inference MCP is optional/legacy for IDE LM; still valid for server-side automation if needed
   - cipher.yml comments: ensure references remain MCP-tooling focused only
5) Migration guide
   - docs/migration/migrate-to-openrouter-direct.md: steps for disabling provider, setting OpenRouter keys in Cline, and verification checks
6) Validation
   - Smoke test Cipher Aggregator tool path (tools/list, tools/call) to confirm MCP functionality remains intact
   - Manual verification that no IDE LM traffic is routed through Cipher

Migration Steps (high level)
- Disable any experimental or feature-flagged LM provider selections in IDE.
- Configure OpenRouter keys in Cline per the Cline documentation.
- Validate completions flow through OpenRouter with expected caching and cost behavior.
- Ensure MCP tool flows (non-LM) continue to operate via Cipher Aggregator.

Risks and Mitigations
- Risk: Stranded references encouraging IDE → Cipher LM routing
  - Mitigation: Deprecation banner in spec and README; explicit migration guide
- Risk: Team workflows relying on provider UX
  - Mitigation: Provide concrete Cline configuration steps and quick checks
- Risk: Loss of centralized LM telemetry in Cipher
  - Mitigation: Utilize OpenRouter logs/telemetry; consider lightweight client-side tracking in Cline if needed

References
- Deprecated spec to update: open/specs/vscode-lm-provider-cipher/spec.md
- Cipher config reference: cipher.yml
- MCP llm-inference server reference: servers/llm-inference-mcp.py
- Integration harness to archive: clients/vscode-lm-provider/test/integration.md, clients/vscode-lm-provider/scripts/run_integration.js

Acceptance Criteria
- Spec clearly marked as Deprecated and linked to this ADR.
- Provider package archived and excluded from CI.
- Docs updated to reflect Cline → OpenRouter direct LM flow.
- Migration guide available and validated with at least one successful migration.