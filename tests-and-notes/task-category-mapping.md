# Task Category Mapping & Preferred Tool Lists

## Executive Summary

This document provides **intelligent task categorization** for the cipher aggregator system, mapping common AI agent tasks to their **optimal tool selections**. This ensures agents always choose the best tool for the job while leveraging cipher's dual-layer architecture (built-in tools + external MCP servers).

## Core Architecture Understanding

### Cipher's Dual Tool Layer

**Layer 1: Built-in Cipher Tools** (Memory & Reasoning)
- Memory management and knowledge extraction
- Reasoning and reflection capabilities
- Workspace memory for team collaboration
- Knowledge graph operations

**Layer 2: External MCP Server Tools** (Specialized Functions)
- Domain-specific operations (GitHub, web search, file ops)
- Specialized capabilities (testing, database, monitoring)
- API integrations and data retrieval

## Task Category Framework

### 🔍 **1. Code & Development Analysis**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Code Search** | `code-index` | `filesystem_list`, `cipher_bash` | External (code-index-mcp) |
| **File Analysis** | `code-index_file_summary` | `filesystem_read` | External (code-index-mcp) |
| **Pattern Recognition** | `code-index_search` | Manual file scanning | External (code-index-mcp) |
| **Code Modification** | `morph` | `filesystem_write` | Built-in vs External |
| **Memory Storage** | `cipher_extract_and_operate_memory` | `memory_bank_store` | Built-in vs External |

**Decision Rules:**
- ✅ **Always use `code-index`** for code search/analysis (indexed, fast)
- ✅ **Use `morph`** for non-trivial code edits (built-in, intelligent)
- ✅ **Use `filesystem`** only for simple file operations
- ❌ **Never** use manual file scanning for code search

### 🧠 **2. Knowledge & Memory Management**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Store Knowledge** | `cipher_extract_and_operate_memory` | `memory_bank_store` | Built-in |
| **Search Knowledge** | `cipher_memory_search` | `memory_bank_search` | Built-in |
| **Reasoning Storage** | `cipher_store_reasoning_memory` | N/A | Built-in |
| **Team Memory** | `cipher_workspace_search` | N/A | Built-in |
| **Knowledge Graph** | `cipher_search_graph` | N/A | Built-in |

**Decision Rules:**
- ✅ **Always prefer built-in Cipher tools** for memory operations
- ✅ **Use `cipher_extract_and_operate_memory`** for intelligent knowledge extraction
- ✅ **Use `cipher_memory_search`** for semantic knowledge retrieval
- ✅ **Use workspace tools** when `USE_WORKSPACE_MEMORY=true`

### 🐙 **3. Repository & Version Control**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **GitHub Operations** | `github` MCP tools | `fetch` | External |
| **Repository Management** | `github_list_repos`, `github_create_repo` | Web scraping | External |
| **Issue Management** | `github_search_issues`, `github_create_pr` | Manual GitHub | External |
| **Code Review** | `github` + `code-index` | Manual review | External + External |

**Decision Rules:**
- ✅ **Always use `github` MCP server** for GitHub operations
- ❌ **Never use `fetch`** for GitHub operations
- ✅ **Combine with `code-index`** for comprehensive code review

### 🌐 **4. Web Operations & Data Retrieval**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Web Search** | `brave-search` | `firecrawl_search` | External |
| **Web Scraping** | `firecrawl` | `fetch` | External |
| **Content Extraction** | `firecrawl_extract` | `firecrawl_scrape` | External |
| **Simple HTTP** | `httpie` | `fetch` | External |
| **Complex HTTP** | `schemathesis` | `httpie` | External |

**Decision Rules:**
- ✅ **Use `brave-search`** for general web search (search engine)
- ✅ **Use `firecrawl`** for web scraping and content extraction
- ✅ **Use `httpie`** for simple HTTP requests
- ✅ **Use `schemathesis`** for API testing with schemas
- ❌ **Never use `fetch`** for complex operations (limited functionality)

### 🗃️ **5. File & Data Operations**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Local Files** | `filesystem` | `fetch` | External |
| **File Reading** | `filesystem_read` | `code-index_file_summary` | External |
| **File Writing** | `filesystem_write` | `morph` | External |
| **Batch Operations** | `file-batch` | Multiple `filesystem` calls | External |
| **Database** | `sql` | `fetch` | External |

**Decision Rules:**
- ✅ **Use `filesystem`** for all local file operations
- ✅ **Use `file-batch`** for multiple file operations
- ✅ **Use `sql`** for database operations
- ❌ **Never use `fetch`** for local files

### 🧪 **6. Testing & Quality Assurance**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Unit Testing** | `pytest` | Manual execution | External |
| **API Testing** | `schemathesis` | `httpie` | External |
| **Coverage Analysis** | `pytest` with coverage | Manual analysis | External |
| **Test Reports** | `get_test_report` | Manual reporting | External |

**Decision Rules:**
- ✅ **Use `pytest`** for unit testing and coverage
- ✅ **Use `schemathesis`** for API schema-based testing
- ✅ **Use dedicated testing tools** over manual approaches

### 📊 **7. Documentation & Research**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Library Docs** | `context7` | `firecrawl` | External |
| **API Documentation** | `context7_resolve_library_id` | Web scraping | External |
| **Research** | `brave-search` + `firecrawl` | Manual research | External + External |
| **Documentation Storage** | `cipher_extract_and_operate_memory` | File storage | Built-in |

**Decision Rules:**
- ✅ **Use `context7`** for library and API documentation
- ✅ **Combine `brave-search` + `firecrawl`** for comprehensive research
- ✅ **Store findings** with built-in memory tools

### 🔧 **8. System Operations & Infrastructure**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **Command Execution** | `cipher_bash` | `fetch` | Built-in |
| **Docker Operations** | `docker` | `cipher_bash` | External |
| **Browser Automation** | `playwright` | `firecrawl` | External |
| **Monitoring** | `prometheus` | Custom scripts | External |

**Decision Rules:**
- ✅ **Use `cipher_bash`** for system commands
- ✅ **Use specialized tools** for infrastructure operations
- ❌ **Avoid custom scripts** when MCP tools available

### 🎨 **9. UI & Development Tools**

| Sub-Category | Preferred Tool | Alternative Tools | Built-in vs External |
|--------------|----------------|-------------------|---------------------|
| **TUI Development** | `textual-devtools` | Manual development | External |
| **Svelte Development** | `svelte` | Manual setup | External |
| **UI Magic** | `magic-mcp` | Custom solutions | External |

**Decision Rules:**
- ✅ **Use specialized development tools** for framework-specific work
- ✅ **Leverage MCP tools** over manual setup

## Priority-Based Decision Matrix

### Level 1: Critical (Never Override)
- **GitHub Operations** → `github` MCP server
- **Web Search** → `brave-search`
- **API Testing** → `schemathesis`
- **Database Operations** → `sql`
- **Memory Management** → Built-in Cipher tools

### Level 2: High (Strong Preference)
- **Code Analysis** → `code-index`
- **Web Scraping** → `firecrawl`
- **File Operations** → `filesystem`
- **Testing** → `pytest`

### Level 3: Medium (Context Dependent)
- **HTTP Requests** → `httpie` vs `fetch`
- **Documentation** → `context7` vs `firecrawl`
- **Command Execution** → `cipher_bash` vs alternatives

## Task Classification Algorithm

```python
def classify_task_and_route(task_description):
    """
    Intelligent task classification and routing for cipher aggregator
    """
    task_lower = task_description.lower()

    # Level 1: Critical domain-specific routing
    if any(keyword in task_lower for keyword in ['github', 'repository', 'pull request', 'issue']):
        return 'github', ['github_list_repos', 'github_search_issues', 'github_create_pr']

    elif any(keyword in task_lower for keyword in ['search for', 'find information', 'google']):
        return 'brave-search', ['brave_web_search']

    elif any(keyword in task_lower for keyword in ['api test', 'schema test', 'openapi']):
        return 'schemathesis', ['schemathesis_test_api', 'schemathesis_load_schema']

    # Level 2: High-priority routing
    elif any(keyword in task_lower for keyword in ['find code', 'search code', 'analyze code']):
        return 'code-index', ['code_index_search', 'code_index_file_summary']

    elif any(keyword in task_lower for keyword in ['scrape', 'extract data', 'crawl website']):
        return 'firecrawl', ['firecrawl_extract', 'firecrawl_scrape']

    # Level 3: Context-dependent routing
    elif any(keyword in task_lower for keyword in ['file', 'read', 'write', 'directory']):
        return 'filesystem', ['filesystem_read', 'filesystem_write', 'filesystem_list']

    # Default to memory tools for knowledge tasks
    else:
        return 'cipher-memory', ['cipher_memory_search', 'cipher_extract_and_operate_memory']
```

## Performance Optimization Rules

### Speed Hierarchy (Fastest to Slowest)
1. **Built-in Cipher Tools** (in-memory operations)
2. **Local Tools** (`filesystem`, `code-index`)
3. **API Tools** (`github`, `brave-search`)
4. **Scraping Tools** (`firecrawl`)
5. **Monitoring Tools** (`prometheus`)

### Reliability Hierarchy (Most to Least Reliable)
1. **Built-in Tools** (local, no network)
2. **Official APIs** (GitHub, Brave)
3. **Specialized MCP Tools** (maintained servers)
4. **Generic Tools** (`fetch`, basic HTTP)

### Cost Efficiency
- **Prefer local operations** over API calls when possible
- **Batch operations** to minimize API usage
- **Cache results** using memory tools
- **Use appropriate timeouts** to avoid hanging

---

*This mapping should be used in conjunction with the Tool Routing Guide for comprehensive intelligent routing.*