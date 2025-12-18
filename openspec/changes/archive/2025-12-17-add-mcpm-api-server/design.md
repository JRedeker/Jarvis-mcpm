# Design: MCPM API Server

## Overview
This document describes the technical design for the MCPM API server that will replace CLI subprocess calls from Jarvis.

## Architecture

### Component Diagram
```
┌─────────────────────────────────────────────────────────────────┐
│                        mcpm-daemon container                     │
│  ┌─────────────────┐    ┌─────────────────┐                     │
│  │  Profile Proxy  │    │  MCPM API Server │ ← NEW              │
│  │  (port 6276+)   │    │  (port 6275)     │                    │
│  └────────┬────────┘    └────────┬─────────┘                    │
│           │                      │                               │
│           ▼                      ▼                               │
│  ┌─────────────────────────────────────────┐                    │
│  │           MCPM Python Library            │                    │
│  │  GlobalConfigManager, ProfileConfigManager│                   │
│  │  RepositoryManager, ClientRegistry       │                    │
│  └─────────────────────────────────────────┘                    │
└─────────────────────────────────────────────────────────────────┘
           ▲
           │ HTTP/JSON
           │
┌──────────┴──────────┐
│   Jarvis (Go MCP)   │
│   HTTPMcpmRunner    │
└─────────────────────┘
```

## API Design

### Framework Choice: FastAPI
- Already familiar in Python ecosystem
- Automatic OpenAPI/JSON Schema generation
- Pydantic integration matches MCPM's existing models
- Async support for future optimization

### Base URL
`http://localhost:6275/api/v1`

### Endpoints

#### Server Management
| Method | Path | Description | Maps to |
|--------|------|-------------|---------|
| GET | `/health` | Health check | `mcpm doctor` |
| GET | `/servers` | List all servers | `mcpm ls` |
| GET | `/servers/{name}` | Get server info | `mcpm info {name}` |
| POST | `/servers/{name}/install` | Install server | `mcpm install {name}` |
| DELETE | `/servers/{name}` | Uninstall server | `mcpm uninstall {name}` |
| GET | `/search?q={query}` | Search registry | `mcpm search {query}` |
| PUT | `/servers/{name}` | Edit server config | `mcpm edit {name}` |
| POST | `/servers` | Create new server | `mcpm new` |

#### Profile Management
| Method | Path | Description | Maps to |
|--------|------|-------------|---------|
| GET | `/profiles` | List profiles | `mcpm profile ls` |
| POST | `/profiles` | Create profile | `mcpm profile create` |
| PUT | `/profiles/{name}` | Edit profile | `mcpm profile edit` |
| DELETE | `/profiles/{name}` | Delete profile | `mcpm profile rm` |

#### Client Management
| Method | Path | Description | Maps to |
|--------|------|-------------|---------|
| GET | `/clients` | List clients | `mcpm client ls` |
| PUT | `/clients/{name}` | Edit client | `mcpm client edit` |

#### System Operations
| Method | Path | Description | Maps to |
|--------|------|-------------|---------|
| GET | `/usage` | Usage statistics | `mcpm usage` |
| POST | `/migrate` | Migrate config | `mcpm migrate` |

### Response Format

All responses follow a consistent structure:

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

Error responses:
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "SERVER_NOT_FOUND",
    "message": "Server 'foo' not found in registry",
    "details": { ... }
  }
}
```

## Jarvis Integration

### Interface Change
The existing `McpmRunner` interface remains, but we add an HTTP implementation:

```go
// Existing interface (unchanged)
type McpmRunner interface {
    Run(args ...string) (string, error)
}

// New HTTP-based implementation
type HTTPMcpmRunner struct {
    BaseURL    string
    HTTPClient *http.Client
}

func (r *HTTPMcpmRunner) Run(args ...string) (string, error) {
    // Map CLI args to HTTP calls
    // Return formatted text output for compatibility
}
```

### Gradual Migration Strategy
1. **Phase 1**: Add HTTPMcpmRunner, keep RealMcpmRunner as fallback
2. **Phase 2**: Switch default to HTTP, CLI as fallback
3. **Phase 3**: Remove CLI fallback once stable

### Configuration
```go
// Environment variable to select transport
// JARVIS_MCPM_TRANSPORT=http (default) | cli
```

## Deployment

### mcpm-daemon Changes
The API server runs as an additional process in the mcpm-daemon container:

```bash
# entrypoint.sh addition
mcpm serve --port 6275 --host 0.0.0.0 &
```

### Port Allocation
- 6275: MCPM API Server (new)
- 6276: p-pokeedge profile
- 6277: memory profile
- 6278: morph profile
- 6279: qdrant profile

### Health Checks
Docker health check updated to verify API server:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:6275/api/v1/health"]
```

## Security Considerations

### Localhost Only
- API server binds to container-internal interface
- No authentication required (same trust model as current CLI)
- Docker network isolation provides security boundary

### Input Validation
- All inputs validated via Pydantic models
- Server names sanitized before use
- Path traversal prevention in file operations

## Testing Strategy

### API Server Tests (Python)
- Unit tests for each endpoint
- Integration tests with real MCPM config
- pytest-asyncio for async endpoint testing

### Jarvis Tests (Go)
- Mock HTTP server for unit tests
- Integration tests against real API server
- Backward compatibility tests with CLI fallback

## Alternatives Considered

### gRPC
- Pro: Strongly typed, efficient binary protocol
- Con: More complex setup, less debugging visibility
- Decision: HTTP/JSON is simpler and sufficient for our needs

### Unix Socket
- Pro: Slightly faster than TCP
- Con: More complex in Docker networking
- Decision: TCP on localhost is simple and fast enough

### Embedded Python (CGO)
- Pro: No network overhead
- Con: Complex build, CGO dependencies
- Decision: HTTP keeps clean separation

## Open Questions

1. **Versioning**: Should we version the API (e.g., `/api/v1/`)?
   - Recommendation: Yes, for future compatibility

2. **Pagination**: Should list endpoints support pagination?
   - Recommendation: No initially, add if server counts grow large

3. **Caching**: Should responses be cached?
   - Recommendation: No, config changes should reflect immediately
