package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"jarvis",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Tool: list_servers
	s.AddTool(mcp.NewTool("list_servers",
		mcp.WithDescription("List all installed MCP servers managed by MCPM"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := runMcpmCommand("list")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list servers: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Tool: install_server
	s.AddTool(mcp.NewTool("install_server",
		mcp.WithDescription("Install a new MCP server using MCPM"),
		mcp.WithString("name",
			mcp.Description("Name of the server to install"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	})

	// Tool: server_info
	s.AddTool(mcp.NewTool("server_info",
		mcp.WithDescription("Get detailed information about a specific MCP server"),
		mcp.WithString("name",
			mcp.Description("Name of the server"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	})

	// Tool: check_status
	s.AddTool(mcp.NewTool("check_status",
		mcp.WithDescription("Check the health and status of the MCPM system"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := runMcpmCommand("doctor")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to check status: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Tool: search_servers
	s.AddTool(mcp.NewTool("search_servers",
		mcp.WithDescription("Search available MCP servers"),
		mcp.WithString("query",
			mcp.Description("Search query"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	})

	// Tool: uninstall_server
	s.AddTool(mcp.NewTool("uninstall_server",
		mcp.WithDescription("Remove an installed MCP server"),
		mcp.WithString("name",
			mcp.Description("Name of the server to uninstall"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	})

	// Tool: edit_server
	s.AddTool(mcp.NewTool("edit_server",
		mcp.WithDescription("Edit a server configuration"),
		mcp.WithString("name",
			mcp.Description("Name of the server to edit"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Note: 'edit' is interactive by default. We might need to handle this carefully or assume non-interactive mode via env vars if supported.
		// For now, we'll just run it, but it might fail if it prompts for input.
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResultError("name argument is required"), nil
		}

		output, err := runMcpmCommand("edit", name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to edit server %s: %v", name, err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

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
		mcp.WithString("url",
			mcp.Description("Server URL (required for remote)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if url, ok := args["url"].(string); ok && url != "" {
			cmdArgs = append(cmdArgs, "--url", url)
		}

		output, err := runMcpmCommand(cmdArgs...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create server %s: %v", name, err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Tool: usage_stats
	s.AddTool(mcp.NewTool("usage_stats",
		mcp.WithDescription("Display comprehensive analytics and usage data"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := runMcpmCommand("usage")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get usage stats: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Tool: manage_client
	s.AddTool(mcp.NewTool("manage_client",
		mcp.WithDescription("Manage MCP client configurations"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|edit|import)"), // Simplified for now
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}
		action, ok := args["action"].(string)
		if !ok {
			return mcp.NewToolResultError("action argument is required"), nil
		}

		// Note: 'client' command might have subcommands. Adjusting based on typical usage.
		// Assuming 'mcpm client <action>' works or similar.
		// If 'mcpm client' is interactive, this might fail.
		// Let's assume basic 'mcpm client' lists clients if no args, but here we take an action.
		// Actually, 'mcpm client' usually takes subcommands. Let's try passing the action.
		output, err := runMcpmCommand("client", action)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to manage client: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Tool: manage_profile
	s.AddTool(mcp.NewTool("manage_profile",
		mcp.WithDescription("Manage MCPM profiles"),
		mcp.WithString("action",
			mcp.Description("Action to perform (ls|create|switch|delete)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Profile name (required for create/switch/delete)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}
		action, ok := args["action"].(string)
		if !ok {
			return mcp.NewToolResultError("action argument is required"), nil
		}

		cmdArgs := []string{"profile", action}
		if name, ok := args["name"].(string); ok && name != "" {
			cmdArgs = append(cmdArgs, name)
		}

		output, err := runMcpmCommand(cmdArgs...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to manage profile: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	})

	// Tool: migrate_config
	s.AddTool(mcp.NewTool("migrate_config",
		mcp.WithDescription("Migrate v1 configuration to v2"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := runMcpmCommand("migrate")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to migrate config: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	})

	// Start the server using Stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Helper function to run mcpm commands
func runMcpmCommand(args ...string) (string, error) {
	// We assume 'mcpm' is available in the PATH or mounted into the container
	cmd := exec.Command("mcpm", args...)
	// Ensure non-interactive mode where possible
	cmd.Env = append(cmd.Env, "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, output: %s", err, string(output))
	}
	return string(output), nil
}
