---
name: git-branch
description: Create a git branch. Branches are only created when the user explicitly requests one. Use when the user asks to create or switch to a new branch.
allowed-tools: Bash
---

# Git Branch Skill

Create git branches, but only on explicit user request.

## Rules

- **NEVER create a new git branch unless the user explicitly asks for one.**
- **What counts as "explicit"**: the user says something like "create a branch", "make a branch", or "branch off". Anything ambiguous is NOT explicit.
- Clarify ambiguous requests with the user instead of assuming. Do not create a branch "to be helpful".
- If work is needed and no branch was requested, stay on the current branch (or ask the user which branch to use).

## Branch Creation Workflow

Only when the user explicitly asks for a new branch.

1. Check the current branch and working tree state:

```bash
git status
git branch --show-current
```

2. Do not branch off uncommitted, unrelated work. If there are unrelated changes, ask the user how to proceed (stash, commit, or discard).

3. Create the branch from an up-to-date `main`:

```bash
git checkout main
git pull
git checkout -b <branch-name>
```

4. Choose a sensible branch name: short, kebab-case, and descriptive of the work (e.g. `fix-login-redirect`). If an issue number is known, prefix it (e.g. `259-keyword-search`).
