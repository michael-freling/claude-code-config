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

## Testing

- Use want/got (never expected/actual)
- Use **assert** when test can continue, **require** when it should stop
- Use **go.uber.org/gomock** for mocks:
  ```go
  ctrl := gomock.NewController(t)
  mock := NewMockInterface(ctrl)
  mock.EXPECT().Method(gomock.Any(), arg).Return(result, nil)
  ```
