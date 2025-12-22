package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock API Server for testing HTTPMcpmRunner
func setupMockAPIServer(t *testing.T, responses map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find matching response
		key := r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			key += "?" + r.URL.RawQuery
		}

		resp, ok := responses[key]
		if !ok {
			// Try without query string
			resp, ok = responses[r.Method+" "+r.URL.Path]
		}
		if !ok {
			// Default 404
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"data":    nil,
				"error": map[string]interface{}{
					"code":    "NOT_FOUND",
					"message": "Endpoint not found",
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestHTTPMcpmRunner_Health(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/health": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"status":    "healthy",
				"timestamp": "2025-01-01T00:00:00Z",
				"version":   "1.0.0",
				"checks": map[string]interface{}{
					"node": map[string]interface{}{
						"status":  "ok",
						"version": "v20.0.0",
					},
					"registry": map[string]interface{}{
						"status": "ok",
						"path":   "/config/technologies.toml",
					},
					"docker": map[string]interface{}{
						"status": "detected",
					},
				},
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("doctor")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "MCPM System Status") {
		t.Errorf("Expected status header, got: %s", output)
	}
	if !strings.Contains(output, "Node.js") {
		t.Errorf("Expected Node.js info, got: %s", output)
	}
	if !strings.Contains(output, "All systems healthy") {
		t.Errorf("Expected healthy message, got: %s", output)
	}
}

func TestHTTPMcpmRunner_ListServers(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/servers": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"count": 2,
				"servers": []map[string]interface{}{
					{
						"name":        "context7",
						"group":       "documentation",
						"description": "Documentation lookup",
						"installed":   true,
					},
					{
						"name":        "brave-search",
						"group":       "search",
						"description": "Web search",
						"installed":   false,
					},
				},
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("ls")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Installed MCP Servers") {
		t.Errorf("Expected servers header, got: %s", output)
	}
	if !strings.Contains(output, "context7") {
		t.Errorf("Expected context7 in output, got: %s", output)
	}
}

func TestHTTPMcpmRunner_ServerInfo(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/servers/context7": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"name":        "context7",
				"description": "Documentation lookup",
				"type":        "stdio",
				"source":      "documentation",
				"installed":   true,
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("info", "context7")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "context7") {
		t.Errorf("Expected server name in output, got: %s", output)
	}
	if !strings.Contains(output, "Documentation lookup") {
		t.Errorf("Expected description in output, got: %s", output)
	}
}

func TestHTTPMcpmRunner_InstallServer(t *testing.T) {
	responses := map[string]interface{}{
		"POST /api/v1/servers/context7/install": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"name":    "context7",
				"method":  "npm",
				"message": "Successfully installed context7",
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("install", "context7")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Successfully installed") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestHTTPMcpmRunner_UninstallServer(t *testing.T) {
	responses := map[string]interface{}{
		"DELETE /api/v1/servers/context7": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"name":    "context7",
				"message": "Successfully uninstalled context7",
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("uninstall", "context7")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Successfully uninstalled") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestHTTPMcpmRunner_SearchServers(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/search": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"query": "memory",
				"count": 2,
				"results": []map[string]interface{}{
					{"name": "basic-memory", "group": "memory", "description": "Memory storage"},
					{"name": "mem0-mcp", "group": "memory", "description": "Mem0 integration"},
				},
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("search", "memory")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Search results") {
		t.Errorf("Expected search header, got: %s", output)
	}
	if !strings.Contains(output, "basic-memory") {
		t.Errorf("Expected basic-memory in results, got: %s", output)
	}
}

func TestHTTPMcpmRunner_ListProfiles(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/profiles": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"count": 2,
				"profiles": []map[string]interface{}{
					{"name": "toolbox", "servers": []string{"context7", "brave-search"}},
					{"name": "memory", "servers": []string{"basic-memory"}},
				},
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("profile", "ls")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Profiles") {
		t.Errorf("Expected profiles header, got: %s", output)
	}
	if !strings.Contains(output, "toolbox") {
		t.Errorf("Expected toolbox in output, got: %s", output)
	}
}

func TestHTTPMcpmRunner_ListClients(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/clients": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"count": 2,
				"clients": []map[string]interface{}{
					{"name": "opencode", "displayName": "OpenCode", "detected": true, "configPath": "/home/user/.config/opencode/opencode.json"},
					{"name": "claude-code", "displayName": "Claude Code CLI", "detected": false, "configPath": ""},
				},
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("client", "ls")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "MCP Clients") {
		t.Errorf("Expected clients header, got: %s", output)
	}
	if !strings.Contains(output, "OpenCode") {
		t.Errorf("Expected OpenCode in output, got: %s", output)
	}
}

func TestHTTPMcpmRunner_Usage(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/usage": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"servers": map[string]interface{}{
					"total":      10,
					"installed":  3,
					"byCategory": map[string]int{"documentation": 2, "search": 3},
				},
				"profiles": map[string]interface{}{
					"total": 2,
					"list":  []string{"toolbox", "memory"},
				},
				"clients": map[string]interface{}{
					"configured": 1,
					"list":       []string{"opencode"},
				},
				"configDir": "/home/user/.config/mcpm",
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("usage")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Output is pretty-printed JSON
	if !strings.Contains(output, "servers") {
		t.Errorf("Expected servers in output, got: %s", output)
	}
}

func TestHTTPMcpmRunner_Migrate(t *testing.T) {
	responses := map[string]interface{}{
		"POST /api/v1/migrate": map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"migrations": []string{"Created config directory", "Created servers.json"},
				"configDir":  "/home/user/.config/mcpm",
				"message":    "Migration complete",
			},
			"error": nil,
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	output, err := runner.Run("migrate")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(output, "Migration Results") {
		t.Errorf("Expected migration header, got: %s", output)
	}
	if !strings.Contains(output, "Created config directory") {
		t.Errorf("Expected migration details, got: %s", output)
	}
}

func TestHTTPMcpmRunner_NoCommand(t *testing.T) {
	runner := NewHTTPMcpmRunner("http://localhost:6275")
	_, err := runner.Run()

	if err == nil {
		t.Error("Expected error for no command")
	}
	if !strings.Contains(err.Error(), "no command") {
		t.Errorf("Expected 'no command' error, got: %v", err)
	}
}

func TestHTTPMcpmRunner_UnknownCommand(t *testing.T) {
	runner := NewHTTPMcpmRunner("http://localhost:6275")
	_, err := runner.Run("unknown-command-xyz")

	if err == nil {
		t.Error("Expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Expected 'unknown command' error, got: %v", err)
	}
}

func TestHTTPMcpmRunner_MissingRequiredArgs(t *testing.T) {
	runner := NewHTTPMcpmRunner("http://localhost:6275")

	tests := []struct {
		name string
		args []string
	}{
		{"info without name", []string{"info"}},
		{"install without name", []string{"install"}},
		{"uninstall without name", []string{"uninstall"}},
		{"search without query", []string{"search"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runner.Run(tt.args...)
			if err == nil {
				t.Errorf("Expected error for %s", tt.name)
			}
		})
	}
}

func TestHTTPMcpmRunner_APIError(t *testing.T) {
	responses := map[string]interface{}{
		"GET /api/v1/servers/nonexistent": map[string]interface{}{
			"success": false,
			"data":    nil,
			"error": map[string]interface{}{
				"code":    "SERVER_NOT_FOUND",
				"message": "Server 'nonexistent' not found",
			},
		},
	}

	server := setupMockAPIServer(t, responses)
	defer server.Close()

	runner := NewHTTPMcpmRunner(server.URL)
	_, err := runner.Run("info", "nonexistent")

	if err == nil {
		t.Error("Expected error for not found server")
	}
	if !strings.Contains(err.Error(), "SERVER_NOT_FOUND") {
		t.Errorf("Expected SERVER_NOT_FOUND error code, got: %v", err)
	}
}

func TestHTTPMcpmRunner_DefaultURL(t *testing.T) {
	runner := NewHTTPMcpmRunner("")

	if runner.BaseURL != "http://localhost:6275" {
		t.Errorf("Expected default URL, got: %s", runner.BaseURL)
	}
}

func TestNewMcpmRunner_DefaultsToHTTPWithFallback(t *testing.T) {
	// This test checks that NewMcpmRunner returns a runner
	// (it will fall back to CLI since no API server is running)
	runner := NewMcpmRunner()

	if runner == nil {
		t.Error("Expected a runner to be returned")
	}

	// Should be RealMcpmRunner since API is not available
	if _, ok := runner.(*RealMcpmRunner); !ok {
		// Could also be HTTPMcpmRunner if API happens to be running
		if _, ok := runner.(*HTTPMcpmRunner); !ok {
			t.Error("Expected RealMcpmRunner or HTTPMcpmRunner")
		}
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ok", "OK"},
		{"healthy", "OK"},
		{"detected", "OK"},
		{"missing", "Not Found"},
		{"not_found", "Not Found"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatStatus(tt.input)
			if result != tt.expected {
				t.Errorf("formatStatus(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
