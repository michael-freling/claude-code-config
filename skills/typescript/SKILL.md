---
name: typescript
description: TypeScript feature development following best practices for type safety, code organization, error handling, and testing. Use when adding or updating TypeScript code in projects.
---

# TypeScript Development Skill

A comprehensive skill for adding or updating features in TypeScript projects following modern best practices, emphasizing type safety, simplicity, consistency, and robust testing.

## When to Use

Use this skill when:
- Adding new features to TypeScript projects
- Updating or enhancing existing TypeScript functionality
- Need to follow TypeScript-specific patterns and conventions
- Working with type-safe code

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
- Project-specific naming conventions
- Directory structure and organization rules
- Error handling patterns
- Testing conventions
- Architecture patterns
- Type/interface conventions
- Code examples showing preferred style

### 1. Analyze Project Structure

- Check for `package.json`, `tsconfig.json`
- Identify directory structure (src/, lib/, types/, etc.)
- Review existing TypeScript conventions (interfaces, types, enums)
- Check for state management libraries
- Identify testing framework (Jest, Vitest, etc.)

### 2. Search for Relevant Code

- Search for similar implementations
- Find existing types and interfaces
- Look for related utilities and helpers
- Identify styling patterns

## Core Principles

### Simplicity First
- **DRY (Don't Repeat Yourself)**: Always update existing code for reusability rather than creating new functions
- **Prefer Breaking Changes**: Unless explicitly mentioned, prefer simplicity over backward compatibility
- **Early Returns**: Use early returns and continue statements instead of nested conditionals
- **Minimal Public API**: Keep functions and classes private/unexported unless explicitly needed for encapsulation

### Consistency
- Follow project-specific patterns from `.claude/design.md`
- Match existing code style and naming conventions
- Use consistent error handling patterns throughout

### Fail-Fast
- Validate inputs early and throw meaningful errors
- Never silently swallow errors
- Use type guards to ensure type safety at runtime

### Latest Versions
- Prefer latest versions of packages unless compatibility issues arise
- If dependencies conflict, downgrade major version to maintain compatibility

## TypeScript Best Practices

### Package Management
**CRITICAL: Always use pnpm, never npm**

```bash
# Good: Use pnpm
pnpm install
pnpm add <package>
pnpm update

# Bad: Do NOT use npm
npm install
npm install --legacy-peer-deps  # NEVER use this option
```

### Code Organization

**Directory Structure Approaches:**

1. **Feature-Based** (Recommended for medium-large apps):
```
src/
  features/
    auth/
      components/
      hooks/
      types.ts
      auth.service.ts
      auth.test.ts
    user/
      components/
      types.ts
      user.service.ts
      user.test.ts
  shared/
    components/
    utils/
    types/
```

2. **Type-Based** (For smaller projects):
```
src/
  components/
  services/
  types/
  utils/
```

3. **Hybrid Approach** (Balanced):
```
src/
  components/
  features/
    auth/
      types.ts
      auth.service.ts
  types/        # Shared types only
  utils/
```

**File Organization Best Practices:**
- Co-locate related files (file.ts, file.test.ts)
- Use barrel exports (index.ts) sparingly, only for public API
- Keep component-specific types in the same file or nearby
- Use a shared types directory only for truly shared types

Example:
```typescript
// src/features/user/types.ts
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

// src/features/user/index.ts (barrel export for public API only)
export { UserService } from './user.service';
export type { User, UserRole, CreateUserDto } from './types';
```

### Type Safety

**Strict Configuration (tsconfig.json):**
```json
{
  "compilerOptions": {
    "strict": true,
    "strictNullChecks": true,
    "noImplicitAny": true,
    "moduleResolution": "Bundler",  // For modern bundlers (Vite, ESBuild)
    // or "NodeNext" for Node.js projects
  }
}
```

**Type System Best Practices:**
- **Always enable strict mode**: `"strict": true` in tsconfig.json
- **Avoid `any`**: Use `unknown` for truly unknown types (forces type checking before use)
- **Prefer `interface` for object shapes**: More extensible, better for public APIs
- **Prefer `type` for unions/intersections**: More flexible for complex type operations
- **Use const assertions**: `as const` for literal types and immutable arrays
- **Leverage utility types**: Built-in types like Partial, Pick, Omit, Record, Required
- **Advanced types**: Template literal types, mapped types, conditional types

Example:
```typescript
// Good: Avoid 'any', use 'unknown'
function processData(data: unknown): string {
  // Must narrow type before use
  if (typeof data === 'string') {
    return data.toUpperCase();
  }
  if (typeof data === 'object' && data !== null && 'toString' in data) {
    return data.toString();
  }
  throw new Error('Invalid data type');
}

// Good: Type-safe API response handler with generics
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

// Good: Utility types for type transformations
type PartialUser = Partial<User>;              // All properties optional
type UserPreview = Pick<User, 'id' | 'name'>;  // Only specific properties
type UserUpdate = Omit<User, 'id' | 'createdAt'>; // Exclude specific properties
type UserRecord = Record<string, User>;         // Dictionary type

// Good: Template literal types for type-safe strings
type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';
type ApiEndpoint = `/api/${string}`;
type HttpUrl = `http${'s' | ''}://${string}`;

// Good: Mapped types for transformations
type Readonly<T> = {
  readonly [P in keyof T]: T[P];
};

type Nullable<T> = {
  [P in keyof T]: T[P] | null;
};

// Good: Conditional types for dynamic type assignments
type IsString<T> = T extends string ? true : false;
type NonNullable<T> = T extends null | undefined ? never : T;

// Good: Const assertions for immutable literal types
const ROLES = ['admin', 'user', 'guest'] as const;
type Role = typeof ROLES[number]; // 'admin' | 'user' | 'guest'

const CONFIG = {
  apiUrl: 'https://api.example.com',
  timeout: 5000,
} as const;
// CONFIG properties are readonly and have literal types
```

### Function Patterns

- **Arrow functions** for inline callbacks and short functions
- **Function declarations** for top-level exported functions
- **Explicit return types** for all exported functions (public APIs)
- **async/await** over raw promises for better readability
- **Early returns** to reduce nesting (fail-fast principle)

Example:
```typescript
// Good: Explicit return types for clarity
export function calculateTotal(items: CartItem[]): number {
  return items.reduce((sum, item) => sum + item.price * item.quantity, 0);
}

// Good: Early returns instead of nested conditions
export function getUserRole(user: User | null): UserRole {
  if (!user) {
    return 'guest';
  }

  if (user.isAdmin) {
    return 'admin';
  }

  return 'user';
}

// Bad: Nested conditions
export function getUserRole(user: User | null): UserRole {
  if (user) {
    if (user.isAdmin) {
      return 'admin';
    } else {
      return 'user';
    }
  } else {
    return 'guest';
  }
}

// Good: Async/await with proper error handling and fail-fast
export async function getUserById(id: string): Promise<User> {
  if (!id) {
    throw new ValidationError('User ID is required', 'id', id);
  }

  const response = await fetchData<User>(`/api/users/${id}`);
  if (!response.data) {
    throw new NotFoundError('User', id);
  }

  return response.data;
}

// Good: Type inference for simple arrow functions
const double = (n: number) => n * 2;
const greet = (name: string) => `Hello, ${name}`;

// Good: Dynamic imports for code splitting
export async function loadUserModule(): Promise<typeof import('./user')> {
  return import('./user');
}
```

### Error Handling

**Fail-Fast Principle: NEVER silently swallow errors**

- **Create custom error classes** for specific error types
- **Use type guards** for safe error handling
- **Provide meaningful error messages** with context
- **Validate inputs early** and throw errors immediately
- **Never use empty catch blocks** or silent error handling

Example:
```typescript
// Good: Custom error classes with context
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

export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public endpoint: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Good: Type guard for error handling
function isValidationError(error: unknown): error is ValidationError {
  return error instanceof ValidationError;
}

function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

// Good: Proper error handling with type guards
async function createUser(data: CreateUserDto): Promise<User> {
  try {
    return await userService.create(data);
  } catch (error) {
    if (isValidationError(error)) {
      console.error(`Validation failed for ${error.field}: ${error.value}`);
      throw error; // Re-throw, don't swallow
    }
    if (isApiError(error)) {
      console.error(`API error at ${error.endpoint}: ${error.statusCode}`);
      throw error;
    }
    if (error instanceof Error) {
      console.error('Unexpected error:', error.message);
      throw error;
    }
    // Handle unknown errors
    throw new Error('An unknown error occurred');
  }
}

// Bad: Silently swallowing errors
try {
  await saveUser(user);
} catch (e) {
  console.log('error'); // Don't do this!
}

// Bad: Returning null instead of throwing
async function getUser(id: string): Promise<User | null> {
  try {
    return await fetchUser(id);
  } catch (error) {
    return null; // Loses error context
  }
}

// Good: Fail-fast with early validation
export function processOrder(order: Order): ProcessedOrder {
  if (!order) {
    throw new ValidationError('Order is required', 'order', order);
  }
  if (!order.items || order.items.length === 0) {
    throw new ValidationError('Order must have items', 'items', order.items);
  }
  if (order.total <= 0) {
    throw new ValidationError('Order total must be positive', 'total', order.total);
  }

  // Process order only after validation
  return {
    ...order,
    processedAt: new Date(),
  };
}
```

### Testing

**CRITICAL: Always use table-driven tests (test.each) as the primary testing approach.**

Table-driven testing is a strategy that encourages reuse of test logic by defining test cases as array entries with inputs and expected results, then executing them against a single generic test function.

**Benefits:**
- **Reduce code duplication**: Write test logic once, reuse for all cases
- **Improve maintainability**: New scenarios simply append to the cases array
- **Enhance readability**: Immediately clear what inputs/outputs are tested
- **Ensure consistency**: All test cases follow the same structure
- **Make gaps obvious**: Missing test cases are easy to spot
- **Type safety**: TypeScript ensures test case data matches expected types

**Use table-driven tests for:**
- All functions with multiple test cases (2+ cases)
- Error handling scenarios
- Edge cases and boundary conditions
- Different input/output combinations
- Validation logic with various inputs

**When NOT to use table-driven tests:**
- Single test case with complex mocking
- Tests requiring significantly different setup per case
- Integration tests with heavy external dependencies
- Tests where table structure becomes too complex

**Testing Guidelines from Coding Standards:**
1. **Use table-driven testing** - Primary approach for all suitable tests
2. **Avoid redundant tests** - Don't add meaningless test cases with same purpose
3. **Prefer injection over global state** - Use constructor/method injection instead of environment variables or global state
4. **Test inputs as fields** - Define test inputs as test case fields, not function arguments

**Table-Driven Test Examples:**

```typescript
// Good: Basic table-driven test using test.each
describe('calculateTotal', () => {
  test.each([
    { items: [{ price: 10, quantity: 2 }], expected: 20 },
    { items: [{ price: 5, quantity: 3 }], expected: 15 },
    { items: [], expected: 0 },
    { items: [{ price: 10, quantity: 0 }], expected: 0 },
  ])('returns $expected for items $items', ({ items, expected }) => {
    expect(calculateTotal(items)).toBe(expected);
  });
});

// Good: Table-driven test with descriptive test names
describe('validateEmail', () => {
  test.each([
    ['valid@example.com', true, 'valid email'],
    ['test@test.co.uk', true, 'valid email with subdomain'],
    ['invalid-email', false, 'missing @ symbol'],
    ['@example.com', false, 'missing local part'],
    ['user@', false, 'missing domain'],
    ['', false, 'empty string'],
  ])('validateEmail(%s) returns %s (%s)', (email, expected, description) => {
    expect(validateEmail(email)).toBe(expected);
  });
});

// Good: Table-driven test with objects for complex cases
describe('User.validate', () => {
  test.each([
    {
      name: 'valid user',
      user: { name: 'John Doe', email: 'john@example.com', role: 'user' },
      expectedError: null,
    },
    {
      name: 'invalid email',
      user: { name: 'John Doe', email: 'invalid', role: 'user' },
      expectedError: ValidationError,
    },
    {
      name: 'empty name',
      user: { name: '', email: 'john@example.com', role: 'user' },
      expectedError: ValidationError,
    },
    {
      name: 'invalid role',
      user: { name: 'John Doe', email: 'john@example.com', role: 'superuser' },
      expectedError: ValidationError,
    },
  ])('$name', ({ user, expectedError }) => {
    if (expectedError) {
      expect(() => User.validate(user)).toThrow(expectedError);
    } else {
      expect(() => User.validate(user)).not.toThrow();
    }
  });
});

// Good: Table-driven test with async functions
describe('fetchUser', () => {
  test.each([
    {
      name: 'existing user',
      userId: '123',
      mockResponse: { id: '123', name: 'John' },
      expectedResult: { id: '123', name: 'John' },
    },
    {
      name: 'non-existent user',
      userId: '999',
      mockResponse: null,
      expectedResult: null,
    },
  ])('$name', async ({ userId, mockResponse, expectedResult }) => {
    // Mock the API call
    jest.spyOn(global, 'fetch').mockResolvedValue({
      ok: mockResponse !== null,
      json: async () => ({ data: mockResponse }),
    } as Response);

    const result = await fetchUser(userId);
    expect(result).toEqual(expectedResult);
  });
});

// Good: Table-driven test with setup/teardown per case
describe('UserRepository', () => {
  let repository: UserRepository;

  beforeEach(() => {
    repository = new UserRepository();
  });

  describe('findById', () => {
    test.each([
      {
        name: 'finds existing user',
        setup: (repo: UserRepository) => repo.create({ id: '1', name: 'Alice' }),
        userId: '1',
        expected: { id: '1', name: 'Alice' },
      },
      {
        name: 'returns null for non-existent user',
        setup: () => {},
        userId: '999',
        expected: null,
      },
    ])('$name', async ({ setup, userId, expected }) => {
      await setup(repository);
      const result = await repository.findById(userId);
      expect(result).toEqual(expected);
    });
  });
});

// Good: Alternative syntax using array of arrays
describe('add', () => {
  test.each([
    [1, 2, 3],
    [0, 0, 0],
    [-1, 1, 0],
    [100, 200, 300],
  ])('add(%i, %i) returns %i', (a, b, expected) => {
    expect(add(a, b)).toBe(expected);
  });
});

// Good: Table-driven test with custom error messages
describe('parseJSON', () => {
  test.each([
    {
      input: '{"name": "John"}',
      expected: { name: 'John' },
      shouldThrow: false,
    },
    {
      input: '{invalid json}',
      expected: null,
      shouldThrow: true,
    },
    {
      input: '',
      expected: null,
      shouldThrow: true,
    },
  ])('parseJSON($input)', ({ input, expected, shouldThrow }) => {
    if (shouldThrow) {
      expect(() => parseJSON(input)).toThrow();
    } else {
      expect(parseJSON(input)).toEqual(expected);
    }
  });
});
```

**For Vitest:**
```typescript
import { describe, test, expect } from 'vitest';

// Vitest also supports test.each with the same syntax
describe('calculator', () => {
  test.each([
    { a: 1, b: 2, expected: 3 },
    { a: -1, b: -2, expected: -3 },
    { a: 0, b: 0, expected: 0 },
  ])('adds $a + $b to equal $expected', ({ a, b, expected }) => {
    expect(a + b).toBe(expected);
  });
});
```

**When NOT to use table-driven tests:**
- Single test case with complex mocking
- Tests that require significantly different setup per case
- Integration tests with heavy external dependencies
- Tests where the table structure becomes too complex

**Best Practices:**
1. Use descriptive test case names (use the `name` field or `$variable` syntax)
2. Keep test cases focused and related
3. Use objects for complex test cases, arrays for simple ones
4. Group related test.each blocks under describe blocks
5. Consider extracting complex setup logic to helper functions
6. Use TypeScript types for test case objects to ensure type safety

### Modern Tooling and Best Practices

**Build Tools:**
- **Vite**: Modern, fast bundler with TypeScript support out-of-the-box (recommended)
- **ESBuild**: Extremely fast bundler and minifier
- **Webpack 5**: Mature bundler with extensive plugin ecosystem
- **Turbo/Turborepo**: For monorepo management with built-in caching

**Code Quality Tools:**
- **ESLint with TypeScript**: Catch potential issues early
  ```bash
  pnpm add -D @typescript-eslint/parser @typescript-eslint/eslint-plugin
  ```
- **Prettier**: Consistent code formatting
  ```bash
  pnpm add -D prettier
  ```
- **TypeDoc**: Auto-generate documentation from TypeScript code

**Testing Frameworks:**
- **Vitest**: Fast, Vite-native test runner (recommended for Vite projects)
- **Jest**: Popular, mature testing framework
- **Playwright/Cypress**: E2E testing

**Type-Safe Libraries:**
- **Zod**: Runtime type validation with TypeScript inference
- **tRPC**: End-to-end type-safe APIs
- **Prisma**: Type-safe database ORM

Example tsconfig.json for modern projects:
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "Bundler", // For Vite/ESBuild
    "resolveJsonModule": true,
    "allowJs": true,
    "checkJs": false,
    "jsx": "react-jsx",
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "outDir": "./dist",
    "removeComments": true,
    "noEmit": true,
    "isolatedModules": true,
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "skipLibCheck": true
  },
  "include": ["src"],
  "exclude": ["node_modules", "dist"]
}
```

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract project-specific patterns
- Note architecture constraints
- Identify preferred libraries

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into small, testable units
- Identify dependencies

**Step 3: Read Existing Code**
- Review similar implementations
- Understand existing patterns
- Note utilities and helpers available

**Step 4: Implement Changes**
- Follow patterns from design.md
- Use naming conventions from project
- Write clean, readable code
- Add appropriate comments
- Match code style of existing files

**Step 5: Ensure Quality and Fix All Issues**

**CRITICAL: All quality checks must pass before considering the task complete.**

1. **Run type checking**
   ```bash
   npm run type-check
   # or
   npx tsc --noEmit
   ```
   - If type errors exist, fix them immediately
   - Do not proceed until type checking passes

2. **Run linter**
   ```bash
   npm run lint
   ```
   - If lint errors exist, try auto-fix first: `npm run lint -- --fix`
   - Fix any remaining issues manually
   - Do not proceed until linter passes

3. **Run formatter (if separate from linter)**
   ```bash
   npx prettier --check .
   # or
   npm run format:check
   ```
   - If formatting issues exist, auto-fix: `npx prettier --write .`
   - Review the changes

4. **Build the project**
   ```bash
   npm run build
   ```
   - If build fails, fix the errors immediately
   - Do not proceed until build succeeds

5. **Run tests**
   ```bash
   npm test
   ```
   - If tests fail, fix them immediately
   - Add new tests (preferably table-driven) if implementing new features
   - Update existing tests if modifying behavior
   - Do not proceed until all tests pass

6. **Update tests for your changes**
   - If you added a new function/component, add corresponding tests
   - If you modified behavior, update existing tests
   - Ensure test coverage is maintained or improved

**Iterative Fix Process:**
- If any step fails, fix the issues and re-run ALL previous steps
- Continue iterating until ALL checks pass:
  - ✅ No type errors
  - ✅ No lint issues
  - ✅ Code properly formatted
  - ✅ Build succeeds
  - ✅ All tests pass

## Commands Reference

**CRITICAL: Always use pnpm, never npm**

```bash
# Package management
pnpm install
pnpm add <package>
pnpm add -D <dev-package>
pnpm update

# Type checking
pnpm run type-check
pnpm exec tsc --noEmit

# Linting and formatting
pnpm run lint
pnpm run lint -- --fix
pnpm exec prettier --check .
pnpm exec prettier --write .

# Build
pnpm run build

# Tests
pnpm test
pnpm run test:watch
pnpm run test:coverage

# Development
pnpm run dev
```

## Common Pitfalls to Avoid

### Type Safety Issues
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

### Error Handling
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

### Testing
```typescript
// Bad: Duplicated test code
test('adds 1 + 2', () => {
  expect(add(1, 2)).toBe(3);
});
test('adds 5 + 3', () => {
  expect(add(5, 3)).toBe(8);
});
test('adds 0 + 0', () => {
  expect(add(0, 0)).toBe(0);
});

// Good: Table-driven test
test.each([
  [1, 2, 3],
  [5, 3, 8],
  [0, 0, 0],
])('add(%i, %i) returns %i', (a, b, expected) => {
  expect(add(a, b)).toBe(expected);
});
```

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Extract project-specific conventions
- [ ] Analyze project structure (check for pnpm-lock.yaml, package.json)
- [ ] Search for similar implementations
- [ ] Review existing patterns
- [ ] Verify package manager (MUST be pnpm, not npm)

### During Implementation
- [ ] Define all types/interfaces first (strict mode, no `any`)
- [ ] Use TypeScript strict mode configuration
- [ ] Implement fail-fast error handling (early validation, custom errors)
- [ ] Follow existing naming conventions
- [ ] Use early returns instead of nested conditions
- [ ] Keep functions/classes private unless needed for public API
- [ ] Write table-driven tests alongside code
- [ ] Prefer updating existing code over creating new functions (DRY principle)

### After Implementation - MUST ALL PASS
- [ ] Run type checking (`pnpm run type-check`) - **FIX ALL TYPE ERRORS**
- [ ] Run linters (`pnpm run lint`) - **FIX ALL LINT ISSUES**
- [ ] Run formatter if applicable - **FIX ALL FORMATTING ISSUES**
- [ ] Build successfully (`pnpm run build`) - **MUST SUCCEED**
- [ ] Run tests (`pnpm test`) - **ALL TESTS MUST PASS**
- [ ] Add/update table-driven tests for new or modified code
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation if needed

## Key Principles

1. **Simplicity First**: DRY principle, early returns, minimal public API, prefer breaking changes for simplicity
2. **Project Guidelines**: Always read and follow `.claude/design.md`
3. **Type Safety**: Strict mode, no `any`, use `unknown`, leverage advanced types
4. **Fail-Fast**: Early validation, never swallow errors, meaningful error messages
5. **Consistency**: Match existing code style and patterns from design.md
6. **Latest Versions**: Use latest package versions unless compatibility issues
7. **Testing**: Table-driven tests as primary approach, injection over global state
8. **Package Manager**: MUST use pnpm, NEVER use npm or legacy-peer-deps

## Version History

### Version 1.1.0 (2025-01-16)
- Added core principles: Simplicity, Consistency, Fail-Fast, Latest Versions
- Enhanced type safety section with modern TypeScript features (template literals, mapped types, conditional types, const assertions)
- Added strict tsconfig.json examples with moduleResolution options
- Expanded error handling with fail-fast principle and comprehensive examples
- Enhanced testing section with table-driven testing best practices from Internet research
- Added testing guidelines: injection over global state, test inputs as fields
- Added package management section emphasizing pnpm requirement
- Added directory structure approaches (feature-based, type-based, hybrid)
- Expanded function patterns with early returns and fail-fast examples
- Added modern tooling section (Vite, ESBuild, Vitest, Zod, tRPC)
- Updated commands reference to use pnpm exclusively
- Updated checklist to reflect all new requirements
- Updated key principles to emphasize simplicity, fail-fast, and pnpm usage

### Version 1.0.0 (Initial)
- Initial TypeScript skill with basic best practices
- Type safety, code organization, error handling, testing guidelines
- Table-driven testing examples with Jest and Vitest
