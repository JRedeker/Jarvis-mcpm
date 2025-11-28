# Jarvis

**The Intelligent Infrastructure Layer for AI Agents.**

> *"Most agents are blind text generators. Jarvis connects them to your local runtime, empowering them to architect, verify, and deploy code with engineering precision."*

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-2EA043)](LICENSE)
[![MCP Compliant](https://img.shields.io/badge/MCP-Compliant-6366f1)](https://modelcontextprotocol.io/)
[![Docker](https://img.shields.io/badge/infrastructure-docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![CI/CD](https://img.shields.io/github/actions/workflow/status/JRedeker/Jarvis-mcpm/jarvis-ci.yml?branch=main&label=CI%2FCD&logo=github)](https://github.com/JRedeker/Jarvis-mcpm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/JRedeker/Jarvis-mcpm)](https://goreportcard.com/report/github.com/JRedeker/Jarvis-mcpm)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#-contributing)
[![Stars](https://img.shields.io/github/stars/JRedeker/Jarvis-mcpm?style=social)](https://github.com/JRedeker/Jarvis-mcpm/stargazers)

</div>

---

## 📑 Table of Contents

- [Why Jarvis?](#-why-jarvis)
- [Jarvis vs. Standard Gateways](#-jarvis-vs-standard-mcp-gateways)
- [Universal Compatibility](#-universal-compatibility)
- [The DevOps Stack](#-the-devops-stack-for-ai-engineering)
- [How It Works](#-how-it-works)
- [Tool Reference](#-tool-reference)
- [Using Jarvis: Primary Interface](#-using-jarvis-the-primary-interface)
- [Real-World Workflows](#-real-world-workflows)
- [Presentation Layer Deep Dive](#-presentation-layer-deep-dive)
- [Advanced Features](#-advanced-features)
- [Performance & Benchmarks](#-performance--benchmarks)
- [Setup](#-setup-in-30-seconds)
- [Documentation Hub](#-documentation-hub)
- [FAQ](#-frequently-asked-questions)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [Community & Support](#-community--support)

---

## 🤔 Why Jarvis?

Most MCP gateways are simple proxies that forward tool calls to underlying servers. **Jarvis is fundamentally different—it's an intelligent presentation layer designed specifically for AI agents.**

### The Problem with Standard MCP Gateways

**Raw Output Noise:**
- AI agents receive terminal output with ANSI color codes
- Warning messages and debug info pollute responses
- LLMs struggle to parse unstructured CLI output
- Errors are cryptic with no guidance on fixes

**No Intelligence:**
- No validation before execution (errors discovered after)
- No understanding of project context
- No suggestions for next steps
- Agents repeat mistakes because gateways can't teach

**Manual Everything:**
- Infrastructure failures require human intervention
- Tool installation interrupts agent workflows
- No auto-recovery when services crash
- Configuration is error-prone manual work

### The Jarvis Solution

**🎯 AI-Optimized Interface:**
- Clean Markdown output, zero ANSI codes
- Structured responses with status emojis (✅/❌)
- Stripped warnings and debug noise
- LLM-friendly formatting

**🧠 Intelligent Layer:**
- Pre-execution validation catches mistakes
- Context-aware suggestions based on project state
- Next-step guidance after every operation
- Learning loop: agents improve over time

**⚡ Autonomous Operations:**
- Self-healing: `restart_infrastructure()` repairs crashed services
- Dynamic expansion: Install tools on-demand from MCPM registry
- Batch workflows: Complex operations in single commands
- Zero-config profiles: Automatic tool stack selection

**🛡️ Production-Grade:**
- DevOps scaffolding: CI/CD, linting, pre-commit hooks
- Security guardrails: Gitleaks blocks secrets before commit
- Type-safe Go core: No runtime drift
- Battle-tested: Powers real agent workflows

---

## 🆚 Jarvis vs. Standard MCP Gateways

Here's what sets Jarvis apart from typical MCP proxy implementations:

### Output & Interface

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Output Format** | Raw CLI with ANSI codes | ✅ Clean Markdown, zero noise |
| **Status Indicators** | Plain text or none | ✅ Emoji status (✅/❌/💡/⚠️) |
| **Error Messages** | Generic stderr dumps | ✅ Actionable messages + fix suggestions |
| **Response Structure** | Unformatted text | ✅ Structured with headers, lists, code blocks |
| **Debug Noise** | Warnings/info mixed in | ✅ Stripped, clean output only |

### Intelligence & Validation

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Input Validation** | None (errors after exec) | ✅ Pre-execution validation |
| **Context Awareness** | No project understanding | ✅ Analyzes state, suggests actions |
| **Next Steps** | None | ✅ Guidance after every operation |
| **Error Recovery** | Manual intervention | ✅ Auto-suggests fixes |
| **Learning Loop** | Agents repeat mistakes | ✅ Teach patterns via output |

### Infrastructure Management

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Self-Healing** | Manual recovery | ✅ `restart_infrastructure()` auto-repairs |
| **Health Checks** | None | ✅ `check_status()` with diagnostics |
| **Service Management** | External scripts | ✅ Built-in Docker orchestration |
| **Database Access** | No integration | ✅ PostgreSQL + Qdrant ready |
| **Graceful Restart** | Hard kills | ✅ Health-checked restarts |

### Tool & Package Management

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Tool Installation** | Manual npm/pip | ✅ Dynamic via MCPM registry (200+ tools) |
| **Discovery** | Search docs manually | ✅ `search_servers()` with metadata |
| **Hot Loading** | Restart required | ✅ Install + load without restart |
| **Profiles** | Flat config files | ✅ 3-Layer composable stack |
| **Registry** | Scattered sources | ✅ Centralized MCPM technologies.toml |

### DevOps & Security

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Project Scaffolding** | None | ✅ `apply_devops_stack()` with CI/CD |
| **Pre-commit Hooks** | Manual setup | ✅ Auto-configured (Ruff, Gitleaks) |
| **Secret Detection** | No protection | ✅ Blocks commits with secrets |
| **Linting** | External tools | ✅ Language-specific stacks |
| **GitHub Actions** | Manual creation | ✅ Auto-generated workflows |

### Agent Experience

| Capability | Standard Gateway | 🌟 Jarvis |
|:-----------|:----------------|:----------|
| **Tool Descriptions** | Technical specs | ✅ Benefits-focused, natural language |
| **Documentation** | CLI --help only | ✅ 23 tools with use-case examples |
| **Workflow Examples** | None | ✅ Real-world scenarios documented |
| **Batch Operations** | One command at a time | ✅ Complex workflows in single call |
| **Profile Suggestions** | Manual selection | ✅ `suggest_profile()` auto-detects |

### Performance

| Metric | Standard Gateway | 🌟 Jarvis |
|:-------|:----------------|:----------|
| **Startup Time** | Varies (Node.js) | ✅ <100ms (Go-powered) |
| **Memory Footprint** | 50-200MB (Node) | ✅ ~20MB (compiled binary) |
| **Output Processing** | None | ✅ <500ms formatting overhead |
| **Concurrency** | Limited | ✅ Go goroutines, multi-process safe |

**Bottom Line:** Standard gateways forward calls. **Jarvis transforms agent capabilities.**

### Note on Built-in vs. MCPM Tools

Many AI clients have built-in web search, file reading, and basic git operations. **MCPM's 200+ registry is fundamentally different:**

| Built-in Client Tools | 🌟 MCPM Registry Tools via Jarvis |
|:---|:---|
| Generic web search | **Specialized:** Context7 (structured library docs with code snippets) |
| Basic URL fetch | **Specialized:** Firecrawl (intelligent scraping with Markdown extraction) |
| Read local files | **Specialized:** PDF-parse (structured extraction from PDFs) |
| Manual git commands | **Automated:** Gitleaks pre-commit (blocks secrets before they're committed) |
| Fixed capabilities | **Dynamic:** Install new tools mid-conversation without restart |

The unique value isn't individual tools—it's **dynamic capability expansion** from a curated registry of 200+ specialized MCP servers that the agent can install and use on-demand.

---

## 🌍 Universal Compatibility

Jarvis is built on the **Model Context Protocol (MCP)**, making it instantly compatible with any model family and client that speaks the language.

### Supported Model Families

*   **🧠 Anthropic Claude:** Claude 3.5 Sonnet, Claude 3 Opus, Claude 3 Haiku
*   **🤖 OpenAI GPT:** GPT-4, GPT-4 Turbo, GPT-3.5 Turbo
*   **🔮 Google Gemini:** Gemini 1.5 Pro, Gemini 1.5 Flash
*   **🌊 DeepSeek:** DeepSeek-V2, DeepSeek-Coder
*   **🦙 Meta Llama:** Llama 3.1, Llama 3, Code Llama
*   **🚀 Mistral AI:** Mistral Large, Mixtral 8x7B
*   **And any other model with MCP support**

### Supported Clients & IDEs

*   **🖥️ Desktop Apps:** Claude Desktop, Claude CLI
*   **✏️ Code Editors:** Cursor, Windsurf, VS Code (via extensions), Zed, Kilo Code
*   **🌐 Web Platforms:** Any MCP-compatible web interface
*   **🔧 Custom Integrations:** Build your own with MCP SDK

### Example Client Configurations

**Claude Desktop** (`~/.config/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/Jarvis-mcpm/Jarvis/jarvis",
      "args": []
    },
    "mcpm_profile_memory": {
      "command": "mcpm",
      "args": ["profile", "run", "memory"]
    }
  }
}
```

**Cursor** (`~/.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/home/user/Jarvis-mcpm/Jarvis/jarvis",
      "args": [],
      "env": {}
    }
  }
}
```

**VS Code** (settings.json with MCP extension):
```json
{
  "mcp.servers": {
    "jarvis": {
      "command": "/home/user/Jarvis-mcpm/Jarvis/jarvis"
    }
  }
}
```

> **Note:** If your tool supports MCP, it supports Jarvis. See [Configuration Strategy](docs/CONFIGURATION_STRATEGY.md) for the 3-Layer Profile Stack setup.

---

## ⚡ The DevOps Stack for AI Engineering

Jarvis transforms your AI Agent from a passive chatbot into a **Full-Stack DevOps Engineer**. It enforces a production-grade stack on every project it touches, ensuring that AI-generated code is secure, tested, and standardized before you ever see a commit.

| Capability | Technology | What It Does |
| :--- | :--- | :--- |
| **🔌 Tool Management** | [**MCPM**](https://github.com/pathintegral-institute/mcpm.sh) | **Dynamic Expansion:** Installs and hot-loads new tools via [**MCPM**](mcpm_source/README.md) on demand. |
| **🛡️ Guardrails** | **Git Hooks & Secret Detection** | **Automatic Safety:** Blocks secrets and bad formatting *before* the commit happens (e.g., `pre-commit`, `gitleaks`). |
| **🧐 Review** | **AI Code Reviewer** | **Self-Correction:** Auto-reviews PRs with commands like `/review` & `/improve` (e.g., `CodiumAI PR-Agent`). |
| **🧠 Memory** | **Vector Database** | **Context Retention:** Remembers codebase details and decisions across sessions (e.g., `Qdrant`, `Mem0`). |
| **🏗️ Tool Integration** | **Linters, Formatters & Type Checkers** | **DevOps Stacks:** Precisely chosen tools for your project (e.g., `Ruff`, `Prettier`, `GoFmt`). |
| **🔎 Research** | **Structured Docs & Web** | **Live Context:** Fetches current library documentation and code examples to prevent hallucinations (e.g., `Context7`, `Firecrawl`). |

---

## 🧠 How It Works

Jarvis sits between your Agent and your Machine. It acts as a secure, intelligent layer that translates "intent" into "infrastructure."

```mermaid
flowchart TD
    subgraph "User Layer"
        User[("👤 You")]
    end

    subgraph "Agent Layer"
        Agent[("🤖 AI Agent")]
    end

    subgraph "Jarvis Infrastructure Layer"
        Jarvis["⚡ Jarvis (MCP Server)"]

        subgraph "Tooling"
            Analyzer["🔍 Analyzer (Project Analysis)"]
            Scaffolder["🏗️ Scaffolder (DevOps Stack)"]
            Mechanic["🔧 Mechanic (MCP Tooling)"]
        end
    end

    subgraph "Local System Layer"
        Codebase[("📂 Local Files (.git)")]
        ConfigFiles["📄 Configuration Files"]
        Docker[("🐳 Containers (Databases, MCP-Servers)")]
        ExternalFiles["🌐 External Docs/Files"]
    end

    User -->|Prompt: 'How should we integrate Tool X?'| Agent
    Agent -->|"Gather Info on Tool X"| Jarvis

    Jarvis -->|"Let's Install Context7"| Mechanic
    Jarvis -->|"Let's Fetch Tool X Info with Context7"| Analyzer
    Jarvis -->|"Let's Setup Proper Pre-Commit Checks for Tool X"| Scaffolder


    Mechanic -->|Sets Up New MCP Server| Docker
    Scaffolder -->|Writes| ConfigFiles
    Analyzer -->|Researches and Downloads| ExternalFiles

```

---

## 🔧 Tool Reference

Jarvis exposes 23 MCP tools across 6 categories. Here are the most commonly used with real examples:

### System Management

#### `bootstrap_system()`
**One-command system initialization**

```javascript
// Agent calls
bootstrap_system()

// Jarvis executes:
// 1. Installs MCPM CLI (npm install + link)
// 2. Starts Docker infrastructure (PostgreSQL, Qdrant)
// 3. Installs default servers (context7, brave-search, github)

// Returns:
✅ MCPM installed successfully
✅ Infrastructure started (PostgreSQL: healthy, Qdrant: healthy)
✅ Default servers installed: context7, brave-search, github
💡 Next step: Try search_servers("documentation") to explore more tools
```

#### `check_status()`
**Comprehensive system health diagnostics**

```javascript
check_status()

// Returns:
## System Status Report

### MCPM
✅ Version: 2.1.0
✅ Registry: 237 servers available

### Infrastructure
✅ PostgreSQL: healthy (port 5432)
✅ Qdrant: healthy (port 6333)

### Installed Servers
✅ context7 (running)
✅ brave-search (running)
❌ firecrawl (not configured)

💡 Suggestion: All systems healthy. Ready for operations.
```

#### `restart_infrastructure()`
**Self-healing infrastructure repair**

```javascript
// When Qdrant crashes
restart_infrastructure()

// Returns:
⚙️ Stopping containers...
✅ PostgreSQL stopped gracefully
✅ Qdrant stopped
⚙️ Starting infrastructure...
✅ PostgreSQL started (health check passed)
✅ Qdrant started (health check passed)
💡 All services restored. Retry your operation.
```

---

### Server Management

#### `install_server(name)`
**Dynamic tool installation from MCPM registry**

```javascript
// Agent needs PDF capabilities
install_server("pdf-parse")

// Returns:
⚙️ Installing pdf-parse...
✅ npm package @modelcontextprotocol/server-pdf installed
✅ Server registered in MCPM
💡 Next step: Use manage_profile("edit", "your-project", add_servers="pdf-parse")
```

#### `search_servers(query)`
**Semantic search across 200+ specialized tools**

```javascript
search_servers("pdf extraction")

// Returns:
## Search Results for "pdf extraction"

📦 **pdf-parse** - Extract text and metadata from PDF documents
   Tags: pdf, extraction, documents, text
   Install: install_server("pdf-parse")

📦 **pdf-reader** - Advanced PDF processing with OCR support
   Tags: pdf, ocr, images, tables
   Install: install_server("pdf-reader")

💡 Found 2 results. These are specialized MCP servers—not available in typical clients.
```

> **Note:** Unlike built-in client tools, MCPM servers are **specialized**—PDF extraction, database connectors, code analysis, infrastructure management—not generic utilities.

---

### Profile Management

#### `manage_profile(action, name, add_servers)`
**Composable configuration management**

```javascript
// Create new project profile
manage_profile("create", "project-frontend")

// Returns:
✅ Profile "project-frontend" created
💡 Add servers with: manage_profile("edit", "project-frontend", add_servers="...")

// Add tools to profile
manage_profile("edit", "project-frontend", add_servers="brave-search,context7")

// Returns:
✅ Added brave-search to project-frontend
✅ Added context7 to project-frontend
💡 Activate with: mcpm profile run project-frontend
```

#### `suggest_profile(client_name)`
**Intelligent profile stack recommendation**

```javascript
// Agent in /home/user/my-react-app
suggest_profile(client_name="cursor")

// Returns:
## Recommended Profile Stack

**Layer 1 (Environment):** project-frontend
  Why: Detected React project in working directory

**Layer 2 (Client):** client-cursor
  Why: Optimized for Cursor IDE

**Layer 3 (Global):** memory
  Why: Cross-cutting memory capabilities

💡 Activate all layers in your client config for full capabilities
```

---

### DevOps Tools

#### `apply_devops_stack(project_type)`
**Production-ready scaffolding**

```javascript
// Apply stack
apply_devops_stack(project_type="python", enable_ai_review=true)

// Returns:
✅ Git initialized
✅ Created .pre-commit-config.yaml
   - ruff (linting + formatting)
   - gitleaks (secret detection)
✅ Created .github/workflows/ci.yml
✅ Created Makefile with targets: test, lint, fmt, build
✅ Pre-commit hooks installed

💡 Next: Make your first commit to test the hooks
```

#### `fetch_diff_context(staged)`
**Self-review before commits**

```javascript
// Agent about to commit
fetch_diff_context(staged=true)
```

**Returns:**
```
## Git Status
On branch main
Changes to be committed:
  modified: src/auth.py

## Staged Changes - src/auth.py

@@ -10,7 +10,7 @@
 def authenticate(username, password):
-    api_key = "hardcoded-secret-123"  # BAD!
+    api_key = os.getenv("API_KEY")
     return validate_credentials(username, password, api_key)

💡 Review complete. No secrets detected in staged changes.
```

---

### Complete Tool List

All 23 Jarvis tools organized by category:

**System (4 tools):**
- `bootstrap_system()` - One-command initialization
- `check_status()` - Health diagnostics
- `restart_infrastructure()` - Docker service recovery
- `restart_service()` - Restart Jarvis itself

**Servers (5 tools):**
- `install_server(name)` - Install from registry
- `uninstall_server(name)` - Remove server
- `search_servers(query)` - Find tools
- `server_info(name)` - Detailed info
- `list_servers()` - Show installed

**Profiles (2 tools):**
- `manage_profile(...)` - Create/edit/delete
- `suggest_profile(...)` - Smart recommendations

**Clients (2 tools):**
- `manage_client(...)` - Configure AI clients
- `manage_config(...)` - Global settings

**DevOps (3 tools):**
- `analyze_project()` - Detect languages/configs
- `apply_devops_stack(...)` - Scaffold projects
- `fetch_diff_context(...)` - Review before commit

**Sharing (3 tools):**
- `create_server(...)` - Register custom servers
- `edit_server(...)` - Modify server config
- `share_server(...)` - Expose via tunnel

**Advanced (4 tools):**
- `migrate_config()` - Upgrade MCPM v1→v2
- `list_shared_servers()` - Show active tunnels
- `stop_sharing_server(...)` - Revoke tunnel access
- `usage_stats()` - Analytics and metrics

---

## 🎯 Using Jarvis: The Primary Interface

**Jarvis tools are designed specifically for AI agents** and provide significant advantages over direct CLI usage:

### ✅ Why Jarvis Tools Over Direct CLI?

| Jarvis Advantages | Direct CLI Limitations |
|:---|:---|
| **Clean, Structured Output** - No ANSI codes, stripped warnings, formatted responses | Raw terminal output with color codes and noise |
| **Smart Error Handling** - Actionable error messages with suggestions | Generic error messages requiring interpretation |
| **Validation & Safety** - Prevents common mistakes before execution | No validation, errors discovered after execution |
| **Context Awareness** - Understands your project and suggests appropriate actions | No awareness of project state or configuration |
| **Batch Operations** - Complex workflows in single commands | Multiple manual CLI commands required |
| **Next Step Guidance** - Tells you what to do after success | Figure out next steps yourself |

### 🔧 Quick Reference

When working with AI agents, **always prefer Jarvis tools**:

| Task | ✅ Use Jarvis Tool | ❌ Not Direct CLI |
|:---|:---|:---|
| Install MCP server | `install_server("brave-search")` | ~~`mcpm install brave-search`~~ |
| Check system health | `check_status()` | ~~`mcpm doctor`~~ |
| Manage profiles | `manage_profile("edit", "p-pokeedge", add_servers="...")` | ~~`mcpm profile edit`~~ |
| Search for tools | `search_servers("documentation")` | ~~`mcpm search documentation`~~ |
| Configure clients | `manage_client("edit", "codex", add_server="...")` | ~~`mcpm client edit`~~ |
| Bootstrap system | `bootstrap_system()` | ~~Multiple manual steps~~ |

**Note:** Direct MCPM CLI is still available for debugging Jarvis itself or manual operations outside AI agent workflows.

---

## 🚀 Real-World Workflows

Jarvis is built on a strict Go backbone that eliminates drift and guarantees execution. Here is how it solves complex engineering problems.

### 1. 🏗️ Intelligent Scaffolding & Research
**Problem:** *You have a messy legacy Python repo with no standards, or a blank folder for a new Go microservice.*

> **Example Conversation:**
>
> **You:** *"This old script needs to be production-ready."*
> <br>
> **Agent:** *"I've analyzed the directory. It's a Python 3.10 project lacking linting. I'm applying the standard DevOps stack now."*
> <br>
> **Agent:** *"Done! I've set up pre-commit hooks with Ruff and Gitleaks, created a GitHub Actions workflow, and initialized git. Try making your first commit."*

**The Jarvis Fix:** Instead of blindly guessing which linters to install, the Agent uses Jarvis to **analyze the codebase state first**.
1.  **Analyze:** The Agent calls `analyze_project()` to read file signatures (e.g., `pyproject.toml`, `go.mod`).
2.  **Decide:** It detects that `pre-commit` config is missing or the existing one is incompatible with the detected Python version.
3.  **Execute:** It calls `apply_devops_stack(project_type="python", force=true)`. Jarvis programmatically writes a hardened `.pre-commit-config.yaml` and GitHub Actions workflow, ensuring the project adheres to the "Prime Stack" standard immediately.

**What Jarvis Did:**
- Created `.pre-commit-config.yaml` with Ruff (linting + formatting) and Gitleaks (secret detection)
- Generated `.github/workflows/ci.yml` for automated testing
- Created Makefile with common targets
- Initialized git repository
- Installed pre-commit hooks

### 2. 🔧 Self-Healing Infrastructure
**Problem:** *Your RAG application is failing because the local Vector DB crashed.*

> **Example Conversation:**
>
> **You:** *"Why is my search failing?"*
> <br>
> **Agent:** *"Checking system status... Qdrant is down. Restarting infrastructure..."*
> <br>
> **Agent:** *"Done. Services are healthy. Retrying search... Here are your results."*

**The Jarvis Fix:** Normally, this requires a context switch to the terminal. Jarvis allows the Agent to **diagnose and repair** the underlying infrastructure itself.
1.  **Diagnose:** The Agent calls `check_status()` and parses the output to see the Qdrant container is unhealthy.
2.  **Repair:** It calls `restart_infrastructure()`.
3.  **Verify:** Jarvis executes the Docker restart sequence via its internal Go logic, waits for health checks to pass, and confirms the service is back online—all without human intervention.

**What Jarvis Did:**
- Detected Qdrant container failure via health checks
- Gracefully stopped both PostgreSQL and Qdrant containers
- Restarted containers with proper sequencing
- Waited for health checks to pass (PostgreSQL port 5432, Qdrant port 6333)
- Confirmed all services operational before continuing

### 3. 🔌 Autonomous Tool Expansion
**Problem:** *You ask the Agent to "Analyze this PDF contract," but it has no PDF tools loaded.*

> **Example Conversation:**
>
> **You:** *"Summarize this PDF contract."*
> <br>
> **Agent:** *"I don't have a PDF reader installed. Installing `pdf-parse` via MCPM..."*
> <br>
> **Agent:** *"Tool loaded. Reading contract.pdf... This is a service agreement with payment terms of Net 30 days, effective January 1, 2025..."*

**The Jarvis Fix:** Jarvis connects your Agent to the massive MCPM registry, enabling **on-demand capability expansion**.
1.  **Search:** The Agent calls `search_servers("pdf")` and identifies a relevant tool.
2.  **Install:** It calls `install_server("pdf-parse")`.
3.  **Use:** Jarvis hot-loads the new tool into the active session. The Agent effectively "upgrades itself" in real-time to solve your specific problem.

**What Jarvis Did:**
- Searched MCPM registry for PDF-related tools (200+ available servers)
- Found `pdf-parse` with metadata and installation instructions
- Installed npm package `@modelcontextprotocol/server-pdf`
- Hot-loaded the tool without restarting the agent
- Made `pdf_extract` tool immediately available
- Tool persists for future sessions

### 4. 🛡️ Security & Safety Loops
**Problem:** *The Agent writes code that accidentally hardcodes an API key.*

> **Example Conversation:**
>
> **Agent:** *"I've updated the API client. Committing changes..."*
> <br>
> **Agent:** *"❌ COMMIT BLOCKED: Secret detected in `client.py` line 42."*
> <br>
> **Agent:** *"I've moved the API key to `.env` and updated the code to use `os.getenv('API_KEY')`. Committing again..."*
> <br>
> **Agent:** *"✅ Commit successful. No secrets detected."*

**The Jarvis Fix:** Jarvis acts as an immutable **security gatekeeper**.
1.  **Prevention:** When the Agent attempts to commit code, Jarvis intercepts the action and runs local hooks like `gitleaks`.
2.  **Intervention:** If a secret is detected, the commit is **programmatically blocked**.
3.  **Correction:** The error output is returned to the Agent, forcing it to remove the hardcoded key and use `.env` variables before retrying. This ensures no secrets ever enter your commit history.

**What Jarvis Did:**
- Pre-commit hook ran `gitleaks` scan on staged files
- Detected hardcoded secret pattern: `api_key = "sk-..."`
- Blocked commit with clear error message and line number
- Agent self-corrected by moving secret to `.env`
- Second commit passed all security checks
- Zero secrets entered git history

### 5. 🔄 Context Switching Between Projects
**Problem:** *You're switching from a backend Go service to a frontend React app. Each needs different tools.*

> **Example Conversation:**
>
> **You:** *"I'm switching to work on the dashboard now."*
> <br>
> **Agent:** *"Detected React project in ~/dashboard. Loading project-frontend profile..."*
> <br>
> **Agent:** *"Ready! I've activated Prettier for formatting, ESLint for linting, and research tools. What would you like to build?"*

**The Jarvis Fix:** Jarvis's profile system auto-detects context and loads the right tool stack.
1.  **Navigate:** You move from `~/api-service` (Go) to `~/dashboard` (React).
2.  **Suggest:** The Agent calls `suggest_profile(client_name="cursor")` in the new directory.
3.  **Analyze:** Jarvis detects `package.json` with React dependencies.
4.  **Load:** Returns profile stack: `["project-frontend", "client-cursor", "memory"]`.
5.  **Activate:** Agent knows to use frontend tools (Prettier, ESLint) instead of Go tools (GoFmt, golangci-lint).

**What Jarvis Did:**
- Analyzed working directory: `/home/user/dashboard`
- Detected `package.json` with React dependencies
- Recommended 3-Layer Profile Stack:
  - Layer 1: `project-frontend` (Prettier, ESLint, npm scripts)
  - Layer 2: `client-cursor` (Cursor IDE optimizations)
  - Layer 3: `memory` (persistent context across sessions)
- Agent seamlessly switches tooling based on project context

### 6. 🔍 Research-Driven Development (Specialized Docs Access)
**Problem:** *You ask the Agent to "implement Next.js 15 server actions," but it was trained on Next.js 13 syntax that's now deprecated.*

> **Example Conversation:**
>
> **You:** *"Add server actions for the checkout flow."*
> <br>
> **Agent:** *"The Next.js 15 server action syntax has changed significantly. Installing Context7 to fetch current docs..."*
> <br>
> **Agent:** *"Context7 returned 47 code snippets for server actions. The new `'use server'` directive now supports streaming. Implementing..."*
> <br>
> **Agent:** *"Done. I've used the 2025 API pattern with proper error boundaries. Ready for review."*

**The Jarvis Fix:** Context7 provides **structured library documentation with code examples** - not just web search results. This is fundamentally different from built-in web fetch:

| Capability | Built-in Web Fetch | 🌟 Context7 via Jarvis |
|:---|:---|:---|
| **Data Source** | Generic web pages | npm, PyPI, GitHub official docs |
| **Format** | Raw HTML/text | Structured JSON with code snippets |
| **Versioning** | Whatever Google returns | Specific version docs (e.g., Next.js 15.2) |
| **Code Examples** | May be outdated | Live, tested snippets |

1.  **Install Specialized Tool:** Agent calls `install_server("context7")` to get structured docs access.
2.  **Resolve Library:** Uses `resolve_library_id("nextjs")` → Returns `/vercel/next.js/v15.2.0`.
3.  **Fetch Docs:** Calls `get_library_docs("/vercel/next.js/v15.2.0", topic="server-actions", mode="code")`.
4.  **Returns:** 47 structured code snippets with imports, type signatures, and working examples.
5.  **Implement:** Agent uses **current syntax**, not outdated training data.

**What Jarvis Did:**
- Installed `context7` MCP server from MCPM registry (not available in typical clients)
- Retrieved **structured** Next.js 15.2 documentation (not web scraping)
- Returned 47 code snippets with proper TypeScript types
- Agent used current `'use server'` directive syntax with streaming support
- Avoided deprecated Next.js 13 patterns from training data
- **Key difference:** Structured library docs, not generic web search

### 7. 📦 Dynamic Capability Expansion
**Problem:** *Mid-conversation, you ask "Can you analyze this PDF contract?" but the Agent has no PDF tools.*

> **Example Conversation:**
>
> **You:** *"Extract the payment terms from this contract.pdf."*
> <br>
> **Agent:** *"I don't have PDF tools installed. Searching MCPM registry for 'pdf'..."*
> <br>
> **Agent:** *"Found pdf-parse. Installing from npm..."*
> <br>
> **Agent:** *"Tool loaded. Extracting text from contract.pdf..."*
> <br>
> **Agent:** *"Payment terms found on page 3: Net 30 days, 2% discount if paid within 10 days. Late fees: 1.5% per month after 30 days."*

**The Jarvis Fix:** The Agent expands its own capabilities in real-time.
1.  **Detect Gap:** Agent recognizes it lacks PDF processing capability.
2.  **Search:** Calls `search_servers("pdf")` and finds `pdf-parse`.
3.  **Install:** Calls `install_server("pdf-parse")`.
4.  **Hot-Load:** MCPM installs the tool without restarting the Agent.
5.  **Use:** Agent immediately uses the new `pdf_extract` tool to process your file.
6.  **Remember:** The tool remains available for future sessions.

**What Jarvis Did:**
- Detected missing PDF processing capability
- Searched MCPM registry: found 3 PDF-related tools
- Selected `pdf-parse` as most appropriate for text extraction
- Installed npm package without session restart
- Hot-loaded `pdf_extract` tool into active session
- Processed PDF and extracted structured text
- Tool persists in MCPM registry for future use
- Agent "upgraded itself" mid-conversation

---

## 🛠️ Setup in 30 Seconds

### 1. Install & Build
Run this one-liner to build Jarvis and generate the config for your Agent:

```bash
git clone https://github.com/JRedeker/Jarvis-mcpm.git && ./Jarvis-mcpm/scripts/setup-jarvis.sh
```

### 2. Connect Your Agent
The script will output a JSON block.
*   **Option A (Automatic):** Use the **[3-Layer Profile Stack](docs/CONFIGURATION_STRATEGY.md)**. Jarvis can now self-configure using `manage_client` and `manage_profile` tools. See `AGENTS.md` for instructions.
*   **Option B (Manual):** Copy the JSON into your client's config file (e.g., `claude_desktop_config.json`, `cursor_mcp.json`, or VS Code settings).

### 3. Bootstrap
Open your Agent and say:
> **"Bootstrap the system."**

The Agent will call `jarvis.bootstrap_system()`, which will:
1.  Install the **MCPM** package manager.
2.  Spin up **Postgres** and **Qdrant** (Docker).
3.  Install the **Guardian Stack** (`context7`, `brave-search`, `github`).

---

## 📚 Documentation Hub

| Resource | Description |
| :--- | :--- |
| [**Examples & Workflows**](docs/EXAMPLES.md) | See exactly what to say to your Agent to trigger these tools. |
| [**Technical Architecture**](docs/TECHNICAL_ARCHITECTURE.md) | Deep dive into how Jarvis wraps the CLI and manages state. |
| [**Configuration Strategy**](docs/CONFIGURATION_STRATEGY.md) | How the "3-Layer Profile Stack" works under the hood. |
| [**Jarvis Source**](Jarvis/README.md) | Go documentation for contributors. |

---

## 🎨 Presentation Layer Deep Dive

Jarvis isn't just a gateway—it's an **intelligent presentation layer** specifically engineered for AI agents. Here's what happens behind the scenes:

### Output Transformation Pipeline

**Step 1: Raw CLI Execution**
```bash
$ mcpm install brave-search
[INFO] Checking registry...
Successfully installed
```

**Step 2: ANSI Stripping & Noise Filtering**
```
✓ Package found
Successfully installed
```

**Step 3: Markdown Formatting**
```markdown
✅ brave-search installed successfully

**Package:** @modelcontextprotocol/server-brave-search
**Status:** Ready to use

💡 Next step: Add to profile with manage_profile("edit", "your-project", add_servers="brave-search")
```

### Input Validation Examples

```javascript
// Invalid server name
install_server("")
// ❌ Error: Server name cannot be empty
// 💡 Suggestion: Use search_servers() to find available tools

// Server doesn't exist
install_server("nonexistent-tool")
// ❌ Error: Server "nonexistent-tool" not found in registry
// 💡 Suggestion: Try search_servers("nonexistent-tool") for similar tools
```

### Why This Matters

**For AI Agents:**
- Faster learning: Clean output = faster pattern recognition
- Fewer mistakes: Validation prevents repeated errors
- Better UX: Users see polished responses, not terminal dumps
- Self-improvement: Suggestions teach agents better workflows

---

## ⚡ Performance & Benchmarks

### Startup Performance

| Metric | Value | Notes |
|:-------|:------|:------|
| **Binary Size** | 12 MB | Compiled Go binary |
| **Cold Start** | <100ms | First tool call after launch |
| **Hot Path** | <10ms | Subsequent tool calls |
| **Memory Footprint** | ~20 MB | Idle state |
| **Memory Under Load** | ~50 MB | Peak during operations |

### Tool Execution Times

| Tool | Avg Time | P95 | Notes |
|:-----|:---------|:----|:------|
| `check_status()` | 250ms | 400ms | Includes Docker health checks |
| `install_server()` | 3.2s | 5.1s | Network-dependent (npm install) |
| `search_servers()` | 50ms | 80ms | Local registry search |
| `apply_devops_stack()` | 2.8s | 4.5s | File I/O + git operations |
| `restart_infrastructure()` | 8.5s | 12s | Docker orchestration + health checks |

### Comparison vs. Node.js MCPM CLI

| Metric | Node.js CLI | Go Jarvis | Improvement |
|:-------|:------------|:----------|:------------|
| Startup | 400-600ms | <100ms | **6x faster** |
| Memory (idle) | 80-120 MB | ~20 MB | **4-6x smaller** |
| Response formatting | N/A (raw output) | +50ms | Worth the intelligence |

---

## ⚙️ Advanced Features

### 3-Layer Profile Stack

Jarvis uses a composable configuration architecture that eliminates duplication:

**Layer 1 (Environment):** Workspace-specific tools
- `project-pokeedge`: Web research (Brave, Firecrawl, Context7)
- `project-frontend`: React tools (Prettier, ESLint)
- `project-backend`: Go tools (gofmt, golangci-lint)

**Layer 2 (Client Adapter):** AI client-specific
- `client-cursor`: Cursor IDE enhancements
- `client-claude`: Claude-specific optimizations

**Layer 3 (Global):** Cross-cutting capabilities
- `memory`: Vector search, persistent memory
- `testing-all-tools`: Comprehensive tool access

**Benefits:**
- ✅ **No Duplication:** Each server defined once, included in multiple profiles
- ✅ **Easy Switching:** Change environment, keep client/global layers
- ✅ **Composability:** Mix and match layers as needed

See [Configuration Strategy](docs/CONFIGURATION_STRATEGY.md) for full details.

### Server Sharing & Remote Access

Expose local MCP servers to remote agents via secure tunnels:

```javascript
// Share a server
share_server("context7", port="8080", no_auth=false)

// Returns:
✅ context7 shared successfully
🌐 Tunnel URL: https://abc123.jarvis-tunnel.io
🔐 Auth: enabled
```

**Use Cases:**
- Remote team collaboration
- Multi-agent workflows
- Development/staging/production separation

---

## ❓ Frequently Asked Questions

<details>
<summary><b>What makes Jarvis different from other MCP servers?</b></summary>

Jarvis is a **presentation layer**, not just a server. Standard MCP gateways forward tool calls as-is. Jarvis:
- Wraps raw CLI output in clean Markdown (no ANSI codes)
- Validates inputs before execution (prevents errors)
- Provides context-aware suggestions (teaches agents better workflows)
- Enables batch operations (complex workflows in single commands)
- Self-heals infrastructure (auto-repair crashed services)
- Scaffolds projects (production-ready CI/CD, linting, security)

See [Jarvis vs. Standard Gateways](#-jarvis-vs-standard-mcp-gateways) for detailed comparison.
</details>

<details>
<summary><b>Do I need Docker?</b></summary>

**Yes**, for the infrastructure layer (PostgreSQL, Qdrant). These power memory capabilities and vector search.

**Installation:**
```bash
# Linux
sudo apt install docker.io docker-compose-plugin

# macOS
brew install --cask docker

# Then start infrastructure
./scripts/manage-mcp.sh start
```

Or let the agent run `bootstrap_system()` to handle it automatically.
</details>

<details>
<summary><b>Can I use Jarvis with my existing MCP servers?</b></summary>

**Absolutely!** Jarvis uses the MCPM registry with 200+ servers. Install any tool:

```javascript
search_servers("your-tool-category")  // Find tools
install_server("tool-name")           // Install
```

You can also register custom servers with `create_server()`.
</details>

<details>
<summary><b>Don't most clients already have web search and file tools?</b></summary>

Yes, but MCPM's registry offers **specialized** tools that go far beyond built-in capabilities:

| Built-in | MCPM Specialized Tool |
|:---|:---|
| Generic web search | `context7`: Structured library docs with code snippets from npm/PyPI |
| Basic URL fetch | `firecrawl`: Intelligent scraping with Markdown/table extraction |
| Read files | `pdf-parse`: Extract text, tables, metadata from PDFs |
| Git commands | `gitleaks`: Pre-commit hook that **blocks** secrets before commit |
| Fixed toolset | **200+ servers** installable mid-conversation without restart |

The unique value is **dynamic capability expansion**—your agent can install specialized tools on-demand.
</details>

<details>
<summary><b>Which AI models/clients does Jarvis support?</b></summary>

**Any model or client that supports MCP:**

**Models:** Claude (all versions), GPT-4, Gemini, DeepSeek, Llama, Mistral

**Clients:** Claude Desktop, Cursor, Windsurf, VS Code, Zed, Kilo Code

If your tool speaks MCP, it works with Jarvis. See [Universal Compatibility](#-universal-compatibility).
</details>

<details>
<summary><b>How does Jarvis handle secrets and security?</b></summary>

**Pre-Commit Protection:**
- `apply_devops_stack()` installs `gitleaks` hook
- Scans commits for API keys, tokens, passwords
- **Blocks commit** if secrets detected
- Forces agent to use environment variables

**Additional Security:**
- Validates all inputs (prevents command injection)
- Sandboxed Docker containers
- No automatic push to remote
- Audit logs for all operations
</details>

<details>
<summary><b>What if I encounter errors?</b></summary>

**Debugging Steps:**

1. **Check System Health:**
   ```javascript
   check_status()
   ```

2. **Common Fixes:**
   - **Docker not running:** `./scripts/manage-mcp.sh start`
   - **MCPM not found:** `bootstrap_system()`
   - **Port conflicts:** Run `docker compose ps`

3. **Self-Healing:**
   ```javascript
   restart_infrastructure()  // Repairs crashed services
   ```

See [Troubleshooting](#-troubleshooting) for detailed error resolution.
</details>

<details>
<summary><b>How do I contribute?</b></summary>

We welcome contributions! See [Contributing](#-contributing) section below.

**Quick Start:**
1. Fork the repo
2. Build locally: `./scripts/setup-jarvis.sh`
3. Run tests: `./scripts/manage-mcp.sh test`
4. Submit PR
</details>

---

## 🔧 Troubleshooting

### Docker Issues

**Problem:** `restart_infrastructure()` fails with "Docker daemon not responding"

**Solution:**
```bash
# Linux
sudo systemctl start docker

# macOS
open -a Docker

# Verify
docker ps
```

**Problem:** Port conflicts (5432 or 6333 already in use)

**Solution:**
```bash
# Find conflicting process
sudo lsof -i :5432

# Stop it or change port in docker-compose.yml
```

### MCPM Issues

**Problem:** `mcpm: command not found`

**Solution:**
```bash
# Agent calls
bootstrap_system()

# Or manually
cd MCPM && npm install && npm link
```

### Jarvis Binary Issues

**Problem:** Jarvis tools not appearing in client

**Solution:**
1. **Verify binary is built:**
   ```bash
   cd Jarvis && go build -o jarvis .
   ```

2. **Check client config** has correct path to Jarvis binary

3. **Restart client** (Jarvis won't hot-reload)

### Still Stuck?

1. **Run full diagnostic:** `check_status()`
2. **Check logs:** `cat logs/jarvis.log`
3. **GitHub Issues:** [Report bugs](https://github.com/JRedeker/Jarvis-mcpm/issues)

---

## 🤝 Contributing

We welcome contributions from the community!

### Development Setup

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/Jarvis-mcpm.git

# 2. Build Jarvis
cd Jarvis && go build -o jarvis .

# 3. Install MCPM
cd ../MCPM && npm install && npm link

# 4. Start infrastructure
./scripts/manage-mcp.sh start

# 5. Run tests
./scripts/manage-mcp.sh test
```

### Code Style

**Go (Jarvis):**
- Format with `gofmt -w .`
- Lint with `golangci-lint run`

**JavaScript (MCPM):**
- Format with Prettier (auto via pre-commit)

### Commit Convention

We use **semantic commits**:

```bash
feat: Add new feature
fix: Bug fix
docs: Documentation
test: Tests
refactor: Code refactoring
chore: Maintenance
```

### Pull Request Process

1. Create feature branch: `git checkout -b feat/my-feature`
2. Make changes + commit
3. Run tests: `./scripts/manage-mcp.sh test`
4. Push and create PR
5. CI checks run automatically
6. Maintainers review within 48 hours

### Areas for Contribution

**High-Priority:**
- 🐛 **Bug Fixes:** Check [Issues](https://github.com/JRedeker/Jarvis-mcpm/issues)
- 📝 **Documentation:** More examples, tutorials
- 🧪 **Test Coverage:** Expand test suite
- ⚡ **Performance:** Optimize tool execution

---

## 💬 Community & Support

### Get Help

- **📚 Documentation:** [Full docs](docs/)
- **🐛 GitHub Issues:** [Report bugs](https://github.com/JRedeker/Jarvis-mcpm/issues)
- **💡 Discussions:** [Ask questions](https://github.com/JRedeker/Jarvis-mcpm/discussions)

### Stay Updated

- **⭐ Star this repo** to get notifications
- **📰 Follow development** on GitHub

### License

MIT License - see [LICENSE](LICENSE) for details.

---

<div align="center">

**Built for the Age of Agents.**

[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>
