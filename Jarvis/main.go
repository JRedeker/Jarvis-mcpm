package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	sharedServers      = make(map[string]*exec.Cmd)
	sharedServersMutex sync.Mutex
	logFile            *os.File
)

func setupLogging() {
	// Determine project root (assuming Jarvis runs from Jarvis/ or project root)
	// We'll try to find the 'logs' directory in the parent or current directory
	logDir := "logs"
	if _, err := os.Stat("../logs"); err == nil {
		logDir = "../logs"
	} else {
		os.MkdirAll("logs", 0755)
	}

	logPath := filepath.Join(logDir, "jarvis.log")
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stderr if file creation fails
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return
	}

	// Create a MultiWriter to write to both the log file and stderr (so IDEs still see errors)
	// However, for pure logging, we might want just the file to avoid polluting the protocol stream
	// if logging libraries print to stdout/stderr by default.
	// mcp-go/server.WithLogging() uses stderr by default.

	// We will set the global logger to write to the file
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf(">> Jarvis Logging Initialized <<")
}

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
	setupLogging()
	printBanner()

	// Create a new MCP server
	s := server.NewMCPServer(
		"jarvis",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(), // This logs MCP protocol messages to stderr
	)

	// Tool: bootstrap_system
	s.AddTool(mcp.NewTool("bootstrap_system",
		mcp.WithDescription("Initialize the MCP environment (install dependencies, start infrastructure)"),
	), handleBootstrapSystem)

	// Tool: restart_service
	s.AddTool(mcp.NewTool("restart_service",
		mcp.WithDescription("Gracefully restarts the Jarvis MCP server"),
	), handleRestartService)

	// Tool: restart_infrastructure
	s.AddTool(mcp.NewTool("restart_infrastructure",
		mcp.WithDescription("Restarts the underlying MCP infrastructure (Docker containers)"),
	), handleRestartInfrastructure)

	// Tool: suggest_profile
	s.AddTool(mcp.NewTool("suggest_profile",
		mcp.WithDescription("Suggests the appropriate MCPM profile stack based on context, client, and mode"),
		mcp.WithBoolean("testing",
			mcp.Description("Whether testing mode is active"),
		),
		mcp.WithString("client_name",
			mcp.Description("The name of the client requesting the profile (e.g., 'gemini', 'codex')"),
		),
	), handleSuggestProfile)

	// Tool: fetch_diff_context
	s.AddTool(mcp.NewTool("fetch_diff_context",
		mcp.WithDescription("Retrieves git diff and status for local AI code review"),
		mcp.WithBoolean("staged",
			mcp.Description("If true, only show staged changes. If false, show all changes."),
		),
	), handleFetchDiffContext)

	// Tool: scaffold_project
	s.AddTool(mcp.NewTool("scaffold_project",
		mcp.WithDescription("Scaffolds a new project with standard dev tooling (git, pre-commit, AI review)"),
		mcp.WithString("project_type",
			mcp.Description("Type of project (python, go, node, general)"),
			mcp.Required(),
		),
		mcp.WithBoolean("enable_ai_review",
			mcp.Description("Setup GitHub Actions for AI PR review"),
			mcp.DefaultString("true"),
		),
	), handleScaffoldProject)

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
			mcp.Description("New command (for stdio servers)")),
		mcp.WithString("args",
			mcp.Description("New arguments (space-separated)")),
		mcp.WithString("env",
			mcp.Description("New environment variables (KEY=value, அண)")),
		mcp.WithString("url",
			mcp.Description("New URL (for remote servers)")),
		mcp.WithString("headers",
			mcp.Description("New headers (KEY=value, அண)")),
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
			mcp.Description("Command to execute (required for stdio)")),
		mcp.WithString("args",
			mcp.Description("Command arguments")),
		mcp.WithString("env",
			mcp.Description("Environment variables (KEY=value, அண) (stdio only)")),
		mcp.WithString("url",
			mcp.Description("Server URL (required for remote)")),
		mcp.WithString("headers",
			mcp.Description("HTTP headers (KEY=value, அண) (remote only)")),
	), handleCreateServer)

	// Tool: usage_stats
	s.AddTool(mcp.NewTool("usage_stats",
		mcp.WithDescription("Display comprehensive analytics and usage data"),
	), handleUsageStats)

	// Tool: manage_client
	s.AddTool(mcp.NewTool("manage_client",
		mcp.WithDescription("Manage MCP client configurations"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|edit|import|config)"),
			mcp.Required(),
		),
		mcp.WithString("client_name",
			mcp.Description("Client name (required for edit/import/config)")),
		mcp.WithString("config_path",
			mcp.Description("Custom configuration path (for config action)")),
		mcp.WithString("add_server",
			mcp.Description("Servers to add (comma-separated, edit only)")),
		mcp.WithString("remove_server",
			mcp.Description("Servers to remove (comma-separated, edit only)")),
		mcp.WithString("add_profile",
			mcp.Description("Profiles to add (comma-separated, edit only)")),
		mcp.WithString("remove_profile",
			mcp.Description("Profiles to remove (comma-separated, edit only)")),
	), handleManageClient)

	// Tool: manage_profile
	s.AddTool(mcp.NewTool("manage_profile",
		mcp.WithDescription("Manage MCPM profiles"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|create|edit|delete)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Profile name (required for create/edit/delete)")),
		mcp.WithString("new_name",
			mcp.Description("New profile name (only for edit action to rename)")),
		mcp.WithString("add_servers",
			mcp.Description("Comma-separated list of servers to add (only for edit action)")),
		mcp.WithString("remove_servers",
			mcp.Description("Comma-separated list of servers to remove (only for edit action)")),
	), handleManageProfile)

	// Tool: manage_config
	s.AddTool(mcp.NewTool("manage_config",
		mcp.WithDescription("Manage MCPM configuration"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|set|unset)"),
			mcp.Required(),
		),
		mcp.WithString("key",
			mcp.Description("Config key (required for set/unset)")),
		mcp.WithString("value",
			mcp.Description("Config value (required for set)")),
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
			mcp.Description("Port to run the shared server on")),
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
	log.Printf("Executing MCPM command: %v", args)
	// mcpm is now available in PATH
	cmd := exec.Command("mcpm", args...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %v. Output: %s", err, string(output))
		return "", fmt.Errorf("🚫 command failed: %s, output: %s", err, string(output))
	}
	log.Printf("Command success. Output length: %d", len(output))
	return string(output), nil
}
