#!/usr/bin/env python3
"""
Bridge for gpt-researcher-mcp:
- Convert parent MCP Content-Length framing to plain JSON lines for the child
- Convert child's plain JSON lines back to Content-Length framed responses
"""

import json
import os
import subprocess
import sys
import threading


def read_framed_request():
    """Read a single MCP message from stdin (Content-Length or plain JSON line)."""
    header = sys.stdin.readline()
    if not header:
        return None
    if not header.startswith("Content-Length"):
        # Treat as plain JSON line from parent
        return header.strip()
    try:
        length = int(header.split(":", 1)[1].strip())
    except Exception:
        return None
    _blank = sys.stdin.readline()  # consume blank line
    body = sys.stdin.read(length)
    return body


def forward_to_child(child_stdin):
    """Read framed requests from parent, send plain JSON lines to child."""
    while True:
        msg = read_framed_request()
        if msg is None:
            break
        try:
            child_stdin.write(msg)
            if not msg.endswith("\n"):
                child_stdin.write("\n")
            child_stdin.flush()
        except BrokenPipeError:
            break


def frame_and_emit(payload: str):
    data = payload.encode("utf-8")
    sys.stdout.write(f"Content-Length: {len(data)}\n\n")
    sys.stdout.flush()
    sys.stdout.buffer.write(data)
    sys.stdout.flush()


def forward_from_child(child_stdout):
    """Read plain JSON lines from child, emit Content-Length framed to parent."""
    for line in child_stdout:
        if not line:
            break
        text = line if isinstance(line, str) else line.decode("utf-8", "ignore")
        text = text.strip()
        if not text:
            continue
        frame_and_emit(text)


def forward_stderr(child_stderr):
    for line in child_stderr:
        sys.stderr.buffer.write(line if isinstance(line, (bytes, bytearray)) else line.encode())
        sys.stderr.buffer.flush()


def main():
    env = os.environ.copy()
    env["PATH"] = os.pathsep.join(
        [
            os.path.expanduser("~/.local/bin"),
            os.path.expanduser("~/bin"),
            "/usr/local/bin",
            "/usr/bin",
            "/bin",
        ]
    )
    cmd = ["npx", "-y", "gpt-researcher-mcp", "--stdio"]
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
        threading.Thread(target=forward_to_child, args=(child.stdin,), daemon=True),
        threading.Thread(target=forward_from_child, args=(child.stdout,), daemon=True),
        threading.Thread(target=forward_stderr, args=(child.stderr,), daemon=True),
    ]
    for t in threads:
        t.start()
    child.wait()


if __name__ == "__main__":
    main()
