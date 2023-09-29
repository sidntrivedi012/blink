package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
)

const BlinkPort int = 8080

type Server struct {
	*echo.Echo
}

// NewServer initializes a new instance of Echo server.
func NewServer() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return &Server{Echo: e}
}

// Start starts the echo server.
func (s *Server) Start() error {
	err := s.Echo.Start(fmt.Sprintf(":%d", BlinkPort))
	if err != nil {
		slog.Error("failed to start server", slog.Any("err", err))
		return err
	}
	return nil
}

// AddRoutes adds routes to the server which need to be served along with
// the respective handler function.
func (s *Server) AddRoutes() {
	s.Echo.POST("/shorten", shortenURL)
	s.Echo.GET("/metrics", getShortenedURLMetrics)
	s.Echo.GET("/*", routeShortenedURL)
}
