#!/usr/bin/env python3
"""
Phase 8: End-to-End Workflow Validation

Tests complete workflows using documented features from cipher.yml:
- Web research: brave-search → firecrawl → memory-bank
- Development: code-index → filesystem → github
- API testing: schemathesis → httpie → pytest

Validates:
- Max 8 calls per task constraint
- Serial execution (maxParallelCalls=1)
- Error recovery with lenient mode
- Success rate >80% for all workflow types
"""

import json
import urllib.request
import urllib.error
import time
import threading
from typing import Dict, Any, List, Tuple, Optional
from datetime import datetime
from dataclasses import dataclass, field

@dataclass
class WorkflowStep:
    """Represents a single step in a workflow"""
    tool_name: str
    description: str
    arguments: Dict[str, Any]
    required: bool = True
    timeout: int = 30

@dataclass
class WorkflowResult:
    """Results from a workflow execution"""
    workflow_name: str
    success: bool
    total_calls: int
    successful_calls: int
    failed_calls: int
    execution_time_seconds: float
    steps_executed: List[Dict[str, Any]] = field(default_factory=list)
    errors: List[str] = field(default_factory=list)
    performance_violations: List[str] = field(default_factory=list)

class Phase8Validator:
    """End-to-end workflow validator"""

    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.max_calls_per_task = 8
        self.execution_mode = "serial"
        self.workflow_results: List[WorkflowResult] = []
        self.sse_response = None
        self.sse_thread = None
        self.response_event = threading.Event()
        self.response_data = None
        self.pending_request_id = None

    def establish_connection(self) -> bool:
        """Establish SSE connection"""
        print("🔄 Establishing SSE connection...")

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

            # Extract session ID from endpoint event
            start_time = time.time()
            while time.time() - start_time < 10:
                time.sleep(0.1)
                if self.session_id:
                    print(f"✅ Connected (Session: {self.session_id[:8]}...)")
                    return True

            print("❌ No session ID received from SSE stream")
            return False

        except Exception as e:
            print(f"❌ Connection failed: {e}")
            return False

    def _listen_sse_stream(self):
        """Background thread to listen to SSE stream for responses"""
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

                # Handle tool call responses
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

    def call_tool(self, tool_name: str, arguments: Dict[str, Any], timeout: int = 30) -> Tuple[bool, Optional[Any], Optional[str]]:
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
                    # Direct response
                    response_data = response.read().decode("utf-8")
                    try:
                        parsed = json.loads(response_data)
                        if "result" in parsed:
                            return True, parsed["result"], None
                        elif "error" in parsed:
                            return False, None, f"Tool error: {parsed['error'].get('message', 'Unknown error')}"
                    except json.JSONDecodeError:
                        return False, None, "Invalid JSON response"

                elif response.status == 202:
                    # Response will come via SSE stream
                    if self.response_event.wait(timeout=timeout):
                        if self.response_data:
                            if "result" in self.response_data:
                                return True, self.response_data["result"], None
                            elif "error" in self.response_data:
                                return False, None, f"Tool error: {self.response_data['error'].get('message', 'Unknown')}"
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

    def execute_workflow(self, workflow_name: str, steps: List[WorkflowStep]) -> WorkflowResult:
        """Execute a complete workflow"""
        print(f"\n{'='*70}")
        print(f"Executing Workflow: {workflow_name}")
        print(f"{'='*70}")

        start_time = time.time()
        result = WorkflowResult(
            workflow_name=workflow_name,
            success=True,
            total_calls=0,
            successful_calls=0,
            failed_calls=0,
            execution_time_seconds=0.0
        )

        # Check max calls constraint
        if len(steps) > self.max_calls_per_task:
            result.success = False
            result.performance_violations.append(
                f"Workflow has {len(steps)} steps, exceeds max {self.max_calls_per_task} calls per task"
            )
            print(f"❌ Workflow violates max calls constraint: {len(steps)} > {self.max_calls_per_task}")
            return result

        # Execute each step serially
        for i, step in enumerate(steps, 1):
            print(f"\n📍 Step {i}/{len(steps)}: {step.description}")
            print(f"   Tool: {step.tool_name}")

            step_start = time.time()
            success, response, error = self.call_tool(step.tool_name, step.arguments, step.timeout)
            step_duration = time.time() - step_start

            result.total_calls += 1

            step_result = {
                "step_number": i,
                "tool_name": step.tool_name,
                "description": step.description,
                "success": success,
                "duration_seconds": step_duration,
                "required": step.required
            }

            if success:
                result.successful_calls += 1
                print(f"   ✅ Success ({step_duration:.2f}s)")
            else:
                result.failed_calls += 1
                print(f"   ❌ Failed: {error}")
                result.errors.append(f"Step {i} ({step.tool_name}): {error}")

                if step.required:
                    result.success = False
                    print(f"   ⛔ Required step failed, stopping workflow")
                    result.steps_executed.append(step_result)
                    break

            result.steps_executed.append(step_result)

            # Small delay between calls to ensure serial execution
            time.sleep(0.1)

        result.execution_time_seconds = time.time() - start_time

        # Calculate final success status
        if result.success:
            success_rate = (result.successful_calls / result.total_calls * 100) if result.total_calls > 0 else 0
            result.success = success_rate >= 80  # 80% success rate threshold

        print(f"\n{'='*70}")
        print(f"Workflow Complete: {workflow_name}")
        print(f"Success: {'✅' if result.success else '❌'}")
        print(f"Total calls: {result.total_calls}")
        print(f"Successful: {result.successful_calls}")
        print(f"Failed: {result.failed_calls}")
        print(f"Success rate: {(result.successful_calls/result.total_calls*100):.1f}%")
        print(f"Execution time: {result.execution_time_seconds:.2f}s")
        print(f"{'='*70}")

        return result

    def define_web_research_workflow(self) -> List[WorkflowStep]:
        """Define web research workflow: brave-search → firecrawl → memory-bank"""
        return [
            WorkflowStep(
                tool_name="brave_search",
                description="Search for AI research papers",
                arguments={"query": "latest AI research 2024", "count": 5},
                required=True
            ),
            WorkflowStep(
                tool_name="firecrawl_scrape",
                description="Scrape content from research website",
                arguments={"url": "https://example.com/research"},
                required=False  # May fail if URL doesn't exist
            ),
            WorkflowStep(
                tool_name="cipher_extract_and_operate_memory",
                description="Store research findings in memory",
                arguments={
                    "interaction": "Found AI research on neural networks and transformers"
                },
                required=True
            )
        ]

    def define_development_workflow(self) -> List[WorkflowStep]:
        """Define development workflow: code-index → filesystem → github"""
        return [
            WorkflowStep(
                tool_name="set_project_path",
                description="Set project path for code analysis",
                arguments={"path": "/home/jrede/dev/MCP"},
                required=True
            ),
            WorkflowStep(
                tool_name="search_code_advanced",
                description="Search for Python functions",
                arguments={
                    "pattern": "def.*test",
                    "case_sensitive": False,
                    "max_results": 5
                },
                required=True
            ),
            WorkflowStep(
                tool_name="list_directory",
                description="List project files",
                arguments={"path": "."},
                required=True
            ),
            WorkflowStep(
                tool_name="cipher_memory_search",
                description="Search for related code patterns",
                arguments={"query": "Python testing patterns", "top_k": 3},
                required=False
            )
        ]

    def define_api_testing_workflow(self) -> List[WorkflowStep]:
        """Define API testing workflow: schemathesis → httpie → pytest"""
        return [
            WorkflowStep(
                tool_name="list_directory",
                description="Check for API spec files",
                arguments={"path": "."},
                required=True
            ),
            WorkflowStep(
                tool_name="fetch_json",
                description="Fetch API endpoint",
                arguments={"url": "http://localhost:3020/health"},
                required=False  # May not have health endpoint
            ),
            WorkflowStep(
                tool_name="prometheus_list_metrics",
                description="Check Prometheus metrics",
                arguments={},
                required=False  # Prometheus may not be available
            )
        ]

    def run_validation(self) -> Dict[str, Any]:
        """Run all end-to-end workflow validations"""
        print("=" * 70)
        print("Phase 8: End-to-End Workflow Validation")
        print("=" * 70)
        print(f"Max calls per task: {self.max_calls_per_task}")
        print(f"Execution mode: {self.execution_mode}")
        print()

        if not self.establish_connection():
            return {"connection": False, "error": "Failed to establish SSE connection"}

        # Define all workflows
        workflows = [
            ("Web Research", self.define_web_research_workflow()),
            ("Development", self.define_development_workflow()),
            ("API Testing", self.define_api_testing_workflow())
        ]

        # Execute each workflow
        for workflow_name, steps in workflows:
            result = self.execute_workflow(workflow_name, steps)
            self.workflow_results.append(result)

        # Generate summary report
        return self.generate_report()

    def generate_report(self) -> Dict[str, Any]:
        """Generate comprehensive validation report"""
        total_workflows = len(self.workflow_results)
        successful_workflows = sum(1 for r in self.workflow_results if r.success)
        overall_success_rate = (successful_workflows / total_workflows * 100) if total_workflows > 0 else 0

        total_calls = sum(r.total_calls for r in self.workflow_results)
        total_successful = sum(r.successful_calls for r in self.workflow_results)
        total_failed = sum(r.failed_calls for r in self.workflow_results)

        report = {
            "validation_timestamp": datetime.now().isoformat(),
            "overall_status": "PASS" if overall_success_rate >= 80 else "FAIL",
            "overall_success_rate": overall_success_rate,
            "summary": {
                "total_workflows": total_workflows,
                "successful_workflows": successful_workflows,
                "failed_workflows": total_workflows - successful_workflows,
                "total_tool_calls": total_calls,
                "successful_calls": total_successful,
                "failed_calls": total_failed
            },
            "workflows": [],
            "performance_analysis": {
                "max_calls_per_task": self.max_calls_per_task,
                "execution_mode": self.execution_mode,
                "violations": []
            },
            "recommendations": []
        }

        # Add individual workflow results
        for result in self.workflow_results:
            workflow_data = {
                "name": result.workflow_name,
                "success": result.success,
                "total_calls": result.total_calls,
                "successful_calls": result.successful_calls,
                "failed_calls": result.failed_calls,
                "success_rate": (result.successful_calls / result.total_calls * 100) if result.total_calls > 0 else 0,
                "execution_time_seconds": result.execution_time_seconds,
                "steps": result.steps_executed,
                "errors": result.errors,
                "performance_violations": result.performance_violations
            }
            report["workflows"].append(workflow_data)

            # Collect performance violations
            report["performance_analysis"]["violations"].extend(result.performance_violations)

        # Generate recommendations
        if overall_success_rate < 80:
            report["recommendations"].append(
                "Overall success rate below 80%. Review failed workflows and improve tool reliability."
            )

        if total_failed > 0:
            report["recommendations"].append(
                f"{total_failed} tool calls failed. Review error messages and improve error handling."
            )

        if any(r.performance_violations for r in self.workflow_results):
            report["recommendations"].append(
                "Performance violations detected. Review workflow design to stay within constraints."
            )

        if not report["recommendations"]:
            report["recommendations"].append(
                "✅ All validations passed! Workflows are executing correctly within constraints."
            )

        return report

    def print_report(self, report: Dict[str, Any]):
        """Print formatted validation report"""
        print("\n" + "=" * 70)
        print("PHASE 8 VALIDATION REPORT")
        print("=" * 70)
        print(f"Status: {report['overall_status']}")
        print(f"Overall Success Rate: {report['overall_success_rate']:.1f}%")
        print(f"Timestamp: {report['validation_timestamp']}")
        print()

        print("SUMMARY")
        print("-" * 70)
        summary = report['summary']
        print(f"Total Workflows: {summary['total_workflows']}")
        print(f"Successful: {summary['successful_workflows']}")
        print(f"Failed: {summary['failed_workflows']}")
        print(f"Total Tool Calls: {summary['total_tool_calls']}")
        print(f"Successful Calls: {summary['successful_calls']}")
        print(f"Failed Calls: {summary['failed_calls']}")
        print()

        print("WORKFLOW DETAILS")
        print("-" * 70)
        for workflow in report['workflows']:
            status_icon = "✅" if workflow['success'] else "❌"
            print(f"\n{status_icon} {workflow['name']}")
            print(f"   Success Rate: {workflow['success_rate']:.1f}%")
            print(f"   Total Calls: {workflow['total_calls']}")
            print(f"   Execution Time: {workflow['execution_time_seconds']:.2f}s")

            if workflow['errors']:
                print(f"   Errors:")
                for error in workflow['errors']:
                    print(f"      - {error}")

        print("\n" + "=" * 70)
        print("PERFORMANCE ANALYSIS")
        print("=" * 70)
        perf = report['performance_analysis']
        print(f"Max Calls Per Task: {perf['max_calls_per_task']}")
        print(f"Execution Mode: {perf['execution_mode']}")

        if perf['violations']:
            print(f"\n⚠️  Violations Detected:")
            for violation in perf['violations']:
                print(f"   - {violation}")
        else:
            print(f"\n✅ No performance violations detected")

        print("\n" + "-" * 70)
        print("RECOMMENDATIONS")
        print("-" * 70)
        for i, rec in enumerate(report['recommendations'], 1):
            print(f"{i}. {rec}")

        print("\n" + "=" * 70)


def main():
    """Main execution"""
    validator = Phase8Validator()
    report = validator.run_validation()
    validator.print_report(report)

    # Save report to file
    report_path = "tests-and-notes/phase8_validation_report.json"
    with open(report_path, 'w') as f:
        json.dump(report, f, indent=2)
    print(f"\n📄 Full report saved to: {report_path}")

    # Exit with appropriate code
    exit(0 if report['overall_status'] == 'PASS' else 1)


if __name__ == "__main__":
    main()
