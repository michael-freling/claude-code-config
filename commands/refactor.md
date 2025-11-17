---
description: Refactor the codebase to improve code quality, structure, or performance
argument-hint: "describe what improvements to make"
allowed-tools: ["*"]
---

# Refactor

Refactor the codebase to improve code quality, structure, or performance.

## Arguments

$ARGUMENTS

## Workflow

### 1. Analysis and Planning (Architect Agent)
- Analyze existing codebase structure and patterns
- Identify code to be refactored
- Clean up unnecessary code in target areas
- Plan refactoring changes
- Identify opportunities for parallel work
- **Confirm plan with user before proceeding**

### 2. Implementation (Software Engineer Agents)
For each task, the appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills:

a. **Update Code**
   - Implement refactoring following language/framework skill guidelines
   - Maintain existing functionality
   - Improve code structure and readability
   - Remove duplication (DRY principle)

b. **Verify Changes**
   - Run linters and formatters
   - Execute all tests to ensure no regressions
   - Verify behavior unchanged

c. **Peer Review**
   - Different software engineer agent reviews refactoring
   - Validates improvements made
   - Ensures no functionality broken

### 3. Final Review (Code Reviewer Agent)
- Comprehensive quality validation
- Verify refactoring improves codebase
- Check for any regressions
- Final approval

## Notes

- Agents written in `~/.claude/roles` work as subagents on each step
- Multiple subagents can work in parallel when possible
- Each agent can commit each change on behalf of the user and push it to a remote repository **ONLY WHEN** the user instructs to do so
- The project may be a single project or monorepo

## Guidelines

- Clean up unnecessary code before refactoring
- Maintain existing functionality (tests must pass)
- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Prefer breaking backward compatibility unless explicitly prohibited
- Keep refactoring focused and incremental
- Ensure all tests pass after refactoring
