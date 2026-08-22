package auth_test

import (
	"errors"
	errors2 "marketplace/internal/error"
	"marketplace/internal/auth"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	code := m.Run()
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func TestGenerateParseRoundtrip(t *testing.T) {
	token, err := auth.GenerateToken(42, "customer")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, ожидали 42", claims.UserID)
	}
	if claims.Role != "customer" {
		t.Errorf("Role = %q, ожидали \"customer\"", claims.Role)
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("токен должен иметь срок жизни в будущем")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := auth.ParseToken("not-a-jwt"); err == nil {
		t.Error("ожидали ошибку для битого токена")
	}
}

func TestParseTokenExpired(t *testing.T) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	if _, err := auth.ParseToken(signed); !errors.Is(err, errors2.ErrTokenExpired) {
		t.Errorf("ожидали ErrTokenExpired, получили %v", err)
	}
}

func TestParseTokenWrongKey(t *testing.T) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("other-key"))
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	if _, err := auth.ParseToken(signed); err == nil {
		t.Error("ожидали ошибку для токена с чужой подписью")
	}
}
