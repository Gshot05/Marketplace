package utils_test

import (
	errors2 "marketplace/internal/error"
	"marketplace/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestValidateIncomingRegistration(t *testing.T) {
	cases := []struct {
		name    string
		email   string
		user    string
		role    string
		wantErr error
	}{
		{"валидные данные", "a@b.com", "Islam", "customer", nil},
		{"невалидный email", "not-an-email", "Islam", "performer", errors2.ErrBadToken},
		{"пустое имя", "a@b.com", "", "customer", errors2.ErrEmptyName},
		{"имя из пробела", "a@b.com", " ", "customer", errors2.ErrEmptyName},
		{"пустая роль", "a@b.com", "Islam", "", errors2.ErrEmptyRole},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateIncomingRegistration(tc.email, tc.user, tc.role)
			if tc.name == "невалидный email" {
				if err == nil || err == tc.wantErr {
					t.Errorf("ожидали ошибку формата email, получили: %v", err)
				}
				return
			}
			if err != tc.wantErr {
				t.Errorf("ожидали %v, получили %v", tc.wantErr, err)
			}
		})
	}
}

func TestValidateBearerToken(t *testing.T) {
	cases := []struct {
		name    string
		header  string
		want    string
		wantErr error
	}{
		{"корректный заголовок", "Bearer abc123", "abc123", nil},
		{"пустой заголовок", "", "", errors2.ErrNoAuth},
		{"только пробелы", "   ", "", errors2.ErrNoAuth},
		{"без префикса Bearer", "abc123", "", errors2.ErrBadToken},
		{"обрезка пробелов", "  Bearer tok ", "tok", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := utils.ValidateBearerToken(tc.header)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("не ожидали ошибку: %v", err)
				}
				if got != tc.want {
					t.Errorf("ожидали токен %q, получили %q", tc.want, got)
				}
				return
			}
			if err != tc.wantErr {
				t.Errorf("ожидали %v, получили %v", tc.wantErr, err)
			}
		})
	}
}

func TestIncomingCreationValidation(t *testing.T) {
	cases := []struct {
		name        string
		title       string
		description string
		price       float64
		wantErr     error
	}{
		{"валидные данные", "Title", "Desc", 100, nil},
		{"пустой заголовок", "", "Desc", 100, errors2.ErrEmptyTitle},
		{"заголовок из пробелов", "  ", "Desc", 100, errors2.ErrEmptyTitle},
		{"пустое описание", "Title", "", 100, errors2.ErrEmptyDescription},
		{"нулевая цена", "Title", "Desc", 0, errors2.ErrEmptyPrice},
		{"отрицательная цена", "Title", "Desc", -5, errors2.ErrEmptyPrice},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.IncomingCreationValidation(tc.title, tc.description, tc.price)
			if err != tc.wantErr {
				t.Errorf("ожидали %v, получили %v", tc.wantErr, err)
			}
		})
	}
}

func TestIsDataExpired(t *testing.T) {
	if err := utils.IsDataExpired(time.Now().Add(-time.Minute), time.Hour); err != nil {
		t.Errorf("свежие данные не должны истекать: %v", err)
	}
	if err := utils.IsDataExpired(time.Now().Add(-2*time.Hour), time.Hour); err != errors2.ErrDataExpired {
		t.Errorf("ожидали ErrDataExpired, получили %v", err)
	}
}

func TestGenerateVerificationCode(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := utils.GenerateVerificationCode()
		if len(code) != 6 {
			t.Fatalf("код %q имеет длину %d, а не 6", code, len(code))
		}
		for _, r := range code {
			if r < '0' || r > '9' {
				t.Fatalf("код %q содержит нецифровой символ %q", code, r)
			}
		}
	}
}

func newTestContext(t *testing.T) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

func TestCheckRoles(t *testing.T) {
	c := newTestContext(t)
	c.Set("user_id", uint(7))
	c.Set("role", "customer")

	uid, err := utils.CheckCustomerRole(c)
	if err != nil || uid != 7 {
		t.Errorf("customer должен пройти проверку: uid=%d err=%v", uid, err)
	}
	if _, err := utils.CheckPerformerRole(c); err != errors2.ErrNotPerformer {
		t.Errorf("ожидали ErrNotPerformer, получили %v", err)
	}

	c.Set("role", "performer")
	if _, err := utils.CheckCustomerRole(c); err != errors2.ErrNotCustomer {
		t.Errorf("ожидали ErrNotCustomer, получили %v", err)
	}
	if uid, err := utils.CheckPerformerRole(c); err != nil || uid != 7 {
		t.Errorf("performer должен пройти проверку: uid=%d err=%v", uid, err)
	}
}

func TestBindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"title":"T","price":10}`))
	c.Request.Header.Set("Content-Type", "application/json")

	req, err := utils.BindJSON[struct {
		Title string  `json:"title"`
		Price float64 `json:"price"`
	}](c)
	if err != nil {
		t.Fatalf("не ожидали ошибку: %v", err)
	}
	if req.Title != "T" || req.Price != 10 {
		t.Errorf("неверно распарсен JSON: %+v", req)
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid`))

	if _, err := utils.BindJSON[struct{ Title string }](c2); err != errors2.ErrWrongJson {
		t.Errorf("ожидали ErrWrongJson, получили %v", err)
	}
}
