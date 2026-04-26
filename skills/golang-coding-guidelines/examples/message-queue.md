# Message Queue Worker Patterns

## Consumer with Graceful Shutdown

```go
type Worker struct {
    consumer MessageConsumer
    handler  MessageHandler
    logger   *slog.Logger
    tracer   trace.Tracer
}

type MessageConsumer interface {
    Receive(ctx context.Context) (Message, error)
    Ack(ctx context.Context, msg Message) error
    Nack(ctx context.Context, msg Message) error
}

type MessageHandler interface {
    Handle(ctx context.Context, msg Message) error
}

func NewWorker(consumer MessageConsumer, handler MessageHandler, logger *slog.Logger) *Worker {
    return &Worker{
        consumer: consumer,
        handler:  handler,
        logger:   logger,
        tracer:   otel.Tracer("worker"),
    }
}
```

## Consumer Loop with Tracing

```go
func (w *Worker) Run(ctx context.Context) error {
    w.logger.InfoContext(ctx, "worker started")

    for {
        select {
        case <-ctx.Done():
            w.logger.InfoContext(ctx, "worker shutting down")
            return nil
        default:
        }

        msg, err := w.consumer.Receive(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return nil
            }
            w.logger.ErrorContext(ctx, "receive message", "error", err)
            continue
        }

        if err := w.processMessage(ctx, msg); err != nil {
            w.logger.ErrorContext(ctx, "process message",
                "message_id", msg.ID,
                "error", err,
            )
        }
    }
}

func (w *Worker) processMessage(ctx context.Context, msg Message) error {
    ctx, span := w.tracer.Start(ctx, "ProcessMessage",
        trace.WithAttributes(attribute.String("message.id", msg.ID)),
    )
    defer span.End()

    if err := w.handler.Handle(ctx, msg); err != nil {
        span.RecordError(err)
        if nackErr := w.consumer.Nack(ctx, msg); nackErr != nil {
            w.logger.ErrorContext(ctx, "nack message", "message_id", msg.ID, "error", nackErr)
        }
        return fmt.Errorf("handle message %s: %w", msg.ID, err)
    }

    if err := w.consumer.Ack(ctx, msg); err != nil {
        span.RecordError(err)
        return fmt.Errorf("ack message %s: %w", msg.ID, err)
    }
    return nil
}
```

## Concurrent Worker Pool

```go
func (w *Worker) RunPool(ctx context.Context, concurrency int) error {
    eg, ctx := errgroup.WithContext(ctx)

    for range concurrency {
        eg.Go(func() error {
            return w.Run(ctx)
        })
    }

    return eg.Wait()
}
```

## Batch Consumer

```go
func (w *BatchWorker) Run(ctx context.Context) error {
    batch := make([]Message, 0, w.batchSize)
    ticker := time.NewTicker(w.flushInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            if len(batch) > 0 {
                w.flush(context.Background(), batch)
            }
            return nil

        case <-ticker.C:
            if len(batch) > 0 {
                w.flush(ctx, batch)
                batch = batch[:0]
            }

        default:
            msg, err := w.consumer.Receive(ctx)
            if err != nil {
                if ctx.Err() != nil {
                    return nil
                }
                continue
            }
            batch = append(batch, msg)
            if len(batch) >= w.batchSize {
                w.flush(ctx, batch)
                batch = batch[:0]
            }
        }
    }
}

func (w *BatchWorker) flush(ctx context.Context, batch []Message) {
    ctx, span := w.tracer.Start(ctx, "FlushBatch",
        trace.WithAttributes(attribute.Int("batch.size", len(batch))),
    )
    defer span.End()

    if err := w.handler.HandleBatch(ctx, batch); err != nil {
        span.RecordError(err)
        w.logger.ErrorContext(ctx, "flush batch", "size", len(batch), "error", err)
        for _, msg := range batch {
            w.consumer.Nack(ctx, msg)
        }
        return
    }

    for _, msg := range batch {
        w.consumer.Ack(ctx, msg)
    }
}
```

## Testing with Mock Consumer

```go
func TestWorker_ProcessMessage(t *testing.T) {
    ctrl := gomock.NewController(t)

    consumer := NewMockMessageConsumer(ctrl)
    handler := NewMockMessageHandler(ctrl)
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))

    worker := NewWorker(consumer, handler, logger)

    tests := []struct {
        name    string
        msg     Message
        setup   func()
        wantErr bool
    }{
        {
            name: "success",
            msg:  Message{ID: "msg-1", Body: []byte(`{"action":"create"}`)},
            setup: func() {
                handler.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(nil)
                consumer.EXPECT().Ack(gomock.Any(), gomock.Any()).Return(nil)
            },
        },
        {
            name: "handler error nacks message",
            msg:  Message{ID: "msg-2", Body: []byte(`{"action":"invalid"}`)},
            setup: func() {
                handler.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(errors.New("invalid action"))
                consumer.EXPECT().Nack(gomock.Any(), gomock.Any()).Return(nil)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            err := worker.processMessage(context.Background(), tt.msg)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```
