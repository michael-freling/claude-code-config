# Claude Code Configuration

Personal Claude Code configurations including commands, skills, rules, and agents.

## Installation

```bash
git clone https://github.com/michael-freling/claude-code-config.git
cd claude-code-config
make install
```

This creates symlinks in `~/.claude/` pointing to this repository.

To uninstall:

```bash
make uninstall
```

## What's Included

### Commands

Slash commands for common workflows:

| Command | Purpose |
|---------|---------|
| `/implement` | Implement features, bug fixes, or refactoring with structured workflow |
| `/pr-summary` | Generate PR title and description from commits |
| `/write-product-spec` | Write product requirements and user stories |
| `/design-ux` | Create UI flows and wireframes |
| `/design-system-architecture` | Design API, database, and security |
| `/design-frontend-details` | Design frontend implementation |
| `/design-backend-details` | Design backend implementation |
| `/document-guideline` | Analyze codebase and create project guidelines |

### Skills

Reusable guidance for specific tasks:

| Skill | Purpose |
|-------|---------|
| `write-tests` | Testing patterns for Go and TypeScript |
| `troubleshoot` | General diagnostic principles |
| `github-actions-troubleshoot` | Fix GitHub Actions CI errors |
| `kubernetes-troubleshoot` | Diagnose K8s pod failures, ArgoCD sync issues |
| `kubernetes-manifests` | K8s manifest best practices |
| `pr-feedback` | Address PR review feedback |
| `github-pr-link` | Link issues to PRs with closing keywords |
| `github-sub-issues` | Create parent-child issue relationships |

### Rules

Coding guidelines applied to all projects:

| Rule | Purpose |
|------|---------|
| `coding-guidelines` | General coding principles (simplicity, readability, testing) |
| `simplicity` | Simplicity principle details |
| `design-workflow` | Design phase ordering conventions |
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
5. /implement                 → Implementation
```

## Implementation Workflow

The `/implement` command follows this workflow:

1. **Understand** - Clarify requirements, identify change type (feature/fix/refactor)
2. **Explore** - Analyze codebase, find root cause (for fixes), evaluate options
3. **Present Plan** - Show plan and wait for user approval
4. **Implement** - Write code and tests
5. **Verify** - Run tests, linters, manual verification
6. **PR Workflow** (if requested) - Commit, create PR, verify CI, fix until CI passes
