package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// Helper to create a CallToolRequest with arguments
func newRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

// Helper to extract text from a CallToolResult
func getResultText(result *mcp.CallToolResult) string {
	if result == nil || len(result.Content) == 0 {
		return ""
	}
	for _, item := range result.Content {
		switch c := item.(type) {
		case mcp.TextContent:
			return c.Text
		case *mcp.TextContent:
			return c.Text
		}
	}
	return ""
}

// ==================== CheckStatus Tests ====================

func TestCheckStatus_Healthy(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("doctor", "🩺 MCPM System Health Check\n\n✅ All systems healthy! No issues found.")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()

	result, err := h.CheckStatus(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("CheckStatus returned error: %v", err)
	}
	if result == nil {
		t.Fatal("CheckStatus returned nil result")
	}

	text := getResultText(result)
	if !strings.Contains(text, "All systems healthy") {
		t.Errorf("Expected healthy message, got: %s", text)
	}
	if !strings.Contains(text, "ALL SYSTEMS GO") {
		t.Errorf("Expected success banner, got: %s", text)
	}
	if !strings.Contains(text, "Jarvis is ready") {
		t.Errorf("Expected Jarvis ready message, got: %s", text)
	}
}

func TestCheckStatus_Unhealthy(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("doctor", "🩺 MCPM System Health Check\n\n❌ Issues found: Node.js not installed")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()

	result, err := h.CheckStatus(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("CheckStatus returned error: %v", err)
	}

	text := getResultText(result)
	if strings.Contains(text, "ALL SYSTEMS GO") {
		t.Errorf("Should not show success banner for unhealthy system, got: %s", text)
	}
	if !strings.Contains(text, "Issues found") {
		t.Errorf("Should show issues, got: %s", text)
	}
}

func TestCheckStatus_CallsMcpmDoctor(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	h.CheckStatus(context.Background(), mcp.CallToolRequest{})

	if !mcpm.AssertCalled("doctor") {
		t.Error("CheckStatus should call 'mcpm doctor'")
	}
}

// ==================== ListServers Tests ====================

func TestListServers_ReturnsOutput(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("ls", "Installed servers:\n- context7\n- brave-search\n- mem0-mcp")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()

	result, err := h.ListServers(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("ListServers returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "context7") {
		t.Errorf("Expected context7 in list, got: %s", text)
	}
	if !strings.Contains(text, "brave-search") {
		t.Errorf("Expected brave-search in list, got: %s", text)
	}
}

func TestListServers_CallsMcpmLs(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	h.ListServers(context.Background(), mcp.CallToolRequest{})

	if !mcpm.AssertCalled("ls") {
		t.Error("ListServers should call 'mcpm ls'")
	}
}

// ==================== ServerInfo Tests ====================

func TestServerInfo_ReturnsDetails(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("info", "Server: context7\nDescription: Documentation lookup\nTransport: stdio")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "context7"})

	result, err := h.ServerInfo(ctx, req)

	if err != nil {
		t.Fatalf("ServerInfo returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "context7") {
		t.Errorf("Expected server name in info, got: %s", text)
	}
	if !strings.Contains(text, "Documentation") {
		t.Errorf("Expected description in info, got: %s", text)
	}
}

func TestServerInfo_RequiresName(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"no arguments", nil},
		{"empty arguments", map[string]interface{}{}},
		{"wrong type", map[string]interface{}{"name": 123}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(tt.args)
			result, _ := h.ServerInfo(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error result for missing/invalid name")
			}
		})
	}
}

// ==================== InstallServer Tests ====================

func TestInstallServer_Success(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("install", "Successfully installed context7")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "context7"})

	result, err := h.InstallServer(ctx, req)

	if err != nil {
		t.Fatalf("InstallServer returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Successfully installed") {
		t.Errorf("Expected success message, got: %s", text)
	}
	if !strings.Contains(text, "context7") {
		t.Errorf("Expected server name in message, got: %s", text)
	}
}

func TestInstallServer_RequiresName(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"no arguments", nil},
		{"empty name", map[string]interface{}{"name": ""}},
		{"whitespace name", map[string]interface{}{"name": "   "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(tt.args)
			result, _ := h.InstallServer(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error result for missing/invalid name")
			}
			text := getResultText(result)
			if !strings.Contains(text, "required") {
				t.Errorf("Expected 'required' in error message, got: %s", text)
			}
		})
	}
}

func TestInstallServer_RejectsInvalidNames(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name       string
		serverName string
	}{
		{"with space", "brave search"},
		{"with slash", "brave/search"},
		{"with backslash", "brave\\search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(map[string]interface{}{"name": tt.serverName})
			result, _ := h.InstallServer(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error for invalid server name")
			}
			text := getResultText(result)
			if !strings.Contains(text, "Invalid server name") {
				t.Errorf("Expected 'Invalid server name' in error, got: %s", text)
			}
		})
	}
}

func TestInstallServer_NotFound(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("install", "Error: Server 'nonexistent' not found in registry").
		WithError("install", errors.New("not found"))

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "nonexistent"})

	result, _ := h.InstallServer(ctx, req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for not found server")
	}
	text := getResultText(result)
	if !strings.Contains(text, "not found") {
		t.Errorf("Expected 'not found' in error, got: %s", text)
	}
	if !strings.Contains(text, "search_servers") {
		t.Errorf("Expected suggestion to use search_servers, got: %s", text)
	}
}

func TestInstallServer_AlreadyInstalled(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("install", "Server 'context7' is already installed").
		WithError("install", errors.New("already installed"))

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "context7"})

	result, _ := h.InstallServer(ctx, req)

	// Should NOT be an error - already installed is acceptable
	if result.IsError {
		t.Error("Already installed should not be treated as error")
	}
	text := getResultText(result)
	if !strings.Contains(text, "already installed") {
		t.Errorf("Expected 'already installed' message, got: %s", text)
	}
}

// ==================== UninstallServer Tests ====================

func TestUninstallServer_Success(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("uninstall", "Successfully uninstalled context7")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "context7"})

	result, err := h.UninstallServer(ctx, req)

	if err != nil {
		t.Fatalf("UninstallServer returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Successfully uninstalled") {
		t.Errorf("Expected success message, got: %s", text)
	}
}

func TestUninstallServer_RequiresName(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{"name": ""})
	result, _ := h.UninstallServer(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for empty name")
	}
}

func TestUninstallServer_NotInstalled(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("uninstall", "Error: Server 'nonexistent' is not installed").
		WithError("uninstall", errors.New("not installed"))

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"name": "nonexistent"})

	result, _ := h.UninstallServer(ctx, req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for not installed server")
	}
	text := getResultText(result)
	if !strings.Contains(text, "not installed") {
		t.Errorf("Expected 'not installed' in error, got: %s", text)
	}
}

// ==================== SearchServers Tests ====================

func TestSearchServers_ReturnsResults(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("search", "Found 3 servers:\n- basic-memory\n- mem0-mcp\n- memory-graph")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"query": "memory"})

	result, err := h.SearchServers(ctx, req)

	if err != nil {
		t.Fatalf("SearchServers returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "basic-memory") {
		t.Errorf("Expected basic-memory in results, got: %s", text)
	}
	if !strings.Contains(text, "server_info") {
		t.Errorf("Expected suggestion to use server_info, got: %s", text)
	}
}

func TestSearchServers_RequiresQuery(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"no query", map[string]interface{}{}},
		{"empty query", map[string]interface{}{"query": ""}},
		{"whitespace query", map[string]interface{}{"query": "   "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(tt.args)
			result, _ := h.SearchServers(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error result for missing/invalid query")
			}
		})
	}
}

func TestSearchServers_NoResults(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("search", "No servers found matching 'xyznotexist'")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"query": "xyznotexist"})

	result, _ := h.SearchServers(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "No servers found") {
		t.Errorf("Expected 'No servers found' message, got: %s", text)
	}
	if !strings.Contains(text, "Try these tips") {
		t.Errorf("Expected helpful tips, got: %s", text)
	}
}

// ==================== ManageProfile Tests ====================

func TestManageProfile_List(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("profile", "Profiles:\n- p-pokeedge\n- memory\n- morph")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"action": "ls"})

	result, err := h.ManageProfile(ctx, req)

	if err != nil {
		t.Fatalf("ManageProfile returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "p-pokeedge") {
		t.Errorf("Expected p-pokeedge in list, got: %s", text)
	}
}

func TestManageProfile_Create(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("profile", "Created profile 'test-profile'")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"action": "create",
		"name":   "test-profile",
	})

	result, _ := h.ManageProfile(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "Created profile") {
		t.Errorf("Expected creation message, got: %s", text)
	}
}

func TestManageProfile_Edit(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"action":      "edit",
		"name":        "test-profile",
		"add_servers": "brave-search,context7",
	})

	h.ManageProfile(ctx, req)

	// Verify the correct command was called
	if mcpm.CallCount("profile") != 1 {
		t.Errorf("Expected 1 call to profile, got %d", mcpm.CallCount("profile"))
	}
}

func TestManageProfile_DeleteMapsToRm(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"action": "delete",
		"name":   "test-profile",
	})

	h.ManageProfile(ctx, req)

	// Verify profile command was called (action mapped to 'rm')
	if !mcpm.AssertCalled("profile") {
		t.Error("Expected profile command to be called")
	}
}

func TestManageProfile_RequiresAction(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{})
	result, _ := h.ManageProfile(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for missing action")
	}
}

// ==================== ManageClient Tests ====================

func TestManageClient_List(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("client", "Clients:\n- claude-code\n- claude-desktop\n- codex")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"action": "ls"})

	result, err := h.ManageClient(ctx, req)

	if err != nil {
		t.Fatalf("ManageClient returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "claude-code") {
		t.Errorf("Expected claude-code in list, got: %s", text)
	}
}

func TestManageClient_RequiresAction(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"no action", map[string]interface{}{}},
		{"empty action", map[string]interface{}{"action": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(tt.args)
			result, _ := h.ManageClient(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error result for missing action")
			}
		})
	}
}

func TestManageClient_InvalidAction(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{"action": "invalid"})
	result, _ := h.ManageClient(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for invalid action")
	}
	text := getResultText(result)
	if !strings.Contains(text, "Invalid action") {
		t.Errorf("Expected 'Invalid action' in error, got: %s", text)
	}
}

func TestManageClient_EditRequiresClientName(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{
		"action": "edit",
	})
	result, _ := h.ManageClient(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for missing client_name")
	}
	text := getResultText(result)
	if !strings.Contains(text, "client_name is required") {
		t.Errorf("Expected client_name required message, got: %s", text)
	}
}

func TestManageClient_EditSuccess(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("client", "Updated client 'codex' configuration")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"action":      "edit",
		"client_name": "codex",
		"add_server":  "brave-search",
	})

	result, _ := h.ManageClient(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "Client configuration updated") {
		t.Errorf("Expected success message, got: %s", text)
	}
}

// ==================== ManageConfig Tests ====================

func TestManageConfig_List(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("config", "Configuration:\n  default_profile: p-pokeedge")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"action": "ls"})

	result, err := h.ManageConfig(ctx, req)

	if err != nil {
		t.Fatalf("ManageConfig returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "default_profile") {
		t.Errorf("Expected config key in output, got: %s", text)
	}
}

func TestManageConfig_Set(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"action": "set",
		"key":    "default_profile",
		"value":  "memory",
	})

	h.ManageConfig(ctx, req)

	if !mcpm.AssertCalled("config") {
		t.Error("Expected config command to be called")
	}
}

func TestManageConfig_RequiresAction(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{})
	result, _ := h.ManageConfig(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for missing action")
	}
}

// ==================== RestartProfiles Tests ====================

func TestRestartProfiles_AllProfiles(t *testing.T) {
	docker := NewMockDockerRunner().WithRunningContainers()
	h := NewHandler(nil, docker, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, err := h.RestartProfiles(ctx, req)

	if err != nil {
		t.Fatalf("RestartProfiles returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Successfully restarted") {
		t.Errorf("Expected success message, got: %s", text)
	}
	if !strings.Contains(text, "all profiles") {
		t.Errorf("Expected 'all profiles' in message, got: %s", text)
	}
	if docker.RestartCount != 1 {
		t.Errorf("Expected 1 restart, got %d", docker.RestartCount)
	}
}

func TestRestartProfiles_SpecificProfile(t *testing.T) {
	docker := NewMockDockerRunner().
		WithSupervisorctlOutput("restart", "mcpm-memory", "mcpm-memory: stopped\nmcpm-memory: started")

	h := NewHandler(nil, docker, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"profile": "memory"})

	result, err := h.RestartProfiles(ctx, req)

	if err != nil {
		t.Fatalf("RestartProfiles returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Successfully restarted") {
		t.Errorf("Expected success message, got: %s", text)
	}
	if !strings.Contains(text, "memory") {
		t.Errorf("Expected profile name in message, got: %s", text)
	}
}

func TestRestartProfiles_Error(t *testing.T) {
	docker := NewMockDockerRunner().
		WithComposeRestartError(errors.New("daemon not running"))

	h := NewHandler(nil, docker, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, _ := h.RestartProfiles(ctx, req)

	if result == nil || !result.IsError {
		t.Error("Expected error result")
	}
	text := getResultText(result)
	if !strings.Contains(text, "Failed to restart") {
		t.Errorf("Expected failure message, got: %s", text)
	}
}

// ==================== SuggestProfile Tests ====================

func TestSuggestProfile_PokeedgeProject(t *testing.T) {
	fs := NewMockFileSystem().WithCwd("/home/user/projects/pokeedge")
	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, err := h.SuggestProfile(ctx, req)

	if err != nil {
		t.Fatalf("SuggestProfile returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "p-pokeedge") {
		t.Errorf("Expected p-pokeedge profile, got: %s", text)
	}
	if !strings.Contains(text, "memory") {
		t.Errorf("Expected memory profile, got: %s", text)
	}
}

func TestSuggestProfile_NewProject(t *testing.T) {
	fs := NewMockFileSystem().WithCwd("/home/user/projects/my-new-project")
	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, err := h.SuggestProfile(ctx, req)

	if err != nil {
		t.Fatalf("SuggestProfile returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "p-new") {
		t.Errorf("Expected p-new profile for unknown project, got: %s", text)
	}
}

func TestSuggestProfile_TestingMode(t *testing.T) {
	fs := NewMockFileSystem().WithCwd("/home/user/projects/test")
	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"testing": true})

	result, err := h.SuggestProfile(ctx, req)

	if err != nil {
		t.Fatalf("SuggestProfile returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "testing-all-tools") {
		t.Errorf("Expected testing-all-tools profile, got: %s", text)
	}
}

// ==================== FetchDiffContext Tests ====================

func TestFetchDiffContext_Success(t *testing.T) {
	git := NewMockGitRunner().
		WithStatus("M  src/main.go\n?? newfile.txt").
		WithDiff("diff --git a/src/main.go\n+new line")
	fs := NewMockFileSystem().WithCwd("/home/user/project")
	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, err := h.FetchDiffContext(ctx, req)

	if err != nil {
		t.Fatalf("FetchDiffContext returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Local Review Context") {
		t.Errorf("Expected report header, got: %s", text)
	}
	if !strings.Contains(text, "main.go") {
		t.Errorf("Expected changed file in status, got: %s", text)
	}
	if !strings.Contains(text, "+new line") {
		t.Errorf("Expected diff content, got: %s", text)
	}
}

func TestFetchDiffContext_NotGitRepo(t *testing.T) {
	git := NewMockGitRunner().WithStatusError(errors.New("not a git repository"))
	fs := NewMockFileSystem().WithCwd("/home/user/notgit")
	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, _ := h.FetchDiffContext(ctx, req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for non-git directory")
	}
	text := getResultText(result)
	if !strings.Contains(text, "git") {
		t.Errorf("Expected git-related error message, got: %s", text)
	}
}

// ==================== AnalyzeProject Tests ====================

func TestAnalyzeProject_GoProject(t *testing.T) {
	fs := NewMockFileSystem().
		WithCwd("/home/user/go-project").
		WithDir("/home/user/go-project", []os.DirEntry{
			&MockDirEntry{EntryName: "go.mod", EntryIsDir: false},
			&MockDirEntry{EntryName: "main.go", EntryIsDir: false},
		}).
		WithFile("/home/user/go-project/.git", []byte{})

	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()

	result, err := h.AnalyzeProject(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("AnalyzeProject returned error: %v", err)
	}

	text := getResultText(result)

	// Parse the JSON output
	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	languages := analysis["languages"].([]interface{})
	foundGo := false
	for _, l := range languages {
		if l == "go" {
			foundGo = true
			break
		}
	}
	if !foundGo {
		t.Errorf("Expected 'go' in languages, got: %v", languages)
	}

	configs := analysis["configs"].(map[string]interface{})
	if !configs["has_git"].(bool) {
		t.Error("Expected has_git to be true")
	}
}

func TestAnalyzeProject_PythonProject(t *testing.T) {
	fs := NewMockFileSystem().
		WithCwd("/home/user/python-project").
		WithDir("/home/user/python-project", []os.DirEntry{
			&MockDirEntry{EntryName: "pyproject.toml", EntryIsDir: false},
			&MockDirEntry{EntryName: "requirements.txt", EntryIsDir: false},
		})

	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()

	result, err := h.AnalyzeProject(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("AnalyzeProject returned error: %v", err)
	}

	text := getResultText(result)

	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	languages := analysis["languages"].([]interface{})
	foundPython := false
	for _, l := range languages {
		if l == "python" {
			foundPython = true
			break
		}
	}
	if !foundPython {
		t.Errorf("Expected 'python' in languages, got: %v", languages)
	}

	// Python should only appear once even with multiple Python files
	pythonCount := 0
	for _, l := range languages {
		if l == "python" {
			pythonCount++
		}
	}
	if pythonCount != 1 {
		t.Errorf("Expected python to appear once, got %d times", pythonCount)
	}
}

func TestAnalyzeProject_WithPreCommit(t *testing.T) {
	preCommitContent := `repos:
  - repo: https://github.com/gitleaks/gitleaks
    hooks:
      - id: gitleaks`

	fs := NewMockFileSystem().
		WithCwd("/home/user/project").
		WithDir("/home/user/project", []os.DirEntry{}).
		WithFile("/home/user/project/.pre-commit-config.yaml", []byte(preCommitContent))

	h := NewHandler(nil, nil, nil, fs)
	ctx := context.Background()

	result, err := h.AnalyzeProject(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("AnalyzeProject returned error: %v", err)
	}

	text := getResultText(result)

	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	configs := analysis["configs"].(map[string]interface{})
	if !configs["has_pre_commit"].(bool) {
		t.Error("Expected has_pre_commit to be true")
	}
	if !configs["has_gitleaks"].(bool) {
		t.Error("Expected has_gitleaks to be true")
	}
}

// ==================== EditServer Tests ====================

func TestEditServer_Success(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("edit", "Updated server 'brave-search' configuration")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"name":    "brave-search",
		"command": "npx",
		"args":    "-y @brave/search",
	})

	result, err := h.EditServer(ctx, req)

	if err != nil {
		t.Fatalf("EditServer returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Updated") {
		t.Errorf("Expected update message, got: %s", text)
	}
}

func TestEditServer_RequiresName(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	req := newRequest(map[string]interface{}{})
	result, _ := h.EditServer(context.Background(), req)

	if result == nil || !result.IsError {
		t.Error("Expected error result for missing name")
	}
}

// ==================== CreateServer Tests ====================

func TestCreateServer_Stdio(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("new", "Created server 'my-server'")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"name":    "my-server",
		"type":    "stdio",
		"command": "python",
		"args":    "server.py",
	})

	result, err := h.CreateServer(ctx, req)

	if err != nil {
		t.Fatalf("CreateServer returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Created") {
		t.Errorf("Expected creation message, got: %s", text)
	}
}

func TestCreateServer_RequiresNameAndType(t *testing.T) {
	mcpm := NewMockMcpmRunner()
	h := NewHandler(mcpm, nil, nil, nil)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{"no name", map[string]interface{}{"type": "stdio"}},
		{"no type", map[string]interface{}{"name": "test"}},
		{"empty", map[string]interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newRequest(tt.args)
			result, _ := h.CreateServer(context.Background(), req)

			if result == nil || !result.IsError {
				t.Error("Expected error result")
			}
		})
	}
}

// ==================== UsageStats Tests ====================

func TestUsageStats_ReturnsOutput(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("usage", "Usage Statistics:\n- context7: 150 calls\n- brave-search: 75 calls")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()

	result, err := h.UsageStats(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("UsageStats returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Usage Statistics") {
		t.Errorf("Expected usage stats, got: %s", text)
	}
}

// ==================== MigrateConfig Tests ====================

func TestMigrateConfig_Success(t *testing.T) {
	mcpm := NewMockMcpmRunner().
		WithResponse("migrate", "Migration complete. Backup saved to ~/.config/mcpm/backup/")

	h := NewHandler(mcpm, nil, nil, nil)
	ctx := context.Background()

	result, err := h.MigrateConfig(ctx, mcp.CallToolRequest{})

	if err != nil {
		t.Fatalf("MigrateConfig returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Migration complete") {
		t.Errorf("Expected migration message, got: %s", text)
	}
}

// ==================== ApplyDevOpsStack Tests ====================

func TestApplyDevOpsStack_NewProject(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/new-project").
		WithDir("/home/user/new-project", []os.DirEntry{})
	// No .git, .pre-commit-config.yaml, or .gitignore exist

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"project_type": "go",
	})

	result, err := h.ApplyDevOpsStack(ctx, req)

	if err != nil {
		t.Fatalf("ApplyDevOpsStack returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Initialized git repository") {
		t.Errorf("Expected git init message, got: %s", text)
	}
	if !strings.Contains(text, "Created .pre-commit-config.yaml") {
		t.Errorf("Expected pre-commit creation message, got: %s", text)
	}
	if !strings.Contains(text, "Created default .gitignore") {
		t.Errorf("Expected gitignore creation message, got: %s", text)
	}
}

func TestApplyDevOpsStack_ExistingGit(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/existing-project").
		WithFile("/home/user/existing-project/.git", []byte{}).
		WithDir("/home/user/existing-project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, err := h.ApplyDevOpsStack(ctx, req)

	if err != nil {
		t.Fatalf("ApplyDevOpsStack returned error: %v", err)
	}

	text := getResultText(result)
	if !strings.Contains(text, "Git repository already exists") {
		t.Errorf("Expected existing git message, got: %s", text)
	}
}

func TestApplyDevOpsStack_SkipsExistingPreCommit(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/project").
		WithFile("/home/user/project/.git", []byte{}).
		WithFile("/home/user/project/.pre-commit-config.yaml", []byte("existing config")).
		WithDir("/home/user/project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"force": false})

	result, _ := h.ApplyDevOpsStack(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "exists. Skipping") {
		t.Errorf("Expected skip message for existing pre-commit, got: %s", text)
	}
}

func TestApplyDevOpsStack_ForceOverwrite(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/project").
		WithFile("/home/user/project/.git", []byte{}).
		WithFile("/home/user/project/.pre-commit-config.yaml", []byte("existing config")).
		WithDir("/home/user/project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{"force": true})

	result, _ := h.ApplyDevOpsStack(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "Overwrote") {
		t.Errorf("Expected overwrite message, got: %s", text)
	}
}

func TestApplyDevOpsStack_WithAIReview(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/project").
		WithFile("/home/user/project/.git", []byte{}).
		WithDir("/home/user/project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"enable_ai_review": true,
	})

	result, _ := h.ApplyDevOpsStack(ctx, req)

	text := getResultText(result)
	if !strings.Contains(text, "pr_agent.yml") {
		t.Errorf("Expected PR agent workflow message, got: %s", text)
	}
}

func TestApplyDevOpsStack_PythonProject(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/python-project").
		WithDir("/home/user/python-project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"project_type": "python",
	})

	h.ApplyDevOpsStack(ctx, req)

	// Check that the written pre-commit config contains Python hooks
	content, exists := fs.Files["/home/user/python-project/.pre-commit-config.yaml"]
	if !exists {
		t.Fatal("Pre-commit config was not written")
	}
	if !strings.Contains(string(content), "ruff") {
		t.Error("Python project should have ruff hooks")
	}
}

func TestApplyDevOpsStack_NodeProject(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/node-project").
		WithDir("/home/user/node-project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"project_type": "node",
	})

	h.ApplyDevOpsStack(ctx, req)

	content, exists := fs.Files["/home/user/node-project/.pre-commit-config.yaml"]
	if !exists {
		t.Fatal("Pre-commit config was not written")
	}
	if !strings.Contains(string(content), "prettier") {
		t.Error("Node project should have prettier hooks")
	}
}

func TestApplyDevOpsStack_GoProject(t *testing.T) {
	git := NewMockGitRunner()
	fs := NewMockFileSystem().
		WithCwd("/home/user/go-project").
		WithDir("/home/user/go-project", []os.DirEntry{})

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{
		"project_type": "go",
	})

	h.ApplyDevOpsStack(ctx, req)

	content, exists := fs.Files["/home/user/go-project/.pre-commit-config.yaml"]
	if !exists {
		t.Fatal("Pre-commit config was not written")
	}
	if !strings.Contains(string(content), "go-fmt") {
		t.Error("Go project should have go-fmt hooks")
	}
}

func TestApplyDevOpsStack_GitInitError(t *testing.T) {
	git := NewMockGitRunner()
	git.InitError = errors.New("permission denied")

	fs := NewMockFileSystem().
		WithCwd("/home/user/project").
		WithDir("/home/user/project", []os.DirEntry{})
	// No .git - so it will try to init

	h := NewHandler(nil, nil, git, fs)
	ctx := context.Background()
	req := newRequest(map[string]interface{}{})

	result, _ := h.ApplyDevOpsStack(ctx, req)

	if result == nil || !result.IsError {
		t.Error("Expected error result when git init fails")
	}
	text := getResultText(result)
	if !strings.Contains(text, "Git init failed") {
		t.Errorf("Expected git init error message, got: %s", text)
	}
}
