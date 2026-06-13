package golog

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is go-log's public log level type. It aliases zapcore.Level because
// go-log intentionally wraps zap while keeping call sites out of zapcore.
type Level = zapcore.Level

const (
	DebugLevel Level = zap.DebugLevel
	InfoLevel  Level = zap.InfoLevel
	WarnLevel  Level = zap.WarnLevel
	ErrorLevel Level = zap.ErrorLevel
)

// ParseLevel parses common textual log levels accepted by go-log.
func ParseLevel(v string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	default:
		return InfoLevel, fmt.Errorf("golog: invalid log level %q (valid: debug, info, warn, error)", v)
	}
}

// MustParseLevel parses v and panics if it is invalid. It is useful in config
// loaders that already panic on invalid environment values.
func MustParseLevel(v string) Level {
	level, err := ParseLevel(v)
	if err != nil {
		panic(err)
	}
	return level
}
