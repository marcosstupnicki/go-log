package golog

import "go.uber.org/zap"

type opts struct {
	fields []zap.Field
}

func applyOptions(options ...Option) opts {
	var res opts
	for _, opt := range options {
		opt(&res)
	}

	return res
}

type Option = func(*opts)

func WithField(key string, value interface{}) Option {
	return func(s *opts) {
		s.fields = append(s.fields, zap.Any(key, value))
	}
}
