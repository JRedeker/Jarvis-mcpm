# Plan: "Smart Scaffold" for Existing Projects

**Objective:** Enable Jarvis to intelligently analyze an existing codebase (like `pokeedge`) and apply the "Prime DevOps Stack" (AI Review, Pre-commit, Dependabot) without destructive overwrites, leveraging the LLM "Brain" for decision-making.

## Core Philosophy: "Tooling for the Brain"
Instead of hardcoding complex "merge logic" or "auto-detection" into the Go binary (Jarvis), we will expose **granular tools** that allow the LLM to:
1.  **See** the state of the project.
2.  **Read** conflicting files.
3.  **Decide** how to merge/update.
4.  **Execute** the precise write operations.

## 1. New Tool: `analyze_project_structure`
**Purpose:** Give the LLM a high-level map of the project to determine its type and existing tooling status.
**Implementation:** Go (Jarvis)
**Input:** `None` (operates on CWD)
**Output:** JSON Object
```json
{
  "languages": ["go", "python", "javascript"], // Detected via file signatures
  "configs": {
    "has_git": true,
    "has_pre_commit": true, // .pre-commit-config.yaml exists
    "has_github_workflows": true,
    "has_pr_agent": false
  },
  "key_files": ["go.mod", "pyproject.toml", ".github/workflows/test.yml"]
}
```

## 2. Upgrade Tool: `scaffold_project` -> `apply_devops_stack`
**Purpose:** Execute the write operations. We will move away from "one-click magic" that blindly overwrites, to a more flexible "applier" that can take custom content or explicit flags.
**Implementation:** Go (Jarvis)
**Inputs:**
*   `target_file`: (Enum: "pre-commit", "pr-agent", "dependabot")
*   `content`: (String, Optional) If provided, write THIS content. If null, use the internal "Prime" template.
*   `strategy`: (Enum: "create_new", "overwrite", "skip_if_exists")

**Why this change?**
If the LLM sees an existing `.pre-commit-config.yaml`, it can:
1.  Read it (using standard `read_file`).
2.  Compose a *new* YAML that merges existing hooks with our Prime hooks.
3.  Call `apply_devops_stack(target="pre-commit", content="<merged_yaml>", strategy="overwrite")`.

## 3. The Agent Workflow (The "Brain" Script)

When the user says "Setup DevOps for this repo":

1.  **Phase 1: Discovery**
    *   Agent calls `analyze_project_structure`.
    *   Agent identifies: "This is a Python project. It already has a pre-commit config but no PR agent."

2.  **Phase 2: Strategy**
    *   **PR Agent:** Missing. Action -> Create standard Prime config.
    *   **Pre-commit:** Exists. Action -> Read existing file.
    *   **Dependabot:** Missing. Action -> Create standard Prime config.

3.  **Phase 3: Execution**
    *   Agent calls `apply_devops_stack(target="pr-agent", strategy="create_new")`.
    *   Agent calls `read_file(".pre-commit-config.yaml")`.
    *   Agent *internally merges* the Prime hooks (Ruff, Gitleaks) into the text it just read.
    *   Agent calls `apply_devops_stack(target="pre-commit", content="<merged_content>", strategy="overwrite")`.

## 4. Implementation Plan

### Step 1: Implement `analyze_project_structure` in Jarvis
*   Scans for: `go.mod`, `package.json`, `pom.xml`, `pyproject.toml`, `requirements.txt`, `Gemfile`, `composer.json`.
*   Checks for: `.git`, `.pre-commit-config.yaml`, `.github/workflows`.

### Step 2: Refactor `scaffold_project`
*   Rename to `apply_devops_stack` (or keep generic `manage_devops_files`).
*   Remove the rigid "all-or-nothing" logic.
*   Add the `content` override parameter to enable LLM-driven merging.

### Step 3: Verify
*   Test on the current `MCP` repo (simulating an "existing" repo).
*   Ensure the LLM correctly identifies it as Go/Node/Python mix and sees existing configs.
