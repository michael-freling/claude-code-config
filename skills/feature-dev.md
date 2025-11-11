# Feature Development Skill

A comprehensive skill for adding or updating features in TypeScript, Next.js, Golang, and Protocol Buffers projects following best practices.

## Overview

This skill helps implement features by:
1. Reading project-specific design documentation
2. Analyzing the codebase structure and conventions
3. Identifying relevant source files
4. Following language-specific and project-specific best practices
5. Ensuring consistency with existing code patterns
6. Implementing changes systematically

## When to Use

Use this skill when:
- Adding new features to an existing project
- Updating or enhancing existing functionality
- Working with TypeScript, Next.js, Go, or Protocol Buffers
- Need to follow project-specific patterns and conventions

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation, check if the project has design documentation:

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
   - Document patterns you discover for future reference

**What to extract from design.md:**
- Project-specific naming conventions
- Directory structure and organization rules
- Error handling patterns used in this project
- Testing conventions and structure
- Architecture patterns (layering, dependency rules)
- Type/interface conventions
- API design patterns
- Code examples showing the preferred style

### 1. Analyze Project Structure

After reading design.md (if available), verify and identify the project type and structure:

**TypeScript/Next.js Projects:**
- Check for `package.json`, `tsconfig.json`, `next.config.js`
- Identify directory structure (src/, pages/, app/, components/, etc.)
- Review existing TypeScript conventions (interfaces, types, enums)
- Check for state management (Redux, Zustand, Context API)

**Go Projects:**
- Check for `go.mod` to identify module name
- Identify package structure and naming conventions
- Review existing patterns (interfaces, error handling, context usage)
- Check for common Go project layouts (cmd/, pkg/, internal/)

**Protocol Buffers:**
- Locate `.proto` files
- Identify proto package naming conventions
- Check for existing message patterns and service definitions
- Review import patterns and dependencies

### 2. Search for Relevant Code

Use targeted searches to find related implementations:

```
# For TypeScript/Next.js
- Search for similar components/pages
- Look for related API routes or server actions
- Find existing hooks or utilities
- Identify styling patterns (CSS modules, Tailwind, styled-components)

# For Go
- Search for related packages
- Find similar function implementations
- Look for interface definitions
- Identify error handling patterns

# For Protocol Buffers
- Find similar message types
- Look for related service definitions
- Identify field naming conventions
```

### 3. Follow Best Practices

#### TypeScript Best Practices

**Code Organization:**
- Follow existing directory structure
- Use barrel exports (index.ts) for clean imports
- Separate concerns: components, hooks, utils, types, constants
- Co-locate related files (component.tsx, component.test.tsx, component.module.css)

Example:
```typescript
// src/types/user.ts
export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
}

export type UserRole = 'admin' | 'user' | 'guest';

export interface CreateUserDto {
  name: string;
  email: string;
  role: UserRole;
}

// src/types/index.ts (barrel export)
export * from './user';
export * from './product';
```

**Type Safety:**
- Define interfaces/types for all data structures
- Use strict TypeScript settings (strict: true in tsconfig.json)
- Prefer `interface` for object shapes, `type` for unions/intersections
- Avoid `any`; use `unknown` for truly unknown types
- Use generics for reusable type-safe functions
- Leverage utility types (Partial, Pick, Omit, Record)

Example:
```typescript
// Good: Type-safe API response handler
interface ApiResponse<T> {
  data: T;
  status: number;
  message: string;
}

async function fetchData<T>(url: string): Promise<ApiResponse<T>> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  return response.json();
}

// Usage with type inference
const userResponse = await fetchData<User>('/api/users/1');
// userResponse.data is typed as User

// Good: Use utility types
type PartialUser = Partial<User>;
type UserPreview = Pick<User, 'id' | 'name'>;
type UserUpdate = Omit<User, 'id' | 'createdAt'>;
```

**Function Patterns:**
- Use arrow functions for inline callbacks
- Use function declarations for top-level functions
- Explicit return types for public APIs
- Use async/await over raw promises

Example:
```typescript
// Good: Explicit return types for clarity
export function calculateTotal(items: CartItem[]): number {
  return items.reduce((sum, item) => sum + item.price * item.quantity, 0);
}

// Good: Async/await with proper error handling
export async function getUserById(id: string): Promise<User | null> {
  try {
    const response = await fetchData<User>(`/api/users/${id}`);
    return response.data;
  } catch (error) {
    console.error('Failed to fetch user:', error);
    return null;
  }
}
```

**Error Handling:**
- Create custom error classes for specific error types
- Use type guards for error handling
- Provide meaningful error messages

Example:
```typescript
// Good: Custom error classes
export class ValidationError extends Error {
  constructor(
    message: string,
    public field: string,
    public value: unknown
  ) {
    super(message);
    this.name = 'ValidationError';
  }
}

export class NotFoundError extends Error {
  constructor(resource: string, id: string) {
    super(`${resource} with id ${id} not found`);
    this.name = 'NotFoundError';
  }
}

// Good: Type guard for error handling
function isValidationError(error: unknown): error is ValidationError {
  return error instanceof ValidationError;
}

try {
  await createUser(data);
} catch (error) {
  if (isValidationError(error)) {
    console.error(`Validation failed for ${error.field}`);
  } else if (error instanceof Error) {
    console.error(error.message);
  }
}
```

#### Next.js Best Practices

**App Router (Next.js 13+):**
- Use Server Components by default
- Add 'use client' only when needed (interactivity, browser APIs)
- Leverage Server Actions for mutations
- Use async/await in Server Components

Example:
```typescript
// app/users/page.tsx - Server Component (default)
import { fetchUsers } from '@/lib/api';

export default async function UsersPage() {
  const users = await fetchUsers(); // Direct async call

  return (
    <div>
      <h1>Users</h1>
      <UserList users={users} />
    </div>
  );
}

// app/users/user-list.tsx - Client Component (interactive)
'use client';

import { useState } from 'react';
import type { User } from '@/types';

interface UserListProps {
  users: User[];
}

export function UserList({ users }: UserListProps) {
  const [filter, setFilter] = useState('');

  const filtered = users.filter(u =>
    u.name.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div>
      <input
        type="text"
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
        placeholder="Filter users..."
      />
      {filtered.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
}
```

**Server Actions:**
- Define server actions in separate files or at top of server components
- Use proper form validation
- Return structured data with success/error states

Example:
```typescript
// app/actions/user.ts
'use server';

import { revalidatePath } from 'next/cache';
import { z } from 'zod';

const createUserSchema = z.object({
  name: z.string().min(1).max(100),
  email: z.string().email(),
});

export async function createUser(formData: FormData) {
  const parsed = createUserSchema.safeParse({
    name: formData.get('name'),
    email: formData.get('email'),
  });

  if (!parsed.success) {
    return { success: false, errors: parsed.error.flatten() };
  }

  try {
    await db.user.create({ data: parsed.data });
    revalidatePath('/users');
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to create user' };
  }
}
```

**Pages Router (Legacy):**
- Use getServerSideProps for dynamic data
- Use getStaticProps + ISR for static content
- Implement getStaticPaths for dynamic routes

Example:
```typescript
// pages/users/[id].tsx
import type { GetServerSideProps } from 'next';
import type { User } from '@/types';

interface UserPageProps {
  user: User;
}

export default function UserPage({ user }: UserPageProps) {
  return (
    <div>
      <h1>{user.name}</h1>
      <p>{user.email}</p>
    </div>
  );
}

export const getServerSideProps: GetServerSideProps<UserPageProps> = async (context) => {
  const { id } = context.params!;

  try {
    const user = await fetchUser(id as string);

    if (!user) {
      return { notFound: true };
    }

    return { props: { user } };
  } catch (error) {
    return { notFound: true };
  }
};
```

**Component Patterns:**
- Use functional components with hooks
- Implement proper prop typing with interfaces
- Use composition over inheritance
- Memoize expensive computations and callbacks

Example:
```typescript
// Good: Well-typed reusable component
import { memo } from 'react';

interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
}

export const Button = memo<ButtonProps>(({
  variant = 'primary',
  size = 'md',
  disabled = false,
  onClick,
  children
}) => {
  return (
    <button
      className={`btn btn-${variant} btn-${size}`}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
});

Button.displayName = 'Button';
```

**Optimization:**
- Use next/image for automatic optimization
- Use next/link for client-side navigation
- Implement dynamic imports for code splitting
- Use React.lazy for component lazy loading

Example:
```typescript
import Image from 'next/image';
import Link from 'next/link';
import dynamic from 'next/dynamic';

// Good: Optimized image
<Image
  src="/profile.jpg"
  alt="Profile"
  width={200}
  height={200}
  priority // for LCP images
/>

// Good: Client-side navigation
<Link href="/dashboard">Dashboard</Link>

// Good: Dynamic import with loading state
const HeavyChart = dynamic(() => import('@/components/HeavyChart'), {
  loading: () => <div>Loading chart...</div>,
  ssr: false, // disable SSR if not needed
});
```

#### Go Best Practices

**Code Organization:**
- Follow standard Go project layout
- One package per directory
- Keep packages focused and cohesive
- Use `internal/` for private packages

Example structure:
```
myproject/
├── cmd/
│   └── myapp/
│       └── main.go          # Application entrypoint
├── internal/
│   ├── user/
│   │   ├── user.go          # User domain logic
│   │   ├── user_test.go     # Tests
│   │   └── repository.go    # Data access
│   └── http/
│       └── handler.go       # HTTP handlers
├── pkg/
│   └── logger/
│       └── logger.go        # Reusable public packages
├── go.mod
└── go.sum
```

**Naming Conventions:**
- Use camelCase for unexported, CamelCase for exported
- Interface names: Reader, Writer, Handler (noun or noun+er)
- Keep names concise but meaningful
- Package names: lowercase, no underscores

Example:
```go
// Good: Clear, concise names
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
}

type userRepositoryImpl struct {
    db *sql.DB
}

// Good: Package-level exported function
func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepositoryImpl{db: db}
}
```

**Error Handling:**
- Return errors as the last return value
- Wrap errors with context using fmt.Errorf with %w
- Don't panic in library code
- Handle all errors explicitly
- Define custom error types for domain errors

Example:
```go
// Good: Custom error types
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

// Good: Error wrapping with context
func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id).Scan(
        &user.ID, &user.Name, &user.Email,
    )
    if err == sql.ErrNoRows {
        return nil, &NotFoundError{Resource: "user", ID: id}
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query user %s: %w", id, err)
    }
    return &user, nil
}

// Good: Explicit error handling
user, err := repo.FindByID(ctx, "123")
if err != nil {
    var notFound *NotFoundError
    if errors.As(err, &notFound) {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return nil, fmt.Errorf("unexpected error: %w", err)
}
```

**Interfaces:**
- Keep interfaces small (1-3 methods ideal)
- Define interfaces at point of use (consumer side)
- Accept interfaces, return structs
- Use interface{} sparingly; prefer concrete types or generics (Go 1.18+)

Example:
```go
// Good: Small, focused interface
type UserStore interface {
    Get(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Good: Accept interface, return struct
func NewUserService(store UserStore, logger Logger) *UserService {
    return &UserService{
        store:  store,
        logger: logger,
    }
}

// Good: Define interface where it's used
// In service package:
type Logger interface {
    Info(msg string, fields ...any)
    Error(msg string, fields ...any)
}

// Multiple logger implementations can satisfy this
```

**Struct Design:**
- Use struct tags for JSON, DB, validation
- Group related fields
- Consider embedding for composition
- Use pointers for large structs or when you need nil

Example:
```go
// Good: Well-structured domain model
type User struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name" validate:"required,min=1,max=100"`
    Email     string    `json:"email" db:"email" validate:"required,email"`
    Role      UserRole  `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
    UserRoleAdmin UserRole = "admin"
    UserRoleUser  UserRole = "user"
    UserRoleGuest UserRole = "guest"
)

// Good: Embedding for composition
type AuditLog struct {
    User                    // Embedded
    Action    string       `json:"action"`
    Timestamp time.Time    `json:"timestamp"`
}
```

**Concurrency:**
- Use context.Context for cancellation and timeouts
- Avoid goroutine leaks with proper cleanup
- Use channels or sync primitives appropriately
- Close channels when done sending

Example:
```go
// Good: Context-aware concurrent processing
func (s *UserService) ProcessUsers(ctx context.Context, ids []string) error {
    eg, ctx := errgroup.WithContext(ctx)

    for _, id := range ids {
        id := id // Capture loop variable
        eg.Go(func() error {
            return s.processUser(ctx, id)
        })
    }

    return eg.Wait()
}

// Good: Channel-based worker pool
func (s *UserService) ProcessWithWorkers(ctx context.Context, users []User) error {
    jobs := make(chan User, len(users))
    results := make(chan error, len(users))

    // Start workers
    for w := 0; w < 5; w++ {
        go func() {
            for user := range jobs {
                select {
                case <-ctx.Done():
                    results <- ctx.Err()
                    return
                default:
                    results <- s.processUser(ctx, user.ID)
                }
            }
        }()
    }

    // Send jobs
    for _, user := range users {
        jobs <- user
    }
    close(jobs)

    // Collect results
    for range users {
        if err := <-results; err != nil {
            return err
        }
    }

    return nil
}
```

**Function Design:**
- Use functional options pattern for complex constructors
- Return early to reduce nesting
- Keep functions small and focused

Example:
```go
// Good: Functional options pattern
type ServerOptions struct {
    Port    int
    Timeout time.Duration
    Logger  Logger
}

type ServerOption func(*ServerOptions)

func WithPort(port int) ServerOption {
    return func(o *ServerOptions) {
        o.Port = port
    }
}

func WithTimeout(timeout time.Duration) ServerOption {
    return func(o *ServerOptions) {
        o.Timeout = timeout
    }
}

func NewServer(opts ...ServerOption) *Server {
    options := &ServerOptions{
        Port:    8080,
        Timeout: 30 * time.Second,
    }

    for _, opt := range opts {
        opt(options)
    }

    return &Server{options: options}
}

// Usage:
server := NewServer(
    WithPort(9000),
    WithTimeout(60 * time.Second),
)
```

**Testing:**
- Write table-driven tests
- Use t.Helper() for test helpers
- Test exported APIs, not implementation details
- Use subtests for better organization

Example:
```go
// Good: Table-driven tests
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name: "valid user",
            user: User{
                Name:  "John Doe",
                Email: "john@example.com",
                Role:  UserRoleUser,
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: User{
                Name:  "John Doe",
                Email: "invalid-email",
                Role:  UserRoleUser,
            },
            wantErr: true,
        },
        {
            name: "empty name",
            user: User{
                Name:  "",
                Email: "john@example.com",
                Role:  UserRoleUser,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// Good: Test helper
func createTestUser(t *testing.T, db *sql.DB, name string) *User {
    t.Helper()

    user := &User{
        ID:    uuid.New().String(),
        Name:  name,
        Email: fmt.Sprintf("%s@example.com", name),
    }

    err := db.QueryRow(
        "INSERT INTO users (id, name, email) VALUES (?, ?, ?) RETURNING created_at",
        user.ID, user.Name, user.Email,
    ).Scan(&user.CreatedAt)
    if err != nil {
        t.Fatalf("failed to create test user: %v", err)
    }

    return user
}
```

#### Protocol Buffers Best Practices

**Message Design:**
- Use singular for non-repeated fields
- Use plural for repeated fields
- Reserve field numbers for deleted fields
- Keep messages focused and composable

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

**Field Naming:**
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

**Versioning and Evolution:**
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

**Service Design:**
- Use clear, action-oriented RPC names (GetUser, CreateOrder, ListProducts)
- Design for forward/backward compatibility
- Include proper request/response messages (never reuse)
- Use streaming where appropriate

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

**Common Patterns:**

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

**Best Practices Summary:**
- Always include a zero value for enums (UNSPECIFIED)
- Use consistent naming across services
- Document all fields with comments
- Group field numbers logically (1-10 for core, 10-20 for metadata, etc.)
- Use well-known types from google.protobuf
- Version your APIs (package user.v1, user.v2)
- Never break backward compatibility in the same version

### 4. Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract project-specific patterns and conventions
- Note any architecture constraints or rules
- Identify preferred libraries or utilities

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into small, testable units
- Identify dependencies between tasks
- Ensure plan aligns with project architecture from design.md

**Step 3: Read Existing Code**
- Review similar implementations mentioned in design.md
- Understand existing patterns in actual code
- Note any utilities or helpers available
- Verify design.md documentation matches reality

**Step 4: Implement Changes**
- Follow patterns from design.md first, then general best practices
- Use naming conventions specified in design.md
- Follow directory structure rules from design.md
- Write clean, readable code
- Add appropriate comments for complex logic
- Match code style of existing files

**Step 5: Ensure Quality**
- Run type checking (TypeScript)
- Run linters/formatters specified in design.md
- Build the project
- Run existing tests
- Verify adherence to design.md guidelines

### 5. Language-Specific Commands

**TypeScript/Next.js:**
```bash
# Type checking
npm run type-check || npx tsc --noEmit

# Linting
npm run lint

# Build
npm run build

# Tests
npm test
```

**Go:**
```bash
# Format
go fmt ./...

# Vet
go vet ./...

# Build
go build ./...

# Test
go test ./...

# Mod tidy
go mod tidy
```

**Protocol Buffers:**
```bash
# Generate code (varies by project)
protoc --go_out=. --go-grpc_out=. *.proto
buf generate
```

## Example Workflows

### TypeScript: Adding a New API Route Handler

**Scenario**: Add a user update endpoint

1. **Analyze existing structure**
   - Search for similar API routes: `pages/api/users/**` or `app/api/users/**`
   - Review existing user-related types

2. **Define types**
```typescript
// types/user.ts
export interface UpdateUserDto {
  name?: string;
  email?: string;
  role?: UserRole;
}

export interface UpdateUserResponse {
  user: User;
  message: string;
}
```

3. **Implement handler** (App Router)
```typescript
// app/api/users/[id]/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { z } from 'zod';

const updateUserSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  email: z.string().email().optional(),
  role: z.enum(['admin', 'user', 'guest']).optional(),
});

export async function PATCH(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const body = await request.json();
    const validated = updateUserSchema.parse(body);

    const user = await db.user.update({
      where: { id: params.id },
      data: validated,
    });

    return NextResponse.json({ user, message: 'User updated' });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return NextResponse.json(
        { error: 'Validation failed', details: error.errors },
        { status: 400 }
      );
    }
    return NextResponse.json(
      { error: 'Failed to update user' },
      { status: 500 }
    );
  }
}
```

4. **Verify**
   - Run `npm run type-check`
   - Run `npm run lint`
   - Run `npm run build`

### Next.js: Adding a New Page Component

**Scenario**: Create a user profile page

1. **Search for similar pages**
   - Look for existing detail pages
   - Check routing patterns

2. **Create Server Component** (App Router)
```typescript
// app/users/[id]/page.tsx
import { notFound } from 'next/navigation';
import { getUserById } from '@/lib/api/users';
import { UserProfile } from '@/components/user-profile';

interface UserPageProps {
  params: { id: string };
}

export async function generateMetadata({ params }: UserPageProps) {
  const user = await getUserById(params.id);
  return {
    title: user ? `${user.name} - Profile` : 'User Not Found',
  };
}

export default async function UserPage({ params }: UserPageProps) {
  const user = await getUserById(params.id);

  if (!user) {
    notFound();
  }

  return (
    <div>
      <h1>{user.name}</h1>
      <UserProfile user={user} />
    </div>
  );
}
```

3. **Create Client Component**
```typescript
// components/user-profile.tsx
'use client';

import { useState } from 'react';
import type { User } from '@/types';

interface UserProfileProps {
  user: User;
}

export function UserProfile({ user }: UserProfileProps) {
  const [isEditing, setIsEditing] = useState(false);

  return (
    <div>
      <p>Email: {user.email}</p>
      <p>Role: {user.role}</p>
      <button onClick={() => setIsEditing(!isEditing)}>
        {isEditing ? 'Cancel' : 'Edit'}
      </button>
    </div>
  );
}
```

4. **Test and verify**
   - Navigate to `/users/[id]` in browser
   - Check for console errors
   - Verify type safety

### Go: Adding a New Service Method

**Scenario**: Add user listing with filters

1. **Define interface in service package**
```go
// internal/user/service.go
type ListUsersFilter struct {
    Role   *UserRole
    Limit  int
    Offset int
}

func (s *Service) ListUsers(ctx context.Context, filter ListUsersFilter) ([]User, error) {
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    users, err := s.repo.List(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    return users, nil
}
```

2. **Implement repository method**
```go
// internal/user/repository.go
func (r *repository) List(ctx context.Context, filter ListUsersFilter) ([]User, error) {
    query := "SELECT id, name, email, role, created_at FROM users WHERE 1=1"
    args := []interface{}{}
    argPos := 1

    if filter.Role != nil {
        query += fmt.Sprintf(" AND role = $%d", argPos)
        args = append(args, *filter.Role)
        argPos++
    }

    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
    args = append(args, filter.Limit, filter.Offset)

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }

    return users, rows.Err()
}
```

3. **Add tests**
```go
// internal/user/service_test.go
func TestService_ListUsers(t *testing.T) {
    tests := []struct {
        name    string
        filter  ListUsersFilter
        want    int
        wantErr bool
    }{
        {
            name:   "list all users",
            filter: ListUsersFilter{Limit: 10},
            want:   3,
        },
        {
            name: "filter by role",
            filter: ListUsersFilter{
                Role:  ptr(UserRoleAdmin),
                Limit: 10,
            },
            want: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test database
            db := setupTestDB(t)
            defer db.Close()

            repo := NewRepository(db)
            service := NewService(repo)

            users, err := service.ListUsers(context.Background(), tt.filter)
            if (err != nil) != tt.wantErr {
                t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(users) != tt.want {
                t.Errorf("ListUsers() got %d users, want %d", len(users), tt.want)
            }
        })
    }
}
```

4. **Verify**
   - Run `go fmt ./...`
   - Run `go vet ./...`
   - Run `go test ./...`
   - Run `go build ./...`

### Protocol Buffers: Adding a New Service

**Scenario**: Add a notification service

1. **Find similar services**
   - Review existing service patterns
   - Check field numbering conventions

2. **Define messages and service**
```protobuf
// api/notification/v1/notification.proto
syntax = "proto3";

package notification.v1;

option go_package = "github.com/example/api/notification/v1;notificationv1";

import "google/protobuf/timestamp.proto";

service NotificationService {
  rpc SendNotification(SendNotificationRequest) returns (SendNotificationResponse) {}
  rpc ListNotifications(ListNotificationsRequest) returns (ListNotificationsResponse) {}
  rpc MarkAsRead(MarkAsReadRequest) returns (MarkAsReadResponse) {}
}

message Notification {
  string id = 1;
  string user_id = 2;
  string title = 3;
  string message = 4;
  NotificationType type = 5;
  bool is_read = 6;
  google.protobuf.Timestamp created_at = 7;
}

enum NotificationType {
  NOTIFICATION_TYPE_UNSPECIFIED = 0;
  NOTIFICATION_TYPE_INFO = 1;
  NOTIFICATION_TYPE_WARNING = 2;
  NOTIFICATION_TYPE_ERROR = 3;
}

message SendNotificationRequest {
  string user_id = 1;
  string title = 2;
  string message = 3;
  NotificationType type = 4;
}

message SendNotificationResponse {
  Notification notification = 1;
}

message ListNotificationsRequest {
  string user_id = 1;
  int32 page_size = 2;
  string page_token = 3;
}

message ListNotificationsResponse {
  repeated Notification notifications = 1;
  string next_page_token = 2;
}

message MarkAsReadRequest {
  string id = 1;
}

message MarkAsReadResponse {
  Notification notification = 1;
}
```

3. **Generate code**
   - Run `buf generate` or `protoc --go_out=. --go-grpc_out=. api/notification/v1/*.proto`

4. **Implement service**
```go
// internal/notification/service.go
type Service struct {
    notificationv1.UnimplementedNotificationServiceServer
    repo Repository
}

func (s *Service) SendNotification(
    ctx context.Context,
    req *notificationv1.SendNotificationRequest,
) (*notificationv1.SendNotificationResponse, error) {
    notification := &notificationv1.Notification{
        Id:        uuid.New().String(),
        UserId:    req.UserId,
        Title:     req.Title,
        Message:   req.Message,
        Type:      req.Type,
        IsRead:    false,
        CreatedAt: timestamppb.Now(),
    }

    if err := s.repo.Create(ctx, notification); err != nil {
        return nil, fmt.Errorf("failed to create notification: %w", err)
    }

    return &notificationv1.SendNotificationResponse{
        Notification: notification,
    }, nil
}
```

5. **Verify**
   - Check generated code compiles
   - Implement tests
   - Update API documentation

## Documentation

After implementing features:
- Update relevant documentation in .claude/docs/
- Document new APIs or interfaces
- Add usage examples where helpful
- Update architecture docs if structure changed

## Common Pitfalls to Avoid

### TypeScript/Next.js

**Type Safety Issues:**
```typescript
// Bad: Using 'any'
function processData(data: any) {
  return data.name; // No type safety
}

// Good: Proper typing
interface Data {
  name: string;
}
function processData(data: Data): string {
  return data.name;
}
```

**Component Mistakes:**
```typescript
// Bad: Client component doing data fetching
'use client';
export default function Page() {
  const [data, setData] = useState(null);
  useEffect(() => {
    fetch('/api/data').then(r => r.json()).then(setData);
  }, []);
  // ...
}

// Good: Server component with direct data fetching
export default async function Page() {
  const data = await fetchData();
  return <ClientComponent data={data} />;
}
```

**Error Handling:**
```typescript
// Bad: Swallowing errors
try {
  await saveUser(user);
} catch (e) {
  console.log('error');
}

// Good: Proper error handling
try {
  await saveUser(user);
} catch (error) {
  if (error instanceof ValidationError) {
    return { error: error.message };
  }
  throw error;
}
```

### Go

**Error Handling:**
```go
// Bad: Ignoring errors
data, _ := readFile(path)

// Good: Handling errors
data, err := readFile(path)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

**Goroutine Leaks:**
```go
// Bad: Goroutine leak
func process(data []string) {
    for _, item := range data {
        go processItem(item) // No way to stop these
    }
}

// Good: Controlled goroutines
func process(ctx context.Context, data []string) error {
    eg, ctx := errgroup.WithContext(ctx)
    for _, item := range data {
        item := item
        eg.Go(func() error {
            return processItem(ctx, item)
        })
    }
    return eg.Wait()
}
```

**Interface Design:**
```go
// Bad: Large interface
type UserService interface {
    Create(User) error
    Update(User) error
    Delete(string) error
    FindByID(string) (*User, error)
    FindByEmail(string) (*User, error)
    List() ([]User, error)
}

// Good: Small, focused interfaces
type UserCreator interface {
    Create(ctx context.Context, user User) error
}

type UserFinder interface {
    FindByID(ctx context.Context, id string) (*User, error)
}
```

### Protocol Buffers

**Field Number Reuse:**
```protobuf
// Bad: Reusing field numbers
message User {
  string id = 1;
  string name = 2;
  string email = 3; // Changed from 'old_field'
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

**Breaking Changes:**
```protobuf
// Bad: Changing field type
message User {
  string id = 1;
  int32 age = 2; // Was string before - breaks compatibility!
}

// Good: Add new field
message User {
  string id = 1;
  string age_string = 2 [deprecated = true];
  int32 age = 3; // New field with new number
}
```

**Missing Zero Values:**
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

## Quick Reference Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL - do this first!)
- [ ] Extract project-specific conventions from design.md
- [ ] Analyze project structure and identify language/framework
- [ ] Search for similar implementations in the codebase
- [ ] Review existing patterns and conventions (verify against design.md)
- [ ] Check for utility functions or helpers already available
- [ ] Note any architecture constraints from design.md

### During Implementation

**TypeScript/Next.js:**
- [ ] Define all types/interfaces first
- [ ] Use strict TypeScript (no `any`)
- [ ] Choose Server vs Client Components appropriately
- [ ] Implement proper error boundaries
- [ ] Follow existing naming conventions
- [ ] Use Next.js optimization features (Image, Link, etc.)

**Go:**
- [ ] Define interfaces at point of use
- [ ] Include context.Context in all long-running operations
- [ ] Return errors as last value
- [ ] Wrap errors with context using %w
- [ ] Write table-driven tests
- [ ] Keep functions small and focused

**Protocol Buffers:**
- [ ] Use next available field numbers
- [ ] Add comments for all messages and fields
- [ ] Include zero value for enums (UNSPECIFIED)
- [ ] Define separate Request/Response messages
- [ ] Use well-known types (Timestamp, Duration, etc.)
- [ ] Version your package (v1, v2, etc.)

### After Implementation
- [ ] Run language-specific formatters and linters
- [ ] Run type checking (TypeScript)
- [ ] Build the project successfully
- [ ] Run existing tests
- [ ] Add new tests for your changes
- [ ] Update documentation
- [ ] Verify no console errors or warnings

## Language-Specific Commands Reference

### TypeScript/Next.js
```bash
# Type checking
npm run type-check
npx tsc --noEmit

# Linting and formatting
npm run lint
npm run lint -- --fix
npx prettier --write .

# Build
npm run build

# Tests
npm test
npm run test:watch

# Development
npm run dev
```

### Go
```bash
# Format code
go fmt ./...
gofmt -s -w .

# Vet code
go vet ./...

# Lint (if golangci-lint installed)
golangci-lint run

# Build
go build ./...
go build -o bin/app ./cmd/app

# Test
go test ./...
go test -v ./...
go test -cover ./...
go test -race ./...

# Mod management
go mod tidy
go mod download
go mod verify
```

### Protocol Buffers
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

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md` before anything else
2. **Consistency**: Match existing code style and patterns in the project
3. **Simplicity**: Write clear, maintainable code that others can understand
4. **Safety**: Leverage type systems and error handling properly
5. **Testing**: Ensure changes work and don't break existing functionality
6. **Documentation**: Keep documentation current with all changes

## When Using This Skill

1. **Read design.md first**: Check for `.claude/design.md` and read it thoroughly (MANDATORY)
2. **Analyze the codebase**: Use grep/search to find similar patterns and verify design.md
3. **Plan the implementation**: Break down into small, testable units aligned with architecture
4. **Follow project patterns**: Use conventions from design.md, then general best practices
5. **Verify continuously**: Run checks specified in design.md after each significant change
6. **Update documentation**: Keep design.md and other docs current with changes

## Relationship with `/document-design` Command

If `.claude/design.md` doesn't exist or is outdated:

1. Run `/document-design` to create or update it
2. Review the generated design.md
3. Then use this skill to implement features following those guidelines

This creates a workflow:
```
/document-design → Review design.md → Use feature-dev skill → Implement features
```
