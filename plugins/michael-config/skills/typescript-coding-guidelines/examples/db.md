# Database Patterns

## Prisma Client Setup with Connection Pool

```typescript
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient({
  datasourceUrl: process.env.DATABASE_URL,
  log: process.env.LOG_LEVEL === "debug" ? ["query"] : [],
});

export { prisma };
```

## Repository with Interface

```typescript
interface UserRepository {
  getById(id: string): Promise<User | null>;
  create(data: CreateUserInput): Promise<User>;
  listByIds(ids: string[]): Promise<User[]>;
}

class PrismaUserRepository implements UserRepository {
  constructor(private readonly prisma: PrismaClient) {}

  async getById(id: string): Promise<User | null> {
    return tracer.startActiveSpan("UserRepository.getById", async (span) => {
      try {
        span.setAttribute("user.id", id);
        return await this.prisma.user.findUnique({ where: { id } });
      } finally {
        span.end();
      }
    });
  }

  async create(data: CreateUserInput): Promise<User> {
    return tracer.startActiveSpan("UserRepository.create", async (span) => {
      try {
        return await this.prisma.user.create({ data });
      } catch (err) {
        span.recordException(err instanceof Error ? err : new Error(String(err)));
        if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === "P2002") {
          throw new ConflictError(`user with email ${data.email} already exists`);
        }
        throw err;
      } finally {
        span.end();
      }
    });
  }

  async listByIds(ids: string[]): Promise<User[]> {
    if (ids.length === 0) return [];
    return this.prisma.user.findMany({ where: { id: { in: ids } } });
  }
}
```

## Transaction — One Per Request

```typescript
async function placeOrder(input: PlaceOrderInput): Promise<Order> {
  return tracer.startActiveSpan("placeOrder", async (span) => {
    try {
      return await prisma.$transaction(async (tx) => {
        const order = await tx.order.create({
          data: {
            userId: input.userId,
            total: input.total,
            status: "pending",
            createdAt: input.now,
          },
        });

        await tx.orderItem.createMany({
          data: input.items.map((item) => ({
            orderId: order.id,
            productId: item.productId,
            quantity: item.quantity,
            price: item.price,
          })),
        });

        await tx.inventory.updateMany({
          where: { productId: { in: input.items.map((i) => i.productId) } },
          data: { stock: { decrement: 1 } },
        });

        return order;
      });
    } catch (err) {
      span.recordException(err instanceof Error ? err : new Error(String(err)));
      throw new Error(`place order for user ${input.userId}`, { cause: err });
    } finally {
      span.end();
    }
  });
}
```

## Drizzle Alternative — Schema and Query

```typescript
import { pgTable, text, integer, timestamp } from "drizzle-orm/pg-core";

export const users = pgTable("users", {
  id: text("id").primaryKey(),
  name: text("name").notNull(),
  email: text("email").notNull().unique(),
  createdAt: timestamp("created_at").notNull(),
});

export const orders = pgTable("orders", {
  id: text("id").primaryKey(),
  userId: text("user_id")
    .notNull()
    .references(() => users.id),
  total: integer("total").notNull(),
  status: text("status").notNull(),
  createdAt: timestamp("created_at").notNull(),
});
```

```typescript
import { db } from "@/lib/db";
import { eq } from "drizzle-orm";

async function getUserById(id: string): Promise<User | undefined> {
  return tracer.startActiveSpan("getUserById", async (span) => {
    try {
      span.setAttribute("user.id", id);
      const rows = await db.select().from(users).where(eq(users.id, id)).limit(1);
      return rows[0];
    } finally {
      span.end();
    }
  });
}
```

## Drizzle Transaction

```typescript
async function placeOrder(input: PlaceOrderInput): Promise<void> {
  await db.transaction(async (tx) => {
    await tx.insert(orders).values({
      id: input.orderId,
      userId: input.userId,
      total: input.total,
      status: "pending",
      createdAt: input.now,
    });

    await tx.insert(orderItems).values(
      input.items.map((item) => ({
        orderId: input.orderId,
        productId: item.productId,
        quantity: item.quantity,
        price: item.price,
      })),
    );
  });
}
```

## Testing with testcontainers

```typescript
import { PostgreSqlContainer, StartedPostgreSqlContainer } from "@testcontainers/postgresql";
import { PrismaClient } from "@prisma/client";

let container: StartedPostgreSqlContainer;
let prisma: PrismaClient;

beforeAll(async () => {
  container = await new PostgreSqlContainer("postgres:16-alpine").start();
  prisma = new PrismaClient({
    datasourceUrl: container.getConnectionUri(),
  });
  await prisma.$executeRawUnsafe(
    await fs.readFile("prisma/migrations/init/migration.sql", "utf-8"),
  );
}, 30_000);

afterAll(async () => {
  await prisma.$disconnect();
  await container.stop();
});

beforeEach(async () => {
  await prisma.user.deleteMany();
});
```
