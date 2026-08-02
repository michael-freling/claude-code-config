# gRPC Server Patterns

## Service Implementation with Tracing

```go
type UserServer struct {
    pb.UnimplementedUserServiceServer
    userService UserService
    logger      *slog.Logger
}

func NewUserServer(userService UserService, logger *slog.Logger) *UserServer {
    return &UserServer{
        userService: userService,
        logger:      logger,
    }
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    ctx, span := otel.Tracer("grpc").Start(ctx, "GetUser")
    defer span.End()

    if req.GetUserId() == 0 {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }

    user, err := s.userService.GetByID(ctx, req.GetUserId())
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, status.Errorf(codes.NotFound, "user %d not found", req.GetUserId())
        }
        span.RecordError(err)
        s.logger.ErrorContext(ctx, "get user", "user_id", req.GetUserId(), "error", err)
        return nil, status.Error(codes.Internal, "internal error")
    }

    return &pb.GetUserResponse{
        User: toProtoUser(user),
    }, nil
}
```

## Error Mapping Helper

```go
func domainToGRPCError(err error) error {
    switch {
    case errors.Is(err, ErrNotFound):
        return status.Error(codes.NotFound, err.Error())
    case errors.Is(err, ErrInvalidInput):
        return status.Error(codes.InvalidArgument, err.Error())
    case errors.Is(err, ErrConflict):
        return status.Error(codes.AlreadyExists, err.Error())
    case errors.Is(err, ErrUnauthorized):
        return status.Error(codes.Unauthenticated, err.Error())
    default:
        return status.Error(codes.Internal, "internal error")
    }
}
```

## Unary Interceptor for Logging and Metrics

```go
func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req any,
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (any, error) {
        start := time.Now()

        resp, err := handler(ctx, req)

        duration := time.Since(start)
        code := status.Code(err)

        logger.InfoContext(ctx, "grpc request",
            "method", info.FullMethod,
            "code", code.String(),
            "duration_ms", duration.Milliseconds(),
        )

        return resp, err
    }
}
```

## Server Setup with Interceptors

```go
func NewGRPCServer(
    userServer pb.UserServiceServer,
    orderServer pb.OrderServiceServer,
    logger *slog.Logger,
    tp trace.TracerProvider,
) *grpc.Server {
    server := grpc.NewServer(
        grpc.StatsHandler(otelgrpc.NewServerHandler(
            otelgrpc.WithTracerProvider(tp),
        )),
        grpc.ChainUnaryInterceptor(
            LoggingInterceptor(logger),
        ),
    )

    pb.RegisterUserServiceServer(server, userServer)
    pb.RegisterOrderServiceServer(server, orderServer)

    return server
}
```

## Server-Side Streaming

```go
func (s *UserServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    ctx, span := otel.Tracer("grpc").Start(stream.Context(), "ListUsers")
    defer span.End()

    cursor := ""
    for {
        users, nextCursor, err := s.userService.List(ctx, req.GetPageSize(), cursor)
        if err != nil {
            span.RecordError(err)
            return status.Error(codes.Internal, "internal error")
        }

        for _, user := range users {
            if err := stream.Send(&pb.ListUsersResponse{User: toProtoUser(user)}); err != nil {
                return fmt.Errorf("send user: %w", err)
            }
        }

        if nextCursor == "" {
            break
        }
        cursor = nextCursor
    }

    return nil
}
```

## Graceful Shutdown

```go
func RunGRPC(ctx context.Context, server *grpc.Server, addr string) error {
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("listen %s: %w", addr, err)
    }

    errCh := make(chan error, 1)
    go func() {
        errCh <- server.Serve(lis)
    }()

    select {
    case err := <-errCh:
        return fmt.Errorf("serve: %w", err)
    case <-ctx.Done():
        server.GracefulStop()
        return nil
    }
}
```
