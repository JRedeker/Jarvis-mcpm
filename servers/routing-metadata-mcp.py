#!/usr/bin/env python3
"""
Routing Metadata MCP Server (SDK-based)

Refactored to use the official MCP Python SDK for robust stdio lifecycle
and async handling. Preserves the original tool APIs and output shapes.

Tools:
  - validate_tool_selection
  - track_tool_execution
  - initialize_session
  - get_routing_analytics
"""

import asyncio
import json
import logging
import os
import sys
from typing import Any, Optional, Dict, List
from datetime import datetime, timezone

# Add project root to sys.path for local imports
PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if PROJECT_ROOT not in sys.path:
    sys.path.insert(0, PROJECT_ROOT)

# Local modules
from cipher_routing_middleware import get_middleware
from servers.otel_logs import setup_otel_logging, emit_structured_log
from servers.log_setup import init_logging

# MCP SDK
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent

# Configure logging (preserve file + stderr like the original)
LOG_PATH = "/home/jrede/dev/MCP/logs/routing-metadata.log"
DECISIONS_JSONL_PATH = "/home/jrede/dev/MCP/logs/routing-decisions.jsonl"
os.makedirs(os.path.dirname(LOG_PATH), exist_ok=True)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[
        logging.FileHandler(LOG_PATH),
        logging.StreamHandler(sys.stderr),
    ],
)
logger = logging.getLogger("routing-metadata-mcp")
# OTel provider (initialized in main)
otel_provider: Optional[Any] = None

# Middleware instance
middleware = get_middleware()

# MCP App
app = Server("routing-metadata-mcp")


def _tool_defs() -> List[Tool]:
    """Return MCP Tool definitions matching original schemas."""
    return [
        Tool(
            name="validate_tool_selection",
            description="Validate a tool selection against routing rules and return metadata",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier",
                    },
                    "agent_type": {
                        "type": "string",
                        "description": "Type of agent (cline, kilocode, general)",
                        "default": "general",
                    },
                    "task_description": {
                        "type": "string",
                        "description": "Description of the task being performed",
                    },
                    "selected_tool": {
                        "type": "string",
                        "description": "Tool selected by the agent",
                    },
                    "context": {
                        "type": "object",
                        "description": "Additional context for routing decision",
                        "additionalProperties": True,
                    },
                },
                "required": ["session_id", "task_description", "selected_tool"],
            },
        ),
        Tool(
            name="track_tool_execution",
            description="Track tool execution metrics for performance monitoring",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier",
                    },
                    "tool_name": {
                        "type": "string",
                        "description": "Name of the executed tool",
                    },
                    "execution_time_ms": {
                        "type": "integer",
                        "description": "Execution time in milliseconds",
                    },
                    "success": {
                        "type": "boolean",
                        "description": "Whether execution was successful",
                    },
                    "error_message": {
                        "type": "string",
                        "description": "Error message if execution failed",
                    },
                },
                "required": ["session_id", "tool_name", "execution_time_ms", "success"],
            },
        ),
        Tool(
            name="initialize_session",
            description="Initialize session tracking with performance constraints",
            inputSchema={
                "type": "object",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier to initialize",
                    },
                    "mode": {
                        "type": "string",
                        "description": "Execution mode (serial/parallel)",
                        "default": "serial",
                    },
                    "max_calls": {
                        "type": "integer",
                        "description": "Maximum calls allowed for this session",
                        "default": 8,
                    },
                },
                "required": ["session_id"],
            },
        ),
        Tool(
            name="get_routing_analytics",
            description="Get routing analytics for the specified time period",
            inputSchema={
                "type": "object",
                "properties": {
                    "days_back": {
                        "type": "integer",
                        "description": "Number of days to analyze",
                        "default": 30,
                    }
                },
            },
        ),
    ]


app.list_tools()
async def list_tools() -> List[Tool]:
    """Expose tools via MCP SDK (mirrors original list)."""
    return _tool_defs()


def _json_text(payload: Dict[str, Any], indent: Optional[int] = 2) -> List[TextContent]:
    """Wrap payload as MCP TextContent with JSON text."""
    return [TextContent(type="text", text=json.dumps(payload, indent=indent, default=str))]


def _append_jsonl(path: str, record: Dict[str, Any]) -> None:
    """Append a single JSON object as JSONL to path; errors are logged but never raise."""
    try:
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "a", encoding="utf-8") as f:
            f.write(json.dumps(record, ensure_ascii=False) + "\n")
    except Exception as e:
        logger.error("Failed to write JSONL '%s': %s", path, e)


@app.call_tool()
async def call_tool(name: str, arguments: Dict[str, Any]) -> List[TextContent]:
    """Dispatch tool calls; preserve original response shapes."""
    try:
        if name == "validate_tool_selection":
            session_id = arguments["session_id"]
            agent_type = arguments.get("agent_type", "general")
            task_description = arguments["task_description"]
            selected_tool = arguments["selected_tool"]
            context = arguments.get("context")

            # Validate tool selection via middleware
            should_allow, suggested_tool, validation = middleware.validate_tool_call(
                session_id=session_id,
                agent_type=agent_type,
                task_description=task_description,
                selected_tool=selected_tool,
                context=context,
            )

            if validation is None:
                raise ValueError("Validation result is None")

            logger.info(
                "Tool validation - Session: %s, Tool: %s, Status: %s, Domain: %s",
                session_id,
                selected_tool,
                validation.get("status"),
                validation.get("detected_domain"),
            )

            # OTel log
            if otel_provider:
                emit_structured_log(
                    otel_provider,
                    "routing-metadata-mcp",
                    "validate_tool_selection",
                    level=logging.INFO,
                    session_id=session_id,
                    selected_tool=validation.get("selected_tool"),
                    recommended_tool=validation.get("recommended_tool"),
                    detected_domain=validation.get("detected_domain"),
                    is_compliant=validation.get("is_compliant"),
                    status=validation.get("status"),
                )
            # Emit JSONL fallback (local, cheap) for ROI analysis
            try:
                record = {
                    "ts": datetime.now(timezone.utc).isoformat(),
                    "service_name": "routing-metadata-mcp",
                    "event": "routing_decision",
                    "session_id": session_id,
                    "agent_type": agent_type,
                    "task_description": task_description,
                    "selected_tool": validation.get("selected_tool"),
                    "recommended_tool": validation.get("recommended_tool") or suggested_tool,
                    "detected_domain": validation.get("detected_domain"),
                    "is_compliant": validation.get("is_compliant"),
                    "status": validation.get("status"),
                    "should_allow": should_allow,
                    "fields": {
                        "context_keys": list((context or {}).keys()) if isinstance(context, dict) else None
                    },
                }
                _append_jsonl(DECISIONS_JSONL_PATH, record)
            except Exception as e:
                logger.error("Failed to emit routing_decision JSONL: %s", e)

            result: Dict[str, Any] = {
                "should_allow": should_allow,
                "suggested_tool": suggested_tool,
                "validation": validation,
                "routing_metadata": {
                    "status": validation.get("status"),
                    "detected_domain": validation.get("detected_domain"),
                    "selected_tool": validation.get("selected_tool"),
                    "recommended_tool": validation.get("recommended_tool"),
                    "is_compliant": validation.get("is_compliant"),
                    "timestamp": validation.get("timestamp"),
                },
            }

            if suggested_tool and "routing_metadata" in result:
                result["routing_metadata"]["suggestion"] = {
                    "recommended_tool": suggested_tool,
                    "reason": f"Better tool for {validation.get('detected_domain')} domain",
                }

            return _json_text(result)

        elif name == "track_tool_execution":
            session_id = arguments["session_id"]
            tool_name = arguments["tool_name"]
            execution_time_ms = arguments["execution_time_ms"]
            success = arguments["success"]
            error_message = arguments.get("error_message")

            middleware.track_tool_execution(
                session_id=session_id,
                tool_name=tool_name,
                execution_time_ms=execution_time_ms,
                success=success,
                error_message=error_message,
            )

            if otel_provider:
                emit_structured_log(
                    otel_provider,
                    "routing-metadata-mcp",
                    "track_tool_execution",
                    level=logging.INFO if success else logging.ERROR,
                    session_id=session_id,
                    tool_name=tool_name,
                    execution_time_ms=execution_time_ms,
                    success=success,
                    error_message=error_message or "",
                    is_error=not success,
                )

            logger.info(
                "Tracked execution - Session: %s, Tool: %s, Success: %s, Time: %sms",
                session_id,
                tool_name,
                success,
                execution_time_ms,
            )

            return _json_text(
                {
                    "tracked": True,
                    "session_id": session_id,
                    "tool_name": tool_name,
                },
                indent=None,
            )

        elif name == "initialize_session":
            session_id = arguments["session_id"]
            mode = arguments.get("mode", "serial")
            max_calls = arguments.get("max_calls", 8)

            middleware.initialize_session(session_id=session_id, mode=mode, max_calls=max_calls)
            logger.info("Initialized session - ID: %s, Mode: %s, Max calls: %s", session_id, mode, max_calls)

            return _json_text(
                {
                    "initialized": True,
                    "session_id": session_id,
                    "mode": mode,
                    "max_calls": max_calls,
                },
                indent=None,
            )

        elif name == "get_routing_analytics":
            days_back = arguments.get("days_back", 30)
            analytics = middleware.get_routing_analytics(days_back)
            logger.info("Retrieved routing analytics for %s days", days_back)
            return _json_text(analytics)

        else:
            # Unknown tool: mirror old behavior by returning error payload as text
            error_payload = {"error": f"Unknown tool: {name}", "tool": name}
            logger.error("Unknown tool called: %s", name)
            return _json_text(error_payload, indent=None)

    except Exception as e:
        logger.error("Error in tool %s: %s", name, str(e), exc_info=True)
        # Preserve previous style: return a JSON error payload inside text content
        return _json_text({"error": str(e), "tool": name}, indent=None)


async def main() -> None:
    global otel_provider

    # Initialize ultra-simple JSONL logging (loguru) with rotation/gzip
    # Configurable via env:
    #   LOG_PATH (default ./logs/cipher.jsonl), LOG_LEVEL, LOG_ROTATE_SIZE, LOG_BACKUPS, LOG_GZIP, LOG_CONSOLE
    # Here we prefer a per-service path; override with ROUTING_LOG_PATH if provided.
    slogger = init_logging(
        service_name="routing-metadata-mcp",
    )
    slogger.info("Server starting: routing-metadata-mcp")

    # Initialize OpenTelemetry logging
    otel_provider = setup_otel_logging("routing-metadata-mcp", enable_console_bridge=True)

    logger.info("Starting Routing Metadata MCP Server (SDK-based)")

    # Use the same stdio_server + app.run pattern as llm-inference (known-good)
    async with stdio_server() as (read_stream, write_stream):
        await app.run(read_stream, write_stream, app.create_initialization_options())


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Server stopped by user (KeyboardInterrupt)")
    except Exception as e:
        logger.error("Fatal error: %s", str(e), exc_info=True)
        sys.exit(1)
