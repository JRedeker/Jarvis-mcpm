package mocks

import (
	"context"
	"errors"
	"testing"
)

func TestMockDockerClient_ComposeUp(t *testing.T) {
	ctx := context.Background()

	t.Run("succeeds by default", func(t *testing.T) {
		mock := NewMockDockerClient()

		err := mock.ComposeUp(ctx)

		if err != nil {
			t.Errorf("ComposeUp() returned error: %v", err)
		}
		if !mock.ContainersRunning {
			t.Error("ComposeUp() should set ContainersRunning to true")
		}
		if !mock.AssertCalled("ComposeUp") {
			t.Error("ComposeUp() was not recorded")
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		expectedErr := errors.New("docker daemon not running")
		mock := NewMockDockerClient().WithComposeUpError(expectedErr)

		err := mock.ComposeUp(ctx)

		if err != expectedErr {
			t.Errorf("ComposeUp() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("records services argument", func(t *testing.T) {
		mock := NewMockDockerClient()

		err := mock.ComposeUp(ctx, "postgres", "qdrant")

		if err != nil {
			t.Errorf("ComposeUp() returned error: %v", err)
		}
		if mock.CallCount("ComposeUp") != 1 {
			t.Errorf("ComposeUp() call count = %d, want 1", mock.CallCount("ComposeUp"))
		}
	})
}

func TestMockDockerClient_ComposeDown(t *testing.T) {
	ctx := context.Background()

	t.Run("succeeds and stops containers", func(t *testing.T) {
		mock := NewMockDockerClient().WithRunningContainers()

		err := mock.ComposeDown(ctx)

		if err != nil {
			t.Errorf("ComposeDown() returned error: %v", err)
		}
		if mock.ContainersRunning {
			t.Error("ComposeDown() should set ContainersRunning to false")
		}
	})
}

func TestMockDockerClient_ComposeRestart(t *testing.T) {
	ctx := context.Background()

	t.Run("succeeds and increments restart count", func(t *testing.T) {
		mock := NewMockDockerClient()

		err := mock.ComposeRestart(ctx)

		if err != nil {
			t.Errorf("ComposeRestart() returned error: %v", err)
		}
		if mock.RestartCount != 1 {
			t.Errorf("RestartCount = %d, want 1", mock.RestartCount)
		}

		// Second restart
		mock.ComposeRestart(ctx)
		if mock.RestartCount != 2 {
			t.Errorf("RestartCount = %d, want 2", mock.RestartCount)
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		expectedErr := errors.New("restart failed")
		mock := NewMockDockerClient().WithComposeRestartError(expectedErr)

		err := mock.ComposeRestart(ctx)

		if err != expectedErr {
			t.Errorf("ComposeRestart() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestMockDockerClient_ComposePs(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty when containers stopped", func(t *testing.T) {
		mock := NewMockDockerClient().WithStoppedContainers()

		result, err := mock.ComposePs(ctx)

		if err != nil {
			t.Errorf("ComposePs() returned error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("ComposePs() = %d containers, want 0", len(result))
		}
	})

	t.Run("returns default containers when running", func(t *testing.T) {
		mock := NewMockDockerClient().WithRunningContainers()

		result, err := mock.ComposePs(ctx)

		if err != nil {
			t.Errorf("ComposePs() returned error: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("ComposePs() = %d containers, want 3", len(result))
		}
	})

	t.Run("returns configured healthy containers", func(t *testing.T) {
		mock := NewMockDockerClient().WithHealthyContainers()

		result, err := mock.ComposePs(ctx)

		if err != nil {
			t.Errorf("ComposePs() returned error: %v", err)
		}

		for _, container := range result {
			if container.Health != "healthy" {
				t.Errorf("Container %s health = %q, want %q", container.Name, container.Health, "healthy")
			}
		}
	})

	t.Run("returns configured unhealthy container", func(t *testing.T) {
		mock := NewMockDockerClient().WithUnhealthyContainer("mcp-daemon")

		result, err := mock.ComposePs(ctx)

		if err != nil {
			t.Errorf("ComposePs() returned error: %v", err)
		}

		foundUnhealthy := false
		for _, container := range result {
			if container.Name == "mcp-daemon" && container.Health == "unhealthy" {
				foundUnhealthy = true
				break
			}
		}
		if !foundUnhealthy {
			t.Error("Expected mcp-daemon to be unhealthy")
		}
	})
}

func TestMockDockerClient_ExecSupervisorctl(t *testing.T) {
	ctx := context.Background()

	t.Run("restart returns success message", func(t *testing.T) {
		mock := NewMockDockerClient()

		output, err := mock.ExecSupervisorctl(ctx, "restart", "mcpm-memory")

		if err != nil {
			t.Errorf("ExecSupervisorctl() returned error: %v", err)
		}
		if output == "" {
			t.Error("ExecSupervisorctl() returned empty output")
		}
		if !mock.AssertCalled("ExecSupervisorctl", "restart") {
			t.Error("ExecSupervisorctl() was not recorded correctly")
		}
	})

	t.Run("status returns running message", func(t *testing.T) {
		mock := NewMockDockerClient()

		output, err := mock.ExecSupervisorctl(ctx, "status", "mcpm-memory")

		if err != nil {
			t.Errorf("ExecSupervisorctl() returned error: %v", err)
		}
		if output == "" {
			t.Error("ExecSupervisorctl() returned empty output")
		}
	})

	t.Run("returns configured output", func(t *testing.T) {
		mock := NewMockDockerClient().WithSupervisorctlOutput("restart", "all", "all processes restarted")

		output, err := mock.ExecSupervisorctl(ctx, "restart", "all")

		if err != nil {
			t.Errorf("ExecSupervisorctl() returned error: %v", err)
		}
		if output != "all processes restarted" {
			t.Errorf("ExecSupervisorctl() output = %q, want %q", output, "all processes restarted")
		}
	})

	t.Run("returns error for unknown action", func(t *testing.T) {
		mock := NewMockDockerClient()

		_, err := mock.ExecSupervisorctl(ctx, "unknown", "target")

		if err == nil {
			t.Error("ExecSupervisorctl() should return error for unknown action")
		}
	})
}

func TestMockDockerClient_CallTracking(t *testing.T) {
	ctx := context.Background()
	mock := NewMockDockerClient()

	// Make some calls
	mock.ComposeUp(ctx)
	mock.ComposePs(ctx)
	mock.ComposeRestart(ctx)
	mock.ComposeRestart(ctx)
	mock.ComposeDown(ctx)

	t.Run("CallCount works correctly", func(t *testing.T) {
		if mock.CallCount("ComposeUp") != 1 {
			t.Errorf("CallCount(ComposeUp) = %d, want 1", mock.CallCount("ComposeUp"))
		}
		if mock.CallCount("ComposeRestart") != 2 {
			t.Errorf("CallCount(ComposeRestart) = %d, want 2", mock.CallCount("ComposeRestart"))
		}
		if mock.CallCount("ExecSupervisorctl") != 0 {
			t.Errorf("CallCount(ExecSupervisorctl) = %d, want 0", mock.CallCount("ExecSupervisorctl"))
		}
	})

	t.Run("ResetCalls clears history", func(t *testing.T) {
		mock.ResetCalls()

		if mock.CallCount("ComposeUp") != 0 {
			t.Error("ResetCalls did not clear call history")
		}
	})
}
