# API Route Patterns (Next.js Route Handlers)

## Route Handler with Tracing and Validation

```typescript
import { tracer } from "@/lib/telemetry";
import { NextRequest, NextResponse } from "next/server";

const createUserSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
});

export async function POST(request: NextRequest) {
  return tracer.startActiveSpan("POST /api/users", async (span) => {
    try {
      const body = await request.json();
      const input = createUserSchema.safeParse(body);
      if (!input.success) {
        return NextResponse.json(
          { error: { message: input.error.issues[0].message, code: "VALIDATION_ERROR" } },
          { status: 400 },
        );
      }

      const user = await userService.create(input.data);

      return NextResponse.json(user, { status: 201 });
    } catch (err) {
      span.recordException(err instanceof Error ? err : new Error(String(err)));
      logger.error("create user", { error: err });
      return NextResponse.json(
        { error: { message: "Internal server error", code: "INTERNAL_ERROR" } },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## GET with Path Parameters

```typescript
interface RouteParams {
  params: Promise<{ id: string }>;
}

export async function GET(request: NextRequest, { params }: RouteParams) {
  return tracer.startActiveSpan("GET /api/users/:id", async (span) => {
    try {
      const { id } = await params;
      span.setAttribute("user.id", id);

      const user = await userService.getById(id);
      if (!user) {
        return NextResponse.json(
          { error: { message: "User not found", code: "NOT_FOUND" } },
          { status: 404 },
        );
      }

      return NextResponse.json(user);
    } catch (err) {
      span.recordException(err instanceof Error ? err : new Error(String(err)));
      logger.error("get user", { error: err });
      return NextResponse.json(
        { error: { message: "Internal server error", code: "INTERNAL_ERROR" } },
        { status: 500 },
      );
    } finally {
      span.end();
    }
  });
}
```

## Domain Error Mapping

```typescript
function mapErrorToResponse(err: unknown): NextResponse {
  if (err instanceof NotFoundError) {
    return NextResponse.json(
      { error: { message: err.message, code: "NOT_FOUND" } },
      { status: 404 },
    );
  }
  if (err instanceof ValidationError) {
    return NextResponse.json(
      { error: { message: err.message, code: "VALIDATION_ERROR" } },
      { status: 400 },
    );
  }
  if (err instanceof ConflictError) {
    return NextResponse.json(
      { error: { message: err.message, code: "CONFLICT" } },
      { status: 409 },
    );
  }
  logger.error("unhandled error", { error: err });
  return NextResponse.json(
    { error: { message: "Internal server error", code: "INTERNAL_ERROR" } },
    { status: 500 },
  );
}
```

## List with Query Parameters and Pagination

```typescript
const listUsersSchema = z.object({
  page: z.coerce.number().int().positive().default(1),
  limit: z.coerce.number().int().min(1).max(100).default(20),
  q: z.string().optional(),
});

export async function GET(request: NextRequest) {
  return tracer.startActiveSpan("GET /api/users", async (span) => {
    try {
      const searchParams = Object.fromEntries(request.nextUrl.searchParams);
      const query = listUsersSchema.safeParse(searchParams);
      if (!query.success) {
        return NextResponse.json(
          { error: { message: "Invalid query parameters", code: "VALIDATION_ERROR" } },
          { status: 400 },
        );
      }

      const { users, total } = await userService.list(query.data);

      return NextResponse.json({
        data: users,
        pagination: {
          page: query.data.page,
          limit: query.data.limit,
          total,
        },
      });
    } finally {
      span.end();
    }
  });
}
```

## Middleware for Authentication

```typescript
import { NextRequest, NextResponse } from "next/server";

export function middleware(request: NextRequest) {
  const token = request.headers.get("authorization")?.replace("Bearer ", "");
  if (!token) {
    return NextResponse.json(
      { error: { message: "Authentication required", code: "UNAUTHENTICATED" } },
      { status: 401 },
    );
  }

  try {
    const payload = verifyToken(token);
    const headers = new Headers(request.headers);
    headers.set("x-user-id", payload.userId);
    return NextResponse.next({ request: { headers } });
  } catch {
    return NextResponse.json(
      { error: { message: "Invalid token", code: "UNAUTHENTICATED" } },
      { status: 401 },
    );
  }
}

export const config = {
  matcher: "/api/:path*",
};
```

## Testing API Routes

```typescript
import { POST } from "./route";

describe("POST /api/users", () => {
  it("creates a user with valid input", async () => {
    const request = new NextRequest("http://localhost/api/users", {
      method: "POST",
      body: JSON.stringify({ name: "Alice", email: "alice@example.com" }),
    });

    const response = await POST(request);
    const body = await response.json();

    expect(response.status).toBe(201);
    expect(body.name).toBe("Alice");
  });

  it("returns 400 for invalid email", async () => {
    const request = new NextRequest("http://localhost/api/users", {
      method: "POST",
      body: JSON.stringify({ name: "Alice", email: "not-an-email" }),
    });

    const response = await POST(request);
    expect(response.status).toBe(400);
  });
});
```
