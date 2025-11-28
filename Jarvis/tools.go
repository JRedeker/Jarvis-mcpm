package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func handleBootstrapSystem(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Find the project root. Assuming Jarvis is run from within Jarvis/ directory or root.
	// We'll try to locate "MCPM" directory.
	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current working directory: %v", err)), nil
	}

	// Logic to find root
	var rootDir string
	if _, err := os.Stat(filepath.Join(cwd, "MCPM")); err == nil {
		rootDir = cwd
	} else if _, err := os.Stat(filepath.Join(cwd, "..", "MCPM")); err == nil {
		rootDir = filepath.Join(cwd, "..")
	} else {
		return mcp.NewToolResultError("Could not locate MCPM directory. Please run Jarvis from the project root or Jarvis subdirectory."), nil
	}

	mcpmDir := filepath.Join(rootDir, "MCPM")

	// 1. Install MCPM dependencies
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = mcpmDir
	if out, err := cmdInstall.CombinedOutput(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run npm install in %s: %v\nOutput: %s", mcpmDir, err, string(out))), nil
	}

	// 2. Link MCPM
	cmdLink := exec.Command("npm", "link")
	cmdLink.Dir = mcpmDir
	if out, err := cmdLink.CombinedOutput(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run npm link in %s: %v\nOutput: %s", mcpmDir, err, string(out))), nil
	}

	// 3. Install Default Servers (The Guardian Stack)
	defaultServers := []string{"context7", "brave-search", "github"}
	for _, server := range defaultServers {
		if _, err := runMcpmCommand("install", server); err != nil {
			// Log but continue, as some might fail if not configured or already installed
			// In a real scenario, we might want to be stricter or check existance first
			// For now, we attempt installation.
			fmt.Printf("Warning: Failed to install default server %s: %v\n", server, err)
		}
	}

	// 4. Start Infrastructure
	// Use the manage-mcp.sh script if available, otherwise fallback to docker-compose
	scriptPath := filepath.Join(rootDir, "scripts", "manage-mcp.sh")
	if _, err := os.Stat(scriptPath); err == nil {
		cmd := exec.Command(scriptPath, "start")
		if out, err := cmd.CombinedOutput(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start infrastructure via script: %v\nOutput: %s", err, string(out))), nil
		}
	} else {
		cmdCompose := exec.Command("docker", "compose", "up", "-d")
		cmdCompose.Dir = rootDir
		if out, err := cmdCompose.CombinedOutput(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to run docker compose up in %s: %v\nOutput: %s", rootDir, err, string(out))), nil
		}
	}

	return mcp.NewToolResultText("System bootstrapped successfully! MCPM installed, default servers (context7, brave-search, github) set up, and Infrastructure started."), nil
}

func handleRestartService(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	return mcp.NewToolResultText("Restarting Jarvis service..."), nil
}

func handleRestartInfrastructure(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Locate script
	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	var rootDir string
	if _, err := os.Stat(filepath.Join(cwd, "MCPM")); err == nil {
		rootDir = cwd
	} else if _, err := os.Stat(filepath.Join(cwd, "..", "MCPM")); err == nil {
		rootDir = filepath.Join(cwd, "..")
	} else {
		return mcp.NewToolResultError("Could not locate project root."), nil
	}

	scriptPath := filepath.Join(rootDir, "scripts", "manage-mcp.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return mcp.NewToolResultError("Management script not found at " + scriptPath), nil
	}

	// Run restart
	cmd := exec.Command(scriptPath, "restart")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restart failed: %v\nOutput: %s", err, string(out))), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Infrastructure restarted successfully.\nOutput:\n%s", string(out))), nil
}

func handleSuggestProfile(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	testingMode, _ := args["testing"].(bool)
	clientName, _ := args["client_name"].(string)

	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current working directory: %v", err)), nil
	}

	// Normalize path
	path := strings.ToLower(cwd)
	var profiles []string

	// LAYER 1: ENVIRONMENT (Base)
	// Mutually exclusive. Determines the workspace context.
	if strings.Contains(path, "pokeedge") {
		profiles = append(profiles, "project-pokeedge")
	} else if strings.Contains(path, "codex") {
		profiles = append(profiles, "project-codex") // Legacy support if folder exists
	} else {
		// Fallback for new/unrecognized projects
		profiles = append(profiles, "project-new")
	}

	// LAYER 2: CLIENT ADAPTERS (Additive)
	// Adds client-specific capabilities (e.g., morph-fast-apply)
	if clientName != "" {
		// Normalize client name
		cn := strings.ToLower(clientName)
		if strings.Contains(cn, "codex") {
			profiles = append(profiles, "client-codex")
		} else if strings.Contains(cn, "gemini") {
			profiles = append(profiles, "client-gemini")
		}
	}

	// LAYER 3: GLOBAL CAPABILITIES (Augment)
	// Always active layers (like memory) or toggles (like testing)

	// Memory is standard for all our agents
	profiles = append(profiles, "memory")

	if testingMode {
		profiles = append(profiles, "testing-all-tools")
	}

	// Format as JSON-like string for easy parsing by agents
	// e.g., "[\"pokeedge\", \"testing-all-tools\"]"
	result := "["
	for i, p := range profiles {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", p)
	}
	result += "]"

	return mcp.NewToolResultText(result), nil
}

func handleFetchDiffContext(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	staged, _ := args["staged"].(bool)

	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	// 1. Get Status
	statusCmd := exec.Command("git", "status", "--short")
	statusOut, err := statusCmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get git status (is this a git repo?): %v", err)), nil
	}

	// 2. Get Diff
	diffArgs := []string{"diff"}
	if staged {
		diffArgs = append(diffArgs, "--staged")
	} else {
		// If not staged, we want HEAD to Working Tree (everything)
		diffArgs = append(diffArgs, "HEAD")
	}

	diffCmd := exec.Command("git", diffArgs...)
	diffOut, err := diffCmd.CombinedOutput()
	if err != nil {
		// Fallback: maybe no commits yet? try just diff
		diffCmd = exec.Command("git", "diff")
		diffOut, _ = diffCmd.CombinedOutput()
	}

	// 3. Format Report
	report := fmt.Sprintf("# Local Review Context\n\n## Working Directory\n`%s`\n\n## Git Status\n```\n%s\n```\n\n## Diff\n```diff\n%s\n```",
		cwd, string(statusOut), string(diffOut))

	return mcp.NewToolResultText(report), nil
}

func handleAnalyzeProject(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	analysis := map[string]interface{}{
		"path":       cwd,
		"languages":  []string{},
		"frameworks": []string{},
		"configs": map[string]bool{
			"has_git":              false,
			"has_pre_commit":       false,
			"has_github_workflows": false,
			"has_pr_agent":         false,
			"has_dependabot":       false,
			"has_gitleaks":         false, // New check
		},
		"key_files": []string{},
	}

	// Check for Git
	if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
		analysis["configs"].(map[string]bool)["has_git"] = true
	}

	// Check configs
	preCommitPath := filepath.Join(cwd, ".pre-commit-config.yaml")
	if _, err := os.Stat(preCommitPath); err == nil {
		analysis["configs"].(map[string]bool)["has_pre_commit"] = true
		// Check for gitleaks in pre-commit
		content, _ := os.ReadFile(preCommitPath)
		if strings.Contains(string(content), "gitleaks") {
			analysis["configs"].(map[string]bool)["has_gitleaks"] = true
		}
	}
	if _, err := os.Stat(filepath.Join(cwd, ".github", "workflows")); err == nil {
		analysis["configs"].(map[string]bool)["has_github_workflows"] = true
		if _, err := os.Stat(filepath.Join(cwd, ".github", "workflows", "pr_agent.yml")); err == nil {
			analysis["configs"].(map[string]bool)["has_pr_agent"] = true
		}
	}
	if _, err := os.Stat(filepath.Join(cwd, ".github", "dependabot.yml")); err == nil {
		analysis["configs"].(map[string]bool)["has_dependabot"] = true
	}

	// Detect Languages & Frameworks
	files, _ := os.ReadDir(cwd)
	for _, f := range files {
		name := f.Name()
		switch name {
		case "go.mod":
			analysis["languages"] = append(analysis["languages"].([]string), "go")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "package.json":
			analysis["languages"] = append(analysis["languages"].([]string), "javascript/typescript")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "pyproject.toml", "requirements.txt", "setup.py":
			found := false
			for _, l := range analysis["languages"].([]string) {
				if l == "python" {
					found = true
					break
				}
			}
			if !found {
				analysis["languages"] = append(analysis["languages"].([]string), "python")
			}
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "pom.xml", "build.gradle":
			analysis["languages"] = append(analysis["languages"].([]string), "java")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "Gemfile":
			analysis["languages"] = append(analysis["languages"].([]string), "ruby")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "composer.json":
			analysis["languages"] = append(analysis["languages"].([]string), "php")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		case "Cargo.toml":
			analysis["languages"] = append(analysis["languages"].([]string), "rust")
			analysis["key_files"] = append(analysis["key_files"].([]string), name)
		}
	}

	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal analysis: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func handleApplyDevOpsStack(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	projectType, _ := args["project_type"].(string)
	enableAiReview, _ := args["enable_ai_review"].(bool)
	force, _ := args["force"].(bool)

	cwd, err := os.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	logs := []string{fmt.Sprintf("🚀 Applying DevOps Stack in %s...", cwd)}
	if projectType != "" {
		logs = append(logs, fmt.Sprintf("Project Type: %s", projectType))
	}

	// 1. Initialize Git
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		if out, err := cmd.CombinedOutput(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Git init failed: %s", out)), nil
		}
		logs = append(logs, "✅ Initialized git repository")
	} else {
		logs = append(logs, "ℹ️ Git repository already exists")
	}

	// 2. Write .pre-commit-config.yaml
	preCommitFile := ".pre-commit-config.yaml"
	if _, err := os.Stat(preCommitFile); err == nil && !force {
		logs = append(logs, fmt.Sprintf("⚠️ %s exists. Skipping. Use 'force=true' to overwrite.", preCommitFile))
	} else {
		// Default config content
		preCommitConfig := `# See https://pre-commit.com for more information
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
      - id: check-merge-conflict

  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.2
    hooks:
      - id: gitleaks
`
		// Language specific hooks
		switch projectType {
		case "python":
			preCommitConfig += `
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.3.0
    hooks:
      - id: ruff
        args: [ --fix ]
      - id: ruff-format
`
		case "go":
			preCommitConfig += `
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      # - id: go-imports # Uncomment if installed locally
`
		case "node", "typescript", "javascript":
			preCommitConfig += `
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v4.0.0
    hooks:
      - id: prettier
`
		}

		if err := os.WriteFile(preCommitFile, []byte(preCommitConfig), 0644); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write pre-commit config: %v", err)), nil
		}
		if force && err == nil { // File existed and we forced
			logs = append(logs, fmt.Sprintf("♻️ Overwrote %s", preCommitFile))
		} else {
			logs = append(logs, fmt.Sprintf("✅ Created %s", preCommitFile))
		}
	}

	// 3. Install pre-commit (Command execution)
	// Only run if we actually have a config or forced
	// We check for pip availability
	if _, err := exec.LookPath("pip"); err == nil {
		// Install package
		exec.Command("pip", "install", "pre-commit").Run()
		// Install hooks
		cmd := exec.Command("pre-commit", "install")
		if out, err := cmd.CombinedOutput(); err != nil {
			logs = append(logs, fmt.Sprintf("⚠️ Failed to install git hooks (is pre-commit in PATH?): %s", out))
		} else {
			logs = append(logs, "✅ Installed git hooks")
		}
	} else {
		logs = append(logs, "⚠️ 'pip' not found. Please install 'pre-commit' manually.")
	}

	// 4. AI Review (GitHub Actions)
	if enableAiReview {
		workflowsDir := filepath.Join(".github", "workflows")
		os.MkdirAll(workflowsDir, 0755)
		workflowFile := filepath.Join(workflowsDir, "pr_agent.yml")

		if _, err := os.Stat(workflowFile); err == nil && !force {
			logs = append(logs, fmt.Sprintf("⚠️ %s exists. Skipping. Use 'force=true' to overwrite.", workflowFile))
		} else {
			prAgentConfig := `name: AI Code Review

on:
  pull_request:
    types: [opened, reopened, ready_for_review, synchronize]
  issue_comment:
    types: [created, edited]

permissions:
  issues: write
  pull-requests: write
  contents: read

jobs:
  pr_agent:
    runs-on: ubuntu-latest
    name: PR Agent
    if: ${{ github.event.sender.type != 'Bot' }}
    steps:
      - id: pr-agent
        uses: Codium-ai/pr-agent@main
        env:
          OPENAI_KEY: ${{ secrets.OPENAI_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

          # --- Configuration ---
          # General
          PR_REVIEW__EXTRA_INSTRUCTIONS: "Focus on architectural consistency, 3-Layer Profile compliance, and security. CHECK DOCUMENTATION: If code logic is changed, ensure comments and markdown docs are updated."
          PR_REVIEW__REQUIRE_TESTS_REVIEW: "true"
          PR_CODE_SUGGESTIONS__NUM_CODE_SUGGESTIONS: 4

          # Commands enabled
          PR_COMMANDS__ALLOW_REVIEW: "true"
          PR_COMMANDS__ALLOW_DESCRIBE: "true"
          PR_COMMANDS__ALLOW_IMPROVE: "true"
`
			if err := os.WriteFile(workflowFile, []byte(prAgentConfig), 0644); err != nil {
				logs = append(logs, fmt.Sprintf("⚠️ Failed to write workflow: %v", err))
			} else {
				if force && err == nil {
					logs = append(logs, fmt.Sprintf("♻️ Overwrote %s", workflowFile))
				} else {
					logs = append(logs, fmt.Sprintf("✅ Created %s", workflowFile))
				}
			}
		}
	}

	// 5. Gitignore (Basic)
	gitignoreFile := ".gitignore"
	if _, err := os.Stat(gitignoreFile); os.IsNotExist(err) {
		gitignore := ".env\n.venv/\nnode_modules/\ndist/\n*.log\n.DS_Store\n"
		os.WriteFile(gitignoreFile, []byte(gitignore), 0644)
		logs = append(logs, fmt.Sprintf("✅ Created default %s", gitignoreFile))
	} else {
		// We rarely overwrite gitignore completely, usually append logic is needed but let's just skip for now to be safe
		logs = append(logs, fmt.Sprintf("ℹ️ %s exists (skipping)", gitignoreFile))
	}

	return mcp.NewToolResultText(strings.Join(logs, "\n")), nil
}

func handleListServers(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, err := runMcpmCommand("ls")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list servers: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleInstallServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	output, err := runMcpmCommand("install", name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to install server %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleServerInfo(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	output, err := runMcpmCommand("info", name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get info for server %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleCheckStatus(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, err := runMcpmCommand("doctor")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to check status: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleSearchServers(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	query, ok := args["query"].(string)
	if !ok {
		return mcp.NewToolResultError("query argument is required"), nil
	}

	output, err := runMcpmCommand("search", query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search servers: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleUninstallServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	output, err := runMcpmCommand("uninstall", name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to uninstall server %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleEditServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	cmdArgs := []string{"edit", name}
	if val, ok := args["command"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--command", val)
	}
	if val, ok := args["args"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--args", val)
	}
	if val, ok := args["env"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--env", val)
	}
	if val, ok := args["url"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--url", val)
	}
	if val, ok := args["headers"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--headers", val)
	}

	output, err := runMcpmCommand(cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to edit server %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleCreateServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}
	serverType, ok := args["type"].(string)
	if !ok {
		return mcp.NewToolResultError("type argument is required"), nil
	}

	cmdArgs := []string{"new", name, "--type", serverType, "--force"} // Force non-interactive

	if cmd, ok := args["command"].(string); ok && cmd != "" {
		cmdArgs = append(cmdArgs, "--command", cmd)
	}
	if argStr, ok := args["args"].(string); ok && argStr != "" {
		cmdArgs = append(cmdArgs, "--args", argStr)
	}
	if envStr, ok := args["env"].(string); ok && envStr != "" {
		cmdArgs = append(cmdArgs, "--env", envStr)
	}
	if url, ok := args["url"].(string); ok && url != "" {
		cmdArgs = append(cmdArgs, "--url", url)
	}
	if headersStr, ok := args["headers"].(string); ok && headersStr != "" {
		cmdArgs = append(cmdArgs, "--headers", headersStr)
	}

	output, err := runMcpmCommand(cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create server %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleUsageStats(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, err := runMcpmCommand("usage")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get usage stats: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func buildManageClientArgs(args map[string]interface{}) []string {
	cmdArgs := []string{}
	if val, ok := args["add_server"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--add-server", val)
	}
	if val, ok := args["remove_server"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--remove-server", val)
	}
	if val, ok := args["add_profile"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--add-profile", val)
	}
	if val, ok := args["remove_profile"].(string); ok && val != "" {
		cmdArgs = append(cmdArgs, "--remove-profile", val)
	}
	return cmdArgs
}

func handleManageClient(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	action, ok := args["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action argument is required"), nil
	}

	cmdArgs := []string{"client", action}

	if action == "edit" || action == "import" || action == "config" {
		clientName, ok := args["client_name"].(string)
		if !ok || clientName == "" {
			return mcp.NewToolResultError("client_name argument is required for this action"), nil
		}
		cmdArgs = append(cmdArgs, clientName)
	}

	if action == "edit" {
		cmdArgs = append(cmdArgs, buildManageClientArgs(args)...)
	}

	if action == "config" {
		if path, ok := args["config_path"].(string); ok && path != "" {
			cmdArgs = append(cmdArgs, "--set-path", path)
		} else {
			// If no path provided, assume get path
			cmdArgs = append(cmdArgs, "--get-path")
		}
	}

	output, err := runMcpmCommand(cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to manage client: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleManageProfile(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	action, ok := args["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action argument is required"), nil
	}

	// Map 'delete' to 'rm' for the CLI
	cliAction := action
	if action == "delete" {
		cliAction = "rm"
	}

	cmdArgs := []string{"profile", cliAction}
	if name, ok := args["name"].(string); ok && name != "" {
		cmdArgs = append(cmdArgs, name)
	}

	// Handle edit arguments
	if action == "edit" {
		if newName, ok := args["new_name"].(string); ok && newName != "" {
			cmdArgs = append(cmdArgs, "--name", newName)
		}
		if add, ok := args["add_servers"].(string); ok && add != "" {
			cmdArgs = append(cmdArgs, "--add-server", add)
		}
		if remove, ok := args["remove_servers"].(string); ok && remove != "" {
			cmdArgs = append(cmdArgs, "--remove-server", remove)
		}
	}

	// For delete/rm, we might need force to avoid interactive prompt if it exists (though profile rm might not prompt by default, safe to add if supported)
	// Checking help, 'rm' has --force.
	if cliAction == "rm" {
		cmdArgs = append(cmdArgs, "--force")
	}

	output, err := runMcpmCommand(cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to manage profile: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleManageConfig(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	action, ok := args["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action argument is required"), nil
	}

	cmdArgs := []string{"config", action}
	if key, ok := args["key"].(string); ok && key != "" {
		cmdArgs = append(cmdArgs, key)
	}
	if value, ok := args["value"].(string); ok && value != "" {
		cmdArgs = append(cmdArgs, value)
	}

	output, err := runMcpmCommand(cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to manage config: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func handleMigrateConfig(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, err := runMcpmCommand("migrate")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to migrate config: %v", err)), nil
	}
	return mcp.NewToolResultText(output), nil
}

func monitorShareProcess(stdout io.Reader, success, failure chan string) {
	scanner := bufio.NewScanner(stdout)
	outputBuilder := strings.Builder{}
	// Set a timeout for startup
	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-timeout:
			failure <- fmt.Sprintf("Timed out waiting for share URL. Output so far:\n%s", outputBuilder.String())
			return
		default:
			if !scanner.Scan() {
				failure <- fmt.Sprintf("Process exited unexpectedly. Output:\n%s", outputBuilder.String())
				return
			}
			line := scanner.Text()
			outputBuilder.WriteString(line + "\n")

			// Look for success indicators
			// Adjust these checks based on actual mcpm share output
			if strings.Contains(line, "Public URL:") || strings.Contains(line, "https://") {
				success <- outputBuilder.String()
				return
			}
		}
	}
}

func handleShareServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	sharedServersMutex.Lock()
	if _, exists := sharedServers[name]; exists {
		sharedServersMutex.Unlock()
		return mcp.NewToolResultError(fmt.Sprintf("Server %s is already being shared", name)), nil
	}
	sharedServersMutex.Unlock()

	cmdArgs := []string{"share", name}
	if port, ok := args["port"].(string); ok && port != "" {
		cmdArgs = append(cmdArgs, "--port", port)
	}
	if noAuth, ok := args["no_auth"].(bool); ok && noAuth {
		cmdArgs = append(cmdArgs, "--no-auth")
	}

	// Run mcpm share in background
	cmd := exec.Command("mcpm", cmdArgs...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true")

	// Create pipes for stdout/stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create stdout pipe: %v", err)), nil
	}
	// stderr, err := cmd.StderrPipe() // Optional: capture stderr for logging

	if err := cmd.Start(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start share command: %v", err)), nil
	}

	// Register the process immediately
	sharedServersMutex.Lock()
	sharedServers[name] = cmd
	sharedServersMutex.Unlock()

	// Read stdout to find the URL
	success := make(chan string)
	failure := make(chan string)

	go monitorShareProcess(stdout, success, failure)

	select {
	case output := <-success:
		return mcp.NewToolResultText(output), nil
	case errStr := <-failure:
		// Cleanup if failed
		_ = cmd.Process.Kill()
		sharedServersMutex.Lock()
		delete(sharedServers, name)
		sharedServersMutex.Unlock()
		return mcp.NewToolResultError(errStr), nil
	}
}

func handleStopSharingServer(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	sharedServersMutex.Lock()
	cmd, exists := sharedServers[name]
	if !exists {
		sharedServersMutex.Unlock()
		return mcp.NewToolResultError(fmt.Sprintf("Server %s is not currently shared", name)), nil
	}
	delete(sharedServers, name)
	sharedServersMutex.Unlock()

	if err := cmd.Process.Kill(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop sharing server %s: %v", name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped sharing server %s", name)), nil
}

func handleListSharedServers(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sharedServersMutex.Lock()
	defer sharedServersMutex.Unlock()

	if len(sharedServers) == 0 {
		return mcp.NewToolResultText("No servers are currently being shared."), nil
	}

	var builder strings.Builder
	builder.WriteString("Currently shared servers:\n")
	for name := range sharedServers {
		builder.WriteString(fmt.Sprintf("- %s\n", name))
	}

	return mcp.NewToolResultText(builder.String()), nil
}
