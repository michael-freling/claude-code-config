---
name: git-branch-pr
description: Create a git branch and GitHub PR. Branches and PRs are only created when the user explicitly requests them. Use when the user asks to create a branch, open a pull request, or push work for review.
allowed-tools: Bash
---

# Git Branch and PR Skill

Create git branches and GitHub pull requests, but only on explicit user request.

## Rules

- **NEVER create a new git branch unless the user explicitly asks for one.**
- **NEVER create a new GitHub PR unless the user explicitly asks for one.**
- If work is needed and no branch was requested, stay on the current branch (or ask the user which branch to use).
- **What counts as "explicit"**: the user says something like "create a branch", "make a PR", "open a pull request", "push this for review". Anything ambiguous is NOT explicit.
- Clarify ambiguous requests with the user instead of assuming. Do not create a branch or PR "to be helpful".

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

## PR Creation Workflow

Only when the user explicitly asks for a new PR.

1. Ensure commits are pushed to the remote:

```bash
git push -u origin <branch-name>
```

2. Create the PR with a clear title and body using `gh`:

```bash
gh pr create --title "Title of the change" --body "$(cat <<'EOF'
## Summary
- Description of changes

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

3. Always include the `🤖 Generated with [Claude Code](https://claude.com/claude-code)` trailer in the PR body.

4. To link and auto-close associated issues, use the github-pr-link skill.
