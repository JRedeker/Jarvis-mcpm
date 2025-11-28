package smoketests

import "time"

// TestStatus represents the outcome of a smoke test
type TestStatus string

const (
	StatusPass    TestStatus = "pass"
	StatusFail    TestStatus = "fail"
	StatusSkip    TestStatus = "skip"
	StatusTimeout TestStatus = "timeout"
)

// TestResult captures the outcome of a single test
type TestResult struct {
	ServerName    string        `json:"server_name"`
	TestType      string        `json:"test_type"` // "config", "connectivity", "health"
	TestName      string        `json:"test_name"`
	Status        TestStatus    `json:"status"`
	ErrorMessage  string        `json:"error_message,omitempty"`
	FixSuggestion string        `json:"fix_suggestion,omitempty"`
	Duration      time.Duration `json:"duration"`
	Timestamp     time.Time     `json:"timestamp"`
	Details       string        `json:"details,omitempty"`
}

// TestSuite is the interface that all test suites must implement
type TestSuite interface {
	Run() []TestResult
	Name() string
	Enabled() bool
}

// SmokeTestReport aggregates all test results
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

// Config defines smoke test configuration
type Config struct {
	Enabled           bool
	RunOnBoot         bool
	TestLevels        TestLevels
	ParallelExecution bool
	MaxParallelTests  int
	GlobalTimeout     time.Duration
	FailFast          bool
	ExcludeServers    []string
}

// TestLevels controls which test categories to run
type TestLevels struct {
	Config       bool
	Connectivity bool
	Health       bool
}

// DefaultConfig returns sensible defaults for smoke testing
func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		RunOnBoot:         true,
		TestLevels:        TestLevels{Config: true, Connectivity: true, Health: false},
		ParallelExecution: true,
		MaxParallelTests:  5,
		GlobalTimeout:     10 * time.Second,
		FailFast:          false,
		ExcludeServers:    []string{},
	}
}
