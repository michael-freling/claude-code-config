# HTTP Client Patterns

Idiomatic Go HTTP client with context propagation, retries, and inline OpenTelemetry instrumentation.

## Client struct with connection pooling

```go
import (
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type VendorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewVendorClient(baseURL string, opts ...Option) *VendorClient {
	o := options{timeout: 10 * time.Second}
	for _, opt := range opts {
		opt(&o)
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &VendorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   o.timeout,
			Transport: otelhttp.NewTransport(transport),
		},
	}
}
```

## GET request with typed response and OTel span

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (c *VendorClient) GetProduct(ctx context.Context, id string) (*Product, error) {
	ctx, span := otel.Tracer("vendor").Start(ctx, "GetProduct")
	defer span.End()
	span.SetAttributes(attribute.String("product.id", id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/products/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("new request for product %s: %w", id, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("GET product %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET product %s: status %d", id, resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("decode product %s: %w", id, err)
	}
	return &product, nil
}
```

## POST request with JSON body

```go
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type CreateOrderRequest struct {
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}

func (c *VendorClient) CreateOrder(ctx context.Context, order CreateOrderRequest) (*CreateOrderResponse, error) {
	ctx, span := otel.Tracer("vendor").Start(ctx, "CreateOrder")
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
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("POST order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("POST order: status %d", resp.StatusCode)
	}

	var result CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode order response: %w", err)
	}
	return &result, nil
}
```

## Retry with exponential backoff

```go
import (
	"context"
	"math"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func (c *VendorClient) doWithRetry(ctx context.Context, req *http.Request, cfg RetryConfig) (*http.Response, error) {
	var lastErr error
	for attempt := range cfg.MaxAttempts {
		resp, err := c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			resp.Body.Close()
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		}

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		delay := time.Duration(float64(cfg.BaseDelay) * math.Pow(2, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}
```
