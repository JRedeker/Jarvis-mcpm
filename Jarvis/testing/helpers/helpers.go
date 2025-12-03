// Package helpers provides test utilities for Jarvis handler tests
package helpers

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestContext returns a context with a reasonable timeout for tests
func TestContext() context.Context {
	return context.Background()
}

// NewToolRequest creates a CallToolRequest with the given arguments
func NewToolRequest(args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

// EmptyRequest returns an empty CallToolRequest
func EmptyRequest() mcp.CallToolRequest {
	return mcp.CallToolRequest{}
}

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

// AssertResultNotNil fails if result is nil
func AssertResultNotNil(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result == nil {
		t.Fatal("result is nil")
	}
}

// AssertResultContains checks if the result text contains the expected string
func AssertResultContains(t *testing.T, result *mcp.CallToolResult, expected string) {
	t.Helper()
	if result == nil {
		t.Fatal("result is nil")
		return
	}

	text := GetResultText(result)
	if text == "" {
		t.Fatalf("result has no text content")
	}

	if !containsString(text, expected) {
		t.Errorf("result text does not contain %q\nGot: %s", expected, text)
	}
}

// AssertResultNotContains checks if the result text does NOT contain the string
func AssertResultNotContains(t *testing.T, result *mcp.CallToolResult, notExpected string) {
	t.Helper()
	if result == nil {
		return
	}

	text := GetResultText(result)
	if containsString(text, notExpected) {
		t.Errorf("result text should not contain %q\nGot: %s", notExpected, text)
	}
}

// AssertResultIsError checks if the result indicates an error
func AssertResultIsError(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result == nil {
		t.Fatal("result is nil")
		return
	}

	if !result.IsError {
		t.Error("expected result to be an error")
	}
}

// AssertResultIsSuccess checks if the result indicates success
func AssertResultIsSuccess(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result == nil {
		t.Fatal("result is nil")
		return
	}

	if result.IsError {
		t.Errorf("expected success but got error: %s", GetResultText(result))
	}
}

// GetResultText extracts the text content from a CallToolResult
func GetResultText(result *mcp.CallToolResult) string {
	if result == nil {
		return ""
	}

	// Check Content array (standard MCP response)
	if len(result.Content) > 0 {
		for _, item := range result.Content {
			// Use type switch for mcp.Content interface
			switch c := item.(type) {
			case mcp.TextContent:
				return c.Text
			case *mcp.TextContent:
				return c.Text
			}
		}
	}

	return ""
}

// containsString is a simple substring check
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// HandlerTestCase represents a test case for a handler
type HandlerTestCase struct {
	Name              string
	Args              map[string]interface{}
	ExpectError       bool
	ExpectContains    []string
	ExpectNotContains []string
	Setup             func()
	Validate          func(t *testing.T, result *mcp.CallToolResult, err error)
}

// RunHandlerTests runs a set of handler test cases
func RunHandlerTests(
	t *testing.T,
	handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error),
	tests []HandlerTestCase,
) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Run setup if provided
			if tt.Setup != nil {
				tt.Setup()
			}

			// Create request
			req := NewToolRequest(tt.Args)

			// Call handler
			result, err := handler(TestContext(), req)

			// Check error expectation
			if tt.ExpectError && err == nil && (result == nil || !result.IsError) {
				t.Error("expected error but got success")
			}
			if !tt.ExpectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check contains
			for _, expected := range tt.ExpectContains {
				AssertResultContains(t, result, expected)
			}

			// Check not contains
			for _, notExpected := range tt.ExpectNotContains {
				AssertResultNotContains(t, result, notExpected)
			}

			// Run custom validation
			if tt.Validate != nil {
				tt.Validate(t, result, err)
			}
		})
	}
}

// MockableHandler wraps a handler function to allow dependency injection
type MockableHandler struct {
	Handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// Call invokes the handler
func (h *MockableHandler) Call(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.Handler(ctx, req)
}
