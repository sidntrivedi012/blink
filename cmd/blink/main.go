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

// ParseFlags parses the command line arguments to blink.
func ParseFlags() {
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	Debug = *debug
}

// SetLogLevel configures the logging behaviour for the application.
func SetLogLevel() {
	if Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
}

func main() {
	SetLogLevel()
	ParseFlags()
	defer func() {
		os.Exit(exitCode)
	}()

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
