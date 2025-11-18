#!/usr/bin/env python3
"""
Text-mode wrapper for @modelcontextprotocol/server-memory.
Passes stdin/stdout through and drops noisy banner lines.
"""
import os
import subprocess
import sys
import threading

BANNER_PREFIXES = (
    "Knowledge Graph MCP Server",
)


def forward_stdin(child):
    """Forward parent stdin to child stdin (text mode)."""
    for line in sys.stdin:
        try:
            child.stdin.write(line)
            child.stdin.flush()
        except BrokenPipeError:
            break


def forward_stdout(child):
    """Forward child stdout to parent, skipping banner lines."""
    for line in child.stdout:
        if any(line.startswith(pref) for pref in BANNER_PREFIXES):
            continue
        sys.stdout.write(line)
        sys.stdout.flush()


def forward_stderr(child):
    for line in child.stderr:
        sys.stderr.write(line)
        sys.stderr.flush()


def main():
    env = os.environ.copy()
    env.setdefault("DISABLE_BANNER", "1")
    env["PATH"] = os.pathsep.join(
        [
            os.path.expanduser("~/.local/bin"),
            os.path.expanduser("~/bin"),
            "/usr/local/bin",
            "/usr/bin",
            "/bin",
        ]
    )
    cmd = ["npx", "-y", "@modelcontextprotocol/server-memory", "--stdio"]
    child = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=1,
        env=env,
    )

    threads = [
        threading.Thread(target=forward_stdin, args=(child,), daemon=True),
        threading.Thread(target=forward_stdout, args=(child,), daemon=True),
        threading.Thread(target=forward_stderr, args=(child,), daemon=True),
    ]
    for t in threads:
        t.start()
    child.wait()


if __name__ == "__main__":
    main()
