package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

const (
	blinkPort    int    = 8080
	redisPort    int    = 6379
	hostName     string = "localhost:8080"
	serverScheme string = "http"
)

type Server struct {
	*redis.Client
	*echo.Echo
	ctx context.Context
}

// NewServer initializes a new instance of Echo server.
func NewServer() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", redisPort),
		Password: "",
		DB:       0,
	})
	return &Server{Echo: e, Client: redisClient}
}

// Start starts the echo server.
func (s *Server) Start() error {
	err := s.Echo.Start(fmt.Sprintf(":%d", blinkPort))
	if err != nil {
		slog.Error("failed to start server", slog.Any("err", err))
		return err
	}
	return nil
}

// AddRoutes adds routes to the server which need to be served along with
// the respective handler function.
func (s *Server) AddRoutes() {
	s.Echo.POST("/api/shorten", s.shortenURL)
	s.Echo.GET("/api/metrics", s.getShortenedURLMetrics)
	s.Echo.GET("/*", s.redirectShortenedURL)
}
