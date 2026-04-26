# Next.js Pages & Layouts (App Router)

## Page with Server Component Data Fetching

```typescript
import { tracer } from "@/lib/telemetry";

interface PageProps {
  params: Promise<{ id: string }>;
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { id } = await params;
  const product = await productService.getById(id);
  return {
    title: product?.name ?? "Product Not Found",
  };
}

export default async function ProductPage({ params }: PageProps) {
  return tracer.startActiveSpan("ProductPage", async (span) => {
    try {
      const { id } = await params;
      const product = await productService.getById(id);
      if (!product) notFound();

      return (
        <div>
          <h1>{product.name}</h1>
          <p>{product.description}</p>
          <ProductActions productId={product.id} />
        </div>
      );
    } finally {
      span.end();
    }
  });
}
```

## Layout with Shared Providers

```typescript
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Header />
        <main className="mx-auto max-w-7xl px-4 py-8">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
```

## Loading State

```typescript
export default function Loading() {
  return (
    <div className="space-y-4">
      <div className="h-8 w-64 animate-pulse rounded bg-gray-200" />
      <div className="h-4 w-full animate-pulse rounded bg-gray-200" />
      <div className="h-4 w-3/4 animate-pulse rounded bg-gray-200" />
    </div>
  );
}
```

## Error Boundary

```typescript
"use client";

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function ErrorBoundary({ error, reset }: ErrorProps) {
  useEffect(() => {
    logger.error("page error", { error: error.message, digest: error.digest });
  }, [error]);

  return (
    <div role="alert" className="rounded-lg border border-red-200 p-6">
      <h2 className="text-lg font-semibold text-red-800">Something went wrong</h2>
      <p className="mt-2 text-gray-600">Please try again or contact support.</p>
      <button
        type="button"
        onClick={reset}
        className="mt-4 rounded bg-red-600 px-4 py-2 text-white"
      >
        Try again
      </button>
    </div>
  );
}
```

## Server Action

```typescript
"use server";

import { tracer } from "@/lib/telemetry";
import { revalidatePath } from "next/cache";

export async function updateProduct(formData: FormData) {
  return tracer.startActiveSpan("updateProduct", async (span) => {
    try {
      const id = formData.get("id") as string;
      const name = formData.get("name") as string;
      const price = Number(formData.get("price"));

      const input = updateProductSchema.parse({ id, name, price });
      await productService.update(input.id, { name: input.name, price: input.price });

      revalidatePath(`/products/${id}`);
    } catch (err) {
      span.recordException(err instanceof Error ? err : new Error(String(err)));
      throw err;
    } finally {
      span.end();
    }
  });
}
```

## Client Component with Interactivity

```typescript
"use client";

interface ProductActionsProps {
  productId: string;
}

function ProductActions({ productId }: ProductActionsProps) {
  const [isPending, startTransition] = useTransition();

  function handleDelete() {
    startTransition(async () => {
      await deleteProduct(productId);
    });
  }

  return (
    <div className="flex gap-2">
      <Link href={`/products/${productId}/edit`} className="rounded border px-3 py-1">
        Edit
      </Link>
      <button
        type="button"
        onClick={handleDelete}
        disabled={isPending}
        className="rounded bg-red-600 px-3 py-1 text-white disabled:opacity-50"
      >
        {isPending ? "Deleting..." : "Delete"}
      </button>
    </div>
  );
}
```

## List Page with Pagination

```typescript
interface PageProps {
  searchParams: Promise<{ page?: string; q?: string }>;
}

export default async function ProductsPage({ searchParams }: PageProps) {
  return tracer.startActiveSpan("ProductsPage", async (span) => {
    try {
      const { page: pageParam, q } = await searchParams;
      const page = Math.max(1, Number(pageParam) || 1);
      const { products, totalPages } = await productService.list({ page, query: q });

      return (
        <div>
          <h1>Products</h1>
          <SearchBar defaultValue={q} />
          <ProductGrid products={products} />
          <Pagination currentPage={page} totalPages={totalPages} />
        </div>
      );
    } finally {
      span.end();
    }
  });
}
```
