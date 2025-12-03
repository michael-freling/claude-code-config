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

## TODO: update coding better
cat <<EOF
# Claude Code Slash Command and Skill Guidelines for Coding

This is coding rules to generate Claude Codes' commands and skills.

The commands and skills have to be created with following documents:

## General Claude Code's slash commands and skill guides

1. Create the slash commands of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/slash-commands
    1. Add `$ARGUMENTS` or `$1`, `$2` for each parameter
    2. Add `allowed-tools`, `argument-hint`, and `description` for each slash command at least.
2. Create the skills of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/skills
    1. Add `name`, `allowed-tools` and `description` frontmatter
    2. Add `examples.md` for complex examples
    3. Add `Version History`
3. Ask users to install all MCP servers that users should install on each slash command or skill.


### Slash command: document-guideline

Create a Claude Code's slash command `document-guideline` to write a document `.claude/docs/guideline.md` by following codes

1. Architect sub agent must analyze codebase and output a guideline
2. Guideline must include at least an architecture, API designs, data models, design patterns, coding conventions, and coding best practices.


### Slash command: document-guideline-monorepo

Create a Claude Code's slash command `document-guideline-monorepo` to write documents `.claude/docs/guideline.md` on each subproject of this monorepo by following codes

1. Architect sub agent must analyze codebase and output a guideline
2. Guideline must include at least an architecture, API designs, data models, design patterns, coding conventions, and coding best practices.
3. Write a guideline document on each subproject instead of a root project. The guideline on the root project must not include guideilne specific to each subproject



### Slash commands: feature, fix, and refactor

Create following Claude Code commands, to write an instruction to Claude Code for each workflow.
The overview of each command is following:

1. **feature** command: add or update a feature.
2. **fix** command: fix a bug. Before starting implementation, reproduce errors at first and understand the root cause of an error.
4. **refactor** command: refactor the codebase. Before doing the refactoring, cleaning unnecessary codes where agents will change.

**Fix command**
1. Analyze existing codebase, and understand where an error happens. Architect Agent must do this.
2. Reproduce errors and understand the root cause of an error. A software engineer with the right techstack must do this.
3. Based on the analysis and the root cause of an error, plan changes to make changes. Confirm the plan to the user to see if the plan is good. Do not start until you get an approval.
4. Once you get an approval, make sure you update the local default branch is the same as the remote one. If not, recreate the local default branch from remote one.
5. If there are multiple phases
    1. Create an epic PR to the default branch with an empty commit.
    2. For each phase, subagents must make each change in a worktree. New worktree msut be under `../worktrees`. And create a sub PR to the branch of the epic PR.
    3. Fix any CI errors until CI passes. To confirm the result of CI, wait for a long time because CI is slow. For example, wait for a minute to start a job, and wait for at least every 5 minutes to complete.
6. If there is a single phase subagents must make each change in a worktree. New worktree must be under `../worktrees`. Then create a PR and fix any CI errors until CI passes while waiting for a long time.


**feature and refactor commands**

1. Analyze existing codebase and architecture.
2. Create new design and plan changes for new feature or refactoring.
3. Get a review for the design and the plan.
4. Confirm the design, whether backward compatibility is required or not, and plan to the user to see if the plan is good. Do not start until you get an approval.
5. Once you get an approval, make sure you update a local main branch is the same as the remote main branch. If not, recreate the local main branch.
6. If there are multiple phases
    1. Create an epic PR to the default branch with an empty commit.
    2. For each phase, subagents must make each change in a worktree. New worktree msut be under `../worktrees`. And create a sub PR to the branch of the epic PR.
    3. Fix any CI errors until CI passes. To confirm the result of CI, wait for a long time because CI is slow. For example, wait for a minute to start a job, and wait for at least every 5 minutes to complete.
7. If there is a single phase subagents must make each change in a worktree. New worktree must be under `../worktrees`. Then create a PR and fix any CI errors until CI passes while waiting for a long time.


#### The details of all commands.
##### Parameters

Each slach command has parameters from a user, which describes what kind of changes the user wants to make. For slash command arguments:
- Only use $ARGUMENTS (or $1, $2 etc.) without additional explanation
- The argument-hint in frontmatter is sufficient to guide users
EOF
