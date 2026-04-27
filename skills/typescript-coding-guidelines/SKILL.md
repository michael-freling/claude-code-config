---
name: typescript-coding-guidelines
description: TypeScript coding guidelines covering language idioms, React, Next.js, and application patterns
---

# TypeScript Coding Guidelines

## Application Pattern Examples

- [React Components](examples/react-components.md) — functional components, hooks, composition, event handling
- [Next.js Pages & Layouts](examples/nextjs-pages.md) — App Router pages, layouts, loading/error states, metadata
- [API Routes](examples/api-routes.md) — Next.js Route Handlers, middleware, request validation
- [Database](examples/db.md) — Prisma/Drizzle patterns, transactions, migrations, connection pooling
- [HTTP Client](examples/http-client.md) — `fetch` wrappers, retry logic, type-safe responses
- [Message Queue Worker](examples/message-queue.md) — consumer loops, acknowledgment, graceful shutdown

## TypeScript Idioms

**Errors:**
- Throw typed errors that extend `Error` — include the operation and identifying input in the message:
  - Wrong: `throw new Error("failed")`
  - Right: `throw new AppError("find user by ID " + userId, { cause: err })`
- Use the `cause` property (ES2022) to chain errors: `new Error("create order", { cause: dbError })`
- Define domain error classes for expected failures:
  ```typescript
  class NotFoundError extends Error {
    constructor(resource: string, id: string) {
      super(`${resource} ${id} not found`);
      this.name = "NotFoundError";
    }
  }
  ```
- Use `instanceof` checks for typed error handling — never compare `error.message` strings
- For functions that can fail without throwing, use a discriminated union return type:
  ```typescript
  type Result<T> = { ok: true; value: T } | { ok: false; error: Error };
  ```

**Type System:**
- Use `as const` for literal tuples and object literals to get the narrowest types
- Use discriminated unions over optional fields when properties are mutually exclusive:
  ```typescript
  type Event =
    | { type: "click"; x: number; y: number }
    | { type: "keypress"; key: string };
  ```
- Use `satisfies` to validate a value matches a type while preserving the narrower inferred type:
  ```typescript
  const config = { port: 3000, host: "localhost" } satisfies ServerConfig;
  ```
- Use branded types for domain IDs to prevent mixing:
  ```typescript
  type UserId = string & { readonly __brand: "UserId" };
  type OrderId = string & { readonly __brand: "OrderId" };
  ```
- Prefer `unknown` over `any` — force callers to narrow the type before use
- Use `Record<K, V>` for typed dictionaries, `Map<K, V>` when keys are not strings

**Initialization:**
- Use a config object with required + optional fields for constructors with many parameters:
  ```typescript
  interface ClientOptions {
    timeout?: number;
    retries?: number;
  }
  function createClient(baseUrl: string, options?: ClientOptions): Client
  ```
- Do not use the functional options pattern (Go idiom) — use optional fields in an options object instead
- Use factory functions over classes when instances don't need methods or `this` binding

**Async Patterns:**
- Always `await` promises — never leave a floating promise without `void` annotation
- Use `Promise.all()` for independent concurrent operations, `Promise.allSettled()` when partial failure is acceptable
- Use `AbortController` for cancellation — pass `AbortSignal` to `fetch` and long-running operations
- Never use `new Promise()` when an `async` function suffices — the constructor is for wrapping callback APIs

**Testing:**
- Use **Vitest** as the test runner — compatible with Jest API, native TypeScript support, fast
- Use `describe`/`it` blocks with descriptive names; group related cases with `describe`, separate success from error cases
- Use `vi.fn()` for function mocks, `vi.spyOn()` for method spies
- Use `vi.useFakeTimers()` for time-dependent tests — call `vi.setSystemTime()` for deterministic dates
- Use `@testing-library/react` for component tests — query by role, label, and text, not by CSS selectors or test IDs
- Use `msw` (Mock Service Worker) to intercept HTTP requests in tests — provides realistic network mocking

## Principles

- **Environment-agnostic code** — use `process.env` via a typed config module, never branch on `process.env.NODE_ENV` in application code; inject different behavior through configuration objects or dependency injection

## Simplicity

**Simplicity is the most important thing.** Minimize the final complexity of the codebase — not the size of the change.

- **Breaking changes are acceptable** unless backward compatibility is explicitly required by published package consumers or project constraints
- **Avoid premature abstraction** — do not create a generic component, hook, or utility until there are 3+ concrete use cases; three similar JSX blocks are better than a premature `<GenericCard>` abstraction
- **Minimize mutable state** — prefer `const`, derived values with `useMemo`/computed getters, and pure functions over `let` and accumulated state; keep `useState` as local as possible

### High Cohesion

Each function should do one thing. When a new requirement applies to only some callers, split the function rather than adding conditional logic.

- **Wrong**: `sendNotification()` handles email, SMS, and push — adding rate limiting for SMS forces all paths through the rate limiter
- **Right**: `sendEmail()`, `sendSMS()`, `sendPush()` — rate limiting lives only in `sendSMS()`

### Separate Functions over Boolean Arguments

When a boolean argument causes a function to take substantially different paths, split into separate functions.

- **Wrong**: `completeJob(success: boolean)` with branching logic based on the flag
- **Right**: `completeJob()` and `failJob()` as separate exported functions

### Narrow Types

Accept the minimum type a function needs — TypeScript's structural typing makes this natural.

- **Wrong**: a function that only needs a user's name and email accepts `User` (which also has `password`, `createdAt`, etc.)
- **Right**: the function accepts `Pick<User, "name" | "email">` or a dedicated interface with just the fields it needs

## Error Handling

- **Always handle errors explicitly** — every `try/catch` must handle or rethrow with context; never write empty `catch {}` blocks
- **Wrap errors with context** using the `cause` property: `throw new Error("create order for user " + userId, { cause: err })`
- **Fail fast on invalid state** — validate at entry points and throw immediately; use TypeScript's type system to make invalid states unrepresentable
- **Type-narrow `unknown` in catch blocks** — `catch (err)` gives `unknown`; narrow with `instanceof Error` before accessing properties

### Server-Side Errors

- Map domain errors to HTTP status codes — `NotFoundError` to `404`, `ValidationError` to `400`, unexpected errors to `500`
- Return structured error responses: `{ error: { message: string; code: string } }` — never expose stack traces or internal messages
- Log the full error server-side with structured logger before returning the sanitized response

### Client-Side Errors

- Retry on `429`, `502`, `503`, `504` and network errors (`TypeError` from `fetch`) with exponential backoff
- Make retry count, backoff base, and retryable status codes configurable via options object
- If the infrastructure layer (API gateway, service mesh) already retries, do not duplicate

## Dependencies

- Use `npm install <package>@latest` and verify compatibility — check `package.json` `engines` field and peer dependencies
- Use exact versions (`"5.3.2"` not `"^5.3.2"`) in application projects; use ranges only in published libraries
- `package-lock.json` must be committed — it ensures reproducible installs

## Code Readability

- **Specific names** — `userCount` not `n`, `fetchProducts` not `getData`. Exception: arrow function parameters in short callbacks (`items.map(x => x.id)`) and loop indices
- **Early returns** — guard clauses at the top of a function; keep the happy path at the base indentation:
  ```typescript
  if (!user) throw new NotFoundError("user", id);
  if (!user.isActive) throw new ForbiddenError("user is inactive");
  // main logic here
  ```
- **Single level of abstraction** — a route handler calls service functions, it does not also write SQL queries
- **Delete dead code** — remove unused functions, unreachable branches, commented-out code; TypeScript compiler flags unused locals, but watch for unused exports
- **No redundant initializations** — do not write `let count: number = 0` when `let count = 0` is inferred; do not initialize optional fields to `undefined`
- **Options object over long parameter lists** — more than 3 parameters signal the need for a named options type:
  ```typescript
  interface CreateOrderOptions {
    userId: string;
    items: OrderItem[];
    couponId?: string;
    note?: string;
  }
  ```
- **Load binary data from files** — import images, fonts, and certificates from files via bundler imports or `fs.readFile`; never inline base64 strings in source
- **Validate only at boundaries** — validate HTTP request bodies (with Zod), external API responses, and user input; do not add runtime checks for types already guaranteed by TypeScript's compiler

## Code Reuse

- **Single definition, owned by the provider** — define types, interfaces, and functions in the module that implements the behavior; consumers import, not redefine
- **Named constants** — `const MAX_RETRIES = 3` instead of magic `3` throughout the code
- **Search before writing** — use project-wide search for existing helpers, types, or components before creating new ones
- **Refactor before extending** — understand the existing code, extract reusable parts, then add the new feature
- **Refactor after QA** — once the feature works, clean up for reusability and consistency

## Code Organization

- **Minimal entry point** — `main.ts` or `index.ts` initializes config, database connections, and server startup; no business logic in the entry file
- **Boundary layers by type** — `api/` for route handlers, `client/` for external API wrappers; **domain logic by feature** — `user/`, `billing/`, `order/`
- **Dependency direction** — handlers import services, services import repositories — never the reverse; no circular imports (TypeScript will error on circular type-only imports but not value imports — enforce with ESLint `import/no-cycle`)
- **Minimize exports** — default to non-exported; export only what other modules consume; use barrel files (`index.ts`) sparingly — they break tree-shaking and cause circular import issues

## Type Safety

- **Discriminated unions for known sets** — use string literal unions with a `type` discriminant field instead of raw strings:
  ```typescript
  type OrderStatus = "pending" | "confirmed" | "shipped" | "delivered";
  ```
- **Use `enum` only for numeric flags** — prefer string literal unions for most cases; `enum` creates runtime objects and has surprising behavior with reverse mappings
- **`satisfies` for config objects** — validates shape at compile time while preserving literal types for autocomplete

## Initialization

- **Required parameters as positional arguments**, optional parameters via options object:
  ```typescript
  function createServer(port: number, options?: ServerOptions): Server
  ```
- **Validate invariants in factories/constructors** — throw if required fields are missing or invalid:
  ```typescript
  function createClient(baseUrl: string): Client {
    if (!baseUrl) throw new Error("baseUrl is required");
    ...
  }
  ```
- Do not use the builder pattern in TypeScript — options objects are simpler and type-checked at the call site

## External I/O

- **Batch operations** — use Prisma's `createMany`/`findMany`, Drizzle's bulk insert, or SQL `INSERT ... VALUES` with arrays to reduce round-trips
- **Transactions** — use `prisma.$transaction()` or Drizzle's `db.transaction()` for multi-table writes within a single request
- **Ensure indexes** — add `@@index` in Prisma schema or `CREATE INDEX` for columns in `WHERE`, `JOIN`, `ORDER BY` on tables with 10K+ rows
- **Pass values from code** — compute `new Date()` in TypeScript and pass as a parameter; do not use `NOW()` in SQL — enables deterministic tests with injected clocks
- **Avoid redundant queries** — do not refetch data already available in the current request context
- **Parameterized queries** — use Prisma's query builder, Drizzle's `sql` template tag, or `$1` placeholders; never interpolate user input into SQL strings

## Logging

- **Use `pino` or `winston`** for structured logging — key-value JSON output, not `console.log` string concatenation
- **Log levels**: `debug` (development traces), `info` (operational events), `warn` (expected/recoverable issues), `error` (unexpected failures)
- **Correlation IDs** — extract trace/span IDs from the OpenTelemetry context and include in log entries; use a request-scoped logger with `requestId` attached
- **Never log secrets** — no tokens, passwords, API keys, PII, or full request/response bodies

## Testing

- **Deterministic tests** — no `setTimeout` waits, no reliance on wall-clock time; use `vi.useFakeTimers()` and `vi.setSystemTime()` for time-dependent logic
- **Inject dependencies** — pass services and clients as constructor/function parameters; use Vitest's `vi.fn()` to create typed mocks:
  ```typescript
  const mockUserService = { getById: vi.fn<[string], Promise<User>>() };
  ```
- **Extract interfaces for external dependencies** — define a `UserRepository` interface and mock it in tests instead of mocking Prisma internals
- **Clean up state** — use `beforeEach` to reset mocks (`vi.resetAllMocks()`), truncate test database tables, and reinitialize fixtures
- **Test case arrays** — group related cases in arrays with descriptive names:
  ```typescript
  const cases = [
    { name: "valid email", input: "a@b.com", expected: true },
    { name: "missing @", input: "ab.com", expected: false },
  ];
  cases.forEach(({ name, input, expected }) => {
    it(name, () => { expect(isValidEmail(input)).toBe(expected); });
  });
  ```
- **Real databases over mocks** — use `testcontainers` to spin up PostgreSQL in CI; mock only third-party HTTP APIs with `msw`
- **Shared test setup** — define factory functions (`createTestUser()`, `createTestOrder()`) with sensible defaults in a shared `test/factories.ts`
- **Never `it.skip()` failures** — fix flaky tests immediately
- **Load test data from files** — use `fs.readFile("test/fixtures/sample.json")` instead of inline JSON blobs

### What to Test

- All exported functions and React components
- Error handling paths — each `catch` branch, each error response
- Edge cases — empty arrays, `null`, `undefined`, zero, `Number.MAX_SAFE_INTEGER`
- Boundary conditions — pagination limits, timeout edges, max string length

### What NOT to Test

- TypeScript type correctness (the compiler handles this)
- Third-party library behavior (`zod.parse` correctness, React rendering internals)
- Trivial re-exports or type-only modules

### Coverage Gaps to Look For

- Untested `catch` blocks
- Missing edge cases for `null`/`undefined` inputs
- `switch`/`if-else` branches not fully covered

### Writing Testable Code

1. **Analyze testability** — identify direct imports of singletons (`prisma`, `fetch`), global state, and tight coupling
2. **Refactor first** — extract interfaces, inject dependencies via function parameters, split large functions
3. **Write tests** — use test case arrays, cover success/error/edge cases
4. **Verify coverage** — `vitest --coverage` and check branch coverage

## Concurrency

- **`Promise.all()`** for independent concurrent operations — cancels nothing on failure (use `Promise.allSettled()` if partial results are acceptable)
- **`Promise.all()` with mapped arrays** when a batch API is not available:
  ```typescript
  const users = await Promise.all(ids.map(id => fetchUser(id)));
  ```
- **`AbortController`** for cancellation — pass `signal` to `fetch`, database queries, and long-running loops that check `signal.aborted`
- **Worker threads** (`worker_threads`) for CPU-bound work — avoid blocking the event loop; use `Atomics` and `SharedArrayBuffer` only when absolutely necessary
- **Message queues** (BullMQ, SQS) for fire-and-forget work — decouple notification sending, report generation, and other async tasks from the request path

## Security

- **No hardcoded secrets** — use `process.env` via a typed config module; never commit `.env` files with real credentials
- **Sanitize at trust boundaries** — validate all user input with Zod schemas in route handlers; escape output for HTML with framework-provided utilities (React auto-escapes JSX, but `dangerouslySetInnerHTML` bypasses this); use parameterized SQL

## Monitoring

- **OpenTelemetry SDK** — initialize `NodeTracerProvider`, `MeterProvider` in the application bootstrap; use `@opentelemetry/api` for `tracer.startActiveSpan()` in application code
- **RED metrics** at minimum for every server:
  - **Rate**: requests per second (counter)
  - **Errors**: error responses per second (counter with status label)
  - **Duration**: request latency histogram
- **Avoid high-cardinality labels** — do not use user IDs, request IDs, or unbounded strings as metric attributes; use status codes, method names, and route patterns

## Frontend Guidelines

- **Mobile-first approach** — write CSS/Tailwind starting from the smallest viewport; add `md:`, `lg:` breakpoints for larger screens
- **Import assets from files** — never inline SVG markup, base64 data URIs, or XML in component files; import from `.svg` files or use an icon library component
- **Minimize hooks per component** — consolidate related `useEffect` calls that depend on the same data source into a single effect; extract complex hook logic into custom hooks
- **`useMemo` only with justification** — only when profiling shows an expensive computation or to prevent unnecessary child re-renders via referential equality; the default is no memoization
- **Verify UI changes visually** — after modifying components, check layouts at mobile (375px), tablet (768px), and desktop (1280px) breakpoints for alignment, spacing, overflow, and z-index issues
