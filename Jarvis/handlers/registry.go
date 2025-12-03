package handlers

import (
	"context"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

// ToolHandler is the function signature for MCP tool handlers
type ToolHandler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

// HandlerFactory creates a ToolHandler from a Handler instance
// This allows dependency injection - the Handler holds the dependencies,
// and the factory extracts the specific method to call
type HandlerFactory func(h *Handler) ToolHandler

// Registry manages all available tool handlers
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFactory
}

// NewRegistry creates a new empty handler registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]HandlerFactory),
	}
}

// Register adds a handler factory to the registry
func (r *Registry) Register(name string, factory HandlerFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = factory
}

// Get retrieves a handler factory by name
func (r *Registry) Get(name string) (HandlerFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, exists := r.handlers[name]
	return factory, exists
}

// List returns all registered handler names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	return names
}

// RegisterAllHandlers registers all core Jarvis tool handlers
func RegisterAllHandlers(reg *Registry) {
	// System Management
	reg.Register("check_status", func(h *Handler) ToolHandler {
		return h.CheckStatus
	})

	// Server Management
	reg.Register("list_servers", func(h *Handler) ToolHandler {
		return h.ListServers
	})
	reg.Register("server_info", func(h *Handler) ToolHandler {
		return h.ServerInfo
	})
	reg.Register("install_server", func(h *Handler) ToolHandler {
		return h.InstallServer
	})
	reg.Register("uninstall_server", func(h *Handler) ToolHandler {
		return h.UninstallServer
	})
	reg.Register("search_servers", func(h *Handler) ToolHandler {
		return h.SearchServers
	})
	reg.Register("edit_server", func(h *Handler) ToolHandler {
		return h.EditServer
	})
	reg.Register("create_server", func(h *Handler) ToolHandler {
		return h.CreateServer
	})
	reg.Register("usage_stats", func(h *Handler) ToolHandler {
		return h.UsageStats
	})

	// Profile Management
	reg.Register("manage_profile", func(h *Handler) ToolHandler {
		return h.ManageProfile
	})
	reg.Register("suggest_profile", func(h *Handler) ToolHandler {
		return h.SuggestProfile
	})
	reg.Register("restart_profiles", func(h *Handler) ToolHandler {
		return h.RestartProfiles
	})

	// Client Management
	reg.Register("manage_client", func(h *Handler) ToolHandler {
		return h.ManageClient
	})

	// Configuration
	reg.Register("manage_config", func(h *Handler) ToolHandler {
		return h.ManageConfig
	})
	reg.Register("migrate_config", func(h *Handler) ToolHandler {
		return h.MigrateConfig
	})

	// Project Analysis
	reg.Register("analyze_project", func(h *Handler) ToolHandler {
		return h.AnalyzeProject
	})
	reg.Register("fetch_diff_context", func(h *Handler) ToolHandler {
		return h.FetchDiffContext
	})
	reg.Register("apply_devops_stack", func(h *Handler) ToolHandler {
		return h.ApplyDevOpsStack
	})
}
