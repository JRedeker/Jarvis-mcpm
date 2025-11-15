#!/usr/bin/env python3
"""
Httpie MCP Server - Provides HTTP client capabilities through MCP protocol
"""

import asyncio
import json
import logging
import sys
import tempfile
import time
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional

# Import the middleware for agent logging
sys.path.append("/home/jrede/dev/MCP/tests-and-notes")
from cipher_routing_middleware import CipherRoutingMiddleware

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize middleware
middleware = CipherRoutingMiddleware()


# JSON-RPC response helpers
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


def list_tools():
    """Return list of available Httpie tools"""
    return [
        {
            "name": "make_request",
            "description": "HTTP request. params:method,url,headers,data,query,auth,timeout",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "method": {
                        "type": "string",
                        "description": "HTTP method",
                        "enum": [
                            "GET",
                            "POST",
                            "PUT",
                            "DELETE",
                            "PATCH",
                            "HEAD",
                            "OPTIONS",
                        ],
                        "default": "GET",
                    },
                    "url": {"type": "string", "description": "request URL"},
                    "headers": {"type": "object", "description": "custom headers"},
                    "data": {"type": "object", "description": "request data"},
                    "query_params": {"type": "object", "description": "query params"},
                    "auth": {
                        "type": "object",
                        "description": "auth config",
                        "properties": {
                            "type": {
                                "type": "string",
                                "enum": ["basic", "bearer", "digest"],
                                "default": "basic",
                            },
                            "username": {"type": "string", "description": "username"},
                            "password": {"type": "string", "description": "password"},
                            "token": {"type": "string", "description": "bearer token"},
                        },
                    },
                    "timeout": {
                        "type": "integer",
                        "description": "seconds",
                        "default": 30,
                    },
                    "verify_ssl": {
                        "type": "boolean",
                        "description": "verify SSL",
                        "default": True,
                    },
                    "follow_redirects": {
                        "type": "boolean",
                        "description": "follow redirects",
                        "default": False,
                    },
                    "output_format": {
                        "type": "string",
                        "description": "output format",
                        "enum": ["json", "headers", "body", "verbose", "meta"],
                        "default": "json",
                    },
                },
                "required": ["url"],
            },
        },
        {
            "name": "upload_file",
            "description": "Upload file. params:url,file_path,field",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "description": "upload URL"},
                    "file_path": {"type": "string", "description": "path to file"},
                    "file_field_name": {
                        "type": "string",
                        "description": "form field name",
                        "default": "file",
                    },
                    "additional_fields": {
                        "type": "object",
                        "description": "form fields",
                    },
                    "headers": {"type": "object", "description": "custom headers"},
                    "auth": {
                        "type": "object",
                        "description": "auth config",
                        "properties": {
                            "type": {
                                "type": "string",
                                "enum": ["basic", "bearer", "digest"],
                            },
                            "username": {"type": "string"},
                            "password": {"type": "string"},
                            "token": {"type": "string"},
                        },
                    },
                },
                "required": ["url", "file_path"],
            },
        },
        {
            "name": "download_file",
            "description": "Download file. params:url,output_path,continue",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "description": "download URL"},
                    "output_path": {"type": "string", "description": "output path"},
                    "continue_download": {
                        "type": "boolean",
                        "description": "resume download",
                        "default": False,
                    },
                    "headers": {"type": "object", "description": "custom headers"},
                    "auth": {
                        "type": "object",
                        "description": "auth config",
                        "properties": {
                            "type": {
                                "type": "string",
                                "enum": ["basic", "bearer", "digest"],
                            },
                            "username": {"type": "string"},
                            "password": {"type": "string"},
                            "token": {"type": "string"},
                        },
                    },
                },
                "required": ["url"],
            },
        },
        {
            "name": "test_api_endpoint",
            "description": "Test endpoint. params:url,method,expected_status,timeout",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "description": "endpoint URL"},
                    "method": {
                        "type": "string",
                        "description": "HTTP method",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                        "default": "GET",
                    },
                    "expected_status": {
                        "type": "integer",
                        "description": "expected HTTP status",
                        "default": 200,
                    },
                    "headers": {"type": "object", "description": "request headers"},
                    "data": {"type": "object", "description": "request data"},
                    "validate_response": {
                        "type": "boolean",
                        "description": "validate response",
                        "default": True,
                    },
                    "timeout": {
                        "type": "integer",
                        "description": "seconds",
                        "default": 30,
                    },
                },
                "required": ["url"],
            },
        },
        {
            "name": "manage_session",
            "description": "Manage HTTP session. params:action,session_name,session_file,url,method",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "action": {
                        "type": "string",
                        "description": "action",
                        "enum": ["create", "use", "list", "delete"],
                        "default": "create",
                    },
                    "session_name": {"type": "string", "description": "session name"},
                    "session_file": {
                        "type": "string",
                        "description": "session file path",
                    },
                    "url": {"type": "string", "description": "URL"},
                    "method": {
                        "type": "string",
                        "description": "HTTP method",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                        "default": "GET",
                    },
                    "data": {"type": "object", "description": "request data"},
                    "headers": {"type": "object", "description": "headers"},
                    "auth": {
                        "type": "object",
                        "description": "auth config",
                        "properties": {
                            "type": {
                                "type": "string",
                                "enum": ["basic", "bearer", "digest"],
                            },
                            "username": {"type": "string"},
                            "password": {"type": "string"},
                            "token": {"type": "string"},
                        },
                    },
                },
                "required": ["action"],
            },
        },
        {
            "name": "handle_authentication",
            "description": "Handle auth. params:auth_type,url,method",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "auth_type": {
                        "type": "string",
                        "description": "auth type",
                        "enum": ["basic", "bearer", "digest", "netrc"],
                        "default": "basic",
                    },
                    "username": {"type": "string", "description": "username"},
                    "password": {"type": "string", "description": "password"},
                    "token": {"type": "string", "description": "bearer token"},
                    "url": {"type": "string", "description": "auth URL"},
                    "test_endpoint": {"type": "string", "description": "test endpoint"},
                    "method": {
                        "type": "string",
                        "description": "HTTP method",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                        "default": "GET",
                    },
                },
                "required": ["auth_type", "url"],
            },
        },
        {
            "name": "format_response",
            "description": "Format response. params:url,output_type,pretty_print,style",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "description": "URL"},
                    "output_type": {
                        "type": "string",
                        "description": "output type",
                        "enum": ["headers", "body", "meta", "verbose", "quiet"],
                        "default": "json",
                    },
                    "pretty_print": {
                        "type": "boolean",
                        "description": "pretty print",
                        "default": True,
                    },
                    "style": {
                        "type": "string",
                        "description": "style",
                        "enum": ["auto", "pie-dark", "pie-light", "monokai", "fruity"],
                        "default": "auto",
                    },
                    "headers": {"type": "object", "description": "headers"},
                    "data": {"type": "object", "description": "data"},
                    "method": {
                        "type": "string",
                        "description": "HTTP method",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                        "default": "GET",
                    },
                },
                "required": ["url"],
            },
        },
        {
            "name": "test_connectivity",
            "description": "Test connectivity. params:url,methods,timeout,check_ssl",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "description": "URL"},
                    "test_methods": {
                        "type": "array",
                        "description": "HTTP methods",
                        "items": {
                            "type": "string",
                            "enum": [
                                "GET",
                                "POST",
                                "PUT",
                                "DELETE",
                                "PATCH",
                                "HEAD",
                                "OPTIONS",
                            ],
                        },
                        "default": ["GET"],
                    },
                    "timeout": {
                        "type": "integer",
                        "description": "seconds",
                        "default": 30,
                    },
                    "check_ssl": {
                        "type": "boolean",
                        "description": "check SSL",
                        "default": True,
                    },
                    "headers": {"type": "object", "description": "headers"},
                    "expected_status": {
                        "type": "integer",
                        "description": "expected status",
                        "default": 200,
                    },
                },
                "required": ["url"],
            },
        },
    ]


async def run_httpie_command(args: List[str], **kwargs) -> Dict[str, Any]:
    """Run httpie command and return results"""
    try:
        # Build command
        cmd = ["python3", "-m", "httpie"] + args

        # Add additional arguments based on kwargs
        if kwargs.get("timeout"):
            cmd.extend(["--timeout", str(kwargs["timeout"])])
        if kwargs.get("verify_ssl") is False:
            cmd.append("--verify=no")
        if kwargs.get("follow_redirects"):
            cmd.append("--follow")
        if kwargs.get("output_format"):
            format_map = {
                "json": "--pretty=all",
                "headers": "--print=h",
                "body": "--print=b",
                "verbose": "--verbose",
                "meta": "--print=m",
                "quiet": "--quiet",
            }
            cmd.append(format_map.get(kwargs["output_format"], "--pretty=all"))

        # Add style for pretty printing
        if kwargs.get("pretty_print") and kwargs.get("style"):
            cmd.extend(["--style", kwargs["style"]])

        # Run command
        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd="/home/jrede/dev/MCP",
        )

        stdout, stderr = await process.communicate()

        result = {
            "exit_code": process.returncode,
            "stdout": stdout.decode("utf-8") if stdout else "",
            "stderr": stderr.decode("utf-8") if stderr else "",
        }

        return result

    except Exception as e:
        return {
            "exit_code": -1,
            "stdout": "",
            "stderr": f"Httpie command execution failed: {str(e)}",
        }


def build_httpie_args(
    method: str,
    url: str,
    headers: Optional[Dict] = None,
    data: Optional[Dict] = None,
    query_params: Optional[Dict] = None,
    auth: Optional[Dict] = None,
    **kwargs,
) -> List[str]:
    """Build httpie command arguments from parameters"""
    args = []

    # Add method if not GET
    if method.upper() != "GET":
        args.append(method.upper())

    # Add URL
    args.append(url)

    # Add headers
    if headers:
        for key, value in headers.items():
            args.append(f"{key}:{value}")

    # Add query parameters
    if query_params:
        for key, value in query_params.items():
            args.append(f"{key}=={value}")

    # Add data based on content type
    if data:
        # Check if data contains file uploads (indicated by @ symbol in any value)
        has_files = any(str(v).startswith("@") for v in data.values() if v)
        if has_files or kwargs.get("multipart"):
            args.append("--form")
            for key, value in data.items():
                if str(value).startswith("@"):
                    args.append(f"{key}@{value[1:]}")  # Remove @ prefix for file upload
                else:
                    args.append(f"{key}={value}")
        else:
            # JSON data
            args.append("--json")
            for key, value in data.items():
                args.append(f"{key}={json.dumps(value)}")

    # Add authentication
    if auth:
        auth_type = auth.get("type", "basic")
        if auth_type == "basic":
            if auth.get("username") and auth.get("password"):
                args.extend([f"{auth['username']}:{auth['password']}"])
        elif auth_type == "bearer":
            if auth.get("token"):
                args.append(f"Authorization:Bearer {auth['token']}")
        elif auth_type == "digest":
            # Digest auth is more complex, might need special handling
            logger.warning("Digest auth not fully implemented in httpie wrapper")
            if auth.get("username") and auth.get("password"):
                args.extend([f"{auth['username']}:{auth['password']}"])
    return args


def handle_tool_call(params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle individual tool calls - wrapped by middleware"""
    name = params.get("name")
    arguments = params.get("arguments", {})

    try:
        if name == "make_request":
            method = arguments.get("method", "GET").upper()
            url = arguments["url"]
            headers = arguments.get("headers", {})
            data = arguments.get("data", {})
            query_params = arguments.get("query_params", {})
            auth = arguments.get("auth", {})
            timeout = arguments.get("timeout", 30)
            verify_ssl = arguments.get("verify_ssl", True)
            follow_redirects = arguments.get("follow_redirects", False)
            output_format = arguments.get("output_format", "json")

            # Build httpie arguments
            args = build_httpie_args(
                method=method,
                url=url,
                headers=headers,
                data=data,
                query_params=query_params,
                auth=auth,
                multipart=False,
            )

            # Run command
            result = asyncio.run(
                run_httpie_command(
                    args,
                    timeout=timeout,
                    verify_ssl=verify_ssl,
                    follow_redirects=follow_redirects,
                    output_format=output_format,
                )
            )

            if result["exit_code"] == 0:
                # Parse stdout for response
                stdout = result["stdout"]
                if output_format == "json":
                    try:
                        response_data = json.loads(stdout)
                    except json.JSONDecodeError:
                        response_data = {"body": stdout, "status": "success"}
                else:
                    response_data = {
                        "response": stdout,
                        "status": "success",
                        "exit_code": 0,
                    }
                return {"result": response_data, "status": "success"}
            else:
                return {
                    "error": result["stderr"],
                    "status": "error",
                    "exit_code": result["exit_code"],
                }

        elif name == "upload_file":
            url = arguments["url"]
            file_path = arguments["file_path"]
            file_field_name = arguments.get("file_field_name", "file")
            additional_fields = arguments.get("additional_fields", {})
            headers = arguments.get("headers", {})
            auth = arguments.get("auth", {})

            if not Path(file_path).exists():
                return {"error": f"File not found: {file_path}", "status": "error"}

            # Build multipart upload command
            args = ["POST", url, f"{file_field_name}@{file_path}"]

            # Add additional form fields
            for key, value in additional_fields.items():
                args.append(f"{key}={value}")

            # Add headers
            for key, value in headers.items():
                args.append(f"{key}:{value}")

            # Add auth
            if auth:
                auth_type = auth.get("type", "basic")
                if (
                    auth_type == "basic"
                    and auth.get("username")
                    and auth.get("password")
                ):
                    args.extend([f"{auth['username']}:{auth['password']}"])
                elif auth_type == "bearer" and auth.get("token"):
                    args.append(f"Authorization:Bearer {auth['token']}")

            # Run multipart upload
            result = asyncio.run(run_httpie_command(args, multipart=True))

            if result["exit_code"] == 0:
                return {
                    "status": "success",
                    "uploaded_file": file_path,
                    "response": result["stdout"],
                }
            else:
                return {
                    "error": result["stderr"],
                    "status": "error",
                    "exit_code": result["exit_code"],
                }

        elif name == "download_file":
            url = arguments["url"]
            output_path = arguments.get("output_path", None)
            continue_download = arguments.get("continue_download", False)
            headers = arguments.get("headers", {})
            auth = arguments.get("auth", {})

            # Use temporary file if no output path specified
            if not output_path:
                temp_file = tempfile.NamedTemporaryFile(
                    delete=False, suffix=Path(url).suffix
                )
                output_path = temp_file.name
                temp_file.close()

            # Build download command
            args = ["GET", url, "--download", f"--output={output_path}"]

            if continue_download:
                args.append("--continue")

            # Add headers
            for key, value in headers.items():
                args.append(f"{key}:{value}")

            # Add auth
            if auth:
                auth_type = auth.get("type", "basic")
                if (
                    auth_type == "basic"
                    and auth.get("username")
                    and auth.get("password")
                ):
                    args.extend([f"{auth['username']}:{auth['password']}"])
                elif auth_type == "bearer" and auth.get("token"):
                    args.append(f"Authorization:Bearer {auth['token']}")

            # Run download
            result = asyncio.run(run_httpie_command(args))

            if result["exit_code"] == 0:
                file_size = Path(output_path).stat().st_size
                return {
                    "status": "success",
                    "downloaded_to": output_path,
                    "file_size_bytes": file_size,
                    "response": result["stdout"],
                }
            else:
                # Clean up temp file on error
                if not arguments.get("output_path"):
                    Path(output_path).unlink(missing_ok=True)
                return {
                    "error": result["stderr"],
                    "status": "error",
                    "exit_code": result["exit_code"],
                }

        elif name == "test_api_endpoint":
            url = arguments["url"]
            method = arguments.get("method", "GET").upper()
            expected_status = arguments.get("expected_status", 200)
            headers = arguments.get("headers", {})
            data = arguments.get("data", {})
            validate_response = arguments.get("validate_response", True)
            timeout = arguments.get("timeout", 30)

            # Build test command
            args = [method, url]

            # Add data if present
            if data and method in ["POST", "PUT", "PATCH"]:
                args.append("--json")
                for key, value in data.items():
                    args.append(f"{key}={json.dumps(value)}")

            # Add headers
            for key, value in headers.items():
                args.append(f"{key}:{value}")

            # Run test
            result = asyncio.run(
                run_httpie_command(args, timeout=timeout, output_format="verbose")
            )

            # Parse response for status code
            response_lines = result["stdout"].split("\n")
            status_line = next(
                (line for line in response_lines if "HTTP/" in line), None
            )

            if status_line:
                # Extract status code from "HTTP/1.1 200 OK"
                status_code = int(status_line.split()[-2])
                test_passed = status_code == expected_status
            else:
                status_code = None
                test_passed = False

            test_result = {
                "url": url,
                "method": method,
                "expected_status": expected_status,
                "actual_status": status_code,
                "test_passed": test_passed,
                "response_time": "unknown",  # Httpie doesn't provide timing in stdout
                "response_body": result["stdout"] if test_passed else None,
                "error": result["stderr"] if not test_passed else None,
                "exit_code": result["exit_code"],
            }

            return {
                "result": test_result,
                "status": "success" if test_passed else "warning",
            }

        elif name == "manage_session":
            action = arguments.get("action", "create")

            if action == "create":
                session_name = arguments.get(
                    "session_name", f"session-{int(time.time())}"
                )
                session_dir = Path("/home/jrede/dev/MCP/data/sessions")
                session_dir.mkdir(exist_ok=True)
                session_file = session_dir / f"{session_name}.json"

                # Create initial session data
                session_data = {
                    "name": session_name,
                    "created": datetime.now().isoformat(),
                    "cookies": {},
                    "headers": {},
                    "auth": arguments.get("auth", {}),
                    "last_used": datetime.now().isoformat(),
                }

                with open(session_file, "w") as f:
                    json.dump(session_data, f, indent=2)

                return {
                    "status": "success",
                    "session_created": session_name,
                    "session_file": str(session_file),
                }

            elif action == "list":
                session_dir = Path("/home/jrede/dev/MCP/data/sessions")
                sessions = []
                if session_dir.exists():
                    for session_file in session_dir.glob("*.json"):
                        with open(session_file, "r") as f:
                            session_data = json.load(f)
                        sessions.append(
                            {
                                "name": session_data["name"],
                                "created": session_data["created"],
                                "last_used": session_data.get("last_used", "unknown"),
                            }
                        )
                return {"status": "success", "sessions": sessions}

            elif action == "delete":
                session_name = arguments.get("session_name")
                if session_name:
                    session_file = (
                        Path("/home/jrede/dev/MCP/data/sessions")
                        / f"{session_name}.json"
                    )
                    if session_file.exists():
                        session_file.unlink()
                        return {"status": "success", "session_deleted": session_name}
                    else:
                        return {"error": "Session not found", "status": "error"}
                return {"error": "Session name required", "status": "error"}

            else:
                return {"error": f"Unknown action: {action}", "status": "error"}

        elif name == "handle_authentication":
            auth_type = arguments.get("auth_type", "basic")
            url = arguments["url"]
            test_endpoint = arguments.get("test_endpoint", url)
            method = arguments.get("method", "GET").upper()

            if auth_type == "basic":
                username = arguments.get("username")
                password = arguments.get("password")
                if username and password:
                    auth_str = f"{username}:{password}"
                    args = [method, test_endpoint, auth_str]
                else:
                    return {
                        "error": "Username and password required for basic auth",
                        "status": "error",
                    }

            elif auth_type == "bearer":
                token = arguments.get("token")
                if token:
                    args = [method, test_endpoint, f"Authorization:Bearer {token}"]
                else:
                    return {
                        "error": "Token required for bearer auth",
                        "status": "error",
                    }

            elif auth_type == "digest":
                username = arguments.get("username")
                password = arguments.get("password")
                if username and password:
                    auth_str = f"{username}:{password}"
                    args = [method, test_endpoint, auth_str]
                    # Note: Httpie handles digest auth automatically when credentials provided
                else:
                    return {
                        "error": "Username and password required for digest auth",
                        "status": "error",
                    }

            elif auth_type == "netrc":
                args = [
                    method,
                    test_endpoint,
                    "--auth-type=ntlm",
                ]  # Httpie uses netrc for NTLM
                # Netrc file should be configured in ~/.netrc

            else:
                return {"error": f"Unknown auth type: {auth_type}", "status": "error"}

            # Test authentication
            result = asyncio.run(run_httpie_command(args, output_format="verbose"))

            if result["exit_code"] == 0:
                # Check if authentication was successful (status 200 or expected)
                response_lines = result["stdout"].split("\n")
                status_line = next(
                    (line for line in response_lines if "HTTP/" in line), None
                )
                if status_line and "200" in status_line:
                    return {
                        "status": "success",
                        "auth_type": auth_type,
                        "test_url": test_endpoint,
                        "authenticated": True,
                        "response": result["stdout"],
                    }
                else:
                    return {
                        "status": "warning",
                        "auth_type": auth_type,
                        "test_url": test_endpoint,
                        "authenticated": False,
                        "response_status": status_line,
                        "response": result["stdout"],
                    }
            else:
                return {
                    "error": result["stderr"],
                    "status": "error",
                    "exit_code": result["exit_code"],
                }

        elif name == "format_response":
            url = arguments["url"]
            output_type = arguments.get("output_type", "json")
            pretty_print = arguments.get("pretty_print", True)
            style = arguments.get("style", "auto")
            headers = arguments.get("headers", {})
            data = arguments.get("data", {})
            method = arguments.get("method", "GET").upper()

            args = [method, url]

            if data and method in ["POST", "PUT", "PATCH"]:
                args.append("--json")
                for key, value in data.items():
                    args.append(f"{key}={json.dumps(value)}")

            for key, value in headers.items():
                args.append(f"{key}:{value}")

            # Set output format
            if output_type == "headers":
                args.append("--print=h")
            elif output_type == "body":
                args.append("--print=b")
            elif output_type == "meta":
                args.append("--print=m")
            elif output_type == "verbose":
                args.append("--verbose")
            elif output_type == "quiet":
                args.append("--quiet")
            else:  # json
                args.append("--pretty=all")

            if pretty_print and style != "auto":
                args.extend(["--style", style])

            result = asyncio.run(run_httpie_command(args))

            if result["exit_code"] == 0:
                return {
                    "status": "success",
                    "formatted_response": result["stdout"],
                    "exit_code": 0,
                }
            else:
                return {
                    "error": result["stderr"],
                    "status": "error",
                    "exit_code": result["exit_code"],
                }

        elif name == "test_connectivity":
            url = arguments["url"]
            test_methods = arguments.get("test_methods", ["GET"])
            timeout = arguments.get("timeout", 30)
            check_ssl = arguments.get("check_ssl", True)
            headers = arguments.get("headers", {})
            expected_status = arguments.get("expected_status", 200)

            connectivity_results = {}

            for method in test_methods:
                args = [method.upper(), url]

                for key, value in headers.items():
                    args.append(f"{key}:{value}")

                if not check_ssl:
                    args.append("--verify=no")

                result = asyncio.run(
                    run_httpie_command(args, timeout=timeout, output_format="verbose")
                )

                # Parse status code
                response_lines = result["stdout"].split("\n")
                status_line = next(
                    (line for line in response_lines if "HTTP/" in line), None
                )
                status_code = None
                if status_line:
                    try:
                        status_code = int(status_line.split()[-2])
                    except (ValueError, IndexError):
                        pass

                test_passed = (
                    result["exit_code"] == 0 and status_code == expected_status
                )
                connectivity_results[method] = {
                    "method": method,
                    "status_code": status_code,
                    "expected_status": expected_status,
                    "test_passed": test_passed,
                    "response_time": "unknown",  # Would need timing implementation
                    "exit_code": result["exit_code"],
                    "error": result["stderr"] if not test_passed else None,
                }

            overall_passed = all(
                r["test_passed"] for r in connectivity_results.values()
            )
            return {
                "status": "success" if overall_passed else "warning",
                "url": url,
                "test_methods": test_methods,
                "results": connectivity_results,
                "overall_passed": overall_passed,
                "ssl_checked": check_ssl,
            }

        else:
            return {"error": f"Unknown tool: {name}", "status": "error"}

    except Exception as e:
        return {"error": str(e), "status": "error"}


def main():
    """Main MCP server loop"""
    print(json.dumps({"jsonrpc": "2.0", "result": list_tools(), "id": None}))
    sys.stdout.flush()

    for line in sys.stdin:
        try:
            request = json.loads(line.strip())

            if request.get("method") == "tools/list":
                response = {
                    "jsonrpc": "2.0",
                    "result": list_tools(),
                    "id": request.get("id"),
                }
                send_response(response)

            elif request.get("method") == "tools/call":
                params = request.get("params", {})
                session_id = params.get("sessionId", "unknown")
                name = params.get("name")
                arguments = params.get("arguments", {})

                # Extract task description from arguments if available
                task_desc = arguments.get("task", f"Execute {name} tool")

                # Use middleware wrapper
                wrapped_handler = middleware.handle_tool_call_wrapper(handle_tool_call)
                result = wrapped_handler(params)

                response = {"jsonrpc": "2.0", "result": result, "id": request.get("id")}
                send_response(response)

            else:
                send_error(-32601, f"Unknown method: {request.get('method')}")

        except json.JSONDecodeError:
            send_error(-32700, "Invalid JSON")
        except Exception as e:
            send_error(-32000, f"Server error: {str(e)}")


if __name__ == "__main__":
    main()
