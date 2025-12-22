# Jarvis Smoke Testing System Design

**Version:** 1.0
**Date:** 2025-11-28
**Status:** Design Proposal
**Related:** Phase 3 of Tool Discoverability Enhancement Plan

---

## Problem Statement

**Current Issue:** Configuration problems (like missing API keys) are only discovered when tools fail at runtime, leading to poor agent UX and debugging overhead.

**Example:** Firecrawl returned 404 errors because `FIRECRAWL_API_KEY` was missing - only discovered after attempting to use the tool.

**Goal:** Detect configuration, connectivity, and health issues at boot time (or on-demand) with clear, actionable error reporting.

---

## Design Principles

1. **Fast First** - Smoke tests must be lightweight (<5 seconds total boot time impact)
2. **Non-Blocking** - Critical failures warn but don't prevent Jarvis from starting
3. **Actionable** - Every failure includes fix suggestions
4. **Extensible** - Easy to add tests for new servers
5. **Observable** - Test results logged and accessible via tools
6. **Configurable** - Users can disable tests or adjust thoroughness

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Jarvis Startup                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Smoke Test Orchestrator                     │
│  • Discovers active profiles                            │
│  • Identifies servers to test                           │
│  • Runs tests in parallel (with timeout)                │
│  • Aggregates results                                   │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐
   │ Config  │  │ Connect │  │ Health  │
   │ Tests   │  │ Tests   │  │ Tests   │
   └─────────┘  └─────────┘  └─────────┘
        │            │            │
        └────────────┴────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │   Test Report         │
         │   • Pass/Fail/Skip    │
         │   • Latency metrics   │
         │   • Fix suggestions   │
         └───────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
    ┌──────┐    ┌──────┐    ┌──────┐
    │ Log  │    │ Tool │    │ Boot │
    │ File │    │ API  │    │ Warn │
    └──────┘    └──────┘    └──────┘
```

---

## Test Categories

### 1. Configuration Tests (Always Run)

**Purpose:** Validate required environment variables and configuration files exist

**Examples:**
- API key presence (not validity - just existence)
- Config file syntax validation
- Required directories exist
- File permissions correct

**Speed:** <100ms per server

```go
type ConfigTest struct {
    ServerName    string
    RequiredEnvs  []string  // e.g., ["FIRECRAWL_API_KEY"]
    RequiredFiles []string  // e.g., [".mcpm/servers/foo/config.json"]
    RequiredPerms map[string]os.FileMode
}

func (t *ConfigTest) Run() TestResult {
    // Check env vars exist (don't validate values)
    // Check files exist and are readable
    // Check permissions
}
```

### 2. Connectivity Tests (Opt-in, Default: On)

**Purpose:** Verify network-dependent services are reachable

**Examples:**
- API endpoint reachability (HTTP HEAD request)
- Database connection (lightweight ping)
- Docker container running
- External service responding

**Speed:** <1s per server (with timeout)

```go
type ConnectivityTest struct {
    ServerName  string
    TestType    string  // "http", "tcp", "docker", "process"
    Endpoint    string  // URL, host:port, container name
    Timeout     time.Duration
    ExpectedCode int   // For HTTP tests
}

func (t *ConnectivityTest) Run() TestResult {
    // Lightweight connectivity check
    // For APIs: HEAD request or OPTIONS
    // For DBs: PING command
    // For Docker: container status check
}
```

### 3. Health Tests (Opt-in, Default: Off)

**Purpose:** Validate actual functionality with minimal resource usage

**Examples:**
- API key validity (make actual authenticated request)
- Database query execution
- Search endpoint returns results
- Tool can be invoked successfully

**Speed:** <2s per server (with timeout)

**Note:** These consume API credits/resources, so default OFF

```go
type HealthTest struct {
    ServerName  string
    ToolName    string  // Which tool to test
    TestPayload map[string]interface{}
    Validator   func(response interface{}) error
    CostWarning string  // "Uses 1 API credit"
}

func (t *HealthTest) Run() TestResult {
    // Actually call a tool with minimal payload
    // Validate response structure
    // Check for expected data
}
```

---

## Test Specification Format

### Server-Level Test Definitions

Tests are defined in server metadata and can be extended by users:

```json
// ~/.mcpm/servers/firecrawl/smoke_tests.json
{
  "server": "firecrawl",
  "tests": {
    "config": {
      "enabled": true,
      "checks": [
        {
          "type": "env_var",
          "name": "FIRECRAWL_API_KEY",
          "required": true,
          "pattern": "^fc-[a-f0-9]{32}$",
          "error_message": "FIRECRAWL_API_KEY must start with 'fc-' and be 35 characters",
          "fix_suggestion": "Get an API key from https://firecrawl.dev and run: mcpm edit firecrawl --env FIRECRAWL_API_KEY=your-key"
        }
      ]
    },
    "connectivity": {
      "enabled": true,
      "timeout": "2s",
      "checks": [
        {
          "type": "http",
          "endpoint": "https://api.firecrawl.dev/v1",
          "method": "HEAD",
          "expected_status": [200, 401, 403],
          "error_message": "Firecrawl API is unreachable",
          "fix_suggestion": "Check internet connection or verify Firecrawl service status"
        }
      ]
    },
    "health": {
      "enabled": false,
      "cost_warning": "Uses 1 API credit per test",
      "timeout": "5s",
      "checks": [
        {
          "type": "tool_call",
          "tool": "firecrawl_scrape",
          "payload": {
            "url": "https://example.com",
            "formats": ["markdown"]
          },
          "validator": "response.markdown != null && response.markdown.length > 0",
          "error_message": "Firecrawl scrape failed",
          "fix_suggestion": "Verify API key is valid and has available credits"
        }
      ]
    }
  }
}
```

### Global Test Configuration

```json
// ~/.mcpm/config.json additions
{
  "smoke_tests": {
    "enabled": true,
    "run_on_boot": true,
    "run_on_profile_change": true,
    "test_levels": {
      "config": true,
      "connectivity": true,
      "health": false
    },
    "parallel_execution": true,
    "max_parallel_tests": 5,
    "global_timeout": "10s",
    "fail_fast": false,
    "report_format": "summary",  // "summary", "detailed", "json"
    "exclude_servers": []  // Skip tests for specific servers
  }
}
```

---

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)

**Files to Create:**
- `Jarvis/smoketests/orchestrator.go`
- `Jarvis/smoketests/types.go`
- `Jarvis/smoketests/config_test.go`
- `Jarvis/smoketests/connectivity_test.go`
- `Jarvis/smoketests/health_test.go`
- `Jarvis/smoketests/registry.go`

**Implementation:**

```go
// Jarvis/smoketests/types.go
package smoketests

import "time"

type TestStatus string

const (
    StatusPass    TestStatus = "pass"
    StatusFail    TestStatus = "fail"
    StatusSkip    TestStatus = "skip"
    StatusTimeout TestStatus = "timeout"
)

type TestResult struct {
    ServerName     string        `json:"server_name"`
    TestType       string        `json:"test_type"` // "config", "connectivity", "health"
    TestName       string        `json:"test_name"`
    Status         TestStatus    `json:"status"`
    ErrorMessage   string        `json:"error_message,omitempty"`
    FixSuggestion  string        `json:"fix_suggestion,omitempty"`
    Duration       time.Duration `json:"duration"`
    Timestamp      time.Time     `json:"timestamp"`
    Details        string        `json:"details,omitempty"`
}

type TestSuite interface {
    Run() []TestResult
    Name() string
    Enabled() bool
}

type SmokeTestReport struct {
    TotalServers     int           `json:"total_servers"`
    TestedServers    int           `json:"tested_servers"`
    PassedTests      int           `json:"passed_tests"`
    FailedTests      int           `json:"failed_tests"`
    SkippedTests     int           `json:"skipped_tests"`
    TotalDuration    time.Duration `json:"total_duration"`
    Results          []TestResult  `json:"results"`
    CriticalFailures []TestResult  `json:"critical_failures"`
    Warnings         []string      `json:"warnings"`
    Timestamp        time.Time     `json:"timestamp"`
}
```

```go
// Jarvis/smoketests/orchestrator.go
package smoketests

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type Orchestrator struct {
    config       *Config
    registry     *TestRegistry
    results      []TestResult
    mu           sync.Mutex
}

func NewOrchestrator(config *Config) *Orchestrator {
    return &Orchestrator{
        config:   config,
        registry: NewTestRegistry(),
        results:  []TestResult{},
    }
}

// RunAll executes smoke tests for all active servers
func (o *Orchestrator) RunAll(ctx context.Context, servers []string) *SmokeTestReport {
    startTime := time.Now()

    if !o.config.Enabled {
        return &SmokeTestReport{
            TotalServers: len(servers),
            Results:      []TestResult{{Status: StatusSkip, TestName: "Smoke tests disabled"}},
        }
    }

    // Create context with global timeout
    ctx, cancel := context.WithTimeout(ctx, o.config.GlobalTimeout)
    defer cancel()

    var wg sync.WaitGroup
    semaphore := make(chan struct{}, o.config.MaxParallelTests)

    for _, serverName := range servers {
        // Check if server is excluded
        if contains(o.config.ExcludeServers, serverName) {
            o.addResult(TestResult{
                ServerName: serverName,
                Status:     StatusSkip,
                TestName:   "Server excluded from testing",
            })
            continue
        }

        wg.Add(1)
        go func(server string) {
            defer wg.Done()

            // Rate limiting
            semaphore <- struct{}{}
            defer func() { <-semaphore }()

            o.runServerTests(ctx, server)
        }(serverName)
    }

    wg.Wait()

    return o.generateReport(startTime)
}

func (o *Orchestrator) runServerTests(ctx context.Context, serverName string) {
    // 1. Config Tests (always run if enabled)
    if o.config.TestLevels.Config {
        suite := o.registry.GetConfigTests(serverName)
        if suite != nil && suite.Enabled() {
            results := suite.Run()
            for _, r := range results {
                o.addResult(r)

                // Fail fast if critical config missing
                if o.config.FailFast && r.Status == StatusFail && r.TestType == "config" {
                    return
                }
            }
        }
    }

    // 2. Connectivity Tests
    if o.config.TestLevels.Connectivity {
        suite := o.registry.GetConnectivityTests(serverName)
        if suite != nil && suite.Enabled() {
            results := suite.Run()
            for _, r := range results {
                o.addResult(r)
            }
        }
    }

    // 3. Health Tests (opt-in)
    if o.config.TestLevels.Health {
        suite := o.registry.GetHealthTests(serverName)
        if suite != nil && suite.Enabled() {
            results := suite.Run()
            for _, r := range results {
                o.addResult(r)
            }
        }
    }
}

func (o *Orchestrator) addResult(result TestResult) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.results = append(o.results, result)
}

func (o *Orchestrator) generateReport(startTime time.Time) *SmokeTestReport {
    o.mu.Lock()
    defer o.mu.Unlock()

    report := &SmokeTestReport{
        Results:       o.results,
        TotalDuration: time.Since(startTime),
        Timestamp:     time.Now(),
    }

    // Aggregate statistics
    serversSeen := make(map[string]bool)
    for _, r := range o.results {
        serversSeen[r.ServerName] = true

        switch r.Status {
        case StatusPass:
            report.PassedTests++
        case StatusFail:
            report.FailedTests++
            if r.TestType == "config" {
                report.CriticalFailures = append(report.CriticalFailures, r)
            }
        case StatusSkip:
            report.SkippedTests++
        }
    }

    report.TestedServers = len(serversSeen)
    report.TotalServers = len(serversSeen)

    return report
}
```

```go
// Jarvis/smoketests/config_test.go
package smoketests

import (
    "fmt"
    "os"
    "regexp"
    "time"
)

type ConfigTestSuite struct {
    serverName string
    checks     []ConfigCheck
    enabled    bool
}

type ConfigCheck struct {
    Type           string  // "env_var", "file", "permission"
    Name           string
    Required       bool
    Pattern        string  // Regex pattern for validation
    ErrorMessage   string
    FixSuggestion  string
}

func (s *ConfigTestSuite) Run() []TestResult {
    var results []TestResult

    for _, check := range s.checks {
        startTime := time.Now()
        result := TestResult{
            ServerName: s.serverName,
            TestType:   "config",
            TestName:   check.Name,
            Timestamp:  startTime,
        }

        switch check.Type {
        case "env_var":
            value := os.Getenv(check.Name)
            if value == "" {
                if check.Required {
                    result.Status = StatusFail
                    result.ErrorMessage = check.ErrorMessage
                    result.FixSuggestion = check.FixSuggestion
                } else {
                    result.Status = StatusSkip
                }
            } else {
                // Validate pattern if provided
                if check.Pattern != "" {
                    matched, _ := regexp.MatchString(check.Pattern, value)
                    if !matched {
                        result.Status = StatusFail
                        result.ErrorMessage = fmt.Sprintf("%s (invalid format)", check.ErrorMessage)
                        result.FixSuggestion = check.FixSuggestion
                    } else {
                        result.Status = StatusPass
                    }
                } else {
                    result.Status = StatusPass
                }
            }

        case "file":
            if _, err := os.Stat(check.Name); os.IsNotExist(err) {
                if check.Required {
                    result.Status = StatusFail
                    result.ErrorMessage = fmt.Sprintf("Required file not found: %s", check.Name)
                    result.FixSuggestion = check.FixSuggestion
                } else {
                    result.Status = StatusSkip
                }
            } else {
                result.Status = StatusPass
            }
        }

        result.Duration = time.Since(startTime)
        results = append(results, result)
    }

    return results
}

func (s *ConfigTestSuite) Name() string {
    return fmt.Sprintf("%s_config", s.serverName)
}

func (s *ConfigTestSuite) Enabled() bool {
    return s.enabled
}
```

```go
// Jarvis/smoketests/connectivity_test.go
package smoketests

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

type ConnectivityTestSuite struct {
    serverName string
    checks     []ConnectivityCheck
    enabled    bool
    timeout    time.Duration
}

type ConnectivityCheck struct {
    Type           string // "http", "tcp", "docker"
    Endpoint       string
    Method         string // For HTTP
    ExpectedStatus []int  // For HTTP
    ErrorMessage   string
    FixSuggestion  string
}

func (s *ConnectivityTestSuite) Run() []TestResult {
    var results []TestResult

    for _, check := range s.checks {
        startTime := time.Now()
        result := TestResult{
            ServerName: s.serverName,
            TestType:   "connectivity",
            TestName:   fmt.Sprintf("%s_%s", check.Type, check.Endpoint),
            Timestamp:  startTime,
        }

        ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
        defer cancel()

        switch check.Type {
        case "http":
            err := s.testHTTP(ctx, check)
            if err != nil {
                result.Status = StatusFail
                result.ErrorMessage = check.ErrorMessage
                result.FixSuggestion = check.FixSuggestion
                result.Details = err.Error()
            } else {
                result.Status = StatusPass
            }

        case "tcp":
            // TCP connection test
            result.Status = StatusSkip
            result.Details = "TCP tests not yet implemented"

        case "docker":
            // Docker container status test
            result.Status = StatusSkip
            result.Details = "Docker tests not yet implemented"
        }

        result.Duration = time.Since(startTime)
        results = append(results, result)
    }

    return results
}

func (s *ConnectivityTestSuite) testHTTP(ctx context.Context, check ConnectivityCheck) error {
    req, err := http.NewRequestWithContext(ctx, check.Method, check.Endpoint, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    client := &http.Client{
        Timeout: s.timeout,
    }

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Check if status code is in expected range
    for _, expected := range check.ExpectedStatus {
        if resp.StatusCode == expected {
            return nil
        }
    }

    return fmt.Errorf("unexpected status code: %d (expected one of %v)", resp.StatusCode, check.ExpectedStatus)
}

func (s *ConnectivityTestSuite) Name() string {
    return fmt.Sprintf("%s_connectivity", s.serverName)
}

func (s *ConnectivityTestSuite) Enabled() bool {
    return s.enabled
}
```

---

### Phase 2: Integration with Jarvis (Week 1-2)

**Files to Modify:**
- `Jarvis/main.go` - Add boot-time smoke tests
- `Jarvis/tools.go` - Add `run_smoke_tests()` tool

**Boot Integration:**

```go
// Jarvis/main.go additions

func main() {
    setupLogging()
    printBanner()

    // NEW: Run smoke tests on boot if enabled
    if shouldRunSmokeTests() {
        runBootSmokeTests()
    }

    // ... existing server setup
}

func shouldRunSmokeTests() bool {
    // Check config or env var
    if os.Getenv("JARVIS_SKIP_SMOKE_TESTS") == "true" {
        return false
    }

    // Check MCPM config
    // config := loadMcpmConfig()
    // return config.SmokeTests.RunOnBoot

    return true // Default: run on boot
}

func runBootSmokeTests() {
    log.Println("🔍 Running smoke tests...")

    orchestrator := smoketests.NewOrchestrator(loadSmokeTestConfig())

    // Get active servers from current profiles
    servers := getActiveServers()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    report := orchestrator.RunAll(ctx, servers)

    // Log results
    logSmokeTestReport(report)

    // Print warning to stderr if failures
    if len(report.CriticalFailures) > 0 {
        fmt.Fprintln(os.Stderr, "\n⚠️  SMOKE TEST FAILURES DETECTED:")
        for _, failure := range report.CriticalFailures {
            fmt.Fprintf(os.Stderr, "  ❌ %s: %s\n", failure.ServerName, failure.ErrorMessage)
            if failure.FixSuggestion != "" {
                fmt.Fprintf(os.Stderr, "     💡 %s\n", failure.FixSuggestion)
            }
        }
        fmt.Fprintln(os.Stderr, "\n  Use `check_status()` for full diagnostics\n")
    } else if report.FailedTests == 0 {
        fmt.Fprintln(os.Stderr, "✅ All smoke tests passed\n")
    }
}

func getActiveServers() []string {
    // Parse which profiles are active based on cwd, client, etc.
    // For now, simple implementation:

    servers := []string{}

    // Check toolbox profile (example)
    profileServers := []string{
        "firecrawl",
        "context7",
        "kagimcp",
        "brave-search",
        "fetch-mcp",
        "time",
    }

    servers = append(servers, profileServers...)

    // Add memory profile servers
    servers = append(servers, "basic-memory", "mem0-mcp")

    return servers
}

func logSmokeTestReport(report *smoketests.SmokeTestReport) {
    log.Printf("Smoke Test Report:")
    log.Printf("  Total Servers: %d", report.TotalServers)
    log.Printf("  Passed: %d, Failed: %d, Skipped: %d",
        report.PassedTests, report.FailedTests, report.SkippedTests)
    log.Printf("  Duration: %v", report.TotalDuration)

    for _, result := range report.Results {
        if result.Status == StatusFail {
            log.Printf("  FAIL: %s/%s - %s", result.ServerName, result.TestName, result.ErrorMessage)
        }
    }
}
```

**Tool Integration:**

```go
// Jarvis/main.go - Add tool
s.AddTool(mcp.NewTool("run_smoke_tests",
    mcp.WithDescription(`Execute smoke tests on configured MCP servers to validate configuration, connectivity, and health. Returns detailed report with pass/fail status and fix suggestions.

Use this tool when:
  • You suspect configuration issues
  • After adding/updating servers
  • Before starting critical workflows
  • For regular health monitoring

Test Levels:
  • config: Validate environment variables and files (fast, always safe)
  • connectivity: Test network reachability (fast, no API costs)
  • health: Execute actual tool calls (slow, may consume API credits)

Examples:
  • Run all tests: run_smoke_tests()
  • Config only: run_smoke_tests(test_level="config")
  • Specific server: run_smoke_tests(servers=["firecrawl"])
  • Skip tests: run_smoke_tests(exclude=["slow-service"])`),
    mcp.WithString("test_level",
        mcp.Description("Test thoroughness: 'config', 'connectivity', or 'health' (default: 'connectivity')"),
    ),
    mcp.WithString("servers",
        mcp.Description("Comma-separated list of servers to test (default: all active servers)"),
    ),
    mcp.WithString("exclude",
        mcp.Description("Comma-separated list of servers to exclude"),
    ),
    mcp.WithInteger("timeout",
        mcp.Description("Timeout in seconds (default: 10)"),
        mcp.DefaultInt(10),
    ),
    mcp.WithOutputSchema(getOutputSchema(smoketests.SmokeTestReport{})),
), handleRunSmokeTests)

// Jarvis/tools.go - Handler
func handleRunSmokeTests(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args, _ := request.Params.Arguments.(map[string]interface{})

    testLevel, _ := args["test_level"].(string)
    if testLevel == "" {
        testLevel = "connectivity"
    }

    timeout, _ := args["timeout"].(int)
    if timeout == 0 {
        timeout = 10
    }

    // Parse servers and exclusions
    var servers []string
    if serversStr, ok := args["servers"].(string); ok && serversStr != "" {
        servers = strings.Split(serversStr, ",")
    } else {
        servers = getActiveServers()
    }

    if excludeStr, ok := args["exclude"].(string); ok && excludeStr != "" {
        exclude := strings.Split(excludeStr, ",")
        servers = filterServers(servers, exclude)
    }

    // Create config
    config := &smoketests.Config{
        Enabled: true,
        TestLevels: smoketests.TestLevels{
            Config:       testLevel == "config" || testLevel == "connectivity" || testLevel == "health",
            Connectivity: testLevel == "connectivity" || testLevel == "health",
            Health:       testLevel == "health",
        },
        GlobalTimeout:     time.Duration(timeout) * time.Second,
        MaxParallelTests:  5,
        FailFast:          false,
    }

    orchestrator := smoketests.NewOrchestrator(config)

    ctx, cancel := context.WithTimeout(context.Background(), config.GlobalTimeout)
    defer cancel()

    report := orchestrator.RunAll(ctx, servers)

    // Format output
    output := formatSmokeTestReport(report)

    // Return both human-readable and structured
    jsonData, _ := json.MarshalIndent(report, "", "  ")

    return &mcp.CallToolResult{
        Content: []interface{}{
            mcp.NewTextContent(output),
        },
        Meta: map[string]interface{}{
            "structured": report,
        },
    }, nil
}

func formatSmokeTestReport(report *smoketests.SmokeTestReport) string {
    var output strings.Builder

    output.WriteString(fmt.Sprintf("# Smoke Test Report\n\n"))
    output.WriteString(fmt.Sprintf("**Duration:** %v\n", report.TotalDuration))
    output.WriteString(fmt.Sprintf("**Servers Tested:** %d\n\n", report.TestedServers))

    // Summary
    output.WriteString(fmt.Sprintf("## Summary\n\n"))
    output.WriteString(fmt.Sprintf("- ✅ Passed: %d\n", report.PassedTests))
    output.WriteString(fmt.Sprintf("- ❌ Failed: %d\n", report.FailedTests))
    output.WriteString(fmt.Sprintf("- ⏭️  Skipped: %d\n\n", report.SkippedTests))

    // Critical Failures
    if len(report.CriticalFailures) > 0 {
        output.WriteString("## ❌ Critical Failures\n\n")
        for _, failure := range report.CriticalFailures {
            output.WriteString(fmt.Sprintf("### %s\n", failure.ServerName))
            output.WriteString(fmt.Sprintf("**Issue:** %s\n\n", failure.ErrorMessage))
            if failure.FixSuggestion != "" {
                output.WriteString(fmt.Sprintf("**Fix:** %s\n\n", failure.FixSuggestion))
            }
        }
    }

    // All Results (by server)
    serverGroups := groupResultsByServer(report.Results)
    output.WriteString("## Detailed Results\n\n")

    for server, results := range serverGroups {
        passed := 0
        failed := 0
        for _, r := range results {
            if r.Status == StatusPass {
                passed++
            } else if r.Status == StatusFail {
                failed++
            }
        }

        status := "✅"
        if failed > 0 {
            status = "❌"
        }

        output.WriteString(fmt.Sprintf("%s **%s** - %d passed, %d failed\n", status, server, passed, failed))

        // Show failures
        for _, r := range results {
            if r.Status == StatusFail {
                output.WriteString(fmt.Sprintf("  - ❌ %s: %s\n", r.TestName, r.ErrorMessage))
            }
        }
        output.WriteString("\n")
    }

    return output.String()
}
```

---

## Phase 3: Transport Protocol Migration (SSE → Streamable HTTP)

**Create test specs for existing servers:**

**Firecrawl:**
```json
{
  "server": "firecrawl",
  "tests": {
    "config": {
      "enabled": true,
      "checks": [
        {
          "type": "env_var",
          "name": "FIRECRAWL_API_KEY",
          "required": true,
          "pattern": "^fc-[a-f0-9]{32}$",
          "error_message": "Missing or invalid FIRECRAWL_API_KEY",
          "fix_suggestion": "Get API key from https://firecrawl.dev → Settings → API Keys, then run: jarvis edit_server('firecrawl', env='FIRECRAWL_API_KEY=fc-YOUR_KEY')"
        }
      ]
    },
    "connectivity": {
      "enabled": true,
      "timeout": "3s",
      "checks": [
        {
          "type": "http",
          "endpoint": "https://api.firecrawl.dev/v1",
          "method": "HEAD",
          "expected_status": [200, 401, 403],
          "error_message": "Firecrawl API unreachable",
          "fix_suggestion": "Check internet connection or visit https://status.firecrawl.dev"
        }
      ]
    }
  }
}
```

**Brave Search:**
```json
{
  "server": "brave-search",
  "tests": {
    "config": {
      "enabled": true,
      "checks": [
        {
          "type": "env_var",
          "name": "BRAVE_API_KEY",
          "required": true,
          "error_message": "Missing BRAVE_API_KEY",
          "fix_suggestion": "Get API key from https://brave.com/search/api/, then configure with: jarvis edit_server('brave-search', env='BRAVE_API_KEY=BSA...')"
        }
      ]
    }
  }
}
```

**Context7:**
```json
{
  "server": "context7",
  "tests": {
    "connectivity": {
      "enabled": true,
      "timeout": "2s",
      "checks": [
        {
          "type": "http",
          "endpoint": "https://api.context7.ai",
          "method": "HEAD",
          "expected_status": [200, 401],
          "error_message": "Context7 API unreachable",
          "fix_suggestion": "Check internet connection"
        }
      ]
    }
  }
}
```

**Qdrant (Docker):**
```json
{
  "server": "mcp-server-qdrant",
  "tests": {
    "connectivity": {
      "enabled": true,
      "timeout": "2s",
      "checks": [
        {
          "type": "http",
          "endpoint": "http://localhost:6333/health",
          "method": "GET",
          "expected_status": [200],
          "error_message": "Qdrant is not running",
          "fix_suggestion": "Start Qdrant with: docker compose up -d mcp-qdrant"
        }
      ]
    }
  }
}
```

**PostgreSQL (Docker):**
```json
{
  "server": "postgresql",
  "tests": {
    "connectivity": {
      "enabled": true,
      "timeout": "2s",
      "checks": [
        {
          "type": "docker",
          "container": "mcp-postgres",
          "expected_status": "running",
          "error_message": "PostgreSQL container not running",
          "fix_suggestion": "Start PostgreSQL with: docker compose up -d mcp-postgres"
        }
      ]
    }
  }
}
```

---

## User Experience

### Boot Sequence with Smoke Tests

```
$ claude-code

     ██╗ █████╗ ██████╗ ██╗   ██╗██╗███████╗
     ██║██╔══██╗██╔══██╗██║   ██║██║██╔════╝
     ██║███████║██████╔╝██║   ██║██║███████╗
██   ██║██╔══██║██╔══██╗╚██╗ ██╔╝██║╚════██║
╚█████╔╝██║  ██║██║  ██║ ╚████╔╝ ██║███████║
 ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝

>> JARVIS MCP Gateway v1.0.0 initialized <<

🔍 Running smoke tests...
  ✅ firecrawl (config, connectivity) - 2/2 passed
  ✅ context7 (connectivity) - 1/1 passed
  ✅ kagimcp (config, connectivity) - 2/2 passed
  ❌ brave-search (config) - FAILED
  ✅ fetch-mcp (connectivity) - 1/1 passed
  ✅ basic-memory (config, connectivity) - 2/2 passed

⚠️  SMOKE TEST FAILURES DETECTED:
  ❌ brave-search: Missing BRAVE_API_KEY
     💡 Get API key from https://brave.com/search/api/,
        then configure: jarvis edit_server('brave-search', env='BRAVE_API_KEY=BSA...')

  Use `check_status()` or `run_smoke_tests()` for full diagnostics

✅ Jarvis ready (7ms smoke tests)
```

### On-Demand Testing

```
user> Run smoke tests to check everything is configured correctly

assistant> I'll run comprehensive smoke tests on all active servers.

[Uses run_smoke_tests tool]

assistant> Smoke test results:

✅ **All systems operational!**

Tested 8 servers in 1.2 seconds:
- ✅ firecrawl - Configuration and connectivity verified
- ✅ context7 - API reachable
- ✅ kagimcp - API key valid, service responding
- ✅ fetch-mcp - Ready
- ✅ time - Ready
- ✅ basic-memory - Database connected
- ✅ mem0-mcp - Vector store accessible
- ✅ morph-fast-apply - Service healthy

All 15 tests passed. Your Jarvis environment is fully operational.
```

---

## Configuration Management

### Adding Tests for New Servers

**Option 1: Automatic (Recommended)**
When installing a server, check if upstream provides smoke test spec:

```bash
mcpm install new-server
# Downloads test spec from https://github.com/author/new-server/smoke_tests.json
# Saves to ~/.mcpm/servers/new-server/smoke_tests.json
```

**Option 2: Manual**
Create `~/.mcpm/servers/new-server/smoke_tests.json`:

```json
{
  "server": "new-server",
  "tests": {
    "config": {
      "enabled": true,
      "checks": [/* your checks */]
    }
  }
}
```

**Option 3: Programmatic (Advanced)**
```go
// In server metadata
smoketests.RegisterTests("new-server", &smoketests.ConfigTestSuite{
    serverName: "new-server",
    checks: []smoketests.ConfigCheck{
        {Type: "env_var", Name: "NEW_SERVER_KEY", Required: true},
    },
    enabled: true,
})
```

---

## Performance Considerations

### Boot Time Impact

**Target:** <2 seconds additional boot time
**Actual (estimated):**
- Config tests: ~100ms for 10 servers (parallel)
- Connectivity tests: ~500ms for 10 servers (parallel, with timeouts)
- **Total:** ~600ms

**Mitigation Strategies:**
1. **Parallel execution** - Test multiple servers concurrently
2. **Aggressive timeouts** - 2s max per connectivity test
3. **Skip slow tests** - Health tests disabled by default
4. **Caching** - Cache successful results for 5 minutes
5. **Progressive loading** - Start MCP server before tests finish

### Memory Impact

**Per test result:** ~200 bytes
**100 tests:** ~20KB
**Negligible impact** on Jarvis memory footprint

---

## Integration with Existing Tools

### Enhanced `check_status()`

Merge smoke test results into existing diagnostics:

```go
func handleCheckStatus(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... existing MCPM doctor logic

    // NEW: Append smoke test results
    smokeReport := runQuickSmokeTests()

    output += formatSmokeTestSummary(smokeReport)

    return mcp.NewToolResultText(output), nil
}
```

Output:
```markdown
# System Health Check

## MCPM Status
✅ MCPM v2.1.0 installed
✅ Configuration valid
...

## Smoke Test Results (Last Run: 2 minutes ago)
✅ 8/8 servers configured correctly
✅ 7/8 servers reachable
⚠️  1 warning: brave-search missing API key

Run `run_smoke_tests()` for detailed diagnostics.
```

---

## Testing the Smoke Testing System

### Unit Tests

```go
// Jarvis/smoketests/orchestrator_test.go

func TestConfigTests(t *testing.T) {
    // Set up test environment
    os.Setenv("TEST_API_KEY", "test-key")
    defer os.Unsetenv("TEST_API_KEY")

    suite := &ConfigTestSuite{
        serverName: "test-server",
        checks: []ConfigCheck{
            {
                Type:     "env_var",
                Name:     "TEST_API_KEY",
                Required: true,
            },
        },
        enabled: true,
    }

    results := suite.Run()

    if len(results) != 1 {
        t.Fatalf("Expected 1 result, got %d", len(results))
    }

    if results[0].Status != StatusPass {
        t.Errorf("Expected PASS, got %s", results[0].Status)
    }
}

func TestConnectivityTestsTimeout(t *testing.T) {
    suite := &ConnectivityTestSuite{
        serverName: "test-server",
        checks: []ConnectivityCheck{
            {
                Type:     "http",
                Endpoint: "http://localhost:99999", // Unreachable
                Method:   "HEAD",
            },
        },
        enabled: true,
        timeout: 100 * time.Millisecond,
    }

    results := suite.Run()

    if results[0].Status != StatusFail {
        t.Errorf("Expected timeout failure, got %s", results[0].Status)
    }

    if results[0].Duration > 200*time.Millisecond {
        t.Errorf("Test took too long: %v", results[0].Duration)
    }
}
```

### Integration Tests

```bash
# test/integration/smoke_tests.sh

#!/bin/bash
set -e

echo "Testing smoke test system..."

# 1. Test with missing API key
unset FIRECRAWL_API_KEY
./jarvis run_smoke_tests servers=firecrawl | grep "Missing.*FIRECRAWL_API_KEY"

# 2. Test with valid config
export FIRECRAWL_API_KEY="fc-test"
./jarvis run_smoke_tests servers=firecrawl test_level=config | grep "PASS"

# 3. Test timeout handling
./jarvis run_smoke_tests timeout=1 | grep "timeout"

echo "✅ All integration tests passed"
```

---

## Rollout Plan

### Week 1: Foundation
- [ ] Implement core types and orchestrator
- [ ] Build config test suite
- [ ] Build connectivity test suite
- [ ] Unit test coverage >80%

### Week 2: Integration
- [ ] Integrate with Jarvis main.go
- [ ] Add `run_smoke_tests()` tool
- [ ] Create test specs for 5 core servers
- [ ] Boot-time smoke tests working

### Week 3: Expansion
- [ ] Add test specs for remaining servers
- [ ] Implement health tests (opt-in)
- [ ] Add caching for test results
- [ ] Performance optimization

### Week 4: Polish
- [ ] Documentation (user guide, developer guide)
- [ ] E2E testing with real agents
- [ ] CI/CD integration
- [ ] Release

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Boot time impact | <2s | Time from start to "ready" message |
| Config detection rate | >95% | % of missing configs detected |
| False positive rate | <5% | % of failed tests that aren't real issues |
| Agent satisfaction | >90% | Survey: "Smoke tests are helpful" |
| MTTR (Mean Time To Recovery) | -50% | Time to diagnose and fix config issues |

---

## Future Enhancements

### v2.0 Features
1. **Continuous Monitoring** - Run smoke tests every N minutes in background
2. **Smart Alerts** - Proactive warnings when services degrade
3. **Auto-Remediation** - Attempt automatic fixes (e.g., restart containers)
4. **Test History** - Track test results over time, show trends
5. **Custom Test Plugins** - User-defined test logic
6. **MCP Protocol Tests** - Validate MCP spec compliance

### v3.0 Features
1. **Distributed Testing** - Test servers across multiple regions
2. **Load Testing** - Verify servers can handle concurrent requests
3. **Security Testing** - Check for common vulnerabilities
4. **Performance Profiling** - Detailed latency breakdowns
5. **AI-Powered Diagnostics** - LLM analyzes failures and suggests fixes

---

## Appendix

### Sample Test Definitions

See `config/smoke_tests/` directory for complete examples:
- `firecrawl.json`
- `brave-search.json`
- `context7.json`
- `kagimcp.json`
- `qdrant.json`
- `postgresql.json`
- `basic-memory.json`
- `mem0-mcp.json`

### CLI Commands

```bash
# Run all smoke tests
jarvis run_smoke_tests

# Run config tests only (fast)
jarvis run_smoke_tests test_level=config

# Test specific server
jarvis run_smoke_tests servers=firecrawl

# Skip slow servers
jarvis run_smoke_tests exclude=slow-service,another-service

# Run with custom timeout
jarvis run_smoke_tests timeout=30

# Health tests (uses API credits)
jarvis run_smoke_tests test_level=health

# Disable boot-time tests
export JARVIS_SKIP_SMOKE_TESTS=true

# View last test results
jarvis check_status
```

---

**Document Status:** Ready for Implementation
**Estimated Effort:** 2-3 weeks
**Priority:** High (addresses critical operational need)
**Dependencies:** None (can start immediately)
