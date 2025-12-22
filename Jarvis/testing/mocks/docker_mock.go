package mocks

import (
	"context"
	"fmt"
	"sync"
)

// MockDockerClient is a mock implementation of DockerClient for testing
type MockDockerClient struct {
	mu sync.Mutex

	// Configured responses
	ComposeUpError      error
	ComposeDownError    error
	ComposeRestartError error
	ComposePsResponse   []ContainerStatus
	ComposePsError      error
	SupervisorctlOutput map[string]string
	SupervisorctlError  error
	// Phase 1: Enhanced Docker Operations
	ComposeBuildError error
	ComposeStopError  error
	ComposeStartError error
	ComposeLogsOutput string
	ComposeLogsError  error

	// State tracking
	ContainersRunning bool
	RestartCount      int
	BuildCount        int

	// Call tracking
	Calls []MockCall
}

// NewMockDockerClient creates a new mock Docker client with sensible defaults
func NewMockDockerClient() *MockDockerClient {
	return &MockDockerClient{
		SupervisorctlOutput: make(map[string]string),
		Calls:               make([]MockCall, 0),
		ContainersRunning:   false,
	}
}

// recordCall records a method call for later verification
func (m *MockDockerClient) recordCall(method string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: method, Args: args})
}

// ResetCalls clears the call history
func (m *MockDockerClient) ResetCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = make([]MockCall, 0)
}

// AssertCalled verifies a method was called
func (m *MockDockerClient) AssertCalled(method string, args ...interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, call := range m.Calls {
		if call.Method == method {
			if len(args) == 0 {
				return true
			}
			// Simple arg matching for first arg
			if len(call.Args) > 0 && len(args) > 0 {
				if fmt.Sprintf("%v", call.Args[0]) == fmt.Sprintf("%v", args[0]) {
					return true
				}
			}
		}
	}
	return false
}

// CallCount returns the number of times a method was called
func (m *MockDockerClient) CallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for _, call := range m.Calls {
		if call.Method == method {
			count++
		}
	}
	return count
}

// ComposeUp implements DockerClient.ComposeUp
func (m *MockDockerClient) ComposeUp(ctx context.Context, services ...string) error {
	m.recordCall("ComposeUp", services)

	if m.ComposeUpError != nil {
		return m.ComposeUpError
	}

	m.mu.Lock()
	m.ContainersRunning = true
	m.mu.Unlock()

	return nil
}

// ComposeDown implements DockerClient.ComposeDown
func (m *MockDockerClient) ComposeDown(ctx context.Context) error {
	m.recordCall("ComposeDown")

	if m.ComposeDownError != nil {
		return m.ComposeDownError
	}

	m.mu.Lock()
	m.ContainersRunning = false
	m.mu.Unlock()

	return nil
}

// ComposeRestart implements DockerClient.ComposeRestart
func (m *MockDockerClient) ComposeRestart(ctx context.Context, services ...string) error {
	m.recordCall("ComposeRestart", services)

	if m.ComposeRestartError != nil {
		return m.ComposeRestartError
	}

	m.mu.Lock()
	m.RestartCount++
	m.mu.Unlock()

	return nil
}

// ComposePs implements DockerClient.ComposePs
func (m *MockDockerClient) ComposePs(ctx context.Context) ([]ContainerStatus, error) {
	m.recordCall("ComposePs")

	if m.ComposePsError != nil {
		return nil, m.ComposePsError
	}

	if m.ComposePsResponse != nil {
		return m.ComposePsResponse, nil
	}

	// Default response based on running state
	if m.ContainersRunning {
		return []ContainerStatus{
			{Name: "mcp-postgres", Status: "Up", Health: "healthy", Running: true},
			{Name: "mcp-qdrant", Status: "Up", Health: "healthy", Running: true},
			{Name: "mcp-daemon", Status: "Up", Health: "healthy", Running: true},
		}, nil
	}

	return []ContainerStatus{}, nil
}

// ExecSupervisorctl implements DockerClient.ExecSupervisorctl
func (m *MockDockerClient) ExecSupervisorctl(ctx context.Context, action, target string) (string, error) {
	m.recordCall("ExecSupervisorctl", action, target)

	if m.SupervisorctlError != nil {
		return "", m.SupervisorctlError
	}

	key := fmt.Sprintf("%s:%s", action, target)
	if output, ok := m.SupervisorctlOutput[key]; ok {
		return output, nil
	}

	// Default success output
	switch action {
	case "restart":
		return fmt.Sprintf("%s: stopped\n%s: started", target, target), nil
	case "status":
		return fmt.Sprintf("%s                          RUNNING   pid 123, uptime 1:00:00", target), nil
	case "stop":
		return fmt.Sprintf("%s: stopped", target), nil
	case "start":
		return fmt.Sprintf("%s: started", target), nil
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}
}

// Helper methods for setting up mock state

// WithRunningContainers sets the mock to report containers as running
func (m *MockDockerClient) WithRunningContainers() *MockDockerClient {
	m.ContainersRunning = true
	return m
}

// WithStoppedContainers sets the mock to report containers as stopped
func (m *MockDockerClient) WithStoppedContainers() *MockDockerClient {
	m.ContainersRunning = false
	return m
}

// WithContainerStatus sets a specific container status response
func (m *MockDockerClient) WithContainerStatus(containers ...ContainerStatus) *MockDockerClient {
	m.ComposePsResponse = containers
	return m
}

// WithHealthyContainers sets up all containers as healthy
func (m *MockDockerClient) WithHealthyContainers() *MockDockerClient {
	m.ComposePsResponse = []ContainerStatus{
		{Name: "mcp-postgres", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"5432:5432"}},
		{Name: "mcp-qdrant", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"6333:6333", "6334:6334"}},
		{Name: "mcp-daemon", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"6276:6276", "6277:6277", "6278:6278"}},
	}
	m.ContainersRunning = true
	return m
}

// WithUnhealthyContainer sets up a specific container as unhealthy
func (m *MockDockerClient) WithUnhealthyContainer(name string) *MockDockerClient {
	m.ComposePsResponse = []ContainerStatus{
		{Name: "mcp-postgres", Status: "Up 2 hours", Health: "healthy", Running: true},
		{Name: "mcp-qdrant", Status: "Up 2 hours", Health: "healthy", Running: true},
		{Name: name, Status: "Up 2 hours", Health: "unhealthy", Running: true},
	}
	m.ContainersRunning = true
	return m
}

// WithComposeUpError configures ComposeUp to return an error
func (m *MockDockerClient) WithComposeUpError(err error) *MockDockerClient {
	m.ComposeUpError = err
	return m
}

// WithComposeRestartError configures ComposeRestart to return an error
func (m *MockDockerClient) WithComposeRestartError(err error) *MockDockerClient {
	m.ComposeRestartError = err
	return m
}

// WithSupervisorctlOutput configures a specific supervisorctl output
func (m *MockDockerClient) WithSupervisorctlOutput(action, target, output string) *MockDockerClient {
	key := fmt.Sprintf("%s:%s", action, target)
	m.SupervisorctlOutput[key] = output
	return m
}

// ComposeBuild implements DockerClient.ComposeBuild
func (m *MockDockerClient) ComposeBuild(ctx context.Context, noCache bool, services ...string) error {
	m.recordCall("ComposeBuild", noCache, services)

	if m.ComposeBuildError != nil {
		return m.ComposeBuildError
	}

	m.mu.Lock()
	m.BuildCount++
	m.mu.Unlock()

	return nil
}

// ComposeStop implements DockerClient.ComposeStop
func (m *MockDockerClient) ComposeStop(ctx context.Context, services ...string) error {
	m.recordCall("ComposeStop", services)

	if m.ComposeStopError != nil {
		return m.ComposeStopError
	}

	m.mu.Lock()
	m.ContainersRunning = false
	m.mu.Unlock()

	return nil
}

// ComposeStart implements DockerClient.ComposeStart
func (m *MockDockerClient) ComposeStart(ctx context.Context, services ...string) error {
	m.recordCall("ComposeStart", services)

	if m.ComposeStartError != nil {
		return m.ComposeStartError
	}

	m.mu.Lock()
	m.ContainersRunning = true
	m.mu.Unlock()

	return nil
}

// ComposeLogs implements DockerClient.ComposeLogs
func (m *MockDockerClient) ComposeLogs(ctx context.Context, service string, lines int) (string, error) {
	m.recordCall("ComposeLogs", service, lines)

	if m.ComposeLogsError != nil {
		return "", m.ComposeLogsError
	}

	if m.ComposeLogsOutput != "" {
		return m.ComposeLogsOutput, nil
	}

	// Default mock output
	return fmt.Sprintf("[%s] Mock log output line 1\n[%s] Mock log output line 2\n", service, service), nil
}

// WithComposeBuildError configures ComposeBuild to return an error
func (m *MockDockerClient) WithComposeBuildError(err error) *MockDockerClient {
	m.ComposeBuildError = err
	return m
}

// WithComposeStopError configures ComposeStop to return an error
func (m *MockDockerClient) WithComposeStopError(err error) *MockDockerClient {
	m.ComposeStopError = err
	return m
}

// WithComposeStartError configures ComposeStart to return an error
func (m *MockDockerClient) WithComposeStartError(err error) *MockDockerClient {
	m.ComposeStartError = err
	return m
}

// WithComposeLogsOutput configures ComposeLogs to return specific output
func (m *MockDockerClient) WithComposeLogsOutput(output string) *MockDockerClient {
	m.ComposeLogsOutput = output
	return m
}

// WithComposeLogsError configures ComposeLogs to return an error
func (m *MockDockerClient) WithComposeLogsError(err error) *MockDockerClient {
	m.ComposeLogsError = err
	return m
}
