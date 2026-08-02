---
name: github-pr-link
description: Link GitHub issues to pull requests using closing keywords. Use when creating or editing PRs to automatically link and close associated issues.
allowed-tools: Bash
---

# GitHub PR Link Skill

Link GitHub issues to pull requests using closing keywords in the PR body.

## How to Extract the Issue Number

1. **From branch name**: Extract the leading number from branch names like `259-keyword-search` → issue #259
2. **From commit messages**: Look for `#<number>` patterns in commit messages
3. **User-provided**: If the user explicitly mentions an issue number, use that

## Linking Keywords

GitHub recognizes these keywords to link issues to PRs:

- `Closes #<number>` - Links and auto-closes issue when PR merges
- `Fixes #<number>` - Links and auto-closes issue when PR merges
- `Resolves #<number>` - Links and auto-closes issue when PR merges

For cross-repository references, use: `Closes owner/repo#123`

Place the keyword at the start of the PR body for visibility.

## During PR Creation

Include the closing keyword in the PR body when creating:

```bash
gh pr create --title "#360 Feature title" --body "$(cat <<'EOF'
Closes #360

## Summary
- Description of changes

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

## After PR Creation

If the PR was already created without a closing keyword, update the body:

```bash
# Get current body
CURRENT_BODY=$(gh pr view <PR_NUMBER> --json body --jq '.body')

# Prepend closing keyword and update
gh pr edit <PR_NUMBER> --body "Closes #<ISSUE_NUMBER>

$CURRENT_BODY"
```

## Verifying the Link

After linking, verify the issue is connected:

```bash
gh pr view <PR_NUMBER> --json closingIssuesReferences --jq '.closingIssuesReferences[].number'
```

## Notes

- The issue number is typically the prefix of the branch name (before the first hyphen)
- If no issue number can be determined, ask the user which issue to link
- The `--add-issue` flag does not exist in `gh pr edit` - use body keywords instead
