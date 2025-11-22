# Fetch MCP Server

**GitHub**: https://github.com/zcaceres/fetch-mcp
**Package**: mcp-server-fetch
**Language**: Python
**Transport**: stdio
**Status**: ✅ Official MCP Reference Server

---

## Overview

A Model Context Protocol server that provides web content fetching capabilities. This server enables LLMs to retrieve and process content from web pages, converting HTML to markdown for easier consumption.

**Key Features:**
- Web content fetching with HTML to markdown conversion
- Configurable content truncation with start_index support
- Raw content option for unprocessed HTML
- Security warnings for local/internal IP access

---

## Installation

### Using uv (recommended)
```bash
# Install uv if not already installed
curl -LsSf https://astral.sh/uv/install.sh | sh

# Run directly with uvx
uvx mcp-server-fetch
```

### Using pip
```bash
pip install mcp-server-fetch
```

### Using npm (via npx)
```bash
npx mcp-server-fetch
```

---

## Configuration

The fetch server requires no environment variables or API keys. It runs with default settings.

**Security Note**: This server can access local/internal IP addresses and may represent a security risk. Exercise caution when using this MCP server to ensure this does not expose any sensitive data.

---

## Available Tools

### `fetch`
Fetches a URL from the internet and extracts its contents as markdown.

**Parameters:**
- `url` (string, required): URL to fetch
- `max_length` (integer, optional): Maximum number of characters to return (default: 5000)
- `start_index` (integer, optional): Start content from this character index (default: 0)
- `raw` (boolean, optional): Get raw content without markdown conversion (default: false)

**Example Usage:**
```json
{
  "tool": "fetch",
  "arguments": {
    "url": "https://example.com",
    "max_length": 3000,
    "raw": false
  }
}
```

---

## Usage Examples

### Basic Web Page Fetching
```bash
# Start the server
mcp-server-fetch

# Test with MCP client
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | mcp-server-fetch
```

### Fetching with Content Truncation
The fetch tool supports chunked reading for large pages:
```json
{
  "tool": "fetch",
  "arguments": {
    "url": "https://long-article.com",
    "max_length": 5000,
    "start_index": 0
  }
}
```

Then continue reading from where you left off:
```json
{
  "tool": "fetch",
  "arguments": {
    "url": "https://long-article.com",
    "max_length": 5000,
    "start_index": 5000
  }
}
```

---

## Known Limitations

1. **Content Truncation**: Default maximum length is 5000 characters
2. **Security Risk**: Can access local/internal IP addresses
3. **HTML Processing**: Conversion to markdown may lose some formatting
4. **No Authentication**: Does not support authenticated requests
5. **Rate Limiting**: No built-in rate limiting for requests

---

## Testing

```bash
# Test server startup
mcp-server-fetch --help

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | mcp-server-fetch

# Test fetch functionality (requires MCP client)
```

---

## Related Documentation

- [Official MCP Servers Repository](https://github.com/modelcontextprotocol/servers)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Fetch Server Source Code](https://github.com/zcaceres/fetch-mcp)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
