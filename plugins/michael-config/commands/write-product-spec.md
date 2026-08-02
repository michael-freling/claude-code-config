---
description: Write product specification with user stories (no technical details)
argument-hint: "<feature-name> <feature description>"
allowed-tools: ["Read", "Glob", "Grep", "Task", "Write"]
---

# Write Product Specification

$ARGUMENTS

## Workflow

### 1. Gather Context (Product Manager Agent)

Use the product-manager agent to:
- Understand the feature requirements
- Review any existing documentation
- Identify target users

### 2. Create Product Specification (Product Manager Agent)

Use the product-manager agent to:
- Write user stories with acceptance criteria
- Define non-functional requirements
- Document constraints and assumptions
- Identify out-of-scope items

### 3. Review Specification (Product Spec Reviewer Agent)

Use the product-spec-reviewer agent to:
- Review completeness
- Check clarity and testability
- Identify gaps or ambiguities
- Provide improvement suggestions

### 4. Confirm with User

**IMPORTANT**: Present the specification and review feedback to the user. Confirm the specification is acceptable before finalizing. Do not finalize until you get approval from the user.

### 5. Write Output

Write the product specification to: `docs/<feature-name>/product-spec.md`

## Output Structure

The output file should follow the Product Specification template:
- Problem Statement
- Goals and Success Metrics
- Target Users
- User Stories with Acceptance Criteria
- Non-Functional Requirements
- Constraints and Assumptions
- Out of Scope
- Open Questions

## Guidelines

- Focus on WHAT and WHY, not HOW
- NO technical implementation details
- NO UI mockups or wireframes
- NO database schemas or API specs
- Use clear, measurable acceptance criteria
