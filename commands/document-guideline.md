---
description: Analyze the codebase and create comprehensive project guidelines as an Architect
argument-hint: ""
allowed-tools: ["*"]
---

# Document Guideline

Analyze the codebase and create comprehensive project guidelines as an Architect.

## Instructions

1. Analyze the repository structure to determine if it's a single project or monorepo
2. For each project (or subproject in a monorepo):
   - Create a guideline document in `.claude/docs/` directory
   - Include the following sections:
     - **Architecture**: System design, components, and their relationships
     - **API Design**: Endpoint patterns, request/response formats, versioning
     - **Data Models**: Database schemas, entity relationships, data flow
     - **Design Patterns**: Used patterns and when to apply them
     - **Coding Best Practices**: Project-specific standards and conventions
3. If monorepo: Create separate guidelines per subproject, keep root guideline general
4. Output guidelines to `.claude/docs/guideline.md` (or per subproject)

## Requirements

- Keep guidelines concise and actionable
- Include examples where helpful
- Focus on project-specific patterns and conventions
- Update existing guidelines rather than creating duplicates
