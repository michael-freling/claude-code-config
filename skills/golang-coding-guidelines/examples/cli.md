# CLI Tool Patterns

## Root Command with cobra

```go
func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mytool",
        Short: "A tool for managing resources",
    }

    cmd.AddCommand(
        NewSyncCmd(),
        NewMigrateCmd(),
    )

    return cmd
}

func main() {
    if err := NewRootCmd().Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Subcommand with Flags and Tracing

```go
func NewSyncCmd() *cobra.Command {
    var (
        dryRun  bool
        workers int
    )

    cmd := &cobra.Command{
        Use:   "sync",
        Short: "Sync resources from source to destination",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := cmd.Context()

            shutdown, err := initTracer(ctx)
            if err != nil {
                return fmt.Errorf("init tracer: %w", err)
            }
            defer shutdown(ctx)

            ctx, span := otel.Tracer("cli").Start(ctx, "sync")
            defer span.End()

            cfg, err := loadConfig()
            if err != nil {
                return fmt.Errorf("load config: %w", err)
            }

            syncer, err := NewSyncer(cfg)
            if err != nil {
                return fmt.Errorf("create syncer: %w", err)
            }

            result, err := syncer.Run(ctx, SyncOptions{
                DryRun:  dryRun,
                Workers: workers,
            })
            if err != nil {
                span.RecordError(err)
                return fmt.Errorf("sync: %w", err)
            }

            fmt.Fprintf(cmd.OutOrStdout(), "Synced %d resources\n", result.Count)
            return nil
        },
    }

    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying")
    cmd.Flags().IntVar(&workers, "workers", 4, "number of concurrent workers")

    return cmd
}
```

## Config Loading from Environment and Flags

```go
type Config struct {
    DatabaseURL string
    APIBaseURL  string
    LogLevel    slog.Level
}

func loadConfig() (*Config, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, errors.New("DATABASE_URL is required")
    }

    apiURL := os.Getenv("API_BASE_URL")
    if apiURL == "" {
        return nil, errors.New("API_BASE_URL is required")
    }

    return &Config{
        DatabaseURL: dbURL,
        APIBaseURL:  apiURL,
        LogLevel:    parseLogLevel(os.Getenv("LOG_LEVEL")),
    }, nil
}

func parseLogLevel(s string) slog.Level {
    switch strings.ToLower(s) {
    case "debug":
        return slog.LevelDebug
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo
    }
}
```

## Structured Output for Piping

```go
func NewListCmd() *cobra.Command {
    var outputJSON bool

    cmd := &cobra.Command{
        Use:   "list",
        Short: "List all resources",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := cmd.Context()

            resources, err := listResources(ctx)
            if err != nil {
                return fmt.Errorf("list resources: %w", err)
            }

            if outputJSON {
                enc := json.NewEncoder(cmd.OutOrStdout())
                enc.SetIndent("", "  ")
                return enc.Encode(resources)
            }

            w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
            fmt.Fprintln(w, "ID\tNAME\tSTATUS")
            for _, r := range resources {
                fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID, r.Name, r.Status)
            }
            return w.Flush()
        },
    }

    cmd.Flags().BoolVar(&outputJSON, "json", false, "output as JSON")

    return cmd
}
```

## Signal Handling and Graceful Shutdown

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    cmd := NewRootCmd()
    if err := cmd.ExecuteContext(ctx); err != nil {
        os.Exit(1)
    }
}
```

## Testing CLI Commands

```go
func TestSyncCmd(t *testing.T) {
    tests := []struct {
        name       string
        args       []string
        env        map[string]string
        wantOutput string
        wantErr    bool
    }{
        {
            name: "dry run",
            args: []string{"sync", "--dry-run"},
            env: map[string]string{
                "DATABASE_URL": "postgres://localhost/test",
                "API_BASE_URL": "http://localhost:8080",
            },
            wantOutput: "Synced 0 resources\n",
        },
        {
            name:    "missing database url",
            args:    []string{"sync"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            for k, v := range tt.env {
                t.Setenv(k, v)
            }

            buf := new(bytes.Buffer)
            cmd := NewRootCmd()
            cmd.SetOut(buf)
            cmd.SetErr(buf)
            cmd.SetArgs(tt.args)

            err := cmd.Execute()
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.wantOutput, buf.String())
        })
    }
}
```
