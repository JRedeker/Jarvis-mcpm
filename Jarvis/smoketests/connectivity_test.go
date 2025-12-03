package smoketests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestConnectivityTestSuite_HTTPCheck_Success(t *testing.T) {
	// Create a test server that returns 200
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{
			Type:     "http",
			Endpoint: ts.URL,
			Method:   "GET",
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass, got %v: %s", results[0].Status, results[0].Details)
	}
}

func TestConnectivityTestSuite_HTTPCheck_Failure(t *testing.T) {
	// Create a test server that returns 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{
			Type:          "http",
			Endpoint:      ts.URL,
			Method:        "GET",
			ErrorMessage:  "Server is down",
			FixSuggestion: "Check server logs",
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if results[0].Status != StatusFail {
		t.Errorf("Expected StatusFail for 500 response, got %v", results[0].Status)
	}
	if results[0].ErrorMessage != "Server is down" {
		t.Errorf("Expected custom error message, got %q", results[0].ErrorMessage)
	}
}

func TestConnectivityTestSuite_HTTPCheck_ExpectedStatus(t *testing.T) {
	// Server returns 404, but we expect it
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{
			Type:           "http",
			Endpoint:       ts.URL,
			ExpectedStatus: []int{404, 200}, // Accept both
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass when 404 is expected, got %v", results[0].Status)
	}
}

func TestConnectivityTestSuite_HTTPCheck_UnexpectedStatus(t *testing.T) {
	// Server returns 400, but we expect 200
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{
			Type:           "http",
			Endpoint:       ts.URL,
			ExpectedStatus: []int{200},
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if results[0].Status != StatusFail {
		t.Errorf("Expected StatusFail for unexpected status, got %v", results[0].Status)
	}
}

func TestConnectivityTestSuite_HTTPCheck_DefaultMethod(t *testing.T) {
	methodReceived := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodReceived = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{
			Type:     "http",
			Endpoint: ts.URL,
			// No method specified
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	suite.Run()

	if methodReceived != "HEAD" {
		t.Errorf("Expected default method HEAD, got %s", methodReceived)
	}
}

func TestConnectivityTestSuite_HTTPCheck_AuthErrors(t *testing.T) {
	// 401 and 403 should be treated as "pass" (server is responding)
	tests := []struct {
		name   string
		status int
	}{
		{"401 Unauthorized", http.StatusUnauthorized},
		{"403 Forbidden", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer ts.Close()

			checks := []ConnectivityCheck{
				{
					Type:     "http",
					Endpoint: ts.URL,
				},
			}

			suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
			results := suite.Run()

			if results[0].Status != StatusPass {
				t.Errorf("Expected StatusPass for %d, got %v", tt.status, results[0].Status)
			}
		})
	}
}

func TestConnectivityTestSuite_UnknownCheckType(t *testing.T) {
	checks := []ConnectivityCheck{
		{
			Type:     "unknown_type",
			Endpoint: "localhost:9999",
		},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if results[0].Status != StatusSkip {
		t.Errorf("Expected StatusSkip for unknown type, got %v", results[0].Status)
	}
}

func TestConnectivityTestSuite_DefaultTimeout(t *testing.T) {
	suite := NewConnectivityTestSuite("test-server", nil, true, 0)

	// Default timeout should be 3 seconds
	if suite.timeout != 3*time.Second {
		t.Errorf("Expected default timeout of 3s, got %v", suite.timeout)
	}
}

func TestConnectivityTestSuite_Name(t *testing.T) {
	suite := NewConnectivityTestSuite("my-server", nil, true, 0)

	if suite.Name() != "my-server_connectivity" {
		t.Errorf("Expected 'my-server_connectivity', got %q", suite.Name())
	}
}

func TestConnectivityTestSuite_Enabled(t *testing.T) {
	enabledSuite := NewConnectivityTestSuite("test", nil, true, 0)
	disabledSuite := NewConnectivityTestSuite("test", nil, false, 0)

	if !enabledSuite.Enabled() {
		t.Error("Expected enabled suite to return true")
	}
	if disabledSuite.Enabled() {
		t.Error("Expected disabled suite to return false")
	}
}

func TestConnectivityTestSuite_SanitizeEndpoint(t *testing.T) {
	suite := NewConnectivityTestSuite("test", nil, true, 0)

	tests := []struct {
		input    string
		expected string
	}{
		{"http://localhost:8080", "localhost_8080"},
		{"https://example.com/path", "example.com_path"},
		{"host:port", "host_port"},
	}

	for _, tt := range tests {
		result := suite.sanitizeEndpoint(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeEndpoint(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestConnectivityTestSuite_SanitizeEndpoint_Truncates(t *testing.T) {
	suite := NewConnectivityTestSuite("test", nil, true, 0)

	longEndpoint := "http://example.com/" + string(make([]byte, 100))
	result := suite.sanitizeEndpoint(longEndpoint)

	if len(result) > 50 {
		t.Errorf("Expected sanitized endpoint to be <= 50 chars, got %d", len(result))
	}
}

func TestConnectivityTestSuite_MultipleChecks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{Type: "http", Endpoint: ts.URL, Method: "GET"},
		{Type: "http", Endpoint: ts.URL, Method: "HEAD"},
		{Type: "unknown", Endpoint: "ignored"},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// First two should pass (http), third should skip (unknown)
	if results[0].Status != StatusPass {
		t.Errorf("Expected first check to pass")
	}
	if results[1].Status != StatusPass {
		t.Errorf("Expected second check to pass")
	}
	if results[2].Status != StatusSkip {
		t.Errorf("Expected third check to skip")
	}
}

func TestConnectivityTestSuite_RecordsDuration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	checks := []ConnectivityCheck{
		{Type: "http", Endpoint: ts.URL},
	}

	suite := NewConnectivityTestSuite("test-server", checks, true, 5*time.Second)
	results := suite.Run()

	if results[0].Duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", results[0].Duration)
	}
}
