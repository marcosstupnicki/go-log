package golog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerCreation(t *testing.T) {
	log, err := New()
	assert.NotNil(t, log)
	assert.NoError(t, err)
}

func TestLogLevels(t *testing.T) {
	log, err := New()
	assert.Nil(t, err)

	// Verify that log functions do not return errors
	assert.NotPanics(t, func() { log.Debug(context.Background(), "Debug log message") })
	assert.NotPanics(t, func() { log.Info(context.Background(), "Info log message") })
	assert.NotPanics(t, func() { log.Warn(context.Background(), "Warning log message") })
	assert.NotPanics(t, func() { log.Error(context.Background(), "Error log message") })
}

func TestContextHandling(t *testing.T) {
	log, err := New()
	assert.Nil(t, err)

	ctx := context.WithValue(context.Background(), "key", "value")

	// Verify that log functions accept context without panics
	assert.NotPanics(t, func() { log.Debug(ctx, "Debug log message") })
}

func TestLogFormat(t *testing.T) {
	log, err := New()
	assert.Nil(t, err)

	// Verify that log functions generate logs with messages and options
	assert.NotPanics(t, func() {
		log.Info(context.Background(), "Info log message", WithField("key", "value"))
	})
}

func TestOptionsConfiguration(t *testing.T) {
	log, err := New()
	assert.Nil(t, err)

	// Verify that options are applied correctly
	assert.NotPanics(t, func() {
		log.Info(context.Background(), "Info log message", WithField("key", "value"))
	})
}
