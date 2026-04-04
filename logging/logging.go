// Package logging provides a pre-configured slog.Logger for tfcp-site services.
package logging

import (
	"log/slog"
	"os"
)

// New returns a JSON slog.Logger that writes to stdout.
// The time field is omitted — Promtail records ingestion time automatically.
// The service field is pre-attached to every log line.
func New(service string, level slog.Level) *slog.Logger {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	return slog.New(h).With("service", service)
}
