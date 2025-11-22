# Playwright MCP Server

**GitHub**: https://github.com/microsoft/playwright-mcp
**Package**: @playwright/mcp
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Official Microsoft Server

---

## Overview

A Model Context Protocol server that provides browser automation capabilities using Microsoft Playwright. This server enables LLMs to interact with web pages, take screenshots, execute JavaScript, and perform comprehensive browser automation tasks.

**Key Features:**
- Full browser automation with Playwright
- Screenshot capture (full page, viewport, or specific elements)
- JavaScript execution in browser context
- Form interaction and navigation
- Multi-browser support (Chromium, Firefox, WebKit)
- Mobile device emulation
- Network interception and monitoring
- PDF generation from web pages

---

## Installation

### Using npm (via npx)
```bash
npx @playwright/mcp
```

### Direct installation
```bash
npm install -g @playwright/mcp
```

### Install Playwright browsers
```bash
# Install all supported browsers
npx playwright install

# Install specific browsers
npx playwright install chromium
npx playwright install firefox
npx playwright install webkit
```

---

## Configuration

The Playwright server requires no environment variables or API keys. It runs with default settings.

**Browser Requirements**: Playwright browsers will be automatically downloaded when first used.

**Security Note**: This server can access local files and local/internal IP addresses since it runs a browser on your machine. Exercise caution when using this MCP server to ensure this does not expose any sensitive data.

---

## Available Tools

### Navigation and Page Management
1. `playwright_navigate`
   - Navigate to any URL in the browser
   - Inputs: url (required), wait_until (optional), timeout (optional)

2. `playwright_screenshot`
   - Capture screenshots of pages or elements
   - Inputs: name (required), selector (optional), full_page (optional), viewport (optional)

3. `playwright_get_content`
   - Get page content as text or HTML
   - Inputs: format (optional: 'text' or 'html')

### Element Interaction
4. `playwright_click`
   - Click on elements using CSS selectors
   - Inputs: selector (required), timeout (optional)

5. `playwright_fill`
   - Fill form fields with text
   - Inputs: selector (required), text (required), clear (optional)

6. `playwright_select`
   - Select options from dropdowns
   - Inputs: selector (required), value (required)

7. `playwright_hover`
   - Hover over elements
   - Inputs: selector (required), timeout (optional)

### JavaScript Execution
8. `playwright_evaluate`
   - Execute JavaScript in the browser context
   - Inputs: script (required), args (optional)

### Browser Context Management
9. `playwright_new_page`
   - Create a new browser page/tab
   - Inputs: None

10. `playwright_close_page`
    - Close the current browser page
    - Inputs: None

### Mobile and Device Emulation
11. `playwright_emulate_device`
    - Emulate a specific mobile device
    - Inputs: device (required) - e.g., 'iPhone 13', 'Pixel 5'

### Network and Performance
12. `playwright_wait_for_selector`
    - Wait for element to appear
    - Inputs: selector (required), timeout (optional)

13. `playwright_pdf`
    - Generate PDF from current page
    - Inputs: path (required), format (optional), margin (optional)

---

## Usage Examples

### Basic Navigation and Screenshot
```json
{
  "tool": "playwright_navigate",
  "arguments": {
    "url": "https://example.com",
    "wait_until": "networkidle"
  }
}
```

```json
{
  "tool": "playwright_screenshot",
  "arguments": {
    "name": "example-homepage",
    "full_page": true
  }
}
```

### Form Interaction
```json
{
  "tool": "playwright_fill",
  "arguments": {
    "selector": "#search-input",
    "text": "Playwright MCP",
    "clear": true
  }
}
```

```json
{
  "tool": "playwright_click",
  "arguments": {
    "selector": "#search-button"
  }
}
```

### JavaScript Execution
```json
{
  "tool": "playwright_evaluate",
  "arguments": {
    "script": "document.title",
    "args": []
  }
}
```

### Mobile Device Emulation
```json
{
  "tool": "playwright_emulate_device",
  "arguments": {
    "device": "iPhone 13"
  }
}
```

### Element Screenshot
```json
{
  "tool": "playwright_screenshot",
  "arguments": {
    "name": "header-section",
    "selector": "header",
    "viewport": {"width": 1200, "height": 800}
  }
}
```

### PDF Generation
```json
{
  "tool": "playwright_pdf",
  "arguments": {
    "path": "page.pdf",
    "format": "A4",
    "margin": {"top": "1cm", "bottom": "1cm"}
  }
}
```

---

## Supported Devices for Emulation

### Mobile Devices
- iPhone SE, iPhone 12, iPhone 13, iPhone 14
- Samsung Galaxy S9, S20, S21
- Google Pixel 5, Pixel 7
- iPad, iPad Pro, iPad Mini

### Desktop Viewports
- Desktop Chrome, Firefox, Safari
- Custom viewport sizes

---

## Browser Launch Options

The server supports various browser configurations:

**Supported Browsers:**
- **chromium** (default) - Google Chrome/Chromium
- **firefox** - Mozilla Firefox
- **webkit** - Apple Safari/WebKit

**Common Options:**
- `headless`: Run browser in headless mode (default: true)
- `viewport`: Set viewport dimensions
- `user_agent`: Custom user agent string
- `timeout`: Page load timeout

---

## Known Limitations

1. **Resource Intensive**: Browser automation consumes significant memory/CPU
2. **Single Context**: Limited browser context management
3. **Network Dependency**: Requires internet access for most operations
4. **File Uploads**: Limited support for file upload scenarios
5. **Authentication**: No built-in authentication handling

---

## Testing

```bash
# Test server startup
npx @playwright/mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @playwright/mcp

# Test browser automation (requires MCP client)
```

---

## Use Cases

1. **Web Testing**: Automated testing of web applications
2. **Screenshot Testing**: Visual regression testing
3. **Form Automation**: Complex form interactions
4. **Content Extraction**: Data extraction from web pages
5. **PDF Generation**: Convert web pages to PDF documents
6. **Mobile Testing**: Test responsive designs

---

## Security Considerations

- **Local File Access**: Can access local files through file:// URLs
- **Internal Network**: Can access internal IP addresses
- **Code Execution**: JavaScript execution in browser context
- **Resource Usage**: Monitor system resources during automation

**Best Practices:**
- Run in isolated environment when possible
- Limit network access in production
- Monitor resource usage
- Use headless mode for automation

---

## Related Documentation

- [Playwright Documentation](https://playwright.dev/)
- [Playwright MCP Server](https://github.com/microsoft/playwright-mcp)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
