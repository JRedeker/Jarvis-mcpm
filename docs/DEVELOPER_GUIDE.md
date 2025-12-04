# Jarvis Developer Guide

This guide covers how to develop, test, and contribute to Jarvis.

---

## Architecture Overview

Jarvis is an MCP (Model Context Protocol) server written in Go that provides intelligent tool management for AI agents.

### Component Structure

```
Jarvis-Dev/
├── Jarvis/                 # Go MCP server
│   ├── main.go            # Entry point, server setup
│   ├── handlers/          # Tool handlers (DI-based)
│   │   ├── handlers.go    # Handler implementations
│   │   ├── server.go      # MCP tool definitions
│   │   └── registry.go    # Handler registration
│   ├── testing/           # Test utilities
│   │   ├── mocks/         # Mock implementations
│   │   ├── helpers/       # Test helpers
│   │   └── fixtures/      # Test data
│   └── smoketests/        # Integration tests
├── MCPM/                   # Node.js CLI
├── mcpm_source/            # Python reference (archived)
├── scripts/                # Management scripts
└── docs/                   # Documentation
```

### Design Principles

1. **Dependency Injection**: All handlers receive dependencies via interfaces
2. **Testability**: Mocks available for all external dependencies
3. **Single Responsibility**: Each handler does one thing well
4. **Clean Output**: Strip ANSI codes, format as Markdown

---

## Development Setup

### Prerequisites

- Go 1.23+
- Node.js 18+
- Docker & Docker Compose
- Git

### Quick Start

```bash
# Clone
git clone https://github.com/JRedeker/Jarvis-mcpm.git
cd Jarvis-mcpm

# Build Jarvis
cd Jarvis && go build -o jarvis .

# Install MCPM
cd ../MCPM && npm install && npm link

# Start infrastructure
./scripts/manage-mcp.sh start

# Run tests
./scripts/manage-mcp.sh test
```

---

## Handler Architecture

### Handler Interface Pattern

All handlers follow the dependency injection pattern:

```go
// handlers/handlers.go

// Interfaces for dependencies
type McpmRunner interface {
    Run(args ...string) (string, error)
}

type DockerRunner interface {
    ComposeUp(ctx context.Context, services ...string) error
    ComposeDown(ctx context.Context) error
    ComposeRestart(ctx context.Context, services ...string) error
    ComposePs(ctx context.Context) ([]ContainerStatus, error)
    ExecSupervisorctl(ctx context.Context, action, target string) (string, error)
}

type GitRunner interface {
    Status(ctx context.Context) (string, error)
    Diff(ctx context.Context, staged bool) (string, error)
    Init(ctx context.Context) error
}

type FileSystem interface {
    Stat(name string) (os.FileInfo, error)
    ReadFile(name string) ([]byte, error)
    WriteFile(name string, data []byte, perm os.FileMode) error
    MkdirAll(path string, perm os.FileMode) error
    ReadDir(name string) ([]os.DirEntry, error)
    Getwd() (string, error)
}

type CommandRunner interface {
    Run(ctx context.Context, name string, args ...string) (string, error)
    RunInDir(ctx context.Context, dir, name string, args ...string) (string, error)
    StartBackground(ctx context.Context, name string, args ...string) (Process, error)
}

type ProcessManager interface {
    Register(name string, proc Process)
    Get(name string) (Process, bool)
    Remove(name string) bool
    List() []string
}
```

### Handler Structure

```go
// Handler holds all dependencies
type Handler struct {
    Mcpm        McpmRunner
    Docker      DockerRunner
    Git         GitRunner
    FS          FileSystem
    Cmd         CommandRunner
    Processes   ProcessManager
    ExitProcess ExitFunc
}

// Example handler implementation
func (h *Handler) CheckStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    output, _ := h.Mcpm.Run("doctor")

    if strings.Contains(output, "All systems healthy") {
        output += "\n\n🚀 **ALL SYSTEMS GO!** 🚀"
    }

    return mcp.NewToolResultText(output), nil
}
```

### Creating Production Handlers

```go
// handlers/server.go

func CreateProductionHandler() *Handler {
    return NewHandler(
        &RealMcpmRunner{},
        &RealDockerRunner{},
        &RealGitRunner{},
        &RealFileSystem{},
    )
}
```

---

## Adding a New Tool

### Step 1: Define the Handler

Add to `handlers/handlers.go`:

```go
// MyNewTool handles the my_new_tool command
func (h *Handler) MyNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args, ok := request.Params.Arguments.(map[string]interface{})
    if !ok {
        return mcp.NewToolResultError("invalid arguments"), nil
    }

    // Get required parameter
    name, ok := args["name"].(string)
    if !ok || strings.TrimSpace(name) == "" {
        return mcp.NewToolResultError("name is required"), nil
    }

    // Implement tool logic
    result := fmt.Sprintf("Processing: %s", name)

    return mcp.NewToolResultText(result), nil
}
```

### Step 2: Register the Tool Definition

Add to `handlers/server.go` in `GetToolDefinitions()`:

```go
{
    Tool: mcp.NewTool("my_new_tool",
        mcp.WithDescription("Does something useful with the given name"),
        mcp.WithString("name",
            mcp.Description("The name to process"),
            mcp.Required(),
        ),
        mcp.WithBoolean("verbose",
            mcp.Description("Enable verbose output"),
        ),
    ),
    Handler: h.MyNewTool,
},
```

### Step 3: Add to Registry

Add to `handlers/registry.go` in `RegisterAllHandlers()`:

```go
reg.Register("my_new_tool", func(h *Handler) ToolHandler {
    return h.MyNewTool
})
```

### Step 4: Write Tests

Create test in `handlers/handlers_test.go`:

```go
func TestMyNewTool_Success(t *testing.T) {
    h := setupTestHandler()

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{
        "name": "test-value",
    }

    result, err := h.MyNewTool(context.Background(), request)

    require.NoError(t, err)
    require.NotNil(t, result)
    assert.Contains(t, getResultText(result), "test-value")
}

func TestMyNewTool_RequiresName(t *testing.T) {
    h := setupTestHandler()

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{}

    result, err := h.MyNewTool(context.Background(), request)

    require.NoError(t, err)
    assert.True(t, result.IsError)
}
```

### Step 5: Update Documentation

Run the API doc generator:

```bash
./scripts/generate-api-docs.sh
```

---

## Testing

### Test Structure

```
Jarvis/
├── handlers/
│   └── handlers_test.go    # Handler unit tests
├── testing/
│   ├── mocks/
│   │   ├── docker_mock.go  # MockDockerClient
│   │   └── mcpm_mock.go    # MockMcpmClient
│   ├── helpers/
│   │   └── helpers.go      # Test utilities
│   └── fixtures/           # Test data files
├── smoketests/
│   └── smoke_test.go       # Integration tests
├── main_test.go            # Main package tests
└── go.mod
```

### Running Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out  # Interactive report

# Specific package
go test -v ./handlers/

# Specific test
go test -v -run TestCheckStatus ./handlers/

# With race detection
go test -race ./...
```

### Using Mocks

```go
import "jarvis/testing/mocks"

func TestWithMocks(t *testing.T) {
    // Create mock
    mockMcpm := mocks.NewMockMcpmClient()
    mockMcpm.SetDoctorResponse("All systems healthy", nil)

    // Create handler with mock
    h := handlers.NewHandler(
        mockMcpm,
        &mocks.MockDockerClient{},
        &mocks.MockGitRunner{},
        &mocks.MockFileSystem{},
    )

    // Test
    result, err := h.CheckStatus(context.Background(), mcp.CallToolRequest{})

    // Assert mock was called
    mockMcpm.AssertCalled(t, "doctor")
}
```

### Mock Capabilities

**MockMcpmClient:**
- `SetDoctorResponse(output, err)` - Configure health check response
- `SetInstallResult(server, output, err)` - Configure install response
- `SetSearchResults(query, results)` - Configure search results
- `WithProfiles(profiles)` - Set available profiles
- `AssertCalled(t, method)` - Verify method was called

**MockDockerClient:**
- `SetContainerStatuses(statuses)` - Configure container states
- `SetHealthy(healthy)` - Set overall health
- `SetRestartError(err)` - Simulate restart failure

**MockFileSystem:**
- `SetFile(path, content)` - Create virtual file
- `SetDir(path, entries)` - Create virtual directory
- `SetCwd(path)` - Set working directory

---

## Code Style

### Go Formatting

```bash
# Format all Go code
gofmt -w .

# Lint
golangci-lint run
```

### Commit Messages

Follow semantic commit format:

```
feat: Add new my_new_tool handler
fix: Handle empty server names in install_server
docs: Update API reference
test: Add coverage for ManageProfile
refactor: Extract common validation to helper
chore: Update dependencies
```

### PR Guidelines

1. Write tests for new functionality
2. Maintain or improve coverage
3. Update API documentation
4. Follow existing patterns
5. Keep PRs focused and small

---

## Debugging

### Logging

```go
import "log"

func (h *Handler) MyTool(...) {
    log.Printf("MyTool called with args: %v", args)
    // ...
}
```

Logs go to `logs/jarvis.log`.

### Debug Environment

```bash
export MCPM_DEBUG=1
export MCPM_NON_INTERACTIVE=true
```

### Common Issues

1. **Handler not registered**: Check `registry.go` has entry
2. **Tool not appearing**: Verify `server.go` definition
3. **Tests failing**: Check mock configuration
4. **Build errors**: Run `go mod tidy`

---

## Release Process

1. Update version in `main.go`
2. Run full test suite: `./scripts/manage-mcp.sh test`
3. Generate docs: `./scripts/generate-api-docs.sh`
4. Create PR with changelog
5. After merge, create GitHub release

---

## Resources

- [MCP Protocol Spec](https://modelcontextprotocol.io/)
- [mcp-go Library](https://github.com/mark3labs/mcp-go)
- [API Reference](API_REFERENCE.md)
- [Troubleshooting Guide](TROUBLESHOOTING.md)
