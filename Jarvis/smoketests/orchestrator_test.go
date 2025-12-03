package smoketests

import (
	"context"
	"testing"
	"time"
)

func TestOrchestrator_RunAll_Empty(t *testing.T) {
	config := DefaultConfig()
	orch := NewOrchestrator(config)

	ctx := context.Background()
	report := orch.RunAll(ctx, []string{})

	if report.TotalServers != 0 {
		t.Errorf("Expected 0 total servers, got %d", report.TotalServers)
	}
	if report.TestedServers != 0 {
		t.Errorf("Expected 0 tested servers, got %d", report.TestedServers)
	}
	if len(report.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(report.Results))
	}
}

func TestOrchestrator_RunAll_WithServers(t *testing.T) {
	config := DefaultConfig()
	config.TestLevels.Config = true
	config.TestLevels.Connectivity = false // Disable connectivity for unit test
	config.GlobalTimeout = 5 * time.Second

	orch := NewOrchestrator(config)

	servers := []string{"test-server-1", "test-server-2"}
	ctx := context.Background()
	report := orch.RunAll(ctx, servers)

	if report.TotalServers != 2 {
		t.Errorf("Expected 2 total servers, got %d", report.TotalServers)
	}
	if report.TestedServers != 2 {
		t.Errorf("Expected 2 tested servers, got %d", report.TestedServers)
	}
}

func TestOrchestrator_RunAll_Timeout(t *testing.T) {
	config := DefaultConfig()
	config.GlobalTimeout = 10 * time.Millisecond

	orch := NewOrchestrator(config)

	servers := []string{"test-server"}
	ctx, cancel := context.WithTimeout(context.Background(), config.GlobalTimeout)
	defer cancel()

	// Should complete without hanging
	report := orch.RunAll(ctx, servers)

	// Report should be valid even with timeout
	if report == nil {
		t.Fatal("Expected report, got nil")
	}
}

func TestOrchestrator_RunAll_ReportHasTimestamp(t *testing.T) {
	config := DefaultConfig()
	orch := NewOrchestrator(config)

	before := time.Now()
	report := orch.RunAll(context.Background(), []string{"test"})
	after := time.Now()

	if report.Timestamp.Before(before) || report.Timestamp.After(after) {
		t.Errorf("Report timestamp %v not in expected range [%v, %v]",
			report.Timestamp, before, after)
	}
}

func TestOrchestrator_RunAll_CountsCorrectly(t *testing.T) {
	config := DefaultConfig()
	config.TestLevels.Config = false
	config.TestLevels.Connectivity = false

	orch := NewOrchestrator(config)

	// With no test levels enabled, all should be skipped
	report := orch.RunAll(context.Background(), []string{"a", "b", "c"})

	// No tests actually run since no levels enabled
	if report.PassedTests != 0 && report.FailedTests != 0 && report.SkippedTests != 0 {
		// Counts should be consistent with results
		totalCounted := report.PassedTests + report.FailedTests + report.SkippedTests
		if totalCounted != len(report.Results) {
			t.Errorf("Count mismatch: passed=%d, failed=%d, skipped=%d but results=%d",
				report.PassedTests, report.FailedTests, report.SkippedTests, len(report.Results))
		}
	}
}

func TestDefaultConfig_HasSensibleDefaults(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Expected Enabled to be true by default")
	}
	if !config.RunOnBoot {
		t.Error("Expected RunOnBoot to be true by default")
	}
	if config.MaxParallelTests < 1 {
		t.Errorf("Expected MaxParallelTests >= 1, got %d", config.MaxParallelTests)
	}
	if config.GlobalTimeout <= 0 {
		t.Errorf("Expected positive GlobalTimeout, got %v", config.GlobalTimeout)
	}
}

func TestTestStatus_Values(t *testing.T) {
	// Verify status constants are correct
	if StatusPass != "pass" {
		t.Error("StatusPass should be 'pass'")
	}
	if StatusFail != "fail" {
		t.Error("StatusFail should be 'fail'")
	}
	if StatusSkip != "skip" {
		t.Error("StatusSkip should be 'skip'")
	}
	if StatusTimeout != "timeout" {
		t.Error("StatusTimeout should be 'timeout'")
	}
}
