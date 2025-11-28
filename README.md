# Jarvis

**The local infrastructure layer for AI agents.**

Jarvis is a Model Context Protocol (MCP) server that gives your AI agent "hands." It provides secure, local control over your development environment, enabling agents to scaffold projects, manage Docker infrastructure, and enforce strict engineering standards.

![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![MCP Compliant](https://img.shields.io/badge/MCP-Compliant-blue)

---

## The Problem: Agents Break Things
LLMs are excellent at generating code logic but struggle with the "last mile" of engineering:
*   They hallucinate file paths.
*   They ignore linter rules and break builds.
*   They lack context on your local environment (databases, running services).

## The Solution: Enforced Engineering
Jarvis acts as a **local runtime** that enforces DevOps best practices. Instead of letting an agent wildly edit files, you give it tools to:
1.  **Analyze** the project structure first.
2.  **Scaffold** with strict, pre-configured templates (Git, Pre-commit, GitHub Actions).
3.  **Verify** its own work by checking diffs and linter outputs before committing.

---

## 🚀 Core Capabilities

Jarvis organizes its tools into specialized roles to guide the agent's workflow:

### 🏗️ The Architect (`apply_devops_stack`)
> *"Upgrade this legacy Python repo with modern tooling."*

*   **Smart Analysis:** Detects existing languages and configs via `analyze_project`.
*   **Safe Scaffolding:** Initializes standard tooling (Git, Pre-commit, GitHub Workflows) without destroying custom setups.
*   **Guardrails:** Auto-configures language-specific hooks (Ruff for Python, Gofmt for Go) to catch errors early.

### 🔧 The Mechanic (`restart_infrastructure`)
> *"Restart the vector database and check health."*

*   **Infrastructure Management:** Controls your local Docker stack (Postgres, Qdrant) via a unified management script.
*   **Self-Healing:** Agents can detect service failures and trigger restarts automatically.
*   **Health Checks:** Runs `doctor` diagnostics to ensure the environment is stable.

### 🧐 The Critic (`fetch_diff_context`)
> *"Review my changes before I commit."*

*   **Local Feedback Loop:** Allows the agent to inspect `git status` and `git diff` in real-time.
*   **Self-Correction:** Enables the agent to catch logic errors, debug prints, or secret leaks *before* they enter the commit history.

---

## 🛠️ Getting Started

### 1. Build the Server
Jarvis is a static Go binary. Build it once, run it anywhere.

```bash
git clone https://github.com/JRedeker/Jarvis-mcpm.git
cd Jarvis-mcpm/Jarvis
go build -o jarvis .
```

### 2. Connect Your Agent
Configure your MCP client (Claude Desktop, Cursor, etc.) to use the binary.

**`claude_desktop_config.json`**:
```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/absolute/path/to/Jarvis-mcpm/Jarvis/jarvis",
      "args": []
    }
  }
}
```

### 3. Bootstrap
Open your agent and say:
> **"Bootstrap the system."**

Jarvis will automatically:
*   Install the internal `mcpm` CLI.
*   Launch the local infrastructure (Postgres & Qdrant) via Docker Compose.

---

## ⚡ Advanced Usage

### Infrastructure Management
You can manually control the underlying services using the master script:

```bash
./scripts/manage-mcp.sh [start|stop|restart|logs|status]
```

### 3-Layer Configuration Strategy
Jarvis uses a dynamic "Profile Stacking" system to load tools based on your context:
1.  **Layer 1 (Environment):** Project-specific tools (e.g., `project-pokeedge`).
2.  **Layer 2 (Client):** Adapter tools for your specific LLM (e.g., `client-codex`).
3.  **Layer 3 (Global):** Always-on tools like `memory` and `filesystem`.

[Read the full Configuration Strategy](docs/CONFIGURATION_STRATEGY.md)

---

## 📚 Documentation

*   [**Technical Architecture**](docs/TECHNICAL_ARCHITECTURE.md) - System design and component map.
*   [**Jarvis Development**](Jarvis/README.md) - Internals of the Go server.
*   [**MCPM Source**](mcpm_source/README.md) - The underlying CLI engine.

## 📜 License

MIT License.
