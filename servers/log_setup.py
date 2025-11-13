from __future__ import annotations

import os
import sys
from pathlib import Path
from typing import Optional, TYPE_CHECKING

if TYPE_CHECKING:
    # Avoid mypy/pylance missing-stub errors at runtime by importing only for typing
    from loguru import Logger  # type: ignore[unused-ignore]
else:
    Logger = object  # runtime placeholder


def _env_bool(name: str, default: bool) -> bool:
    val = os.getenv(name)
    if val is None:
        return default
    return val.strip().lower() in ("1", "true", "yes", "on")


def _env_int(name: str, default: int) -> int:
    try:
        return int(os.getenv(name, str(default)).replace("_", ""))
    except Exception:
        return default


def init_logging(service_name: str = "cipher", log_path: Optional[str] = None):
    """
    Initialize process-wide JSONL logging.
    Imports loguru at runtime to avoid static-checker stub errors.

    Environment variables:
      LOG_PATH          Default: ./logs/cipher.jsonl
      LOG_LEVEL         Default: INFO (e.g., DEBUG, INFO, WARNING, ERROR)
      LOG_ROTATE_SIZE   Default: 10_485_760 (10MB)
      LOG_BACKUPS       Default: 5
      LOG_GZIP          Default: true (gzip rotated segments)
      LOG_CONSOLE       Default: true (human-readable console sink)
    """
    try:
        from loguru import logger as _logger  # runtime import
    except Exception as e:
        raise RuntimeError(
            "loguru is required for servers.log_setup.init_logging(); install with 'pip install loguru'"
        ) from e
    # Defaults
    log_path = log_path or os.getenv("LOG_PATH", "./logs/cipher.jsonl")
    log_level = os.getenv("LOG_LEVEL", "INFO")
    rotate_bytes = _env_int("LOG_ROTATE_SIZE", 10_485_760)  # 10 MB
    backups = _env_int("LOG_BACKUPS", 5)
    gzip_enabled = _env_bool("LOG_GZIP", True)
    console_enabled = _env_bool("LOG_CONSOLE", True)

    # Ensure directory exists
    path_str = str(log_path)
    path = Path(path_str)
    path.parent.mkdir(parents=True, exist_ok=True)

    # Reset default handlers
    _logger.remove()

    # Console (human-friendly)
    if console_enabled:
        _logger.add(
            sys.stderr,
            level=log_level,
            backtrace=False,
            diagnose=False,
            enqueue=True,  # non-blocking console as well
            format="<green>{time:YYYY-MM-DDTHH:mm:ss.SSSZ}</green> | "
                   "<level>{level: <8}</level> | "
                   f"{service_name} | "
                   "<cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - "
                   "<level>{message}</level>",
        )

    # File JSONL (serialized)
    rotation_str = f"{rotate_bytes} bytes"
    retention_str = backups  # Loguru accepts int for count-based retention

    _logger.add(
        str(path),
        level=log_level,
        serialize=True,          # JSONL output
        backtrace=False,
        diagnose=False,
        rotation=rotation_str,   # size-based rotation
        retention=retention_str, # keep N files
        compression="gz" if gzip_enabled else None,
        enqueue=True,            # non-blocking (background thread)
    )

    # Attach basic context fields via bind (optional)
    _logger.bind(service_name=service_name)

    # Initial message
    _logger.opt(lazy=True).info({
        "event": "logging_initialized",
        "service_name": service_name,
        "fields": {
            "path": str(path),
            "level": log_level,
            "rotation_bytes": rotate_bytes,
            "retention_files": backups,
            "gzip": gzip_enabled,
            "console": console_enabled,
        },
    })

    return _logger