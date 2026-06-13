package golog

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultLevel(t *testing.T) {
	logger, err := New("local")
	require.NoError(t, err)
	require.NotNil(t, logger.log)
	assert.Equal(t, "env", logger.envField.Key)
}

func TestNew_WithLevel(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New("test", WithLevel(tt.level))
			require.NoError(t, err)
			assert.NotNil(t, logger.log)
		})
	}
}

func TestNew_EnvironmentInEnvField(t *testing.T) {
	logger, err := New("staging")
	require.NoError(t, err)
	assert.Equal(t, "env", logger.envField.Key)
}

func TestLogMethods_NoPanic(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	require.NoError(t, err)
	ctx := context.Background()

	assert.NotPanics(t, func() {
		logger.Debug(ctx, "debug message")
		logger.Info(ctx, "info message")
		logger.Warn(ctx, "warn message")
		logger.Error(ctx, "error message")
	})
}

func TestLogMethods_WithFields(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	require.NoError(t, err)
	ctx := context.Background()

	assert.NotPanics(t, func() {
		logger.Info(ctx, "with fields",
			Field("name", "test"),
			Field("count", 42),
			Field("active", true),
			Field("elapsed", 150*time.Millisecond),
			Field("err", errors.New("test error")),
			Field("data", map[string]string{"key": "value"}),
		)
	})
}

func TestLogMethods_FieldTypes(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	require.NoError(t, err)
	ctx := context.Background()

	now := time.Now()
	assert.NotPanics(t, func() {
		logger.Info(ctx, "all field types",
			Field("str", "hello"),
			Field("int", 42),
			Field("int64", int64(999)),
			Field("float64", 3.14),
			Field("bool", true),
			Field("error", errors.New("oops")),
			Field("duration", 5*time.Second),
			Field("time", now),
			Field("slice", []string{"a", "b"}),
		)
	})
}

func TestMergeAll_EnvPlusContextPlusCallSite(t *testing.T) {
	logger, err := New("test")
	require.NoError(t, err)
	ctx := Enrich(context.Background(), Field("request_id", "abc-123"))

	// mergeAll should combine env + context (request_id) + call (status)
	fields := logger.mergeAll(ctx, Field("status", 200))

	// env + request_id + status = 3
	assert.Len(t, fields, 3)
}

func TestMergeAll_EnvAlwaysPresent(t *testing.T) {
	logger, err := New("prod")
	require.NoError(t, err)
	ctx := context.Background()

	fields := logger.mergeAll(ctx)
	// At minimum the env field is always present
	assert.Len(t, fields, 1)
	assert.Equal(t, "env", fields[0].Key)
}
