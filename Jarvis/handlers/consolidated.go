// Package handlers provides consolidated MCP tool handlers
// This file implements the action-based tool consolidation pattern
// to reduce context token usage from ~2,750 tokens to ~750 tokens
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Server handles the consolidated jarvis_server tool
// Actions: list, info, install, uninstall, search, edit, create, usage
func (h *Handler) Server(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: list|info|install|uninstall|search|edit|create|usage"), nil
	}

	switch action {
	case "list":
		return h.ListServers(ctx, req)
	case "info":
		return h.ServerInfo(ctx, req)
	case "install":
		return h.InstallServer(ctx, req)
	case "uninstall":
		return h.UninstallServer(ctx, req)
	case "search":
		return h.SearchServers(ctx, req)
	case "edit":
		return h.EditServer(ctx, req)
	case "create":
		return h.CreateServer(ctx, req)
	case "usage":
		return h.UsageStats(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: list|info|install|uninstall|search|edit|create|usage", action)), nil
	}
}

// Profile handles the consolidated jarvis_profile tool
// Actions: list, create, edit, delete, suggest, restart
func (h *Handler) Profile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: list|create|edit|delete|suggest|restart"), nil
	}

	switch action {
	case "list":
		// Map to manage_profile with ls action
		args["action"] = "ls"
		return h.ManageProfile(ctx, req)
	case "create":
		return h.ManageProfile(ctx, req)
	case "edit":
		return h.ManageProfile(ctx, req)
	case "delete":
		return h.ManageProfile(ctx, req)
	case "suggest":
		return h.SuggestProfile(ctx, req)
	case "restart":
		return h.RestartProfiles(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: list|create|edit|delete|suggest|restart", action)), nil
	}
}

// Client handles the consolidated jarvis_client tool
// Actions: list, edit, import, config
func (h *Handler) Client(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: list|edit|import|config"), nil
	}

	// Map action names for ManageClient
	switch action {
	case "list":
		args["action"] = "ls"
	case "edit", "import", "config":
		// Keep as-is
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: list|edit|import|config", action)), nil
	}

	return h.ManageClient(ctx, req)
}

// Config handles the consolidated jarvis_config tool
// Actions: get, set, list, migrate
func (h *Handler) Config(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: get|set|list|migrate"), nil
	}

	switch action {
	case "get", "set":
		return h.ManageConfig(ctx, req)
	case "list":
		args["action"] = "ls"
		return h.ManageConfig(ctx, req)
	case "migrate":
		return h.MigrateConfig(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: get|set|list|migrate", action)), nil
	}
}

// Project handles the consolidated jarvis_project tool
// Actions: analyze, diff, devops
func (h *Handler) Project(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: analyze|diff|devops"), nil
	}

	switch action {
	case "analyze":
		return h.AnalyzeProject(ctx, req)
	case "diff":
		return h.FetchDiffContext(ctx, req)
	case "devops":
		return h.ApplyDevOpsStack(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: analyze|diff|devops", action)), nil
	}
}

// System handles the consolidated jarvis_system tool
// Actions: bootstrap, restart, restart_infra
func (h *Handler) System(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: bootstrap|restart|restart_infra"), nil
	}

	switch action {
	case "bootstrap":
		return h.BootstrapSystem(ctx, req)
	case "restart":
		return h.RestartService(ctx, req)
	case "restart_infra":
		return h.RestartInfrastructure(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: bootstrap|restart|restart_infra", action)), nil
	}
}

// Share handles the consolidated jarvis_share tool
// Actions: start, stop, list
func (h *Handler) Share(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: start|stop|list"), nil
	}

	switch action {
	case "start":
		return h.ShareServer(ctx, req)
	case "stop":
		return h.StopSharingServer(ctx, req)
	case "list":
		return h.ListSharedServers(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: start|stop|list", action)), nil
	}
}

// Diagnose handles the jarvis_diagnose tool for debugging MCP profile issues
// Actions: profile_health, test_endpoint, logs, full, config_sync
func (h *Handler) Diagnose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	action, ok := args["action"].(string)
	if !ok || action == "" {
		return mcp.NewToolResultError("action is required. Valid: profile_health|test_endpoint|logs|full|config_sync"), nil
	}

	switch action {
	case "profile_health":
		return h.DiagnoseProfileHealth(ctx, req)
	case "test_endpoint":
		return h.DiagnoseTestEndpoint(ctx, req)
	case "logs":
		return h.DiagnoseLogs(ctx, req)
	case "full":
		return h.DiagnoseFull(ctx, req)
	case "config_sync":
		return h.DiagnoseConfigSync(ctx, req)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action '%s'. Valid: profile_health|test_endpoint|logs|full|config_sync", action)), nil
	}
}
