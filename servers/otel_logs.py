from typing import Any, Optional
import logging


def setup_otel_logging(service_name: str, enable_console_bridge: bool = False) -> Optional[Any]:
    """Minimal shim for OpenTelemetry logging used in dev environments.

    Returns a dummy provider object (the logger) so callers can still emit structured logs.
    """
    logger = logging.getLogger(service_name + "-otel-shim")
    logger.setLevel(logging.INFO)
    if not logger.handlers:
        handler = logging.StreamHandler()
        formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
        handler.setFormatter(formatter)
        logger.addHandler(handler)
    return logger


def emit_structured_log(provider: Any, service: str, operation: str, level: int = logging.INFO, **kwargs) -> None:
    """Emit a simple structured log via the provided logger (shim).

    This keeps behavior similar to real OTel emission without requiring full OTel setup.
    """
    try:
        if provider:
            provider.log(level, f"{service}.{operation} - {kwargs}")
    except Exception:
        # Don't let logging errors break the host process
        pass
