---
description: Document coding standards, conventions, and architecture for the project
---

# Document Project Design

Analyze the current project or subproject and create a comprehensive design document at `.claude/design.md` that captures:

## Your Task

1. **Identify Project Type and Structure**
   - Detect the primary language(s) and frameworks
   - Analyze the directory structure
   - Identify if this is a monorepo with subprojects

2. **Analyze Existing Code Patterns**
   - Search for common patterns in the codebase
   - Identify naming conventions actually used
   - Find architectural patterns (MVC, Clean Architecture, etc.)
   - Review how errors are handled
   - Check how tests are structured

3. **Document the Following in `.claude/design.md`:**

   **Project Overview:**
   - Project name and description
   - Primary languages and frameworks
   - Key dependencies

   **Directory Structure:**
   - Explain the purpose of each major directory
   - Document any special conventions

   **Coding Standards:**
   - Naming conventions (files, variables, functions, types)
   - Code organization principles
   - Formatting rules (if using prettier, gofmt, etc.)

   **Architecture:**
   - Overall architecture pattern
   - How layers/modules are organized
   - Data flow and dependencies
   - Key design decisions

   **Best Practices:**
   - Error handling patterns
   - Testing conventions
   - Type usage (TypeScript)
   - Interface design (Go)
   - Concurrency patterns
   - API design

   **Language-Specific Conventions:**
   - TypeScript: Type definitions, async patterns
   - Go: Interface usage, error wrapping
   - Protocol Buffers: Versioning, field numbering

   **Examples from the Codebase:**
   - Include actual code snippets showing the patterns
   - Reference specific files as examples

## Process

1. Use the Task tool with subagent_type=Explore to analyze the codebase thoroughly
2. Search for configuration files (tsconfig.json, .eslintrc, .prettierrc, go.mod, etc.)
3. Read representative files from each major component
4. Identify patterns by searching for common implementations
5. Create or update `.claude/design.md` with comprehensive documentation
6. If this is a monorepo, consider creating design.md in subproject directories

## Output Format

Create `.claude/design.md` with clear sections, concrete examples, and references to actual files in the codebase. Make it detailed enough that Claude Code can follow these patterns when implementing new features.

**Template**: Use `.claude/templates/design.md.template` as a reference for structure, but adapt to the actual project.

## Notes

- If `.claude/design.md` already exists, update it rather than replacing it
- Focus on documenting what IS, not what SHOULD BE
- Include file references with line numbers where helpful
- Make examples concrete and specific to this project
- Keep it practical and actionable
