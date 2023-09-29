package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func shortenURL(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func getShortenedURLMetrics(c echo.Context) error {
	return nil
}

func routeShortenedURL(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}
