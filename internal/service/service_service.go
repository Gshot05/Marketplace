package service

import (
	"context"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
)

type ServiceService struct {
	repo repository.IServiceRepo
}

func NewServiceService(repo repository.IServiceRepo) *ServiceService {
	return &ServiceService{
		repo: repo,
	}
}

func (s *ServiceService) CreateService(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error) {
	service, err := s.repo.Create(ctx, performerID, title, description, price)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) UpdateService(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error) {
	service, err := s.repo.Update(ctx, serviceID, performerID, title, description, price)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceService) DeleteService(ctx context.Context, serviceID, performerID uint) (bool, error) {
	deleted, err := s.repo.Delete(ctx, serviceID, performerID)
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

func (s *ServiceService) ListServices(ctx context.Context) ([]model.Service, error) {
	service, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return service, nil
}
