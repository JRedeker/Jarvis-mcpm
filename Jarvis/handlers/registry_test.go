package handlers

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TDD: Write failing tests first for the Registry

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()

	// Register a handler
	reg.Register("check_status", func(h *Handler) ToolHandler {
		return h.CheckStatus
	})

	// Get the handler
	handler, exists := reg.Get("check_status")
	if !exists {
		t.Fatal("Expected check_status handler to exist")
	}
	if handler == nil {
		t.Fatal("Expected handler to not be nil")
	}
}

func TestRegistry_GetNonExistent(t *testing.T) {
	reg := NewRegistry()

	_, exists := reg.Get("nonexistent")
	if exists {
		t.Error("Expected nonexistent handler to not exist")
	}
}

func TestRegistry_ListAllTools(t *testing.T) {
	reg := NewRegistry()

	reg.Register("check_status", func(h *Handler) ToolHandler {
		return h.CheckStatus
	})
	reg.Register("list_servers", func(h *Handler) ToolHandler {
		return h.ListServers
	})

	tools := reg.List()
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	// Check both are present
	found := make(map[string]bool)
	for _, name := range tools {
		found[name] = true
	}
	if !found["check_status"] || !found["list_servers"] {
		t.Errorf("Expected check_status and list_servers, got %v", tools)
	}
}

func TestRegistry_ExecuteHandler(t *testing.T) {
	// Create mocks
	mcpm := NewMockMcpmRunner().
		WithResponse("doctor", "✅ All systems healthy")

	// Create handler with dependencies
	h := NewHandler(mcpm, nil, nil, nil)

	// Create registry and register
	reg := NewRegistry()
	reg.Register("check_status", func(h *Handler) ToolHandler {
		return h.CheckStatus
	})

	// Get and execute
	handlerFunc, exists := reg.Get("check_status")
	if !exists {
		t.Fatal("Handler not found")
	}

	// Create the actual handler function with our Handler instance
	result, err := handlerFunc(h)(context.Background(), mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestRegistry_RegisterAllCoreHandlers(t *testing.T) {
	reg := NewRegistry()
	RegisterAllHandlers(reg)

	// Verify all 24 core tools are registered
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
		"edit_server",
		"create_server",
		"usage_stats",
		"migrate_config",
		"restart_profiles",
		"suggest_profile",
		"fetch_diff_context",
		"analyze_project",
		"apply_devops_stack",
	}

	for _, name := range expectedTools {
		if _, exists := reg.Get(name); !exists {
			t.Errorf("Expected handler %s to be registered", name)
		}
	}
}

func TestRegistry_CreateWithDependencies(t *testing.T) {
	// Test that we can create a fully-wired registry
	mcpm := NewMockMcpmRunner()
	docker := NewMockDockerRunner()
	git := NewMockGitRunner()
	fs := NewMockFileSystem()

	reg := NewRegistry()
	RegisterAllHandlers(reg)

	// Create handler with all dependencies
	h := NewHandler(mcpm, docker, git, fs)

	// Verify we can get and call any handler
	for _, name := range reg.List() {
		handlerFactory, exists := reg.Get(name)
		if !exists {
			t.Errorf("Handler %s not found after registration", name)
			continue
		}

		// Should be able to create the handler function
		handlerFunc := handlerFactory(h)
		if handlerFunc == nil {
			t.Errorf("Handler factory for %s returned nil", name)
		}
	}
}
