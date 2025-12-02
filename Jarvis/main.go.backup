package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		mcp.WithDescription("Complete system initialization: installs MCPM, sets up default servers (context7, brave-search, github), and starts Docker infrastructure (PostgreSQL, Qdrant). One command to get fully operational."),
	), handleBootstrapSystem)

	// Tool: restart_service
	s.AddTool(mcp.NewTool("restart_service",
		mcp.WithDescription("Gracefully restarts Jarvis to apply configuration changes or resolve stuck states. Automatically saves state and reconnects active sessions. Use after editing server configs or when tools become unresponsive."),
	), handleRestartService)

	// Tool: restart_infrastructure
	s.AddTool(mcp.NewTool("restart_infrastructure",
		mcp.WithDescription("Safely reboots Docker infrastructure (PostgreSQL, Qdrant) with health checks and automatic reconnection. Resolves database connection issues, clears stale locks, and ensures all services are healthy. Zero data loss."),
	), handleRestartInfrastructure)

	// Tool: suggest_profile
	s.AddTool(mcp.NewTool("suggest_profile",
		mcp.WithDescription("Intelligently determines optimal MCPM profile stack by analyzing working directory, client type, and mode. Returns recommended profiles with explanations for why each is needed. Prevents profile conflicts and missing dependencies."),
		mcp.WithBoolean("testing",
			mcp.Description("Whether testing mode is active"),
		),
		mcp.WithString("client_name",
			mcp.Description("The name of the client requesting the profile (e.g., 'gemini', 'codex')"),
		),
	), handleSuggestProfile)

	// Tool: fetch_diff_context
	s.AddTool(mcp.NewTool("fetch_diff_context",
		mcp.WithDescription("Retrieves git status and diff for AI-powered code review before commits. Includes both staged and unstaged changes, file status (modified/added/deleted), and conflict markers. Essential for self-review workflows."),
		mcp.WithBoolean("staged",
			mcp.Description("If true, only show staged changes. If false, show all changes."),
		),
	), handleFetchDiffContext)

	// Tool: apply_devops_stack
	s.AddTool(mcp.NewTool("apply_devops_stack",
		mcp.WithDescription("Scaffold production-ready DevOps: git initialization, pre-commit hooks, language-specific linting (Python/Go/Node), secret detection, and optional AI-powered PR reviews via GitHub Actions."),
		mcp.WithString("project_type",
			mcp.Description("Type of project (python, go, node, general). Auto-detected if omitted."),
		),
		mcp.WithBoolean("enable_ai_review",
			mcp.Description("Setup GitHub Actions for AI PR review (recommended)"),
			mcp.DefaultString("true"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing configuration files (use with caution)"),
			mcp.DefaultString("false"),
		),
	), handleApplyDevOpsStack)

	// Tool: analyze_project
	s.AddTool(mcp.NewTool("analyze_project",
		mcp.WithDescription("Intelligent project analysis: detects languages (Python/Go/Node/etc.), identifies existing config files (pre-commit, linters, CI/CD), and returns structured JSON. Use before apply_devops_stack."),
	), handleAnalyzeProject)

	// Tool: list_servers
	s.AddTool(mcp.NewTool("list_servers",
		mcp.WithDescription("Displays all installed MCP servers with installation methods, status, and profile memberships. Shows both global servers and profile-specific servers. Useful for auditing your MCP stack and identifying unused servers."),
	), handleListServers)

	// Tool: install_server
	s.AddTool(mcp.NewTool("install_server",
		mcp.WithDescription("Install MCP servers with automatic dependency resolution, validation, and clean error messages. Handles Docker, npm, and pip installations seamlessly."),
		mcp.WithString("name",
			mcp.Description("Name of the server to install"),
			mcp.Required(),
		),
	), handleInstallServer)

	// Tool: server_info
	s.AddTool(mcp.NewTool("server_info",
		mcp.WithDescription("Detailed server documentation including description, installation methods, environment variables, usage examples, and links. Use this before installing to understand what you're getting."),
		mcp.WithString("name",
			mcp.Description("Name of the server"),
			mcp.Required(),
		),
	), handleServerInfo)

	// Tool: check_status
	s.AddTool(mcp.NewTool("check_status",
		mcp.WithDescription("Comprehensive system diagnostics: MCPM version, Python/Node environments, config validation, client status, and profile health checks. Your go-to tool for troubleshooting."),
	), handleCheckStatus)

	// Tool: search_servers
	s.AddTool(mcp.NewTool("search_servers",
		mcp.WithDescription("Find MCP servers across the registry with fuzzy matching. Returns rich metadata including categories, tags, examples, and installation methods. Perfect for discovering new capabilities."),
		mcp.WithString("query",
			mcp.Description("Search query (supports fuzzy matching)"),
			mcp.Required(),
		),
	), handleSearchServers)

	// Tool: uninstall_server
	s.AddTool(mcp.NewTool("uninstall_server",
		mcp.WithDescription("Cleanly removes MCP servers and updates all affected profiles and client configurations. Prevents broken references and orphaned dependencies. Shows impact analysis before removal."),
		mcp.WithString("name",
			mcp.Description("Name of the server to uninstall"),
			mcp.Required(),
		),
	), handleUninstallServer)

	// Tool: edit_server
	s.AddTool(mcp.NewTool("edit_server",
		mcp.WithDescription("Modify server configurations with validation and automatic client config updates. Change commands, arguments, environment variables, or remote URLs. Useful for updating API keys or fixing connection issues without reinstalling."),
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
		mcp.WithDescription("Register custom MCP servers not in the official registry. Supports both stdio (local) and remote (SSE) transports. Validates configuration before saving. Perfect for internal tools or development servers."),
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
		mcp.WithDescription("Analyzes MCP usage patterns with tool call frequency, server popularity, error rates, and performance metrics. Identifies underutilized servers and optimization opportunities. Helps justify infrastructure costs."),
	), handleUsageStats)

	// Tool: manage_client
	s.AddTool(mcp.NewTool("manage_client",
		mcp.WithDescription("Configure AI clients (Claude Code, Claude Desktop, Codex, Gemini, etc.) with servers and profiles. Lists installed clients, adds/removes configurations, and persists paths automatically."),
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
		mcp.WithDescription("Create, edit, and manage MCPM profiles (collections of servers). Supports listing all profiles, adding/removing servers, renaming, and deletion. Essential for organizing your MCP stack."),
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
		mcp.WithDescription("View and modify MCPM global settings including default profiles, installation preferences, and behavior flags. Changes persist across sessions. Use 'ls' to see all current settings before making changes."),
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
		mcp.WithDescription("Upgrades MCPM v1 configurations to v2 format with automatic backup and validation. Preserves all server definitions, profiles, and client configurations. Run this once when upgrading from legacy MCPM installations."),
	), handleMigrateConfig)

	// Tool: share_server
	s.AddTool(mcp.NewTool("share_server",
		mcp.WithDescription("Exposes local MCP servers via secure tunnels with optional authentication. Enables remote teams to access your tools without VPN or port forwarding. Auto-generates shareable URLs with configurable access controls."),
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
		mcp.WithDescription("Revokes tunnel access and terminates shared server sessions. Immediately disconnects all remote clients. Changes are logged for security auditing."),
		mcp.WithString("name",
			mcp.Description("Name of the server to stop sharing"),
			mcp.Required(),
		),
	), handleStopSharingServer)

	// Tool: list_shared_servers
	s.AddTool(mcp.NewTool("list_shared_servers",
		mcp.WithDescription("Shows all active server shares with tunnel URLs, authentication status, connected clients, and uptime. Useful for monitoring remote access and identifying security risks."),
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
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true", "NO_COLOR=true")

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Strip common noise from MCPM output
	outputStr = stripMcpmNoise(outputStr)

	if err != nil {
		log.Printf("Command failed: %v. Output: %s", err, outputStr)
		return fmt.Sprintf("Error: %v\n\n%s", err, outputStr), fmt.Errorf("command failed: %v", err)
	}
	log.Printf("Command success. Output length: %d", len(output))

	return strings.TrimSpace(outputStr), nil
}

// stripMcpmNoise removes common warnings and noise from MCPM output
func stripMcpmNoise(output string) string {
	lines := strings.Split(output, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		// Skip warning lines
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
