#!/usr/bin/env python3
import json
import os
import sys
import time

# Configuration from environment
MAX_FILES = int(os.getenv("MAX_FILES", "5"))
MAX_BYTES = int(os.getenv("MAX_BYTES", str(1024 * 1024)))  # 1MB
DELAY_MS = int(os.getenv("DELAY_MS", "300"))
HEAD_BYTES = int(os.getenv("HEAD_BYTES", str(4096)))  # 4KB
TAIL_BYTES = int(os.getenv("TAIL_BYTES", "0"))

def send_response(result):
    """Send JSON-RPC response to stdout"""
    print(json.dumps(result))
    sys.stdout.flush()

def send_error(code, message, data=None):
    """Send JSON-RPC error response"""
    error_response = {
        "jsonrpc": "2.0",
        "error": {
            "code": code,
            "message": message
        }
    }
    if data:
        error_response["error"]["data"] = data
    print(json.dumps(error_response))
    sys.stdout.flush()

def handle_initialize():
    """Handle MCP initialize request"""
    return {
        "jsonrpc": "2.0",
        "id": 1,
        "result": {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "tools": {}
            },
            "serverInfo": {
                "name": "file-batch-mcp",
                "version": "1.0.0"
            }
        }
    }

def handle_list_tools():
    """Handle MCP list_tools request"""
    return {
        "jsonrpc": "2.0",
        "id": 2,
        "result": {
            "tools": [
                {
                    "name": "read_files_batched",
                    "description": "Read files batched. params:paths[]",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "paths": {
                                "type": "array",
                                "items": {"type": "string"},
                                "description": "file paths[]"
                            }
                        },
                        "required": ["paths"]
                    }
                }
            ]
        }
    }

def handle_call_tool(name, arguments):
    """Handle MCP call_tool request"""
    if name != "read_files_batched":
        return send_error(-32601, f"Unknown tool: {name}")

    paths = arguments.get("paths", [])
    if not isinstance(paths, list):
        return send_error(-32602, "Invalid arguments: 'paths' must be an array")

    results = []
    total_bytes = 0
    count = 0

    for path in paths:
        if count >= MAX_FILES or total_bytes >= MAX_BYTES:
            break

        try:
            with open(path, "rb") as f:
                data = f.read()

            # Apply size limits
            truncated = False
            if len(data) > HEAD_BYTES + TAIL_BYTES:
                head = data[:HEAD_BYTES]
                tail = data[-TAIL_BYTES:] if TAIL_BYTES else b""
                data = head + b"\n...[truncated]...\n" + tail
                truncated = True

            results.append({
                "path": path,
                "size": len(data),
                "truncated": truncated,
                "content": data.decode(errors="ignore")
            })

            total_bytes += len(data)
            count += 1

        except Exception as e:
            results.append({
                "path": path,
                "error": str(e)
            })

        # Apply delay between files
        if count < len(paths):
            time.sleep(DELAY_MS / 1000.0)

    return {
        "jsonrpc": "2.0",
        "id": 3,
        "result": {
            "content": [
                {
                    "type": "text",
                    "text": json.dumps({
                        "files": results,
                        "summary": f"Read {count} files, {total_bytes} bytes",
                        "limits": {
                            "max_files": MAX_FILES,
                            "max_bytes": MAX_BYTES,
                            "head_bytes": HEAD_BYTES,
                            "tail_bytes": TAIL_BYTES,
                            "delay_ms": DELAY_MS
                        }
                    })
                }
            ]
        }
    }

def main():
    """Main MCP server loop"""
    for line in sys.stdin:
        try:
            request = json.loads(line.strip())
            method = request.get("method")
            request_id = request.get("id")

            if method == "initialize":
                response = handle_initialize()
                response["id"] = request_id
                send_response(response)

            elif method == "tools/list":
                response = handle_list_tools()
                response["id"] = request_id
                send_response(response)

            elif method == "tools/call":
                tool_name = request["params"]["name"]
                arguments = request["params"]["arguments"]
                response = handle_call_tool(tool_name, arguments)
                if response and "id" in response:
                    response["id"] = request_id
                    send_response(response)

            elif method == "ping":
                send_response({"jsonrpc": "2.0", "id": request_id, "result": "pong"})

            else:
                send_error(-32601, f"Method not found: {method}")

        except json.JSONDecodeError:
            send_error(-32700, "Parse error")
        except Exception as e:
            send_error(-32603, f"Internal error: {str(e)}")

if __name__ == "__main__":
    main()