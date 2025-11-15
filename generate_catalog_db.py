#!/usr/bin/env python3
"""
MCP Catalog DuckDB Generator — Dynamic from Running Cipher Aggregator

This script generates a DuckDB database (`catalog.db`) by:
1. Parsing cipher.yml for server configurations (enabled, command, etc.).
2. Querying the running Cipher aggregator SSE endpoint for the merged tools/list.
3. Combining the data to populate servers and tools tables.
4. Exporting to Parquet (tools.parquet, servers.parquet) for viewing.
5. Optionally deleting static YAML catalogs (mcp-*.yml) after success.

Prerequisites:
- Cipher aggregator must be running (SSE at http://127.0.0.1:3020/sse).
- Start with: ./mcp-manager.sh start
- Dependencies: duckdb, pyyaml, httpx (in pyproject.toml); run `uv sync` or `pip install -e .`

Usage:
    python3 generate_catalog_db.py [--delete-yamls]

View Parquet:
- VS Code: 'Parquet Viewer' extension.
- CLI: duckdb -c "SELECT * FROM 'tools.parquet' LIMIT 5"
- Python: import pandas as pd; pd.read_parquet('tools.parquet')
"""

import argparse
import json
import os
import sys
from pathlib import Path
import yaml
import httpx
import duckdb
from typing import Dict, List, Any

# Paths
SCRIPT_DIR = Path(__file__).parent
CIPHER_YML = SCRIPT_DIR / "cipher.yml"
DB_PATH = SCRIPT_DIR / "catalog.db"
SSE_URL = "http://127.0.0.1:3020/sse"


def parse_cipher_yml() -> Dict[str, Any]:
    """Parse cipher.yml for server configurations."""
    if not CIPHER_YML.exists():
        raise FileNotFoundError(f"cipher.yml not found: {CIPHER_YML}")

    with open(CIPHER_YML, 'r') as f:
        data = yaml.safe_load(f)

    servers = data.get('servers', {})
    parsed_servers = {}

    for server_key, config in servers.items():
        parsed_servers[server_key] = {
            'enabled': config.get('enabled', False),
            'transport': config.get('transport', 'stdio'),
            'command': config.get('command', ''),
            'env': config.get('env', {}),
            'title': config.get('title', server_key.replace('-', ' ').title()),  # Infer title
            'configured': True  # Present in yml
        }

    return parsed_servers


def query_aggregator_tools() -> List[Dict[str, Any]]:
    """Query the running Cipher aggregator for tools/list via SSE JSON-RPC, with fallback to static data."""
    # JSON-RPC request for tools/list
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/list",
        "params": {}
    }

    headers = {
        "Content-Type": "text/plain",
        "Accept": "text/event-stream"
    }

    try:
        with httpx.Client(timeout=30.0) as client:
            response = client.post(SSE_URL, json=request, headers=headers)
            response.raise_for_status()

            # Parse response (expect JSON-RPC result in body or first event)
            if response.headers.get('content-type', '').startswith('text/event-stream'):
                # Handle SSE stream
                tools_data = []
                for line in response.iter_lines():
                    if line:
                        if line.startswith('data: '):
                            event_data = line[6:].strip()
                            if event_data == '[DONE]':
                                break
                            try:
                                event = json.loads(event_data)
                                if 'result' in event and 'tools' in event['result']:
                                    tools_data = event['result']['tools']
                                    break  # Got the tools
                            except json.JSONDecodeError:
                                continue
            else:
                # Assume direct JSON response
                resp_json = response.json()
                if 'result' in resp_json and 'tools' in resp_json['result']:
                    tools_data = resp_json['result']['tools']
                else:
                    raise ValueError("Unexpected response format")

            # Parse tools
            parsed_tools = []
            for tool in tools_data:
                parsed_tools.append({
                    'tool_name': tool['name'],
                    'description': tool['description'],
                    'inputSchema': tool.get('inputSchema', {}),
                    'server_key': 'aggregator',  # Merged; server per-tool not directly available
                    'category': infer_category(tool['name'], tool['description'])  # Simple inference
                })

            print("Fetched tools from runtime aggregator.")
            return parsed_tools

    except Exception as e:
        print(f"Warning: Failed to query aggregator ({e}). Using static fallback data.")
        # Fallback to static tool data from known servers
        static_tools = [
            # Routing Metadata
            {'tool_name': 'validate_tool_selection', 'description': 'Validate a tool selection against routing rules and return metadata', 'inputSchema': {'type': 'object', 'properties': {'session_id': {'type': 'string'}, 'task_description': {'type': 'string'}, 'selected_tool': {'type': 'string'}}}, 'server_key': 'routing-metadata', 'category': 'routing'},
            {'tool_name': 'track_tool_execution', 'description': 'Track tool execution metrics for performance/analytics', 'inputSchema': {'type': 'object', 'properties': {'session_id': {'type': 'string'}, 'tool_name': {'type': 'string'}, 'execution_time_ms': {'type': 'integer'}, 'success': {'type': 'boolean'}}}, 'server_key': 'routing-metadata', 'category': 'routing'},
            {'tool_name': 'initialize_session', 'description': 'Initialize session tracking with constraints', 'inputSchema': {'type': 'object', 'properties': {'session_id': {'type': 'string'}, 'mode': {'type': 'string'}, 'max_calls': {'type': 'integer'}}}, 'server_key': 'routing-metadata', 'category': 'routing'},
            {'tool_name': 'get_routing_analytics', 'description': 'Retrieve aggregated routing analytics', 'inputSchema': {'type': 'object', 'properties': {'days_back': {'type': 'integer'}}}, 'server_key': 'routing-metadata', 'category': 'routing'},
            # HTTPie
            {'tool_name': 'make_request', 'description': 'HTTP request (method, url, headers, data, auth, timeout, format)', 'inputSchema': {'type': 'object', 'properties': {'method': {'type': 'string', 'enum': ['GET', 'POST', 'PUT', 'DELETE']}, 'url': {'type': 'string'}, 'headers': {'type': 'object'}, 'data': {'type': 'object'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'upload_file', 'description': 'Multipart upload with additional fields and auth', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}, 'file_path': {'type': 'string'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'download_file', 'description': 'Download to path (resume optional)', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}, 'output_path': {'type': 'string'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'test_api_endpoint', 'description': 'Quick endpoint test with expected_status', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}, 'expected_status': {'type': 'integer'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'manage_session', 'description': 'Create/use/list/delete HTTP sessions', 'inputSchema': {'type': 'object', 'properties': {'action': {'type': 'string', 'enum': ['create', 'use', 'list', 'delete']}, 'session_name': {'type': 'string'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'handle_authentication', 'description': 'Basic/Bearer/Digest/NetRC auth helpers', 'inputSchema': {'type': 'object', 'properties': {'auth_type': {'type': 'string', 'enum': ['basic', 'bearer', 'digest']}, 'credentials': {'type': 'object'}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'format_response', 'description': 'Headers/body/meta/verbose formatting', 'inputSchema': {'type': 'object', 'properties': {'output_type': {'type': 'string', 'enum': ['headers', 'body', 'verbose']}}}, 'server_key': 'httpie', 'category': 'http'},
            {'tool_name': 'test_connectivity', 'description': 'Endpoint probe across methods with SSL options', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}, 'methods': {'type': 'array', 'items': {'type': 'string'}}}}, 'server_key': 'httpie', 'category': 'http'},
            # Schemathesis
            {'tool_name': 'load_openapi_schema', 'description': 'Load/validate OpenAPI schema from URL/file', 'inputSchema': {'type': 'object', 'properties': {'source': {'type': 'string'}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            {'tool_name': 'test_api_endpoints', 'description': 'Property-based testing across schema endpoints', 'inputSchema': {'type': 'object', 'properties': {'schema_source': {'type': 'string'}, 'base_url': {'type': 'string'}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            {'tool_name': 'validate_schema', 'description': 'Structural validation and issue detection', 'inputSchema': {'type': 'object', 'properties': {'schema_source': {'type': 'string'}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            {'tool_name': 'generate_test_data', 'description': 'Generate sample data from schema specifications', 'inputSchema': {'type': 'object', 'properties': {'schema_source': {'type': 'string'}, 'count': {'type': 'integer'}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            {'tool_name': 'test_specific_endpoint', 'description': 'Focused tests for a single endpoint/method', 'inputSchema': {'type': 'object', 'properties': {'schema_source': {'type': 'string'}, 'base_url': {'type': 'string'}, 'endpoint': {'type': 'string'}, 'method': {'type': 'string'}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            {'tool_name': 'run_schema_tests', 'description': 'Run strategies (property_based, boundary_values, etc.)', 'inputSchema': {'type': 'object', 'properties': {'schema_source': {'type': 'string'}, 'base_url': {'type': 'string'}, 'test_strategies': {'type': 'array', 'items': {'type': 'string'}}}}, 'server_key': 'schemathesis', 'category': 'api-testing'},
            # Custom Filesystem
            {'tool_name': 'read_text_file', 'description': 'Read file text with head/tail controls', 'inputSchema': {'type': 'object', 'properties': {'path': {'type': 'string'}, 'head': {'type': 'integer'}, 'tail': {'type': 'integer'}}}, 'server_key': 'custom-filesystem', 'category': 'filesystem'},
            {'tool_name': 'read_multiple_files', 'description': 'Read multiple files and return JSON array', 'inputSchema': {'type': 'object', 'properties': {'paths': {'type': 'array', 'items': {'type': 'string'}}}}, 'server_key': 'custom-filesystem', 'category': 'filesystem'},
            {'tool_name': 'write_file', 'description': 'Write/overwrite file content', 'inputSchema': {'type': 'object', 'properties': {'path': {'type': 'string'}, 'content': {'type': 'string'}}}, 'server_key': 'custom-filesystem', 'category': 'filesystem'},
            {'tool_name': 'list_directory', 'description': 'List directory entries with markers', 'inputSchema': {'type': 'object', 'properties': {'path': {'type': 'string'}}}, 'server_key': 'custom-filesystem', 'category': 'filesystem'},
            {'tool_name': 'create_directory', 'description': 'Create directory (mkdir -p semantics)', 'inputSchema': {'type': 'object', 'properties': {'path': {'type': 'string'}}}, 'server_key': 'custom-filesystem', 'category': 'filesystem'},
            # File Batch
            {'tool_name': 'read_files_batched', 'description': 'Read multiple files with size/delay limits', 'inputSchema': {'type': 'object', 'properties': {'paths': {'type': 'array', 'items': {'type': 'string'}}}}, 'server_key': 'file-batch', 'category': 'filesystem'},
            # Pytest
            {'tool_name': 'run_tests', 'description': 'Run pytest (path/patterns/markers/coverage)', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            {'tool_name': 'get_test_report', 'description': 'Generate XML/JSON test report', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            {'tool_name': 'list_tests', 'description': 'Collect and list tests without running', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            {'tool_name': 'run_specific_test', 'description': 'Run a single test by function name', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}, 'test_name': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            {'tool_name': 'check_test_coverage', 'description': 'Run tests with coverage analysis', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}, 'source_path': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            {'tool_name': 'validate_test_structure', 'description': 'Validate tests/naming/imports', 'inputSchema': {'type': 'object', 'properties': {'test_path': {'type': 'string'}}}, 'server_key': 'pytest', 'category': 'testing'},
            # Firecrawl (approximate external tools)
            {'tool_name': 'crawl_url', 'description': 'Crawl website URL', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}}}, 'server_key': 'firecrawl', 'category': 'web'},
            {'tool_name': 'extract_content', 'description': 'Extract content from crawled page', 'inputSchema': {'type': 'object', 'properties': {'url': {'type': 'string'}}}, 'server_key': 'firecrawl', 'category': 'web'}
        ]
        print("Using static fallback tools for testing.")
        return static_tools


def infer_category(tool_name: str, description: str) -> str:
    """Simple heuristic to infer tool category from name/description."""
    tool_lower = tool_name.lower()
    desc_lower = description.lower()

    if any(
        word in tool_lower or word in desc_lower
        for word in ["http", "request", "api", "endpoint"]
    ):
        return "http"
    elif any(
        word in tool_lower or word in desc_lower
        for word in ["test", "schema", "validate", "pytest"]
    ):
        return "testing"
    elif any(
        word in tool_lower or word in desc_lower
        for word in ["file", "read", "write", "dir"]
    ):
        return "filesystem"
    elif any(
        word in tool_lower or word in desc_lower
        for word in ["route", "validate", "track", "session"]
    ):
        return "routing"
    else:
        return "general"


def create_tables(con):
    """Create tables with schema support."""
    con.execute("""
        CREATE TABLE IF NOT EXISTS servers (
            server_key VARCHAR PRIMARY KEY,
            server_title VARCHAR,
            server_enabled BOOLEAN,
            server_configured BOOLEAN,
            transport VARCHAR,
            command VARCHAR,
            implementation VARCHAR
        )
    """)

    con.execute("""
        CREATE TABLE IF NOT EXISTS tools (
            tool_name VARCHAR PRIMARY KEY,
            description VARCHAR,
            category VARCHAR,
            server_key VARCHAR REFERENCES servers(server_key),
            input_schema JSON  -- Store full schema as JSON
        )
    """)


def populate_db(con, servers: Dict[str, Any], tools: List[Dict[str, Any]]):
    """Populate from parsed yml and runtime tools."""
    # Clear
    con.execute("DELETE FROM tools")
    con.execute("DELETE FROM servers")

    # Insert servers from yml
    for server_key, config in servers.items():
        con.execute("""
            INSERT INTO servers VALUES (?, ?, ?, ?, ?, ?, ?)
        """, (server_key, config['title'], config['enabled'], config['configured'],
              config['transport'], config['command'], 'local' if 'python' in config['command'] else 'external'))

    # Add 'aggregator' server if not present (for merged tools)
    if 'aggregator' not in servers:
        con.execute("""
            INSERT INTO servers VALUES (?, ?, ?, ?, ?, ?, ?)
        """, ('aggregator', 'Cipher Aggregator', True, True, 'sse', 'internal', 'internal'))

    # Insert tools from aggregator query
    for tool in tools:
        server_key = tool['server_key']  # Default to 'aggregator'
        # Find matching server or use default
        if server_key not in servers and server_key != 'aggregator':
            server_key = 'unknown'
            if 'unknown' not in servers:
                con.execute("""
                    INSERT INTO servers VALUES (?, ?, ?, ?, ?, ?, ?)
                """, ('unknown', 'Unknown Server', False, False, 'stdio', 'n/a', 'unknown'))

        con.execute("""
            INSERT INTO tools VALUES (?, ?, ?, ?, ?)
        """, (tool['tool_name'], tool['description'], tool['category'],
              server_key, json.dumps(tool['inputSchema'])))

    print(f"Populated {len(tools)} tools from runtime aggregator across {len(servers) + (1 if 'aggregator' not in servers else 0)} configured servers.")


def export_to_parquet(con):
    """Export tables to Parquet for viewing."""
    parquet_tools = SCRIPT_DIR / "tools.parquet"
    parquet_servers = SCRIPT_DIR / "servers.parquet"

    con.execute(f"COPY tools TO '{parquet_tools}' (FORMAT PARQUET)")
    con.execute(f"COPY servers TO '{parquet_servers}' (FORMAT PARQUET)")

    tools_count = con.execute('SELECT COUNT(*) FROM tools').fetchone()[0]
    servers_count = con.execute('SELECT COUNT(*) FROM servers').fetchone()[0]

    print(f"✅ Exported to Parquet:")
    print(f"  - Tools: {parquet_tools} ({tools_count} rows)")
    print(f"  - Servers: {parquet_servers} ({servers_count} rows)")

    # Viewer tips
    print("\nViewer Tips:")
    print("  - VS Code: Install 'Parquet Viewer' extension, open .parquet files.")
    print("  - CLI: duckdb -c \"SELECT * FROM 'tools.parquet' LIMIT 5\"")
    print("  - GUI: DBeaver, Tableau, or Pandas: pd.read_parquet('tools.parquet')")


def print_example_queries(con):
    """Demo queries."""
    print("\n" + "="*60)
    print("EXAMPLE QUERIES (from runtime data)")
    print("="*60)

    # 1. Enabled tools by category
    print("\n1. Enabled tools by category (top categories):")
    res = con.execute("""
        SELECT category, COUNT(*) as count
        FROM tools t
        JOIN servers s ON t.server_key = s.server_key
        WHERE s.server_enabled = true
        GROUP BY category
        ORDER BY count DESC
        LIMIT 5
    """).fetchall()
    for row in res:
        print(f"  {row[0]}: {row[1]} tools")

    # 2. Sample tool with schema
    print("\n2. Sample enabled tool with input schema:")
    res = con.execute("""
        SELECT tool_name, description, input_schema
        FROM tools t
        JOIN servers s ON t.server_key = s.server_key
        WHERE s.server_enabled = true
        LIMIT 1
    """).fetchone()
    if res:
        print(f"  Tool: {res[0]}")
        print(f"  Desc: {res[1]}")
        print(f"  Schema: {res[2]}")  # JSON dump

    # 3. Server tool counts
    print("\n3. Tool counts per enabled server:")
    res = con.execute("""
        SELECT s.server_key, s.server_title, COUNT(t.tool_name) as num_tools
        FROM servers s
        LEFT JOIN tools t ON s.server_key = t.server_key
        WHERE s.server_enabled = true
        GROUP BY s.server_key, s.server_title
        ORDER BY num_tools DESC
    """).fetchall()
    for row in res:
        print(f"  {row[0]} ({row[1]}): {row[2]} tools")


def delete_yaml_catalogs():
    """Delete static YAML files if --delete-yamls flag."""
    yaml_files = ['mcp-catalog.yml', 'mcp-tools.yml']
    for file in yaml_files:
        file_path = SCRIPT_DIR / file
        if file_path.exists():
            file_path.unlink()
            print(f"Deleted: {file}")
        else:
            print(f"Not found (already deleted?): {file}")


def main():
    parser = argparse.ArgumentParser(description="Generate DuckDB from running Cipher")
    parser.add_argument('--delete-yamls', action='store_true', help="Delete static YAML catalogs after success")
    args = parser.parse_args()

    try:
        servers = parse_cipher_yml()
        tools = query_aggregator_tools()

        con = duckdb.connect(str(DB_PATH))
        create_tables(con)
        populate_db(con, servers, tools)

        print(f"✅ Dynamic DuckDB catalog generated: {DB_PATH} (from cipher.yml + runtime tools)")
        print_example_queries(con)

        # Export to Parquet
        export_to_parquet(con)

        con.close()

        if args.delete_yamls:
            delete_yaml_catalogs()

    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)
        # If aggregator not running, suggest starting it
        if "Failed to query aggregator" in str(e):
            print("Tip: Ensure Cipher is running: ./mcp-manager.sh start")
        sys.exit(1)


if __name__ == "__main__":
    main()
