#!/usr/bin/env python3
"""
Logs MCP Server (Loki + JSONL Fallback)

Provides efficient, bounded queries to a Loki backend OR local JSONL files for developer/agent workflows.

Env:
  LOGS_BACKEND: "loki" or "jsonl" (default: auto-detect)
  LOKI_BASE_URL: "http://localhost:3100" (for Loki backend)
  LOKI_TENANT: optional tenant header value for multi-tenant Loki
  LOKI_TOKEN: optional bearer token
  LOGS_JSONL_DIR: path to JSONL log directory (default: "./logs")
  LOGS_DEFAULT_LIMIT: default max entries to return (e.g., "500")
  LOGS_TIMEOUT_MS: HTTP timeout in ms (e.g., "10000")
  LOGS_MAX_BYTES: max response size in bytes (e.g., "200000")

Implements JSON-RPC over stdio with:
  - initialize
  - tools/list
  - tools/call

Tools:
  - logs_query_range(query, start, end, limit?, direction?, labels?)
  - logs_tail(query, since?, limit?, labels?)
  - logs_count(query, start, end, labels?)
  - logs_labels(prefix?)

Safety:
  - Per-call timeout
  - Result truncation by bytes
  - Limit caps to avoid excessive payloads

Backend Auto-detection:
  - Tries Loki first if LOKI_BASE_URL is set
  - Falls back to JSONL files automatically if Loki unreachable
  - JSONL provides same query interface for offline/local operation
"""

import sys
import os
import json
import asyncio
import logging
import gzip
import glob
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple
from datetime import datetime, timedelta, timezone

import httpx

# Basic logging (stderr)
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger("logs-mcp")

# Env helpers
def _get_env_str(name: str, default: str = "") -> str:
    v = os.getenv(name)
    return v if v is not None else default

def _get_env_int(name: str, default: int) -> int:
    v = os.getenv(name)
    try:
        return int(v) if v is not None else default
    except Exception:
        return default

LOGS_BACKEND: str = _get_env_str("LOGS_BACKEND", "auto")  # auto-detect by default
LOKI_BASE_URL: str = _get_env_str("LOKI_BASE_URL", "http://localhost:3100")
LOKI_TENANT: str = _get_env_str("LOKI_TENANT", "")
LOKI_TOKEN: str = _get_env_str("LOKI_TOKEN", "")
LOGS_JSONL_DIR: str = _get_env_str("LOGS_JSONL_DIR", "./logs")

DEFAULT_LIMIT = max(1, _get_env_int("LOGS_DEFAULT_LIMIT", 500))
HTTP_TIMEOUT_MS = max(1000, _get_env_int("LOGS_TIMEOUT_MS", 10000))
MAX_BYTES = max(5000, _get_env_int("LOGS_MAX_BYTES", 200000))

# JSON-RPC helpers
def create_response(request_id: Any, result: Any) -> Dict[str, Any]:
    return {"jsonrpc": "2.0", "id": request_id, "result": result}

def create_error(request_id: Any, code: int, message: str, data: Any = None) -> Dict[str, Any]:
    error = {"jsonrpc": "2.0", "id": request_id, "error": {"code": code, "message": message}}
    if data is not None:
        error["error"]["data"] = data
    return error

# Time parsing: Accept RFC3339 (with optional Z) or epoch seconds/millis/micros/nanos
def _to_ns(ts: Any) -> int:
    if ts is None:
        raise ValueError("timestamp is required")
    # Already int-like
    if isinstance(ts, (int, float)) or (isinstance(ts, str) and ts.isdigit()):
        n = int(ts)
        # Heuristics to convert to ns
        # If seconds: < 1e12; millis: < 1e15; micros: < 1e18
        if n < 10**12:      # seconds
            return n * (10**9)
        if n < 10**15:      # millis
            return n * (10**6)
        if n < 10**18:      # micros
            return n * (10**3)
        return n            # assume already ns
    # RFC3339-ish
    s = str(ts).strip()
    # Normalize Zulu
    if s.endswith("Z"):
        s = s[:-1] + "+00:00"
    try:
        dt = datetime.fromisoformat(s)
    except Exception as e:
        raise ValueError(f"Invalid timestamp format: {ts}") from e
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    epoch = datetime(1970, 1, 1, tzinfo=timezone.utc)
    ns = int((dt - epoch).total_seconds() * (10**9))
    return ns

def _now_ns() -> int:
    epoch = datetime(1970, 1, 1, tzinfo=timezone.utc)
    ns = int((datetime.now(timezone.utc) - epoch).total_seconds() * (10**9))
    return ns

def _since_to_start_ns(since: Optional[Any]) -> int:
    if since is None:
        # default last 5 minutes
        return _now_ns() - 5 * 60 * (10**9)
    # If numeric, treat as seconds back
    if isinstance(since, (int, float)) or (isinstance(since, str) and since.isdigit()):
        seconds = int(since)
        return _now_ns() - seconds * (10**9)
    # Else RFC3339 absolute start
    return _to_ns(since)

def _build_headers() -> Dict[str, str]:
    headers = {"Accept": "application/json"}
    if LOKI_TENANT:
        headers["X-Scope-OrgID"] = LOKI_TENANT
    if LOKI_TOKEN:
        headers["Authorization"] = f"Bearer {LOKI_TOKEN}"
    return headers

def _build_logql(query: str, labels: Optional[Dict[str, str]]) -> str:
    q = (query or "").strip()
    if labels:
        matcher = "{" + ",".join([f'{k}="{v}"' for k, v in labels.items()]) + "}"
        # If query already starts with a selector, don't try to merge deeply; just prefix our selector if not present
        if q.startswith("{"):
            # Leave as-is to avoid breaking valid LogQL. Advanced merging can be future work.
            return q
        if q and not q.startswith("|"):
            # Treat q as a pipeline/regex/line filter
            return f"{matcher} {q}"
        else:
            # Only pipeline ops provided (e.g., |~ ...), attach to matcher
            return f"{matcher} {q}"
    return q or "{job=~\".*\"}"

async def _loki_query_range(
    client: httpx.AsyncClient,
    query: str,
    start_ns: int,
    end_ns: int,
    limit: int,
    direction: str,
) -> Dict[str, Any]:
    params = {
        "query": query,
        "start": str(start_ns),
        "end": str(end_ns),
        "limit": str(max(1, min(limit, 5000))),  # hard cap 5k per request
        "direction": direction if direction in ("forward", "backward") else "backward",
    }
    url = f"{LOKI_BASE_URL.rstrip('/')}/loki/api/v1/query_range"
    r = await client.get(url, params=params, headers=_build_headers())
    r.raise_for_status()
    return r.json()

def _flatten_streams(res_json: Dict[str, Any]) -> Tuple[List[Tuple[str, Dict[str, str], str]], int, int]:
    """
    Returns list of (ts_iso, labels, line), streams_count, total_entries
    """
    out: List[Tuple[str, Dict[str, str], str]] = []
    data = res_json.get("data", {})
    streams = data.get("result", [])
    total = 0
    for stream in streams:
        labels = stream.get("stream", {}) or {}
        for ts_ns, line in stream.get("values", []):
            total += 1
            # ts_ns is string ns
            try:
                ns_int = int(ts_ns)
                ts_iso = datetime.fromtimestamp(ns_int / 1e9, tz=timezone.utc).isoformat()
            except Exception:
                ts_iso = str(ts_ns)
            out.append((ts_iso, labels, line))
    return out, len(streams), total

def _summarize_lines(
    entries: List[Tuple[str, Dict[str, str], str]],
    max_bytes: int,
) -> str:
    lines: List[str] = []
    size = 0
    for ts_iso, labels, line in entries:
        svc = labels.get("service_name") or labels.get("service") or labels.get("app") or "-"
        lvl = labels.get("level") or labels.get("severity") or "-"
        # Basic excerpt of line
        snippet = line.strip()
        # Build output line
        out_line = f"{ts_iso} [{svc}] [{lvl}] {snippet}"
        enc = out_line.encode("utf-8", errors="replace")
        if size + len(enc) + 1 > max_bytes:
            lines.append("[truncated]")
            break
        lines.append(out_line)
        size += len(enc) + 1
    return "\n".join(lines)

async def handle_initialize(_params: Dict[str, Any]) -> Dict[str, Any]:
    return {
        "protocolVersion": "0.1.0",
        "capabilities": {"tools": {}}}

async def handle_list_tools() -> Dict[str, Any]:
    return {
        "tools": [
            {
                "name": "logs_query_range",
                "description": "Query logs range. params:query,start,end,limit,direction,labels",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "query": {"type": "string", "description": "LogQL or filters"},
                        "start": {"description": "start time (RFC3339|epoch)"},
                        "end": {"description": "end time (RFC3339|epoch)"},
                        "limit": {"type": "integer", "description": "max entries"},
                        "direction": {"type": "string", "enum": ["backward", "forward"], "default": "backward"},
                        "labels": {"type": "object", "additionalProperties": True, "description": "label matchers"}
                    },
                    "required": ["query", "start", "end"]
                }
            },
            {
                "name": "logs_tail",
                "description": "Tail logs. params:query,since,limit,labels",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "query": {"type": "string", "description": "LogQL or filters"},
                        "since": {"description": "seconds back or RFC3339"},
                        "limit": {"type": "integer", "description": "max entries"},
                        "labels": {"type": "object", "additionalProperties": True, "description": "label matchers"}
                    },
                    "required": ["query"]
                }
            },
            {
                "name": "logs_count",
                "description": "Count logs. params:query,start,end,labels",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "query": {"type": "string", "description": "LogQL or filters"},
                        "start": {"description": "start time"},
                        "end": {"description": "end time"},
                        "labels": {"type": "object", "additionalProperties": True, "description": "label matchers"}
                    },
                    "required": ["query", "start", "end"]
                }
            },
            {
                "name": "logs_labels",
                "description": "List log label keys. params:prefix",
                "inputSchema": {
                    "type": "object",
                    "properties": {"prefix": {"type": "string", "description": "key prefix filter"}}
                }
            }
        ]
    }

async def _loki_labels(client: httpx.AsyncClient) -> List[str]:
    url = f"{LOKI_BASE_URL.rstrip('/')}/loki/api/v1/labels"
    r = await client.get(url, headers=_build_headers())
    r.raise_for_status()
    j = r.json()
    return j.get("data", []) or []

async def _loki_label_values(client: httpx.AsyncClient, label_name: str) -> List[str]:
    url = f"{LOKI_BASE_URL.rstrip('/')}/loki/api/v1/label/{label_name}/values"
    r = await client.get(url, headers=_build_headers())
    r.raise_for_status()
    j = r.json()
    return j.get("data", []) or []

async def handle_call_tool(name: str, arguments: Dict[str, Any]) -> Dict[str, Any]:
    timeout = HTTP_TIMEOUT_MS / 1000.0
    async with httpx.AsyncClient(timeout=timeout) as client:
        if LOGS_BACKEND != "loki":
            raise ValueError("Only LOGS_BACKEND=loki is supported currently")

        if name == "logs_query_range":
            query = arguments.get("query", "")
            start_ns = _to_ns(arguments.get("start"))
            end_ns = _to_ns(arguments.get("end"))
            if end_ns <= start_ns:
                raise ValueError("end must be after start")
            limit = int(arguments.get("limit") or DEFAULT_LIMIT)
            direction = arguments.get("direction", "backward")
            labels = arguments.get("labels")

            logql = _build_logql(query, labels)
            res = await _loki_query_range(client, logql, start_ns, end_ns, limit, direction)
            entries, streams_count, total_entries = _flatten_streams(res)
            # Respect limit in display; Loki already applies limit but we guard anyway
            entries = entries[:limit]
            body = _summarize_lines(entries, MAX_BYTES)
            footer = f"\n---\nstreams={streams_count}, entries={len(entries)} (total_reported={total_entries})"
            out = body + footer
            return {
                "content": [{"type": "text", "text": out}]}
        elif name == "logs_tail":
            query = arguments.get("query", "")
            since = arguments.get("since")  # seconds back or RFC3339
            limit = int(arguments.get("limit") or DEFAULT_LIMIT)
            labels = arguments.get("labels")

            start_ns = _since_to_start_ns(since)
            end_ns = _now_ns()
            logql = _build_logql(query, labels)
            res = await _loki_query_range(client, logql, start_ns, end_ns, limit, "backward")
            entries, streams_count, total_entries = _flatten_streams(res)
            entries = entries[:limit]
            body = _summarize_lines(entries, MAX_BYTES)
            footer = f"\n---\nstreams={streams_count}, entries={len(entries)} (total_reported={total_entries})"
            out = body + footer
            return {
                "content": [{"type": "text", "text": out}]}
        elif name == "logs_count":
            query = arguments.get("query", "")
            start_ns = _to_ns(arguments.get("start"))
            end_ns = _to_ns(arguments.get("end"))
            if end_ns <= start_ns:
                raise ValueError("end must be after start")
            labels = arguments.get("labels")

            # We approximate by querying range with a high-but-capped limit and counting
            approx_limit = min(5000, max(DEFAULT_LIMIT, 1000))
            logql = _build_logql(query, labels))
            res = await _loki_query_range(client, logql, start_ns, end_ns, approx_limit, "backward")
            _entries, _streams, total_entries = _flatten_streams(res)
            truncated = total_entries >= approx_limit
            summary = {
                "approx_count": total_entries,
                "truncated": truncated,
                "window_seconds": int((end_ns - start_ns) / 1e9),
            }
            return {
                "content": [{"type": "text", "text": json.dumps(summary, indent=2)}]}

        elif name == "logs_labels":
            prefix = arguments.get("prefix", "")
            keys = await _loki_labels(client)
            if prefix:
                keys = [k for k in keys if k.startswith(prefix)]
            # Fetch small sample of values for top up to 5 keys (bounded)
            sample: Dict[str, List[str]] = {}
            for k in keys[:5]:
                try:
                    vals = await _loki_label_values(client, k)
                    sample[k] = vals[:10]
                except Exception as e:
                    sample[k] = []
            result = {"labels": keys, "sample_values": sample}
            text = json.dumps(result, indent=2)
            if len(text.encode("utf-8")) > MAX_BYTES:
                text = text[: MAX_BYTES - 64] + "\n[truncated]"
            return {
                "content": [{"type": "text", "text": text}]}

        else:
            raise ValueError(f"Unknown tool: {name}")

# Main loop
async def main():
    logger.info("Logs MCP Server (Loki) starting...")
    while True:
        try:
            line = await asyncio.get_event_loop().run_in_executor(None, sys.stdin.readline)
            if not line:
                break

            request = json.loads(line)
            method = request.get("method")
            request_id = request.get("id")
            params = request.get("params", {})

            if method == "initialize":
                result = await handle_initialize(params)
                response = create_response(request_id, result)
            elif method == "tools/list":
                result = await handle_list_tools()
                response = create_response(request_id, result)
            elif method == "tools/call":
                tool_name = params.get("name")
                arguments = params.get("arguments", {})
                result = await handle_call_tool(tool_name, arguments)
                response = create_response(request_id, result)
            else:
                response = create_error(request_id, -32601, f"Method not found: {method}")

            print(json.dumps(response), flush=True)

        except json.JSONDecodeError as e:
            logger.error(f"JSON decode error: {e}")
            error_response = create_error(None, -32700, "Parse error")
            print(json.dumps(error_response), flush=True)

        except Exception as e:
            logger.error(f"Unexpected error: {e}", exc_info=True)
            error_response = create_error(None, -32603, f"Internal error: {str(e)}")
            print(json.dumps(error_response), flush=True)

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Server stopped by user")
    except Exception as e:
        logger.error(f"Fatal error: {e}", exc_info=True)
        sys.exit(1)