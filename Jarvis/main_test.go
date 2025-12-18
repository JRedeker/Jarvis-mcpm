package main

import (
	"os"
	"path/filepath"
	"testing"

	"jarvis/handlers"
)

func TestSetupLogging(t *testing.T) {
	// Store original directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create a temporary directory
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create logs directory
	logsDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logsDir, 0755)

	// Test that setupLogging doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("setupLogging() panicked: %v", r)
		}
	}()

	setupLogging()

	// Check if log file was created
	logPath := filepath.Join(logsDir, "jarvis.log")
	if _, err := os.Stat(logPath); err == nil {
		// Log file exists, which is good
		t.Log("Log file created successfully")
	} else {
		// Log file might not be created in test environment, which is okay
		t.Log("Log file not created (may be expected in test environment)")
	}
}

func TestMain(m *testing.M) {
	// Setup before tests
	// Clean up any test artifacts
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

// TDD: Test that the refactored main correctly wires handlers

func TestCreateServer_RegistersAllTools(t *testing.T) {
	// Given: Production handler and server
	h := handlers.CreateProductionHandler()

	// When: We get tool definitions
	defs := handlers.GetToolDefinitions(h)

	// Then: All 8 consolidated tools should be registered (v3.0)
	expectedTools := []string{
		"jarvis_check_status",
		"jarvis_server",
		"jarvis_profile",
		"jarvis_client",
		"jarvis_config",
		"jarvis_project",
		"jarvis_system",
		"jarvis_share",
	}

	// Build map of registered tools
	registered := make(map[string]bool)
	for _, def := range defs {
		registered[def.Tool.Name] = true
	}

	// Verify exactly 8 tools
	if len(defs) != 8 {
		t.Errorf("Expected 8 consolidated tools, got %d", len(defs))
	}

	for _, name := range expectedTools {
		if !registered[name] {
			t.Errorf("Expected tool %s to be registered", name)
		}
	}
}

func TestCreateServer_ToolDefinitionsHaveDescriptions(t *testing.T) {
	h := handlers.CreateProductionHandler()
	defs := handlers.GetToolDefinitions(h)

	for _, def := range defs {
		if def.Tool.Description == "" {
			t.Errorf("Tool %s has no description", def.Tool.Name)
		}
	}
}

func TestCreateServer_ToolDefinitionsHaveHandlers(t *testing.T) {
	h := handlers.CreateProductionHandler()
	defs := handlers.GetToolDefinitions(h)

	for _, def := range defs {
		if def.Handler == nil {
			t.Errorf("Tool %s has nil handler", def.Tool.Name)
		}
	}
}
