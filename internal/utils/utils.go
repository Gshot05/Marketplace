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
		return errors.New("name cannot be empty")
	}
	return nil
}

func CheckCustomerRole(c *gin.Context) (uint, bool) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "customer" {
		c.JSON(403, gin.H{"error": "only customers can perform this action"})
		return 0, false
	}

	return uid, true
}

func CheckPerformerRole(c *gin.Context) (uint, bool) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "performer" {
		c.JSON(403, gin.H{"error": "only performers can perform this action"})
		return 0, false
	}

	return uid, true
}

func BindJSON[T any](c *gin.Context) (T, bool) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body: " + err.Error()})
		return request, false
	}
	return request, true
}
