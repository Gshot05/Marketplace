package handlers_test

import (
	errors2 "marketplace/internal/error"
	"net/http"
	"testing"
)

func TestCreateOfferRoleCheck(t *testing.T) {
	r := newTestRouter("performer", 1)
	buildOfferRoutes(r)

	w := doJSON(t, r, http.MethodPost, "/offers", `{"title":"T","description":"D","price":10}`)
	if w.Code != http.StatusForbidden {
		t.Errorf("performer на /offers: код %d, ожидали 403", w.Code)
	}

	r = newTestRouter("customer", 1)
	buildOfferRoutes(r)
	w = doJSON(t, r, http.MethodPost, "/offers", `{"title":"T","description":"D","price":10}`)
	if w.Code != http.StatusOK {
		t.Errorf("customer на /offers: код %d, ожидали 200; body=%s", w.Code, w.Body.String())
	}
}

func TestCreateOfferBadJSON(t *testing.T) {
	r := newTestRouter("customer", 1)
	buildOfferRoutes(r)

	w := doJSON(t, r, http.MethodPost, "/offers", `{invalid`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("код %d, ожидали 400", w.Code)
	}
	if w.Body.String() != `{"error":"`+errors2.ErrWrongJson.Error()+`"}` {
		t.Errorf("неожиданное тело: %s", w.Body.String())
	}
}

func TestCreateOfferValidationErrorMapsTo500(t *testing.T) {
	r := newTestRouter("customer", 1)
	buildOfferRoutes(r)

	w := doJSON(t, r, http.MethodPost, "/offers", `{"title":"","description":"","price":0}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("код %d, ожидали 500 (хендлер маппит любые ошибки сервиса на 500)", w.Code)
	}
}

func TestDeleteOfferNotFoundAndSuccess(t *testing.T) {
	r := newTestRouter("customer", 1)
	buildOfferRoutes(r)

	w := doJSON(t, r, http.MethodDelete, "/offers", `{"offerID":99}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("чужой/несуществующий оффер: код %d, ожидали 404", w.Code)
	}

	w = doJSON(t, r, http.MethodDelete, "/offers", `{"offerID":1}`)
	if w.Code != http.StatusOK {
		t.Errorf("свой оффер: код %d, ожидали 200", w.Code)
	}
}

func TestListOffersEmptyMapsTo500(t *testing.T) {
	r := newTestRouter("customer", 1)
	buildOfferRoutes(r)

	w := doJSON(t, r, http.MethodGet, "/offers", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("пустой список: код %d, ожидали 500", w.Code)
	}
}

func TestServiceHandlerPerformerOnly(t *testing.T) {
	r := newTestRouter("customer", 1)
	buildServiceRoutes(r)

	w := doJSON(t, r, http.MethodPost, "/services", `{"title":"T","description":"D","price":10}`)
	if w.Code != http.StatusForbidden {
		t.Errorf("customer на /services: код %d, ожидали 403", w.Code)
	}

	r = newTestRouter("performer", 1)
	buildServiceRoutes(r)
	w = doJSON(t, r, http.MethodPost, "/services", `{"title":"T","description":"D","price":10}`)
	if w.Code != http.StatusOK {
		t.Errorf("performer на /services: код %d, ожидали 200; body=%s", w.Code, w.Body.String())
	}
}

func TestFavoriteHandlerCustomerOnlyAndFlow(t *testing.T) {
	r := newTestRouter("performer", 1)
	buildFavoriteRoutes(r)

	w := doJSON(t, r, http.MethodPost, "/favorites", `{"serviceID":5}`)
	if w.Code != http.StatusForbidden {
		t.Errorf("performer на /favorites: код %d, ожидали 403", w.Code)
	}

	r = newTestRouter("customer", 1)
	buildFavoriteRoutes(r)

	w = doJSON(t, r, http.MethodPost, "/favorites", `{"serviceID":5}`)
	if w.Code != http.StatusOK {
		t.Errorf("добавление в избранное: код %d, ожидали 200; body=%s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodGet, "/favorites", "")
	if w.Code != http.StatusOK {
		t.Errorf("список избранного: код %d, ожидали 200", w.Code)
	}
}
