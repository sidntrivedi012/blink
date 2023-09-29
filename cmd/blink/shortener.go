package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
