# Claude Code Skills

This directory contains custom skills for Claude Code to help with common development tasks.

## Recommended Workflow

For best results, follow this workflow when working on a project:

1. **First time in a project**: Run `/document-design` to analyze and document the codebase
2. **Before implementing features**: Review `.claude/design.md` to understand conventions
3. **When implementing**: Claude Code will automatically use language-specific skills based on your task

```
/document-design → Review .claude/design.md → Implement (skills auto-activate)
```

## Available Commands

### /document-design

Analyzes your codebase and creates `.claude/design.md` documenting:
- Coding standards and naming conventions
- Architecture patterns and directory structure
- Best practices used in the project
- Language-specific conventions
- Concrete examples from the codebase

**Usage:**
```
/document-design
```

**When to use:**
- First time working on a project
- When conventions are unclear
- After major architectural changes
- For monorepos: run in each subproject directory

### /monorepo-change

Applies changes across multiple subprojects in a monorepo with planning, parallel execution, and review.

**Usage:**
```
/monorepo-change [change description]
```

**Examples:**
```
/monorepo-change Add user authentication feature
/monorepo-change Fix race condition in API gateway
/monorepo-change Update all Node.js dependencies
/monorepo-change Refactor error handling for consistency
```

**When to use:**
- Working in a monorepo with multiple subprojects
- Changes that span multiple languages or frameworks
- Features, bugfixes, refactoring, configuration, or dependency updates
- When you need coordinated changes across services

## Available Skills

Skills are automatically activated by Claude Code when working with specific languages or frameworks. Each skill reads and follows `.claude/design.md` if it exists, ensuring all implementations match project-specific conventions.

### typescript

Automatically activated when adding or updating TypeScript code.

**Provides guidance on:**
- Type safety with interfaces, types, and generics
- Error handling with custom error classes
- Function patterns and async/await
- Code organization and barrel exports
- Common TypeScript pitfalls

### nextjs

Automatically activated when working with Next.js applications.

**Provides guidance on:**
- App Router (Server Components, Server Actions, metadata)
- Pages Router (getServerSideProps, getStaticProps)
- Component patterns (Client vs Server)
- Optimization (next/image, next/link, dynamic imports)
- API routes and data fetching

### golang

Automatically activated when adding or updating Go code.

**Provides guidance on:**
- Code organization following standard Go layout
- Error handling with custom error types
- Interface design (small, focused interfaces)
- Concurrency with context and goroutines
- Table-driven testing
- Struct design and composition

### protobuf

Automatically activated when working with Protocol Buffer files (.proto).

**Provides guidance on:**
- Message design and field naming
- Versioning and backward compatibility
- Service design with proper RPC patterns
- Common patterns (pagination, errors, resource names)
- Field numbering best practices

## Creating Your Own Skills

To create a new skill:

1. Create a directory: `skills/your-skill-name/`
2. Create a `SKILL.md` file with frontmatter:
   ```yaml
   ---
   name: your-skill-name
   description: What it does and when Claude should use it (max 1,024 characters)
   ---
   ```

3. Include these sections in the content:
   - Overview and when to use
   - Step-by-step process
   - Best practices with examples
   - Common pitfalls to avoid
   - Implementation checklists

4. Make it actionable with specific code examples
5. Include language-specific guidance where relevant

**Tips:**
- Use clear, specific descriptions so Claude knows when to activate the skill
- Reference `.claude/design.md` in your skill for project-specific patterns
- Keep skills focused on a single capability or language
- Test that the description triggers appropriate activation
