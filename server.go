package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

func NewServer(
	ctx context.Context,
	repo ShortenedURLRepository,
	t Templater,
	l *log.Logger,
	baseUrl *url.URL,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(ctx, repo, t, baseUrl, mux)

	loggingMiddleware := LoggerMiddleware(l)

	return loggingMiddleware(mux)
}
