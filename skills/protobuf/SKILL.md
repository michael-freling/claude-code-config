---
name: protobuf
allowed-tools: Read, Write, Edit, Grep, Glob, Bash
description: Protocol Buffers (protobuf) development following Google's official style guide and best practices for message design, service definitions, versioning, field numbering, and backward compatibility. Emphasizes simplicity, DRY principles, fail-fast error handling, and consistent code patterns. Use when adding or updating .proto files.
---

# Protocol Buffers Development Skill

A comprehensive skill for adding or updating Protocol Buffers definitions following Google's official style guide and industry best practices.

## When to Use

Use this skill when:
- Adding new .proto files or messages
- Updating existing Protocol Buffer definitions
- Defining gRPC services
- Need to ensure backward compatibility
- Working with protobuf versioning

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/docs/guideline.md`**

Before starting any implementation:

1. **Look for `.claude/docs/guideline.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/docs/guideline.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no guideline.md exists:**
   - Consider running `/document-guideline` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from guideline.md:**
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

## Core Principles

1. **Simplicity First**: Follow DRY (Don't Repeat Yourself) - reuse existing messages and patterns rather than creating new ones
2. **Backward Compatibility**: Design for evolution - clients and servers are never updated simultaneously
3. **Fail-Fast**: Use validation and type safety to catch errors early rather than silently failing
4. **Consistency**: Follow Google's official style guide and project-specific conventions
5. **Latest Versions**: Use proto3 syntax and well-known types from the latest Protocol Buffers release
6. **Encapsulation**: Keep message definitions focused and avoid exposing internal implementation details
7. **Comments**: Comments MUST BE about WHY not WHAT - explain the reasoning behind design decisions, not what the fields represent (field names should be self-explanatory)

## File Structure and Formatting

Following Google's official style guide:

**File Naming**: Use `lower_snake_case.proto`

**File Structure Order**:
1. License header (if applicable)
2. File overview comment
3. Syntax declaration (`syntax = "proto3";`)
4. Package statement
5. Imports (sorted alphabetically)
6. File options (go_package, java_package, etc.)
7. Messages, enums, and services

**Formatting**:
- Maximum line length: 80 characters
- Indentation: 2 spaces
- String quotes: Double quotes preferred
- Comments: Use `//` for all comments, add documentation for all public elements

Example of proper file structure:
```protobuf
// Copyright 2025 Example Inc.
// Licensed under Apache 2.0

// User service API definitions for managing user accounts and profiles.
syntax = "proto3";

package user.v1;

import "google/protobuf/timestamp.proto";

option go_package = "github.com/example/api/user/v1;userv1";
option java_package = "com.example.api.user.v1";
option java_outer_classname = "UserProto";

// User represents a user account in the system.
message User {
  // ... message definition
}
```

## Naming Conventions (Google Style Guide)

| Component | Style | Example |
|-----------|-------|---------|
| Files | `lower_snake_case.proto` | `user_service.proto` |
| Packages | dot-delimited `lower_snake_case` | `user.v1`, `my.package.name` |
| Messages | `TitleCase` | `UserProfile`, `CreateUserRequest` |
| Fields | `lower_snake_case`, plural for repeated | `user_name`, `repeated email_addresses` |
| Oneofs | `lower_snake_case` | `auth_method`, `notification_type` |
| Enums | `TitleCase` | `UserRole`, `OrderStatus` |
| Enum values | `UPPER_SNAKE_CASE` with prefix | `USER_ROLE_UNSPECIFIED`, `ORDER_STATUS_PENDING` |
| Services | `TitleCase` with `Service` suffix | `UserService`, `OrderService` |
| RPC methods | `TitleCase`, action-oriented | `GetUser`, `CreateOrder`, `ListProducts` |

**Important Underscore Rules**:
- NEVER use underscores as initial or final characters
- Underscores must be followed by letters (not numbers or additional underscores)
- Treat abbreviations as single words: `GetDnsRequest` (not `GetDNSRequest`), `dns_request` (not `d_n_s_request`)

## Protocol Buffers Best Practices

### Message Design

**Core Guidelines**:
- Use singular names for non-repeated fields
- Use plural names for repeated fields
- Always reserve field numbers for deleted fields
- Keep messages focused and composable (avoid hundreds of fields)
- Add comprehensive comments to all messages and fields
- Define one message per concern - split large messages into smaller, composable ones
- Use well-known types (`google.protobuf.Timestamp`, `Duration`, `FieldMask`, etc.) instead of custom implementations
- Avoid boolean fields for concepts that might expand to multiple states - use enums instead

**Field Encapsulation**:
- Don't expose internal implementation details in public APIs
- Use separate messages for API and storage to enable independent evolution

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

### Versioning and Evolution (CRITICAL)

**NEVER Do These (Breaking Changes)**:
- ❌ NEVER re-use a field number (even from deleted fields) - breaks deserialization completely
- ❌ NEVER change field types (with rare exceptions: int32/uint32/int64/bool are compatible)
- ❌ NEVER add required fields - use optional with `// required` comments instead
- ❌ NEVER change default values - causes inconsistencies across versions
- ❌ NEVER convert between repeated and scalar fields - loses data
- ❌ NEVER re-use enum value numbers
- ❌ NEVER change numeric values of existing enum entries

**Always Do (Safe Changes)**:
- ✅ ALWAYS reserve field numbers AND names when deleting fields
- ✅ ALWAYS reserve enum value numbers when deprecating them
- ✅ ALWAYS mark deprecated fields with `[deprecated = true]` before removing
- ✅ ALWAYS version your packages (v1, v2, etc.) for major API changes
- ✅ ALWAYS add new enum values at the end (after deprecated ones)
- ✅ ALWAYS include backward-compatible changes only
- ✅ Fields can be renamed safely (names don't appear in binary serialization)
- ✅ New fields can be added (ensure clients handle unknown fields gracefully)

**Versioning Strategy**:
- Use package versioning (e.g., `package user.v1`, `package user.v2`) for major changes
- Keep old versions available during migration periods
- Update server schemas first, then gradually update clients
- Design for forward and backward compatibility - clients and servers are NEVER updated simultaneously

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

### Enum Design

**Critical Enum Rules**:
- ALWAYS include a zero-value default with `_UNSPECIFIED` or `_UNKNOWN` suffix
- Prefix ALL enum values with the enum name in `UPPER_SNAKE_CASE` (prevents collisions)
- Zero value MUST be first in the declaration
- Add new values at the END (after any deprecated values)
- NEVER change numeric values of existing entries
- Reserve numbers when deprecating enum values
- Avoid C/C++ reserved keywords (NULL, NAN, etc.)
- Consider nesting enums inside messages for better scoping

Example:
```protobuf
// Good: Properly defined enum with zero value and prefixes
enum OrderStatus {
  // Zero value - required for proto3 default
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING = 1;
  ORDER_STATUS_CONFIRMED = 2;
  ORDER_STATUS_SHIPPED = 3;
  ORDER_STATUS_DELIVERED = 4;

  // Deprecated value - kept for compatibility
  ORDER_STATUS_PROCESSING = 5 [deprecated = true];

  // Reserved for removed values
  reserved 6, 7;
  reserved "ORDER_STATUS_CANCELLED_OLD";
}

// Bad: Missing zero value, no prefixes
enum Status {
  ACTIVE = 1;  // ❌ No zero value!
  INACTIVE = 2;  // ❌ No prefix - risk of collisions
}

// Good: Nested enum for better scoping
message Order {
  enum Status {
    STATUS_UNSPECIFIED = 0;
    STATUS_PENDING = 1;
    STATUS_CONFIRMED = 2;
  }

  Status status = 1;
}
```

### Field Numbering Best Practices

**Field Number Rules**:
- Use 1-15 for frequently used fields (most efficient encoding - 1 byte)
- Use 16-2047 for less frequent fields (2 bytes)
- Group related fields logically (e.g., 1-10 core, 11-20 metadata, 21-30 audit)
- Leave gaps between groups for future expansion
- NEVER re-use deleted field numbers - always reserve them
- Field numbers 19000-19999 are reserved by Protocol Buffers

**Efficient Field Numbering Strategy**:
```
1-15:    Core/frequently accessed fields (1-byte encoding)
16-50:   Common optional fields
51-100:  Extended fields
100+:    Rarely used or future expansion
```

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
- Read `.claude/docs/guideline.md` if it exists (MANDATORY)
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

## Serialization and Interchange

**Best Practices**:
- ✅ Use binary serialization for data interchange (most efficient and stable)
- ❌ NEVER use text format (JSON, text proto) for long-term storage or interchange
- ❌ NEVER rely on serialization stability across builds for cache keys
- ✅ Use different messages for API contracts vs. storage schemas
- ✅ Separate client-facing messages from internal storage to enable independent evolution

**Why Avoid Text Formats for Interchange**:
- Text formats break when fields or enums are renamed
- Binary format is stable and handles unknown fields gracefully
- Significantly better performance and smaller size

## Language-Specific Options

**Java**:
```protobuf
option java_package = "com.example.api.user.v1";
option java_outer_classname = "UserProto";  // TitleCase of filename
option java_multiple_files = true;  // Separate class per message
```

**Go**:
```protobuf
option go_package = "github.com/example/api/user/v1;userv1";
```

**Guidelines**:
- Keep generated code in separate packages from hand-written code
- Derive java_package from proto package to avoid collisions
- Use java_outer_classname to convert filename to TitleCase
- Avoid language keywords for field names (protobuf may rename them)

## Commands and Validation

```bash
# Using buf (recommended for linting and breaking change detection)
buf lint                                    # Lint proto files
buf format -w                               # Format proto files
buf generate                                # Generate code
buf breaking --against '.git#branch=main'  # Check for breaking changes

# Using protoc directly
protoc --go_out=. --go-grpc_out=. api/**/*.proto

# Validation workflow (fail-fast approach)
buf lint && buf generate && buf breaking --against '.git#branch=main'
```

**Recommended Tools**:
- Use `buf` for modern protobuf workflows (linting, formatting, breaking change detection)
- Consider `buf.build` for schema registry and dependency management
- Use language-specific linters for generated code

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

### Using Booleans for Multi-State Concepts
```protobuf
// Bad: Boolean that might need more states later
message Feature {
  bool enabled = 1;  // What if we need "partial" or "testing" mode?
}

// Good: Enum allows future expansion
message Feature {
  enum State {
    STATE_UNSPECIFIED = 0;
    STATE_DISABLED = 1;
    STATE_ENABLED = 2;
    // Easy to add: STATE_TESTING = 3, STATE_PARTIAL = 4
  }
  State state = 1;
}
```

### Missing Zero Values
```protobuf
// Bad: No zero value
enum Status {
  ACTIVE = 1;  // ❌ Clients receiving unknown values can't distinguish from unset
  INACTIVE = 2;
}

// Good: Zero value defined
enum Status {
  STATUS_UNSPECIFIED = 0;  // ✅ Clear handling of unset/unknown values
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
}
```

### Too Many Fields in One Message
```protobuf
// Bad: Hundreds of fields in one message
message UserProfile {
  string id = 1;
  string name = 2;
  // ... 200+ more fields
  // This causes: memory bloat, hard to maintain, generated code size limits
}

// Good: Split into focused, composable messages
message User {
  string id = 1;
  string name = 2;
  ContactInfo contact = 3;
  Preferences preferences = 4;
}

message ContactInfo {
  string email = 1;
  string phone = 2;
  Address address = 3;
}

message Preferences {
  string language = 1;
  string timezone = 2;
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
- [ ] **Read `.claude/docs/guideline.md` if it exists** (CRITICAL)
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

## Key Principles Summary

1. **Project Guidelines First**: Always read and follow `.claude/docs/guideline.md` (MANDATORY)
2. **Simplicity & DRY**: Reuse existing messages and patterns; avoid creating redundant definitions
3. **Backward Compatibility**: NEVER reuse field numbers or change types; design for evolution
4. **Fail-Fast**: Use validation, type safety, and strict linting to catch errors early
5. **Google Style Guide**: Follow official naming conventions and file structure
6. **Latest Versions**: Use proto3 syntax and latest well-known types
7. **Zero Values**: Always define UNSPECIFIED (0) as first enum value
8. **Separation**: Use separate Request/Response messages for each RPC; separate API from storage
9. **Well-Known Types**: Use google.protobuf types (Timestamp, Duration, etc.) for common patterns
10. **Field Numbering**: Use 1-15 for frequent fields; group logically; reserve deleted numbers
11. **Documentation**: Comment all public messages, fields, services, and RPCs
12. **Encapsulation**: Keep messages focused; avoid exposing internal implementation details
13. **Binary Serialization**: Use binary format for interchange; avoid text formats for storage

## Version History

### Version 2.0 (2025-01-16)
- Added comprehensive Google Style Guide naming conventions (files, packages, messages, fields, enums, services)
- Added Core Principles section emphasizing simplicity, DRY, fail-fast, and consistency
- Enhanced versioning section with explicit NEVER/ALWAYS rules for backward compatibility
- Added detailed enum design guidelines with zero-value requirements and prefixing rules
- Expanded field numbering with encoding efficiency details (1-15 = 1 byte)
- Added serialization best practices (binary vs text format)
- Added language-specific options for Java and Go
- Added new pitfall examples: booleans for multi-state, too many fields, etc.
- Enhanced file structure section with proper ordering and formatting rules
- Added underscore naming rules from official style guide
- Improved commands section with buf tooling recommendations
- Updated all examples to follow Google's official style guide
- Aligned with coding-guideline.md principles: simplicity, DRY, fail-fast, latest versions, encapsulation

### Version 1.0 (Initial)
- Basic Protocol Buffers guidelines
- Message design patterns
- Service definitions
- Field numbering basics
