# Jarvis
00002|
00003| ```
00004|      ██╗ █████╗ ██████╗ ██╗   ██╗██╗███████╗
00005|      ██║██╔══██╗██╔══██╗██║   ██║██║██╔════╝
00006|      ██║███████║██████╔╝██║   ██║██║███████╗
00007| ██   ██║██╔══██║██╔══██╗╚██╗ ██╔╝██║╚════██║
00008| ╚█████╔╝██║  ██║██║  ██║ ╚████╔╝ ██║███████║
00009|  ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝
00010| ```
00011|
00012| **The Agentic MCP Server for Dynamic Tool Management**
00013|
00014| <div align="center">
00015|
00016| [![Go](https://img.shields.io/badge/go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
00017| [![License](https://img.shields.io/badge/license-MIT-2EA043)](LICENSE)
00018| [![MCP](https://img.shields.io/badge/MCP-Compliant-6366f1)](https://modelcontextprotocol.io/)
00019| [![Tests](https://img.shields.io/github/actions/workflow/status/JRedeker/Jarvis-mcpm/test.yml?branch=main&label=tests)](https://github.com/JRedeker/Jarvis-mcpm/actions)
00020| [![Go Report](https://goreportcard.com/badge/github.com/JRedeker/Jarvis-mcpm)](https://goreportcard.com/report/github.com/JRedeker/Jarvis-mcpm)
00021|
00022| </div>
00023|
00024| > **Jarvis lets your AI agent manage its own tools.** Install servers, switch profiles, configure clients—all through natural language. One MCP server to rule them all.
00025|
00026| ---
00027|
00028| ## The Problem
00029|
00030| You're using Claude, Cursor, or another AI client. You have MCP servers for memory, search, code tools. But:
00031|
00032| - **Static configs** — Every new tool requires manual JSON editing and client restart
00033| - **Per-client duplication** — Same servers configured separately in Claude, Cursor, VS Code
00034| - **No agent autonomy** — Your agent can't install tools it needs mid-conversation
00035| - **Context switching** — Different projects need different tool sets
00036|
00037| ## The Solution
00038|
00039| Jarvis is an **agentic-first MCP server** that gives your AI agent control over its own tooling:
00040|
00041| ```
00042| You: "I need to analyze some PDFs"
00043| Agent: [searches registry → installs pdf-parse → uses it immediately]
00044| Agent: "Done. The contract has a 30-day payment term on page 3."
00045| ```
00046|
00047| **No config editing. No restart. The agent handles it.**
00048|
00049| ---
00050|
00051| ## Core Capabilities
00052|
00053| ### 1. Dynamic Tool Installation
00054|
00055| Your agent discovers and installs tools from a registry of 200+ MCP servers:
00056|
00057| ```javascript
00058| jarvis_server({ action: "search", query: "pdf" })     // Find tools
00059| jarvis_server({ action: "install", name: "pdf-parse" }) // Install
00060| // Tool is immediately available—no restart needed
00061| ```
00062|
00063| ### 2. Profile-Based Tool Sets
00064|
00065| Group tools into composable profiles. Switch entire toolsets based on project context:
00066|
00067| ```javascript
00068| jarvis_profile({ action: "suggest" })  // Auto-detect best profile for current directory
00069| // Returns: "essentials,dev-core,research" for coding projects
00070| ```
00071|
00072| | Profile | Tools | Use Case |
00073| |---------|-------|----------|
00074| | `essentials` | Time, Fetch | Always-on utilities |
00075| | `dev-core` | Context7, Morph | Coding intelligence |
00076| | `research` | Brave, Firecrawl | Web research (Docker) |
00077| | `memory` | Mem0, Basic Memory | Persistence |
00078| | `data` | Qdrant, Postgres | Heavy DBs |
00079|
00080| ### 3. Multi-Client Management
00081|
00082| Configure Claude Desktop, Cursor, VS Code, OpenCode—all from one place:
00083|
00084| ```javascript
00085| jarvis_client({ action: "edit", client_name: "cursor", add_profile: "memory" })
00086| // Cursor now has memory tools. No manual JSON editing.
00087| ```
00088|
00089| ### 4. Self-Healing Infrastructure
00090|
00091| Agent detects and repairs its own infrastructure:
00092|
00093| ```
00094| You: "My search isn't working"
00095| Agent: [checks status → finds Qdrant down → restarts → confirms healthy]
00096| Agent: "Fixed. Qdrant was down, restarted it. Search should work now."
00097| ```
00098|
00099| ---
00100|
00101| ## Installation
00102|
00103| ```bash
00104| git clone https://github.com/JRedeker/Jarvis-mcpm.git
00105| ./Jarvis-mcpm/scripts/setup-jarvis.sh
00106| ```
00107|
00108| Copy the output JSON into your AI client config, or use `--auto-config` to do it automatically.
00109|
00110| ---
00111|
00112| ## How It Works
00113|
00114| ```mermaid
00115| flowchart TB
00116|     subgraph Clients["AI Clients"]
00117|         Claude[Claude Desktop]
00118|         Cursor[Cursor]
00119|         VSCode[VS Code]
00120|         OpenCode[OpenCode]
00121|     end
00122|
00123|     subgraph Jarvis["Jarvis (Agentic MCP Server)"]
00124|         Gateway[Tool Management API]
00125|         Registry[(200+ Server Registry)]
00126|     end
00127|
00128|     subgraph Daemon["MCPM Daemon"]
00129|         P1[Profile: Research]
00130|         P2[Profile: Memory]
00131|         P3[Profile: Dev-Core]
00132|     end
00133|
00134|     Clients -->|"Install/Switch/Configure"| Gateway
00135|     Gateway --> Registry
00136|     Gateway -->|Manage| Daemon
00137|     Daemon --> P1 & P2 & P3
00138| ```
00139|
00140| **Jarvis sits between your AI clients and the MCP ecosystem.** It's the only MCP server your agent needs to manage all other MCP servers.
00141|
00142| ---
00143|
00144| ## Tool Reference
00145|
00146| 9 tools, action-based routing:
00147|
00148| | Tool | What It Does |
00149| |:-----|:-------------|
00150| | `jarvis_server` | Install, uninstall, search, configure MCP servers |
00151| | `jarvis_profile` | Create and switch tool profiles |
00152| | `jarvis_client` | Configure AI clients (Claude, Cursor, etc.) |
00153| | `jarvis_check_status` | System health diagnostics |
00154| | `jarvis_system` | Bootstrap, restart infrastructure |
00155| | `jarvis_project` | Analyze projects, apply DevOps stacks |
00156| | `jarvis_config` | Manage global settings |
00157| | `jarvis_share` | Share servers via tunnels |
00158| | `jarvis_diagnose` | Debug MCP profile issues (NEW) |
00159|
00160| <details>
00161| <summary><b>Example Commands</b></summary>
00162|
00163| ```javascript
00164| // Install a new tool
00165| jarvis_server({ action: "install", name: "brave-search" })
00166|
00167| // Switch project profile
00168| jarvis_profile({ action: "edit", name: "my-project", add_servers: "context7,firecrawl" })
00169|
00167| // Configure a client
00168| jarvis_client({ action: "edit", client_name: "opencode", add_profile: "memory" })
00169|
00170| // Check system health
00171| jarvis_check_status()
00172|
00173| // Bootstrap everything
00174| jarvis_system({ action: "bootstrap" })
00175|
00176| // Debug when tools fail to load
00177| jarvis_diagnose({ action: "profile_health" })
00178| jarvis_diagnose({ action: "logs", profile: "research" })
00179| ```
00180| </details>
00181|
00182| ---
00183|
00184| ## Universal Compatibility
00185|
00186| **Models:** Any MCP-compatible model — Claude, GPT, Gemini, DeepSeek, Llama
00187|
00188| **Clients:** Claude Desktop, Claude CLI, Cursor, Windsurf, VS Code, Zed, OpenCode, Kilo Code
00189|
00190| <details>
00191| <summary><b>Client Configuration</b></summary>
00192|
00193| ```json
00194| {
00195|   "mcpServers": {
00196|     "jarvis": {
00197|       "command": "/path/to/Jarvis/jarvis",
00198|       "args": []
00199|     }
00200|   }
00201| }
00202| ```
00203|
00204| That's it. Jarvis manages everything else.
00205| </details>
00206|
00207| ---
00208|
00209| ## Why Jarvis?
00210|
00211| | Without Jarvis | With Jarvis |
00212| |----------------|-------------|
00213| | Edit JSON configs manually | Agent installs tools via natural language |
00214| | Restart client for new tools | Hot-load tools mid-conversation |
00215| | Duplicate configs per client | One source of truth, multi-client |
00216| | Fixed tool set per session | Dynamic capabilities on-demand |
00217| | Manual infrastructure repair | Self-healing |
00218|
00219| ---
00220|
00221| ## Documentation
00222|
00223| | Doc | Description |
00224| |-----|-------------|
00225| | [Examples](docs/EXAMPLES.md) | Workflow examples |
00226| | [FAQ](docs/FAQ.md) | Common questions |
00227| | [Architecture](docs/TECHNICAL_ARCHITECTURE.md) | Technical deep dive |
00228| | [Configuration](docs/CONFIGURATION_STRATEGY.md) | Micro-Profile Strategy |
00229| | [Troubleshooting](docs/TROUBLESHOOTING.md) | Issue resolution |
00230|
00231| ---
00232|
00233| ## Contributing
00234|
00235| ```bash
00236| git clone https://github.com/YOUR_USERNAME/Jarvis-mcpm.git
00237| cd Jarvis && go build -o jarvis . && go test -v ./...
00238| ```
00239|
00240| ---
00241|
00242| <div align="center">
00243|
00244| **MIT License** · [Issues](https://github.com/JRedeker/Jarvis-mcpm/issues) · [Discussions](https://github.com/JRedeker/Jarvis-mcpm/discussions)
00245|
00246| </div>
00247|
