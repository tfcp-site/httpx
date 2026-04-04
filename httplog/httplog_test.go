package httplog_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tfcp-site/httpx/correlation"
	"github.com/tfcp-site/httpx/httplog"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(h).With("service", "test-svc")
}

func parseLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var record map[string]any
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("invalid JSON log output: %v\nraw: %s", err, buf.String())
	}
	return record
}

func TestMiddleware_normalRequest(t *testing.T) {
	var buf bytes.Buffer
	h := httplog.Middleware(newTestLogger(&buf))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/chat", nil))

	record := parseLog(t, &buf)

	for _, field := range []string{"method", "path", "status", "duration_ms", "request_id"} {
		if _, ok := record[field]; !ok {
			t.Errorf("missing field %q", field)
		}
	}
	if _, ok := record["cached"]; ok {
		t.Error("cached field must be absent for normal requests")
	}
	if record["method"] != http.MethodPost {
		t.Errorf("method = %v, want POST", record["method"])
	}
	if record["path"] != "/api/chat" {
		t.Errorf("path = %v, want /api/chat", record["path"])
	}
	if record["status"] != float64(http.StatusCreated) {
		t.Errorf("status = %v, want %d", record["status"], http.StatusCreated)
	}
}

func TestMiddleware_cachedRequest(t *testing.T) {
	var buf bytes.Buffer
	h := httplog.Middleware(newTestLogger(&buf))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httplog.MarkCached(r.Context())
	}))

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/article/1", nil))

	record := parseLog(t, &buf)

	if record["cached"] != true {
		t.Errorf("cached = %v, want true", record["cached"])
	}
	if _, ok := record["duration_ms"]; ok {
		t.Error("duration_ms must be absent for cached requests")
	}
}

func TestMiddleware_readsRequestID(t *testing.T) {
	const wantID = "test-request-id"
	var buf bytes.Buffer
	h := correlation.Middleware(
		httplog.Middleware(newTestLogger(&buf))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(correlation.Header, wantID)
	h.ServeHTTP(httptest.NewRecorder(), req)

	record := parseLog(t, &buf)
	if record["request_id"] != wantID {
		t.Errorf("request_id = %v, want %q", record["request_id"], wantID)
	}
}

func TestMiddleware_defaultStatus200(t *testing.T) {
	var buf bytes.Buffer
	h := httplog.Middleware(newTestLogger(&buf))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no WriteHeader call — defaults to 200
	}))

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	record := parseLog(t, &buf)
	if record["status"] != float64(http.StatusOK) {
		t.Errorf("status = %v, want 200", record["status"])
	}
}
