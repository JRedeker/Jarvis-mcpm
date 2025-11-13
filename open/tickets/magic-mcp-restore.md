Title: Restore magic-mcp MCP server (npm 404)

Status: Open
Priority: Medium
Created: 2025-11-12
Component: cipher-aggregator, magic-mcp
Labels: bug, mcp-integration, npm

Description:
The aggregator fails to start the `magic-mcp` server because the configured npm package cannot be found in the registry (404). This prevents the aggregator from connecting to the server and causes repeated connection failures in the logs.

Relevant logs:
- Repeated npm 404 entries during aggregator startup:
  - "npm error 404 Not Found - GET https://registry.npmjs.org/@21stdev%2fmagic-mcp - Not found"
  - "ERROR: MCP Connection: Failed to connect to MCP server: magic-mcp"

Impact:
- Aggregator attempts to start the server on each startup, producing errors and adding noise to logs.
- Tool availability may be reduced for workflows that depended on magic-mcp.

Suggested remediation steps:
1. Confirm correct package identity
   - Verify intended package name (e.g., the org/name spelling and dash vs dot variations).
   - If package moved or renamed, update cipher.yml with correct package name.

2. If repository is the source (preferred when npm package unavailable)
   - Replace npm install spec with a git dependency, e.g.:
     - args:
       - -y
       - "github:user/magic-mcp"
     - Or use an exact tarball/git URL.

3. Private/internal package options
   - Configure npm auth / registry access for the environment or publish a scoped package to the configured registry.

4. Temporary mitigations (applied)
   - Server was disabled in cipher.yml to stabilize aggregator startup.
   - Keep disabled until a confirmed replacement or package source is provided.

Owner: ops / integrator
Attachments: aggregator log snippets showing npm 404s (search for "magic-mcp" in logs/cipher-aggregator-*.log)

Notes:
- I disabled this server in cipher.yml. Once a source is confirmed, I can re-enable and verify startup.
