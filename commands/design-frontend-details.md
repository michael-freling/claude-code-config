---
description: Design frontend implementation details by exploring codebase and identifying reusable components
argument-hint: "<feature-name> <frontend feature to design>"
allowed-tools: ["Read", "Glob", "Grep", "Task", "Write", "LSP"]
---

# Design Frontend Details

$ARGUMENTS

## Workflow

### 1. Gather Context (Frontend Design Engineer Agent)

Use the frontend-design-engineer agent to:
- Review product specification (if exists at `docs/<feature-name>/product-spec.md`)
- Review UX design (if exists at `docs/<feature-name>/ux-design.md`)
- Review architecture (if exists at `docs/<feature-name>/architecture.md`)

### 2. Explore Codebase (Frontend Design Engineer Agent)

Use the frontend-design-engineer agent to:
- Analyze existing component library
- Identify design system patterns
- Map state management approach
- Document data fetching patterns

### 3. Analyze Reusability (Frontend Design Engineer Agent)

Use the frontend-design-engineer agent to:
- List reusable components
- Identify components needing extension
- Note new components required

### 4. Design Components (Frontend Design Engineer Agent)

Use the frontend-design-engineer agent to:
- Design component hierarchy
- Define props interfaces
- Plan state management
- Document data flow

### 5. Review Design (Frontend Design Reviewer Agent)

Use the frontend-design-reviewer agent to:
- Check pattern consistency
- Review reusability
- Assess performance considerations
- Validate testability

### 6. Confirm with User

**IMPORTANT**: Present the frontend design and review feedback to the user. Confirm the design is acceptable before finalizing. Do not finalize until you get approval from the user.

### 7. Write Output

Write the frontend design to: `docs/<feature-name>/frontend-design.md`

## Output Structure

The output file should include:
- Design Overview
- Codebase Analysis (existing patterns)
- Reusability Assessment
- Component Architecture (tree, props, state)
- Data Flow
- Implementation Considerations (accessibility, performance, testing)
- Implementation Order

## Guidelines

- NO implementation code
- Reference existing patterns in codebase
- Prioritize component reuse
- Consider accessibility
- Design for testability
