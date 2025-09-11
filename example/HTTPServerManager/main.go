package main

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
	rtr "github.com/DilemaFixer/Cmd/router"
)

// Example 1: HTTPServerManager
// This example demonstrates basic command definition with different option types
// Usage examples:
//   myapp server start --host=localhost --port=8080 --debug
//   myapp server stop --force
//   myapp server status

func main() {
	// Simulate command line input - in real app you would use os.Args
	input := "server start --host=localhost --port=8080 --debug"

	// Parse the input
	parsedInput, err := p.ParseInput(input)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Create context and routing
	context := ctx.NewContext(parsedInput)
	iterator := rtr.NewRoutingIterator(context)
	router := rtr.NewRouter()

	// Define server management commands
	router.NewCmd("server").
		// Start server endpoint
		Endpoint("start").
		Description("Start the HTTP server").
		StringOption("host").   // Optional host (default will be used if not provided)
		IntOption("port").      // Optional port
		BoolOption("debug").    // Debug mode flag
		BoolOption("ssl").      // Enable SSL flag
		StringOption("config"). // Optional config file path
		Handler(startServerHandler).
		Build().

		// Stop server endpoint
		Endpoint("stop").
		Description("Stop the HTTP server").
		BoolOption("force").  // Force stop without graceful shutdown
		IntOption("timeout"). // Timeout in seconds for graceful shutdown
		Handler(stopServerHandler).
		Build().

		// Status endpoint
		Endpoint("status").
		Description("Show server status").
		BoolOption("detailed"). // Show detailed status
		Handler(statusHandler).
		Build().
		Register()

	// Execute the command
	router.Route(*context, iterator)
}

// Handler for server start command
func startServerHandler(ctx ctx.Context) error {
	fmt.Println("🚀 Starting HTTP Server...")

	// Get configuration values with defaults
	host := ctx.GetValueOrDefault("host", "localhost")
	port := ctx.GetValueOrDefault("port", "3000")

	// Check boolean flags
	debug := ctx.GetValueAsBool("debug")
	ssl := ctx.GetValueAsBool("ssl")

	// Get optional config file
	configFile := ctx.GetValueOrDefault("config", "server.conf")

	fmt.Printf("📍 Host: %s\n", host)
	fmt.Printf("🔌 Port: %s\n", port)
	fmt.Printf("🐛 Debug Mode: %v\n", debug)
	fmt.Printf("🔒 SSL Enabled: %v\n", ssl)
	fmt.Printf("⚙️  Config File: %s\n", configFile)

	if debug {
		fmt.Println("🔍 Debug logging enabled")
		fmt.Printf("📋 All flags: %v\n", ctx.GetFlagsAsMap())
	}

	fmt.Println("✅ Server started successfully!")
	return nil
}

// Handler for server stop command
func stopServerHandler(ctx ctx.Context) error {
	fmt.Println("🛑 Stopping HTTP Server...")

	force := ctx.GetValueAsBool("force")
	timeout := ctx.GetValueOrDefault("timeout", "30")

	if force {
		fmt.Println("⚡ Force stopping server (no graceful shutdown)")
	} else {
		fmt.Printf("⏰ Graceful shutdown with %s seconds timeout\n", timeout)
	}

	fmt.Println("✅ Server stopped successfully!")
	return nil
}

// Handler for server status command
func statusHandler(ctx ctx.Context) error {
	fmt.Println("📊 Server Status")
	fmt.Println("Status: Running")
	fmt.Println("Uptime: 2h 34m 12s")
	fmt.Println("Active Connections: 45")

	detailed := ctx.GetValueAsBool("detailed")
	if detailed {
		fmt.Println("\n📈 Detailed Status:")
		fmt.Println("  Memory Usage: 128MB")
		fmt.Println("  CPU Usage: 12%")
		fmt.Println("  Request Rate: 150/min")
		fmt.Println("  Error Rate: 0.02%")
	}

	return nil
}
