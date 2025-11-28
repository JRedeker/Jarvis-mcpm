# Jarvis

**The Intelligent Infrastructure Layer for AI Agents.**

> *"Most agents are blind text generators. Jarvis connects them to your local runtime, empowering them to architect, verify, and deploy code with engineering precision."*

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-2EA043)](LICENSE)
[![MCP Compliant](https://img.shields.io/badge/MCP-Compliant-6366f1)](https://modelcontextprotocol.io/)
[![Docker](https://img.shields.io/badge/infrastructure-docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

</div>

---

## 🌍 Universal Compatibility

Jarvis is built on the **Model Context Protocol (MCP)**, making it instantly compatible with any model family and client that speaks the language.

*   **🧠 Model Families:** Anthropic Claude, OpenAI GPT, Google Gemini, DeepSeek, Meta Llama.
*   **🖥️ Clients:** Claude Desktop, Cursor, Windsurf, VS Code (via extensions), Zed, Kilo Code, and more.

If your tool supports MCP, it supports Jarvis.

---

## ⚡ The DevOps Stack for AI Engineering

Jarvis transforms your AI Agent from a passive chat bot into a **Full-Stack DevOps Engineer**. It enforces a production-grade stack on every project it touches, ensuring that AI-generated code is secure, tested, and standardized before you ever see a commit.

| Capability | Technology | What It Does |
| :--- | :--- | :--- |
| **🔌 Tool Management** | [**MCPM**](https://github.com/pathintegral-institute/mcpm.sh) | **Dynamic Expansion:** Installs and hot-loads new tools via [**MCPM**](mcpm_source/README.md) on demand. |
| **🛡️ Guardrails** | **Git Hooks & Secret Detection** | **Automatic Safety:** Blocks secrets and bad formatting *before* the commit happens (e.g., `pre-commit`, `gitleaks`). |
| **🧐 Review** | **AI Code Reviewer** | **Self-Correction:** Auto-reviews PRs with commands like `/review` & `/improve` (e.g., `CodiumAI PR-Agent`). |
| **🧠 Memory** | **Vector Database** | **Context Retention:** Remembers codebase details and decisions across sessions (e.g., `Qdrant`, `Mem0`). |
| **🏗️ Tool Integration** | **Linters, Formatters & Type Checkers** | **DevOps Stacks:** Precisely chosen tools for your project (e.g., `Ruff`, `Prettier`, `GoFmt`). |
| **🔎 Research** | **Web Search Index** | **Research:** Fetches live docs via the web to prevent hallucinations (e.g., `Brave Search`). |

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
            Static Analysis ["🔍 Tools for project analysis"]
            DevOps Pipelines ["🏗️ Intelligent DevOps Stack"]
            MCP Servers["🔧 Dynamic MCP Tooling"]
        end
    end

    subgraph "Local System Layer"
        Codebase[("📂 Local Files (.git, configs)")]
        Docker[("🐳 Containers (Databases, MCP-Servers)")]
    end

    User -->|Prompt: 'How should we integrate Tool X?'| Agent
    Agent -->|"Gather Info on Tool X"| Jarvis
    Jarvis -->|"Let's Install Context7"| Mechanic
    Jarvis -->|"Let's Fetch Tool X Info with Context7"| Analyzer
    Jarvis -->|"Let's Setup Proper Pre-Commit Checks for Tool X"| Scaffolder


    Mechanic -->|Sets Up New MCP Server| Docker
    Scaffolder -->|Writes| Configuration Files
    Analyzer -->|Researches and Downloads| Files

```

---

## 🚀 Real-World Workflows

Jarvis is built on a strict Go backbone that eliminates drift and guarantees execution. Here is how it solves complex engineering problems.

### 1. 🏗️ Intelligent Scaffolding & Research
**Problem:** *You have a messy legacy Python repo with no standards, or a blank folder for a new Go microservice.*

**The Jarvis Fix:** Instead of blindly guessing which linters to install, the Agent uses Jarvis to **analyze the codebase state first**.
1.  **Analyze:** The Agent calls `analyze_project()` to read file signatures (e.g., `pyproject.toml`, `go.mod`).
2.  **Decide:** It detects that `pre-commit` config is missing or the existing one is incompatible with the detected Python version.
3.  **Execute:** It calls `apply_devops_stack(project_type="python", force=true)`. Jarvis programmatically writes a hardened `.pre-commit-config.yaml` and GitHub Actions workflow, ensuring the project adheres to the "Prime Stack" standard immediately.

> **Scenario Example:**
>
> **You:** *"This old script needs to be production-ready."*
> <br>
> **Agent:** *"I've analyzed the directory. It's a Python 3.10 project lacking linting. I'm applying the standard DevOps stack now."*
> <br>
> **Jarvis:** *Writes `.pre-commit-config.yaml` with Ruff and Gitleaks.*

### 2. 🔧 Self-Healing Infrastructure
**Problem:** *Your RAG application is failing because the local Vector DB crashed.*

**The Jarvis Fix:** Normally, this requires a context switch to the terminal. Jarvis allows the Agent to **diagnose and repair** the underlying infrastructure itself.
1.  **Diagnose:** The Agent calls `check_status()` and parses the output to see the Qdrant container is unhealthy.
2.  **Repair:** It calls `restart_infrastructure()`.
3.  **Verify:** Jarvis executes the Docker restart sequence via its internal Go logic, waits for health checks to pass, and confirms the service is back online—all without human intervention.

> **Scenario Example:**
>
> **You:** *"Why is my search failing?"*
> <br>
> **Agent:** *"Checking system status... Qdrant is down. Restarting infrastructure... Done. Services are healthy. Retrying search."*

### 3. 🔌 Autonomous Tool Expansion
**Problem:** *You ask the Agent to "Analyze this PDF contract," but it has no PDF tools loaded.*

**The Jarvis Fix:** Jarvis connects your Agent to the massive MCPM registry, enabling **on-demand capability expansion**.
1.  **Search:** The Agent calls `search_servers("pdf")` and identifies a relevant tool.
2.  **Install:** It calls `install_server("pdf-parse")`.
3.  **Use:** Jarvis hot-loads the new tool into the active session. The Agent effectively "upgrades itself" in real-time to solve your specific problem.

> **Scenario Example:**
>
> **You:** *"Summarize this PDF."*
> <br>
> **Agent:** *"I don't have a PDF reader installed. Installing `pdf-parse` via MCPM... Tool loaded. Reading PDF now..."*

### 4. 🛡️ Security & Safety Loops
**Problem:** *The Agent writes code that accidentally hardcodes an API key.*

**The Jarvis Fix:** Jarvis acts as an immutable **security gatekeeper**.
1.  **Prevention:** When the Agent attempts to commit code, Jarvis intercepts the action and runs local hooks like `gitleaks`.
2.  **Intervention:** If a secret is detected, the commit is **programmatically blocked**.
3.  **Correction:** The error output is returned to the Agent, forcing it to remove the hardcoded key and use `.env` variables before retrying. This ensures no secrets ever enter your commit history.

> **Scenario Example:**
>
> **Agent:** *"Committing fix for API client..."*
> <br>
> **Jarvis:** *❌ COMMIT BLOCKED: Secret detected in `client.py`.*
> <br>
> **Agent:** *"Apologies. I've moved the API key to `.env` and am retrying the commit."*

---

## 🛠️ Setup in 30 Seconds

### 1. Install & Build
Run this one-liner to build Jarvis and generate the config for your Agent:

```bash
git clone https://github.com/JRedeker/Jarvis-mcpm.git && ./Jarvis-mcpm/scripts/setup-jarvis.sh
```

### 2. Connect Your Agent
The script will output a JSON block.
*   **Option A (Manual):** Copy the JSON into your client's config file (e.g., `claude_desktop_config.json`, `cursor_mcp.json`, or VS Code settings).
*   **Option B (Agentic):** Just paste the output to your Agent and say:
    > *"Configure yourself to use this MCP server."*

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

<div align="center">

**Built for the Age of Agents.**

[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>
