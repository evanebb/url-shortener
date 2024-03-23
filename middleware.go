package main

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (rw *responseWriterWrapper) Status() int {
	return rw.status
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// LoggerMiddleware will log information about the HTTP request that was made.
func LoggerMiddleware(l *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := newResponseWriterWrapper(w)
			next.ServeHTTP(wrapped, r)
			l.Printf("%s %d %s %s %s", r.RemoteAddr, wrapped.status, r.Method, r.URL.EscapedPath(), time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}
