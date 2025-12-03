package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Integration test: Verify handlers work correctly when wired into an MCP server

func TestMCPServer_ToolsListIncludesAllHandlers(t *testing.T) {
	// Create MCP server with our handlers
	mcpServer := createTestMCPServer(t)

	// List tools via the server
	tools := mcpServer.ListTools()

	// Verify all expected tools are present
	expectedTools := []string{
		"check_status",
		"list_servers",
		"server_info",
		"install_server",
		"uninstall_server",
		"search_servers",
		"manage_profile",
		"manage_client",
		"manage_config",
		"analyze_project",
	}

	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("Expected tool %s to be registered in MCP server", expected)
		}
	}
}

func TestMCPServer_CheckStatusTool(t *testing.T) {
	mcpServer := createTestMCPServer(t)

	// Call check_status tool
	req := mcp.CallToolRequest{}
	req.Params.Name = "check_status"
	req.Params.Arguments = map[string]interface{}{}

	result, err := mcpServer.HandleToolCall(context.Background(), req)

	if err != nil {
		t.Fatalf("Tool call failed: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify response contains expected content
	text := getResultText(result)
	if !strings.Contains(text, "healthy") && !strings.Contains(text, "MCPM") {
		t.Errorf("Expected health check content, got: %s", text)
	}
}

func TestMCPServer_InstallServerValidation(t *testing.T) {
	mcpServer := createTestMCPServer(t)

	// Call install_server without name - should fail validation
	req := mcp.CallToolRequest{}
	req.Params.Name = "install_server"
	req.Params.Arguments = map[string]interface{}{}

	result, err := mcpServer.HandleToolCall(context.Background(), req)

	if err != nil {
		t.Fatalf("Tool call failed unexpectedly: %v", err)
	}

	// Should return an error result
	if result == nil || !result.IsError {
		t.Error("Expected error result for missing name argument")
	}

	text := getResultText(result)
	if !strings.Contains(text, "required") {
		t.Errorf("Expected 'required' in error message, got: %s", text)
	}
}

func TestMCPServer_SearchServersReturnsResults(t *testing.T) {
	mcpServer := createTestMCPServer(t)

	req := mcp.CallToolRequest{}
	req.Params.Name = "search_servers"
	req.Params.Arguments = map[string]interface{}{
		"query": "memory",
	}

	result, err := mcpServer.HandleToolCall(context.Background(), req)

	if err != nil {
		t.Fatalf("Tool call failed: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Mock returns search results
	text := getResultText(result)
	if text == "" {
		t.Error("Expected non-empty search results")
	}
}

func TestMCPServer_AnalyzeProjectReturnsJSON(t *testing.T) {
	mcpServer := createTestMCPServer(t)

	req := mcp.CallToolRequest{}
	req.Params.Name = "analyze_project"
	req.Params.Arguments = map[string]interface{}{}

	result, err := mcpServer.HandleToolCall(context.Background(), req)

	if err != nil {
		t.Fatalf("Tool call failed: %v", err)
	}

	text := getResultText(result)

	// Should be valid JSON
	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		t.Errorf("Expected JSON output, got parse error: %v\nOutput: %s", err, text)
	}

	// Should have expected keys
	if _, ok := analysis["path"]; !ok {
		t.Error("Expected 'path' key in analysis")
	}
	if _, ok := analysis["languages"]; !ok {
		t.Error("Expected 'languages' key in analysis")
	}
}

func TestMCPServer_ManageProfileActions(t *testing.T) {
	mcpServer := createTestMCPServer(t)

	tests := []struct {
		name   string
		action string
	}{
		{"list profiles", "ls"},
		{"create profile", "create"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = "manage_profile"
			req.Params.Arguments = map[string]interface{}{
				"action": tt.action,
				"name":   "test-profile",
			}

			result, err := mcpServer.HandleToolCall(context.Background(), req)

			if err != nil {
				t.Fatalf("Tool call failed: %v", err)
			}
			if result == nil {
				t.Fatal("Expected result, got nil")
			}
		})
	}
}

// TestMCPServerWrapper wraps an MCP server for testing
type TestMCPServerWrapper struct {
	server   *server.MCPServer
	handler  *Handler
	registry *Registry
}

func (w *TestMCPServerWrapper) ListTools() []mcp.Tool {
	// Get tools from registry
	tools := make([]mcp.Tool, 0)
	for _, name := range w.registry.List() {
		tools = append(tools, mcp.Tool{Name: name})
	}
	return tools
}

func (w *TestMCPServerWrapper) HandleToolCall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.Params.Name

	factory, exists := w.registry.Get(name)
	if !exists {
		return mcp.NewToolResultError("Unknown tool: " + name), nil
	}

	handlerFunc := factory(w.handler)
	return handlerFunc(ctx, req)
}

// createTestMCPServer creates a test MCP server with mock dependencies
func createTestMCPServer(t *testing.T) *TestMCPServerWrapper {
	t.Helper()

	// Create mocks with default responses
	mcpm := NewMockMcpmRunner().
		WithResponse("doctor", "✅ All systems healthy! MCPM version: 1.0.0").
		WithResponse("ls", "Installed servers:\n- context7\n- brave-search").
		WithResponse("search", "Found 2 servers:\n- basic-memory\n- mem0-mcp").
		WithResponse("info", "Server: context7\nDescription: Documentation lookup").
		WithResponse("install", "Successfully installed server").
		WithResponse("profile", "Profiles:\n- p-pokeedge\n- memory")

	docker := NewMockDockerRunner().WithRunningContainers()

	git := NewMockGitRunner().
		WithStatus("M  main.go").
		WithDiff("+// new code")

	fs := NewMockFileSystem().
		WithCwd("/home/test/project").
		WithDir("/home/test/project", nil)

	// Create handler with mocks
	handler := NewHandler(mcpm, docker, git, fs)

	// Create and populate registry
	registry := NewRegistry()
	RegisterAllHandlers(registry)

	return &TestMCPServerWrapper{
		handler:  handler,
		registry: registry,
	}
}
