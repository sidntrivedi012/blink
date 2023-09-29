package main

import (
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

	// Establish connection with Redis server.
	_ = initRedisClient()

	// Start the HTTP server.
	server := NewServer()
	server.AddRoutes()
	err := server.Start()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
