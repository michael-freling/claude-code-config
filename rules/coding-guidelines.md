# Coding Guidelines

## Principles

- Simplicity is the most important thing. Follow simplicity rule.
- Explicit is better than implicit.

## Dependencies

- When installing applications, libraries, or tools, always check and use the most latest and stable version with compatibility with existing systems.

## Code Readability

- Use specific names, not generic ones (shared, common, utils, info).
- Write the code with minimal comments — only high-level explanations of purpose, architecture, or non-obvious decisions. No line-by-line comments
- Delete deadcodes.
- Delete assignments of the default or zero values.
- **Prefer to continue or return early** than nesting code
   - "if is bad, else is worse"

## External I/O

- **Batch operations** instead of loops for databases, APIs, and other external calls
  - Reduces network round-trips
  - Enables optimized bulk queries and writes
- **Transactions** for multi-step writes to ensure atomicity

## Logging

- Use structured logging with key-value pairs
- Log levels: Debug (detailed), Info (normal), Warn (recoverable), Error (unexpected)
- Never log sensitive data (passwords, tokens, PII)

## Testing

### Writing Testable Code

- **Use dependency injection** for external dependencies (databases, APIs, etc.)
- **Avoid global state** (global variables, environment variables, external files). Pass dependencies explicitly instead.
- **Avoid environment-specific branching** in core logic
- Prefer not to use coverage ignore comments; refactor to make code testable

### Writing Effective Tests

- DO NOT SKIP or IGNORE errors from tests, pre-commits, linters, or any verification. Fix them properly.
- **Use table-driven testing**
    - Split happy and error test sets if complicated
    - Reduces code duplication and improves maintainability
- Prefer scalar or object values as test case fields, not functions (except for setup/mock functions)
- Prefer not to lower coverage thresholds; write more tests instead
