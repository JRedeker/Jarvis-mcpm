#!/usr/bin/env python3
"""
Agent Configuration System for Cipher-Aggregator Phase 5 Implementation

This system provides routing rules to agents and monitors their tool selection
decisions to ensure compliance with cipher.yml specifications.

Key Features:
- Dynamic system prompt injection with routing rules
- Real-time monitoring and validation of tool selections
- Performance constraint enforcement (max 8 calls, serial execution)
- Agent configuration verification
- Routing pattern analysis and optimization recommendations
"""

import sys
import os
sys.path.append('/home/jrede/dev/MCP')

from routing_enforcement_system import CipherRoutingEngine, TaskCategory, Domain, ToolSelectionStatus
from typing import Dict, List, Optional, Any, Tuple
from datetime import datetime
import json
import logging
import uuid
from dataclasses import asdict

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class AgentConfigurationSystem:
    """Main system for configuring agents with routing rules and monitoring selections."""

    def __init__(self, cipher_config_path: str = "/home/jrede/dev/MCP/cipher.yml"):
        """Initialize the agent configuration system."""
        self.routing_engine = CipherRoutingEngine(cipher_config_path)
        self.active_agents: Dict[str, Dict] = {}
        self.monitoring_enabled = True

        logger.info("Agent Configuration System initialized")

    def create_agent_configuration(self, agent_type: str = "general",
                                 session_id: Optional[str] = None,
                                 include_performance_constraints: bool = True) -> Dict[str, Any]:
        """Create agent configuration with routing rules."""

        if session_id is None:
            session_id = f"session-{datetime.now().timestamp()}-{str(uuid.uuid4())[:8]}"

        # Generate system prompt with routing rules
        system_prompt = self.routing_engine.generate_system_prompt(
            agent_type=agent_type,
            include_rules=True
        )

        # Create agent configuration
        configuration = {
            "session_id": session_id,
            "agent_type": agent_type,
            "system_prompt": system_prompt,
            "performance_constraints": {
                "max_calls": 8,
                "execution_mode": "serial",
                "timeout_warning_threshold": 10000  # 10 seconds
            } if include_performance_constraints else {},
            "routing_validation": {
                "enabled": True,
                "strict_mode": False,
                "performance_monitoring": True
            },
            "monitoring_config": {
                "track_decisions": True,
                "log_compliance": True,
                "alert_on_violations": True
            },
            "created_at": datetime.now().isoformat(),
            "cipher_version": "2.0",
            "routing_rules_version": "1.0"
        }

        # Store agent configuration
        self.active_agents[session_id] = {
            "configuration": configuration,
            "metrics": self.routing_engine.track_session_performance(session_id),
            "created_at": datetime.now()
        }

        logger.info(f"Created agent configuration for {agent_type} (session: {session_id})")
        return configuration

    def monitor_tool_selection(self, session_id: str, task_description: str,
                             selected_tool: str, execution_time_ms: int = 0,
                             success: bool = True, context: Optional[Dict] = None) -> Dict[str, Any]:
        """Monitor and validate agent tool selection decisions."""

        if session_id not in self.active_agents:
            logger.warning(f"Unknown session ID: {session_id}")
            return {"error": "Unknown session"}

        agent_config = self.active_agents[session_id]
        agent_type = agent_config["configuration"]["agent_type"]

        # Validate the tool selection using the routing engine
        tool_selection = self.routing_engine.validate_tool_selection(
            session_id=session_id,
            agent_type=agent_type,
            task_description=task_description,
            selected_tool=selected_tool,
            context=context
        )

        # Update session performance metrics
        self.routing_engine.update_session_call(session_id, execution_time_ms, success)

        # Generate monitoring report
        monitoring_report = {
            "session_id": session_id,
            "task_description": task_description,
            "selected_tool": selected_tool,
            "validation_result": {
                "is_compliant": tool_selection.selection_status == ToolSelectionStatus.COMPLIANT,
                "status": tool_selection.selection_status.value,
                "detected_domain": tool_selection.detected_domain.value,
                "recommended_tool": tool_selection.recommended_tool,
                "performance_score": tool_selection.performance_score
            },
            "execution_metrics": {
                "execution_time_ms": execution_time_ms,
                "success": success,
                "performance_score": tool_selection.performance_score
            },
            "timestamp": datetime.now().isoformat()
        }

        # Log compliance issues
        if tool_selection.selection_status != ToolSelectionStatus.COMPLIANT:
            logger.warning(f"Routing violation in session {session_id}: "
                          f"Expected {tool_selection.recommended_tool}, "
                          f"got {selected_tool} for {tool_selection.detected_domain.value}")

        return monitoring_report

    def get_agent_metrics(self, session_id: str) -> Optional[Dict[str, Any]]:
        """Get current metrics for an agent session."""
        if session_id not in self.active_agents:
            return None

        agent_data = self.active_agents[session_id]
        session_metrics = self.routing_engine.session_metrics.get(session_id)

        return {
            "session_id": session_id,
            "agent_type": agent_data["configuration"]["agent_type"],
            "created_at": agent_data["created_at"].isoformat(),
            "performance_metrics": asdict(session_metrics) if session_metrics else {},
            "configuration": agent_data["configuration"]
        }

    def check_performance_constraints(self, session_id: str) -> Dict[str, Any]:
        """Check if agent is within performance constraints."""
        if session_id not in self.active_agents:
            return {"error": "Unknown session"}

        return self.routing_engine.check_performance_constraints(session_id)

    def analyze_routing_patterns(self, days_back: int = 30) -> Dict[str, Any]:
        """Analyze routing patterns for optimization recommendations."""
        return self.routing_engine.analyze_routing_patterns(days_back)

    def generate_agent_report(self, session_id: str) -> Optional[Dict[str, Any]]:
        """Generate comprehensive report for an agent session."""
        if session_id not in self.active_agents:
            return None

        metrics = self.get_agent_metrics(session_id)
        performance = self.check_performance_constraints(session_id)
        patterns = self.analyze_routing_patterns(7)  # Last 7 days

        return {
            "session_report": metrics,
            "performance_analysis": performance,
            "routing_patterns": patterns,
            "generated_at": datetime.now().isoformat()
        }

    def cleanup_session(self, session_id: str):
        """Clean up an agent session."""
        if session_id in self.active_agents:
            del self.active_agents[session_id]
            logger.info(f"Cleaned up session: {session_id}")

    def get_domain_specific_rules(self, domain: Domain) -> Dict[str, Any]:
        """Get routing rules for a specific domain."""
        domain_rules = []
        for rule in self.routing_engine.routing_rules:
            if rule.domain == domain:
                domain_rules.append(asdict(rule))

        return {
            "domain": domain.value,
            "rules": domain_rules,
            "total_rules": len(domain_rules)
        }


# Agent-specific configuration templates
AGENT_CONFIGURATIONS = {
    "cline": {
        "description": "VSCode Cline agent configuration",
        "system_prompt_additions": [
            "You are Cline, an AI coding assistant integrated with VSCode.",
            "Always prioritize domain-specific tools over generic ones.",
            "Follow cipher routing rules for optimal tool selection."
        ],
        "performance_mode": "efficient"
    },
    "kilocode": {
        "description": "KiloCode agent configuration",
        "system_prompt_additions": [
            "You are Kilo Code, a highly skilled software engineer.",
            "Use the best tool for each specific domain.",
            "Follow documented routing patterns and performance constraints."
        ],
        "performance_mode": "comprehensive"
    },
    "general": {
        "description": "General agent configuration",
        "system_prompt_additions": [
            "You are an AI assistant with access to specialized tools.",
            "Select the most appropriate tool for each task.",
            "Follow routing guidelines for optimal results."
        ],
        "performance_mode": "balanced"
    }
}


def create_specialized_agent_config(agent_type: str, domain: Optional[Domain] = None,
                                  session_id: Optional[str] = None) -> Dict[str, Any]:
    """Create a specialized agent configuration for specific domains."""

    config_system = AgentConfigurationSystem()

    # Start with base configuration
    config = config_system.create_agent_configuration(
        agent_type=agent_type,
        session_id=session_id
    )

    # Add specialized domain rules if specified
    if domain:
        domain_rules = config_system.get_domain_specific_rules(domain)
        config["domain_specialization"] = domain_rules
        config["system_prompt"] += f"\n\n## Domain Specialization: {domain.value}\n"
        config["system_prompt"] += f"You are specialized in {domain.value} operations.\n"
        config["system_prompt"] += f"Preferred tools: {', '.join([rule['preferred_tool'] for rule in domain_rules['rules']])}"

    return config


def main():
    """Demonstration of the agent configuration system."""
    print("=== Cipher Agent Configuration System Demo ===\n")

    # Initialize the system
    config_system = AgentConfigurationSystem()

    # Create configurations for different agent types
    agent_types = ["cline", "kilocode", "general"]

    for agent_type in agent_types:
        print(f"Creating configuration for {agent_type} agent...")
        config = config_system.create_agent_configuration(agent_type)
        print(f"✅ Configuration created (session: {config['session_id']})\n")

    # Simulate some tool selections
    print("Simulating tool selection monitoring...")
    session_id = list(config_system.active_agents.keys())[0]

    test_selections = [
        ("Search for Python tutorials online", "brave_web_search", 1500, True),
        ("Read a file from workspace", "filesystem_read", 200, True),
        ("Scrape website content", "firecrawl_scrape", 3000, True),
        ("Make HTTP request", "fetch", 800, False)  # This should trigger a warning
    ]

    for task, tool, time_ms, success in test_selections:
        result = config_system.monitor_tool_selection(
            session_id=session_id,
            task_description=task,
            selected_tool=tool,
            execution_time_ms=time_ms,
            success=success
        )

        status = "✅" if result["validation_result"]["is_compliant"] else "⚠️"
        print(f"{status} {task} → {tool} ({time_ms}ms)")

    # Generate final report
    print(f"\nGenerating final report...")
    report = config_system.generate_agent_report(session_id)
    if report:
        print("✅ Report generated successfully")
        print(f"   Performance Score: {report['performance_analysis'].get('overall_score', 'N/A')}")
        print(f"   Compliance Rate: {report['routing_patterns'].get('compliance_rate', 'N/A')}")


if __name__ == "__main__":
    main()