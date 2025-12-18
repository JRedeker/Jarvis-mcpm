# Spec: MCPM API Server

## Overview
The MCPM API Server exposes MCPM's internal APIs via HTTP/JSON endpoints, enabling Jarvis (Go) to call MCPM functionality without subprocess overhead.

## ADDED Requirements

### Requirement: API Server Module
The MCPM package SHALL include an `api` module that provides an HTTP server exposing management operations.

#### Scenario: Start API server
- **WHEN** user runs `mcpm serve --port 6275`
- **THEN** an HTTP server SHALL start on the specified port
- **AND** the server SHALL respond to `/api/v1/health` with status 200

#### Scenario: Health check endpoint
- **WHEN** a GET request is made to `/api/v1/health`
- **THEN** the response SHALL contain `{"success": true, "data": {"status": "healthy", ...}}`

---

### Requirement: Server Management Endpoints
The API server SHALL expose endpoints for managing MCP servers.

#### Scenario: List all servers
- **WHEN** a GET request is made to `/api/v1/servers`
- **THEN** the response SHALL contain a JSON array of all installed servers
- **AND** each server object SHALL include name, type, command/url, and profile_tags

#### Scenario: Get server info
- **WHEN** a GET request is made to `/api/v1/servers/:name` for an installed server
- **THEN** the response SHALL contain the server's full configuration

#### Scenario: Install server
- **WHEN** a POST request is made to `/api/v1/servers/:name/install`
- **THEN** the server SHALL be installed from the registry
- **AND** the response SHALL contain `{"success": true, "data": {"name": "...", ...}}`

#### Scenario: Install non-existent server
- **WHEN** a POST request is made to `/api/v1/servers/nonexistent/install`
- **THEN** the response SHALL contain `{"success": false, "error": {"code": "SERVER_NOT_FOUND", ...}}`

#### Scenario: Uninstall server
- **WHEN** a DELETE request is made to `/api/v1/servers/:name` for an installed server
- **THEN** the server SHALL be removed from the global configuration
- **AND** the response SHALL contain `{"success": true}`

#### Scenario: Search registry
- **WHEN** a GET request is made to `/api/v1/servers/search?q=github`
- **THEN** the response SHALL contain matching servers from the registry

---

### Requirement: Profile Management Endpoints
The API server SHALL expose endpoints for managing MCPM profiles.

#### Scenario: List profiles
- **WHEN** a GET request is made to `/api/v1/profiles`
- **THEN** the response SHALL contain all profiles and their associated servers

#### Scenario: Create profile
- **WHEN** a POST request is made to `/api/v1/profiles` with `{"name": "test-profile"}`
- **THEN** a new profile SHALL be created
- **AND** the response SHALL contain `{"success": true}`

#### Scenario: Edit profile - add server
- **WHEN** a PUT request is made to `/api/v1/profiles/:name` with `{"add_servers": ["context7"]}`
- **THEN** the server SHALL be added to the profile

#### Scenario: Delete profile
- **WHEN** a DELETE request is made to `/api/v1/profiles/:name`
- **THEN** the profile SHALL be removed

---

### Requirement: Client Management Endpoints
The API server SHALL expose endpoints for managing client configurations.

#### Scenario: List clients
- **WHEN** a GET request is made to `/api/v1/clients`
- **THEN** the response SHALL contain all known MCP clients and their configurations

#### Scenario: Edit client
- **WHEN** a PUT request is made to `/api/v1/clients/:name` with `{"add_profile": "memory"}`
- **THEN** the profile SHALL be added to the client's configuration

---

### Requirement: System Operation Endpoints
The API server SHALL expose endpoints for system-level operations.

#### Scenario: Get usage statistics
- **WHEN** a GET request is made to `/api/v1/usage`
- **THEN** the response SHALL contain tool usage statistics

#### Scenario: Migrate configuration
- **WHEN** a POST request is made to `/api/v1/migrate`
- **THEN** configuration migration SHALL be performed
- **AND** the response SHALL indicate migration status

---

### Requirement: Consistent Response Format
All API responses SHALL follow a consistent JSON structure.

#### Scenario: Successful response
- **WHEN** any API operation succeeds
- **THEN** the response SHALL have structure `{"success": true, "data": {...}, "error": null}`

#### Scenario: Error response
- **WHEN** any API operation fails
- **THEN** the response SHALL have structure `{"success": false, "data": null, "error": {"code": "...", "message": "..."}}`

#### Scenario: HTTP status codes
- **WHEN** an API operation completes
- **THEN** successful operations SHALL return 200 (OK) or 201 (Created)
- **AND** validation errors SHALL return 400 (Bad Request)
- **AND** resource not found SHALL return 404 (Not Found)
- **AND** internal errors SHALL return 500 (Internal Server Error)
