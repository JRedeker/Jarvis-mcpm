# MCPM Bundle: Modifications & Fork Details

This document details the architectural changes and modifications made to the bundled version of **MCPM** included in the `MCPM/` directory of the Jarvis repository.

**Note:** The original source code (`mcpm_source/`) is provided for reference, but Jarvis utilizes a **specialized Node.js implementation** (`MCPM/`) to ensure stability and agent-compatibility.

## 1. Architectural Shift: Python to Node.js

### The "Pydantic Problem"
The original MCPM CLI (written in Python) relies heavily on `pydantic` for schema validation. However, the Python ecosystem currently suffers from significant fragmentation between Pydantic v1 and v2.
- **Issue:** Integrating MCPM into diverse environments often leads to `ImportError` or `VersionConflict` due to clashing Pydantic requirements from other tools.
- **Resolution:** To guarantee that Jarvis can bootstrap itself on *any* user machine without complex virtual environment debugging, we **re-implemented the Core CLI in Node.js**.

### The Solution
We replaced the Python runtime requirements with a lightweight Node.js CLI.
- **Source File:** [`MCPM/index.js`](./index.js)
- **Dependencies:** Minimal (`commander`, `toml`, `chalk`).

## 2. New Features for Agents

We added features specifically requested by the Agentic workflow that are not yet in the upstream Python version.

### A. Docker-First Installation
**Location:** [`MCPM/index.js` (Line 85)](./index.js)

We introduced logic to check for a `docker` field in the registry.
- **Behavior:** If a tool supports Docker, Jarvis prefers `docker pull` over `npm install`.
- **Reasoning:** Docker provides better isolation and security for agents running unknown tools.

```javascript
// MCPM/index.js
if (dockerImage) {
    console.log(chalk.blue(`Docker image found: ${dockerImage}. Preferring Docker installation...`));
    // ... executes docker pull
}
```

### B. JSON Configuration Snippets
**Location:** [`MCPM/index.js` (Line 65)](./index.js)

After installation, the CLI now generates the exact JSON snippet needed for the client configuration.
- **Behavior:** Outputs a ready-to-copy JSON block.
- **Impact:** Allows the Agent to "read" the output and immediately configure itself without hallucinating the config format.

```javascript
// MCPM/index.js
const printConfigSnippet = (toolName, type, target) => {
    // ... generates config object
    console.log(JSON.stringify({ [toolName]: config }, null, 2));
};
```

## 3. Registry Adaptations

**Location:** [`MCPM/config/technologies.toml`](./config/technologies.toml)

We maintained compatibility with the TOML registry format but extended it:
- **Docker Field:** Added `docker = "image:tag"` support (e.g., for `mindsdb`).
- **Simplification:** Removed Python-specific class bindings in favor of command-line arguments.

## 4. Summary of Changes

| Feature | Original (Python) | Jarvis Bundle (Node.js) |
| :--- | :--- | :--- |
| **Runtime** | Python 3.10+ | Node.js 18+ |
| **Schema Validation** | Pydantic (Complex) | Native JS Checks (Simple) |
| **Installation** | `pipx install` | `npm link` |
| **Docker Support** | Manual | **Automatic** |
| **Config Output** | Manual Lookup | **Auto-Generated JSON** |
