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
For each task, the appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills:

a. **Update Code**
   - Follow language/framework skill guidelines
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

### 3. Final Review (Code Reviewer Agent)
- Comprehensive quality validation
- Standards adherence check
- Security and performance review
- Final approval

## Notes

- Agents written in `~/.claude/roles` work as subagents on each step
- Multiple subagents can work in parallel when possible
- Each agent can commit each change on behalf of the user and push it to a remote repository **ONLY WHEN** the user instructs to do so
- The project may be a single project or monorepo

## Guidelines

- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Use appropriate language/framework skills
- Maintain test coverage
- Ensure backward compatibility unless explicitly allowed to break it
