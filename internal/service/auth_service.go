package service

import (
	"context"
	"errors"
	"marketplace/internal/auth"
	"marketplace/internal/model"
	"marketplace/internal/notifications"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"
)

type AuthService struct {
	repo          repository.IAuthRepo
	notifications notifications.INotifications
}

func NewAuthService(
	repo repository.IAuthRepo,
	notifications notifications.INotifications,
) *AuthService {
	return &AuthService{
		repo:          repo,
		notifications: notifications,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password, role, name string) error {
	err := utils.ValidateIncomingRegistration(email, name, role)
	if err != nil {
		return err
	}

	code := utils.GenerateVerificationCode()

	err = s.repo.SaveVerificationCode(ctx, email, code)
	if err != nil {
		return err
	}

	go func() {
		err := s.notifications.SendVerificationCode(context.Background(), email, code)
		if err != nil {
			return
		}
	}()

	return nil
}

func (s *AuthService) ConfirmEmail(ctx context.Context, email, code, password, role, name string) (string, error) {
	valid, err := s.repo.VerifyCode(ctx, email, code)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("неверный или просроченный код подтверждения")
	}

	userID, err := s.repo.RegisterUser(ctx, email, password, role, name)
	if err != nil {
		return "", err
	}

	err = s.repo.DeleteUsedCode(ctx, email)
	if err != nil {
	}

	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		return "", err
	}

	return token, nil
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
