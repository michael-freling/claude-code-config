---
description: Refactor the codebase following structured workflow with architecture design and review
argument-hint: "describe what improvements to make"
allowed-tools: ["*"]
---

# Refactor

$ARGUMENTS

## Workflow

### 1. Analyze Existing Codebase and Architecture (Architect Agent)

Analyze the current codebase to understand:
- Current structure and patterns
- Areas that need refactoring
- Potential improvements

### 2. Clean Up Unnecessary Code (Architect Agent)

Before planning implementation, clean up unnecessary code in areas where changes will be made:
- Remove dead code and unused imports
- Remove redundant logic
- Identify and document technical debt

### 3. Create Design and Plan Changes (Architect Agent)

Design the refactoring including:
- Architecture improvements
- Code structure changes
- Implementation plan with specific tasks
- Identify opportunities for parallel work

### 4. Get Design Review (Code Reviewer Agent)

Have the refactoring plan reviewed for:
- Soundness of approach
- Potential risks or issues
- Better alternatives

### 5. Confirm Plan with User

**IMPORTANT**: Present the design, whether backward compatibility is required or not, and plan to the user. Confirm the plan is good before proceeding. Do not start implementation until you get approval from the user.

### 6. Create Git Worktrees

Create new git worktrees in the `../worktrees` directory for the actual implementation. The worktree names must include the ticket number provided.

If tasks can be worked on in parallel, create multiple worktrees to enable concurrent development.

### 7. Implementation (Software Engineer Agents)

Subagents make each change with review by other agents in the new worktrees. For each task, follow this process:

a. **Write Code with Tests**
   - The appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills
   - Implement refactoring following language/framework skill guidelines
   - Maintain existing functionality
   - Improve code structure and readability
   - Remove duplication (DRY principle)
   - Add or update tests as needed

b. **Review Changes**
   - A different agent reviews the refactoring
   - Validates improvements made
   - Ensures no functionality broken
   - Verifies test coverage

c. **Commit Changes**
   - Commit the incremental change before moving to next task

### 8. Create GitHub PR

Once all changes are completed in the worktrees, create a GitHub Pull Request.

## Guidelines

- Clean up unnecessary code before refactoring
- Maintain existing functionality (tests must pass)
- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Prefer breaking backward compatibility unless explicitly prohibited
- Keep refactoring focused and incremental
- Ensure all tests pass after refactoring
- Include ticket number in commit messages
