# Register MCP Server with jarvis - Operational Runbook

**Purpose**: Add a new MCP server to jarvis registry
**Last Updated**: 2025-11-18
**Estimated Time**: 3-5 minutes

---

## 🚀 Quick Registration

```bash
# Register using configuration file
./mcpjungle register -c config.json

# Or register inline (for simple servers)
./mcpjungle register --name <name> --url <url> --description <desc>
```

---

## 📋 Prerequisites Checklist

- [ ] jarvis server is running (`curl http://localhost:8080/health`)
- [ ] MCP server is accessible (test connection first)
- [ ] Configuration file is valid JSON
- [ ] Required environment variables are set
- [ ] Server name is unique (not already registered)

---

## 🔧 Registration Procedures

### Step 1: Verify jarvis is Running
```bash
# Check health
curl http://localhost:8080/health

# List current servers
./mcpjungle list servers
```

### Step 2: Test MCP Server Connectivity

#### For HTTP Servers
```bash
# Test basic connectivity
curl <server-url>/health || curl <server-url>

# Example: Test context7
curl https://mcp.context7.com/mcp
```

#### For STDIO Servers
```bash
# Test command execution
<command> <args> --help

# Example: Test brave-search
npx -y @brave/brave-search-mcp-server --help
```

### Step 3: Create Configuration File

#### HTTP Server Template
```json
{
  "name": "server-name",
  "transport": "streamable_http",
  "description": "Description of what this server does",
  "url": "https://server-url.com/mcp",
  "bearer_token": "optional-auth-token"
}
```

#### STDIO Server Template
```json
{
  "name": "server-name",
  "transport": "stdio",
  "description": "Description of what this server does",
  "command": "npx",
  "args": ["-y", "package-name"],
  "env": {
    "API_KEY": "${API_KEY}"
  },
  "timeout": 60
}
```

### Step 4: Register Server
```bash
# Register with config file
./mcpjungle register -c your-server.json

# Expected output: "Server <name> registered successfully!"
```

---

## ✅ Verification Steps

### Check Registration
```bash
# List all servers
./mcpjungle list servers

# Verify your server appears in list
./mcpjungle list servers | grep <your-server-name>
```

### Discover Tools
```bash
# List tools from your server
./mcpjungle list tools | grep <your-server-name>

# Count tools from your server
./mcpjungle list tools | grep <your-server-name> | wc -l
```

### Test Tool Invocation
```bash
# Test a simple tool (replace with actual tool name)
./mcpjungle invoke <server-name>__<tool-name> --input '{}'

# Example: Test filesystem
./mcpjungle invoke filesystem__list_allowed_directories
```

---

## 📊 Registration Examples

### Example 1: HTTP Server (context7)
```json
{
  "name": "context7",
  "transport": "streamable_http",
  "description": "Documentation lookup via llms.txt",
  "url": "https://mcp.context7.com/mcp"
}
```

### Example 2: STDIO Server (brave-search)
```json
{
  "name": "brave-search",
  "transport": "stdio",
  "description": "Web search via Brave Search API",
  "command": "npx",
  "args": ["-y", "@brave/brave-search-mcp-server"],
  "env": {"BRAVE_API_KEY": "${BRAVE_API_KEY}"},
  "timeout": 60
}
```

### Example 3: Python Server (gpt-researcher)
```json
{
  "name": "gpt-researcher",
  "transport": "stdio",
  "description": "AI-powered research with deep web analysis",
  "command": "/home/jrede/dev/MCP/.venv/bin/python3",
  "args": ["gptr-mcp/server.py"],
  "env": {
    "OPENAI_API_KEY": "${OPENAI_API_KEY}",
    "TAVILY_API_KEY": "${TAVILY_API_KEY}"
  },
  "timeout": 60
}
```

---

## 🔍 Troubleshooting

### Registration Fails
```bash
# Check error message
./mcpjungle register -c config.json 2>&1

# Common issues:
# - Server not accessible: Test connectivity first
# - Invalid JSON: Validate with jq: cat config.json | jq .
# - Missing env vars: Check .env file
# - Timeout: Increase timeout value
```

### Server Shows as Registered but Tools Don't Work
```bash
# Check server logs
./mcpjungle list servers --verbose

# Test individual tool
./mcpjungle invoke <server>__<tool> --input '{}' --verbose

# Check jarvis logs for server stderr
grep -A 10 -B 5 "<server-name>" jarvis.log
```

### STDIO Server Issues
```bash
# Check if command works manually
<command> <args>

# Check environment variables
env | grep <VAR_NAME>

# Test with shorter timeout first
# Increase timeout if server is slow to start
```

---

## 🛠️ Advanced Configuration

### Custom Timeouts
```json
{
  "name": "slow-server",
  "transport": "stdio",
  "timeout": 300,
  "command": "python",
  "args": ["slow-server.py"]
}
```

### Complex Environment
```json
{
  "name": "complex-server",
  "transport": "stdio",
  "env": {
    "API_KEY": "${API_KEY}",
    "DEBUG": "true",
    "CUSTOM_VAR": "value"
  },
  "command": "npx",
  "args": ["-y", "complex-package"]
}
```

### Bearer Token Authentication
```json
{
  "name": "authenticated-server",
  "transport": "streamable_http",
  "url": "https://api.example.com/mcp",
  "bearer_token": "your-secret-token"
}
```

---

## 🧹 Cleanup (If Registration Fails)

### Remove Failed Registration
```bash
# Deregister server
./mcpjungle deregister <server-name>

# Verify removal
./mcpjungle list servers | grep <server-name> || echo "Server removed"
```

### Start Fresh
```bash
# Stop jarvis
pkill mcpjungle

# Clear database (WARNING: removes all servers)
rm mcpjungle.db

# Restart jarvis
./mcpjungle start --port 8080
```

---

## 📈 Success Metrics

**Registration Time**: 1-30 seconds (depends on server type)
**Tool Discovery**: Immediate after registration
**First Tool Test**: <5 seconds
**Error Rate**: Should be 0% with proper configuration

---

## 📝 Post-Registration Checklist

- [ ] Server appears in `./mcpjungle list servers`
- [ ] Tools are discoverable via `./mcpjungle list tools`
- [ ] At least one tool invocation succeeds
- [ ] Server description is accurate
- [ ] No errors in jarvis logs

---

## 📞 Support

If registration fails:
1. Check prerequisites checklist
2. Verify server connectivity independently
3. Review error messages carefully
4. Check jarvis logs for server stderr output
5. Test configuration manually before registration

**Escalation**: Document failed attempts and check `docs/troubleshooting.md`
