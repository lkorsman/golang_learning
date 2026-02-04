package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleAuth(t *testing.T) {
	tests := []struct {
		name		string
		apiKey		string
		wantStatus 	int
		wantCalled bool
	}{
		{
			name:		"missing API key",
			apiKey: 	"", 
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: 		"invalid API key",
			apiKey: 	"wrong",
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: 		"valid API key",
			apiKey:		"secret",
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			handler := SimpleAuth(next)

			req := httptest.NewRequest(http.MethodPost, "/products", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, rec.Code)
			}

			if called != tt.wantCalled {
				t.Fatalf("expected called=%v, got %v", tt.wantCalled, called)
			}
		})
	}
}