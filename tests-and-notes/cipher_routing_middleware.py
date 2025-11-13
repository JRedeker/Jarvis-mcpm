#!/usr/bin/env python3
"""
Cipher Routing Middleware for Agent Logging
Integrates with AgentConfigurationSystem to log AI agent attributes during MCP tool calls.
"""

import json
import logging
from typing import Dict, Any, Optional
from datetime import datetime

# Import the AgentConfigurationSystem
from agent_configuration_system import AgentConfigurationSystem

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class CipherRoutingMiddleware:
    """Middleware for logging agent attributes in MCP tool handlers."""

    def __init__(self):
        self.config_system = AgentConfigurationSystem()

    def log_agent_call(self, session_id: str, task_description: str, selected_tool: str,
                       execution_time_ms: int = 0, success: bool = True, context: Optional[Dict] = None):
        """Log the calling agent's attributes and monitor tool selection."""
        try:
            # Get agent metrics
            agent_metrics = self.config_system.get_agent_metrics(session_id)

            if agent_metrics:
                logger.info(f"Agent calling tool: {agent_metrics['agent_type']} (session: {session_id})")
                logger.info(f"Agent configuration: {json.dumps(agent_metrics['configuration'], indent=2)}")
                logger.info(f"Performance metrics: {json.dumps(agent_metrics['performance_metrics'], indent=2)}")

            # Monitor the tool selection
            monitoring_report = self.config_system.monitor_tool_selection(
                session_id=session_id,
                task_description=task_description,
                selected_tool=selected_tool,
                execution_time_ms=execution_time_ms,
                success=success,
                context=context
            )
            logger.info(f"Tool selection monitoring report: {json.dumps(monitoring_report, indent=2)}")

            # Check performance constraints
            constraints = self.config_system.check_performance_constraints(session_id)
            if constraints.get('violations'):
                logger.warning(f"Performance constraints violated: {json.dumps(constraints, indent=2)}")

            return {
                "logged": True,
                "agent_type": agent_metrics['agent_type'] if agent_metrics else "unknown",
                "compliance": monitoring_report["validation_result"]["is_compliant"],
                "performance_score": monitoring_report["validation_result"]["performance_score"]
            }

        except Exception as e:
            logger.error(f"Error logging agent call: {e}")
            return {"logged": False, "error": str(e)}

    def handle_tool_call_wrapper(self, tool_handler):
        """Decorator to wrap MCP tool handlers with logging."""
        def wrapper(params: Dict[str, Any]):
            session_id = params.get("sessionId", "unknown")
            tool_name = params.get("name", "unknown")
            arguments = params.get("arguments", {})

            # Infer task description from arguments or use default
            task_desc = arguments.get("task", "Tool execution") if isinstance(arguments, dict) else "Tool execution"

            start_time = datetime.now()

            try:
                # Call the actual tool handler
                result = tool_handler(params)

                execution_time_ms = int((datetime.now() - start_time).total_seconds() * 1000)

                # Log the agent call
                log_result = self.log_agent_call(
                    session_id=session_id,
                    task_description=task_desc,
                    selected_tool=tool_name,
                    execution_time_ms=execution_time_ms,
                    success=True,
                    context=arguments
                )

                logger.info(f"Tool {tool_name} executed successfully. Log result: {log_result}")

                return result

            except Exception as e:
                execution_time_ms = int((datetime.now() - start_time).total_seconds() * 1000)

                # Log the failed call
                log_result = self.log_agent_call(
                    session_id=session_id,
                    task_description=task_desc,
                    selected_tool=tool_name,
                    execution_time_ms=execution_time_ms,
                    success=False,
                    context={"error": str(e)}
                )

                logger.error(f"Tool {tool_name} failed. Log result: {log_result}")
                raise e

        return wrapper

# Example usage in an MCP server
def example_mcp_tool_handler(params: Dict[str, Any]) -> Dict[str, Any]:
    """Example tool handler that can be wrapped with middleware."""
    # Simulate tool execution
    return {"result": "Tool executed", "params": params}

if __name__ == "__main__":
    middleware = CipherRoutingMiddleware()

    # Wrap the example handler
    wrapped_handler = middleware.handle_tool_call_wrapper(example_mcp_tool_handler)

    # Simulate a tool call
    test_params = {
        "sessionId": "test-session-123",
        "name": "example_tool",
        "arguments": {"task": "Test agent logging"}
    }

    result = wrapped_handler(test_params)
    print(f"Tool result: {result}")