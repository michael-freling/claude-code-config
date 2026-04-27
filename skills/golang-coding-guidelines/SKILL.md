---
name: golang-coding-guidelines
description: Go coding guidelines covering language idioms, error handling, concurrency, and application patterns for HTTP servers, HTTP clients, gRPC, databases, message queues, and CLI tools
---

# Go Coding Guidelines

## Application Pattern Examples

- [HTTP Server](examples/http-server.md) — `net/http` handlers, middleware, routing, request/response patterns
- [HTTP Client](examples/http-client.md) — outbound HTTP calls with retries, timeouts, connection pooling
- [gRPC Server](examples/grpc-server.md) — protobuf services, interceptors, streaming
- [gRPC Client](examples/grpc-client.md) — dial options, interceptors, retries
- [Database](examples/db.md) — `database/sql` with connection pooling, transactions, migrations
- [Message Queue Worker](examples/message-queue.md) — consumer loops, acknowledgment, graceful shutdown
- [CLI Tool](examples/cli.md) — `cobra` commands, flag parsing, exit codes

## Go Idioms

**Errors:**
- Every error must be checked or returned — never use `_` to discard errors
- Wrap errors with context using `%w` — include the operation and identifying input:
  - Wrong: `return fmt.Errorf("failed: %w", err)`
  - Right: `return fmt.Errorf("find user by ID %d: %w", userID, err)`
- Use `errors.Is()` to check for specific error values, `errors.As()` to extract error types — never compare error strings
- Define sentinel errors with `var ErrNotFound = errors.New("not found")` — use `errors.Is(err, ErrNotFound)` at call sites
- Define custom error types when callers need to extract structured data: implement the `error` interface and unwrap with `errors.As()`

**Type Safety:**
- Use `iota` for enums — define a named type and a `String()` method:
  ```go
  type Status int
  const (
      StatusPending Status = iota
      StatusActive
      StatusClosed
  )
  ```
- Use type aliases sparingly — prefer new named types for domain concepts (`type UserID int64`) to prevent mixing IDs from different domains

**Code Organization:**
- Package names: short, lowercase, no underscores — `user`, not `user_service` or `userService`
- Use `internal/` packages to restrict visibility to the parent module
- One package per directory — no shared directories between packages
- Avoid `package models` or `package types` — define types in the package that uses them
- `cmd/<binary>/main.go` for each executable; `internal/` for private packages; top-level packages for the public API

**Initialization:**
- Functional options for optional constructor parameters:
  ```go
  type Option func(*config)
  func WithTimeout(d time.Duration) Option {
      return func(c *config) { c.timeout = d }
  }
  func NewClient(baseURL string, opts ...Option) *Client { ... }
  ```
- Do not use functional options when all parameters are required — use plain arguments
- `init()` is only acceptable for registering drivers (`database/sql`, `image`) — never for business logic or I/O

**Concurrency:**
- Prefer `golang.org/x/sync/errgroup` over `sync.WaitGroup` — built-in error propagation, context cancellation via `errgroup.WithContext`, and cleaner API with `eg.Go()`
- Pass `context.Context` as the first parameter to any function that does I/O or may block
- Never store `context.Context` in a struct — pass it as a function parameter
- Use channels for communication between goroutines, `sync.Mutex` for protecting shared state
- Always select on `ctx.Done()` in goroutines that loop or block

**Testing:**
- Use want/got ordering (never expected/actual) in assertion messages: `"got %v, want %v"`
- Use **`assert`** when a test can continue after failure, **`require`** when it must stop
- Use **`go.uber.org/gomock`** for generating mock implementations of interfaces
- Use `t.Helper()` in test helper functions so failures report the caller's line
- Use `t.Cleanup()` for teardown instead of `defer` in test functions — works correctly with subtests
- Use `testcontainers-go` for integration tests with real databases and services

## Principles

- **Environment-agnostic code** — use `os.Getenv`, config structs, or `flag` for environment-specific values; never branch on `os.Getenv("ENV") == "production"` in application code; inject different configurations via dependency injection

## Simplicity

**Simplicity is the most important thing.** Minimize the final complexity of the codebase — not the size of the change.

- **Breaking changes are acceptable** unless backward compatibility is explicitly required by a published module's semver contract or project constraints
- **Avoid premature abstraction** — do not define an interface until there are 2+ implementations or a concrete testing need; three similar lines are better than a premature `Doer` interface
- **Minimize mutable state** — prefer stateless functions, derived values, and passing values over pointers-to-pointers; keep `sync.Mutex`-protected state in the smallest possible struct

### High Cohesion

Each function should do one thing. When a new requirement applies to only some callers, split the function rather than adding conditional logic.

- **Wrong**: `SendNotification()` handles email, SMS, and push — adding rate limiting for SMS forces all paths through the rate limiter
- **Right**: `SendEmail()`, `SendSMS()`, `SendPush()` — rate limiting lives only in `SendSMS()`

### Separate Methods over Boolean Arguments

When a `bool` argument causes a function to take substantially different paths, split into separate functions.

- **Wrong**: `CompleteJob(success bool)` with branching logic based on the flag
- **Right**: `CompleteJob()` and `FailJob()` as separate exported methods

### Narrow Interfaces

Accept the minimum interface a function needs — Go's implicit interface satisfaction makes this natural.

- **Wrong**: a function that only reads accepts `*os.File`
- **Right**: the function accepts `io.Reader`
- Define small interfaces in the consumer package: `type Reader interface { Read(p []byte) (n int, err error) }`

## Error Handling

- **Always handle errors explicitly** — every `if err != nil` must either return the error (wrapped) or handle it; never `_ = someFunc()`
- **Wrap errors with `fmt.Errorf("operation %s: %w", id, err)`** — include what operation failed and with what input
- **Fail fast on invalid state** — validate inputs at entry points and return errors immediately with `if err != nil { return }` guard clauses
- **Use sentinel errors** for expected conditions: `var ErrNotFound = errors.New("not found")`; check with `errors.Is(err, ErrNotFound)`
- **Use custom error types** when callers need structured data from the error (status codes, retry hints)

### Server-Side Errors

- Map domain errors to HTTP/gRPC status codes — `ErrNotFound` to `404`/`codes.NotFound`, `ErrInvalidInput` to `400`/`codes.InvalidArgument`, unexpected errors to `500`/`codes.Internal`
- Return user-friendly error messages via response bodies; never expose stack traces or internal error strings to clients
- Log the full error with `slog.ErrorContext(ctx, "operation", "error", err)` for server-side diagnosis

### Client-Side Errors

- Retry on `429`, `502`, `503`, `504` and network timeouts with exponential backoff — use `golang.org/x/time/rate` or a retry library
- Make retry count, backoff base, and retryable status codes configurable via functional options or config struct
- If the infrastructure layer (service mesh, load balancer) already retries, do not duplicate

## Dependencies

- Run `go get <module>@latest` to add dependencies — Go modules pin the exact version in `go.sum`
- Use `go mod tidy` to remove unused dependencies
- Pin tool versions with `//go:generate go run <tool>@v<version>` or a `tools.go` file

## Code Readability

- **Specific names** — `userCount` not `n`, `httpClient` not `c`. Exception: loop variables (`i`, `k`, `v`) and receivers (`s`, `h`) scoped to a few lines
- **Receiver names** — one or two letters derived from the type name, consistent across all methods: `func (s *Server) Handle(...)`, not `func (srv *Server) Handle(...)`
- **Early returns** — guard clauses with `if err != nil { return }` at the top; keep the happy path at the base indentation
- **Single level of abstraction** — a handler function calls service methods, it does not also write SQL
- **Delete dead code** — remove unused functions, unreachable branches, commented-out blocks, and unused imports (Go enforces unused imports, but watch for unused exported symbols)
- **No zero-value assignments** — `var count int` is already `0`; do not write `count := 0` unless it clarifies intent in a reset context
- **Avoid long parameter lists** — more than 5 parameters or 3 return values signal the need for an options struct or a result struct:
  ```go
  type CreateOrderParams struct {
      UserID    int64
      Items     []Item
      CouponID  *string
      Note      string
  }
  ```
- **Load binary data from files** — use `os.ReadFile` or `//go:embed` for certificates, images, and fixtures; never inline base64 strings
- **Validate only at boundaries** — validate HTTP/gRPC request payloads and external API responses; do not add validation in internal functions or for records your own application manages

## Code Reuse

- **Single definition, owned by the provider** — define interfaces, types, and functions in the package that implements them; consumers import, not redefine
- **Named constants** — `const maxRetries = 3` instead of magic `3` scattered across the code
- **Search before writing** — `grep -r` for existing helpers, types, or constants before creating new ones
- **Refactor before extending** — understand the existing code, extract reusable parts, then add the new feature
- **Refactor after QA** — once the feature works, clean up for reusability and consistency

## Code Organization

- **Minimal `main` package** — `func main()` initializes config, database, clients, wires dependencies, and calls `Run(ctx)` or `ListenAndServe`; no business logic in `main`
- **Boundary layers by type** — `handler/` for HTTP/gRPC, `client/` for outbound SDKs; **domain logic by feature** — `user/`, `billing/`, `order/`
- **Dependency direction** — `handler` imports `user`, `user` imports `db` interface — never the reverse; Go's compiler enforces no circular imports
- **Minimize exports** — unexported by default; export only what other packages consume; use `internal/` to enforce

## Type Safety

- **Named types for domain IDs** — `type UserID int64`, `type OrderID string` to prevent accidentally passing a `UserID` where an `OrderID` is expected
- **`iota` enums with a `String()` method** for known sets of values; use `go generate` with `stringer` to auto-generate
- **Avoid `interface{}`/`any`** unless writing generic library code — prefer concrete types or type parameters (generics)

## Initialization

- **Required parameters as positional arguments**, optional parameters via functional options:
  ```go
  func NewServer(addr string, handler http.Handler, opts ...Option) *Server
  ```
- **Validate invariants in the constructor** — return an error if required fields are empty or invalid:
  ```go
  func NewClient(baseURL string) (*Client, error) {
      if baseURL == "" {
          return nil, errors.New("baseURL is required")
      }
      ...
  }
  ```

## External I/O

- **Batch operations** — use `sqlx.In()` for bulk `INSERT`/`SELECT`, batch gRPC calls, or bulk API endpoints to reduce round-trips
- **Transactions** — `db.BeginTx(ctx, nil)` for multi-table writes within a single request; `defer tx.Rollback()` immediately after `BeginTx` (no-op after `Commit`)
- **Ensure indexes** — add indexes for columns used in `WHERE`, `JOIN`, and `ORDER BY` on tables with 10K+ rows
- **Pass values from Go code** — compute `time.Now()` in Go and pass it as a parameter; do not use `NOW()` in SQL — enables deterministic tests with injected clocks
- **Avoid redundant queries** — do not `SELECT` a row you just `INSERT`ed if you already have the data; do not query the same record twice in one handler
- **Parameterized queries** — always use `?` or `$1` placeholders; never `fmt.Sprintf` user input into SQL strings

## Logging

- **`log/slog`** as the standard structured logger — use `slog.With("key", value)` for context, not `fmt.Sprintf`
- **Log levels**: `slog.Debug` (development traces), `slog.Info` (operational events), `slog.Warn` (expected/recoverable issues), `slog.Error` (unexpected failures needing investigation)
- **Correlation IDs** — extract trace/span IDs from `ctx` and include in log entries; use `slog.ErrorContext(ctx, ...)` to propagate context automatically
- **Never log secrets** — no tokens, passwords, API keys, PII, or full request/response bodies; redact sensitive fields before logging

## Testing

- **Deterministic tests** — no `time.Sleep`, no reliance on wall-clock ordering; inject `time.Time` or `func() time.Time` for time-dependent logic
- **Inject dependencies** — pass interfaces via constructors; use `gomock` to generate mocks for external service interfaces
- **`t.TempDir()`** for file-based tests — automatically cleaned up
- **Table-driven tests** — group related cases in `[]struct{ name string; ... }` slices; separate success cases from error cases:
  ```go
  tests := []struct {
      name    string
      input   string
      want    int
      wantErr bool
  }{...}
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) { ... })
  }
  ```
- **Real databases over mocks** — use `testcontainers-go` to spin up PostgreSQL/MySQL in tests; mock only third-party HTTP APIs with `httptest.NewServer`
- **Shared test setup** — define factory functions that return populated structs with sensible defaults; share `*sql.DB` and mock servers across subtests via a test helper struct
- **Never `t.Skip()` failures** — fix flaky tests immediately
- **Load test data from files** — use `testdata/` directory and `os.ReadFile("testdata/fixture.json")`

### What to Test

- All exported functions and methods
- Error paths — each `if err != nil` branch
- Edge cases — empty slices, nil pointers, zero values, max `int64`
- Boundary conditions — off-by-one, pagination limits, timeout edges

### What NOT to Test

- Standard library behavior (`json.Marshal` correctness, `http.StatusOK` value)
- Trivial getters that return a struct field
- Generated code (`protoc` output, `gomock` mocks)

### Coverage Gaps to Look For

- Untested error branches after I/O calls
- Missing edge cases for slice/map operations
- Conditional logic (`switch`, `if/else`) not fully covered

### Writing Testable Code

1. **Analyze testability** — identify concrete dependencies (`*sql.DB`, HTTP clients), global state, `init()` side effects
2. **Refactor first** — extract interfaces at the boundary, inject via constructor, split large functions
3. **Write tests** — table-driven, success and error cases, using `assert`/`require`
4. **Verify coverage** — `go test -coverprofile=coverage.out ./...` and check with `go tool cover -html=coverage.out`

## Concurrency

- **`sync.Mutex`** for protecting shared mutable state; **`sync.RWMutex`** when reads vastly outnumber writes; **`sync/atomic`** for simple counters
- **`errgroup.WithContext`** for fan-out — cancels remaining goroutines on first error:
  ```go
  eg, ctx := errgroup.WithContext(ctx)
  for _, item := range items {
      eg.Go(func() error { return process(ctx, item) })
  }
  if err := eg.Wait(); err != nil { ... }
  ```
- **Channels for coordination** — use buffered channels for producer/consumer patterns, unbuffered for synchronization points
- **`context.WithCancel`/`context.WithTimeout`** to propagate cancellation — every goroutine should select on `ctx.Done()`

## Security

- **No hardcoded secrets** — use `os.Getenv` or a secret manager (`hashicorp/vault`, cloud KMS); never commit `.env` files with real credentials
- **Sanitize at trust boundaries** — validate and sanitize all user input in HTTP/gRPC handlers; use parameterized SQL; escape output for HTML with `html/template`

## Monitoring

- **OpenTelemetry SDK** — initialize `TracerProvider`, `MeterProvider`, and `LoggerProvider` in `main`; use `otel.Tracer("pkg")` and `otel.Meter("pkg")` in application code
- **RED metrics** at minimum for every server:
  - **Rate**: requests per second (counter)
  - **Errors**: error responses per second (counter with status label)
  - **Duration**: request latency histogram
- **Avoid high-cardinality labels** — do not use user IDs, request IDs, or unbounded strings as metric labels; use status codes, method names, and endpoint patterns
