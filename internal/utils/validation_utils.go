package utils

import (
	errors2 "marketplace/internal/error"
	"net/mail"
	"strings"
)

func ValidateIncomingRegistration(email, name, role string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	if name == "" || name == " " {
		return errors2.ErrEmptyName
	}

	if role == "" || role == " " {
		return errors2.ErrEmptyRole
	}

	return nil
}

func ValidateBearerToken(token string) (string, error) {
	const bearerPrefix = "Bearer "

	token = strings.TrimSpace(token)
	if token == "" || token == " " {
		return "", errors2.ErrNoAuth
	}

	if !strings.HasPrefix(token, bearerPrefix) {
		return "", errors2.ErrBadToken
	}

	return strings.TrimPrefix(token, bearerPrefix), nil
}

func IncomingCreationValidation(title, description string, price float64) error {
	title = strings.TrimSpace(title)
	if title == "" || title == " " {
		return errors2.ErrEmptyTitle
	}

	description = strings.TrimSpace(description)
	if description == "" || description == " " {
		return errors2.ErrEmptyDescription
	}

	if price <= 0 {
		return errors2.ErrEmptyPrice
	}

	return nil
}
