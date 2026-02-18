package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"lukekorsman.com/store/internal/auth"
	"lukekorsman.com/store/internal/cache"
)

type Handler struct {
	store Store
	cache *cache.RedisCache
}

func NewHandler(store Store, redisCache *cache.RedisCache) *Handler {
	return &Handler{
		store: store,
		cache: redisCache,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cacheKey := "products:list"

	var products []Product

	if h.cache != nil {
		err := h.cache.Get(ctx, cacheKey, &products)

		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			writeJSON(w, http.StatusOK, products)
			return
		}
	}

	products, err := h.store.List(ctx)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	if h.cache != nil {
		if err := h.cache.Set(ctx, cacheKey, products, 5*time.Minute); err != nil {
			fmt.Printf("Failed to cache products: %v\n", err)
		}
		w.Header().Set("X-Cache", "MISS")
	} else {
		w.Header().Set("X-Cache", "DISABLED")
	}

    writeJSON(w, http.StatusOK, products)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
    if !ok {
        http.Error(w, "user not found", http.StatusUnauthorized)
        return
    }

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if errs := ValidateProduct(p); len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": errs,
		})
		return
	}

	fmt.Printf("User %s is creating a product\n", user.Email)
	created, err := h.store.Create(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.cache != nil {
		if err := h.cache.Delete(r.Context(), "products:list"); err != nil {
        	fmt.Printf("Failed to invalidate cache: %v\n", err)
    	}
	}
	
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("product:%d", id)

	var product Product
	err = h.cache.Get(ctx, cacheKey, &product)

	if err == nil {
		w.Header().Set("X-Cache", "HIT")
		writeJSON(w, http.StatusOK, product)
		return
	}

	product, err = h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.cache.Set(ctx, cacheKey, product, 10*time.Minute); err != nil {
        fmt.Printf("Failed to cache product: %v\n", err)
    }
    
    w.Header().Set("X-Cache", "MISS")
    writeJSON(w, http.StatusOK, product)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if errs := ValidateProduct(p); len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": errs,
		})
		return
	}

	updated, err := h.store.Update(r.Context(), id, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cacheKeys := []string{
        "products:list",
        fmt.Sprintf("product:%d", id),
    }

    if err := h.cache.Delete(r.Context(), cacheKeys...); err != nil {
        fmt.Printf("Failed to invalidate cache: %v\n", err)
    }

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cacheKeys := []string{
        "products:list",
        fmt.Sprintf("product:%d", id),
    }
    if err := h.cache.Delete(r.Context(), cacheKeys...); err != nil {
        fmt.Printf("Failed to invalidate cache: %v\n", err)
    }

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
