package httpx

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorDetail struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ValidationErrorResponse struct {
	Message string                  `json:"message"`
	Errors  []ValidationErrorDetail `json:"errors"`
}

func ParseValidationErrors(err error) ValidationErrorResponse {
	var validationErrors []ValidationErrorDetail

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			validationErrors = append(validationErrors, ValidationErrorDetail{
				Field: e.Field(),
				Tag:   e.Tag(),
				Value: e.Param(),
			})
		}
	}
	return ValidationErrorResponse{
		Message: "validation errors",
		Errors:  validationErrors,
	}
}
