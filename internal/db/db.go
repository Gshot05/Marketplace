package db

import (
	"log"
	"marketplace/internal/model"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to datrabse: v%", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Offer{},
		&model.Service{},
		&model.Favorite{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
