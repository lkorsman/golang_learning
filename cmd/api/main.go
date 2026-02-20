package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lukekorsman.com/store/internal/auth"
	"lukekorsman.com/store/internal/cache"
	"lukekorsman.com/store/internal/config"
	apphttp "lukekorsman.com/store/internal/http"
	"lukekorsman.com/store/internal/product"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(apphttp.RequestTimer)
	r.Use(apphttp.MetricsMiddleware)

	var store product.Store
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		mysqlStore, err := product.NewMySQLStore(dbURL)
		if err != nil {
			panic(err)
		}
		defer mysqlStore.Close()
		store = mysqlStore
	} else {
		store = product.NewMemoryStore()
		fmt.Println("Using in-memory store")
	}

    redisCache, err := cache.NewRedisCache(cfg.RedisURL)
    if err != nil {
        fmt.Printf("Redis unavailable, running without cache: %v\n", err)
    }
    if redisCache != nil {
        defer redisCache.Close()
    }

	userStore := auth.NewMemoryUserStore()
    jwtManager := auth.NewJWTManager(cfg.JWTSecret, "store-api")
    authHandler := auth.NewHandler(userStore, jwtManager)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/auth", func(r chi.Router) {
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
    })

	productHandler := product.NewHandler(store, redisCache)
	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.Get)

		r.Group(func(r chi.Router) {
			r.Use(apphttp.JWTAuth(jwtManager, userStore))
			r.Post("/", productHandler.Create)
			r.Put("/{id}", productHandler.Update)
			r.Delete("/{id}", productHandler.Delete)
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		fmt.Println("Listening on :" + cfg.Port)
		fmt.Println("Metrics available at http://localhost:" + cfg.Port + "/metrics")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("shutdown error:", err)
	}

	fmt.Println("Server gracefully stopped")
}
