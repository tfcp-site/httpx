# Log Format Specification

## Standard fields

Field order is fixed across all log lines.

| # | Field | Type | Example |
|---|-------|------|---------|
| 1 | `level` | string | `INFO`, `WARN`, `ERROR` |
| 2 | `service` | string | `"mentor"` |
| 3 | `msg` | string | `"request"` |
| 4 | _additional fields_ | — | — |

`time` is intentionally absent — Promtail records ingestion time automatically.

## Client request context

All log lines produced during a client request carry `request_id` — a correlation identifier that flows through the entire chain: middleware → handler → services → outgoing calls.

`request_id` is always the first additional field (position 4).

## Log types

### HTTP

Two log lines per incoming request.

**Arrival** (`msg: "request"`):

| Field | Type | Example |
|-------|------|---------|
| `request_id` | string | `"550e8400-..."` |
| `method` | string | `"POST"` |
| `path` | string | `"/api/chat"` |

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"550e8400-...","method":"POST","path":"/api/chat"}
```

**Completion** (`msg: "response"`):

| Field | Type | Example |
|-------|------|---------|
| `request_id` | string | `"550e8400-..."` |
| `method` | string | `"POST"` |
| `path` | string | `"/api/chat"` |
| `status` | int | `200` |
| `duration` | string | `"1.5s"`, `"150ms"` |
| `cached` | bool | `true` |

`duration` and `cached` are mutually exclusive: normal response → `duration`, cache hit → `cached: true`.

```json
{"level":"INFO","service":"mentor","msg":"response","request_id":"550e8400-...","method":"POST","path":"/api/chat","status":200,"duration":"1.240s"}
```

### App

Arbitrary application logs. Additional fields are defined per log site.

**Within a client request** — carry `request_id`:

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"550e8400-...","method":"POST","path":"/api/chat"}
{"level":"INFO","service":"mentor","msg":"cache miss","request_id":"550e8400-...","key":"chat:abc"}
{"level":"INFO","service":"mentor","msg":"llm response","request_id":"550e8400-...","tokens":312}
{"level":"INFO","service":"mentor","msg":"response","request_id":"550e8400-...","method":"POST","path":"/api/chat","status":200,"duration":"1.240s"}
```

**Outside a client request** (startup, background jobs, etc.) — no `request_id`:

```json
{"level":"INFO","service":"mentor","msg":"server started","addr":":8080"}
{"level":"INFO","service":"mentor","msg":"cleanup done","job":"cleanup","deleted":42}
```
