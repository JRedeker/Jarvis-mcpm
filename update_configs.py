import json
import os


def update_vscode_settings():
    path = os.path.expanduser("~/.config/Code/User/settings.json")
    if not os.path.exists(path):
        print(f"VS Code settings not found at {path}")
        return

    try:
        with open(path, "r") as f:
            # basic cleaning of potential comments if any (simple lines starting with //)
            content = f.read()
            # valid json check
            data = json.loads(content)

        # Update MCP section
        # We want to replace the stdio config with SSE config for p-pokeedge
        if "mcp" not in data:
            data["mcp"] = {"servers": {}}

        if "servers" not in data["mcp"]:
            data["mcp"]["servers"] = {}

        # Remove old stdio entry if it exists (key might be different or same)
        # The user had "mcpm_profile_p-pokeedge"
        keys_to_remove = [
            k for k in data["mcp"]["servers"] if "pokeedge" in k or "mcpm" in k
        ]
        for k in keys_to_remove:
            del data["mcp"]["servers"][k]

        # Add new SSE entry
        data["mcp"]["servers"]["p-pokeedge"] = {
            "url": "http://localhost:6276/sse",
            "transport": "sse",
            "alwaysAllow": [
                "brave-search_brave_web_search",
                "brave-search_brave_local_search",
                "brave-search_brave_news_search",
                "brave-search_brave_summarizer",
                "context7_resolve-library-id",
                "context7_get-library-docs",
                "firecrawl_firecrawl_scrape",
                "firecrawl_firecrawl_map",
                "firecrawl_firecrawl_search",
                "firecrawl_firecrawl_crawl",
                "firecrawl_firecrawl_check_crawl_status",
                "firecrawl_firecrawl_extract",
                "fetch-mcp_fetch_html",
                "fetch-mcp_fetch_json",
                "fetch-mcp_fetch_markdown",
                "fetch-mcp_fetch_txt",
                "kagimcp_kagi_search_fetch",
                "kagimcp_kagi_summarizer",
                "time_get_current_time",
                "time_convert_time",
            ],
        }

        # Create backup
        with open(path + ".backup", "w") as f:
            f.write(content)

        # Write back
        with open(path, "w") as f:
            json.dump(data, f, indent=2)

        print("✅ VS Code settings updated.")

    except json.JSONDecodeError as e:
        print(f"❌ Failed to parse VS Code settings (might contain comments): {e}")
    except Exception as e:
        print(f"❌ Error updating VS Code settings: {e}")


def update_claude_desktop_config():
    path = os.path.expanduser("~/.config/Claude/claude_desktop_config.json")
    if not os.path.exists(path):
        print(f"Claude Desktop config not found at {path}")
        return

    try:
        with open(path, "r") as f:
            data = json.load(f)

        # Define the new configuration mapping
        new_configs = {
            "p-pokeedge": {
                "url": "http://localhost:6276/sse",
                "transport": "sse",
                "alwaysAllow": data.get("mcpServers", {})
                .get("p-pokeedge", {})
                .get("alwaysAllow", []),
            },
            "memory": {"url": "http://localhost:6277/sse", "transport": "sse"},
            "morph": {"url": "http://localhost:6278/sse", "transport": "sse"},
            "qdrant": {"url": "http://localhost:6279/sse", "transport": "sse"},
        }

        # Update the mcpServers dict
        if "mcpServers" not in data:
            data["mcpServers"] = {}

        # Ensure Jarvis is kept/added as stdio
        data["mcpServers"]["jarvis"] = {
            "command": "/home/jrede/dev/MCP/Jarvis/jarvis",
            "args": [],
        }

        # Apply SSE configs
        for name, config in new_configs.items():
            data["mcpServers"][name] = config

        # Create backup
        with open(path + ".backup", "w") as f:
            with open(path, "r") as original:
                f.write(original.read())

        # Write back
        with open(path, "w") as f:
            json.dump(data, f, indent=2)

        print("✅ Claude Desktop config updated.")

    except Exception as e:
        print(f"❌ Error updating Claude Desktop config: {e}")


if __name__ == "__main__":
    print("Starting config updates...")
    update_vscode_settings()
    update_claude_desktop_config()
    print("Done.")
