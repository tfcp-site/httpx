# httpx

Shared HTTP utilities for tfcp-site services.

## Packages

| Package | Purpose | Docs |
|---|---|---|
| `correlation` | Propagates `request_id` across service boundaries | [docs/correlation.md](docs/correlation.md) |
| `logging` | Pre-configured structured logger | [docs/logging.md](docs/logging.md) |
| `httplog` | HTTP middleware for request logging | [docs/httplog.md](docs/httplog.md) |

Log format specification: [docs/spec.md](docs/spec.md)

## Installation

```bash
go get github.com/tfcp-site/httpx@latest
```

Requires `GOPRIVATE` for the private module:

```bash
go env -w GOPRIVATE="github.com/tfcp-site/*"
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
