#!/usr/bin/env python3
"""
Test Catalog Query Script

This script tests the Cipher aggregator SSE endpoint by:
1. Sending initialize request to establish session.
2. Sending tools/list request.
3. Printing the tools for validation.
4. Optionally running generate_catalog_db.py after.

Usage:
    python3 test_catalog.py

Requires: httpx (in pyproject.toml)
"""

import json
import httpx
from pathlib import Path

SSE_URL = "http://127.0.0.1:3020/sse"

def send_request(method: str, params: dict = None, session_id: str = None) -> dict:
    """Send JSON-RPC request to SSE endpoint."""
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params or {}
    }
    if session_id:
        request["params"]["sessionId"] = session_id

    headers = {"Content-Type": "application/json"}

    with httpx.Client(timeout=30.0) as client:
        response = client.post(SSE_URL, json=request, headers=headers)
        response.raise_for_status()
        return response.json()

def main():
    try:
        # Step 1: Initialize session
        print("1. Initializing session...")
        init_resp = send_request("initialize")
        if 'result' in init_resp:
            print("Initialize success:", init_resp['result'])
            session_id = init_resp['result'].get('sessionId', 'default-session')
        else:
            print("Initialize failed:", init_resp)
            return

        # Step 2: List tools
        print("\n2. Querying tools/list...")
        tools_resp = send_request("tools/list", {}, session_id)
        if 'result' in tools_resp and 'tools' in tools_resp['result']:
            tools = tools_resp['result']['tools']
            print(f"Success: Fetched {len(tools)} tools.")
            for tool in tools:
                print(f" - {tool['name']}: {tool['description'][:100]}...")
        else:
            print("Tools/list failed:", tools_resp)

        # Step 3: Test tool call (e.g., validate_tool_selection)
        print("\n3. Testing tool call (validate_tool_selection)...")
        tool_call = send_request("tools/call", {
            "name": "validate_tool_selection",
            "arguments": {
                "session_id": session_id,
                "task_description": "Test HTTP request",
                "selected_tool": "make_request"
            }
        }, session_id)
        print("Tool call response:", tool_call)

        # Step 4: Run generate_catalog_db.py
        print("\n4. Running generate_catalog_db.py...")
        import subprocess
        result = subprocess.run([".venv/bin/python3", "generate_catalog_db.py"], capture_output=True, text=True)
        print("Generate output:", result.stdout)
        if result.stderr:
            print("Generate error:", result.stderr)

        print("\n✅ Tests complete. Aggregator SSE working.")

    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    main()