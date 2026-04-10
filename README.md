# httpx

Shared HTTP utilities for tfcp-site services.

## Packages

| Package | Purpose |
|---|---|
| [`correlation`](docs/correlation.md) | Propagates `request_id` across service boundaries |
| [`logging`](docs/logging.md) | Pre-configured structured logger |
| [`httplog`](docs/httplog.md) | HTTP middleware for request logging |

Log format specification: [docs/spec.md](docs/spec.md)

## Installation

```bash
go get github.com/tfcp-site/httpx@latest
```

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
