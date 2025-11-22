#!/bin/bash
set -euo pipefail

# SessionEnd hook to clean up git worktrees created during the session
# This script reads worktree information from .claude/worktree-session.json
# and removes the worktree and its branch if they exist

# Read input from stdin (SessionEnd hook data)
INPUT=$(cat)

# Extract cwd from hook input
CWD=$(echo "$INPUT" | jq -r '.cwd // empty')
if [[ -z "$CWD" ]]; then
    echo "Error: Could not determine working directory from hook input" >&2
    exit 1
fi

# Path to worktree session file
WORKTREE_SESSION_FILE="$CWD/.claude/worktree-session.json"

# Check if worktree session file exists
if [[ ! -f "$WORKTREE_SESSION_FILE" ]]; then
    # No worktree was created in this session, nothing to clean up
    exit 0
fi

# Read worktree information
WORKTREE_PATH=$(jq -r '.worktree_path // empty' "$WORKTREE_SESSION_FILE")
BRANCH_NAME=$(jq -r '.branch_name // empty' "$WORKTREE_SESSION_FILE")
TICKET_NUMBER=$(jq -r '.ticket_number // empty' "$WORKTREE_SESSION_FILE")

if [[ -z "$WORKTREE_PATH" ]] || [[ -z "$BRANCH_NAME" ]]; then
    echo "Error: Invalid worktree session file" >&2
    exit 1
fi

echo "Cleaning up worktree for ticket #$TICKET_NUMBER..." >&2

# Change to the original repository directory
cd "$CWD" || exit 1

# Remove the worktree if it exists
if git worktree list | grep -q "$WORKTREE_PATH"; then
    echo "Removing worktree at $WORKTREE_PATH..." >&2
    git worktree remove "$WORKTREE_PATH" --force 2>/dev/null || {
        echo "Warning: Failed to remove worktree, attempting cleanup..." >&2
        # If worktree removal fails, try pruning and removing manually
        rm -rf "$WORKTREE_PATH" 2>/dev/null || true
        git worktree prune 2>/dev/null || true
    }
fi

# Delete the branch if it exists and has no commits beyond main
if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
    echo "Deleting branch $BRANCH_NAME..." >&2
    git branch -D "$BRANCH_NAME" 2>/dev/null || {
        echo "Warning: Failed to delete branch $BRANCH_NAME" >&2
    }
fi

# Remove the session file
rm -f "$WORKTREE_SESSION_FILE"

echo "Worktree cleanup completed" >&2

# Return success
exit 0
