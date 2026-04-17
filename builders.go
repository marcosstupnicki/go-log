package golog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger is the interface satisfied by *zap.Logger.
type logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// newZapLogger creates a production zap.Logger at the given level.
// Returns an error if the logger cannot be built.
func newZapLogger(level zapcore.Level) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("golog: failed to build zap logger: %w", err)
	}

	return logger, nil
}
