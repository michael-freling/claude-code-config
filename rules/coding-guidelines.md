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
- **Load binary data from files** instead of embedding base64 or other encoded strings inline (images, fonts, certificates, etc.)

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

- **Use dependency injection** for external dependencies
- **Avoid global state** - pass dependencies explicitly
- **Use table-driven tests** - group related cases, separate success from error cases
- **Never skip or ignore test failures** - fix them properly
- Prefer not to lower coverage thresholds; write more tests instead

For detailed patterns and examples, use the `/write-tests` skill.
