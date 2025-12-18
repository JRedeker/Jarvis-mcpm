# Jarvis Troubleshooting Guide (v3.0)

This guide covers common issues and their solutions when using Jarvis.

---

## Quick Diagnostics

Always start with the `jarvis_check_status()` tool:

```javascript
jarvis_check_status()  // Returns comprehensive system health report
```

This checks:
- MCPM CLI availability
- Docker daemon status
- Daemon process status (supervisor managed services)
- Container health (PostgreSQL, Qdrant)
- HTTP endpoint connectivity
- Node.js and Python availability

---

## Docker Issues

### Docker Daemon Not Running

**Symptoms:**
- `restart_infrastructure()` fails
- "Cannot connect to Docker daemon"
- Container commands timeout

**Solutions:**

```bash
# Linux (systemd)
sudo systemctl start docker
sudo systemctl enable docker  # Start on boot

# macOS
open -a Docker  # Or start from Applications

# Verify
docker ps
```

### Port Conflicts

**Symptoms:**
- "Address already in use" errors
- Services fail to start
- Health checks fail

**Common Ports:**
- `5432` - PostgreSQL
- `6333` - Qdrant HTTP
- `6334` - Qdrant gRPC
- `6275-6278` - Profile HTTP endpoints

**Solutions:**

```bash
# Find process using a port
sudo lsof -i :5432
# or
sudo ss -tlnp | grep 5432

# Kill the process or change ports in docker-compose.yml
```

### Containers Unhealthy

**Symptoms:**
- `jarvis_check_status()` shows containers as unhealthy
- Database connections fail
- Memory tools return errors

**Solutions:**

```javascript
// Let Jarvis fix it
jarvis_system({ action: "restart_infra" })

// If that fails, manually:
```

```bash
cd /path/to/Jarvis-Dev
docker compose down
docker compose up -d
docker compose ps  # Verify health
```

### Container Logs

**Viewing logs for debugging:**

```bash
./scripts/manage-mcp.sh logs
# or
docker compose logs -f postgres
docker compose logs -f qdrant
```

---

## MCPM API Server Issues

### API Server Not Reachable

**Symptoms:**
- `check_status()` shows API server as not reachable
- HTTP transport falls back to CLI
- "Connection refused" on port 6275

**Solutions:**

1. **Check if API server is running:**
   ```bash
   curl http://localhost:6275/api/v1/health
   ```

2. **Start API server manually:**
   ```bash
   cd MCPM
   npm run serve
   # or
   mcpm serve --port 6275
   ```

3. **Via Docker (recommended):**
   ```bash
   docker compose up mcpm-daemon
   ```

4. **Check supervisor status (in container):**
   ```bash
   docker exec mcp-daemon supervisorctl status
   ```

### Transport Selection

**Environment Variables:**

| Variable | Default | Description |
|:---------|:--------|:------------|
| `JARVIS_MCPM_TRANSPORT` | `http` | Transport: `http` or `cli` |
| `MCPM_API_URL` | `http://localhost:6275` | API server URL |

**Force CLI Transport (fallback):**
```bash
export JARVIS_MCPM_TRANSPORT=cli
```

---

## MCPM Issues

### Command Not Found

**Symptoms:**
- `mcpm: command not found`
- Jarvis tools fail with MCPM errors

**Solutions:**

```javascript
// Use Jarvis to fix
jarvis_system({ action: "bootstrap" })
```

Or manually:

```bash
cd MCPM
npm install
npm link
which mcpm  # Should show path
```

### Server Installation Fails

**Symptoms:**
- `jarvis_server(action="install")` returns error
- "Server not found in registry"

**Solutions:**

1. **Check spelling** - Server names are case-sensitive:
   ```javascript
   jarvis_server({ action: "search", query: "brave" })  // Find correct name
   ```

2. **Check registry** - Verify server exists:
   ```javascript
   jarvis_server({ action: "info", name: "brave-search" })  // Get details
   ```

3. **Network issues** - npm install may fail:
   ```bash
   npm cache clean --force
   ```

### Profile Issues

**Symptoms:**
- Profile doesn't activate
- Servers missing from profile
- HTTP connection fails

**Solutions:**

```javascript
// List all profiles
jarvis_profile({ action: "list" })

// Check specific profile
jarvis_server({ action: "info", name: "profile-name" })

// Restart profiles
jarvis_profile({ action: "restart" })
```

---

## Jarvis Binary Issues

### Tools Not Appearing in Client

**Symptoms:**
- Client doesn't show Jarvis tools
- MCP connection fails

**Solutions:**

1. **Verify binary exists:**
   ```bash
   cd Jarvis
   go build -o jarvis .
   ./jarvis -h  # Should show help
   ```

2. **Check client config** - Ensure path is correct:
   ```json
   {
     "mcpServers": {
       "jarvis": {
         "command": "/absolute/path/to/Jarvis/jarvis",
         "args": []
       }
     }
   }
   ```

3. **Restart client** - Changes require client restart

### Build Errors

**Symptoms:**
- `go build` fails
- Import errors

**Solutions:**

```bash
cd Jarvis
go mod tidy     # Fix dependencies
go mod verify   # Verify checksums
go build -v .   # Verbose build
```

### Tests Failing

**Symptoms:**
- `go test ./...` fails
- Coverage dropped

**Solutions:**

```bash
# Run with verbose output
go test -v ./...

# Run specific failing test
go test -v -run TestCheckStatus ./handlers/

# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Connection Issues

### HTTP Connections Failing

**Symptoms:**
- "Connection refused" on HTTP endpoints
- Tools timeout
- Agent loses connection

**Solutions:**

1. **Check daemon is running:**
   ```bash
   docker compose ps | grep daemon
   ```

2. **Verify HTTP ports:**
   ```bash
   curl http://localhost:6276/health  # Should return OK
   ```

3. **Restart daemon:**
   ```javascript
   jarvis_profile({ action: "restart" })
   ```

### Agent Can't Reach Jarvis

**Symptoms:**
- Jarvis tools don't appear
- "MCP server disconnected"

**Solutions:**

1. **Check Jarvis is running** (stdio mode):
   - Jarvis runs as subprocess of your client
   - Check client config has correct path

2. **Check Jarvis is running** (HTTP mode):
   ```bash
   ./Jarvis/jarvis -http -port 6275
   curl http://localhost:6275/health
   ```

---

## Performance Issues

### Slow Tool Execution

**Symptoms:**
- Tools take >5 seconds
- Timeouts

**Common Causes:**
1. **Network latency** - npm/registry calls are slow
2. **Docker startup** - First call after boot is slow
3. **Large operations** - `jarvis_project(action="devops")` does file I/O

**Solutions:**

```bash
# Pre-warm Docker
docker compose up -d

# Check container resources
docker stats
```

### High Memory Usage

**Symptoms:**
- System becomes sluggish
- Docker crashes

**Solutions:**

```bash
# Check container memory
docker stats --no-stream

# Limit container memory in docker-compose.yml:
# services:
#   postgres:
#     mem_limit: 512m
```

---

## Logging and Debugging

### Jarvis Logs

```bash
# View logs
cat logs/jarvis.log

# Tail logs
tail -f logs/jarvis.log
```

### Debug Environment Variables

```bash
export MCPM_NON_INTERACTIVE=true
export MCPM_FORCE=true
export MCPM_DEBUG=1
```

### Verbose Mode

For debugging specific tools, check the handler implementation in `handlers/handlers.go`.

---

## Common Error Messages

### "invalid arguments"

**Cause:** Missing required parameter

**Solution:** Check tool signature in API reference

### "Server not found in registry"

**Cause:** Server name doesn't exist

**Solution:** Use `jarvis_server({ action: "search", query: "..." })` to find correct name

### "Profile does not exist"

**Cause:** Referencing non-existent profile

**Solution:** Use `jarvis_profile({ action: "list" })` to list profiles

### "Docker daemon not responding"

**Cause:** Docker not running or permission issue

**Solution:** Start Docker, check user is in docker group

---

## Getting Help

1. **Check Status First:**
   ```javascript
   jarvis_check_status()
   ```

2. **Check Logs:**
   ```bash
   cat logs/jarvis.log | tail -50
   ```

3. **Run Diagnostics:**
   ```bash
   ./scripts/manage-mcp.sh status
   docker compose ps
   ```

4. **Report Issues:**
   - [GitHub Issues](https://github.com/JRedeker/Jarvis-mcpm/issues)
   - Include `jarvis_check_status()` output
   - Include relevant log snippets

---

## Reset to Clean State

When all else fails, reset everything:

```bash
# Stop all containers
docker compose down -v

# Clean npm cache
npm cache clean --force

# Rebuild
cd Jarvis && go build -o jarvis .
cd ../MCPM && npm install && npm link

# Restart
./scripts/manage-mcp.sh start

# Bootstrap
# (in your agent)
bootstrap_system()
```
