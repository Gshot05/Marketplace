package model

import "time"

type User struct {
	ID           uint
	Email        string
	PasswordHash string
	Role         string // customer | performer
	Name         string
	CreatedAt    time.Time
}

type Offer struct {
	ID          uint
	CustomerID  uint
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type Service struct {
	ID          uint
	PerformerID uint
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type Favorite struct {
	ID         uint
	CustomerID uint
	ServiceID  uint
	CreatedAt  time.Time
}
