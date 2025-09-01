package service

import (
	"context"
	"marketplace/internal/model"
)

type IOfferService interface {
	CreateOffer(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error)
	UpdateOffer(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error)
	DeleteOffer(ctx context.Context, offerID, customerID uint) (bool, error)
	ListOffers(ctx context.Context) ([]model.Offer, error)
}

type IServiceService interface {
	CreateService(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error)
	UpdateService(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error)
	DeleteService(ctx context.Context, serviceID, performerID uint) (bool, error)
	ListServices(ctx context.Context) ([]model.Service, error)
}

type IFavoriteService interface {
	AddFavorite(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error)
	DeleteFavorite(ctx context.Context, customerID, serviceID uint) (bool, error)
	ListFavorites(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error)
}

type IAuthService interface {
	RegisterUser(ctx context.Context, email, password, role, name string) (uint, error)
	LoginUser(ctx context.Context, email string) (*model.User, error)
}
