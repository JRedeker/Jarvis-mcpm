from datetime import datetime, timezone
from typing import Any, Dict, Optional, Tuple


class _Middleware:
    def validate_tool_call(
        self,
        session_id: str,
        agent_type: str = "general",
        task_description: str = "",
        selected_tool: str = "",
        context: Optional[Dict[str, Any]] = None,
    ) -> Tuple[bool, Optional[str], Dict[str, Any]]:
        """Return a permissive validation result used by routing-metadata.
        Keeps the shape expected by callers in servers/
        """
        now = datetime.now(timezone.utc).isoformat()
        validation = {
            "status": "allowed",
            "detected_domain": "general",
            "selected_tool": selected_tool,
            "recommended_tool": selected_tool,
            "is_compliant": True,
            "timestamp": now,
        }
        # should_allow, suggested_tool, validation
        return True, None, validation

    def track_tool_execution(
        self, session_id: str, tool_name: str, execution_time_ms: int, success: bool, error_message: Optional[str] = None
    ) -> None:
        # No-op tracking for the shim
        return None

    def initialize_session(self, session_id: str, mode: str = "serial", max_calls: int = 8) -> None:
        # No-op initialization
        return None

    def get_routing_analytics(self, days_back: int = 30) -> Dict[str, Any]:
        return {"days_back": days_back, "summary": {"total_calls": 0}}


def get_middleware() -> _Middleware:
    """Return a middleware instance compatible with the real implementation.

    This shim allows dev/testing environments to run routing-metadata without
    pulling in the full production middleware package.
    """
    return _Middleware()
