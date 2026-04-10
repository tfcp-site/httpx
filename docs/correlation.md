# correlation

Ensures end-to-end request traceability across services using the `X-Request-ID` header.

## How it works

1. The middleware reads `X-Request-ID` from the incoming request. If absent, it generates a UUID.
2. The ID is stored in the request context and echoed back in the response header.
3. When making outgoing HTTP calls, `Transport` automatically injects the ID from context into the request header.

```
client → service-a → service-b → service-c
          X-Request-ID: f4a1-...  ─────────►
```

## Usage

**Incoming requests — register the middleware:**

```go
handler = correlation.Middleware(handler)
```

**Outgoing requests — wrap the HTTP client transport:**

```go
client := &http.Client{
    Transport: &correlation.Transport{Base: http.DefaultTransport},
}

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
client.Do(req)
```

**Inject `request_id` into logs automatically:**

```go
logger := logging.New("mentor", slog.LevelInfo,
    correlation.Extractor(),
)
```

**Manual context injection** (useful in tests or background jobs):

```go
ctx = correlation.WithContext(ctx, "my-custom-id")
```

## API

| Symbol | Description |
|---|---|
| `Middleware(next) http.Handler` | Extracts or generates `request_id` and stores it in context |
| `Transport` | `http.RoundTripper` that injects `request_id` into outgoing requests |
| `FromContext(ctx) string` | Returns the `request_id` from context |
| `WithContext(ctx, id) context.Context` | Stores a `request_id` in context manually |
| `Extractor() func(context.Context) []slog.Attr` | Ready-to-use extractor for `logging.New` |
| `Header` | Header name: `"X-Request-ID"` |
