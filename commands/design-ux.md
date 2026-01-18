---
description: Design UX with user flow diagrams and SVG wireframes
argument-hint: "<feature-name> <UI to design>"
allowed-tools: ["Read", "Glob", "Grep", "Task", "Write"]
---

# Design UX

$ARGUMENTS

## Workflow

### 1. Gather Context (UX Designer Agent)

Use the ux-designer agent to:
- Review product specification (if exists at `docs/<feature-name>/product-spec.md`)
- Analyze existing UI patterns in codebase
- Understand design system conventions

### 2. Create User Flows (UX Designer Agent)

Use the ux-designer agent to:
- Create user flow diagrams (Mermaid syntax)
- Map primary and alternative paths
- Document error flows

### 3. Create Wireframes (UX Designer Agent)

Use the ux-designer agent to:
- Create SVG wireframes for key screens
- Define component layout
- Document interaction states
- Note accessibility requirements

### 4. Review Design (UX Design Reviewer Agent)

Use the ux-design-reviewer agent to:
- Review usability
- Check accessibility compliance
- Verify consistency with design system
- Check completeness of states

### 5. Confirm with User

**IMPORTANT**: Present the UX design and review feedback to the user. Confirm the design is acceptable before finalizing. Do not finalize until you get approval from the user.

### 6. Write Output

Write the UX design to: `docs/<feature-name>/ux-design.md`
Write SVG wireframes to: `docs/<feature-name>/wireframes/`

## Output Structure

The output file should include:
- Design Overview
- User Flows (Mermaid diagrams)
- Screen Layouts (SVG wireframes)
- Component Specifications
- Interaction Patterns
- Responsive Design Notes
- Accessibility Requirements

## Guidelines

- Start with user flows before layouts
- Use Mermaid syntax for flow diagrams
- Create clean, minimal SVG wireframes
- Document all component states
- Consider accessibility from the start
