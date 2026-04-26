# HTTP Client Patterns

Idiomatic TypeScript HTTP client using `fetch` with inline OpenTelemetry instrumentation.

## Client class with base URL, timeout, and AbortSignal composition

```typescript
import { trace, SpanStatusCode } from "@opentelemetry/api";
import { z } from "zod";

const tracer = trace.getTracer("vendor-client");

class VendorClient {
  constructor(
    private readonly baseURL: string,
    private readonly timeoutMs: number = 10_000,
  ) {}

  private createAbortSignal(callerSignal?: AbortSignal): AbortSignal {
    const timeout = AbortSignal.timeout(this.timeoutMs);
    return callerSignal ? AbortSignal.any([callerSignal, timeout]) : timeout;
  }
}
```

## GET request with typed response and OTel span

```typescript
const UserSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string().email(),
});
type User = z.infer<typeof UserSchema>;

class VendorClient {
  // ... constructor and createAbortSignal from above

  async getUser(
    userId: string,
    options?: { signal?: AbortSignal },
  ): Promise<User> {
    return tracer.startActiveSpan("VendorClient.getUser", async (span) => {
      span.setAttribute("user.id", userId);
      try {
        const response = await fetch(`${this.baseURL}/users/${userId}`, {
          method: "GET",
          headers: { Accept: "application/json" },
          signal: this.createAbortSignal(options?.signal),
        });

        if (!response.ok) {
          throw new Error(
            `GET /users/${userId} failed with status ${response.status}`,
          );
        }

        const body = await response.json();
        const user = UserSchema.parse(body);

        span.setStatus({ code: SpanStatusCode.OK });
        return user;
      } catch (error) {
        span.setStatus({
          code: SpanStatusCode.ERROR,
          message: error instanceof Error ? error.message : String(error),
        });
        span.recordException(
          error instanceof Error ? error : new Error(String(error)),
        );
        throw new Error(`get user ${userId}: ${error}`, { cause: error });
      } finally {
        span.end();
      }
    });
  }
}
```

## POST request with request body and OTel span

```typescript
const CreateOrderResponseSchema = z.object({
  orderId: z.string(),
  status: z.enum(["pending", "confirmed"]),
  createdAt: z.string().datetime(),
});
type CreateOrderResponse = z.infer<typeof CreateOrderResponseSchema>;

interface CreateOrderRequest {
  userId: string;
  items: Array<{ productId: string; quantity: number }>;
}

class VendorClient {
  // ... constructor and createAbortSignal from above

  async createOrder(
    request: CreateOrderRequest,
    options?: { signal?: AbortSignal },
  ): Promise<CreateOrderResponse> {
    return tracer.startActiveSpan("VendorClient.createOrder", async (span) => {
      span.setAttribute("user.id", request.userId);
      span.setAttribute("order.item_count", request.items.length);
      try {
        const response = await fetch(`${this.baseURL}/orders`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
          },
          body: JSON.stringify(request),
          signal: this.createAbortSignal(options?.signal),
        });

        if (!response.ok) {
          throw new Error(
            `POST /orders failed with status ${response.status}`,
          );
        }

        const body = await response.json();
        const order = CreateOrderResponseSchema.parse(body);

        span.setAttribute("order.id", order.orderId);
        span.setStatus({ code: SpanStatusCode.OK });
        return order;
      } catch (error) {
        span.setStatus({
          code: SpanStatusCode.ERROR,
          message: error instanceof Error ? error.message : String(error),
        });
        span.recordException(
          error instanceof Error ? error : new Error(String(error)),
        );
        throw new Error(
          `create order for user ${request.userId}: ${error}`,
          { cause: error },
        );
      } finally {
        span.end();
      }
    });
  }
}
```

## Retry with exponential backoff

Use `p-retry` for transient failures (429, 5xx, network errors). Wrap the
retryable call inside the OTel span so the span captures the final outcome.

```typescript
import pRetry, { AbortError } from "p-retry";

class VendorClient {
  // ... constructor and createAbortSignal from above

  async getUserWithRetry(
    userId: string,
    options?: { signal?: AbortSignal; maxRetries?: number },
  ): Promise<User> {
    return tracer.startActiveSpan(
      "VendorClient.getUserWithRetry",
      async (span) => {
        span.setAttribute("user.id", userId);
        try {
          const user = await pRetry(
            async () => {
              const response = await fetch(
                `${this.baseURL}/users/${userId}`,
                {
                  method: "GET",
                  headers: { Accept: "application/json" },
                  signal: this.createAbortSignal(options?.signal),
                },
              );

              if (response.status === 429 || response.status >= 500) {
                throw new Error(
                  `GET /users/${userId} failed with status ${response.status}`,
                );
              }

              if (!response.ok) {
                // 4xx errors other than 429 are not transient; abort retry
                throw new AbortError(
                  `GET /users/${userId} failed with status ${response.status}`,
                );
              }

              const body = await response.json();
              return UserSchema.parse(body);
            },
            {
              retries: options?.maxRetries ?? 3,
              onFailedAttempt: (error) => {
                span.addEvent("retry_attempt", {
                  attempt: error.attemptNumber,
                  retriesLeft: error.retriesLeft,
                });
              },
            },
          );

          span.setStatus({ code: SpanStatusCode.OK });
          return user;
        } catch (error) {
          span.setStatus({
            code: SpanStatusCode.ERROR,
            message: error instanceof Error ? error.message : String(error),
          });
          span.recordException(
            error instanceof Error ? error : new Error(String(error)),
          );
          throw new Error(`get user with retry ${userId}: ${error}`, {
            cause: error,
          });
        } finally {
          span.end();
        }
      },
    );
  }
}
```

## Caller-controlled cancellation

The caller creates an `AbortController` and passes its signal. The client
composes it with the timeout signal via `AbortSignal.any()`.

```typescript
async function fetchUserProfile(client: VendorClient): Promise<User> {
  const controller = new AbortController();

  // Cancel the request after an external event, e.g. user navigation
  window.addEventListener("beforeunload", () => controller.abort(), {
    once: true,
  });

  // The client composes this signal with its own timeout
  return client.getUser("user-123", { signal: controller.signal });
}
```

Programmatic cancellation from a parent operation:

```typescript
async function syncUsers(
  client: VendorClient,
  userIds: string[],
): Promise<User[]> {
  const controller = new AbortController();

  try {
    const users: User[] = [];
    for (const id of userIds) {
      const user = await client.getUser(id, { signal: controller.signal });
      users.push(user);
    }
    return users;
  } catch (error) {
    // Cancel any remaining in-flight requests
    controller.abort();
    throw new Error(`sync users: ${error}`, { cause: error });
  }
}
```
