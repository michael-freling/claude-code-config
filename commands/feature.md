---
description: Add or update a feature following a structured workflow with role-based agents
argument-hint: "describe the feature to add or update"
allowed-tools: ["*"]
---

# Feature

Add or update a feature following a structured workflow with role-based agents.

## Arguments

$ARGUMENTS

## Workflow

### 1. Analysis and Planning (Architect Agent)
- Analyze existing codebase structure and patterns
- Plan implementation changes
- Identify opportunities for parallel development
- **Confirm plan with user before proceeding**

### 2. Implementation (Software Engineer Agents)
Implement with incremental changes per task repeatedly until complete feature is implemented. For each task, the appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills:

a. **Write Code with Tests**
   - Follow language/framework skill guidelines
   - Implement feature code and corresponding tests
   - Maintain consistency with existing patterns
   - Keep changes simple and focused

b. **Verify Changes**
   - Run linters and formatters
   - Execute relevant tests
   - Fix any issues found

c. **Peer Review**
   - Different software engineer agent reviews implementation
   - Validates code quality and standards
   - Suggests improvements if needed

d. **Commit Change**
   - Commit the incremental change before moving to next change

### 3. Final Review (Code Reviewer Agent)
- Comprehensive quality validation
- Verify all changes work together correctly
- Standards adherence check
- Security and performance review
- Final approval

## Guidelines

- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Use appropriate language/framework skills
- Maintain test coverage
- Ensure backward compatibility unless explicitly allowed to break it
