---
description: Add or update a feature following a structured workflow with architecture design and review
argument-hint: "describe the feature to add or update"
allowed-tools: ["*"]
---

# Feature

$ARGUMENTS

## Workflow

### 1. Analyze Existing Codebase and Architecture (Architect Agent)

Analyze the current codebase structure, patterns, and architecture to understand the context for the new feature.

### 2. Create Design and Plan Changes (Architect Agent)

Design the new feature including:
- Architecture changes required
- API designs (if applicable)
- Data models (if applicable)
- Implementation plan with specific tasks
- Opportunities for parallel development

### 3. Get Design Review (Code Reviewer Agent)

Have the design and plan reviewed for:
- Architectural soundness
- Alignment with existing patterns
- Potential issues or improvements

### 4. Confirm Plan with User

**IMPORTANT**: Present the design and plan to the user. Do not proceed with implementation until you get approval from the user.

### 5. Create Git Worktree

Create a new git worktree in the `../worktrees` directory for the implementation. The worktree name must include the ticket number provided.

### 6. Implementation (Software Engineer Agents)

In the new worktree, implement changes with the following process for each task:

a. **Write Code with Tests**
   - The appropriate software engineer agent (Golang, TypeScript, Next.js, etc.) implements using Claude Code Skills
   - Follow language/framework skill guidelines
   - Implement feature code and corresponding tests
   - Maintain consistency with existing patterns
   - Keep changes simple and focused

b. **Review Changes**
   - A different agent reviews the implementation
   - Validates code quality and standards
   - Suggests improvements if needed
   - Ensures tests are adequate

c. **Commit Changes**
   - Commit the incremental change before moving to next task

### 7. Create GitHub PR

Once all changes are completed in the worktree, create a GitHub Pull Request.

## Guidelines

- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
- Use appropriate language/framework skills
- Maintain test coverage
- Prefer breaking backward compatibility unless explicitly prohibited
- Include ticket number in commit messages
