"""
Logging setup module for MCP servers
Provides consistent logging configuration across all Python MCP servers
"""

import logging
import sys
from typing import Optional


def init_logging(
    name: str = "mcp-server",
    level: int = logging.INFO,
    format_string: Optional[str] = None,
    service_name: Optional[str] = None
) -> logging.Logger:
    """
    Initialize logging for MCP servers with consistent formatting

    Args:
        name: Logger name (default: "mcp-server")
        level: Logging level (default: logging.INFO)
        format_string: Custom format string (optional)
        service_name: Service name for logging (optional, uses name if not provided)

    Returns:
        Configured logger instance
    """
    # Use service_name if provided, otherwise use name
    logger_name = service_name if service_name else name
    logger = logging.getLogger(logger_name)
    logger.setLevel(level)

    # Avoid duplicate handlers
    if logger.handlers:
        return logger

    # Create console handler
    handler = logging.StreamHandler(sys.stderr)
    handler.setLevel(level)

    # Set format
    if format_string is None:
        format_string = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'

    formatter = logging.Formatter(format_string)
    handler.setFormatter(formatter)

    logger.addHandler(handler)

    return logger


def get_logger(name: str = "mcp-server") -> logging.Logger:
    """
    Get or create a logger instance

    Args:
        name: Logger name

    Returns:
        Logger instance
    """
    return logging.getLogger(name)
