---
description: Document coding standards, conventions, and architecture for the project
---

# Document Project Design

Analyze the current project or subproject and create a comprehensive design document at `.claude/design.md` that captures:

## Your Task

1. **Identify Project Type and Structure**
   - Detect the primary language(s) and frameworks
   - Analyze the directory structure
   - **Identify if this is a monorepo with subprojects**
     - Check for multiple go.mod, package.json, or similar project markers
     - Look for workspace configurations (pnpm-workspace.yaml, lerna.json, etc.)
     - If monorepo detected, create design.md for EACH subproject

2. **Analyze Existing Code Patterns**
   - Search for common patterns in the codebase
   - Identify naming conventions actually used
   - Find architectural patterns (MVC, Clean Architecture, etc.)
   - Review how errors are handled
   - Check how tests are structured

3. **Document Design Information:**

   **For Monorepo Root (`.claude/design.md`):**
   - Monorepo structure and organization
   - List of subprojects and their purposes
   - Shared tooling and infrastructure
   - Inter-project dependencies
   - Monorepo-level build and deployment processes
   - **DO NOT include project-specific coding standards or best practices**

   **For Each Subproject (`<subproject>/.claude/design.md`):**

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

   **For Single-Project Repositories:**
   - Include all sections in the root `.claude/design.md`

## Process

1. **Detect Repository Type:**
   - Use the Task tool with subagent_type=Explore to analyze the codebase thoroughly
   - Look for monorepo indicators:
     - Multiple go.mod, package.json, Cargo.toml, etc.
     - Workspace configuration files (pnpm-workspace.yaml, lerna.json, go.work, etc.)
     - Multiple independent projects in subdirectories

2. **For Monorepo:**
   - Create root `.claude/design.md` with monorepo-level information only
   - Identify all subprojects (typically in apps/, packages/, services/, or similar)
   - For EACH subproject, create `<subproject>/.claude/design.md`:
     - Analyze subproject configuration files
     - Read representative files from the subproject
     - Identify patterns specific to that subproject
     - Document coding standards and best practices for that subproject

3. **For Single-Project Repository:**
   - Search for configuration files (tsconfig.json, .eslintrc, .prettierrc, go.mod, etc.)
   - Read representative files from each major component
   - Identify patterns by searching for common implementations
   - Create or update `.claude/design.md` with comprehensive documentation

## Output Format

**For Monorepo:**
- Root `.claude/design.md`: High-level monorepo structure, subproject list, shared infrastructure
- Each subproject `.claude/design.md`: Detailed coding standards, architecture, and best practices for that specific project

**For Single-Project:**
- Root `.claude/design.md`: Comprehensive design document with all sections

All design documents should have clear sections, concrete examples, and references to actual files in the codebase. Make them detailed enough that Claude Code can follow these patterns when implementing new features.

**Template**: Use `.claude/templates/design.md.template` as a reference for structure, but adapt to the actual project.

## Notes

- If design.md files already exist, update them rather than replacing them
- Focus on documenting what IS, not what SHOULD BE
- Include file references with line numbers where helpful
- Make examples concrete and specific to each project/subproject
- Keep it practical and actionable
- **In monorepos, ensure root design.md stays high-level and delegates project-specific details to subproject design.md files**
