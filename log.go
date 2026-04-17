package golog

import (
	"context"

	"go.uber.org/zap"
)

// Logger is the main logging structure. It is safe for concurrent use.
// Pass it as a dependency (constructor injection), not via context.
// Only log FIELDS go in context (via Enrich).
type Logger struct {
	log      logger    // zap interface
	envField zap.Field // pre-computed "env" field
}

// New creates a Logger for the given environment. The environment string
// is added as an "env" field to every log entry. Returns an error if the
// underlying zap logger cannot be built.
//
//	logger, err := golog.New("local")
//	logger, err := golog.New("prod", golog.WithLevel(golog.DebugLevel))
func New(environment string, opts ...LogOption) (Logger, error) {
	cfg := defaultLogConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	zapLog, err := newZapLogger(cfg.level)
	if err != nil {
		return Logger{}, err
	}

	return Logger{
		log:      zapLog,
		envField: zap.String("env", environment),
	}, nil
}

// Debug logs a debug-level message. Context fields (from Enrich) are
// automatically merged with the env field and call-site fields.
func (l Logger) Debug(ctx context.Context, msg string, fields ...Option) {
	l.log.Debug(msg, l.mergeAll(ctx, fields...)...)
}

// Info logs an info-level message.
func (l Logger) Info(ctx context.Context, msg string, fields ...Option) {
	l.log.Info(msg, l.mergeAll(ctx, fields...)...)
}

// Warn logs a warn-level message.
func (l Logger) Warn(ctx context.Context, msg string, fields ...Option) {
	l.log.Warn(msg, l.mergeAll(ctx, fields...)...)
}

// Error logs an error-level message.
func (l Logger) Error(ctx context.Context, msg string, fields ...Option) {
	l.log.Error(msg, l.mergeAll(ctx, fields...)...)
}

// mergeAll combines env field + context fields + call-site fields.
// Order: env, context (Enrich), call-site (per-log-call).
// Later fields with the same key override earlier ones in zap.
func (l Logger) mergeAll(ctx context.Context, callFields ...Option) []zap.Field {
	applied := applyOptions(callFields...)
	ctxFields := fieldsFromContext(ctx)

	total := 1 + len(ctxFields) + len(applied.fields) // 1 for envField
	all := make([]zap.Field, 0, total)
	all = append(all, l.envField)
	all = append(all, ctxFields...)
	all = append(all, applied.fields...)
	return all
}
