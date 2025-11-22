import json
import os
import subprocess
import sys
from typing import Optional, Dict, Any

# Configuration
MCP_JUNGLE_CONFIG_DIR = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "config", "jarvis", "servers")
MCPM_EXECUTABLE = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), ".mcpm_venv", "bin", "mcpm")

class Router:
    def __init__(self):
        self.jungle_registry = self._load_jungle_registry()

    def _load_jungle_registry(self) -> Dict[str, Any]:
        """Loads available servers from the static mcp-jungle configuration."""
        registry = {}
        if not os.path.exists(MCP_JUNGLE_CONFIG_DIR):
            print(f"Warning: Config directory {MCP_JUNGLE_CONFIG_DIR} not found.", file=sys.stderr)
            return registry

        for filename in os.listdir(MCP_JUNGLE_CONFIG_DIR):
            if filename.endswith(".json"):
                server_name = filename[:-5]
                try:
                    with open(os.path.join(MCP_JUNGLE_CONFIG_DIR, filename), 'r') as f:
                        config = json.load(f)
                        registry[server_name] = config
                except json.JSONDecodeError:
                    print(f"Error decoding {filename}", file=sys.stderr)
        return registry

    def _check_mcpm_registry(self, tool_name: str) -> bool:
        """Checks if a tool is installed via mcpm."""
        try:
            # mcpm ls returns a list of installed servers
            result = subprocess.run([MCPM_EXECUTABLE, "ls"], capture_output=True, text=True)
            if result.returncode == 0:
                # Simple check: look for the tool name in the output
                return tool_name in result.stdout
            return False
        except FileNotFoundError:
            return False

    def _install_via_mcpm(self, tool_name: str) -> bool:
        """Attempts to install a tool via mcpm."""
        print(f"Attempting to install {tool_name} via mcpm...", file=sys.stderr)
        try:
            result = subprocess.run([MCPM_EXECUTABLE, "install", tool_name], capture_output=True, text=True)
            if result.returncode == 0:
                print(f"Successfully installed {tool_name}.", file=sys.stderr)
                return True
            else:
                print(f"Failed to install {tool_name}. Error: {result.stderr}", file=sys.stderr)
                return False
        except Exception as e:
            print(f"Error running mcpm: {e}", file=sys.stderr)
            return False

    def get_tool(self, tool_name: str, allow_dynamic_install: bool = True) -> Optional[Dict[str, Any]]:
        """
        Routes the request for a tool/server.
        1. Check Core Registry (mcp-jungle)
        2. Check Dynamic Registry (mcpm)
        3. Attempt Dynamic Install
        """
        # 1. Check Core Registry
        if tool_name in self.jungle_registry:
            return {
                "source": "mcp-jungle",
                "name": tool_name,
                "config": self.jungle_registry[tool_name]
            }

        # 2. Check Dynamic Registry (mcpm)
        if self._check_mcpm_registry(tool_name):
            return {"source": "mcpm", "name": tool_name, "status": "installed"}

        # 3. Attempt Dynamic Install
        if allow_dynamic_install:
            if self._install_via_mcpm(tool_name):
                return {"source": "mcpm", "name": tool_name, "status": "just_installed"}

        return None

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python router.py <tool_name>")
        sys.exit(1)

    tool = sys.argv[1]
    router = Router()
    result = router.get_tool(tool)
    if result:
        print(json.dumps(result, indent=2))
    else:
        print(f"Tool '{tool}' not found.", file=sys.stderr)
        sys.exit(1)