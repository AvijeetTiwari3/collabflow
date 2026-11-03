package unit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/auth"
)

func TestTokenGenerationAndValidation(t *testing.T) {
	mgr := auth.NewTokenManager("test-jwt-secret-key-32-bytes-long!", 15, 7)
	userID := uuid.New()
	email := "alex@google.com"
	name := "Alex SDE"

	// 1. Generate access token
	token, err := mgr.GenerateAccessToken(userID, email, name)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// 2. Validate token
	claims, err := mgr.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Name != name {
		t.Errorf("expected name %s, got %s", name, claims.Name)
	}

	// 3. Password hashing
	pw := "SecretP@ssword2026!"
	hash, err := auth.HashPassword(pw)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !auth.CheckPasswordHash(pw, hash) {
		t.Fatalf("password check failed for valid password")
	}
	if auth.CheckPasswordHash("wrong-password", hash) {
		t.Fatalf("password check succeeded for invalid password")
	}
}
