---
description: Fix a bug by first reproducing the error and understanding its root cause
argument-hint: "describe the error and what fix is needed"
allowed-tools: ["*"]
---

# Fix

Fix a bug by first reproducing the error and understanding its root cause.

## Arguments

$ARGUMENTS

## Workflow

### 1. Analysis (Architect Agent)
- Analyze existing codebase
- Understand where the error occurs
- Identify affected components

### 2. Reproduction (Software Engineer Agent)
- Software engineer agent with appropriate tech stack reproduces the error
- Understand the root cause
- Document findings

### 3. Planning (Architect Agent)
- Based on analysis and root cause, plan fixes
- Identify opportunities for parallel changes
- **Confirm plan with user before proceeding**

### 4. Implementation (Software Engineer Agents)
For each task, the appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills:

a. **Update Code**
   - Implement fix following language/framework skill guidelines
   - Address root cause, not just symptoms
   - Keep changes minimal and focused

b. **Verify Changes**
   - Run linters and formatters
   - Execute tests to confirm fix
   - Verify no regressions introduced

c. **Peer Review**
   - Different software engineer agent reviews fix
   - Validates solution addresses root cause
   - Ensures no side effects

### 5. Final Review (Code Reviewer Agent)
- Comprehensive quality validation
- Verify fix resolves the issue
- Check for edge cases
- Final approval

## Notes

- Agents written in `~/.claude/roles` work as subagents on each step
- Multiple subagents can work in parallel when possible
- Each agent can commit each change on behalf of the user and push it to a remote repository **ONLY WHEN** the user instructs to do so
- The project may be a single project or monorepo

## Guidelines

- Always reproduce the error first
- Understand root cause before fixing
- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Add tests to prevent regression
- Document fix rationale if complex
