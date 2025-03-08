package golog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
