# MCPJungle Technology Documentation

**GitHub**: https://github.com/wong2/mcp-jungle
**Stars**: ~692
**Status**: Selected Hub (Decision Made)
**Purpose**: Primary MCP hub for managing all external MCP servers and tool groups
**Instance Name**: `jarvis` (our MCPJungle hub)

---

## Overview

MCPJungle is a simpler alternative to MetaMCP that provides:
- 📦 Tool group organization
- 🔧 Basic server management
- ⚡ Lightweight design
- 🚀 Easy setup

---

## Research Required

**⚠️ Documentation to be completed using Context7 + GPT-Researcher**

### Context7 Query Plan
```
Library: wong2/mcp-jungle
Focus Areas:
- Installation procedures
- Configuration format
- Tool group setup
- Comparison with MetaMCP
- Use cases and limitations
```

### GPT-Researcher Query Plan
```json
{
  "task": "Research MCPJungle features, setup, and comparison with MetaMCP",
  "report_type": "research_report"
}
```

---

## Why Chosen for Our Architecture

Decision update (2025-11-17):
- ✅ Simpler and lighter than MetaMCP
- ✅ Aligns with current needs (CLI-first, no Web UI required)
- ✅ Tool groups are sufficient for our server organization
- ✅ Easier to operate alongside Cipher as a tools/memory layer
- ⚠️ No built-in Web UI (acceptable tradeoff)
- ⚠️ Less advanced middleware than MetaMCP, but good enough for our scale

**Decision**: MCPJungle is our primary MCP hub. Future MCP servers (httpie, schemathesis, pytest, gpt-researcher, memory-bank, Context7, Morph, etc.) will be aggregated behind MCPJungle once the hub is deployed.

---

## Research Deliverables

- [ ] Installation guide
- [ ] Configuration examples
- [ ] Feature comparison with MetaMCP
- [ ] Migration path (if needed later)
