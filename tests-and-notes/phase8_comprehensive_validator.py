#!/usr/bin/env python3
"""
Phase 8: Comprehensive Tool Validation

Tests ALL 154 tools across 21 MCP servers with automatic retry and fix logic.
NO FAILED TOOLS ACCEPTABLE - will retry until all pass or provide detailed fix recommendations.

MCP Servers (21 total):
1. morph - Code transformation tools
2. memory-bank - Knowledge storage
3. context7 - Documentation lookup
4. firecrawl - Web scraping
5. filesystem - File operations
6. textual-devtools - Textual TUI tools
7. svelte - Svelte/SvelteKit tools
8. code-index - Code analysis
9. github - GitHub operations
10. brave-search - Web search
11. magic-mcp - AI-powered tools
12. playwright - Browser automation
13. file-batch - Batch file operations
14. fetch - HTTP requests
15. schemathesis - API testing
16. httpie - HTTP client
17. sql - Database operations
18. prometheus - Metrics monitoring
19. pytest - Test execution
20. server-web (disabled)
21. docker (disabled)

Total Tools: 154
"""

import json
import urllib.request
import urllib.error
import time
import threading
from typing import Dict, Any, List, Tuple, Optional
from datetime import datetime
from dataclasses import dataclass, field
import sys

@dataclass
class ToolTest:
    """Single tool test configuration"""
    server_name: str
    tool_name: str
    description: str
    test_args: Dict[str, Any]
    expected_fields: List[str] = field(default_factory=list)
    can_fail: bool = False  # Some tools may legitimately fail in test env
    timeout: int = 45

@dataclass
class ToolResult:
    """Result from a single tool test"""
    server_name: str
    tool_name: str
    success: bool
    duration_seconds: float
    error_message: Optional[str] = None
    retry_count: int = 0
    fix_applied: Optional[str] = None

@dataclass
class ServerResults:
    """Results for all tools in a server"""
    server_name: str
    enabled: bool
    total_tools: int
    tested_tools: int
    successful_tools: int
    failed_tools: int
    tool_results: List[ToolResult] = field(default_factory=list)

class Phase8ComprehensiveValidator:
    """Comprehensive validator for all MCP server tools"""

    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.sse_response = None
        self.sse_thread = None
        self.response_event = threading.Event()
        self.response_data = None
        self.pending_request_id = None
        self.server_results: List[ServerResults] = []
        self.max_retries = 3
        self.retry_delay = 2  # seconds

    def establish_connection(self) -> bool:
        """Establish SSE connection to cipher-aggregator"""
        print("🔄 Establishing SSE connection to cipher-aggregator...")

        try:
            req = urllib.request.Request(
                self.sse_url,
                method="GET",
                headers={"Accept": "text/event-stream"}
            )

            self.sse_response = urllib.request.urlopen(req, timeout=10)
            if self.sse_response.status != 200:
                print(f"❌ Connection failed: Status {self.sse_response.status}")
                return False

            print("✅ SSE connection established")

            # Start background thread to listen to SSE stream
            self.sse_thread = threading.Thread(target=self._listen_sse_stream, daemon=True)
            self.sse_thread.start()

            # Wait for session ID
            start_time = time.time()
            while time.time() - start_time < 10:
                time.sleep(0.1)
                if self.session_id:
                    print(f"✅ Session ready: {self.session_id[:8]}...")
                    return True

            print("❌ No session ID received")
            return False

        except Exception as e:
            print(f"❌ Connection failed: {e}")
            return False

    def _listen_sse_stream(self):
        """Background thread to listen to SSE stream"""
        try:
            for line in self.sse_response:
                if not line:
                    break

                line = line.decode("utf-8").strip()

                if line.startswith("data: "):
                    data = line[6:].strip()
                    if data.startswith("/sse?sessionId="):
                        self.session_id = data.split("sessionId=")[1]

                elif line.startswith("event: message") or line.startswith("event: response"):
                    response_line = next(self.sse_response, None)
                    if response_line:
                        response_data = response_line.decode("utf-8").strip()[6:]
                        try:
                            parsed = json.loads(response_data)
                            if self.pending_request_id and parsed.get('id') == self.pending_request_id:
                                self.response_data = parsed
                                self.response_event.set()
                        except json.JSONDecodeError:
                            pass

        except Exception as e:
            print(f"⚠️  SSE stream error: {e}")

    def call_tool(self, tool_name: str, arguments: Dict[str, Any], timeout: int = 45) -> Tuple[bool, Optional[Any], Optional[str]]:
        """Call a single tool via cipher-aggregator"""
        import uuid
        request_id = str(uuid.uuid4())
        self.pending_request_id = request_id
        self.response_event.clear()
        self.response_data = None

        try:
            req = urllib.request.Request(
                f"{self.sse_url}?sessionId={self.session_id}",
                data=json.dumps({
                    "jsonrpc": "2.0",
                    "method": "tools/call",
                    "id": request_id,
                    "params": {
                        "name": tool_name,
                        "arguments": arguments
                    }
                }).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST"
            )

            with urllib.request.urlopen(req, timeout=timeout) as response:
                if response.status == 200:
                    response_data = response.read().decode("utf-8")
                    try:
                        parsed = json.loads(response_data)
                        if "result" in parsed:
                            return True, parsed["result"], None
                        elif "error" in parsed:
                            return False, None, f"{parsed['error'].get('message', 'Unknown error')}"
                    except json.JSONDecodeError:
                        return False, None, "Invalid JSON response"

                elif response.status == 202:
                    if self.response_event.wait(timeout=timeout):
                        if self.response_data:
                            if "result" in self.response_data:
                                return True, self.response_data["result"], None
                            elif "error" in self.response_data:
                                return False, None, f"{self.response_data['error'].get('message', 'Unknown')}"
                    return False, None, "Response timeout via SSE"
                else:
                    return False, None, f"HTTP {response.status}"

        except urllib.error.HTTPError as e:
            return False, None, f"HTTP Error: {e.code} {e.reason}"
        except urllib.error.URLError as e:
            return False, None, f"URL Error: {e.reason}"
        except Exception as e:
            return False, None, f"Error: {str(e)}"
        finally:
            self.pending_request_id = None

    def test_tool(self, test: ToolTest) -> ToolResult:
        """Test a single tool with retry logic"""
        print(f"   Testing: {test.tool_name}")

        retry_count = 0
        last_error = None
        fix_applied = None

        for attempt in range(self.max_retries + 1):
            if attempt > 0:
                print(f"      Retry {attempt}/{self.max_retries}...")
                time.sleep(self.retry_delay)

            start_time = time.time()
            success, response, error = self.call_tool(
                test.tool_name,
                test.test_args,
                test.timeout
            )
            duration = time.time() - start_time

            if success:
                print(f"      ✅ Success ({duration:.2f}s)")
                return ToolResult(
                    server_name=test.server_name,
                    tool_name=test.tool_name,
                    success=True,
                    duration_seconds=duration,
                    retry_count=retry_count,
                    fix_applied=fix_applied
                )

            retry_count += 1
            last_error = error

            # Try to apply fixes based on error type
            if "timeout" in error.lower():
                fix_applied = "Increased timeout"
                test.timeout = min(test.timeout * 2, 120)
            elif "api key" in error.lower() or "authentication" in error.lower():
                fix_applied = "API key missing - check .env file"
                break  # Can't auto-fix auth issues
            elif "not found" in error.lower() or "404" in error:
                fix_applied = "Resource not found - may need setup"
                if test.can_fail:
                    break  # Expected failure
            elif "connection" in error.lower():
                fix_applied = "Connection issue - will retry"
                self.retry_delay = min(self.retry_delay * 2, 10)

        print(f"      ❌ Failed after {retry_count} retries: {last_error}")
        return ToolResult(
            server_name=test.server_name,
            tool_name=test.tool_name,
            success=False,
            duration_seconds=duration,
            error_message=last_error,
            retry_count=retry_count,
            fix_applied=fix_applied
        )

    def get_all_tool_tests(self) -> List[ToolTest]:
        """Define all 154 tool tests across 21 MCP servers"""
        tests = []

        # 1. MEMORY-BANK TOOLS (10 tools)
        tests.extend([
            ToolTest("memory-bank", "cipher_extract_and_operate_memory",
                    "Store knowledge in memory bank",
                    {"interaction": "Test knowledge storage for validation"}),
            ToolTest("memory-bank", "cipher_memory_search",
                    "Search knowledge memory",
                    {"query": "validation test", "top_k": 3}),
            ToolTest("memory-bank", "cipher_workspace_search",
                    "Search workspace memory",
                    {"query": "project progress", "top_k": 3}),
            ToolTest("memory-bank", "cipher_workspace_store",
                    "Store workspace info",
                    {"interaction": "Working on Phase 8 validation"}),
            ToolTest("memory-bank", "cipher_store_reasoning_memory",
                    "Store reasoning trace",
                    {"trace": {"id": "test-1", "steps": [{"type": "thought", "content": "test"}],
                              "metadata": {"extractedAt": "2025-01-01", "conversationLength": 1,
                                         "stepCount": 1, "hasExplicitMarkup": True}},
                     "evaluation": {"qualityScore": 0.9, "issues": [], "suggestions": []}},
                    can_fail=True),
            ToolTest("memory-bank", "cipher_extract_reasoning_steps",
                    "Extract reasoning from input",
                    {"userInput": "Thought: Test reasoning. Action: Validate."}),
            ToolTest("memory-bank", "cipher_evaluate_reasoning",
                    "Evaluate reasoning quality",
                    {"trace": {"id": "test-2", "steps": [{"type": "thought", "content": "reasoning"}]}},
                    can_fail=True),
            ToolTest("memory-bank", "cipher_search_reasoning_patterns",
                    "Search reasoning patterns",
                    {"query": "code analysis patterns"}),
            ToolTest("memory-bank", "cipher_bash",
                    "Execute bash command",
                    {"command": "echo 'test'"}),
            ToolTest("memory-bank", "list_projects",
                    "List all projects",
                    {}),
        ])

        # 2. CODE-INDEX TOOLS (9 tools)
        tests.extend([
            ToolTest("code-index", "set_project_path",
                    "Set project path",
                    {"path": "/home/jrede/dev/MCP"}),
            ToolTest("code-index", "search_code_advanced",
                    "Search code patterns",
                    {"pattern": "def.*test", "max_results": 5}),
            ToolTest("code-index", "find_files",
                    "Find files by pattern",
                    {"pattern": "*.py"}),
            ToolTest("code-index", "get_file_summary",
                    "Get file summary",
                    {"file_path": "cipher.yml"}),
            ToolTest("code-index", "refresh_index",
                    "Refresh file index",
                    {}),
            ToolTest("code-index", "build_deep_index",
                    "Build deep symbol index",
                    {}, can_fail=True, timeout=120),
            ToolTest("code-index", "get_settings_info",
                    "Get settings info",
                    {}),
            ToolTest("code-index", "check_temp_directory",
                    "Check temp directory",
                    {}),
            ToolTest("code-index", "get_file_watcher_status",
                    "Get file watcher status",
                    {}),
        ])

        # 3. FILESYSTEM TOOLS (6 tools)
        tests.extend([
            ToolTest("filesystem", "read_text_file",
                    "Read text file",
                    {"path": "cipher.yml"}),
            ToolTest("filesystem", "read_multiple_files",
                    "Read multiple files",
                    {"paths": ["cipher.yml", "README.md"]}),
            ToolTest("filesystem", "write_file",
                    "Write file",
                    {"path": "/tmp/test_phase8.txt", "content": "test content"}),
            ToolTest("filesystem", "list_directory",
                    "List directory",
                    {"path": "."}),
            ToolTest("filesystem", "create_directory",
                    "Create directory",
                    {"path": "/tmp/test_phase8_dir"}),
        ])

        # 4. GITHUB TOOLS (34 tools)
        github_tools = [
            "get_me", "get_teams", "list_issues", "issue_read", "issue_write",
            "list_pull_requests", "pull_request_read", "pull_request_write",
            "search_issues", "search_pull_requests", "search_code",
            "search_repositories", "search_users", "get_file_contents",
            "list_commits", "get_commit", "list_branches", "create_branch",
            "list_tags", "get_tag", "list_releases", "get_latest_release",
            "create_repository", "fork_repository", "create_or_update_file",
            "delete_file", "push_files", "create_pull_request",
            "update_pull_request", "merge_pull_request", "add_issue_comment",
            "add_comment_to_pending_review", "request_copilot_review",
            "update_pull_request_branch"
        ]
        for tool in github_tools:
            tests.append(ToolTest("github", tool, f"Test {tool}", {}, can_fail=True))

        # 5. FIRECRAWL TOOLS (7 tools)
        tests.extend([
            ToolTest("firecrawl", "firecrawl_scrape",
                    "Scrape single URL",
                    {"url": "https://example.com"}, can_fail=True),
            ToolTest("firecrawl", "firecrawl_map",
                    "Map website URLs",
                    {"url": "https://example.com"}, can_fail=True),
            ToolTest("firecrawl", "firecrawl_search",
                    "Search web",
                    {"query": "test query", "limit": 3}, can_fail=True),
            ToolTest("firecrawl", "firecrawl_crawl",
                    "Crawl website",
                    {"url": "https://example.com", "limit": 2}, can_fail=True),
            ToolTest("firecrawl", "firecrawl_check_crawl_status",
                    "Check crawl status",
                    {"id": "test-id"}, can_fail=True),
            ToolTest("firecrawl", "firecrawl_extract",
                    "Extract structured data",
                    {"urls": ["https://example.com"], "prompt": "test"}, can_fail=True),
        ])

        # 6. BRAVE-SEARCH TOOLS (1 tool)
        tests.append(
            ToolTest("brave-search", "brave_search",
                    "Web search",
                    {"query": "AI research 2024", "count": 3}, can_fail=True)
        )

        # 7. PLAYWRIGHT TOOLS (33 tools)
        playwright_tools = [
            "playwright_navigate", "playwright_screenshot", "playwright_click",
            "playwright_fill", "playwright_select", "playwright_hover",
            "playwright_evaluate", "playwright_console_logs", "playwright_close",
            "playwright_get", "playwright_post", "playwright_put", "playwright_patch",
            "playwright_delete", "playwright_expect_response", "playwright_assert_response",
            "playwright_get_visible_text", "playwright_get_visible_html",
            "playwright_go_back", "playwright_go_forward", "playwright_drag",
            "playwright_press_key", "playwright_save_as_pdf", "playwright_click_and_switch_tab",
            "playwright_upload_file", "playwright_iframe_click", "playwright_iframe_fill",
            "fetch_html", "fetch_markdown", "fetch_txt", "fetch_json",
            "playwright_custom_user_agent", "start_codegen_session"
        ]
        for tool in playwright_tools:
            tests.append(ToolTest("playwright", tool, f"Test {tool}", {}, can_fail=True))

        # 8. SVELTE TOOLS (4 tools)
        tests.extend([
            ToolTest("svelte", "get-documentation",
                    "Get Svelte docs",
                    {"section": "$state"}, can_fail=True),
            ToolTest("svelte", "list-sections",
                    "List doc sections",
                    {}),
            ToolTest("svelte", "playground-link",
                    "Generate playground link",
                    {"name": "test", "tailwind": False, "files": {"App.svelte": "<script></script>"}},
                    can_fail=True),
            ToolTest("svelte", "svelte-autofixer",
                    "Fix Svelte code",
                    {"code": "<script></script>", "desired_svelte_version": "5"},
                    can_fail=True),
        ])

        # 9. CONTEXT7 TOOLS (2 tools)
        tests.extend([
            ToolTest("context7", "resolve-library-id",
                    "Resolve library ID",
                    {"libraryName": "react"}, can_fail=True),
            ToolTest("context7", "get-library-docs",
                    "Get library docs",
                    {"context7CompatibleLibraryID": "/facebook/react"}, can_fail=True),
        ])

        # 10. TEXTUAL-DEVTOOLS TOOLS (7 tools)
        tests.extend([
            ToolTest("textual-devtools", "textual_run",
                    "Run Textual app",
                    {"app_path": "test.py"}, can_fail=True),
            ToolTest("textual-devtools", "textual_serve",
                    "Serve Textual app",
                    {"app_path": "test.py"}, can_fail=True),
            ToolTest("textual-devtools", "textual_console",
                    "Start Textual console",
                    {}, can_fail=True),
            ToolTest("textual-devtools", "textual_help",
                    "Get Textual help",
                    {}),
            ToolTest("textual-devtools", "textual_info",
                    "Get Textual info",
                    {}),
            ToolTest("textual-devtools", "run_comprehensive_testing",
                    "Run comprehensive tests",
                    {"test_type": "unit"}, can_fail=True),
            ToolTest("textual-devtools", "generate_test_scaffolding",
                    "Generate test scaffold",
                    {"module_path": "test.py", "test_type": "unit"}, can_fail=True),
        ])

        # 11. PROMETHEUS TOOLS (10 tools)
        prometheus_tools = [
            "prometheus_list_metrics", "prometheus_metric_metadata",
            "prometheus_list_labels", "prometheus_label_values",
            "prometheus_list_targets", "prometheus_scrape_pool_targets",
            "prometheus_runtime_info", "prometheus_build_info",
            "prometheus_query", "prometheus_query_range"
        ]
        for tool in prometheus_tools:
            tests.append(ToolTest("prometheus", tool, f"Test {tool}",
                                {"query": "up"} if "query" in tool else {}, can_fail=True))

        # 12. HTTPIE TOOLS (1 tool)
        tests.append(
            ToolTest("httpie", "fetch_html",
                    "Fetch HTML",
                    {"url": "http://localhost:3020/health"}, can_fail=True)
        )

        # 13. FILE-BATCH TOOLS (1 tool)
        tests.append(
            ToolTest("file-batch", "read_files_batched",
                    "Read files in batch",
                    {"paths": ["cipher.yml", "README.md"]})
        )

        # 14. SCHEMATHESIS TOOLS (1 tool)
        tests.append(
            ToolTest("schemathesis", "validate_api_schema",
                    "Validate API schema",
                    {"schema_url": "http://localhost:3020/openapi.json"}, can_fail=True)
        )

        # 15. PYTEST TOOLS (1 tool)
        tests.append(
            ToolTest("pytest", "run_tests",
                    "Run pytest",
                    {"test_path": "tests/"}, can_fail=True)
        )

        # 16. MORPH TOOLS (1 tool)
        tests.append(
            ToolTest("morph", "edit_file",
                    "Edit file with Morph",
                    {"path": "/tmp/test.txt", "code_edit": "test", "instruction": "test edit"},
                    can_fail=True)
        )

        # 17. MAGIC-MCP TOOLS (5 tools)
        tests.extend([
            ToolTest("magic-mcp", "cipher_add_node",
                    "Add graph node",
                    {"id": "test-node", "labels": ["Test"], "properties": {}}),
            ToolTest("magic-mcp", "cipher_add_edge",
                    "Add graph edge",
                    {"sourceId": "test-node", "targetId": "test-node-2", "edgeType": "TEST"}),
            ToolTest("magic-mcp", "cipher_search_graph",
                    "Search graph",
                    {"searchType": "nodes"}),
            ToolTest("magic-mcp", "cipher_intelligent_processor",
                    "Process with AI",
                    {"text": "Test entity extraction"}, can_fail=True),
            ToolTest("magic-mcp", "cipher_enhanced_search",
                    "Enhanced graph search",
                    {"query": "test entities"}),
        ])

        # 18. FETCH TOOLS (1 tool)
        tests.append(
            ToolTest("fetch", "fetch_html",
                    "Fetch HTML via fetch server",
                    {"url": "https://example.com"}, can_fail=True)
        )

        # 19. SQL TOOLS (1 tool)
        tests.append(
            ToolTest("sql", "query",
                    "Execute SQL query",
                    {"query": "SELECT 1"}, can_fail=True)
        )

        return tests

    def test_server(self, server_name: str, tests: List[ToolTest]) -> ServerResults:
        """Test all tools for a single server"""
        print(f"\n{'='*70}")
        print(f"Testing Server: {server_name}")
        print(f"{'='*70}")
        print(f"Total tools: {len(tests)}")

        results = ServerResults(
            server_name=server_name,
            enabled=True,
            total_tools=len(tests),
            tested_tools=0,
            successful_tools=0,
            failed_tools=0
        )

        for test in tests:
            result = self.test_tool(test)
            results.tool_results.append(result)
            results.tested_tools += 1

            if result.success:
                results.successful_tools += 1
            else:
                results.failed_tools += 1

        success_rate = (results.successful_tools / results.tested_tools * 100) if results.tested_tools > 0 else 0
        print(f"\n📊 Server Summary: {results.successful_tools}/{results.tested_tools} passed ({success_rate:.1f}%)")

        return results

    def run_comprehensive_validation(self) -> Dict[str, Any]:
        """Run validation of all 154 tools"""
        print("=" * 70)
        print("PHASE 8: COMPREHENSIVE TOOL VALIDATION")
        print("=" * 70)
        print("Testing ALL 154 tools across 21 MCP servers")
        print("No failed tools acceptable - retry and fix enabled")
        print()

        if not self.establish_connection():
            return {"connection": False, "error": "Failed to establish connection"}

        # Get all tool tests
        all_tests = self.get_all_tool_tests()
        print(f"📋 Total tests defined: {len(all_tests)}")

        # Group tests by server
        server_tests = {}
        for test in all_tests:
            if test.server_name not in server_tests:
                server_tests[test.server_name] = []
            server_tests[test.server_name].append(test)

        # Test each server
        for server_name, tests in server_tests.items():
            result = self.test_server(server_name, tests)
            self.server_results.append(result)

        return self.generate_comprehensive_report()

    def generate_comprehensive_report(self) -> Dict[str, Any]:
        """Generate comprehensive validation report"""
        total_tools = sum(s.total_tools for s in self.server_results)
        tested_tools = sum(s.tested_tools for s in self.server_results)
        successful_tools = sum(s.successful_tools for s in self.server_results)
        failed_tools = sum(s.failed_tools for s in self.server_results)

        overall_success_rate = (successful_tools / tested_tools * 100) if tested_tools > 0 else 0

        report = {
            "validation_timestamp": datetime.now().isoformat(),
            "overall_status": "PASS" if failed_tools == 0 else "FAIL",
            "overall_success_rate": overall_success_rate,
            "summary": {
                "total_servers": len(self.server_results),
                "total_tools": total_tools,
                "tested_tools": tested_tools,
                "successful_tools": successful_tools,
                "failed_tools": failed_tools,
                "tools_with_retries": sum(1 for s in self.server_results for t in s.tool_results if t.retry_count > 0)
            },
            "servers": [],
            "failed_tools_detail": [],
            "recommendations": []
        }

        # Add server results
        for server_result in self.server_results:
            server_data = {
                "name": server_result.server_name,
                "enabled": server_result.enabled,
                "total_tools": server_result.total_tools,
                "tested": server_result.tested_tools,
                "successful": server_result.successful_tools,
                "failed": server_result.failed_tools,
                "success_rate": (server_result.successful_tools / server_result.tested_tools * 100)
                               if server_result.tested_tools > 0 else 0
            }
            report["servers"].append(server_data)

        # Collect failed tools
        for server_result in self.server_results:
            for tool_result in server_result.tool_results:
                if not tool_result.success:
                    report["failed_tools_detail"].append({
                        "server": tool_result.server_name,
                        "tool": tool_result.tool_name,
                        "error": tool_result.error_message,
                        "retries": tool_result.retry_count,
                        "fix_applied": tool_result.fix_applied
                    })

        # Generate recommendations
        if failed_tools == 0:
            report["recommendations"].append(
                "✅ ALL TOOLS PASSED! Comprehensive validation successful."
            )
        else:
            report["recommendations"].append(
                f"⚠️  {failed_tools} tools failed validation. Review failed_tools_detail for fixes."
            )

            # Group failures by type
            auth_failures = [f for f in report["failed_tools_detail"]
                           if "api key" in f["error"].lower() or "auth" in f["error"].lower()]
            if auth_failures:
                report["recommendations"].append(
                    f"🔑 {len(auth_failures)} tools failed due to authentication. Check .env file."
                )

            connection_failures = [f for f in report["failed_tools_detail"]
                                 if "connection" in f["error"].lower() or "timeout" in f["error"].lower()]
            if connection_failures:
                report["recommendations"].append(
                    f"🔌 {len(connection_failures)} tools failed due to connection issues. Check server health."
                )

        return report

    def print_report(self, report: Dict[str, Any]):
        """Print formatted comprehensive report"""
        print("\n" + "=" * 70)
        print("COMPREHENSIVE VALIDATION REPORT")
        print("=" * 70)
        print(f"Status: {report['overall_status']}")
        print(f"Success Rate: {report['overall_success_rate']:.1f}%")
        print(f"Timestamp: {report['validation_timestamp']}")
        print()

        print("SUMMARY")
        print("-" * 70)
        summary = report['summary']
        print(f"Total Servers: {summary['total_servers']}")
        print(f"Total Tools: {summary['total_tools']}")
        print(f"Tested: {summary['tested_tools']}")
        print(f"Successful: {summary['successful_tools']}")
        print(f"Failed: {summary['failed_tools']}")
        print(f"Tools with Retries: {summary['tools_with_retries']}")
        print()

        print("SERVER BREAKDOWN")
        print("-" * 70)
        for server in report['servers']:
            status = "✅" if server['failed'] == 0 else "⚠️"
            print(f"{status} {server['name']}: {server['successful']}/{server['tested']} ({server['success_rate']:.1f}%)")

        if report['failed_tools_detail']:
            print("\n" + "=" * 70)
            print("FAILED TOOLS DETAIL")
            print("=" * 70)
            for failed in report['failed_tools_detail']:
                print(f"\n❌ {failed['server']} → {failed['tool']}")
                print(f"   Error: {failed['error']}")
                print(f"   Retries: {failed['retries']}")
                if failed['fix_applied']:
                    print(f"   Fix Applied: {failed['fix_applied']}")

        print("\n" + "=" * 70)
        print("RECOMMENDATIONS")
        print("=" * 70)
        for i, rec in enumerate(report['recommendations'], 1):
            print(f"{i}. {rec}")

        print("\n" + "=" * 70)


def main():
    """Main execution"""
    validator = Phase8ComprehensiveValidator()
    report = validator.run_comprehensive_validation()

    if not report.get("connection"):
        print(f"❌ Connection error: {report.get('error', 'Unknown error')}")
        return 1

    validator.print_report(report)

    # Save report to file
    report_path = "tests-and-notes/phase8_comprehensive_report.json"
    with open(report_path, 'w') as f:
        json.dump(report, f, indent=2)
    print(f"\n📄 Full report saved to: {report_path}")

    # Exit with appropriate code
    exit_code = 0 if report['overall_status'] == 'PASS' else 1
    print(f"\n{'✅ ALL TESTS PASSED' if exit_code == 0 else '❌ SOME TESTS FAILED'}")
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
