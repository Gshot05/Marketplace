package handlers_test

import (
	"context"
	errors2 "marketplace/internal/error"
	"marketplace/internal/handlers"
	"marketplace/internal/logger"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/service"
	"marketplace/internal/workerpool"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

type silentNotifier struct{}

func (silentNotifier) SendVerificationCode(ctx context.Context, to, code string) error {
	return nil
}

func (silentNotifier) SendLoginNotification(ctx context.Context, to string) error {
	return nil
}

type memAuthRepo struct {
	mu    sync.Mutex
	codes map[string]string
}

func newMemAuthRepo() *memAuthRepo {
	return &memAuthRepo{codes: make(map[string]string)}
}

func (m *memAuthRepo) RegisterUser(ctx context.Context, email, password, role, name string) (uint, error) {
	return 1, nil
}

func (m *memAuthRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	return &model.User{ID: 1, Email: email, PasswordHash: string(hash), Role: "customer"}, nil
}

func (m *memAuthRepo) VerifyCode(ctx context.Context, email, code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	saved, ok := m.codes[email]
	if !ok || saved != code {
		return errors2.ErrWrongVerify
	}
	return nil
}

func (m *memAuthRepo) DeleteUsedCode(ctx context.Context, email string) error {
	return nil
}

func (m *memAuthRepo) SaveVerificationCode(ctx context.Context, email, code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.codes[email] = code
	return nil
}

var _ repository.IAuthRepo = (*memAuthRepo)(nil)

func buildAuthRouter(repo *memAuthRepo) *gin.Engine {
	r := gin.New()
	authService := service.NewAuthService(repo, silentNotifier{}, workerpool.New(1, 8))
	authHandler := handlers.NewAuthHandler(authService, logger.NewLogger(nil, workerpool.New(1, 8)))

	g := r.Group("/auth")
	g.POST("/register", authHandler.Register())
	g.POST("/login", authHandler.Login())
	g.POST("/confirm", authHandler.ConfirmEmail())
	return r
}

func TestRegisterAndConfirmFlow(t *testing.T) {
	repo := newMemAuthRepo()
	r := buildAuthRouter(repo)

	w := doJSON(t, r, http.MethodPost, "/auth/register", `{"email":"flow@test.com","password":"pass123","role":"customer","name":"Islam"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("register: код %d, body=%s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodPost, "/auth/confirm", `{"email":"flow@test.com","verifyCode":"000000"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("confirm с неверным кодом: код %d, ожидали 400; body=%s", w.Code, w.Body.String())
	}

	code := repo.codes["flow@test.com"]
	if len(code) != 6 {
		t.Fatalf("сохранённый код %q должен быть из 6 цифр", code)
	}

	body := `{"email":"flow@test.com","verifyCode":"` + code + `"}`
	w = doJSON(t, r, http.MethodPost, "/auth/confirm", body)
	if w.Code != http.StatusOK {
		t.Errorf("confirm с верным кодом: код %d, body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"token"`) {
		t.Errorf("ожидали JWT в ответе: %s", w.Body.String())
	}
}

func TestLoginSuccessAndWrongCreds(t *testing.T) {
	r := buildAuthRouter(newMemAuthRepo())

	w := doJSON(t, r, http.MethodPost, "/auth/login", `{"email":"user@test.com","password":"pass123"}`)
	if w.Code != http.StatusOK {
		t.Errorf("корректный логин: код %d, body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"token"`) {
		t.Errorf("ожидали JWT в ответе: %s", w.Body.String())
	}

	w = doJSON(t, r, http.MethodPost, "/auth/login", `{"email":"user@test.com","password":"wrong"}`)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("неверный пароль: код %d, ожидали 401", w.Code)
	}
}

func TestRegisterBadJSON(t *testing.T) {
	r := buildAuthRouter(newMemAuthRepo())

	w := doJSON(t, r, http.MethodPost, "/auth/register", `{invalid`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("код %d, ожидали 400", w.Code)
	}
	if !strings.Contains(w.Body.String(), "error") {
		t.Errorf("ожидали поле error в ответе: %s", w.Body.String())
	}
}
