package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type contextKey string

const ValidatedBodyKey contextKey = "validatedBody"

var validate = validator.New()

// FormatValidationError translates technical tags into readable messages
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	// Type assertion to extract structural field errors
	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "Validation failed"}
	}

	for _, f := range valErrors {
		// Convert field name to lowercase to match typical JSON keys
		field := strings.ToLower(f.Field())

		switch f.Tag() {
		case "required":
			errors[field] = "This field is required"
		case "email":
			errors[field] = "Must be a valid email address"
		case "min":
			errors[field] = "Must be at least " + f.Param() + " characters long"
		case "gte":
			errors[field] = "Must be greater than or equal to " + f.Param()
		case "datetime":
			errors[field] = "Must be a valid date and time"
		case "boolean":
			errors[field] = "Must be a valid boolean value"
		default:
			errors[field] = "Invalid value (failed " + f.Tag() + " check)"
		}
	}
	return errors
}

// ValidateBody middleware returning structured JSON errors
func ValidateBody[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload T

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Malformed or invalid JSON payload"})
				return
			}

			if err := validate.Struct(payload); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnprocessableEntity)

				// Return the clean map to the client
				json.NewEncoder(w).Encode(map[string]any{
					"errors": FormatValidationError(err),
				})
				return
			}

			ctx := context.WithValue(r.Context(), ValidatedBodyKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
