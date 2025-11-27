# 🤖 Jarvis: The Autonomous DevOps Companion

**The Universal API for Agentic Engineering.**

Jarvis is an MCP Server that transforms any LLM—from Gemini Pro to local Llama—into a **Staff DevOps Engineer**. It bridges the gap between *High-Level Intelligence* (your agent) and *Low-Level Infrastructure* (your machine).

For the **AI-Native Engineer**, context switching is dead. Jarvis handles the ops so you can handle the architecture.

---

## ⚡ The Philosophy: 10x Leverage

We built Jarvis for the engineer who demands the output of a department from a single laptop.

*   **Model Agnostic:** Bring your own brains. Jarvis provides the hands. Whether you use Claude, Codex, or Cursor, Jarvis creates a standard interface to reality.
*   **Local Sovereignty:** Your tools, your Docker containers, your git history. Jarvis runs locally and securely manages your environment.
*   **The "Check-Engine" Loop:** Agents are great at writing code but bad at compiling it. Jarvis provides the **Hard Guardrails** (Linters, Pre-commit hooks) that allow Agents to self-correct without human intervention.

---

## 🚀 Capabilities: The DevOps Stack

Jarvis organizes its tools into four distinct "Roles" for your Agent to assume:

### 1. The Architect (`scaffold_project`)
> *"Jarvis, spin up a new Python service with strict security."*

*   **Instant Setup:** Initializes Git, standard ignores, and directory structures.
*   **Guardrails:** Auto-configures `pre-commit` hooks (Ruff, Gitleaks) tailored to the language.
*   **CI/CD:** Generates GitHub Actions for automated AI Code Reviews (`pr-agent`).

### 2. The Strategist (`suggest_profile`)
> *"Context switching to the PokeEdge project."*

*   **Smart Stacking:** Dynamically assembles the perfect toolstack based on your directory.
    *   *Layer 1:* **Environment** (DBs, APIs specific to the project)
    *   *Layer 2:* **Client Adapter** (Tools specific to your LLM interface)
    *   *Layer 3:* **Global** (Memory, Testing)

### 3. The Critic (`fetch_diff_context`)
> *"Review my changes before I commit."*

*   **Local Feedback:** Your Agent can "see" its own diffs and `git status` in real-time.
*   **Self-Correction:** Enables the Agent to catch logic errors or debug prints before they ever hit the repo.

### 4. The Mechanic (MCPM Integration)
*   **Infrastructure:** Installs and manages other MCP servers (Vector DBs, Search Tools).
*   **Health:** Runs `doctor` checks to keep the system green.

---

## 🛠️ Quick Start

### 1. Build the Runtime
```bash
cd Jarvis
go build -o jarvis .
```

### 2. Connect Your Brain
Tell your AI client (Claude/Codex/Gemini) where Jarvis lives.
*Example for **Claude Desktop**:*

```json
{
  "mcpServers": {
    "jarvis": {
      "command": "/absolute/path/to/MCP/Jarvis/jarvis",
      "args": []
    }
  }
}
```

### 3. Bootstrap
Open your Agent and say: **"Bootstrap the system."**

---

## 📚 Documentation
*   [**Configuration Strategy**](docs/CONFIGURATION_STRATEGY.md) - The 3-Layer Stack Architecture.
*   [**Jarvis Development**](Jarvis/README.md) - Source code & Tool definitions.
*   [**MCPM Source**](mcpm_source/README.md) - Our custom, automation-optimized CLI.

## 📜 License
MIT License. Built for the builders.