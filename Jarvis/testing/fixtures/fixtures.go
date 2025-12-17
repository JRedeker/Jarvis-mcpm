// Package fixtures provides test data for Jarvis handler tests
package fixtures

import "jarvis/testing/mocks"

// Common server definitions used across tests

// Context7Server is the context7 documentation server
var Context7Server = mocks.ServerInfo{
	Name:        "context7",
	Description: "Documentation lookup for libraries and frameworks",
	Installed:   true,
	Profiles:    []string{"p-pokeedge", "p-new", "testing-all-tools"},
	Command:     "npx",
	Args:        []string{"-y", "@anthropic-ai/context7"},
	Transport:   "stdio",
}

// BraveSearchServer is the brave search server
var BraveSearchServer = mocks.ServerInfo{
	Name:        "brave-search",
	Description: "Web search using Brave Search API",
	Installed:   true,
	Profiles:    []string{"testing-all-tools"},
	Command:     "npx",
	Args:        []string{"-y", "@anthropic-ai/brave-search"},
	Env:         map[string]string{"BRAVE_API_KEY": "test-key"},
	Transport:   "stdio",
}

// BasicMemoryServer is the basic memory server
var BasicMemoryServer = mocks.ServerInfo{
	Name:        "basic-memory",
	Description: "Simple persistent memory storage",
	Installed:   true,
	Profiles:    []string{"memory", "testing-all-tools"},
	Command:     "uvx",
	Args:        []string{"basic-memory", "mcp"},
	Transport:   "stdio",
}

// Mem0Server is the mem0 memory server
var Mem0Server = mocks.ServerInfo{
	Name:        "mem0-mcp",
	Description: "Advanced memory with semantic search",
	Installed:   true,
	Profiles:    []string{"memory"},
	Command:     "python",
	Args:        []string{"/path/to/stdio_client.py"},
	Transport:   "stdio",
}

// FirecrawlServer is the firecrawl web scraping server
var FirecrawlServer = mocks.ServerInfo{
	Name:        "firecrawl",
	Description: "Web scraping and content extraction",
	Installed:   true,
	Profiles:    []string{"p-pokeedge", "p-new", "testing-all-tools"},
	Env:         map[string]string{"FIRECRAWL_API_KEY": "test-key"},
	Transport:   "stdio",
}

// MorphServer is the morph fast apply server
var MorphServer = mocks.ServerInfo{
	Name:        "morph-fast-apply",
	Description: "AI-powered code refactoring",
	Installed:   true,
	Profiles:    []string{"morph", "testing-all-tools"},
	Transport:   "stdio",
}

// UninstalledServer represents a server that exists in registry but not installed
var UninstalledServer = mocks.ServerInfo{
	Name:           "playwright",
	Description:    "Browser automation and testing",
	Installed:      false,
	Profiles:       []string{},
	InstallMethods: []string{"npm", "docker"},
}

// AllInstalledServers is a collection of typical installed servers
var AllInstalledServers = []mocks.ServerInfo{
	Context7Server,
	BraveSearchServer,
	BasicMemoryServer,
	Mem0Server,
	FirecrawlServer,
	MorphServer,
}

// Common profile definitions

// PokeedgeProfile is the pokeedge project profile
var PokeedgeProfile = mocks.ProfileInfo{
	Name:    "p-pokeedge",
	Servers: []string{"context7", "kagimcp", "time", "firecrawl", "fetch-mcp"},
}

// MemoryProfile is the memory layer profile
var MemoryProfile = mocks.ProfileInfo{
	Name:    "memory",
	Servers: []string{"basic-memory", "mem0-mcp"},
}

// MorphProfile is the morph capability profile
var MorphProfile = mocks.ProfileInfo{
	Name:    "morph",
	Servers: []string{"morph-fast-apply"},
}

// TestingProfile is the testing-all-tools profile
var TestingProfile = mocks.ProfileInfo{
	Name: "testing-all-tools",
	Servers: []string{
		"brave-search", "context7", "kagimcp", "magic-mcp", "arxiv-mcp",
		"basic-memory", "time", "morph-fast-apply", "firecrawl", "mcp-server-qdrant",
	},
}

// AllProfiles is a collection of typical profiles
var AllProfiles = []mocks.ProfileInfo{
	PokeedgeProfile,
	MemoryProfile,
	MorphProfile,
	TestingProfile,
}

// Common client definitions

// OpenCodeClient is the OpenCode AI coding agent
var OpenCodeClient = mocks.ClientInfo{
	Name:       "opencode",
	ConfigPath: "/home/test/.config/opencode/opencode.json",
	Installed:  true,
	Servers:    []string{"jarvis"},
	Profiles:   []string{"p-pokeedge", "memory", "morph"},
}

// ClaudeCodeClient is the Claude Code CLI client
var ClaudeCodeClient = mocks.ClientInfo{
	Name:       "claude-code",
	ConfigPath: "/home/test/.claude.json",
	Installed:  true,
	Servers:    []string{"jarvis"},
	Profiles:   []string{"p-pokeedge", "memory"},
}

// ClaudeDesktopClient is the Claude Desktop client
var ClaudeDesktopClient = mocks.ClientInfo{
	Name:       "claude-desktop",
	ConfigPath: "/home/test/.config/Claude/claude_desktop_config.json",
	Installed:  true,
	Servers:    []string{"jarvis"},
	Profiles:   []string{"p-pokeedge", "memory"},
}

// CodexClient is the Codex CLI client
var CodexClient = mocks.ClientInfo{
	Name:       "codex",
	ConfigPath: "/home/test/.codex/config.json",
	Installed:  false,
	Servers:    []string{},
	Profiles:   []string{},
}

// AllClients is a collection of typical clients
var AllClients = []mocks.ClientInfo{
	OpenCodeClient,
	ClaudeCodeClient,
	ClaudeDesktopClient,
	CodexClient,
}

// Doctor output fixtures

// HealthyDoctorResult is a fully healthy doctor result
var HealthyDoctorResult = &mocks.DoctorResult{
	MCPMInstalled: true,
	MCPMVersion:   "2.9.0",
	PythonOK:      true,
	PythonVersion: "3.13.9",
	NodeOK:        true,
	NodeVersion:   "v22.21.0",
	NpmVersion:    "11.6.2",
	ConfigOK:      true,
	ConfigPath:    "/home/test/.config/mcpm/config.json",
	AllHealthy:    true,
	Issues:        []string{},
	RawOutput: `🩺 MCPM System Health Check

📦 MCPM Installation
  ✅ MCPM version: 2.9.0
🐍 Python Environment
  ✅ Python version: 3.13.9
📊 Node.js Environment
  ✅ Node.js version: v22.21.0
  ✅ npm version: 11.6.2
⚙️  MCPM Configuration
  ✅ Config file: /home/test/.config/mcpm/config.json

✅ All systems healthy! No issues found.`,
}

// UnhealthyDoctorResult has issues
var UnhealthyDoctorResult = &mocks.DoctorResult{
	MCPMInstalled: true,
	MCPMVersion:   "2.9.0",
	PythonOK:      true,
	PythonVersion: "3.13.9",
	NodeOK:        false,
	NodeVersion:   "",
	NpmVersion:    "",
	ConfigOK:      true,
	AllHealthy:    false,
	Issues:        []string{"Node.js not installed", "npm not available"},
	RawOutput: `🩺 MCPM System Health Check

📦 MCPM Installation
  ✅ MCPM version: 2.9.0
🐍 Python Environment
  ✅ Python version: 3.13.9
📊 Node.js Environment
  ❌ Node.js not installed
  ❌ npm not available

❌ Issues found: 2`,
}

// Container status fixtures

// HealthyContainers is all containers healthy
var HealthyContainers = []mocks.ContainerStatus{
	{Name: "mcp-postgres", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"5432:5432"}},
	{Name: "mcp-qdrant", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"6333:6333", "6334:6334"}},
	{Name: "mcp-daemon", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"6276:6276", "6277:6277", "6278:6278"}},
}

// UnhealthyDaemon is daemon unhealthy
var UnhealthyDaemon = []mocks.ContainerStatus{
	{Name: "mcp-postgres", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"5432:5432"}},
	{Name: "mcp-qdrant", Status: "Up 2 hours", Health: "healthy", Running: true, Ports: []string{"6333:6333", "6334:6334"}},
	{Name: "mcp-daemon", Status: "Up 2 hours", Health: "unhealthy", Running: true, Ports: []string{"6276:6276"}},
}

// NoContainers is no containers running
var NoContainers = []mocks.ContainerStatus{}

// Search result fixtures

// MemorySearchResults is results for "memory" query
var MemorySearchResults = []mocks.ServerInfo{
	BasicMemoryServer,
	Mem0Server,
	{
		Name:        "memory-graph",
		Description: "Graph-based memory storage",
		Installed:   false,
	},
}

// DocumentationSearchResults is results for "documentation" query
var DocumentationSearchResults = []mocks.ServerInfo{
	Context7Server,
	{
		Name:        "readme-gen",
		Description: "Generate documentation from code",
		Installed:   false,
	},
}

// EmptySearchResults is no results found
var EmptySearchResults = []mocks.ServerInfo{}

// Install result fixtures

// SuccessInstall is successful installation
var SuccessInstall = &mocks.InstallResult{
	Success: true,
	Message: "Successfully installed",
}

// AlreadyInstalledResult is already installed
var AlreadyInstalledResult = &mocks.InstallResult{
	Success:          true,
	AlreadyInstalled: true,
	Message:          "Already installed",
}

// FailedInstallResult is installation failed
var FailedInstallResult = &mocks.InstallResult{
	Success: false,
	Message: "Installation failed: package not found",
}
