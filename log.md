# Формат лога

## Стандартный набор полей

Порядок фиксирован для всех логов.

| # | Поле | Тип | Пример |
|---|------|-----|--------|
| 1 | `level` | string | `INFO`, `WARN`, `ERROR` |
| 2 | `service` | string | `"mentor"` |
| 3 | `msg` | string | `"request"` |
| 4 | _доп. поля_ | — | — |

`time` намеренно отсутствует — Promtail записывает время ingestion автоматически.

## Клиентский запрос

Все логи, порождённые в рамках обработки клиентского запроса, несут `request_id` — сквозной идентификатор, который проходит через всю цепочку: middleware → handler → сервисы → исходящие вызовы.

`request_id` идёт первым дополнительным полем (позиция 4).

## Типы логов

### HTTP

Два лога на каждый входящий запрос.

**Приход** (`msg: "request"`):

| Поле | Тип | Пример |
|------|-----|--------|
| `request_id` | string | `"550e8400-..."` |
| `method` | string | `"POST"` |
| `path` | string | `"/api/chat"` |

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"550e8400-...","method":"POST","path":"/api/chat"}
```

**Уход** (`msg: "response"`):

| Поле | Тип | Пример |
|------|-----|--------|
| `request_id` | string | `"550e8400-..."` |
| `method` | string | `"POST"` |
| `path` | string | `"/api/chat"` |
| `status` | int | `200` |
| `duration` | string | `"1.5s"`, `"150ms"` |
| `cached` | bool | `true` |

`duration` и `cached` взаимоисключающие: обработан → `duration`, из кеша → `cached: true`.

```json
{"level":"INFO","service":"mentor","msg":"response","request_id":"550e8400-...","method":"POST","path":"/api/chat","status":200,"duration":"1.240s"}
```

### App

Произвольные логи приложения. Дополнительные поля — по смыслу конкретного лога.

**В рамках клиентского запроса** — несут `request_id`:

```json
{"level":"INFO","service":"mentor","msg":"request","request_id":"550e8400-...","method":"POST","path":"/api/chat"}
{"level":"INFO","service":"mentor","msg":"cache miss","request_id":"550e8400-...","key":"chat:abc"}
{"level":"INFO","service":"mentor","msg":"llm response","request_id":"550e8400-...","tokens":312}
{"level":"INFO","service":"mentor","msg":"response","request_id":"550e8400-...","method":"POST","path":"/api/chat","status":200,"duration":"1.240s"}
```

**Вне клиентского запроса** (startup, background job и т.п.) — без `request_id`:

```json
{"level":"INFO","service":"mentor","msg":"server started","addr":":8080"}
{"level":"INFO","service":"mentor","msg":"cleanup done","job":"cleanup","deleted":42}
```
