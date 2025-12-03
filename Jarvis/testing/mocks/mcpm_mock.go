package mocks

import (
	"context"
	"fmt"
	"sync"
)

// MockMcpmClient is a mock implementation of McpmClient for testing
type MockMcpmClient struct {
	mu sync.Mutex

	// Configured responses
	DoctorResponse      *DoctorResult
	DoctorError         error
	InstallResponses    map[string]*InstallResult
	InstallError        error
	UninstallError      error
	ListResponse        []ServerInfo
	ListError           error
	SearchResponses     map[string][]ServerInfo
	SearchError         error
	InfoResponses       map[string]*ServerInfo
	InfoError           error
	ProfileListResponse []ProfileInfo
	ProfileListError    error
	ProfileCreateError  error
	ProfileEditError    error
	ProfileDeleteError  error
	ClientListResponse  []ClientInfo
	ClientListError     error
	ClientEditError     error
	ConfigValues        map[string]string
	ConfigGetError      error
	ConfigSetError      error
	MigrateResponse     *MigrateResult
	MigrateError        error

	// Call tracking
	Calls []MockCall
}

// MockCall records a call to a mock method
type MockCall struct {
	Method string
	Args   []interface{}
}

// NewMockMcpmClient creates a new mock MCPM client with sensible defaults
func NewMockMcpmClient() *MockMcpmClient {
	return &MockMcpmClient{
		InstallResponses: make(map[string]*InstallResult),
		SearchResponses:  make(map[string][]ServerInfo),
		InfoResponses:    make(map[string]*ServerInfo),
		ConfigValues:     make(map[string]string),
		Calls:            make([]MockCall, 0),
	}
}

// recordCall records a method call for later verification
func (m *MockMcpmClient) recordCall(method string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockCall{Method: method, Args: args})
}

// ResetCalls clears the call history
func (m *MockMcpmClient) ResetCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = make([]MockCall, 0)
}

// AssertCalled verifies a method was called with the given arguments
func (m *MockMcpmClient) AssertCalled(method string, args ...interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, call := range m.Calls {
		if call.Method == method {
			if len(args) == 0 {
				return true
			}
			if len(call.Args) == len(args) {
				match := true
				for i, arg := range args {
					if call.Args[i] != arg {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// AssertNotCalled verifies a method was NOT called
func (m *MockMcpmClient) AssertNotCalled(method string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, call := range m.Calls {
		if call.Method == method {
			return false
		}
	}
	return true
}

// CallCount returns the number of times a method was called
func (m *MockMcpmClient) CallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for _, call := range m.Calls {
		if call.Method == method {
			count++
		}
	}
	return count
}

// Doctor implements McpmClient.Doctor
func (m *MockMcpmClient) Doctor(ctx context.Context) (*DoctorResult, error) {
	m.recordCall("Doctor")

	if m.DoctorError != nil {
		return nil, m.DoctorError
	}

	if m.DoctorResponse != nil {
		return m.DoctorResponse, nil
	}

	// Default healthy response
	return &DoctorResult{
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
	}, nil
}

// Install implements McpmClient.Install
func (m *MockMcpmClient) Install(ctx context.Context, name string) (*InstallResult, error) {
	m.recordCall("Install", name)

	if m.InstallError != nil {
		return nil, m.InstallError
	}

	if result, ok := m.InstallResponses[name]; ok {
		return result, nil
	}

	// Default success response
	return &InstallResult{
		Success: true,
		Message: fmt.Sprintf("Successfully installed %s", name),
	}, nil
}

// Uninstall implements McpmClient.Uninstall
func (m *MockMcpmClient) Uninstall(ctx context.Context, name string) error {
	m.recordCall("Uninstall", name)
	return m.UninstallError
}

// List implements McpmClient.List
func (m *MockMcpmClient) List(ctx context.Context) ([]ServerInfo, error) {
	m.recordCall("List")

	if m.ListError != nil {
		return nil, m.ListError
	}

	if m.ListResponse != nil {
		return m.ListResponse, nil
	}

	// Default empty list
	return []ServerInfo{}, nil
}

// Search implements McpmClient.Search
func (m *MockMcpmClient) Search(ctx context.Context, query string) ([]ServerInfo, error) {
	m.recordCall("Search", query)

	if m.SearchError != nil {
		return nil, m.SearchError
	}

	if results, ok := m.SearchResponses[query]; ok {
		return results, nil
	}

	// Default empty results
	return []ServerInfo{}, nil
}

// Info implements McpmClient.Info
func (m *MockMcpmClient) Info(ctx context.Context, name string) (*ServerInfo, error) {
	m.recordCall("Info", name)

	if m.InfoError != nil {
		return nil, m.InfoError
	}

	if info, ok := m.InfoResponses[name]; ok {
		return info, nil
	}

	return nil, fmt.Errorf("server '%s' not found", name)
}

// ProfileList implements McpmClient.ProfileList
func (m *MockMcpmClient) ProfileList(ctx context.Context) ([]ProfileInfo, error) {
	m.recordCall("ProfileList")

	if m.ProfileListError != nil {
		return nil, m.ProfileListError
	}

	if m.ProfileListResponse != nil {
		return m.ProfileListResponse, nil
	}

	return []ProfileInfo{}, nil
}

// ProfileCreate implements McpmClient.ProfileCreate
func (m *MockMcpmClient) ProfileCreate(ctx context.Context, name string) error {
	m.recordCall("ProfileCreate", name)
	return m.ProfileCreateError
}

// ProfileEdit implements McpmClient.ProfileEdit
func (m *MockMcpmClient) ProfileEdit(ctx context.Context, name string, opts ProfileEditOpts) error {
	m.recordCall("ProfileEdit", name, opts)
	return m.ProfileEditError
}

// ProfileDelete implements McpmClient.ProfileDelete
func (m *MockMcpmClient) ProfileDelete(ctx context.Context, name string) error {
	m.recordCall("ProfileDelete", name)
	return m.ProfileDeleteError
}

// ClientList implements McpmClient.ClientList
func (m *MockMcpmClient) ClientList(ctx context.Context) ([]ClientInfo, error) {
	m.recordCall("ClientList")

	if m.ClientListError != nil {
		return nil, m.ClientListError
	}

	if m.ClientListResponse != nil {
		return m.ClientListResponse, nil
	}

	return []ClientInfo{}, nil
}

// ClientEdit implements McpmClient.ClientEdit
func (m *MockMcpmClient) ClientEdit(ctx context.Context, name string, opts ClientEditOpts) error {
	m.recordCall("ClientEdit", name, opts)
	return m.ClientEditError
}

// ConfigGet implements McpmClient.ConfigGet
func (m *MockMcpmClient) ConfigGet(ctx context.Context, key string) (string, error) {
	m.recordCall("ConfigGet", key)

	if m.ConfigGetError != nil {
		return "", m.ConfigGetError
	}

	if value, ok := m.ConfigValues[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("config key '%s' not found", key)
}

// ConfigSet implements McpmClient.ConfigSet
func (m *MockMcpmClient) ConfigSet(ctx context.Context, key, value string) error {
	m.recordCall("ConfigSet", key, value)

	if m.ConfigSetError != nil {
		return m.ConfigSetError
	}

	m.ConfigValues[key] = value
	return nil
}

// ConfigList implements McpmClient.ConfigList
func (m *MockMcpmClient) ConfigList(ctx context.Context) (map[string]string, error) {
	m.recordCall("ConfigList")

	if m.ConfigGetError != nil {
		return nil, m.ConfigGetError
	}

	return m.ConfigValues, nil
}

// Migrate implements McpmClient.Migrate
func (m *MockMcpmClient) Migrate(ctx context.Context) (*MigrateResult, error) {
	m.recordCall("Migrate")

	if m.MigrateError != nil {
		return nil, m.MigrateError
	}

	if m.MigrateResponse != nil {
		return m.MigrateResponse, nil
	}

	return &MigrateResult{
		Needed:  false,
		Success: true,
		Message: "Already up to date",
	}, nil
}

// Helper methods for setting up mock responses

// WithHealthyDoctor configures a healthy doctor response
func (m *MockMcpmClient) WithHealthyDoctor() *MockMcpmClient {
	m.DoctorResponse = &DoctorResult{
		MCPMInstalled: true,
		MCPMVersion:   "2.9.0",
		PythonOK:      true,
		PythonVersion: "3.13.9",
		NodeOK:        true,
		NodeVersion:   "v22.21.0",
		NpmVersion:    "11.6.2",
		ConfigOK:      true,
		AllHealthy:    true,
	}
	return m
}

// WithUnhealthyDoctor configures an unhealthy doctor response
func (m *MockMcpmClient) WithUnhealthyDoctor(issues ...string) *MockMcpmClient {
	m.DoctorResponse = &DoctorResult{
		MCPMInstalled: true,
		MCPMVersion:   "2.9.0",
		PythonOK:      true,
		PythonVersion: "3.13.9",
		NodeOK:        false,
		AllHealthy:    false,
		Issues:        issues,
	}
	return m
}

// WithServers configures the list of installed servers
func (m *MockMcpmClient) WithServers(servers ...ServerInfo) *MockMcpmClient {
	m.ListResponse = servers
	return m
}

// WithProfiles configures the list of profiles
func (m *MockMcpmClient) WithProfiles(profiles ...ProfileInfo) *MockMcpmClient {
	m.ProfileListResponse = profiles
	return m
}

// WithClients configures the list of clients
func (m *MockMcpmClient) WithClients(clients ...ClientInfo) *MockMcpmClient {
	m.ClientListResponse = clients
	return m
}

// WithInstallResult configures a specific install result
func (m *MockMcpmClient) WithInstallResult(name string, result *InstallResult) *MockMcpmClient {
	m.InstallResponses[name] = result
	return m
}

// WithServerInfo configures a specific server info
func (m *MockMcpmClient) WithServerInfo(name string, info *ServerInfo) *MockMcpmClient {
	m.InfoResponses[name] = info
	return m
}

// WithSearchResults configures search results for a query
func (m *MockMcpmClient) WithSearchResults(query string, results []ServerInfo) *MockMcpmClient {
	m.SearchResponses[query] = results
	return m
}
