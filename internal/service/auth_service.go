package service

import (
	"context"
	"marketplace/internal/model"
	"marketplace/internal/notifications"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"
)

type AuthService struct {
	repo          repository.IAuthRepo
	notifications notifications.INotifications
}

func NewAuthService(repo repository.IAuthRepo,
	notifications notifications.INotifications) *AuthService {
	return &AuthService{
		repo:          repo,
		notifications: notifications,
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

	go func() {
		err := s.notifications.SendRegistrationSuccess(context.Background(), email)
		if err != nil {
			return
		}
	}()

	return userID, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.notifications.SendLoginNotification(context.Background(), email)
		if err != nil {
			return
		}
	}()

	return user, nil
}
