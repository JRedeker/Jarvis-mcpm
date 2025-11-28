# Jarvis

**The DevOps Hand of the AI Agent.**

Jarvis is an MCP server that allows an AI Agent (like Claude, Codex, or Gemini) to securely control your local development environment.

It is not a CLI tool for humans. It is an API for your Agent.

![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![MCP Compliant](https://img.shields.io/badge/MCP-Compliant-blue)

---

## The Concept: Agent-Driven Engineering

In a traditional workflow, you are the engineer and the AI is just a text generator.
With Jarvis, you become the **Architect**, and the AI Agent becomes the **Engineer**.

**Jarvis provides the tools that the Agent needs to do the job.**

```mermaid
graph LR
    User[You (The Architect)] -->|Prompt: 'Setup this project'| Agent[AI Agent]
    Agent -->|Tool Call: 'apply_devops_stack'| Jarvis[Jarvis (MCP Server)]
    Jarvis -->|Executes| Infrastructure[Local System (Git, Docker, Files)]
```

---

## 🤖 How Your Agent Uses Jarvis

You don't run Jarvis commands. You just talk to your Agent. Jarvis works silently in the background, empowering the Agent to perform complex DevOps tasks.

### 1. As an Architect
**You say:** *"Analyze this codebase and set up the DevOps stack."*
**The Agent calls:** `jarvis.analyze_project()` then `jarvis.apply_devops_stack()`
**Jarvis executes:**
*   Scans your files to detect languages (Python, Go, Node).
*   Initializes Git and strictly configured `.gitignore`.
*   Installs `pre-commit` hooks (Ruff, Gitleaks) to prevent bad code from being saved.
*   Sets up GitHub Actions for automated AI code review.

### 2. As a Mechanic
**You say:** *"The vector database is acting up. Restart it."*
**The Agent calls:** `jarvis.restart_infrastructure()`
**Jarvis executes:**
*   Gracefully stops the local Docker containers (Postgres, Qdrant).
*   Reboots the services and verifies health checks.
*   Returns the logs to the Agent so it can confirm the fix.

### 3. As a Critic
**You say:** *"Review my changes before I commit."*
**The Agent calls:** `jarvis.fetch_diff_context()`
**Jarvis executes:**
*   Retrieves the real-time `git diff` and `git status`.
*   Formats it for the Agent to analyze.
*   The Agent then gives you a critique based on the actual state of your disk.

---

## 🛠️ Setup

### 1. Build the Server
Jarvis is a static Go binary that runs on your machine.

```bash
git clone https://github.com/JRedeker/Jarvis-mcpm.git
cd Jarvis-mcpm/Jarvis
go build -o jarvis .
```

### 2. Give Your Agent Access
Configure your AI client (Claude Desktop, Cursor, etc.) to see Jarvis.

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

### 3. Activate
Open your Agent and say:
> **"Bootstrap the system."**

The Agent will call `jarvis.bootstrap_system()`, which will automatically install dependencies and spin up the local Docker infrastructure. You are now ready to code.

---

## 📚 Documentation

*   [**Technical Architecture**](docs/TECHNICAL_ARCHITECTURE.md) - How Jarvis, MCPM, and the Agent interact.
*   [**Configuration Strategy**](docs/CONFIGURATION_STRATEGY.md) - How the Agent decides which tools to load.
*   [**Jarvis Development**](Jarvis/README.md) - Source code documentation.

## 📜 License

MIT License.
