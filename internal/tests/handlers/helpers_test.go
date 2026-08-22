package handlers_test

import (
	"context"
	"marketplace/internal/handlers"
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/workerpool"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func testLogger() *logger.Logger {
	return logger.NewLogger(nil, workerpool.New(1, 8))
}

func newTestRouter(role string, uid uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("user_id", uid)
			c.Set("role", role)
		}
		c.Next()
	})
	return r
}

func doJSON(t *testing.T, r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type stubOfferRepo struct {
	list []model.Offer
}

func (s *stubOfferRepo) Create(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error) {
	return &model.Offer{ID: 1, CustomerID: customerID, Title: title}, nil
}

func (s *stubOfferRepo) Update(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error) {
	return &model.Offer{ID: offerID, CustomerID: customerID, Title: title}, nil
}

func (s *stubOfferRepo) Delete(ctx context.Context, offerID, customerID uint) (bool, error) {
	return offerID == 1 && customerID == 1, nil
}

func (s *stubOfferRepo) List(ctx context.Context) ([]model.Offer, error) {
	return s.list, nil
}

type stubServiceRepo struct{}

func (s *stubServiceRepo) Create(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error) {
	return &model.Service{ID: 1, PerformerID: performerID, Title: title}, nil
}

func (s *stubServiceRepo) Update(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error) {
	return &model.Service{ID: serviceID, PerformerID: performerID, Title: title}, nil
}

func (s *stubServiceRepo) Delete(ctx context.Context, serviceID, performerID uint) (bool, error) {
	return serviceID == 1 && performerID == 1, nil
}

func (s *stubServiceRepo) List(ctx context.Context) ([]model.Service, error) {
	return nil, nil
}

type stubFavoriteRepo struct{}

func (s *stubFavoriteRepo) Add(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error) {
	return &model.FavoriteReq{ID: 1, CustomerID: customerID, ServiceID: serviceID}, nil
}

func (s *stubFavoriteRepo) Delete(ctx context.Context, customerID, serviceID uint) (bool, error) {
	return serviceID == 1 && customerID == 1, nil
}

func (s *stubFavoriteRepo) List(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error) {
	return []model.FavoriteInfoReq{{ID: 1}}, nil
}

func buildOfferRoutes(r *gin.Engine) {
	offerHandler := handlers.NewOfferHandler(service.NewOfferService(&stubOfferRepo{}), testLogger())
	r.POST("/offers", offerHandler.CreateOffer())
	r.PATCH("/offers", offerHandler.UpdateOffer())
	r.DELETE("/offers", offerHandler.DeleteOffer())
	r.GET("/offers", offerHandler.ListOffers())
}

func buildServiceRoutes(r *gin.Engine) {
	serviceHandler := handlers.NewServiceHandler(service.NewServiceService(&stubServiceRepo{}), testLogger())
	r.POST("/services", serviceHandler.CreateService())
	r.DELETE("/services", serviceHandler.DeleteService())
}

func buildFavoriteRoutes(r *gin.Engine) {
	favoriteHandler := handlers.NewFavoriteHandler(service.NewFavoriteService(&stubFavoriteRepo{}), testLogger())
	r.POST("/favorites", favoriteHandler.AddFavorite())
	r.GET("/favorites", favoriteHandler.ListFavorites())
}
