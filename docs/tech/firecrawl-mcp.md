# Firecrawl MCP Server

**GitHub**: https://github.com/firecrawl/firecrawl-mcp-server

---

## Overview

A Model Context Protocol server that provides advanced web scraping and content extraction capabilities using Firecrawl. This server enables LLMs to extract structured data from websites, crawl multiple pages, and convert web content to various formats.

**Key Features:**
- Advanced web scraping with JavaScript rendering
- Multi-page crawling and sitemap generation
- Content extraction to markdown, HTML, or text
- Structured data extraction
- Proxy support and rate limiting
- JavaScript execution and screenshot capture

---

## Installation

### Using npm (via npx)
```bash
npx @firecrawl/mcp-server
```

### Direct installation
```bash
npm install -g @firecrawl/mcp-server
```

---

## Configuration

**Required Environment Variable:**
- `FIRECRAWL_API_KEY`: Firecrawl API key

**Get a Firecrawl API Key:**
1. Visit [Firecrawl](https://www.firecrawl.dev/)
2. Sign up for an account
3. Generate an API key from the dashboard
4. Set as environment variable: `export FIRECRAWL_API_KEY="your_api_key"`

**Optional Configuration:**
- `FIRECRAWL_BASE_URL`: Custom base URL (default: https://api.firecrawl.dev)
- `FIRECRAWL_TIMEOUT`: Request timeout in seconds (default: 30)

---

## Available Tools

### Content Scraping
1. `scrape_url`
   - Scrape a single URL with various options
   - Inputs: `url` (required), `formats` (optional), `include_tags` (optional), `exclude_tags` (optional), `timeout` (optional)
   - Returns: Scraped content in requested formats

### Web Crawling
2. `crawl_url`
   - Crawl multiple pages from a starting URL
   - Inputs: `url` (required), `max_depth` (optional), `max_pages` (optional), `formats` (optional)
   - Returns: Crawled pages content

### Site Mapping
3. `map_url`
   - Generate a sitemap of all accessible pages
   - Inputs: `url` (required), `max_depth` (optional), `include_paths` (optional), `exclude_paths` (optional)
   - Returns: List of discovered URLs

### Advanced Features
4. `scrape_with_js`
   - Scrape with custom JavaScript execution
   - Inputs: `url` (required), `js_code` (required), `wait_for` (optional), `formats` (optional)
   - Returns: Page content plus JavaScript execution results

5. `take_screenshot`
   - Capture screenshots of web pages
   - Inputs: `url` (required), `width` (optional), `height` (optional), `full_page` (optional)
   - Returns: Screenshot as base64 encoded image

---

## Usage Examples

### Basic Scraping
```json
{
  "tool": "scrape_url",
  "arguments": {
    "url": "https://example.com",
    "formats": ["markdown", "html"]
  }
}
```

### Web Crawling
```json
{
  "tool": "crawl_url",
  "arguments": {
    "url": "https://docs.example.com",
    "max_depth": 2,
    "max_pages": 10,
    "formats": ["markdown"]
  }
}
```

### Site Mapping
```json
{
  "tool": "map_url",
  "arguments": {
    "url": "https://example.com",
    "max_depth": 3,
    "include_paths": ["/docs/*", "/api/*"]
  }
}
```

### JavaScript Execution
```json
{
  "tool": "scrape_with_js",
  "arguments": {
    "url": "https://example.com",
    "js_code": "document.querySelector('h1').textContent",
    "wait_for": "h1"
  }
}
```

### Screenshot Capture
```json
{
  "tool": "take_screenshot",
  "arguments": {
    "url": "https://example.com",
    "width": 1200,
    "height": 800,
    "full_page": true
  }
}
```

---

## Content Formats

### Available Formats
- `markdown`: Clean markdown conversion (default)
- `html`: Raw HTML content
- `text`: Plain text extraction
- `screenshot`: Visual capture
- `links`: Extract all links
- `images`: Extract image URLs

### Format Combinations
```json
{
  "formats": ["markdown", "links", "images"]
}
```

---

## Advanced Options

### Content Filtering
```json
{
  "include_tags": ["article", "main", "section"],
  "exclude_tags": ["nav", "footer", "aside"],
  "timeout": 45
}
```

### Crawling Configuration
```json
{
  "max_depth": 3,
  "max_pages": 50,
  "include_paths": ["/blog/*", "/docs/*"],
  "exclude_paths": ["/admin/*", "/private/*"]
}
```

### Proxy and Headers
```json
{
  "headers": {
    "User-Agent": "MCP-Firecrawl-Bot/1.0",
    "Accept-Language": "en-US"
  },
  "proxy": "http://proxy.example.com:8080"
}
```

---

## Known Limitations

1. **API Rate Limits**: Subject to Firecrawl API rate limiting
2. **JavaScript Required**: Some sites require JS execution for full content
3. **Anti-bot Measures**: May be blocked by sophisticated anti-bot systems
4. **Dynamic Content**: Single-page applications may need special handling
5. **Large Sites**: Crawling very large sites may hit limits
6. **Authentication**: No built-in support for authenticated content

---

## Testing

```bash
# Set environment variable
export FIRECRAWL_API_KEY="your_api_key"

# Test server startup
npx @firecrawl/mcp-server

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @firecrawl/mcp-server

# Test scraping functionality (requires MCP client)
```

---

## Use Cases

1. **Documentation Scraping**: Extract technical documentation
2. **Content Aggregation**: Collect articles and blog posts
3. **Data Extraction**: Extract structured data from websites
4. **Competitive Analysis**: Monitor competitor websites
5. **Research**: Gather information from multiple sources
6. **Testing**: Capture screenshots for visual testing

---

## Performance Tips

### Efficient Scraping
- Use specific `include_tags` and `exclude_tags`
- Set appropriate timeouts
- Request only needed formats
- Use `max_pages` limits for crawling

### Rate Limiting
- Respect rate limits
- Use delays between requests
- Monitor API usage
- Implement retry logic

---

## Comparison with Other Scraping Tools

| Feature              | Firecrawl | BeautifulSoup | Puppeteer | Playwright |
| -------------------- | --------- | ------------- | --------- | ---------- |
| JavaScript Rendering | ✅         | ❌             | ✅         | ✅          |
| API Interface        | ✅         | ❌             | ❌         | ❌          |
| Multi-format Output  | ✅         | ❌             | ❌         | ❌          |
| Built-in Crawling    | ✅         | ❌             | ❌         | ❌          |
| Rate Limiting        | ✅         | Manual        | Manual    | Manual     |
| Screenshot Capture   | ✅         | ❌             | ✅         | ✅          |

---

## Related Documentation

- [Firecrawl Documentation](https://docs.firecrawl.dev/)
- [Firecrawl MCP Server Documentation](https://docs.firecrawl.dev/mcp-server)
- [Firecrawl MCP Server](https://github.com/firecrawl/firecrawl-mcp-server)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Firecrawl API Reference](https://docs.firecrawl.dev/api-reference)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
