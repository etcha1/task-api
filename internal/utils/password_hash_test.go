package utils

import "testing"

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "secret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("CheckPasswordHash should return true for the matching password")
	}

	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash should return false for a mismatched password")
	}
}
