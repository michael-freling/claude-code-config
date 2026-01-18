---
description: Implement code changes (features, bug fixes, or refactoring) following a structured workflow
argument-hint: "describe the feature, bug fix, or refactoring needed"
allowed-tools: ["*"]
---

# Implement

$ARGUMENTS

## Change Types

| Type | Goal | Key Constraint |
|------|------|----------------|
| **Feature** | Add new functionality | - |
| **Fix** | Fix broken behavior | Must reproduce bug first |
| **Refactor** | Improve structure, same behavior | Must preserve existing behavior |

## Workflow

### Phase 1: Understand

Clarify what the user wants:
- Identify the type of change (feature, fix, or refactor)
- Ask clarifying questions if requirements are ambiguous
- Understand scope and success criteria
- For fixes: identify reproduction steps and expected vs actual behavior

### Phase 2: Explore

Investigate the codebase:
- Analyze existing structure and patterns
- For features: identify implementation approaches and trade-offs
- For fixes: reproduce the error and find root cause
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

### Phase 6: PR Workflow (If Requested)

If the user asks to create a PR:

**6a. Commit with Appropriate Granularity**
- Group related changes into logical commits
- Write clear commit messages describing the "why"
- Separate refactoring from feature/fix changes when possible

**6b. Create PR**
- Push changes to remote
- Create PR with clear summary and test plan

**6c. Verify CI Passes**
- Wait for CI to complete (at least 1 minute for jobs to start)
- Check CI status every 5 minutes until complete

**6d. Fix Until CI Passes**
- If CI fails, analyze the errors
- Fix issues and push new commits
- Repeat until CI passes

## Guidelines

- Follow coding guidelines (DRY, simplicity, explicit over implicit)
- Adhere to project-specific guidelines if available
- Maintain or increase test coverage
- For fixes: always reproduce first, understand root cause
- For refactors: never change behavior, only structure
