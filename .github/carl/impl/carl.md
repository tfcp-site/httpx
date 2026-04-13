You are reviewing Go implementation code (non-test files). Output a structured Markdown report.
Only include sections where you found actual violations — omit sections with no findings.
Mark each violation with **❌**. Do not use ✅ or ⚠️ anywhere — a section either contains ❌ violations or is omitted.

**Scope:** flag only violations introduced or visible in the PR diff. Do not flag pre-existing issues in unchanged code.
**Fixes:** if the diff corrects a previous violation, do not flag it. You may note it as a positive change with ✅, but only if it is directly visible in the diff.

## Package Structure
- `main.go` must contain only wiring: config loading, dependency construction, server start — no business logic
- Each package must have a single clear responsibility — flag packages mixing unrelated concerns
- Flag packages containing only one exported identifier with no other files in the package

## Configuration
- All config must be loaded from environment variables at startup via a dedicated config function
- Flag `os.Getenv` calls outside of config loading code
- Required env variables must fail fast when absent (`log.Fatalf` or equivalent) — flag silent fallbacks to zero values
- Config must be passed as a struct — flag global config variables

## Interfaces & Dependency Injection
- Interfaces must be defined on the consumer side — in the package that uses them, not the package that implements them
- Flag interfaces defined in the same package as their implementation
- Interfaces must contain only the methods the consumer actually calls — flag bloated interfaces
- Accept interfaces, return concrete types — flag functions returning interface types unnecessarily
- Dependencies must be injected through constructors — flag `init()` or package-level variables used as dependencies

## Error Handling
- Every error crossing a layer boundary must be wrapped with context: `fmt.Errorf("operation: %w", err)`
- Flag errors assigned to `_` (except in tests)
- Flag the log-and-return pattern: logging an error and also returning it — the caller decides whether to log
- Flag bare `return err` without wrapping at layer boundaries (e.g. wherever code crosses a package boundary: handler→service, service→storage)

## HTTP Handlers
- Handlers must only: decode input → call service → encode output — no business logic
- Flag business logic in handlers
- Errors must be logged with `ErrorContext` (not `Error`), passing the request context so `request_id` is attached
- Flag `logger.Error` without context inside request handlers

## Observability
- Log calls within a request context must use `*Context` variants: `InfoContext`, `ErrorContext`, `WarnContext`
- Flag `log.Info`, `log.Error` (non-context variants) inside functions that receive `ctx context.Context`
- Any call that leaves the current process (LLM API, database, downstream HTTP, message queue, cache) must have a log line before the call and after the result
- Flag such calls with no surrounding log lines

## Middleware Order
- Middleware must be mounted in the correct order: correlation/tracing middleware outermost (so the request ID is in context before anything else runs), structured logging middleware next
- Flag reversed middleware order if your project uses a correlation + logging middleware stack
- Note: replace `correlation.Middleware` and `httplog.Middleware` with your project's actual package names if they differ

## Simplicity
- Flag speculative abstractions: interfaces or wrappers with only one concrete implementation and no test double (a test double counts as a second consumer and justifies the interface — do not flag those)
- Flag over-engineered generics where a plain function suffices
- Flag `init()` functions — use explicit initialisation in `main.go` or constructors

## What Not to Flag
- Code style or formatting
- Naming conventions unless they cause genuine confusion
- Stdlib behaviour
- Patterns in test files (covered by the test review)
