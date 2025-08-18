package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger struct {
	pool *pgxpool.Pool
}

// New создает логгер с pgxpool.Pool
func New(pool *pgxpool.Pool) *Logger {
	return &Logger{pool: pool}
}

func (l *Logger) Log(level string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)

	_, err := l.pool.Exec(
		context.Background(),
		"INSERT INTO logs (level, message, created_at) VALUES ($1, $2, $3)",
		level,
		message,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("logger failed: %w", err)
	}
	return nil
}

// Удобные методы-обертки
func (l *Logger) Info(format string, args ...interface{}) {
	_ = l.Log("INFO", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	_ = l.Log("ERROR", format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	_ = l.Log("DEBUG", format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	_ = l.Log("WARN", format, args...)
}
