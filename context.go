package golog

import (
	"context"

	"go.uber.org/zap"
)

type ctxFieldsKey struct{}

// Enrich returns a new context with the given fields appended to any
// fields already stored. Typically called once in middleware to add
// request-scoped fields (request_id, method, path). Every subsequent
// log call that receives this context will include these fields
// automatically — no need for FromContext or WithContext.
//
//	ctx = golog.Enrich(ctx,
//	    golog.Field("request_id", rid),
//	    golog.Field("path", r.URL.Path),
//	)
func Enrich(ctx context.Context, fields ...Option) context.Context {
	applied := applyOptions(fields...)
	if len(applied.fields) == 0 {
		return ctx
	}

	existing := fieldsFromContext(ctx)
	merged := make([]zap.Field, 0, len(existing)+len(applied.fields))
	merged = append(merged, existing...)
	merged = append(merged, applied.fields...)

	return context.WithValue(ctx, ctxFieldsKey{}, merged)
}

// fieldsFromContext extracts log fields previously stored via Enrich.
// Returns nil if no fields are present.
func fieldsFromContext(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}
	fields, _ := ctx.Value(ctxFieldsKey{}).([]zap.Field)
	return fields
}
