package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type LongURL struct {
	URL string `json:"long_url"`
}

// shortenURL is the handler that takes care of the URL shortening API.
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

// redirectShortenedURL is the handler that redirects the shortened URL to the long form URL if
// an entry for it is stored in the database.
func (s *Server) redirectShortenedURL(c echo.Context) error {
	urlHash := strings.TrimLeft(c.Request().RequestURI, "/")

	// Check that if a long URL corresponding to this short URL exists in the
	// database. If yes, redirect to that.
	exists, longURL, err := getLongURLFromCache(s.ctx, s.Client, urlHash)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return c.String(http.StatusNotFound, "found no entry of this short url")
	}
	return c.Redirect(http.StatusTemporaryRedirect, longURL)
}

// getShortenedURLMetrics is the handler that returns the top three most shortened URL hostnames in the response
// along with the count.
func (s *Server) getShortenedURLMetrics(c echo.Context) error {
	topDomains, err := fetchMostShortenedURLs(s.ctx, s.Client)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, topDomains)
}
