package model

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // customer | performer
	Name         string `json:"name"`
}

type Offer struct {
	ID          uint    `json:"id"`
	CustomerID  uint    `json:"customer_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Service struct {
	ID          uint    `json:"id"`
	PerformerID uint    `json:"performer_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Favorite struct {
	ID         uint `json:"id"`
	CustomerID uint `json:"customer_id"`
	ServiceID  uint `json:"service_id"`
}

type FavoriteInfo struct {
	ID                 uint   `json:"id"`
	CustomerName       string `json:"customer_name"`
	ServiceTitle       string `json:"service_title"`
	ServiceDescription string `json:"service_description"`
}
