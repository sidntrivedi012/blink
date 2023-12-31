package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
)

var (
	Debug    bool
	exitCode int = 0
)

// parseFlags parses the command line arguments to blink.
func parseFlags() {
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	Debug = *debug
}

// setLogLevel configures the logging behaviour for the application.
func setLogLevel() {
	if Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
}

func main() {
	defer func() {
		os.Exit(exitCode)
	}()
	setLogLevel()
	parseFlags()

	// Initialize the HTTP and Redis server.
	server := NewServer()
	server.ctx = context.Background()
	defer server.Client.Close()

	// Add routes and start the server.
	server.AddRoutes()
	err := server.Start()
	if err != nil {
		slog.Error(err.Error())
		exitCode = 1
		return
	}
}
