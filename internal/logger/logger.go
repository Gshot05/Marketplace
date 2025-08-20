package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger struct {
	pool *pgxpool.Pool
}

func NewLogger(pool *pgxpool.Pool) *Logger {
	return &Logger{pool: pool}
}

func (l *Logger) LogAsync(level string, format string, args ...interface{}) {
	go func() {
		message := fmt.Sprintf(format, args...)
		_, err := l.pool.Exec(
			context.Background(),
			"INSERT INTO logs (level, message, created_at) VALUES ($1, $2, $3)",
			level,
			message,
			time.Now(),
		)
		if err != nil {
			log.Printf("Async log failed: %v", err)
		}
	}()
}

// Удобные методы-обертки
func (l *Logger) Info(format string, args ...interface{}) {
	l.LogAsync("INFO", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.LogAsync("ERROR", format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.LogAsync("DEBUG", format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.LogAsync("WARN", format, args...)
}
