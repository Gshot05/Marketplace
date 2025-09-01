package repository

import (
	"context"
	"marketplace/internal/model"
)

type IOfferRepo interface {
	Create(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error)
	Update(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error)
	Delete(ctx context.Context, offerID, customerID uint) (bool, error)
	List(ctx context.Context) ([]model.Offer, error)
}

type IServiceRepo interface {
	Create(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error)
	Update(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error)
	Delete(ctx context.Context, serviceID, performerID uint) (bool, error)
	List(ctx context.Context) ([]model.Service, error)
}

type IFavoriteRepo interface {
	Add(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error)
	Delete(ctx context.Context, customerID, serviceID uint) (bool, error)
	List(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error)
}

type IAuthRepo interface {
	RegisterUser(ctx context.Context, email, password, role, name string) (uint, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}
