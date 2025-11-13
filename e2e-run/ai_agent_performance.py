#!/usr/bin/env python3
"""
AI Agent E2E Performance Profiler

Tracks detailed performance metrics for AI agent workflows and tool calls.
Collects latency, success rates, routing compliance, and system metrics.
"""

import time
import json
import os
from datetime import datetime
from typing import Dict, List, Any, Optional
from dataclasses import dataclass, asdict
import statistics

@dataclass
class ToolCallMetric:
    """Individual tool call performance metric"""
    tool_name: str
    duration_seconds: float
    success: bool
    retries: int
    timestamp: str
    error_type: Optional[str] = None
    result_size_bytes: Optional[int] = None
    routing_compliant: bool = True

@dataclass
class ScenarioMetric:
    """Complete scenario performance metric"""
    scenario_name: str
    start_time: str
    end_time: str
    duration_seconds: float
    tool_count: int
    success: bool
    success_rate: float
    routing_rate: float
    performance_score: float
    tool_metrics: List[ToolCallMetric]

@dataclass
class SystemMetric:
    """System-level performance metrics"""
    total_duration: float
    total_tool_calls: int
    successful_calls: int
    average_latency: float
    p95_latency: float
    throughput_per_second: float
    routing_compliance_rate: float
    retry_success_rate: float

class PerformanceProfiler:
    """Main performance profiling system"""

    def __init__(self, output_dir: str):
        self.output_dir = output_dir
        self.tool_metrics: List[ToolCallMetric] = []
        self.scenario_metrics: List[ScenarioMetric] = []
        self.scenario_start_times: Dict[str, float] = {}
        self.scenario_tool_counts: Dict[str, int] = {}

        print(f"📊 Performance profiler initialized")
        print(f"   Output directory: {self.output_dir}")

    def start_scenario(self, scenario_name: str):
        """Start timing a scenario"""
        self.scenario_start_times[scenario_name] = time.time()
        self.scenario_tool_counts[scenario_name] = 0

    def end_scenario(
        self,
        scenario_name: str,
        success: bool,
        duration: float,
        tool_count: int,
        success_rate: float,
        routing_rate: float
    ):
        """End timing a scenario and record metrics"""
        if scenario_name not in self.scenario_start_times:
            print(f"⚠️  Scenario '{scenario_name}' start time not recorded")
            start_time = time.time() - duration
        else:
            start_time = self.scenario_start_times[scenario_name]

        # Collect tool metrics for this scenario
        scenario_tool_metrics = [
            metric for metric in self.tool_metrics
            if metric.tool_name in [m.name for m in self.scenario_metrics]  # Simplified for now
        ]

        # Filter by recent time window (last 5 minutes)
        cutoff_time = time.time() - 300
        recent_tool_metrics = [
            metric for metric in self.tool_metrics
            if datetime.fromisoformat(metric.timestamp.replace('Z', '+00:00')).timestamp() > cutoff_time
        ]

        # If we can't find specific metrics, use all recent ones
        if not scenario_tool_metrics:
            scenario_tool_metrics = recent_tool_metrics

        scenario_metric = ScenarioMetric(
            scenario_name=scenario_name,
            start_time=datetime.fromtimestamp(start_time).isoformat(),
            end_time=datetime.fromtimestamp(time.time()).isoformat(),
            duration_seconds=duration,
            tool_count=tool_count,
            success=success,
            success_rate=success_rate,
            routing_rate=routing_rate,
            performance_score=success_rate * routing_rate * 100,
            tool_metrics=scenario_tool_metrics
        )

        self.scenario_metrics.append(scenario_metric)

        # Clean up
        if scenario_name in self.scenario_start_times:
            del self.scenario_start_times[scenario_name]
        if scenario_name in self.scenario_tool_counts:
            del self.scenario_tool_counts[scenario_name]

    def record_tool_call(
        self,
        tool_name: str,
        duration: float,
        success: bool,
        retries: int = 0,
        error: Optional[str] = None,
        result_size: Optional[int] = None,
        routing_compliant: bool = True
    ):
        """Record performance metric for a single tool call"""

        # Determine error type if present
        error_type = None
        if error:
            if "timeout" in error.lower():
                error_type = "timeout"
            elif "authentication" in error.lower() or "api key" in error.lower():
                error_type = "auth"
            elif "not found" in error.lower() or "404" in error:
                error_type = "not_found"
            elif "connection" in error.lower():
                error_type = "connection"
            else:
                error_type = "other"

        metric = ToolCallMetric(
            tool_name=tool_name,
            duration_seconds=duration,
            success=success,
            retries=retries,
            timestamp=datetime.now().isoformat(),
            error_type=error_type,
            result_size_bytes=result_size,
            routing_compliant=routing_compliant
        )

        self.tool_metrics.append(metric)

    def get_tool_performance_summary(self) -> Dict[str, Any]:
        """Get performance summary for all tool calls"""
        if not self.tool_metrics:
            return {"message": "No tool metrics recorded"}

        # Group by tool name
        tool_groups = {}
        for metric in self.tool_metrics:
            tool_name = metric.tool_name
            if tool_name not in tool_groups:
                tool_groups[tool_name] = []
            tool_groups[tool_name].append(metric)

        # Calculate summary statistics
        summary = {}
        for tool_name, metrics in tool_groups.items():
            durations = [m.duration_seconds for m in metrics]
            success_count = sum(1 for m in metrics if m.success)
            retry_counts = [m.retries for m in metrics]

            summary[tool_name] = {
                "total_calls": len(metrics),
                "successful_calls": success_count,
                "success_rate": (success_count / len(metrics)) * 100,
                "average_duration": statistics.mean(durations),
                "min_duration": min(durations),
                "max_duration": max(durations),
                "p95_duration": self._calculate_percentile(durations, 95),
                "total_retries": sum(retry_counts),
                "average_retries": statistics.mean(retry_counts) if retry_counts else 0,
                "error_types": self._get_error_type_breakdown(metrics),
                "routing_compliant_rate": (sum(1 for m in metrics if m.routing_compliant) / len(metrics)) * 100
            }

        return summary

    def get_scenario_performance_summary(self) -> Dict[str, Any]:
        """Get performance summary for all scenarios"""
        if not self.scenario_metrics:
            return {"message": "No scenario metrics recorded"}

        summary = {}
        for metric in self.scenario_metrics:
            summary[metric.scenario_name] = {
                "duration_seconds": metric.duration_seconds,
                "tool_count": metric.tool_count,
                "success": metric.success,
                "success_rate": metric.success_rate * 100,
                "routing_rate": metric.routing_rate * 100,
                "performance_score": metric.performance_score,
                "average_tool_latency": statistics.mean([m.duration_seconds for m in metric.tool_metrics]) if metric.tool_metrics else 0
            }

        return summary

    def get_system_performance_metrics(self) -> SystemMetric:
        """Calculate system-wide performance metrics"""
        if not self.tool_metrics:
            return SystemMetric(0, 0, 0, 0, 0, 0, 0, 0)

        durations = [m.duration_seconds for m in self.tool_metrics]
        successful_calls = sum(1 for m in self.tool_metrics if m.success)
        total_retries = sum(m.retries for m in self.tool_metrics)
        routing_compliant = sum(1 for m in self.tool_metrics if m.routing_compliant)

        # Calculate retry success rate
        total_call_attempts = len(self.tool_metrics) + total_retries
        retry_success_rate = (successful_calls / total_call_attempts * 100) if total_call_attempts > 0 else 0

        # Estimate throughput (simplified)
        total_duration = sum(durations)
        throughput = (len(self.tool_metrics) / total_duration) if total_duration > 0 else 0

        return SystemMetric(
            total_duration=total_duration,
            total_tool_calls=len(self.tool_metrics),
            successful_calls=successful_calls,
            average_latency=statistics.mean(durations) if durations else 0,
            p95_latency=self._calculate_percentile(durations, 95) if durations else 0,
            throughput_per_second=throughput,
            routing_compliance_rate=(routing_compliant / len(self.tool_metrics)) * 100 if self.tool_metrics else 0,
            retry_success_rate=retry_success_rate
        )

    def generate_performance_charts(self, results: Dict[str, Any]):
        """Generate performance charts (simplified ASCII)"""
        if not self.tool_metrics:
            return

        charts = []

        # Tool success rate chart
        charts.append("=" * 60)
        charts.append("TOOL SUCCESS RATES")
        charts.append("=" * 60)

        tool_summary = self.get_tool_performance_summary()
        for tool_name, stats in tool_summary.items():
            bar_length = int(stats["success_rate"] / 2)  # Scale for display
            bar = "█" * bar_length + "░" * (50 - bar_length)
            charts.append(f"{tool_name[:20]:<20} |{bar}| {stats['success_rate']:.1f}%")

        # Tool latency chart
        charts.append("\n" + "=" * 60)
        charts.append("TOOL LATENCY (AVERAGE SECONDS)")
        charts.append("=" * 60)

        for tool_name, stats in tool_summary.items():
            bar_length = int(stats["average_duration"] * 10)  # Scale for display
            bar = "█" * min(bar_length, 50)
            charts.append(f"{tool_name[:20]:<20} |{bar}| {stats['average_duration']:.2f}s")

        # Save chart to file
        chart_file = os.path.join(self.output_dir, "performance_charts.txt")
        with open(chart_file, "w") as f:
            f.write("\n".join(charts))

        print(f"📈 Performance charts saved to: {chart_file}")

    def export_detailed_metrics(self, results: Dict[str, Any]):
        """Export detailed performance metrics to JSON"""
        export_data = {
            "export_timestamp": datetime.now().isoformat(),
            "test_run_metadata": results.get("test_run", {}),
            "performance_summary": {
                "system_metrics": asdict(self.get_system_performance_metrics()),
                "tool_performance": self.get_tool_performance_summary(),
                "scenario_performance": self.get_scenario_performance_summary()
            },
            "raw_metrics": {
                "tool_metrics": [asdict(metric) for metric in self.tool_metrics],
                "scenario_metrics": [asdict(metric) for metric in self.scenario_metrics]
            }
        }

        metrics_file = os.path.join(self.output_dir, "detailed_metrics.json")
        with open(metrics_file, "w") as f:
            json.dump(export_data, f, indent=2)

        print(f"📄 Detailed metrics exported to: {metrics_file}")

    def _calculate_percentile(self, values: List[float], percentile: int) -> float:
        """Calculate percentile of a list of values"""
        if not values:
            return 0

        sorted_values = sorted(values)
        index = (percentile / 100) * (len(sorted_values) - 1)

        if index.is_integer():
            return sorted_values[int(index)]
        else:
            lower_index = int(index)
            upper_index = min(lower_index + 1, len(sorted_values) - 1)
            weight = index - lower_index
            return sorted_values[lower_index] * (1 - weight) + sorted_values[upper_index] * weight

    def _get_error_type_breakdown(self, metrics: List[ToolCallMetric]) -> Dict[str, int]:
        """Get breakdown of error types from metrics"""
        error_counts = {}
        for metric in metrics:
            if not metric.success and metric.error_type:
                error_counts[metric.error_type] = error_counts.get(metric.error_type, 0) + 1
        return error_counts

    def print_performance_summary(self):
        """Print performance summary to console"""
        if not self.tool_metrics:
            print("📊 No performance metrics available")
            return

        print("\n" + "=" * 60)
        print("PERFORMANCE SUMMARY")
        print("=" * 60)

        system_metrics = self.get_system_performance_metrics()

        print(f"Total Tool Calls: {system_metrics.total_tool_calls}")
        print(f"Successful Calls: {system_metrics.successful_calls} ({system_metrics.total_tool_calls and system_metrics.successful_calls/system_metrics.total_tool_calls*100:.1f}%)")
        print(f"Average Latency: {system_metrics.average_latency:.2f}s")
        print(f"P95 Latency: {system_metrics.p95_latency:.2f}s")
        print(f"Throughput: {system_metrics.throughput_per_second:.2f} calls/sec")
        print(f"Routing Compliance: {system_metrics.routing_compliance_rate:.1f}%")
        print(f"Retry Success Rate: {system_metrics.retry_success_rate:.1f}%")

        # Top 5 slowest tools
        tool_summary = self.get_tool_performance_summary()
        slowest_tools = sorted(tool_summary.items(), key=lambda x: x[1]["average_duration"], reverse=True)[:5]

        print(f"\n🐌 TOP 5 SLOWEST TOOLS:")
        for tool_name, stats in slowest_tools:
            print(f"  {tool_name:<25} {stats['average_duration']:.2f}s avg")

        # Error types
        all_error_counts = {}
        for metrics in tool_summary.values():
            for error_type, count in metrics["error_types"].items():
                all_error_counts[error_type] = all_error_counts.get(error_type, 0) + count

        if all_error_counts:
            print(f"\n❌ ERROR BREAKDOWN:")
            for error_type, count in sorted(all_error_counts.items(), key=lambda x: x[1], reverse=True):
                print(f"  {error_type:<20} {count} occurrences")
