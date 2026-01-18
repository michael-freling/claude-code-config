# Simplicity Principle

Simplicity means minimizing the **final complexity** of the codebase—not the size of the change.

## Core Principles

- Optimize for a clean end state, not a minimal diff
- **DRY (Don't Repeat Yourself)**: Extract shared logic, configs, and patterns into reusable components
- Minimize redundancy and fragmentation
- Group related things together; specify differences only where they exist
- Larger refactoring is acceptable if it results in a simpler final state
- **Breaking changes are acceptable** unless backward compatibility is explicitly required by users or project constraints

## Examples

### Code: High Cohesion

Split functions by actual responsibility, not by caller convenience.

**Scenario**: Function A is used by B, C, and D, but D only uses part of A. A new feature is needed for B and C only.

- **Wrong**: Add feature to A, making it complex for D's use case
- **Right**: Split A into X (core, used by all) and Y (extended, used by B/C only)

### Tests: Table-Driven Consolidation

Consolidate test cases into fewer test functions using table-driven tests.

- **Wrong**: Many separate test functions testing different aspects of the same function
- **Right**: One table-driven test function covering all cases (happy path and error cases)

