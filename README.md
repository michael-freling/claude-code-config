# Claude Code Configuration

Personal Claude Code configurations including commands, skills, rules, and agents.

## Installation

```bash
git clone https://github.com/michael-freling/claude-code-plugins.git
cd claude-code-plugins
make install
```

This creates symlinks in `~/.claude/` pointing to this repository.

To uninstall:

```bash
make uninstall
```

## Repository Structure

This repository is a Claude Code plugin marketplace. The marketplace manifest
stays at the repo root, and the plugin itself lives under `plugins/michael-config/`:

```
.claude-plugin/
└── marketplace.json          # marketplace manifest (repo root)
plugins/
└── michael-config/
    ├── .claude-plugin/
    │   └── plugin.json        # plugin manifest
    ├── agents/
    ├── commands/
    ├── skills/
    ├── rules/
    └── templates/
```

## What's Included

### Commands

Slash commands for common workflows:

| Command | Purpose |
|---------|---------|
| `/implement-pr` | Implement code changes and create a PR, fixing CI errors until it passes |
| `/pr-summary` | Generate PR title and description from commits |
| `/write-product-spec` | Write product requirements and user stories |
| `/design-ux` | Create UI flows and wireframes |
| `/design-system-architecture` | Design API, database, and security |
| `/design-frontend-details` | Design frontend implementation |
| `/design-backend-details` | Design backend implementation |
| `/document-guideline` | Analyze codebase and create project guidelines |
| `/document-guideline-monorepo` | Create project guidelines for monorepo setups |
| `/pr-review` | Review a pull request for code quality |
| `/pr-feedback` | Address PR review feedback and fix issues |

### Skills

Reusable guidance for specific tasks:

| Skill | Purpose |
|-------|---------|
| `troubleshoot` | General diagnostic principles |
| `github-actions-troubleshoot` | Fix GitHub Actions CI errors |
| `kubernetes-troubleshoot` | Diagnose K8s pod failures, ArgoCD sync issues |
| `kubernetes-manifests` | K8s manifest best practices |
| `pr-feedback` | Address PR review feedback |
| `github-pr-link` | Link issues to PRs with closing keywords |
| `github-sub-issues` | Create parent-child issue relationships |
| `git-branch` | Create a git branch (only on explicit request) |
| `github-pr` | Create a GitHub PR (only on explicit request) |

### Rules

Coding guidelines applied to all projects:

| Rule | Purpose |
|------|---------|
| `coding-guidelines` | General coding principles (simplicity, readability, testing) |
| `design-workflow` | Design phase ordering conventions |
| `system-operations` | System operations guidelines |
| `github-actions` | CI workflow best practices (quiet flags) |
| `docker-commands` | Docker command best practices |
| `golang` | Go-specific guidelines |
| `typescript` | TypeScript-specific guidelines |
| `protobuf` | Protocol Buffers conventions |

### Agents

Specialized agents for different tasks:

**Design Agents:**
- `product-manager` - Product requirements
- `ux-designer` - UI flows and wireframes
- `software-architect` - System architecture
- `frontend-design-engineer` / `backend-design-engineer` - Implementation design

**Review Agents:**
- `architecture-reviewer` - Review system designs
- `product-spec-reviewer` - Review product specs
- `ux-design-reviewer` - Review UX designs
- `frontend-design-reviewer` / `backend-design-reviewer` - Review implementation designs
- `code-reviewer` - Review code changes

**Implementation Agents:**
- `code-developer` - Implement code changes
- `github-actions-workflow-engineer` - Create/fix GitHub Actions workflows
- `kubernetes-engineer` - Create/troubleshoot K8s resources

## Design Workflow

For larger features, follow this sequence:

```
1. /write-product-spec        → Product requirements (WHAT and WHY)
         ↓
2. /design-ux                 → User flows and wireframes
         ↓
3. /design-system-architecture → API, database, security design
         ↓
4. /design-frontend-details   → Frontend implementation design
   /design-backend-details    → Backend implementation design
         ↓
5. /implement-pr              → Implementation + PR
```

## Implementation Workflow

The `/implement-pr` command follows this workflow:

1. **Understand** - Clarify requirements, identify change type (feature/fix/refactor)
2. **Explore** - Analyze codebase, find root cause (for fixes), evaluate options
3. **Present Plan** - Show plan and wait for user approval
4. **Implement** - Write code and tests
5. **Verify** - Run tests, linters, manual verification
6. **Commit** - Group related changes into logical commits
7. **Create PR** - Push and create PR with summary and test plan
8. **Fix CI Until Green** - Monitor CI, fix failures, repeat until passing
