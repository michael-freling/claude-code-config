# Message Queue Patterns

Idiomatic Go message queue consumer and producer using NATS JetStream with inline OpenTelemetry instrumentation. Patterns apply to any message broker (Kafka, RabbitMQ, SQS).

## Consumer with graceful shutdown

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type OrderConsumer struct {
	consumer jetstream.Consumer
	handler  OrderHandler
	logger   *slog.Logger
}

type OrderHandler interface {
	ProcessOrder(ctx context.Context, order Order) error
}

func NewOrderConsumer(js jetstream.JetStream, handler OrderHandler, logger *slog.Logger) (*OrderConsumer, error) {
	consumer, err := js.CreateOrUpdateConsumer(context.Background(), "ORDERS", jetstream.ConsumerConfig{
		Durable:       "order-processor",
		FilterSubject: "orders.created",
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	return &OrderConsumer{
		consumer: consumer,
		handler:  handler,
		logger:   logger,
	}, nil
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	for {
		msg, err := c.consumer.Next(jetstream.FetchMaxWait(5 * time.Second))
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			c.logger.WarnContext(ctx, "fetch message", "error", err)
			continue
		}

		if err := c.processMessage(ctx, msg); err != nil {
			c.logger.ErrorContext(ctx, "process message",
				"subject", msg.Subject(),
				"error", err,
			)
			msg.Nak()
			continue
		}
		msg.Ack()
	}
}
```

## Message processing with OTel span

```go
func (c *OrderConsumer) processMessage(ctx context.Context, msg jetstream.Msg) error {
	ctx, span := otel.Tracer("consumer").Start(ctx, "ProcessOrderMessage")
	defer span.End()
	span.SetAttributes(attribute.String("subject", msg.Subject()))

	var order Order
	if err := json.Unmarshal(msg.Data(), &order); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unmarshal failed")
		return fmt.Errorf("unmarshal order message: %w", err)
	}

	span.SetAttributes(attribute.Int64("order.id", order.ID))

	if err := c.handler.ProcessOrder(ctx, order); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("process order %d: %w", order.ID, err)
	}

	c.logger.InfoContext(ctx, "order processed",
		"order_id", order.ID,
		"item_count", len(order.Items),
	)
	return nil
}
```

## Producer with publish confirmation

```go
type OrderPublisher struct {
	js     jetstream.JetStream
	logger *slog.Logger
}

func NewOrderPublisher(js jetstream.JetStream, logger *slog.Logger) *OrderPublisher {
	return &OrderPublisher{js: js, logger: logger}
}

func (p *OrderPublisher) PublishOrderCreated(ctx context.Context, order Order) error {
	ctx, span := otel.Tracer("publisher").Start(ctx, "PublishOrderCreated")
	defer span.End()
	span.SetAttributes(attribute.Int64("order.id", order.ID))

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal order %d: %w", order.ID, err)
	}

	ack, err := p.js.Publish(ctx, "orders.created", data)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("publish order %d: %w", order.ID, err)
	}

	p.logger.InfoContext(ctx, "order published",
		"order_id", order.ID,
		"stream", ack.Stream,
		"sequence", ack.Sequence,
	)
	return nil
}
```

## Batch publish

```go
func (p *OrderPublisher) PublishBatch(ctx context.Context, orders []Order) error {
	ctx, span := otel.Tracer("publisher").Start(ctx, "PublishBatch")
	defer span.End()
	span.SetAttributes(attribute.Int("batch.size", len(orders)))

	for _, order := range orders {
		if err := p.PublishOrderCreated(ctx, order); err != nil {
			return err
		}
	}
	return nil
}
```

## Worker main with graceful shutdown

```go
import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func runWorker() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return fmt.Errorf("connect nats: %w", err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("create jetstream: %w", err)
	}

	consumer, err := NewOrderConsumer(js, orderHandler, slog.Default())
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	slog.Info("worker starting")
	return consumer.Run(ctx)
}
```
