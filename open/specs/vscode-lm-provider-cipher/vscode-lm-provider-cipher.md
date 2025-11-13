# VS Code LM Provider (Cipher-first) - Specification

**Status**: Proposed
**Priority**: High
**Owner**: Platform/Agents

## Overview

This specification defines a VS Code Language Model Provider extension that routes all IDE language model calls through Cipher Aggregator to OpenRouter, enabling centralized authentication, intelligent tier-based routing, cost tracking, and metadata display.

## Directory Structure

```
open/specs/vscode-lm-provider-cipher/
├── vscode-lm-provider-cipher.md  # This file - spec index & overview
├── spec.md                       # Main specification document
├── tickets/                      # Related tickets (organized in subdirectories)
│   ├── implement-response-metadata-enrichment/
│   ├── metadata-display-verification-report/
│   ├── metadata-display-implementation-summary/
│   └── metadata-display-architecture-issue/
├── notes/                        # Implementation notes, decisions, research
└── assets/                       # Diagrams, mockups, examples
```

## Main Specification

**File**: `./spec.md`

The main specification document covers:
- Architecture and data flow
- Component design
- Security and privacy
- Performance requirements
- Testing strategy
- Rollout plan
- Milestones (M0-M5)

## Key Use Cases

### 1. Automatic Metadata Display ⭐
**Ticket**: `tickets/implement-response-metadata-enrichment.md`

Display routing metadata at the end of every Cline response:
```
---
🤖 [Tier: m2 | Model: minimax-01 | Cost: $0.001827 | Tokens: 731→1009 (1740)]
```

**Status**: Requires VS Code LM Provider implementation
- Server-side ready (cost calculation, metadata generation)
- Client-side blocked (need provider to intercept LLM calls)
- See: `tickets/metadata-display-architecture-issue.md` for details

### 2. Centralized LLM Routing
Route all Cline LLM requests through:
```
Cline → VS Code LM API → Cipher Provider → Cipher Aggregator → llm-inference-mcp → OpenRouter
```

### 3. Cost Tracking & Analytics
- Per-request cost logging (JSONL format)
- Tier-based cost attribution
- Session-level analytics
- Budget monitoring

### 4. Tier Selection Control
- **Auto**: Cipher selects optimal tier based on task
- **Manual**: User/developer overrides tier per request
- **Models**: l0 (speed), m1 (light), m2 (balanced), m3 (reasoning), m4 (max)

## Implementation Status

### ✅ Server-Side Complete
- `servers/llm-inference-mcp.py` - Tier routing, cost calculation, metadata generation
- `cipher.yml` - Configuration and routing rules
- Cost logging to `/home/jrede/dev/MCP/logs/openrouter-costs.jsonl`

### ⏳ Client-Side Pending
- VS Code extension (TypeScript)
- LM Provider registration
- Cipher Aggregator client
- Metadata extraction and display

## Milestones

- **M0**: ✅ Spec approval (this document)
- **M1**: ⏳ Provider skeleton, non-streaming happy path
- **M2**: ⏳ Override tier UI + Model info command
- **M3**: ⏳ Telemetry + performance validation
- **M4**: ⏳ Streaming end-to-end
- **M5**: ⏳ General availability

## Dependencies

**Server Components**:
- ✅ `servers/llm-inference-mcp.py` - LLM inference with tier routing
- ✅ Cipher Aggregator - MCP tool orchestration
- ✅ OpenRouter API - LLM provider

**Client Components**:
- ⏳ VS Code Extension API - LM Provider interface
- ⏳ TypeScript Cipher client - JSON-RPC over SSE
- ⏳ Metadata extraction logic

## Related Documentation

- **Main Spec**: `./vscode-lm-provider-cipher.md`
- **LLM Inference**: `../../../docs/llm-inference-setup.md`
- **Cipher Config**: `../../../cipher.yml`
- **Routing Patterns**: `../../../data/memory-bank/routing-patterns/vscode-lm-to-cipher-routing.md`

## Getting Started

### For Implementers

1. Review main spec: `./spec.md`
2. Study VS Code LM API: https://code.visualstudio.com/api/extension-guides/language-model
3. Review server implementation: `../../../servers/llm-inference-mcp.py`
4. Check existing client patterns in Cline source

### For Stakeholders

1. Read architecture issue: `tickets/metadata-display-architecture-issue.md`
2. Review original requirement: `tickets/implement-response-metadata-enrichment.md`
3. Understand scope and effort in main spec

## Decision Log

### Why This Approach?

**Problem**: Metadata display requires intercepting ALL Cline LLM calls
**Options Considered**:
1. ❌ Prompt engineering (.clinerules) - Can't force tool usage
2. ❌ MCP tool pattern - Cline chooses when to call tools
3. ✅ LM Provider - Intercepts all LLM calls transparently

**Decision**: Implement VS Code LM Provider extension
**Rationale**:
- Only way to guarantee metadata on every response
- Enables centralized cost tracking
- Supports tier-based routing
- Preserves OpenRouter caching benefits

See: `tickets/metadata-display-architecture-issue.md` (Section: Working Solutions)

## Quick Links

- [Main Specification](./spec.md)
- [Architecture Issue Analysis](./tickets/metadata-display-architecture-issue/metadata-display-architecture-issue.md)
- [Server Implementation Summary](./tickets/metadata-display-implementation-summary/metadata-display-implementation-summary.md)
- [Original Metadata Requirement](./tickets/implement-response-metadata-enrichment/implement-response-metadata-enrichment.md)

## Contact & Ownership

**Owners**: Platform/Agents team
**Status**: Proposed - awaiting M0 approval
**Next Steps**: Allocate resources for M1 (Provider skeleton)
