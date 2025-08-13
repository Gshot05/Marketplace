package utils

import (
	"errors"
	"net/mail"

	"github.com/gin-gonic/gin"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func ValidateName(name string) error {
	if name == "" || name == " " {
		return errors.New("Имя не может быть пустым")
	}
	return nil
}

func CheckRole(role string) error {
	if role == "" || role == " " {
		return errors.New("Некорректная роль")
	}
	return nil
}

func CheckCustomerRole(c *gin.Context) (uint, bool) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "customer" {
		c.JSON(403, gin.H{"error": "Только заказчики имеют доступ к этой функции"})
		return 0, false
	}

	return uid, true
}

func CheckPerformerRole(c *gin.Context) (uint, bool) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "performer" {
		c.JSON(403, gin.H{"error": "Только исполнители имеют доступ к этой функции"})
		return 0, false
	}

	return uid, true
}

func BindJSON[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, err
	}
	return request, nil
}
