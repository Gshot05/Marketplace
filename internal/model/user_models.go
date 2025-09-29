package model

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // customer | performer
	Name         string `json:"name"`
}

type VerifyUser struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verifyCode"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Name       string `json:"name"`
}
