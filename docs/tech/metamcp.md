# MetaMCP Technology Documentation

**GitHub**: https://github.com/wong2/mcp-manager
**Stars**: ~1,700
**Status**: Reference Only (Not Selected)
**Purpose**: Enterprise MCP server aggregator with Web UI and advanced features (kept as a reference architecture; MCPJungle is our chosen hub)

---

## Overview

MetaMCP (also known as MCP Manager) is a sophisticated Model Context Protocol server aggregator that provides:
- 🌐 Web-based UI for managing MCP servers
- 📦 Namespace-based server organization
- 🔧 Advanced middleware pipeline
- 🚀 Dynamic tool activation rules
- 🔐 Access control and authentication
- 📊 Monitoring and metrics

---

## Research Required

**⚠️ This documentation will be completed using Context7 + GPT-Researcher once API keys are configured.**

### Context7 Query Plan
```
Library: wong2/mcp-manager
Focus Areas:
- Installation and setup procedures
- Configuration file structure
- Namespace configuration
- Middleware options
- WebSocket server setup
- Integration with streamable-HTTP servers
- Production deployment best practices
- Known issues and limitations
```

### GPT-Researcher Query Plan
```json
{
  "task": "Comprehensive research on MetaMCP (MCP Manager) deployment and configuration",
  "report_type": "research_report",
  "sources": [
    "MetaMCP GitHub repository and documentation",
    "MetaMCP issues and discussions",
    "Community deployment experiences",
    "MCP protocol specification for aggregators"
  ]
}
```

---

## Key Features (Pending Research Validation)

### Web UI
- Visual server management
- Real-time tool availability display
- Namespace configuration interface
- Middleware pipeline editor

### Namespaces
- Organize servers by project type, team, or environment
- Dynamic activation based on context
- Tool group management
- Isolated tool scopes

### Middleware
- Logging and metrics collection
- Rate limiting
- Authentication/authorization
- Request/response transformation

---

## Installation (Pending Research)

```bash
# Placeholder - to be filled from Context7/GPT-Researcher
npm install -g mcp-manager  # Verify correct package name
```

---

## Configuration (Pending Research)

```json
{
  "// Note": "This is a placeholder structure - actual schema TBD from research",
  "servers": {},
  "namespaces": {},
  "middleware": []
}
```

---

## Integration with Cipher

**Plan**: Cipher runs in default mode as an MCP server that MetaMCP aggregates

```yaml
# MetaMCP config (TBD)
servers:
  cipher-memory:
    type: streamable-http
    url: http://127.0.0.1:3021/http
    namespace: default
```

---

## Research Deliverables

Once research completes:
- [  ] Complete installation instructions
- [ ] Full configuration schema
- [ ] Namespace best practices
- [ ] Middleware cookbook
- [ ] WebSocket client configuration
- [ ] Production deployment guide
- [ ] Troubleshooting playbook
- [ ] Performance benchmarks

---

## Next Steps

1. Configure real TAVILY_API_KEY and OPENAI_API_KEY
2. Execute Context7 query on wong2/mcp-manager
3. Execute GPT-Researcher comprehensive report
4. Fill in this documentation with findings
5. Store research in Cipher memory for future reference
