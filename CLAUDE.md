# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code across all repositories.

## Core Principles

### Document Organization
1. **.claude/docs/** directory for human-readable documentation
   - Output all latest project information here
   - Make it useful for both Claude Code and humans
   - Keep documentation current and well-organized

2. **.claude/archive/** directory for Claude Code-specific content
   - Store temporary documents and analysis here
   - These files are read only by Claude Code
   - Can contain detailed technical information, code analysis, and working notes

### Content Generation Guidelines
- **Simplicity and Consistency**: Generate documents, code, and all other content as simple and consistent as possible
- **Exception**: Documents in .claude/archive/ can be complex and detailed since they're only for Claude Code
- **Clarity**: Prioritize readability and maintainability in all human-facing content

## Working Practices

### When Starting a New Project
1. Analyze the repository structure
2. Create or update CLAUDE.md in the repository root with project-specific guidance
3. Output human-readable summaries to .claude/docs/
4. Store detailed analysis in .claude/archive/

### Documentation Standards
- Use Markdown for documentation files
- Include clear headings and structure
- Add code examples where helpful

### File Organization
- **.claude/docs/**: README.md, ARCHITECTURE.md, API.md, CONTRIBUTING.md, etc.
- **.claude/archive/**: analysis-[date].md, code-map.json, dependencies.txt, notes.md, etc.

### Best Practices
- Always check for existing documentation before creating new files
- Update documentation after significant changes
- Use consistent naming conventions in a project
- Keep the main project CLAUDE.md focused on project-specific information
- Ensure consistency in documentation style and structure in a project
