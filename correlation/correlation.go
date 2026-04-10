// Package correlation provides HTTP middleware and transport for propagating
// a correlation ID (X-Request-ID) across service boundaries.
package correlation

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// Header is the HTTP header used to carry the correlation ID.
const Header = "X-Request-ID"

type contextKey struct{}

// FromContext returns the correlation ID stored in ctx.
// Returns an empty string if no ID is present.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}

// WithContext returns a new context carrying the given correlation ID.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// Middleware extracts the correlation ID from the X-Request-ID request header.
// If the header is absent, a new UUID is generated. The ID is stored in the
// request context and echoed back in the response header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(Header)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set(Header, id)
		next.ServeHTTP(w, r.WithContext(WithContext(r.Context(), id)))
	})
}

// Extractor returns a context extractor that injects the correlation ID into
// log records as a "request_id" attribute. Pass it to logging.New.
func Extractor() func(context.Context) []slog.Attr {
	return func(ctx context.Context) []slog.Attr {
		if id := FromContext(ctx); id != "" {
			return []slog.Attr{slog.String("request_id", id)}
		}
		return nil
	}
}

// Transport is an http.RoundTripper that injects the correlation ID from the
// request context into the outgoing X-Request-ID header.
type Transport struct {
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	if id := FromContext(req.Context()); id != "" {
		req = req.Clone(req.Context())
		req.Header.Set(Header, id)
	}
	return base.RoundTrip(req)
}
