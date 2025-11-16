# Claude Code Slash Command and Skill Guidelines for Coding

This is coding rules to generate Claude Codes' commands and skills.

The commands and skills have to be created with following documents:

## General Claude Code's slash commands and skill guides

1. Keep commands and skills instructions as simple, concise and brief as possible
2. Create the slash commands of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/slash-commands
    1. Add `$ARGUMENTS` or `$1`, `$2` for each parameter
    2. Add `allowed-tools`, `argument-hint`, and `description` for each slash command at least.
3. Create the skills of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/skills
4. Ask users to install all MCP servers that users should install on each slash command or skill.


### Slash command: document-guideline

Create a Claude Code's slash command `document-guideline` by following codes

1. Architect role agent must analyze codebase and output a guideline
2. Guideline must include at least an architecture, API designs, data models, design patterns, and coding best practices.
3. If a repository is a monorepo, output a guideline document on each subproject instead of a root project. The guideline on the root project must not include guideilne specific to subprojects


### Slash command: document-roles

Create a Claude Code command document-roles under `.claude/roles` to create following roles

1. Architect - design high level API designs, an architecture, and data models.
2. Software engineer - writes clean, efficient, and maintainable code and verify them. Based on languages and framework, there are each type of software engineer, including
    1. Golang's Software Engineer
    2. TypeScript's Software Engineer
    2. Next.js's Software Engineer
3. Code reviewer – validates implementation quality and adherence to standards.


### Slash command: create-coding-skills

Create following claude code skills by software engineer agents on each language/framework.

1. Golang
2. TypeScript - Prefer pnpm instead of npm for new project
3. Next.js (TypeScript framework)
4. Protocol buffers
5. GitHub Actions
6. Dockerfile, Containerfile, or Compose file
7. Bash

Each skill has to follow following rules
- Create the skills of Claude Code by analyzing the usages and best practices on the Internet at first, including https://code.claude.com/docs/en/skills
- Follow general coding guidelines and guidelines of corresponding language/framework below.
- Include examples on each guideline when it helps an agent to remove unambiguity.
- Review each skill and update based on feedback


### Slash commands: feature, fix, and refactor

Create following Claude Code commands, to write an instruction to Claude Code for each workflow.
Each slach command has parameters from a user, which describes what kind of changes the user wants to make. For slash command arguments:
- Only use $ARGUMENTS (or $1, $2 etc.) without additional explanation
- The argument-hint in frontmatter is sufficient to guide users

The details of each command is following:

1. **feature** command: add or update a feature.
2. **fix** command: fix a bug. Before starting implementation, reproduce errors at first and understand the root cause of an error.
4. **refactor** command: refactor the codebase. Before doing the refactoring, cleaning unnecessary codes where agents will change.

**Fix command**
1. Analyze existing codebase, and understand where an error happens. Architect Agent role must do this.
2. Reproduce errors and understand the root cause of an error. A software engineer with the right techstack must do this.
3. Based on the analysis and the root cause of an error, plan changes in order to make some changes in parallel later, and confirm the plan to the user to see if the plan is good. Architect Agent role must do this.
4. Start implementation. On each task, software engineer agents should implement by following order, on each language or framework with Claude Code Skills.
    a. Updating codes
    b. Verify changes with lint or tests
    c. Review changes by different software engineer agent
5. Reviewer agent finally reviews changes

**feature and refactor commands**
1. Analyze existing codebase, plan changes in order to make some changes in parallel later, and confirm the plan to the user to see if the plan is good. Architect Agent role must do this.
2. Start implementation. On each task, software engineer agents should implement by following order, on each language or framework with Claude Code Skills.
    a. Updating codes
    b. Verify changes with lint or tests
    c. Review changes by different software engineer agent
3. Reviewer agent finally reviews changes

#### Notes

Each agent can commit each change on behalf of the user and push it to a remote repository **ONLY WHEN** a user instructs to do so.
Also, the project may be a single project or monorepo.

## General coding guideline

First of all, follow best practices of coding and each language.
That includings following guideline:

1. Simplicity is the most important. This means
    1. Instead of creating new functions, always plan to update existing codes for reusability. Follow a DRY principal
    2. Prefer to break backward compatibility unless users explicitly mention
2. Ensure consistency by following a document created by document-guideline Claude Code commands
3. Ensure fail-fast instead of silently killing errors
4. Prefer using the latest versions to older versions of modules/packages. But if there is a compatibility problem to the dependency on the latest version, downgrade the major version to support the dependency.
5. Prefer not to make anything public/exported outside of packages for encapsulation.
6. Prefer to continue or return early than nesting codes as much as possible. "if is bad, else is worse


### Testing

1. Use table-driven testings. Split happy and error test sets if they are complicated
2. Do not add redundant, meaningless test case if the purpose is the same as other test cases.
4. Prefer injecting values to changing global states. For example, instead of updating environment variables, updating global variables, changing working directories, updating existing functions to allow a constructor or method injection.
3. Avoid defining arguments of each test as function arguments, instead, define test inputs as test case's fields


## Golang's guideline

1. Use `cobra` for a command line and `viper` for configuration management
2. For testing, use shuffle and race options
3. Prefer to have only one test function for each function and have a subtest inside the function.


## TypeScript guideline

1. Prefer `pnpm` to `npm`
4. MUST NOT use `legacy-peer-deps` option of npm


## Next.js guideline

1. Follow all guideline of TypeScript
2. Use a cypress for testing and include a guideline of Cypress


## GitHub Actions guideline

1. Identify if a repository includes a single project or it's in a monorepo. If it's a monorepo, follow a few rules:
    1. separate workflows per subproject and run workflows only when the project or its dependency is updated.
    2. create a reusable workflow, for example, each language, and use the workflow from each subproject instead of duplicating the same workflows on each project's workflow.
2. Use `gh act` to verify new changes locally


## Base guideline

1. Handle errors always. For example, add `set -euo pipefail` at the beginning of a script.
