// Package handlers provides OpenCode client configuration support
package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OpenCodeConfig represents the OpenCode configuration file structure
type OpenCodeConfig struct {
	Schema string                       `json:"$schema,omitempty"`
	MCP    map[string]OpenCodeMCPServer `json:"mcp,omitempty"`
	// Preserve other fields
	Other map[string]interface{} `json:"-"`
}

// OpenCodeMCPServer represents an MCP server in OpenCode config
type OpenCodeMCPServer struct {
	Type        string            `json:"type"`                  // "local" or "remote"
	Command     []string          `json:"command,omitempty"`     // for local type
	URL         string            `json:"url,omitempty"`         // for remote type
	Enabled     *bool             `json:"enabled,omitempty"`     // optional, defaults to true
	Environment map[string]string `json:"environment,omitempty"` // env vars for local
	Headers     map[string]string `json:"headers,omitempty"`     // headers for remote
}

// ClientConfig represents a known client's configuration details
type ClientConfig struct {
	Name        string
	DisplayName string
	ConfigPaths []string // Ordered by priority
	Format      string   // "claude", "opencode", etc.
}

// KnownClients contains all supported AI clients
var KnownClients = map[string]ClientConfig{
	"opencode": {
		Name:        "opencode",
		DisplayName: "OpenCode",
		ConfigPaths: []string{
			"$OPENCODE_CONFIG",
			"./opencode.json",
			"~/.config/opencode/opencode.json",
		},
		Format: "opencode",
	},
	"claude-code": {
		Name:        "claude-code",
		DisplayName: "Claude Code CLI",
		ConfigPaths: []string{
			"~/.claude.json",
		},
		Format: "claude",
	},
	"claude-desktop": {
		Name:        "claude-desktop",
		DisplayName: "Claude Desktop",
		ConfigPaths: []string{
			"~/.config/Claude/claude_desktop_config.json",
		},
		Format: "claude",
	},
}

// ProfilePorts maps profile names to their HTTP endpoints
var ProfilePorts = map[string]int{
	"p-pokeedge": 6276,
	"p-new":      6280,
	"memory":     6277,
	"morph":      6278,
	"qdrant":     6279,
}

// DetectOpenCodeConfig finds the OpenCode configuration file
func DetectOpenCodeConfig(fs FileSystem) (string, error) {
	// Priority 1: Environment variable
	if envPath := os.Getenv("OPENCODE_CONFIG"); envPath != "" {
		expandedPath := expandPath(envPath)
		if _, err := fs.Stat(expandedPath); err == nil {
			return expandedPath, nil
		}
	}

	// Priority 2: Project-local config
	cwd, err := fs.Getwd()
	if err == nil {
		localPath := filepath.Join(cwd, "opencode.json")
		if _, err := fs.Stat(localPath); err == nil {
			return localPath, nil
		}
	}

	// Priority 3: Global config
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(home, ".config", "opencode", "opencode.json")
		if _, err := fs.Stat(globalPath); err == nil {
			return globalPath, nil
		}
	}

	return "", fmt.Errorf("OpenCode configuration not found. Expected locations:\n" +
		"  1. $OPENCODE_CONFIG environment variable\n" +
		"  2. ./opencode.json (project directory)\n" +
		"  3. ~/.config/opencode/opencode.json (global)")
}

// ReadOpenCodeConfig reads and parses an OpenCode configuration file
func ReadOpenCodeConfig(fs FileSystem, path string) (*OpenCodeConfig, map[string]interface{}, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// First unmarshal to get the full raw data
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Then unmarshal the MCP section specifically
	config := &OpenCodeConfig{
		MCP: make(map[string]OpenCodeMCPServer),
	}

	if schema, ok := raw["$schema"].(string); ok {
		config.Schema = schema
	}

	if mcpRaw, ok := raw["mcp"].(map[string]interface{}); ok {
		for name, serverRaw := range mcpRaw {
			serverMap, ok := serverRaw.(map[string]interface{})
			if !ok {
				continue
			}

			server := OpenCodeMCPServer{}
			if t, ok := serverMap["type"].(string); ok {
				server.Type = t
			}
			if url, ok := serverMap["url"].(string); ok {
				server.URL = url
			}
			if cmd, ok := serverMap["command"].([]interface{}); ok {
				for _, c := range cmd {
					if s, ok := c.(string); ok {
						server.Command = append(server.Command, s)
					}
				}
			}
			if enabled, ok := serverMap["enabled"].(bool); ok {
				server.Enabled = &enabled
			}
			if env, ok := serverMap["environment"].(map[string]interface{}); ok {
				server.Environment = make(map[string]string)
				for k, v := range env {
					if s, ok := v.(string); ok {
						server.Environment[k] = s
					}
				}
			}
			if headers, ok := serverMap["headers"].(map[string]interface{}); ok {
				server.Headers = make(map[string]string)
				for k, v := range headers {
					if s, ok := v.(string); ok {
						server.Headers[k] = s
					}
				}
			}

			config.MCP[name] = server
		}
	}

	return config, raw, nil
}

// WriteOpenCodeConfig writes an OpenCode configuration file, preserving other fields
func WriteOpenCodeConfig(fs FileSystem, path string, config *OpenCodeConfig, raw map[string]interface{}) error {
	// Update the raw map with our changes
	if raw == nil {
		raw = make(map[string]interface{})
	}

	if config.Schema != "" {
		raw["$schema"] = config.Schema
	}

	// Convert MCP servers to raw format
	mcpRaw := make(map[string]interface{})
	for name, server := range config.MCP {
		serverMap := make(map[string]interface{})
		serverMap["type"] = server.Type

		if server.Type == "local" && len(server.Command) > 0 {
			serverMap["command"] = server.Command
		}
		if server.Type == "remote" && server.URL != "" {
			serverMap["url"] = server.URL
		}
		if server.Enabled != nil {
			serverMap["enabled"] = *server.Enabled
		}
		if len(server.Environment) > 0 {
			serverMap["environment"] = server.Environment
		}
		if len(server.Headers) > 0 {
			serverMap["headers"] = server.Headers
		}

		mcpRaw[name] = serverMap
	}
	raw["mcp"] = mcpRaw

	// Marshal with indentation
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := fs.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddProfileToOpenCode adds an MCP profile to OpenCode configuration
func AddProfileToOpenCode(fs FileSystem, configPath, profileName, jarvisPath string) error {
	config, raw, err := ReadOpenCodeConfig(fs, configPath)
	if err != nil {
		// If file doesn't exist, create new config
		if os.IsNotExist(err) {
			config = &OpenCodeConfig{
				Schema: "https://opencode.ai/config.json",
				MCP:    make(map[string]OpenCodeMCPServer),
			}
			raw = make(map[string]interface{})
		} else {
			return err
		}
	}

	enabled := true

	// Handle special cases
	switch profileName {
	case "jarvis":
		// Jarvis uses local stdio transport
		if jarvisPath == "" {
			return fmt.Errorf("jarvisPath is required for adding jarvis")
		}
		config.MCP["jarvis"] = OpenCodeMCPServer{
			Type:    "local",
			Command: []string{jarvisPath},
			Enabled: &enabled,
		}
	default:
		// Check if it's a known profile with a port
		port, ok := ProfilePorts[profileName]
		if !ok {
			return fmt.Errorf("unknown profile '%s'. Known profiles: %v", profileName, getProfileNames())
		}
		config.MCP[profileName] = OpenCodeMCPServer{
			Type:    "remote",
			URL:     fmt.Sprintf("http://localhost:%d/mcp", port),
			Enabled: &enabled,
		}
	}

	return WriteOpenCodeConfig(fs, configPath, config, raw)
}

// RemoveProfileFromOpenCode removes an MCP profile from OpenCode configuration
func RemoveProfileFromOpenCode(fs FileSystem, configPath, profileName string) error {
	config, raw, err := ReadOpenCodeConfig(fs, configPath)
	if err != nil {
		return err
	}

	if _, exists := config.MCP[profileName]; !exists {
		return fmt.Errorf("profile '%s' not found in OpenCode configuration", profileName)
	}

	delete(config.MCP, profileName)
	return WriteOpenCodeConfig(fs, configPath, config, raw)
}

// ListOpenCodeServers returns a formatted list of MCP servers in OpenCode config
func ListOpenCodeServers(fs FileSystem, configPath string) (string, error) {
	config, _, err := ReadOpenCodeConfig(fs, configPath)
	if err != nil {
		return "", err
	}

	if len(config.MCP) == 0 {
		return "No MCP servers configured in OpenCode.", nil
	}

	var sb strings.Builder
	sb.WriteString("## OpenCode MCP Servers\n\n")

	for name, server := range config.MCP {
		enabled := "enabled"
		if server.Enabled != nil && !*server.Enabled {
			enabled = "disabled"
		}

		if server.Type == "local" {
			sb.WriteString(fmt.Sprintf("- **%s** (local, %s)\n", name, enabled))
			sb.WriteString(fmt.Sprintf("  Command: `%s`\n", strings.Join(server.Command, " ")))
		} else {
			sb.WriteString(fmt.Sprintf("- **%s** (remote, %s)\n", name, enabled))
			sb.WriteString(fmt.Sprintf("  URL: `%s`\n", server.URL))
		}
	}

	return sb.String(), nil
}

// Helper functions

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func getProfileNames() []string {
	names := make([]string, 0, len(ProfilePorts))
	for name := range ProfilePorts {
		names = append(names, name)
	}
	return names
}

// GenerateOpenCodeTemplate generates a starter OpenCode configuration
func GenerateOpenCodeTemplate(jarvisPath string) string {
	enabled := true
	config := OpenCodeConfig{
		Schema: "https://opencode.ai/config.json",
		MCP: map[string]OpenCodeMCPServer{
			"jarvis": {
				Type:    "local",
				Command: []string{jarvisPath},
				Enabled: &enabled,
			},
			"p-pokeedge": {
				Type:    "remote",
				URL:     "http://localhost:6276/mcp",
				Enabled: &enabled,
			},
			"memory": {
				Type:    "remote",
				URL:     "http://localhost:6277/mcp",
				Enabled: &enabled,
			},
			"morph": {
				Type:    "remote",
				URL:     "http://localhost:6278/mcp",
				Enabled: &enabled,
			},
		},
	}

	raw := map[string]interface{}{
		"$schema": config.Schema,
		"mcp":     config.MCP,
	}

	data, _ := json.MarshalIndent(raw, "", "  ")
	return string(data)
}
