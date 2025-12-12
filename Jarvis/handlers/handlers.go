// Package handlers provides MCP tool handlers with dependency injection for testing
package handlers

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
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// McpmRunner is the interface for running MCPM commands
type McpmRunner interface {
	Run(args ...string) (string, error)
}

// DockerRunner is the interface for Docker operations
type DockerRunner interface {
	ComposeUp(ctx context.Context, services ...string) error
	ComposeDown(ctx context.Context) error
	ComposeRestart(ctx context.Context, services ...string) error
	ComposePs(ctx context.Context) ([]ContainerStatus, error)
	ExecSupervisorctl(ctx context.Context, action, target string) (string, error)
}

// ContainerStatus represents a Docker container's status
type ContainerStatus struct {
	Name    string
	Status  string
	Health  string
	Running bool
	Ports   []string
}

// GitRunner is the interface for Git operations
type GitRunner interface {
	Status(ctx context.Context) (string, error)
	Diff(ctx context.Context, staged bool) (string, error)
	Init(ctx context.Context) error
}

// FileSystem is the interface for file system operations
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	Getwd() (string, error)
}

// CommandRunner is the interface for running shell commands
type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
	RunInDir(ctx context.Context, dir, name string, args ...string) (string, error)
	StartBackground(ctx context.Context, name string, args ...string) (Process, error)
}

// Process represents a running process
type Process interface {
	Kill() error
	Wait() error
	Stdout() io.Reader
}

// ProcessManager manages shared server processes
type ProcessManager interface {
	Register(name string, proc Process)
	Get(name string) (Process, bool)
	Remove(name string) bool
	List() []string
}

// ExitFunc is a function that exits the process (for testing)
type ExitFunc func(code int)

// Handler contains dependencies for all tool handlers
type Handler struct {
	Mcpm        McpmRunner
	Docker      DockerRunner
	Git         GitRunner
	FS          FileSystem
	Cmd         CommandRunner
	Processes   ProcessManager
	ExitProcess ExitFunc
}

// NewHandler creates a new Handler with the given dependencies
func NewHandler(mcpm McpmRunner, docker DockerRunner, git GitRunner, fs FileSystem) *Handler {
	return &Handler{
		Mcpm:        mcpm,
		Docker:      docker,
		Git:         git,
		FS:          fs,
		Cmd:         &RealCommandRunner{},
		Processes:   NewInMemoryProcessManager(),
		ExitProcess: os.Exit,
	}
}

// NewHandlerWithAll creates a Handler with all dependencies explicitly provided
func NewHandlerWithAll(mcpm McpmRunner, docker DockerRunner, git GitRunner, fs FileSystem, cmd CommandRunner, procs ProcessManager, exit ExitFunc) *Handler {
	return &Handler{
		Mcpm:        mcpm,
		Docker:      docker,
		Git:         git,
		FS:          fs,
		Cmd:         cmd,
		Processes:   procs,
		ExitProcess: exit,
	}
}

// RealCommandRunner implements CommandRunner using os/exec
type RealCommandRunner struct{}

func (r *RealCommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (r *RealCommandRunner) RunInDir(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (r *RealCommandRunner) StartBackground(ctx context.Context, name string, args ...string) (Process, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &RealProcess{cmd: cmd, stdout: stdout}, nil
}

// RealProcess wraps an exec.Cmd as a Process
type RealProcess struct {
	cmd    *exec.Cmd
	stdout io.Reader
}

func (p *RealProcess) Kill() error {
	if p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}

func (p *RealProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *RealProcess) Stdout() io.Reader {
	return p.stdout
}

// InMemoryProcessManager manages processes in memory
type InMemoryProcessManager struct {
	mu        sync.RWMutex
	processes map[string]Process
}

// NewInMemoryProcessManager creates a new process manager
func NewInMemoryProcessManager() *InMemoryProcessManager {
	return &InMemoryProcessManager{
		processes: make(map[string]Process),
	}
}

func (m *InMemoryProcessManager) Register(name string, proc Process) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processes[name] = proc
}

func (m *InMemoryProcessManager) Get(name string) (Process, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.processes[name]
	return p, ok
}

func (m *InMemoryProcessManager) Remove(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.processes[name]; ok {
		delete(m.processes, name)
		return true
	}
	return false
}

func (m *InMemoryProcessManager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.processes))
	for name := range m.processes {
		names = append(names, name)
	}
	return names
}

// CheckStatus handles the check_status tool
func (h *Handler) CheckStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, _ := h.Mcpm.Run("doctor")

	// Check supervisor status if available
	supOutput, err := h.Docker.ExecSupervisorctl(ctx, "status", "")
	if err == nil && supOutput != "" {
		// Filter out "No token data found" warnings
		lines := strings.Split(supOutput, "\n")
		var filteredLines []string
		for _, line := range lines {
			if !strings.Contains(line, "No token data found") && strings.TrimSpace(line) != "" {
				filteredLines = append(filteredLines, line)
			}
		}

		if len(filteredLines) > 0 {
			output += "\n\n## 🕵️ Daemon Process Status (Supervisor)\n"
			output += "```\n" + strings.Join(filteredLines, "\n") + "\n```"
		}
	}

	// Check if the output indicates success
	if strings.Contains(output, "All systems healthy") {
		output += "\n\n🚀 **ALL SYSTEMS GO!** 🚀\n**Jarvis is ready to assist.**"
	}

	return mcp.NewToolResultText(output), nil
}

// ListServers handles the list_servers tool
func (h *Handler) ListServers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, _ := h.Mcpm.Run("ls")
	return mcp.NewToolResultText(output), nil
}

// ServerInfo handles the server_info tool
func (h *Handler) ServerInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	output, _ := h.Mcpm.Run("info", name)
	return mcp.NewToolResultText(output), nil
}

// InstallServer handles the install_server tool
func (h *Handler) InstallServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return mcp.NewToolResultError("❌ Server name is required and cannot be empty.\n\n💡 Tip: Use search_servers(query) to find available servers, or server_info(name) for installation details."), nil
	}

	// Validate server name doesn't contain invalid characters
	if strings.ContainsAny(name, " /\\") {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Invalid server name '%s' - server names cannot contain spaces or slashes.\n\n💡 Tip: Server names use hyphens, e.g., 'brave-search' not 'brave search'", name)), nil
	}

	output, err := h.Mcpm.Run("install", name)

	// Add helpful context on common errors
	if err != nil {
		if strings.Contains(output, "not found") || strings.Contains(output, "404") {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Server '%s' not found in registry.\n\n💡 Next steps:\n1. Use search_servers(\"%s\") to find similar servers\n2. Check spelling - server names are case-sensitive\n3. Visit MCPM registry for full list of available servers\n\n%s", name, name, output)), nil
		}
		if strings.Contains(output, "already installed") {
			return mcp.NewToolResultText(fmt.Sprintf("✅ Server '%s' is already installed.\n\n💡 Next step: Use manage_profile() to add it to a profile, or restart_service() if you just updated it.", name)), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully installed '%s'\n\n%s\n\n💡 Next step: Use manage_profile(\"edit\", \"your-profile\", add_servers=\"%s\") to add it to a profile.", name, output, name)), nil
}

// UninstallServer handles the uninstall_server tool
func (h *Handler) UninstallServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return mcp.NewToolResultError("❌ Server name is required and cannot be empty.\n\n💡 Tip: Use list_servers() to see all installed servers."), nil
	}

	output, err := h.Mcpm.Run("uninstall", name)

	// Add helpful context on common errors
	if err != nil {
		if strings.Contains(output, "not found") || strings.Contains(output, "not installed") {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Server '%s' is not installed.\n\n💡 Next steps:\n1. Use list_servers() to see all installed servers\n2. Check if the server is in a profile instead of globally installed\n3. Verify the server name is spelled correctly\n\n%s", name, output)), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully uninstalled '%s'\n\n%s\n\n⚠️  Remember: If this server was in any profiles, you may need to update those profiles using manage_profile().", name, output)), nil
}

// SearchServers handles the search_servers tool
func (h *Handler) SearchServers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	query, ok := args["query"].(string)
	if !ok || strings.TrimSpace(query) == "" {
		return mcp.NewToolResultError("❌ Search query is required and cannot be empty.\n\n💡 Tip: Try broad searches like 'database', 'web', or 'documentation' to discover servers by category."), nil
	}

	output, err := h.Mcpm.Run("search", query)

	// Add helpful context for no results
	if err == nil && (strings.Contains(output, "No servers found") || strings.Contains(output, "0 results")) {
		return mcp.NewToolResultText(fmt.Sprintf("❌ No servers found matching '%s'\n\n💡 Try these tips:\n1. Use broader search terms (e.g., 'web' instead of 'web scraping')\n2. Try common categories: database, file, api, documentation, testing\n3. Use list_servers() to see all installed servers\n4. Check MCPM registry for the full catalog\n\n%s", query, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s\n\n💡 Next step: Use server_info(name) to see installation details before installing.", output)), nil
}

// ManageProfile handles the manage_profile tool
func (h *Handler) ManageProfile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// For delete/rm, add force to avoid interactive prompt
	if cliAction == "rm" {
		cmdArgs = append(cmdArgs, "--force")
	}

	output, _ := h.Mcpm.Run(cmdArgs...)
	return mcp.NewToolResultText(output), nil
}

// buildManageClientArgs extracts client edit arguments
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

// ManageClient handles the manage_client tool
func (h *Handler) ManageClient(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	action, ok := args["action"].(string)
	if !ok || strings.TrimSpace(action) == "" {
		return mcp.NewToolResultError("❌ Action is required and cannot be empty.\n\n💡 Valid actions: ls, edit, import, config\n\nExamples:\n- manage_client(\"ls\") - List all configured clients\n- manage_client(\"edit\", client_name=\"codex\", add_server=\"brave-search\")"), nil
	}

	// Validate action is one of the allowed values
	validActions := map[string]bool{"ls": true, "edit": true, "import": true, "config": true}
	if !validActions[action] {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Invalid action '%s'\n\n💡 Valid actions: ls, edit, import, config\n\n- ls: List all configured clients\n- edit: Add/remove servers or profiles from a client\n- import: Import client configuration\n- config: Get or set client config path", action)), nil
	}

	cmdArgs := []string{"client", action}

	if action == "edit" || action == "import" || action == "config" {
		clientName, ok := args["client_name"].(string)
		if !ok || strings.TrimSpace(clientName) == "" {
			return mcp.NewToolResultError(fmt.Sprintf("❌ client_name is required for action '%s'\n\n💡 Tip: Use manage_client(\"ls\") to see available client names\n\nCommon clients: codex, claude-code, claude-desktop, gemini, kilocode", action)), nil
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
			cmdArgs = append(cmdArgs, "--get-path")
		}
	}

	output, err := h.Mcpm.Run(cmdArgs...)

	// Add helpful context on success
	if err == nil && action == "edit" {
		return mcp.NewToolResultText(fmt.Sprintf("✅ Client configuration updated\n\n%s\n\n💡 Next step: Restart the client to apply changes, or use restart_service() if modifying Jarvis itself.", output)), nil
	}

	return mcp.NewToolResultText(output), nil
}

// ManageConfig handles the manage_config tool
func (h *Handler) ManageConfig(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	output, _ := h.Mcpm.Run(cmdArgs...)
	return mcp.NewToolResultText(output), nil
}

// EditServer handles the edit_server tool
func (h *Handler) EditServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	output, _ := h.Mcpm.Run(cmdArgs...)
	return mcp.NewToolResultText(output), nil
}

// CreateServer handles the create_server tool
func (h *Handler) CreateServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	cmdArgs := []string{"new", name, "--type", serverType, "--force"}

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

	output, _ := h.Mcpm.Run(cmdArgs...)
	return mcp.NewToolResultText(output), nil
}

// UsageStats handles the usage_stats tool
func (h *Handler) UsageStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, _ := h.Mcpm.Run("usage")
	return mcp.NewToolResultText(output), nil
}

// MigrateConfig handles the migrate_config tool
func (h *Handler) MigrateConfig(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	output, _ := h.Mcpm.Run("migrate")
	return mcp.NewToolResultText(output), nil
}

// RestartProfiles handles the restart_profiles tool
func (h *Handler) RestartProfiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	profile, _ := args["profile"].(string)

	var output string
	var err error
	var actionDesc string

	if profile != "" {
		// Restart specific profile via supervisorctl inside the container
		output, err = h.Docker.ExecSupervisorctl(ctx, "restart", "mcpm-"+profile)
		actionDesc = fmt.Sprintf("profile '%s'", profile)
	} else {
		// Restart the entire mcpm-daemon container
		err = h.Docker.ComposeRestart(ctx, "mcpm-daemon")
		actionDesc = "all profiles (mcpm-daemon container)"
		if err == nil {
			output = "Container restarted successfully"
		}
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to restart %s: %v\nOutput: %s", actionDesc, err, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully restarted %s\n\nOutput:\n%s", actionDesc, output)), nil
}

// SuggestProfile handles the suggest_profile tool
func (h *Handler) SuggestProfile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	testingMode, _ := args["testing"].(bool)

	cwd, err := h.FS.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current working directory: %v", err)), nil
	}

	// Normalize path
	path := strings.ToLower(cwd)
	var profiles []string

	// LAYER 1: PROJECT (Auto-detected)
	if strings.Contains(path, "pokeedge") {
		profiles = append(profiles, "p-pokeedge")
	} else if strings.Contains(path, "codex") {
		profiles = append(profiles, "p-codex")
	} else {
		profiles = append(profiles, "p-new")
	}

	// LAYER 3: ENVIRONMENT (Auto-applied globals)
	profiles = append(profiles, "memory")

	if testingMode {
		profiles = append(profiles, "testing-all-tools")
	}

	// Format as JSON-like string
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

// FetchDiffContext handles the fetch_diff_context tool
func (h *Handler) FetchDiffContext(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	staged, _ := args["staged"].(bool)

	cwd, err := h.FS.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	// Get Status
	statusOut, err := h.Git.Status(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get git status (is this a git repo?): %v", err)), nil
	}

	// Get Diff
	diffOut, _ := h.Git.Diff(ctx, staged)

	// Format Report
	report := fmt.Sprintf("# Local Review Context\n\n## Working Directory\n`%s`\n\n## Git Status\n```\n%s\n```\n\n## Diff\n```diff\n%s\n```",
		cwd, statusOut, diffOut)

	return mcp.NewToolResultText(report), nil
}

// AnalyzeProject handles the analyze_project tool
func (h *Handler) AnalyzeProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cwd, err := h.FS.Getwd()
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
			"has_gitleaks":         false,
		},
		"key_files": []string{},
	}

	// Check for Git
	if _, err := h.FS.Stat(filepath.Join(cwd, ".git")); err == nil {
		analysis["configs"].(map[string]bool)["has_git"] = true
	}

	// Check configs
	preCommitPath := filepath.Join(cwd, ".pre-commit-config.yaml")
	if _, err := h.FS.Stat(preCommitPath); err == nil {
		analysis["configs"].(map[string]bool)["has_pre_commit"] = true
		// Check for gitleaks in pre-commit
		content, _ := h.FS.ReadFile(preCommitPath)
		if strings.Contains(string(content), "gitleaks") {
			analysis["configs"].(map[string]bool)["has_gitleaks"] = true
		}
	}
	if _, err := h.FS.Stat(filepath.Join(cwd, ".github", "workflows")); err == nil {
		analysis["configs"].(map[string]bool)["has_github_workflows"] = true
		if _, err := h.FS.Stat(filepath.Join(cwd, ".github", "workflows", "pr_agent.yml")); err == nil {
			analysis["configs"].(map[string]bool)["has_pr_agent"] = true
		}
	}
	if _, err := h.FS.Stat(filepath.Join(cwd, ".github", "dependabot.yml")); err == nil {
		analysis["configs"].(map[string]bool)["has_dependabot"] = true
	}

	// Detect Languages & Frameworks
	files, _ := h.FS.ReadDir(cwd)
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

// RealMcpmRunner implements McpmRunner by executing real MCPM commands
type RealMcpmRunner struct{}

// Run executes an MCPM command
func (r *RealMcpmRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("mcpm", args...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true", "NO_COLOR=true")

	output, err := cmd.CombinedOutput()
	outputStr := stripMcpmNoise(string(output))

	if err != nil {
		return fmt.Sprintf("Error: %v\n\n%s", err, outputStr), fmt.Errorf("command failed: %v", err)
	}

	return strings.TrimSpace(outputStr), nil
}

// stripMcpmNoise removes common warnings and noise from MCPM output
func stripMcpmNoise(output string) string {
	lines := strings.Split(output, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(line, "Warning: Input is not a terminal") {
			continue
		}
		if strings.Contains(line, "(fd=0)") && strings.Contains(line, "Warning:") {
			continue
		}
		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}

// RealFileSystem implements FileSystem using the real os package
type RealFileSystem struct{}

func (r *RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (r *RealFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (r *RealFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (r *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (r *RealFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (r *RealFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

// ApplyDevOpsStack handles the apply_devops_stack tool
func (h *Handler) ApplyDevOpsStack(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	projectType, _ := args["project_type"].(string)
	enableAiReview, _ := args["enable_ai_review"].(bool)
	force, _ := args["force"].(bool)

	cwd, err := h.FS.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	logs := []string{fmt.Sprintf("🚀 Applying DevOps Stack in %s...", cwd)}
	if projectType != "" {
		logs = append(logs, fmt.Sprintf("Project Type: %s", projectType))
	}

	// 1. Initialize Git
	if _, err := h.FS.Stat(filepath.Join(cwd, ".git")); os.IsNotExist(err) {
		if err := h.Git.Init(ctx); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Git init failed: %v", err)), nil
		}
		logs = append(logs, "✅ Initialized git repository")
	} else {
		logs = append(logs, "ℹ️ Git repository already exists")
	}

	// 2. Write .pre-commit-config.yaml
	preCommitFile := filepath.Join(cwd, ".pre-commit-config.yaml")
	if _, err := h.FS.Stat(preCommitFile); err == nil && !force {
		logs = append(logs, fmt.Sprintf("⚠️ %s exists. Skipping. Use 'force=true' to overwrite.", ".pre-commit-config.yaml"))
	} else {
		preCommitConfig := generatePreCommitConfig(projectType)
		if err := h.FS.WriteFile(preCommitFile, []byte(preCommitConfig), 0644); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write pre-commit config: %v", err)), nil
		}
		if force {
			logs = append(logs, "♻️ Overwrote .pre-commit-config.yaml")
		} else {
			logs = append(logs, "✅ Created .pre-commit-config.yaml")
		}
	}

	// 3. AI Review (GitHub Actions)
	if enableAiReview {
		workflowsDir := filepath.Join(cwd, ".github", "workflows")
		h.FS.MkdirAll(workflowsDir, 0755)
		workflowFile := filepath.Join(workflowsDir, "pr_agent.yml")

		if _, err := h.FS.Stat(workflowFile); err == nil && !force {
			logs = append(logs, fmt.Sprintf("⚠️ %s exists. Skipping. Use 'force=true' to overwrite.", "pr_agent.yml"))
		} else {
			prAgentConfig := generatePRAgentConfig()
			if err := h.FS.WriteFile(workflowFile, []byte(prAgentConfig), 0644); err != nil {
				logs = append(logs, fmt.Sprintf("⚠️ Failed to write workflow: %v", err))
			} else {
				if force {
					logs = append(logs, "♻️ Overwrote pr_agent.yml")
				} else {
					logs = append(logs, "✅ Created pr_agent.yml")
				}
			}
		}
	}

	// 4. Gitignore
	gitignoreFile := filepath.Join(cwd, ".gitignore")
	if _, err := h.FS.Stat(gitignoreFile); os.IsNotExist(err) {
		gitignore := ".env\n.venv/\nnode_modules/\ndist/\n*.log\n.DS_Store\n"
		h.FS.WriteFile(gitignoreFile, []byte(gitignore), 0644)
		logs = append(logs, "✅ Created default .gitignore")
	} else {
		logs = append(logs, "ℹ️ .gitignore exists (skipping)")
	}

	return mcp.NewToolResultText(strings.Join(logs, "\n")), nil
}

// generatePreCommitConfig creates pre-commit config based on project type
func generatePreCommitConfig(projectType string) string {
	config := `# See https://pre-commit.com for more information
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
	switch projectType {
	case "python":
		config += `
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.3.0
    hooks:
      - id: ruff
        args: [ --fix ]
      - id: ruff-format
`
	case "go":
		config += `
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
`
	case "node", "typescript", "javascript":
		config += `
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v4.0.0
    hooks:
      - id: prettier
`
	}
	return config
}

// generatePRAgentConfig creates the PR Agent workflow config
func generatePRAgentConfig() string {
	return `name: AI Code Review

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
          PR_REVIEW__EXTRA_INSTRUCTIONS: "Focus on architectural consistency and security."
          PR_REVIEW__REQUIRE_TESTS_REVIEW: "true"
          PR_CODE_SUGGESTIONS__NUM_CODE_SUGGESTIONS: 4
`
}

// BootstrapSystem handles the bootstrap_system tool
func (h *Handler) BootstrapSystem(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cwd, err := h.FS.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current working directory: %v", err)), nil
	}

	// Find project root with MCPM directory
	var rootDir string
	if _, err := h.FS.Stat(filepath.Join(cwd, "MCPM")); err == nil {
		rootDir = cwd
	} else if _, err := h.FS.Stat(filepath.Join(cwd, "..", "MCPM")); err == nil {
		rootDir = filepath.Join(cwd, "..")
	} else {
		return mcp.NewToolResultError("Could not locate MCPM directory. Please run Jarvis from the project root or Jarvis subdirectory."), nil
	}

	mcpmDir := filepath.Join(rootDir, "MCPM")

	// 1. Install MCPM dependencies
	output, err := h.Cmd.RunInDir(ctx, mcpmDir, "npm", "install")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run npm install in %s: %v\nOutput: %s", mcpmDir, err, output)), nil
	}

	// 2. Link MCPM
	output, err = h.Cmd.RunInDir(ctx, mcpmDir, "npm", "link")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run npm link in %s: %v\nOutput: %s", mcpmDir, err, output)), nil
	}

	// 3. Install Default Servers (The Guardian Stack)
	defaultServers := []string{"context7", "brave-search", "github"}
	var warnings []string
	for _, server := range defaultServers {
		if _, err := h.Mcpm.Run("install", server); err != nil {
			warnings = append(warnings, fmt.Sprintf("Warning: Failed to install default server %s: %v", server, err))
		}
	}

	// 4. Start Infrastructure
	scriptPath := filepath.Join(rootDir, "scripts", "manage-mcp.sh")
	if _, err := h.FS.Stat(scriptPath); err == nil {
		output, err = h.Cmd.Run(ctx, scriptPath, "start")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start infrastructure via script: %v\nOutput: %s", err, output)), nil
		}
	} else {
		output, err = h.Cmd.RunInDir(ctx, rootDir, "docker", "compose", "up", "-d")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to run docker compose up in %s: %v\nOutput: %s", rootDir, err, output)), nil
		}
	}

	result := "System bootstrapped successfully! MCPM installed, default servers (context7, brave-search, github) set up, and Infrastructure started."
	if len(warnings) > 0 {
		result += "\n\n" + strings.Join(warnings, "\n")
	}

	return mcp.NewToolResultText(result), nil
}

// RestartService handles the restart_service tool
func (h *Handler) RestartService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	go func() {
		time.Sleep(1 * time.Second)
		h.ExitProcess(0)
	}()
	return mcp.NewToolResultText("Restarting Jarvis service..."), nil
}

// RestartInfrastructure handles the restart_infrastructure tool
func (h *Handler) RestartInfrastructure(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cwd, err := h.FS.Getwd()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CWD: %v", err)), nil
	}

	// Find project root
	var rootDir string
	if _, err := h.FS.Stat(filepath.Join(cwd, "MCPM")); err == nil {
		rootDir = cwd
	} else if _, err := h.FS.Stat(filepath.Join(cwd, "..", "MCPM")); err == nil {
		rootDir = filepath.Join(cwd, "..")
	} else {
		return mcp.NewToolResultError("Could not locate project root."), nil
	}

	scriptPath := filepath.Join(rootDir, "scripts", "manage-mcp.sh")
	if _, err := h.FS.Stat(scriptPath); os.IsNotExist(err) {
		return mcp.NewToolResultError("Management script not found at " + scriptPath), nil
	}

	output, err := h.Cmd.Run(ctx, scriptPath, "restart")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restart failed: %v\nOutput: %s", err, output)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Infrastructure restarted successfully.\nOutput:\n%s", output)), nil
}

// ShareServer handles the share_server tool
func (h *Handler) ShareServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	// Check if already shared
	if _, exists := h.Processes.Get(name); exists {
		return mcp.NewToolResultError(fmt.Sprintf("Server %s is already being shared", name)), nil
	}

	cmdArgs := []string{"share", name}
	if port, ok := args["port"].(string); ok && port != "" {
		cmdArgs = append(cmdArgs, "--port", port)
	}
	if noAuth, ok := args["no_auth"].(bool); ok && noAuth {
		cmdArgs = append(cmdArgs, "--no-auth")
	}

	// Start mcpm share in background
	proc, err := h.Cmd.StartBackground(ctx, "mcpm", cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start share command: %v", err)), nil
	}

	// Register the process
	h.Processes.Register(name, proc)

	// Monitor for URL or failure
	success := make(chan string)
	failure := make(chan string)

	go h.monitorShareProcess(proc.Stdout(), success, failure)

	select {
	case output := <-success:
		return mcp.NewToolResultText(output), nil
	case errStr := <-failure:
		_ = proc.Kill()
		h.Processes.Remove(name)
		return mcp.NewToolResultError(errStr), nil
	case <-time.After(30 * time.Second):
		_ = proc.Kill()
		h.Processes.Remove(name)
		return mcp.NewToolResultError("Timeout waiting for share URL"), nil
	}
}

func (h *Handler) monitorShareProcess(stdout io.Reader, success, failure chan<- string) {
	scanner := bufio.NewScanner(stdout)
	var output strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")

		// Look for URL pattern
		if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
			success <- output.String()
			return
		}

		// Look for error patterns
		if strings.Contains(strings.ToLower(line), "error") || strings.Contains(strings.ToLower(line), "failed") {
			failure <- output.String()
			return
		}
	}

	// If we reach here without finding URL, it's a failure
	if output.Len() > 0 {
		failure <- output.String()
	} else {
		failure <- "No output received from share command"
	}
}

// StopSharingServer handles the stop_sharing_server tool
func (h *Handler) StopSharingServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	name, ok := args["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return mcp.NewToolResultError("name argument is required"), nil
	}

	proc, exists := h.Processes.Get(name)
	if !exists {
		return mcp.NewToolResultError(fmt.Sprintf("Server %s is not currently shared", name)), nil
	}

	h.Processes.Remove(name)

	if err := proc.Kill(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stop sharing server %s: %v", name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Stopped sharing server %s", name)), nil
}

// ListSharedServers handles the list_shared_servers tool
func (h *Handler) ListSharedServers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	names := h.Processes.List()

	if len(names) == 0 {
		return mcp.NewToolResultText("No servers are currently being shared."), nil
	}

	var builder strings.Builder
	builder.WriteString("Currently shared servers:\n")
	for _, name := range names {
		builder.WriteString(fmt.Sprintf("- %s\n", name))
	}

	return mcp.NewToolResultText(builder.String()), nil
}
