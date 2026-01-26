# Coding Guidelines

## Principles

- **Simplicity is the most important thing.** Minimize the final complexity of the codebase—not the size of the change.
  - **Breaking changes are acceptable** unless backward compatibility is explicitly required by users or project constraints


### Simplicity Example: High Cohesion

Split functions by actual responsibility, not by caller convenience.

**Scenario**: Function A is used by B, C, and D, but D only uses part of A. A new feature is needed for B and C only.

- **Wrong**: Add feature to A, making it complex for D's use case
- **Right**: Split A into X (core, used by all) and Y (extended, used by B/C only)


### Simplicity Example: Separate Success and Failure Paths

When a method handles both success and failure outcomes with substantially different logic, split into separate methods.

- **Wrong**: `CompleteJob(success bool)` with branching logic based on the flag
- **Right**: `CompleteJob()` and `FailJob()` as separate methods with clear responsibilities


## Dependencies

- When installing applications, libraries, or tools, always check and use the most latest and stable version with compatibility with existing systems.

## Code Readability

- Use specific names, not generic ones (shared, common, utils, info).
- Write the code with minimal comments — only high-level explanations of purpose, architecture, or non-obvious decisions. No line-by-line comments
- Delete deadcodes.
- Delete assignments of the default or zero values.
- **Load binary data from files** instead of embedding base64 or other encoded strings inline (images, fonts, certificates, etc.)
- **Validate only at context boundaries** - add defensive checks (nil checks, validation, error handling) only where data crosses trust domains: incoming client requests and outgoing third-party API responses; internal code should trust internal code

## Code Reuse

- **Single definition** - each abstraction (interface, type, function) should have exactly one definition in the codebase; if multiple packages need it, define once and import
- **Provider owns shared abstractions** - define abstractions in the provider/implementation package, not in each consumer
- **Check existing content before writing** - search for similar code, logic, or documentation before adding new content; reuse or extend existing patterns rather than duplicating; when consolidating files, read the destination thoroughly first

## External I/O

- **Batch operations** instead of loops for databases, APIs, and other external calls
  - Reduces network round-trips
  - Enables optimized bulk queries and writes
- **Transactions** for multi-step writes to ensure atomicity
- **Pass values from program, not SQL functions** - compute values in application code instead of using SQL functions like `NOW()` or `CURRENT_TIMESTAMP`; enables tests to inject controlled values
- **Avoid redundant queries** - don't query external systems for data already available in the current context (message payload, request, event object); don't query the same record multiple times in a single operation

## Logging

- Log levels: Debug (detailed), Info (normal), Warn (recoverable), Error (unexpected)

## Testing

- **Avoid global state** - use dependency injection instead of relying on package-level variables, struct fields with shared state, databases, or environment variables
  - Prevents flaky tests from shared mutable state
  - Enables concurrent test execution (`t.Parallel()`)
- **Use table-driven tests** - group related cases, separate success from error cases
- **Never skip or ignore test failures** - fix them properly
- Prefer not to lower coverage thresholds; write more tests instead

For detailed patterns and examples, use the `/write-tests` skill.

## Concurrency

- **Synchronize shared state** - when multiple goroutines/threads access shared mutable state, use proper locking (mutex, channels, atomic operations)
- **Document retry policies** - when implementing retry logic, document or make configurable: max retry count, backoff strategy, which errors are retryable
