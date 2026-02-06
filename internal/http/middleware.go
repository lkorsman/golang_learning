package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func RequestTimer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("%s %s took %v\n", r.Method, r.URL.Path, duration)
	})
}

func SimpleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user := User{
			ID:   1,
			Name: "Alice",
		}

		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(ww, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.status).
				Dur("duration_ms", time.Since(start)).
				Msg("request completed")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
