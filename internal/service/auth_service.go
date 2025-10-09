package service

import (
	"context"
	"marketplace/internal/auth"
	errors2 "marketplace/internal/error"
	"marketplace/internal/model"
	"marketplace/internal/notifications"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
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

var pendingRegistrations sync.Map

func (s *AuthService) RegisterUser(ctx context.Context, email, password, role, name string) error {
	err := utils.ValidateIncomingRegistration(email, name, role)
	if err != nil {
		return err
	}

	pendingUser := model.PendingUser{
		Password:  password,
		Role:      role,
		Name:      name,
		CreatedAt: time.Now(),
	}
	pendingRegistrations.Store(email, pendingUser)

	code := utils.GenerateVerificationCode()

	err = s.repo.SaveVerificationCode(ctx, email, code)
	if err != nil {
		pendingRegistrations.Delete(email)
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

func (s *AuthService) ConfirmEmail(ctx context.Context, email, code string) (string, error) {
	valid, err := s.repo.VerifyCode(ctx, email, code)
	if err != nil {
		return "", errors2.ErrWrongConfirmData
	}
	if !valid {
		return "", errors2.ErrWrongVerify
	}

	data, exists := pendingRegistrations.Load(email)
	if !exists {
		return "", errors2.ErrWrongConfirmData
	}

	pendingUser := data.(model.PendingUser)

	if time.Since(pendingUser.CreatedAt) > time.Hour {
		pendingRegistrations.Delete(email)
		return "", errors2.ErrWrongConfirmData
	}

	userID, err := s.repo.RegisterUser(ctx, email, pendingUser.Password, pendingUser.Role, pendingUser.Name)
	if err != nil {
		return "", err
	}

	pendingRegistrations.Delete(email)
	err = s.repo.DeleteUsedCode(ctx, email)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(userID, pendingUser.Role)
	if err != nil {
		return "", errors2.ErrCreateToken
	}

	return token, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", errors2.ErrWrongPassOrLog
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors2.ErrCreateToken
	}

	go func() {
		err := s.notifications.SendLoginNotification(context.Background(), email)
		if err != nil {
			return
		}
	}()

	return token, nil
}
