#!/usr/bin/env python3
"""
Unit Test MCP Server
Provides test execution and analysis capabilities via MCP protocol
"""

import argparse
import json
import os
import subprocess
import sys
from typing import Any, Dict, Optional


def run_command(cmd: str, cwd: Optional[str] = None) -> Dict[str, Any]:
    """Execute a command and return results"""
    try:
        result = subprocess.run(
            cmd, shell=True, capture_output=True, text=True, cwd=cwd
        )
        return {
            "success": result.returncode == 0,
            "stdout": result.stdout,
            "stderr": result.stderr,
            "returncode": result.returncode,
        }
    except Exception as e:
        return {"success": False, "stdout": "", "stderr": str(e), "returncode": -1}


def handle_list_frameworks() -> Dict[str, Any]:
    """List available test frameworks"""
    frameworks = []

    # Check for pytest
    result = run_command("python -c 'import pytest; print(pytest.__version__)'")
    if result["success"]:
        frameworks.append(
            {
                "name": "pytest",
                "version": result["stdout"].strip(),
                "description": "Python testing framework",
            }
        )

    # Check for unittest
    result = run_command("python -c 'import unittest; print(unittest.__version__)'")
    if result["success"]:
        frameworks.append(
            {
                "name": "unittest",
                "version": result["stdout"].strip(),
                "description": "Built-in Python testing framework",
            }
        )

    # Check for node-based testing
    result = run_command("npm list -g --depth=0 | grep -E '(jest|mocha|vitest)'")
    if result["success"] and result["stdout"].strip():
        node_tests = [
            line.strip() for line in result["stdout"].split("\n") if line.strip()
        ]
        for test in node_tests:
            frameworks.append(
                {
                    "name": test.split("@")[0] if "@" in test else test,
                    "version": test.split("@")[1] if "@" in test else "unknown",
                    "description": "Node.js testing framework",
                }
            )

    return {
        "content": [
            {"type": "text", "text": json.dumps({"frameworks": frameworks}, indent=2)}
        ]
    }


def handle_run_tests(args: Dict[str, Any]) -> Dict[str, Any]:
    """Run tests in specified directory"""
    test_dir = args.get("test_directory", ".")
    framework = args.get("framework", "pytest")
    pattern = args.get("pattern", "test_*.py")
    verbose = args.get("verbose", True)

    if framework == "pytest":
        cmd = f"python -m pytest {test_dir} -k '{pattern}' --tb=short"
        if verbose:
            cmd += " -v"
    elif framework == "unittest":
        cmd = f"python -m unittest discover {test_dir} '{pattern}'"
    else:
        return {
            "content": [{"type": "text", "text": f"Unsupported framework: {framework}"}]
        }

    result = run_command(cmd, cwd=test_dir)

    # Parse results
    output = result["stdout"] + result["stderr"]
    if "FAILED" in output:
        status = "failed"
    elif "passed" in output.lower():
        status = "passed"
    else:
        status = "unknown"

    return {
        "content": [
            {
                "type": "text",
                "text": json.dumps(
                    {
                        "command": cmd,
                        "status": status,
                        "output": output,
                        "framework": framework,
                        "test_directory": test_dir,
                    },
                    indent=2,
                ),
            }
        ]
    }


def handle_analyze_results(args: Dict[str, Any]) -> Dict[str, Any]:
    """Analyze test results from a report file"""
    report_file = args.get("report_file")

    if not report_file or not os.path.exists(report_file):
        return {
            "content": [
                {"type": "text", "text": f"Report file not found: {report_file}"}
            ]
        }

    try:
        with open(report_file, "r") as f:
            content = f.read()

        # Basic analysis
        lines = content.split("\n")
        passed = len([l for l in lines if "PASSED" in l])
        failed = len([l for l in lines if "FAILED" in l])
        errors = len([l for l in lines if "ERROR" in l])

        analysis = {
            "total_lines": len(lines),
            "passed_tests": passed,
            "failed_tests": failed,
            "error_tests": errors,
            "success_rate": round(passed / (passed + failed + errors) * 100, 2)
            if (passed + failed + errors) > 0
            else 0,
        }

        return {"content": [{"type": "text", "text": json.dumps(analysis, indent=2)}]}
    except Exception as e:
        return {
            "content": [{"type": "text", "text": f"Error analyzing results: {str(e)}"}]
        }


def main():
    """Main MCP server function"""
    parser = argparse.ArgumentParser(description="Unit Test MCP Server")

    while True:
        try:
            line = sys.stdin.readline()
            if not line:
                break

            request = json.loads(line.strip())

            if request.get("method") == "tools/list":
                tools = [
                    {
                        "name": "list_frameworks",
                        "description": "List available test frameworks",
                        "inputSchema": {"type": "object", "properties": {}},
                    },
                    {
                        "name": "run_tests",
                        "description": "Run tests in a directory",
                        "inputSchema": {
                            "type": "object",
                            "properties": {
                                "test_directory": {"type": "string", "default": "."},
                                "framework": {"type": "string", "default": "pytest"},
                                "pattern": {"type": "string", "default": "test_*.py"},
                                "verbose": {"type": "boolean", "default": True},
                            },
                        },
                    },
                    {
                        "name": "analyze_results",
                        "description": "Analyze test results from a report file",
                        "inputSchema": {
                            "type": "object",
                            "properties": {"report_file": {"type": "string"}},
                        },
                    },
                ]

                response = {
                    "jsonrpc": "2.0",
                    "id": request.get("id"),
                    "result": {"tools": tools},
                }
                print(json.dumps(response))
                sys.stdout.flush()

            elif request.get("method") == "tools/call":
                tool_name = request.get("params", {}).get("name")
                args = request.get("params", {}).get("arguments", {})

                if tool_name == "list_frameworks":
                    result = handle_list_frameworks()
                elif tool_name == "run_tests":
                    result = handle_run_tests(args)
                elif tool_name == "analyze_results":
                    result = handle_analyze_results(args)
                else:
                    result = {
                        "content": [
                            {"type": "text", "text": f"Unknown tool: {tool_name}"}
                        ]
                    }

                response = {"jsonrpc": "2.0", "id": request.get("id"), "result": result}
                print(json.dumps(response))
                sys.stdout.flush()

        except KeyboardInterrupt:
            break
        except Exception as e:
            error_response = {
                "jsonrpc": "2.0",
                "id": request.get("id") if "request" in locals() else None,
                "error": {"code": -32603, "message": f"Server error: {str(e)}"},
            }
            print(json.dumps(error_response))
            sys.stdout.flush()


if __name__ == "__main__":
    main()
