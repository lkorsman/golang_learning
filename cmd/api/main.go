package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lukekorsman.com/store/internal/product"
	apphttp "lukekorsman.com/store/internal/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(apphttp.RequestTimer)

	store := product.NewMemoryStore()
	handler := product.NewHandler(store)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", handler.List)
		r.With(apphttp.SimpleAuth).Post("/", handler.Create)
	})

	srv := &http.Server{
		Addr:	":8080",
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