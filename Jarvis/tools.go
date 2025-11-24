package main

import (
	"bufio"
	"context"
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

	// 3. Start Infrastructure
	cmdCompose := exec.Command("docker-compose", "up", "-d")
	cmdCompose.Dir = rootDir
	if out, err := cmdCompose.CombinedOutput(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run docker-compose up in %s: %v\nOutput: %s", rootDir, err, string(out))), nil
	}

	return mcp.NewToolResultText("System bootstrapped successfully! MCPM installed and Infrastructure started."), nil
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

	if action == "edit" || action == "import" {
		clientName, ok := args["client_name"].(string)
		if !ok || clientName == "" {
			return mcp.NewToolResultError("client_name argument is required for this action"), nil
		}
		cmdArgs = append(cmdArgs, clientName)
	}

	if action == "edit" {
		cmdArgs = append(cmdArgs, buildManageClientArgs(args)...)
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
