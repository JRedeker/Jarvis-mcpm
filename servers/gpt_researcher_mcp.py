#!/usr/bin/env python3
"""
Lightweight GPT Researcher MCP server (stdio).

This server talks directly to the `gpt_researcher` library and exposes a small
set of MCP tools:

- deep_research(query): run full web research and return context + sources
- quick_search(query): fast search results optimized for speed
- write_report(research_id, custom_prompt?): generate a report from prior research

It implements the minimal JSON-RPC / MCP handshake used by cipher-aggregator:
- initialize
- tools/list
- tools/call
- resources/list, prompts/list (empty)
"""

import asyncio
import json
import logging
import sys
import uuid
from typing import Any, Dict, List, Optional

from gpt_researcher import GPTResearcher


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


# In-memory store tying research_id -> GPTResearcher instance
_RESEARCHERS: Dict[str, GPTResearcher] = {}


def _send_framed(payload: Dict[str, Any]) -> None:
    """Emit a Content-Length framed JSON-RPC message to stdout."""
    data = json.dumps(payload).encode("utf-8")
    sys.stdout.write(f"Content-Length: {len(data)}\n\n")
    sys.stdout.flush()
    sys.stdout.buffer.write(data)
    sys.stdout.flush()


def send_error(code: int, message: str, req_id: Optional[Any] = None, data: Any = None) -> None:
    error: Dict[str, Any] = {
        "jsonrpc": "2.0",
        "error": {"code": code, "message": message},
    }
    if req_id is not None:
        error["id"] = req_id
    if data is not None:
        error["error"]["data"] = data
    _send_framed(error)


def read_request() -> Optional[Dict[str, Any]]:
    """Read a request from stdin supporting MCP Content-Length framing or newline-delimited JSON."""
    while True:
        header = sys.stdin.readline()
        if not header:
            return None
        if header.startswith("Content-Length"):
            try:
                length = int(header.split(":", 1)[1].strip())
            except Exception:
                continue
            # consume blank separator
            sys.stdin.readline()
            body = sys.stdin.read(length)
            if not body:
                return None
            try:
                return json.loads(body)
            except json.JSONDecodeError:
                send_error(-32700, "Invalid JSON")
                continue
        else:
            header = header.strip()
            if not header:
                continue
            try:
                return json.loads(header)
            except json.JSONDecodeError:
                send_error(-32700, "Invalid JSON")
                continue


def list_tools() -> List[Dict[str, Any]]:
    """Describe available GPT Researcher tools."""
    return [
        {
            "name": "deep_research",
            "description": (
                "Run deep web research on a query using GPT Researcher and return "
                "context plus sources and URLs."
            ),
            "inputSchema": {
                "type": "object",
                "properties": {
                    "query": {
                        "type": "string",
                        "description": "Research topic or question.",
                    }
                },
                "required": ["query"],
            },
        },
        {
            "name": "quick_search",
            "description": (
                "Run a fast web search via GPT Researcher; optimized for speed over depth."
            ),
            "inputSchema": {
                "type": "object",
                "properties": {
                    "query": {
                        "type": "string",
                        "description": "Search query.",
                    }
                },
                "required": ["query"],
            },
        },
        {
            "name": "write_report",
            "description": (
                "Generate a research report for a previously-run deep_research call."
            ),
            "inputSchema": {
                "type": "object",
                "properties": {
                    "research_id": {
                        "type": "string",
                        "description": "ID returned from deep_research.",
                    },
                    "custom_prompt": {
                        "type": "string",
                        "description": "Optional custom report prompt.",
                    },
                },
                "required": ["research_id"],
            },
        },
    ]


async def _deep_research_async(query: str) -> Dict[str, Any]:
    researcher = GPTResearcher(query)
    logger.info("Starting deep research for query: %s", query)
    await researcher.conduct_research()
    research_id = str(uuid.uuid4())
    _RESEARCHERS[research_id] = researcher

    context = researcher.get_research_context()
    sources = researcher.get_research_sources()
    source_urls = researcher.get_source_urls()
    try:
        costs = researcher.get_costs()
    except Exception:
        costs = None

    return {
        "status": "success",
        "research_id": research_id,
        "query": query,
        "context": context,
        "sources": sources,
        "source_urls": source_urls,
        "costs": costs,
    }


async def _quick_search_async(query: str) -> Dict[str, Any]:
    researcher = GPTResearcher(query)
    logger.info("Starting quick search for query: %s", query)
    results = await researcher.quick_search(query=query)
    return {
        "status": "success",
        "query": query,
        "results": results,
        "result_count": len(results) if results else 0,
    }


async def _write_report_async(research_id: str, custom_prompt: Optional[str]) -> Dict[str, Any]:
    researcher = _RESEARCHERS.get(research_id)
    if not researcher:
        return {
            "status": "error",
            "error": f"Unknown research_id: {research_id}",
        }
    logger.info("Generating report for research_id: %s", research_id)
    report = await researcher.write_report(custom_prompt=custom_prompt)
    try:
        sources = researcher.get_research_sources()
        costs = researcher.get_costs()
    except Exception:
        sources = None
        costs = None
    return {
        "status": "success",
        "research_id": research_id,
        "report": report,
        "sources": sources,
        "costs": costs,
    }


def handle_tool_call(params: Dict[str, Any]) -> Dict[str, Any]:
    """Dispatch MCP tools to GPT Researcher operations."""
    name = params.get("name")
    arguments = params.get("arguments", {}) or {}

    try:
        if name == "deep_research":
            query = arguments.get("query")
            if not query:
                return {"status": "error", "error": "Missing required argument: query"}
            return asyncio.run(_deep_research_async(query))

        if name == "quick_search":
            query = arguments.get("query")
            if not query:
                return {"status": "error", "error": "Missing required argument: query"}
            return asyncio.run(_quick_search_async(query))

        if name == "write_report":
            research_id = arguments.get("research_id")
            if not research_id:
                return {
                    "status": "error",
                    "error": "Missing required argument: research_id",
                }
            custom_prompt = arguments.get("custom_prompt")
            return asyncio.run(_write_report_async(research_id, custom_prompt))

        return {"status": "error", "error": f"Unknown tool: {name}"}
    except Exception as exc:
        logger.exception("Error in GPT Researcher tool %s", name)
        return {"status": "error", "error": str(exc)}


def main() -> None:
    """Main MCP stdio loop."""
    while True:
        request = read_request()
        if request is None:
            break

        try:
            method = request.get("method")
            req_id = request.get("id")

            if method == "initialize":
                response = {
                    "jsonrpc": "2.0",
                    "result": {
                        "protocolVersion": "2024-11-05",
                        "capabilities": {
                            "tools": {"listChanged": False},
                            "resources": {"listChanged": False},
                            "prompts": {"listChanged": False},
                        },
                        "serverInfo": {
                            "name": "gpt-researcher-mcp-light",
                            "version": "0.1.0",
                        },
                    },
                    "id": req_id,
                }
                _send_framed(response)
            elif method == "tools/list":
                _send_framed(
                    {
                        "jsonrpc": "2.0",
                        "result": {"tools": list_tools()},
                        "id": req_id,
                    }
                )
            elif method == "tools/call":
                params = request.get("params", {})
                params.setdefault("arguments", {})
                result = handle_tool_call(params)
                _send_framed({"jsonrpc": "2.0", "result": result, "id": req_id})
            elif method in ("resources/list", "prompts/list"):
                key = "resources" if method == "resources/list" else "prompts"
                _send_framed({"jsonrpc": "2.0", "result": {key: []}, "id": req_id})
            elif method == "notifications/message":
                _send_framed({"jsonrpc": "2.0", "result": {}, "id": req_id})
            else:
                send_error(-32601, f"Unknown method: {method}", req_id=req_id)
        except Exception as exc:
            send_error(-32000, f"Server error: {exc}", req_id=request.get("id"))


if __name__ == "__main__":
    main()
