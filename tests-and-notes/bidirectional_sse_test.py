#!/usr/bin/env python3
"""
Bidirectional SSE test for cipher-aggregator
Maintains SSE connection while sending requests
"""

import json
import urllib.error
import urllib.request
import urllib.parse
import threading
import time

class BidirectionalSSETest:
    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.response_event = threading.Event()
        self.response_data = None
        self.pending_request_id = None

    def establish_sse_connection(self):
        """Establish SSE connection and extract session ID"""
        print("🔄 Step 1: Establishing persistent SSE connection...")

        try:
            # Create SSE connection request
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
                        print(f"🎯 Session ID extracted: {self.session_id}")

                # Handle responses (for future use)
                elif line.startswith("event: response"):
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
                        except json.JSONDecodeError:
                            pass

        except Exception as e:
            print(f"SSE stream error: {e}")

    def send_request(self, method, params=None):
        """Send request while maintaining SSE connection"""
        if not self.session_id:
            print("❌ No session ID available")
            return None

        print(f"🔄 Step 2: Sending {method} request...")

        # Generate unique request ID
        import uuid
        request_id = str(uuid.uuid4())
        self.pending_request_id = request_id

        request_data = {
            "jsonrpc": "2.0",
            "method": method,
            "id": request_id,
            "params": params or {}
        }

        # Add session ID as URL parameter
        url_with_session = f"{self.sse_url}?sessionId={self.session_id}"

        try:
            req = urllib.request.Request(
                url_with_session,
                data=json.dumps(request_data).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST",
            )

            print(f"POST to: {url_with_session}")
            print(f"Request: {json.dumps(request_data, indent=2)}")

            # Make request while SSE connection is still active
            with urllib.request.urlopen(req, timeout=30) as response:
                print(f"Status: {response.status}")

                if response.status == 200:
                    # Direct response
                    response_data = response.read().decode("utf-8")
                    print(f"Response: {response_data}")
                    try:
                        parsed_response = json.loads(response_data)
                        print(f"Parsed response: {json.dumps(parsed_response, indent=2)}")
                        return parsed_response
                    except json.JSONDecodeError:
                        print("⚠️  Response is not valid JSON")
                        return None

                elif response.status == 202:
                    # Response will come via SSE stream
                    print("📨 202 Accepted - waiting for response via SSE stream...")

                    # Wait for response from SSE stream
                    if self.response_event.wait(timeout=30):
                        print(f"✅ Response received via SSE: {json.dumps(self.response_data, indent=2)}")
                        return self.response_data
                    else:
                        print("❌ No response received from SSE stream within timeout")
                        return None
                else:
                    print(f"❌ Unexpected status code: {response.status}")
                    return None

        except urllib.error.HTTPError as e:
            print(f"❌ HTTP Error {e.code}: {e.reason}")
            error_body = e.read().decode("utf-8") if hasattr(e, 'read') else "No error body"
            print(f"Error body: {error_body}")
            return None
        except Exception as e:
            print(f"❌ Request error: {e}")
            return None
        finally:
            # Clear pending request
            self.pending_request_id = None
            self.response_event.clear()
            self.response_data = None

    def test_tools_list(self):
        """Test tools/list with bidirectional SSE"""

        # Step 1: Establish persistent SSE connection
        if not self.establish_sse_connection():
            return False, []

        # Step 2: Send tools/list request
        result = self.send_request("tools/list")

        if result and "result" in result and "tools" in result["result"]:
            tools = result["result"]["tools"]
            print(f"\n✅ SUCCESS: Found {len(tools)} tools")

            # Show some examples
            print("\nSample tools:")
            for i, tool in enumerate(tools[:5]):
                name = tool.get("name", "Unknown")
                desc = tool.get("description", "No description")
                print(f"  {i+1}. {name}: {desc}")

            if len(tools) > 5:
                print(f"  ... and {len(tools) - 5} more tools")

            return True, tools
        else:
            print("❌ Failed to get tools list")
            return False, []

def main():
    print("=== Bidirectional SSE Test ===")

    tester = BidirectionalSSETest()
    success, tools = tester.test_tools_list()

    if success:
        print("\n🎉 Bidirectional SSE test PASSED!")
        print(f"Total tools available: {len(tools)}")
    else:
        print("\n❌ Bidirectional SSE test FAILED!")

    print("\n=== Test Complete ===")

if __name__ == "__main__":
    main()
