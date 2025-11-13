# Policy: Routing + Knowledge Split for Cline

*Working with Cipher MCP aggregator hub for intelligent routing and knowledge management*

## 🎯 Core Working Model

Cline works through a single MCP hub: **Cipher (aggregator)**
- Exposes tools from 9+ connected MCP servers
- Provides unified interface for code editing, search, knowledge, and filesystem
- Enables shared project knowledge across IDEs and sessions

## 📋 Decision Framework: Where to Put Knowledge

### Put in **Cipher / memory-bank (SHARED)**
**Purpose**: Cross-IDE persistence and shared learning

**Types of knowledge to store:**
- Architecture overviews and domain concepts
- API contracts and event flows
- Naming conventions and invariants
- "How to contribute" guidance
- Code review checklists
- Stable troubleshooting runbooks
- Design decisions and rationale
- Root cause analysis summaries
- Tool usage policies (global)

**Storage pattern**: Use `memory-bank.upsert` with concise title + summary (< 500 words)

### Keep **local in Cline rules (PER-REPO)**
**Purpose**: Repo-specific commands and behaviors

**Types of knowledge to keep local:**
- Test/build/dev commands (`pnpm test`, `just dev`, `docker compose up`)
- Folder layout and entry points
- Environment quirks specific to repo
- IDE/client-specific behaviors
- Local file paths and patterns
- Repo-specific shortcuts and aliases

**Storage pattern**: Edit `.clinerules` with short, factual lines

## 🔧 Tool Routing Rules

### Always follow these routing decisions:

**For Code Edits:**
- Use `morph.*` for any non-trivial change (>5 lines, multiple files/regions)
- Use inline edit only for trivial one-liners
- Morph handles Fast Apply for complex refactoring

**For Search/Understanding:**
- Use `code-index.*` to locate functions, usages, patterns
- Use `code-index.build_deep_index` before major work
- Use `code-index.search_code_advanced` for complex searches
- **NEVER** manually scan directories

**For File Operations:**
- Use `filesystem.*` for read/write operations
- Let Morph handle file modifications during code changes
- Use filesystem tools for configuration and documentation updates

**For Knowledge Management:**
- Use `memory-bank.*` after completing meaningful tasks
- Store durable learning and insights
- Update existing entries rather than duplicating

**For External Data:**
- Use dedicated MCPs (firecrawl, context7) for web data
- Avoid web calls unless specifically needed
- Route through appropriate MCP servers

## 🔄 Operational Checklist (Per Repository Session)

### Session Start
1. **`code-index.set_project_path(<repo-root>)`** - Establish workspace context
2. **`code-index.build_deep_index`** - Build searchable code index
3. Check for existing relevant knowledge in `memory-bank.*`

### During Work
4. For edits → `morph.*` (non-trivial changes)
5. For search → `code-index.*` (understanding and navigation)
6. For files → `filesystem.*` (when needed by other tools)
7. For knowledge → `memory-bank.upsert` (after significant learnings)

### Session End
8. Summarize key learnings to `memory-bank.*`
9. Clean up temporary knowledge entries
10. Document any new patterns or decisions discovered

## 🚦 Context/Token Management

### Best Practices:
- **Fetch top-K snippets** via `code-index.*` rather than pasting large files
- **Summarize long conversations** locally and store to `memory-bank.*`
- **Use short, targeted tool calls** with explicit arguments
- **Prefer iterative refinement** over bulk operations

### When Conversations Get Long:
1. Summarize current state locally
2. Store summary to `memory-bank.*`
3. Drop older conversation turns
4. Continue with fresh context

## ⚖️ Conflict Resolution

### Hierarchy (Higher wins):
1. **Cipher global policies** (memory-bank knowledge)
2. **Project-specific knowledge** (in memory-bank)
3. **Local Cline rules** (.clinerules)

### When Conflicts Occur:
- **Prefer Cipher policies** over local rules
- **If tool fails**, report error and continue with available tools
- **Never switch to manual bulk scanning** as fallback
- **Escalate unclear conflicts** to human for guidance

## 🛡️ Safety Guidelines

### Never Store in Memory-Bank:
- Secrets, API keys, or credentials
- Customer data or PII
- Private implementation details
- Temporary debugging information

### Always Redact:
- Credentials and tokens
- Customer identifiers
- Internal system details
- Sensitive operational data

### Keep Entries Focused:
- **One idea per entry** in memory-bank
- **Under 500 words** per knowledge entry
- **Concise summaries** with actionable insights
- **Link to documentation** rather than duplicate content

## 🎯 Success Metrics

### Knowledge Quality:
- Memory-bank entries provide value across sessions
- Local rules stay minimal and specific
- No duplicate knowledge across layers
- Clear routing decisions for tools

### Efficiency:
- Reduced manual file scanning
- Faster code navigation and editing
- Shared learning across IDEs
- Consistent project context

---

**Remember**: The goal is intelligent routing and shared learning. Use Cipher's tools strategically and keep knowledge appropriately placed for maximum reuse and minimal duplication.
