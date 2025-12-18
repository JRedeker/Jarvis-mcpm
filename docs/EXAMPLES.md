## Examples of Agent-Driven Workflows

Here are concrete examples of how an Agent uses Jarvis v3.0 to solve real-world problems.

> **Note:** Jarvis v3.0 consolidated 24 tools into 8 action-based tools for 52% context token reduction.

### Scenario 1: Upgrading a Legacy Project
**You:** "Take a look at this Python repo. It's messy. Set up a proper CI pipeline and linting."
**Agent:**
1.  Calls `jarvis_project(action="analyze")` -> Detects `requirements.txt` and no git config.
2.  Calls `jarvis_project(action="devops", project_type="python")`.
3.  Jarvis initializes Git, creates `.pre-commit-config.yaml` with Ruff, and sets up GitHub Actions.
4.  **Result:** The project is now strictly engineered.

### Scenario 2: The "Broken" Database
**You:** "My vector search is failing. I think Qdrant crashed."
**Agent:**
1.  Calls `jarvis_check_status()` -> Confirms Qdrant container is unhealthy.
2.  Calls `jarvis_system(action="restart_infra")`.
3.  Jarvis gracefully restarts the Docker containers.
4.  **Result:** Service restored without you touching the terminal.

### Scenario 3: Context Switching
**You:** "I'm switching from the backend (Go) to the frontend (React). Load the right tools."
**Agent:**
1.  Calls `jarvis_profile(action="suggest")` in the new directory.
2.  Jarvis analyzes the path and client context.
3.  **Result:** Returns `["project-frontend", "client-claude", "memory"]`. The agent now knows exactly which 3 layers to activate.

### Scenario 4: Intelligent Refactoring
**You:** "Refactor the user authentication logic in `auth.py`. It's messy."
**Agent:**
1.  Checks loaded tools -> Sees `morph-fast-apply` (loaded via `project-pokeedge` profile).
2.  Calls `morph.transform_code(path="auth.py", transformation="...")`.
3.  **Result:** The tool applies the semantic refactor without breaking the file, bypassing the need for fragile line-number based diffs.

### Scenario 5: Installing a New MCP Server
**You:** "I need the GitHub MCP server for PR management."
**Agent:**
1.  Calls `jarvis_server(action="search", query="github")` -> Finds `github` server.
2.  Calls `jarvis_server(action="install", name="github")`.
3.  Calls `jarvis_server(action="info", name="github")` -> Shows configuration details.
4.  **Result:** GitHub MCP server installed and ready to use.

### Scenario 6: Managing Client Configurations
**You:** "Set up OpenCode with my project profile and memory."
**Agent:**
1.  Calls `jarvis_client(action="import", client_name="opencode")` -> Imports starter config.
2.  Calls `jarvis_client(action="edit", client_name="opencode", add_profile="p-pokeedge,memory")`.
3.  **Result:** OpenCode configured with Layer 1 (project) and Layer 3 (memory) profiles.

---
