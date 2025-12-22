package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// MockFileSystemForOpenCode implements FileSystem for testing OpenCode functionality
type MockFileSystemForOpenCode struct {
	files    map[string][]byte
	dirs     map[string]bool
	cwd      string
	statErr  map[string]error
	readErr  map[string]error
	writeErr map[string]error
}

func NewMockFileSystemForOpenCode() *MockFileSystemForOpenCode {
	return &MockFileSystemForOpenCode{
		files:    make(map[string][]byte),
		dirs:     make(map[string]bool),
		cwd:      "/home/test/project",
		statErr:  make(map[string]error),
		readErr:  make(map[string]error),
		writeErr: make(map[string]error),
	}
}

func (m *MockFileSystemForOpenCode) Stat(name string) (os.FileInfo, error) {
	if err, ok := m.statErr[name]; ok {
		return nil, err
	}
	if _, ok := m.files[name]; ok {
		return nil, nil // File exists
	}
	if _, ok := m.dirs[name]; ok {
		return nil, nil // Dir exists
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystemForOpenCode) ReadFile(name string) ([]byte, error) {
	if err, ok := m.readErr[name]; ok {
		return nil, err
	}
	if data, ok := m.files[name]; ok {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystemForOpenCode) WriteFile(name string, data []byte, perm os.FileMode) error {
	if err, ok := m.writeErr[name]; ok {
		return err
	}
	m.files[name] = data
	return nil
}

func (m *MockFileSystemForOpenCode) MkdirAll(path string, perm os.FileMode) error {
	m.dirs[path] = true
	return nil
}

func (m *MockFileSystemForOpenCode) ReadDir(name string) ([]os.DirEntry, error) {
	return nil, nil
}

func (m *MockFileSystemForOpenCode) Getwd() (string, error) {
	return m.cwd, nil
}

func (m *MockFileSystemForOpenCode) SetFile(path string, content []byte) {
	m.files[path] = content
}

func (m *MockFileSystemForOpenCode) SetDir(path string) {
	m.dirs[path] = true
}

func TestDetectOpenCodeConfig(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*MockFileSystemForOpenCode)
		envVar      string
		wantPath    string
		wantErr     bool
		errContains string
	}{
		{
			name: "detect from environment variable",
			setup: func(fs *MockFileSystemForOpenCode) {
				fs.SetFile("/custom/path/opencode.json", []byte(`{}`))
			},
			envVar:   "/custom/path/opencode.json",
			wantPath: "/custom/path/opencode.json",
			wantErr:  false,
		},
		{
			name: "detect from project local",
			setup: func(fs *MockFileSystemForOpenCode) {
				fs.cwd = "/home/test/myproject"
				fs.SetFile("/home/test/myproject/opencode.json", []byte(`{}`))
			},
			wantPath: "/home/test/myproject/opencode.json",
			wantErr:  false,
		},
		{
			name: "detect from global config",
			setup: func(fs *MockFileSystemForOpenCode) {
				home, _ := os.UserHomeDir()
				globalPath := filepath.Join(home, ".config", "opencode", "opencode.json")
				fs.SetFile(globalPath, []byte(`{}`))
			},
			wantPath: filepath.Join(os.Getenv("HOME"), ".config", "opencode", "opencode.json"),
			wantErr:  false,
		},
		{
			name:        "no config found",
			setup:       func(fs *MockFileSystemForOpenCode) {},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewMockFileSystemForOpenCode()
			tt.setup(fs)

			// Set environment variable if specified
			if tt.envVar != "" {
				os.Setenv("OPENCODE_CONFIG", tt.envVar)
				defer os.Unsetenv("OPENCODE_CONFIG")
			} else {
				os.Unsetenv("OPENCODE_CONFIG")
			}

			path, err := DetectOpenCodeConfig(fs)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error should contain '%s', got: %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if path != tt.wantPath {
				t.Errorf("Expected path %s, got %s", tt.wantPath, path)
			}
		})
	}
}

func TestReadOpenCodeConfig(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantServers int
		wantErr     bool
	}{
		{
			name: "valid config with multiple servers",
			content: `{
				"$schema": "https://opencode.ai/config.json",
				"mcp": {
					"jarvis": {
						"type": "local",
						"command": ["/path/to/jarvis"],
						"enabled": true
					},
					"memory": {
						"type": "remote",
						"url": "http://localhost:6277/mcp",
						"enabled": true
					}
				}
			}`,
			wantServers: 2,
			wantErr:     false,
		},
		{
			name:        "empty config",
			content:     `{}`,
			wantServers: 0,
			wantErr:     false,
		},
		{
			name:    "invalid json",
			content: `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewMockFileSystemForOpenCode()
			fs.SetFile("/test/config.json", []byte(tt.content))

			config, _, err := ReadOpenCodeConfig(fs, "/test/config.json")

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(config.MCP) != tt.wantServers {
				t.Errorf("Expected %d servers, got %d", tt.wantServers, len(config.MCP))
			}
		})
	}
}

func TestWriteOpenCodeConfig(t *testing.T) {
	fs := NewMockFileSystemForOpenCode()
	enabled := true

	config := &OpenCodeConfig{
		Schema: "https://opencode.ai/config.json",
		MCP: map[string]OpenCodeMCPServer{
			"jarvis": {
				Type:    "local",
				Command: []string{"/path/to/jarvis"},
				Enabled: &enabled,
			},
		},
	}

	err := WriteOpenCodeConfig(fs, "/test/opencode.json", config, nil)
	if err != nil {
		t.Fatalf("WriteOpenCodeConfig failed: %v", err)
	}

	// Verify file was written
	data, ok := fs.files["/test/opencode.json"]
	if !ok {
		t.Fatal("Config file was not written")
	}

	// Verify content
	content := string(data)
	if !strings.Contains(content, "jarvis") {
		t.Error("Config should contain jarvis server")
	}
	if !strings.Contains(content, "local") {
		t.Error("Config should contain local type")
	}
}

func TestAddProfileToOpenCode(t *testing.T) {
	tests := []struct {
		name        string
		profile     string
		jarvisPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "add jarvis profile",
			profile:    "jarvis",
			jarvisPath: "/path/to/jarvis",
			wantErr:    false,
		},
		{
			name:    "add memory profile",
			profile: "memory",
			wantErr: false,
		},
		{
			name:    "add toolbox profile",
			profile: "toolbox",
			wantErr: false,
		},
		{
			name:        "add unknown profile",
			profile:     "unknown-profile",
			wantErr:     true,
			errContains: "unknown profile",
		},
		{
			name:        "add jarvis without path",
			profile:     "jarvis",
			jarvisPath:  "",
			wantErr:     true,
			errContains: "jarvisPath is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewMockFileSystemForOpenCode()
			fs.SetFile("/test/opencode.json", []byte(`{"mcp": {}}`))

			err := AddProfileToOpenCode(fs, "/test/opencode.json", tt.profile, tt.jarvisPath)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error should contain '%s', got: %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify the profile was added
			config, _, err := ReadOpenCodeConfig(fs, "/test/opencode.json")
			if err != nil {
				t.Fatalf("Failed to read config: %v", err)
			}

			if _, exists := config.MCP[tt.profile]; !exists {
				t.Errorf("Profile %s was not added to config", tt.profile)
			}
		})
	}
}

func TestRemoveProfileFromOpenCode(t *testing.T) {
	fs := NewMockFileSystemForOpenCode()
	fs.SetFile("/test/opencode.json", []byte(`{
		"mcp": {
			"memory": {
				"type": "remote",
				"url": "http://localhost:6277/mcp"
			}
		}
	}`))

	// Remove existing profile
	err := RemoveProfileFromOpenCode(fs, "/test/opencode.json", "memory")
	if err != nil {
		t.Errorf("RemoveProfileFromOpenCode failed: %v", err)
	}

	// Verify removal
	config, _, _ := ReadOpenCodeConfig(fs, "/test/opencode.json")
	if _, exists := config.MCP["memory"]; exists {
		t.Error("Profile 'memory' should have been removed")
	}

	// Try to remove non-existent profile
	err = RemoveProfileFromOpenCode(fs, "/test/opencode.json", "nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent profile")
	}
}

func TestListOpenCodeServers(t *testing.T) {
	fs := NewMockFileSystemForOpenCode()
	fs.SetFile("/test/opencode.json", []byte(`{
		"mcp": {
			"jarvis": {
				"type": "local",
				"command": ["/path/to/jarvis"],
				"enabled": true
			},
			"memory": {
				"type": "remote",
				"url": "http://localhost:6277/mcp",
				"enabled": true
			}
		}
	}`))

	output, err := ListOpenCodeServers(fs, "/test/opencode.json")
	if err != nil {
		t.Fatalf("ListOpenCodeServers failed: %v", err)
	}

	if !strings.Contains(output, "jarvis") {
		t.Error("Output should contain 'jarvis'")
	}
	if !strings.Contains(output, "memory") {
		t.Error("Output should contain 'memory'")
	}
	if !strings.Contains(output, "local") {
		t.Error("Output should contain 'local' type")
	}
	if !strings.Contains(output, "remote") {
		t.Error("Output should contain 'remote' type")
	}
}

func TestGenerateOpenCodeTemplate(t *testing.T) {
	template := GenerateOpenCodeTemplate("/path/to/jarvis")

	if !strings.Contains(template, "jarvis") {
		t.Error("Template should contain jarvis")
	}
	if !strings.Contains(template, "/path/to/jarvis") {
		t.Error("Template should contain jarvis path")
	}
	if !strings.Contains(template, "toolbox") {
		t.Error("Template should contain toolbox")
	}
	if !strings.Contains(template, "memory") {
		t.Error("Template should contain memory")
	}
	if !strings.Contains(template, "6276") {
		t.Error("Template should contain port 6276")
	}
	if !strings.Contains(template, "opencode.ai/config.json") {
		t.Error("Template should contain schema URL")
	}
}

func TestKnownClients(t *testing.T) {
	// Verify OpenCode is in the known clients
	opencode, ok := KnownClients["opencode"]
	if !ok {
		t.Fatal("OpenCode should be in KnownClients")
	}

	if opencode.DisplayName != "OpenCode" {
		t.Errorf("Expected DisplayName 'OpenCode', got '%s'", opencode.DisplayName)
	}

	if opencode.Format != "opencode" {
		t.Errorf("Expected Format 'opencode', got '%s'", opencode.Format)
	}

	if len(opencode.ConfigPaths) != 3 {
		t.Errorf("Expected 3 config paths, got %d", len(opencode.ConfigPaths))
	}
}

func TestProfilePorts(t *testing.T) {
	expectedPorts := map[string]int{
		"toolbox": 6276,
		"memory":  6277,
		"morph":   6278,
		"qdrant":  6279,
	}

	for profile, expectedPort := range expectedPorts {
		if port, ok := ProfilePorts[profile]; !ok {
			t.Errorf("Profile %s should be in ProfilePorts", profile)
		} else if port != expectedPort {
			t.Errorf("Profile %s should have port %d, got %d", profile, expectedPort, port)
		}
	}
}
