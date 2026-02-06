package product

import (
	"testing"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		wantErrs int
	}{
		{
			name:     "validate product",
			product:  Product{Name: "Book", Price: 10.99},
			wantErrs: 0,
		},
		{
			name:     "empty name",
			product:  Product{Name: "", Price: 10},
			wantErrs: 1,
		},
		{
			name:     "zero price",
			product:  Product{Name: "Book", Price: 0},
			wantErrs: 1,
		},
		{
			name:     "multiple errors",
			product:  Product{Name: "", Price: -5},
			wantErrs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateProduct(tt.product)
			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
			}
		})
	}
}
