# Claude Code Skills

This directory contains custom skills for Claude Code to help with common development tasks.

## Recommended Workflow

For best results, follow this workflow when working on a project:

1. **First time in a project**: Run `/document-design` to analyze and document the codebase
2. **Before implementing features**: Review `.claude/design.md` to understand conventions
3. **When implementing**: Use `/skill feature-dev` which automatically follows design.md

```
/document-design → Review .claude/design.md → /skill feature-dev → Implement
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

## Available Skills

### feature-dev

A comprehensive skill for adding or updating features in TypeScript, Next.js, Golang, and Protocol Buffers projects.

**IMPORTANT**: This skill automatically reads and follows `.claude/design.md` if it exists, ensuring all implementations match project-specific conventions.

**Usage:**
```
/skill feature-dev
```

**What it covers:**

- **TypeScript Best Practices**
  - Type safety with interfaces, types, and generics
  - Error handling with custom error classes
  - Function patterns and async/await
  - Code organization and barrel exports

- **Next.js Best Practices**
  - App Router (Server Components, Server Actions)
  - Pages Router (getServerSideProps, getStaticProps)
  - Component patterns (Client vs Server)
  - Optimization (next/image, next/link, dynamic imports)

- **Go Best Practices**
  - Code organization following standard layout
  - Error handling with custom error types
  - Interface design (small, focused interfaces)
  - Concurrency with context and goroutines
  - Table-driven testing

- **Protocol Buffers Best Practices**
  - Message design and field naming
  - Versioning and backward compatibility
  - Service design with proper RPC patterns
  - Common patterns (pagination, errors, resource names)

**Features:**

1. Language-specific code examples
2. Complete workflow demonstrations
3. Common pitfalls with bad vs good examples
4. Quick reference checklists
5. Command reference for each language

**When to use:**

- Adding new features to existing projects
- Updating or enhancing functionality
- Need to follow project-specific patterns
- Working with TypeScript, Next.js, Go, or Protocol Buffers

## Creating Your Own Skills

To create a new skill:

1. Create a new `.md` file in `.claude/skills/`
2. Follow this structure:
   - Title and overview
   - When to use this skill
   - Step-by-step process
   - Examples
   - Best practices
   - Common pitfalls

3. Make it actionable with specific code examples
4. Include language-specific guidance where relevant
