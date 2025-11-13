#!/usr/bin/env python3
"""
Simple test to check if cipher memory tools are available
"""

import os
import sys
import json
import asyncio
from pathlib import Path

# Add project root to path for imports
sys.path.insert(0, str(Path(__file__).parent))

# Load environment variables from .env file
try:
    from dotenv import load_dotenv
    load_dotenv()
except ImportError:
    print("⚠️  Warning: python-dotenv not available, using manual env loading")
    # Fallback to manual environment loading
    def load_env():
        env_vars = {}
        try:
            with open('.env', 'r') as f:
                for line in f:
                    if '=' in line and not line.strip().startswith('#'):
                        key, value = line.strip().split('=', 1)
                        env_vars[key.strip()] = value.strip()
        except FileNotFoundError:
            print("⚠️  Warning: .env file not found")
        return env_vars

async def test_cipher_memory_tools():
    """Test if cipher memory tools are available"""
    print("🔍 Testing Cipher Memory Tools Availability...")

    try:
        # Simple HTTP request to tools/list
        import urllib.request
        import urllib.error

        req = urllib.request.Request(
            "http://localhost:3020/sse",
            method="GET",
            headers={"Accept": "text/event-stream"}
        )

        with urllib.request.urlopen(req, timeout=10) as response:
            if response.status == 200:
                # Look for session ID in SSE stream
                for line in response:
                    line = line.decode('utf-8').strip()
                    if line.startswith("data: "):
                        data = line[6:].strip()
                        if data.startswith("/sse?sessionId="):
                            session_id = data.split("sessionId=")[1]
                            print(f"✅ Session ID: {session_id}")

                            # Now test tools list
                            tools_req = urllib.request.Request(
                                f"http://localhost:3020/sse?sessionId={session_id}",
                                method="POST",
                                headers={"Content-Type": "application/json"},
                                data=json.dumps({
                                    "jsonrpc": "2.0",
                                    "method": "tools/list",
                                    "id": 1,
                                    "params": {}
                                }).encode("utf-8")
                            )

                            with urllib.request.urlopen(tools_req, timeout=10) as tools_response:
                                if tools_response.status == 200:
                                    tools_data = tools_response.read().decode("utf-8")
                                    tools = json.loads(tools_data)

                                    # Check for cipher memory tools
                                    expected_tools = [
                                        "cipher_memory_search",
                                        "cipher_workspace_search",
                                        "cipher_extract_and_operate_memory",
                                        "cipher_store_reasoning_memory",
                                        "cipher_search_reasoning_patterns",
                                        "cipher_workspace_store",
                                        "cipher_extract_entities",
                                        "cipher_intelligent_processor"
                                    ]

                                    tool_names = [tool.get("name", "") for tool in tools.get("tools", [])]

                                    print(f"📋 Found {len(tool_names)} total tools")

                                    found_memory_tools = []
                                    missing_memory_tools = []

                                    for tool_name in expected_tools:
                                        if tool_name in tool_names:
                                            found_memory_tools.append(tool_name)
                                            print(f"  ✅ {tool_name}")
                                        else:
                                            missing_memory_tools.append(tool_name)
                                            print(f"  ❌ {tool_name} - MISSING")

                                    print(f"\n📊 Memory Tools Summary:")
                                    print(f"  Found: {len(found_memory_tools)}/{len(expected_tools)}")
                                    print(f"  Missing: {len(missing_memory_tools)}")

                                    if missing_memory_tools:
                                        print(f"  ⚠️  Missing tools: {', '.join(missing_memory_tools)}")
                                        return False
                                    else:
                                        print("✅ All cipher memory tools are available!")
                                        return True
                                else:
                                    print(f"❌ Failed to get tools list: {tools_response.status}")
                                    return False
                        else:
                            continue

                print("❌ No session ID received from SSE stream")
                return False

            else:
                print(f"❌ SSE connection failed: {response.status}")
                return False

    except Exception as e:
        print(f"❌ Error: {e}")
        return False


async def main():
    """Main execution function"""
    print("🧠 Simple Cipher Memory Tools Test")

    success = await test_cipher_memory_tools()

    if success:
        print("\n✅ Cipher memory tools are available!")
        print("📝 Ready to proceed with Phase 3 implementation")
    else:
        print("\n❌ Cipher memory tools test failed")

    return 0 if success else 1


if __name__ == "__main__":
    exit_code = asyncio.run(main())
    sys.exit(exit_code)
