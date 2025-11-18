# Start jarvis (MCPJungle) - Operational Runbook

**Purpose**: Start the jarvis MCP server and verify it's operational
**Last Updated**: 2025-11-18
**Estimated Time**: 2-3 minutes

---

## 🚀 Quick Start

```bash
# Start jarvis server
./mcpjungle start --port 8080

# In another terminal, verify it's running
curl http://localhost:8080/health
```

**Expected Output**: `{"status":"ok"}`

---

## 📋 Prerequisites Checklist

- [ ] `./mcpjungle` binary exists and is executable
- [ ] Port 8080 is available (not in use by other services)
- [ ] All required API keys are in `.env` file:
  - `BRAVE_API_KEY`
  - `FIRECRAWL_API_KEY`
  - `MORPH_API_KEY`
  - `TAVILY_API_KEY`
  - `OPENAI_API_KEY`

---

## 🔧 Detailed Startup Procedure

### Step 1: Environment Verification
```bash
# Check binary exists
ls -la ./mcpjungle

# Verify version
./mcpjungle version

# Check API keys
grep -E "BRAVE_API_KEY|FIRECRAWL_API_KEY|MORPH_API_KEY|TAVILY_API_KEY|OPENAI_API_KEY" .env
```

### Step 2: Port Availability Check
```bash
# Check if port 8080 is in use
netstat -tuln | grep :8080 || echo "Port 8080 is available"

# If port is in use, choose alternative:
# ./mcpjungle start --port 8081
```

### Step 3: Start jarvis Server
```bash
# Start in foreground (recommended for development)
./mcpjungle start --port 8080

# Or start in background
nohup ./mcpjungle start --port 8080 > jarvis.log 2>&1 &
```

**Expected Output**:
```
███╗   ███╗ ██████╗██████╗       ██╗██╗   ██╗███╗   ██╗ ██████╗ ██╗     ███████╗
████╗ ████║██╔════╝██╔══██╗      ██║██║   ██║████╗  ██║██╔════╝ ██║     ██╔════╝
██╔████╔██║██║     ██████╔╝      ██║██║   ██║██╔██╗ ██║██║  ███╗██║     █████╗
██║╚██╔╝██║██║     ██╔═══╝  ██   ██║██║   ██║██║╚██╗██║██║   ██║██║     ██╔══╝
██║ ╚═╝ ██║╚██████╗██║      ╚█████╔╝╚██████╔╝██║ ╚████║╚██████╔╝███████╗███████╗
╚═╝     ╚═╝ ╚═════╝╚═╝       ╚════╝  ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚══════╝

MCPJungle HTTP server listening on :8080
```

---

## ✅ Verification Steps

### Health Check
```bash
# Basic health check
curl http://localhost:8080/health

# Detailed health with timing
curl -w "Time: %{time_total}s\n" http://localhost:8080/health
```

**Expected**: `{"status":"ok"}`

### Server Discovery
```bash
# List registered servers
./mcpjungle list servers

# Count total tools
./mcpjungle list tools | wc -l
```

**Expected**: 6 servers, 34 tools

### Tool Test
```bash
# Test a simple tool
./mcpjungle invoke filesystem__list_allowed_directories
```

---

## 🔍 Monitoring & Troubleshooting

### Check Process Status
```bash
# If running in background
ps aux | grep mcpjungle

# Check logs
tail -f jarvis.log
```

### Common Issues

| Issue | Symptom | Solution |
|-------|---------|----------|
| **Port in use** | `bind: address already in use` | Use different port: `--port 8081` |
| **Binary not found** | `command not found: mcpjungle` | Check path: `ls -la ./mcpjungle` |
| **Health check fails** | `Connection refused` | Wait 10s, check if process started |
| **API key missing** | Tool registration fails | Check `.env` file contents |

### Log Analysis
```bash
# Check for errors
grep -i error jarvis.log

# Check for warnings
grep -i warn jarvis.log

# Monitor real-time
tail -f jarvis.log | grep -E "(ERROR|WARN|INFO)"
```

---

## 🛑 Shutdown Procedure

### Graceful Shutdown
```bash
# If running in foreground: Ctrl+C

# If running in background:
pkill -TERM mcpjungle

# Verify shutdown
curl http://localhost:8080/health || echo "Server stopped"
```

### Force Shutdown (if needed)
```bash
pkill -9 mcpjungle
```

---

## 📊 Performance Metrics

**Startup Time**: ~2-3 seconds
**Memory Usage**: ~50-100MB
**CPU Usage**: <1% when idle
**Health Check Response**: <100ms

---

## 🔐 Security Notes

- jarvis binds to `0.0.0.0:8080` (all interfaces) by default
- No authentication in development mode
- Consider firewall rules for production deployments
- API keys are loaded from environment variables

---

## 📞 Support

If jarvis fails to start:
1. Check prerequisites checklist above
2. Review troubleshooting section
3. Check logs for specific error messages
4. Verify all 6 servers are registered after startup

**Emergency Contact**: Check `docs/troubleshooting.md` for escalation procedures
