#!/usr/bin/env python3
"""
AI Agent E2E Report Generator

Generates comprehensive reports from E2E test results including:
- JSON machine-readable reports
- HTML human-readable dashboard
- Console summary reports
- Performance analysis
"""

import json
import os
from datetime import datetime
from typing import Dict, List, Any
from pathlib import Path

class ReportGenerator:
    """Generate comprehensive test reports"""

    def __init__(self, output_dir: str):
        self.output_dir = output_dir
        self.report_timestamp = datetime.now().isoformat()

        # Ensure output directory exists
        os.makedirs(self.output_dir, exist_ok=True)

        print(f"📋 Report generator initialized")
        print(f"   Output directory: {self.output_dir}")

    def generate_json_report(self, results: Dict[str, Any]):
        """Generate machine-readable JSON report"""
        report = {
            "report_metadata": {
                "generator": "AI Agent E2E Test Suite",
                "version": "1.0.0",
                "timestamp": self.report_timestamp,
                "output_directory": self.output_dir
            },
            **results
        }

        json_file = os.path.join(self.output_dir, "ai_agent_e2e_report.json")
        with open(json_file, "w") as f:
            json.dump(report, f, indent=2)

        print(f"📄 JSON report saved to: {json_file}")
        return json_file

    def generate_html_report(self, results: Dict[str, Any]):
        """Generate human-readable HTML dashboard"""
        html_content = self._generate_html_template(results)

        html_file = os.path.join(self.output_dir, "ai_agent_e2e_report.html")
        with open(html_file, "w") as f:
            f.write(html_content)

        print(f"🌐 HTML report saved to: {html_file}")
        return html_file

    def generate_console_summary(self, results: Dict[str, Any]):
        """Generate console summary report"""
        summary_lines = []

        summary_lines.append("=" * 70)
        summary_lines.append("AI AGENT E2E TEST SUITE - EXECUTIVE SUMMARY")
        summary_lines.append("=" * 70)

        # Overall status
        test_run = results.get("test_run", {})
        overall_status = test_run.get("overall_status", "UNKNOWN")
        status_icon = "✅" if overall_status == "PASS" else "❌"

        summary_lines.append(f"{status_icon} Overall Status: {overall_status}")
        summary_lines.append(f"⏱️  Runtime: {test_run.get('duration_seconds', 0):.1f}s / {test_run.get('target_seconds', 1200)}s")
        summary_lines.append(f"📊 Success Rate: {test_run.get('success_rate', 0):.1f}% ({test_run.get('scenarios_passed', 0)}/{test_run.get('scenarios_tested', 0)})")
        summary_lines.append(f"🎯 Performance Score: {test_run.get('performance_score', 0):.1f}/100")
        summary_lines.append(f"🔄 Routing Compliance: {test_run.get('routing_compliance', 0):.1f}%")

        # Performance metrics
        performance = results.get("performance", {})
        summary_lines.append("")
        summary_lines.append("PERFORMANCE METRICS")
        summary_lines.append("-" * 70)
        summary_lines.append(f"Total Tool Calls: {performance.get('total_tool_calls', 0)}")
        summary_lines.append(f"Successful Calls: {performance.get('successful_calls', 0)} ({performance.get('success_rate', 0):.1f}%)")
        summary_lines.append(f"Performance Score: {performance.get('performance_score', 0):.1f}/100")

        # Scenario breakdown
        scenarios = results.get("scenarios", [])
        if scenarios:
            summary_lines.append("")
            summary_lines.append("SCENARIO RESULTS")
            summary_lines.append("-" * 70)

            for scenario in scenarios:
                status_icon = "✅" if scenario.get("success") else "❌"
                routing_icon = "🔄" if scenario.get("routing_compliant") else "⚠️"
                summary_lines.append(
                    f"{status_icon} {scenario.get('scenario_name', 'Unknown'):<30} "
                    f"{scenario.get('duration_seconds', 0):.1f}s "
                    f"{routing_icon} {scenario.get('performance_score', 0):.1f}"
                )

        # Recommendations
        summary_lines.append("")
        summary_lines.append("RECOMMENDATIONS")
        summary_lines.append("-" * 70)

        recommendations = self._generate_recommendations(results)
        for i, rec in enumerate(recommendations, 1):
            summary_lines.append(f"{i}. {rec}")

        summary_lines.append("")
        summary_lines.append("=" * 70)
        summary_lines.append(f"Generated: {self.report_timestamp}")
        summary_lines.append("=" * 70)

        # Save console summary
        summary_file = os.path.join(self.output_dir, "console_summary.txt")
        with open(summary_file, "w") as f:
            f.write("\n".join(summary_lines))

        print(f"📋 Console summary saved to: {summary_file}")

        # Also print to console
        print("\n" + "\n".join(summary_lines))

    def _generate_html_template(self, results: Dict[str, Any]) -> str:
        """Generate HTML template with embedded CSS and charts"""

        test_run = results.get("test_run", {})
        performance = results.get("performance", {})
        scenarios = results.get("scenarios", [])

        # Calculate success/failure counts
        passed_scenarios = [s for s in scenarios if s.get("success")]
        failed_scenarios = [s for s in scenarios if not s.get("success")]

        html_template = f"""
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Agent E2E Test Report</title>
    <style>
        body {{
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }}
        .container {{
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }}
        .header {{
            text-align: center;
            border-bottom: 3px solid #007acc;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }}
        .header h1 {{
            color: #007acc;
            margin: 0;
            font-size: 2.5em;
        }}
        .header .timestamp {{
            color: #666;
            font-size: 0.9em;
            margin-top: 10px;
        }}
        .metrics-grid {{
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin: 30px 0;
        }}
        .metric-card {{
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            border-left: 4px solid #007acc;
        }}
        .metric-value {{
            font-size: 2em;
            font-weight: bold;
            color: #007acc;
        }}
        .metric-label {{
            color: #666;
            font-size: 0.9em;
            margin-top: 5px;
        }}
        .status-pass {{
            color: #28a745;
        }}
        .status-fail {{
            color: #dc3545;
        }}
        .scenario-table {{
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }}
        .scenario-table th, .scenario-table td {{
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }}
        .scenario-table th {{
            background-color: #007acc;
            color: white;
        }}
        .scenario-table tr:nth-child(even) {{
            background-color: #f2f2f2;
        }}
        .progress-bar {{
            background-color: #e9ecef;
            border-radius: 10px;
            overflow: hidden;
            height: 20px;
        }}
        .progress-fill {{
            background-color: #007acc;
            height: 100%;
            transition: width 0.3s ease;
        }}
        .chart-container {{
            margin: 30px 0;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
        }}
        .recommendations {{
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }}
        .recommendations h3 {{
            color: #856404;
            margin-top: 0;
        }}
        .footer {{
            text-align: center;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #ddd;
            color: #666;
            font-size: 0.9em;
        }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 AI Agent E2E Test Report</h1>
            <div class="timestamp">Generated: {test_run.get('start_time', 'Unknown')}</div>
        </div>

        <!-- Overall Metrics -->
        <div class="metrics-grid">
            <div class="metric-card">
                <div class="metric-value status-{'pass' if test_run.get('overall_status') == 'PASS' else 'fail'}">{test_run.get('overall_status', 'UNKNOWN')}</div>
                <div class="metric-label">Overall Status</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{test_run.get('duration_seconds', 0):.1f}s</div>
                <div class="metric-label">Total Runtime</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{test_run.get('success_rate', 0):.1f}%</div>
                <div class="metric-label">Success Rate</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{test_run.get('performance_score', 0):.1f}/100</div>
                <div class="metric-label">Performance Score</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{test_run.get('routing_compliance', 0):.1f}%</div>
                <div class="metric-label">Routing Compliance</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{performance.get('total_tool_calls', 0)}</div>
                <div class="metric-label">Tool Calls</div>
            </div>
        </div>

        <!-- Scenario Results -->
        <h2>📊 Scenario Results</h2>
        <table class="scenario-table">
            <thead>
                <tr>
                    <th>Scenario</th>
                    <th>Status</th>
                    <th>Duration</th>
                    <th>Tool Count</th>
                    <th>Success Rate</th>
                    <th>Performance</th>
                    <th>Routing</th>
                </tr>
            </thead>
            <tbody>
"""

        for scenario in scenarios:
            status_class = "status-pass" if scenario.get("success") else "status-fail"
            status_text = "PASS" if scenario.get("success") else "FAIL"
            routing_icon = "✅" if scenario.get("routing_compliant") else "⚠️"

            html_template += f"""
                <tr>
                    <td>{scenario.get('scenario_name', 'Unknown')}</td>
                    <td class="{status_class}">{status_text}</td>
                    <td>{scenario.get('duration_seconds', 0):.1f}s</td>
                    <td>{scenario.get('tool_count', 0)}</td>
                    <td>{scenario.get('metrics', {}).get('success_rate', 0)*100:.1f}%</td>
                    <td>{scenario.get('performance_score', 0):.1f}/100</td>
                    <td>{routing_icon}</td>
                </tr>
"""

        html_template += f"""
            </tbody>
        </table>

        <!-- Performance Charts -->
        <div class="chart-container">
            <h2>📈 Performance Overview</h2>
            <p><strong>Scenarios Passed:</strong> {len(passed_scenarios)} | <strong>Failed:</strong> {len(failed_scenarios)}</p>

            <div style="margin: 20px 0;">
                <p><strong>Success Rate Progress</strong></p>
                <div class="progress-bar">
                    <div class="progress-fill" style="width: {test_run.get('success_rate', 0)}%"></div>
                </div>
                <p style="text-align: center; margin-top: 5px;">{test_run.get('success_rate', 0):.1f}%</p>
            </div>

            <div style="margin: 20px 0;">
                <p><strong>Performance Score</strong></p>
                <div class="progress-bar">
                    <div class="progress-fill" style="width: {test_run.get('performance_score', 0)}%"></div>
                </div>
                <p style="text-align: center; margin-top: 5px;">{test_run.get('performance_score', 0):.1f}/100</p>
            </div>
        </div>

        <!-- Recommendations -->
        <div class="recommendations">
            <h3>💡 Recommendations</h3>
            <ul>
"""

        recommendations = self._generate_recommendations(results)
        for rec in recommendations:
            html_template += f"                <li>{rec}</li>\n"

        html_template += f"""
            </ul>
        </div>

        <!-- Technical Details -->
        <h2>🔧 Technical Details</h2>
        <div class="metrics-grid">
            <div class="metric-card">
                <div class="metric-value">{performance.get('successful_calls', 0)}</div>
                <div class="metric-label">Successful Tool Calls</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{len(scenarios)}</div>
                <div class="metric-label">Total Scenarios</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{len(passed_scenarios)}</div>
                <div class="metric-label">Passed Scenarios</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{len(failed_scenarios)}</div>
                <div class="metric-label">Failed Scenarios</div>
            </div>
        </div>

        <div class="footer">
            <p>Generated by AI Agent E2E Test Suite v1.0.0 | Timestamp: {self.report_timestamp}</p>
            <p>Output Directory: {self.output_dir}</p>
        </div>
    </div>
</body>
</html>
"""

        return html_template

    def _generate_recommendations(self, results: Dict[str, Any]) -> List[str]:
        """Generate intelligent recommendations based on test results"""
        recommendations = []

        test_run = results.get("test_run", {})
        scenarios = results.get("scenarios", [])
        performance = results.get("performance", {})

        # Performance recommendations
        if test_run.get("duration_seconds", 0) > 1200:
            recommendations.append("Runtime exceeded 20 minutes. Consider optimizing slow scenarios or reducing test scope.")

        if test_run.get("success_rate", 0) < 80:
            recommendations.append("Success rate below 80%. Review failed scenarios and fix underlying issues.")

        if test_run.get("performance_score", 0) < 70:
            recommendations.append("Performance score below 70%. Investigate slow tools and optimize workflows.")

        if test_run.get("routing_compliance", 0) < 75:
            recommendations.append("Routing compliance below 75%. Review tool selection patterns and update routing rules.")

        # Scenario-specific recommendations
        failed_scenarios = [s for s in scenarios if not s.get("success")]
        if failed_scenarios:
            recommendations.append(f"Focus on fixing {len(failed_scenarios)} failed scenarios: {', '.join([s.get('scenario_name', 'Unknown') for s in failed_scenarios[:3]])}")

        # Tool performance recommendations
        total_calls = performance.get("total_tool_calls", 0)
        successful_calls = performance.get("successful_calls", 0)
        if total_calls > 0:
            tool_success_rate = (successful_calls / total_calls) * 100
            if tool_success_rate < 85:
                recommendations.append(f"Tool call success rate is {tool_success_rate:.1f}%. Check tool availability and API credentials.")

        # Positive recommendations
        if test_run.get("overall_status") == "PASS":
            recommendations.append("✅ All tests passed! The AI agent E2E suite is working correctly.")

        if test_run.get("routing_compliance", 0) > 90:
            recommendations.append("🎯 Excellent routing compliance! Tool selection is optimal.")

        # Default recommendation if none generated
        if not recommendations:
            recommendations.append("Review the detailed metrics to identify areas for improvement.")

        return recommendations

    def generate_performance_comparison(self, previous_results: List[Dict[str, Any]]) -> Dict[str, Any]:
        """Generate performance comparison with previous runs"""
        if not previous_results:
            return {"message": "No previous results available for comparison"}

        latest_results = previous_results[-1]
        comparison = {
            "comparison_timestamp": self.report_timestamp,
            "latest_run": latest_results.get("test_run", {}),
            "improvements": [],
            "regressions": [],
            "trends": {}
        }

        # Compare key metrics
        latest_duration = latest_results.get("test_run", {}).get("duration_seconds", 0)
        latest_success_rate = latest_results.get("test_run", {}).get("success_rate", 0)
        latest_performance_score = latest_results.get("test_run", {}).get("performance_score", 0)

        # If we have multiple previous runs, show trends
        if len(previous_results) > 1:
            prev_duration = previous_results[-2].get("test_run", {}).get("duration_seconds", 0)
            prev_success_rate = previous_results[-2].get("test_run", {}).get("success_rate", 0)
            prev_performance_score = previous_results[-2].get("test_run", {}).get("performance_score", 0)

            # Duration trend
            if latest_duration < prev_duration:
                comparison["improvements"].append(f"Runtime improved by {prev_duration - latest_duration:.1f}s")
            elif latest_duration > prev_duration:
                comparison["regressions"].append(f"Runtime increased by {latest_duration - prev_duration:.1f}s")

            # Success rate trend
            if latest_success_rate > prev_success_rate:
                comparison["improvements"].append(f"Success rate improved by {latest_success_rate - prev_success_rate:.1f}%")
            elif latest_success_rate < prev_success_rate:
                comparison["regressions"].append(f"Success rate decreased by {prev_success_rate - latest_success_rate:.1f}%")

            # Performance score trend
            if latest_performance_score > prev_performance_score:
                comparison["improvements"].append(f"Performance score improved by {latest_performance_score - prev_performance_score:.1f}")
            elif latest_performance_score < prev_performance_score:
                comparison["regressions"].append(f"Performance score decreased by {prev_performance_score - latest_performance_score:.1f}")

        # Save comparison report
        comparison_file = os.path.join(self.output_dir, "performance_comparison.json")
        with open(comparison_file, "w") as f:
            json.dump(comparison, f, indent=2)

        print(f"📊 Performance comparison saved to: {comparison_file}")
        return comparison
