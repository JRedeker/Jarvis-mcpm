# Ticket: Fix Cipher SSE Server Connection to routing-metadata MCP

**Status**: Open
**Priority**: High
**Created**: 2025-11-12
**Component**: cipher-aggregator, routing-metadata-mcp
**Labels**: bug, mcp-integration, sse-transport

## Problem Description

The `routing-metadata` MCP server is failing to connect to cipher-aggregator via SSE transport. The server starts successfully when run standalone but cipher cannot establish the stdio connection.

### Symptoms

```
13:34:47 INFO: MCP Manager: Connecting to server: routing-metadata
13:34:47 INFO: MCP Manager: Registered client: routing-metadata
13:34:47 INFO: MCP Connection: Connecting to routing-metadata (stdio)
INFO:cipher_routing_middleware:Enforcement mode: warn
INFO:__main__:Routing Metadata MCP Server starting...
13:34:47 ERROR: MCP Connection: Failed to connect to MCP server: routing-metadata
13:34:47 ERROR: MCP Manager: Failed to connect to server: routing-metadata
```

### Impact

- Routing metadata cannot be retrieved for tool calls
- No routing validation metadata in MCP responses
- Analytics and tracking disabled
- Logging to database works but MCP exposure fails

## Environment

- **Cipher Routing Layer**: cipher-aggregator (SSE)
- **Transport**: SSE (Server-Sent Events)
- **Host**: 127.0.0.1:3020
- **Python Version**: 3.12
- **Server Path**: `/home/jrede/dev/MCP/servers/routing-metadata-mcp.py`

## Configuration

### cipher.yml Entry
```yaml
routing-metadata:
  type: stdio
  command: python3
  args:
    - /home/jrede/dev/MCP/servers/routing-metadata-mcp.py
  enabled: true
  timeout: 10000
  connectionMode: lenient
```

## Root Cause Analysis

### Hypothesis 1: Server Exits Before Connection
The Python asyncio event loop may be exiting before cipher completes the handshake.

**Evidence**:
- Server logs show "Routing Metadata MCP Server starting..."
- No subsequent logs indicating it's processing requests
- Other Python MCP servers (llm-inference, file-batch) connect successfully

**Possible Causes**:
1. `asyncio.run(main())` exits immediately if stdin is empty
2. No keep-alive mechanism in the event loop
3. Stdin not properly connected by cipher

### Hypothesis 2: Stdin/Stdout Buffering
Python's default buffering may prevent proper communication.

**Evidence**:
- `print(json.dumps(response), flush=True)` uses flush
- But stdin reading may still be buffered
- Server uses `sys.stdin.readline()` which may block incorrectly

### Hypothesis 3: Import Error During Startup
The middleware import might fail silently in cipher's environment.

**Evidence**:
- Server starts standalone with imports working
- Middleware initialization logs appear
- But may fail when run as child process

## Debugging Steps

### Step 1: Test Standalone Server Communication

```bash
# Test basic JSON-RPC over stdin/stdout
echo '{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"protocolVersion":"0.1.0"}}' | python3 servers/routing-metadata-mcp.py
```

**Expected Output**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "0.1.0",
    "capabilities": {"tools": {}},
    "serverInfo": {"name": "routing-metadata", "version": "1.0.0"}
  }
}
```

### Step 2: Add Debug Logging

Modify `servers/routing-metadata-mcp.py`:

```python
# Add at the top of main()
logger.info(f"STDIN isatty: {sys.stdin.isatty()}")
logger.info(f"STDOUT isatty: {sys.stdout.isatty()}")
logger.info(f"STDERR isatty: {sys.stderr.isatty()}")
logger.info(f"Environment: {os.environ.get('MCP_SERVER', 'not set')}")

# Add before while loop
logger.info("Entering main event loop, waiting for stdin...")
sys.stderr.flush()
```

### Step 3: Compare with Working MCP Server

Check differences between `routing-metadata-mcp.py` and `llm-inference-mcp.py`:

```bash
# Check if llm-inference uses different async pattern
cat servers/llm-inference-mcp.py | grep -A10 "async def main"
cat servers/routing-metadata-mcp.py | grep -A10 "async def main"
```

### Step 4: Test with MCP SDK

Consider using official MCP Python SDK instead of raw asyncio:

```python
from mcp.server import Server
from mcp.server.stdio import stdio_server

# Official SDK handles stdio lifecycle better
```

### Step 5: Compare Aggregator Health and Routing

```bash
# Check aggregator health endpoint
curl -s http://localhost:3020/health || true

# Inspect aggregator logs
tail -n 200 logs/cipher-aggregator.log || true

# Verify routing-metadata server process is running
ps aux | grep -i "routing-metadata-mcp.py" | grep -v grep || true
```

## Potential Solutions

### Solution 1: Use MCP SDK (Recommended)

Replace custom asyncio implementation with official MCP SDK:

```python
#!/usr/bin/env python3
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent

app = Server("routing-metadata")

@app.list_tools()
async def list_tools() -> list[Tool]:
    return [...]

@app.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    # Implementation
    pass

if __name__ == "__main__":
    import asyncio
    asyncio.run(stdio_server(app))
```

**Pros**:
- Battle-tested by MCP team
- Handles stdio lifecycle correctly
- Better error handling
- More maintainable

**Cons**:
- Requires dependency: `pip install mcp`
- Need to refactor existing code

### Solution 2: Fix Event Loop Lifecycle

Add proper stdin monitoring to keep event loop alive:

```python
async def main():
    logger.info("Routing Metadata MCP Server starting...")

    # Create stdin reader
    reader = asyncio.StreamReader()
    protocol = asyncio.StreamReaderProtocol(reader)
    await loop.connect_read_pipe(lambda: protocol, sys.stdin)

    while True:
        try:
            line = await reader.readline()
            if not line:  # EOF
                logger.info("Received EOF, shutting down")
                break
            # ... process request
        except Exception as e:
            logger.error(f"Error: {e}", exc_info=True)
```

### Solution 3: Add Unbuffered I/O

Force unbuffered I/O at Python level:

```python
# Add shebang with unbuffered flag
#!/usr/bin/env python3 -u

# Or programmatically
sys.stdin.reconfigure(encoding='utf-8', newline='', line_buffering=False)
sys.stdout.reconfigure(encoding='utf-8', newline='', line_buffering=False)
```

### Solution 4: Match llm-inference-mcp Pattern

Copy the exact implementation pattern from `llm-inference-mcp.py` since it connects successfully:

```bash
# Extract the working pattern
cp servers/llm-inference-mcp.py servers/routing-metadata-mcp-v2.py
# Modify to add routing tools while keeping structure identical
```

## Testing Checklist

- [ ] Server responds to echo pipe test
- [ ] Server logs show stdin connected
- [ ] Server doesn't exit immediately
- [ ] Cipher connects without errors
- [ ] Tools appear in cipher's tool list
- [ ] validate_tool_selection can be called
- [ ] Response metadata is properly formatted
- [ ] Errors are logged appropriately
- [ ] Server handles graceful shutdown
- [ ] Works after cipher restart

## Acceptance Criteria

1. ✅ `routing-metadata` appears in cipher's connected servers list
2. ✅ No connection errors in cipher-aggregator.log
3. ✅ Server processes initialize, list_tools, and call_tool requests
4. ✅ Logs show successful validation calls
5. ✅ Metadata appears in tool responses with `_meta.routing` field

## Related Files

- Server: `servers/routing-metadata-mcp.py`
- Config: `cipher.yml`
- Docs: `docs/routing-metadata-usage.md`
- Middleware: `cipher_routing_middleware.py`
- Logs: `logs/routing-metadata.log`, `logs/cipher-aggregator.log`

## References

- MCP Python SDK: https://github.com/modelcontextprotocol/python-sdk
- MCP Spec: https://spec.modelcontextprotocol.io/
- Working example: `servers/llm-inference-mcp.py`
- Aggregator config: [cipher.yml](cipher.yml)
- Routing rules: [.kilocode/rules/cipher-routing-rules.md](.kilocode/rules/cipher-routing-rules.md)

## Next Actions

1. **Immediate**: Run Step 1 debug test to verify basic functionality
2. **Short-term**: Compare with llm-inference-mcp.py implementation
3. **Medium-term**: Consider migrating to official MCP SDK
4. **Long-term**: Document MCP server best practices for cipher integration

## Notes

- This is P0 for routing metadata visibility
- Server functionality is sound (middleware works standalone)
- Issue is specifically with stdio transport lifecycle
- Other stdio servers work, so solution exists in codebase
- May be related to asyncio vs SDK event loop differences
