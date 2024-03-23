package main

import (
	"context"
	"errors"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io/fs"
	"log"
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

	boltDb, err := bolt.Open(c.dbFile, 0600, nil)
	if err != nil {
		return err
	}

	repo, err := newBoltShortenedURLRepository(boltDb)
	if err != nil {
		return err
	}

	logger := log.New(os.Stdout, "http: ", log.Ltime)

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
			log.Printf("listening on %s (HTTPS enabled)\n", httpServer.Addr)
			if err := httpServer.ListenAndServeTLS(c.httpsCertFile, c.httpsKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				_, _ = fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
			}
		} else {
			log.Printf("listening on %s\n", httpServer.Addr)
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				_, _ = fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
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
			_, _ = fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
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
