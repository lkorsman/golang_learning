package product

import (
	"strings"
)

type ValidationError struct {
	Field string
	Message string
}

func (v ValidationError) Error() string {
	return v.Field + ": " + v.Message
}

func ValidateProduct(p Product) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(p.Name) == "" {
		errs = append(errs, ValidationError{
			Field: "name",
			Message: "name is require",
		})
	}

	if len(p.Name) > 100 {
		errs = append(errs, ValidationError{
			Field: "name",
			Message: "name must be less than 100 characters",
		})
	}

	if p.Price <= 0 {
		errs = append(errs, ValidationError{
			Field: "price",
			Message: "price must be greater than 0",
		})
	}

	if p.Price > 999999.99 {
		errs = append(errs, ValidationError{
			Field: "price",
			Message: "price is too large",
		})
	}

	return errs
}