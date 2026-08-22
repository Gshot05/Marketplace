package services_test

import (
	"context"
	errors2 "marketplace/internal/error"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/workerpool"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	code := m.Run()
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func newAuthService(repo *fakeAuthRepo, notifier *fakeNotifier) *service.AuthService {
	return service.NewAuthService(repo, notifier, workerpool.New(1, 8))
}

func TestRegisterUserSuccess(t *testing.T) {
	repo := &fakeAuthRepo{}
	notifier := newFakeNotifier()
	s := newAuthService(repo, notifier)

	err := s.RegisterUser(context.Background(), "a@b.com", "pass123", "customer", "Islam")
	if err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	if repo.savedEmail != "a@b.com" {
		t.Errorf("код сохранён для %q, ожидали a@b.com", repo.savedEmail)
	}
	if len(repo.savedCode) != 6 || !isDigits(repo.savedCode) {
		t.Errorf("код %q должен быть строкой из 6 цифр", repo.savedCode)
	}

	select {
	case code := <-notifier.verifyCh:
		if code != repo.savedCode {
			t.Errorf("отправленный код %q не совпадает с сохранённым %q", code, repo.savedCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("письмо с кодом не отправлено через worker pool")
	}
}

func TestRegisterUserInvalidInputSkipsRepo(t *testing.T) {
	repo := &fakeAuthRepo{}
	s := newAuthService(repo, newFakeNotifier())

	if err := s.RegisterUser(context.Background(), "bad-email", "pass", "customer", "Islam"); err == nil {
		t.Error("ожидали ошибку для невалидного email")
	}
	if err := s.RegisterUser(context.Background(), "a@b.com", "pass", "customer", ""); err != errors2.ErrEmptyName {
		t.Errorf("ожидали ErrEmptyName, получили %v", err)
	}
	if err := s.RegisterUser(context.Background(), "a@b.com", "pass", "", "Islam"); err != errors2.ErrEmptyRole {
		t.Errorf("ожидали ErrEmptyRole, получили %v", err)
	}

	if repo.savedEmail != "" {
		t.Error("репозиторий не должен вызываться при невалидных данных")
	}
}

func TestConfirmEmailSuccess(t *testing.T) {
	var mu sync.Mutex
	repo := &fakeAuthRepo{registerID: 42}
	notifier := newFakeNotifier()
	s := newAuthService(repo, notifier)

	ctx := context.Background()
	if err := s.RegisterUser(ctx, "confirm@test.com", "secret", "performer", "Islam"); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}
	mu.Lock()
	code := repo.savedCode
	mu.Unlock()

	token, err := s.ConfirmEmail(ctx, "confirm@test.com", code)
	if err != nil {
		t.Fatalf("ConfirmEmail: %v", err)
	}
	if token == "" {
		t.Fatal("ожидали JWT после подтверждения")
	}
	if strings.Count(token, ".") != 2 {
		t.Errorf("строка %q не похожа на JWT", token)
	}
}

func TestConfirmEmailWrongCode(t *testing.T) {
	repo := &fakeAuthRepo{verifyErr: errors2.ErrWrongVerify}
	s := newAuthService(repo, newFakeNotifier())

	_, err := s.ConfirmEmail(context.Background(), "x@y.com", "000000")
	if err != errors2.ErrWrongVerify {
		t.Errorf("ожидали ErrWrongVerify, получили %v", err)
	}
}

func TestConfirmEmailUnknownEmail(t *testing.T) {
	repo := &fakeAuthRepo{}
	s := newAuthService(repo, newFakeNotifier())

	_, err := s.ConfirmEmail(context.Background(), "unknown@test.com", "123456")
	if err != errors2.ErrWrongConfirmData {
		t.Errorf("ожидали ErrWrongConfirmData, получили %v", err)
	}
}

func TestLoginUserSuccess(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("right-pass"), bcrypt.DefaultCost)
	repo := &fakeAuthRepo{
		user: &model.User{ID: 7, Email: "login@test.com", PasswordHash: string(hash), Role: "customer"},
	}
	notifier := newFakeNotifier()
	s := newAuthService(repo, notifier)

	token, err := s.LoginUser(context.Background(), "login@test.com", "right-pass")
	if err != nil {
		t.Fatalf("LoginUser: %v", err)
	}
	if strings.Count(token, ".") != 2 {
		t.Errorf("строка %q не похожа на JWT", token)
	}

	select {
	case to := <-notifier.loginCh:
		if to != "login@test.com" {
			t.Errorf("уведомление ушло на %q", to)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("уведомление о входе не отправлено через worker pool")
	}
}

func TestLoginUserWrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("right-pass"), bcrypt.DefaultCost)
	repo := &fakeAuthRepo{
		user: &model.User{ID: 7, Email: "login@test.com", PasswordHash: string(hash), Role: "customer"},
	}
	s := newAuthService(repo, newFakeNotifier())

	if _, err := s.LoginUser(context.Background(), "login@test.com", "wrong-pass"); err != errors2.ErrWrongPassOrLog {
		t.Errorf("ожидали ErrWrongPassOrLog, получили %v", err)
	}
}

func TestLoginUserUnknownEmail(t *testing.T) {
	repo := &fakeAuthRepo{getUserErr: errors2.ErrWrongPassOrLog}
	s := newAuthService(repo, newFakeNotifier())

	if _, err := s.LoginUser(context.Background(), "ghost@test.com", "any"); err != errors2.ErrWrongPassOrLog {
		t.Errorf("ожидали ErrWrongPassOrLog, получили %v", err)
	}
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
