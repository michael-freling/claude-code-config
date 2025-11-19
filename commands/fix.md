---
description: Fix a bug by reproducing the error, understanding root cause, and planning fixes
argument-hint: "describe the error and what fix is needed"
allowed-tools: ["*"]
---

# Fix

$ARGUMENTS

## Workflow

### 1. Analyze Existing Codebase (Architect Agent)

Analyze the codebase to understand where the error happens:
- Identify affected components
- Understand code flow leading to the error
- Map out relevant parts of the codebase

### 2. Reproduce Error and Find Root Cause (Software Engineer Agent)

A software engineer with the appropriate tech stack must:
- Reproduce the error to confirm the issue
- Understand the root cause of the error
- Document findings and analysis

### 3. Plan Changes (Architect Agent)

Based on the analysis and root cause:
- Plan the necessary fixes
- Identify opportunities for parallel changes
- Ensure the fix addresses the root cause, not just symptoms

### 4. Confirm Plan with User

**IMPORTANT**: Present the analysis and plan to the user. Confirm the plan is good before proceeding with implementation.

### 5. Implementation (Software Engineer Agents)

For each task, the appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills:

a. **Update Code**
   - Implement fix following language/framework skill guidelines
   - Address root cause, not just symptoms
   - Keep changes minimal and focused
   - Add tests to prevent regression

b. **Verify Changes**
   - Run linters and formatters
   - Execute tests to confirm fix
   - Verify no regressions introduced

c. **Review Changes**
   - Different software engineer agent reviews fix
   - Validates solution addresses root cause
   - Ensures no side effects
   - Verifies test coverage

d. **Commit Changes**
   - Commit the incremental change

## Guidelines

- Always reproduce the error first before attempting fixes
- Understand root cause before fixing
- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Add tests to prevent regression
- Document fix rationale if complex
- Include ticket number in commit messages
