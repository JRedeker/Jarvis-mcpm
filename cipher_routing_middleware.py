#!/usr/bin/env python3
"""
Cipher Routing Middleware - Phase 5 Global Deployment

Integrates routing enforcement into cipher-aggregator's tool call flow.
Uses evalLlm for intelligent routing decisions and validates all tool selections.

Architecture:
    Tool Call Request
         ↓
    Routing Middleware (this module)
         ↓
    Validate against routing rules
         ↓
    Log to routing_decisions.db
         ↓
    Execute tool OR suggest better tool
"""

import sys
import os
from typing import Dict, Any, Optional, Tuple
import logging
from datetime import datetime
import json

# Add project root to path for imports
project_root = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, project_root)

# Import from tests-and-notes directory (note: hyphenated directory requires different import approach)
import importlib.util
spec = importlib.util.spec_from_file_location(
    "routing_enforcement_system",
    os.path.join(project_root, "tests-and-notes", "routing_enforcement_system.py")
)
routing_module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(routing_module)

CipherRoutingEngine = routing_module.CipherRoutingEngine
ToolSelectionStatus = routing_module.ToolSelectionStatus

logger = logging.getLogger(__name__)


class CipherRoutingMiddleware:
    """
    Middleware for cipher-aggregator that enforces routing rules on all tool calls.

    Integrates with cipher's existing evalLlm architecture for intelligent routing.
    """

    def __init__(self, config_path: str = "/home/jrede/dev/MCP/cipher.yml"):
        """Initialize routing middleware with cipher configuration."""
        self.routing_engine = CipherRoutingEngine(config_path)
        self.enforcement_mode = "warn"  # Options: "warn", "block", "suggest"

        logger.info("Cipher Routing Middleware initialized")
        logger.info(f"Enforcement mode: {self.enforcement_mode}")

    def validate_tool_call(
        self,
        session_id: str,
        agent_type: str,
        task_description: str,
        selected_tool: str,
        context: Optional[Dict] = None
    ) -> Tuple[bool, Optional[str], Optional[Dict]]:
        """
        Validate a tool call against routing rules.

        Args:
            session_id: Current session identifier
            agent_type: Type of agent making the call (cline, kilocode, etc.)
            task_description: Description of the task being performed
            selected_tool: Tool selected by the agent
            context: Additional context for routing decision

        Returns:
            Tuple of (should_allow, suggested_tool, validation_result)
            - should_allow: Whether to allow the tool call
            - suggested_tool: Better tool suggestion (if non-compliant)
            - validation_result: Full validation details
        """

        # Validate selection
        selection = self.routing_engine.validate_tool_selection(
            session_id=session_id,
            agent_type=agent_type,
            task_description=task_description,
            selected_tool=selected_tool,
            context=context
        )

        # Build validation result
        validation_result = {
            "status": selection.selection_status.value,
            "detected_domain": selection.detected_domain.value,
            "recommended_tool": selection.recommended_tool,
            "selected_tool": selection.selected_tool,
            "is_compliant": selection.selection_status == ToolSelectionStatus.COMPLIANT,
            "timestamp": selection.selection_timestamp.isoformat()
        }

        # Determine action based on enforcement mode and status
        if selection.selection_status == ToolSelectionStatus.COMPLIANT:
            return True, None, validation_result

        # Non-compliant - choose action based on mode
        if self.enforcement_mode == "block":
            # Block execution and force recommended tool
            logger.warning(
                f"BLOCKED non-compliant tool selection: {selected_tool} "
                f"(recommended: {selection.recommended_tool})"
            )
            return False, selection.recommended_tool, validation_result

        elif self.enforcement_mode == "suggest":
            # Allow but suggest better tool
            logger.info(
                f"SUGGESTING better tool: {selection.recommended_tool} "
                f"(selected: {selected_tool})"
            )
            return True, selection.recommended_tool, validation_result

        else:  # warn mode (default)
            # Allow but log warning
            logger.warning(
                f"WARNING: Non-compliant tool selection: {selected_tool} "
                f"(recommended: {selection.recommended_tool}) - domain: {selection.detected_domain.value}"
            )
            return True, None, validation_result

    def track_tool_execution(
        self,
        session_id: str,
        tool_name: str,
        execution_time_ms: int,
        success: bool,
        error_message: Optional[str] = None
    ):
        """
        Track tool execution metrics for performance monitoring.

        Args:
            session_id: Session identifier
            tool_name: Name of the executed tool
            execution_time_ms: Execution time in milliseconds
            success: Whether execution was successful
            error_message: Error message if execution failed
        """
        try:
            self.routing_engine.update_session_call(
                session_id=session_id,
                execution_time_ms=execution_time_ms,
                success=success
            )

            # Check performance constraints
            constraints = self.routing_engine.check_performance_constraints(session_id)

            if constraints["violation"]:
                logger.error(
                    f"Performance constraint violation in session {session_id}: "
                    f"{constraints['violations']}"
                )

        except Exception as e:
            logger.error(f"Error tracking tool execution: {e}")

    def initialize_session(
        self,
        session_id: str,
        mode: str = "serial",
        max_calls: int = 8
    ):
        """
        Initialize session tracking with performance constraints.

        Args:
            session_id: Session identifier to initialize
            mode: Execution mode (serial/parallel)
            max_calls: Maximum calls allowed for this session
        """
        try:
            self.routing_engine.track_session_performance(
                session_id=session_id,
                mode=mode,
                max_calls=max_calls
            )
            logger.info(
                f"Initialized session {session_id}: mode={mode}, max_calls={max_calls}"
            )
        except Exception as e:
            logger.error(f"Error initializing session: {e}")

    def get_routing_analytics(self, days_back: int = 30) -> Dict[str, Any]:
        """
        Get routing analytics for the specified time period.

        Args:
            days_back: Number of days to analyze

        Returns:
            Dictionary with routing analytics and recommendations
        """
        return self.routing_engine.analyze_routing_patterns(days_back)

    def set_enforcement_mode(self, mode: str):
        """
        Set enforcement mode for routing validation.

        Args:
            mode: One of "warn", "block", "suggest"
        """
        if mode not in ["warn", "block", "suggest"]:
            raise ValueError(f"Invalid enforcement mode: {mode}")

        self.enforcement_mode = mode
        logger.info(f"Enforcement mode changed to: {mode}")

    def generate_system_prompt_injection(self, agent_type: str = "general") -> str:
        """
        Generate system prompt injection with routing rules for agents.

        This should be added to agent configurations to provide routing guidance.

        Args:
            agent_type: Type of agent (cline, kilocode, general)

        Returns:
            System prompt text with routing rules
        """
        return self.routing_engine.generate_system_prompt(
            agent_type=agent_type,
            include_rules=True
        )


# Global middleware instance
_middleware_instance: Optional[CipherRoutingMiddleware] = None


def get_middleware() -> CipherRoutingMiddleware:
    """Get or create the global middleware instance."""
    global _middleware_instance
    if _middleware_instance is None:
        _middleware_instance = CipherRoutingMiddleware()
    return _middleware_instance


def validate_tool_call(
    session_id: str,
    agent_type: str,
    task_description: str,
    selected_tool: str,
    context: Optional[Dict] = None
) -> Tuple[bool, Optional[str], Optional[Dict]]:
    """
    Convenience function for validating tool calls.

    See CipherRoutingMiddleware.validate_tool_call for details.
    """
    middleware = get_middleware()
    return middleware.validate_tool_call(
        session_id, agent_type, task_description, selected_tool, context
    )


def track_tool_execution(
    session_id: str,
    tool_name: str,
    execution_time_ms: int,
    success: bool,
    error_message: Optional[str] = None
):
    """
    Convenience function for tracking tool execution.

    See CipherRoutingMiddleware.track_tool_execution for details.
    """
    middleware = get_middleware()
    middleware.track_tool_execution(
        session_id, tool_name, execution_time_ms, success, error_message
    )


def initialize_session(session_id: str, mode: str = "serial", max_calls: int = 8):
    """
    Convenience function for initializing sessions.

    See CipherRoutingMiddleware.initialize_session for details.
    """
    middleware = get_middleware()
    middleware.initialize_session(session_id, mode, max_calls)


# Example integration hooks for cipher-aggregator
def example_tool_call_handler(request: Dict[str, Any]) -> Dict[str, Any]:
    """
    Example of how to integrate middleware into cipher's tool call handler.

    This would be integrated into cipher-aggregator's actual tool call processing.
    """
    # Extract request parameters
    session_id = request.get("sessionId", "unknown")
    tool_name = request["params"]["name"]
    task_description = request["params"].get("task", "")
    agent_type = request.get("agent", "general")

    # Initialize session if not exists
    initialize_session(session_id)

    # Validate tool selection
    should_allow, suggested_tool, validation = validate_tool_call(
        session_id=session_id,
        agent_type=agent_type,
        task_description=task_description,
        selected_tool=tool_name
    )

    if not should_allow:
        # Block execution, return suggestion
        return {
            "jsonrpc": "2.0",
            "id": request["id"],
            "error": {
                "code": -32001,
                "message": f"Tool selection blocked. Use {suggested_tool} instead.",
                "data": validation
            }
        }

    # Execute tool (this would call actual cipher tool execution)
    start_time = datetime.now()
    try:
        # ... actual tool execution ...
        result = {"success": True}
        success = True
        error_msg = None
    except Exception as e:
        result = {"error": str(e)}
        success = False
        error_msg = str(e)

    # Track execution
    execution_time = int((datetime.now() - start_time).total_seconds() * 1000)
    track_tool_execution(
        session_id=session_id,
        tool_name=tool_name,
        execution_time_ms=execution_time,
        success=success,
        error_message=error_msg
    )

    # Add routing info to response if suggestion exists
    if suggested_tool:
        result["routing_suggestion"] = {
            "recommended_tool": suggested_tool,
            "reason": f"Better tool for {validation['detected_domain']} domain"
        }

    return {
        "jsonrpc": "2.0",
        "id": request["id"],
        "result": result
    }


if __name__ == "__main__":
    # Test the middleware
    print("=== CIPHER ROUTING MIDDLEWARE TEST ===\n")

    middleware = get_middleware()

    # Test 1: Compliant selection
    print("Test 1: Compliant Tool Selection")
    should_allow, suggestion, validation = middleware.validate_tool_call(
        session_id="test-session-1",
        agent_type="cline",
        task_description="Search GitHub for repositories",
        selected_tool="github"
    )
    print(f"  Allow: {should_allow}")
    print(f"  Suggestion: {suggestion}")
    print(f"  Status: {validation['status']}\n")

    # Test 2: Non-compliant selection
    print("Test 2: Non-Compliant Tool Selection")
    should_allow, suggestion, validation = middleware.validate_tool_call(
        session_id="test-session-2",
        agent_type="cline",
        task_description="Search GitHub for repositories",
        selected_tool="fetch"  # Wrong tool!
    )
    print(f"  Allow: {should_allow}")
    print(f"  Suggestion: {suggestion}")
    print(f"  Status: {validation['status']}")
    print(f"  Recommended: {validation['recommended_tool']}\n")

    # Test 3: Performance tracking
    print("Test 3: Performance Tracking")
    middleware.initialize_session("perf-test-1", mode="serial", max_calls=8)
    middleware.track_tool_execution("perf-test-1", "github", 250, True)
    middleware.track_tool_execution("perf-test-1", "code-index", 1500, True)
    constraints = middleware.routing_engine.check_performance_constraints("perf-test-1")
    print(f"  Calls remaining: {constraints['calls_remaining']}")
    print(f"  Avg execution time: {constraints['average_execution_time']:.0f}ms")
    print(f"  Success rate: {constraints['success_rate']:.1f}%\n")

    print("=== MIDDLEWARE READY FOR INTEGRATION ===")
