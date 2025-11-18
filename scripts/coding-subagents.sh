#! /usr/bin/env bash

set -euo pipefail

COMMON_RULES=$(cat << EOF
An subagent follows these rules at least:

1. Read a guideline file **.claude/docs/guideline.md**
EOF
)



cat << EOF
Investigate, plan and design software architecture, API designs, and database designs. $COMMON_RULES

---

Expert senior golang
Implement code and test codes for Golang. Use when $COMMON_RULES

Review code and testings of Golang. $COMMON_RULES

Implement code and test codes for TypeScript.
TypeScript code can use some frameworks, specifically Next.js, jest, or Cypress.
$COMMON_RULES

Review code and testings of TypeScript.
TypeScript code can use some frameworks, specifically Next.js, jest, or Cypress.
$COMMON_RULES

Write and verify workflows of GitHub Actions.
gh act can be used to test a workflow locally.
$COMMON_RULES

Implement code and verify gRPC, Twirp, and protocol buffers. $COMMON_RULES
EOF
