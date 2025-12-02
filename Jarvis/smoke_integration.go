package main

import (
	"context"
	"fmt"
	"jarvis/smoketests"
	"log"
	"os"
	"time"
)

func shouldRunSmokeTests() bool {
	// Check if smoke tests are explicitly disabled
	if os.Getenv("JARVIS_SKIP_SMOKE_TESTS") == "true" {
		return false
	}
	return true
}

func getActiveServers() []string {
	// Return list of servers that should be tested
	// This could be enhanced to detect based on cwd, profiles, etc.
	return []string{
		"firecrawl",
		"context7",
		"kagimcp",
		"brave-search",
		"fetch-mcp",
		"time",
		"basic-memory",
		"mem0-mcp",
	}
}

func runBootSmokeTests() {
	log.Println("🔍 Running smoke tests...")
	fmt.Fprintln(os.Stderr, "\n🔍 Running smoke tests...")

	// Create orchestrator with default config
	config := smoketests.DefaultConfig()
	config.GlobalTimeout = 5 * time.Second // Fast boot-time tests

	orchestrator := smoketests.NewOrchestrator(config)

	// Get servers to test
	servers := getActiveServers()

	ctx, cancel := context.WithTimeout(context.Background(), config.GlobalTimeout)
	defer cancel()

	report := orchestrator.RunAll(ctx, servers)

	// Log summary
	log.Printf("Smoke tests completed: %d passed, %d failed, %d skipped (%.2fs)",
		report.PassedTests, report.FailedTests, report.SkippedTests,
		report.TotalDuration.Seconds())

	// Print warnings to stderr if failures detected
	if len(report.CriticalFailures) > 0 {
		fmt.Fprintln(os.Stderr, "\n⚠️  SMOKE TEST FAILURES DETECTED:")
		for _, failure := range report.CriticalFailures {
			fmt.Fprintf(os.Stderr, "  ❌ %s: %s\n", failure.ServerName, failure.ErrorMessage)
			if failure.FixSuggestion != "" {
				fmt.Fprintf(os.Stderr, "     💡 %s\n", failure.FixSuggestion)
			}
		}
		fmt.Fprintln(os.Stderr, "\n  Use check_status() or run_smoke_tests() for full diagnostics")
	} else if report.FailedTests == 0 {
		fmt.Fprintln(os.Stderr, "✅ All smoke tests passed")
	} else {
		fmt.Fprintf(os.Stderr, "⚠️  %d non-critical test(s) failed\n\n", report.FailedTests)
	}
}
