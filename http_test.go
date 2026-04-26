package golog

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPMiddleware_EnrichesRequestContext(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	handler := HTTPMiddleware(logger, func(*http.Request) string {
		return "req-123"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := fieldsFromContext(r.Context())
		keys := make(map[string]bool, len(fields))
		for _, field := range fields {
			keys[field.Key] = true
		}

		for _, key := range []string{"request_id", "method", "path", "remote_addr", "user_agent"} {
			if !keys[key] {
				t.Fatalf("context fields missing %q: %#v", key, keys)
			}
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/hooks/source-1", nil)
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
}

func TestHTTPMiddleware_DoesNotAdvertiseUnsupportedFlusher(t *testing.T) {
	logger, err := New("test", WithLevel(DebugLevel))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	var sawFlusher bool
	handler := HTTPMiddleware(logger, nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, sawFlusher = w.(http.Flusher)
		w.WriteHeader(http.StatusNoContent)
	}))

	rw := &basicResponseWriter{header: make(http.Header)}
	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	handler.ServeHTTP(rw, req)

	if sawFlusher {
		t.Fatal("middleware exposed http.Flusher even though the original writer did not support it")
	}
	if rw.status != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rw.status, http.StatusNoContent)
	}
}

type basicResponseWriter struct {
	header http.Header
	status int
	body   []byte
}

func (w *basicResponseWriter) Header() http.Header {
	return w.header
}

func (w *basicResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *basicResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = append(w.body, b...)
	return len(b), nil
}
