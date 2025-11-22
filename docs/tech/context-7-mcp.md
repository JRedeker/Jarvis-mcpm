
**github**: https://github.com/upstash/context7

---

# Context7 MCP Server Documentation

## Introduction

Context7 MCP is a Model Context Protocol (MCP) server that provides AI coding assistants with up-to-date, version-specific documentation and code examples for any library or framework. Instead of relying on outdated training data or hallucinated APIs, Context7 fetches real-time documentation directly from the source, allowing LLMs to generate accurate, working code. The server operates in two modes: as a local stdio server integrated directly into MCP clients like Cursor, Claude Code, and VS Code, or as a remote HTTP server accessible via REST API endpoints.

The project implements a TypeScript-based MCP server that exposes two primary tools: `resolve-library-id` for searching and identifying libraries, and `get-library-docs` for retrieving documentation content. It supports flexible authentication via API keys, proxy configurations for enterprise environments, and intelligent rate limiting. The server handles search queries, resolves library identifiers, and returns formatted documentation optimized for LLM consumption, complete with code snippets, source reputation scores, benchmark scores, and version information.

## API Reference

### MCP Tool: resolve-library-id

Resolves a package or library name into a Context7-compatible library ID. This tool must be called before fetching documentation unless the user provides an explicit library ID in the format `/org/project` or `/org/project/version`.

```typescript
// Example usage within an MCP client (e.g., Cursor, Claude Code)
// User prompt: "How do I use Next.js routing? use context7"

// The LLM automatically calls this tool:
{
  "tool": "resolve-library-id",
  "arguments": {
    "libraryName": "next.js"
  }
}

// Response format:
{
  "content": [
    {
      "type": "text",
      "text": `Available Libraries (top matches):

Each result includes:
- Library ID: Context7-compatible identifier (format: /org/project)
- Name: Library or package name
- Description: Short summary
- Code Snippets: Number of available code examples
- Source Reputation: Authority indicator (High, Medium, Low, or Unknown)
- Benchmark Score: Quality indicator (100 is the highest score)
- Versions: List of versions if available. Use one of those versions if the user provides a version in their query. The format of the version is /org/project/version.

For best results, select libraries based on name match, source reputation, snippet coverage, benchmark score, and relevance to your use case.

----------

- Title: Next.js
- Context7-compatible library ID: /vercel/next.js
- Description: The React Framework for Production
- Code Snippets: 1247
- Source Reputation: High
- Benchmark Score: 95
- Versions: v14.3.0-canary.87, v15.0.0, v13.5.6

----------

- Title: Next.js Documentation
- Context7-compatible library ID: /vercel/nextjs.org
- Description: Official Next.js documentation site
- Code Snippets: 892
- Source Reputation: High
- Benchmark Score: 92`
    }
  ]
}

// Error handling example (rate limited):
{
  "content": [
    {
      "type": "text",
      "text": "Rate limited due to too many requests. You can create a free API key at https://context7.com/dashboard for higher rate limits."
    }
  ]
}
```

### MCP Tool: get-library-docs

Fetches up-to-date documentation for a specific library using its Context7-compatible library ID. Requires the exact ID obtained from `resolve-library-id` unless provided directly by the user. Uses pagination for efficient documentation retrieval.

```typescript
// Example usage within an MCP client
// After resolving library ID, fetch documentation:
{
  "tool": "get-library-docs",
  "arguments": {
    "context7CompatibleLibraryID": "/vercel/next.js",
    "topic": "routing",
    "page": 1
  }
}

// Response contains formatted documentation with code examples:
{
  "content": [
    {
      "type": "text",
      "text": `# Next.js Routing

## App Router (Recommended)

The App Router uses the \`app/\` directory and supports React Server Components...

### Example: Creating a Page

\`\`\`typescript
// app/dashboard/page.tsx
export default function DashboardPage() {
  return <h1>Dashboard</h1>
}
\`\`\`

### Dynamic Routes

Create dynamic route segments using [brackets]...

\`\`\`typescript
// app/blog/[slug]/page.tsx
export default function BlogPost({ params }: { params: { slug: string } }) {
  return <article>{params.slug}</article>
}
\`\`\`
...`
    }
  ]
}

// Pagination example (requesting additional context):
{
  "tool": "get-library-docs",
  "arguments": {
    "context7CompatibleLibraryID": "/mongodb/docs",
    "topic": "aggregation",
    "page": 2  // Fetch next page if more context is needed
  }
}

// Error handling example (invalid library ID):
{
  "content": [
    {
      "type": "text",
      "text": "Documentation not found or not finalized for this library. This might have happened because you used an invalid Context7-compatible library ID. To get a valid Context7-compatible library ID, use the 'resolve-library-id' with the package name you wish to retrieve documentation for."
    }
  ]
}
```

### REST API: Search Libraries

Direct HTTP endpoint for searching libraries when running Context7 MCP in HTTP mode.

```bash
# Basic search without authentication
curl "https://context7.com/api/v1/search?query=react"

# Search with API key authentication
curl -H "Authorization: Bearer ctx7sk_your_api_key" \
  "https://context7.com/api/v1/search?query=supabase"

# Response format:
{
  "results": [
    {
      "id": "/supabase/supabase",
      "title": "Supabase",
      "description": "Open source Firebase alternative",
      "branch": "main",
      "lastUpdateDate": "2025-01-15T10:30:00Z",
      "state": "finalized",
      "totalTokens": 45000,
      "totalSnippets": 320,
      "totalPages": 89,
      "trustScore": 9,
      "benchmarkScore": 88,
      "versions": ["v2.0.0", "v1.8.1"]
    }
  ]
}

# Error response (rate limited):
{
  "results": [],
  "error": "Rate limited due to too many requests. You can create a free API key at https://context7.com/dashboard for higher rate limits."
}

# Error response (unauthorized):
{
  "results": [],
  "error": "Unauthorized. Please check your API key. The API key you provided (possibly incorrect) is: ctx7sk_invalid. API keys should start with 'ctx7sk'"
}
```

### REST API: Fetch Library Documentation

Direct HTTP endpoint for retrieving library documentation when running Context7 MCP in HTTP mode. Now uses v2 API with pagination support.

```bash
# Fetch documentation with default settings
curl "https://context7.com/api/v2/docs/code/mongodb/docs?type=txt"

# Fetch with topic focus and pagination
curl -H "Authorization: Bearer ctx7sk_your_api_key" \
  "https://context7.com/api/v2/docs/code/vercel/next.js?page=1&limit=10&topic=routing&type=txt"

# Fetch specific version
curl "https://context7.com/api/v2/docs/code/vercel/next.js/v14.3.0-canary.87?page=1&limit=10&type=txt"

# Fetch second page for additional context
curl "https://context7.com/api/v2/docs/code/vercel/next.js?page=2&limit=10&topic=routing&type=txt"

# Response: Plain text documentation with code examples
# Example output:
# Next.js Documentation
#
# ## Getting Started
#
# Create a new Next.js app:
# ```bash
# npx create-next-app@latest
# ```
# ...

# Error response (not found):
curl -w "\n%{http_code}\n" \
  "https://context7.com/api/v2/docs/code/invalid/library?type=txt"
# Output:
# The library you are trying to access does not exist. Please try with a different library ID.
# 404

# Using proxy for corporate environments
export HTTPS_PROXY=http://proxy.corp.com:8080
curl "https://context7.com/api/v2/docs/code/supabase/supabase?type=txt"
```

### CLI: Running Context7 MCP Server

Command-line interface for starting the Context7 MCP server in different modes.

```bash
# Install globally
npm install -g @upstash/context7-mcp

# Run in stdio mode (default, for MCP client integration)
context7-mcp --transport stdio --api-key ctx7sk_your_api_key

# Alternative: Use via npx without installation
npx -y @upstash/context7-mcp --transport stdio

# Run in HTTP mode (for remote MCP server)
context7-mcp --transport http --port 3000

# HTTP mode on custom port
context7-mcp --transport http --port 8080

# Using environment variable for API key
export CONTEXT7_API_KEY=ctx7sk_your_api_key
context7-mcp --transport stdio

# Development: Run from source with Bun
git clone https://github.com/upstash/context7.git
cd context7
bun install
bun run build
bun run dist/index.js --transport http --port 3000

# Testing with MCP Inspector
npx -y @modelcontextprotocol/inspector npx @upstash/context7-mcp

# Common error handling:
# Invalid transport
context7-mcp --transport invalid
# Output: Invalid --transport value: 'invalid'. Must be one of: stdio, http.
# Exit code: 1

# Incompatible flags (HTTP + API key)
context7-mcp --transport http --api-key ctx7sk_key
# Output: The --api-key flag is not allowed when using --transport http.
# Exit code: 1

# Incompatible flags (stdio + port)
context7-mcp --transport stdio --port 8080
# Output: The --port flag is not allowed when using --transport stdio.
# Exit code: 1
```

### Configuration: MCP Client Setup

Configuration examples for integrating Context7 MCP into various AI coding assistants.

```json
// Cursor configuration (~/.cursor/mcp.json)
// Remote HTTP server connection
{
  "mcpServers": {
    "context7": {
      "url": "https://mcp.context7.com/mcp",
      "headers": {
        "CONTEXT7_API_KEY": "ctx7sk_your_api_key"
      }
    }
  }
}

// Cursor configuration - Local stdio connection
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp", "--api-key", "ctx7sk_your_api_key"]
    }
  }
}

// VS Code configuration (settings.json)
{
  "mcp": {
    "servers": {
      "context7": {
        "type": "stdio",
        "command": "npx",
        "args": ["-y", "@upstash/context7-mcp", "--api-key", "ctx7sk_your_api_key"]
      }
    }
  }
}

// VS Code with environment variable
{
  "mcp": {
    "servers": {
      "context7": {
        "type": "stdio",
        "command": "npx",
        "args": ["-y", "@upstash/context7-mcp"],
        "env": {
          "CONTEXT7_API_KEY": "ctx7sk_your_api_key"
        }
      }
    }
  }
}

// Claude Code CLI setup
// Remote connection:
claude mcp add --transport http context7 https://mcp.context7.com/mcp \
  --header "CONTEXT7_API_KEY: ctx7sk_your_api_key"

// Local connection:
claude mcp add context7 -- npx -y @upstash/context7-mcp \
  --api-key ctx7sk_your_api_key

// Windsurf configuration with proxy support
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp", "--api-key", "ctx7sk_your_api_key"],
      "env": {
        "HTTPS_PROXY": "http://proxy.company.com:8080"
      }
    }
  }
}

// Docker-based setup (cline_mcp_settings.json)
{
  "mcpServers": {
    "context7": {
      "autoApprove": [],
      "disabled": false,
      "timeout": 60,
      "command": "docker",
      "args": ["run", "-i", "--rm", "context7-mcp"],
      "transportType": "stdio"
    }
  }
}
```

### Configuration: context7.json Schema

Project configuration file for customizing how Context7 parses and indexes your library's documentation.

```json
// Basic configuration example
{
  "projectTitle": "My Awesome Library",
  "description": "A powerful library for building modern web applications",
  "folders": ["docs", "examples"],
  "excludeFolders": ["src", "tests"],
  "excludeFiles": ["CHANGELOG.md", "LICENSE"],
  "rules": [
    "Always use TypeScript for better type safety",
    "Import from the main package, not submodules"
  ],
  "previousVersions": [
    {
      "tag": "v2.0.0",
      "title": "Version 2.0.0 (Legacy)"
    }
  ]
}

// Advanced configuration with patterns
{
  "projectTitle": "Upstash Ratelimit",
  "description": "Ratelimiting library based on Upstash Redis",
  "folders": ["documentation/**", "guides/*"],
  "excludeFolders": ["**/test/**", "**/tests/**", "node_modules"],
  "excludeFiles": ["README-dev.md", "CONTRIBUTING.md", "package.json"],
  "rules": [
    "Use Upstash Redis as a database",
    "Use single region setup for optimal performance",
    "Always handle rate limit errors gracefully"
  ],
  "previousVersions": [
    {
      "tag": "v1.2.1",
      "title": "version 1.2.1"
    },
    {
      "tag": "v1.0.0",
      "title": "version 1.0.0 (Initial Release)"
    }
  ]
}

// Minimal configuration
{
  "projectTitle": "Simple Library",
  "description": "Minimal library with auto-discovery of documentation"
}

// Full schema reference available at:
// https://context7.com/schema/context7.json
```

### Client IP Encryption

Internal encryption mechanism for client IP addresses when communicating with Context7 API.

```typescript
// Encryption configuration via environment variable
// Set a 64-character hex key (32 bytes for AES-256)
export CLIENT_IP_ENCRYPTION_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// Start server with custom encryption key
CLIENT_IP_ENCRYPTION_KEY="your_64_char_hex_key" \
  npx @upstash/context7-mcp --transport http

// Internal implementation (from encryption.ts):
// AES-256-CBC encryption with random IV per request
// Format: "iv:encrypted_data" (hex encoded)
// Example encrypted IP: "a1b2c3d4e5f6....:4f5e6d7c8b9a...."

// Headers sent to Context7 API:
// {
//   "mcp-client-ip": "iv:encrypted_ip",
//   "Authorization": "Bearer ctx7sk_your_api_key",
//   "X-Context7-Source": "mcp-server"
// }

// Note: Encryption key validation
// - Must be exactly 64 hexadecimal characters
// - Represents 32 bytes for AES-256
// - Falls back to unencrypted if key is invalid
// - Default key used if CLIENT_IP_ENCRYPTION_KEY not set
```

## Summary and Integration Patterns

Context7 MCP serves as a critical bridge between AI coding assistants and up-to-date library documentation, eliminating the common issues of outdated training data and hallucinated APIs. The primary use cases include real-time documentation retrieval during code generation, version-specific API reference lookup, and guided implementation following library best practices. Developers integrate Context7 by adding "use context7" to their prompts in supported MCP clients like Cursor, Claude Code, VS Code, Windsurf, and others, or by configuring automatic invocation through client rules. The server handles authentication transparently, supports both free and API-key-based tiers with different rate limits, and works seamlessly behind corporate proxies.

Integration patterns fall into two main categories: local stdio integration for direct MCP client embedding, and remote HTTP integration for centralized deployment scenarios. Local integration provides the lowest latency and simplest setup, while remote integration enables sharing a single Context7 instance across teams and centralized rate limit management. The documentation retrieval system uses pagination for efficient content delivery, allowing LLMs to request additional pages when more context is needed for comprehensive answers. Advanced users can customize documentation indexing through the `context7.json` configuration file, specifying which folders to include, files to exclude, library-specific rules for AI assistants, and supporting multiple versions. The server's architecture supports enterprise requirements including proxy authentication, encrypted client IP transmission via AES-256-CBC, comprehensive error handling with actionable error messages, and intelligent source reputation and benchmark scoring to help LLMs select the most authoritative documentation sources.
