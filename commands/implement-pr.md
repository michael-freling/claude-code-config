---
description: Implement code changes and create a PR, fixing CI errors until it passes
argument-hint: "describe the feature, bug fix, or refactoring needed"
allowed-tools: ["*"]
---

# Implement PR

$ARGUMENTS

## Change Types

| Type | Goal | Key Constraint |
|------|------|----------------|
| **Feature** | Add new functionality | - |
| **Fix** | Fix broken behavior | Must reproduce bug first |
| **Flaky Test Fix** | Eliminate test non-determinism | Must run tests multiple times to verify stability |
| **Refactor** | Improve structure, same behavior | Must preserve existing behavior |

## Workflow

### Phase 1: Understand

Clarify what the user wants:
- Identify the type of change (feature, fix, flaky test fix, or refactor)
- Ask clarifying questions if requirements are ambiguous
- Understand scope and success criteria
- For fixes: identify reproduction steps and expected vs actual behavior
- For flaky test fixes: identify which test(s) are flaky and their failure patterns

### Phase 2: Explore

Investigate the codebase:
- Analyze existing structure and patterns
- For features: identify implementation approaches and trade-offs
- For fixes: reproduce the error and find root cause
- For flaky test fixes: identify the source of non-determinism (timing, ordering, shared state, external dependencies, race conditions)
- For refactors: identify code smells and improvement opportunities

### Phase 3: Present Plan

**IMPORTANT**: Present the plan to the user before coding:
- Analysis findings (root cause for fixes, options for features)
- Proposed approach
- Files to modify or create
- Potential risks or side effects

Wait for user approval before proceeding.

### Phase 4: Implement

Once approved, implement the changes:

**For fixes:**
1. Write a failing test that reproduces the bug
2. Implement the fix
3. Verify the test passes

**For features:**
1. Implement the feature
2. Write tests for new functionality

**For flaky test fixes:**
1. Run the flaky test multiple times (5-10 runs) to confirm flakiness and understand failure rate
2. Identify the root cause of non-determinism
3. Implement the fix (e.g., add proper synchronization, remove timing dependencies, isolate shared state)
4. Run the test multiple times (at least 10 runs) to verify it's now stable

**For refactors:**
1. Ensure existing tests pass before changes
2. Make structural changes
3. Verify tests still pass after changes

Use appropriate skills (e.g., `/write-tests`).

### Phase 5: Verify

Verify the implementation is correct:
- Run all tests and ensure they pass
- Run linters and fix any issues
- Manually verify the change works as expected
- Check for regressions

**For flaky test fixes - additional verification:**
- Run the previously flaky test at least 10 consecutive times
- All runs must pass to confirm the fix is stable
- Example verification commands:
  - Go: `for i in {1..10}; do go test -v -run TestName ./path/... || exit 1; done`
  - Jest: `for i in {1..10}; do npx jest --testNamePattern="test name" || exit 1; done`
  - pytest: `for i in {1..10}; do pytest -v test_file.py::test_name || exit 1; done`
- If any run fails, the fix is incomplete - investigate further

### Phase 6: Commit

Group related changes into logical commits:
- Write clear commit messages describing the "why"
- Separate refactoring from feature/fix changes when possible

### Phase 7: Create PR

Push changes and create the PR:
- Push changes to remote
- Create PR with clear summary and test plan

### Phase 8: Fix CI Until Green

Monitor and fix CI:
1. Wait for CI to complete (at least 1 minute for jobs to start)
2. Check CI status every 5 minutes until complete
3. If CI fails:
   - Analyze the errors
   - Fix issues and push new commits
   - Repeat until CI passes

## Guidelines

- Follow coding guidelines (DRY, simplicity, explicit over implicit)
- Adhere to project-specific guidelines if available
- Maintain or increase test coverage
- For fixes: always reproduce first, understand root cause
- For flaky test fixes: always verify stability with multiple runs (minimum 10)
- For refactors: never change behavior, only structure
