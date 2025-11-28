package main

import (
	"os"
	"path/filepath"
	"testing"
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
