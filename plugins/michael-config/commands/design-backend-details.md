---
description: Design backend implementation details by exploring codebase and identifying reusable modules
argument-hint: "<feature-name> <backend feature to design>"
allowed-tools: ["Read", "Glob", "Grep", "Task", "Write", "LSP"]
---

# Design Backend Details

$ARGUMENTS

## Workflow

### 1. Gather Context (Backend Design Engineer Agent)

Use the backend-design-engineer agent to:
- Review product specification (if exists at `docs/<feature-name>/product-spec.md`)
- Review architecture (if exists at `docs/<feature-name>/architecture.md`)

### 2. Explore Codebase (Backend Design Engineer Agent)

Use the backend-design-engineer agent to:
- Analyze existing service structure
- Identify patterns and conventions
- Map dependency injection approach
- Document error handling patterns

### 3. Analyze Reusability (Backend Design Engineer Agent)

Use the backend-design-engineer agent to:
- List reusable modules
- Identify interfaces to implement
- Note new modules required

### 4. Design Modules (Backend Design Engineer Agent)

Use the backend-design-engineer agent to:
- Design package/module structure
- Define interfaces
- Plan data layer
- Document error handling strategy

### 5. Review Design (Backend Design Reviewer Agent)

Use the backend-design-reviewer agent to:
- Check scalability
- Review maintainability
- Assess testability
- Validate performance considerations
- Check security

### 6. Confirm with User

**IMPORTANT**: Present the backend design and review feedback to the user. Confirm the design is acceptable before finalizing. Do not finalize until you get approval from the user.

### 7. Write Output

Write the backend design to: `docs/<feature-name>/backend-design.md`

## Output Structure

The output file should include:
- Design Overview
- Codebase Analysis (existing patterns)
- Reusability Assessment
- Module Architecture (package structure, interfaces)
- Data Layer
- Quality Considerations (performance, testability, readability)
- Error Handling Strategy
- Implementation Order

## Guidelines

- NO implementation code
- Reference existing patterns in codebase
- Prioritize module reuse
- Design for testability (dependency injection)
- Consider performance implications
