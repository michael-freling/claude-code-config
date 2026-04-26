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

- Typed error classes extending `Error` for domain errors:
  ```typescript
  class AppError extends Error {
    constructor(message: string, options?: ErrorOptions) {
      super(message, options);
      this.name = "AppError";
    }
  }
  class NotFoundError extends AppError {
    constructor(entity: string, id: string, options?: ErrorOptions) {
      super(`${entity} not found: ${id}`, options);
      this.name = "NotFoundError";
    }
  }
  ```
- Always set `cause` when wrapping: `throw new AppError("find user", { cause: err })`
- Discriminated unions for expected failures:
  ```typescript
  type Result<T> = { ok: true; value: T } | { ok: false; error: AppError }
  ```
- `unknown` in catch clauses, narrow with `instanceof`:
  ```typescript
  catch (err: unknown) {
    if (err instanceof NotFoundError) { /* handle missing entity */ }
    throw new AppError("unexpected", { cause: err });
  }
  ```

### Type Safety

- `as const` for literal types:
  ```typescript
  const ROLES = ["admin", "user", "guest"] as const;
  type Role = (typeof ROLES)[number]; // "admin" | "user" | "guest"
  ```
- Discriminated unions over raw string unions for exhaustive checking
- Branded types for type-safe IDs:
  ```typescript
  type UserId = string & { readonly __brand: "UserId" }
  function userId(raw: string): UserId { return raw as UserId; }
  ```
- `satisfies` to validate without widening:
  ```typescript
  const config = { port: 3000, host: "localhost" } satisfies ServerConfig;
  ```
- `readonly` and `Readonly<T>` for immutable data:
  ```typescript
  function processItems(items: readonly OrderItem[]): void { /* ... */ }
  ```

### Code Organization

- Barrel exports (`index.ts`) sparingly — only for public API surfaces
- Named exports over default exports
- `internal/` convention or `package.json` `exports` field to restrict visibility:
  ```json
  { "exports": { ".": "./src/index.ts", "./internal": null } }
  ```

### Initialization

- Options object for optional params:
  ```typescript
  interface CreateClientOptions { timeout?: number; retries?: number; logger?: Logger }
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
    toString(): string { return this.value; }
  }
  ```

### Concurrency

- `Promise.all()` for independent parallel ops
- `Promise.allSettled()` for partial-failure tolerance
- `p-limit` for concurrency control:
  ```typescript
  import pLimit from "p-limit";
  const limit = pLimit(5);
  const results = await Promise.all(urls.map(u => limit(() => fetch(u))));
  ```
- `AbortController`/`AbortSignal` for cancellation:
  ```typescript
  const controller = new AbortController();
  const response = await fetch(url, { signal: AbortSignal.timeout(5000) });
  ```

### Testing

- **Vitest** with `describe`/`it`, `vi.fn()`/`vi.spyOn()`, `test.each`
- **testcontainers** for real DB testing
- `beforeEach` for isolation, `vi.restoreAllMocks()` in `afterEach`

### Dependencies

- `strict: true`, `noUncheckedIndexedAccess: true` in `tsconfig.json`
- `--save-exact` for pinned versions or `"exact": true` in `.npmrc`

---

## Principles

- DI via constructor/factory — not scattered `process.env` reads. Parse env once at entry with Zod, pass typed config:
  ```typescript
  import { z } from "zod";

  const envSchema = z.object({
    DATABASE_URL: z.string().url(),
    LOG_LEVEL: z.enum(["debug", "info", "warn", "error"]).default("info"),
    PORT: z.coerce.number().default(3000),
  });
  export type Env = z.infer<typeof envSchema>;
  export const config = envSchema.parse(process.env);
  ```
- `tsconfig.json` `paths` aliases over environment-conditional imports:
  ```json
  { "compilerOptions": { "paths": { "@/lib/*": ["./src/lib/*"], "@/components/*": ["./src/components/*"] } } }
  ```
- Inject dependencies through constructors or factory functions — test by passing stubs, not by patching globals

## Simplicity

- `readonly` and `Readonly<T>` over mutable state — derive with `.map()`, `.filter()` not mutable accumulation:
  ```typescript
  // Wrong
  const result: string[] = [];
  for (const u of users) { if (u.active) result.push(u.name); }

  // Right
  const result = users.filter(u => u.active).map(u => u.name);
  ```
- Discriminated unions make illegal states unrepresentable — eliminate runtime checks:
  ```typescript
  type Order =
    | { status: "pending"; items: OrderItem[] }
    | { status: "shipped"; items: OrderItem[]; shippedAt: Date; trackingNumber: string }
    | { status: "cancelled"; items: OrderItem[]; cancelledAt: Date; reason: string };

  function getTrackingUrl(order: Extract<Order, { status: "shipped" }>): string {
    return `https://track.example.com/${order.trackingNumber}`;
  }
  ```
- Three similar `.map()` chains beat a premature `genericTransform<T>()`
- Breaking changes acceptable unless backward compat explicitly required

### High Cohesion

- Wrong: `sendNotification()` handles email/SMS/push — right: `sendEmail()`, `sendSMS()`, `sendPush()`
- React: wrong is `<FormField type="text"|"select"|"checkbox">` branching — right: `<TextField>`, `<SelectField>`, `<CheckboxField>`

### Separate Methods over Boolean Arguments

- Wrong: `completeJob(success: boolean)` — right: `completeJob()` and `failJob()`
- React: wrong is `<Button primary={true}>` branching entirely — right: `<PrimaryButton>` and `<SecondaryButton>` if logic differs substantially

### Narrow Interfaces

- `Pick<T, 'needed' | 'fields'>` to narrow params:
  ```typescript
  // Wrong
  function sendWelcome(user: User): void { /* only uses name and email */ }

  // Right
  function sendWelcome(user: Pick<User, "name" | "email">): void
  ```
- `Readonly<T>` when only reading
- React: `ReactNode` over specific types; `ComponentPropsWithoutRef<"button">` for wrappers:
  ```tsx
  interface ButtonProps extends ComponentPropsWithoutRef<"button"> {
    variant?: "primary" | "secondary";
  }
  ```

## Error Handling

- Always `catch (err: unknown)` — narrow with `instanceof`:
  ```typescript
  try {
    await saveOrder(order);
  } catch (err: unknown) {
    if (err instanceof ValidationError) return { error: err.message };
    throw new AppError("save order", { cause: err });
  }
  ```
- `{ cause: err }` for stack trace chaining (ES2022+)
- `Result<T>` discriminated union for expected failures:
  ```typescript
  type Result<T> = { ok: true; value: T } | { ok: false; error: AppError };

  function parseConfig(raw: string): Result<Config> {
    try {
      return { ok: true, value: configSchema.parse(JSON.parse(raw)) };
    } catch (err) {
      return { ok: false, error: new AppError("parse config", { cause: err }) };
    }
  }
  ```
- Server: Next.js `error.tsx` boundaries, server actions return `{ success, error }` to client:
  ```typescript
  "use server";
  export async function createOrder(data: FormData): Promise<{ success: boolean; error?: string }> {
    try {
      await orderService.create(parseOrderForm(data));
      return { success: true };
    } catch (err: unknown) {
      if (err instanceof ValidationError) return { success: false, error: err.message };
      throw err; // let error.tsx handle unexpected errors
    }
  }
  ```
- Client: check `response.ok` after fetch, `p-retry` for exponential backoff:
  ```typescript
  import pRetry from "p-retry";
  const data = await pRetry(() => fetchOrder(orderId), { retries: 3 });
  ```

## Dependencies

- `npm install --save-exact` or `"exact": true` in `.npmrc`
- `npm audit` / `pnpm audit` for vulnerabilities
- Prefer packages with built-in TS types over `@types/*`
- `strict: true`, `noUncheckedIndexedAccess: true` in `tsconfig.json`:
  ```json
  {
    "compilerOptions": {
      "strict": true,
      "noUncheckedIndexedAccess": true,
      "exactOptionalPropertyTypes": true
    }
  }
  ```

## Code Readability

- Specific names: `userRepository` not `repo`, `OrderStatus` not `Status`
- Guard clauses with early `return`/`throw`:
  ```typescript
  function getDiscount(user: User): number {
    if (!user.active) return 0;
    if (!user.subscription) return 0;
    return calculatePlanDiscount(user.subscription);
  }
  ```
- `satisfies` for inline type validation: `const config = { port: 3000 } satisfies Config`
- Options object when >5 params:
  ```typescript
  // Wrong
  function createUser(name: string, email: string, role: Role, team: string, active: boolean, verified: boolean): User

  // Right
  interface CreateUserParams { name: string; email: string; role: Role; team: string; active: boolean; verified: boolean }
  function createUser(params: CreateUserParams): User
  ```
- `fs.readFile()` for binary data, not embedded base64
- Zod at API boundaries; trust internal types elsewhere:
  ```typescript
  // API boundary — validate
  const body = createOrderSchema.parse(await request.json());
  // Internal call — trust types
  await orderService.create(body);
  ```
- Minimal comments, delete dead code

## Code Reuse

- Types defined in provider module — consumers import, never redefine
- `as const` objects for shared constants:
  ```typescript
  export const HTTP_STATUS = { OK: 200, NOT_FOUND: 404, INTERNAL: 500 } as const;
  type HttpStatus = (typeof HTTP_STATUS)[keyof typeof HTTP_STATUS]; // 200 | 404 | 500
  ```
- React: shared hooks in `hooks/`, shared UI in `components/`:
  ```typescript
  // hooks/use-debounce.ts
  export function useDebounce<T>(value: T, delayMs: number): T { /* ... */ }

  // components/data-table.tsx
  export function DataTable<T>({ columns, data }: DataTableProps<T>): ReactNode { /* ... */ }
  ```
- Search before adding; refactor before and after changes

## Code Organization

- Entry point (`main.ts`, `server.ts`, Next.js `instrumentation.ts`) only initializes — wire dependencies and start:
  ```typescript
  // instrumentation.ts
  import { NodeSDK } from "@opentelemetry/sdk-node";
  const sdk = new NodeSDK({ /* ... */ });
  sdk.start();
  ```
- Boundary layers: `app/api/` (HTTP), `lib/db/` (persistence), `lib/vendors/` (3rd party)
- Domain by feature: `lib/user/`, `lib/billing/` — not `lib/services/`, `lib/repositories/`:
  ```
  lib/
    user/
      user-service.ts
      user-repository.ts
      user.ts          # types
    billing/
      billing-service.ts
      invoice-repository.ts
      billing.ts       # types
  ```
- React/Next.js: `components/` shared UI, `app/` pages/routes, `lib/` business logic
- `package.json` `exports` or `internal/` for visibility control
- No circular dependencies — enforce with `eslint-plugin-import` `no-cycle` rule

## Type Safety

- Discriminated unions for state machines:
  ```typescript
  type ConnectionState =
    | { status: "connecting" }
    | { status: "connected"; socket: WebSocket }
    | { status: "error"; error: Error };

  function renderStatus(state: ConnectionState): string {
    switch (state.status) {
      case "connecting": return "Connecting...";
      case "connected": return `Connected to ${state.socket.url}`;
      case "error": return `Error: ${state.error.message}`;
    }
    // exhaustive — no default needed; TS errors if a case is missed
  }
  ```
- Branded types for IDs:
  ```typescript
  type UserId = string & { readonly __brand: "UserId" };
  type OrderId = string & { readonly __brand: "OrderId" };
  // Prevents: findUser(orderId) — type error
  ```
- `satisfies` to narrow without widening
- `readonly` arrays and `Readonly<T>`:
  ```typescript
  interface Config {
    readonly host: string;
    readonly port: number;
    readonly allowedOrigins: readonly string[];
  }
  ```
- `unknown` over `any` — narrow with type guards:
  ```typescript
  function isApiError(err: unknown): err is { code: number; message: string } {
    return typeof err === "object" && err !== null && "code" in err && "message" in err;
  }
  ```
- `as const` objects or enums over raw strings

## Initialization

- Required positional, optional in options object:
  ```typescript
  interface ClientOptions { timeout?: number; retries?: number }
  function createClient(baseUrl: string, options?: ClientOptions): Client {
    const { timeout = 5000, retries = 3 } = options ?? {};
    return new Client(baseUrl, timeout, retries);
  }
  ```
- Private constructor + `create()` factory:
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
- Don't use options pattern when all params are required

## External I/O

- Prisma `createMany()`, `findMany({ where: { id: { in: ids } } })` for batches — not N+1 loops:
  ```typescript
  // Wrong — N+1 queries
  const users = await Promise.all(ids.map(id => prisma.user.findUnique({ where: { id } })));

  // Right — single query
  const users = await prisma.user.findMany({ where: { id: { in: ids } } });
  ```
- `prisma.$transaction()` interactive transactions for multi-table writes:
  ```typescript
  await prisma.$transaction(async (tx) => {
    const order = await tx.order.create({ data: orderData });
    await tx.orderItem.createMany({ data: items.map(i => ({ ...i, orderId: order.id })) });
    await tx.inventory.updateMany({ where: { productId: { in: productIds } }, data: { reserved: true } });
  });
  ```
- `new Date()` from app code, not `NOW()` — inject `Clock` for testing:
  ```typescript
  interface Clock { now(): Date }
  const realClock: Clock = { now: () => new Date() };
  ```
- `Prisma.sql` tagged template for raw SQL — never string concatenation:
  ```typescript
  const users = await prisma.$queryRaw(Prisma.sql`SELECT * FROM "User" WHERE email = ${email}`);
  ```
- `AbortSignal.timeout()` for fetch timeouts, `AbortSignal.any()` to compose:
  ```typescript
  const response = await fetch(url, { signal: AbortSignal.timeout(5000) });
  // Compose with user cancellation
  const signal = AbortSignal.any([AbortSignal.timeout(5000), userController.signal]);
  ```

## Logging

- `pino` structured JSON:
  ```typescript
  import pino from "pino";
  const logger = pino({ level: config.LOG_LEVEL });
  logger.info({ userId, orderId, itemCount: order.items.length }, "order placed");
  ```
- Levels: `debug` (dev detail), `info` (operational events), `warn` (recoverable), `error` (needs investigation)
- Trace context:
  ```typescript
  const childLogger = logger.child({ traceId: span.spanContext().traceId, spanId: span.spanContext().spanId });
  ```
- Never log tokens, passwords, `process.env`, full request/response bodies

## Testing

- Vitest `test.each` for table-driven tests:
  ```typescript
  test.each([
    { input: "", expected: false },
    { input: "valid@email.com", expected: true },
    { input: "no-at-sign", expected: false },
  ])("validates $input -> $expected", ({ input, expected }) => {
    expect(isValidEmail(input)).toBe(expected);
  });
  ```
- DI via constructor — mock only 3rd party APIs:
  ```typescript
  const mockEmailSender: EmailSender = { send: vi.fn().mockResolvedValue(undefined) };
  const service = new UserService(mockEmailSender);
  ```
- testcontainers for real Postgres:
  ```typescript
  import { PostgreSqlContainer } from "@testcontainers/postgresql";
  let container: StartedPostgreSqlContainer;
  beforeAll(async () => { container = await new PostgreSqlContainer().start(); });
  afterAll(async () => { await container.stop(); });
  ```
- React: `@testing-library/react` — `render()`, `screen.getByRole()`, `userEvent` — test behavior, not implementation:
  ```typescript
  import { render, screen } from "@testing-library/react";
  import userEvent from "@testing-library/user-event";

  it("submits form with valid data", async () => {
    const onSubmit = vi.fn();
    render(<OrderForm onSubmit={onSubmit} />);
    await userEvent.type(screen.getByRole("textbox", { name: /email/i }), "user@example.com");
    await userEvent.click(screen.getByRole("button", { name: /submit/i }));
    expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({ email: "user@example.com" }));
  });
  ```
- `Clock` interface for time — never global `Date.now()` mock:
  ```typescript
  interface Clock { now(): Date }
  const fakeClock: Clock = { now: () => new Date("2024-01-15T10:00:00Z") };
  ```
- `beforeEach` reset DB, `afterEach` `vi.restoreAllMocks()`

## Concurrency

- JS is single-threaded — no locks unless Worker threads
- `Promise.all()` when batching impossible:
  ```typescript
  const [user, orders, notifications] = await Promise.all([
    userService.findById(userId),
    orderService.findByUser(userId),
    notificationService.findUnread(userId),
  ]);
  ```
- `p-limit` for controlled concurrency:
  ```typescript
  import pLimit from "p-limit";
  const limit = pLimit(5);
  const results = await Promise.all(urls.map(u => limit(() => fetch(u))));
  ```
- `Promise.allSettled()` for partial failures:
  ```typescript
  const results = await Promise.allSettled(users.map(u => sendWelcomeEmail(u)));
  const failures = results.filter((r): r is PromiseRejectedResult => r.status === "rejected");
  if (failures.length > 0) logger.warn({ count: failures.length }, "some welcome emails failed");
  ```
- BullMQ + Redis for fire-and-forget (emails, reports):
  ```typescript
  import { Queue } from "bullmq";
  const emailQueue = new Queue("email", { connection: redis });
  await emailQueue.add("welcome", { userId, email });
  ```

## Security

- `process.env` once at startup, validate with Zod:
  ```typescript
  const env = envSchema.parse(process.env); // throws on missing/invalid vars
  ```
- Zod schemas at API boundaries:
  ```typescript
  const createUserSchema = z.object({
    name: z.string().min(1).max(100),
    email: z.string().email(),
    role: z.enum(["admin", "user"]),
  });
  export async function POST(request: Request) {
    const body = createUserSchema.parse(await request.json());
    // body is fully validated and typed
  }
  ```
- `DOMPurify` for `dangerouslySetInnerHTML`:
  ```typescript
  import DOMPurify from "dompurify";
  <div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(htmlContent) }} />
  ```

## Monitoring

- `@opentelemetry/sdk-node` at bootstrap (`instrumentation.ts`):
  ```typescript
  import { NodeSDK } from "@opentelemetry/sdk-node";
  import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
  const sdk = new NodeSDK({ instrumentations: [getNodeAutoInstrumentations()] });
  sdk.start();
  ```
- RED metrics: counter for requests, histogram for duration, counter for errors
- No high-cardinality labels (no user IDs/request IDs as metric labels)
- `tracer.startActiveSpan()` for custom spans:
  ```typescript
  import { trace, SpanStatusCode } from "@opentelemetry/api";
  const tracer = trace.getTracer("order-service");

  async function processOrder(order: Order): Promise<void> {
    await tracer.startActiveSpan("processOrder", async (span) => {
      try {
        span.setAttribute("order.id", order.id);
        await validateOrder(order);
        await chargePayment(order);
        await fulfillOrder(order);
      } catch (err) {
        span.recordException(err as Error);
        span.setStatus({ code: SpanStatusCode.ERROR, message: (err as Error).message });
        throw err;
      } finally {
        span.end();
      }
    });
  }
  ```

## Frontend Guidelines

- **Mobile-first** — start from mobile layout, add responsive breakpoints with Tailwind `sm:`, `md:`, `lg:`:
  ```tsx
  <div className="flex flex-col gap-4 md:flex-row md:gap-8 lg:max-w-6xl">
    <aside className="w-full md:w-64 lg:w-80">{/* sidebar */}</aside>
    <main className="flex-1">{/* content */}</main>
  </div>
  ```
- **Import assets from files** — never inline SVG/base64:
  ```tsx
  // Wrong
  const Icon = () => <svg viewBox="0 0 24 24"><path d="M12..."/></svg>;

  // Right
  import CheckIcon from "./check-icon.svg";
  <Image src={CheckIcon} alt="check" width={24} height={24} />
  ```
- **Minimize `useEffect`** — consolidate related effects; prefer derived state and event handlers over effects:
  ```tsx
  // Wrong — unnecessary effect
  const [fullName, setFullName] = useState("");
  useEffect(() => { setFullName(`${firstName} ${lastName}`); }, [firstName, lastName]);

  // Right — derived state
  const fullName = `${firstName} ${lastName}`;
  ```
- **`useMemo` only with justification** — expensive computation or preventing re-renders via `React.memo`:
  ```tsx
  // Justified — expensive filtering and sorting
  const sortedItems = useMemo(
    () => items.filter(i => i.active).sort((a, b) => a.name.localeCompare(b.name)),
    [items],
  );

  // Not justified — cheap operation
  const label = `${user.name} (${user.role})`; // no useMemo needed
  ```
- **Verify UI visually** — screenshot at desktop (1280px) and mobile (375px) breakpoints after changes
