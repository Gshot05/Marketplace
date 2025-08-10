package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func RunInitSQL(pool *pgxpool.Pool, filepath string) error {
	ctx := context.Background()

	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	sql := string(sqlBytes)

	_, err = pool.Exec(ctx, sql)
	if err != nil {
		log.Printf("Failed to execute init SQL: %v", err)
		return err
	}

	return nil
}

func Connect() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = RunInitSQL(pool, os.Getenv("MIGRATIONS"))
	if err != nil {
		log.Fatalf("Failed to run init.sql: %v", err)
	}

	return pool
}
