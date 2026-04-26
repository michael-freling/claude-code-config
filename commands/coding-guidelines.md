---
description: Generate a language-specific coding guidelines skill based on the standard template
argument-hint: "<language or framework, e.g., Go, TypeScript, Python, Rust>"
allowed-tools: ["*"]
---

# Generate Coding Guidelines Skill

Generate a coding guidelines skill for: $ARGUMENTS

## Key Rules

- The generated SKILL.md must cover every section from the coding guidelines template below, but rewritten as a native style guide for the target language — using that language's types, patterns, standard libraries, and idioms throughout. Every bullet point must include language-specific examples, patterns, or library references. Do not copy the template with minor word swaps
- Language idioms go in SKILL.md alongside the template; application-level patterns go in separate example files — one file per application type most common in that language's ecosystem (e.g., Go: HTTP server, HTTP client, gRPC server, gRPC client, DB, message queue worker, one-time job; TypeScript: React components, Next.js pages, API routes, DB)
- Each example file must include OpenTelemetry instrumentation inline where applicable — not in a separate file
- Example files complement the template — they must NOT repeat general principles already in the template
- Use the language's standard tooling and most widely adopted libraries
- Do not include guidelines for things that linters or formatters automatically fix

## Workflow

### Phase 1: Check Existing

Check if `skills/<language>-coding-guidelines/` already exists. If it does, update it rather than creating from scratch.

### Phase 2: Research

Research language-specific conventions for:

**Language idioms:**
- Error handling patterns (e.g., Result types, exceptions, error values)
- Constructor/initialization patterns (e.g., builders, factory functions, functional options)
- Concurrency primitives and patterns
- Testing framework, assertion style, mock libraries
- Type system features and idioms
- Database libraries and connection patterns

**Application-level patterns** — identify the most common application types for this language's ecosystem and research idiomatic patterns for each (with OpenTelemetry instrumentation inline where applicable)

### Phase 3: Generate

Create the following structure:

```
skills/<language>-coding-guidelines/
├── SKILL.md
└── examples/
    ├── <application-type-1>.md
    ├── <application-type-2>.md
    └── ...
```

**SKILL.md:**
1. Frontmatter with `name: <language>-coding-guidelines` and `description`
2. Pointers to each example file in `examples/`
3. Language idioms section — error patterns, type system, initialization, concurrency, testing conventions specific to this language
4. Every section from the coding guidelines template below, rewritten as native guidance for the target language. Each bullet point should reference specific types, functions, libraries, or patterns from that language's ecosystem. The template provides the structure and topics to cover — the output must read as a language-native style guide, not a copy of the template

**examples/** — one file per application type most common in the language's ecosystem. Each file should include OpenTelemetry instrumentation inline where it naturally applies. Follow the format in the "Example Output" section below.

### Phase 4: Review

- Does the SKILL.md contain language idioms and the full adapted template?
- Does each example file contain idiomatic, runnable code patterns?
- Does each example file include OpenTelemetry instrumentation inline where applicable?
- Are example files scoped to their application type?
- Do example files complement the template without repeating general principles?
- For Go: does the SKILL.md include the Go Guidelines section?
- For frontend languages: does the SKILL.md include the Frontend Guidelines section?

## Coding Guidelines Template

### Principles

- **Environment-agnostic code** — do not add environment-specific branches (dev/staging/prod checks, OS-specific paths, CI-specific logic) in application code; use configuration, environment variables, or dependency injection so the same code runs consistently across all environments

### Simplicity

**Simplicity is the most important thing.** Minimize the final complexity of the codebase — not the size of the change.

- **Breaking changes are acceptable** unless backward compatibility is explicitly required by users or project constraints
- **Avoid premature abstraction** — do not create an abstraction until there are 3 or more concrete use cases; three similar lines are better than a wrong abstraction
- **Minimize mutable state** — prefer stateless functions and derived values over cached or stored state; keep mutable state as local as possible

#### High Cohesion

Each function should do one thing. When a new requirement applies to only some callers, split the function rather than adding conditional logic that complicates it for all callers.

- **Wrong**: `sendNotification()` handles email, SMS, and push — adding rate limiting for SMS forces all paths through the rate limiter check
- **Right**: `sendEmail()`, `sendSMS()`, `sendPush()` — rate limiting lives only in `sendSMS()` where it belongs

#### Separate Methods over Boolean Arguments

When a boolean argument causes a function to branch into substantially different logic, split into separate methods with clear responsibilities.

- **Wrong**: `CompleteJob(success bool)` with branching logic based on the flag
- **Right**: `CompleteJob()` and `FailJob()` as separate methods

#### Narrow Interfaces

Accept the minimum type or interface a function needs. This reduces coupling and makes the function usable in more contexts.

- **Wrong**: A function that only reads accepts a `File` (which also supports write, seek, close)
- **Right**: The function accepts a `Reader` interface — the narrowest contract that satisfies its needs

### Error Handling

- **Always handle errors explicitly** — never silently ignore or swallow errors
- **Wrap errors with context** — include what operation failed and with what input, enabling diagnosis without reading the source
- **Fail fast on invalid state** — detect and report errors at the earliest point, before they propagate and become harder to diagnose

#### Server-Side Errors

- Return appropriate status codes — 4XX for client errors (bad input, auth failures), 5XX for server-side failures
- Include user-friendly error messages with actionable guidance so users know how to fix the error
- Never expose internal error details (stack traces, internal messages) to end users — use generic messages for server errors

#### Client-Side Errors

- Retry transient failures (rate limits, server errors, timeouts) with exponential backoff
- Document or make configurable: max retry count, backoff strategy, which errors are retryable
- Some systems handle retries at the infrastructure layer — do not duplicate if already handled

### Dependencies

- When installing applications, libraries, or tools, always check and use the latest stable version with compatibility with existing systems
- **Pin exact versions** for reproducible builds

### Code Readability

- **Use specific names, not generic ones** (shared, common, utils, info). Exception: short-lived variables scoped within ~5 lines can use short names like `i` or `p`
- **Minimal comments** — only high-level explanations of purpose, architecture, or non-obvious decisions. No line-by-line comments
- **Early returns over nesting** — use guard clauses at the top of a function; keep the main logic at the base indentation level
- **Single level of abstraction** — each function should operate at one level; do not mix high-level orchestration with low-level details
- **Delete dead code** — remove unused functions, unreachable branches, commented-out code, and unused imports/variables when encountered during the current task; do not leave them for later cleanup
- **Delete assignments of default or zero values**
- **Avoid long parameter lists** — if a function has more than 5 parameters or more than 3 return values, define a custom type instead; this also reduces caller changes when the signature evolves
- **Load binary data from files** instead of embedding base64 or other encoded strings inline (images, fonts, certificates, etc.)
- **Validate only at context boundaries** — validate incoming client requests and outgoing 3rd party API responses. Do not add defensive validation in internal code, database records managed by the same application, or private functions. For libraries, validate in public functions only

### Code Reuse

- **Single definition, owned by the provider** — each abstraction (interface, type, function) should have exactly one definition in the codebase, defined in the provider/implementation module, not in each consumer. Do not define duplicate abstractions that describe the same behavior in separate modules — search for and reuse the existing definition
- **Named constants over magic values** — extract magic numbers and strings into named constants when used in multiple places
- **Check existing content before writing** — search for similar code, logic, or documentation before adding new content; reuse or extend existing patterns rather than duplicating; when consolidating files, read the destination thoroughly first
- **Refactor before adding new changes** — understand existing code, identify reusable code, refactor to enable reuse, then add new functionality
- **Refactor after QA passes** — once new changes work correctly, refactor for code reusability and consistency

### Code Organization

- **Minimal main package** — the entry point should only initialize global objects (DB connections, clients, config) and start the main process; keep business logic out of main
- **Boundary layers by type, domain logic by domain** — group boundary/infrastructure code (HTTP handlers, gRPC handlers, DB clients, 3rd party SDKs) into one package per layer; group domain/business logic by feature (e.g., `user/`, `billing/`) rather than by technical role (e.g., `services/`, `repositories/`)
- **Dependency direction** — higher-level modules depend on lower-level modules, never the reverse; avoid circular dependencies
- **Minimize exports** — expose only what consumers need; keep implementation details internal

### Type Safety

- **Prefer enums or sum types over raw strings** for known sets of values

### Initialization

- **Separate required from optional parameters** in constructors — required parameters are positional arguments, optional parameters use a configuration pattern (options struct, builder, functional options). Do not introduce optional parameter patterns when all parameters are required
- **Validate invariants at construction time** — ensure objects are always in a valid state after creation, not at every call site

### External I/O

- **Batch operations** instead of loops for databases, APIs, and other external calls — reduces network round-trips and enables optimized bulk queries and writes
- **Transactions** — use one transaction per request context when writing to multiple tables; for batch processing, use one transaction per batch unit and repeat
- **Ensure indexes** for queries on large tables (10K+ records)
- **Pass values from program, not SQL functions** — compute values in application code instead of using SQL functions like `NOW()` or `CURRENT_TIMESTAMP`; enables tests to inject controlled values
- **Avoid redundant queries** — don't query external systems for data already available in the current context (message payload, request, event object); don't query the same record multiple times in a single operation
- **Use parameterized queries** — never concatenate user input into SQL or command strings

### Logging

- **Log levels**: Debug (development only, not production), Info (important operational details), Warn (expected/recoverable errors), Error (unexpected errors requiring investigation)
- **Use structured logging** — key-value pairs, not string interpolation
- **Include correlation IDs** — attach trace span IDs, user IDs, and relevant request parameters for traceability
- **Never log sensitive data** — tokens, passwords, API keys, PII, environment variables, or full HTTP request/response bodies must never appear in logs; redact if partial output is needed

### Testing

- **Tests must be deterministic** — no flaky tests; aim for at least 95% branch coverage
- **Inject dependencies** — pass dependencies via constructors or method parameters to enable mocking and control. Inject time values or time-generating functions instead of using sleep or reading system time directly
- **Extract interfaces for external dependencies** — define interfaces for databases, APIs, and other external services to enable mock implementations in tests
- **Clean up global state before each test** — reset database records, use temporary files/directories that auto-delete, initialize mock servers fresh
- **Use table-driven tests** — group related cases, separate success from error cases
- **Prefer real implementations over mocks** — use real databases via containers; mock only 3rd party APIs and external services that cannot run locally
- **Share test setup** — use factory functions with sensible defaults for test data (test data builders); use shared initialization objects that hold DB clients, mock servers, and test infrastructure across test cases
- **Never skip or ignore test failures** — fix them properly
- **Load test data from files** — do not embed base64 or encoded strings for test data like images or certificates; generate files and read them from test code

#### What to Test

- All public functions
- Error handling paths
- Edge cases (empty, null, zero, max values)
- Boundary conditions

#### What NOT to Test

- Third-party library internals
- Simple getters/setters with no logic
- Framework internals

#### Coverage Gaps to Look For

- Untested error branches
- Missing edge cases
- Conditional logic not fully covered

#### Writing Testable Code

When adding tests to existing code, follow this order:
1. **Analyze testability** — identify dependencies, global state, and coupling that prevent testing
2. **Refactor first** — extract interfaces, inject dependencies, split large functions before writing tests
3. **Write tests** — use table-driven patterns, cover success/error/edge cases
4. **Verify coverage** — ensure coverage is maintained or increased

### Concurrency

- **Synchronize shared state** — when multiple concurrent tasks access shared mutable state, use proper locking (mutex, channels, atomic operations)
- **Use concurrency when batching is not possible** — when an external API does not support batch requests, send requests concurrently to reduce end-to-end latency
- **Producer-subscriber for fire-and-forget** — when a response is not needed (e.g., sending notifications), consider message queues to decouple the work from the request path

### Security

- **No hardcoded secrets** — never commit credentials, API keys, or tokens; use environment variables or secret management
- **Sanitize at trust boundaries** — validate and sanitize all input crossing trust domains (user input, external API responses)

### Monitoring

- **OpenTelemetry** as the standard for metrics, logs, and traces — set up the SDK at bootstrap and add instrumentation as needed
- **RED metrics** (Rate, Errors, Duration) for HTTP and gRPC servers at minimum
- **Avoid high cardinality** in metric labels — TSDB backends handle high cardinality poorly

### Go Guidelines

Include this section in the SKILL.md only when the target language is Go. These are Go-specific idioms that must appear in the generated skill.

- **Every error must be checked or returned** — never use `_` to discard errors
- **Wrap errors with `%w`** — include what operation failed and with what input:
  - Wrong: `return fmt.Errorf("failed: %w", err)`
  - Right: `return fmt.Errorf("find user by ID %d: %w", userID, err)`
- **Use `errors.Is()` and `errors.As()`** — never compare error strings
- **Functional options** for optional constructor parameters: `type Option func(*options)`, `func New(required string, opts ...Option) *T`
- **Prefer `golang.org/x/sync/errgroup`** over `sync.WaitGroup` — built-in error propagation, context cancellation via `errgroup.WithContext`, cleaner API with `eg.Go()`
- **Configure connection pool** when using `database/sql`: `SetConnMaxLifetime`, `SetMaxOpenConns`, `SetMaxIdleConns`
- **Testing**: use want/got (never expected/actual), **assert** when test can continue / **require** when it should stop, **go.uber.org/gomock** for mocks

### Frontend Guidelines

Include this section in the SKILL.md only when the target language is used for frontend development (e.g., TypeScript, JavaScript). Rewrite with language-specific patterns.

- **Mobile-first approach** — implement UI components starting from mobile layout, then add responsive breakpoints for larger screens
- **Import assets from files** — never inline SVG, base64, or XML data in component code; import from separate files
- **Minimize React hooks per component** — consolidate related `useEffect` hooks that depend on the same data source instead of creating separate effects for each field
- **Use `useMemo` only with justification** — only when there is a real performance need (expensive computation or preventing unnecessary child re-renders)
- **Verify UI changes visually** — after implementing or modifying UI components, take screenshots at key breakpoints (desktop, mobile) to check for layout issues (alignment, spacing, overflow, z-index)

## Example Output

The following are examples using Go. Use these as a reference for the expected format and level of detail when generating for any language.

### Example: SKILL.md language idioms section

This section goes in the SKILL.md alongside the adapted template.

**Errors:**
- Every error must be checked or returned
- Wrap errors with context using `%w`:
  - Wrong: `return fmt.Errorf("failed: %w", err)`
  - Right: `return fmt.Errorf("find user by ID %d: %w", userID, err)`
- Use `errors.Is()` to check for specific error values, `errors.As()` to extract error types — never compare error strings

**Type Safety:**
- Use `iota` for enums instead of raw constants

**Code Organization:**
- Package names: short, lowercase, no underscores (e.g., `user`, not `user_service`)
- Use `internal/` packages to restrict visibility to the parent module

**Initialization:**
- Functional options for optional constructor parameters:
  - `type Option func(*options)`
  - `func New(required string, opts ...Option) *T { ... }`

**Concurrency:**
- Prefer `golang.org/x/sync/errgroup` over `sync.WaitGroup` — provides built-in error propagation, context cancellation via `errgroup.WithContext`, and cleaner API with `eg.Go()`

**Testing:**
- Use want/got (never expected/actual)
- Use **assert** when test can continue, **require** when it should stop
- Use **go.uber.org/gomock** for mocks

### Example: examples/http-server.md

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
        slog.ErrorContext(ctx, "create user", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### Example: examples/http-client.md

```go
func (c *VendorClient) GetProduct(ctx context.Context, id string) (*Product, error) {
    ctx, span := otel.Tracer("vendor").Start(ctx, "GetProduct")
    defer span.End()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/products/"+id, nil)
    if err != nil {
        return nil, fmt.Errorf("new request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        span.RecordError(err)
        return nil, fmt.Errorf("GET product %s: %w", id, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("GET product %s: status %d", id, resp.StatusCode)
    }

    var product Product
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, fmt.Errorf("decode product %s: %w", id, err)
    }
    return &product, nil
}
```

### Example: examples/db.md

```go
// Connection pool configuration
db.SetConnMaxLifetime(maxLifetime)
db.SetMaxOpenConns(maxOpen)
db.SetMaxIdleConns(maxIdle)

// One transaction per request for multi-table writes
func (s *OrderService) PlaceOrder(ctx context.Context, order Order) error {
    ctx, span := otel.Tracer("order").Start(ctx, "PlaceOrder")
    defer span.End()

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    if err := s.createOrder(ctx, tx, order); err != nil {
        return fmt.Errorf("create order: %w", err)
    }
    if err := s.updateInventory(ctx, tx, order.Items); err != nil {
        return fmt.Errorf("update inventory: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    return nil
}
```
