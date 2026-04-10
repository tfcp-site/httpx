# httpx

Shared HTTP utilities for tfcp-site services.

## Packages

| Package | Purpose |
|---|---|
| [`correlation`](#correlation) | Propagates `request_id` across service boundaries |
| [`logging`](#logging) | Pre-configured structured logger |
| [`httplog`](#httplog) | HTTP middleware for request logging |

## Installation

```bash
go get github.com/tfcp-site/httpx@latest
```

Requires `GOPRIVATE` for the private module:

```bash
go env -w GOPRIVATE="github.com/tfcp-site/*"
```

---

## correlation

Ensures end-to-end request traceability across services using the `X-Request-ID` header.

**How it works:**

1. The middleware reads `X-Request-ID` from the incoming request. If absent, it generates a UUID.
2. The ID is stored in the request context and echoed back in the response header.
3. When making outgoing HTTP calls, `Transport` automatically injects the ID from context into the request header.

```
client → service-a → service-b → service-c
          X-Request-ID: f4a1-...  ─────────►
```

**Usage:**

```go
// Incoming requests — register the middleware
handler = correlation.Middleware(handler)

// Outgoing requests — wrap the HTTP client transport
client := &http.Client{
    Transport: &correlation.Transport{Base: http.DefaultTransport},
}
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
client.Do(req)
```

**API:**

| Symbol | Description |
|---|---|
| `Middleware(next) http.Handler` | Extracts or generates `request_id` and stores it in context |
| `Transport` | `http.RoundTripper` that injects `request_id` into outgoing requests |
| `FromContext(ctx) string` | Returns the `request_id` from context |
| `WithContext(ctx, id) context.Context` | Stores a `request_id` in context manually |
| `Extractor() func(context.Context) []slog.Attr` | Ready-to-use extractor for `logging.New` |
| `Header` | Header name: `"X-Request-ID"` |

---

## logging

Creates a `*slog.Logger` configured for tfcp-site conventions: JSON to stdout, no `time` field (Promtail records ingestion time), and a static `service` field on every line.

Supports dynamic log enrichment from context via `ContextExtractor`. Each package provides its own extractor — the service wires up the ones it needs at logger construction time.

**Usage:**

```go
logger := logging.New("mentor", slog.LevelInfo,
    correlation.Extractor(), // injects request_id from context
)

// request_id appears automatically on every log line
logger.InfoContext(ctx, "cache miss", "key", "chat:123")
// {"level":"INFO","service":"mentor","msg":"cache miss","request_id":"f4a1-...","key":"chat:123"}
```

**API:**

| Symbol | Description |
|---|---|
| `New(service, level, extractors...) *slog.Logger` | Creates a logger |
| `ContextExtractor` | `func(context.Context) []slog.Attr` — extractor function type |

---

## httplog

HTTP middleware that logs each incoming request as two structured lines: arrival and completion.

The `request_id` field is injected automatically via the logger (see `logging` + `correlation.Extractor`) — `httplog` has no dependency on `correlation`.

**Arrival** — logged immediately when the request is received:

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"f4a1-...","method":"POST","path":"/api/chat"}
```

**Completion** — logged after the handler returns:

```json
{"level":"INFO","service":"mentor","msg":"response","request_id":"f4a1-...","method":"POST","path":"/api/chat","status":200,"duration_ms":45}
```

For cached responses, `cached: true` is written instead of `duration_ms`.

**Usage:**

```go
// Mark a request as served from cache inside a handler
httplog.MarkCached(ctx)
```

**API:**

| Symbol | Description |
|---|---|
| `Middleware(log) func(http.Handler) http.Handler` | Logs request arrival and completion |
| `MarkCached(ctx)` | Signals to the middleware that the response was served from cache |

---

## Wiring

Typical `main.go`:

```go
logger := logging.New("mentor", slog.LevelInfo,
    correlation.Extractor(),
)

mux := http.NewServeMux()
mux.HandleFunc("/api/chat", chatHandler)

handler := correlation.Middleware(
    httplog.Middleware(logger)(mux),
)

http.ListenAndServe(":8080", handler)
```
