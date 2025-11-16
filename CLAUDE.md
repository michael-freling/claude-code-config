# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## General Coding Guidelines

Follow best practices of coding and each language:

1. **Simplicity is the most important**
   - Instead of creating new functions, always update existing code for reusability (DRY principle)
   - Prefer to break backward compatibility unless users explicitly mention otherwise

2. **Ensure consistency** by following `.claude/design.md`
   - **CRITICAL**: Always read `.claude/design.md` if it exists before starting work
   - Project-specific guidelines override general best practices

3. **Ensure fail-fast** instead of silently killing errors
   - Never swallow errors
   - Validate inputs early
   - Throw meaningful error messages

4. **Prefer incremental changes** to Big-Bang changes

5. **Prefer using the latest versions** of modules/packages
   - If compatibility problems arise, downgrade the major version to support dependencies

6. **Comments MUST BE about WHY not WHAT**
   - Explain reasoning behind decisions, not what the code does
   - Code should be self-explanatory through good naming and structure

7. **Prefer not to make anything public/exported** outside of packages for encapsulation

8. **Prefer to continue or return early** than nesting code
   - "if is bad, else is worse"

## Testing Guidelines

1. **Use table-driven testing**
   - Split happy and error test sets if complicated
   - Reduces code duplication and improves maintainability

2. **Do not add redundant, meaningless test cases** if the purpose is the same as other test cases

3. **Prefer injecting values to changing global states**
   - Instead of updating environment variables, global variables, or changing working directories
   - Update existing functions to allow constructor or method injection

4. **Define test inputs as test case fields**, not as function arguments

## Working with Agents and Roles

**Use subagents and roles under `~/.claude/roles` as much as possible:**

- **Architect**: Design high-level API designs, architecture, and data models
- **Software Engineer** (language-specific): Write clean, efficient, maintainable code
- **Code Reviewer**: Validate implementation quality and adherence to standards

## Document Organization

1. **.claude/docs/** - Human-readable documentation
   - Output latest project information here
   - Make it useful for both Claude Code and humans

2. **.claude/archive/** - Claude Code-specific content
   - Temporary documents and detailed analysis
   - Read only by Claude Code

3. **.claude/design.md** - Project-specific design guidelines
   - **CRITICAL**: Always read if it exists before making changes
   - Contains project conventions and patterns that override general best practices
