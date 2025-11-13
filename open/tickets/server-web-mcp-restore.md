# Ticket: Restore server-web-mcp MCP Server

## Issue Summary
The server-web-mcp server fails to connect due to npm package issues. The package `@modelcontextprotocol/server-web` may not be available or have dependency problems.

## Error Details
Based on cipher.yml configuration, this server is currently disabled due to npm package availability issues.

## Current Status
- **Server**: server-web-mcp
- **Status**: Disabled in cipher.yml
- **Command**: `npx -y @modelcontextprotocol/server-web`
- **Timeout**: 60000ms
- **Connection Mode**: lenient

## Root Cause
The npm package `@modelcontextprotocol/server-web` may have:
1. Been removed from the npm registry
2. Dependency conflicts with other packages
3. Version compatibility issues
4. Platform-specific build problems

## Proposed Solutions

### Option 1: Verify Package Availability
- Check if `@modelcontextprotocol/server-web` exists in npm registry
- Verify package version and dependencies
- Test manual installation to identify specific issues

### Option 2: Alternative Web Server Package
- Look for alternative web server MCP packages
- Consider `@modelcontextprotocol/server-http` or similar
- Evaluate third-party web server MCP alternatives

### Option 3: Custom Web Server Implementation
- Create a custom web server using Python/Node.js
- Implement basic HTTP server functionality
- Add web-based tools and utilities

### Option 4: Use Existing Server Functionality
- Integrate web server features into existing servers
- Add web-related tools to routing-metadata or other servers

## Investigation Steps
1. **Package Verification**
   ```bash
   npm view @modelcontextprotocol/server-web
   npm info @modelcontextprotocol/server-web
   ```

2. **Manual Installation Test**
   ```bash
   npx -y @modelcontextprotocol/server-web
   ```

3. **Dependency Analysis**
   - Check package.json for dependency conflicts
   - Review npm install logs for specific errors
   - Test with different Node.js versions

4. **Alternative Package Search**
   ```bash
   npm search @modelcontextprotocol server
   npm search web mcp server
   ```

## Expected Functionality
- Web server management tools
- HTTP request/response handling
- Web-based utilities and operations
- Server status and monitoring

## Testing Commands
```bash
# Test package availability
npm view @modelcontextprotocol/server-web

# Test installation
npx -y @modelcontextprotocol/server-web

# Check for alternatives
npm search @modelcontextprotocol server
npm search web mcp
```

## Impact
- **Severity**: Medium
- **Affected Users**: Users needing web server management tools
- **Workaround**: Use manual web server operations or alternative tools

## Dependencies
- npm package registry
- Node.js runtime
- Package dependencies (to be analyzed)

## Next Steps
1. Verify package availability and specific error details
2. Test manual installation to identify root cause
3. Research alternative packages if needed
4. Implement custom solution if no packages are available
5. Update cipher.yml configuration
6. Test connection and basic functionality

---
*Created: 2025-11-12 16:55*
*Priority: Medium*
