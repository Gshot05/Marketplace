package model

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"` // customer | performer
	Name         string
	CreatedAt    time.Time
}

type Offer struct {
	ID          uint `gorm:"primaryKey"`
	CustomerID  uint
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type Service struct {
	ID          uint `gorm:"primaryKey"`
	PerformerID uint
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type Favorite struct {
	ID         uint `gorm:"primaryKey"`
	CustomerID uint
	ServiceID  uint
	CreatedAt  time.Time
}
