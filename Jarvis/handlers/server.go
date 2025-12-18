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
// CONSOLIDATED: 24 tools -> 8 tools for ~73% context token reduction
func GetToolDefinitions(h *Handler) []ToolDefinition {
	return []ToolDefinition{
		// 1. jarvis_check_status - System health (no action param, single purpose)
		{
			Tool: mcp.NewTool("jarvis_check_status",
				mcp.WithDescription("System health check for MCPM, Docker, and services."),
			),
			Handler: h.CheckStatus,
		},

		// 2. jarvis_server - Server management (8 actions)
		{
			Tool: mcp.NewTool("jarvis_server",
				mcp.WithDescription("Manage MCP servers: list, info, install, uninstall, search, edit, create, usage."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("list", "info", "install", "uninstall", "search", "edit", "create", "usage"),
				),
				mcp.WithString("name",
					mcp.Description("Server name (for info/install/uninstall/edit/create)"),
				),
				mcp.WithString("query",
					mcp.Description("Search query (for search action)"),
				),
				mcp.WithString("type",
					mcp.Description("Transport: 'stdio' or 'streamable-http' (for create)"),
				),
				mcp.WithString("command",
					mcp.Description("Command to run (for edit/create)"),
				),
				mcp.WithString("args",
					mcp.Description("Arguments, space-separated (for edit/create)"),
				),
				mcp.WithString("env",
					mcp.Description("Env vars KEY=val, comma-separated (for edit/create)"),
				),
				mcp.WithString("url",
					mcp.Description("URL for HTTP servers (for edit/create)"),
				),
				mcp.WithString("headers",
					mcp.Description("HTTP headers KEY=val, comma-separated (for edit/create)"),
				),
			),
			Handler: h.Server,
		},

		// 3. jarvis_profile - Profile management (6 actions)
		{
			Tool: mcp.NewTool("jarvis_profile",
				mcp.WithDescription("Manage MCPM profiles: list, create, edit, delete, suggest, restart."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("list", "create", "edit", "delete", "suggest", "restart"),
				),
				mcp.WithString("name",
					mcp.Description("Profile name"),
				),
				mcp.WithString("new_name",
					mcp.Description("New name when renaming (for edit)"),
				),
				mcp.WithString("add_servers",
					mcp.Description("Servers to add, comma-separated (for edit)"),
				),
				mcp.WithString("remove_servers",
					mcp.Description("Servers to remove, comma-separated (for edit)"),
				),
				mcp.WithString("profile",
					mcp.Description("Specific profile to restart (for restart)"),
				),
				mcp.WithBoolean("testing",
					mcp.Description("Include testing profile (for suggest)"),
				),
			),
			Handler: h.Profile,
		},

		// 4. jarvis_client - Client management (4 actions)
		{
			Tool: mcp.NewTool("jarvis_client",
				mcp.WithDescription("Configure MCP clients: list, edit, import, config."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("list", "edit", "import", "config"),
				),
				mcp.WithString("client_name",
					mcp.Description("Client: opencode, claude-code, claude-desktop"),
				),
				mcp.WithString("add_server",
					mcp.Description("Server to add (for edit)"),
				),
				mcp.WithString("remove_server",
					mcp.Description("Server to remove (for edit)"),
				),
				mcp.WithString("add_profile",
					mcp.Description("Profile to add (for edit)"),
				),
				mcp.WithString("remove_profile",
					mcp.Description("Profile to remove (for edit)"),
				),
				mcp.WithString("config_path",
					mcp.Description("Config file path (for config)"),
				),
			),
			Handler: h.Client,
		},

		// 5. jarvis_config - Configuration management (4 actions)
		{
			Tool: mcp.NewTool("jarvis_config",
				mcp.WithDescription("Manage MCPM config: get, set, list, migrate."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("get", "set", "list", "migrate"),
				),
				mcp.WithString("key",
					mcp.Description("Config key (for get/set)"),
				),
				mcp.WithString("value",
					mcp.Description("Value to set (for set)"),
				),
			),
			Handler: h.Config,
		},

		// 6. jarvis_project - Project analysis (3 actions)
		{
			Tool: mcp.NewTool("jarvis_project",
				mcp.WithDescription("Project tools: analyze, diff, devops."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("analyze", "diff", "devops"),
				),
				mcp.WithBoolean("staged",
					mcp.Description("Show only staged changes (for diff)"),
				),
				mcp.WithString("project_type",
					mcp.Description("Override: python, go, node, typescript (for devops)"),
				),
				mcp.WithBoolean("force",
					mcp.Description("Overwrite existing configs (for devops)"),
				),
				mcp.WithBoolean("enable_ai_review",
					mcp.Description("Add PR Agent workflow (for devops)"),
				),
			),
			Handler: h.Project,
		},

		// 7. jarvis_system - System operations (3 actions)
		{
			Tool: mcp.NewTool("jarvis_system",
				mcp.WithDescription("System ops: bootstrap, restart, restart_infra."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("bootstrap", "restart", "restart_infra"),
				),
			),
			Handler: h.System,
		},

		// 8. jarvis_share - Server sharing (3 actions)
		{
			Tool: mcp.NewTool("jarvis_share",
				mcp.WithDescription("Share MCP servers: start, stop, list."),
				mcp.WithString("action",
					mcp.Description("Operation to perform"),
					mcp.Required(),
					mcp.Enum("start", "stop", "list"),
				),
				mcp.WithString("name",
					mcp.Description("Server name (for start/stop)"),
				),
				mcp.WithString("port",
					mcp.Description("Port for shared server (for start)"),
				),
				mcp.WithBoolean("no_auth",
					mcp.Description("Disable authentication (for start)"),
				),
			),
			Handler: h.Share,
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
