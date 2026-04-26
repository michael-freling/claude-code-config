---
name: golang-coding-guidelines
description: Go coding guidelines covering language idioms, error handling, concurrency, and application patterns for HTTP servers, HTTP clients, gRPC, databases, message queues, and CLI tools
---

# Go Coding Guidelines

## Example Files

- `examples/http-server.md` — HTTP handler patterns with chi/stdlib, middleware, request validation, error responses
- `examples/http-client.md` — HTTP client with context, retries, typed responses, connection pooling
- `examples/grpc-server.md` — gRPC service implementation, interceptors, error mapping
- `examples/grpc-client.md` — gRPC client with retries, deadlines, connection management
- `examples/db.md` — database/sql and sqlx patterns, connection pool, transactions, batch queries
- `examples/message-queue.md` — Message queue consumer/producer patterns with graceful shutdown
- `examples/cli.md` — CLI tool patterns with cobra, configuration, graceful shutdown

## Go Language Idioms

### Errors

- Every error must be checked or returned — never use `_` to discard errors
- Wrap errors with context using `%w` — include what operation failed and with what input:
  - Wrong: `return fmt.Errorf("failed: %w", err)`
  - Right: `return fmt.Errorf("find user by ID %d: %w", userID, err)`
- Use `errors.Is()` to check for specific error values, `errors.As()` to extract error types — never compare error strings
- Define sentinel errors with `var ErrNotFound = errors.New("not found")` at package level
- Custom error types implement `error` interface:
  ```go
  type ValidationError struct {
      Field   string
      Message string
  }
  func (e *ValidationError) Error() string {
      return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
  }
  ```

### Type Safety

- Use `iota` for enums:
  ```go
  type OrderStatus int
  const (
      OrderStatusPending OrderStatus = iota + 1
      OrderStatusShipped
      OrderStatusCancelled
  )
  ```
- Implement `String()` via `stringer` or manually for readable output
- Use named types to prevent argument swaps: `type UserID int64`, `type OrderID int64`

### Code Organization

- Package names: short, lowercase, no underscores (e.g., `user`, not `user_service`)
- Use `internal/` packages to restrict visibility to the parent module
- `cmd/<binary>/main.go` for entry points, `internal/` for business logic:
  ```
  cmd/
    server/main.go
    worker/main.go
  internal/
    user/
      service.go
      repository.go
    billing/
      service.go
      invoice.go
  pkg/           # only for genuinely reusable library code
  ```
- One package per domain — not `services/`, `repositories/`, `models/`
- Boundary code (HTTP handlers, gRPC handlers, DB clients) in separate packages from domain logic

### Initialization

- Functional options for optional constructor parameters:
  ```go
  type Option func(*options)

  type options struct {
      timeout time.Duration
      logger  *slog.Logger
  }

  func WithTimeout(d time.Duration) Option {
      return func(o *options) { o.timeout = d }
  }

  func WithLogger(l *slog.Logger) Option {
      return func(o *options) { o.logger = l }
  }

  func NewClient(baseURL string, opts ...Option) *Client {
      o := options{
          timeout: 10 * time.Second,
          logger:  slog.Default(),
      }
      for _, opt := range opts {
          opt(&o)
      }
      return &Client{baseURL: baseURL, timeout: o.timeout, logger: o.logger}
  }
  ```
- Do not use functional options when all parameters are required — use plain constructor arguments

### Concurrency

- Prefer `golang.org/x/sync/errgroup` over `sync.WaitGroup` — built-in error propagation, context cancellation via `errgroup.WithContext`, cleaner API with `eg.Go()`:
  ```go
  eg, ctx := errgroup.WithContext(ctx)
  eg.Go(func() error {
      return fetchUser(ctx, userID)
  })
  eg.Go(func() error {
      return fetchOrders(ctx, userID)
  })
  if err := eg.Wait(); err != nil {
      return fmt.Errorf("fetch user data: %w", err)
  }
  ```
- `eg.SetLimit(n)` to bound concurrent goroutines
- `sync.Mutex` for shared mutable state — keep the critical section minimal
- Channels for communication between goroutines, mutexes for protecting shared state

### Testing

- Use want/got (never expected/actual):
  ```go
  if got := result; got != want {
      t.Errorf("GetUser() = %v, want %v", got, want)
  }
  ```
- **assert** when test can continue, **require** when it should stop:
  ```go
  require.NoError(t, err)
  assert.Equal(t, want, got)
  ```
- **go.uber.org/gomock** for mocks:
  ```go
  ctrl := gomock.NewController(t)
  mock := NewMockUserRepository(ctrl)
  mock.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
  ```
- `t.Run()` for subtests, `t.Parallel()` for parallel execution
- `t.TempDir()` for temporary files — auto-cleaned
- `t.Cleanup()` for teardown — runs in LIFO order

### Dependencies

- `go.sum` provides reproducible builds by default
- Use `go get package@latest` to install latest; pin with `go get package@v1.2.3`
- `go mod tidy` to clean unused dependencies

### Struct Tags

- Use `json:"field_name"` for JSON serialization, `db:"column_name"` for database mapping
- `json:"-"` to exclude fields, `json:"field,omitempty"` to omit zero values

---

## Principles

- DI via constructor — parse config once at startup, pass typed config struct:
  ```go
  type Config struct {
      DatabaseURL string
      LogLevel    slog.Level
      Port        int
  }

  func LoadConfig() (*Config, error) {
      cfg := &Config{
          DatabaseURL: os.Getenv("DATABASE_URL"),
          Port:        8080,
      }
      if cfg.DatabaseURL == "" {
          return nil, errors.New("DATABASE_URL is required")
      }
      return cfg, nil
  }
  ```
- Inject dependencies through constructor functions — test by passing fakes, not by patching globals
- No environment-specific branches in application code — use configuration

## Simplicity

- Prefer value types and stateless functions — minimize mutable state:
  ```go
  // Wrong — mutable accumulation
  var names []string
  for _, u := range users {
      if u.Active {
          names = append(names, u.Name)
      }
  }

  // Right for Go — this is idiomatic; use slices.Collect for Go 1.23+
  names := make([]string, 0, len(users))
  for _, u := range users {
      if u.Active {
          names = append(names, u.Name)
      }
  }
  ```
- Three similar functions are better than a premature generic — use generics only when the abstraction is proven across 3+ concrete types
- Breaking changes are acceptable unless backward compatibility is explicitly required

### High Cohesion

- Wrong: `sendNotification()` handles email/SMS/push — right: `sendEmail()`, `sendSMS()`, `sendPush()`
- Each function does one thing; split rather than add conditional branches

### Separate Methods over Boolean Arguments

- Wrong: `CompleteJob(success bool)` — right: `CompleteJob()` and `FailJob()`

### Narrow Interfaces

- Accept the minimum interface:
  ```go
  // Wrong — accepts *os.File when only reading
  func processData(f *os.File) error

  // Right — accepts io.Reader
  func processData(r io.Reader) error
  ```
- Define interfaces in the consumer package, not the provider:
  ```go
  // In package handler (consumer), not package user (provider)
  type UserFinder interface {
      FindByID(ctx context.Context, id int64) (*user.User, error)
  }
  ```

## Error Handling

- Always handle errors explicitly — every `if err != nil` block must wrap and return or handle:
  ```go
  user, err := s.repo.FindByID(ctx, id)
  if err != nil {
      return nil, fmt.Errorf("find user by ID %d: %w", id, err)
  }
  ```
- Fail fast on invalid state — validate at the start, return early:
  ```go
  func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
      if len(req.Items) == 0 {
          return nil, &ValidationError{Field: "items", Message: "at least one item required"}
      }
      // ... main logic at base indentation
  }
  ```
- Server: return appropriate status codes — 4XX for client errors, 5XX for server errors. Include user-friendly messages. Never expose stack traces or internal error details:
  ```go
  if errors.Is(err, user.ErrNotFound) {
      http.Error(w, "user not found", http.StatusNotFound)
      return
  }
  slog.ErrorContext(ctx, "create user", "error", err)
  http.Error(w, "internal server error", http.StatusInternalServerError)
  ```
- Client: retry transient failures with exponential backoff (see `examples/http-client.md`)

## Dependencies

- `go get package@latest` to install the latest stable version
- `go.sum` provides build reproducibility by default
- `go mod tidy` to remove unused dependencies
- Check `go vet` and `staticcheck` for correctness issues

## Code Readability

- Specific names: `userRepo` not `repo`, `OrderStatus` not `Status`. Short names (`i`, `r`, `w`) are fine within ~5 lines of scope
- Early returns with guard clauses:
  ```go
  func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
      if id <= 0 {
          return nil, &ValidationError{Field: "id", Message: "must be positive"}
      }
      user, err := s.repo.FindByID(ctx, id)
      if err != nil {
          return nil, fmt.Errorf("find user by ID %d: %w", id, err)
      }
      return user, nil
  }
  ```
- Use structs when a function has more than 5 parameters or more than 3 return values:
  ```go
  type CreateUserParams struct {
      Name     string
      Email    string
      Role     Role
      TeamID   int64
      Settings UserSettings
      Verified bool
  }
  func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
  ```
- Delete dead code — unused functions, unreachable branches, commented-out code
- Delete assignments of zero values — Go initializes to zero values automatically
- Load binary data from files with `os.ReadFile()` — never embed base64 inline
- Validate only at boundaries — incoming HTTP/gRPC requests and outgoing 3rd-party API responses. Trust internal types

## Code Reuse

- Single definition, owned by the provider — types defined in provider package, consumers import
- Named constants over magic values:
  ```go
  const (
      maxRetries    = 3
      retryInterval = 500 * time.Millisecond
  )
  ```
- Search for existing code before adding new — reuse or extend existing patterns
- Refactor before adding new changes; refactor again after QA passes

## Code Organization

- `cmd/<binary>/main.go` only initializes — DB connections, clients, config — and starts the server/worker:
  ```go
  func main() {
      cfg, err := config.Load()
      if err != nil {
          log.Fatal(err)
      }
      db, err := sql.Open("postgres", cfg.DatabaseURL)
      if err != nil {
          log.Fatal(err)
      }
      defer db.Close()

      repo := user.NewRepository(db)
      svc := user.NewService(repo)
      handler := api.NewHandler(svc)

      srv := &http.Server{Addr: cfg.Addr, Handler: handler.Routes()}
      log.Fatal(srv.ListenAndServe())
  }
  ```
- Boundary layers by type: `internal/api/` (HTTP), `internal/grpcapi/` (gRPC), `internal/store/` (DB)
- Domain logic by feature: `internal/user/`, `internal/billing/` — not `internal/services/`, `internal/repositories/`
- Dependency direction: handlers → services → repositories. Never the reverse
- Minimize exports — unexported by default, export only what consumers need

## Type Safety

- `iota` enums with `String()` method:
  ```go
  type Status int
  const (
      StatusPending Status = iota + 1
      StatusActive
      StatusDisabled
  )
  //go:generate stringer -type=Status
  ```
- Named types prevent argument swaps:
  ```go
  type UserID int64
  type OrderID int64
  func FindOrder(userID UserID, orderID OrderID) // compiler catches FindOrder(orderID, userID)
  ```

## Initialization

- Required positional, optional via functional options:
  ```go
  func NewService(repo Repository, opts ...Option) *Service
  ```
- Validate invariants in the constructor:
  ```go
  func NewService(repo Repository) (*Service, error) {
      if repo == nil {
          return nil, errors.New("repo is required")
      }
      return &Service{repo: repo}, nil
  }
  ```
- Do not use functional options when all parameters are required

## External I/O

- Batch operations — `IN` clauses over N+1 loops:
  ```go
  // Wrong — N+1 queries
  for _, id := range ids {
      user, err := repo.FindByID(ctx, id)
  }

  // Right — single query
  users, err := repo.FindByIDs(ctx, ids)
  ```
- One transaction per request for multi-table writes (see `examples/db.md`)
- Pass values from application code, not SQL functions — `$1` not `NOW()`:
  ```go
  _, err := db.ExecContext(ctx, `INSERT INTO orders (created_at) VALUES ($1)`, time.Now())
  ```
- Parameterized queries — never concatenate user input into SQL:
  ```go
  // Wrong
  query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)

  // Right
  rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE name = $1", name)
  ```
- Configure connection pool for `database/sql`:
  ```go
  db.SetConnMaxLifetime(5 * time.Minute)
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(10)
  ```

## Logging

- `log/slog` structured logging (stdlib, Go 1.21+):
  ```go
  slog.InfoContext(ctx, "order placed",
      "user_id", userID,
      "order_id", orderID,
      "item_count", len(order.Items),
  )
  ```
- Levels: `Debug` (dev detail), `Info` (operational events), `Warn` (recoverable), `Error` (needs investigation)
- Include trace context:
  ```go
  span := trace.SpanFromContext(ctx)
  slog.InfoContext(ctx, "processing order",
      "trace_id", span.SpanContext().TraceID().String(),
      "span_id", span.SpanContext().SpanID().String(),
  )
  ```
- Never log tokens, passwords, API keys, PII, or full request/response bodies

## Testing

- Table-driven tests with `t.Run()`:
  ```go
  tests := []struct {
      name    string
      input   string
      want    int
      wantErr bool
  }{
      {name: "valid", input: "42", want: 42},
      {name: "empty", input: "", wantErr: true},
      {name: "negative", input: "-1", want: -1},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          got, err := Parse(tt.input)
          if tt.wantErr {
              require.Error(t, err)
              return
          }
          require.NoError(t, err)
          assert.Equal(t, tt.want, got)
      })
  }
  ```
- Inject dependencies — pass via constructors or method parameters. Inject `time.Time` or `func() time.Time` instead of calling `time.Now()` directly
- Clean up global state: `t.TempDir()` for temp files, truncate DB tables in `TestMain` or `t.Cleanup`
- Prefer real implementations over mocks — use testcontainers for real DB:
  ```go
  container, err := postgres.Run(ctx, "postgres:16-alpine")
  require.NoError(t, err)
  t.Cleanup(func() { container.Terminate(ctx) })
  ```
- Mock only 3rd party APIs with gomock or httptest:
  ```go
  srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      json.NewEncoder(w).Encode(expectedResponse)
  }))
  t.Cleanup(srv.Close)
  ```
- Share test setup with helper functions:
  ```go
  func newTestService(t *testing.T) *Service {
      t.Helper()
      db := newTestDB(t)
      repo := NewRepository(db)
      return NewService(repo)
  }
  ```
- Never skip or ignore test failures — fix them properly

## Concurrency

- `sync.Mutex` for shared mutable state — keep critical sections minimal:
  ```go
  type Counter struct {
      mu    sync.Mutex
      count int
  }
  func (c *Counter) Increment() {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.count++
  }
  ```
- `errgroup` with `SetLimit` for bounded concurrency when batching is not possible:
  ```go
  eg, ctx := errgroup.WithContext(ctx)
  eg.SetLimit(10)
  for _, item := range items {
      eg.Go(func() error {
          return process(ctx, item)
      })
  }
  if err := eg.Wait(); err != nil {
      return fmt.Errorf("process items: %w", err)
  }
  ```
- Channels for producer-consumer patterns; `context.Context` for cancellation and timeouts

## Security

- No hardcoded secrets — use environment variables or secret management
- Validate and sanitize all input at trust boundaries
- Use `crypto/rand` for secure random values, never `math/rand`

## Monitoring

- OpenTelemetry SDK at bootstrap:
  ```go
  import (
      "go.opentelemetry.io/otel"
      "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
      sdktrace "go.opentelemetry.io/otel/sdk/trace"
  )

  func initTracer(ctx context.Context) (func(), error) {
      exporter, err := otlptracegrpc.New(ctx)
      if err != nil {
          return nil, fmt.Errorf("create exporter: %w", err)
      }
      tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
      otel.SetTracerProvider(tp)
      return func() { tp.Shutdown(ctx) }, nil
  }
  ```
- RED metrics (Rate, Errors, Duration) for HTTP and gRPC servers — use `otelhttp` and `otelgrpc` instrumentation
- No high-cardinality labels — no user IDs or request IDs as metric labels
- Custom spans with `otel.Tracer()`:
  ```go
  ctx, span := otel.Tracer("service").Start(ctx, "ProcessOrder")
  defer span.End()
  span.SetAttributes(attribute.Int64("order.id", orderID))
  ```

## Go Guidelines

- Every error must be checked or returned — never `_` for errors
- Wrap errors with `%w` — include what operation failed and with what input:
  - Wrong: `return fmt.Errorf("failed: %w", err)`
  - Right: `return fmt.Errorf("find user by ID %d: %w", userID, err)`
- Use `errors.Is()` and `errors.As()` — never compare error strings
- Functional options for optional constructor parameters: `type Option func(*options)`, `func New(required string, opts ...Option) *T`
- Prefer `golang.org/x/sync/errgroup` over `sync.WaitGroup` — built-in error propagation, context cancellation via `errgroup.WithContext`, cleaner API with `eg.Go()`
- Configure connection pool when using `database/sql`: `SetConnMaxLifetime`, `SetMaxOpenConns`, `SetMaxIdleConns`
- Testing: use want/got (never expected/actual), **assert** when test can continue / **require** when it should stop, **go.uber.org/gomock** for mocks
- `context.Context` is always the first parameter: `func DoWork(ctx context.Context, ...) error`
- Return `error` as the last return value: `func DoWork(ctx context.Context) (Result, error)`
- Use `defer` for cleanup (close files, unlock mutexes, end spans)
- Accept interfaces, return structs
- `log/slog` for structured logging (Go 1.21+)
