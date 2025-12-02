package smoketests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// ConnectivityTestSuite runs network connectivity tests
type ConnectivityTestSuite struct {
	serverName string
	checks     []ConnectivityCheck
	enabled    bool
	timeout    time.Duration
}

// ConnectivityCheck defines a network connectivity check
type ConnectivityCheck struct {
	Type           string // "http", "tcp", "docker"
	Endpoint       string // URL or host:port
	Method         string // For HTTP: GET, HEAD, POST, etc.
	ExpectedStatus []int  // For HTTP: acceptable status codes
	ErrorMessage   string
	FixSuggestion  string
}

// NewConnectivityTestSuite creates a new connectivity test suite
func NewConnectivityTestSuite(serverName string, checks []ConnectivityCheck, enabled bool, timeout time.Duration) *ConnectivityTestSuite {
	if timeout == 0 {
		timeout = 3 * time.Second
	}

	return &ConnectivityTestSuite{
		serverName: serverName,
		checks:     checks,
		enabled:    enabled,
		timeout:    timeout,
	}
}

// Run executes all connectivity checks
func (s *ConnectivityTestSuite) Run() []TestResult {
	var results []TestResult

	for _, check := range s.checks {
		startTime := time.Now()
		result := TestResult{
			ServerName: s.serverName,
			TestType:   "connectivity",
			TestName:   fmt.Sprintf("%s_%s", check.Type, s.sanitizeEndpoint(check.Endpoint)),
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
				if result.ErrorMessage == "" {
					result.ErrorMessage = fmt.Sprintf("HTTP check failed: %s", check.Endpoint)
				}
				result.FixSuggestion = check.FixSuggestion
				result.Details = err.Error()
			} else {
				result.Status = StatusPass
				result.Details = fmt.Sprintf("HTTP %s to %s succeeded", check.Method, check.Endpoint)
			}

		default:
			result.Status = StatusSkip
			result.Details = fmt.Sprintf("Unknown check type: %s", check.Type)
		}

		result.Duration = time.Since(startTime)

		if ctx.Err() == context.DeadlineExceeded {
			result.Status = StatusTimeout
			result.ErrorMessage = fmt.Sprintf("Test timed out after %v", s.timeout)
		}

		results = append(results, result)
	}

	return results
}

func (s *ConnectivityTestSuite) testHTTP(ctx context.Context, check ConnectivityCheck) error {
	method := check.Method
	if method == "" {
		method = "HEAD"
	}

	req, err := http.NewRequestWithContext(ctx, method, check.Endpoint, nil)
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

	if len(check.ExpectedStatus) > 0 {
		for _, expected := range check.ExpectedStatus {
			if resp.StatusCode == expected {
				return nil
			}
		}
		return fmt.Errorf("unexpected status code: %d (expected one of %v)", resp.StatusCode, check.ExpectedStatus)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil
	}

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (s *ConnectivityTestSuite) sanitizeEndpoint(endpoint string) string {
	sanitized := endpoint
	sanitized = regexp.MustCompile(`^https?://`).ReplaceAllString(sanitized, "")
	sanitized = regexp.MustCompile(`[^a-zA-Z0-9-_.]`).ReplaceAllString(sanitized, "_")
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}
	return sanitized
}

func (s *ConnectivityTestSuite) Name() string {
	return fmt.Sprintf("%s_connectivity", s.serverName)
}

func (s *ConnectivityTestSuite) Enabled() bool {
	return s.enabled
}
