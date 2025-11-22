# Tool Groups Configuration

**Purpose**: Organize MCP servers into logical collections for different use cases
**Status**: Phase 1C Implementation Plan
**Target**: Universal group as Priority #1

---

## Overview

Tool groups allow exposing only specific subsets of tools to different clients, enabling:
- **Security**: Limit tool access per client/application
- **Simplicity**: Reduce tool clutter in IDEs
- **Specialization**: Tailor tools for specific workflows
- **Performance**: Faster tool discovery and loading

---

## Universal Tool Group (Priority #1)

**Purpose**: Core tools for all AI workflows and general development
**Servers**: 7 total servers
**Expected Tools**: ~33 tools
**Use Cases**: General development assistance, research, file management, web interaction

### Included Servers

1. **context7** (Documentation Lookup)
   - Tools: `get_library_docs`, `search_libraries`
   - Purpose: Technical documentation and API reference

2. **brave-search** (Web Search)
   - Tools: `brave_web_search`, `brave_local_search`
   - Purpose: Web search and local business search

3. **filesystem** (File Operations)
   - Tools: `read_file`, `write_file`, `list_directory`, `create_directory`
   - Purpose: File system operations and code management

4. **firecrawl** (Web Scraping)
   - Tools: `scrape_url`, `crawl_url`, `map_url`
   - Purpose: Advanced web content extraction

5. **fetch** (HTTP Requests)
   - Tools: `fetch`
   - Purpose: Direct HTTP content fetching

6. **github** (Repository Management)
   - Tools: `create_or_update_file`, `search_repositories`, `create_issue`
   - Purpose: GitHub repository and project management

7. **memory** (Persistent Storage)
   - Tools: `memory_create_entity`, `memory_search`, `memory_create_relation`
   - Purpose: Cross-session knowledge persistence

---

## Backend Development Group

**Purpose**: Universal + database and research tools for backend development
**Additional Servers**: 3 servers
**Total Tools**: ~45 tools
**Use Cases**: Database development, API development, system administration

### Additional Servers (beyond Universal)

8. **sqlite** (Database Operations)
   - Tools: `read_query`, `write_query`, `describe_table`, `list_tables`
   - Purpose: SQLite database operations and analytics

9. **gpt-researcher** (Deep Research)
   - Tools: `research_topic`, `generate_report`
   - Purpose: Comprehensive research and report generation

10. **morph-fast-apply** (Code Transformation)
    - Tools: `apply_diff`, `transform_code`
    - Purpose: Automated code modifications and refactoring

---

## Frontend Development Group

**Purpose**: Universal + browser automation and code transformation
**Additional Servers**: 2 servers
**Total Tools**: ~40 tools
**Use Cases**: Web development, UI testing, frontend automation

### Additional Servers (beyond Universal)

8. **playwright** (Browser Automation)
   - Tools: `puppeteer_navigate`, `puppeteer_screenshot`, `puppeteer_click`, `puppeteer_fill`
   - Purpose: Browser automation, E2E testing, web scraping

9. **morph-fast-apply** (Code Transformation)
   - Tools: `apply_diff`, `transform_code`
   - Purpose: Automated code modifications and refactoring

---

## Implementation Strategy

### Phase 1C: Tool Group Creation

1. **Create Universal Group First**
   ```bash
   # Register universal tool group with jarvis
   mcpjungle create group -c ./config/jarvis/groups/universal.json
   ```

2. **Test Group Functionality**
   ```bash
   # List tools in universal group
   mcpjungle list tools --group universal

   # Test tool invocation in group context
   mcpjungle invoke filesystem__read_file --group universal --input '{"path": "README.md"}'
   ```

3. **Create Specialized Groups**
   - Backend development group
   - Frontend development group
   - Data science group (future)
   - DevOps group (future)

### Configuration File Structure

**File**: `config/jarvis/groups/universal.json`
```json
{
  "name": "universal",
  "description": "Core tools for all AI workflows",
  "included_servers": [
    "context7",
    "brave-search",
    "filesystem",
    "firecrawl",
    "fetch",
    "github",
    "memory"
  ]
}
```

---

## IDE Configuration Examples

### Claude Desktop (Universal Group)
```json
{
  "mcpServers": {
    "mcpjungle-universal": {
      "command": "npx",
      "args": ["mcp-remote", "http://localhost:8080/v0/groups/universal/mcp", "--allow-http"]
    }
  }
}
```

### Cursor (Universal Group)
```json
{
  "mcpServers": {
    "mcpjungle-universal": {
      "url": "http://localhost:8080/v0/groups/universal/mcp"
    }
  }
}
```

---

## Tool Group Benefits

### For Developers
- **Reduced Complexity**: Only see relevant tools for your task
- **Faster Discovery**: Find the right tool quickly
- **Better Performance**: Smaller tool lists load faster

### For Organizations
- **Security**: Limit tool access per team/role
- **Compliance**: Audit and control tool usage
- **Training**: Easier onboarding with focused toolsets

---

## Future Groups (Phase 6)

### Data Science Group
- **Servers**: sqlite, gpt-researcher, memory, fetch
- **Focus**: Data analysis, research, visualization

### DevOps Group
- **Servers**: github, filesystem, fetch, memory
- **Focus**: Infrastructure, deployment, monitoring

### Security Group
- **Servers**: fetch, memory, filesystem
- **Focus**: Security analysis, vulnerability scanning

### Content Creation Group
- **Servers**: firecrawl, fetch, memory, gpt-researcher
- **Focus**: Content research, creation, management

---

## Migration Strategy

### From Individual Servers to Groups
1. **Gradual Migration**: Start with universal group
2. **Backward Compatibility**: Keep individual server access
3. **User Choice**: Let users choose individual vs. group access
4. **Performance Monitoring**: Track usage and performance

### Group Evolution
1. **User Feedback**: Collect feedback on tool groupings
2. **Usage Analytics**: Analyze which tools are used together
3. **Dynamic Groups**: Consider user-customizable groups
4. **AI-Assisted Grouping**: Use AI to suggest optimal groupings

---

## Related Documentation

- [MCPJungle Tool Groups](https://github.com/mcpjungle/MCPJungle#tool-groups)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Tool Group Configuration Guide](docs/guides/tool-group-usage.md)

---

**Last Updated**: 2025-11-18
**Status**: Phase 1C Implementation Plan
**Next Steps**: Create universal group configuration and register with jarvis
