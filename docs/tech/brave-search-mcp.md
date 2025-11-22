# Brave Search MCP Server

**GitHub**: https://github.com/brave/brave-search-mcp-server
**Package**: @brave/brave-search-mcp-server
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides web search capabilities using Brave's Search API. This server enables LLMs to perform web searches and access local business information.

**Key Features:**
- Web search via Brave Search API
- Local business search
- Privacy-focused search results
- Real-time web content access

---

## Installation

### Using npm (via npx)
```bash
npx @brave/brave-search-mcp-server
```

### Direct installation
```bash
npm install -g @brave/brave-search-mcp-server
```

---

## Configuration

**Required Environment Variable:**
- `BRAVE_API_KEY`: Brave Search API key

**Get a Brave API Key:**
1. Visit [Brave Search API](https://api.search.brave.com/)
2. Sign up for an account
3. Generate an API key
4. Set as environment variable: `export BRAVE_API_KEY="your_api_key"`

---

## Available Tools

### Web Search
1. `brave_web_search`
   - Perform web searches using Brave Search
   - Inputs: `query` (string, required), `count` (optional, default: 10)
   - Returns: Search results with titles, URLs, descriptions

### Local Search
2. `brave_local_search`
   - Search for local businesses and places
   - Inputs: `query` (string, required), `location` (optional)
   - Returns: Local business listings with addresses, phone numbers, ratings

---

## Usage Examples

### Web Search
```json
{
  "tool": "brave_web_search",
  "arguments": {
    "query": "Model Context Protocol latest developments",
    "count": 5
  }
}
```

### Local Business Search
```json
{
  "tool": "brave_local_search",
  "arguments": {
    "query": "coffee shops",
    "location": "San Francisco, CA"
  }
}
```

### Combined Search Strategy
```json
// First, get general information
{
  "tool": "brave_web_search",
  "arguments": {
    "query": "best practices for MCP server development"
  }
}

// Then find local resources
{
  "tool": "brave_local_search",
  "arguments": {
    "query": "software development meetups"
  }
}
```

---

## Search Result Format

### Web Search Results
```json
{
  "results": [
    {
      "title": "Model Context Protocol Guide",
      "url": "https://example.com/mcp-guide",
      "description": "Comprehensive guide to MCP server development...",
      "age": "2 days ago"
    }
  ]
}
```

### Local Search Results
```json
{
  "results": [
    {
      "name": "Blue Bottle Coffee",
      "address": "315 Linden St, San Francisco, CA 94102",
      "phone": "(510) 653-3394",
      "rating": 4.5,
      "hours": "Mon-Fri: 7:00 AM - 7:00 PM"
    }
  ]
}
```

---

## Known Limitations

1. **API Rate Limits**: Subject to Brave Search API rate limiting
2. **Result Count**: Limited number of results per query
3. **Geographic Coverage**: Local search may have limited coverage in some regions
4. **Real-time Updates**: Search results may not reflect real-time changes
5. **Privacy Mode**: Some search results may be privacy-filtered

---

## Testing

```bash
# Set environment variable
export BRAVE_API_KEY="your_api_key"

# Test server startup
npx @brave/brave-search-mcp-server

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @brave/brave-search-mcp-server

# Test search functionality (requires MCP client)
```

---

## Use Cases

1. **Research**: Find latest information on technical topics
2. **Local Discovery**: Find businesses, restaurants, services
3. **Fact Checking**: Verify current information
4. **Trend Analysis**: Search for trending topics and news
5. **Competitive Analysis**: Research companies and products

---

## Privacy Features

- **No Tracking**: Brave Search doesn't track users
- **Private Results**: Search queries are not stored or profiled
- **Independent Index**: Uses Brave's own search index
- **Ad-free**: No targeted advertising based on search history

---

## Comparison with Other Search Servers

| Feature | Brave Search | Google Search | Bing Search |
|---------|-------------|---------------|-------------|
| Privacy | ✅ No tracking | ❌ Tracks users | ❌ Tracks users |
| API Access | ✅ Available | ✅ Available | ✅ Available |
| Local Search | ✅ Included | ✅ Available | ✅ Available |
| Rate Limits | Moderate | High | Moderate |
| Cost | Free tier available | Paid | Paid |

---

## Related Documentation

- [Brave Search API Documentation](https://api.search.brave.com/)
- [Brave Search MCP Server](https://github.com/brave/brave-search-mcp-server)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Brave Search Official Site](https://search.brave.com/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
