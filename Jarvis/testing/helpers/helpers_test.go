package helpers

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewToolRequest(t *testing.T) {
	args := map[string]interface{}{
		"name":  "test-server",
		"force": true,
	}

	req := NewToolRequest(args)

	if req.Params.Arguments == nil {
		t.Fatal("Arguments should not be nil")
	}

	argsMap, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		t.Fatal("Arguments should be map[string]interface{}")
	}

	if argsMap["name"] != "test-server" {
		t.Errorf("name = %v, want test-server", argsMap["name"])
	}

	if argsMap["force"] != true {
		t.Errorf("force = %v, want true", argsMap["force"])
	}
}

func TestEmptyRequest(t *testing.T) {
	req := EmptyRequest()

	if req.Params.Arguments != nil {
		t.Error("EmptyRequest should have nil Arguments")
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "foo", false},
		{"hello", "", true},
		{"abc", "abc", true},
		{"abc", "abcd", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := containsString(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("containsString(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestGetResultText(t *testing.T) {
	t.Run("returns empty for nil result", func(t *testing.T) {
		text := GetResultText(nil)
		if text != "" {
			t.Errorf("GetResultText(nil) = %q, want empty", text)
		}
	})

	t.Run("extracts text from TextContent", func(t *testing.T) {
		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Hello, World!",
				},
			},
		}

		text := GetResultText(result)
		if text != "Hello, World!" {
			t.Errorf("GetResultText() = %q, want %q", text, "Hello, World!")
		}
	})
}

func TestAssertFunctions(t *testing.T) {
	// These are helper assertions, so we test them indirectly
	// by ensuring they don't panic on valid input

	t.Run("AssertNoError does not panic on nil", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNoError(mockT, nil)
		// If we got here, it didn't panic
	})

	t.Run("AssertResultNotNil does not panic on non-nil", func(t *testing.T) {
		mockT := &testing.T{}
		result := &mcp.CallToolResult{}
		AssertResultNotNil(mockT, result)
	})
}
