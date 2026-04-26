---
name: typescript-coding-guidelines
description: TypeScript coding guidelines covering language idioms, React, Next.js, and application patterns
---

# TypeScript Coding Guidelines

## Example Files

- `examples/react-components.md` — React component patterns, hooks, state management, performance
- `examples/nextjs-pages.md` — Next.js App Router, server components, server actions, data fetching
- `examples/api-routes.md` — Next.js API route handlers, middleware, validation, error responses
- `examples/db.md` — Prisma patterns, transactions, batch queries
- `examples/http-client.md` — fetch-based client with AbortSignal, typed responses

## TypeScript Language Idioms

### Errors

- Use typed error classes extending `Error` for domain errors
- Always set `cause` when wrapping: `throw new AppError("find user", { cause: err })`
- Use discriminated unions for expected failures:
  ```typescript
  type Result<T> = { ok: true; value: T } | { ok: false; error: AppError }
  ```
- Use `unknown` in catch clauses, narrow before accessing:
  ```typescript
  catch (err: unknown) {
    if (err instanceof AppError) { /* handle */ }
    throw new AppError("unexpected", { cause: err });
  }
  ```
- Use `instanceof` checks for error class hierarchy — the TypeScript equivalent of `errors.Is()`/`errors.As()`

### Type Safety

- Use `as const` for literal types and const enums:
  ```typescript
  const ROLES = ["admin", "user", "guest"] as const;
  type Role = (typeof ROLES)[number]; // "admin" | "user" | "guest"
  ```
- Prefer discriminated unions over raw string unions for exhaustive checking
- Use branded types for type-safe IDs:
  ```typescript
  type UserId = string & { readonly __brand: "UserId" }
  ```
- Use `satisfies` operator to validate types without widening:
  ```typescript
  const config = { port: 3000, host: "localhost" } satisfies ServerConfig;
  ```
- Use `readonly` and `Readonly<T>` for immutable data

### Code Organization

- Use barrel exports (`index.ts`) sparingly — only for public API surfaces
- Prefer named exports over default exports
- Use `internal/` convention or `package.json` `exports` field to restrict visibility

### Initialization

- Options object for optional parameters:
  ```typescript
  interface CreateClientOptions { timeout?: number; retries?: number }
  function createClient(baseUrl: string, options?: CreateClientOptions): Client
  ```
- Private constructor + static factory for validated construction:
  ```typescript
  class EmailAddress {
    private constructor(private readonly value: string) {}
    static create(input: string): EmailAddress {
      if (!EMAIL_REGEX.test(input)) throw new ValidationError("invalid email");
      return new EmailAddress(input);
    }
  }
  ```

### Concurrency

- `Promise.all()` for independent parallel operations
- `Promise.allSettled()` when partial failures are acceptable
- `p-limit` for concurrency control:
  ```typescript
  import pLimit from "p-limit";
  const limit = pLimit(5);
  await Promise.all(urls.map(url => limit(() => fetch(url))));
  ```
- `AbortController`/`AbortSignal` for cancellation

### Testing

- **Vitest** as test runner, `describe`/`it` with descriptive names
- `vi.fn()`/`vi.spyOn()` for mocks; prefer dependency injection over module mocking
- `test.each` for table-driven tests:
  ```typescript
  test.each([
    { input: "", expected: false },
    { input: "valid@email.com", expected: true },
  ])("validates $input -> $expected", ({ input, expected }) => {
    expect(isValidEmail(input)).toBe(expected);
  });
  ```
- **testcontainers** for real database testing
- `beforeEach` for test isolation, `vi.restoreAllMocks()` in `afterEach`

### Dependencies

- `strict: true`, `noUncheckedIndexedAccess: true` in `tsconfig.json`
- Pin exact versions with `npm install --save-exact` or `"exact": true` in `.npmrc`

---

## Principles

- Use dependency injection via constructor parameters or factory functions — not `process.env` reads scattered through business logic. Access env vars once in the entry point (`main.ts`, `server.ts`, or Next.js `instrumentation.ts`), validate with Zod, and pass typed config objects down:
  ```typescript
  const envSchema = z.object({
    DATABASE_URL: z.string().url(),
    REDIS_URL: z.string().url(),
    LOG_LEVEL: z.enum(["debug", "info", "warn", "error"]).default("info"),
  });
  const config = envSchema.parse(process.env);
  const userService = new UserService(prisma, config.LOG_LEVEL);
  ```
- Use `tsconfig.json` `paths` aliases instead of environment-conditional imports:
  ```json
  { "compilerOptions": { "paths": { "@/lib/*": ["./lib/*"], "@/components/*": ["./components/*"] } } }
  ```

## Simplicity

Simplicity is the most important thing. Minimize the final complexity of the codebase — not the size of the change.

- Prefer `readonly` properties and `Readonly<T>` over mutable state — derive values with `.map()`, `.filter()`, `.reduce()` instead of accumulating into mutable arrays:
  ```typescript
  // Wrong: mutable accumulation
  const result: string[] = [];
  for (const user of users) {
    if (user.active) result.push(user.name);
  }
  // Right: derived value
  const result = users.filter(u => u.active).map(u => u.name);
  ```
- Use TypeScript's type system to make illegal states unrepresentable — discriminated unions eliminate entire categories of runtime checks:
  ```typescript
  // Wrong: nullable fields with implicit constraints
  interface Order { status: string; shippedAt: Date | null; trackingNumber: string | null }

  // Right: discriminated union — shipped orders always have tracking
  type Order =
    | { status: "pending"; items: OrderItem[] }
    | { status: "shipped"; items: OrderItem[]; shippedAt: Date; trackingNumber: string }
    | { status: "cancelled"; items: OrderItem[]; cancelledAt: Date; reason: string }
  ```
- Three similar `.map()` chains are better than a premature `genericTransform<T>()` abstraction
- Breaking changes are acceptable unless backward compatibility is explicitly required

### High Cohesion

Each function should do one thing. When a new requirement applies to only some callers, split the function rather than adding conditional logic.

- Wrong: `sendNotification()` handles email, SMS, push — adding rate limiting for SMS forces all paths through the rate limiter check
- Right: separate `sendEmail()`, `sendSMS()`, `sendPush()` — rate limiting lives only in `sendSMS()`
- In React: wrong is a `<FormField>` that handles text, select, checkbox, date via a `type` prop with internal branching — right is separate `<TextField>`, `<SelectField>`, `<CheckboxField>`, `<DateField>` components

### Separate Methods over Boolean Arguments

When a boolean argument causes a function to branch into substantially different logic, split into separate functions.

- Wrong: `completeJob(success: boolean)` with branching logic based on the flag
- Right: `completeJob()` and `failJob()` as separate functions
- In React: wrong is `<Button primary={true}>` that branches into completely different rendering and event handling — right is `<PrimaryButton>` and `<SecondaryButton>` if the logic differs substantially

### Narrow Interfaces

Accept the minimum type a function needs. This reduces coupling and makes the function usable in more contexts.

- Use `Pick<T, 'needed' | 'fields'>` to narrow function parameters to only what's needed:
  ```typescript
  // Wrong: accepts full User when it only reads name and email
  function sendWelcome(user: User): void

  // Right: accepts only what it needs
  function sendWelcome(user: Pick<User, "name" | "email">): void
  ```
- Use `Readonly<T>` when a function only reads data
- In React: accept `ReactNode` instead of specific component types; accept `ComponentPropsWithoutRef<"button">` for wrapper components

## Error Handling

- Never use `any` in catch blocks — always `catch (err: unknown)` and narrow with `instanceof`:
  ```typescript
  try {
    await saveOrder(order);
  } catch (err: unknown) {
    if (err instanceof ValidationError) {
      return { error: err.message };
    }
    throw new AppError("save order", { cause: err });
  }
  ```
- Wrap errors with `{ cause: err }` for stack trace chaining (ES2022+)
- Use `Result<T>` discriminated union for expected failures instead of exceptions:
  ```typescript
  type Result<T> = { ok: true; value: T } | { ok: false; error: AppError }

  function parseConfig(raw: string): Result<Config> {
    try {
      const parsed = JSON.parse(raw);
      const config = configSchema.parse(parsed);
      return { ok: true, value: config };
    } catch (err: unknown) {
      return { ok: false, error: new AppError("parse config", { cause: err }) };
    }
  }
  ```
- Server-side: Next.js `error.tsx` boundaries for page-level errors, try-catch in server actions returning `{ error: string }` to the client
- Client-side: check `response.ok` after `fetch`, retry transient failures with exponential backoff via `p-retry`:
  ```typescript
  import pRetry from "p-retry";
  const data = await pRetry(() => fetchUser(userId), { retries: 3 });
  ```

## Dependencies

- Pin exact versions: `npm install --save-exact` or set `"exact": true` in `.npmrc`
- Use `npm audit` or `pnpm audit` to check for known vulnerabilities
- Prefer packages with TypeScript types built-in (check for `types` field in `package.json`) over separate `@types/*` packages when available
- Enable `strict: true` and `noUncheckedIndexedAccess: true` in `tsconfig.json`

## Code Readability

- Use specific names: `userRepository` not `repo`, `OrderStatus` not `Status`, `formatCurrency()` not `format()`
- Use guard clauses with early `return` or early `throw` — avoid nested `if/else` chains:
  ```typescript
  function getDiscount(user: User): number {
    if (!user.active) return 0;
    if (!user.subscription) return 0;
    if (user.subscription.plan === "free") return 0;
    return calculatePlanDiscount(user.subscription);
  }
  ```
- Use `satisfies` for inline type validation without losing literal types:
  ```typescript
  const config = { port: 3000, host: "localhost" } satisfies ServerConfig;
  ```
- Avoid parameter lists longer than 5 — use an options object:
  ```typescript
  interface CreateUserOptions {
    name: string;
    email: string;
    role: Role;
    teamId?: string;
    avatarUrl?: string;
    locale?: string;
  }
  function createUser(options: CreateUserOptions): Promise<User>
  ```
- Use `fs.readFile()` to load binary data (images, certificates) instead of embedding base64 strings inline
- Validate with Zod at API boundaries (incoming requests, external API responses); trust internal types elsewhere:
  ```typescript
  const createUserSchema = z.object({
    name: z.string().min(1).max(100),
    email: z.string().email(),
    role: z.enum(["admin", "user"]),
  });
  // In API route handler:
  const body = createUserSchema.parse(await request.json());
  ```
- Minimal comments — only high-level explanations of purpose or non-obvious decisions. No line-by-line comments
- Delete dead code — remove unused functions, unreachable branches, commented-out code, and unused imports when encountered

## Code Reuse

- Define types in the module that implements them — consumers import, never redefine:
  ```typescript
  // lib/user/types.ts — defines User, UserRole
  // lib/billing/service.ts — imports User from "@/lib/user/types", never redefines it
  ```
- Use `as const` objects for shared constants:
  ```typescript
  export const HTTP_STATUS = { OK: 200, NOT_FOUND: 404, INTERNAL: 500 } as const;
  type HttpStatus = (typeof HTTP_STATUS)[keyof typeof HTTP_STATUS];
  ```
- In React: extract shared hooks (`useDebounce`, `usePagination`) into a `hooks/` directory; extract shared UI into a `components/` directory
- Search for existing code before adding new — reuse or extend existing patterns rather than duplicating
- Refactor before adding new changes — understand existing code, identify reusable modules, refactor to enable reuse, then add new functionality

## Code Organization

- Entry point (`main.ts`, `server.ts`, or Next.js `instrumentation.ts`) only initializes: DB connections, OTel SDK, config parsing — keep business logic out
- Boundary layers: `app/api/` (Next.js API routes), `lib/db/` (Prisma client and migrations), `lib/vendors/` (3rd party SDK wrappers)
- Domain logic by feature: `lib/user/`, `lib/billing/`, `lib/order/` — not `lib/services/`, `lib/repositories/`
- In React/Next.js: `components/` for shared UI, `app/` for pages and routes, `lib/` for business logic
- Use `package.json` `exports` field or `internal/` directories to control module visibility
- Dependency direction: higher-level modules depend on lower-level modules, never the reverse; avoid circular dependencies

## Type Safety

- Use discriminated unions for state machines:
  ```typescript
  type ConnectionState =
    | { status: "connecting" }
    | { status: "connected"; socket: WebSocket }
    | { status: "error"; error: Error }
  ```
- Use branded types for IDs to prevent mixing:
  ```typescript
  type UserId = string & { readonly __brand: "UserId" }
  type OrderId = string & { readonly __brand: "OrderId" }

  function createUserId(id: string): UserId { return id as UserId; }
  function getUser(id: UserId): Promise<User> // cannot accidentally pass OrderId
  ```
- Use `satisfies` to narrow without widening:
  ```typescript
  const routes = { home: "/", about: "/about" } satisfies Record<string, string>;
  // typeof routes.home is "/" (literal), not string
  ```
- Use `readonly` arrays and `Readonly<T>` for data that should not be mutated:
  ```typescript
  function processItems(items: readonly OrderItem[]): Summary
  ```
- Use `unknown` instead of `any` — narrow with type guards:
  ```typescript
  function isAppError(err: unknown): err is AppError {
    return err instanceof AppError;
  }
  ```
- Prefer enums or `as const` objects over raw strings for known sets of values

## Initialization

- Required parameters as positional args, optional as options object:
  ```typescript
  interface ClientOptions { timeout?: number; retries?: number; signal?: AbortSignal }
  function createClient(baseUrl: string, options?: ClientOptions): Client
  ```
- Use private constructor + static `create()` factory for validated construction:
  ```typescript
  class EmailAddress {
    private constructor(private readonly value: string) {}
    static create(input: string): EmailAddress {
      if (!EMAIL_REGEX.test(input)) throw new ValidationError("invalid email");
      return new EmailAddress(input);
    }
    toString(): string { return this.value; }
  }
  ```
- Do not introduce optional parameter patterns when all parameters are required

## External I/O

- Use Prisma's `createMany()`, `findMany({ where: { id: { in: ids } } })` for batch operations — not loops of `findUnique()`:
  ```typescript
  // Wrong: N+1 queries
  const users = await Promise.all(ids.map(id => prisma.user.findUnique({ where: { id } })));
  // Right: single batch query
  const users = await prisma.user.findMany({ where: { id: { in: ids } } });
  ```
- Use `prisma.$transaction()` interactive transactions for multi-table writes:
  ```typescript
  await prisma.$transaction(async (tx) => {
    const order = await tx.order.create({ data: orderData });
    await tx.inventory.updateMany({
      where: { productId: { in: order.items.map(i => i.productId) } },
      data: { reserved: { increment: 1 } },
    });
  });
  ```
- Pass `new Date()` from application code, not `Prisma.sql\`NOW()\`` — inject a `Clock` interface for testing
- Use parameterized queries with `Prisma.sql` tagged template for raw SQL — never string concatenation:
  ```typescript
  const users = await prisma.$queryRaw(Prisma.sql`SELECT * FROM users WHERE id = ${userId}`);
  ```
- Use `AbortSignal.timeout()` for fetch timeouts, `AbortSignal.any()` to compose with caller cancellation:
  ```typescript
  const response = await fetch(url, {
    signal: AbortSignal.any([AbortSignal.timeout(5000), callerSignal]),
  });
  ```

## Logging

- Use `pino` for structured JSON logging:
  ```typescript
  import pino from "pino";
  const logger = pino({ level: process.env.LOG_LEVEL ?? "info" });
  logger.info({ userId, orderId }, "order placed");
  ```
- Log levels: `debug` (dev only), `info` (operational events), `warn` (recoverable issues), `error` (requires investigation)
- Attach trace context for distributed tracing:
  ```typescript
  import { trace } from "@opentelemetry/api";
  const span = trace.getActiveSpan();
  const childLogger = logger.child({ traceId: span?.spanContext().traceId });
  ```
- Never log: tokens, passwords, `process.env`, full request/response bodies, PII

## Testing

- Use Vitest with `test.each` for table-driven tests:
  ```typescript
  test.each([
    { input: "", expected: false },
    { input: "valid@email.com", expected: true },
  ])("validates $input -> $expected", ({ input, expected }) => {
    expect(isValidEmail(input)).toBe(expected);
  });
  ```
- Inject dependencies via constructor — mock only 3rd party APIs:
  ```typescript
  // Production
  const service = new OrderService(prisma, stripeClient, logger);
  // Test
  const service = new OrderService(testPrisma, mockStripe, testLogger);
  ```
- Use testcontainers for real Postgres in integration tests:
  ```typescript
  import { PostgreSqlContainer } from "@testcontainers/postgresql";
  const container = await new PostgreSqlContainer().start();
  const prisma = new PrismaClient({ datasources: { db: { url: container.getConnectionUri() } } });
  ```
- In React: use `@testing-library/react` with `render()`, `screen.getByRole()`, `userEvent` — test behavior, not implementation details
- Inject time via a `Clock` interface — never mock `Date.now()` globally:
  ```typescript
  interface Clock { now(): Date }
  const realClock: Clock = { now: () => new Date() };
  const testClock: Clock = { now: () => new Date("2024-01-15T10:00:00Z") };
  ```
- Clean up: `beforeEach` to reset DB state, `afterEach` with `vi.restoreAllMocks()`

## Concurrency

- JavaScript is single-threaded — shared state does not need locks unless using Worker threads
- Use `Promise.all()` when batching is not possible:
  ```typescript
  const results = await Promise.all(ids.map(id => fetchUser(id)));
  ```
- Use `p-limit` to control concurrency and prevent overwhelming external services:
  ```typescript
  import pLimit from "p-limit";
  const limit = pLimit(5);
  await Promise.all(urls.map(url => limit(() => fetch(url))));
  ```
- Use `Promise.allSettled()` when partial failures are acceptable:
  ```typescript
  const results = await Promise.allSettled(ids.map(id => fetchUser(id)));
  const succeeded = results.filter((r): r is PromiseFulfilledResult<User> => r.status === "fulfilled");
  ```
- Use message queues (BullMQ with Redis) for fire-and-forget work like sending emails or generating reports

## Security

- Never commit secrets — use `process.env` read once at startup, validate with Zod:
  ```typescript
  const envSchema = z.object({
    DATABASE_URL: z.string().url(),
    JWT_SECRET: z.string().min(32),
  });
  const env = envSchema.parse(process.env);
  ```
- Sanitize user input at API boundaries with Zod schemas — reject, do not coerce
- Use `DOMPurify` for HTML sanitization in React when rendering user content with `dangerouslySetInnerHTML`:
  ```typescript
  import DOMPurify from "dompurify";
  <div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(userHtml) }} />
  ```

## Monitoring

- Use `@opentelemetry/sdk-node` for auto-instrumentation at bootstrap (`instrumentation.ts`):
  ```typescript
  import { NodeSDK } from "@opentelemetry/sdk-node";
  import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
  const sdk = new NodeSDK({ instrumentations: [getNodeAutoInstrumentations()] });
  sdk.start();
  ```
- RED metrics via `@opentelemetry/api`: counter for requests, histogram for duration, counter for errors:
  ```typescript
  import { metrics } from "@opentelemetry/api";
  const meter = metrics.getMeter("api");
  const requestCounter = meter.createCounter("http.requests");
  const requestDuration = meter.createHistogram("http.request.duration", { unit: "ms" });
  ```
- Avoid high-cardinality labels — do not use user IDs or request IDs as metric labels
- Use `trace.getTracer("service-name")` and `tracer.startActiveSpan()` for custom spans:
  ```typescript
  const tracer = trace.getTracer("order-service");
  await tracer.startActiveSpan("processOrder", async (span) => {
    try {
      const result = await processOrder(orderId);
      span.setAttributes({ "order.id": orderId, "order.total": result.total });
      return result;
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw err;
    } finally {
      span.end();
    }
  });
  ```
