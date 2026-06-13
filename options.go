package golog

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// --- Per-log-call field options ---

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

// Option adds a field to a log entry or to context via Enrich.
type Option = func(*opts)

// Field creates a structured log field. The type of value is detected
// automatically to use the most efficient zap encoder.
//
// For error values the key is ignored — zap always uses "error" as the
// key (standard convention). So Field("whatever", err) logs as "error": "...".
func Field(key string, value interface{}) Option {
	return func(o *opts) {
		switch v := value.(type) {
		case string:
			o.fields = append(o.fields, zap.String(key, v))
		case int:
			o.fields = append(o.fields, zap.Int(key, v))
		case int64:
			o.fields = append(o.fields, zap.Int64(key, v))
		case float64:
			o.fields = append(o.fields, zap.Float64(key, v))
		case bool:
			o.fields = append(o.fields, zap.Bool(key, v))
		case error:
			o.fields = append(o.fields, zap.Error(v))
		case time.Duration:
			o.fields = append(o.fields, zap.Duration(key, v))
		case time.Time:
			o.fields = append(o.fields, zap.Time(key, v))
		case fmt.Stringer:
			o.fields = append(o.fields, zap.String(key, v.String()))
		default:
			o.fields = append(o.fields, zap.Any(key, v))
		}
	}
}

// --- Constructor options (LogOption) ---

// logConfig holds configuration for the logger builder.
type logConfig struct {
	level Level
}

// defaultLogConfig returns the default logger configuration.
func defaultLogConfig() logConfig {
	return logConfig{
		level: InfoLevel,
	}
}

// LogOption is a functional option for New().
type LogOption func(*logConfig)

// WithLevel sets the minimum log level. Use the package-level constants:
// golog.DebugLevel, golog.InfoLevel, golog.WarnLevel, golog.ErrorLevel.
// Default is InfoLevel.
func WithLevel(level Level) LogOption {
	return func(c *logConfig) {
		c.level = level
	}
}
