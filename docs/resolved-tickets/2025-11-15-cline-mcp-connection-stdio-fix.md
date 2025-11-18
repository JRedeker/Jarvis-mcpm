# Cline MCP Connection Fix: Streamable-HTTP to stdio Migration

**Date**: 2025-11-15
**Status**: ✅ Resolved
**Component**: Cline MCP Integration
**Severity**: High (blocked all Cline → cipher-aggregator communication)

## Summary

Fixed Cline's inability to connect to cipher-aggregator by migrating from incompatible streamable-HTTP protocol to stdio-based communication and properly configuring environment variables.

## Problem Statement

### Symptoms
1. **Schema Validation Error**: Cline reported "Invalid MCP settings schema" with the initial configuration
2. **SSE 400 Error**: After schema fix, connection attempts returned `SSE error: Non-200 status code (400)`
3. **Missing API Keys**: After protocol fix, cipher failed with "No API key or Ollama configuration found"

### Error Logs
```
2025-11-15 17:52:49 [ERROR]: [MCP Streamable-HTTP Server] Invalid or missing session ID for SSE stream
```

```
[CIPHER-MCP] ERROR: No API key or Ollama configuration found
Available providers: OpenAI, Anthropic, OpenRouter, Ollama, Qwen, Gemini, Azure, AWS, LM Studio
MCP error -32000: Connection closed
```

## Root Cause Analysis

### Issue 1: Protocol Mismatch

**Original Configuration**:
```json
{
  "mcpServers": {
    "cipher-aggregator": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:3020/http",
      "timeout": 60
    }
  }
}
```

**Problem**: Cipher's streamable-HTTP protocol requires a two-step handshake:
1. `POST /http` → receive session ID
2. `GET /http?sessionId=xxx` → establish SSE stream

**Cline's MCP Client** only supports:
- Simple SSE (single GET establishes stream)
- stdio (stdin/stdout pipe communication)

Cline was attempting direct GET requests without session initialization, causing 400 errors.

### Issue 2: Unsupported Schema Fields

Cline's MCP settings schema doesn't recognize:
- `type` field (infers from `url` vs `command` presence)
- `timeout` field (not part of spec)

### Issue 3: Environment Variable Isolation

When Cline spawns subprocesses via `command` + `args`, they run in **isolated environments** that don't inherit:
- Shell environment variables
- Project `.env` file contents
- Parent process environment

The `.env` file at `/home/jrede/dev/MCP/.env` is only loaded by:
- Interactive shells
- `mcp-manager.sh` (explicitly sources it)
- **NOT** by Cline's child processes

## Solution

### Final Working Configuration

**File**: `~/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`

```json
{
  "mcpServers": {
    "cipher-aggregator": {
      "command": "node",
      "args": [
        "/home/jrede/.npm-global/bin/cipher",
        "--mode", "mcp",
        "--agent", "/home/jrede/dev/MCP/cipher.yml",
        "--mcp-transport-type", "stdio"
      ],
      "env": {
        "MCP_SERVER_MODE": "aggregator",
        "OPENROUTER_API_KEY": "sk-or-v1-dummy123",
        "OPENAI_API_KEY": "sk-dummy-openai-key",
        "BRAVE_API_KEY": "dummy-brave-key",
        "TAVILY_API_KEY": "dummy-tavily-key",
        "AGGREGATOR_CONFLICT_RESOLUTION": "prefix"
      },
      "disabled": false,
      "autoApprove": [
        "cipher_bash",
        "cipher_extract_and_operate_memory",
        "cipher_memory_search",
        "brave_web_search",
        ...
      ]
    }
  }
}
```

### Key Changes

1. **Protocol Migration**: `url` → `command` + `args` (HTTP → stdio)
2. **Removed Unsupported Fields**: Deleted `type` and `timeout`
3. **Explicit Environment Variables**: Added all API keys from `.env` to `env` object
4. **stdio Transport**: Added `--mcp-transport-type stdio` argument

## Technical Details

### stdio vs streamable-HTTP

| Feature | stdio | streamable-HTTP |
|---------|-------|-----------------|
| Communication | stdin/stdout pipes | HTTP GET/POST + SSE |
| Session Management | N/A (direct pipe) | Explicit session IDs |
| Cline Support | ✅ Full | ❌ Not supported |
| Use Case | IDE integrations | Web clients, multi-connection |
| Process Model | 1:1 (dedicated instance) | 1:N (shared instance) |

### Environment Variable Inheritance

```
Cline Extension Process
  └─> spawns subprocess via command + args
       └─> NEW isolated environment
            └─> ONLY gets variables from "env" object
            └─> DOES NOT inherit shell env
            └─> DOES NOT read .env files
```

### Why It Works Now

1. **Direct Communication**: Cline writes JSON-RPC to cipher's stdin, reads responses from stdout
2. **No Session Handshake**: stdio protocol doesn't require session initialization
3. **Explicit Credentials**: All API keys passed directly to subprocess environment
4. **Dedicated Instance**: Each Cline creates its own cipher instance (no conflicts)

## Verification

```bash
# Verify cipher can start with stdio
node /home/jrede/.npm-global/bin/cipher \
  --mode mcp \
  --agent /home/jrede/dev/MCP/cipher.yml \
  --mcp-transport-type stdio

# Test with MCP client
echo '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | \
  node /home/jrede/.npm-global/bin/cipher --mode mcp --mcp-transport-type stdio
```

## Related Files

- [`cline_mcp_settings.json`](../../.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json) - Cline MCP configuration
- [`cipher.yml`](../../cipher.yml) - Cipher agent configuration
- [`.env`](../../.env) - Environment variables (not auto-loaded by Cline)
- [`docs/tech/cipher-aggregator.md`](../tech/cipher-aggregator.md) - Aggregator architecture

## Lessons Learned

1. **Protocol Compatibility**: Always verify client supports server's transport protocol
2. **Environment Isolation**: Subprocess environment must be explicitly configured
3. **Schema Validation**: Use only documented fields for third-party tools
4. **stdio for IDEs**: IDE integrations benefit from stdio over HTTP-based protocols

## Future Considerations

- Consider adding environment variable loading to cipher startup
- Document supported MCP transport types for different clients
- Create wrapper script that auto-loads `.env` for stdio clients