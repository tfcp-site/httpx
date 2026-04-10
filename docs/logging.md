# logging

Creates a `*slog.Logger` configured for tfcp-site conventions: JSON to stdout, no `time` field (Promtail records ingestion time), and a static `service` field on every line.

## Context enrichment

Supports dynamic log enrichment from context via `ContextExtractor`. Each package provides its own extractor — the service wires up the ones it needs at logger construction time.

When a log record is written with a context (e.g. `logger.InfoContext(ctx, ...)`), all registered extractors are called and their attributes are injected into the record automatically.

## Usage

```go
logger := logging.New("mentor", slog.LevelInfo,
    correlation.Extractor(), // injects request_id from context
)

// request_id appears automatically on every log line
logger.InfoContext(ctx, "cache miss", "key", "chat:123")
// {"level":"INFO","service":"mentor","msg":"cache miss","request_id":"f4a1-...","key":"chat:123"}
```

Adding more extractors as the service grows:

```go
logger := logging.New("mentor", slog.LevelInfo,
    correlation.Extractor(),
    auth.Extractor(), // injects user_id from context
)
```

## API

| Symbol | Description |
|---|---|
| `New(service, level, extractors...) *slog.Logger` | Creates a logger |
| `ContextExtractor` | `func(context.Context) []slog.Attr` — extractor function type |

## Providing an extractor from a package

Any package can expose an extractor without importing `logging`. Return `func(context.Context) []slog.Attr` — Go will implicitly convert it to `logging.ContextExtractor` at the call site:

```go
func Extractor() func(context.Context) []slog.Attr {
    return func(ctx context.Context) []slog.Attr {
        if val := fromContext(ctx); val != "" {
            return []slog.Attr{slog.String("my_field", val)}
        }
        return nil
    }
}
```
