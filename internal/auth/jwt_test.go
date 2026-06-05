package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	role := "client"
	expirationHours := 24

	token, err := GenerateToken(userID, email, role, secret, expirationHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}
}

func TestGenerateTokenZeroExpiration(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	role := "client"
	expirationHours := 0

	token, err := GenerateToken(userID, email, role, secret, expirationHours)
	if err != nil {
		t.Fatalf("GenerateToken with 0 expiration failed: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	role := "admin"
	expirationHours := 24

	token, err := GenerateToken(userID, email, role, secret, expirationHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("Email mismatch: got %s, want %s", claims.Email, email)
	}

	if claims.Role != role {
		t.Errorf("Role mismatch: got %s, want %s", claims.Role, role)
	}
}

func TestValidateTokenInvalidSecret(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	role := "client"
	expirationHours := 24

	token, err := GenerateToken(userID, email, role, secret, expirationHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	wrongSecret := "wrong-secret-key"
	_, err = ValidateToken(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error with wrong secret")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	role := "client"
	expirationHours := -1 // Expired 1 hour ago

	token, err := GenerateToken(userID, email, role, secret, expirationHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	time.Sleep(1 * time.Second) // Ensure token is expired

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func TestValidateTokenInvalidFormat(t *testing.T) {
	secret := "test-secret-key"
	invalidToken := "invalid.token.format"

	_, err := ValidateToken(invalidToken, secret)
	if err == nil {
		t.Fatal("Expected error for invalid token format")
	}
}
