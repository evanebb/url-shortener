package main

import (
	"context"
	"errors"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Set up the dependencies (database, logger, etc.)
	c, err := getAppConfiguration()
	if err != nil {
		return err
	}

	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(c.logLevel)); err != nil {
		return err
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	boltDb, err := bolt.Open(c.dbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer func(boltDb *bolt.DB) {
		if err := boltDb.Close(); err != nil {
			logger.Error("error closing BoltDB database", "error", err)
		}
	}(boltDb)

	repo, err := newBoltShortenedURLRepository(boltDb)
	if err != nil {
		return err
	}

	templateFs, err := fs.Sub(resources, "resources")
	if err != nil {
		return err
	}

	t := NewTemplater(templateFs)

	srv := NewServer(ctx, repo, t, logger, c.baseURL)

	httpServer := &http.Server{
		Addr:    c.listenAddress,
		Handler: srv,
	}

	go func() {
		if c.httpsEnabled {
			logger.Info(fmt.Sprintf("listening on %s (HTTPS enabled)", httpServer.Addr))
			if err := httpServer.ListenAndServeTLS(c.httpsCertFile, c.httpsKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("error listening and serving", "error", err)
			}
		} else {
			logger.Info(fmt.Sprintf("listening on %s", httpServer.Addr))
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("error listening and serving", "error", err)
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", "error", err)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
