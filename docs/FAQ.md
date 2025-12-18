# Frequently Asked Questions

Common questions about Jarvis and MCPM.

---

## General

### What makes Jarvis different from other MCP servers?

Unlike gateways that just forward calls, Jarvis is an **intelligent presentation layer**:

| Feature | Standard Gateway | Jarvis |
|---------|-----------------|--------|
| Output | Raw CLI + ANSI codes | Clean Markdown with emojis |
| Validation | None (errors after execution) | Pre-execution checks with suggestions |
| DevOps | Manual setup | Auto CI/CD, pre-commit, secrets |
| Recovery | Requires human intervention | Self-healing |
| Tools | Fixed toolset | 200+ installable on-demand |

**Bottom line:** Other gateways forward calls. Jarvis transforms raw CLI into agent-optimized responses.

### Which AI models does Jarvis support?

**Any model with MCP support:**
- Anthropic Claude (all versions)
- OpenAI GPT-4, GPT-4o
- Google Gemini
- DeepSeek
- Meta Llama
- Mistral
- Local models via Ollama

### Which AI clients does Jarvis support?

**Any client that speaks MCP:**
- Claude Desktop
- Claude CLI
- Cursor
- Windsurf
- VS Code (with MCP extension)
- Zed
- OpenCode
- Kilo Code

---

## Installation & Setup

### Do I need Docker?

**Yes**, for the infrastructure layer (PostgreSQL, Qdrant). These power memory capabilities and vector search.

```bash
# Linux
sudo apt install docker.io docker-compose-plugin

# macOS
brew install --cask docker

# Start infrastructure
./scripts/manage-mcp.sh start
```

Or let the agent run `jarvis_system({ action: "bootstrap" })` to handle setup.

### Can I use Jarvis without Docker?

Partially. Jarvis itself works without Docker, but you lose:
- Persistent memory (PostgreSQL)
- Vector search (Qdrant)
- MCPM daemon (profile HTTP endpoints)

For full functionality, Docker is required.

### How do I update Jarvis?

```bash
cd Jarvis-mcpm
git pull
cd Jarvis && go build -o jarvis .
```

Then restart your AI client.

---

## Tools & Capabilities

### Don't most clients already have web search and file tools?

Yes, but MCPM's registry offers **specialized** tools beyond built-in capabilities:

| Built-in | MCPM Specialized Tool |
|----------|----------------------|
| Generic web search | `context7`: Structured library docs with code snippets |
| Basic URL fetch | `firecrawl`: Intelligent scraping with table extraction |
| Read files | `pdf-parse`: Extract text, tables, metadata from PDFs |
| Git commands | `gitleaks`: Pre-commit hook that blocks secrets |
| Fixed toolset | **200+ servers** installable mid-conversation |

The unique value is **dynamic capability expansion** - your agent installs tools on-demand.

### How do I find available tools?

```javascript
// Search by capability
jarvis_server({ action: "search", query: "pdf" })

// List all installed
jarvis_server({ action: "list" })

// Get details
jarvis_server({ action: "info", name: "pdf-parse" })
```

### Can I create custom MCP servers?

Yes. Register a custom server:

```javascript
jarvis_server({
  action: "create",
  name: "my-server",
  command: "node",
  args: "/path/to/server.js"
})
```

Or for HTTP servers:

```javascript
jarvis_server({
  action: "create",
  name: "my-api",
  type: "streamable-http",
  url: "http://localhost:3000/mcp"
})
```

---

## Security

### How does Jarvis handle secrets?

**Pre-Commit Protection:**
1. `jarvis_project({ action: "devops" })` installs `gitleaks` hook
2. Scans commits for API keys, tokens, passwords
3. Blocks commit if secrets detected
4. Forces agent to use environment variables

**Additional Security:**
- Input validation (prevents command injection)
- Sandboxed Docker containers
- No automatic push to remote
- Audit logs for all operations

### Is my data sent anywhere?

**No.** Jarvis runs entirely locally:
- All tools execute on your machine
- Docker containers are local
- No telemetry or data collection
- Network calls only when you use web tools (Brave, Context7, etc.)

---

## Performance

### How fast is Jarvis?

Go-powered performance:

| Metric | Value |
|--------|-------|
| Cold Start | <100ms |
| Hot Path | <10ms |
| Memory (idle) | ~20 MB |
| Memory (load) | ~50 MB |

### How does Jarvis compare to Node.js MCPM CLI?

| Metric | Node.js CLI | Go Jarvis | Improvement |
|--------|-------------|-----------|-------------|
| Startup | 400-600ms | <100ms | **6x faster** |
| Memory | 80-120 MB | ~20 MB | **4-6x smaller** |

### Why is my first command slow?

Common causes:
1. **Docker cold start** - First command after boot warms up containers
2. **npm install** - Network-dependent tool installation
3. **Registry search** - First search may be slower

Subsequent commands use cached resources and are much faster.

---

## Profiles & Configuration

### What is the 3-Layer Profile Stack?

Composable configuration that eliminates duplication:

| Layer | Purpose | Example |
|-------|---------|---------|
| **Layer 1: Environment** | Workspace-specific | `project-frontend`, `project-backend` |
| **Layer 2: Client** | AI client-specific | `client-cursor`, `client-claude` |
| **Layer 3: Global** | Always-on capabilities | `memory`, `testing` |

Mix and match layers per project. See [Configuration Strategy](CONFIGURATION_STRATEGY.md).

### How do I switch profiles?

```javascript
// Get recommendations for current directory
jarvis_profile({ action: "suggest" })

// Manually activate
jarvis_client({
  action: "edit",
  client_name: "cursor",
  add_profile: "project-frontend,memory"
})
```

### What ports does Jarvis use?

| Port | Service |
|------|---------|
| 5432 | PostgreSQL |
| 6333 | Qdrant HTTP |
| 6334 | Qdrant gRPC |
| 6275 | MCPM API Server |
| 6276 | p-pokeedge profile |
| 6277 | memory profile |
| 6278 | morph profile |
| 6279 | qdrant profile |

---

## Troubleshooting

### What if I encounter errors?

1. **Run diagnostics:**
   ```javascript
   jarvis_check_status()
   ```

2. **Common fixes:**
   - Docker not running: `sudo systemctl start docker`
   - MCPM not found: `jarvis_system({ action: "bootstrap" })`
   - Port conflicts: Check `docker compose ps`

3. **Self-healing:**
   ```javascript
   jarvis_system({ action: "restart_infra" })
   ```

See [Troubleshooting Guide](TROUBLESHOOTING.md) for detailed solutions.

### How do I reset everything?

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

# Bootstrap (in agent)
jarvis_system({ action: "bootstrap" })
```

---

## Contributing

### How do I contribute?

1. Fork the repo
2. Build locally: `./scripts/setup-jarvis.sh`
3. Run tests: `cd Jarvis && go test -v ./...`
4. Submit PR

See [Developer Guide](DEVELOPER_GUIDE.md) for full setup.

### What contributions are most needed?

- Bug fixes ([Issues](https://github.com/JRedeker/Jarvis-mcpm/issues))
- Documentation improvements
- Test coverage expansion
- Performance optimizations
- New MCP server integrations

---

## See Also

- [Examples](EXAMPLES.md) - Workflow examples
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
- [Configuration Strategy](CONFIGURATION_STRATEGY.md) - 3-Layer Profile Stack
- [API Reference](API_REFERENCE.md) - Full tool documentation
