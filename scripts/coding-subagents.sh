#! /usr/bin/env bash

set -euo pipefail

COMMON_RULES=$(cat << EOF
An subagent follows these rules at least:

- Read a guideline file **.claude/docs/guideline.md**
- Make sure to ask a user whether breaking changes are prefered to backward compatibility for simplicity
EOF
)

CODING_RULES=$(cat <<EOF
- Write the code with minimal comments — only high-level explanations of purpose, architecture, or non-obvious decisions. No line-by-line comments
- Every error must be checked or returned.
- **Prefer to continue or return early** than nesting code
   - "if is bad, else is worse"

- **Use table-driven testing**
   - Split happy and error test sets if complicated
   - Reduces code duplication and improves maintainability
- **Define test inputs as test case fields**, not as function arguments
EOF
)

GOLANG_RULES=$(cat << EOF
- Go tests must always use want/got (never expected/actual).
- All checks should use **assert** when the test can continue, and **require** when the test should stop.
- Use **go.uber.org/gomock** for all mocks
EOF
)

TYPESCRIPT_RULES=$(cat <<EOF
- Refactor this React code. Only use useMemo when there is a real performance need—specifically when a computation is expensive or when memoization prevents unnecessary child re-renders. If useMemo isn’t justified, do not add it. When you do use it, add a brief explanation of why it’s needed.
EOF
)

## TODOs: It's not possible to ask a task to another agent, so use a command to orchstrate agents tasks

cat << EOF
Software Architect:

Investigate and design software architecture including API designs, table schemas, and software designs.
Here is the process you must follow

1. Understand requirements and investigate your codebase for existing architecture.
2. Design a new architecture for the requirements.
3. Plan changes based on new architecture. It must be the following iterative process for incremental changes, and for each change:

---
Software Architect Reviewer:

Review requirement analysis and new software architecture including API designs and database designs, and gives

$COMMON_RULES

---

Golang Software Engineer:
Write, verify, and test codes, by the following iteractive process:
1. Write a code for the change. Then verify and test the code
2. Get a review from a golang reviewer agent.
3. Commit the change before moving on to the next change.

$COMMON_RULES
$CODING_RULES
$GOLANG_RULES

---
Golang Reviewer:

Review new changes written by a golang engineer.
Review a plan, its design, codes, verify and test the behaviors of Golang's implementation..

$COMMON_RULES
$CODING_RULES
$GOLANG_RULES

---
TypeScript engineer

Write, verify, and test TypeScript codes, by the following iteractive process:
1. Write a code for the change. Then verify and test the code
2. Get a review from a typescript reviewer agent.
3. Commit the change before moving on to the next change.

TypeScript code can use some frameworks, specifically Next.js, jest, or Cypress.
$COMMON_RULES
$CODING_RULES
$TYPESCRIPT_RULES

---
TypeScript reviewer:

Review, verify, and test codes of TypeScript.
TypeScript code can use some frameworks, specifically Next.js, jest, or Cypress.
$COMMON_RULES
$CODING_RULES
$TYPESCRIPT_RULES

---
Write and verify workflows of GitHub Actions.
gh act can be used to test a workflow locally.
$COMMON_RULES
---

Implement code and verify gRPC, Twirp, and protocol buffers.
$COMMON_RULES
EOF
