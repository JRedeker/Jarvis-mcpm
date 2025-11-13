# Ticket: Restore time-mcp MCP Server

## Issue Summary
The time-mcp server fails to connect due to npm 404 error - the package `@modelcontextprotocol/server-time` is not found in the npm registry.

## Error Details
```
npm error code E404
npm error 404 Not Found - GET https://registry.npmjs.org/@modelcontextprotocol%2fserver-time - Not found
npm error 404
npm error 404  '@modelcontextprotocol/server-time@*' is not in this registry.
```

## Current Status
- **Server**: time-mcp
- **Status**: Disabled in cipher.yml
- **Command**: `npx -y @modelcontextprotocol/server-time`
- **Timeout**: 60000ms
- **Connection Mode**: lenient

## Root Cause
The npm package `@modelcontextprotocol/server-time` is not available in the npm registry. This package may have been:
1. Moved to a different repository
2. Renamed
3. Deprecated
4. Never officially published

## Proposed Solutions

### Option 1: Find Alternative Package
- Search for alternative time/date MCP servers
- Look for `@modelcontextprotocol/server-datetime` or similar packages
- Check if functionality is available in other servers

### Option 2: Implement Custom Time Server
- Create a custom time server script
- Use existing Python/Node.js time libraries
- Integrate with the MCP protocol

### Option 3: Use Built-in Time Functionality
- Implement time functions in existing servers
- Add time-based tools to routing-metadata or other servers

## Next Steps
1. Research available time/date MCP servers in npm registry
2. Test alternative packages if found
3. If no alternatives exist, implement custom solution
4. Update cipher.yml configuration
5. Test connection and basic functionality

## Testing Commands
```bash
# Test npm package availability
npm view @modelcontextprotocol/server-time

# Check for alternatives
npm search @modelcontextprotocol server
npm search datetime mcp

# Test manual installation (once package is found)
npx -y @modelcontextprotocol/server-time
```

## Impact
- **Severity**: Low
- **Affected Users**: Users needing current time/date functionality
- **Workaround**: Use system time or manual date commands

## Dependencies
- npm package registry
- Node.js runtime

---
*Created: 2025-11-12 16:55*
*Priority: Medium*
