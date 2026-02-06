package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apphttp "lukekorsman.com/store/internal/http"
)

func TestListProducts_JSON(t *testing.T) {
	tests := []struct {
		name			string
		seedProducts	[]Product
		wantStatus		int
		wantCount		int
		wantFirstName 	string
	
	}{
		{
			name: 		"empty list",
			wantStatus: http.StatusOK,
			wantCount: 	0,
		},
		{
			name:			"single product",
			seedProducts: 	[]Product{
				{Name: "Book", Price: 10},
			},
			wantStatus:		http.StatusOK,
			wantCount: 		1,
			wantFirstName:	"Book",
		},
		{
			name:			"multiple products",
			seedProducts: 	[]Product{
				{Name: "Book", Price: 10},
				{Name: "Laptop", Price: 1200},
			},
			wantStatus:		http.StatusOK,
			wantCount: 		2,
			wantFirstName:	"Book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			for _, p := range tt.seedProducts {
				store.Create(context.Background(), p)
			}

			handler := NewHandler(store)

			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()

			handler.List(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("expected Content-Type application/json, got %s", ct)
			}

			var products []Product
			if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
				t.Fatalf("failed to decode JSON: %v", err)
			}

			if len(products) != tt.wantCount {
				t.Fatalf("expected %d products, got %d", tt.wantCount, len(products))
			}

			if tt.wantCount > 0 && products[0].Name != tt.wantFirstName {
				t.Fatalf("expected first product %s, got %s", tt.wantFirstName, products[0].Name)
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name		string
		apiKey		string
		body		string
		wantStatus	int
	} {
		{
			name: 		"unauthorized",
			apiKey:		"",
			body:		`{"name":"Book","Price":10}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: 		"invalid JSON",
			apiKey: 	"secret",
			body: 		`{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:		"valid request",
			apiKey: 	"secret",
			body: 		`{"name":"Book","Price":10}`,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			handler := NewHandler(store)

			protected := apphttp.SimpleAuth(
				http.HandlerFunc(handler.Create),
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/products",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			rec := httptest.NewRecorder()
			protected.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func BenchmarkListProducts_Sizes(b *testing.B) {
	sizes := []int{0, 10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			store := NewMemoryStore()
			for i := 0; i < size; i++ {
				store.Create(context.Background(), Product{Name: "Item", Price: 10})
			}

			handler := NewHandler(store)
			req := httptest.NewRequest(http.MethodGet, "/products", nil)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				rec := httptest.NewRecorder()
				handler.List(rec, req)
			}
		})
	}
}
