package smoketests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigTestSuite_EnvVarCheck_Required(t *testing.T) {
	// Set a test env var
	os.Setenv("TEST_REQUIRED_VAR", "test-value")
	defer os.Unsetenv("TEST_REQUIRED_VAR")

	checks := []ConfigCheck{
		{
			Type:     "env_var",
			Name:     "TEST_REQUIRED_VAR",
			Required: true,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_EnvVarCheck_RequiredMissing(t *testing.T) {
	// Ensure the var is not set
	os.Unsetenv("TEST_MISSING_VAR")

	checks := []ConfigCheck{
		{
			Type:          "env_var",
			Name:          "TEST_MISSING_VAR",
			Required:      true,
			ErrorMessage:  "Missing required var",
			FixSuggestion: "Set TEST_MISSING_VAR",
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != StatusFail {
		t.Errorf("Expected StatusFail, got %v", results[0].Status)
	}
	if results[0].ErrorMessage != "Missing required var" {
		t.Errorf("Expected custom error message, got %q", results[0].ErrorMessage)
	}
	if results[0].FixSuggestion != "Set TEST_MISSING_VAR" {
		t.Errorf("Expected custom fix suggestion, got %q", results[0].FixSuggestion)
	}
}

func TestConfigTestSuite_EnvVarCheck_OptionalMissing(t *testing.T) {
	os.Unsetenv("TEST_OPTIONAL_VAR")

	checks := []ConfigCheck{
		{
			Type:     "env_var",
			Name:     "TEST_OPTIONAL_VAR",
			Required: false,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != StatusSkip {
		t.Errorf("Expected StatusSkip for optional missing var, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_EnvVarCheck_PatternMatch(t *testing.T) {
	os.Setenv("TEST_PATTERN_VAR", "abc-123")
	defer os.Unsetenv("TEST_PATTERN_VAR")

	checks := []ConfigCheck{
		{
			Type:     "env_var",
			Name:     "TEST_PATTERN_VAR",
			Required: true,
			Pattern:  `^[a-z]+-\d+$`,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass for matching pattern, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_EnvVarCheck_PatternNoMatch(t *testing.T) {
	os.Setenv("TEST_PATTERN_VAR", "invalid")
	defer os.Unsetenv("TEST_PATTERN_VAR")

	checks := []ConfigCheck{
		{
			Type:     "env_var",
			Name:     "TEST_PATTERN_VAR",
			Required: true,
			Pattern:  `^\d+$`, // Only digits
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusFail {
		t.Errorf("Expected StatusFail for non-matching pattern, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_FileCheck_Exists(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	checks := []ConfigCheck{
		{
			Type:     "file",
			Name:     tmpFile,
			Required: true,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass for existing file, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_FileCheck_Missing(t *testing.T) {
	checks := []ConfigCheck{
		{
			Type:     "file",
			Name:     "/nonexistent/path/to/file",
			Required: true,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusFail {
		t.Errorf("Expected StatusFail for missing file, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_FileCheck_OptionalMissing(t *testing.T) {
	checks := []ConfigCheck{
		{
			Type:     "file",
			Name:     "/nonexistent/optional/file",
			Required: false,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusSkip {
		t.Errorf("Expected StatusSkip for optional missing file, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_PermissionCheck_Readable(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "readable.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	checks := []ConfigCheck{
		{
			Type: "permission",
			Name: tmpFile,
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusPass {
		t.Errorf("Expected StatusPass for readable file, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_UnknownCheckType(t *testing.T) {
	checks := []ConfigCheck{
		{
			Type: "unknown_type",
			Name: "test",
		},
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if results[0].Status != StatusSkip {
		t.Errorf("Expected StatusSkip for unknown check type, got %v", results[0].Status)
	}
}

func TestConfigTestSuite_Name(t *testing.T) {
	suite := NewConfigTestSuite("my-server", nil, true)

	if suite.Name() != "my-server_config" {
		t.Errorf("Expected 'my-server_config', got %q", suite.Name())
	}
}

func TestConfigTestSuite_Enabled(t *testing.T) {
	enabledSuite := NewConfigTestSuite("test", nil, true)
	disabledSuite := NewConfigTestSuite("test", nil, false)

	if !enabledSuite.Enabled() {
		t.Error("Expected enabled suite to return true")
	}
	if disabledSuite.Enabled() {
		t.Error("Expected disabled suite to return false")
	}
}

func TestConfigTestSuite_MultipleChecks(t *testing.T) {
	os.Setenv("TEST_MULTI_1", "value1")
	os.Setenv("TEST_MULTI_2", "value2")
	defer os.Unsetenv("TEST_MULTI_1")
	defer os.Unsetenv("TEST_MULTI_2")

	checks := []ConfigCheck{
		{Type: "env_var", Name: "TEST_MULTI_1", Required: true},
		{Type: "env_var", Name: "TEST_MULTI_2", Required: true},
		{Type: "env_var", Name: "TEST_MULTI_3", Required: false}, // Optional, will skip
	}

	suite := NewConfigTestSuite("test-server", checks, true)
	results := suite.Run()

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// First two should pass, third should skip
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
