// Package logging provides a pre-configured slog.Logger for tfcp-site services.
package logging

import (
	"context"
	"log/slog"
	"os"
)

// ContextExtractor extracts slog attributes from a context.
// It is called on every log record written with a context.
type ContextExtractor func(context.Context) []slog.Attr

// New returns a JSON slog.Logger that writes to stdout.
// The time field is omitted — Promtail records ingestion time automatically.
// The service field is pre-attached to every log line.
// Extractors are called on each log record to inject attributes from context.
func New(service string, level slog.Level, extractors ...ContextExtractor) *slog.Logger {
	inner := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(&contextHandler{inner: inner, extractors: extractors}).With("service", service)
}

type contextHandler struct {
	inner      slog.Handler
	extractors []ContextExtractor
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, ext := range h.extractors {
		r.AddAttrs(ext(ctx)...)
	}
	return h.inner.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{inner: h.inner.WithAttrs(attrs), extractors: h.extractors}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{inner: h.inner.WithGroup(name), extractors: h.extractors}
}
