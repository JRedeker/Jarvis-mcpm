# Agent Configuration Guide for Cipher-Aggregator

## ⚠️ ARCHITECTURAL PRINCIPLE: NEVER Access Individual MCP Servers

Agents (Cline, Kilo) are getting errors like:
```
Error: MCP server 'firecrawl' is not configured. Available servers: cipher-aggregator
Error: MCP server 'memory-bank' is not configured. Available servers: cipher-aggregator
```

**THIS IS CORRECT BEHAVIOR!** Agents should NEVER access individual MCP servers directly.

## ✅ CORRECT ARCHITECTURE

```
Agent → cipher-aggregator (tool catalog + intelligent routing) → MCP servers
         ↓                    ↓                            ↓
    Tool Discovery      Semantic Routing             Actual Execution
    (AI-Optimized)   (Embedding-Based)           (Fast & Efficient)
```

## 🔍 What Cipher Provides

### **Built-in Semantic Tool Discovery**
Cipher already has **AI-powered tool routing** with embeddings:

#### **Embedding-Dependent Tools** (Auto-Discovered)
- `cipher_memory_search` - Semantic search over stored knowledge
- `cipher_extract_and_operate_memory` - Knowledge extraction with embeddings
- `cipher_workspace_search` - Team/project context search
- `cipher_search_reasoning_patterns` - Reasoning pattern discovery

#### **Memory Tools** (Knowledge Management)
- `cipher_extract_and_operate_memory` - Extracts knowledge from interactions and immediately applies ADD/UPDATE/DELETE/NONE as one atomic operation
- `cipher_memory_search` - Semantic search over stored knowledge to retrieve relevant facts/code patterns
- `cipher_store_reasoning_memory` - Stores high-quality reasoning traces for future analysis (append-only reflection memory)

#### **Reasoning Tools** (Reflection & Analysis)
- `cipher_extract_reasoning_steps` (internal) - Extracts structured reasoning steps from user input
- `cipher_evaluate_reasoning` (internal) - Evaluates reasoning trace quality and generates improvement suggestions
- `cipher_search_reasoning_patterns` - Searches reflection memory for relevant reasoning patterns

#### **Workspace Tools** (Team Context)
- `cipher_workspace_search` - Searches team/project workspace memory for progress, bugs, PR summaries, and collaboration context
- `cipher_workspace_store` - Background tool capturing team and project signals into workspace memory

#### **Knowledge Graph Tools** (Entity Relationships)
- `cipher_add_node`, `cipher_update_node`, `cipher_delete_node` - Manage entities in the knowledge graph
- `cipher_add_edge` - Create relationships between entities
- `cipher_search_graph`, `cipher_enhanced_search` - Search the graph with basic and enhanced strategies
- `cipher_get_neighbors` - Retrieve related entities around a node
- `cipher_extract_entities` - Extract entities for graph insertion from text
- `cipher_query_graph` - Run graph queries and retrieve structured results
- `cipher_relationship_manager` - Higher-level relationship operations and maintenance

#### **System Tools** (Direct Operations)
- `cipher_bash` - Execute bash commands. Supports one-off or persistent sessions with working dir and timeout controls

### **Why This Architecture Works**

1. **Agents ONLY see cipher-aggregator** (never individual servers)
2. **Cipher provides intelligent tool catalog** with semantic search
3. **Embedding-dependent tools** are automatically excluded when embeddings disabled
4. **Agent-accessible vs internal** - Clear tool visibility separation
5. **Fast lookup** - Pre-categorized tools for rapid discovery
6. **Environment-based inclusion** - Tools respect USE_WORKSPACE_MEMORY=true and KNOWLEDGE_GRAPH_ENABLED=true

## Solution: Proper Agent Configuration

### For Cline Configuration

Update your Cline MCP configuration to point to the cipher-aggregator:

```json
{
  "mcpServers": {
    "cipher-aggregator": {
      "command": "sse",
      "args": ["http://127.0.0.1:3020/sse"]
    }
  }
}
```

### For Kilo Configuration

Update your Kilo MCP configuration:

```json
{
  "mcp": {
    "servers": {
      "cipher-aggregator": {
        "transport": "sse",
        "url": "http://127.0.0.1:3020/sse"
      }
    }
  }
}
```

### Generic MCP Client Configuration

For any MCP client, use this configuration:

```json
{
  "mcpServers": {
    "cipher-aggregator": {
      "command": "sse",
      "args": [
        "http://127.0.0.1:3020/sse"
      ]
    }
  }
}
```

## Tool Discovery & Usage Pattern

### **AI-Optimized Tool Discovery**

Agents should discover tools by:

1. **Connect to cipher-aggregator only**
2. **Request tools/list** to get the complete catalog
3. **Parse tool descriptions** to understand capabilities
4. **Use semantic routing** - let cipher handle intelligent tool selection

### **Correct Tool Access Pattern**

```javascript
// ✅ CORRECT - Use cipher's semantic routing
{
  "method": "tools/call",
  "params": {
    "name": "cipher_memory_search",  // Agent-accessible tool
    "arguments": {
      "query": "configuration patterns",
      "top_k": 5
    }
  }
}

// ✅ CORRECT - Built-in tools with clear intent
{
  "method": "tools/call",
  "params": {
    "name": "cipher_workspace_search",
    "arguments": {
      "query": "team progress on MCP integration",
      "project": "current-project"
    }
  }
}
```

### **What NOT to Do**

```javascript
// ❌ WRONG - Never access individual MCP servers
{
  "method": "tools/call",
  "params": {
    "server": "firecrawl",
    "name": "search",
    "arguments": {...}
  }
}

// ❌ WRONG - Direct server access attempts
{
  "method": "tools/call",
  "params": {
    "name": "firecrawl_search",  // This won't work
    "arguments": {...}
  }
}
```

## Available Tools Through Cipher

### **MCP Aggregated Tools** (Unified Interface)
All external MCP tools are available through cipher with unified naming:

- **Firecrawl**: `firecrawl_scrape`, `firecrawl_search`, `firecrawl_crawl`, `firecrawl_extract`
- **Memory Bank**: `memory_bank_store`, `memory_bank_search`, `memory_bank_list_projects`
- **Httpie**: `make_request`, `upload_file`, `download_file`, `test_api_endpoint`
- **Schemathesis**: `load_openapi_schema`, `test_api_endpoints`, `validate_schema`
- **Pytest**: `pytest_run_tests`, `pytest_coverage`
- **GitHub**: `github_list_repos`, `github_create_repo`
- **Brave Search**: `brave_web_search`, `brave_local_search`, `brave_video_search`
- **Filesystem**: `filesystem_list`, `filesystem_read`, `filesystem_write`
- **Code Index**: `code_index_search`, `code_index_find_files`
- **Context7**: `context7_list_sections`, `context7_get_documentation`

### **Cipher's Built-in AI Tools** (Semantic)
- **Memory & Knowledge**: `cipher_memory_search`, `cipher_extract_and_operate_memory`
- **Workspace Context**: `cipher_workspace_search`, `cipher_workspace_store`
- **Reasoning Patterns**: `cipher_search_reasoning_patterns`
- **Knowledge Graph**: `cipher_search_graph`, `cipher_enhanced_search`
- **System Operations**: `cipher_bash`

## Connection Steps

1. **Establish SSE Connection**: Connect to `http://127.0.0.1:3020/sse`
2. **Get Session ID**: Parse session from SSE response
3. **Initialize Connection**: Send initialize request with session
4. **Discover Tools**: Request `tools/list` to see complete catalog
5. **Parse Tool Catalog**: Understand available capabilities
6. **Call Tools**: Use unified tool names through cipher's routing

## Testing the Configuration

### **Test cipher-aggregator connection:**
```bash
curl -X POST http://127.0.0.1:3020/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "id": 1,
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'
```

### **Test tool discovery:**
```bash
curl -X POST "http://127.0.0.1:3020/sse?sessionId=YOUR_SESSION_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 2
  }'
```

## Environment Configuration

### **Required Environment Variables**
```bash
# For embedding-dependent tools
OPENAI_API_KEY=your_openai_key

# For workspace memory tools
USE_WORKSPACE_MEMORY=true

# For knowledge graph tools
KNOWLEDGE_GRAPH_ENABLED=true

# For specific MCP servers
FIRECRAWL_API_KEY=your_firecrawl_key
BRAVE_API_KEY=your_brave_key
GITHUB_TOKEN=your_github_token
```

## Troubleshooting

### **If tools aren't visible:**

1. **Check cipher-aggregator status**:
   ```bash
   ./mcp-manager.sh status
   ```

2. **Verify embedding configuration**:
   ```bash
   # Ensure embeddings are working for AI-powered discovery
   echo $OPENAI_API_KEY
   ```

3. **Check environment flags**:
   ```bash
   echo $USE_WORKSPACE_MEMORY  # Should be "true"
   echo $KNOWLEDGE_GRAPH_ENABLED  # Should be "true"
   ```

4. **Monitor cipher logs**:
   ```bash
   tail -f logs/cipher-aggregator.log
   ```

### **Common Issues:**

- **"MCP server not found"** - ✅ This is correct! Use cipher-aggregator only
- **Embedding tools missing** - Check OPENAI_API_KEY and embedding configuration
- **Workspace tools missing** - Ensure USE_WORKSPACE_MEMORY=true
- **Knowledge graph tools missing** - Ensure KNOWLEDGE_GRAPH_ENABLED=true
- **Timeout errors** - Check cipher-aggregator startup and API keys

## Configuration Files Location

Update these files with the correct cipher-aggregator configuration:

- **Cline**: `~/.config/cline/mcp.json`
- **Kilo**: `~/.config/kilo/mcp.json`
- **Generic**: Your MCP client's configuration file

## Next Steps

1. **Update agent MCP config** to use `cipher-aggregator` only
2. **Test tool discovery** by requesting `tools/list`
3. **Parse tool descriptions** for AI-optimized understanding
4. **Use semantic routing** through cipher's intelligent catalog
5. **Monitor cipher logs** for embedding and environment status

This configuration leverages cipher's built-in **AI-powered tool routing** and **semantic discovery** while maintaining the correct architectural separation between agents and individual MCP servers.
