package golog

import (
	"net/http"

	"github.com/felixge/httpsnoop"
)

const defaultHTTPLogMessage = "http_request"

// HTTPMiddleware enriches the request context with request-scoped log fields
// and emits one access log after the downstream handler completes.
// requestID may be nil. It exists so callers can integrate request-id
// middleware from chi, gin, echo or other HTTP stacks without go-log depending
// on any of them.
func HTTPMiddleware(log Logger, requestID func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fields := []Option{
				Field("method", r.Method),
				Field("path", r.URL.Path),
				Field("remote_addr", r.RemoteAddr),
				Field("user_agent", r.UserAgent()),
			}
			if requestID != nil {
				if rid := requestID(r); rid != "" {
					fields = append(fields, Field("request_id", rid))
				}
			}

			ctx := Enrich(r.Context(), fields...)
			metrics := httpsnoop.CaptureMetrics(next, w, r.WithContext(ctx))

			log.Info(ctx, defaultHTTPLogMessage,
				Field("status", metrics.Code),
				Field("bytes", metrics.Written),
				Field("duration_ms", metrics.Duration.Milliseconds()),
			)
		})
	}
}
