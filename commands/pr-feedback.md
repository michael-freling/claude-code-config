---
description: Address PR review feedback. Use when fixing PR comments, responding to code review, or implementing reviewer suggestions.
argument-hint: "[PR number] (optional, defaults to current branch's PR)"
allowed-tools: ["*"]
---

# PR Feedback

Address reviewer feedback on a pull request by reading comments, implementing fixes, and replying to reviewers.

$ARGUMENTS

## Workflow

### Phase 1: Identify PR and Gather Feedback

Determine the PR number and fetch all review feedback:

```bash
# Get PR number from argument or current branch
gh pr view --json number,headRefName,url

# Fetch all review data
gh pr view <PR> --json number,title,reviews,reviewRequests,comments,reviewDecision

# Get review comments (inline comments on code)
gh api repos/{owner}/{repo}/pulls/<PR>/comments

# Get PR reviews (overall review with body)
gh api repos/{owner}/{repo}/pulls/<PR>/reviews
```

Parse the feedback to understand:
- Review threads and conversations
- Inline comments on specific lines
- Overall review comments
- Requested changes vs suggestions

### Phase 2: Analyze and Categorize Feedback

Categorize each piece of feedback:

| Category | Description | Priority |
|----------|-------------|----------|
| **Must fix** | Bugs, security issues, correctness problems, requested changes | High |
| **Should fix** | Style issues, improvements, reasonable suggestions | Medium |
| **Discussion** | Questions, clarifications, design discussions needing user input | Requires user |

For each comment, determine:
- What file and line it references
- What change is being requested
- Whether it's actionable or needs discussion
- If related comments should be grouped together

Identify any conflicting feedback (two reviewers suggesting opposite things).

### Phase 3: Present Plan to User

Present the categorized feedback and proposed approach:

```
## Feedback Summary

### Must Fix (X items)
1. [file:line] Description of issue - Proposed fix
2. ...

### Should Fix (X items)
1. [file:line] Description - Proposed fix
2. ...

### Discussion Needed (X items)
1. [file:line] Question/concern - Options for user to decide
2. ...

### Conflicting Feedback (if any)
- Reviewer A says X, Reviewer B says Y - Need user decision
```

**Wait for user approval before implementing.** Resolve any discussions or conflicts first.

### Phase 4: Implement Changes

Address feedback in logical order (group related changes):

For each piece of feedback:
1. Read the relevant code for full context
2. Implement the requested change
3. Verify the fix is correct

Group related changes together rather than addressing comments in isolation.

### Phase 5: Reply to Comments

After implementing fixes, reply to each addressed comment:

```bash
# Reply to a review comment thread
gh api repos/{owner}/{repo}/pulls/<PR>/comments/<COMMENT_ID>/replies \
  -f body="Fixed in <commit>. <brief explanation of what was done>"
```

Reply format:
- Keep replies concise
- Reference the commit or change made
- Explain the approach if it differs from the suggestion
- For discussions, explain the decision made

### Phase 6: Commit and Push

Group related changes into logical commits:
- Each commit should address a coherent set of related feedback
- Reference feedback in commit messages where relevant
- Write clear commit messages describing the "why"

```bash
git add <files>
git commit -m "$(cat <<'EOF'
Address review feedback: <summary>

<details of changes>

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
git push
```

### Phase 7: Fix CI Until Green

Monitor and fix CI after pushing:

1. Wait for CI to start (at least 1 minute)
2. Check CI status:
   ```bash
   gh pr checks <PR> --json name,state,conclusion
   ```
3. If CI fails:
   - Analyze the errors
   - Fix issues and push new commits
   - Repeat until all checks pass
4. Check status every 5 minutes until complete

## Key Rules

1. **Always present plan first** - Get user approval before implementing changes
2. **Reply to individual comments** - Use threaded replies, not general PR comments
3. **Prioritize must-fix items** - Address critical issues before nice-to-have improvements
4. **Group related changes** - Commit related feedback fixes together
5. **Handle conflicts explicitly** - If reviewers disagree, ask user to decide
6. **Keep CI green** - Don't leave PR with failing checks

## Example Commands

```bash
# Get PR for current branch
gh pr view --json number,url

# Get all review comments
gh api repos/owner/repo/pulls/123/comments --jq '.[] | {id, path, line, body, user: .user.login}'

# Get reviews with their states
gh api repos/owner/repo/pulls/123/reviews --jq '.[] | {id, state, body, user: .user.login}'

# Reply to a comment
gh api repos/owner/repo/pulls/123/comments/456/replies -f body="Fixed - moved the validation logic as suggested."

# Check CI status
gh pr checks 123 --json name,state,conclusion
```
