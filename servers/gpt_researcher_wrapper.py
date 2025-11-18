#!/usr/bin/env python3
"""
Plain passthrough wrapper for gpt-researcher-mcp: no framing edits, just banner drop.
(Currently unused when cipher.yml points directly to npx.)
"""
import os
import subprocess
import sys
import threading

BANNER_PREFIXES = (
    "None of PyTorch",
    "Content-Length: 0",
)


def forward_stdin(child):
    for chunk in iter(lambda: sys.stdin.buffer.read1(4096), b""):
        try:
            child.stdin.buffer.write(chunk)
            child.stdin.buffer.flush()
        except BrokenPipeError:
            break


def forward_stdout(child):
    for chunk in iter(lambda: child.stdout.buffer.read1(4096), b""):
        text = chunk.decode(errors="ignore")
        if any(text.startswith(pref) for pref in BANNER_PREFIXES):
            continue
        sys.stdout.buffer.write(chunk)
        sys.stdout.buffer.flush()


def forward_stderr(child):
    for chunk in iter(lambda: child.stderr.buffer.read1(4096), b""):
        sys.stderr.buffer.write(chunk)
        sys.stderr.buffer.flush()


def main():
    env = os.environ.copy()
    env["PATH"] = os.pathsep.join([
        os.path.expanduser("~/.local/bin"),
        os.path.expanduser("~/bin"),
        "/usr/local/bin",
        "/usr/bin",
        "/bin",
    ])
    cmd = ["npx", "-y", "gpt-researcher-mcp", "--stdio"]
    child = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        bufsize=0,
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
