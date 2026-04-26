# gRPC Server Patterns

Idiomatic Go gRPC server using `google.golang.org/grpc` with inline OpenTelemetry instrumentation.

## Server bootstrap with interceptors

```go
import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func newGRPCServer(userSvc *user.Service) *grpc.Server {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			recoveryInterceptor,
		),
	)

	pb.RegisterUserServiceServer(srv, &UserServer{userSvc: userSvc})
	reflection.Register(srv)
	return srv
}

func listenGRPC(srv *grpc.Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	return srv.Serve(lis)
}
```

## Service implementation with error mapping

```go
import (
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userSvc *user.Service
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	ctx, span := otel.Tracer("grpc").Start(ctx, "GetUser")
	defer span.End()
	span.SetAttributes(attribute.Int64("user.id", req.GetUserId()))

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	u, err := s.userSvc.FindByID(ctx, req.GetUserId())
	if err != nil {
		return nil, s.mapError(ctx, span, err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := otel.Tracer("grpc").Start(ctx, "CreateUser")
	defer span.End()

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	u, err := s.userSvc.Create(ctx, user.CreateParams{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, s.mapError(ctx, span, err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{Id: u.ID, Name: u.Name, Email: u.Email},
	}, nil
}
```

## Error mapping — domain errors to gRPC status codes

```go
func (s *UserServer) mapError(ctx context.Context, span trace.Span, err error) error {
	var validationErr *user.ValidationError
	switch {
	case errors.Is(err, user.ErrNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, user.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.As(err, &validationErr):
		return status.Error(codes.InvalidArgument, validationErr.Error())
	default:
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		slog.ErrorContext(ctx, "internal error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}
```

## Logging interceptor

```go
import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	st, _ := status.FromError(err)
	slog.InfoContext(ctx, "grpc request",
		"method", info.FullMethod,
		"code", st.Code().String(),
		"duration_ms", duration.Milliseconds(),
	)
	return resp, err
}
```

## Recovery interceptor

```go
import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "panic recovered",
				"method", info.FullMethod,
				"panic", r,
				"stack", string(debug.Stack()),
			)
			err = status.Error(codes.Internal, "internal error")
		}
	}()
	return handler(ctx, req)
}
```

## Graceful shutdown

```go
func shutdownGRPC(srv *grpc.Server) {
	stopped := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(30 * time.Second):
		srv.Stop()
	}
}
```
