# Message Queue Worker Patterns

## BullMQ Worker with Tracing

```typescript
import { Worker, Job } from "bullmq";
import { tracer } from "@/lib/telemetry";

interface OrderEventPayload {
  orderId: string;
  userId: string;
  action: "created" | "shipped" | "cancelled";
}

function createOrderWorker(
  orderService: OrderService,
  logger: Logger,
  connection: ConnectionOptions,
): Worker {
  return new Worker<OrderEventPayload>(
    "order-events",
    async (job: Job<OrderEventPayload>) => {
      await tracer.startActiveSpan(`process:${job.name}`, async (span) => {
        try {
          span.setAttribute("job.id", job.id ?? "unknown");
          span.setAttribute("order.id", job.data.orderId);

          await processOrderEvent(orderService, job.data);

          logger.info("processed order event", {
            jobId: job.id,
            orderId: job.data.orderId,
            action: job.data.action,
          });
        } catch (err) {
          span.recordException(err instanceof Error ? err : new Error(String(err)));
          logger.error("failed to process order event", {
            jobId: job.id,
            orderId: job.data.orderId,
            error: err,
          });
          throw err;
        } finally {
          span.end();
        }
      });
    },
    {
      connection,
      concurrency: 5,
      limiter: { max: 100, duration: 60_000 },
    },
  );
}
```

## Event Handler Dispatch

```typescript
async function processOrderEvent(
  orderService: OrderService,
  payload: OrderEventPayload,
): Promise<void> {
  switch (payload.action) {
    case "created":
      await orderService.sendConfirmation(payload.orderId);
      break;
    case "shipped":
      await orderService.sendShipmentNotification(payload.orderId);
      break;
    case "cancelled":
      await orderService.refund(payload.orderId);
      break;
  }
}
```

## Producer — Enqueue Jobs

```typescript
import { Queue } from "bullmq";

class OrderEventProducer {
  private readonly queue: Queue<OrderEventPayload>;

  constructor(connection: ConnectionOptions) {
    this.queue = new Queue("order-events", { connection });
  }

  async enqueue(payload: OrderEventPayload): Promise<void> {
    await this.queue.add(payload.action, payload, {
      attempts: 3,
      backoff: { type: "exponential", delay: 1000 },
      removeOnComplete: 1000,
      removeOnFail: 5000,
    });
  }

  async close(): Promise<void> {
    await this.queue.close();
  }
}
```

## Graceful Shutdown

```typescript
function startWorker(worker: Worker, logger: Logger): void {
  const shutdown = async (signal: string) => {
    logger.info("received signal, shutting down", { signal });
    await worker.close();
    logger.info("worker stopped");
    process.exit(0);
  };

  process.on("SIGINT", () => shutdown("SIGINT"));
  process.on("SIGTERM", () => shutdown("SIGTERM"));

  worker.on("error", (err) => {
    logger.error("worker error", { error: err });
  });

  worker.on("failed", (job, err) => {
    logger.warn("job failed", {
      jobId: job?.id,
      attemptsMade: job?.attemptsMade,
      error: err.message,
    });
  });
}
```

## Batch Processing with Rate Limiting

```typescript
function createBatchWorker(
  batchProcessor: BatchProcessor,
  logger: Logger,
  connection: ConnectionOptions,
): Worker {
  return new Worker<BatchPayload>(
    "batch-import",
    async (job: Job<BatchPayload>) => {
      await tracer.startActiveSpan("batch-import", async (span) => {
        try {
          const { items } = job.data;
          span.setAttribute("batch.size", items.length);

          const chunks = chunk(items, 100);
          for (let i = 0; i < chunks.length; i++) {
            await batchProcessor.process(chunks[i]);
            await job.updateProgress(Math.round(((i + 1) / chunks.length) * 100));
          }
        } finally {
          span.end();
        }
      });
    },
    { connection, concurrency: 2 },
  );
}

function chunk<T>(array: T[], size: number): T[][] {
  const chunks: T[][] = [];
  for (let i = 0; i < array.length; i += size) {
    chunks.push(array.slice(i, i + size));
  }
  return chunks;
}
```

## Testing Workers

```typescript
describe("processOrderEvent", () => {
  const orderService = {
    sendConfirmation: vi.fn(),
    sendShipmentNotification: vi.fn(),
    refund: vi.fn(),
  };

  beforeEach(() => {
    vi.resetAllMocks();
  });

  const cases = [
    {
      name: "created event sends confirmation",
      payload: { orderId: "o-1", userId: "u-1", action: "created" as const },
      expectedCall: "sendConfirmation",
    },
    {
      name: "shipped event sends notification",
      payload: { orderId: "o-2", userId: "u-1", action: "shipped" as const },
      expectedCall: "sendShipmentNotification",
    },
    {
      name: "cancelled event triggers refund",
      payload: { orderId: "o-3", userId: "u-1", action: "cancelled" as const },
      expectedCall: "refund",
    },
  ];

  cases.forEach(({ name, payload, expectedCall }) => {
    it(name, async () => {
      await processOrderEvent(orderService as unknown as OrderService, payload);
      expect(orderService[expectedCall]).toHaveBeenCalledWith(payload.orderId);
    });
  });
});
```
