# Next.js App Router Patterns

## Server Component with Data Fetching

```typescript
// app/users/[id]/page.tsx

interface PageProps {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ tab?: string }>;
}

async function UserPage({ params, searchParams }: PageProps) {
  const { id } = await params;
  const { tab } = await searchParams;

  const user = await db.user.findUniqueOrThrow({
    where: { id },
    include: { posts: tab === "posts" },
  });

  return (
    <main>
      <h1>{user.name}</h1>
      {tab === "posts" && <PostList posts={user.posts} />}
    </main>
  );
}

export default UserPage;
```

```typescript
// app/users/[id]/loading.tsx

export default function Loading() {
  return <UserSkeleton />;
}
```

```typescript
// app/users/[id]/error.tsx
"use client";

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function Error({ error, reset }: ErrorProps) {
  return (
    <div>
      <h2>Something went wrong</h2>
      <p>{error.message}</p>
      <button onClick={reset}>Try again</button>
    </div>
  );
}
```

## Server Action with Form Handling

```typescript
// app/users/actions.ts
"use server";

import { z } from "zod";
import { trace } from "@opentelemetry/api";

const CreateUserSchema = z.object({
  name: z.string().min(1).max(100),
  email: z.string().email(),
  role: z.enum(["admin", "member"]),
});

type CreateUserResult =
  | { success: true; data: User }
  | { success: false; error: string };

export async function createUser(
  _prevState: CreateUserResult | null,
  formData: FormData,
): Promise<CreateUserResult> {
  const span = trace.getTracer("actions").startSpan("createUser");

  try {
    const parsed = CreateUserSchema.safeParse({
      name: formData.get("name"),
      email: formData.get("email"),
      role: formData.get("role"),
    });

    if (!parsed.success) {
      return {
        success: false,
        error: parsed.error.issues.map((i) => i.message).join(", "),
      };
    }

    const user = await db.user.create({ data: parsed.data });
    revalidatePath("/users");

    return { success: true, data: user };
  } catch (err) {
    span.recordException(err as Error);
    return { success: false, error: "Failed to create user" };
  } finally {
    span.end();
  }
}
```

```typescript
// app/users/create-user-form.tsx
"use client";

import { useActionState } from "react";
import { createUser } from "./actions";

export function CreateUserForm() {
  const [state, formAction, isPending] = useActionState(createUser, null);

  return (
    <form action={formAction}>
      <input name="name" required />
      <input name="email" type="email" required />
      <select name="role">
        <option value="member">Member</option>
        <option value="admin">Admin</option>
      </select>

      <button type="submit" disabled={isPending}>
        {isPending ? "Creating..." : "Create User"}
      </button>

      {state && !state.success && (
        <p role="alert">{state.error}</p>
      )}
    </form>
  );
}
```

## Client Component with Interactivity

```typescript
// app/users/user-card.tsx
"use client";

import { useState, useTransition } from "react";
import { deleteUser } from "./actions";

interface UserCardProps {
  user: User;
}

export function UserCard({ user }: UserCardProps) {
  const [showConfirm, setShowConfirm] = useState(false);
  const [isPending, startTransition] = useTransition();

  function handleDelete() {
    startTransition(async () => {
      await deleteUser(user.id);
      setShowConfirm(false);
    });
  }

  return (
    <div>
      <h3>{user.name}</h3>
      <p>{user.email}</p>

      {showConfirm ? (
        <div>
          <p>Delete this user?</p>
          <button onClick={handleDelete} disabled={isPending}>
            {isPending ? "Deleting..." : "Confirm"}
          </button>
          <button onClick={() => setShowConfirm(false)}>Cancel</button>
        </div>
      ) : (
        <button onClick={() => setShowConfirm(true)}>Delete</button>
      )}
    </div>
  );
}
```

## Layout with Shared Data

```typescript
// app/layout.tsx

import type { Metadata } from "next";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/components/auth-provider";

export const metadata: Metadata = {
  title: {
    template: "%s | MyApp",
    default: "MyApp",
  },
  description: "Application description",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <AuthProvider>
          <ThemeProvider>
            {children}
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
```

```typescript
// app/dashboard/layout.tsx

import { redirect } from "next/navigation";
import { getSession } from "@/lib/auth";
import { DashboardNav } from "./dashboard-nav";

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await getSession();
  if (!session) {
    redirect("/login");
  }

  return (
    <div>
      <DashboardNav user={session.user} />
      <main>{children}</main>
    </div>
  );
}
```

```typescript
// app/dashboard/settings/page.tsx

import type { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  const session = await getSession();

  return {
    title: `Settings - ${session?.user.name}`,
  };
}

export default async function SettingsPage() {
  const session = await getSession();

  return <SettingsForm user={session!.user} />;
}
```

## Middleware

```typescript
// middleware.ts

import { NextRequest, NextResponse } from "next/server";
import { verifyToken } from "@/lib/auth";

const PUBLIC_PATHS = ["/login", "/signup", "/api/health"];

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (PUBLIC_PATHS.some((path) => pathname.startsWith(path))) {
    return NextResponse.next();
  }

  const token = request.cookies.get("session")?.value;
  if (!token) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("callbackUrl", pathname);
    return NextResponse.redirect(loginUrl);
  }

  const payload = await verifyToken(token);
  if (!payload) {
    const loginUrl = new URL("/login", request.url);
    return NextResponse.redirect(loginUrl);
  }

  // Pass user info to downstream server components via headers
  const headers = new Headers(request.headers);
  headers.set("x-user-id", payload.userId);
  headers.set("x-user-role", payload.role);

  return NextResponse.next({ request: { headers } });
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
```

## Error Boundary

```typescript
// app/error.tsx
"use client";

import { useEffect } from "react";

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function GlobalError({ error, reset }: ErrorProps) {
  useEffect(() => {
    console.error("Unhandled error:", error);
  }, [error]);

  return (
    <div>
      <h2>Something went wrong</h2>
      <p>An unexpected error occurred. Please try again.</p>
      <button onClick={reset}>Try again</button>
    </div>
  );
}
```

```typescript
// app/users/[id]/not-found.tsx

import Link from "next/link";

export default function UserNotFound() {
  return (
    <div>
      <h2>User not found</h2>
      <p>The user you are looking for does not exist.</p>
      <Link href="/users">Back to users</Link>
    </div>
  );
}
```

Trigger `notFound()` from the server component:

```typescript
// app/users/[id]/page.tsx

import { notFound } from "next/navigation";

async function UserPage({ params }: PageProps) {
  const { id } = await params;

  const user = await db.user.findUnique({ where: { id } });
  if (!user) {
    notFound();
  }

  return <UserProfile user={user} />;
}
```
