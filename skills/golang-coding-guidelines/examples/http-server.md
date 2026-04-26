# HTTP Server Patterns

## Handler with Dependency Injection

```go
type UserHandler struct {
    userService UserService
    logger      *slog.Logger
}

func NewUserHandler(userService UserService, logger *slog.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}
```

## Request Handling with Tracing

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("handler").Start(r.Context(), "CreateUser")
    defer span.End()

    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    if err := req.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.userService.Create(ctx, req)
    if err != nil {
        span.RecordError(err)
        h.logger.ErrorContext(ctx, "create user", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

## Middleware with OpenTelemetry Metrics

```go
func MetricsMiddleware(meter metric.Meter) func(http.Handler) http.Handler {
    requestCount, _ := meter.Int64Counter("http.server.request.count")
    requestDuration, _ := meter.Float64Histogram("http.server.request.duration",
        metric.WithUnit("ms"),
    )

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

            next.ServeHTTP(rw, r)

            attrs := metric.WithAttributes(
                attribute.String("method", r.Method),
                attribute.String("path", r.Pattern),
                attribute.Int("status", rw.statusCode),
            )
            requestCount.Add(r.Context(), 1, attrs)
            requestDuration.Record(r.Context(), float64(time.Since(start).Milliseconds()), attrs)
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## Router Setup

```go
func NewRouter(userHandler *UserHandler, orderHandler *OrderHandler, meter metric.Meter) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("POST /users", userHandler.CreateUser)
    mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
    mux.HandleFunc("POST /orders", orderHandler.CreateOrder)

    var handler http.Handler = mux
    handler = MetricsMiddleware(meter)(handler)

    return handler
}
```

## Graceful Shutdown

```go
func Run(ctx context.Context, handler http.Handler, addr string) error {
    server := &http.Server{
        Addr:              addr,
        Handler:           handler,
        ReadHeaderTimeout: 10 * time.Second,
    }

    errCh := make(chan error, 1)
    go func() {
        errCh <- server.ListenAndServe()
    }()

    select {
    case err := <-errCh:
        return fmt.Errorf("server listen: %w", err)
    case <-ctx.Done():
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        if err := server.Shutdown(shutdownCtx); err != nil {
            return fmt.Errorf("server shutdown: %w", err)
        }
        return nil
    }
}
```
