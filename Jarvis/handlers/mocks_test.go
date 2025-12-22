package handlers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"
)

// MockMcpmRunner is a mock implementation of McpmRunner for testing
type MockMcpmRunner struct {
	mu sync.Mutex

	// Configured responses (command -> output)
	Responses map[string]string
	Errors    map[string]error

	// Call tracking
	Calls []MockCall
}

// MockCall represents a recorded method call
type MockCall struct {
	Method string
	Args   []interface{}
}

// NewMockMcpmRunner creates a new mock MCPM runner
func NewMockMcpmRunner() *MockMcpmRunner {
	return &MockMcpmRunner{
		Responses: make(map[string]string),
		Errors:    make(map[string]error),
		Calls:     make([]MockCall, 0),
	}
}

// Run implements McpmRunner.Run
func (m *MockMcpmRunner) Run(args ...string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	callArgs := make([]interface{}, len(args))
	for i, a := range args {
		callArgs[i] = a
	}
	m.Calls = append(m.Calls, MockCall{Method: "Run", Args: callArgs})

	// Build key from args
	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	// Check for configured error
	if err, ok := m.Errors[key]; ok {
		return m.Responses[key], err
	}

	// Return configured response
	if resp, ok := m.Responses[key]; ok {
		return resp, nil
	}

	// Default response
	return fmt.Sprintf("Mock response for: %v", args), nil
}

// WithResponse configures a response for a command
func (m *MockMcpmRunner) WithResponse(command, response string) *MockMcpmRunner {
	m.Responses[command] = response
	return m
}

// WithError configures an error for a command
func (m *MockMcpmRunner) WithError(command string, err error) *MockMcpmRunner {
	m.Errors[command] = err
	return m
}

// AssertCalled checks if a method was called with specific args
func (m *MockMcpmRunner) AssertCalled(command string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, call := range m.Calls {
		if len(call.Args) > 0 && call.Args[0] == command {
			return true
		}
	}
	return false
}

// CallCount returns the number of times Run was called with a specific command
func (m *MockMcpmRunner) CallCount(command string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for _, call := range m.Calls {
		if len(call.Args) > 0 && call.Args[0] == command {
			count++
		}
	}
	return count
}

// MockDockerRunner is a mock implementation of DockerRunner for testing
type MockDockerRunner struct {
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

// NewMockDockerRunner creates a new mock Docker runner
func NewMockDockerRunner() *MockDockerRunner {
	return &MockDockerRunner{
		SupervisorctlOutput: make(map[string]string),
		Calls:               make([]MockCall, 0),
	}
}

func (m *MockDockerRunner) recordCall(method string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: method, Args: args})
}

func (m *MockDockerRunner) ComposeUp(ctx context.Context, services ...string) error {
	m.recordCall("ComposeUp", services)
	if m.ComposeUpError != nil {
		return m.ComposeUpError
	}
	m.mu.Lock()
	m.ContainersRunning = true
	m.mu.Unlock()
	return nil
}

func (m *MockDockerRunner) ComposeDown(ctx context.Context) error {
	m.recordCall("ComposeDown")
	if m.ComposeDownError != nil {
		return m.ComposeDownError
	}
	m.mu.Lock()
	m.ContainersRunning = false
	m.mu.Unlock()
	return nil
}

func (m *MockDockerRunner) ComposeRestart(ctx context.Context, services ...string) error {
	m.recordCall("ComposeRestart", services)
	if m.ComposeRestartError != nil {
		return m.ComposeRestartError
	}
	m.mu.Lock()
	m.RestartCount++
	m.mu.Unlock()
	return nil
}

func (m *MockDockerRunner) ComposePs(ctx context.Context) ([]ContainerStatus, error) {
	m.recordCall("ComposePs")
	if m.ComposePsError != nil {
		return nil, m.ComposePsError
	}
	if m.ComposePsResponse != nil {
		return m.ComposePsResponse, nil
	}
	if m.ContainersRunning {
		return []ContainerStatus{
			{Name: "mcp-postgres", Status: "Up", Health: "healthy", Running: true},
			{Name: "mcp-qdrant", Status: "Up", Health: "healthy", Running: true},
			{Name: "mcp-daemon", Status: "Up", Health: "healthy", Running: true},
		}, nil
	}
	return []ContainerStatus{}, nil
}

func (m *MockDockerRunner) ExecSupervisorctl(ctx context.Context, action, target string) (string, error) {
	m.recordCall("ExecSupervisorctl", action, target)
	if m.SupervisorctlError != nil {
		return "", m.SupervisorctlError
	}
	key := fmt.Sprintf("%s:%s", action, target)
	if output, ok := m.SupervisorctlOutput[key]; ok {
		return output, nil
	}
	switch action {
	case "restart":
		return fmt.Sprintf("%s: stopped\n%s: started", target, target), nil
	case "status":
		return fmt.Sprintf("%s                          RUNNING   pid 123, uptime 1:00:00", target), nil
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}
}

// WithRunningContainers sets up containers as running
func (m *MockDockerRunner) WithRunningContainers() *MockDockerRunner {
	m.ContainersRunning = true
	return m
}

// WithComposeRestartError configures an error for ComposeRestart
func (m *MockDockerRunner) WithComposeRestartError(err error) *MockDockerRunner {
	m.ComposeRestartError = err
	return m
}

// WithSupervisorctlOutput configures output for a specific action/target
func (m *MockDockerRunner) WithSupervisorctlOutput(action, target, output string) *MockDockerRunner {
	key := fmt.Sprintf("%s:%s", action, target)
	m.SupervisorctlOutput[key] = output
	return m
}

// WithSupervisorctlError configures an error for supervisorctl
func (m *MockDockerRunner) WithSupervisorctlError(err error) *MockDockerRunner {
	m.SupervisorctlError = err
	return m
}

// ComposeBuild implements DockerRunner.ComposeBuild
func (m *MockDockerRunner) ComposeBuild(ctx context.Context, noCache bool, services ...string) error {
	m.recordCall("ComposeBuild", noCache, services)
	if m.ComposeBuildError != nil {
		return m.ComposeBuildError
	}
	m.mu.Lock()
	m.BuildCount++
	m.mu.Unlock()
	return nil
}

// ComposeStop implements DockerRunner.ComposeStop
func (m *MockDockerRunner) ComposeStop(ctx context.Context, services ...string) error {
	m.recordCall("ComposeStop", services)
	if m.ComposeStopError != nil {
		return m.ComposeStopError
	}
	m.mu.Lock()
	m.ContainersRunning = false
	m.mu.Unlock()
	return nil
}

// ComposeStart implements DockerRunner.ComposeStart
func (m *MockDockerRunner) ComposeStart(ctx context.Context, services ...string) error {
	m.recordCall("ComposeStart", services)
	if m.ComposeStartError != nil {
		return m.ComposeStartError
	}
	m.mu.Lock()
	m.ContainersRunning = true
	m.mu.Unlock()
	return nil
}

// ComposeLogs implements DockerRunner.ComposeLogs
func (m *MockDockerRunner) ComposeLogs(ctx context.Context, service string, lines int) (string, error) {
	m.recordCall("ComposeLogs", service, lines)
	if m.ComposeLogsError != nil {
		return "", m.ComposeLogsError
	}
	if m.ComposeLogsOutput != "" {
		return m.ComposeLogsOutput, nil
	}
	return fmt.Sprintf("[%s] Mock log output line 1\n[%s] Mock log output line 2\n", service, service), nil
}

// WithComposeBuildError configures an error for ComposeBuild
func (m *MockDockerRunner) WithComposeBuildError(err error) *MockDockerRunner {
	m.ComposeBuildError = err
	return m
}

// WithComposeStopError configures an error for ComposeStop
func (m *MockDockerRunner) WithComposeStopError(err error) *MockDockerRunner {
	m.ComposeStopError = err
	return m
}

// WithComposeStartError configures an error for ComposeStart
func (m *MockDockerRunner) WithComposeStartError(err error) *MockDockerRunner {
	m.ComposeStartError = err
	return m
}

// WithComposeLogsOutput configures output for ComposeLogs
func (m *MockDockerRunner) WithComposeLogsOutput(output string) *MockDockerRunner {
	m.ComposeLogsOutput = output
	return m
}

// WithComposeLogsError configures an error for ComposeLogs
func (m *MockDockerRunner) WithComposeLogsError(err error) *MockDockerRunner {
	m.ComposeLogsError = err
	return m
}

// MockGitRunner is a mock implementation of GitRunner for testing
type MockGitRunner struct {
	mu sync.Mutex

	StatusOutput string
	StatusError  error
	DiffOutput   string
	DiffError    error
	InitError    error

	Calls []MockCall
}

// NewMockGitRunner creates a new mock Git runner
func NewMockGitRunner() *MockGitRunner {
	return &MockGitRunner{
		Calls: make([]MockCall, 0),
	}
}

func (m *MockGitRunner) Status(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "Status"})
	return m.StatusOutput, m.StatusError
}

func (m *MockGitRunner) Diff(ctx context.Context, staged bool) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "Diff", Args: []interface{}{staged}})
	return m.DiffOutput, m.DiffError
}

func (m *MockGitRunner) Init(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "Init"})
	return m.InitError
}

// WithStatus configures the Status output
func (m *MockGitRunner) WithStatus(output string) *MockGitRunner {
	m.StatusOutput = output
	return m
}

// WithDiff configures the Diff output
func (m *MockGitRunner) WithDiff(output string) *MockGitRunner {
	m.DiffOutput = output
	return m
}

// WithStatusError configures an error for Status
func (m *MockGitRunner) WithStatusError(err error) *MockGitRunner {
	m.StatusError = err
	return m
}

// MockFileSystem is a mock implementation of FileSystem for testing
type MockFileSystem struct {
	mu sync.Mutex

	Files    map[string][]byte
	DirFiles map[string][]os.DirEntry
	Cwd      string
	StatErr  map[string]error

	Calls []MockCall
}

// NewMockFileSystem creates a new mock file system
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:    make(map[string][]byte),
		DirFiles: make(map[string][]os.DirEntry),
		StatErr:  make(map[string]error),
		Cwd:      "/home/test/project",
		Calls:    make([]MockCall, 0),
	}
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "Stat", Args: []interface{}{name}})

	if err, ok := m.StatErr[name]; ok {
		return nil, err
	}
	if _, ok := m.Files[name]; ok {
		return &mockFileInfo{name: name, isDir: false}, nil
	}
	if _, ok := m.DirFiles[name]; ok {
		return &mockFileInfo{name: name, isDir: true}, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystem) ReadFile(name string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "ReadFile", Args: []interface{}{name}})

	if content, ok := m.Files[name]; ok {
		return content, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "WriteFile", Args: []interface{}{name, data, perm}})
	m.Files[name] = data
	return nil
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "MkdirAll", Args: []interface{}{path, perm}})
	return nil
}

func (m *MockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "ReadDir", Args: []interface{}{name}})

	if entries, ok := m.DirFiles[name]; ok {
		return entries, nil
	}
	return []os.DirEntry{}, nil
}

func (m *MockFileSystem) Getwd() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: "Getwd"})
	return m.Cwd, nil
}

// WithFile adds a file to the mock filesystem
func (m *MockFileSystem) WithFile(path string, content []byte) *MockFileSystem {
	m.Files[path] = content
	return m
}

// WithDir adds a directory with entries to the mock filesystem
func (m *MockFileSystem) WithDir(path string, entries []os.DirEntry) *MockFileSystem {
	m.DirFiles[path] = entries
	return m
}

// WithCwd sets the current working directory
func (m *MockFileSystem) WithCwd(cwd string) *MockFileSystem {
	m.Cwd = cwd
	return m
}

// WithStatError configures an error for Stat on a specific path
func (m *MockFileSystem) WithStatError(path string, err error) *MockFileSystem {
	m.StatErr[path] = err
	return m
}

// mockFileInfo implements os.FileInfo for testing
type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// mockDirEntry implements os.DirEntry for testing
type MockDirEntry struct {
	EntryName  string
	EntryIsDir bool
}

func (m *MockDirEntry) Name() string               { return m.EntryName }
func (m *MockDirEntry) IsDir() bool                { return m.EntryIsDir }
func (m *MockDirEntry) Type() fs.FileMode          { return 0 }
func (m *MockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

// MockCommandRunner is a mock implementation of CommandRunner for testing
type MockCommandRunner struct {
	Output string
	Error  error
	Calls  []MockCall
}

func (m *MockCommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	m.Calls = append(m.Calls, MockCall{Method: "Run", Args: []interface{}{name, args}})
	return m.Output, m.Error
}

func (m *MockCommandRunner) RunInDir(ctx context.Context, dir, name string, args ...string) (string, error) {
	m.Calls = append(m.Calls, MockCall{Method: "RunInDir", Args: []interface{}{dir, name, args}})
	return m.Output, m.Error
}

func (m *MockCommandRunner) StartBackground(ctx context.Context, name string, args ...string) (Process, error) {
	m.Calls = append(m.Calls, MockCall{Method: "StartBackground", Args: []interface{}{name, args}})
	return nil, m.Error
}
