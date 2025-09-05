package utils

import (
	errors2 "marketplace/internal/error"

	"github.com/gin-gonic/gin"
)

func CheckCustomerRole(c *gin.Context) (uint, error) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "customer" {
		return 0, errors2.ErrNotCustomer
	}
	return uid, nil
}

func CheckPerformerRole(c *gin.Context) (uint, error) {
	uid := c.GetUint("user_id")
	role := c.GetString("role")

	if role != "performer" {
		return 0, errors2.ErrNotPerformer
	}
	return uid, nil
}

func BindJSON[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, errors2.ErrWrongJson
	}
	return request, nil
}
