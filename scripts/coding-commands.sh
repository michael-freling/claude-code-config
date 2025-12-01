#! /usr/env/bin bash

set -euo pipefail

COMMON_RULES=$(cat <<'EOF'
1. Create the slash commands of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/slash-commands
    1. Add **$ARGUMENTS** or **$1**, **$2** for each parameter
    2. Add **allowed-tools**, **argument-hint**, and **description** for each slash command at least.
3. Ask users to install all MCP servers that users should install on each slash command or skill.
EOF
)

cat <<EOF
Create a command to create new PR from another big PR which needs to be split.

The new PR should be consist of two parts of PRs:
1. A parent PR, which has to be created at first, and probably with an empty commit against the default branch.
2. Child PRs, that are against the parent PR. Include the link or PR number of the parent PR in the description of each child PR.

After all child PRs are created, the description of a parent PR should be updated to include all of PR numbers or the links of child PRs for document purposes.

Child PRs are split based on the meaningful parts of the original PR, which can be easily reviewed.
It might be split based on a set of some commits for reviewing.
But if the original commits contain a lot of noisy commits, each PR might be from a set of some files.

Follow the following Claude Code Command rules:
$COMMON_RULES
EOF
