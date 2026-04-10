# httplog

HTTP middleware that logs each incoming request as two structured lines: arrival and completion.

`request_id` and other context attributes are injected automatically by the logger — `httplog` has no dependency on `correlation`. See [logging.md](logging.md) and [correlation.md](correlation.md).

## Log output

**Arrival** — logged immediately when the request is received (`msg: "request"`):

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"f4a1-...","method":"POST","path":"/api/chat"}
```

**Completion** — logged after the handler returns (`msg: "response"`):

```json
{"level":"INFO","service":"mentor","msg":"response","request_id":"f4a1-...","method":"POST","path":"/api/chat","status":200,"duration":"45ms"}
```

For cached responses, `cached: true` is written instead of `duration`.

See [spec.md](spec.md) for the full log format specification.

## Usage

```go
handler := httplog.Middleware(logger)(mux)
```

**Mark a response as served from cache inside a handler:**

```go
func articleHandler(w http.ResponseWriter, r *http.Request) {
    if hit, val := cache.Get(r.URL.Path); hit {
        httplog.MarkCached(r.Context())
        w.Write(val)
        return
    }
    // ...
}
```

## API

| Symbol | Description |
|---|---|
| `Middleware(log) func(http.Handler) http.Handler` | Logs request arrival and completion |
| `MarkCached(ctx)` | Signals to the middleware that the response was served from cache |
