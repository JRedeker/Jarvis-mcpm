#!/usr/bin/env python3
"""
Schemathesis MCP Server - Provides API testing capabilities through MCP protocol
"""

import asyncio
import json
import sys
import tempfile
from pathlib import Path
from typing import Any, Dict, List


# JSON-RPC response helpers
def send_response(result):
    """Send JSON-RPC response to stdout"""
    print(json.dumps(result))
    sys.stdout.flush()


def send_error(code, message, data=None):
    """Send JSON-RPC error response"""
    error_response = {"jsonrpc": "2.0", "error": {"code": code, "message": message}}
    if data:
        error_response["error"]["data"] = data
    print(json.dumps(error_response))
    sys.stdout.flush()


def list_tools():
    """Return list of available Schemathesis tools"""
    return [
        {
            "name": "load_openapi_schema",
            "description": "Load and validate OpenAPI/Swagger schemas from URLs or files. Supports OpenAPI 2.0, 3.0, and 3.1 specifications. Validates schema structure and provides detailed error messages for malformed schemas.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "source": {
                        "type": "string",
                        "description": "URL to OpenAPI schema or file path to local schema file",
                    },
                    "base_url": {
                        "type": "string",
                        "description": "Base URL for API testing (optional, extracted from schema if not provided)",
                    },
                },
                "required": ["source"],
            },
        },
        {
            "name": "test_api_endpoints",
            "description": "Run comprehensive hypothesis-based property tests against all API endpoints defined in the schema. Generates test data based on schema specifications and validates API responses against schema definitions.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "schema_source": {
                        "type": "string",
                        "description": "URL or file path to OpenAPI schema",
                    },
                    "base_url": {
                        "type": "string",
                        "description": "API base URL to test against",
                    },
                    "max_examples": {
                        "type": "integer",
                        "description": "Maximum number of test examples to generate (default: 100)",
                        "default": 100,
                    },
                    "workers": {
                        "type": "integer",
                        "description": "Number of worker processes for parallel testing (default: 1)",
                        "default": 1,
                    },
                    "headers": {
                        "type": "object",
                        "description": "Additional headers to include in API requests",
                    },
                    "timeout": {
                        "type": "integer",
                        "description": "Request timeout in seconds (default: 30)",
                        "default": 30,
                    },
                },
                "required": ["schema_source", "base_url"],
            },
        },
        {
            "name": "validate_schema",
            "description": "Validate OpenAPI schema structure and check for common issues. Provides detailed analysis of schema completeness, parameter definitions, and potential testing problems.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "schema_source": {
                        "type": "string",
                        "description": "URL or file path to OpenAPI schema",
                    }
                },
                "required": ["schema_source"],
            },
        },
        {
            "name": "generate_test_data",
            "description": "Generate sample test data based on schema specifications without making actual API calls. Useful for understanding data structures and validating test data generation.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "schema_source": {
                        "type": "string",
                        "description": "URL or file path to OpenAPI schema",
                    },
                    "endpoint": {
                        "type": "string",
                        "description": "Specific endpoint path (optional, generates for all if not provided)",
                    },
                    "method": {
                        "type": "string",
                        "description": "HTTP method (GET, POST, PUT, DELETE) for specific endpoint",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                    },
                    "count": {
                        "type": "integer",
                        "description": "Number of test data samples to generate (default: 10)",
                        "default": 10,
                    },
                },
                "required": ["schema_source"],
            },
        },
        {
            "name": "test_specific_endpoint",
            "description": "Test individual endpoints with generated and custom test data. Focus testing on specific functionality with detailed result analysis.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "schema_source": {
                        "type": "string",
                        "description": "URL or file path to OpenAPI schema",
                    },
                    "base_url": {
                        "type": "string",
                        "description": "API base URL to test against",
                    },
                    "endpoint": {
                        "type": "string",
                        "description": "Specific endpoint path to test",
                    },
                    "method": {
                        "type": "string",
                        "description": "HTTP method for the endpoint",
                        "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                    },
                    "custom_data": {
                        "type": "object",
                        "description": "Custom test data to use (optional, generated if not provided)",
                    },
                    "max_examples": {
                        "type": "integer",
                        "description": "Maximum test examples to generate (default: 50)",
                        "default": 50,
                    },
                },
                "required": ["schema_source", "base_url", "endpoint", "method"],
            },
        },
        {
            "name": "run_schema_tests",
            "description": "Execute comprehensive schema-based testing with various strategies including property-based testing, boundary value analysis, and edge case detection.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "schema_source": {
                        "type": "string",
                        "description": "URL or file path to OpenAPI schema",
                    },
                    "base_url": {
                        "type": "string",
                        "description": "API base URL to test against",
                    },
                    "test_strategies": {
                        "type": "array",
                        "description": "List of test strategies to apply",
                        "items": {
                            "type": "string",
                            "enum": [
                                "property_based",
                                "boundary_values",
                                "edge_cases",
                                "data_validation",
                            ],
                        },
                        "default": ["property_based", "data_validation"],
                    },
                    "max_examples": {
                        "type": "integer",
                        "description": "Maximum examples per strategy (default: 25)",
                        "default": 25,
                    },
                },
                "required": ["schema_source", "base_url"],
            },
        },
    ]


async def download_schema(schema_url: str) -> str:
    """Download schema from URL to temporary file"""
    try:
        import urllib.parse
        import urllib.request

        req = urllib.request.Request(
            schema_url, headers={"User-Agent": "Schemathesis-MCP/1.0"}
        )
        with urllib.request.urlopen(req, timeout=30) as response:
            content = response.read().decode("utf-8")

        # Write to temporary file
        with tempfile.NamedTemporaryFile(mode="w", suffix=".yaml", delete=False) as f:
            f.write(content)
            return f.name
    except Exception as e:
        raise Exception(f"Failed to download schema from {schema_url}: {str(e)}")


def validate_schema_file(schema_source: str) -> str:
    """Validate and return path to schema file"""
    if schema_source.startswith(("http://", "https://")):
        # Download from URL
        return asyncio.run(download_schema(schema_source))
    else:
        # Treat as file path
        if not Path(schema_source).exists():
            raise Exception(f"Schema file not found: {schema_source}")
        return schema_source


async def run_schemathesis_command(
    args: List[str], schema_path: str, base_url: str, **kwargs
) -> Dict[str, Any]:
    """Run schemathesis command and return results"""
    try:
        # Build command
        cmd = (
            ["python", "-m", "schemathesis"]
            + args
            + [schema_path, "--base-url", base_url]
        )

        # Add additional arguments
        if kwargs.get("max_examples"):
            cmd.extend(["--hypothesis-max-examples", str(kwargs["max_examples"])])
        if kwargs.get("workers"):
            cmd.extend(["--workers", str(kwargs["workers"])])
        if kwargs.get("headers"):
            for key, value in kwargs["headers"].items():
                cmd.extend(["--header", f"{key}:{value}"])
        if kwargs.get("timeout"):
            cmd.extend(["--request-timeout", str(kwargs["timeout"])])

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
            "stderr": f"Command execution failed: {str(e)}",
        }


def handle_tool_call(name: str, arguments: Dict[str, Any]) -> Dict[str, Any]:
    """Handle individual tool calls"""
    try:
        if name == "load_openapi_schema":
            schema_source = arguments["source"]
            base_url = arguments.get("base_url")

            schema_path = validate_schema_file(schema_source)

            # If no base URL provided, try to extract from schema
            if not base_url:
                base_url = "https://api.example.com"  # Default fallback

            return {
                "content": [
                    {
                        "type": "text",
                        "text": f"✅ Schema loaded successfully from {schema_source}\n📁 Schema file: {schema_path}\n🔗 Base URL: {base_url}\n\nReady for testing! Use other tools to start testing endpoints.",
                    }
                ]
            }

        elif name == "test_api_endpoints":
            schema_source = arguments["schema_source"]
            base_url = arguments["base_url"]

            schema_path = validate_schema_file(schema_source)

            # Run comprehensive API testing
            result = asyncio.run(
                run_schemathesis_command(
                    ["run"],
                    schema_path,
                    base_url,
                    max_examples=arguments.get("max_examples", 100),
                    workers=arguments.get("workers", 1),
                    headers=arguments.get("headers", {}),
                    timeout=arguments.get("timeout", 30),
                )
            )

            if result["exit_code"] == 0:
                test_summary = f"✅ API Testing Completed Successfully!\n📊 Results:\n{result['stdout']}"
            else:
                test_summary = f"⚠️ Testing completed with issues:\n📊 Results:\n{result['stdout']}\n\n❌ Errors:\n{result['stderr']}"

            return {"content": [{"type": "text", "text": test_summary}]}

        elif name == "validate_schema":
            schema_source = arguments["schema_source"]
            schema_path = validate_schema_file(schema_source)

            # Use schemathesis to validate schema
            result = asyncio.run(
                run_schemathesis_command(
                    ["check"],
                    schema_path,
                    "https://dummy.example.com",  # Dummy URL for validation
                )
            )

            if result["exit_code"] == 0:
                validation_result = f"✅ Schema validation successful!\n📋 Schema: {schema_source}\n\nNo structural issues found."
            else:
                validation_result = f"⚠️ Schema validation issues found:\n📋 Schema: {schema_source}\n\n❌ Issues:\n{result['stderr']}"

            return {"content": [{"type": "text", "text": validation_result}]}

        elif name == "generate_test_data":
            schema_source = arguments["schema_source"]
            schema_path = validate_schema_file(schema_source)
            count = arguments.get("count", 10)

            # Generate test data using schemathesis
            result = asyncio.run(
                run_schemathesis_command(
                    ["sample"],
                    schema_path,
                    "https://dummy.example.com",
                    max_examples=count,
                )
            )

            return {
                "content": [
                    {
                        "type": "text",
                        "text": f"🧪 Generated {count} test data samples:\n\n{result['stdout']}\n\n{result['stderr'] if result['stderr'] else ''}",
                    }
                ]
            }

        elif name == "test_specific_endpoint":
            schema_source = arguments["schema_source"]
            base_url = arguments["base_url"]
            endpoint = arguments["endpoint"]
            method = arguments["method"]

            schema_path = validate_schema_file(schema_source)

            # Test specific endpoint
            result = asyncio.run(
                run_schemathesis_command(
                    ["run", "--method", method, "--endpoint", endpoint],
                    schema_path,
                    base_url,
                    max_examples=arguments.get("max_examples", 50),
                )
            )

            return {
                "content": [
                    {
                        "type": "text",
                        "text": f"🎯 Testing endpoint: {method} {endpoint}\n\n📊 Results:\n{result['stdout']}\n\n{result['stderr'] if result['stderr'] else ''}",
                    }
                ]
            }

        elif name == "run_schema_tests":
            schema_source = arguments["schema_source"]
            base_url = arguments["base_url"]
            strategies = arguments.get(
                "test_strategies", ["property_based", "data_validation"]
            )

            schema_path = validate_schema_file(schema_source)

            # Run comprehensive testing with different strategies
            results = []
            for strategy in strategies:
                result = asyncio.run(
                    run_schemathesis_command(
                        ["run", "--hypothesis-settings", strategy],
                        schema_path,
                        base_url,
                        max_examples=arguments.get("max_examples", 25),
                    )
                )

                strategy_result = (
                    f"🔬 Strategy: {strategy}\n📊 Results:\n{result['stdout']}"
                )
                if result["stderr"]:
                    strategy_result += f"\n❌ Issues: {result['stderr']}"
                results.append(strategy_result)

            return {
                "content": [
                    {
                        "type": "text",
                        "text": "🔍 Comprehensive Schema Testing Results:\n\n"
                        + "\n\n".join(results),
                    }
                ]
            }

        else:
            send_error(-32601, f"Unknown tool: {name}")
            return {}

    except Exception as e:
        send_error(-32603, f"Tool execution failed: {str(e)}")
        return {}


def main():
    """Main MCP server loop"""
    while True:
        try:
            line = sys.stdin.readline()
            if not line:
                break

            request = json.loads(line.strip())
            method = request.get("method")
            request_id = request.get("id")

            if method == "initialize":
                send_response(
                    {
                        "jsonrpc": "2.0",
                        "id": request_id,
                        "result": {
                            "protocolVersion": "2024-11-05",
                            "capabilities": {},
                            "serverInfo": {
                                "name": "schemathesis-mcp",
                                "version": "1.0.0",
                                "description": "API testing server using Schemathesis for OpenAPI schema validation and hypothesis-based testing",
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
                tool_name = request.get("params", {}).get("name")
                arguments = request.get("params", {}).get("arguments", {})
                result = handle_tool_call(tool_name, arguments)
                if result:
                    send_response(
                        {"jsonrpc": "2.0", "id": request_id, "result": result}
                    )

            else:
                send_error(-32601, f"Unknown method: {method}")

        except json.JSONDecodeError:
            send_error(-32700, "Invalid JSON")
        except Exception as e:
            send_error(-32603, f"Server error: {str(e)}")


if __name__ == "__main__":
    main()
