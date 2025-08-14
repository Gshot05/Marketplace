package utils

import (
	errors "marketplace/internal/error"
	"net/mail"

	"github.com/gin-gonic/gin"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func ValidateName(name string) error {
	if name == "" || name == " " {
		return errors.ErrEmptyName
	}
	return nil
}

func CheckRole(role string) error {
	if role == "" || role == " " {
		return errors.ErrEmptyRole
	}
	return nil
}

func CheckCustomerRole(c *gin.Context) (uint, error) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "customer" {
		return 0, errors.ErrNotCustomer
	}
	return uid, nil
}

func CheckPerformerRole(c *gin.Context) (uint, error) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "performer" {
		return 0, errors.ErrNotPerformer
	}
	return uid, nil
}

func BindJSON[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, errors.ErrWrongJson
	}
	return request, nil
}
