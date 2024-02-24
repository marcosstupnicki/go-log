package golog

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, opt ...Option)
	Info(ctx context.Context, opt ...Option)
	Warn(ctx context.Context, opt ...Option)
	Error(ctx context.Context, opt ...Option)
}

type logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type Log struct {
	log logger
}

func New() (Log, error) {
	res, err := newZapLogger()
	if err != nil {
		return Log{}, err
	}

	return Log{
		log: res,
	}, nil

}

func (r Log) Debug(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Debug(msg, opts.fields...)
}

func (r Log) Info(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Info(msg, opts.fields...)
}

func (r Log) Warn(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Warn(msg, opts.fields...)
}

func (r Log) Error(ctx context.Context, msg string, options ...Option) {
	opts := applyOptions(options...)
	r.log.Error(msg, opts.fields...)
}
