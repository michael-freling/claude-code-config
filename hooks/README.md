# Hooks

This directory contains Claude Code hooks that execute in response to various events.

## cleanup-worktree.sh

**Event**: SessionEnd

**Purpose**: Automatically clean up git worktrees created during a Claude Code session when the session ends.

**How it works**:
1. When a session ends, the hook reads `.claude/worktree-session.json` in the project directory
2. If a worktree was created during the session, it removes:
   - The worktree directory (e.g., `../worktrees/issue-123`)
   - The associated git branch
3. Cleans up the session tracking file

**Dependencies**:
- `jq` - Required for parsing JSON session data
- Git worktree support

**Configuration**:
This hook is configured in `settings.json` under the `hooks.SessionEnd` section.

## Usage

The cleanup hook works automatically with the `/worktree` command. No manual intervention is needed.

When you create a worktree using `/worktree <ticket-number>`, the session information is stored and automatically cleaned up when the Claude Code session ends.
