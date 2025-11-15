#!/usr/bin/env python3
"""
Pytest MCP Server - Provides pytest testing capabilities through MCP protocol
"""

import asyncio
import json
import os
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Any, Dict, List, Optional
import xml.etree.ElementTree as ET


# JSON-RPC response helpers
def send_response(result: Dict[str, Any]) -> None:
    """Send JSON-RPC response to stdout"""
    print(json.dumps(result))
    sys.stdout.flush()


def send_error(code: int, message: str, data: Optional[Dict[str, Any]] = None) -> None:
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
    """Return list of available Pytest tools"""
    return [
        {
            "name": "run_tests",
            "description": "Run pytest tests in specified directory or with specific test files. Supports various pytest options including verbosity, markers, and output formats.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test directory or specific test file (default: current directory)",
                    },
                    "test_pattern": {
                        "type": "string",
                        "description": "Test name pattern to match (e.g., 'test_login*')",
                    },
                    "markers": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Pytest markers to run (e.g., ['slow', 'unit'])",
                    },
                    "verbose": {
                        "type": "boolean",
                        "description": "Enable verbose output (default: true)",
                    },
                    "tb": {
                        "type": "string",
                        "description": "Traceback style (short, long, line, no)",
                    },
                    "max_failures": {
                        "type": "integer",
                        "description": "Stop after N failures (default: unlimited)",
                    },
                    "cov": {
                        "type": "boolean",
                        "description": "Run with coverage reporting",
                    },
                },
                "required": ["test_path"],
            },
        },
        {
            "name": "get_test_report",
            "description": "Generate detailed test report in XML or JSON format. Includes test results, duration, and error details.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test directory or specific test file",
                    },
                    "output_format": {
                        "type": "string",
                        "enum": ["xml", "json"],
                        "description": "Report format (xml or json)",
                    },
                    "output_file": {
                        "type": "string",
                        "description": "Path to save report file (optional)",
                    },
                    "failed_only": {
                        "type": "boolean",
                        "description": "Generate report for failed tests only",
                    },
                },
                "required": ["test_path"],
            },
        },
        {
            "name": "list_tests",
            "description": "List all available tests in specified directory. Shows test names, locations, and markers without running them.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test directory or specific test file",
                    },
                    "collect_only": {
                        "type": "boolean",
                        "description": "List only test collection info",
                    },
                    "markers": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Filter by markers",
                    },
                },
                "required": ["test_path"],
            },
        },
        {
            "name": "run_specific_test",
            "description": "Run a specific test by name or function. Useful for running individual tests during development.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test file containing the test",
                    },
                    "test_name": {
                        "type": "string",
                        "description": "Specific test function name (e.g., 'test_user_login')",
                    },
                    "verbose": {
                        "type": "boolean",
                        "description": "Enable verbose output",
                    },
                    "tb": {
                        "type": "string",
                        "description": "Traceback style",
                    },
                },
                "required": ["test_path", "test_name"],
            },
        },
        {
            "name": "check_test_coverage",
            "description": "Run tests with coverage analysis. Requires pytest-cov package to be installed.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test directory or specific test file",
                    },
                    "source_path": {
                        "type": "string",
                        "description": "Path to source code for coverage analysis",
                    },
                    "output_format": {
                        "type": "string",
                        "enum": ["term", "html", "xml"],
                        "description": "Coverage report format",
                    },
                    "fail_under": {
                        "type": "number",
                        "description": "Fail if coverage is below this percentage",
                    },
                },
                "required": ["test_path", "source_path"],
            },
        },
        {
            "name": "validate_test_structure",
            "description": "Validate pytest test files and structure. Checks for common issues like missing test functions or incorrect naming conventions.",
            "inputSchema": {
                "type": "object",
                "properties": {
                    "test_path": {
                        "type": "string",
                        "description": "Path to test directory or specific test file",
                    },
                    "check_conventions": {
                        "type": "boolean",
                        "description": "Check pytest naming conventions",
                    },
                    "check_imports": {
                        "type": "boolean",
                        "description": "Check for import issues",
                    },
                },
                "required": ["test_path"],
            },
        },
    ]


def build_pytest_args(test_path, test_pattern=None, markers=None, verbose=True, tb="short", max_failures=None, cov=False):
    """Build pytest command arguments"""
    args = ["python3", "-m", "pytest", test_path]

    if verbose:
        args.append("-v")

    if tb:
        args.extend(["-tb", tb])

    if test_pattern:
        args.extend(["-k", test_pattern])

    if markers:
        for marker in markers:
            args.extend(["-m", marker])

    if max_failures:
        args.extend(["--maxfail", str(max_failures)])

    if cov:
        args.extend(["--cov", "--cov-report", "term-missing"])

    return args


def parse_pytest_output(output_text: str) -> Dict[str, Any]:
    """Parse pytest output to extract test results"""
    result: Dict[str, Any] = {
        "total_tests": 0,
        "passed": 0,
        "failed": 0,
        "errors": 0,
        "skipped": 0,
        "warnings": 0,
        "test_results": [],
        "summary": "",
        "exit_code": 0,
    }

    lines = output_text.split("\n")

    # Parse session summary
    for raw_line in lines:
        line = raw_line.strip()
        if "===" in line and "passed" in line:
            result["summary"] = line
            # Extract numbers from summary
            import re

            passed_match = re.search(r"(\\d+) passed", line)
            failed_match = re.search(r"(\\d+) failed", line)
            error_match = re.search(r"(\\d+) error", line)
            skipped_match = re.search(r"(\\d+) skipped", line)

            if passed_match:
                result["passed"] = int(passed_match.group(1))
            if failed_match:
                result["failed"] = int(failed_match.group(1))
            if error_match:
                result["errors"] = int(error_match.group(1))
            if skipped_match:
                result["skipped"] = int(skipped_match.group(1))

    # Use local ints for arithmetic to keep types precise
    passed = int(result.get("passed", 0) or 0)
    failed = int(result.get("failed", 0) or 0)
    errors = int(result.get("errors", 0) or 0)
    skipped = int(result.get("skipped", 0) or 0)
    result["total_tests"] = passed + failed + errors + skipped

    return result


def run_pytest_command(args, cwd=None):
    """Run pytest command and return result"""
    try:
        process = subprocess.run(
            args,
            cwd=cwd,
            capture_output=True,
            text=True,
            timeout=300  # 5 minute timeout
        )

        return {
            "success": True,
            "stdout": process.stdout,
            "stderr": process.stderr,
            "exit_code": process.returncode,
            "result": parse_pytest_output(process.stdout)
        }
    except subprocess.TimeoutExpired:
        return {
            "success": False,
            "error": "Test execution timeout (>5 minutes)"
        }
    except Exception as e:
        return {
            "success": False,
            "error": f"Failed to run tests: {str(e)}"
        }


async def handle_run_tests(args_dict):
    """Handle run_tests tool request"""
    try:
        test_path = args_dict.get("test_path", ".")
        test_pattern = args_dict.get("test_pattern")
        markers = args_dict.get("markers", [])
        verbose = args_dict.get("verbose", True)
        tb = args_dict.get("tb", "short")
        max_failures = args_dict.get("max_failures")
        cov = args_dict.get("cov", False)

        # Validate test path
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")

        # Build pytest command
        pytest_args = build_pytest_args(
            test_path, test_pattern, markers, verbose, tb, max_failures, cov
        )

        # Run pytest
        result = run_pytest_command(pytest_args)

        if result["success"]:
            response = {
                "jsonrpc": "2.0",
                "id": 1,
                "result": {
                    "success": result["exit_code"] == 0,
                    "command": " ".join(pytest_args),
                    "test_results": result["result"],
                    "output": result["stdout"],
                    "errors": result["stderr"]
                }
            }
            send_response(response)
        else:
            send_error(-32000, f"Test execution failed: {result.get('error')}")

    except Exception as e:
        send_error(-32000, f"Error running tests: {str(e)}")


async def handle_get_test_report(args_dict):
    """Handle get_test_report tool request"""
    try:
        test_path = args_dict.get("test_path", ".")
        output_format = args_dict.get("output_format", "xml")
        output_file = args_dict.get("output_file")
        failed_only = args_dict.get("failed_only", False)

        # Validate test path
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")

        # Build pytest command
        args = ["python3", "-m", "pytest", test_path, "--tb=short"]

        if output_format == "xml":
            args.append("--junit-xml")
            if output_file:
                args.extend(["--junit-logging", "system-out"])
                args.extend(["--junit-logging", "system-err"])
        elif output_format == "json":
            args.extend(["--json-report", "--json-report-file"])
            if output_file:
                args.append(output_file)
            else:
                args.append("test_report.json")

        if failed_only:
            args.append("-x")  # Exit on first failure

        # Run pytest
        result = run_pytest_command(args)

        if result["success"]:
            # Try to parse XML report if available
            xml_data: Dict[str, Any] = {}
            if output_format == "xml" and output_file and os.path.exists(output_file):
                try:
                    with open(output_file, "r", encoding="utf-8") as f:
                        tree = ET.parse(f)
                    root = tree.getroot()

                    # Initialize explicit container for testcases
                    testcases: List[Dict[str, Any]] = []

                    # Extract test cases
                    for testsuite in root.findall("testsuite"):
                        for testcase in testsuite.findall("testcase"):
                            test_data: Dict[str, Any] = {
                                "name": testcase.get("name"),
                                "classname": testcase.get("classname"),
                                "file": testcase.get("file"),
                                "line": testcase.get("line"),
                                "status": "passed",
                                "time": testcase.get("time"),
                            }

                            if len(testcase.findall("failure")) > 0:
                                failure = testcase.find("failure")
                                if failure is not None:
                                    test_data["status"] = "failed"
                                    test_data["failure"] = failure.get("message")
                            elif len(testcase.findall("error")) > 0:
                                test_data["status"] = "error"
                                error_elem = testcase.find("error")
                                if error_elem is not None:
                                    test_data["error"] = error_elem.get("message")
                            elif len(testcase.findall("skipped")) > 0:
                                test_data["status"] = "skipped"

                            testcases.append(test_data)

                    xml_data = {
                        "testsuite": {
                            "name": root.get("name"),
                            "tests": root.get("tests"),
                            "failures": root.get("failures"),
                            "errors": root.get("errors"),
                            "skipped": root.get("skipped"),
                            "time": root.get("time"),
                            "testcases": testcases,
                        }
                    }

                except Exception:
                    xml_data = {"error": "Could not parse XML report"}

            response = {
                "jsonrpc": "2.0",
                "id": 1,
                "result": {
                    "success": True,
                    "command": " ".join(args),
                    "output": result["stdout"],
                    "errors": result["stderr"],
                    "exit_code": result["exit_code"],
                    "report_data": xml_data,
                },
            }
            send_response(response)
        else:
            send_error(-32000, f"Report generation failed: {result.get('error')}")

    except Exception as e:
        send_error(-32000, f"Error generating report: {str(e)}")


async def handle_list_tests(args_dict):
    """Handle list_tests tool request"""
    try:
        test_path = args_dict.get("test_path", ".")
        collect_only = args_dict.get("collect_only", False)
        markers = args_dict.get("markers", [])

        # Validate test path
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")

        # Build pytest command
        args = ["python3", "-m", "pytest", test_path, "--collect-only", "-q"]

        if markers:
            for marker in markers:
                args.extend(["-m", marker])

        # Run pytest
        result = run_pytest_command(args)

        if result["success"]:
            # Parse collected tests
            collected_tests: List[Dict[str, Any]] = []
            stdout = result.get("stdout", "")
            if not isinstance(stdout, str):
                stdout = str(stdout)
            lines = stdout.split("\n")

            for line in lines:
                line = line.strip()
                if "::" in line and " <" not in line:
                    # This looks like a test path
                    parts = line.split("::")
                    file_path = parts[0]
                    test_name = parts[-1]

                    collected_tests.append({"file": file_path, "name": test_name})

            response = {
                "jsonrpc": "2.0",
                "id": 1,
                "result": {
                    "success": True,
                    "command": " ".join(args),
                    "output": result["stdout"],
                    "errors": result["stderr"],
                    "exit_code": result["exit_code"],
                    "tests": collected_tests,
                },
            }
            send_response(response)
        else:
            send_error(-32000, f"Test listing failed: {result.get('error')}")

    except Exception as e:
        send_error(-32000, f"Error listing tests: {str(e)}")


async def handle_run_specific_test(args_dict):
    """Handle run_specific_test tool request"""
    try:
        test_path = args_dict.get("test_path")
        test_name = args_dict.get("test_name")
        verbose = args_dict.get("verbose", True)
        tb = args_dict.get("tb", "short")

        # Validate required parameters
        if not test_path or not test_name:
            return send_error(-32602, "test_path and test_name are required")

        # Validate test path
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")

        # Build pytest command
        args = ["python3", "-m", "pytest", test_path, "-k", test_name]

        if verbose:
            args.append("-v")
        if tb:
            args.extend(["-tb", tb])

        # Run pytest
        result = run_pytest_command(args)

        if result["success"]:
            response = {
                "jsonrpc": "2.0",
                "id": 1,
                "result": {
                    "success": result["exit_code"] == 0,
                    "test_name": test_name,
                    "test_path": test_path,
                    "command": " ".join(args),
                    "test_results": result["result"],
                    "output": result["stdout"],
                    "errors": result["stderr"]
                }
            }
            send_response(response)
        else:
            send_error(-32000, f"Specific test execution failed: {result.get('error')}")

    except Exception as e:
        send_error(-32000, f"Error running specific test: {str(e)}")


async def handle_check_test_coverage(args_dict):
    """Handle check_test_coverage tool request"""
    try:
        test_path = args_dict.get("test_path", ".")
        source_path = args_dict.get("source_path")
        output_format = args_dict.get("output_format", "term")
        fail_under = args_dict.get("fail_under")

        # Validate required parameters
        if not test_path or not source_path:
            return send_error(-32602, "test_path and source_path are required")

        # Validate paths
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")
        if not os.path.exists(source_path):
            return send_error(-32602, f"Source path does not exist: {source_path}")

        # Build pytest command
        args = [
            "python3", "-m", "pytest", test_path,
            "--cov", source_path,
            "--cov-report", output_format
        ]

        if fail_under:
            args.extend(["--cov-fail-under", str(fail_under)])

        # Run pytest
        result = run_pytest_command(args)

        if result["success"]:
            response = {
                "jsonrpc": "2.0",
                "id": 1,
                "result": {
                    "success": result["exit_code"] == 0,
                    "test_path": test_path,
                    "source_path": source_path,
                    "command": " ".join(args),
                    "coverage_report": result["stdout"],
                    "test_results": result["result"],
                    "errors": result["stderr"]
                }
            }
            send_response(response)
        else:
            send_error(-32000, f"Coverage analysis failed: {result.get('error')}")

    except Exception as e:
        send_error(-32000, f"Error checking test coverage: {str(e)}")


async def handle_validate_test_structure(args_dict):
    """Handle validate_test_structure tool request"""
    try:
        test_path = args_dict.get("test_path", ".")
        check_conventions = args_dict.get("check_conventions", True)
        check_imports = args_dict.get("check_imports", True)

        # Validate test path
        if not os.path.exists(test_path):
            return send_error(-32602, f"Test path does not exist: {test_path}")

        # Initialize validation results
        validation_results: Dict[str, Any] = {
            "valid": True,
            "issues": [],
            "warnings": [],
            "recommendations": [],
        }

        # Walk through test directory
        for root_dir, dirs, files in os.walk(test_path):
            for file in files:
                if file.endswith(".py") and file.startswith("test_"):
                    file_path = os.path.join(root_dir, file)

                    try:
                        with open(file_path, "r", encoding="utf-8") as f:
                            content = f.read()

                        # Check naming conventions
                        if check_conventions:
                            if "def test_" not in content and "class Test" not in content:
                                issues = validation_results.setdefault("issues", [])
                                if isinstance(issues, list):
                                    issues.append(
                                        f"File {file_path} appears to contain no valid test functions"
                                    )
                                validation_results["valid"] = False

                        # Check imports
                        if check_imports:
                            import_lines = [
                                line
                                for line in content.split("\n")
                                if line.strip().startswith("import")
                            ]
                            if not import_lines:
                                warnings = validation_results.setdefault("warnings", [])
                                if isinstance(warnings, list):
                                    warnings.append(
                                        f"File {file_path} has no imports"
                                    )

                    except Exception as e:
                        issues = validation_results.setdefault("issues", [])
                        if isinstance(issues, list):
                            issues.append(
                                f"Could not read file {file_path}: {str(e)}"
                            )
                        validation_results["valid"] = False

        # Add recommendations
        recommendations = validation_results.setdefault("recommendations", [])
        if not isinstance(recommendations, list):
            recommendations = []
            validation_results["recommendations"] = recommendations

        if check_conventions:
            recommendations.append(
                "Ensure test files are named with 'test_' prefix"
            )
            recommendations.append(
                "Ensure test functions are named with 'test_' prefix"
            )

        recommendations.append(
            "Use descriptive test names that explain what is being tested"
        )
        recommendations.append(
            "Group related tests in classes with 'Test' prefix"
        )
        recommendations.append(
            "Use pytest fixtures for setup and teardown"
        )

        response = {
            "jsonrpc": "2.0",
            "id": 1,
            "result": {
                "success": True,
                "validation_results": validation_results,
            },
        }
        send_response(response)

    except Exception as e:
        send_error(-32000, f"Error validating test structure: {str(e)}")


async def handle_request(request):
    """Handle incoming JSON-RPC request"""
    method = request.get("method")
    params = request.get("params", {})

    try:
        if method == "tools/list":
            tools = list_tools()
            send_response({
                "jsonrpc": "2.0",
                "id": request.get("id"),
                "result": {"tools": tools}
            })
        elif method == "tools/call":
            tool_name = params.get("name")
            arguments = params.get("arguments", {})

            if tool_name == "run_tests":
                await handle_run_tests(arguments)
            elif tool_name == "get_test_report":
                await handle_get_test_report(arguments)
            elif tool_name == "list_tests":
                await handle_list_tests(arguments)
            elif tool_name == "run_specific_test":
                await handle_run_specific_test(arguments)
            elif tool_name == "check_test_coverage":
                await handle_check_test_coverage(arguments)
            elif tool_name == "validate_test_structure":
                await handle_validate_test_structure(arguments)
            else:
                send_error(-32601, f"Unknown tool: {tool_name}")
        else:
            send_error(-32601, f"Unknown method: {method}")

    except Exception as e:
        send_error(-32000, f"Internal error: {str(e)}")


async def main():
    """Main MCP server loop"""
    # Send capabilities
    send_response({
        "jsonrpc": "2.0",
        "id": None,
        "result": {
            "capabilities": {
                "tools": {
                    "listChanged": False
                }
            }
        }
    })

    # Read requests from stdin
    for line in sys.stdin:
        if line.strip():
            try:
                request = json.loads(line.strip())
                await handle_request(request)
            except json.JSONDecodeError:
                send_error(-32700, "Invalid JSON")
            except Exception as e:
                send_error(-32000, f"Unexpected error: {str(e)}")


if __name__ == "__main__":
    asyncio.run(main())