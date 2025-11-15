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
from typing import Any, Dict, Optional

import httpx

SSE_URL = "http://127.0.0.1:3020/sse"


def send_request(method: str, params: Optional[Dict[str, Any]] = None, session_id: Optional[str] = None) -> Dict[str, Any]:
    """Send JSON-RPC request to the SSE endpoint and parse the SSE response."""
    payload: Dict[str, Any] = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params.copy() if params else {}
    }
    if session_id:
        payload["params"]["sessionId"] = session_id

    headers = {
        "Content-Type": "application/json",
        "Accept": "text/event-stream"
    }

    with httpx.Client(timeout=None) as client:
        with client.stream("POST", SSE_URL, json=payload, headers=headers) as response:
            response.raise_for_status()
            for line in response.iter_text(chunk_size=1024):
                line = line.strip()
                if not line:
                    continue
                if line.startswith("data:"):
                    data = line[len("data:"):].strip()
                    if not data:
                        continue
                    return json.loads(data)
            raise RuntimeError(f"No SSE data received for method {method}")


def main() -> None:
    try:
        print("1. Initializing session...")
        init_resp = send_request("initialize")
        if "result" in init_resp:
            print("Initialize success:", init_resp["result"])
            session_id = init_resp["result"].get("sessionId", "default-session")
        else:
            print("Initialize failed:", init_resp)
            return

        print("\n2. Querying tools/list...")
        tools_resp = send_request("tools/list", {}, session_id)
        if "result" in tools_resp and "tools" in tools_resp["result"]:
            tools = tools_resp["result"]["tools"]
            print(f"Success: Fetched {len(tools)} tools.")
            for tool in tools:
                description = tool.get("description", "")
                print(f" - {tool['name']}: {description[:100]}...")
        else:
            print("Tools/list failed:", tools_resp)

        print("\n3. Testing tool call (validate_tool_selection)...")
        tool_call = send_request(
            "tools/call",
            {
                "name": "validate_tool_selection",
                "arguments": {
                    "session_id": session_id,
                    "task_description": "Test HTTP request",
                    "selected_tool": "make_request"
                }
            },
            session_id
        )
        print("Tool call response:", tool_call)

        print("\n4. Running generate_catalog_db.py...")
        import subprocess

        result = subprocess.run([".venv/bin/python3", "generate_catalog_db.py"], capture_output=True, text=True)
        print("Generate output:", result.stdout)
        if result.stderr:
            print("Generate error:", result.stderr)

        print("\n✅ Tests complete. Aggregator SSE working.")

    except Exception as exc:
        print(f"❌ Error: {exc}")


if __name__ == "__main__":
    main()
