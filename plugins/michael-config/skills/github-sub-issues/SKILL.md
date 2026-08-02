---
name: github-sub-issues
description: Create parent-child relationships between GitHub issues using the sub-issues feature. Use when organizing issues into hierarchies.
allowed-tools: Bash
---

# GitHub Sub-Issues Skill

Create parent-child relationships between GitHub issues using the sub-issues feature.

## Prerequisites

- GitHub CLI (`gh`) must be authenticated
- Both parent and child issues must exist in the same repository

## Adding a Sub-Issue

### Step 1: Get Issue Node IDs

```bash
# Get parent issue node ID
PARENT_ID=$(gh api graphql -f query='
  query($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
      issue(number: $number) { id }
    }
  }
' -F owner='{owner}' -F repo='{repo}' -F number={parent_number} --jq '.data.repository.issue.id')

# Get child issue node ID
CHILD_ID=$(gh api graphql -f query='
  query($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
      issue(number: $number) { id }
    }
  }
' -F owner='{owner}' -F repo='{repo}' -F number={child_number} --jq '.data.repository.issue.id')
```

### Step 2: Create Sub-Issue Relationship

```bash
gh api graphql \
  -H "GraphQL-Features: sub_issues" \
  -f query='
    mutation($parentId: ID!, $childId: ID!) {
      addSubIssue(input: {issueId: $parentId, subIssueId: $childId}) {
        issue { title }
        subIssue { title }
      }
    }
  ' -F parentId="$PARENT_ID" -F childId="$CHILD_ID"
```

## Removing a Sub-Issue Relationship

```bash
gh api graphql \
  -H "GraphQL-Features: sub_issues" \
  -f query='
    mutation($parentId: ID!, $childId: ID!) {
      removeSubIssue(input: {issueId: $parentId, subIssueId: $childId}) {
        issue { title }
        subIssue { title }
      }
    }
  ' -F parentId='{parent_node_id}' -F childId='{child_node_id}'
```

## Listing Sub-Issues

```bash
gh api graphql \
  -H "GraphQL-Features: sub_issues" \
  -f query='
    query($owner: String!, $repo: String!, $number: Int!) {
      repository(owner: $owner, name: $repo) {
        issue(number: $number) {
          title
          subIssues(first: 50) {
            nodes {
              number
              title
              state
            }
          }
        }
      }
    }
  ' -F owner='{owner}' -F repo='{repo}' -F number={parent_number}
```

## Extracting Repository Info

Get owner and repo from current git remote:

```bash
gh repo view --json owner,name --jq '"\(.owner.login)/\(.name)"'
```

## Error Handling

- **Issue not found**: Verify issue numbers exist in the repository
- **Permission denied**: Ensure `gh` is authenticated with write access
- **Already a sub-issue**: The child may already be linked to a parent
- **Circular reference**: An issue cannot be both parent and child of another
