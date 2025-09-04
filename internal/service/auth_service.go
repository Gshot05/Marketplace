package service

import (
	"context"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"
)

type AuthService struct {
	repo repository.IAuthRepo
}

func NewAuthService(repo repository.IAuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password, role, name string) (uint, error) {
	err := utils.ValidateIncomingRegistration(email, name, role)
	if err != nil {
		return 0, err
	}

	userID, err := s.repo.RegisterUser(ctx, email, password, role, name)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
