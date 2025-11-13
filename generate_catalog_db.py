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
    python generate_catalog_db.py [--delete-yamls]

View Parquet:
- VS Code: 'Parquet Viewer' extension.
- CLI: duckdb -c "SELECT * FROM 'tools.parquet' LIMIT 5"
- Python: import pandas as pd; pd.read_parquet('tools.parquet')
"""

import argparse
import json
import sys
from pathlib import Path
from typing import Any, Dict, List

import duckdb
import httpx
import yaml

# Paths
SCRIPT_DIR = Path(__file__).parent
CIPHER_YML = SCRIPT_DIR / "cipher.yml"
DB_PATH = SCRIPT_DIR / "catalog.db"
SSE_URL = "http://127.0.0.1:3020/sse"


def parse_cipher_yml() -> Dict[str, Any]:
    """Parse cipher.yml for server configurations."""
    if not CIPHER_YML.exists():
        raise FileNotFoundError(f"cipher.yml not found: {CIPHER_YML}")

    with open(CIPHER_YML, "r") as f:
        data = yaml.safe_load(f)

    servers = data.get("servers", {})
    parsed_servers = {}

    for server_key, config in servers.items():
        parsed_servers[server_key] = {
            "enabled": config.get("enabled", False),
            "transport": config.get("transport", "stdio"),
            "command": config.get("command", ""),
            "env": config.get("env", {}),
            "title": config.get(
                "title", server_key.replace("-", " ").title()
            ),  # Infer title
            "configured": True,  # Present in yml
        }

    return parsed_servers


def query_aggregator_tools() -> List[Dict[str, Any]]:
    """Query the running Cipher aggregator for tools/list via SSE JSON-RPC."""
    # JSON-RPC request for tools/list
    request = {"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}

    headers = {"Content-Type": "text/plain", "Accept": "text/event-stream"}

    try:
        with httpx.Client(timeout=30.0) as client:
            # POST the request to SSE endpoint
            response = client.post(SSE_URL, json=request, headers=headers)
            response.raise_for_status()

            # Parse response
            if response.headers.get("content-type", "").startswith("text/event-stream"):
                # Handle SSE stream
                tools_data = []
                for line in response.iter_lines():
                    if line:
                        if line.startswith("data: "):
                            event_data = line[6:].strip()
                            if event_data == "[DONE]":
                                break
                            try:
                                event = json.loads(event_data)
                                if "result" in event and "tools" in event["result"]:
                                    tools_data = event["result"]["tools"]
                                    break  # Got the tools
                            except json.JSONDecodeError:
                                continue
            else:
                # Assume direct JSON response
                resp_json = response.json()
                if "result" in resp_json and "tools" in resp_json["result"]:
                    tools_data = resp_json["result"]["tools"]
                else:
                    raise ValueError("Unexpected response format")

            # Parse tools
            parsed_tools = []
            for tool in tools_data:
                parsed_tools.append(
                    {
                        "tool_name": tool["name"],
                        "description": tool["description"],
                        "inputSchema": tool.get("inputSchema", {}),
                        "server_key": "aggregator",  # Merged; server per-tool not directly available
                        "category": infer_category(
                            tool["name"], tool["description"]
                        ),  # Simple inference
                    }
                )

            return parsed_tools

    except httpx.RequestError as e:
        raise RuntimeError(f"Failed to query aggregator at {SSE_URL}: {e}")
    except ValueError as e:
        raise RuntimeError(f"Failed to parse aggregator response: {e}")


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

    # Insert tools from aggregator query
    for tool in tools:
        server_key = tool['server_key']  # Default to 'aggregator'; could map if per-server query
        # Find matching server or use default
        if server_key not in servers:
            server_key = 'unknown'

        con.execute("""
            INSERT INTO tools VALUES (?, ?, ?, ?, ?)
        """, (tool['tool_name'], tool['description'], tool['category'],
              server_key, json.dumps(tool['inputSchema'])))

    print(f"Populated {len(tools)} tools from runtime aggregator across {len(servers)} configured servers.")


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
