package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	blinkPort    string
	hostName     string
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
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})
	return &Server{Echo: e, Client: redisClient}
}

// Start starts the echo server.
func (s *Server) Start() error {
	blinkPort = os.Getenv("APP_PORT")
	if blinkPort == "" {
		blinkPort = "8080"
	}
	hostName = fmt.Sprintf("localhost:%s", blinkPort)
	err := s.Echo.Start(fmt.Sprintf(":%s", blinkPort))
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
