# MCP-MASTER: Dynamic MCP Bus Architecture

**Version:** 2.1
**Date:** 2025-11-16
**Status:** Planning Phase - Todo Tracking System Implemented

--

## Executive Summary

This document defines the architecture for a **MCPJungle (jarvis) + Cipher solution** where our MCPJungle hub instance (`jarvis`) is the primary MCP hub aggregating external MCP servers, and Cipher provides tools and memory capabilities behind that hub.

--

## 📊 Master Todo Tracking

**Overall Project Status**: 🟢 Ready for Phase 1 - API Keys Configured

### Current Phase
**Phase 0: Pre-Research** - ✅ API Keys Complete, Documentation Review Pending

### Phase Progress Summary

| Phase | Status | Progress | Estimated Duration | Target Date |
|-------|--------|----------|-------------------|-------------|
| **Phase 0: Pre-Research** | 🟡 In Progress | 25% (4/16) | 1-2 days | TBD |
| **Phase 1: Infrastructure** | ⏸️ Ready to Start | 0% | 1 week | TBD |
| **Phase 2: Configuration** | ⏸️ Not Started | 0% | 1 week | TBD |
| **Phase 3: Client Migration** | ⏸️ Not Started | 0% | 1 week | TBD |
| **Phase 4: Advanced Features** | ⏸️ Not Started | 0% | 1 week | TBD |
| **Phase 5: Cutover & Cleanup** | ⏸️ Not Started | 0% | 1 week | TBD |

### Quick Stats
- **Total Tasks**: 89
- **Completed**: 4 (API keys configured)
- **In Progress**: 0
- **Blocked**: 0
- **Remaining**: 85

### Critical Blockers
✅ **RESOLVED**: All API keys are configured in `.env` file
- ✅ TAVILY_API_KEY configured
- ✅ OPENAI_API_KEY configured
- ✅ BRAVE_API_KEY configured
- ✅ OPENROUTER_API_KEY configured

**No current blockers** - Ready to proceed with Phase 0 validation and Phase 1 implementation

### Next Actions
1. ✅ ~~Configure API keys in `.env` file~~ (COMPLETE)
2. [ ] Verify gpt-researcher and Context7 are working with current keys
3. [ ] Complete Phase 0 documentation review tasks
4. [ ] Begin Phase 1: Install MetaMCP (`npm install -g mcp-manager`)
5. [ ] Create new `cipher-default.yml` configuration file

--

## 📋 Pre-Research Checklist

> **📊 NOTE**: All tasks are now tracked in the [Master Todo Tracking](#-master-todo-tracking) section above. This checklist remains for reference but should be updated via the master tracking section.

**Status Update**: ✅ API keys are configured in `.env` - Ready to proceed with validation and Phase 1

**Before beginning Phase 0 research, complete these prerequisites:**

--

## Testing & Validation

### Component Testing

**Test Cipher Default Mode (Port 3021)**

```bash
# Start Cipher in default mode
./mcp-manager.sh start-default

# Test memory tools
curl -X POST http://127.0.0.1:3021/http \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'

# Should return: ask_cipher, cipher_memory_search, cipher_memory_store, etc.
```

**Test MetaMCP (Port 3000)**

```bash
# Start MetaMCP
mcp-manager start --config metamcp.config.json

# Test WebSocket connection
wscat -c ws://localhost:3000

# Test HTTP endpoint
curl http://localhost:3000/health
```

### Integration Testing

**Full Stack Test:**

```bash
# 1. Start all services
./mcp-manager.sh start-all

# 2. Run integration tests
python tests/test_mcp_integration.py

# 3. Verify tool availability
python tests/verify_all_tools.py
```

**Performance Benchmarks:**

```bash
# Tool discovery time
python tests/benchmark_tool_discovery.py

# Memory search performance
python tests/benchmark_memory_search.py

# Concurrent request handling
python tests/benchmark_concurrent.py
```

### Validation Checklist

**Pre-Migration Validation:**

- [ ] All API keys configured in environment
- [ ] Qdrant vector store accessible
- [ ] All MCP servers installable
- [ ] Network ports available (3000, 3021)
- [ ] Sufficient disk space for logs

**Post-Migration Validation:**

- [ ] All tools available in IDE
- [ ] Memory operations functional
- [ ] Search operations working
- [ ] File operations validated
- [ ] Performance acceptable
- [ ] Error handling working
- [ ] Logging capturing all events
- [ ] Metrics endpoint accessible


--

## Rollback Procedures

### Emergency Rollback (15 minutes)

**If critical issues discovered:**

1. **Stop MetaMCP:**

```bash
mcp-manager stop
```

2. **Restart old Cipher aggregator:**

```bash
./mcp-manager.sh start
```

3. **Update IDE settings:**

- Revert to old MCP settings (port 3020)
- Restart IDE

4. **Verify rollback:**

```bash
./mcp-manager.sh status
# Should show Cipher aggregator running on port 3020
```

### Phased Rollback

**If issues discovered in specific phase:**

- **Phase 1-2 issues**: Continue using old Cipher aggregator, debug MetaMCP in parallel
- **Phase 3 issues**: Revert IDE settings, keep both systems running
- **Phase 4 issues**: Disable advanced features, keep basic functionality
- **Phase 5 issues**: Immediate rollback to Phase 3 state

### Rollback Decision Matrix

| Issue Type | Severity | Action | Timeline |
|------------|----------|--------|----------|
| All tools unavailable | Critical | Emergency rollback | 15 minutes |
| Some tools missing | High | Phased rollback | 1 hour |
| Performance degradation | Medium | Debug in parallel | 1 day |
| Minor bugs | Low | Continue, fix forward | 1 week |

--

## Operational Guide

### Day-to-Day Operations

**Starting the System:**

```bash
# Start all services
./mcp-manager.sh start-all

# Or start individually
./mcp-manager.sh start-cipher-default  # Port 3021
mcp-manager start --config metamcp.config.json  # Port 3000
```

**Monitoring:**

```bash
# Check health
curl http://localhost:3000/health
curl http://127.0.0.1:3021/health

# View logs
tail -f logs/metamcp-requests.log
tail -f logs/cipher-default.log

# Check metrics
curl http://localhost:3000/metrics
```

**Maintenance:**

```bash
# Rotate logs
./mcp-manager.sh rotate-logs

# Backup configuration
./mcp-manager.sh backup-config

# Update servers
./mcp-manager.sh update-servers
```

### Troubleshooting

**Common Issues:**

1. **MetaMCP won't start:**

   - Check port 3000 availability: `lsof -i :3000`
   - Verify config file syntax: `jsonlint metamcp.config.json`
   - Check logs: `tail -f logs/metamcp-error.log`

2. **Cipher tools not visible:**

   - Verify Cipher running: `curl http://127.0.0.1:3021/health`
   - Check MetaMCP server config for cipher-memory
   - Verify namespace configuration

3. **IDE connection issues:**

   - Check WebSocket URL: `ws://localhost:3000`
   - Verify firewall not blocking port 3000
   - Check IDE MCP settings format

4. **Performance issues:**

   - Review metrics endpoint for bottlenecks
   - Check connection pooling settings
   - Verify Qdrant performance

--

## Appendices

### Appendix A: Tool Reference

**Cipher Internal Tools (13 tools):**

- `cipher_extract_and_operate_memory` - Store knowledge in memory
- `cipher_memory_search` - Search memory for relevant information
- `cipher_memory_store` - Add documents to memory
- `cipher_store_reasoning_memory` - Store reasoning traces
- `cipher_extract_reasoning_steps` - Extract reasoning from text
- `cipher_evaluate_reasoning` - Evaluate reasoning quality
- `cipher_search_reasoning_patterns` - Search for reasoning patterns
- `cipher_bash` - Execute bash commands
- `brave_web_search` - Web search via brave
- `brave_local_search` - Local business search
- `brave_video_search` - Video search
- `brave_image_search` - Image search
- `brave_news_search` - News search
- `brave_summarizer` - AI summarization

**Aggregated MCP Servers (current):**

- `filesystem` - File operations
- `file-batch` - Batch file operations
- `brave-search` - Web search
- `routing-metadata` - Tool metadata

**Standalone MCP Servers (temporary until MetaMCP hub):**

- `httpie` - HTTP API testing (run as direct MCP server)
- `pytest` - Python testing (run as direct MCP server)
- `schemathesis` - API schema testing (run as direct MCP server)
- `gpt-researcher-mcp` - AI research (direct MCP server via `npx gpt-researcher-mcp --stdio`)
- `memory-bank` - Persistent memory (`npx @modelcontextprotocol/server-memory --stdio`)
- `context7` - Documentation lookup (remote streamable-HTTP)
- `morph-fast-apply` - Code editing (remote streamable-HTTP)

**Migration Plan:** Once the new MetaMCP/MetaHub is in place, all of the above standalone MCP servers will be moved behind that hub and the separate IDE registrations will be removed so that Cipher-aggregator (or the new hub) becomes the single MCP entrypoint again.

### Appendix B: Configuration Management

**Environment Variables:**

```bash
# Required
export BRAVE_API_KEY="your_key_here"
export TAVILY_API_KEY="your_key_here"
export OPENAI_API_KEY="your_key_here"
export OPENROUTER_API_KEY="your_key_here"
export QDRANT_URL="http://localhost:6333"

# Optional
export QDRANT_API_KEY=""
export MCP_LOG_LEVEL="info"
```

**Configuration Files:**

- `metamcp.config.json` - Main MetaMCP configuration
- `cipher-default.yml` - Cipher default mode configuration
- `.env` - Environment variables
- `logs/` - Log directory
- `data/` - Data directory (Qdrant, SQLite)

### Appendix C: Security Considerations

**Access Control:**

- MetaMCP runs on localhost only by default
- No authentication required for local access
- For multi-user environments, consider:
  - VPN access only
  - Reverse proxy with authentication
  - Network policies

**API Key Management:**

- Store keys in environment variables, not config files
- Use secret management for production
- Rotate keys regularly
- Limit key permissions

**Data Privacy:**

- Cipher memory stored locally in Qdrant
- No data sent to external services for memory operations
- Search queries go to brave (privacy-focused)
- Review each MCP server’s privacy policy

--

## Document Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 2.1 | 2025-11-16 | Kilo Code | Added Master Todo Tracking system with 89 consolidated tasks, Modification Log, and cross-references |
| 2.0 | 2025-11-16 | Kilo Code | Complete rewrite with hybrid architecture |
| 1.0 | 2025-11-15 | Kilo Code | Initial stdio-based architecture |

**Status**: ✅ **TODO TRACKING SYSTEM IMPLEMENTED** - Ready for Phase 0 (API key configuration)

**Next Steps**:

1. Configure API keys in `.env`
2. Install MetaMCP: `npm install -g mcp-manager`
3. Create `cipher-default.yml` configuration
4. Start parallel testing
