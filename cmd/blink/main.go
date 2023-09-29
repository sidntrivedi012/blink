package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
)

var Debug bool

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

	// Initialize the HTTP and Redis server.
	server := NewServer()
	server.ctx = context.Background()

	// Add routes and start the server.
	server.AddRoutes()
	err := server.Start()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
