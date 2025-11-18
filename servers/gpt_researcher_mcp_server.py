#!/usr/bin/env python3
"""
Wrapper to run the gpt-researcher-mcp Python MCP server via stdio.

This loads the packaged server from node_modules/gpt-researcher-mcp/src
and runs its async main() using the local virtualenv Python.
"""

import asyncio
import sys
from pathlib import Path


def _add_gpt_researcher_mcp_to_path() -> None:
    repo_root = Path(__file__).resolve().parents[1]
    mcp_src = repo_root / "node_modules" / "gpt-researcher-mcp" / "src"
    if mcp_src.is_dir():
        sys.path.insert(0, str(mcp_src))


def main() -> None:
    _add_gpt_researcher_mcp_to_path()
    from gpt_researcher_mcp.server import main as server_main

    asyncio.run(server_main())


if __name__ == "__main__":
    main()

