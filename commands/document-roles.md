---
description: Create role definitions for specialized agents in .claude/roles/ directory
argument-hint: ""
allowed-tools: ["*"]
---

# Document Roles

Create role definitions for specialized agents in `.claude/roles/` directory.

## Instructions

Create the following role documents in `.claude/roles/`:

### 1. Architect Role (`.claude/roles/architect.md`)
- Design high-level API designs, architecture, and data models
- Analyze codebase structure and patterns
- Plan implementation strategies
- Break down complex tasks into parallel workstreams

### 2. Software Engineer Roles (`.claude/roles/engineers/`)
Create specialized engineers per language/framework:
- **Golang Engineer** - Go development following best practices
- **TypeScript Engineer** - TypeScript development, prefer pnpm
- **Next.js Engineer** - Next.js/React development with TypeScript

Each engineer role should:
- Write clean, efficient, and maintainable code
- Verify changes with linting and tests
- Follow language-specific best practices

### 3. Code Reviewer Role (`.claude/roles/code-reviewer.md`)
- Validate implementation quality
- Check adherence to standards and guidelines
- Verify test coverage and code quality
- Ensure consistency with project conventions

## Output Structure

```
.claude/roles/
├── architect.md
├── code-reviewer.md
└── engineers/
    ├── golang.md
    ├── typescript.md
    └── nextjs.md
```

## Requirements

- Keep role definitions concise and clear
- Focus on responsibilities and expectations
- Include quality standards for each role
