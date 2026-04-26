# Next.js API Route Handlers with OpenTelemetry

## GET handler with query parameters and OTel span

```typescript
import { trace, SpanStatusCode } from "@opentelemetry/api";
import { NextRequest, NextResponse } from "next/server";
import { z } from "zod";
import { UserService } from "@/services/user";

const tracer = trace.getTracer("user-api");

const listUsersQuery = z.object({
  page: z.coerce.number().int().positive().default(1),
  limit: z.coerce.number().int().min(1).max(100).default(20),
  role: z.enum(["admin", "member", "viewer"]).optional(),
});

export async function GET(request: NextRequest) {
  return tracer.startActiveSpan("GET /api/users", async (span) => {
    try {
      const params = Object.fromEntries(request.nextUrl.searchParams);
      const parsed = listUsersQuery.safeParse(params);
      if (!parsed.success) {
        span.setStatus({
          code: SpanStatusCode.ERROR,
          message: "validation error",
        });
        return NextResponse.json(
          { error: "Invalid query parameters", details: parsed.error.flatten() },
          { status: 400 },
        );
      }

      span.setAttributes({
        "query.page": parsed.data.page,
        "query.limit": parsed.data.limit,
      });

      const users = await UserService.list(parsed.data);
      return NextResponse.json(users);
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR, message: "internal error" });
      return NextResponse.json(
        { error: "Internal server error" },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## POST handler with body validation, OTel span, and structured logging

```typescript
import { trace, SpanStatusCode, context as otelContext } from "@opentelemetry/api";
import { NextRequest, NextResponse } from "next/server";
import { z } from "zod";
import pino from "pino";
import { UserService } from "@/services/user";

const tracer = trace.getTracer("user-api");
const logger = pino({ name: "user-api" });

const createUserBody = z.object({
  name: z.string().min(1).max(200),
  email: z.string().email(),
  role: z.enum(["admin", "member", "viewer"]),
});

export async function POST(request: NextRequest) {
  return tracer.startActiveSpan("POST /api/users", async (span) => {
    const traceId = span.spanContext().traceId;
    const spanId = span.spanContext().spanId;
    const log = logger.child({ traceId, spanId });

    try {
      const body = await request.json();
      const parsed = createUserBody.safeParse(body);
      if (!parsed.success) {
        log.warn({ errors: parsed.error.flatten() }, "validation failed");
        span.setStatus({ code: SpanStatusCode.ERROR, message: "validation error" });
        return NextResponse.json(
          { error: "Invalid request body", details: parsed.error.flatten() },
          { status: 400 },
        );
      }

      span.setAttributes({ "user.email": parsed.data.email, "user.role": parsed.data.role });

      const user = await UserService.create(parsed.data);
      log.info({ userId: user.id }, "user created");

      return NextResponse.json(user, { status: 201 });
    } catch (err) {
      log.error({ err }, "failed to create user");
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR, message: "internal error" });
      return NextResponse.json(
        { error: "Internal server error" },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## PATCH and DELETE handlers with path parameters

```typescript
// app/api/users/[id]/route.ts

import { trace, SpanStatusCode } from "@opentelemetry/api";
import { NextRequest, NextResponse } from "next/server";
import { z } from "zod";
import { UserService } from "@/services/user";

const tracer = trace.getTracer("user-api");

const updateUserBody = z.object({
  name: z.string().min(1).max(200).optional(),
  role: z.enum(["admin", "member", "viewer"]).optional(),
});

type RouteParams = { params: Promise<{ id: string }> };

export async function PATCH(request: NextRequest, { params }: RouteParams) {
  const { id } = await params;

  return tracer.startActiveSpan(`PATCH /api/users/${id}`, async (span) => {
    try {
      span.setAttribute("user.id", id);

      const existing = await UserService.findById(id);
      if (!existing) {
        span.setStatus({ code: SpanStatusCode.ERROR, message: "not found" });
        return NextResponse.json({ error: "User not found" }, { status: 404 });
      }

      const body = await request.json();
      const parsed = updateUserBody.safeParse(body);
      if (!parsed.success) {
        span.setStatus({ code: SpanStatusCode.ERROR, message: "validation error" });
        return NextResponse.json(
          { error: "Invalid request body", details: parsed.error.flatten() },
          { status: 400 },
        );
      }

      const updated = await UserService.update(id, parsed.data);
      return NextResponse.json(updated);
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR, message: "internal error" });
      return NextResponse.json(
        { error: "Internal server error" },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}

export async function DELETE(_request: NextRequest, { params }: RouteParams) {
  const { id } = await params;

  return tracer.startActiveSpan(`DELETE /api/users/${id}`, async (span) => {
    try {
      span.setAttribute("user.id", id);

      const existing = await UserService.findById(id);
      if (!existing) {
        span.setStatus({ code: SpanStatusCode.ERROR, message: "not found" });
        return NextResponse.json({ error: "User not found" }, { status: 404 });
      }

      await UserService.delete(id);
      return new NextResponse(null, { status: 204 });
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR, message: "internal error" });
      return NextResponse.json(
        { error: "Internal server error" },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## Middleware-style auth validation

```typescript
// lib/auth.ts

import { NextRequest, NextResponse } from "next/server";
import { verifyToken } from "@/services/auth";

type UserContext = {
  userId: string;
  email: string;
  role: "admin" | "member" | "viewer";
};

type AuthResult =
  | { ok: true; user: UserContext }
  | { ok: false; response: NextResponse };

export async function validateAuth(request: NextRequest): Promise<AuthResult> {
  const header = request.headers.get("authorization");
  if (!header?.startsWith("Bearer ")) {
    return {
      ok: false,
      response: NextResponse.json(
        { error: "Missing or malformed authorization header" },
        { status: 401 },
      ),
    };
  }

  const token = header.slice(7);
  const user = await verifyToken(token);
  if (!user) {
    return {
      ok: false,
      response: NextResponse.json(
        { error: "Invalid or expired token" },
        { status: 401 },
      ),
    };
  }

  return { ok: true, user };
}
```

Usage in a route handler:

```typescript
import { trace, SpanStatusCode } from "@opentelemetry/api";
import { NextRequest, NextResponse } from "next/server";
import { validateAuth } from "@/lib/auth";
import { OrderService } from "@/services/order";

const tracer = trace.getTracer("order-api");

export async function GET(request: NextRequest) {
  return tracer.startActiveSpan("GET /api/orders", async (span) => {
    try {
      const auth = await validateAuth(request);
      if (!auth.ok) {
        span.setStatus({ code: SpanStatusCode.ERROR, message: "auth failed" });
        return auth.response;
      }

      span.setAttribute("user.id", auth.user.userId);

      const orders = await OrderService.listByUser(auth.user.userId);
      return NextResponse.json(orders);
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR, message: "internal error" });
      return NextResponse.json(
        { error: "Internal server error" },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## Error response patterns

All error responses use a consistent shape:

```typescript
type ErrorResponse = {
  error: string;
  details?: unknown;
};
```

| Status | When to use                        | `error` field                          | `details` field              |
|--------|------------------------------------|----------------------------------------|------------------------------|
| 400    | Zod validation fails               | `"Invalid request body"` or `"Invalid query parameters"` | `parsed.error.flatten()` |
| 401    | Missing or invalid auth token      | `"Missing or malformed authorization header"` or `"Invalid or expired token"` | omitted |
| 404    | Resource not found by ID           | `"User not found"` (use the resource name) | omitted |
| 500    | Unhandled exception in catch block | `"Internal server error"`              | omitted (never expose internals) |

Key rules:
- 400 responses include `details` with structured Zod errors so clients can display field-level messages
- 401 responses distinguish between missing credentials and invalid credentials
- 500 responses use a fixed generic message and log the real error server-side via `span.recordException()` and structured logger
