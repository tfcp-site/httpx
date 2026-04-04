// Package httplog provides HTTP middleware that logs one structured JSON line
// per incoming request, following the tfcp-site log format spec.
package httplog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/tfcp-site/httpx/correlation"
)

type contextKey struct{}

type requestMeta struct {
	cached bool
}

// MarkCached signals to Middleware that this request was served from cache.
// When marked, Middleware logs cached:true and omits duration_ms.
// It is a no-op when called outside of a Middleware context.
func MarkCached(ctx context.Context) {
	if m, ok := ctx.Value(contextKey{}).(*requestMeta); ok {
		m.cached = true
	}
}

// Middleware logs one JSON line per incoming HTTP request using log.
// Common fields: request_id (from correlation context), method, path, status.
// Timing: duration_ms for normal requests, cached:true for cache hits.
func Middleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			meta := &requestMeta{}
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r.WithContext(
				context.WithValue(r.Context(), contextKey{}, meta),
			))

			attrs := []slog.Attr{
				slog.String("request_id", correlation.FromContext(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.status),
			}
			if meta.cached {
				attrs = append(attrs, slog.Bool("cached", true))
			} else {
				attrs = append(attrs, slog.Int64("duration_ms", time.Since(start).Milliseconds()))
			}

			log.LogAttrs(r.Context(), slog.LevelInfo, "http request", attrs...)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
