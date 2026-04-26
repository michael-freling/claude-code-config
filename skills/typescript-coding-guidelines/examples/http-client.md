# HTTP Client Patterns

## Client with Options and Tracing

```typescript
import { tracer } from "@/lib/telemetry";

interface ClientOptions {
  timeout?: number;
  maxRetries?: number;
}

class ApiClient {
  private readonly baseUrl: string;
  private readonly timeout: number;
  private readonly maxRetries: number;

  constructor(baseUrl: string, options?: ClientOptions) {
    this.baseUrl = baseUrl;
    this.timeout = options?.timeout ?? 30_000;
    this.maxRetries = options?.maxRetries ?? 3;
  }

  async get<T>(path: string, signal?: AbortSignal): Promise<T> {
    return this.request("GET", path, undefined, signal);
  }

  async post<T>(path: string, body: unknown, signal?: AbortSignal): Promise<T> {
    return this.request("POST", path, body, signal);
  }

  private async request<T>(
    method: string,
    path: string,
    body: unknown,
    signal?: AbortSignal,
  ): Promise<T> {
    return tracer.startActiveSpan(`${method} ${path}`, async (span) => {
      try {
        const response = await this.fetchWithRetry(
          `${this.baseUrl}${path}`,
          {
            method,
            headers: { "Content-Type": "application/json" },
            body: body ? JSON.stringify(body) : undefined,
            signal: signal ?? AbortSignal.timeout(this.timeout),
          },
        );

        if (!response.ok) {
          const text = await response.text();
          throw new ApiError(
            `${method} ${path}: status ${response.status}`,
            response.status,
            text,
          );
        }

        return (await response.json()) as T;
      } catch (err) {
        span.recordException(err instanceof Error ? err : new Error(String(err)));
        throw err;
      } finally {
        span.end();
      }
    });
  }
}
```

## Custom Error Class

```typescript
class ApiError extends Error {
  constructor(
    message: string,
    readonly statusCode: number,
    readonly responseBody: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}
```

## Retry with Exponential Backoff

```typescript
private async fetchWithRetry(url: string, init: RequestInit): Promise<Response> {
  let lastError: Error | undefined;

  for (let attempt = 0; attempt < this.maxRetries; attempt++) {
    try {
      const response = await fetch(url, init);
      if (!isRetryableStatus(response.status)) {
        return response;
      }
      lastError = new Error(`status ${response.status}`);
    } catch (err) {
      if (err instanceof DOMException && err.name === "AbortError") throw err;
      lastError = err instanceof Error ? err : new Error(String(err));
    }

    await sleep(backoffDelay(attempt));
  }

  throw new Error(`max retries exceeded for ${url}`, { cause: lastError });
}

function isRetryableStatus(status: number): boolean {
  return status === 429 || status === 502 || status === 503 || status === 504;
}

function backoffDelay(attempt: number): number {
  return Math.min(2 ** attempt * 100, 10_000);
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
```

## Concurrent Fan-Out Requests

```typescript
async function getProducts(ids: string[]): Promise<Product[]> {
  return tracer.startActiveSpan("getProducts", async (span) => {
    try {
      span.setAttribute("product.count", ids.length);
      return await Promise.all(ids.map((id) => client.get<Product>(`/products/${id}`)));
    } finally {
      span.end();
    }
  });
}
```

## Typed Response Validation with Zod

```typescript
const productResponseSchema = z.object({
  id: z.string(),
  name: z.string(),
  price: z.number(),
  inStock: z.boolean(),
});

type ProductResponse = z.infer<typeof productResponseSchema>;

async function getProduct(id: string): Promise<ProductResponse> {
  const data = await client.get(`/products/${id}`);
  return productResponseSchema.parse(data);
}
```

## Testing with msw

```typescript
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("ApiClient", () => {
  const client = new ApiClient("http://api.test");

  it("returns parsed product", async () => {
    server.use(
      http.get("http://api.test/products/p-1", () =>
        HttpResponse.json({ id: "p-1", name: "Widget", price: 999, inStock: true }),
      ),
    );

    const product = await client.get<Product>("/products/p-1");
    expect(product.name).toBe("Widget");
  });

  it("throws ApiError on 404", async () => {
    server.use(
      http.get("http://api.test/products/missing", () =>
        HttpResponse.json({ error: "not found" }, { status: 404 }),
      ),
    );

    await expect(client.get("/products/missing")).rejects.toThrow(ApiError);
  });

  it("retries on 503", async () => {
    let callCount = 0;
    server.use(
      http.get("http://api.test/products/flaky", () => {
        callCount++;
        if (callCount < 3) {
          return HttpResponse.json(null, { status: 503 });
        }
        return HttpResponse.json({ id: "flaky", name: "Flaky", price: 100, inStock: true });
      }),
    );

    const product = await client.get<Product>("/products/flaky");
    expect(product.name).toBe("Flaky");
    expect(callCount).toBe(3);
  });
});
```
