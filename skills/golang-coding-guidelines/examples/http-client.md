# HTTP Client Patterns

## Client with Functional Options

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    tracer     trace.Tracer
}

type Option func(*clientConfig)

type clientConfig struct {
    timeout    time.Duration
    transport  http.RoundTripper
    maxRetries int
}

func WithTimeout(d time.Duration) Option {
    return func(c *clientConfig) { c.timeout = d }
}

func WithMaxRetries(n int) Option {
    return func(c *clientConfig) { c.maxRetries = n }
}

func NewClient(baseURL string, opts ...Option) *Client {
    cfg := clientConfig{
        timeout:    30 * time.Second,
        maxRetries: 3,
    }
    for _, opt := range opts {
        opt(&cfg)
    }

    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout:   cfg.timeout,
            Transport: cfg.transport,
        },
        tracer: otel.Tracer("client"),
    }
}
```

## Request with Tracing and Error Handling

```go
func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
    ctx, span := c.tracer.Start(ctx, "GetProduct")
    defer span.End()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/products/"+id, nil)
    if err != nil {
        return nil, fmt.Errorf("new request for product %s: %w", id, err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        span.RecordError(err)
        return nil, fmt.Errorf("GET product %s: %w", id, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("GET product %s: status %d, body: %s", id, resp.StatusCode, body)
    }

    var product Product
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, fmt.Errorf("decode product %s: %w", id, err)
    }
    return &product, nil
}
```

## POST with Request Body

```go
func (c *Client) CreateOrder(ctx context.Context, order CreateOrderRequest) (*Order, error) {
    ctx, span := c.tracer.Start(ctx, "CreateOrder")
    defer span.End()

    body, err := json.Marshal(order)
    if err != nil {
        return nil, fmt.Errorf("marshal order: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/orders", bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("new request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        span.RecordError(err)
        return nil, fmt.Errorf("POST order: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("POST order: status %d, body: %s", resp.StatusCode, body)
    }

    var result Order
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode order response: %w", err)
    }
    return &result, nil
}
```

## Retry with Exponential Backoff

```go
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
    var lastErr error
    for attempt := range c.maxRetries {
        resp, err := c.httpClient.Do(req)
        if err != nil {
            lastErr = err
            if !isRetryable(err) {
                return nil, err
            }
            backoff(ctx, attempt)
            continue
        }
        if !isRetryableStatus(resp.StatusCode) {
            return resp, nil
        }
        resp.Body.Close()
        lastErr = fmt.Errorf("status %d", resp.StatusCode)
        backoff(ctx, attempt)
    }
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryableStatus(code int) bool {
    return code == http.StatusTooManyRequests ||
        code == http.StatusBadGateway ||
        code == http.StatusServiceUnavailable ||
        code == http.StatusGatewayTimeout
}

func backoff(ctx context.Context, attempt int) {
    delay := time.Duration(1<<attempt) * 100 * time.Millisecond
    timer := time.NewTimer(delay)
    defer timer.Stop()
    select {
    case <-timer.C:
    case <-ctx.Done():
    }
}
```

## Concurrent Fan-Out Requests

```go
func (c *Client) GetProducts(ctx context.Context, ids []string) ([]*Product, error) {
    ctx, span := c.tracer.Start(ctx, "GetProducts")
    defer span.End()

    products := make([]*Product, len(ids))
    eg, ctx := errgroup.WithContext(ctx)

    for i, id := range ids {
        eg.Go(func() error {
            product, err := c.GetProduct(ctx, id)
            if err != nil {
                return fmt.Errorf("get product %s: %w", id, err)
            }
            products[i] = product
            return nil
        })
    }

    if err := eg.Wait(); err != nil {
        return nil, err
    }
    return products, nil
}
```
