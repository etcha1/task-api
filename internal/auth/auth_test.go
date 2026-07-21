package auth

import (
	"testing"
)

func TestInitializeWithSecret(t *testing.T) {
	InitializeWithSecret("test-secret")

	if TokenAuth == nil {
		t.Fatal("TokenAuth should be initialized")
	}
}

func TestNewToken(t *testing.T) {
	InitializeWithSecret("test-secret")

	token, err := NewToken(42)
	if err != nil {
		t.Fatalf("NewToken returned error: %v", err)
	}

	if token == "" {
		t.Fatal("NewToken should return a non-empty token")
	}
}

func TestGetSecretUsesEnvironmentVariable(t *testing.T) {
	t.Setenv("JWT_SECRET", "env-secret")

	if got := getSecret(); got != "env-secret" {
		t.Fatalf("getSecret() = %q, want %q", got, "env-secret")
	}
}
