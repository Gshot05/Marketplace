package middleware_test

import (
	"marketplace/internal/auth"
	"marketplace/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func newAuthRouter() (*gin.Engine, *httptest.ResponseRecorder) {
	r := gin.New()
	r.GET("/protected", middleware.AuthMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r, httptest.NewRecorder()
}

func TestAuthMiddlewareNoHeader(t *testing.T) {
	r, _ := newAuthRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("без заголовка: код %d, ожидали 401", w.Code)
	}
}

func TestAuthMiddlewareBadPrefix(t *testing.T) {
	r, _ := newAuthRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("без Bearer: код %d, ожидали 401", w.Code)
	}
}

func TestAuthMiddlewareGarbageToken(t *testing.T) {
	r, _ := newAuthRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer garbage.token.here")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("битый токен: код %d, ожидали 401", w.Code)
	}
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	token, err := auth.GenerateToken(42, "customer")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	r, _ := newAuthRouter()

	var gotUserID uint
	var gotRole string
	var reached bool
	r.GET("/check", middleware.AuthMiddleware, func(c *gin.Context) {
		reached = true
		gotUserID = c.GetUint("user_id")
		gotRole = c.GetString("role")
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("валидный токен: код %d, ожидали 200; body=%s", w.Code, w.Body.String())
	}
	if !reached {
		t.Fatal("финальный хендлер не достигнут")
	}
	if gotUserID != 42 || gotRole != "customer" {
		t.Errorf("контекст: user_id=%d role=%q, ожидали 42/customer", gotUserID, gotRole)
	}
}

func TestAuthMiddlewareExpiredToken(t *testing.T) {
	r, _ := newAuthRouter()

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("просроченный токен: код %d, ожидали 401", w.Code)
	}
}
