# Jarvis Examples & Workflows

Concrete examples of how an AI Agent uses Jarvis v3.0 to solve real-world problems.

> **Note:** Jarvis v3.0 consolidated 24 tools into 8 action-based tools for ~52% context token reduction.

---

## Quick Reference

| Scenario | Key Tools | Result |
|----------|-----------|--------|
| [DevOps Scaffolding](#1-devops-scaffolding) | `jarvis_project(action="devops")` | Production-ready CI/CD |
| [Self-Healing](#2-self-healing-infrastructure) | `jarvis_check_status()`, `jarvis_system(action="restart_infra")` | Auto-repair services |
| [Tool Installation](#3-dynamic-tool-installation) | `jarvis_server(action="search\|install")` | On-demand capabilities |
| [Security](#4-security-guardrails) | Pre-commit hooks via devops | Block secrets |
| [Context Switching](#5-context-switching) | `jarvis_profile(action="suggest")` | Smart profile selection |
| [Refactoring](#6-intelligent-refactoring) | `morph-fast-apply` | Semantic code transforms |

---

## 1. DevOps Scaffolding

**Problem:** You have a messy legacy repo with no standards, or a blank folder for a new project.

### Conversation

**You:** *"This old script needs to be production-ready."*

**Agent:** *"I've analyzed the directory. It's a Python project lacking linting. I'm applying the standard DevOps stack now."*

**Agent:** *"Done! I've set up pre-commit hooks with Ruff and Gitleaks, created a GitHub Actions workflow, and initialized git. Try making your first commit."*

### How Jarvis Enables This

1. **Analyze:** Jarvis reads file signatures (`pyproject.toml`, `go.mod`, `package.json`) to detect project type
2. **Decide:** Identifies missing configuration (no pre-commit, no CI/CD)
3. **Execute:** Writes production-grade config files

### Tool Calls

```javascript
// Step 1: Analyze project
jarvis_project({ action: "analyze" })
// Returns: { type: "python", missing: ["pre-commit", "ci"] }

// Step 2: Apply DevOps stack
jarvis_project({ action: "devops", project_type: "python" })
```

### Files Created

| File | Purpose |
|------|---------|
| `.pre-commit-config.yaml` | Ruff (linting + formatting), Gitleaks (secrets) |
| `.github/workflows/ci.yml` | Automated testing on push/PR |
| `Makefile` | Common targets: `test`, `lint`, `fmt`, `build` |

---

## 2. Self-Healing Infrastructure

**Problem:** Your RAG application is failing because the local Vector DB crashed.

### Conversation

**You:** *"Why is my search failing?"*

**Agent:** *"Checking system status... Qdrant is down. Restarting infrastructure..."*

**Agent:** *"Done. Services are healthy. Retrying search... Here are your results."*

### How Jarvis Enables This

1. **Diagnose:** Jarvis runs health checks and identifies the unhealthy container
2. **Repair:** Executes Docker restart sequence automatically
3. **Verify:** Waits for health checks to pass before continuing

### Tool Calls

```javascript
// Step 1: Check status
jarvis_check_status()
// Returns: { qdrant: "unhealthy", postgres: "healthy" }

// Step 2: Restart infrastructure
jarvis_system({ action: "restart_infra" })
// Returns: { status: "success", services: ["postgres", "qdrant"] }
```

### What Happens Behind the Scenes

1. Gracefully stops PostgreSQL and Qdrant containers
2. Restarts containers with proper sequencing (DB first)
3. Waits for health checks (PostgreSQL:5432, Qdrant:6333)
4. Confirms all services operational

---

## 3. Dynamic Tool Installation

**Problem:** You ask the Agent to "Analyze this PDF contract," but it has no PDF tools.

### Conversation

**You:** *"Summarize this PDF contract."*

**Agent:** *"I don't have a PDF reader installed. Installing `pdf-parse` via MCPM..."*

**Agent:** *"Tool loaded. Reading contract.pdf... This is a service agreement with payment terms of Net 30 days, auto-renewal clause on page 5..."*

### How Jarvis Enables This

1. **Search:** Agent searches the registry for relevant tools
2. **Install:** Jarvis installs the tool without session restart
3. **Hot-Load:** New capability immediately available

### Tool Calls

```javascript
// Step 1: Search registry
jarvis_server({ action: "search", query: "pdf extraction" })
// Returns: [{ name: "pdf-parse", description: "Extract text from PDFs" }]

// Step 2: Install
jarvis_server({ action: "install", name: "pdf-parse" })
// Returns: { status: "success", package: "@modelcontextprotocol/server-pdf" }

// Step 3: Get info
jarvis_server({ action: "info", name: "pdf-parse" })
// Returns: { tools: ["pdf_extract", "pdf_metadata"], config: {...} }
```

### Why This Matters

- **200+ servers** in the MCPM registry
- Agent "upgrades itself" mid-conversation
- Tools persist for future sessions
- No manual npm/pip installation required

---

## 4. Security Guardrails

**Problem:** The Agent writes code that accidentally hardcodes an API key.

### Conversation

**Agent:** *"I've updated the API client. Committing changes..."*

**Agent:** *"COMMIT BLOCKED: Secret detected in `client.py` line 42."*

**Agent:** *"I've moved the API key to `.env` and updated the code to use `os.getenv('API_KEY')`. Committing again..."*

**Agent:** *"Commit successful. No secrets detected."*

### How Jarvis Enables This

1. **Prevention:** Pre-commit hook runs `gitleaks` scan on staged files
2. **Intervention:** Commit blocked with clear error and line number
3. **Correction:** Agent self-corrects using environment variables

### What Gets Blocked

| Pattern | Example | Action |
|---------|---------|--------|
| API Keys | `api_key = "sk-..."` | Block + report line |
| AWS Secrets | `AWS_SECRET_ACCESS_KEY` | Block + report line |
| Private Keys | `-----BEGIN RSA PRIVATE KEY-----` | Block + report line |
| Passwords | `password = "hunter2"` | Block + report line |

### Result

- Zero secrets in git history
- Agent learns to use `.env` files
- Automatic enforcement (no manual review needed)

---

## 5. Context Switching

**Problem:** You're switching from a backend Go service to a frontend React app. Each needs different tools.

### Conversation

**You:** *"I'm switching to work on the dashboard now."*

**Agent:** *"Detected React project in ~/dashboard. Loading project-frontend profile..."*

**Agent:** *"Ready! I've activated Prettier for formatting, ESLint for linting, and research tools. What would you like to build?"*

### How Jarvis Enables This

1. **Navigate:** You move to a new directory
2. **Detect:** Jarvis analyzes file signatures
3. **Suggest:** Recommends appropriate profile stack
4. **Activate:** Agent switches to context-appropriate tools

### Tool Calls

```javascript
// In new directory
jarvis_profile({ action: "suggest" })
// Returns:
// {
//   layer1: "project-frontend",  // React detected
//   layer2: "client-cursor",     // Current IDE
//   layer3: "memory",            // Always-on
//   reason: "Detected package.json with React dependencies"
// }
```

### 3-Layer Profile Stack

| Layer | Profile | Purpose |
|-------|---------|---------|
| 1 (Environment) | `project-frontend` | Prettier, ESLint, npm |
| 2 (Client) | `client-cursor` | Cursor IDE optimizations |
| 3 (Global) | `memory` | Persistent context |

---

## 6. Intelligent Refactoring

**Problem:** You need to refactor messy authentication logic without breaking things.

### Conversation

**You:** *"Refactor the user authentication logic in `auth.py`. It's messy."*

**Agent:** *"I see `morph-fast-apply` is available. Using semantic transforms..."*

**Agent:** *"Refactored auth.py: extracted validation logic, added type hints, improved error handling. No line-number conflicts."*

### How Jarvis Enables This

1. **Check Tools:** Agent sees `morph-fast-apply` in loaded profile
2. **Transform:** Uses semantic understanding (not line numbers)
3. **Apply:** Changes applied without merge conflicts

### Why Morph vs Traditional Diffs

| Traditional | Morph |
|-------------|-------|
| Line-number based | Semantic understanding |
| Breaks on whitespace | Whitespace-agnostic |
| One change at a time | Batch transforms |
| Manual conflict resolution | Auto-merge |

---

## 7. Research-Driven Development

**Problem:** Agent implements a feature using outdated syntax from training data.

### Conversation

**You:** *"Add server actions for the checkout flow."*

**Agent:** *"The server action syntax has changed since my training. Fetching current docs via Context7..."*

**Agent:** *"Got 47 code snippets for server actions. Implementing with current patterns..."*

**Agent:** *"Done. Used the latest API pattern with proper error boundaries."*

### How Jarvis Enables This

Context7 provides **structured library documentation** - not web search results:

| Built-in Web Fetch | Context7 via Jarvis |
|-------------------|---------------------|
| Generic web pages | npm/PyPI/GitHub docs |
| Raw HTML/text | Structured JSON |
| Random version | Specific version |
| May be outdated | Live, tested snippets |

### Tool Calls

```javascript
// Step 1: Resolve library
context7.resolve_library_id({ libraryName: "next.js" })
// Returns: { id: "/vercel/next.js" }

// Step 2: Get docs
context7.get_library_docs({
  context7CompatibleLibraryID: "/vercel/next.js",
  topic: "server actions"
})
// Returns: { snippets: [...], version: "14.2.0" }
```

---

## 8. Client Configuration

**Problem:** Setting up a new AI client (OpenCode, Cursor, etc.) with proper profiles.

### Conversation

**You:** *"Set up OpenCode with my project profile and memory."*

**Agent:** *"Importing starter config for OpenCode..."*

**Agent:** *"Added p-pokeedge and memory profiles. Config saved to ~/.config/opencode/opencode.json"*

### Tool Calls

```javascript
// Step 1: Import starter config
jarvis_client({
  action: "import",
  client_name: "opencode"
})

// Step 2: Add specific profiles
jarvis_client({
  action: "edit",
  client_name: "opencode",
  add_profile: "p-pokeedge,memory"
})

// Step 3: Verify
jarvis_client({ action: "list" })
```

---

## Tool Quick Reference

### System Operations

```javascript
jarvis_check_status()                              // Health check
jarvis_system({ action: "bootstrap" })             // First-time setup
jarvis_system({ action: "restart" })               // Restart Jarvis
jarvis_system({ action: "restart_infra" })         // Restart Docker services
```

### Server Management

```javascript
jarvis_server({ action: "list" })                  // List all servers
jarvis_server({ action: "search", query: "pdf" })  // Search registry
jarvis_server({ action: "install", name: "..." })  // Install server
jarvis_server({ action: "info", name: "..." })     // Get details
jarvis_server({ action: "uninstall", name: "..." })// Remove server
```

### Profile Management

```javascript
jarvis_profile({ action: "list" })                 // List profiles
jarvis_profile({ action: "suggest" })              // Get recommendations
jarvis_profile({ action: "create", name: "..." })  // Create profile
jarvis_profile({ action: "edit", name: "...", add_servers: "..." })
jarvis_profile({ action: "restart" })              // Hot-reload
```

### Project Operations

```javascript
jarvis_project({ action: "analyze" })              // Analyze codebase
jarvis_project({ action: "devops" })               // Apply DevOps stack
jarvis_project({ action: "diff", staged: true })   // Review changes
```

### Client Configuration

```javascript
jarvis_client({ action: "list" })                  // List clients
jarvis_client({ action: "import", client_name: "..." })
jarvis_client({ action: "edit", client_name: "...", add_profile: "..." })
```

---

## See Also

- [Configuration Strategy](CONFIGURATION_STRATEGY.md) - 3-Layer Profile Stack
- [API Reference](API_REFERENCE.md) - Full tool documentation
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions
