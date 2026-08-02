---
name: pr-feedback
description: Address PR review feedback. Use when fixing PR comments, responding to code review, or implementing reviewer suggestions.
---

# PR Feedback

Address PR feedback following these guidelines.

## Never Postpone Refactoring

**NEVER postpone refactoring from PR feedback for new code.** When a reviewer requests refactoring on new functions or features, implement it in the same PR. This includes:

- Extracting code into separate services/modules
- Renaming functions or variables
- Restructuring code organization
- Adding abstractions or interfaces

Exception: Only postpone if the refactoring affects existing code that's not part of the current PR and would require significant changes to unrelated areas.

## Follow-up Issues

**NEVER create follow-up issues unless explicitly requested by the user.** Address all PR feedback in the current PR. If a task is genuinely out of scope:

1. Ask the reviewer before creating a follow-up issue
2. Get explicit approval before deferring work

## Code Reuse

When adding new functionality that's similar to existing code:

1. Investigate what can be shared before implementation
2. Extract common patterns into shared utilities
3. Document differences if code can't be shared

## Claude Code Configuration

**NEVER postpone improvements to Claude Code configurations** (rules, agents, commands) when requested in PR feedback. Implement them in the same PR. This includes:

- Adding new rules or agents
- Updating existing rules or agents
- Fixing or improving configuration files

These configurations help Claude Code work better on this project, so deferring them reduces productivity on subsequent tasks.
