package golog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrich_AddsFieldsToContext(t *testing.T) {
	ctx := context.Background()
	ctx = Enrich(ctx, Field("request_id", "abc-123"))

	fields := fieldsFromContext(ctx)
	assert.Len(t, fields, 1)
	assert.Equal(t, "request_id", fields[0].Key)
}

func TestEnrich_AccumulatesFields(t *testing.T) {
	ctx := context.Background()
	// First enrichment (e.g. in middleware)
	ctx = Enrich(ctx, Field("request_id", "abc-123"), Field("method", "GET"))
	// Second enrichment (e.g. deeper in the call chain)
	ctx = Enrich(ctx, Field("user_id", "user-456"))

	fields := fieldsFromContext(ctx)
	assert.Len(t, fields, 3)
}

func TestEnrich_EmptyOptionsIsNoop(t *testing.T) {
	ctx := context.Background()
	ctx2 := Enrich(ctx)

	// Should return the same context (no new value stored)
	assert.Equal(t, ctx, ctx2)
}

func TestFieldsFromContext_NilContext(t *testing.T) {
	//nolint:staticcheck
	fields := fieldsFromContext(nil)
	assert.Nil(t, fields)
}

func TestFieldsFromContext_EmptyContext(t *testing.T) {
	fields := fieldsFromContext(context.Background())
	assert.Nil(t, fields)
}

func TestEnrich_IntegrationWithLogger(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	require.NoError(t, err)

	// Simulate middleware enriching context
	ctx := context.Background()
	ctx = Enrich(ctx,
		Field("request_id", "req-001"),
		Field("path", "/api/users"),
	)

	// Simulate handler logging — context fields are included automatically
	assert.NotPanics(t, func() {
		logger.Info(ctx, "handler invoked", Field("user_id", "u-42"))
	})
}

func TestEnrich_MultipleMiddlewares(t *testing.T) {
	// Simulates: auth middleware adds user_id, then tracing middleware adds trace_id
	ctx := context.Background()
	ctx = Enrich(ctx, Field("trace_id", "t-001"))
	ctx = Enrich(ctx, Field("user_id", "u-42"))

	fields := fieldsFromContext(ctx)
	assert.Len(t, fields, 2)

	keys := make([]string, len(fields))
	for i, f := range fields {
		keys[i] = f.Key
	}
	assert.Contains(t, keys, "trace_id")
	assert.Contains(t, keys, "user_id")
}
