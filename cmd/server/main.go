package main

import (
	"context"
	"log"
	"marketplace/internal/db"
	"marketplace/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	pool := db.Connect()

	r := gin.Default()
	logWP, mailWP := router.RegisterRoutes(r, pool)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %s\n", err)
		}
	}()
	log.Println("Сервер запущен на: 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %s\n", err)
	}

	if err := mailWP.Shutdown(ctx); err != nil {
		log.Printf("Mail worker pool shutdown: %v", err)
	}
	if err := logWP.Shutdown(ctx); err != nil {
		log.Printf("Log worker pool shutdown: %v", err)
	}

	pool.Close()
	log.Println("Сервер остановлен")
}
