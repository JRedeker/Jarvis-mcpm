## Phase 1: MCPM API Server (Node.js)

Note: Implementation uses Node.js/Express instead of Python/FastAPI since the existing MCPM CLI is Node.js.

- [x] 1.1 Create `MCPM/api/` module structure with `server.js`, `routes/`, `helpers.js` <!-- id: 0 -->
- [x] 1.2 Add Express dependency to `package.json` <!-- id: 1 -->
- [x] 1.3 Implement `/health` endpoint (wraps `mcpm doctor` logic) <!-- id: 2 -->
- [x] 1.4 Implement server endpoints: `GET/POST/DELETE /servers`, `GET /servers/{name}` <!-- id: 3 -->
- [x] 1.5 Implement search endpoint: `GET /search?q={query}` <!-- id: 4 -->
- [x] 1.6 Implement profile endpoints: `GET/POST/PUT/DELETE /profiles` <!-- id: 5 -->
- [x] 1.7 Implement client endpoints: `GET/PUT /clients` <!-- id: 6 -->
- [x] 1.8 Implement system endpoints: `GET /usage`, `POST /migrate` <!-- id: 7 -->
- [x] 1.9 Add `mcpm serve` CLI command to start API server <!-- id: 8 -->
- [ ] 1.10 Write unit tests for all API endpoints <!-- id: 9 -->

## Phase 2: Jarvis HTTP Client (Go)

- [x] 2.1 Define `HTTPMcpmRunner` struct with BaseURL and http.Client <!-- id: 10 -->
- [x] 2.2 Implement HTTP client methods mapping to API endpoints <!-- id: 11 -->
- [x] 2.3 Add response parsing to extract structured data from JSON <!-- id: 12 -->
- [x] 2.4 Format responses as text output for backward compatibility <!-- id: 13 -->
- [x] 2.5 Add environment variable `JARVIS_MCPM_TRANSPORT` for transport selection <!-- id: 14 -->
- [x] 2.6 Update `NewHandler()` to select transport based on config <!-- id: 15 -->
- [ ] 2.7 Write unit tests for HTTPMcpmRunner with mock HTTP server <!-- id: 16 -->
- [ ] 2.8 Update integration tests to work with both transports <!-- id: 17 -->

## Phase 3: Deployment Integration

- [x] 3.1 Update `mcpm-daemon/entrypoint.sh` to start API server on port 6275 <!-- id: 18 -->
- [x] 3.2 Update `docker-compose.yml` to expose port 6275 <!-- id: 19 -->
- [x] 3.3 Add health check for API server to Docker config <!-- id: 20 -->
- [x] 3.4 Update `check_status` tool to verify API server health <!-- id: 21 -->

## Phase 4: Documentation & Cleanup

- [x] 4.1 Document API endpoints in `docs/API_REFERENCE.md` <!-- id: 22 -->
- [x] 4.2 Update `AGENTS.md` with new architecture <!-- id: 23 -->
- [x] 4.3 Add troubleshooting guide for API server issues <!-- id: 24 -->
- [ ] 4.4 Update project.md to reflect new architecture <!-- id: 25 -->

## Summary

**Completed:** 21/25 tasks (84%)

**Remaining (low priority):**
- Unit tests for API endpoints (1.10)
- Unit tests for HTTPMcpmRunner (2.7)
- Integration test updates (2.8)
- Project.md update (4.4)

These remaining tasks are not blocking the implementation and can be done in follow-up work.
