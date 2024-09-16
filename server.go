package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
)

func NewServer(
	ctx context.Context,
	repo ShortenedURLRepository,
	t Templater,
	l *slog.Logger,
	baseUrl *url.URL,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(ctx, repo, t, l, baseUrl, mux)

	loggingMiddleware := LoggerMiddleware(l)

	return loggingMiddleware(mux)
}
