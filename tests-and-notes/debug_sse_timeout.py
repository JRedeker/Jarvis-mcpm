#!/usr/bin/env python3
"""
Debug SSE stream timeout issue
"""

import json
import urllib.request
import urllib.error
import threading
import time

def debug_sse_timeout():
    """Debug the SSE stream timeout issue"""
    print("=== Debugging SSE Stream Timeout ===")

    # Step 1: Establish SSE connection
    print("🔄 Establishing SSE connection...")
    try:
        req = urllib.request.Request(
            "http://localhost:3020/sse",
            method="GET",
            headers={"Accept": "text/event-stream"}
        )
        response = urllib.request.urlopen(req, timeout=10)

        if response.status == 200:
            print("✅ SSE connection established")

            # Extract session ID from first few lines
            session_id = None
            lines_received = 0

            def read_session_id():
                nonlocal session_id, lines_received
                try:
                    for line in response:
                        if not line:
                            break
                        lines_received += 1
                        decoded_line = line.decode("utf-8").strip()
                        print(f"[{lines_received}] SSE: {decoded_line}")

                        if decoded_line.startswith("data: "):
                            data = decoded_line[6:].strip()
                            if data.startswith("/sse?sessionId="):
                                session_id = data.split("sessionId=")[1]
                                print(f"🎯 Session ID found: {session_id}")
                                return
                        elif decoded_line == "":
                            continue
                        else:
                            print(f"Other SSE data: {decoded_line}")

                        if lines_received > 20:  # Limit debugging output
                            break

                except Exception as e:
                    print(f"SSE read error: {e}")

            # Start reading thread
            read_thread = threading.Thread(target=read_session_id, daemon=True)
            read_thread.start()

            # Wait for session ID
            start_time = time.time()
            while time.time() - start_time < 10 and not session_id:
                time.sleep(0.1)

            if not session_id:
                print("❌ No session ID received")
                return False

            # Step 2: Send tools/list request
            print("📤 Sending tools/list request...")
            request_data = {
                "jsonrpc": "2.0",
                "method": "tools/list",
                "id": "debug-001",
                "params": {}
            }

            url_with_session = f"http://localhost:3020/sse?sessionId={session_id}"
            print(f"POST to: {url_with_session}")

            try:
                post_req = urllib.request.Request(
                    url_with_session,
                    data=json.dumps(request_data).encode("utf-8"),
                    headers={"Content-Type": "application/json"},
                    method="POST"
                )

                post_response = urllib.request.urlopen(post_req, timeout=10)
                print(f"📥 POST Response status: {post_response.status}")

                if post_response.status == 202:
                    print("📨 202 Accepted - checking SSE stream for response...")

                    # Try to read from the original SSE stream
                    end_time = time.time() + 10
                    responses_found = 0

                    while time.time() < end_time:
                        try:
                            # Try to read with timeout
                            response_line = response.readline()
                            if response_line:
                                decoded = response_line.decode("utf-8").strip()
                                print(f"📡 Response line: {decoded}")
                                responses_found += 1
                            else:
                                print("No more data from SSE stream")
                                break
                        except Exception as e:
                            print(f"SSE read error during response: {e}")
                            break

                    print(f"📊 Total SSE responses received: {responses_found}")

                    if responses_found == 0:
                        print("❌ No responses received from SSE stream")
                        return False
                    else:
                        print("✅ Some responses were received from SSE stream")
                        return True

            except urllib.error.HTTPError as e:
                print(f"❌ POST HTTP Error {e.code}: {e.reason}")
                return False
            except Exception as e:
                print(f"❌ POST Error: {e}")
                return False

        else:
            print(f"❌ SSE connection failed: {response.status}")
            return False

    except Exception as e:
        print(f"❌ SSE connection error: {e}")
        return False

if __name__ == "__main__":
    debug_sse_timeout()
