package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "lukekorsman.com/store/internal/http"
	"lukekorsman.com/store/internal/product"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(apphttp.RequestTimer)

	handler := product.NewHandler(store)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", handler.List)
		r.With(apphttp.SimpleAuth).Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
		r.With(apphttp.SimpleAuth).Put("/{id}", handler.Update)
		r.With(apphttp.SimpleAuth).Delete("/{id}", handler.Delete)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Listening on :8080")
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
