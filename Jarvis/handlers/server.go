package handlers

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolDefinition contains the MCP tool definition and its handler
type ToolDefinition struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

// GetToolDefinitions returns all Jarvis tool definitions with their handlers
// This is the production wiring - it creates real dependencies
func GetToolDefinitions(h *Handler) []ToolDefinition {
	return []ToolDefinition{
		// System Management
		{
			Tool: mcp.NewTool("check_status",
				mcp.WithDescription("Comprehensive system health check for MCPM, Docker, and all services. Validates Node.js, Python, dependencies, running containers, and HTTP endpoints. Returns actionable fix suggestions for any issues found."),
			),
			Handler: h.CheckStatus,
		},

		// Server Management
		{
			Tool: mcp.NewTool("list_servers",
				mcp.WithDescription("Shows all installed MCP servers with their status, transport type, and profile associations. Use this to inventory available tools before making changes."),
			),
			Handler: h.ListServers,
		},
		{
			Tool: mcp.NewTool("server_info",
				mcp.WithDescription("Detailed information about a specific server including command, args, environment variables, installation source, and usage statistics. Essential before editing or troubleshooting a server."),
				mcp.WithString("name",
					mcp.Description("Server name to get info for"),
					mcp.Required(),
				),
			),
			Handler: h.ServerInfo,
		},
		{
			Tool: mcp.NewTool("install_server",
				mcp.WithDescription("Installs an MCP server from the registry with automatic dependency resolution. Validates the server exists before installing and suggests alternatives for typos."),
				mcp.WithString("name",
					mcp.Description("Server name to install (e.g., 'context7', 'brave-search')"),
					mcp.Required(),
				),
			),
			Handler: h.InstallServer,
		},
		{
			Tool: mcp.NewTool("uninstall_server",
				mcp.WithDescription("Removes an installed MCP server and cleans up its configuration. Warns about profile impact before removal."),
				mcp.WithString("name",
					mcp.Description("Server name to uninstall"),
					mcp.Required(),
				),
			),
			Handler: h.UninstallServer,
		},
		{
			Tool: mcp.NewTool("search_servers",
				mcp.WithDescription("Search the MCP server registry by keyword, category, or capability. Returns matching servers with descriptions and installation commands."),
				mcp.WithString("query",
					mcp.Description("Search query (e.g., 'memory', 'documentation', 'database')"),
					mcp.Required(),
				),
			),
			Handler: h.SearchServers,
		},
		{
			Tool: mcp.NewTool("edit_server",
				mcp.WithDescription("Modify an installed server's configuration including command, arguments, environment variables, and URL."),
				mcp.WithString("name",
					mcp.Description("Server name to edit"),
					mcp.Required(),
				),
				mcp.WithString("command",
					mcp.Description("New command to run the server"),
				),
				mcp.WithString("args",
					mcp.Description("New arguments (space-separated)"),
				),
				mcp.WithString("env",
					mcp.Description("Environment variables (KEY=value format, comma-separated)"),
				),
				mcp.WithString("url",
					mcp.Description("New URL for HTTP servers"),
				),
				mcp.WithString("headers",
					mcp.Description("HTTP headers (KEY=value format, comma-separated)"),
				),
			),
			Handler: h.EditServer,
		},
		{
			Tool: mcp.NewTool("create_server",
				mcp.WithDescription("Register a new custom MCP server that isn't in the registry."),
				mcp.WithString("name",
					mcp.Description("Unique server name"),
					mcp.Required(),
				),
				mcp.WithString("type",
					mcp.Description("Transport type: 'stdio' or 'streamable-http'"),
					mcp.Required(),
				),
				mcp.WithString("command",
					mcp.Description("Command to run (for stdio type)"),
				),
				mcp.WithString("args",
					mcp.Description("Arguments (space-separated)"),
				),
				mcp.WithString("env",
					mcp.Description("Environment variables (KEY=value, comma-separated)"),
				),
				mcp.WithString("url",
					mcp.Description("URL (for streamable-http type)"),
				),
				mcp.WithString("headers",
					mcp.Description("HTTP headers (KEY=value, comma-separated)"),
				),
			),
			Handler: h.CreateServer,
		},
		{
			Tool: mcp.NewTool("usage_stats",
				mcp.WithDescription("Shows tool usage statistics across all servers and profiles."),
			),
			Handler: h.UsageStats,
		},

		// Profile Management
		{
			Tool: mcp.NewTool("manage_profile",
				mcp.WithDescription("Create, edit, delete, or list MCPM profiles. Profiles group servers for specific use cases."),
				mcp.WithString("action",
					mcp.Description("Action: 'ls', 'create', 'edit', 'delete'"),
					mcp.Required(),
				),
				mcp.WithString("name",
					mcp.Description("Profile name"),
				),
				mcp.WithString("new_name",
					mcp.Description("New name when renaming"),
				),
				mcp.WithString("add_servers",
					mcp.Description("Servers to add (comma-separated)"),
				),
				mcp.WithString("remove_servers",
					mcp.Description("Servers to remove (comma-separated)"),
				),
			),
			Handler: h.ManageProfile,
		},
		{
			Tool: mcp.NewTool("suggest_profile",
				mcp.WithDescription("Intelligently determines optimal MCPM profile stack by analyzing working directory, client type, and mode."),
				mcp.WithBoolean("testing",
					mcp.Description("Include testing-all-tools profile"),
				),
			),
			Handler: h.SuggestProfile,
		},
		{
			Tool: mcp.NewTool("restart_profiles",
				mcp.WithDescription("Restarts the MCPM daemon container to reload all MCP profiles with updated configurations."),
				mcp.WithString("profile",
					mcp.Description("Optional: restart only a specific profile"),
				),
			),
			Handler: h.RestartProfiles,
		},

		// Client Management
		{
			Tool: mcp.NewTool("manage_client",
				mcp.WithDescription("Configure MCP client applications (OpenCode, Claude Code, Claude Desktop) with profiles and servers. OpenCode has native support with automatic config detection."),
				mcp.WithString("action",
					mcp.Description("Action: 'ls', 'edit', 'import', 'config'"),
					mcp.Required(),
				),
				mcp.WithString("client_name",
					mcp.Description("Client to configure (e.g., 'opencode', 'claude-code', 'claude-desktop')"),
				),
				mcp.WithString("add_server",
					mcp.Description("Server to add directly to client"),
				),
				mcp.WithString("remove_server",
					mcp.Description("Server to remove from client"),
				),
				mcp.WithString("add_profile",
					mcp.Description("Profile to add to client"),
				),
				mcp.WithString("remove_profile",
					mcp.Description("Profile to remove from client"),
				),
				mcp.WithString("config_path",
					mcp.Description("Config file path (for 'config' action)"),
				),
			),
			Handler: h.ManageClient,
		},

		// Configuration
		{
			Tool: mcp.NewTool("manage_config",
				mcp.WithDescription("Get, set, or list MCPM configuration values."),
				mcp.WithString("action",
					mcp.Description("Action: 'get', 'set', 'ls'"),
					mcp.Required(),
				),
				mcp.WithString("key",
					mcp.Description("Config key"),
				),
				mcp.WithString("value",
					mcp.Description("Value to set"),
				),
			),
			Handler: h.ManageConfig,
		},
		{
			Tool: mcp.NewTool("migrate_config",
				mcp.WithDescription("Upgrades MCPM configuration to the latest format with automatic backup."),
			),
			Handler: h.MigrateConfig,
		},

		// Project Analysis
		{
			Tool: mcp.NewTool("analyze_project",
				mcp.WithDescription("Analyzes the current project to detect languages, frameworks, and existing DevOps configurations. Returns JSON report."),
			),
			Handler: h.AnalyzeProject,
		},
		{
			Tool: mcp.NewTool("fetch_diff_context",
				mcp.WithDescription("Returns git status and diff for self-review before commits. Helps catch issues before pushing."),
				mcp.WithBoolean("staged",
					mcp.Description("Show only staged changes (default: all changes)"),
				),
			),
			Handler: h.FetchDiffContext,
		},
		{
			Tool: mcp.NewTool("apply_devops_stack",
				mcp.WithDescription("Scaffolds projects with linting, pre-commit hooks, and CI/CD workflows based on detected project type."),
				mcp.WithString("project_type",
					mcp.Description("Override detected type: 'python', 'go', 'node', 'typescript'"),
				),
				mcp.WithBoolean("force",
					mcp.Description("Overwrite existing configuration files"),
				),
				mcp.WithBoolean("enable_ai_review",
					mcp.Description("Add PR Agent workflow for AI code review"),
				),
			),
			Handler: h.ApplyDevOpsStack,
		},

		// System Bootstrap & Infrastructure
		{
			Tool: mcp.NewTool("bootstrap_system",
				mcp.WithDescription("Complete system initialization: installs MCPM, sets up default servers (context7, brave-search, github), and starts Docker infrastructure (PostgreSQL, Qdrant). One command to get fully operational."),
			),
			Handler: h.BootstrapSystem,
		},
		{
			Tool: mcp.NewTool("restart_service",
				mcp.WithDescription("Gracefully restarts Jarvis to apply configuration changes or resolve stuck states. Automatically saves state and reconnects active sessions."),
			),
			Handler: h.RestartService,
		},
		{
			Tool: mcp.NewTool("restart_infrastructure",
				mcp.WithDescription("Safely reboots Docker infrastructure (PostgreSQL, Qdrant) with health checks and automatic reconnection. Resolves database connection issues, clears stale locks, and ensures all services are healthy."),
			),
			Handler: h.RestartInfrastructure,
		},

		// Server Sharing
		{
			Tool: mcp.NewTool("share_server",
				mcp.WithDescription("Exposes local MCP servers via secure tunnels with optional authentication. Enables remote teams to access your tools without VPN or port forwarding."),
				mcp.WithString("name",
					mcp.Description("Name of the server to share"),
					mcp.Required(),
				),
				mcp.WithString("port",
					mcp.Description("Port to run the shared server on"),
				),
				mcp.WithBoolean("no_auth",
					mcp.Description("Disable authentication for the shared server"),
				),
			),
			Handler: h.ShareServer,
		},
		{
			Tool: mcp.NewTool("stop_sharing_server",
				mcp.WithDescription("Revokes tunnel access and terminates shared server sessions. Immediately disconnects all remote clients."),
				mcp.WithString("name",
					mcp.Description("Name of the server to stop sharing"),
					mcp.Required(),
				),
			),
			Handler: h.StopSharingServer,
		},
		{
			Tool: mcp.NewTool("list_shared_servers",
				mcp.WithDescription("Shows all active server shares with tunnel URLs, authentication status, and connected clients."),
			),
			Handler: h.ListSharedServers,
		},
	}
}

// RegisterToolsWithMCPServer registers all Jarvis tools with an MCP server
func RegisterToolsWithMCPServer(s *server.MCPServer, h *Handler) {
	for _, def := range GetToolDefinitions(h) {
		s.AddTool(def.Tool, def.Handler)
	}
}

// CreateProductionHandler creates a Handler with real production dependencies
func CreateProductionHandler() *Handler {
	return NewHandler(
		&RealMcpmRunner{},
		&RealDockerRunner{},
		&RealGitRunner{},
		&RealFileSystem{},
	)
}

// RealDockerRunner implements DockerRunner using real Docker commands
type RealDockerRunner struct{}

func (r *RealDockerRunner) ComposeUp(ctx context.Context, services ...string) error {
	args := []string{"compose", "up", "-d"}
	args = append(args, services...)
	cmd := exec.CommandContext(ctx, "docker", args...)
	return cmd.Run()
}

func (r *RealDockerRunner) ComposeDown(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "compose", "down")
	return cmd.Run()
}

func (r *RealDockerRunner) ComposeRestart(ctx context.Context, services ...string) error {
	args := []string{"compose", "restart"}
	args = append(args, services...)
	cmd := exec.CommandContext(ctx, "docker", args...)
	return cmd.Run()
}

func (r *RealDockerRunner) ComposePs(ctx context.Context) ([]ContainerStatus, error) {
	cmd := exec.CommandContext(ctx, "docker", "compose", "ps", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON output (simplified - real impl would parse properly)
	var statuses []ContainerStatus
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Basic parsing - production would use proper JSON
		statuses = append(statuses, ContainerStatus{
			Name:    line,
			Running: strings.Contains(line, "running"),
		})
	}
	return statuses, nil
}

func (r *RealDockerRunner) ExecSupervisorctl(ctx context.Context, action, target string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "exec", "mcp-daemon", "supervisorctl", action, target)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RealGitRunner implements GitRunner using real git commands
type RealGitRunner struct{}

func (r *RealGitRunner) Status(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w", err)
	}
	return string(output), nil
}

func (r *RealGitRunner) Diff(ctx context.Context, staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--cached")
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}
	return string(output), nil
}

func (r *RealGitRunner) Init(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "init")
	return cmd.Run()
}
