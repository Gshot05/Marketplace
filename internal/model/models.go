package model

type User struct {
	ID           uint
	Email        string
	PasswordHash string
	Role         string // customer | performer
	Name         string
}

type Offer struct {
	ID          uint
	CustomerID  uint
	Title       string
	Description string
	Price       float64
}

type Service struct {
	ID          uint
	PerformerID uint
	Title       string
	Description string
	Price       float64
}

type Favorite struct {
	ID         uint
	CustomerID uint
	ServiceID  uint
}
