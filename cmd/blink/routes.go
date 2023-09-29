package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LongURL struct {
	URL string `json:"long_url"`
}

func (s *Server) shortenURL(c echo.Context) error {
	var longURL LongURL
	if err := c.Bind(&longURL); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	shortenedURL, err := shortenLongURL(s.ctx, longURL.URL, s.Client)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, shortenedURL)
}

func (s *Server) getShortenedURLMetrics(c echo.Context) error {
	return nil
}

func (s *Server) routeShortenedURL(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}
