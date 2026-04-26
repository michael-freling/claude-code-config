# Database (Prisma + OpenTelemetry)

## Connection Pool Configuration

Configure the connection pool via the datasource URL:

```
DATABASE_URL="postgresql://user:pass@host:5432/mydb?connection_limit=20&pool_timeout=10"
```

Singleton PrismaClient to avoid multiple instances during Next.js hot reloads:

```typescript
import { PrismaClient } from "@prisma/client";

const globalForPrisma = globalThis as unknown as { prisma: PrismaClient };

export const prisma =
  globalForPrisma.prisma ??
  new PrismaClient({
    log:
      process.env.NODE_ENV === "development"
        ? ["query", "warn", "error"]
        : ["error"],
  });

if (process.env.NODE_ENV !== "production") {
  globalForPrisma.prisma = prisma;
}
```

## One Transaction per Request for Multi-Table Writes

```typescript
import { PrismaClient } from "@prisma/client";
import { trace, SpanStatusCode } from "@opentelemetry/api";

const tracer = trace.getTracer("order-service");

interface Order {
  id: string;
  customerId: string;
  items: Array<{ productId: string; quantity: number; unitPrice: number }>;
}

async function placeOrder(prisma: PrismaClient, order: Order): Promise<void> {
  return tracer.startActiveSpan("placeOrder", async (span) => {
    span.setAttribute("order.id", order.id);
    span.setAttribute("order.item_count", order.items.length);

    try {
      await prisma.$transaction(async (tx) => {
        await tx.order.create({
          data: {
            id: order.id,
            customerId: order.customerId,
            items: {
              createMany: {
                data: order.items.map((item) => ({
                  productId: item.productId,
                  quantity: item.quantity,
                  unitPrice: item.unitPrice,
                })),
              },
            },
          },
        });

        for (const item of order.items) {
          await tx.product.update({
            where: { id: item.productId },
            data: { stock: { decrement: item.quantity } },
          });
        }
      });
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw new Error(`place order ${order.id}`, { cause: err });
    } finally {
      span.end();
    }
  });
}
```

## Batch Queries

### Batch Create with skipDuplicates

```typescript
import { PrismaClient } from "@prisma/client";
import { trace, SpanStatusCode } from "@opentelemetry/api";

const tracer = trace.getTracer("product-service");

interface ProductTag {
  productId: string;
  tag: string;
}

async function createProductTags(
  prisma: PrismaClient,
  tags: ProductTag[],
): Promise<number> {
  return tracer.startActiveSpan("createProductTags", async (span) => {
    span.setAttribute("batch.size", tags.length);

    try {
      const result = await prisma.productTag.createMany({
        data: tags,
        skipDuplicates: true,
      });

      span.setAttribute("batch.created_count", result.count);
      return result.count;
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw new Error("create product tags", { cause: err });
    } finally {
      span.end();
    }
  });
}
```

### Batch Read with `in`

```typescript
async function getProductsByIds(
  prisma: PrismaClient,
  ids: string[],
): Promise<Product[]> {
  return tracer.startActiveSpan("getProductsByIds", async (span) => {
    span.setAttribute("batch.size", ids.length);

    try {
      const products = await prisma.product.findMany({
        where: { id: { in: ids } },
      });

      span.setAttribute("batch.found_count", products.length);
      return products;
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw new Error("get products by IDs", { cause: err });
    } finally {
      span.end();
    }
  });
}
```

## Raw SQL with Parameterized Queries

Use `Prisma.sql` tagged templates for queries that Prisma ORM cannot express (window functions, CTEs, complex aggregations). Always parameterize -- never concatenate strings.

```typescript
import { PrismaClient, Prisma } from "@prisma/client";
import { trace, SpanStatusCode } from "@opentelemetry/api";

const tracer = trace.getTracer("analytics-service");

interface RankedProduct {
  id: string;
  name: string;
  categoryId: string;
  totalSales: number;
  rankInCategory: number;
}

async function getTopProductsByCategory(
  prisma: PrismaClient,
  categoryIds: string[],
  limit: number,
): Promise<RankedProduct[]> {
  return tracer.startActiveSpan("getTopProductsByCategory", async (span) => {
    span.setAttribute("category_count", categoryIds.length);
    span.setAttribute("limit", limit);

    try {
      const results = await prisma.$queryRaw<RankedProduct[]>(Prisma.sql`
        WITH ranked AS (
          SELECT
            p.id,
            p.name,
            p.category_id AS "categoryId",
            COALESCE(SUM(oi.quantity * oi.unit_price), 0) AS "totalSales",
            ROW_NUMBER() OVER (
              PARTITION BY p.category_id
              ORDER BY COALESCE(SUM(oi.quantity * oi.unit_price), 0) DESC
            ) AS "rankInCategory"
          FROM product p
          LEFT JOIN order_item oi ON oi.product_id = p.id
          WHERE p.category_id IN (${Prisma.join(categoryIds)})
          GROUP BY p.id, p.name, p.category_id
        )
        SELECT * FROM ranked
        WHERE "rankInCategory" <= ${limit}
        ORDER BY "categoryId", "rankInCategory"
      `);

      span.setAttribute("result_count", results.length);
      return results;
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw new Error("get top products by category", { cause: err });
    } finally {
      span.end();
    }
  });
}
```

## Injecting Time Values from Application Code

Pass `new Date()` from application code instead of relying on database functions like `NOW()`. This makes tests deterministic by injecting a `Clock` interface.

```typescript
import { PrismaClient } from "@prisma/client";
import { trace, SpanStatusCode } from "@opentelemetry/api";

const tracer = trace.getTracer("subscription-service");

interface Clock {
  now(): Date;
}

const systemClock: Clock = { now: () => new Date() };

async function cancelSubscription(
  prisma: PrismaClient,
  subscriptionId: string,
  clock: Clock = systemClock,
): Promise<void> {
  return tracer.startActiveSpan("cancelSubscription", async (span) => {
    span.setAttribute("subscription.id", subscriptionId);

    try {
      const cancelledAt = clock.now();

      await prisma.$transaction(async (tx) => {
        await tx.subscription.update({
          where: { id: subscriptionId },
          data: {
            status: "CANCELLED",
            cancelledAt,
          },
        });

        await tx.subscriptionEvent.create({
          data: {
            subscriptionId,
            type: "CANCELLED",
            occurredAt: cancelledAt,
          },
        });
      });
    } catch (err) {
      span.recordException(err as Error);
      span.setStatus({ code: SpanStatusCode.ERROR });
      throw new Error(`cancel subscription ${subscriptionId}`, { cause: err });
    } finally {
      span.end();
    }
  });
}
```

### Testing with a Fake Clock

```typescript
const fakeClock: Clock = {
  now: () => new Date("2025-06-15T10:00:00Z"),
};

test("cancelSubscription sets cancelledAt from clock", async () => {
  await cancelSubscription(prisma, subscriptionId, fakeClock);

  const subscription = await prisma.subscription.findUniqueOrThrow({
    where: { id: subscriptionId },
  });
  expect(subscription.cancelledAt).toEqual(new Date("2025-06-15T10:00:00Z"));
  expect(subscription.status).toBe("CANCELLED");
});
```
