package service

import (
	"context"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
)

type OfferService struct {
	repo repository.IOfferRepo
}

func NewOfferService(repo repository.IOfferRepo) *OfferService {
	return &OfferService{
		repo: repo,
	}
}

func (s *OfferService) CreateOffer(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error) {
	offer, err := s.repo.Create(ctx, customerID, title, description, price)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (s *OfferService) UpdateOffer(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error) {
	offer, err := s.repo.Update(ctx, offerID, customerID, title, description, price)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (s *OfferService) DeleteOffer(ctx context.Context, offerID, customerID uint) (bool, error) {
	deleted, err := s.repo.Delete(ctx, offerID, customerID)
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

func (s *OfferService) ListOffers(ctx context.Context) ([]model.Offer, error) {
	offer, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return offer, nil
}
