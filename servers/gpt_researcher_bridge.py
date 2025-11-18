#!/usr/bin/env python3
"""
Bridge for gpt-researcher-mcp:
- Convert parent MCP Content-Length framing to plain JSON lines for the child
- Convert child's plain JSON lines back to Content-Length framed responses
"""

import datetime
import json
import os
import subprocess
import sys
import threading
from typing import TextIO

LOG_PATH = os.environ.get("GPT_RESEARCHER_BRIDGE_LOG", "/tmp/gpt_researcher_bridge.log")


def log(message: str):
    timestamp = datetime.datetime.utcnow().isoformat()
    line = f"[{timestamp}] {message}"
    try:
        with open(LOG_PATH, "a", encoding="utf-8") as fh:
            fh.write(line + "\n")
    except Exception:
        pass
    try:
        sys.stderr.write(line + "\n")
        sys.stderr.flush()
    except Exception:
        pass


def preview(payload: str, limit: int = 240) -> str:
    payload = payload.replace("\n", "\\n")
    if len(payload) > limit:
        return payload[:limit] + "…"
    return payload


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
            log("Parent stream closed; stopping forward_to_child")
            break
        try:
            rewritten = msg
            if msg.strip().startswith("{"):
                try:
                    payload = json.loads(msg)
                    if payload.get("method") == "initialize":
                        params = payload.setdefault("params", {})
                        original_version = params.get("protocolVersion")
                        if original_version != "2024-11-05":
                            log(
                                "Rewriting protocolVersion "
                                f"{original_version} -> 2024-11-05"
                            )
                        params["protocolVersion"] = "2024-11-05"

                        caps = params.setdefault("capabilities", {})
                        prompts = caps.setdefault("prompts", {})
                        if "listChanged" not in prompts:
                            prompts["listChanged"] = False

                        resources = caps.setdefault("resources", {})
                        if "listChanged" not in resources:
                            resources["listChanged"] = False
                        if "subscribe" not in resources:
                            resources["subscribe"] = False

                        tools = caps.setdefault("tools", {})
                        if "listChanged" not in tools:
                            tools["listChanged"] = False

                        caps.setdefault("experimental", {})

                        payload.setdefault("clientInfo", {}).setdefault("name", "cipher-mcp-client")
                        payload["clientInfo"].setdefault("version", "0.0.0")

                        rewritten = json.dumps(payload)
                except json.JSONDecodeError:
                    log(f"Failed to parse parent message for rewrite: {preview(msg)}")
            child_stdin.write(rewritten)
            if not rewritten.endswith("\n"):
                child_stdin.write("\n")
            child_stdin.flush()
            if rewritten is not msg:
                log(f"Parent → Child (rewritten): {preview(rewritten)}")
            else:
                log(f"Parent → Child: {preview(msg)}")
        except BrokenPipeError:
            log("BrokenPipeError while writing to child stdin")
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
        log(f"Child → Parent: {preview(text)}")
        frame_and_emit(text)
    log("Child stdout forwarder exiting")


def forward_stderr(child_stderr):
    for line in child_stderr:
        data = line if isinstance(line, (bytes, bytearray)) else line.encode()
        try:
            decoded = data.decode("utf-8", "ignore").rstrip()
        except Exception:
            decoded = repr(data)
        log(f"Child STDERR: {decoded}")
        sys.stderr.buffer.write(data)
        sys.stderr.buffer.flush()
    log("Child stderr forwarder exiting")


def main():
    env = os.environ.copy()
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    venv_bin = os.path.join(repo_root, ".venv", "bin")
    env["PATH"] = os.pathsep.join(
        [
            venv_bin,
            os.path.expanduser("~/.local/bin"),
            os.path.expanduser("~/bin"),
            "/usr/local/bin",
            "/usr/bin",
            "/bin",
            env.get("PATH", ""),
        ]
    )
    env.setdefault("PYTHON", os.path.join(venv_bin, "python"))
    env.setdefault("VIRTUAL_ENV", os.path.join(repo_root, ".venv"))
    log(f"Launching gpt-researcher-mcp with PATH={env['PATH']}")
    log(f"PYTHON={env.get('PYTHON')} VIRTUAL_ENV={env.get('VIRTUAL_ENV')}")
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
    log(f"Spawned child PID {child.pid}")

    threads = [
        threading.Thread(target=forward_to_child, args=(child.stdin,), daemon=True),
        threading.Thread(target=forward_from_child, args=(child.stdout,), daemon=True),
        threading.Thread(target=forward_stderr, args=(child.stderr,), daemon=True),
    ]
    for t in threads:
        t.start()
    return_code = child.wait()
    log(f"Child process exited with code {return_code}")


if __name__ == "__main__":
    main()
