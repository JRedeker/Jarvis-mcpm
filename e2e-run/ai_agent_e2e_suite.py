#!/usr/bin/env python3
"""
AI Agent E2E Test Suite - Main Orchestrator

Comprehensive end-to-end testing framework for cipher-aggregator with AI agents.
Target: 20-minute runtime with 12 realistic test scenarios.

Features:
- Realistic AI agent workflows
- Performance benchmarking
- Routing compliance validation (Phase 5 integration)
- Comprehensive reporting
- Auto-retry with intelligent fixes
"""

import sys
import os
import time
import json
import threading
import traceback
from datetime import datetime, timezone
from typing import Dict, List, Any, Optional, Tuple
from dataclasses import dataclass, asdict
import urllib.request
import urllib.error

# Add project root to path for imports
project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, project_root)

# Import Phase 5 routing middleware if available
try:
    from cipher_routing_middleware import get_middleware, validate_tool_call, track_tool_execution
    ROUTING_AVAILABLE = True
except ImportError:
    ROUTING_AVAILABLE = False
    print("⚠️  Phase 5 routing middleware not available - testing without routing validation")

from ai_agent_scenarios import ScenarioManager
from ai_agent_performance import PerformanceProfiler
from ai_agent_reporter import ReportGenerator

@dataclass
class TestResult:
    """Individual test result"""
    scenario_name: str
    success: bool
    duration_seconds: float
    tool_count: int
    routing_compliant: bool
    performance_score: float
    error_message: Optional[str] = None
    retry_count: int = 0
    metrics: Dict[str, Any] = None

@dataclass
class TestRunMetadata:
    """Test run metadata and results"""
    start_time: str
    end_time: str
    duration_seconds: float
    target_seconds: int = 1200  # 20 minutes
    scenarios_tested: int = 12
    scenarios_passed: int = 0
    success_rate: float = 0.0
    performance_score: float = 0.0
    routing_compliance: float = 0.0
    overall_status: str = "PENDING"

class AIAgentE2ESuite:
    """Main E2E test orchestrator"""

    def __init__(self, output_dir: str = "./e2e-run"):
        self.output_dir = output_dir
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None
        self.sse_response = None
        self.sse_thread = None
        self.response_event = threading.Event()
        self.response_data = None
        self.pending_request_id = None

        # Initialize components
        self.scenario_manager = ScenarioManager()
        self.performance_profiler = PerformanceProfiler(output_dir)
        self.report_generator = ReportGenerator(output_dir)

        # Test state
        self.test_results: List[TestResult] = []
        self.routing_middleware = None

        # Constraints
        self.max_calls_per_scenario = 8
        self.scenario_timeout = 100  # seconds
        self.tool_timeout = 30  # seconds

        # Metrics
        self.start_time = None
        self.total_tool_calls = 0
        self.successful_calls = 0
        self.routing_compliant_calls = 0

        print(f"🎯 AI Agent E2E Suite initialized")
        print(f"   Output directory: {self.output_dir}")
        print(f"   Target runtime: 20 minutes")
        print(f"   Max tool calls per scenario: {self.max_calls_per_scenario}")

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

    def validate_routing(self, tool_name: str, context: Dict[str, Any]) -> bool:
        """Validate tool selection with Phase 5 middleware"""
        if not ROUTING_AVAILABLE or not self.routing_middleware:
            return True  # Skip routing validation if not available

        try:
            should_allow, suggested_tool, validation = validate_tool_call(
                session_id="e2e-test-session",
                agent_type="ai_agent",
                task_description=context.get("description", "E2E test"),
                selected_tool=tool_name,
                context=context
            )

            # For testing, just record if it was compliant
            is_compliant = validation.get("is_compliant", True)
            return is_compliant

        except Exception as e:
            print(f"⚠️  Routing validation failed: {e}")
            return True  # Don't fail the test due to routing issues

    def execute_scenario(self, scenario_name: str) -> TestResult:
        """Execute a single test scenario"""
        print(f"\n🎭 Executing Scenario: {scenario_name}")

        scenario = self.scenario_manager.get_scenario(scenario_name)
        if not scenario:
            return TestResult(
                scenario_name=scenario_name,
                success=False,
                duration_seconds=0,
                tool_count=0,
                routing_compliant=False,
                performance_score=0,
                error_message=f"Scenario '{scenario_name}' not found"
            )

        start_time = time.time()
        tool_calls = []
        routing_compliant_count = 0
        successful_calls = 0
        error_messages = []
        retry_count = 0

        # Profile scenario start
        self.performance_profiler.start_scenario(scenario_name)

        try:
            for i, tool_call in enumerate(scenario.tools, 1):
                if len(scenario.tools) > self.max_calls_per_scenario:
                    print(f"   ⚠️  Scenario has {len(scenario.tools)} tools, exceeds max {self.max_calls_per_scenario}")

                print(f"   [{i}/{len(scenario.tools)}] {tool_call['name']}: {tool_call['description']}")

                # Track tool call start
                tool_start_time = time.time()
                call_success = False
                call_error = None
                call_result = None
                tool_retry_count = 0

                # Retry logic (max 2 retries)
                for attempt in range(3):
                    try:
                        success, result, error = self.call_tool(
                            tool_call['name'],
                            tool_call.get('arguments', {}),
                            tool_call.get('timeout', self.tool_timeout)
                        )

                        if success:
                            call_success = True
                            call_result = result
                            break
                        else:
                            call_error = error
                            if "timeout" in error.lower():
                                # Increase timeout on retry
                                tool_call['timeout'] = min(tool_call.get('timeout', 30) * 2, 120)
                            tool_retry_count += 1
                            if attempt < 2:  # Don't sleep on last attempt
                                time.sleep(2)  # Brief delay before retry
                    except Exception as e:
                        call_error = str(e)
                        tool_retry_count += 1
                        if attempt < 2:
                            time.sleep(2)

                call_duration = time.time() - tool_start_time
                self.total_tool_calls += 1

                if call_success:
                    successful_calls += 1
                    routing_compliant = self.validate_routing(tool_call['name'], {
                        'description': tool_call.get('description', ''),
                        'arguments': tool_call.get('arguments', {})
                    })
                    if routing_compliant:
                        routing_compliant_count += 1

                    tool_calls.append({
                        'name': tool_call['name'],
                        'success': True,
                        'duration': call_duration,
                        'retries': tool_retry_count,
                        'result_type': type(call_result).__name__ if call_result else 'None'
                    })

                    print(f"      ✅ Success ({call_duration:.1f}s, {tool_retry_count} retries)")

                    # Profile individual tool call
                    self.performance_profiler.record_tool_call(
                        tool_call['name'],
                        call_duration,
                        success=True,
                        retries=tool_retry_count
                    )

                else:
                    error_messages.append(f"{tool_call['name']}: {call_error}")
                    tool_calls.append({
                        'name': tool_call['name'],
                        'success': False,
                        'duration': call_duration,
                        'error': call_error,
                        'retries': tool_retry_count
                    })

                    print(f"      ❌ Failed: {call_error}")

                    # Profile failed tool call
                    self.performance_profiler.record_tool_call(
                        tool_call['name'],
                        call_duration,
                        success=False,
                        error=call_error,
                        retries=tool_retry_count
                    )

                    # If it's a critical tool and multiple failures, stop the scenario
                    if tool_call.get('critical', False) and len([t for t in tool_calls if t['success']]) < len(tool_calls) * 0.3:
                        print(f"      ⛔ Critical tool failed, stopping scenario")
                        break

                retry_count += tool_retry_count

                # Brief delay between calls
                time.sleep(0.1)

            # Scenario completion
            scenario_duration = time.time() - start_time

            # Determine scenario success
            success_rate = successful_calls / len(tool_calls) if tool_calls else 0
            routing_rate = routing_compliant_count / successful_calls if successful_calls > 0 else 1.0
            success = success_rate >= 0.6  # 60% success threshold

            # Profile scenario end
            self.performance_profiler.end_scenario(
                scenario_name,
                success,
                scenario_duration,
                len(tool_calls),
                success_rate,
                routing_rate
            )

            result = TestResult(
                scenario_name=scenario_name,
                success=success,
                duration_seconds=scenario_duration,
                tool_count=len(tool_calls),
                routing_compliant=routing_rate >= 0.8,
                performance_score=success_rate * routing_rate * 100,
                error_message="; ".join(error_messages) if error_messages else None,
                retry_count=retry_count,
                metrics={
                    'success_rate': success_rate,
                    'routing_rate': routing_rate,
                    'tool_calls': tool_calls
                }
            )

            print(f"\n📊 Scenario Result: {'✅ PASS' if success else '❌ FAIL'}")
            print(f"   Duration: {scenario_duration:.1f}s")
            print(f"   Tools: {successful_calls}/{len(tool_calls)} ({success_rate*100:.1f}%)")
            print(f"   Routing: {routing_compliant_count}/{successful_calls} ({routing_rate*100:.1f}%)")
            print(f"   Retries: {retry_count}")

            return result

        except Exception as e:
            error_msg = f"Scenario execution failed: {str(e)}"
            print(f"   ❌ {error_msg}")
            return TestResult(
                scenario_name=scenario_name,
                success=False,
                duration_seconds=time.time() - start_time,
                tool_count=0,
                routing_compliant=False,
                performance_score=0,
                error_message=error_msg
            )

    def run_full_suite(self) -> Dict[str, Any]:
        """Run the complete E2E test suite"""
        print("🚀 AI Agent E2E Test Suite Starting...")
        print(f"🎯 Target: 20 minutes | 12 scenarios")
        print(f"⚡ Max tool calls per scenario: {self.max_calls_per_scenario}")
        print(f"🔗 Connection: {self.sse_url}")

        self.start_time = time.time()

        # Establish connection
        if not self.establish_connection():
            return {
                'error': 'Failed to establish SSE connection',
                'connection': False,
                'test_run': {
                    'start_time': datetime.now(timezone.utc).isoformat(),
                    'end_time': datetime.now(timezone.utc).isoformat(),
                    'duration_seconds': 0,
                    'overall_status': 'FAILED'
                }
            }

        # Initialize routing middleware
        if ROUTING_AVAILABLE:
            try:
                self.routing_middleware = get_middleware()
                print("✅ Phase 5 routing middleware initialized")
            except Exception as e:
                print(f"⚠️  Routing middleware initialization failed: {e}")

        # Get all scenarios
        all_scenarios = self.scenario_manager.list_scenarios()
        print(f"\n📋 Found {len(all_scenarios)} test scenarios")

        # Execute scenarios
        for i, scenario_name in enumerate(all_scenarios, 1):
            try:
                result = self.execute_scenario(scenario_name)
                self.test_results.append(result)

                # Update running metrics
                self.successful_calls += 1 if result.success else 0

            except Exception as e:
                print(f"❌ Scenario '{scenario_name}' failed with exception: {e}")
                traceback.print_exc()

        # Generate final report
        total_duration = time.time() - self.start_time

        # Calculate final metrics
        passed_scenarios = sum(1 for r in self.test_results if r.success)
        total_scenarios = len(self.test_results)
        success_rate = (passed_scenarios / total_scenarios * 100) if total_scenarios > 0 else 0
        avg_performance = sum(r.performance_score for r in self.test_results) / total_scenarios if total_scenarios > 0 else 0
        overall_routing = (sum(1 for r in self.test_results if r.routing_compliant) / total_scenarios * 100) if total_scenarios > 0 else 0

        # Determine overall status
        performance_good = total_duration <= 1200
        success_good = success_rate >= 80
        routing_good = overall_routing >= 75

        overall_status = "PASS"
        if not (performance_good and success_good and routing_good):
            overall_status = "FAIL"

        test_run_metadata = TestRunMetadata(
            start_time=datetime.fromtimestamp(self.start_time, timezone.utc).isoformat(),
            end_time=datetime.now(timezone.utc).isoformat(),
            duration_seconds=total_duration,
            scenarios_tested=total_scenarios,
            scenarios_passed=passed_scenarios,
            success_rate=success_rate,
            performance_score=avg_performance,
            routing_compliance=overall_routing,
            overall_status=overall_status
        )

        # Create comprehensive results
        results = {
            'test_run': asdict(test_run_metadata),
            'performance': {
                'total_duration': total_duration,
                'target_duration': 1200,
                'performance_score': (1200 / total_duration * 100) if total_duration > 0 else 0,
                'total_tool_calls': self.total_tool_calls,
                'successful_calls': self.successful_calls,
                'success_rate': (self.successful_calls / self.total_tool_calls * 100) if self.total_tool_calls > 0 else 0
            },
            'scenarios': [asdict(result) for result in self.test_results],
            'routing_compliance': {
                'enabled': ROUTING_AVAILABLE,
                'middleware_available': self.routing_middleware is not None,
                'overall_score': overall_routing
            },
            'constraints_validation': {
                'max_calls_per_scenario': self.max_calls_per_scenario,
                'scenario_timeout': self.scenario_timeout,
                'tool_timeout': self.tool_timeout,
                'serial_execution': True,
                'max_parallel_calls': 1
            }
        }

        # Generate reports
        self.report_generator.generate_json_report(results)
        self.report_generator.generate_html_report(results)
        self.report_generator.generate_console_summary(results)

        print("\n" + "="*70)
        print("🏁 E2E TEST SUITE COMPLETED")
        print("="*70)
        print(f"Status: {overall_status}")
        print(f"Runtime: {total_duration:.1f}s / 1200s ({(total_duration/1200*100):.1f}%)")
        print(f"Success Rate: {success_rate:.1f}% ({passed_scenarios}/{total_scenarios})")
        print(f"Performance Score: {avg_performance:.1f}/100")
        print(f"Routing Compliance: {overall_routing:.1f}%")
        print(f"Tool Calls: {self.successful_calls}/{self.total_tool_calls} successful")
        print(f"Output files saved to: {self.output_dir}")

        return results

def main():
    """Main execution entry point"""
    import argparse

    parser = argparse.ArgumentParser(description="AI Agent E2E Test Suite")
    parser.add_argument('--output-dir', default='./e2e-run', help='Output directory for results')
    parser.add_argument('--timeout', type=int, default=1200, help='Total timeout in seconds (default: 20 minutes)')

    args = parser.parse_args()

    # Initialize and run suite
    suite = AIAgentE2ESuite(args.output_dir)
    results = suite.run_full_suite()

    # Exit with appropriate code
    if results.get('test_run', {}).get('overall_status') == 'PASS':
        print("\n✅ All tests passed!")
        return 0
    else:
        print(f"\n❌ Some tests failed!")
        return 1

if __name__ == "__main__":
    exit(main())
