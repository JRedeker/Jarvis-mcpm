#!/usr/bin/env python3
"""
Comprehensive Test Suite for System Prompt Routing Enforcement System

This test suite validates all routing enforcement functionality including:
- Domain-specific routing validation
- Performance constraint enforcement
- Tool selection monitoring
- Agent configuration verification
- Integration with cipher-aggregator systems
"""

import pytest
import tempfile
import sqlite3
import json
import time
import uuid
from datetime import datetime, timedelta
from pathlib import Path
from unittest.mock import patch, MagicMock
import sys
import os

# Add the parent directory to the path to import our module
sys.path.insert(0, str(Path(__file__).parent))

from routing_enforcement_system import (
    RoutingEnforcementEngine,
    RoutingRule,
    ToolSelection,
    PerformanceConstraint,
    PerformanceConstraintViolation,
    RoutingViolation
)


class TestRoutingRule:
    """Test routing rule data structure"""

    def test_routing_rule_creation(self):
        """Test that routing rule can be created properly"""
        rule = RoutingRule(
            domain="github",
            primary_tools=["github"],
            forbidden_tools=["fetch", "curl"],
            fallback_tools=["filesystem"],
            priority_level=1,
            task_patterns=["repository", "github"],
            conflict_resolution="specialized"
        )

        assert rule.domain == "github"
        assert rule.primary_tools == ["github"]
        assert "fetch" in rule.forbidden_tools
        assert rule.priority_level == 1
        assert rule.conflict_resolution == "specialized"

    def test_routing_rule_equality(self):
        """Test routing rule equality comparison"""
        rule1 = RoutingRule(
            domain="github",
            primary_tools=["github"],
            forbidden_tools=["fetch"],
            fallback_tools=["filesystem"],
            priority_level=1,
            task_patterns=["repo"],
            conflict_resolution="specialized"
        )

        rule2 = RoutingRule(
            domain="github",
            primary_tools=["github"],
            forbidden_tools=["fetch"],
            fallback_tools=["filesystem"],
            priority_level=1,
            task_patterns=["repo"],
            conflict_resolution="specialized"
        )

        assert rule1.domain == rule2.domain
        assert rule1.primary_tools == rule2.primary_tools


class TestPerformanceConstraint:
    """Test performance constraint enforcement"""

    def test_default_constraints(self):
        """Test default performance constraints"""
        constraints = PerformanceConstraint()

        assert constraints.max_calls_per_task == 8
        assert constraints.max_parallel_calls == 1
        assert constraints.call_timeout == 45000
        assert constraints.domain_specific_timeout is not None
        assert "github" in constraints.domain_specific_timeout

    def test_custom_constraints(self):
        """Test custom performance constraints"""
        custom_timeouts = {"github": 60000}
        constraints = PerformanceConstraint(
            max_calls_per_task=5,
            max_parallel_calls=2,
            call_timeout=60000,
            domain_specific_timeout=custom_timeouts
        )

        assert constraints.max_calls_per_task == 5
        assert constraints.max_parallel_calls == 2
        assert constraints.call_timeout == 60000
        assert constraints.domain_specific_timeout == custom_timeouts


class TestRoutingEnforcementEngine:
    """Test the main routing enforcement engine"""

    @pytest.fixture
    def temp_engine(self):
        """Create a temporary routing enforcement engine for testing"""
        with tempfile.TemporaryDirectory() as temp_dir:
            engine = RoutingEnforcementEngine(config_path="dummy.yml")
            # Override paths to use temp directory
            engine.rules_db_path = Path(temp_dir) / "rules.db"
            engine.logs_db_path = Path(temp_dir) / "logs.db"
            engine.performance_logs_path = Path(temp_dir) / "performance.jsonl"
            engine._init_databases()
            yield engine

    def test_engine_initialization(self, temp_engine):
        """Test engine initialization"""
        assert temp_engine.routing_rules is not None
        assert "github" in temp_engine.routing_rules
        assert "web_search" in temp_engine.routing_rules
        assert temp_engine.stats['total_validations'] == 0

    def test_task_classification(self, temp_engine):
        """Test task classification into domain categories"""
        # Test GitHub classification
        assert temp_engine._classify_task("create a new repository on github") == "github"
        assert temp_engine._classify_task("search for pull request issues") == "github"

        # Test web search classification
        assert temp_engine._classify_task("search for information online") == "web_search"
        assert temp_engine._classify_task("google the latest news") == "web_search"

        # Test web scraping classification
        assert temp_engine._classify_task("scrape data from website") == "web_scraping"
        assert temp_engine._classify_task("extract content from webpage") == "web_scraping"

        # Test code analysis classification
        assert temp_engine._classify_task("find code patterns in files") == "code_analysis"
        assert temp_engine._classify_task("search for function definitions") == "code_analysis"

        # Test file operations classification
        assert temp_engine._classify_task("read file from disk") == "file_operations"
        assert temp_engine._classify_task("write to local directory") == "file_operations"

        # Test general fallback
        assert temp_engine._classify_task("do something generic") == "general"

    def test_find_applicable_rule(self, temp_engine):
        """Test finding applicable routing rules"""
        # Test exact match
        rule = temp_engine._find_applicable_rule("create repository", "github")
        assert rule is not None
        assert rule.domain == "github"

        # Test pattern match
        rule = temp_engine._find_applicable_rule("search for repo information", "general")
        assert rule is not None
        assert rule.domain == "github"

        # Test no match
        rule = temp_engine._find_applicable_rule("unknown task xyz", "unknown")
        assert rule is None

    def test_valid_tool_selection(self, temp_engine):
        """Test valid tool selection validation"""
        session_id = "test_session"
        task_desc = "create repository on github"
        tool = "github"

        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, task_desc, tool
        )

        assert is_valid is True
        assert recommended_tool == tool
        assert reason is None
        assert temp_engine.stats['total_validations'] == 1

    def test_forbidden_tool_violation(self, temp_engine):
        """Test forbidden tool violation detection"""
        session_id = "test_session"
        task_desc = "create repository on github"
        tool = "fetch"  # Forbidden for github operations

        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, task_desc, tool
        )

        assert is_valid is False
        assert recommended_tool == "github"
        assert "forbidden" in reason.lower()
        assert temp_engine.stats['violations_detected'] == 1

    def test_suboptimal_tool_selection(self, temp_engine):
        """Test suboptimal but acceptable tool selection"""
        session_id = "test_session"
        task_desc = "create repository on github"
        tool = "filesystem"  # Fallback tool, suboptimal

        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, task_desc, tool
        )

        assert is_valid is True  # Acceptable fallback
        assert recommended_tool == "github"  # But better tool recommended
        assert "suboptimal" in reason.lower()

    def test_web_search_validation(self, temp_engine):
        """Test web search specific validation"""
        session_id = "test_session"

        # Test forbidden tool for web search
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "search for information online", "fetch"
        )
        assert is_valid is False

        # Test correct tool
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "search for information online", "brave-search"
        )
        assert is_valid is True

    def test_web_scraping_validation(self, temp_engine):
        """Test web scraping specific validation"""
        session_id = "test_session"

        # Test forbidden tool for web scraping
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "scrape data from website", "curl"
        )
        assert is_valid is False

        # Test correct tool
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "scrape data from website", "firecrawl"
        )
        assert is_valid is True

    def test_code_analysis_validation(self, temp_engine):
        """Test code analysis specific validation"""
        session_id = "test_session"

        # Test forbidden tool for code analysis
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "find code patterns", "manual scanning"
        )
        assert is_valid is False

        # Test correct tool
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "find code patterns", "code-index"
        )
        assert is_valid is True

    def test_api_testing_validation(self, temp_engine):
        """Test API testing specific validation"""
        session_id = "test_session"

        # Test forbidden tool for API testing
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "test API schema", "fetch"
        )
        assert is_valid is False

        # Test correct tool
        is_valid, recommended_tool, reason = temp_engine.validate_tool_selection(
            session_id, "test API schema", "schemathesis"
        )
        assert is_valid is True

    def test_performance_tracking_basic(self, temp_engine):
        """Test basic performance tracking functionality"""
        session_id = "test_session"
        task_id = "test_task"
        tool_name = "github"

        # Test successful tracking
        with temp_engine.track_performance(session_id, task_id, tool_name):
            time.sleep(0.01)  # Small delay to ensure measurable time

        # Check that performance log was created
        assert temp_engine.performance_logs_path.exists()

        with open(temp_engine.performance_logs_path, 'r') as f:
            lines = f.readlines()
            assert len(lines) == 1

            record = json.loads(lines[0])
            assert record['session_id'] == session_id
            assert record['task_id'] == task_id
            assert record['tool_name'] == tool_name
            assert 'duration' in record

    def test_max_calls_constraint(self, temp_engine):
        """Test maximum calls per task constraint"""
        session_id = "test_session"
        task_id = "test_task"
        tool_name = "github"

        # Fill up to the limit
        for i in range(temp_engine.performance_constraints.max_calls_per_task):
            with temp_engine.track_performance(session_id, task_id, tool_name):
                pass

        # Next call should raise exception
        with pytest.raises(PerformanceConstraintViolation):
            with temp_engine.track_performance(session_id, task_id, tool_name):
                pass

    def test_timeout_detection(self, temp_engine):
        """Test timeout violation detection"""
        session_id = "test_session"
        task_id = "test_task"
        tool_name = "github"

        # Mock the timeout to be very small
        original_timeout = temp_engine.performance_constraints.call_timeout
        temp_engine.performance_constraints.call_timeout = 1  # 1ms

        with temp_engine.track_performance(session_id, task_id, tool_name):
            time.sleep(0.01)  # 10ms sleep, exceeds 1ms timeout

        # Restore original timeout
        temp_engine.performance_constraints.call_timeout = original_timeout

        # Check that timeout violation was logged
        assert temp_engine.stats['performance_violations'] > 0

    def test_agent_configuration_verification_valid(self, temp_engine):
        """Test agent configuration verification with valid config"""
        valid_config = {
            'systemPrompt': 'DOMAIN-SPECIFIC routing rules for tool selection',
            'maxParallelCalls': 1,
            'maxCallsPerTask': 8,
            'toolExecution': {'callTimeout': 45000}
        }

        result = temp_engine.verify_agent_configuration(valid_config)

        assert result['is_valid'] is True
        assert len(result['issues']) == 0
        assert len(result['required_actions']) == 0

    def test_agent_configuration_verification_invalid(self, temp_engine):
        """Test agent configuration verification with invalid config"""
        invalid_config = {
            'maxParallelCalls': 5,  # Should be 1
            'maxCallsPerTask': 15,  # Should be <= 8
            # Missing systemPrompt with routing rules
            # Missing toolExecution.callTimeout
        }

        result = temp_engine.verify_agent_configuration(invalid_config)

        assert result['is_valid'] is False
        assert len(result['issues']) > 0
        assert len(result['required_actions']) > 0

        # Check specific issues
        issues_text = ' '.join(result['issues']).lower()
        assert 'parallel' in issues_text
        assert 'calls' in issues_text
        assert 'routing rules' in issues_text
        assert 'timeout' in issues_text

    def test_generate_routing_report(self, temp_engine):
        """Test routing report generation"""
        # Generate some test data
        session_id = "test_session"
        temp_engine.validate_tool_selection(session_id, "github task", "github")
        temp_engine.validate_tool_selection(session_id, "fetch forbidden", "fetch")

        report = temp_engine.generate_routing_report()

        assert 'enforcement_stats' in report
        assert 'recent_selections' in report
        assert 'violations_by_category' in report
        assert 'performance_stats' in report
        assert 'active_sessions' in report
        assert 'generated_at' in report

        assert report['enforcement_stats']['total_validations'] == 2
        assert report['enforcement_stats']['violations_detected'] == 1

    def test_concurrent_sessions(self, temp_engine):
        """Test handling of multiple concurrent sessions"""
        sessions = ["session1", "session2", "session3"]

        for session_id in sessions:
            temp_engine.validate_tool_selection(session_id, "github task", "github")

        assert len(temp_engine.active_sessions) == 3

        # Each session should have its own tracking
        for session_id in sessions:
            assert session_id in temp_engine.active_sessions

    def test_database_operations(self, temp_engine):
        """Test database initialization and operations"""
        # Check that databases were created
        assert temp_engine.rules_db_path.exists()
        assert temp_engine.logs_db_path.exists()

        # Check routing rules database schema
        with sqlite3.connect(str(temp_engine.rules_db_path)) as conn:
            cursor = conn.execute("SELECT name FROM sqlite_master WHERE type='table'")
            tables = [row[0] for row in cursor.fetchall()]
            assert 'routing_rules' in tables
            assert 'tool_selections' in tables

    def test_performance_stats_calculation(self, temp_engine):
        """Test performance statistics calculation"""
        # Create some test performance data
        test_data = [
            {'session_id': 's1', 'task_id': 't1', 'tool_name': 'github', 'duration': 1.0},
            {'session_id': 's1', 'task_id': 't1', 'tool_name': 'github', 'duration': 2.0},
            {'session_id': 's2', 'task_id': 't2', 'tool_name': 'firecrawl', 'duration': 0.5}
        ]

        for record in test_data:
            with open(temp_engine.performance_logs_path, 'a') as f:
                f.write(json.dumps(record) + '\n')

        stats = temp_engine._calculate_performance_stats()

        assert stats['average_duration'] == pytest.approx(1.1667, rel=1e-2)
        assert stats['max_duration'] == 2.0
        assert stats['min_duration'] == 0.5
        assert stats['total_tool_calls'] == 3
        assert 'github' in stats['tool_usage_distribution']
        assert stats['tool_usage_distribution']['github'] == 2


class TestIntegrationScenarios:
    """Test integration scenarios and real-world usage"""

    @pytest.fixture
    def integration_engine(self):
        """Create engine for integration testing"""
        with tempfile.TemporaryDirectory() as temp_dir:
            engine = RoutingEnforcementEngine(config_path="dummy.yml")
            engine.rules_db_path = Path(temp_dir) / "rules.db"
            engine.logs_db_path = Path(temp_dir) / "logs.db"
            engine.performance_logs_path = Path(temp_dir) / "performance.jsonl"
            engine._init_databases()
            yield engine

    def test_github_workflow_validation(self, integration_engine):
        """Test complete GitHub workflow validation"""
        session_id = "github_session"
        task_id = "create_repo_task"

        # Valid GitHub operations
        valid_tools = ["github"]
        for tool in valid_tools:
            is_valid, recommended, reason = integration_engine.validate_tool_selection(
                session_id, "create a new repository", tool
            )
            assert is_valid is True

        # Invalid GitHub operations (forbidden tools)
        invalid_tools = ["fetch", "curl"]
        for tool in invalid_tools:
            is_valid, recommended, reason = integration_engine.validate_tool_selection(
                session_id, "create a new repository", tool
            )
            assert is_valid is False
            assert recommended == "github"

    def test_web_research_workflow_validation(self, integration_engine):
        """Test complete web research workflow validation"""
        session_id = "research_session"

        # Web search followed by scraping
        search_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "search for latest AI research", "brave-search"
        )
        assert search_valid is True

        scrape_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "extract data from research papers", "firecrawl"
        )
        assert scrape_valid is True

        # But using forbidden tools should be caught
        fetch_invalid, _, _ = integration_engine.validate_tool_selection(
            session_id, "search for latest AI research", "fetch"
        )
        assert fetch_invalid is False

    def test_development_workflow_validation(self, integration_engine):
        """Test complete development workflow validation"""
        session_id = "dev_session"

        # Code analysis
        analysis_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "find all Python functions in project", "code-index"
        )
        assert analysis_valid is True

        # File operations
        file_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "read configuration file", "filesystem"
        )
        assert file_valid is True

        # GitHub operations
        git_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "commit changes to repository", "github"
        )
        assert git_valid is True

    def test_api_testing_workflow_validation(self, integration_engine):
        """Test complete API testing workflow validation"""
        session_id = "api_session"

        # API testing tools
        schema_valid, _, _ = integration_engine.validate_tool_selection(
            session_id, "validate API schema compliance", "schemathesis"
        )
        assert schema_valid is True

        # But using wrong tools should be flagged
        fetch_invalid, _, _ = integration_engine.validate_tool_selection(
            session_id, "test API endpoints", "fetch"
        )
        assert fetch_invalid is False

    def test_performance_monitoring_integration(self, integration_engine):
        """Test performance monitoring across multiple operations"""
        session_id = "perf_session"
        task_id = "complex_task"

        # Simulate multiple tool calls
        tool_calls = ["code-index", "filesystem", "memory-bank", "github"]

        for tool in tool_calls:
            # Track performance
            with integration_engine.track_performance(session_id, task_id, tool):
                time.sleep(0.001)  # Minimal delay

            # Validate tool selection
            is_valid, _, _ = integration_engine.validate_tool_selection(
                session_id, f"use {tool} for task", tool
            )
            assert is_valid is True

        # Check final statistics
        assert integration_engine.stats['total_validations'] == len(tool_calls)

        # Check performance data was logged
        with open(integration_engine.performance_logs_path, 'r') as f:
            lines = f.readlines()
            assert len(lines) == len(tool_calls)

    def test_agent_configuration_scenarios(self, integration_engine):
        """Test various agent configuration scenarios"""

        # Scenario 1: Cline-like agent configuration
        cline_config = {
            'systemPrompt': 'Use DOMAIN-SPECIFIC routing rules. GitHub operations should use github MCP.',
            'maxParallelCalls': 1,
            'maxCallsPerTask': 8,
            'toolExecution': {'callTimeout': 45000},
            'mcpServers': {
                'github': {'command': 'python', 'args': ['servers/github-mcp.py']},
                'firecrawl': {'command': 'python', 'args': ['servers/firecrawl-mcp.py']}
            }
        }

        result = integration_engine.verify_agent_configuration(cline_config)
        assert result['is_valid'] is True

        # Scenario 2: Poorly configured agent
        poor_config = {
            'maxParallelCalls': 3,  # Wrong
            'maxCallsPerTask': 20,  # Wrong
            # No routing rules in system prompt
            # No timeout configuration
        }

        result = integration_engine.verify_agent_configuration(poor_config)
        assert result['is_valid'] is False
        assert len(result['required_actions']) > 0


def test_module_integration():
    """Test that the module can be imported and basic functions work"""
    from routing_enforcement_system import main

    # Test that main function exists and doesn't crash
    try:
        main()
    except SystemExit:
        pass  # Expected if main() calls sys.exit()


if __name__ == "__main__":
    # Run the tests
    pytest.main([__file__, "-v", "--tb=short"])