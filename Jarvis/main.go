package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	sharedServers      = make(map[string]*exec.Cmd)
	sharedServersMutex sync.Mutex
)

func printBanner() {
	banner := `
     ██╗ █████╗ ██████╗ ██╗   ██╗██╗███████╗
     ██║██╔══██╗██╔══██╗██║   ██║██║██╔════╝
     ██║███████║██████╔╝██║   ██║██║███████╗
██   ██║██╔══██║██╔══██╗╚██╗ ██╔╝██║╚════██║
╚█████╔╝██║  ██║██║  ██║ ╚████╔╝ ██║███████║
 ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝
`
	// Print to Stderr to avoid interfering with MCP stdio protocol
	fmt.Fprintln(os.Stderr, "\033[36m"+banner+"\033[0m")
	fmt.Fprintln(os.Stderr, "\033[1;32m>> JARVIS MCP Gateway v1.0.0 initialized <<\033[0m")
}

func main() {
	printBanner()

	// Create a new MCP server
	s := server.NewMCPServer(
		"jarvis",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Tool: bootstrap_system
	s.AddTool(mcp.NewTool("bootstrap_system",
		mcp.WithDescription("Initialize the MCP environment (install dependencies, start infrastructure)"),
	), handleBootstrapSystem)

	// Tool: list_servers
	s.AddTool(mcp.NewTool("list_servers",
		mcp.WithDescription("List all installed MCP servers managed by MCPM"),
	), handleListServers)

	// Tool: install_server
	s.AddTool(mcp.NewTool("install_server",
		mcp.WithDescription("Install a new MCP server using MCPM"),
		mcp.WithString("name",
			mcp.Description("Name of the server to install"),
			mcp.Required(),
		),
	), handleInstallServer)

	// Tool: server_info
	s.AddTool(mcp.NewTool("server_info",
		mcp.WithDescription("Get detailed information about a specific MCP server"),
		mcp.WithString("name",
			mcp.Description("Name of the server"),
			mcp.Required(),
		),
	), handleServerInfo)

	// Tool: check_status
	s.AddTool(mcp.NewTool("check_status",
		mcp.WithDescription("Check the health and status of the MCPM system"),
	), handleCheckStatus)

	// Tool: search_servers
	s.AddTool(mcp.NewTool("search_servers",
		mcp.WithDescription("Search available MCP servers"),
		mcp.WithString("query",
			mcp.Description("Search query"),
			mcp.Required(),
		),
	), handleSearchServers)

	// Tool: uninstall_server
	s.AddTool(mcp.NewTool("uninstall_server",
		mcp.WithDescription("Remove an installed MCP server"),
		mcp.WithString("name",
			mcp.Description("Name of the server to uninstall"),
			mcp.Required(),
		),
	), handleUninstallServer)

	// Tool: edit_server
	s.AddTool(mcp.NewTool("edit_server",
		mcp.WithDescription("Edit a server configuration"),
		mcp.WithString("name",
			mcp.Description("Name of the server to edit"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("New command (for stdio servers)"),
		),
		mcp.WithString("args",
			mcp.Description("New arguments (space-separated)"),
		),
		mcp.WithString("env",
			mcp.Description("New environment variables (KEY=value, அண)"),
		),
		mcp.WithString("url",
			mcp.Description("New URL (for remote servers)"),
		),
		mcp.WithString("headers",
			mcp.Description("New headers (KEY=value, அண)"),
		),
	), handleEditServer)

	// Tool: create_server
	s.AddTool(mcp.NewTool("create_server",
		mcp.WithDescription("Create a new server configuration"),
		mcp.WithString("name",
			mcp.Description("Name of the new server"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description("Server type (stdio|remote)"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute (required for stdio)"),
		),
		mcp.WithString("args",
			mcp.Description("Command arguments"),
		),
		mcp.WithString("env",
			mcp.Description("Environment variables (KEY=value, அண) (stdio only)"),
		),
		mcp.WithString("url",
			mcp.Description("Server URL (required for remote)"),
		),
		mcp.WithString("headers",
			mcp.Description("HTTP headers (KEY=value, அண) (remote only)"),
		),
	), handleCreateServer)

	// Tool: usage_stats
	s.AddTool(mcp.NewTool("usage_stats",
		mcp.WithDescription("Display comprehensive analytics and usage data"),
	), handleUsageStats)

	// Tool: manage_client
	s.AddTool(mcp.NewTool("manage_client",
		mcp.WithDescription("Manage MCP client configurations"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|edit|import)"),
			mcp.Required(),
		),
		mcp.WithString("client_name",
			mcp.Description("Client name (required for edit/import)"),
		),
		mcp.WithString("add_server",
			mcp.Description("Servers to add (comma-separated, edit only)"),
		),
		mcp.WithString("remove_server",
			mcp.Description("Servers to remove (comma-separated, edit only)"),
		),
		mcp.WithString("add_profile",
			mcp.Description("Profiles to add (comma-separated, edit only)"),
		),
		mcp.WithString("remove_profile",
			mcp.Description("Profiles to remove (comma-separated, edit only)"),
		),
	), handleManageClient)

	// Tool: manage_profile
	s.AddTool(mcp.NewTool("manage_profile",
		mcp.WithDescription("Manage MCPM profiles"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|create|edit|delete)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Profile name (required for create/edit/delete)"),
		),
		mcp.WithString("new_name",
			mcp.Description("New profile name (only for edit action to rename)"),
		),
		mcp.WithString("add_servers",
			mcp.Description("Comma-separated list of servers to add (only for edit action)"),
		),
		mcp.WithString("remove_servers",
			mcp.Description("Comma-separated list of servers to remove (only for edit action)"),
		),
	), handleManageProfile)

	// Tool: manage_config
	s.AddTool(mcp.NewTool("manage_config",
		mcp.WithDescription("Manage MCPM configuration"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|set|unset)"),
			mcp.Required(),
		),
		mcp.WithString("key",
			mcp.Description("Config key (required for set/unset)"),
		),
		mcp.WithString("value",
			mcp.Description("Config value (required for set)"),
		),
	), handleManageConfig)

	// Tool: migrate_config
	s.AddTool(mcp.NewTool("migrate_config",
		mcp.WithDescription("Migrate v1 configuration to v2"),
	), handleMigrateConfig)

	// Tool: share_server
	s.AddTool(mcp.NewTool("share_server",
		mcp.WithDescription("Share a local server via a secure tunnel"),
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
	), handleShareServer)

	// Tool: stop_sharing_server
	s.AddTool(mcp.NewTool("stop_sharing_server",
		mcp.WithDescription("Stop sharing a server"),
		mcp.WithString("name",
			mcp.Description("Name of the server to stop sharing"),
			mcp.Required(),
		),
	), handleStopSharingServer)

	// Tool: list_shared_servers
	s.AddTool(mcp.NewTool("list_shared_servers",
		mcp.WithDescription("List currently shared servers"),
	), handleListSharedServers)

	// Start the server using Stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Helper function to run mcpm commands
func runMcpmCommand(args ...string) (string, error) {
	// mcpm is now available in PATH
	cmd := exec.Command("mcpm", args...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("🚫 command failed: %s, output: %s", err, string(output))
	}
	return string(output), nil
}
