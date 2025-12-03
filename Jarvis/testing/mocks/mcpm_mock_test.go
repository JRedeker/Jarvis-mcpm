package mocks

import (
	"context"
	"testing"
)

func TestMockMcpmClient_Doctor(t *testing.T) {
	ctx := context.Background()

	t.Run("returns healthy response by default", func(t *testing.T) {
		mock := NewMockMcpmClient()

		result, err := mock.Doctor(ctx)

		if err != nil {
			t.Errorf("Doctor() returned error: %v", err)
		}
		if result == nil {
			t.Fatal("Doctor() returned nil result")
		}
		if !result.AllHealthy {
			t.Error("Doctor() default response should be healthy")
		}
		if !mock.AssertCalled("Doctor") {
			t.Error("Doctor() was not recorded in call history")
		}
	})

	t.Run("returns configured healthy response", func(t *testing.T) {
		mock := NewMockMcpmClient().WithHealthyDoctor()

		result, err := mock.Doctor(ctx)

		if err != nil {
			t.Errorf("Doctor() returned error: %v", err)
		}
		if result.MCPMVersion != "2.9.0" {
			t.Errorf("Doctor() MCPMVersion = %q, want %q", result.MCPMVersion, "2.9.0")
		}
	})

	t.Run("returns configured unhealthy response", func(t *testing.T) {
		mock := NewMockMcpmClient().WithUnhealthyDoctor("Node.js not installed")

		result, err := mock.Doctor(ctx)

		if err != nil {
			t.Errorf("Doctor() returned error: %v", err)
		}
		if result.AllHealthy {
			t.Error("Doctor() should return unhealthy")
		}
		if len(result.Issues) != 1 {
			t.Errorf("Doctor() Issues = %v, want 1 issue", result.Issues)
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		mock := NewMockMcpmClient()
		mock.DoctorError = context.DeadlineExceeded

		_, err := mock.Doctor(ctx)

		if err != context.DeadlineExceeded {
			t.Errorf("Doctor() error = %v, want DeadlineExceeded", err)
		}
	})
}

func TestMockMcpmClient_Install(t *testing.T) {
	ctx := context.Background()

	t.Run("returns success by default", func(t *testing.T) {
		mock := NewMockMcpmClient()

		result, err := mock.Install(ctx, "test-server")

		if err != nil {
			t.Errorf("Install() returned error: %v", err)
		}
		if !result.Success {
			t.Error("Install() should succeed by default")
		}
		if !mock.AssertCalled("Install", "test-server") {
			t.Error("Install() was not recorded correctly")
		}
	})

	t.Run("returns configured already installed", func(t *testing.T) {
		mock := NewMockMcpmClient().WithInstallResult("context7", &InstallResult{
			Success:          true,
			AlreadyInstalled: true,
			Message:          "context7 is already installed",
		})

		result, err := mock.Install(ctx, "context7")

		if err != nil {
			t.Errorf("Install() returned error: %v", err)
		}
		if !result.AlreadyInstalled {
			t.Error("Install() should return already installed")
		}
	})
}

func TestMockMcpmClient_List(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list by default", func(t *testing.T) {
		mock := NewMockMcpmClient()

		result, err := mock.List(ctx)

		if err != nil {
			t.Errorf("List() returned error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("List() = %v, want empty", result)
		}
	})

	t.Run("returns configured servers", func(t *testing.T) {
		mock := NewMockMcpmClient().WithServers(
			ServerInfo{Name: "context7", Installed: true},
			ServerInfo{Name: "brave-search", Installed: true},
		)

		result, err := mock.List(ctx)

		if err != nil {
			t.Errorf("List() returned error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("List() = %d servers, want 2", len(result))
		}
	})
}

func TestMockMcpmClient_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty for unknown query", func(t *testing.T) {
		mock := NewMockMcpmClient()

		result, err := mock.Search(ctx, "unknown")

		if err != nil {
			t.Errorf("Search() returned error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Search() = %v, want empty", result)
		}
	})

	t.Run("returns configured results", func(t *testing.T) {
		mock := NewMockMcpmClient().WithSearchResults("memory", []ServerInfo{
			{Name: "basic-memory", Description: "Basic memory storage"},
			{Name: "mem0-mcp", Description: "Mem0 memory server"},
		})

		result, err := mock.Search(ctx, "memory")

		if err != nil {
			t.Errorf("Search() returned error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Search() = %d results, want 2", len(result))
		}
	})
}

func TestMockMcpmClient_Info(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error for unknown server", func(t *testing.T) {
		mock := NewMockMcpmClient()

		_, err := mock.Info(ctx, "unknown")

		if err == nil {
			t.Error("Info() should return error for unknown server")
		}
	})

	t.Run("returns configured info", func(t *testing.T) {
		mock := NewMockMcpmClient().WithServerInfo("context7", &ServerInfo{
			Name:        "context7",
			Description: "Documentation lookup server",
			Installed:   true,
		})

		result, err := mock.Info(ctx, "context7")

		if err != nil {
			t.Errorf("Info() returned error: %v", err)
		}
		if result.Name != "context7" {
			t.Errorf("Info() Name = %q, want %q", result.Name, "context7")
		}
	})
}

func TestMockMcpmClient_CallTracking(t *testing.T) {
	ctx := context.Background()
	mock := NewMockMcpmClient()

	// Make some calls
	mock.Doctor(ctx)
	mock.Install(ctx, "server1")
	mock.Install(ctx, "server2")
	mock.List(ctx)

	t.Run("AssertCalled works", func(t *testing.T) {
		if !mock.AssertCalled("Doctor") {
			t.Error("Doctor was called but AssertCalled returned false")
		}
		if !mock.AssertCalled("Install", "server1") {
			t.Error("Install(server1) was called but AssertCalled returned false")
		}
		if mock.AssertCalled("Install", "server3") {
			t.Error("Install(server3) was not called but AssertCalled returned true")
		}
	})

	t.Run("AssertNotCalled works", func(t *testing.T) {
		if mock.AssertNotCalled("Doctor") {
			t.Error("Doctor was called but AssertNotCalled returned true")
		}
		if !mock.AssertNotCalled("Uninstall") {
			t.Error("Uninstall was not called but AssertNotCalled returned false")
		}
	})

	t.Run("CallCount works", func(t *testing.T) {
		if mock.CallCount("Install") != 2 {
			t.Errorf("CallCount(Install) = %d, want 2", mock.CallCount("Install"))
		}
		if mock.CallCount("Doctor") != 1 {
			t.Errorf("CallCount(Doctor) = %d, want 1", mock.CallCount("Doctor"))
		}
		if mock.CallCount("Uninstall") != 0 {
			t.Errorf("CallCount(Uninstall) = %d, want 0", mock.CallCount("Uninstall"))
		}
	})

	t.Run("ResetCalls clears history", func(t *testing.T) {
		mock.ResetCalls()

		if mock.CallCount("Doctor") != 0 {
			t.Error("ResetCalls did not clear call history")
		}
	})
}

func TestMockMcpmClient_ProfileOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("ProfileList returns configured profiles", func(t *testing.T) {
		mock := NewMockMcpmClient().WithProfiles(
			ProfileInfo{Name: "p-pokeedge", Servers: []string{"context7", "brave-search"}},
			ProfileInfo{Name: "memory", Servers: []string{"basic-memory", "mem0-mcp"}},
		)

		result, err := mock.ProfileList(ctx)

		if err != nil {
			t.Errorf("ProfileList() returned error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("ProfileList() = %d profiles, want 2", len(result))
		}
	})

	t.Run("ProfileCreate records call", func(t *testing.T) {
		mock := NewMockMcpmClient()

		err := mock.ProfileCreate(ctx, "new-profile")

		if err != nil {
			t.Errorf("ProfileCreate() returned error: %v", err)
		}
		if !mock.AssertCalled("ProfileCreate", "new-profile") {
			t.Error("ProfileCreate was not recorded")
		}
	})

	t.Run("ProfileEdit records call with options", func(t *testing.T) {
		mock := NewMockMcpmClient()
		opts := ProfileEditOpts{
			AddServers:    []string{"new-server"},
			RemoveServers: []string{"old-server"},
		}

		err := mock.ProfileEdit(ctx, "p-pokeedge", opts)

		if err != nil {
			t.Errorf("ProfileEdit() returned error: %v", err)
		}
		if mock.CallCount("ProfileEdit") != 1 {
			t.Error("ProfileEdit was not recorded")
		}
	})
}

func TestMockMcpmClient_ClientOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("ClientList returns configured clients", func(t *testing.T) {
		mock := NewMockMcpmClient().WithClients(
			ClientInfo{Name: "claude-code", Installed: true},
			ClientInfo{Name: "claude-desktop", Installed: true},
		)

		result, err := mock.ClientList(ctx)

		if err != nil {
			t.Errorf("ClientList() returned error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("ClientList() = %d clients, want 2", len(result))
		}
	})
}

func TestMockMcpmClient_ConfigOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("ConfigSet and ConfigGet work together", func(t *testing.T) {
		mock := NewMockMcpmClient()

		err := mock.ConfigSet(ctx, "default_node", "/usr/bin/node")
		if err != nil {
			t.Errorf("ConfigSet() returned error: %v", err)
		}

		value, err := mock.ConfigGet(ctx, "default_node")
		if err != nil {
			t.Errorf("ConfigGet() returned error: %v", err)
		}
		if value != "/usr/bin/node" {
			t.Errorf("ConfigGet() = %q, want %q", value, "/usr/bin/node")
		}
	})

	t.Run("ConfigGet returns error for missing key", func(t *testing.T) {
		mock := NewMockMcpmClient()

		_, err := mock.ConfigGet(ctx, "missing_key")

		if err == nil {
			t.Error("ConfigGet() should return error for missing key")
		}
	})
}
