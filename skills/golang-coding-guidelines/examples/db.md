# Database Patterns (database/sql + sqlx)

Idiomatic Go database patterns with connection pooling, transactions, batch queries, and inline OpenTelemetry instrumentation.

## Connection pool configuration

```go
import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}
```

## Repository with single-row query and OTel span

```go
import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	ctx, span := otel.Tracer("user-repo").Start(ctx, "FindByID")
	defer span.End()
	span.SetAttributes(attribute.Int64("user.id", id))

	var u User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, role, created_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("find user by ID %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find user by ID %d: %w", id, err)
	}
	return &u, nil
}
```

## Batch query with IN clause

```go
import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserRepositorySqlx struct {
	db *sqlx.DB
}

func (r *UserRepositorySqlx) FindByIDs(ctx context.Context, ids []int64) ([]User, error) {
	ctx, span := otel.Tracer("user-repo").Start(ctx, "FindByIDs")
	defer span.End()
	span.SetAttributes(attribute.Int("user.count", len(ids)))

	if len(ids) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT id, name, email, role, created_at FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("build IN query: %w", err)
	}
	query = r.db.Rebind(query)

	var users []User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, fmt.Errorf("find users by IDs: %w", err)
	}
	return users, nil
}
```

## One transaction per request for multi-table writes

```go
import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *OrderService) PlaceOrder(ctx context.Context, order Order) error {
	ctx, span := otel.Tracer("order").Start(ctx, "PlaceOrder")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("order.id", order.ID),
		attribute.Int("order.item_count", len(order.Items)),
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	if err := s.createOrder(ctx, tx, order, now); err != nil {
		return fmt.Errorf("create order %d: %w", order.ID, err)
	}
	if err := s.createOrderItems(ctx, tx, order.ID, order.Items); err != nil {
		return fmt.Errorf("create order items for order %d: %w", order.ID, err)
	}
	if err := s.updateInventory(ctx, tx, order.Items); err != nil {
		return fmt.Errorf("update inventory for order %d: %w", order.ID, err)
	}

	if err := tx.Commit(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("commit order %d: %w", order.ID, err)
	}
	return nil
}
```

## Insert with application-generated timestamps

```go
func (r *UserRepository) Create(ctx context.Context, u *User, now time.Time) error {
	ctx, span := otel.Tracer("user-repo").Start(ctx, "Create")
	defer span.End()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Email, u.Role, now, now,
	)
	if err != nil {
		return fmt.Errorf("insert user %q: %w", u.Email, err)
	}
	return nil
}
```

## Batch insert with sqlx

```go
func (r *OrderItemRepository) CreateBatch(ctx context.Context, tx *sql.Tx, items []OrderItem) error {
	ctx, span := otel.Tracer("order-item-repo").Start(ctx, "CreateBatch")
	defer span.End()
	span.SetAttributes(attribute.Int("item.count", len(items)))

	if len(items) == 0 {
		return nil
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.ExecContext(ctx, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice); err != nil {
			return fmt.Errorf("insert item for order %d, product %s: %w", item.OrderID, item.ProductID, err)
		}
	}
	return nil
}
```

## Testing with testcontainers

```go
import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	return db
}
```
