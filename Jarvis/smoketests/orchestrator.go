package smoketests

import (
	"context"
	"sync"
	"time"
)

// Orchestrator manages the execution of smoke tests
type Orchestrator struct {
	config *Config
}

// NewOrchestrator creates a new test orchestrator
func NewOrchestrator(config *Config) *Orchestrator {
	return &Orchestrator{
		config: config,
	}
}

// RunAll executes tests for all specified servers
func (o *Orchestrator) RunAll(ctx context.Context, servers []string) *SmokeTestReport {
	startTime := time.Now()
	report := &SmokeTestReport{
		TotalServers: len(servers),
		Timestamp:    startTime,
		Results:      []TestResult{},
	}

	var wg sync.WaitGroup
	resultsChan := make(chan []TestResult, len(servers))
	semaphore := make(chan struct{}, o.config.MaxParallelTests)

	for _, server := range servers {
		wg.Add(1)
		go func(srv string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			var serverResults []TestResult

			// 1. Config Tests
			if o.config.TestLevels.Config {
				// Create a default config suite for now
				// In a real implementation, we'd fetch server-specific rules
				suite := NewConfigTestSuite(srv, []ConfigCheck{}, true)
				serverResults = append(serverResults, suite.Run()...)
			}

			// 2. Connectivity Tests
			if o.config.TestLevels.Connectivity {
				// Basic connectivity check
				// In a real implementation, we'd know the port/url
				checks := []ConnectivityCheck{}
				suite := NewConnectivityTestSuite(srv, checks, true, 2*time.Second)
				serverResults = append(serverResults, suite.Run()...)
			}

			resultsChan <- serverResults
		}(server)
	}

	wg.Wait()
	close(resultsChan)

	for res := range resultsChan {
		report.Results = append(report.Results, res...)
		for _, r := range res {
			if r.Status == StatusFail {
				report.FailedTests++
				report.CriticalFailures = append(report.CriticalFailures, r)
			} else if r.Status == StatusPass {
				report.PassedTests++
			} else {
				report.SkippedTests++
			}
		}
	}

	report.TestedServers = len(servers) // Simplified
	report.TotalDuration = time.Since(startTime)

	return report
}
