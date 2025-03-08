package golog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// Logger is the main logging structure.
type Logger struct {
	log logger
}

// New creates a new Logger instance.
func New() (Logger, error) {
	res, err := newZapLogger()
	if err != nil {
		return Logger{}, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return Logger{
		log: res,
	}, nil
}

// Debug logs a debug message.
func (r Logger) Debug(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Debug(msg, opts.fields...)
}

// Info logs an info message.
func (r Logger) Info(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Info(msg, opts.fields...)
}

// Warn logs a warning message.
func (r Logger) Warn(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Warn(msg, opts.fields...)
}

// Error logs an error message.
func (r Logger) Error(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Error(msg, opts.fields...)
}
