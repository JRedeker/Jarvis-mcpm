# 🤖 Jarvis: The Autonomous DevOps Companion

Jarvis is an MCP server that lets any LLM—from Gemini Pro to a local Llama—operate as an automated DevOps engineer.
It connects *High-Level Intelligence* (your agent LLM) with *Low-Level Infrastructure* (your machine).

For engineers building around AI, Jarvis reduces context switching. Jarvis handles the ops so you can focus on architecture.

---

## ⚡ The Philosophy: Leverage Local

Jarvis targets engineers who want to stay focused on functionality while still respecting DevOps.

*   **Model Agnostic:** Use any model. Whether using Claude, Codex, or Cursor, Jarvis exposes a consistent interface to tools and infrastructure.
*   **Local Control:** Your tools, your Docker containers, your git history. Jarvis runs locally and securely manages your environment.
*   **The "Check-Engine" Loop:** Agents are good at writing code but weaker at compiling and enforcing standards. Jarvis provides strict checks (linters, pre-commit hooks) so agents can self-correct with minimal human intervention.

---

## 🚀 Capabilities: The DevOps Toolbox For AI Agents

Jarvis organizes its tools into four distinct roles for the agent:

### 1. Architect (`scaffold_project`)
> *"Using Jarvis to spin up a new Python service with strict security."*

*   **Instant Setup:** Initializes Git, standard ignores, and directory structures.
*   **Guardrails:** Auto-configures `pre-commit` hooks (Ruff, Gitleaks) tailored to the language.
*   **CI/CD:** Generates GitHub Actions for automated AI code reviews (`pr-agent`).

### 2. Strategist (`suggest_profile`)
> *"Using Jarvis to switch context to the Tesseract project."*

*   **Smart Stacking:** Dynamically assembles an appropriate tool stack based on the current directory.
    *   *Layer 1:* **Environment** (DBs, APIs specific to the project)
    *   *Layer 2:* **Client Adapter** (Tools specific to the LLM interface)
    *   *Layer 3:* **Global** (Memory, Testing)

### 3. Critic (`fetch_diff_context`)
> *"Using Jarvis to review changes before committing."*

*   **Local Feedback:** The agent can inspect its own diffs and `git status` in real time.
*   **Self-Correction:** Enables the agent to catch logic errors or debug prints before they are committed.

### 4. Mechanic (MCPM Integration)
*   **Infrastructure:** Installs and manages other MCP servers (vector DBs, search tools).
*   **Health:** Runs `doctor` checks to keep the system healthy.

---

## 🛠️ Quick Start

### 1. Build the Runtime
```bash
cd Jarvis
go build -o jarvis .
````

### 2. Connect Your Model

Configure the AI client (Claude/Codex/Gemini) to point to Jarvis.

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

Open the agent and say: **"Bootstrap the system."**

---

## 📚 Documentation

* [**Configuration Strategy**](docs/CONFIGURATION_STRATEGY.md) - The 3-layer stack architecture.
* [**Jarvis Development**](Jarvis/README.md) - Source code and tool definitions.
* [**MCPM Source**](mcpm_source/README.md) - Custom CLI optimized for automation.

## 📜 License

MIT License.
