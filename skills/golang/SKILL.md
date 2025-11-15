---
name: golang
description: Go (Golang) feature development following best practices for error handling, interfaces, concurrency, testing, and project structure. Use when adding or updating Go code in projects.
---

# Go Development Skill

A comprehensive skill for adding or updating features in Go projects following best practices.

## When to Use

Use this skill when:
- Adding new features to Go projects
- Updating or enhancing existing Go functionality
- Need to follow Go-specific patterns and conventions
- Working with Go interfaces, concurrency, or error handling

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation:

1. **Look for `.claude/design.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/design.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no design.md exists:**
   - Consider running `/document-design` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from design.md:**
- Project-specific naming conventions
- Directory structure and organization rules
- Error handling patterns
- Testing conventions
- Architecture patterns
- Interface design patterns
- Code examples showing preferred style

### 1. Analyze Project Structure

- Check for `go.mod` to identify module name
- Identify package structure and naming conventions
- Review existing patterns (interfaces, error handling, context usage)
- Check for common Go project layouts (cmd/, pkg/, internal/)
- Identify testing patterns

### 2. Search for Relevant Code

- Search for related packages
- Find similar function implementations
- Look for interface definitions
- Identify error handling patterns

## Go Best Practices

### Code Organization

- Follow standard Go project layout
- One package per directory
- Keep packages focused and cohesive
- Use `internal/` for private packages

Example structure:
```
myproject/
├── cmd/
│   └── myapp/
│       └── main.go          # Application entrypoint
├── internal/
│   ├── user/
│   │   ├── user.go          # User domain logic
│   │   ├── user_test.go     # Tests
│   │   └── repository.go    # Data access
│   └── http/
│       └── handler.go       # HTTP handlers
├── pkg/
│   └── logger/
│       └── logger.go        # Reusable public packages
├── go.mod
└── go.sum
```

### Naming Conventions

- Use camelCase for unexported, CamelCase for exported
- Interface names: Reader, Writer, Handler (noun or noun+er)
- Keep names concise but meaningful
- Package names: lowercase, no underscores

Example:
```go
// Good: Clear, concise names
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
}

type userRepositoryImpl struct {
    db *sql.DB
}

// Good: Package-level exported function
func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepositoryImpl{db: db}
}
```

### Error Handling

- Return errors as the last return value
- Wrap errors with context using fmt.Errorf with %w
- Don't panic in library code
- Handle all errors explicitly
- Define custom error types for domain errors

Example:
```go
// Good: Custom error types
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

// Good: Error wrapping with context
func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id).Scan(
        &user.ID, &user.Name, &user.Email,
    )
    if err == sql.ErrNoRows {
        return nil, &NotFoundError{Resource: "user", ID: id}
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query user %s: %w", id, err)
    }
    return &user, nil
}

// Good: Explicit error handling
user, err := repo.FindByID(ctx, "123")
if err != nil {
    var notFound *NotFoundError
    if errors.As(err, &notFound) {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return nil, fmt.Errorf("unexpected error: %w", err)
}
```

### Interfaces

- Keep interfaces small (1-3 methods ideal)
- Define interfaces at point of use (consumer side)
- Accept interfaces, return structs
- Use interface{} sparingly; prefer concrete types or generics (Go 1.18+)

Example:
```go
// Good: Small, focused interface
type UserStore interface {
    Get(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Good: Accept interface, return struct
func NewUserService(store UserStore, logger Logger) *UserService {
    return &UserService{
        store:  store,
        logger: logger,
    }
}

// Good: Define interface where it's used
// In service package:
type Logger interface {
    Info(msg string, fields ...any)
    Error(msg string, fields ...any)
}
```

### Struct Design

- Use struct tags for JSON, DB, validation
- Group related fields
- Consider embedding for composition
- Use pointers for large structs or when you need nil

Example:
```go
// Good: Well-structured domain model
type User struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name" validate:"required,min=1,max=100"`
    Email     string    `json:"email" db:"email" validate:"required,email"`
    Role      UserRole  `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
    UserRoleAdmin UserRole = "admin"
    UserRoleUser  UserRole = "user"
    UserRoleGuest UserRole = "guest"
)

// Good: Embedding for composition
type AuditLog struct {
    User                    // Embedded
    Action    string       `json:"action"`
    Timestamp time.Time    `json:"timestamp"`
}
```

### Concurrency

- Use context.Context for cancellation and timeouts
- Avoid goroutine leaks with proper cleanup
- Use channels or sync primitives appropriately
- Close channels when done sending

Example:
```go
// Good: Context-aware concurrent processing
func (s *UserService) ProcessUsers(ctx context.Context, ids []string) error {
    eg, ctx := errgroup.WithContext(ctx)

    for _, id := range ids {
        id := id // Capture loop variable
        eg.Go(func() error {
            return s.processUser(ctx, id)
        })
    }

    return eg.Wait()
}

// Good: Channel-based worker pool
func (s *UserService) ProcessWithWorkers(ctx context.Context, users []User) error {
    jobs := make(chan User, len(users))
    results := make(chan error, len(users))

    // Start workers
    for w := 0; w < 5; w++ {
        go func() {
            for user := range jobs {
                select {
                case <-ctx.Done():
                    results <- ctx.Err()
                    return
                default:
                    results <- s.processUser(ctx, user.ID)
                }
            }
        }()
    }

    // Send jobs
    for _, user := range users {
        jobs <- user
    }
    close(jobs)

    // Collect results
    for range users {
        if err := <-results; err != nil {
            return err
        }
    }

    return nil
}
```

### Function Design

- Use functional options pattern for complex constructors
- Return early to reduce nesting
- Keep functions small and focused

Example:
```go
// Good: Functional options pattern
type ServerOptions struct {
    Port    int
    Timeout time.Duration
    Logger  Logger
}

type ServerOption func(*ServerOptions)

func WithPort(port int) ServerOption {
    return func(o *ServerOptions) {
        o.Port = port
    }
}

func WithTimeout(timeout time.Duration) ServerOption {
    return func(o *ServerOptions) {
        o.Timeout = timeout
    }
}

func NewServer(opts ...ServerOption) *Server {
    options := &ServerOptions{
        Port:    8080,
        Timeout: 30 * time.Second,
    }

    for _, opt := range opts {
        opt(options)
    }

    return &Server{options: options}
}

// Usage:
server := NewServer(
    WithPort(9000),
    WithTimeout(60 * time.Second),
)
```

### Testing

**CRITICAL: Always use table-driven tests as the primary testing approach.**

Table-driven tests are the idiomatic Go way to write tests and should be used for:
- All functions with multiple test cases
- Error handling scenarios
- Edge cases and boundary conditions
- Different input/output combinations

**Benefits of table-driven tests:**
- Reduce code duplication
- Make it easy to add new test cases
- Improve test readability
- Ensure consistent test structure
- Make test coverage gaps obvious

**Core principles:**
- Define test cases as a slice of structs
- Use t.Run() for subtests with descriptive names
- Use t.Helper() for test helpers
- Test exported APIs, not implementation details
- Use clear field names in test case structs

**Table-Driven Test Examples:**

```go
// Good: Basic table-driven test
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {name: "positive numbers", a: 2, b: 3, want: 5},
        {name: "negative numbers", a: -2, b: -3, want: -5},
        {name: "mixed signs", a: 5, b: -3, want: 2},
        {name: "zero values", a: 0, b: 0, want: 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}

// Good: Table-driven test with error handling
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
        errType error // Optional: specific error type
    }{
        {
            name: "valid user",
            user: User{
                Name:  "John Doe",
                Email: "john@example.com",
                Role:  UserRoleUser,
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: User{
                Name:  "John Doe",
                Email: "invalid-email",
                Role:  UserRoleUser,
            },
            wantErr: true,
            errType: &ValidationError{},
        },
        {
            name: "empty name",
            user: User{
                Name:  "",
                Email: "john@example.com",
                Role:  UserRoleUser,
            },
            wantErr: true,
            errType: &ValidationError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if tt.errType != nil && !errors.As(err, &tt.errType) {
                t.Errorf("Validate() error type = %T, want %T", err, tt.errType)
            }
        })
    }
}

// Good: Table-driven test with context and setup
func TestUserRepository_FindByID(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewUserRepository(db)

    // Pre-populate test data
    testUser := createTestUser(t, db, "John Doe")

    tests := []struct {
        name    string
        id      string
        want    *User
        wantErr bool
    }{
        {
            name:    "existing user",
            id:      testUser.ID,
            want:    testUser,
            wantErr: false,
        },
        {
            name:    "non-existent user",
            id:      "non-existent-id",
            want:    nil,
            wantErr: true,
        },
        {
            name:    "empty id",
            id:      "",
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            got, err := repo.FindByID(ctx, tt.id)

            if (err != nil) != tt.wantErr {
                t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FindByID() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Good: Table-driven test with multiple return values
func TestParseConfig(t *testing.T) {
    tests := []struct {
        name       string
        input      string
        wantConfig *Config
        wantErr    bool
    }{
        {
            name:  "valid config",
            input: `{"host": "localhost", "port": 8080}`,
            wantConfig: &Config{
                Host: "localhost",
                Port: 8080,
            },
            wantErr: false,
        },
        {
            name:       "invalid json",
            input:      `{invalid}`,
            wantConfig: nil,
            wantErr:    true,
        },
        {
            name:       "empty string",
            input:      "",
            wantConfig: nil,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseConfig(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.wantConfig) {
                t.Errorf("ParseConfig() = %v, want %v", got, tt.wantConfig)
            }
        })
    }
}

// Good: Test helper with t.Helper()
func createTestUser(t *testing.T, db *sql.DB, name string) *User {
    t.Helper()

    user := &User{
        ID:    uuid.New().String(),
        Name:  name,
        Email: fmt.Sprintf("%s@example.com", name),
    }

    err := db.QueryRow(
        "INSERT INTO users (id, name, email) VALUES (?, ?, ?) RETURNING created_at",
        user.ID, user.Name, user.Email,
    ).Scan(&user.CreatedAt)
    if err != nil {
        t.Fatalf("failed to create test user: %v", err)
    }

    return user
}

// Good: Setup helper
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test database: %v", err)
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }

    return db
}
```

**When NOT to use table-driven tests:**
- Single test case with complex setup/teardown
- Tests that require significantly different setup per case
- Integration tests with heavy external dependencies
- Tests where the table structure becomes too complex

**Best Practices:**
1. Always use descriptive test case names
2. Group related test cases together
3. Use t.Helper() in helper functions
4. Clean up resources in defer statements
5. Use t.Parallel() for independent tests
6. Consider testify/assert for cleaner assertions if the project uses it

### Mocking External Dependencies

**When to use mocking:**
- Testing code that depends on external systems (databases, APIs, third-party services)
- When httptest is not suitable (non-HTTP external dependencies)
- Testing error scenarios from external systems
- Avoiding slow or flaky external calls in unit tests
- Testing code in isolation

**Approaches to mocking in Go:**

1. **Interface-based mocking** (Preferred, idiomatic Go)
2. **gomock/mockgen** for complex scenarios
3. **httptest** for HTTP-specific testing

**1. Interface-Based Mocking (The Go Way)**

The idiomatic Go approach is to define interfaces and create test implementations:

```go
// Production code - define interface
type WeatherClient interface {
    GetTemperature(ctx context.Context, city string) (float64, error)
}

type WeatherService struct {
    client WeatherClient
}

func NewWeatherService(client WeatherClient) *WeatherService {
    return &WeatherService{client: client}
}

func (s *WeatherService) GetCityStatus(ctx context.Context, city string) (string, error) {
    temp, err := s.client.GetTemperature(ctx, city)
    if err != nil {
        return "", fmt.Errorf("failed to get temperature: %w", err)
    }

    if temp > 30 {
        return "hot", nil
    }
    return "cold", nil
}

// Test code - implement mock
type mockWeatherClient struct {
    temperature float64
    err         error
}

func (m *mockWeatherClient) GetTemperature(ctx context.Context, city string) (float64, error) {
    return m.temperature, m.err
}

// Table-driven test with mock
func TestWeatherService_GetCityStatus(t *testing.T) {
    tests := []struct {
        name        string
        temperature float64
        clientErr   error
        want        string
        wantErr     bool
    }{
        {
            name:        "hot day",
            temperature: 35.0,
            clientErr:   nil,
            want:        "hot",
            wantErr:     false,
        },
        {
            name:        "cold day",
            temperature: 15.0,
            clientErr:   nil,
            want:        "cold",
            wantErr:     false,
        },
        {
            name:        "client error",
            temperature: 0,
            clientErr:   errors.New("network error"),
            want:        "",
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := &mockWeatherClient{
                temperature: tt.temperature,
                err:         tt.clientErr,
            }

            service := NewWeatherService(mockClient)
            got, err := service.GetCityStatus(context.Background(), "TestCity")

            if (err != nil) != tt.wantErr {
                t.Errorf("GetCityStatus() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("GetCityStatus() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**2. Using gomock for Complex Mocking**

For complex interfaces or when you need to verify method calls, use [gomock](https://github.com/uber-go/mock):

```bash
# Install gomock
go get github.com/uber-go/mock/gomock
go install github.com/uber-go/mock/mockgen@latest

# Generate mocks
mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks
```

Example with gomock:

```go
// client.go - production interface
package weather

import "context"

type APIClient interface {
    GetTemperature(ctx context.Context, city string) (float64, error)
    GetHumidity(ctx context.Context, city string) (float64, error)
    GetForecast(ctx context.Context, city string, days int) ([]Forecast, error)
}

type Forecast struct {
    Date        string
    Temperature float64
    Condition   string
}

// Generate mocks with:
// mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks

// service_test.go - test with gomock
package weather_test

import (
    "context"
    "errors"
    "testing"

    "github.com/uber-go/mock/gomock"
    "yourproject/weather"
    "yourproject/weather/mocks"
)

func TestWeatherService_GetDetailedStatus(t *testing.T) {
    tests := []struct {
        name        string
        city        string
        setupMock   func(*mocks.MockAPIClient)
        want        string
        wantErr     bool
    }{
        {
            name: "hot and humid",
            city: "Miami",
            setupMock: func(m *mocks.MockAPIClient) {
                m.EXPECT().
                    GetTemperature(gomock.Any(), "Miami").
                    Return(35.0, nil)
                m.EXPECT().
                    GetHumidity(gomock.Any(), "Miami").
                    Return(80.0, nil)
            },
            want:    "hot and humid",
            wantErr: false,
        },
        {
            name: "temperature call fails",
            city: "Boston",
            setupMock: func(m *mocks.MockAPIClient) {
                m.EXPECT().
                    GetTemperature(gomock.Any(), "Boston").
                    Return(0.0, errors.New("API error"))
            },
            want:    "",
            wantErr: true,
        },
        {
            name: "verify call order",
            city: "Seattle",
            setupMock: func(m *mocks.MockAPIClient) {
                gomock.InOrder(
                    m.EXPECT().GetTemperature(gomock.Any(), "Seattle").Return(20.0, nil),
                    m.EXPECT().GetHumidity(gomock.Any(), "Seattle").Return(60.0, nil),
                )
            },
            want:    "mild",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockClient := mocks.NewMockAPIClient(ctrl)
            tt.setupMock(mockClient)

            service := weather.NewService(mockClient)
            got, err := service.GetDetailedStatus(context.Background(), tt.city)

            if (err != nil) != tt.wantErr {
                t.Errorf("GetDetailedStatus() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("GetDetailedStatus() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**3. Using httptest for HTTP Clients**

When the external dependency is an HTTP API, prefer httptest:

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHTTPClient_GetUser(t *testing.T) {
    tests := []struct {
        name           string
        serverResponse string
        serverStatus   int
        wantUser       *User
        wantErr        bool
    }{
        {
            name:           "successful response",
            serverResponse: `{"id": "123", "name": "John"}`,
            serverStatus:   http.StatusOK,
            wantUser:       &User{ID: "123", Name: "John"},
            wantErr:        false,
        },
        {
            name:           "server error",
            serverResponse: `{"error": "internal error"}`,
            serverStatus:   http.StatusInternalServerError,
            wantUser:       nil,
            wantErr:        true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create test server
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.serverStatus)
                w.Write([]byte(tt.serverResponse))
            }))
            defer server.Close()

            // Create client pointing to test server
            client := NewHTTPClient(server.URL)
            got, err := client.GetUser(context.Background(), "123")

            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.wantUser) {
                t.Errorf("GetUser() = %v, want %v", got, tt.wantUser)
            }
        })
    }
}
```

**Best Practices for Testable Code:**

1. **Design for testability**: Accept interfaces, return concrete types
```go
// Good: Accepts interface
func NewService(db Database, cache Cache, logger Logger) *Service {
    return &Service{db: db, cache: cache, logger: logger}
}

// Bad: Accepts concrete implementation
func NewService(db *sql.DB, cache *redis.Client) *Service {
    return &Service{db: db, cache: cache}
}
```

2. **Keep interfaces small**: Easier to mock and maintain
```go
// Good: Small, focused interface
type UserGetter interface {
    GetUser(ctx context.Context, id string) (*User, error)
}

// Bad: Large interface difficult to mock
type UserRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
    ListUsers(ctx context.Context, filters Filters) ([]*User, error)
}
```

3. **Use constructor functions**: Makes dependency injection clear
```go
type Service struct {
    client ExternalClient
    db     Database
}

func NewService(client ExternalClient, db Database) *Service {
    return &Service{
        client: client,
        db:     db,
    }
}
```

**When to choose each approach:**

- **Interface-based mocking**: For simple interfaces (1-3 methods), straightforward testing
- **gomock**: For complex interfaces, when you need to verify call order or call counts, multiple return values
- **httptest**: Specifically for testing HTTP clients and servers

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract project-specific patterns
- Note architecture constraints

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into small, testable units
- Identify dependencies

**Step 3: Read Existing Code**
- Review similar implementations
- Understand existing patterns
- Note utilities available

**Step 4: Implement Changes**
- Follow patterns from design.md
- Use naming conventions from project
- Write clean, readable code
- Match code style

**Step 5: Ensure Quality and Fix All Issues**

**CRITICAL: All quality checks must pass before considering the task complete.**

1. **Format code**
   ```bash
   go fmt ./...
   ```
   - If files are modified, review changes

2. **Run go vet**
   ```bash
   go vet ./...
   ```
   - If vet reports issues, fix them immediately
   - Do not proceed until vet passes

3. **Run linter (if available)**
   ```bash
   golangci-lint run
   ```
   - If linter reports issues, fix them immediately
   - Do not proceed until linter passes
   - If golangci-lint is not available, skip this step

4. **Build the project**
   ```bash
   go build ./...
   ```
   - If build fails, fix the errors immediately
   - Do not proceed until build succeeds

5. **Run tests**
   ```bash
   go test ./...
   ```
   - If tests fail, fix them immediately
   - Add new tests if implementing new features
   - Update existing tests if modifying behavior
   - Do not proceed until all tests pass

6. **Update tests for your changes**
   - If you added a new function/method, add corresponding tests
   - If you modified behavior, update existing tests
   - Ensure test coverage is maintained or improved

7. **Run go mod tidy**
   ```bash
   go mod tidy
   ```

**Iterative Fix Process:**
- If any step fails, fix the issues and re-run ALL previous steps
- Continue iterating until ALL checks pass:
  - ✅ Code formatted
  - ✅ No vet issues
  - ✅ No lint issues (if linter available)
  - ✅ Build succeeds
  - ✅ All tests pass
  - ✅ Dependencies are tidy

## Commands Reference

```bash
# Format code
go fmt ./...
gofmt -s -w .

# Vet code
go vet ./...

# Lint (if golangci-lint installed)
golangci-lint run

# Build
go build ./...
go build -o bin/app ./cmd/app

# Test
go test ./...
go test -v ./...
go test -cover ./...
go test -race ./...

# Mod management
go mod tidy
go mod download
go mod verify
```

## Common Pitfalls to Avoid

### Error Handling
```go
// Bad: Ignoring errors
data, _ := readFile(path)

// Good: Handling errors
data, err := readFile(path)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

### Goroutine Leaks
```go
// Bad: Goroutine leak
func process(data []string) {
    for _, item := range data {
        go processItem(item) // No way to stop these
    }
}

// Good: Controlled goroutines
func process(ctx context.Context, data []string) error {
    eg, ctx := errgroup.WithContext(ctx)
    for _, item := range data {
        item := item
        eg.Go(func() error {
            return processItem(ctx, item)
        })
    }
    return eg.Wait()
}
```

### Interface Design
```go
// Bad: Large interface
type UserService interface {
    Create(User) error
    Update(User) error
    Delete(string) error
    FindByID(string) (*User, error)
    FindByEmail(string) (*User, error)
    List() ([]User, error)
}

// Good: Small, focused interfaces
type UserCreator interface {
    Create(ctx context.Context, user User) error
}

type UserFinder interface {
    FindByID(ctx context.Context, id string) (*User, error)
}
```

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Extract project-specific conventions
- [ ] Analyze project structure and go.mod
- [ ] Search for similar implementations
- [ ] Review existing patterns

### During Implementation
- [ ] Define interfaces at point of use
- [ ] Include context.Context in long-running operations
- [ ] Return errors as last value
- [ ] Wrap errors with context using %w
- [ ] Write table-driven tests
- [ ] Keep functions small and focused

### After Implementation - MUST ALL PASS
- [ ] Run `go fmt ./...` and review any changes
- [ ] Run `go vet ./...` - **FIX ALL ISSUES**
- [ ] Run `golangci-lint run` if available - **FIX ALL ISSUES**
- [ ] Run `go build ./...` - **MUST SUCCEED**
- [ ] Run `go test ./...` - **ALL TESTS MUST PASS**
- [ ] Add/update tests for new or modified code
- [ ] Run `go mod tidy`
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Simplicity**: Write clear, idiomatic Go code
3. **Error Handling**: Handle all errors explicitly
4. **Interfaces**: Keep them small and define at point of use
5. **Concurrency**: Use context for cancellation and avoid leaks
6. **Testing**: Write table-driven tests for all exported APIs
