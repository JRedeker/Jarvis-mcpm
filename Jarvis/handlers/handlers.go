// Package handlers provides MCP tool handlers with dependency injection for testing
package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// NewMcpmRunner creates an McpmRunner based on the JARVIS_MCPM_TRANSPORT environment variable
// Supported values: "http" (default), "cli"
// If "http" is selected but the API server is not reachable, it falls back to CLI
func NewMcpmRunner() McpmRunner {
	transport := os.Getenv("JARVIS_MCPM_TRANSPORT")
	apiURL := os.Getenv("MCPM_API_URL")

	// Default to HTTP transport
	if transport == "" {
		transport = "http"
	}

	switch transport {
	case "http":
		runner := NewHTTPMcpmRunner(apiURL)
		// Test connectivity - if API is not available, fall back to CLI
		if _, err := runner.Run("doctor"); err != nil {
			// API not available, fall back to CLI
			return &RealMcpmRunner{}
		}
		return runner
	case "cli":
		return &RealMcpmRunner{}
	default:
		// Unknown transport, use CLI as fallback
		return &RealMcpmRunner{}
	}
}

// NewDefaultHandler creates a Handler with default dependencies, auto-selecting MCPM transport
func NewDefaultHandler() *Handler {
	return NewHandler(
		NewMcpmRunner(),
		&RealDockerRunner{},
		&RealGitRunner{},
		&RealFileSystem{},
	)
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

	// Check MCPM API Server health
	apiStatus := checkMcpmAPIHealth()
	if apiStatus != "" {
		output += "\n\n## 🌐 MCPM API Server\n"
		output += apiStatus
	}

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

// checkMcpmAPIHealth checks if the MCPM API server is running
func checkMcpmAPIHealth() string {
	apiURL := os.Getenv("MCPM_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:6275"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL + "/api/v1/health")
	if err != nil {
		return fmt.Sprintf("- Status: ❌ Not reachable (%s)\n- URL: %s", err.Error(), apiURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return fmt.Sprintf("- Status: ✅ Healthy\n- URL: %s/api/v1", apiURL)
	}
	return fmt.Sprintf("- Status: ⚠️ Unhealthy (HTTP %d)\n- URL: %s", resp.StatusCode, apiURL)
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
		return mcp.NewToolResultError("❌ Action is required and cannot be empty.\n\n💡 Valid actions: ls, edit, import, config\n\nExamples:\n- manage_client(\"ls\") - List all configured clients\n- manage_client(\"edit\", client_name=\"opencode\", add_profile=\"memory\")"), nil
	}

	// Validate action is one of the allowed values
	validActions := map[string]bool{"ls": true, "edit": true, "import": true, "config": true}
	if !validActions[action] {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Invalid action '%s'\n\n💡 Valid actions: ls, edit, import, config\n\n- ls: List all configured clients\n- edit: Add/remove servers or profiles from a client\n- import: Import client configuration\n- config: Get or set client config path", action)), nil
	}

	// Handle ls action - list all known clients
	if action == "ls" {
		return h.listClients(ctx)
	}

	// For other actions, client_name is required
	clientName, ok := args["client_name"].(string)
	if !ok || strings.TrimSpace(clientName) == "" {
		return mcp.NewToolResultError(fmt.Sprintf("❌ client_name is required for action '%s'\n\n💡 Tip: Use manage_client(\"ls\") to see available client names\n\nSupported clients: opencode, claude-code, claude-desktop", action)), nil
	}

	// Handle OpenCode specifically with native support
	if clientName == "opencode" {
		return h.manageOpenCodeClient(ctx, action, args)
	}

	// Fall back to MCPM CLI for other clients
	cmdArgs := []string{"client", action, clientName}

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

// listClients lists all known MCP clients and their detection status
func (h *Handler) listClients(ctx context.Context) (*mcp.CallToolResult, error) {
	var sb strings.Builder
	sb.WriteString("## 📱 MCP Clients\n\n")

	for name, client := range KnownClients {
		detected := "❌ Not detected"
		configPath := ""

		// Try to detect the client
		if name == "opencode" {
			if path, err := DetectOpenCodeConfig(h.FS); err == nil {
				detected = "✅ Detected"
				configPath = path
			}
		} else {
			// Check standard paths for other clients
			for _, p := range client.ConfigPaths {
				expandedPath := expandPath(p)
				if _, err := h.FS.Stat(expandedPath); err == nil {
					detected = "✅ Detected"
					configPath = expandedPath
					break
				}
			}
		}

		sb.WriteString(fmt.Sprintf("### %s\n", client.DisplayName))
		sb.WriteString(fmt.Sprintf("- Status: %s\n", detected))
		if configPath != "" {
			sb.WriteString(fmt.Sprintf("- Config: `%s`\n", configPath))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("💡 Use `manage_client(\"config\", client_name=\"opencode\", config_path=\"/path/to/config\")` to set a custom config path.\n")

	return mcp.NewToolResultText(sb.String()), nil
}

// manageOpenCodeClient handles OpenCode-specific client management
func (h *Handler) manageOpenCodeClient(ctx context.Context, action string, args map[string]interface{}) (*mcp.CallToolResult, error) {
	switch action {
	case "config":
		return h.openCodeConfig(ctx, args)
	case "edit":
		return h.openCodeEdit(ctx, args)
	case "import":
		return h.openCodeImport(ctx, args)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown action '%s' for OpenCode", action)), nil
	}
}

// openCodeConfig handles getting/setting OpenCode config path
func (h *Handler) openCodeConfig(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	if configPath, ok := args["config_path"].(string); ok && configPath != "" {
		// Validate the path exists or can be created
		expandedPath := expandPath(configPath)
		dir := filepath.Dir(expandedPath)
		if err := h.FS.MkdirAll(dir, 0755); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Cannot create config directory: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("✅ OpenCode config path set to: `%s`\n\n💡 Use `manage_client(\"edit\", client_name=\"opencode\", add_profile=\"jarvis\")` to configure MCP servers.", expandedPath)), nil
	}

	// Get current config path
	path, err := DetectOpenCodeConfig(h.FS)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("❌ %v\n\n💡 Create a config file at one of the expected locations, or use config_path to specify a custom location.", err)), nil
	}

	// List configured servers
	serverList, err := ListOpenCodeServers(h.FS, path)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("✅ OpenCode config found at: `%s`\n\n⚠️ Could not read servers: %v", path, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ OpenCode config: `%s`\n\n%s", path, serverList)), nil
}

// openCodeEdit handles adding/removing profiles from OpenCode config
func (h *Handler) openCodeEdit(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	// Detect config path
	configPath, err := DetectOpenCodeConfig(h.FS)
	if err != nil {
		// If not found, use the global default
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Cannot determine config path: %v", err)), nil
		}
		configPath = filepath.Join(home, ".config", "opencode", "opencode.json")
	}

	// Override with explicit config_path if provided
	if path, ok := args["config_path"].(string); ok && path != "" {
		configPath = expandPath(path)
	}

	var results []string

	// Handle add_profile
	if addProfile, ok := args["add_profile"].(string); ok && addProfile != "" {
		profiles := strings.Split(addProfile, ",")
		for _, profile := range profiles {
			profile = strings.TrimSpace(profile)
			if profile == "" {
				continue
			}

			// Get jarvis path for jarvis profile
			jarvisPath := ""
			if profile == "jarvis" {
				// Try to find Jarvis binary
				cwd, _ := h.FS.Getwd()
				possiblePaths := []string{
					filepath.Join(cwd, "Jarvis", "jarvis"),
					filepath.Join(cwd, "..", "Jarvis", "jarvis"),
					filepath.Join(cwd, "jarvis"),
				}
				for _, p := range possiblePaths {
					if _, err := h.FS.Stat(p); err == nil {
						jarvisPath = p
						break
					}
				}
				if jarvisPath == "" {
					results = append(results, fmt.Sprintf("⚠️ Could not find Jarvis binary. Please provide the path manually."))
					continue
				}
			}

			if err := AddProfileToOpenCode(h.FS, configPath, profile, jarvisPath); err != nil {
				results = append(results, fmt.Sprintf("❌ Failed to add %s: %v", profile, err))
			} else {
				results = append(results, fmt.Sprintf("✅ Added %s", profile))
			}
		}
	}

	// Handle remove_profile
	if removeProfile, ok := args["remove_profile"].(string); ok && removeProfile != "" {
		profiles := strings.Split(removeProfile, ",")
		for _, profile := range profiles {
			profile = strings.TrimSpace(profile)
			if profile == "" {
				continue
			}

			if err := RemoveProfileFromOpenCode(h.FS, configPath, profile); err != nil {
				results = append(results, fmt.Sprintf("❌ Failed to remove %s: %v", profile, err))
			} else {
				results = append(results, fmt.Sprintf("✅ Removed %s", profile))
			}
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultError("❌ No changes specified. Use add_profile or remove_profile.\n\n💡 Example: manage_client(\"edit\", client_name=\"opencode\", add_profile=\"jarvis,memory,toolbox\")"), nil
	}

	output := fmt.Sprintf("## OpenCode Configuration Updated\n\nConfig: `%s`\n\n%s\n\n💡 Restart OpenCode to apply changes.", configPath, strings.Join(results, "\n"))
	return mcp.NewToolResultText(output), nil
}

// openCodeImport handles importing a template configuration
func (h *Handler) openCodeImport(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	// Detect or create config path
	configPath, err := DetectOpenCodeConfig(h.FS)
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Cannot determine config path: %v", err)), nil
		}
		configPath = filepath.Join(home, ".config", "opencode", "opencode.json")
	}

	// Override with explicit config_path if provided
	if path, ok := args["config_path"].(string); ok && path != "" {
		configPath = expandPath(path)
	}

	// Find Jarvis binary
	jarvisPath := ""
	cwd, _ := h.FS.Getwd()
	possiblePaths := []string{
		filepath.Join(cwd, "Jarvis", "jarvis"),
		filepath.Join(cwd, "..", "Jarvis", "jarvis"),
		filepath.Join(cwd, "jarvis"),
	}
	for _, p := range possiblePaths {
		if _, err := h.FS.Stat(p); err == nil {
			jarvisPath = p
			break
		}
	}

	if jarvisPath == "" {
		return mcp.NewToolResultError("❌ Could not find Jarvis binary. Please build Jarvis first:\n\n```bash\ncd Jarvis && go build -o jarvis .\n```"), nil
	}

	// Generate template
	template := GenerateOpenCodeTemplate(jarvisPath)

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := h.FS.MkdirAll(dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Failed to create config directory: %v", err)), nil
	}

	// Write config
	if err := h.FS.WriteFile(configPath, []byte(template), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Failed to write config: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ OpenCode configuration imported!\n\nConfig: `%s`\n\nConfigured servers:\n- jarvis (local)\n- toolbox (http://localhost:6276/mcp)\n- memory (http://localhost:6277/mcp)\n- morph (http://localhost:6278/mcp)\n\n💡 Make sure Docker infrastructure is running: `./scripts/manage-mcp.sh start`", configPath)), nil
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
		profiles = append(profiles, "toolbox")
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

// =============================================================================
// Diagnostic Handlers - Essential for AI agents debugging MCP issues
// =============================================================================

// DiagnoseProfileHealth checks if MCP profiles are healthy and serving tools
func (h *Handler) DiagnoseProfileHealth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	profileName, _ := args["profile"].(string)

	var builder strings.Builder
	builder.WriteString("## MCP Profile Health Report\n\n")

	// Get supervisor status from mcp-daemon
	output, err := h.Docker.ExecSupervisorctl(ctx, "status", "all")
	if err != nil {
		builder.WriteString("### ⚠️ Daemon Status: Error\n")
		builder.WriteString(fmt.Sprintf("Could not reach mcp-daemon: %v\n", err))
		builder.WriteString("\n**Suggestion:** Run `docker ps` to check if mcp-daemon is running.\n")
		return mcp.NewToolResultText(builder.String()), nil
	}

	builder.WriteString("### Supervisor Status\n```\n")
	builder.WriteString(output)
	builder.WriteString("```\n\n")

	// Parse supervisor output to identify issues
	lines := strings.Split(output, "\n")
	var runningProfiles []string
	var failedProfiles []string

	for _, line := range lines {
		if strings.Contains(line, "mcpm-") {
			if strings.Contains(line, "RUNNING") {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					name := strings.TrimPrefix(parts[0], "mcpm-")
					runningProfiles = append(runningProfiles, name)
				}
			} else if strings.Contains(line, "FATAL") || strings.Contains(line, "STOPPED") || strings.Contains(line, "EXITED") {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					name := strings.TrimPrefix(parts[0], "mcpm-")
					failedProfiles = append(failedProfiles, name)
				}
			}
		}
	}

	builder.WriteString("### Summary\n")
	builder.WriteString(fmt.Sprintf("- ✅ Running profiles: %d (%s)\n", len(runningProfiles), strings.Join(runningProfiles, ", ")))
	if len(failedProfiles) > 0 {
		builder.WriteString(fmt.Sprintf("- ❌ Failed profiles: %d (%s)\n", len(failedProfiles), strings.Join(failedProfiles, ", ")))
		builder.WriteString("\n**Next step:** Use `jarvis_diagnose(action=\"logs\", profile=\"<name>\")` to see error details.\n")
	}

	// If specific profile requested, check its endpoint
	if profileName != "" {
		builder.WriteString(fmt.Sprintf("\n### Profile '%s' Details\n", profileName))

		// Get port mapping
		portMap := map[string]string{
			"essentials": "6276",
			"memory":     "6277",
			"dev-core":   "6278",
			"data":       "6279",
			"p-new":      "6280",
			"research":   "6281",
		}

		if port, ok := portMap[profileName]; ok {
			endpoint := fmt.Sprintf("http://localhost:%s/mcp", port)
			builder.WriteString(fmt.Sprintf("- Endpoint: %s\n", endpoint))
			builder.WriteString(fmt.Sprintf("- Port: %s\n", port))
			builder.WriteString(fmt.Sprintf("\n**Test with:** `jarvis_diagnose(action=\"test_endpoint\", endpoint=\"%s\")`\n", endpoint))
		}
	}

	return mcp.NewToolResultText(builder.String()), nil
}

// DiagnoseTestEndpoint tests an MCP endpoint and reports tool availability
func (h *Handler) DiagnoseTestEndpoint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	endpoint, _ := args["endpoint"].(string)

	if endpoint == "" {
		return mcp.NewToolResultError("endpoint is required. Example: 'http://localhost:6276/mcp'"), nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("## MCP Endpoint Test: %s\n\n", endpoint))

	// Test HTTP connectivity first
	client := &http.Client{Timeout: 10 * time.Second}

	// Check health endpoint
	healthURL := strings.Replace(endpoint, "/mcp", "/health", 1)
	resp, err := client.Get(healthURL)
	if err != nil {
		builder.WriteString("### ❌ Connection Failed\n")
		builder.WriteString(fmt.Sprintf("Error: %v\n", err))
		builder.WriteString("\n**Suggestions:**\n")
		builder.WriteString("1. Check if mcp-daemon container is running\n")
		builder.WriteString("2. Verify port is not blocked\n")
		builder.WriteString("3. Run `jarvis_diagnose(action=\"profile_health\")` for supervisor status\n")
		return mcp.NewToolResultText(builder.String()), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		builder.WriteString("### ✅ Health Check Passed\n")
		body, _ := io.ReadAll(resp.Body)
		builder.WriteString(fmt.Sprintf("Response: %s\n\n", string(body)))
	} else {
		builder.WriteString(fmt.Sprintf("### ⚠️ Health Check: HTTP %d\n\n", resp.StatusCode))
	}

	// Test MCP initialize
	builder.WriteString("### MCP Protocol Test\n")
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"jarvis-diagnose","version":"1.0"}}}`

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(initReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp2, err := client.Do(req)
	if err != nil {
		builder.WriteString(fmt.Sprintf("❌ MCP initialize failed: %v\n", err))
		return mcp.NewToolResultText(builder.String()), nil
	}
	defer resp2.Body.Close()

	sessionID := resp2.Header.Get("mcp-session-id")
	if sessionID != "" {
		builder.WriteString(fmt.Sprintf("✅ Session established: %s\n", sessionID))
	}

	body, _ := io.ReadAll(resp2.Body)
	bodyStr := string(body)

	type initResult struct {
		Result struct {
			ServerInfo struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"serverInfo"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	var initResp initResult
	if err := json.Unmarshal(body, &initResp); err == nil && initResp.Result.ServerInfo.Name != "" {
		builder.WriteString("✅ MCP initialize successful\n")
		serverName := initResp.Result.ServerInfo.Name
		if strings.HasPrefix(serverName, "profile-") {
			serverName = strings.TrimPrefix(serverName, "profile-")
		}
		builder.WriteString(fmt.Sprintf("Server: %s\n", serverName))
	} else if strings.Contains(bodyStr, "serverInfo") {
		builder.WriteString("✅ MCP initialize successful\n")
	} else if strings.Contains(bodyStr, "error") || initResp.Error != nil {
		builder.WriteString("❌ MCP initialize returned error\n")
		builder.WriteString(fmt.Sprintf("Response: %s\n", bodyStr))
	}

	// Now test tools/list
	if sessionID != "" {
		builder.WriteString("\n### Tools Available\n")
		toolsReq := `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`

		req2, _ := http.NewRequest("POST", endpoint, strings.NewReader(toolsReq))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/json, text/event-stream")
		req2.Header.Set("mcp-session-id", sessionID)

		resp3, err := client.Do(req2)
		if err != nil {
			builder.WriteString(fmt.Sprintf("❌ tools/list failed: %v\n", err))
		} else {
			defer resp3.Body.Close()
			body3, _ := io.ReadAll(resp3.Body)
			bodyStr3 := string(body3)

			type toolsResult struct {
				Result struct {
					Tools []struct {
						Name string `json:"name"`
					} `json:"tools"`
				} `json:"result"`
				Error interface{} `json:"error"`
			}

			var toolsResp toolsResult
			if err := json.Unmarshal(body3, &toolsResp); err == nil && len(toolsResp.Result.Tools) > 0 {
				builder.WriteString(fmt.Sprintf("✅ Found %d tools\n", len(toolsResp.Result.Tools)))
				builder.WriteString("\nTools:\n")
				const maxTools = 25
				for i, tool := range toolsResp.Result.Tools {
					if i >= maxTools {
						remaining := len(toolsResp.Result.Tools) - maxTools
						builder.WriteString(fmt.Sprintf("- ...and %d more\n", remaining))
						break
					}
					builder.WriteString(fmt.Sprintf("- %s\n", tool.Name))
				}
			} else if strings.Contains(bodyStr3, "error") || toolsResp.Error != nil {
				builder.WriteString("❌ tools/list returned error\n")
				builder.WriteString(fmt.Sprintf("Error: %s\n", bodyStr3))
				builder.WriteString("\n**This usually means a subprocess failed.**\n")
				builder.WriteString("Run `jarvis_diagnose(action=\"logs\")` to see subprocess errors.\n")
			} else {
				toolCount := strings.Count(bodyStr3, `"name":`) - 1 // -1 for the method name
				if toolCount < 0 {
					toolCount = 0
				}
				if toolCount > 0 {
					builder.WriteString(fmt.Sprintf("✅ Found %d tools\n", toolCount))
				}
			}
		}
	}

	return mcp.NewToolResultText(builder.String()), nil
}

// DiagnoseLogs retrieves subprocess logs from the mcp-daemon
func (h *Handler) DiagnoseLogs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := request.Params.Arguments.(map[string]interface{})
	profileName, _ := args["profile"].(string)
	lines := 50
	if l, ok := args["lines"].(float64); ok && l > 0 {
		lines = int(l)
	}

	var builder strings.Builder
	builder.WriteString("## MCP Subprocess Logs\n\n")

	if profileName == "" {
		// List available profiles
		builder.WriteString("**Available profiles:** essentials, memory, dev-core, research, data, p-new\n\n")
		builder.WriteString("Specify a profile with `jarvis_diagnose(action=\"logs\", profile=\"<name>\")`\n")
		return mcp.NewToolResultText(builder.String()), nil
	}

	// Get stderr logs (where errors appear)
	target := fmt.Sprintf("mcpm-%s", profileName)
	output, err := h.Docker.ExecSupervisorctl(ctx, "tail", fmt.Sprintf("-f %s stderr", target))

	// ExecSupervisorctl doesn't support tail -f well, use docker exec directly
	cmd := fmt.Sprintf("docker exec mcp-daemon tail -%d /var/log/mcpm/%s.err.log 2>/dev/null || docker exec mcp-daemon supervisorctl tail %s stderr", lines, profileName, target)
	output2, _ := exec.CommandContext(ctx, "sh", "-c", cmd).CombinedOutput()

	if len(output2) > 0 {
		output = string(output2)
	}

	builder.WriteString(fmt.Sprintf("### Profile: %s (stderr)\n", profileName))
	builder.WriteString("```\n")
	if err != nil && len(output) == 0 {
		builder.WriteString(fmt.Sprintf("Error retrieving logs: %v\n", err))
	} else if len(output) == 0 {
		builder.WriteString("(no stderr output)\n")
	} else {
		// Limit output to last N lines
		logLines := strings.Split(strings.TrimSpace(output), "\n")
		if len(logLines) > lines {
			logLines = logLines[len(logLines)-lines:]
		}
		builder.WriteString(strings.Join(logLines, "\n"))
	}
	builder.WriteString("\n```\n")

	// Look for common error patterns and provide suggestions
	if strings.Contains(output, "ValueError") || strings.Contains(output, "ImportError") {
		builder.WriteString("\n### ⚠️ Python Error Detected\n")
		builder.WriteString("The subprocess crashed due to a Python error.\n")
		builder.WriteString("**Common causes:**\n")
		builder.WriteString("- Missing environment variables\n")
		builder.WriteString("- Incorrect configuration in servers.json\n")
		builder.WriteString("- Incompatible package versions\n")
	}

	if strings.Contains(output, "Connection refused") || strings.Contains(output, "ECONNREFUSED") {
		builder.WriteString("\n### ⚠️ Connection Error Detected\n")
		builder.WriteString("A subprocess couldn't connect to a dependency.\n")
		builder.WriteString("**Check:**\n")
		builder.WriteString("- Is the target service running? (qdrant, postgres, etc.)\n")
		builder.WriteString("- Is the URL correct in servers.json?\n")
	}

	if strings.Contains(output, "Multiple location") || strings.Contains(output, "Only one of") {
		builder.WriteString("\n### ⚠️ Configuration Conflict Detected\n")
		builder.WriteString("Multiple conflicting options were specified.\n")
		builder.WriteString("**Fix:** Edit ~/.config/mcpm/servers.json and remove conflicting options.\n")
	}

	return mcp.NewToolResultText(builder.String()), nil
}

// DiagnoseFull runs all diagnostics and provides a comprehensive report
func (h *Handler) DiagnoseFull(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var builder strings.Builder
	builder.WriteString("# Full MCP Diagnostic Report\n\n")
	builder.WriteString("Generated by Jarvis at " + time.Now().Format(time.RFC3339) + "\n\n")

	// 1. Profile Health
	builder.WriteString("---\n")
	healthResult, _ := h.DiagnoseProfileHealth(ctx, request)
	if healthResult != nil {
		for _, content := range healthResult.Content {
			if textContent, ok := content.(mcp.TextContent); ok {
				builder.WriteString(textContent.Text)
			}
		}
	}

	// 2. Test each standard endpoint
	builder.WriteString("\n---\n")
	builder.WriteString("## Endpoint Tests\n\n")

	endpoints := map[string]string{
		"essentials": "http://localhost:6276/mcp",
		"memory":     "http://localhost:6277/mcp",
		"dev-core":   "http://localhost:6278/mcp",
		"data":       "http://localhost:6279/mcp",
		"research":   "http://localhost:6281/mcp",
	}

	for name, url := range endpoints {
		builder.WriteString(fmt.Sprintf("### %s\n", name))

		// Quick health check
		client := &http.Client{Timeout: 5 * time.Second}
		healthURL := strings.Replace(url, "/mcp", "/health", 1)
		resp, err := client.Get(healthURL)
		if err != nil {
			builder.WriteString(fmt.Sprintf("❌ Unreachable: %v\n\n", err))
		} else {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				builder.WriteString("✅ Healthy\n\n")
			} else {
				builder.WriteString(fmt.Sprintf("⚠️ HTTP %d\n\n", resp.StatusCode))
			}
		}
	}

	// 3. Configuration check
	builder.WriteString("---\n")
	builder.WriteString("## Configuration\n\n")
	builder.WriteString("- MCPM config: ~/.config/mcpm/servers.json\n")
	builder.WriteString("- OpenCode config: ~/.config/opencode/opencode.json\n")
	builder.WriteString("\n**To update configs after changes:**\n")
	builder.WriteString("```\njarvis_profile(action=\"restart\", profile=\"<name>\")\n```\n")

	return mcp.NewToolResultText(builder.String()), nil
}
