---
name: write-tests
description: Write testable code and implement tests. Use when writing tests, improving test coverage, refactoring for testability, or implementing features that need tests.
---

# Write Tests Skill

Guidance for writing testable code and implementing effective tests.

For language-specific examples, see:
- `references/golang.md` - Go testing patterns
- `references/typescript.md` - TypeScript/React testing patterns

## Writing Testable Code

### Dependency Injection

Pass dependencies as parameters instead of using global access. This enables mocking in tests.

### Extract Interfaces

Define interfaces for external dependencies to enable mocking.

### Delete Dead Code

Before adding tests:
- Remove unused functions
- Delete commented-out code
- Remove unreachable branches

### Split Large Functions

- Each function should do one thing
- Extract distinct operations into separate functions

## Test Structure

### Table-Driven Tests

Always use table-driven tests with separate success and error cases:
- Group related test cases in a table/array
- Use descriptive names for each case
- Separate success cases from error cases for clarity

### Mock Dependencies

Create mock implementations of interfaces for testing:
- Keep mocks simple and focused
- Allow configuring mock behavior per test case
- Create test data builders with sensible defaults

## Coverage Strategy

### What to Test
- All public functions
- Error handling paths
- Edge cases (empty, null, max values)
- Boundary conditions

### What NOT to Test
- Third-party libraries
- Simple getters/setters
- Framework internals

### Coverage Gaps to Look For
- Untested error branches
- Missing edge cases
- Conditional logic not fully covered

## Workflow

When implementing features with tests:

1. **Analyze testability** - Identify dependencies and refactoring needed
2. **Refactor first** - Make code testable before writing tests
3. **Write tests** - Use table-driven patterns, cover success/error/edge cases
4. **Verify coverage** - Ensure coverage is maintained or increased
