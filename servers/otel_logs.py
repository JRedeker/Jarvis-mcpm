#!/usr/bin/env python3
"""
OpenTelemetry Logs Helper Module

Provides shared OTel logging initialization for Cipher MCP servers.
Configures OTLP export via gRPC with batching and optional fallback to JSONL.

Environment Variables:
  OTEL_LOGS_EXPORTER: "otlp" or "none" (default: otlp)
  OTEL_EXPORTER_OTLP_PROTOCOL: "grpc" or "http/protobuf" (default: grpc)
  OTEL_EXPORTER_OTLP_ENDPOINT: Collector endpoint (default: http://localhost:4317)
  OTEL_EXPORTER_OTLP_COMPRESSION: "gzip" or "none" (default: none for local)
  OTEL_BSP_SCHEDULE_DELAY: Batch schedule delay in ms (default: 2000)
  OTEL_BSP_MAX_QUEUE_SIZE: Max queue size (default: 2048)
  OTEL_BSP_MAX_EXPORT_BATCH_SIZE: Max batch size (default: 512)
  OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT: Max attribute length (default: 4096)
"""

import logging
import os
import socket
from typing import Optional

from opentelemetry import _logs as otel_logs
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.sdk.resources import Resource
from opentelemetry.semconv.resource import ResourceAttributes

logger = logging.getLogger(__name__)


def get_git_sha() -> str:
    """Get git SHA for service.version, fallback to 'dev'"""
    try:
        import subprocess

        result = subprocess.run(
            ["git", "rev-parse", "--short", "HEAD"],
            capture_output=True,
            text=True,
            cwd=os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
            timeout=2,
        )
        if result.returncode == 0:
            return result.stdout.strip()
    except Exception:
        pass
    return "dev"


def setup_otel_logging(
    service_name: str,
    enable_console_bridge: bool = True,
) -> Optional[LoggerProvider]:
    """
    Initialize OpenTelemetry logging for a service.

    Args:
        service_name: Name of the service (e.g., "llm-inference-mcp")
        enable_console_bridge: If True, bridge stdlib logging to OTel

    Returns:
        LoggerProvider instance if OTel is enabled, None otherwise
    """
    exporter_type = os.getenv("OTEL_LOGS_EXPORTER", "otlp").lower()

    if exporter_type == "none":
        logger.info(
            f"OTel logging disabled for {service_name} (OTEL_LOGS_EXPORTER=none)"
        )
        return None

    # Resource attributes (service metadata)
    resource = Resource.create(
        {
            ResourceAttributes.SERVICE_NAME: service_name,
            ResourceAttributes.SERVICE_VERSION: get_git_sha(),
            ResourceAttributes.SERVICE_INSTANCE_ID: f"{socket.gethostname()}-{os.getpid()}",
        }
    )

    # Create logger provider
    logger_provider = LoggerProvider(resource=resource)

    # Configure OTLP exporter
    try:
        protocol = os.getenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc").lower()
        endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
        compression_str = os.getenv("OTEL_EXPORTER_OTLP_COMPRESSION", "none").lower()

        # Batch processor settings
        schedule_delay_ms = int(os.getenv("OTEL_BSP_SCHEDULE_DELAY", "2000"))
        max_queue_size = int(os.getenv("OTEL_BSP_MAX_QUEUE_SIZE", "2048"))
        max_export_batch_size = int(os.getenv("OTEL_BSP_MAX_EXPORT_BATCH_SIZE", "512"))

        # Create OTLP exporter (currently only gRPC supported in this helper)
        if protocol == "grpc":
            otlp_exporter = OTLPLogExporter(endpoint=endpoint)
        else:
            logger.warning(
                f"Protocol {protocol} not yet implemented, falling back to gRPC"
            )
            otlp_exporter = OTLPLogExporter(endpoint="http://localhost:4317")

        # Create batch processor
        batch_processor = BatchLogRecordProcessor(
            otlp_exporter,
            schedule_delay_millis=schedule_delay_ms,
            max_queue_size=max_queue_size,
            max_export_batch_size=max_export_batch_size,
        )

        logger_provider.add_log_record_processor(batch_processor)

        logger.info(
            f"OTel logging initialized for {service_name}: "
            f"endpoint={endpoint}, protocol={protocol}, compression={compression_str}"
        )

    except Exception as e:
        logger.error(f"Failed to initialize OTLP exporter for {service_name}: {e}")
        logger.info("Continuing without OTel export (local logs will still work)")
        return None

    # Set global logger provider
    otel_logs.set_logger_provider(logger_provider)

    # Bridge stdlib logging to OTel if requested
    if enable_console_bridge:
        handler = LoggingHandler(
            level=logging.NOTSET,
            logger_provider=logger_provider,
        )
        logging.getLogger().addHandler(handler)
        logger.info(f"Stdlib logging bridged to OTel for {service_name}")

    return logger_provider


def emit_structured_log(
    logger_provider: Optional[LoggerProvider],
    service_name: str,
    event: str,
    level: int = logging.INFO,
    **attributes,
) -> None:
    """
    Emit a structured log with attributes directly to OTel.

    Args:
        logger_provider: The OTel logger provider
        service_name: Service name for the logger
        event: Event name (goes in body)
        level: Log level (logging.INFO, etc.)
        **attributes: Additional attributes as key-value pairs
    """
    if logger_provider is None:
        # OTel not enabled, just use stdlib logger
        std_logger = logging.getLogger(service_name)
        std_logger.log(level, f"{event}: {attributes}")
        return

    # Fully guard OTel emission; fallback to stdlib logger if anything fails
    try:
        otel_logger = logger_provider.get_logger(service_name)

        # Map stdlib levels to OTel severity with broad compatibility
        try:
            from opentelemetry._logs.severity import SeverityNumber  # modern location
        except Exception:
            try:
                from opentelemetry.sdk._logs import SeverityNumber  # older SDK location
            except Exception:
                try:
                    from opentelemetry.sdk._logs.severity import (
                        SeverityNumber,  # alternative SDK location
                    )
                except Exception:
                    from enum import IntEnum

                    class SeverityNumber(IntEnum):
                        DEBUG = 5
                        INFO = 9
                        WARN = 13
                        ERROR = 17
                        FATAL = 21

        severity_map = {
            logging.DEBUG: SeverityNumber.DEBUG,
            logging.INFO: SeverityNumber.INFO,
            logging.WARNING: SeverityNumber.WARN,
            logging.ERROR: SeverityNumber.ERROR,
            logging.CRITICAL: SeverityNumber.FATAL,
        }
        severity = severity_map.get(level, SeverityNumber.INFO)

        # Emit log record
        otel_logger.emit(
            otel_logs.LogRecord(
                timestamp=None,  # auto-generated
                severity_number=severity,
                severity_text=logging.getLevelName(level),
                body=event,
                attributes=attributes,
            )
        )
    except Exception as e:
        std_logger = logging.getLogger(service_name)
        std_logger.log(level, f"{event}: {attributes} (otel_emit_failed={e})")
