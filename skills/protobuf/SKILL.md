---
name: protobuf
description: Protocol Buffers (protobuf) development following best practices for message design, service definitions, versioning, field numbering, and backward compatibility. Use when adding or updating .proto files.
---

# Protocol Buffers Development Skill

A comprehensive skill for adding or updating Protocol Buffers definitions following best practices.

## When to Use

Use this skill when:
- Adding new .proto files or messages
- Updating existing Protocol Buffer definitions
- Defining gRPC services
- Need to ensure backward compatibility
- Working with protobuf versioning

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation:

1. **Look for `.claude/design.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/design.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no design.md exists:**
   - Consider running `/document-design` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from design.md:**
- Proto package naming conventions
- Field numbering conventions
- Service design patterns
- Versioning strategy
- Import patterns
- Code examples showing preferred style

### 1. Analyze Project Structure

- Locate `.proto` files and directory structure
- Identify proto package naming conventions
- Check for existing message patterns
- Review service definitions
- Check import patterns and dependencies
- Review buf.yaml or protoc configuration

### 2. Search for Relevant Code

- Find similar message types
- Look for related service definitions
- Identify field naming conventions
- Check versioning patterns (v1, v2, etc.)

## Protocol Buffers Best Practices

### Message Design

- Use singular for non-repeated fields
- Use plural for repeated fields
- Reserve field numbers for deleted fields
- Keep messages focused and composable
- Always add comments

Example:
```protobuf
// Good: Well-structured message with clear field types
syntax = "proto3";

package user.v1;

option go_package = "github.com/example/api/user/v1;userv1";

import "google/protobuf/timestamp.proto";

// User represents a user in the system
message User {
  // Unique identifier for the user
  string id = 1;

  // User's full name
  string name = 2;

  // User's email address
  string email = 3;

  // User's role in the system
  UserRole role = 4;

  // When the user was created
  google.protobuf.Timestamp created_at = 5;

  // When the user was last updated
  google.protobuf.Timestamp updated_at = 6;

  // Reserved for fields that were removed
  reserved 7, 8;
  reserved "old_field", "deprecated_field";
}

enum UserRole {
  USER_ROLE_UNSPECIFIED = 0;  // Always have a zero value
  USER_ROLE_ADMIN = 1;
  USER_ROLE_USER = 2;
  USER_ROLE_GUEST = 3;
}
```

### Field Naming

- Use snake_case for field names
- Be descriptive but concise
- Group related fields with comments
- Use well-known types (google.protobuf.Timestamp, Duration, etc.)

Example:
```protobuf
message CreateUserRequest {
  // User information
  string name = 1;
  string email = 2;
  UserRole role = 3;

  // Metadata
  string client_id = 10;
  string request_id = 11;
}

message CreateUserResponse {
  // Created user
  User user = 1;

  // Operation metadata
  google.protobuf.Timestamp created_at = 2;
}
```

### Versioning and Evolution

- Never change field numbers
- Never change field types
- Mark deprecated fields with [deprecated = true]
- Use reserved for removed fields
- Version your packages (v1, v2, etc.)

Example:
```protobuf
message LegacyUser {
  string id = 1;
  string name = 2;

  // Deprecated: Use role field instead
  string user_type = 3 [deprecated = true];

  // New field added in v2
  UserRole role = 4;

  // Field 5 was removed, now reserved
  reserved 5;
  reserved "old_password_field";
}
```

### Service Design

- Use clear, action-oriented RPC names (GetUser, CreateOrder, ListProducts)
- Design for forward/backward compatibility
- Include proper request/response messages (never reuse)
- Use streaming where appropriate
- Follow REST-style naming conventions

Example:
```protobuf
service UserService {
  // Get a single user by ID
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}

  // List users with pagination
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {}

  // Create a new user
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}

  // Update an existing user
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse) {}

  // Delete a user
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse) {}

  // Stream user updates
  rpc WatchUsers(WatchUsersRequest) returns (stream User) {}
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  // Pagination
  int32 page_size = 1;
  string page_token = 2;

  // Filtering
  string name_filter = 3;
  UserRole role_filter = 4;
}

message ListUsersResponse {
  repeated User users = 1;
  string next_page_token = 2;
  int32 total_count = 3;
}
```

### Common Patterns

**Pagination:**
```protobuf
message ListRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message ListResponse {
  repeated Item items = 1;
  string next_page_token = 2;
}
```

**Error Details:**
```protobuf
import "google/rpc/status.proto";
import "google/rpc/error_details.proto";

// Use google.rpc.Status for rich error details
message OperationResponse {
  oneof result {
    SuccessData success = 1;
    google.rpc.Status error = 2;
  }
}
```

**Resource Names:**
```protobuf
// Use consistent resource naming
message GetUserRequest {
  // Resource name format: users/{user_id}
  string name = 1;
}

message User {
  // Resource name: users/{user_id}
  string name = 1;

  // Display name for UI
  string display_name = 2;
}
```

**Oneof for Variants:**
```protobuf
message Notification {
  string id = 1;
  google.protobuf.Timestamp created_at = 2;

  oneof notification_type {
    EmailNotification email = 10;
    SmsNotification sms = 11;
    PushNotification push = 12;
  }
}

message EmailNotification {
  string to = 1;
  string subject = 2;
  string body = 3;
}
```

### Field Numbering Best Practices

- Use 1-15 for frequently used fields (more efficient encoding)
- Group related fields (1-10 for core, 10-20 for metadata, etc.)
- Leave gaps for future expansion
- Reserve deleted field numbers

Example:
```protobuf
message User {
  // Core fields: 1-10
  string id = 1;
  string name = 2;
  string email = 3;
  UserRole role = 4;

  // Timestamps: 10-15
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;

  // Metadata: 20-30
  string created_by = 20;
  string updated_by = 21;

  // Reserved for future use or deleted fields
  reserved 5, 6, 7;
}
```

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract proto-specific conventions
- Note versioning strategy
- Check field numbering patterns

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Identify affected services/messages
- Plan field numbers carefully
- Consider backward compatibility

**Step 3: Read Existing Code**
- Review similar message definitions
- Understand existing patterns
- Note package organization

**Step 4: Implement Changes**
- Follow patterns from design.md
- Use consistent naming conventions
- Add comprehensive comments
- Reserve deleted fields properly

**Step 5: Ensure Quality**
- Run buf lint (if using buf)
- Generate code and verify compilation
- Check for breaking changes
- Update documentation

## Commands

```bash
# Using buf (recommended)
buf lint
buf format -w
buf generate
buf breaking --against '.git#branch=main'

# Using protoc directly
protoc --go_out=. --go-grpc_out=. api/**/*.proto

# Validation
buf lint
```

## Common Pitfalls to Avoid

### Field Number Reuse
```protobuf
// Bad: Reusing field numbers
message User {
  string id = 1;
  string name = 2;
  string email = 3; // Changed from 'old_field' - BREAKS COMPATIBILITY!
}

// Good: Reserve deleted fields
message User {
  string id = 1;
  string name = 2;
  reserved 3;
  reserved "old_field";
  string email = 4;
}
```

### Breaking Changes
```protobuf
// Bad: Changing field type
message User {
  string id = 1;
  int32 age = 2; // Was string before - BREAKS COMPATIBILITY!
}

// Good: Add new field
message User {
  string id = 1;
  string age_string = 2 [deprecated = true];
  int32 age = 3; // New field with new number
}
```

### Missing Zero Values
```protobuf
// Bad: No zero value
enum Status {
  ACTIVE = 1;
  INACTIVE = 2;
}

// Good: Zero value defined
enum Status {
  STATUS_UNSPECIFIED = 0;
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
}
```

### Reusing Request/Response Messages
```protobuf
// Bad: Reusing messages
message UserRequest {
  string id = 1;
  string name = 2;
}

service UserService {
  rpc GetUser(UserRequest) returns (User) {}
  rpc CreateUser(UserRequest) returns (User) {}
}

// Good: Separate messages for each RPC
message GetUserRequest {
  string id = 1;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
}
```

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Extract proto-specific conventions
- [ ] Review existing .proto files
- [ ] Identify versioning strategy
- [ ] Check field numbering patterns

### During Implementation
- [ ] Use next available field numbers
- [ ] Add comments for all messages and fields
- [ ] Include zero value for enums (UNSPECIFIED)
- [ ] Define separate Request/Response messages
- [ ] Use well-known types (Timestamp, Duration, etc.)
- [ ] Version your package (v1, v2, etc.)
- [ ] Reserve deleted field numbers
- [ ] Group related fields logically

### After Implementation
- [ ] Run `buf lint` or protoc validation
- [ ] Run `buf generate` or protoc to generate code
- [ ] Check for breaking changes (`buf breaking`)
- [ ] Verify generated code compiles
- [ ] Update API documentation
- [ ] Test with client/server implementations

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Backward Compatibility**: Never reuse field numbers or change types
3. **Documentation**: Comment all messages, fields, and services
4. **Versioning**: Use package versioning (v1, v2) for major changes
5. **Zero Values**: Always define UNSPECIFIED (0) for enums
6. **Separation**: Use separate Request/Response messages for each RPC
7. **Well-Known Types**: Use google.protobuf types for common patterns
8. **Field Numbering**: Group related fields and reserve deleted numbers
