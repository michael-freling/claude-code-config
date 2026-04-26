# gRPC Client Patterns

Idiomatic Go gRPC client with retries, deadlines, and inline OpenTelemetry instrumentation.

## Client setup with connection management

```go
import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewUserClient(addr string, opts ...Option) (*UserClient, error) {
	o := options{timeout: 10 * time.Second}
	for _, opt := range opts {
		opt(&o)
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": ""}],
				"retryPolicy": {
					"maxAttempts": 3,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2,
					"retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
				}
			}]
		}`),
	)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}

	return &UserClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}
```

## RPC call with context deadline and OTel span

```go
import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *UserClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	ctx, span := otel.Tracer("user-client").Start(ctx, "GetUser")
	defer span.End()
	span.SetAttributes(attribute.Int64("user.id", userID))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{UserId: userID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == grpccodes.NotFound {
			return nil, fmt.Errorf("user %d: %w", userID, ErrNotFound)
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("get user %d: %w", userID, err)
	}

	return &User{
		ID:    resp.GetUser().GetId(),
		Name:  resp.GetUser().GetName(),
		Email: resp.GetUser().GetEmail(),
	}, nil
}
```

## Batch RPC calls with errgroup

```go
import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

func (c *UserClient) GetUsers(ctx context.Context, ids []int64) ([]*User, error) {
	ctx, span := otel.Tracer("user-client").Start(ctx, "GetUsers")
	defer span.End()

	var mu sync.Mutex
	users := make([]*User, 0, len(ids))

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(10)
	for _, id := range ids {
		eg.Go(func() error {
			u, err := c.GetUser(ctx, id)
			if err != nil {
				return err
			}
			mu.Lock()
			users = append(users, u)
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}
```

## Error status mapping — gRPC to domain errors

```go
import (
	"errors"

	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapGRPCError(err error, entity string, id any) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case grpccodes.NotFound:
		return fmt.Errorf("%s %v: %w", entity, id, ErrNotFound)
	case grpccodes.InvalidArgument:
		return fmt.Errorf("%s %v: %w: %s", entity, id, ErrValidation, st.Message())
	case grpccodes.AlreadyExists:
		return fmt.Errorf("%s %v: %w", entity, id, ErrAlreadyExists)
	default:
		return fmt.Errorf("%s %v: %w", entity, id, err)
	}
}
```
