#!/usr/bin/env python3
"""
Phase 7: MCP Server Health Validation
Tests all 21 configured MCP servers for proper initialization and tool availability
"""

import json
import urllib.request
import urllib.error
import time
from typing import Dict, Any, List, Tuple

class ServerHealthValidator:
    """Validates health of all configured MCP servers"""

    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.server_results = {}

        # All 21 servers from cipher.yml
        self.servers = [
            "morph", "memory-bank", "context7", "firecrawl", "filesystem",
            "textual-devtools", "svelte", "code-index", "github", "brave-search",
            "server-web", "magic-mcp", "playwright", "docker", "file-batch",
            "fetch", "schemathesis", "httpie", "sql", "prometheus", "pytest"
        ]

        # Known issues from previous investigation
        self.known_issues = {
            "brave-search": "Package name typo in cipher.yml (@brave/brave-search-mCP-server should be mcp not mCP)",
            "sql": "Requires MySQL running on localhost:3306",
            "magic-mcp": "Package @21stdev/magic-mcp not found in npm registry",
            "docker": "Disabled in cipher.yml (enabled: false)",
            "server-web": "Disabled in cipher.yml (enabled: false)"
        }

    def establish_connection(self) -> bool:
        """Establish SSE connection"""
        print("🔄 Establishing SSE connection...")

        try:
            req = urllib.request.Request(
                self.sse_url,
                method="GET",
                headers={"Accept": "text/event-stream"}
            )

            response = urllib.request.urlopen(req, timeout=10)

            for line in response:
                line = line.decode("utf-8").strip()
                if line.startswith("data: /sse?sessionId="):
                    self.session_id = line.split("sessionId=")[1]
                    print(f"✅ Connected (Session: {self.session_id[:8]}...)")
                    return True

        except Exception as e:
            print(f"❌ Connection failed: {e}")
            return False

        return False

    def get_all_tools(self) -> Dict[str, List[str]]:
        """Get list of all available tools grouped by server"""
        print("\n📋 Fetching all available tools...")

        try:
            req = urllib.request.Request(
                f"{self.sse_url}?sessionId={self.session_id}",
                data=json.dumps({
                    "jsonrpc": "2.0",
                    "method": "tools/list",
                    "id": "get-all-tools"
                }).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST"
            )

            with urllib.request.urlopen(req, timeout=30) as response:
                data = json.loads(response.read().decode("utf-8"))

                if "result" in data and "tools" in data["result"]:
                    tools = data["result"]["tools"]
                    print(f"   ✅ Found {len(tools)} total tools")

                    # Group tools by server (inferred from tool name prefixes)
                    tools_by_server = {}
                    for tool in tools:
                        name = tool.get("name", "")

                        # Infer server from tool name patterns
                        if name.startswith("cipher_"):
                            server = "memory-bank"
                        elif name.startswith("github_") or name == "create_repository":
                            server = "github"
                        elif name.startswith("playwright_"):
                            server = "playwright"
                        elif name.startswith("firecrawl_"):
                            server = "firecrawl"
                        elif name.startswith("get-") or name == "playground-link":
                            server = "svelte"
                        elif name in ["read_text_file", "read_multiple_files", "write_file",
                                     "list_directory", "create_directory"]:
                            server = "filesystem"
                        elif name.startswith("set_project_path") or name.startswith("search_code"):
                            server = "code-index"
                        elif name.startswith("prometheus_"):
                            server = "prometheus"
                        elif name.startswith("run_comprehensive_testing") or name.startswith("enforce_coverage"):
                            server = "pytest"
                        elif name == "read_files_batched":
                            server = "file-batch"
                        elif name in ["fetch_html", "fetch_markdown", "fetch_txt", "fetch_json"]:
                            server = "fetch"
                        elif name == "textual_run" or name == "textual_serve":
                            server = "textual-devtools"
                        elif name == "edit_file":
                            server = "morph"
                        elif name.startswith("resolve-library-id"):
                            server = "context7"
                        else:
                            server = "unknown"

                        if server not in tools_by_server:
                            tools_by_server[server] = []
                        tools_by_server[server].append(name)

                    return tools_by_server

        except Exception as e:
            print(f"   ❌ Failed to fetch tools: {e}")
            return {}

        return {}

    def test_server_health(self, server: str, tools: List[str]) -> Tuple[bool, str]:
        """Test if a specific server is healthy by checking if it has tools"""

        if server in self.known_issues:
            return False, f"Known issue: {self.known_issues[server]}"

        if len(tools) == 0:
            return False, "No tools exposed"

        return True, f"{len(tools)} tools available"

    def test_sample_tool(self, tool_name: str) -> Tuple[bool | None, str]:
        """Test a sample tool call to verify it's actually callable"""

        # Define safe test calls for various tools
        safe_tests = {
            "cipher_memory_search": {"query": "test", "top_k": 1},
            "list_directory": {"path": "."},
            "prometheus_list_metrics": {},
            "set_project_path": {"path": "/home/jrede/dev/MCP"},
            "textual_info": {},
            "get-documentation": {"section": "overview"},
        }

        if tool_name not in safe_tests:
            return None, "No safe test available"

        try:
            req = urllib.request.Request(
                f"{self.sse_url}?sessionId={self.session_id}",
                data=json.dumps({
                    "jsonrpc": "2.0",
                    "method": "tools/call",
                    "id": f"test-{tool_name}",
                    "params": {
                        "name": tool_name,
                        "arguments": safe_tests[tool_name]
                    }
                }).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST"
            )

            with urllib.request.urlopen(req, timeout=10) as response:
                data = json.loads(response.read().decode("utf-8"))

                if "result" in data:
                    return True, "Tool callable"
                elif "error" in data:
                    return False, f"Error: {data['error'].get('message', 'Unknown')}"

        except Exception as e:
            return False, f"Call failed: {str(e)[:50]}"

        return False, "Unknown response"

    def run_validation(self) -> Dict[str, Any]:
        """Run complete server health validation"""
        print("=" * 70)
        print("Phase 7: MCP Server Health Validation")
        print("=" * 70)

        if not self.establish_connection():
            return {"connection": False}

        # Get all tools grouped by server
        tools_by_server = self.get_all_tools()

        results = {
            "connection": True,
            "total_servers": len(self.servers),
            "servers_with_tools": 0,
            "servers_without_tools": 0,
            "healthy_servers": [],
            "unhealthy_servers": [],
            "disabled_servers": [],
            "tool_counts": {}
        }

        print("\n" + "=" * 70)
        print("Server Health Analysis")
        print("=" * 70)

        for server in self.servers:
            tools = tools_by_server.get(server, [])
            is_healthy, message = self.test_server_health(server, tools)

            results["tool_counts"][server] = len(tools)

            # Check if server is disabled
            if server in ["docker", "server-web"]:
                results["disabled_servers"].append(server)
                status = "⏸️  DISABLED"
                print(f"{status}: {server:20s} - {message}")
                continue

            if is_healthy:
                results["servers_with_tools"] += 1
                results["healthy_servers"].append(server)
                status = "✅ HEALTHY"
            else:
                results["servers_without_tools"] += 1
                results["unhealthy_servers"].append(server)
                status = "❌ UNHEALTHY"

            print(f"{status}: {server:20s} - {message}")

        # Test sample tools from healthy servers
        print("\n" + "=" * 70)
        print("Sample Tool Testing (Healthy Servers)")
        print("=" * 70)

        sample_tools = {
            "memory-bank": "cipher_memory_search",
            "filesystem": "list_directory",
            "code-index": "set_project_path",
            "prometheus": "prometheus_list_metrics",
            "textual-devtools": "textual_info",
        }

        for server, tool in sample_tools.items():
            if server in results["healthy_servers"]:
                is_callable, message = self.test_sample_tool(tool)
                if is_callable is not None:
                    status = "✅ PASS" if is_callable else "❌ FAIL"
                    print(f"{status}: {server:20s} - {tool:30s} - {message}")

        # Calculate success rate
        active_servers = results["total_servers"] - len(results["disabled_servers"])
        success_rate = (results["servers_with_tools"] / active_servers * 100) if active_servers > 0 else 0

        print("\n" + "=" * 70)
        print("Validation Summary")
        print("=" * 70)
        print(f"Total Servers Configured: {results['total_servers']}")
        print(f"Disabled Servers: {len(results['disabled_servers'])}")
        print(f"Active Servers: {active_servers}")
        print(f"Healthy Servers: {results['servers_with_tools']}")
        print(f"Unhealthy Servers: {results['servers_without_tools']}")
        print(f"Success Rate: {success_rate:.1f}%")

        # List specific issues
        if results["unhealthy_servers"]:
            print("\n" + "=" * 70)
            print("Servers Requiring Attention")
            print("=" * 70)
            for server in results["unhealthy_servers"]:
                if server in self.known_issues:
                    print(f"❌ {server:20s} - {self.known_issues[server]}")
                else:
                    print(f"❌ {server:20s} - No tools exposed (check logs)")

        print("\n" + "=" * 70)
        print("Recommendations")
        print("=" * 70)

        if success_rate >= 90:
            print("✅ Excellent! >90% of active servers are healthy")
        elif success_rate >= 80:
            print("⚠️  Good, but room for improvement (80-90% healthy)")
        else:
            print("❌ Critical: <80% of servers are healthy")

        print("\nNext Steps:")
        if "brave-search" in results["unhealthy_servers"]:
            print("1. Fix brave-search typo: mCP-server → mcp-server in cipher.yml")
        if "sql" in results["unhealthy_servers"]:
            print("2. Either install MySQL or disable sql server")
        if "magic-mcp" in results["unhealthy_servers"]:
            print("3. Remove magic-mcp (package doesn't exist) or find correct package name")

        print("4. Review cipher-aggregator logs for any connection errors")
        print("5. Restart mcp-manager.sh after configuration changes")

        return results

if __name__ == "__main__":
    validator = ServerHealthValidator()
    results = validator.run_validation()

    # Exit with appropriate code
    success_rate = (results["servers_with_tools"] /
                   (results["total_servers"] - len(results.get("disabled_servers", []))) * 100)
    exit(0 if success_rate >= 80 else 1)
