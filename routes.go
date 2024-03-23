package main

import (
	"context"
	"net/http"
	"net/url"
)

func addRoutes(
	ctx context.Context,
	repo ShortenedURLRepository,
	t Templater,
	baseUrl *url.URL,
	mux *http.ServeMux,
) {
	mux.Handle("POST /generate", handleGenerateURL(repo, t, baseUrl))
	mux.Handle("GET /u/{id}", handleRedirectURL(repo, t))
	mux.Handle("GET /", handleIndex(t))
}
