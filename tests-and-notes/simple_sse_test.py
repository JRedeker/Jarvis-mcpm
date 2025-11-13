#!/usr/bin/env python3
"""
Simple SSE test for cipher-aggregator
Tests tools/list endpoint without external dependencies
"""

import json
import urllib.error
import urllib.request
import uuid
from datetime import datetime


def generate_session_id():
    """Generate a unique session ID"""
    return f"session-{int(datetime.now().timestamp())}-{uuid.uuid4().hex[:8]}"


def test_tools_list():
    """Test the tools/list endpoint using proper MCP SSE flow"""

    print("🔄 Step 1: Establishing SSE connection...")

    # Step 1: Connect to SSE to get session ID
    sse_url = "http://localhost:3020/sse"
    session_id = None

    try:
        # Create SSE connection request
        req = urllib.request.Request(
            sse_url,
            method="GET",
            headers={"Accept": "text/event-stream"}
        )

        with urllib.request.urlopen(req, timeout=10) as response:
            if response.status != 200:
                print(f"❌ SSE connection failed: {response.status}")
                return False, []

            print("✅ SSE connection established")

            # Read SSE stream to get session ID
            session_id = None
            for line in response:
                line = line.decode("utf-8").strip()
                if line.startswith("data: "):
                    data = line[6:].strip()
                    if data.startswith("/sse?sessionId="):
                        session_id = data.split("sessionId=")[1]
                        print(f"✅ Session ID received: {session_id}")
                        break

            if not session_id:
                print("❌ No session ID received from SSE stream")
                return False, []

    except Exception as e:
        print(f"❌ SSE connection error: {e}")
        return False, []

    print("🔄 Step 2: Sending tools/list request...")

    # Step 2: Send tools/list request with session ID as URL parameter
    request_data = {
        "jsonrpc": "2.0",
        "method": "tools/list",
        "id": 1,
        "params": {}
    }

    # Use session ID as URL parameter (not in request body)
    url_with_session = f"{sse_url}?sessionId={session_id}"

    try:
        req = urllib.request.Request(
            url_with_session,
            data=json.dumps(request_data).encode("utf-8"),
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        print(f"POST to: {url_with_session}")
        print(f"Request: {json.dumps(request_data, indent=2)}")

        # Make the request
        with urllib.request.urlopen(req, timeout=30) as response:
            response_data = response.read().decode("utf-8")
            print(f"Status: {response.status}")
            print(f"Response: {response_data}")

            # Try to parse as JSON
            try:
                parsed_response = json.loads(response_data)
                print(f"Parsed response: {json.dumps(parsed_response, indent=2)}")

                # Check if we got tools
                if "result" in parsed_response and "tools" in parsed_response["result"]:
                    tools = parsed_response["result"]["tools"]
                    print(f"\n✅ SUCCESS: Found {len(tools)} tools")

                    # Look for built-in AI features
                    built_in_tools = []
                    for tool in tools:
                        tool_name = tool.get("name", "")
                        if (
                            "cipher_memory" in tool_name
                            or "cipher_workspace" in tool_name
                        ):
                            built_in_tools.append(tool_name)

                    if built_in_tools:
                        print(f"✅ Built-in AI tools found: {built_in_tools}")
                    else:
                        print("⚠️  No built-in AI tools detected")

                    return True, tools
                else:
                    print("⚠️  No tools found in response")
                    return False, []

            except json.JSONDecodeError:
                print("⚠️  Response is not valid JSON")
                return False, []

    except urllib.error.HTTPError as e:
        print(f"❌ HTTP Error {e.code}: {e.reason}")
        error_body = e.read().decode("utf-8") if hasattr(e, "read") else "No error body"
        print(f"Error body: {error_body}")
        return False, []
    except urllib.error.URLError as e:
        print(f"❌ URL Error: {e.reason}")
        return False, []
    except Exception as e:
        print(f"❌ Unexpected error: {e}")
        return False, []


if __name__ == "__main__":
    print("=== Cipher-Aggregator SSE Tools List Test ===")
    success, tools = test_tools_list()

    if success:
        print("\n✅ Server is responding correctly!")
        print(f"Total tools available: {len(tools)}")

        # Show some tool examples
        print("\nSample tools:")
        for i, tool in enumerate(tools[:10]):  # Show first 10 tools
            name = tool.get("name", "Unknown")
            desc = tool.get("description", "No description")
            print(f"  {i + 1}. {name}: {desc}")

        if len(tools) > 10:
            print(f"  ... and {len(tools) - 10} more tools")
    else:
        print("\n❌ Server test failed!")

    print("\n=== Test Complete ===")
