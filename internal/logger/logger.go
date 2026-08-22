package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"marketplace/internal/workerpool"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger struct {
	pool *pgxpool.Pool
	wp   *workerpool.Pool
}

func NewLogger(pool *pgxpool.Pool, wp *workerpool.Pool) *Logger {
	return &Logger{pool: pool, wp: wp}
}

func (l *Logger) LogAsync(level string, format string, args ...interface{}) {
	submitted := l.wp.Submit(func() {
		if l.pool == nil {
			return
		}
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
	})
	if !submitted {
		log.Printf("Log queue full, dropped [%s] %s", level, fmt.Sprintf(format, args...))
	}
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
