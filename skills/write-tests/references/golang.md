# Go Testing Patterns

## Pattern Purposes

| Pattern | Purpose |
|---------|---------|
| **Dependency Injection** | Pass dependencies as parameters instead of using global access; enables mocking in tests |
| **Extract Interfaces** | Define interfaces for external dependencies to enable mock implementations |
| **Table-Driven Tests** | Reduces duplication by defining test cases as data; easier to add cases and see all scenarios |
| **Test Data Builders** | Creates consistent test objects with sensible defaults; avoids repeating struct initialization |
| **Mock Dependencies** | Implements interfaces with controllable behavior for testing |
| **HTTP Testing** | Tests HTTP handlers without starting a real server |

## Dependency Injection

```go
// Bad: Hard to test - uses global database connection
func GetUser(id string) (*User, error) {
    db := database.GetConnection()
    return db.Query("SELECT * FROM users WHERE id = ?", id)
}

// Good: Testable - accepts store as parameter
func GetUser(store UserStore, id string) (*User, error) {
    return store.Get(id)
}
```

## Extract Interfaces

```go
// Define interfaces for dependencies to enable mocking
type UserStore interface {
    Get(id string) (*User, error)
    Save(user *User) error
}

type EmailSender interface {
    Send(to, subject, body string) error
}
```

## Test Data Builders

```go
func createUser(opts ...func(*User)) *User {
    u := &User{
        ID:    "1",
        Name:  "Test User",
        Email: "test@example.com",
    }
    for _, opt := range opts {
        opt(u)
    }
    return u
}

// Usage
user := createUser()
admin := createUser(func(u *User) { u.Role = "admin" })
```

## Table-Driven Tests

```go
import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
    successCases := []struct {
        name  string
        email string
    }{
        {name: "simple email", email: "user@example.com"},
        {name: "with subdomain", email: "user@mail.example.com"},
    }

    for _, tc := range successCases {
        t.Run(tc.name, func(t *testing.T) {
            err := ValidateEmail(tc.email)
            require.NoError(t, err)
        })
    }

    errorCases := []struct {
        name  string
        email string
    }{
        {name: "empty email", email: ""},
        {name: "missing @", email: "userexample.com"},
        {name: "missing domain", email: "user@"},
    }

    for _, tc := range errorCases {
        t.Run(tc.name, func(t *testing.T) {
            err := ValidateEmail(tc.email)
            require.Error(t, err)
        })
    }
}
```

## Mock Dependencies

```go
type MockUserStore struct {
    users map[string]*User
    err   error
}

func (m *MockUserStore) Get(id string) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.users[id], nil
}

func (m *MockUserStore) Save(user *User) error {
    if m.err != nil {
        return m.err
    }
    m.users[user.ID] = user
    return nil
}

// Usage in tests
func TestGetUser(t *testing.T) {
    mock := &MockUserStore{
        users: map[string]*User{
            "1": createUser(),
        },
    }

    got, err := GetUser(mock, "1")
    require.NoError(t, err)
    assert.Equal(t, "1", got.ID)
}
```

## HTTP Handler Testing

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetUserHandler(t *testing.T) {
    mock := &MockUserStore{
        users: map[string]*User{
            "1": createUser(),
        },
    }
    handler := NewUserHandler(mock)

    successCases := []struct {
        name       string
        userID     string
        wantStatus int
    }{
        {name: "existing user", userID: "1", wantStatus: http.StatusOK},
    }

    for _, tc := range successCases {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/users/"+tc.userID, nil)
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            assert.Equal(t, tc.wantStatus, rec.Code)
        })
    }

    errorCases := []struct {
        name       string
        userID     string
        wantStatus int
    }{
        {name: "not found", userID: "999", wantStatus: http.StatusNotFound},
    }

    for _, tc := range errorCases {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/users/"+tc.userID, nil)
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            assert.Equal(t, tc.wantStatus, rec.Code)
        })
    }
}
```
