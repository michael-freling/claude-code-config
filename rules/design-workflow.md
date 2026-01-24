# Design Workflow Conventions

This document establishes conventions for the design phase of feature development.

## Design Phase Ordering

Follow this sequence when designing features:

```
1. /write-product-spec    → Product requirements (WHAT and WHY)
         ↓
2. /design-ux             → User flows and wireframes
         ↓
3. /design-system-architecture → API, database, security design
         ↓
4. /design-frontend-details    → Frontend implementation design
   /design-backend-details     → Backend implementation design
         ↓
5. /implement-pr          → Implementation + PR
```

## Output Directory Convention

All design documents are organized by feature:

```
docs/
└── <feature-name>/
    ├── product-spec.md      # Product requirements
    ├── ux-design.md         # UX design document
    ├── wireframes/          # SVG wireframe files
    │   └── *.svg
    ├── architecture.md      # System architecture
    ├── frontend-design.md   # Frontend implementation design
    └── backend-design.md    # Backend implementation design
```

## Review Requirements

Each design phase requires review before proceeding:

| Design Phase | Designer Agent | Reviewer Agent |
|--------------|----------------|----------------|
| Product Spec | product-manager | product-spec-reviewer |
| UX Design | ux-designer | ux-design-reviewer |
| Architecture | software-architect | architecture-reviewer |
| Frontend Details | frontend-design-engineer | frontend-design-reviewer |
| Backend Details | backend-design-engineer | backend-design-reviewer |

## Diagram Conventions

- **User Flows**: Use Mermaid syntax
- **UI Layouts**: Use SVG wireframes
- **Architecture**: Use ASCII or Mermaid diagrams
- **Entity-Relationship**: Use ASCII or Mermaid diagrams

## Design Document Standards

When creating design documents:

1. **No Implementation Code**: Design documents contain specifications, not code
2. **Reference Existing Patterns**: Point to existing codebase patterns
3. **Document Trade-offs**: Explain decisions and alternatives considered
4. **Include Open Questions**: List items needing clarification
5. **User Confirmation**: Always get user approval before finalizing
