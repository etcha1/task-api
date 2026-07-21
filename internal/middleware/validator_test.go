package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testPayload struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func TestValidateBodyRejectsInvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"A","email":"not-an-email"}`))
	rec := httptest.NewRecorder()

	handler := ValidateBody[testPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called for invalid payload")
	}))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	errors, ok := body["errors"].(map[string]any)
	if !ok {
		t.Fatalf("response body missing errors field: %#v", body)
	}

	if errors["name"] != "Must be at least 2 characters long" {
		t.Fatalf("unexpected name error: %#v", errors["name"])
	}

	if errors["email"] != "Must be a valid email address" {
		t.Fatalf("unexpected email error: %#v", errors["email"])
	}
}

func TestValidateBodyStoresValidatedPayloadInContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Alice","email":"alice@example.com"}`))
	rec := httptest.NewRecorder()

	var got testPayload
	handler := ValidateBody[testPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(ValidatedBodyKey).(testPayload)
		if !ok {
			t.Fatal("validated payload was not stored in context")
		}
		got = payload
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if got.Name != "Alice" || got.Email != "alice@example.com" {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestValidateBodyRejectsMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":`))
	rec := httptest.NewRecorder()

	handler := ValidateBody[testPayload]()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for malformed JSON")
	}))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["error"] != "Malformed or invalid JSON payload" {
		t.Fatalf("unexpected error message: %#v", body)
	}
}

func TestFormatValidationErrorReturnsReadableMessages(t *testing.T) {
	var payload struct {
		Name string `validate:"required"`
		Age  int    `validate:"gte=18"`
	}

	err := validate.Struct(payload)
	if err == nil {
		t.Fatal("expected validation error")
	}

	got := FormatValidationError(err)
	if got["name"] != "This field is required" {
		t.Fatalf("unexpected name error: %#v", got["name"])
	}

	if got["age"] != "Must be greater than or equal to 18" {
		t.Fatalf("unexpected age error: %#v", got["age"])
	}
}
