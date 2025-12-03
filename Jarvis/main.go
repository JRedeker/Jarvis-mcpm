package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"jarvis/handlers"
)

// Command-line flags for transport mode
var (
	httpMode = flag.Bool("http", false, "Run as Streamable HTTP server instead of stdio")
	httpPort = flag.String("port", "6275", "Port for HTTP server (default: 6275)")
	httpHost = flag.String("host", "127.0.0.1", "Host for HTTP server (default: 127.0.0.1)")
)

var (
	sharedServers      = make(map[string]*exec.Cmd)
	sharedServersMutex sync.Mutex
	logFile            *os.File
)

func setupLogging() {
	// Determine project root (assuming Jarvis runs from Jarvis/ or project root)
	// We'll try to find the 'logs' directory in the parent or current directory
	logDir := "logs"
	if _, err := os.Stat("../logs"); err == nil {
		logDir = "../logs"
	} else {
		os.MkdirAll("logs", 0755)
	}

	logPath := filepath.Join(logDir, "jarvis.log")
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stderr if file creation fails
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return
	}

	// Create a MultiWriter to write to both the log file and stderr (so IDEs still see errors)
	// However, for pure logging, we might want just the file to avoid polluting the protocol stream
	// if logging libraries print to stdout/stderr by default.
	// mcp-go/server.WithLogging() uses stderr by default.

	// We will set the global logger to write to the file
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf(">> Jarvis Logging Initialized <<")
}

func printBanner() {
	banner := `
     ██╗ █████╗ ██████╗ ██╗   ██╗██╗███████╗
     ██║██╔══██╗██╔══██╗██║   ██║██║██╔════╝
     ██║███████║██████╔╝██║   ██║██║███████╗
██   ██║██╔══██║██╔══██╗╚██╗ ██╔╝██║╚════██║
╚█████╔╝██║  ██║██║  ██║ ╚████╔╝ ██║███████║
 ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝
`
	// Print to Stderr to avoid interfering with MCP stdio protocol
	fmt.Fprintln(os.Stderr, "\033[36m"+banner+"\033[0m")
	fmt.Fprintln(os.Stderr, "\033[1;32m>> JARVIS MCP Gateway v1.0.0 initialized <<\033[0m")
}

func main() {
	flag.Parse()
	setupLogging()
	printBanner()

	// Run smoke tests on boot if enabled
	if shouldRunSmokeTests() {
		runBootSmokeTests()
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"jarvis",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(), // This logs MCP protocol messages to stderr
	)

	// Register all core tools from handlers package (TDD-tested)
	h := handlers.CreateProductionHandler()
	handlers.RegisterToolsWithMCPServer(s, h)

	// Legacy tools that use global state (not yet migrated to handlers package)
	// TODO: Migrate these to handlers package with proper dependency injection

	// Tool: bootstrap_system
	s.AddTool(mcp.NewTool("bootstrap_system",
		mcp.WithDescription("Complete system initialization: installs MCPM, sets up default servers (context7, brave-search, github), and starts Docker infrastructure (PostgreSQL, Qdrant). One command to get fully operational."),
	), handleBootstrapSystem)

	// Tool: restart_service
	s.AddTool(mcp.NewTool("restart_service",
		mcp.WithDescription("Gracefully restarts Jarvis to apply configuration changes or resolve stuck states. Automatically saves state and reconnects active sessions. Use after editing server configs or when tools become unresponsive."),
	), handleRestartService)

	// Tool: restart_infrastructure
	s.AddTool(mcp.NewTool("restart_infrastructure",
		mcp.WithDescription("Safely reboots Docker infrastructure (PostgreSQL, Qdrant) with health checks and automatic reconnection. Resolves database connection issues, clears stale locks, and ensures all services are healthy. Zero data loss."),
	), handleRestartInfrastructure)

	// Tool: share_server (uses global sharedServers map)
	s.AddTool(mcp.NewTool("share_server",
		mcp.WithDescription("Exposes local MCP servers via secure tunnels with optional authentication. Enables remote teams to access your tools without VPN or port forwarding. Auto-generates shareable URLs with configurable access controls."),
		mcp.WithString("name",
			mcp.Description("Name of the server to share"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("Port to run the shared server on")),
		mcp.WithBoolean("no_auth",
			mcp.Description("Disable authentication for the shared server"),
		),
	), handleShareServer)

	// Tool: stop_sharing_server (uses global sharedServers map)
	s.AddTool(mcp.NewTool("stop_sharing_server",
		mcp.WithDescription("Revokes tunnel access and terminates shared server sessions. Immediately disconnects all remote clients. Changes are logged for security auditing."),
		mcp.WithString("name",
			mcp.Description("Name of the server to stop sharing"),
			mcp.Required(),
		),
	), handleStopSharingServer)

	// Tool: list_shared_servers (uses global sharedServers map)
	s.AddTool(mcp.NewTool("list_shared_servers",
		mcp.WithDescription("Shows all active server shares with tunnel URLs, authentication status, connected clients, and uptime. Useful for monitoring remote access and identifying security risks."),
	), handleListSharedServers)

	// Start the server based on transport mode
	if *httpMode {
		// Streamable HTTP mode (MCP 2025-03-26 spec)
		addr := fmt.Sprintf("%s:%s", *httpHost, *httpPort)
		log.Printf("Starting Jarvis in HTTP mode on %s", addr)
		fmt.Fprintf(os.Stderr, "\033[1;33m>> HTTP Mode: http://%s/mcp <<\033[0m\n", addr)

		httpServer := server.NewStreamableHTTPServer(s)

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Println("Shutting down HTTP server...")
			httpServer.Shutdown(nil)
		}()

		if err := httpServer.Start(addr); err != nil {
			log.Printf("HTTP server error: %v", err)
			fmt.Printf("Server error: %v\n", err)
		}
	} else {
		// Stdio mode (default for direct client connections)
		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}
}

// Helper function to run mcpm commands
func runMcpmCommand(args ...string) (string, error) {
	log.Printf("Executing MCPM command: %v", args)
	// mcpm is now available in PATH
	cmd := exec.Command("mcpm", args...)
	cmd.Env = append(os.Environ(), "MCPM_NON_INTERACTIVE=true", "MCPM_FORCE=true", "NO_COLOR=true")

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Strip common noise from MCPM output
	outputStr = stripMcpmNoise(outputStr)

	if err != nil {
		log.Printf("Command failed: %v. Output: %s", err, outputStr)
		return fmt.Sprintf("Error: %v\n\n%s", err, outputStr), fmt.Errorf("command failed: %v", err)
	}
	log.Printf("Command success. Output length: %d", len(output))

	return strings.TrimSpace(outputStr), nil
}

// stripMcpmNoise removes common warnings and noise from MCPM output
func stripMcpmNoise(output string) string {
	lines := strings.Split(output, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		// Skip warning lines
		if strings.Contains(line, "Warning: Input is not a terminal") {
			continue
		}
		if strings.Contains(line, "(fd=0)") && strings.Contains(line, "Warning:") {
			continue
		}
		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
