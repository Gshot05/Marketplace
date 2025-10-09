package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	RegisterReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Claims struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	PendingUser struct {
		Password  string    `json:"-"`
		Role      string    `json:"role"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
)
