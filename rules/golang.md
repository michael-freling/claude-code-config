---
paths:
  - "**/*.go"
---

# Go Coding Guidelines

## Errors

- Every error must be checked or returned
- Wrap errors with context using `%w`:
  ```go
  // Wrong
  return fmt.Errorf("failed: %w", err)
  // Right
  return fmt.Errorf("find user by ID %d: %w", userID, err)
  ```

## Initialization

- **Functional options** for optional constructor parameters:
  ```go
  type Option func(*options)
  func New(required string, opts ...Option) *T { ... }
  ```

## Concurrency

Prefer `golang.org/x/sync/errgroup` over `sync.WaitGroup` for managing concurrent goroutines:

- Provides built-in error propagation
- Supports context cancellation via `errgroup.WithContext`
- Cleaner API with `eg.Go()` instead of manual `wg.Add/Done`

```go
var eg errgroup.Group
eg.Go(func() error {
    // Goroutine work
    return nil
})
if err := eg.Wait(); err != nil {
    return fmt.Errorf("eg.Wait: %w", err)
}
```

## Database

- **Configure connection pool** when using `database/sql`:
  ```go
  db.SetConnMaxLifetime(maxLifetime)
  db.SetMaxOpenConns(maxOpen)
  db.SetMaxIdleConns(maxIdle)
  ```

## Testing

- Use want/got (never expected/actual)
- Use **assert** when test can continue, **require** when it should stop
- Use **go.uber.org/gomock** for mocks:
  ```go
  ctrl := gomock.NewController(t)
  mock := NewMockInterface(ctrl)
  mock.EXPECT().Method(gomock.Any(), arg).Return(result, nil)
  ```
