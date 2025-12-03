// Package mocks provides mock implementations for testing Jarvis handlers
package mocks

import (
	"context"
)

// CommandExecutor defines the interface for executing shell commands
// This allows us to mock mcpm, docker, and git commands in tests
type CommandExecutor interface {
	Execute(ctx context.Context, name string, args ...string) (string, error)
	ExecuteWithEnv(ctx context.Context, name string, env []string, args ...string) (string, error)
}

// McpmClient defines the interface for MCPM operations
type McpmClient interface {
	Doctor(ctx context.Context) (*DoctorResult, error)
	Install(ctx context.Context, name string) (*InstallResult, error)
	Uninstall(ctx context.Context, name string) error
	List(ctx context.Context) ([]ServerInfo, error)
	Search(ctx context.Context, query string) ([]ServerInfo, error)
	Info(ctx context.Context, name string) (*ServerInfo, error)

	// Profile operations
	ProfileList(ctx context.Context) ([]ProfileInfo, error)
	ProfileCreate(ctx context.Context, name string) error
	ProfileEdit(ctx context.Context, name string, opts ProfileEditOpts) error
	ProfileDelete(ctx context.Context, name string) error

	// Client operations
	ClientList(ctx context.Context) ([]ClientInfo, error)
	ClientEdit(ctx context.Context, name string, opts ClientEditOpts) error

	// Config operations
	ConfigGet(ctx context.Context, key string) (string, error)
	ConfigSet(ctx context.Context, key, value string) error
	ConfigList(ctx context.Context) (map[string]string, error)

	// Migration
	Migrate(ctx context.Context) (*MigrateResult, error)
}

// DockerClient defines the interface for Docker operations
type DockerClient interface {
	ComposeUp(ctx context.Context, services ...string) error
	ComposeDown(ctx context.Context) error
	ComposeRestart(ctx context.Context, services ...string) error
	ComposePs(ctx context.Context) ([]ContainerStatus, error)
	ExecSupervisorctl(ctx context.Context, action, target string) (string, error)
}

// GitClient defines the interface for Git operations
type GitClient interface {
	Status(ctx context.Context) (string, error)
	Diff(ctx context.Context, staged bool) (string, error)
	Init(ctx context.Context) error
	IsRepo(ctx context.Context) bool
}

// Result types

// DoctorResult contains the health check results from mcpm doctor
type DoctorResult struct {
	MCPMInstalled bool
	MCPMVersion   string
	PythonOK      bool
	PythonVersion string
	NodeOK        bool
	NodeVersion   string
	NpmVersion    string
	ConfigOK      bool
	ConfigPath    string
	AllHealthy    bool
	Issues        []string
	RawOutput     string
}

// InstallResult contains the result of an install operation
type InstallResult struct {
	Success          bool
	AlreadyInstalled bool
	Message          string
}

// ServerInfo contains information about an MCP server
type ServerInfo struct {
	Name           string
	Description    string
	Installed      bool
	Profiles       []string
	Command        string
	Args           []string
	Env            map[string]string
	URL            string
	Transport      string
	InstallMethods []string
}

// ProfileInfo contains information about an MCPM profile
type ProfileInfo struct {
	Name    string
	Servers []string
}

// ProfileEditOpts contains options for editing a profile
type ProfileEditOpts struct {
	NewName       string
	AddServers    []string
	RemoveServers []string
}

// ClientInfo contains information about an AI client
type ClientInfo struct {
	Name       string
	ConfigPath string
	Installed  bool
	Servers    []string
	Profiles   []string
}

// ClientEditOpts contains options for editing a client
type ClientEditOpts struct {
	AddServers     []string
	RemoveServers  []string
	AddProfiles    []string
	RemoveProfiles []string
}

// MigrateResult contains the result of a config migration
type MigrateResult struct {
	Needed     bool
	Success    bool
	Message    string
	BackupPath string
}

// ContainerStatus contains Docker container status information
type ContainerStatus struct {
	Name    string
	Status  string
	Health  string
	Running bool
	Ports   []string
}
