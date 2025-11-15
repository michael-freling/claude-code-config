---
name: typescript
description: TypeScript feature development following best practices for type safety, code organization, error handling, and testing. Use when adding or updating TypeScript code in projects.
---

# TypeScript Development Skill

A comprehensive skill for adding or updating features in TypeScript projects following best practices.

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

## TypeScript Best Practices

### Code Organization

- Follow existing directory structure
- Use barrel exports (index.ts) for clean imports
- Separate concerns: types, utils, constants
- Co-locate related files (file.ts, file.test.ts)

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

### Type Safety

- Define interfaces/types for all data structures
- Use strict TypeScript settings (strict: true)
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

### Function Patterns

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

### Error Handling

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

### Testing

**CRITICAL: Always use table-driven tests (test.each) as the primary testing approach.**

Table-driven tests in TypeScript/Jest provide the same benefits as in other languages:
- Reduce code duplication
- Make it easy to add new test cases
- Improve test readability
- Ensure consistent test structure
- Make test coverage gaps obvious

**Use table-driven tests for:**
- All functions with multiple test cases
- Error handling scenarios
- Edge cases and boundary conditions
- Different input/output combinations

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

```bash
# Type checking
npm run type-check
npx tsc --noEmit

# Linting and formatting
npm run lint
npm run lint -- --fix
npx prettier --check .
npx prettier --write .

# Build
npm run build

# Tests
npm test
npm run test:watch
npm run test:coverage

# Development
npm run dev
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
- [ ] Analyze project structure
- [ ] Search for similar implementations
- [ ] Review existing patterns

### During Implementation
- [ ] Define all types/interfaces first
- [ ] Use strict TypeScript (no `any`)
- [ ] Implement proper error handling
- [ ] Follow existing naming conventions
- [ ] Write table-driven tests alongside code

### After Implementation - MUST ALL PASS
- [ ] Run type checking (`npm run type-check`) - **FIX ALL TYPE ERRORS**
- [ ] Run linters (`npm run lint`) - **FIX ALL LINT ISSUES**
- [ ] Run formatter if applicable - **FIX ALL FORMATTING ISSUES**
- [ ] Build successfully (`npm run build`) - **MUST SUCCEED**
- [ ] Run tests (`npm test`) - **ALL TESTS MUST PASS**
- [ ] Add/update table-driven tests for new or modified code
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Type Safety**: Leverage TypeScript's type system fully
3. **Consistency**: Match existing code style and patterns
4. **Simplicity**: Write clear, maintainable code
5. **Testing**: Use table-driven tests to ensure changes work and don't break existing functionality
