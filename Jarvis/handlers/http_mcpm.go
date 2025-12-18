// Package handlers provides MCP tool handlers with dependency injection for testing
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// HTTPMcpmRunner implements McpmRunner by calling the MCPM API server
type HTTPMcpmRunner struct {
	BaseURL    string
	HTTPClient *http.Client
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *APIError       `json:"error"`
}

// APIError represents an API error response
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewHTTPMcpmRunner creates a new HTTP-based MCPM runner
func NewHTTPMcpmRunner(baseURL string) *HTTPMcpmRunner {
	if baseURL == "" {
		baseURL = os.Getenv("MCPM_API_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:6275"
	}

	return &HTTPMcpmRunner{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run implements McpmRunner interface by mapping CLI args to API calls
func (r *HTTPMcpmRunner) Run(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command specified")
	}

	cmd := args[0]
	switch cmd {
	case "doctor":
		return r.health()
	case "ls":
		return r.listServers()
	case "info":
		if len(args) < 2 {
			return "", fmt.Errorf("server name required")
		}
		return r.serverInfo(args[1])
	case "install":
		if len(args) < 2 {
			return "", fmt.Errorf("server name required")
		}
		return r.installServer(args[1])
	case "uninstall":
		if len(args) < 2 {
			return "", fmt.Errorf("server name required")
		}
		return r.uninstallServer(args[1])
	case "search":
		if len(args) < 2 {
			return "", fmt.Errorf("search query required")
		}
		return r.searchServers(args[1])
	case "new":
		return r.createServer(args[1:])
	case "edit":
		if len(args) < 2 {
			return "", fmt.Errorf("server name required")
		}
		return r.editServer(args[1], args[2:])
	case "profile":
		return r.handleProfile(args[1:])
	case "client":
		return r.handleClient(args[1:])
	case "usage":
		return r.usage()
	case "migrate":
		return r.migrate()
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

// health calls GET /api/v1/health
func (r *HTTPMcpmRunner) health() (string, error) {
	resp, err := r.get("/api/v1/health")
	if err != nil {
		return "", err
	}

	var data struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
		Checks    struct {
			Node struct {
				Status  string `json:"status"`
				Version string `json:"version"`
			} `json:"node"`
			Registry struct {
				Status string `json:"status"`
				Path   string `json:"path"`
			} `json:"registry"`
			Docker struct {
				Status string `json:"status"`
			} `json:"docker"`
		} `json:"checks"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	// Format as text output matching CLI format
	var sb strings.Builder
	sb.WriteString("MCPM System Status:\n")
	sb.WriteString(fmt.Sprintf("- Node.js: %s %s\n", data.Checks.Node.Version, formatStatus(data.Checks.Node.Status)))
	sb.WriteString(fmt.Sprintf("- Config: %s %s\n", data.Checks.Registry.Path, formatStatus(data.Checks.Registry.Status)))
	sb.WriteString(fmt.Sprintf("- Docker: %s\n", formatStatus(data.Checks.Docker.Status)))

	if data.Status == "healthy" {
		sb.WriteString("\nAll systems healthy.")
	} else {
		sb.WriteString("\nSystem degraded - check above for issues.")
	}

	return sb.String(), nil
}

// listServers calls GET /api/v1/servers
func (r *HTTPMcpmRunner) listServers() (string, error) {
	resp, err := r.get("/api/v1/servers")
	if err != nil {
		return "", err
	}

	var data struct {
		Count   int `json:"count"`
		Servers []struct {
			Name        string `json:"name"`
			Group       string `json:"group"`
			Description string `json:"description"`
			Installed   bool   `json:"installed"`
		} `json:"servers"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	// Format as text output
	var sb strings.Builder
	sb.WriteString("Installed MCP Servers:\n")

	currentGroup := ""
	hasInstalled := false
	for _, s := range data.Servers {
		if s.Installed {
			if s.Group != currentGroup {
				sb.WriteString(fmt.Sprintf("\n[%s]\n", s.Group))
				currentGroup = s.Group
			}
			desc := ""
			if s.Description != "" {
				desc = fmt.Sprintf(": %s", s.Description)
			}
			sb.WriteString(fmt.Sprintf("- %s%s (Installed)\n", s.Name, desc))
			hasInstalled = true
		}
	}

	if !hasInstalled {
		sb.WriteString("No servers currently installed. Use \"mcpm install <name>\" to add one.\n")
	}

	return sb.String(), nil
}

// serverInfo calls GET /api/v1/servers/:name
func (r *HTTPMcpmRunner) serverInfo(name string) (string, error) {
	resp, err := r.get(fmt.Sprintf("/api/v1/servers/%s", url.PathEscape(name)))
	if err != nil {
		return "", err
	}

	// Return raw JSON for compatibility with existing handlers
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, resp.Data, "", "  "); err != nil {
		return string(resp.Data), nil
	}
	return prettyJSON.String(), nil
}

// installServer calls POST /api/v1/servers/:name/install
func (r *HTTPMcpmRunner) installServer(name string) (string, error) {
	resp, err := r.post(fmt.Sprintf("/api/v1/servers/%s/install", url.PathEscape(name)), nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Method  string `json:"method"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// uninstallServer calls DELETE /api/v1/servers/:name
func (r *HTTPMcpmRunner) uninstallServer(name string) (string, error) {
	resp, err := r.delete(fmt.Sprintf("/api/v1/servers/%s", url.PathEscape(name)))
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// searchServers calls GET /api/v1/search?q=query
func (r *HTTPMcpmRunner) searchServers(query string) (string, error) {
	resp, err := r.get(fmt.Sprintf("/api/v1/search?q=%s", url.QueryEscape(query)))
	if err != nil {
		return "", err
	}

	var data struct {
		Query   string `json:"query"`
		Count   int    `json:"count"`
		Results []struct {
			Name        string `json:"name"`
			Group       string `json:"group"`
			Description string `json:"description"`
		} `json:"results"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	// Format as text output
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Search results for \"%s\":\n", data.Query))

	if data.Count == 0 {
		sb.WriteString("No servers found.\n")
	} else {
		for _, r := range data.Results {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", r.Name, r.Group))
		}
	}

	return sb.String(), nil
}

// createServer calls POST /api/v1/servers
func (r *HTTPMcpmRunner) createServer(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("server name required")
	}

	body := map[string]interface{}{
		"name":  args[0],
		"force": true, // Always force in API mode
	}

	// Parse remaining args
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--type":
			if i+1 < len(args) {
				body["type"] = args[i+1]
				i++
			}
		case "--command":
			if i+1 < len(args) {
				body["command"] = args[i+1]
				i++
			}
		case "--args":
			if i+1 < len(args) {
				body["args"] = args[i+1]
				i++
			}
		case "--url":
			if i+1 < len(args) {
				body["url"] = args[i+1]
				i++
			}
		case "--env":
			if i+1 < len(args) {
				body["env"] = args[i+1]
				i++
			}
		case "--headers":
			if i+1 < len(args) {
				body["headers"] = args[i+1]
				i++
			}
		}
	}

	resp, err := r.post("/api/v1/servers", body)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// editServer calls PUT /api/v1/servers/:name
func (r *HTTPMcpmRunner) editServer(name string, args []string) (string, error) {
	body := map[string]interface{}{}

	// Parse args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--command":
			if i+1 < len(args) {
				body["command"] = args[i+1]
				i++
			}
		case "--args":
			if i+1 < len(args) {
				body["args"] = args[i+1]
				i++
			}
		case "--url":
			if i+1 < len(args) {
				body["url"] = args[i+1]
				i++
			}
		case "--env":
			if i+1 < len(args) {
				body["env"] = args[i+1]
				i++
			}
		case "--headers":
			if i+1 < len(args) {
				body["headers"] = args[i+1]
				i++
			}
		}
	}

	resp, err := r.put(fmt.Sprintf("/api/v1/servers/%s", url.PathEscape(name)), body)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// handleProfile handles profile subcommands
func (r *HTTPMcpmRunner) handleProfile(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("profile subcommand required")
	}

	action := args[0]
	switch action {
	case "ls":
		return r.listProfiles()
	case "create":
		if len(args) < 2 {
			return "", fmt.Errorf("profile name required")
		}
		return r.createProfile(args[1], args[2:])
	case "edit":
		if len(args) < 2 {
			return "", fmt.Errorf("profile name required")
		}
		return r.editProfile(args[1], args[2:])
	case "rm", "delete":
		if len(args) < 2 {
			return "", fmt.Errorf("profile name required")
		}
		return r.deleteProfile(args[1])
	default:
		return "", fmt.Errorf("unknown profile subcommand: %s", action)
	}
}

// listProfiles calls GET /api/v1/profiles
func (r *HTTPMcpmRunner) listProfiles() (string, error) {
	resp, err := r.get("/api/v1/profiles")
	if err != nil {
		return "", err
	}

	var data struct {
		Count    int `json:"count"`
		Profiles []struct {
			Name    string   `json:"name"`
			Servers []string `json:"servers"`
		} `json:"profiles"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Profiles:\n")

	if data.Count == 0 {
		sb.WriteString("No profiles configured.\n")
	} else {
		for _, p := range data.Profiles {
			sb.WriteString(fmt.Sprintf("- %s (%d servers)\n", p.Name, len(p.Servers)))
		}
	}

	return sb.String(), nil
}

// createProfile calls POST /api/v1/profiles
func (r *HTTPMcpmRunner) createProfile(name string, args []string) (string, error) {
	body := map[string]interface{}{
		"name": name,
	}

	// Parse args for servers
	for i := 0; i < len(args); i++ {
		if args[i] == "--servers" && i+1 < len(args) {
			body["servers"] = strings.Split(args[i+1], ",")
			i++
		}
	}

	resp, err := r.post("/api/v1/profiles", body)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// editProfile calls PUT /api/v1/profiles/:name
func (r *HTTPMcpmRunner) editProfile(name string, args []string) (string, error) {
	body := map[string]interface{}{}

	// Parse args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--name":
			if i+1 < len(args) {
				body["new_name"] = args[i+1]
				i++
			}
		case "--add-server":
			if i+1 < len(args) {
				body["add_servers"] = args[i+1]
				i++
			}
		case "--remove-server":
			if i+1 < len(args) {
				body["remove_servers"] = args[i+1]
				i++
			}
		}
	}

	resp, err := r.put(fmt.Sprintf("/api/v1/profiles/%s", url.PathEscape(name)), body)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// deleteProfile calls DELETE /api/v1/profiles/:name
func (r *HTTPMcpmRunner) deleteProfile(name string) (string, error) {
	resp, err := r.delete(fmt.Sprintf("/api/v1/profiles/%s", url.PathEscape(name)))
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// handleClient handles client subcommands
func (r *HTTPMcpmRunner) handleClient(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("client subcommand required")
	}

	action := args[0]
	switch action {
	case "ls":
		return r.listClients()
	case "edit":
		if len(args) < 2 {
			return "", fmt.Errorf("client name required")
		}
		return r.editClient(args[1], args[2:])
	default:
		return "", fmt.Errorf("unknown client subcommand: %s", action)
	}
}

// listClients calls GET /api/v1/clients
func (r *HTTPMcpmRunner) listClients() (string, error) {
	resp, err := r.get("/api/v1/clients")
	if err != nil {
		return "", err
	}

	var data struct {
		Count   int `json:"count"`
		Clients []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Detected    bool   `json:"detected"`
			ConfigPath  string `json:"configPath"`
		} `json:"clients"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("MCP Clients:\n")

	for _, c := range data.Clients {
		status := "Not detected"
		if c.Detected {
			status = "Detected"
		}
		sb.WriteString(fmt.Sprintf("\n%s\n", c.DisplayName))
		sb.WriteString(fmt.Sprintf("- Status: %s\n", status))
		if c.ConfigPath != "" {
			sb.WriteString(fmt.Sprintf("- Config: %s\n", c.ConfigPath))
		}
	}

	return sb.String(), nil
}

// editClient calls PUT /api/v1/clients/:name
func (r *HTTPMcpmRunner) editClient(name string, args []string) (string, error) {
	body := map[string]interface{}{}

	// Parse args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--set-path":
			if i+1 < len(args) {
				body["config_path"] = args[i+1]
				i++
			}
		case "--add-server":
			if i+1 < len(args) {
				body["add_server"] = args[i+1]
				i++
			}
		case "--remove-server":
			if i+1 < len(args) {
				body["remove_server"] = args[i+1]
				i++
			}
		case "--add-profile":
			if i+1 < len(args) {
				body["add_profile"] = args[i+1]
				i++
			}
		case "--remove-profile":
			if i+1 < len(args) {
				body["remove_profile"] = args[i+1]
				i++
			}
		}
	}

	resp, err := r.put(fmt.Sprintf("/api/v1/clients/%s", url.PathEscape(name)), body)
	if err != nil {
		return "", err
	}

	var data struct {
		Name    string   `json:"name"`
		Changes []string `json:"changes"`
		Message string   `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.Message, nil
}

// usage calls GET /api/v1/usage
func (r *HTTPMcpmRunner) usage() (string, error) {
	resp, err := r.get("/api/v1/usage")
	if err != nil {
		return "", err
	}

	// Return formatted JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, resp.Data, "", "  "); err != nil {
		return string(resp.Data), nil
	}
	return prettyJSON.String(), nil
}

// migrate calls POST /api/v1/migrate
func (r *HTTPMcpmRunner) migrate() (string, error) {
	resp, err := r.post("/api/v1/migrate", nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Migrations []string `json:"migrations"`
		Message    string   `json:"message"`
	}

	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Migration Results:\n")
	for _, m := range data.Migrations {
		sb.WriteString(fmt.Sprintf("- %s\n", m))
	}

	return sb.String(), nil
}

// HTTP helper methods

func (r *HTTPMcpmRunner) get(path string) (*APIResponse, error) {
	return r.doRequest("GET", path, nil)
}

func (r *HTTPMcpmRunner) post(path string, body interface{}) (*APIResponse, error) {
	return r.doRequest("POST", path, body)
}

func (r *HTTPMcpmRunner) put(path string, body interface{}) (*APIResponse, error) {
	return r.doRequest("PUT", path, body)
}

func (r *HTTPMcpmRunner) delete(path string) (*APIResponse, error) {
	return r.doRequest("DELETE", path, nil)
}

func (r *HTTPMcpmRunner) doRequest(method, path string, body interface{}) (*APIResponse, error) {
	url := r.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(respBody))
	}

	if !apiResp.Success && apiResp.Error != nil {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	return &apiResp, nil
}

// formatStatus formats a status string for CLI output
func formatStatus(status string) string {
	switch status {
	case "ok", "healthy", "detected":
		return "OK"
	case "missing", "not_found":
		return "Not Found"
	default:
		return status
	}
}
