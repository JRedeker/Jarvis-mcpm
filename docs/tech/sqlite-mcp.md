# SQLite MCP Server

**GitHub**: https://github.com/modelcontextprotocol/servers-archived/tree/main/src/sqlite
**Package**: mcp-server-sqlite
**Language**: Python
**Transport**: stdio
**Status**: ⚠️ Archived - Simple SQLite implementation

---

## Overview

A Model Context Protocol (MCP) server implementation that provides database interaction and business intelligence capabilities through SQLite. This server enables running SQL queries, analyzing business data, and automatically generating business insight memos.

**Key Features:**
- SQLite database operations
- Business intelligence and analytics
- Automatic insight generation
- Dynamic resource creation (business memos)
- Interactive demonstration prompts

**⚠️ Important**: This server is archived and provides basic SQLite functionality.

---

## Installation

### Using uv (recommended)
```bash
# Install uv if not already installed
curl -LsSf https://astral.sh/uv/install.sh | sh

# Run directly with uvx
uvx mcp-server-sqlite
```

### Using pip
```bash
pip install mcp-server-sqlite
```

---

## Configuration

The SQLite server requires no environment variables or API keys. It runs with default settings and creates databases locally.

**Database Storage**: SQLite databases are created locally in the working directory.

---

## Components

### Resources
The server exposes a single dynamic resource:
- `memo://insights`: A continuously updated business insights memo that aggregates discovered insights during analysis
  - Auto-updates as new insights are discovered via the append-insight tool

### Prompts
The server provides a demonstration prompt:
- `mcp-demo`: Interactive prompt that guides users through database operations
  - Required argument: `topic` - The business domain to analyze
  - Generates appropriate database schemas and sample data
  - Guides users through analysis and insight generation
  - Integrates with the business insights memo

### Tools
The server offers six core tools:

#### Query Tools
1. `read_query`
   - Execute SELECT queries to read data from the database
   - Input: `query` (string) - The SELECT SQL query to execute
   - Returns: Query results as array of objects

2. `write_query`
   - Execute INSERT, UPDATE, or DELETE queries
   - Input: `query` (string) - The SQL query to execute
   - Returns: Number of affected rows or error message

#### Schema and Analysis
3. `describe_table`
   - Get schema information for a specific table
   - Input: `table_name` (string) - Name of the table to describe
   - Returns: Table schema with column information

4. `list_tables`
   - List all tables in the database
   - Input: None
   - Returns: Array of table names

#### Business Intelligence
5. `append_insight`
   - Append a business insight to the insights memo
   - Input: `insight` (string) - The business insight to append
   - Returns: Success confirmation

6. `get_business_insights`
   - Retrieve the current business insights memo
   - Input: None
   - Returns: Complete insights memo content

---

## Usage Examples

### Basic Database Operations
```json
{
  "tool": "list_tables",
  "arguments": {}
}
```

```json
{
  "tool": "describe_table",
  "arguments": {
    "table_name": "sales"
  }
}
```

### Reading Data
```json
{
  "tool": "read_query",
  "arguments": {
    "query": "SELECT * FROM sales WHERE date >= '2024-01-01' ORDER BY date DESC LIMIT 10"
  }
}
```

### Writing Data
```json
{
  "tool": "write_query",
  "arguments": {
    "query": "INSERT INTO sales (product, amount, date) VALUES ('Widget A', 99.99, '2024-01-15')"
  }
}
```

### Business Intelligence
```json
{
  "tool": "read_query",
  "arguments": {
    "query": "SELECT product, SUM(amount) as total_sales FROM sales GROUP BY product ORDER BY total_sales DESC"
  }
}
```

```json
{
  "tool": "append_insight",
  "arguments": {
    "insight": "Widget A is the top-selling product with $2,450 in total sales this month"
  }
}
```

```json
{
  "tool": "get_business_insights",
  "arguments": {}
}
```

---

## Business Intelligence Workflow

1. **Data Discovery**: Use `list_tables` and `describe_table` to understand available data
2. **Analysis**: Use `read_query` to analyze business data with SQL
3. **Insight Generation**: Identify key business insights from query results
4. **Documentation**: Use `append_insight` to document findings
5. **Reporting**: Use `get_business_insights` to review accumulated insights

---

## SQL Support

The server supports standard SQLite SQL including:

**Data Types**: INTEGER, REAL, TEXT, BLOB, NULL
**Functions**: COUNT, SUM, AVG, MIN, MAX, date functions
**Joins**: INNER, LEFT, RIGHT, FULL OUTER
**Subqueries**: Supported in SELECT, WHERE, and FROM clauses
**Views**: CREATE VIEW and DROP VIEW
**Indexes**: CREATE INDEX and DROP INDEX

---

## Known Limitations

1. **Archived Status**: This server is no longer actively maintained
2. **SQLite Only**: Limited to SQLite functionality, no other database support
3. **Single Database**: Works with one database at a time
4. **Basic BI**: Simple business intelligence, not enterprise-grade analytics
5. **No Transactions**: No explicit transaction management
6. **File-based**: Database is file-based, not server-based

---

## Testing

```bash
# Test server startup
uvx mcp-server-sqlite

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | uvx mcp-server-sqlite

# Test database operations (requires MCP client)
```

---

## Use Cases

1. **Data Analysis**: Analyze CSV files or business data
2. **Prototyping**: Quick database prototyping for applications
3. **Business Intelligence**: Generate insights from business data
4. **Reporting**: Create simple reports and dashboards
5. **Learning**: Learn SQL and database concepts

---

## Advanced Features

### Creating Custom Tables
```sql
CREATE TABLE customers (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Complex Analytics
```sql
-- Monthly sales trends
SELECT
    strftime('%Y-%m', date) as month,
    COUNT(*) as transactions,
    SUM(amount) as total_sales,
    AVG(amount) as avg_transaction
FROM sales
GROUP BY strftime('%Y-%m', date)
ORDER BY month;
```

---

## Related Documentation

- [Official MCP Servers Repository](https://github.com/modelcontextprotocol/servers)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [SQLite Server Source Code](https://github.com/modelcontextprotocol/servers-archived/tree/main/src/sqlite)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
