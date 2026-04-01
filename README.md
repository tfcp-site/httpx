# httpx

Shared HTTP utilities for tfcp-site services.

## Packages

### `correlation`

Propagates a correlation ID (`X-Request-ID`) across service boundaries so that a single request can be traced through multiple services by querying logs.

**How it works:**

1. The server middleware reads `X-Request-ID` from the incoming request header. If absent, it generates a UUID and stores the ID in the request context.
2. When the service makes outgoing HTTP calls, `Transport` reads the ID from the context and injects it into the `X-Request-ID` header automatically.
3. Each service logs the ID as a structured field. All log lines across all services share the same ID for a given request.

```
client → service-a → service-b → service-c
          X-Request-ID: f4a1-...  ─────────►
```

## Installation

```bash
go get github.com/tfcp-site/httpx@latest
```

Requires `GOPRIVATE` to be set since this is a private module:

```bash
go env -w GOPRIVATE="github.com/tfcp-site/*"
```

## Usage

**Server — register middleware**

```go
import "github.com/tfcp-site/httpx/correlation"

handler = correlation.Middleware(handler)
```

**Client — propagate ID to downstream services**

```go
client := &http.Client{
    Transport: &correlation.Transport{Base: http.DefaultTransport},
}

// The ID is read from ctx and set as X-Request-ID automatically.
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
client.Do(req)
```

**Logger — include ID in log fields**

```go
slog.InfoContext(ctx, "processing request",
    "request_id", correlation.FromContext(ctx),
)
```

**Manual context injection** (useful in tests or background jobs)

```go
ctx = correlation.WithContext(ctx, "my-custom-id")
```

## API

| Symbol | Description |
|---|---|
| `Middleware(next http.Handler) http.Handler` | Server middleware — extracts or generates the correlation ID |
| `Transport` | `http.RoundTripper` — injects the ID from context into outgoing requests |
| `FromContext(ctx) string` | Returns the correlation ID stored in the context |
| `WithContext(ctx, id) context.Context` | Returns a new context carrying the given ID |
| `Header` | The header name: `"X-Request-ID"` |
