# Jarvis Developer Guide

This guide covers how to develop, test, and contribute to Jarvis.

**Version:** 3.0 (Tool Consolidation)

---

## Architecture Overview

Jarvis is an MCP (Model Context Protocol) server written in Go that provides intelligent tool management for AI agents.

### v3.0 Consolidated Tools

Jarvis v3.0 consolidated 24 tools into 8 action-based tools for 52% context token reduction:

| Tool | Actions |
|:-----|:--------|
| `jarvis_check_status` | (single purpose) |
| `jarvis_server` | list, info, install, uninstall, search, edit, create, usage |
| `jarvis_profile` | list, create, edit, delete, suggest, restart |
| `jarvis_client` | list, edit, import, config |
| `jarvis_config` | get, set, list, migrate |
| `jarvis_project` | analyze, diff, devops |
| `jarvis_system` | bootstrap, restart, restart_infra |
| `jarvis_share` | start, stop, list |

### Component Structure

```
Jarvis-Dev/
├── Jarvis/                 # Go MCP server
│   ├── main.go            # Entry point, server setup
│   ├── handlers/          # Tool handlers (DI-based)
│   │   ├── consolidated.go # v3.0 consolidated handlers
│   │   ├── handlers.go    # Legacy handler implementations
│   │   ├── server.go      # MCP tool definitions (8 tools)
│   │   └── registry.go    # Handler registration
│   ├── testing/           # Test utilities
│   │   ├── mocks/         # Mock implementations
│   │   ├── helpers/       # Test helpers
│   │   └── fixtures/      # Test data
│   └── smoketests/        # Integration tests
├── MCPM/                   # Node.js CLI + API server
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

## Adding a New Action to an Existing Tool

In v3.0, we use consolidated action-based tools. To add a new action:

### Step 1: Add the Action Handler

Add to `handlers/consolidated.go` in the appropriate switch statement:

```go
// In the Server handler for a new server action
case "my_action":
    name, _ := args["name"].(string)
    if name == "" {
        return mcp.NewToolResultError("name is required for my_action"), nil
    }
    output, _ := h.Mcpm.Run("my-action", name)
    return mcp.NewToolResultText(output), nil
```

### Step 2: Update Tool Definition

Update the action enum in `handlers/server.go`:

```go
mcp.WithString("action",
    mcp.Description("Action: list, info, install, uninstall, search, edit, create, usage, my_action"),
    mcp.Required(),
),
```

---

## Adding a New Consolidated Tool

If you need an entirely new tool category (rare), follow this pattern:

### Step 1: Define the Handler

Add to `handlers/consolidated.go`:

```go
// MyCategory handles jarvis_mycategory actions
func (h *Handler) MyCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args, ok := request.Params.Arguments.(map[string]interface{})
    if !ok {
        return mcp.NewToolResultError("invalid arguments"), nil
    }

    action, _ := args["action"].(string)
    if action == "" {
        return mcp.NewToolResultError("action is required"), nil
    }

    switch action {
    case "action1":
        // Implementation
        return mcp.NewToolResultText("Result"), nil
    case "action2":
        // Implementation
        return mcp.NewToolResultText("Result"), nil
    default:
        return mcp.NewToolResultError("unknown action: " + action), nil
    }
}
```

### Step 2: Register the Tool Definition

Add to `handlers/server.go` in `GetToolDefinitions()`:

```go
{
    Tool: mcp.NewTool("jarvis_mycategory",
        mcp.WithDescription("Manage my category with actions: action1, action2"),
        mcp.WithString("action",
            mcp.Description("Action: action1, action2"),
            mcp.Required(),
        ),
        mcp.WithString("name",
            mcp.Description("Name parameter for actions"),
        ),
    ),
    Handler: h.MyCategory,
},
```

### Step 3: Add to Registry

Add to `handlers/registry.go` in `RegisterAllHandlers()`:

```go
reg.Register("jarvis_mycategory", func(h *Handler) ToolHandler {
    return h.MyCategory
})
```

### Step 4: Write Tests

Create test in `handlers/handlers_test.go`:

```go
func TestMyCategory_Action1_Success(t *testing.T) {
    h := setupTestHandler()

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{
        "action": "action1",
        "name":   "test-value",
    }

    result, err := h.MyCategory(context.Background(), request)

    require.NoError(t, err)
    require.NotNil(t, result)
    assert.Contains(t, getResultText(result), "test-value")
}

func TestMyCategory_RequiresAction(t *testing.T) {
    h := setupTestHandler()

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{
        "name": "test-value",
    }

    result, err := h.MyCategory(context.Background(), request)

    require.NoError(t, err)
    assert.True(t, result.IsError)
    assert.Contains(t, getResultText(result), "action is required")
}

func TestMyCategory_UnknownAction(t *testing.T) {
    h := setupTestHandler()

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{
        "action": "invalid",
    }

    result, err := h.MyCategory(context.Background(), request)

    require.NoError(t, err)
    assert.True(t, result.IsError)
    assert.Contains(t, getResultText(result), "unknown action")
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

    // Test jarvis_check_status (v3.0 consolidated tool)
    result, err := h.CheckStatus(context.Background(), mcp.CallToolRequest{})

    // Assert mock was called
    mockMcpm.AssertCalled(t, "doctor")
}

// Example: Testing jarvis_server consolidated handler
func TestServer_Install(t *testing.T) {
    mockMcpm := mocks.NewMockMcpmClient()
    mockMcpm.SetInstallResult("context7", "Installed successfully", nil)

    h := handlers.NewHandler(mockMcpm, nil, nil, nil)

    request := mcp.CallToolRequest{}
    request.Params.Arguments = map[string]interface{}{
        "action": "install",
        "name":   "context7",
    }

    result, err := h.Server(context.Background(), request)

    require.NoError(t, err)
    assert.Contains(t, getResultText(result), "Installed successfully")
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
