package utils

import (
	"errors"
	"net/mail"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func ValidateName(name string) error {
	if name == "" || name == " " {
		return errors.New("name cannot be empty")
	}
	return nil
}
