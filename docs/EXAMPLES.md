## ⚡ Examples of Agent-Driven Workflows

Here are concrete examples of how an Agent uses Jarvis to solve real-world problems.

### Scenario 1: Upgrading a Legacy Project
**You:** "Take a look at this Python repo. It's messy. Set up a proper CI pipeline and linting."
**Agent:**
1.  Calls `jarvis.analyze_project()` -> Detects `requirements.txt` and no git config.
2.  Calls `jarvis.apply_devops_stack(project_type="python")`.
3.  Jarvis initializes Git, creates `.pre-commit-config.yaml` with Ruff, and sets up GitHub Actions.
4.  **Result:** The project is now strictly engineered.

### Scenario 2: The "Broken" Database
**You:** "My vector search is failing. I think Qdrant crashed."
**Agent:**
1.  Calls `jarvis.check_status()` -> Confirms Qdrant container is unhealthy.
2.  Calls `jarvis.restart_infrastructure()`.
3.  Jarvis gracefully restarts the Docker containers.
4.  **Result:** Service restored without you touching the terminal.

### Scenario 3: Context Switching
**You:** "I'm switching from the backend (Go) to the frontend (React). Load the right tools."
**Agent:**
1.  Calls `jarvis.analyze_project()` in the new directory.
2.  Detects `package.json` and React dependencies.
3.  Calls `jarvis.suggest_profile()` -> Returns `["project-frontend", "client-claude"]`.
4.  **Result:** The Agent effectively "reboots" its toolset to match your new context.

---
