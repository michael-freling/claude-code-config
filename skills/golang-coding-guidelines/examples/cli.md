# CLI Tool Patterns

Idiomatic Go CLI tools using `cobra` with configuration, structured logging, and graceful shutdown.

## Root command with cobra

```go
import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var verbose bool

	root := &cobra.Command{
		Use:   "myapp",
		Short: "My application CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := slog.LevelInfo
			if verbose {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
		},
	}

	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
	root.AddCommand(newServeCmd())
	root.AddCommand(newMigrateCmd())
	return root
}
```

## Subcommand with flag validation

```go
func newServeCmd() *cobra.Command {
	var (
		addr        string
		databaseURL string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if databaseURL == "" {
				return fmt.Errorf("--database-url is required")
			}
			return runServer(cmd.Context(), addr, databaseURL)
		},
	}

	cmd.Flags().StringVar(&addr, "addr", ":8080", "listen address")
	cmd.Flags().StringVar(&databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	return cmd
}
```

## One-time job (migration, data backfill)

```go
import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate(cmd.Context())
		},
	}
}

func runMigrate(ctx context.Context) error {
	ctx, span := otel.Tracer("cli").Start(ctx, "migrate")
	defer span.End()

	db, err := openDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	start := time.Now()
	applied, err := migrate(ctx, db)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	span.SetAttributes(attribute.Int("migrations.applied", applied))
	slog.InfoContext(ctx, "migrations complete",
		"applied", applied,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return nil
}
```

## Graceful shutdown with context

```go
import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
)

func runServer(ctx context.Context, addr, databaseURL string) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := openDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := newHTTPServer(addr, db)

	errCh := make(chan error, 1)
	go func() {
		slog.InfoContext(ctx, "server starting", "addr", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
		slog.InfoContext(ctx, "shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
```

## Data backfill job with progress logging

```go
func newBackfillCmd() *cobra.Command {
	var batchSize int

	cmd := &cobra.Command{
		Use:   "backfill",
		Short: "Backfill user display names",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackfill(cmd.Context(), batchSize)
		},
	}

	cmd.Flags().IntVar(&batchSize, "batch-size", 1000, "number of records per batch")
	return cmd
}

func runBackfill(ctx context.Context, batchSize int) error {
	ctx, span := otel.Tracer("cli").Start(ctx, "backfill")
	defer span.End()

	db, err := openDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	var cursor int64
	var totalProcessed int

	for {
		if ctx.Err() != nil {
			slog.InfoContext(ctx, "backfill interrupted", "processed", totalProcessed)
			return ctx.Err()
		}

		batch, nextCursor, err := fetchBatch(ctx, db, cursor, batchSize)
		if err != nil {
			return fmt.Errorf("fetch batch at cursor %d: %w", cursor, err)
		}
		if len(batch) == 0 {
			break
		}

		if err := processBatch(ctx, db, batch); err != nil {
			return fmt.Errorf("process batch at cursor %d: %w", cursor, err)
		}

		totalProcessed += len(batch)
		cursor = nextCursor

		slog.InfoContext(ctx, "backfill progress",
			"processed", totalProcessed,
			"cursor", cursor,
		)
	}

	span.SetAttributes(attribute.Int("backfill.total", totalProcessed))
	slog.InfoContext(ctx, "backfill complete", "total", totalProcessed)
	return nil
}
```
