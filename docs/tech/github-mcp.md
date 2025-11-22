# GitHub MCP Server

**GitHub**: https://github.com/github/github-mcp-server

---

## Overview

A Model Context Protocol server for the GitHub API, enabling file operations, repository management, search functionality, and more.

**Key Features:**
- Automatic branch creation for file operations
- Comprehensive error handling with clear messages
- Git history preservation without force pushing
- Batch operations for single and multi-file operations
- Advanced search capabilities (code, issues/PRs, users)

---

## Installation

### Using npm (via npx)
```bash
npx @github/github-mcp-server
```

### Direct installation
```bash
npm install -g @github/github-mcp-server
```

---

## Configuration

**Required Environment Variable:**
- `GITHUB_PERSONAL_ACCESS_TOKEN`: GitHub personal access token with appropriate permissions

**Token Permissions Required:**
- `repo`: Full control of private repositories
- `read:org`: Read organization data
- `user`: Update user data (if needed)

**Create a GitHub Personal Access Token:**
1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Click "Generate new token"
3. Select scopes based on your needs
4. Copy the token and set it as environment variable

---

## Available Tools

### Repository Management
1. `create_or_update_file`
   - Create or update a single file in a repository
   - Inputs: owner, repo, path, content, message, branch, sha (optional)

2. `create_repository`
   - Create a new repository
   - Inputs: name, description, private, auto_init

3. `get_file_contents`
   - Get contents of a file or directory
   - Inputs: owner, repo, path, branch (optional)

### Search Operations
4. `search_repositories`
   - Search for repositories
   - Inputs: query, sort, order, per_page, page

5. `search_code`
   - Search for code across repositories
   - Inputs: query, sort, order, per_page, page

6. `search_issues_and_prs`
   - Search for issues and pull requests
   - Inputs: query, sort, order, per_page, page

### Issue/PR Management
7. `create_issue`
   - Create a new issue
   - Inputs: owner, repo, title, body, labels, assignees

8. `create_pull_request`
   - Create a new pull request
   - Inputs: owner, repo, title, body, head, base

9. `get_issue`
   - Get issue details
   - Inputs: owner, repo, issue_number

### User Operations
10. `get_user`
    - Get user information
    - Inputs: username

11. `get_authenticated_user`
    - Get authenticated user information
    - Inputs: none

---

## Usage Examples

### Setting up Environment
```bash
export GITHUB_PERSONAL_ACCESS_TOKEN="your_github_token_here"
```

### Basic Repository Operations
```json
{
  "tool": "create_or_update_file",
  "arguments": {
    "owner": "your-username",
    "repo": "your-repo",
    "path": "README.md",
    "content": "# My Project\n\nThis is my project description.",
    "message": "Initial commit",
    "branch": "main"
  }
}
```

### Searching for Repositories
```json
{
  "tool": "search_repositories",
  "arguments": {
    "query": "language:python machine learning",
    "sort": "stars",
    "order": "desc",
    "per_page": 10
  }
}
```

### Creating an Issue
```json
{
  "tool": "create_issue",
  "arguments": {
    "owner": "your-username",
    "repo": "your-repo",
    "title": "Bug: Application crashes on startup",
    "body": "The application crashes when starting with the following error...",
    "labels": ["bug", "high-priority"]
  }
}
```

---

## Known Limitations

1. **Rate Limiting**: Subject to GitHub API rate limits
2. **Token Permissions**: Requires careful scope management
3. **Branch Auto-creation**: May create branches automatically for file operations
4. **No Webhook Support**: Does not support GitHub webhooks

---

## Testing

```bash
# Set environment variable
export GITHUB_PERSONAL_ACCESS_TOKEN="your_token"

# Test server startup
npx @github/github-mcp-server

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @github/github-mcp-server
```

---

## Related Documentation

- [GitHub MCP Server Repository](http://github.com/github/github-mcp-server)
- [GitHub API Documentation](https://docs.github.com/en/rest)
- [MCP Protocol Specification](https://modelcontextprotocol.io)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
