#!/bin/bash

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command' 2>/dev/null)

# Block --no-verify flag on git commit
if [[ "$command" == *"--no-verify"* ]]; then
  echo "ERROR: Blocked - --no-verify flag is not allowed" >&2
  exit 2
fi

# Block git push to default branches (main, master)
if [[ "$command" =~ ^git[[:space:]]+push ]]; then
  current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
  default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')

  # Fallback to common default branch names if detection fails
  if [[ -z "$default_branch" ]]; then
    default_branch="main"
  fi

  # Check if explicitly pushing TO main/master branch
  # Matches: "git push origin main", "git push origin master", "git push origin main:main"
  if [[ "$command" =~ [[:space:]]origin[[:space:]]+(main|master)[[:space:]]*$ ]] || \
     [[ "$command" =~ [[:space:]]origin[[:space:]]+(main|master): ]] || \
     [[ "$command" =~ :(main|master)[[:space:]]*$ ]]; then
    echo "ERROR: Blocked - Direct push to default branch is not allowed" >&2
    exit 2
  fi

  # Check if implicitly pushing current branch (which is main/master)
  # Matches: "git push", "git push origin", "git push -u origin" when on main/master
  # But NOT: "git push origin feature-branch" (explicit different branch)
  if [[ "$current_branch" == "main" || "$current_branch" == "master" || "$current_branch" == "$default_branch" ]]; then
    # Check if command specifies a different target branch explicitly
    # If it's just "git push" or "git push origin" or "git push -u origin", block it
    if [[ "$command" =~ ^git[[:space:]]+push[[:space:]]*$ ]] || \
       [[ "$command" =~ ^git[[:space:]]+push[[:space:]]+-[a-zA-Z]+[[:space:]]*$ ]] || \
       [[ "$command" =~ ^git[[:space:]]+push[[:space:]]+origin[[:space:]]*$ ]] || \
       [[ "$command" =~ ^git[[:space:]]+push[[:space:]]+-[a-zA-Z]+[[:space:]]+origin[[:space:]]*$ ]] || \
       [[ "$command" =~ ^git[[:space:]]+push[[:space:]]+origin[[:space:]]+-[a-zA-Z]+[[:space:]]*$ ]]; then
      echo "ERROR: Blocked - Direct push to default branch ($current_branch) is not allowed" >&2
      exit 2
    fi
  fi
fi

# Block modification of GitHub branch protection (DELETE, PUT, PATCH)
# This covers all protection endpoints including:
#   /branches/{branch}/protection
#   /branches/{branch}/protection/enforce_admins
#   /branches/{branch}/protection/required_pull_request_reviews
#   /branches/{branch}/protection/required_status_checks
#   /branches/{branch}/protection/restrictions
if [[ "$command" =~ gh[[:space:]]+api ]] && [[ "$command" =~ branches.*/protection ]]; then
  # Block DELETE (removes protection)
  if [[ "$command" =~ -X[[:space:]]+DELETE ]] || [[ "$command" =~ --method[[:space:]]+DELETE ]]; then
    echo "ERROR: Blocked - Removing branch protection is not allowed" >&2
    exit 2
  fi
  # Block PUT (replaces/weakens protection)
  if [[ "$command" =~ -X[[:space:]]+PUT ]] || [[ "$command" =~ --method[[:space:]]+PUT ]]; then
    echo "ERROR: Blocked - Modifying branch protection via PUT is not allowed" >&2
    exit 2
  fi
  # Block PATCH (updates/weakens specific protection settings)
  if [[ "$command" =~ -X[[:space:]]+PATCH ]] || [[ "$command" =~ --method[[:space:]]+PATCH ]]; then
    echo "ERROR: Blocked - Modifying branch protection via PATCH is not allowed" >&2
    exit 2
  fi
fi

# Block modification of GitHub rulesets (newer protection feature)
# Endpoints: /repos/{owner}/{repo}/rulesets and /repos/{owner}/{repo}/rulesets/{id}
if [[ "$command" =~ gh[[:space:]]+api ]] && [[ "$command" =~ repos/[^/]+/[^/]+/rulesets ]]; then
  # Block DELETE (removes ruleset)
  if [[ "$command" =~ -X[[:space:]]+DELETE ]] || [[ "$command" =~ --method[[:space:]]+DELETE ]]; then
    echo "ERROR: Blocked - Deleting repository rulesets is not allowed" >&2
    exit 2
  fi
  # Block PUT (replaces/weakens ruleset)
  if [[ "$command" =~ -X[[:space:]]+PUT ]] || [[ "$command" =~ --method[[:space:]]+PUT ]]; then
    echo "ERROR: Blocked - Modifying repository rulesets via PUT is not allowed" >&2
    exit 2
  fi
fi

exit 0
