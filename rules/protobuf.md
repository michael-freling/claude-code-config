---
paths:
  - "**/*.proto"
---

# Protocol Buffers Guidelines

## Linting

Follow the project's proto linting rules (e.g., buf lint, protolint). Naming conventions like snake_case fields, `Service` suffix, `Request`/`Response` naming, and `_UNSPECIFIED` enum zero values are typically enforced by linters.

## Service Patterns

Choose based on use case:

| Pattern | When to Use | Service Naming | RPC Naming |
|---------|-------------|----------------|------------|
| **BFF** | Frontend-specific APIs, aggregated views | `{Feature}Service` | Action-oriented (e.g., `LoadDashboard`, `SubmitOrder`) |
| **CRUD** | Generic resource APIs, admin interfaces | `{Resource}Service` | `{Verb}{Resource}` (e.g., `GetUser`, `CreateUser`) |

## Message Design

### BFF Pattern

```protobuf
// Tailored for specific frontend needs
service CheckoutService {
  rpc LoadCheckoutPage(LoadCheckoutPageRequest) returns (LoadCheckoutPageResponse);
  rpc SubmitOrder(SubmitOrderRequest) returns (SubmitOrderResponse);
}

message LoadCheckoutPageResponse {
  Cart cart = 1;
  repeated PaymentMethod payment_methods = 2;
  repeated ShippingOption shipping_options = 3;
}
```

### CRUD Pattern

```protobuf
// Generic resource operations
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

## Composition Pattern

Separate widely-used fields from context-specific fields:

```protobuf
message Resource {
  uint64 id = 1;
  ResourceDetails details = 2;  // Only needed in detail views
}

message ResourceDetails {
  string title = 1;
  string description = 2;
}
```

## Validation (protovalidate)

If the project uses [protovalidate](https://github.com/bufbuild/protovalidate), add validation rules to request fields:

```protobuf
import "buf/validate/validate.proto";

message CreateRequest {
  string name = 1 [(buf.validate.field).string = {min_len: 1, max_len: 200}];
  uint64 id = 2 [(buf.validate.field).uint64 = {gt: 0}];
  MyStatus status = 3 [(buf.validate.field).enum = {not_in: [0]}];
  repeated string tags = 4 [(buf.validate.field).repeated = {max_items: 10}];
}
```

## Best Practices

- Reserve field numbers when removing fields
- Use semantic versioning (`service.v1`)
