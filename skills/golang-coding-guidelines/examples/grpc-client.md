# gRPC Client Patterns

## Client with Dial Options and Tracing

```go
type UserClient struct {
    conn   *grpc.ClientConn
    client pb.UserServiceClient
    tracer trace.Tracer
}

func NewUserClient(addr string, tp trace.TracerProvider) (*UserClient, error) {
    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler(
            otelgrpc.WithTracerProvider(tp),
        )),
    )
    if err != nil {
        return nil, fmt.Errorf("dial %s: %w", addr, err)
    }

    return &UserClient{
        conn:   conn,
        client: pb.NewUserServiceClient(conn),
        tracer: tp.Tracer("grpc-client"),
    }, nil
}

func (c *UserClient) Close() error {
    return c.conn.Close()
}
```

## Unary Call with Error Handling

```go
func (c *UserClient) GetUser(ctx context.Context, userID int64) (*User, error) {
    ctx, span := c.tracer.Start(ctx, "GetUser")
    defer span.End()

    resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{UserId: userID})
    if err != nil {
        span.RecordError(err)
        st, ok := status.FromError(err)
        if ok && st.Code() == codes.NotFound {
            return nil, fmt.Errorf("get user %d: %w", userID, ErrNotFound)
        }
        return nil, fmt.Errorf("get user %d: %w", userID, err)
    }

    return fromProtoUser(resp.GetUser()), nil
}
```

## Retry Configuration

```go
retryPolicy := `{
    "methodConfig": [{
        "name": [{"service": "user.v1.UserService"}],
        "retryPolicy": {
            "maxAttempts": 3,
            "initialBackoff": "0.1s",
            "maxBackoff": "1s",
            "backoffMultiplier": 2,
            "retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
        }
    }]
}`

conn, err := grpc.NewClient(addr,
    grpc.WithDefaultServiceConfig(retryPolicy),
)
```

## Consuming Server-Side Streaming

```go
func (c *UserClient) ListUsers(ctx context.Context, pageSize int32) ([]*User, error) {
    ctx, span := c.tracer.Start(ctx, "ListUsers")
    defer span.End()

    stream, err := c.client.ListUsers(ctx, &pb.ListUsersRequest{PageSize: pageSize})
    if err != nil {
        return nil, fmt.Errorf("list users: %w", err)
    }

    var users []*User
    for {
        resp, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            break
        }
        if err != nil {
            span.RecordError(err)
            return nil, fmt.Errorf("recv user: %w", err)
        }
        users = append(users, fromProtoUser(resp.GetUser()))
    }

    return users, nil
}
```

## Concurrent Fan-Out with errgroup

```go
func (c *UserClient) GetUsers(ctx context.Context, ids []int64) ([]*User, error) {
    ctx, span := c.tracer.Start(ctx, "GetUsers")
    defer span.End()

    users := make([]*User, len(ids))
    eg, ctx := errgroup.WithContext(ctx)

    for i, id := range ids {
        eg.Go(func() error {
            user, err := c.GetUser(ctx, id)
            if err != nil {
                return err
            }
            users[i] = user
            return nil
        })
    }

    if err := eg.Wait(); err != nil {
        return nil, err
    }
    return users, nil
}
```
