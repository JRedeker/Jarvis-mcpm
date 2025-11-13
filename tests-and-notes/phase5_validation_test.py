#!/usr/bin/env python3
"""
Phase 5 System Prompt Routing Validation Test Suite

This test suite validates the implementation of system prompt routing according
to Phase 5 requirements from todo.md, ensuring:

1. System prompt injection with routing rules from cipher.yml
2. Agent configuration with domain-specific routing validation
3. Tool selection monitoring and compliance checking
4. Performance constraint enforcement (max 8 calls, serial execution)
5. Domain-specific routing rule adherence

Run this test suite to validate Phase 5 implementation:
    python phase5_validation_test.py
"""

import sys
from datetime import datetime
from typing import Dict, List, Any
import tempfile

# Add current directory to path for imports
sys.path.append('tests-and-notes')

from routing_enforcement_system import CipherRoutingEngine, Domain, TaskCategory
from agent_configuration_system import AgentConfigurationSystem


class TestSystemPromptRouting:
    """Test system prompt routing implementation."""

    def setup_method(self):
        """Set up test fixtures."""
        self.config_system = AgentConfigurationSystem()
        self.test_session_id = "test-session-12345"

    def test_cipher_config_loading(self):
        """Test loading of cipher.yml configuration."""
        engine = CipherRoutingEngine()
        assert engine.config is not None
        assert "systemPrompt" in engine.config
        assert len(engine.routing_rules) > 0
        print("✅ cipher.yml configuration loaded successfully")

    def test_domain_detection(self):
        """Test domain detection from task descriptions."""
        engine = CipherRoutingEngine()

        test_cases = [
            ("Search for Python tutorials online", Domain.WEB_SEARCH),
            ("Create a new repository", Domain.GITHUB),
            ("Scrape website content", Domain.WEB_SCRAPING),
            ("Analyze code structure", Domain.CODE_ANALYSIS),
            ("Test API endpoints", Domain.API_TESTING),
            ("Read file from workspace", Domain.FILE_OPERATIONS),
            ("Run pytest tests", Domain.TEST_EXECUTION),
            ("Make HTTP request", Domain.HTTP_REQUESTS),
        ]

        for task_desc, expected_domain in test_cases:
            detected = engine.detect_domain(task_desc)
            assert detected == expected_domain, f"Expected {expected_domain}, got {detected} for: {task_desc}"

        print("✅ Domain detection working correctly")

    def test_tool_recommendation(self):
        """Test tool recommendation based on routing rules."""
        engine = CipherRoutingEngine()

        test_cases = [
            ("Search for Python tutorials", "brave-search"),
            ("Create repository", "github"),
            ("Scrape website", "firecrawl"),
            ("Analyze code", "code-index"),
        ]

        for task_desc, expected_tool in test_cases:
            domain, recommended = engine.get_routing_recommendation(task_desc)
            assert recommended == expected_tool, f"Expected {expected_tool}, got {recommended} for: {task_desc}"

        print("✅ Tool recommendation working correctly")

    def test_agent_configuration_creation(self):
        """Test agent configuration creation with routing rules."""
        config = self.config_system.create_agent_configuration(
            agent_type="kilocode",
            session_id=self.test_session_id
        )

        assert config["session_id"] == self.test_session_id
        assert config["agent_type"] == "kilocode"
        assert "system_prompt" in config
        assert "routing_rules" in config
        assert "performance_constraints" in config
        assert "monitoring_enabled" in config

        # Check system prompt contains routing rules
        system_prompt = config["system_prompt"]
        assert "DOMAIN-SPECIFIC FIRST" in system_prompt
        assert "github" in system_prompt.lower() or "github" in system_prompt

        print("✅ Agent configuration created successfully")

    def test_tool_selection_monitoring(self):
        """Test tool selection monitoring and validation."""
        # Create agent configuration
        config = self.config_system.create_agent_configuration(
            agent_type="test-agent",
            session_id=self.test_session_id
        )

        # Test compliant tool selection
        result = self.config_system.monitor_tool_selection(
            session_id=self.test_session_id,
            task_description="Search for Python tutorials",
            selected_tool="brave-search",
            execution_time_ms=1000,
            success=True
        )

        assert result["validation_result"]["is_compliant"] == True
        assert result["validation_result"]["selected_tool"] == "brave-search"

        # Test non-compliant tool selection
        result = self.config_system.monitor_tool_selection(
            session_id=self.test_session_id,
            task_description="Search for Python tutorials",
            selected_tool="fetch",  # Wrong tool for web search
            execution_time_ms=1000,
            success=True
        )

        assert result["validation_result"]["is_compliant"] == False
        assert result["validation_result"]["selected_tool"] == "fetch"

        print("✅ Tool selection monitoring working correctly")

    def test_performance_constraints(self):
        """Test performance constraint monitoring."""
        # Test call limit tracking
        metrics = self.config_system.check_performance_constraints(self.test_session_id)
        assert "total_calls" in metrics
        assert "calls_remaining" in metrics
        assert "max_calls" in metrics

        # Verify max 8 calls constraint
        assert metrics["max_calls"] == 8

        print("✅ Performance constraints enforced correctly")

    def test_session_performance_tracking(self):
        """Test session performance tracking."""
        # Create session metrics
        session_metrics = self.config_system.routing_engine.track_session_performance(
            self.test_session_id,
            mode="serial",
            max_calls=8
        )

        assert session_metrics.session_id == self.test_session_id
        assert session_metrics.total_calls == 0
        assert session_metrics.calls_remaining == 8
        assert session_metrics.execution_mode == "serial"

        # Update session with tool calls
        self.config_system.routing_engine.update_session_call(
            self.test_session_id,
            execution_time_ms=1500,
            success=True
        )

        # Verify call was tracked
        updated_metrics = self.config_system.routing_engine.track_session_performance(
            self.test_session_id
        )
        assert updated_metrics.total_calls == 1
        assert updated_metrics.calls_remaining == 7

        print("✅ Session performance tracking working correctly")

    def test_domain_specific_routing_rules(self):
        """Test domain-specific routing rule compliance."""
        engine = CipherRoutingEngine()

        # Test GitHub domain
        github_rule = self.config_system.get_domain_specific_rules(Domain.GITHUB)
        assert github_rule["domain"] == "github"
        assert len(github_rule["rules"]) > 0

        # Verify GitHub operations use github MCP
        tool_selection = engine.validate_tool_selection(
            session_id=self.test_session_id,
            agent_type="test",
            task_description="Create a new repository on GitHub",
            selected_tool="github",
        )
        assert tool_selection.selection_status == ToolSelectionStatus.COMPLIANT

        # Verify forbidden tools are detected
        tool_selection = engine.validate_tool_selection(
            session_id=self.test_session_id,
            agent_type="test",
            task_description="Create a new repository on GitHub",
            selected_tool="fetch",  # Should be forbidden for GitHub
        )
        assert tool_selection.selection_status == ToolSelectionStatus.NON_COMPLIANT

        print("✅ Domain-specific routing rules working correctly")

    def test_system_prompt_injection(self):
        """Test system prompt injection with routing rules."""
        # Test KiloCode agent configuration
        config = self.config_system.create_agent_configuration(
            agent_type="kilocode",
            session_id=self.test_session_id
        )

        system_prompt = config["system_prompt"]

        # Verify routing rules are injected
        assert "DOMAIN-SPECIFIC FIRST (CRITICAL)" in system_prompt
        assert "Serial execution only" in system_prompt
        assert "Max 8 tool calls per task" in system_prompt
        assert "Use domain-specific tools" in system_prompt

        # Verify specific routing instructions
        routing_sections = [
            "GitHub operations",
            "Web scraping",
            "Code analysis",
            "API testing",
            "File operations",
            "Web search"
        ]

        for section in routing_sections:
            assert section.lower() in system_prompt.lower(), f"Missing routing section: {section}"

        print("✅ System prompt injection working correctly")

    def test_routing_pattern_analysis(self):
        """Test routing pattern analysis and optimization."""
        # Create some test data by making tool selections
        test_selections = [
            ("Search web", "brave-search", True),
            ("Create repo", "github", True),
            ("Scrape site", "firecrawl", True),
            ("Wrong tool", "fetch", False),
        ]

        for task, tool, success in test_selections:
            self.config_system.monitor_tool_selection(
                session_id=self.test_session_id,
                task_description=task,
                selected_tool=tool,
                execution_time_ms=1000,
                success=success
            )

        # Analyze patterns
        analysis = self.config_system.analyze_routing_patterns(1)  # Last 1 day

        assert "domain_statistics" in analysis
        assert "tool_effectiveness" in analysis
        assert "optimization_recommendations" in analysis

        print("✅ Routing pattern analysis working correctly")

    def test_agent_type_specific_configurations(self):
        """Test agent type specific configurations."""
        agent_types = ["kilocode", "cline", "general"]

        for agent_type in agent_types:
            config = self.config_system.create_agent_configuration(
                agent_type=agent_type,
                session_id=f"test-{agent_type}"
            )

            assert config["agent_type"] == agent_type
            assert "system_prompt" in config
            assert "routing_rules" in config

            # Each agent type should have specific system prompt additions
            system_prompt = config["system_prompt"]
            assert len(system_prompt) > 1000  # Should be substantial

        print("✅ Agent type specific configurations working correctly")

    def test_comprehensive_end_to_end_scenario(self):
        """Test comprehensive end-to-end scenario."""
        session_id = "e2e-test-session"

        # 1. Create agent configuration
        config = self.config_system.create_agent_configuration(
            agent_type="kilocode",
            session_id=session_id
        )
        assert config is not None

        # 2. Test multiple tool selections
        test_scenarios = [
            ("Search for machine learning tutorials", "brave-search", True, "web_search"),
            ("Create a new GitHub repository", "github", True, "github"),
            ("Extract data from website", "firecrawl", True, "web_scraping"),
            ("Analyze code structure", "code-index", True, "code_analysis"),
            ("Read configuration file", "filesystem_read", True, "file_operations"),
        ]

        for task, tool, success, domain in test_scenarios:
            result = self.config_system.monitor_tool_selection(
                session_id=session_id,
                task_description=task,
                selected_tool=tool,
                execution_time_ms=1500,
                success=success
            )

            assert result["validation_result"]["is_compliant"] == success
            assert result["validation_result"]["detected_domain"] == domain

        # 3. Check performance constraints
        performance = self.config_system.check_performance_constraints(session_id)
        assert performance["total_calls"] == 5
        assert performance["calls_remaining"] == 3  # 8 - 5 = 3

        # 4. Generate comprehensive report
        report = self.config_system.generate_agent_report(session_id)
        assert report is not None
        assert "session_report" in report
        assert "performance_analysis" in report
        assert "routing_patterns" in report

        # 5. Cleanup
        self.config_system.cleanup_session(session_id)

        print("✅ End-to-end scenario completed successfully")


class TestPerformanceConstraints:
    """Test performance constraint enforcement."""

    def test_serial_execution_requirement(self):
        """Test that serial execution is enforced."""
        config_system = AgentConfigurationSystem()

        # Create configuration with performance constraints
        config = config_system.create_agent_configuration(
            agent_type="test",
            session_id="perf-test",
            include_performance_constraints=True
        )

        assert "performance_constraints" in config
        assert config["performance_constraints"]["max_parallel_calls"] == 1
        assert config["performance_constraints"]["max_calls_per_task"] == 8

        print("✅ Serial execution requirements enforced")

    def test_call_limit_tracking(self):
        """Test call limit tracking and warnings."""
        config_system = AgentConfigurationSystem()
        session_id = "call-limit-test"

        # Create session
        config_system.create_agent_configuration(
            agent_type="test",
            session_id=session_id
        )

        # Make calls approaching limit
        for i in range(7):
            config_system.monitor_tool_selection(
                session_id=session_id,
                task_description=f"Test task {i}",
                selected_tool="filesystem_read",
                execution_time_ms=1000,
                success=True
            )

        # Check performance at call 7
        performance = config_system.check_performance_constraints(session_id)
        assert performance["calls_remaining"] == 1
        assert performance["warning_level"] == "caution"

        # Make final call
        config_system.monitor_tool_selection(
            session_id=session_id,
            task_description="Final test task",
            selected_tool="filesystem_read",
            execution_time_ms=1000,
            success=True
        )

        # Check performance at limit
        performance = config_system.check_performance_constraints(session_id)
        assert performance["calls_remaining"] == 0
        assert performance["warning_level"] == "critical"

        print("✅ Call limit tracking and warnings working")


def run_phase5_validation():
    """Run comprehensive Phase 5 validation test suite."""
    print("🚀 Starting Phase 5 System Prompt Routing Validation\n")
    print("=" * 60)

    # Create test suite
    test_suite = TestSystemPromptRouting()
    performance_tests = TestPerformanceConstraints()

    try:
        # Test basic functionality
        print("\n📋 Testing Basic System Functionality")
        print("-" * 40)
        test_suite.setup_method()
        test_suite.test_cipher_config_loading()
        test_suite.test_domain_detection()
        test_suite.test_tool_recommendation()

        # Test agent configuration
        print("\n👤 Testing Agent Configuration System")
        print("-" * 40)
        test_suite.test_agent_configuration_creation()
        test_suite.test_agent_type_specific_configurations()
        test_suite.test_system_prompt_injection()

        # Test monitoring and validation
        print("\n🔍 Testing Monitoring and Validation")
        print("-" * 40)
        test_suite.test_tool_selection_monitoring()
        test_suite.test_domain_specific_routing_rules()

        # Test performance constraints
        print("\n⚡ Testing Performance Constraints")
        print("-" * 40)
        test_suite.test_performance_constraints()
        test_suite.test_session_performance_tracking()

        # Test advanced features
        print("\n🔧 Testing Advanced Features")
        print("-" * 40)
        test_suite.test_routing_pattern_analysis()
        performance_tests.test_serial_execution_requirement()
        performance_tests.test_call_limit_tracking()

        # End-to-end test
        print("\n🎯 Testing End-to-End Scenario")
        print("-" * 40)
        test_suite.test_comprehensive_end_to_end_scenario()

        print("\n" + "=" * 60)
        print("🎉 Phase 5 Validation Completed Successfully!")
        print("✅ All system prompt routing requirements validated")
        print("✅ Domain-specific routing enforcement working")
        print("✅ Performance constraints properly enforced")
        print("✅ Agent configuration system fully functional")

        return True

    except Exception as e:
        print(f"\n❌ Phase 5 Validation Failed: {e}")
        import traceback
        traceback.print_exc()
        return False


if __name__ == "__main__":
    success = run_phase5_validation()
    sys.exit(0 if success else 1)