---
description: Generate PR title and description from commits compared to remote main branch
allowed-tools: ["Bash", "Read", "Grep", "Glob"]
---

# PR Summary

Update the current PR's title and description to reflect the latest changes.

## Workflow

### 1. Get Current Branch and PR

```bash
# Get current branch
git branch --show-current

# Get PR number for current branch
gh pr view --json number,title,body
```

### 2. Analyze Changes

Compare with remote main branch:

```bash
# Fetch latest from remote
git fetch origin main

# Get all commits in this branch
git log origin/main..HEAD --oneline

# Get the full diff
git diff origin/main...HEAD --stat
```

### 3. Generate Title and Description

Based on the analysis:

**Title format:**
- For single-purpose PRs: `<type>: <concise description>`
- Types: `feat`, `fix`, `refactor`, `docs`, `chore`, `test`

**Description format:**
```markdown
## Summary
<2-4 bullet points describing the main changes>
```

### 4. Update the PR

```bash
gh pr edit --title "new title" --body "$(cat <<'EOF'
## Summary
...
EOF
)"
```

## Guidelines

- Read all commits to understand the full scope of changes
- Group related changes together in the description
- Keep the title concise (under 72 characters)
- Focus on the "why" not just the "what"
- Include breaking changes prominently if any
