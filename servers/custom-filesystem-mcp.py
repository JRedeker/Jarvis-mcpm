#!/usr/bin/env python3
import json
import os
import sys
from typing import Any, Dict, List, Optional


def send_response(result):
    """Send JSON-RPC response to stdout"""
    print(json.dumps(result))
    sys.stdout.flush()


def send_error(code: int, message: str, data: Any = None) -> None:
    """Send JSON-RPC error response"""
    error_response: Dict[str, Any] = {
        "jsonrpc": "2.0",
        "error": {"code": code, "message": message},
    }
    if data is not None:
        error_response["error"]["data"] = data
    print(json.dumps(error_response))
    sys.stdout.flush()


def read_file_content(
    path: str, head: Optional[int] = None, tail: Optional[int] = None
) -> str:
    """Read file content with optional head/tail limits"""
    try:
        with open(path, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()

        if tail:
            lines = content.split("\n")
            content = "\n".join(lines[-tail:])
        elif head:
            lines = content.split("\n")
            content = "\n".join(lines[:head])

        return content
    except Exception as e:
        raise Exception(f"Failed to read {path}: {str(e)}")


def list_tools():
    """Return list of available tools (modern ones only)"""
    return [
        {
            "name": "read_text_file",
            "description": "Read file text. params:path,head,tail",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "path": {"type": "string"},
                    "tail": {"type": "number", "description": "last N lines"},
                    "head": {"type": "number", "description": "first N lines"},
                },
                "required": ["path"],
                "additionalProperties": False,
            },
        },
        {
            "name": "read_multiple_files",
            "description": "Read multiple files. params:paths[]",
            "inputSchema": {
                "type": "object",
                "properties": {"paths": {"type": "array", "items": {"type": "string"}}},
                "required": ["paths"],
                "additionalProperties": False,
            },
        },
        {
            "name": "write_file",
            "description": "Write file (overwrite). params:path,content",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "path": {"type": "string"},
                    "content": {"type": "string"},
                },
                "required": ["path", "content"],
                "additionalProperties": False,
            },
        },
        {
            "name": "list_directory",
            "description": "List files/dirs. params:path",
            "inputSchema": {
                "type": "object",
                "properties": {"path": {"type": "string"}},
                "required": ["path"],
                "additionalProperties": False,
            },
        },
        {
            "name": "create_directory",
            "description": "Create directory. params:path",
            "inputSchema": {
                "type": "object",
                "properties": {"path": {"type": "string"}},
                "required": ["path"],
                "additionalProperties": False,
            },
        },
    ]


def handle_tool_call(name: str, arguments: Dict[str, Any]):
    """Handle MCP call_tool request"""
    try:
        if name == "read_text_file":
            path = arguments["path"]
            head = arguments.get("head")
            tail = arguments.get("tail")

            content = read_file_content(path, head, tail)

            return {"content": [{"type": "text", "text": content}]}

        elif name == "read_multiple_files":
            paths = arguments["paths"]
            results = []

            for path in paths:
                try:
                    content = read_file_content(path)
                    results.append({"path": path, "content": content})
                except Exception as e:
                    results.append({"path": path, "error": str(e)})

            return {
                "content": [{"type": "text", "text": json.dumps(results, indent=2)}]
            }

        elif name == "write_file":
            path = arguments["path"]
            content = arguments["content"]

            # Ensure directory exists
            os.makedirs(os.path.dirname(path), exist_ok=True)

            with open(path, "w", encoding="utf-8") as f:
                f.write(content)

            return {
                "content": [
                    {
                        "type": "text",
                        "text": f"Successfully wrote {len(content)} characters to {path}",
                    }
                ]
            }

        elif name == "list_directory":
            path = arguments["path"]
            items: List[str] = []

            try:
                for item in os.listdir(path):
                    item_path = os.path.join(path, item)
                    if os.path.isdir(item_path):
                        items.append(f"[DIR] {item}")
                    else:
                        items.append(f"[FILE] {item}")

                return {"content": [{"type": "text", "text": "\n".join(items)}]}
            except Exception as e:
                raise Exception(f"Failed to list directory {path}: {str(e)}")

        elif name == "create_directory":
            path = arguments["path"]
            os.makedirs(path, exist_ok=True)

            return {"content": [{"type": "text", "text": f"Directory created: {path}"}]}

        else:
            send_error(-32601, f"Unknown tool: {name}")
            return

    except Exception as e:
        send_error(-32603, f"Tool execution failed: {str(e)}")
        return


def read_request() -> Optional[Dict[str, Any]]:
    """Read MCP Content-Length framed or newline-delimited JSON requests."""
    while True:
        header = sys.stdin.readline()
        if not header:
            return None
        if header.startswith("Content-Length"):
            try:
                length = int(header.split(":", 1)[1].strip())
            except Exception:
                continue
            sys.stdin.readline()  # blank line
            body = sys.stdin.read(length)
            if not body:
                return None
            try:
                return json.loads(body)
            except json.JSONDecodeError:
                send_error(-32700, "Invalid JSON")
                continue
        else:
            header = header.strip()
            if not header:
                continue
            try:
                return json.loads(header)
            except json.JSONDecodeError:
                send_error(-32700, "Invalid JSON")
                continue


def main():
    """Main MCP server loop"""
    while True:
        request = read_request()
        if request is None:
            break
        try:
            method = request.get("method")
            request_id = request.get("id")

            if method == "initialize":
                send_response(
                    {
                        "jsonrpc": "2.0",
                        "id": request_id,
                        "result": {
                            "protocolVersion": "2024-11-05",
                            "capabilities": {
                                "tools": {"listChanged": False},
                                "resources": {"listChanged": False},
                                "prompts": {"listChanged": False},
                            },
                            "serverInfo": {
                                "name": "custom-filesystem-mcp",
                                "version": "1.0.0",
                            },
                        },
                    }
                )

            elif method == "tools/list":
                send_response(
                    {
                        "jsonrpc": "2.0",
                        "id": request_id,
                        "result": {"tools": list_tools()},
                    }
                )

            elif method == "tools/call":
                name = request["params"]["name"]
                arguments = request["params"]["arguments"]
                result = handle_tool_call(name, arguments)
                if result:
                    send_response(
                        {"jsonrpc": "2.0", "id": request_id, "result": result}
                    )
            elif method in ("resources/list", "prompts/list"):
                key = "resources" if method == "resources/list" else "prompts"
                send_response({"jsonrpc": "2.0", "id": request_id, "result": {key: []}})
            elif method == "notifications/message":
                send_response({"jsonrpc": "2.0", "id": request_id, "result": {}})
            else:
                send_error(-32601, f"Unknown method: {method}")

        except json.JSONDecodeError:
            continue
        except KeyboardInterrupt:
            break
        except Exception as e:
            send_error(-32603, f"Server error: {str(e)}")


if __name__ == "__main__":
    main()
