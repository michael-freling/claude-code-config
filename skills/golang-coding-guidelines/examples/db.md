# Database Patterns

## Connection Pool Configuration

```go
func NewDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }

    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)

    if err := db.PingContext(context.Background()); err != nil {
        return nil, fmt.Errorf("ping db: %w", err)
    }
    return db, nil
}
```

## Repository with Interface

```go
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, tx *sql.Tx, user *User) error
    ListByIDs(ctx context.Context, ids []int64) ([]*User, error)
}

type userRepository struct {
    db     *sql.DB
    tracer trace.Tracer
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{
        db:     db,
        tracer: otel.Tracer("db"),
    }
}
```

## Query with Tracing

```go
func (r *userRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    ctx, span := r.tracer.Start(ctx, "UserRepository.GetByID")
    defer span.End()

    var user User
    err := r.db.QueryRowContext(ctx,
        `SELECT id, name, email, created_at FROM users WHERE id = $1`, id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("get user by ID %d: %w", id, ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("get user by ID %d: %w", id, err)
    }
    return &user, nil
}
```

## Batch Query with sqlx.In

```go
func (r *userRepository) ListByIDs(ctx context.Context, ids []int64) ([]*User, error) {
    ctx, span := r.tracer.Start(ctx, "UserRepository.ListByIDs")
    defer span.End()

    if len(ids) == 0 {
        return nil, nil
    }

    query, args, err := sqlx.In(
        `SELECT id, name, email, created_at FROM users WHERE id IN (?)`, ids,
    )
    if err != nil {
        return nil, fmt.Errorf("build IN query for %d users: %w", len(ids), err)
    }
    query = r.db.Rebind(query)

    var users []*User
    if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
        return nil, fmt.Errorf("list users by IDs: %w", err)
    }
    return users, nil
}
```

## Transaction — One Per Request

```go
func (s *OrderService) PlaceOrder(ctx context.Context, order Order) error {
    ctx, span := otel.Tracer("order").Start(ctx, "PlaceOrder")
    defer span.End()

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    if err := s.orderRepo.Create(ctx, tx, &order); err != nil {
        return fmt.Errorf("create order: %w", err)
    }
    if err := s.inventoryRepo.DeductStock(ctx, tx, order.Items); err != nil {
        return fmt.Errorf("deduct stock: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit order tx: %w", err)
    }
    return nil
}
```

## Insert within Transaction

```go
func (r *orderRepository) Create(ctx context.Context, tx *sql.Tx, order *Order) error {
    ctx, span := r.tracer.Start(ctx, "OrderRepository.Create")
    defer span.End()

    err := tx.QueryRowContext(ctx,
        `INSERT INTO orders (user_id, total, status, created_at)
         VALUES ($1, $2, $3, $4) RETURNING id`,
        order.UserID, order.Total, order.Status, order.CreatedAt,
    ).Scan(&order.ID)
    if err != nil {
        return fmt.Errorf("insert order for user %d: %w", order.UserID, err)
    }
    return nil
}
```

## Batch Insert

```go
func (r *orderRepository) CreateItems(ctx context.Context, tx *sql.Tx, orderID int64, items []OrderItem) error {
    ctx, span := r.tracer.Start(ctx, "OrderRepository.CreateItems")
    defer span.End()

    if len(items) == 0 {
        return nil
    }

    valueStrings := make([]string, 0, len(items))
    valueArgs := make([]any, 0, len(items)*4)
    for i, item := range items {
        base := i * 4
        valueStrings = append(valueStrings,
            fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4),
        )
        valueArgs = append(valueArgs, orderID, item.ProductID, item.Quantity, item.Price)
    }

    query := fmt.Sprintf(
        `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES %s`,
        strings.Join(valueStrings, ", "),
    )

    if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
        return fmt.Errorf("insert %d items for order %d: %w", len(items), orderID, err)
    }
    return nil
}
```

## Testing with testcontainers-go

```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    ctx := context.Background()
    container, err := postgres.Run(ctx, "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        testcontainers.WithWaitStrategy(
            wait.ForListeningPort("5432/tcp").WithStartupTimeout(30*time.Second),
        ),
    )
    require.NoError(t, err)
    t.Cleanup(func() { container.Terminate(ctx) })

    dsn, err := container.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := sql.Open("postgres", dsn)
    require.NoError(t, err)
    t.Cleanup(func() { db.Close() })

    return db
}
```
