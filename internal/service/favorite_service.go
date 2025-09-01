package service

import (
	"context"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
)

type FavoriteService struct {
	repo repository.IFavoriteRepo
}

func NewFavoriteService(repo repository.IFavoriteRepo) *FavoriteService {
	return &FavoriteService{
		repo: repo,
	}
}

func (s *FavoriteService) AddFavorite(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error) {
	fav, err := s.repo.Add(ctx, customerID, serviceID)
	if err != nil {
		return nil, err
	}
	return fav, nil
}

func (s *FavoriteService) DeleteFavorite(ctx context.Context, customerID, serviceID uint) (bool, error) {
	deleted, err := s.repo.Delete(ctx, customerID, serviceID)
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

func (s *FavoriteService) ListFavorites(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error) {
	fav, err := s.repo.List(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return fav, nil
}
