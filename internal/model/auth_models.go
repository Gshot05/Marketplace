package model

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
)
