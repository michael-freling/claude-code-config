#!/bin/bash

input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command' 2>/dev/null)

# Block --no-verify flag
if [[ "$command" == *"--no-verify"* ]]; then
  echo "ERROR: Blocked - --no-verify flag is not allowed" >&2
  exit 2
fi

exit 0
