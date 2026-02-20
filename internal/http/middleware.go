package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lukekorsman.com/store/internal/auth"
	"lukekorsman.com/store/internal/metrics"

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

func JWTAuth(jwtManager *auth.JWTManager, userStore auth.UserStore) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }
            
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            claims, err := jwtManager.Verify(tokenString)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
   
            user, err := userStore.GetByID(r.Context(), claims.UserID)
            if err != nil {
                http.Error(w, "user not found", http.StatusUnauthorized)
                return
            }
   
            ctx := auth.ContextWithUser(r.Context(), user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
        
        next.ServeHTTP(ww, r)
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(ww.status)
        
        // Record metrics
        metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
    })
}