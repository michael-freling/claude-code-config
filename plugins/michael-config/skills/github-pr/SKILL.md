---
name: github-pr
description: Create a GitHub pull request. PRs are only created when the user explicitly requests one. Use when the user asks to open a pull request or push work for review.
allowed-tools: Bash
---

# GitHub PR Skill

Create GitHub pull requests, but only on explicit user request.

## Rules

- **NEVER create a new GitHub PR unless the user explicitly asks for one.**
- **What counts as "explicit"**: the user says something like "make a PR", "open a pull request", or "push this for review". Anything ambiguous is NOT explicit.
- Clarify ambiguous requests with the user instead of assuming. Do not create a PR "to be helpful".

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
