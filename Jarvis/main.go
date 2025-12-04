package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"jarvis/handlers"
)

// Command-line flags for transport mode
var (
	httpMode = flag.Bool("http", false, "Run as Streamable HTTP server instead of stdio")
	httpPort = flag.String("port", "6275", "Port for HTTP server (default: 6275)")
	httpHost = flag.String("host", "127.0.0.1", "Host for HTTP server (default: 127.0.0.1)")
)

var logFile *os.File

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
