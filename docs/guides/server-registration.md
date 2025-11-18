# MCPJungle Server Registration Guide

**Version:** 1.0
**Date:** 2025-11-18
**Status:** Configuration Reference

---

## Overview

This guide provides detailed instructions for registering various MCP servers with MCPJungle (jarvis), including configuration examples and best practices.

## Registration Methods

### **Method 1: Command Line Registration**
```bash
# Register HTTP server
mcpjungle register --name context7 --url https://mcp.context7.com/mcp

# Register STDIO server with config file
mcpjungle register -c server-config.json
```

### **Method 2: Configuration File Registration**
```bash
# Create config file
mcpjungle register -c ./servers/context7.json

# Register multiple servers
mcpjungle register -c ./servers/brave-search.json
mcpjungle register -c ./servers/filesystem.json
```

---

## Server Registration Templates

### **HTTP/Streamable-HTTP Servers**

#### **context7 (Documentation)**
```json
{
  "name": "context7",
  "transport": "streamable_http",
  "description": "Documentation lookup via llms.txt",
  "url": "https://mcp.context7.com/mcp",
  "bearer_token": "optional-api-token"
}
```

**Registration Command:**
```bash
mcpjungle register -c context7.json
```

#### **Semgrep (Security)**
```json
{
  "name": "semgrep",
  "transport": "streamable_http",
  "description": "Static analysis security scanning",
  "url": "https://mcp.semgrep.ai/mcp/",
  "bearer_token": "your-semgrep-token"
}
```

### **STDIO Servers**

#### **brave-search (Web Search)**
```json
{
  "name": "brave-search",
  "transport": "stdio",
  "description": "Web search via Brave Search API",
  "command": "npx",
  "args": ["-y", "@brave/brave-search-mcp-server"],
  "env": {
    "BRAVE_API_KEY": "${BRAVE_API_KEY}"
  },
  "timeout": 60
}
```

**Registration Command:**
```bash
# Set environment variable first
export BRAVE_API_KEY="your-brave-api-key"

# Register server
mcpjungle register -c brave-search.json
```

#### **filesystem (File Operations)**
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "File system operations with security restrictions",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/host"],
  "timeout": 30
}
```

#### **firecrawl (Web Crawling)**
```json
{
  "name": "firecrawl",
  "transport": "stdio",
  "description": "Web crawling and content extraction",
  "command": "npx",
  "args": ["-y", "firecrawl-mcp"],
  "env": {
    "FIRECRAWL_API_KEY": "${FIRECRAWL_API_KEY}"
  },
  "timeout": 120
}
```

#### **morph-fast-apply (Code Editing)**
```json
{
  "name": "morph-fast-apply",
  "transport": "stdio",
  "description": "AI-powered code editing and refactoring",
  "command": "npx",
  "args": ["-y", "@morph-llm/morph-fast-apply"],
  "env": {
    "MORPH_API_KEY": "${MORPH_API_KEY}",
    "ALL_TOOLS": "false"
  },
  "timeout": 60
}
```

#### **gpt-researcher (AI Research)**
```json
{
  "name": "gpt-researcher",
  "transport": "stdio",
  "description": "AI-powered research and report generation",
  "command": "/home/jrede/dev/MCP/.venv/bin/python3",
  "args": ["servers/gpt_researcher_mcp.py"],
  "env": {
    "TAVILY_API_KEY": "${TAVILY_API_KEY}",
    "OPENAI_API_KEY": "${OPENAI_API_KEY}"
  },
  "timeout": 300
}
```

---

## Complete Registration Process

### **Step 1: Prepare Environment**
```bash
# Create servers directory
mkdir -p config/jarvis/servers

# Set up environment variables
export BRAVE_API_KEY="your-brave-api-key"
export FIRECRAWL_API_KEY="your-firecrawl-api-key"
export MORPH_API_KEY="your-morph-api-key"
export TAVILY_API_KEY="your-tavily-api-key"
export OPENAI_API_KEY="your-openai-api-key"
```

### **Step 2: Create Configuration Files**

Create individual JSON files for each server:

```bash
# Create context7 config
cat > config/jarvis/servers/context7.json << 'EOF'
{
  "name": "context7",
  "transport": "streamable_http",
  "description": "Documentation lookup via llms.txt",
  "url": "https://mcp.context7.com/mcp"
}
EOF

# Create brave-search config
cat > config/jarvis/servers/brave-search.json << 'EOF'
{
  "name": "brave-search",
  "transport": "stdio",
  "description": "Web search via Brave Search API",
  "command": "npx",
  "args": ["-y", "@brave/brave-search-mcp-server"],
  "env": {
    "BRAVE_API_KEY": "${BRAVE_API_KEY}"
  },
  "timeout": 60
}
EOF

# Create filesystem config
cat > config/jarvis/servers/filesystem.json << 'EOF'
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "File system operations with security restrictions",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/host"],
  "timeout": 30
}
EOF

# Create firecrawl config
cat > config/jarvis/servers/firecrawl.json << 'EOF'
{
  "name": "firecrawl",
  "transport": "stdio",
  "description": "Web crawling and content extraction",
  "command": "npx",
  "args": ["-y", "firecrawl-mcp"],
  "env": {
    "FIRECRAWL_API_KEY": "${FIRECRAWL_API_KEY}"
  },
  "timeout": 120
}
EOF

# Create morph-fast-apply config
cat > config/jarvis/servers/morph-fast-apply.json << 'EOF'
{
  "name": "morph-fast-apply",
  "transport": "stdio",
  "description": "AI-powered code editing and refactoring",
  "command": "npx",
  "args": ["-y", "@morph-llm/morph-fast-apply"],
  "env": {
    "MORPH_API_KEY": "${MORPH_API_KEY}",
    "ALL_TOOLS": "false"
  },
  "timeout": 60
}
EOF

# Create gpt-researcher config
cat > config/jarvis/servers/gpt-researcher.json << 'EOF'
{
  "name": "gpt-researcher",
  "transport": "stdio",
  "description": "AI-powered research and report generation",
  "command": "/home/jrede/dev/MCP/.venv/bin/python3",
  "args": ["servers/gpt_researcher_mcp.py"],
  "env": {
    "TAVILY_API_KEY": "${TAVILY_API_KEY}",
    "OPENAI_API_KEY": "${OPENAI_API_KEY}"
  },
  "timeout": 300
}
EOF
```

### **Step 3: Register All Servers**
```bash
# Register HTTP servers
mcpjungle register -c config/jarvis/servers/context7.json

# Register STDIO servers
mcpjungle register -c config/jarvis/servers/brave-search.json
mcpjungle register -c config/jarvis/servers/filesystem.json
mcpjungle register -c config/jarvis/servers/firecrawl.json
mcpjungle register -c config/jarvis/servers/morph-fast-apply.json
mcpjungle register -c config/jarvis/servers/gpt-researcher.json
```

---

## Verification Steps

### **List Registered Servers**
```bash
# List all servers
mcpjungle list servers

# Expected output:
# context7
# brave-search
# filesystem
# firecrawl
# morph-fast-apply
# gpt-researcher
```

### **List Available Tools**
```bash
# List all tools from all servers
mcpjungle list tools

# List tools from specific server
mcpjungle list tools --server context7

# List tools with descriptions
mcpjungle list tools --format detailed
```

### **Test Tool Invocation**
```bash
# Test context7
mcpjungle invoke context7__get-library-docs \
  --input '{"library": "lodash/lodash"}'

# Test brave-search
mcpjungle invoke brave-search__brave_web_search \
  --input '{"query": "MCP protocol specification"}'

# Test filesystem
mcpjungle invoke filesystem__read_file \
  --input '{"path": "README.md"}'
```

---

## Advanced Configuration

### **Server Timeouts**
```json
{
  "name": "gpt-researcher",
  "transport": "stdio",
  "description": "AI research with extended timeout",
  "command": "python3",
  "args": ["servers/gpt_researcher_mcp.py"],
  "timeout": 600,  // 10 minutes for long research
  "env": {
    "TAVILY_API_KEY": "${TAVILY_API_KEY}"
  }
}
```

### **Environment Variables**
```json
{
  "name": "brave-search",
  "transport": "stdio",
  "description": "Brave search with API key",
  "command": "npx",
  "args": ["-y", "@brave/brave-search-mcp-server"],
  "env": {
    "BRAVE_API_KEY": "${BRAVE_API_KEY}",
    "BRAVE_SAFESEARCH": "moderate",  // Additional env vars
    "BRAVE_COUNTRY": "US"
  }
}
```

### **Working Directory**
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "Filesystem with custom working directory",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/custom/path"],
  "working_directory": "/home/user/projects"
}
```

---

## Docker-Specific Configuration

### **Volume Mounts for Filesystem Access**
```json
{
  "name": "filesystem",
  "transport": "stdio",
  "description": "Filesystem with Docker volume access",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/host"],
  "working_directory": "/host"
}
```

**Docker Compose Setup**:
```yaml
services:
  mcpjungle:
    image: mcpjungle/mcpjungle:latest-stdio
    volumes:
      - ./data:/data
      - .:/host  # Mount current directory as /host
```

---

## Server Management

### **Enable/Disable Servers**
```bash
# Disable a server (keeps registration, disables tools)
mcpjungle disable server context7

# Re-enable server
mcpjungle enable server context7

# Check server status
mcpjungle get server context7
```

### **Update Server Configuration**
```bash
# Deregister first
mcpjungle deregister context7

# Register with new config
mcpjungle register -c updated-context7.json
```

### **Bulk Operations**
```bash
# Register multiple servers at once
for config in config/jarvis/servers/*.json; do
  mcpjungle register -c "$config"
done

# Disable all servers (for maintenance)
mcpjungle list servers | xargs -I {} mcpjungle disable server {}
```

---

## Troubleshooting

### **Registration Failures**

#### **HTTP Server Connection Issues**
```bash
# Test HTTP endpoint manually
curl -I https://mcp.context7.com/mcp

# Check if server is accessible
curl -X POST https://mcp.context7.com/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

#### **STDIO Server Issues**
```bash
# Test command manually
npx -y @brave/brave-search-mcp-server --help

# Check if npx is available
which npx
npx --version

# Test with environment variables
BRAVE_API_KEY="your-key" npx -y @brave/brave-search-mcp-server
```

#### **Permission Issues**
```bash
# Check file permissions
ls -la config/jarvis/servers/

# Fix permissions
chmod 644 config/jarvis/servers/*.json

# Check if commands are executable
which python3
which npx
```

### **Runtime Issues**

#### **Tool Invocation Failures**
```bash
# Check server logs
docker compose logs mcpjungle | grep -i error

# Test specific tool
mcpjungle invoke brave-search__brave_web_search \
  --input '{"query": "test"}' \
  --verbose

# Check tool availability
mcpjungle list tools --server brave-search
```

#### **Timeout Issues**
```bash
# Increase timeout in config
# Update timeout value in JSON configuration
# Restart server after changes
```

---

## Best Practices

### **Configuration Management**
- Store configs in version control
- Use environment variables for secrets
- Document server dependencies
- Test configurations before production

### **Security**
- Never commit API keys to version control
- Use environment variable substitution
- Restrict filesystem access appropriately
- Validate server URLs and commands

### **Performance**
- Set appropriate timeouts
- Use connection pooling where available
- Monitor server resource usage
- Implement health checks

### **Maintenance**
- Regular backup of configurations
- Monitor for server updates
- Test new server versions
- Document custom configurations

---

## Configuration Templates Repository

Create a repository structure for server configurations:

```
config/jarvis/
├── servers/
│   ├── http/
│   │   ├── context7.json
│   │   ├── semgrep.json
│   │   └── huggingface.json
│   ├── stdio/
│   │   ├── brave-search.json
│   │   ├── filesystem.json
│   │   ├── firecrawl.json
│   │   ├── morph-fast-apply.json
│   │   └── gpt-researcher.json
│   └── memory/
│       ├── memory-bank.json
│       ├── cipher-default.json
│       └── custom-memory.json
├── groups/
│   ├── development.json
│   ├── research.json
│   └── minimal.json
└── environments/
    ├── development.json
    ├── staging.json
    └── production.json
```

---

## Automation Scripts

### **Bulk Registration Script**
```bash
#!/bin/bash
# register-all-servers.sh

SERVERS_DIR="config/jarvis/servers"
echo "Registering all MCP servers..."

for config in "$SERVERS_DIR"/*.json; do
  if [ -f "$config" ]; then
    echo "Registering: $(basename "$config")"
    mcpjungle register -c "$config"

    if [ $? -eq 0 ]; then
      echo "✅ Success: $(basename "$config")"
    else
      echo "❌ Failed: $(basename "$config")"
    fi
  fi
done

echo "Registration complete!"
```

### **Validation Script**
```bash
#!/bin/bash
# validate-servers.sh

echo "Validating MCP server registrations..."

# List all registered servers
echo "Registered servers:"
mcpjungle list servers

# Test each server
for server in $(mcpjungle list servers | tail -n +2); do
  echo "Testing server: $server"
  mcpjungle list tools --server "$server" > /dev/null 2>&1

  if [ $? -eq 0 ]; then
    echo "✅ $server - OK"
  else
    echo "❌ $server - FAILED"
  fi
done
```

---

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Kilo Code | Initial server registration guide |
| 1.1 | TBD | TBD | Add troubleshooting section |
| 1.2 | TBD | TBD | Add automation scripts |

**Status**: ✅ **Complete** - Server registration procedures documented

**Next Steps**:
1. Test registration procedures during Phase 1
2. Update configurations with actual testing results
3. Add any new servers discovered during implementation
