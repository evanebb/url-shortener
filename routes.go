package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
)

func addRoutes(
	ctx context.Context,
	repo ShortenedURLRepository,
	t Templater,
	l *slog.Logger,
	baseUrl *url.URL,
	mux *http.ServeMux,
) {
	mux.Handle("POST /generate", handleGenerateURL(l, repo, t, baseUrl))
	mux.Handle("GET /u/{id}", handleRedirectURL(l, repo, t))
	mux.Handle("GET /", handleIndex(t))
}
