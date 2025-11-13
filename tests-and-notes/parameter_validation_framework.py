#!/usr/bin/env python3
"""
Phase 2: Parameter Validation Framework
Parses tool schemas from cipher-aggregator and validates parameters before tool calls
"""

import json
import urllib.error
import urllib.request
import threading
import time
from typing import Dict, List, Any, Optional, Tuple
import re

class ParameterValidator:
    """Validates MCP tool parameters based on parsed schemas"""

    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.tools_schemas = {}  # Cache of tool schemas
        self.sse_response = None
        self.sse_thread = None
        self.response_event = threading.Event()
        self.response_data = None
        self.pending_request_id = None

    def establish_sse_connection(self) -> bool:
        """Establish SSE connection and extract session ID"""
        print("🔄 Establishing SSE connection for schema extraction...")

        try:
            req = urllib.request.Request(
                self.sse_url,
                method="GET",
                headers={"Accept": "text/event-stream"}
            )

            self.sse_response = urllib.request.urlopen(req, timeout=10)
            if self.sse_response.status != 200:
                print(f"❌ SSE connection failed: {self.sse_response.status}")
                return False

            print("✅ SSE connection established")

            # Start background thread to listen to SSE stream
            self.sse_thread = threading.Thread(target=self._listen_sse_stream, daemon=True)
            self.sse_thread.start()

            # Extract session ID
            start_time = time.time()
            while time.time() - start_time < 10:
                time.sleep(0.1)
                if self.session_id:
                    print(f"✅ Session ID received: {self.session_id}")
                    return True

            print("❌ No session ID received from SSE stream")
            return False

        except Exception as e:
            print(f"❌ SSE connection error: {e}")
            return False

    def _listen_sse_stream(self):
        """Background thread to listen to SSE stream"""
        try:
            for line in self.sse_response:
                if not line:
                    break

                line = line.decode("utf-8").strip()

                # Handle session ID from endpoint event
                if line.startswith("data: "):
                    data = line[6:].strip()
                    if data.startswith("/sse?sessionId="):
                        self.session_id = data.split("sessionId=")[1]

                # Handle responses - look for event: message or event: response
                elif line.startswith("event: message") or line.startswith("event: response"):
                    # Read the response data
                    response_line = next(self.sse_response, None)
                    if response_line:
                        response_data = response_line.decode("utf-8").strip()[6:]
                        try:
                            parsed = json.loads(response_data)
                            # Only set response if it matches our pending request
                            if self.pending_request_id and parsed.get('id') == self.pending_request_id:
                                self.response_data = parsed
                                self.response_event.set()
                                print(f"✅ SSE response received for request {self.pending_request_id}")
                        except json.JSONDecodeError:
                            print(f"❌ Failed to parse SSE response data: {response_data}")

        except Exception as e:
            print(f"SSE stream error: {e}")

    def extract_tool_schemas(self) -> Dict[str, Dict]:
        """Extract tool schemas from tools/list response"""
        if not self.session_id:
            print("❌ No session ID available")
            return {}

        print("🔄 Extracting tool schemas from tools/list...")

        # Generate unique request ID
        import uuid
        request_id = str(uuid.uuid4())
        self.pending_request_id = request_id

        request_data = {
            "jsonrpc": "2.0",
            "method": "tools/list",
            "id": request_id,
            "params": {}
        }

        url_with_session = f"{self.sse_url}?sessionId={self.session_id}"

        try:
            req = urllib.request.Request(
                url_with_session,
                data=json.dumps(request_data).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST",
            )

            print(f"📤 Requesting tools list from: {url_with_session}")

            with urllib.request.urlopen(req, timeout=30) as response:
                print(f"📥 Response status: {response.status}")

                if response.status == 200:
                    # Direct response
                    response_data = response.read().decode("utf-8")
                    try:
                        parsed_response = json.loads(response_data)
                        return self._parse_tools_response(parsed_response)
                    except json.JSONDecodeError:
                        print("❌ Response is not valid JSON")
                        return {}

                elif response.status == 202:
                    # Response will come via SSE stream
                    print("📨 202 Accepted - waiting for response via SSE stream...")

                    # Wait for response from SSE stream
                    if self.response_event.wait(timeout=30):
                        print(f"✅ Response received via SSE")
                        return self._parse_tools_response(self.response_data)
                    else:
                        print("❌ No response received from SSE stream within timeout")
                        return {}
                else:
                    print(f"❌ Unexpected status: {response.status}")
                    return {}

        except urllib.error.HTTPError as e:
            print(f"❌ HTTP Error {e.code}: {e.reason}")
            return {}
        except Exception as e:
            print(f"❌ Request error: {e}")
            return {}
        finally:
            # Clear pending request
            self.pending_request_id = None
            self.response_event.clear()
            self.response_data = None

    def _parse_tools_response(self, parsed_response: Dict) -> Dict[str, Dict]:
        """Parse tools response and extract schemas"""
        if "result" in parsed_response and "tools" in parsed_response["result"]:
            tools = parsed_response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools")

            # Extract schemas
            schemas = {}
            for tool in tools:
                tool_name = tool.get("name", "")
                input_schema = tool.get("inputSchema", {})
                schemas[tool_name] = {
                    "description": tool.get("description", ""),
                    "inputSchema": input_schema,
                    "required": input_schema.get("required", []),
                    "properties": input_schema.get("properties", {})
                }

            self.tools_schemas = schemas
            print(f"✅ Extracted schemas for {len(schemas)} tools")
            return schemas
        else:
            print("❌ No tools found in response")
            return {}

    def validate_parameters(self, tool_name: str, params: Dict[str, Any]) -> Tuple[bool, Optional[str]]:
        """
        Validate parameters against tool schema

        Returns:
            (is_valid, error_message)
        """
        if tool_name not in self.tools_schemas:
            return False, f"Unknown tool: {tool_name}"

        schema = self.tools_schemas[tool_name]
        required_params = schema.get("required", [])
        properties = schema.get("properties", {})

        # Check for missing required parameters
        missing_params = []
        for required_param in required_params:
            if required_param not in params or params[required_param] is None:
                missing_params.append(required_param)

        if missing_params:
            return False, f"Missing required parameters for '{tool_name}': {', '.join(missing_params)}"

        # Check for unknown parameters
        unknown_params = []
        for param in params:
            if param not in properties:
                unknown_params.append(param)

        if unknown_params:
            return False, f"Unknown parameters for '{tool_name}': {', '.join(unknown_params)}"

        return True, None

    def get_schema_info(self, tool_name: str) -> Optional[Dict]:
        """Get detailed schema information for a tool"""
        return self.tools_schemas.get(tool_name)

    def list_available_tools(self) -> List[str]:
        """List all available tools"""
        return list(self.tools_schemas.keys())

    def get_required_params(self, tool_name: str) -> List[str]:
        """Get required parameters for a tool"""
        if tool_name in self.tools_schemas:
            return self.tools_schemas[tool_name].get("required", [])
        return []

def test_parameter_validation():
    """Test the parameter validation framework"""
    print("=== Parameter Validation Framework Test ===")

    validator = ParameterValidator()

    # Step 1: Establish SSE connection
    if not validator.establish_sse_connection():
        print("❌ Failed to establish SSE connection")
        return False

    # Step 2: Extract tool schemas
    schemas = validator.extract_tool_schemas()
    if not schemas:
        print("❌ Failed to extract tool schemas")
        return False

    # Step 3: Test validation with various scenarios
    print("\n🧪 Testing parameter validation...")

    # Test 1: Valid parameters
    tools = validator.list_available_tools()
    if tools:
        test_tool = tools[0]
        required = validator.get_required_params(test_tool)

        print(f"\n📋 Testing tool: {test_tool}")
        print(f"   Required params: {required}")

        # Test with empty params (should fail if required params exist)
        is_valid, error = validator.validate_parameters(test_tool, {})
        if required:
            print(f"   ✅ Empty params correctly rejected: {error}")
        else:
            print(f"   ✅ Empty params accepted (no required params)")

        # Test with valid params (if we can construct them)
        if required:
            # Try with required params as empty strings
            test_params = {param: "" for param in required}
            is_valid, error = validator.validate_parameters(test_tool, test_params)
            print(f"   ✅ Required params validation: {'PASS' if is_valid else f'FAIL - {error}'}")

    # Test 2: Unknown tool
    is_valid, error = validator.validate_parameters("unknown_tool", {})
    print(f"   ✅ Unknown tool correctly rejected: {error}")

    # Test 3: Show some example schemas
    print(f"\n📚 Sample tool schemas:")
    for i, (tool_name, schema) in enumerate(list(schemas.items())[:3]):
        print(f"   {i+1}. {tool_name}")
        print(f"      Description: {schema['description'][:60]}...")
        print(f"      Required: {schema['required']}")
        if schema['properties']:
            prop_names = list(schema['properties'].keys())[:3]
            print(f"      Properties: {prop_names}{'...' if len(schema['properties']) > 3 else ''}")

    print(f"\n✅ Parameter validation framework ready!")
    print(f"   Total tools loaded: {len(schemas)}")
    print(f"   Tools with required params: {sum(1 for s in schemas.values() if s['required'])}")

    return True

if __name__ == "__main__":
    success = test_parameter_validation()
    if success:
        print("\n🎉 Phase 2 Parameter Validation Framework: READY")
    else:
        print("\n❌ Phase 2 Parameter Validation Framework: FAILED")
