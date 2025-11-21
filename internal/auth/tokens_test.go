package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken("user-1", "test@example.com", []byte("secret"), time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	claims, err := ParseToken(token, []byte("secret"))
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != "user-1" || claims.Email != "test@example.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, err := GenerateToken("user-1", "test@example.com", []byte("secret"), -1*time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if _, err := ParseToken(token, []byte("secret")); err == nil {
		t.Fatalf("expected error for expired token")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("user-1", "test@example.com", []byte("secret"), time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if _, err := ParseToken(token, []byte("other")); err == nil {
		t.Fatalf("expected error for wrong secret")
	}
}
