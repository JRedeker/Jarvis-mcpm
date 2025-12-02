package smoketests

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

// ConfigTestSuite runs configuration validation tests
type ConfigTestSuite struct {
	serverName string
	checks     []ConfigCheck
	enabled    bool
}

// ConfigCheck defines a configuration validation check
type ConfigCheck struct {
	Type          string // "env_var", "file", "permission"
	Name          string
	Required      bool
	Pattern       string // Regex pattern for validation
	ErrorMessage  string
	FixSuggestion string
}

// NewConfigTestSuite creates a new configuration test suite
func NewConfigTestSuite(serverName string, checks []ConfigCheck, enabled bool) *ConfigTestSuite {
	return &ConfigTestSuite{
		serverName: serverName,
		checks:     checks,
		enabled:    enabled,
	}
}

// Run executes all configuration checks
func (s *ConfigTestSuite) Run() []TestResult {
	var results []TestResult

	for _, check := range s.checks {
		startTime := time.Now()
		result := TestResult{
			ServerName: s.serverName,
			TestType:   "config",
			TestName:   check.Name,
			Timestamp:  startTime,
		}

		switch check.Type {
		case "env_var":
			s.checkEnvVar(&result, check)

		case "file":
			s.checkFile(&result, check)

		case "permission":
			s.checkPermission(&result, check)

		default:
			result.Status = StatusSkip
			result.Details = fmt.Sprintf("Unknown check type: %s", check.Type)
		}

		result.Duration = time.Since(startTime)
		results = append(results, result)
	}

	return results
}

func (s *ConfigTestSuite) checkEnvVar(result *TestResult, check ConfigCheck) {
	value := os.Getenv(check.Name)

	if value == "" {
		if check.Required {
			result.Status = StatusFail
			result.ErrorMessage = check.ErrorMessage
			if result.ErrorMessage == "" {
				result.ErrorMessage = fmt.Sprintf("Required environment variable %s is not set", check.Name)
			}
			result.FixSuggestion = check.FixSuggestion
		} else {
			result.Status = StatusSkip
			result.Details = "Optional variable not set"
		}
		return
	}

	// Validate pattern if provided
	if check.Pattern != "" {
		matched, err := regexp.MatchString(check.Pattern, value)
		if err != nil {
			result.Status = StatusFail
			result.ErrorMessage = fmt.Sprintf("Invalid pattern: %v", err)
			return
		}

		if !matched {
			result.Status = StatusFail
			result.ErrorMessage = check.ErrorMessage
			if result.ErrorMessage == "" {
				result.ErrorMessage = fmt.Sprintf("Environment variable %s has invalid format", check.Name)
			}
			result.FixSuggestion = check.FixSuggestion
			result.Details = fmt.Sprintf("Value does not match pattern: %s", check.Pattern)
			return
		}
	}

	result.Status = StatusPass
	result.Details = "Environment variable is set and valid"
}

func (s *ConfigTestSuite) checkFile(result *TestResult, check ConfigCheck) {
	info, err := os.Stat(check.Name)

	if os.IsNotExist(err) {
		if check.Required {
			result.Status = StatusFail
			result.ErrorMessage = check.ErrorMessage
			if result.ErrorMessage == "" {
				result.ErrorMessage = fmt.Sprintf("Required file not found: %s", check.Name)
			}
			result.FixSuggestion = check.FixSuggestion
		} else {
			result.Status = StatusSkip
			result.Details = "Optional file not present"
		}
		return
	}

	if err != nil {
		result.Status = StatusFail
		result.ErrorMessage = fmt.Sprintf("Error accessing file: %v", err)
		return
	}

	result.Status = StatusPass
	result.Details = fmt.Sprintf("File exists (size: %d bytes)", info.Size())
}

func (s *ConfigTestSuite) checkPermission(result *TestResult, check ConfigCheck) {
	info, err := os.Stat(check.Name)

	if os.IsNotExist(err) {
		result.Status = StatusFail
		result.ErrorMessage = fmt.Sprintf("File not found: %s", check.Name)
		return
	}

	if err != nil {
		result.Status = StatusFail
		result.ErrorMessage = fmt.Sprintf("Error accessing file: %v", err)
		return
	}

	// Check if file is readable
	file, err := os.Open(check.Name)
	if err != nil {
		result.Status = StatusFail
		result.ErrorMessage = fmt.Sprintf("File is not readable: %v", err)
		result.FixSuggestion = fmt.Sprintf("chmod +r %s", check.Name)
		return
	}
	file.Close()

	result.Status = StatusPass
	result.Details = fmt.Sprintf("File has correct permissions (mode: %s)", info.Mode())
}

// Name returns the name of this test suite
func (s *ConfigTestSuite) Name() string {
	return fmt.Sprintf("%s_config", s.serverName)
}

// Enabled returns whether this test suite is enabled
func (s *ConfigTestSuite) Enabled() bool {
	return s.enabled
}
