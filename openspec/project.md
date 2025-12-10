# Project Context

## Purpose
**Jarvis** is an intelligent infrastructure layer designed to bridge the gap between Large Language Models (LLMs) and local development environments. It serves as a "Tools Server for Agents," enabling AI assistants to autonomously manage infrastructure, enforce DevOps best practices, and dynamically install and configure other tools (MCP servers).

Its core mission is to empower AI agents to be self-sufficient by providing them with a safe, standard, and powerful interface (MCP) to interact with the underlying system (Docker, Git, CLI tools).

## Tech Stack
- **Core Server (Jarvis):** Go (Golang) - chosen for performance, static typing, and concurrency.
- **CLI Manager (MCPM):** Node.js - chosen for its rich ecosystem and ease of scripting.
- **Reference Implementation (mcpm_source):** Python - likely the original prototype or logic source (FastMCP).
- **Infrastructure:** Docker & Docker Compose.
  - **Databases:** PostgreSQL (relational), Qdrant (vector/memory).
  - **Runtime:** `mcpm-daemon` (Python/Node environment for running tools).
- **Scripts:** Bash (for glue logic and testing).

## Project Conventions

### Code Style
- **Go:** Adheres to standard Go formatting (`gofmt`, `goimports`). Linting via `golangci-lint` (pinned version v1.62.2).
- **Python:** Follows PEP 8, enforced by `ruff` (linting and formatting).
- **JavaScript/Node:** Standard JS conventions, typically enforced by `prettier` (via pre-commit).
- **Security:** Secret scanning via `gitleaks` in pre-commit hooks.

### Architecture Patterns
- **Intelligent Facade:** `Jarvis` acts as a facade over standard CLI tools (`mcpm`, `docker`, `git`). It intercepts tool calls, adds validation, handling, and formats the output to be "AI-friendly" (concise, actionable).
- **Dependency Injection:** The Go codebase uses interface-based dependency injection (in `handlers.go`) to allow robust unit testing by mocking the CLI runners (`McpmRunner`, `DockerRunner`, `GitRunner`).
- **Microservices-lite:** Components are containerized (via `mcpm-daemon`) but managed centrally by Jarvis.

### Testing Strategy
- **Layered Testing:**
  - **Unit Tests (Go):** Comprehensive coverage for handlers and logic (`go test -v ./...`).
  - **Integration Tests (Shell):** `bats` framework used to test shell scripts (`scripts/tests/`).
  - **Linting:** Strict linting for all languages in CI.
  - **Build Verification:** CI ensures Docker images and binaries build successfully.
- **Continuous Integration:** GitHub Actions (`.github/workflows/`) run all tests on Pull Requests.

### Git Workflow
- **Branching:** Feature branches (e.g., `feat/`, `refactor/`, `fix/`) merged into `main` via Pull Request.
- **Commits:** Semantic commit messages (e.g., `feat:`, `fix:`, `chore:`) are encouraged/enforced.
- **Review:** Automated AI code review via PR Agent and standard human review.

## Domain Context
- **Model Context Protocol (MCP):** A standard for connecting AI assistants to systems. Jarvis is an MCP Server *that manages other MCP Servers*.
- **Agent Autonomy:** The design assumes the user is an AI Agent or a human developer working with one. Output is optimized for machine reasoning (clear success/failure signals, suggested next steps).

## Important Constraints
- **Docker Requirement:** The system relies heavily on Docker for the `mcpm-daemon` and databases.
- **Environment:** Designed for Linux/Unix-like environments.
- **Compatibility:** Requires Go 1.24+ and Node.js 18+.

## External Dependencies
- **Docker Hub:** For base images.
- **GitHub:** For hosting, CI/CD, and dependency management.
- **MCP Registry:** The source of truth for installable MCP servers (managed by MCPM).
