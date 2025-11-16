---
name: github-actions
description: GitHub Actions workflow development following industry best practices for CI/CD, monorepo handling, security, reusable workflows, and local testing with gh act. Use when adding or updating GitHub Actions workflows.
---

# GitHub Actions Development Skill

A comprehensive skill for creating and maintaining GitHub Actions workflows with emphasis on industry best practices, security, efficiency, and local testing using `gh act`.

## When to Use

Use this skill when:
- Creating new GitHub Actions workflows
- Updating or debugging existing workflows
- Setting up CI/CD pipelines
- Automating repository tasks
- Testing workflows locally before pushing

## MCP Server Setup (Recommended)

For enhanced GitHub Actions development with real-time workflow information, install the GitHub MCP server:

### Installation

Add to your Claude Code settings (`.claude/settings.json` or global settings):

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-github"
      ],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_github_token_here"
      }
    }
  }
}
```

### Getting a GitHub Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Give it a descriptive name (e.g., "Claude Code MCP")
4. Select scopes:
   - `repo` - Full control of private repositories
   - `workflow` - Update GitHub Actions workflows
   - `read:org` - Read org and team membership
5. Click "Generate token"
6. Copy the token and add it to your settings

### What the MCP Server Provides

With the GitHub MCP server installed, Claude Code can:

- **Read workflow files** directly from repositories
- **View workflow runs** and their status
- **Check workflow logs** for debugging
- **List repository workflows** to understand CI/CD setup
- **Access repository information** (branches, commits, PRs)
- **Read action definitions** from the marketplace

### Usage with This Skill

When the GitHub MCP server is available, this skill can:

1. **Analyze existing workflows** more thoroughly
   ```
   Can you analyze our CI workflow and suggest improvements?
   ```

2. **Debug workflow failures** with access to logs
   ```
   The deploy workflow is failing, can you check what's wrong?
   ```

3. **Review workflow runs** across branches
   ```
   Show me all failed workflow runs from the last week
   ```

4. **Compare workflows** across repositories
   ```
   How does our CI compare to similar projects?
   ```

### Example Workflow

With MCP server installed:

```
You: "Create a new CI workflow for this Node.js project"

Claude Code (with MCP):
1. Checks existing workflows in the repository
2. Reviews recent workflow runs to understand patterns
3. Creates new workflow following your conventions
4. Tests locally with gh act
5. Suggests improvements based on workflow history
```

Without MCP server:
- Still creates excellent workflows
- Uses best practices and this skill's guidance
- Requires you to provide more context about existing workflows

### Optional but Recommended

While the MCP server enhances the experience, this skill works perfectly without it. The core workflow (design, implement, test with `gh act`, push) remains the same.

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation:

1. **Look for `.claude/design.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/design.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no design.md exists:**
   - Consider running `/document-design` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from design.md:**
- Project-specific CI/CD requirements
- Testing and deployment conventions
- Security and permissions policies
- Branch protection rules
- Code examples showing preferred workflow style

### 1. Analyze Project Structure

- Check for existing workflows in `.github/workflows/`
- **Identify if the repository is a monorepo** (multiple subprojects in one repo)
  - Look for directories like `services/`, `apps/`, `packages/`, `go/`, `node/`
  - Check for workspace files: `go.work`, `pnpm-workspace.yaml`, `lerna.json`
  - If monorepo: Plan separate workflows per subproject with path filters
- Identify project type and technology stack
- Review existing workflow patterns
- Check for secrets and environment variables
- Identify testing framework and build tools

### 2. Search for Relevant Code

- Search for similar workflows in the repository
- Look for reusable actions or composite actions
- Identify common patterns (checkout, caching, testing, deployment)
- Review workflow triggers and conditions

## GitHub Actions Best Practices

### Core Principles

**CRITICAL: Follow these core principles in all workflows**

1. **Simplicity First**: Keep workflows simple and easy to understand. Complex workflows are harder to maintain and debug.
2. **DRY (Don't Repeat Yourself)**: Use reusable workflows and composite actions to avoid duplication. Build a library of reusable components.
3. **Consistency**: Follow established patterns within the repository. Analyze existing workflows and match their style.
4. **Fail-Fast**: Configure workflows to fail immediately on errors. Use appropriate timeouts and error handling.
5. **Latest Versions**: Always use the latest stable versions of actions (e.g., `actions/checkout@v4`, not `@v2`). Pin to major versions (e.g., `@v4`) rather than specific commits for easier maintenance.
6. **Encapsulation**: Limit what workflows expose. Use appropriate permissions and minimize secrets exposure.
7. **Early Returns**: Use path filters and conditions to skip unnecessary work early in the workflow.

### Workflow Structure

- Use clear, descriptive workflow names
- Organize workflows by purpose (CI, CD, release, automation)
- Use consistent naming conventions
- Keep workflows focused and modular
- Use reusable workflows for common patterns (up to 10 nested levels, 50 total workflows per run as of Nov 2025)

Example:
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm test
```

### Testing Workflows Locally with `gh act`

**CRITICAL: Always test workflows locally with `gh act` before pushing to GitHub.**

`gh act` is a GitHub CLI extension that runs your workflows locally using Docker containers, allowing you to catch errors early and iterate quickly.

**Installation:**
```bash
# Install gh act extension
gh extension install https://github.com/nektos/gh-act
```

**Basic Usage:**

```bash
# List all workflows and their jobs
gh act -l

# Run all workflows (default: push event)
gh act

# Run specific event
gh act pull_request

# Run specific job
gh act -j test

# Run with specific workflow file
gh act -W .github/workflows/ci.yml

# Dry run (show what would be executed)
gh act -n

# Run with verbose output
gh act -v

# Run specific job with secrets
gh act -j test -s GITHUB_TOKEN=your_token
```

**Advanced Testing:**

```bash
# Test with custom event payload
gh act push --eventpath path/to/event.json

# Test with environment variables
gh act -j test --env MY_VAR=value

# Use custom runner image
gh act -P ubuntu-latest=catthehacker/ubuntu:act-latest

# Test matrix builds
gh act -j test --matrix node-version:18
```

**Common gh act Patterns:**

```bash
# Test before pushing
gh act -j test && git push

# Test all workflows on pull_request event
gh act pull_request

# Test with GitHub token for private repos
gh act -s GITHUB_TOKEN="$(gh auth token)"

# Test workflow with reusable workflows
gh act -j build --remote-name origin
```

### Workflow Triggers

- Use appropriate events (push, pull_request, schedule, workflow_dispatch)
- Add path filters to reduce unnecessary runs
- Use workflow_dispatch for manual triggers
- Consider using concurrency to cancel redundant runs

Example:
```yaml
on:
  push:
    branches: [main]
    paths:
      - 'src/**'
      - 'package.json'
  pull_request:
    branches: [main]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to deploy to'
        required: true
        type: choice
        options:
          - development
          - staging
          - production

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
```

### Job Design

- Use descriptive job names
- Set appropriate timeouts
- Use job dependencies with `needs`
- Run independent jobs in parallel
- Use matrix strategy for multi-version testing

Example:
```yaml
jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      - run: npm ci
      - run: npm run lint

  test:
    name: Test (Node ${{ matrix.node-version }})
    runs-on: ubuntu-latest
    timeout-minutes: 15
    strategy:
      matrix:
        node-version: [18, 20, 22]
      fail-fast: false
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}
          cache: 'npm'
      - run: npm ci
      - run: npm test

  build:
    name: Build
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      - run: npm ci
      - run: npm run build
      - uses: actions/upload-artifact@v4
        with:
          name: build-output
          path: dist/
```

### Caching

- Cache dependencies to speed up workflows
- Use appropriate cache keys
- Consider cache size limits
- Clear old caches periodically

Example:
```yaml
# Built-in caching with setup actions
- uses: actions/setup-node@v4
  with:
    node-version: '20'
    cache: 'npm'

# Manual caching
- name: Cache dependencies
  uses: actions/cache@v4
  with:
    path: ~/.npm
    key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-node-
```

### Security Best Practices

**CRITICAL: Security is non-negotiable. Follow these practices strictly.**

#### GITHUB_TOKEN Permissions (Least Privilege)

**Always set minimum required permissions for GITHUB_TOKEN:**

```yaml
# Repository/Organization Level Setting:
# Settings -> Actions -> Workflow permissions -> "Read repository contents permission"

# Workflow/Job Level (PREFERRED):
permissions:
  contents: read      # Only what's needed
  pull-requests: write  # If creating/updating PRs
  # Never grant more than necessary

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    permissions:
      contents: read
      deployments: write
    environment:
      name: production
      url: https://example.com
    steps:
      - uses: actions/checkout@v4

      - name: Deploy
        env:
          DEPLOY_TOKEN: ${{ secrets.DEPLOY_TOKEN }}
        run: |
          echo "Deploying to production"
          # Never echo secrets!
```

**Key Security Rules:**

1. **Pinned Action Versions**: Use specific version tags (e.g., `@v4`) to protect against supply-chain attacks
   ```yaml
   # Good: Pinned to major version
   - uses: actions/checkout@v4

   # Better: Pinned to specific version
   - uses: actions/checkout@v4.1.1

   # Best for security-critical: SHA pin (harder to maintain)
   - uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11
   ```

2. **Never Expose Secrets**: Never log, echo, or expose secrets in workflow output
   ```yaml
   # Bad: Exposes secret
   - run: echo "Token is ${{ secrets.API_TOKEN }}"

   # Good: Use secrets only where needed
   - env:
       API_TOKEN: ${{ secrets.API_TOKEN }}
     run: ./deploy.sh  # Script uses $API_TOKEN internally
   ```

3. **Least Privilege GITHUB_TOKEN**: Set minimal permissions at workflow or job level
   ```yaml
   # Workflow level (applies to all jobs)
   permissions:
     contents: read

   # Job level (override for specific job)
   jobs:
     deploy:
       permissions:
         contents: read
         deployments: write
   ```

4. **Secure Secrets Management**:
   - Store all sensitive data in GitHub Secrets (Settings -> Secrets and variables -> Actions)
   - Use environment-specific secrets for different deployment targets
   - Rotate secrets regularly
   - Never commit secrets to the repository

5. **Use Environments for Deployment Protection**:
   - Configure environment protection rules (required reviewers, wait timers)
   - Set environment-specific secrets
   - Use branch protection rules

6. **OIDC for Cloud Authentication** (Preferred over long-lived tokens):
   ```yaml
   # Use OpenID Connect for AWS/Azure/GCP authentication
   permissions:
     id-token: write  # Required for OIDC
     contents: read

   - name: Configure AWS credentials
     uses: aws-actions/configure-aws-credentials@v4
     with:
       role-to-assume: arn:aws:iam::123456789012:role/GitHubActionsRole
       aws-region: us-east-1
   ```

7. **Validate Third-Party Actions**:
   - Review source code of third-party actions before use
   - Prefer well-maintained, popular actions
   - Consider security scanning tools for actions

### Reusable Workflows

**CRITICAL: Use reusable workflows to follow DRY principle and standardize CI/CD patterns.**

As of November 2025, GitHub Actions supports:
- Up to **10 nested reusable workflows** (increased from previous limits)
- Up to **50 total workflows** called in a single workflow run

Create reusable workflows for common patterns:

```yaml
# .github/workflows/reusable-test.yml
name: Reusable Test Workflow

on:
  workflow_call:
    inputs:
      node-version:
        required: true
        type: string
        description: 'Node.js version to use'
      working-directory:
        required: false
        type: string
        default: '.'
        description: 'Working directory for commands'
    secrets:
      npm-token:
        required: false
        description: 'NPM authentication token'
    outputs:
      test-result:
        description: 'Test execution result'
        value: ${{ jobs.test.outputs.result }}

jobs:
  test:
    runs-on: ubuntu-latest
    outputs:
      result: ${{ steps.test.outputs.result }}
    defaults:
      run:
        working-directory: ${{ inputs.working-directory }}
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        id: test
        run: |
          npm test
          echo "result=success" >> $GITHUB_OUTPUT
```

Use it in other workflows:
```yaml
# .github/workflows/ci.yml
jobs:
  test:
    uses: ./.github/workflows/reusable-test.yml
    with:
      node-version: '20'
      working-directory: 'apps/web'
    secrets:
      npm-token: ${{ secrets.NPM_TOKEN }}
```

**Best Practices for Reusable Workflows:**

1. **Centralize in Monorepo**: Store reusable workflows in `.github/workflows/` and reference from subproject workflows
2. **Use Semantic Versioning**: Tag reusable workflows when storing in separate repositories
3. **Document Inputs/Outputs**: Provide clear descriptions for all inputs, secrets, and outputs
4. **Pass Environment Variables as Inputs**: Environment variables from caller workflows don't automatically propagate
   ```yaml
   # Caller workflow must explicitly pass values
   with:
     environment: ${{ vars.ENVIRONMENT }}
     region: ${{ vars.AWS_REGION }}
   ```
5. **Use Outputs for Communication**: Return results from reusable workflows to caller workflows

### Composite Actions

Create custom composite actions for repeated step sequences:

```yaml
# .github/actions/setup-node-app/action.yml
name: 'Setup Node.js App'
description: 'Setup Node.js with caching and install dependencies'
inputs:
  node-version:
    description: 'Node.js version'
    required: false
    default: '20'
runs:
  using: 'composite'
  steps:
    - uses: actions/setup-node@v4
      with:
        node-version: ${{ inputs.node-version }}
        cache: 'npm'
    - shell: bash
      run: npm ci
```

Use it:
```yaml
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/setup-node-app
    with:
      node-version: '20'
  - run: npm test
```

### Monorepo Workflow Organization

**CRITICAL: For monorepo projects, organize workflows by subproject and use path filters to minimize unnecessary CI runs.**

When working with monorepos containing multiple subprojects, follow these patterns inspired by industry best practices:

**Key Benefits of Proper Monorepo Organization:**
- **Faster CI/CD**: Only affected subprojects run their workflows
- **Resource Efficiency**: Reduced compute time and costs
- **Better Developer Experience**: Clearer feedback and faster iteration
- **Easier Maintenance**: Focused workflows are simpler to understand and update

#### 1. Detect Monorepo Structure

First, analyze the repository structure to identify if it's a monorepo:

```bash
# Common monorepo patterns
# - Language-specific directories: go/, node/, python/
# - Service directories: services/api/, services/web/
# - App directories: apps/frontend/, apps/backend/
# - Package directories: packages/shared/, packages/ui/
```

#### 2. Create Separate Workflow Files per Subproject

**Pattern: One workflow file per subproject**

Instead of a single monolithic workflow, create focused workflows:

```
.github/
  workflows/
    ci-api.yml          # API service CI
    ci-web.yml          # Web app CI
    ci-shared.yml       # Shared packages CI
    ci-proto.yml        # Proto/schema CI (if applicable)
    cd-api.yml          # API deployment
    cd-web.yml          # Web deployment
```

**Benefits:**
- Easier to understand and maintain
- Faster feedback (only relevant CI runs)
- Clearer failure attribution
- Independent versioning and deployment

#### 3. Use Path Filters to Trigger Only on Relevant Changes

**CRITICAL: Always use path filters to prevent unnecessary workflow runs.**

**Advanced Path Filtering with dorny/paths-filter:**

For more sophisticated change detection (especially useful for matrix strategies), use the `dorny/paths-filter` action:

```yaml
name: CI - Monorepo Smart Detection

on: [push, pull_request]

jobs:
  changes:
    runs-on: ubuntu-latest
    outputs:
      api: ${{ steps.filter.outputs.api }}
      web: ${{ steps.filter.outputs.web }}
      mobile: ${{ steps.filter.outputs.mobile }}
    steps:
      - uses: actions/checkout@v4

      - uses: dorny/paths-filter@v3
        id: filter
        with:
          filters: |
            api:
              - 'services/api/**'
              - 'proto/**'
            web:
              - 'apps/web/**'
              - 'packages/**'
            mobile:
              - 'apps/mobile/**'
              - 'packages/**'

  test-api:
    needs: changes
    if: needs.changes.outputs.api == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: cd services/api && go test ./...

  test-web:
    needs: changes
    if: needs.changes.outputs.web == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: cd apps/web && npm test
```

**Important Notes:**
- The `dorny/paths-filter` action is limited to 300 files in diffs
- For larger changes, consider more specific filters or multiple filter jobs
- Use this approach when you need job-level or step-level filtering (built-in path filters only work at workflow level)

**Basic Path Filtering Example for Go API service:**

```yaml
# .github/workflows/ci-api.yml
name: CI - API Service

on:
  push:
    branches: [main, develop]
    paths:
      - 'services/api/**'           # API source code
      - 'proto/**'                  # Shared proto files
      - 'go.work'                   # Go workspace file
      - 'go.work.sum'               # Go workspace sum
      - '.github/workflows/ci-api.yml'  # This workflow file
  pull_request:
    branches: [main]
    paths:
      - 'services/api/**'
      - 'proto/**'
      - 'go.work'
      - 'go.work.sum'
      - '.github/workflows/ci-api.yml'

jobs:
  test:
    name: Test API
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: services/api
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache-dependency-path: services/api/go.sum

      - name: Run tests
        run: go test ./...

      - name: Run linter
        run: golangci-lint run

      - name: Build
        run: go build ./...
```

Example for Next.js web app:

```yaml
# .github/workflows/ci-web.yml
name: CI - Web App

on:
  push:
    branches: [main, develop]
    paths:
      - 'apps/web/**'
      - 'packages/shared/**'        # Shared packages
      - 'package.json'              # Root package.json
      - 'pnpm-lock.yaml'            # Lock file
      - '.github/workflows/ci-web.yml'
  pull_request:
    branches: [main]
    paths:
      - 'apps/web/**'
      - 'packages/shared/**'
      - 'package.json'
      - 'pnpm-lock.yaml'
      - '.github/workflows/ci-web.yml'

jobs:
  test:
    name: Test Web App
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: apps/web
    steps:
      - uses: actions/checkout@v4

      - uses: pnpm/action-setup@v2
        with:
          version: 8

      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Run tests
        run: pnpm test

      - name: Run linter
        run: pnpm lint

      - name: Build
        run: pnpm build
```

#### 4. Handle Cross-Subproject Dependencies

When changes affect multiple subprojects (e.g., proto changes):

```yaml
# .github/workflows/ci-proto.yml
name: CI - Proto & Generated Code

on:
  push:
    branches: [main, develop]
    paths:
      - 'proto/**'
      - '.github/workflows/ci-proto.yml'
  pull_request:
    branches: [main]
    paths:
      - 'proto/**'
      - '.github/workflows/ci-proto.yml'

jobs:
  validate:
    name: Validate Proto Files
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install buf
        uses: bufbuild/buf-setup-action@v1

      - name: Lint proto files
        run: buf lint
        working-directory: proto

      - name: Check for breaking changes
        run: buf breaking --against '.git#branch=main'
        working-directory: proto

  test-consumers:
    name: Test Proto Consumers
    needs: validate
    runs-on: ubuntu-latest
    strategy:
      matrix:
        consumer:
          - services/api
          - services/worker
    steps:
      - uses: actions/checkout@v4

      - name: Generate proto code
        run: buf generate
        working-directory: proto

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Test consumer
        run: go test ./...
        working-directory: ${{ matrix.consumer }}
```

#### 5. Common Path Filter Patterns

**Go Monorepo:**
```yaml
paths:
  - 'services/myservice/**'    # Service code
  - 'pkg/shared/**'            # Shared packages
  - 'go.work'                  # Workspace file
  - 'go.work.sum'
  - 'services/myservice/go.mod'
  - 'services/myservice/go.sum'
```

**Node.js/TypeScript Monorepo (pnpm/yarn workspaces):**
```yaml
paths:
  - 'apps/myapp/**'            # App code
  - 'packages/**'              # Shared packages
  - 'package.json'             # Root package.json
  - 'pnpm-lock.yaml'           # or yarn.lock, package-lock.json
  - 'pnpm-workspace.yaml'      # Workspace config
```

**Python Monorepo:**
```yaml
paths:
  - 'services/myservice/**'
  - 'libs/shared/**'
  - 'pyproject.toml'
  - 'poetry.lock'              # or requirements.txt
```

#### 6. Shared Workflow Configuration

Use reusable workflows for common patterns:

```yaml
# .github/workflows/reusable-go-ci.yml
name: Reusable Go CI

on:
  workflow_call:
    inputs:
      working-directory:
        required: true
        type: string
      go-version:
        required: false
        type: string
        default: '1.21'

jobs:
  ci:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ${{ inputs.working-directory }}
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
          cache-dependency-path: ${{ inputs.working-directory }}/go.sum

      - run: go test ./...
      - run: golangci-lint run
      - run: go build ./...
```

Use in subproject workflows:

```yaml
# .github/workflows/ci-api.yml
name: CI - API

on:
  push:
    paths: ['services/api/**', 'proto/**', '.github/workflows/ci-api.yml']
  pull_request:
    paths: ['services/api/**', 'proto/**', '.github/workflows/ci-api.yml']

jobs:
  test:
    uses: ./.github/workflows/reusable-go-ci.yml
    with:
      working-directory: services/api
      go-version: '1.21'
```

#### 7. Complete Monorepo Example

Directory structure:
```
myrepo/
├── .github/
│   └── workflows/
│       ├── ci-api.yml
│       ├── ci-web.yml
│       ├── ci-mobile.yml
│       ├── ci-proto.yml
│       └── reusable-go-ci.yml
├── proto/
│   └── api/
│       └── v1/
├── services/
│   ├── api/          # Go service
│   └── worker/       # Go service
├── apps/
│   ├── web/          # Next.js app
│   └── mobile/       # React Native app
├── packages/
│   └── shared/       # Shared TypeScript packages
├── go.work
└── package.json
```

**Key Rules:**
1. ✅ One workflow per subproject
2. ✅ Always include path filters
3. ✅ Include shared dependencies in path filters (proto/, packages/, etc.)
4. ✅ Include the workflow file itself in path filters
5. ✅ Use `working-directory` in jobs or steps
6. ✅ Use appropriate cache keys per subproject
7. ✅ Test with `gh act` to verify path filters work

**Path Filter Testing:**

```bash
# Test that workflow only runs for relevant changes
git checkout -b test-path-filters

# Change only API code
echo "// test" >> services/api/main.go
gh act -l  # Should only show ci-api workflow

# Change only web code
echo "// test" >> apps/web/src/page.tsx
gh act -l  # Should only show ci-web workflow

# Change proto files
echo "// test" >> proto/api/v1/user.proto
gh act -l  # Should show ci-proto AND ci-api workflows (dependency)
```

#### 8. Avoiding Common Monorepo Pitfalls

**Bad: Single workflow for everything**
```yaml
# ❌ Don't do this in monorepos
name: CI - Everything
on: [push, pull_request]  # No path filters!
jobs:
  test-all:
    steps:
      - run: cd services/api && go test ./...
      - run: cd services/worker && go test ./...
      - run: cd apps/web && npm test
      - run: cd apps/mobile && npm test
      # Runs ALL tests even if only one file changed!
```

**Good: Separate workflows with path filters**
```yaml
# ✅ Do this instead
# .github/workflows/ci-api.yml
name: CI - API
on:
  push:
    paths: ['services/api/**', '.github/workflows/ci-api.yml']
jobs:
  test:
    defaults:
      run:
        working-directory: services/api
    steps:
      - run: go test ./...
```

**Bad: Missing shared dependencies in path filters**
```yaml
# ❌ Missing proto dependency
paths:
  - 'services/api/**'
  # Missing: proto/** (API depends on proto files!)
```

**Good: Include all dependencies**
```yaml
# ✅ Include all relevant paths
paths:
  - 'services/api/**'
  - 'proto/**'              # Shared schemas
  - 'pkg/shared/**'         # Shared packages
  - 'go.work'               # Workspace config
  - '.github/workflows/ci-api.yml'  # This workflow
```

### Error Handling

- Use `continue-on-error` for non-critical steps
- Set appropriate `timeout-minutes`
- Use conditional execution with `if`
- Add failure notifications

Example:
```yaml
steps:
  - name: Run linter
    continue-on-error: true
    run: npm run lint

  - name: Run tests
    timeout-minutes: 10
    run: npm test

  - name: Deploy
    if: github.ref == 'refs/heads/main' && success()
    run: npm run deploy

  - name: Notify on failure
    if: failure()
    uses: actions/github-script@v7
    with:
      script: |
        github.rest.issues.createComment({
          issue_number: context.issue.number,
          owner: context.repo.owner,
          repo: context.repo.repo,
          body: '❌ Workflow failed!'
        })
```

### Outputs and Artifacts

Share data between jobs:

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - id: version
        run: echo "version=1.0.0" >> $GITHUB_OUTPUT

      - run: npm run build

      - uses: actions/upload-artifact@v4
        with:
          name: build
          path: dist/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: build
          path: dist/

      - run: echo "Deploying version ${{ needs.build.outputs.version }}"
```

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract CI/CD requirements
- Note security policies
- Identify deployment targets

**Step 2: Plan the Workflow**
- Define workflow purpose
- Identify required jobs
- Determine trigger events
- List required secrets/variables

**Step 3: Review Existing Workflows**
- Check for similar workflows
- Identify reusable patterns
- Note common actions used

**Step 4: Implement the Workflow**
- Create workflow file in `.github/workflows/`
- Use clear naming and structure
- Add appropriate comments
- Follow security best practices

**Step 5: Test Locally with gh act**

**CRITICAL: All workflows must be tested locally before pushing.**

1. **Install gh act** (if not already installed)
   ```bash
   gh extension install https://github.com/nektos/gh-act
   ```

2. **List workflows to test**
   ```bash
   gh act -l
   ```
   - Review all jobs that will be executed
   - Identify which jobs to test

3. **Test the workflow**
   ```bash
   # Test all jobs
   gh act

   # Test specific job
   gh act -j test

   # Test with verbose output
   gh act -v
   ```
   - If the workflow fails, fix the issues
   - Re-run `gh act` until it succeeds

4. **Test different scenarios**
   ```bash
   # Test pull_request event
   gh act pull_request

   # Test with secrets
   gh act -s GITHUB_TOKEN="$(gh auth token)"

   # Test specific matrix combination
   gh act -j test --matrix node-version:18
   ```

5. **Validate workflow syntax**
   ```bash
   # Use actionlint for additional validation
   actionlint .github/workflows/*.yml
   ```

6. **Document test results**
   - Verify all jobs complete successfully
   - Check that artifacts are created correctly
   - Ensure secrets are handled properly
   - Test failure scenarios

**Step 6: Push and Monitor**

After successful local testing:

1. **Push the workflow**
   ```bash
   git add .github/workflows/
   git commit -m "Add/update workflow"
   git push
   ```

2. **Monitor the first run**
   - Watch the workflow execution in GitHub Actions tab
   - Verify all jobs complete as expected
   - Check for any warnings or issues

3. **Iterate if needed**
   - If issues occur, fix and test locally with `gh act` again
   - Push the fixes
   - Monitor the next run

## Commands Reference

```bash
# gh act - Local workflow testing
gh extension install https://github.com/nektos/gh-act  # Install
gh act -l                        # List workflows
gh act                          # Run default (push) event
gh act pull_request             # Run pull_request event
gh act -j test                  # Run specific job
gh act -W .github/workflows/ci.yml  # Run specific workflow
gh act -n                       # Dry run
gh act -v                       # Verbose output
gh act -s SECRET=value          # Pass secret
gh act --env VAR=value          # Pass environment variable
gh act -s GITHUB_TOKEN="$(gh auth token)"  # Use GitHub token

# GitHub CLI - Workflow management
gh workflow list                # List all workflows
gh workflow view                # View workflow details
gh workflow run                 # Trigger workflow_dispatch
gh workflow enable              # Enable workflow
gh workflow disable             # Disable workflow

# GitHub CLI - Run management
gh run list                     # List workflow runs
gh run view                     # View run details
gh run watch                    # Watch a run in real-time
gh run rerun                    # Rerun a workflow

# Validation
actionlint .github/workflows/*.yml  # Lint workflow files
```

## Common Pitfalls to Avoid

### Not Testing Locally
```yaml
# Bad: Push without testing
git add .github/workflows/ci.yml
git commit -m "Add CI"
git push  # Hope it works!

# Good: Test first with gh act
gh act -j test
# Fix any issues
gh act -j test  # Verify fixes
git push
```

### Hardcoded Secrets
```yaml
# Bad: Hardcoded secrets
steps:
  - run: |
      curl -H "Authorization: token ghp_xxx" ...

# Good: Use GitHub secrets
steps:
  - env:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    run: |
      curl -H "Authorization: token $GITHUB_TOKEN" ...
```

### Missing Timeout
```yaml
# Bad: No timeout (could run for hours)
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: npm test

# Good: Set reasonable timeout
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - run: npm test
```

### Inefficient Caching
```yaml
# Bad: No caching
steps:
  - uses: actions/setup-node@v4
  - run: npm install  # Downloads every time

# Good: Use caching
steps:
  - uses: actions/setup-node@v4
    with:
      cache: 'npm'
  - run: npm ci
```

## Testing Strategy

**Table-Driven Workflow Testing:**

While GitHub Actions workflows are YAML-based, you can use table-driven tests for:
- Custom actions (JavaScript/TypeScript)
- Scripts called by workflows
- Workflow generation tools

Example for testing custom JavaScript action:
```typescript
// __tests__/action.test.ts
describe('Custom Action', () => {
  test.each([
    {
      name: 'valid input',
      inputs: { version: '1.0.0', environment: 'production' },
      expected: { success: true },
    },
    {
      name: 'invalid version',
      inputs: { version: 'invalid', environment: 'production' },
      expected: { success: false, error: 'Invalid version format' },
    },
    {
      name: 'missing environment',
      inputs: { version: '1.0.0' },
      expected: { success: false, error: 'Environment required' },
    },
  ])('$name', async ({ inputs, expected }) => {
    const result = await runAction(inputs);
    expect(result).toMatchObject(expected);
  });
});
```

**Local Testing Checklist:**

Before pushing any workflow:
- [ ] Run `gh act -l` to list all jobs
- [ ] Test with `gh act` (default push event)
- [ ] Test pull_request event if workflow uses it: `gh act pull_request`
- [ ] Test with required secrets: `gh act -s SECRET=value`
- [ ] Test matrix builds if applicable
- [ ] Verify all jobs complete successfully
- [ ] Check that outputs and artifacts are correct
- [ ] Validate with `actionlint` if available

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Extract CI/CD requirements
- [ ] **Check if repository is a monorepo** (multiple subprojects)
- [ ] If monorepo: Identify all subprojects and their dependencies
- [ ] Analyze existing workflows
- [ ] Identify required secrets and permissions
- [ ] Review repository settings

### During Implementation
- [ ] Use clear, descriptive names
- [ ] Add appropriate triggers and filters
- [ ] **For monorepos: Create separate workflow per subproject**
- [ ] **For monorepos: Add path filters to only run on relevant changes**
- [ ] **For monorepos: Include shared dependencies in path filters**
- [ ] **For monorepos: Use `working-directory` in jobs/steps**
- [ ] Implement proper caching
- [ ] Set timeout limits
- [ ] Use pinned action versions
- [ ] Handle secrets securely
- [ ] Add error handling

### After Implementation - MUST ALL PASS
- [ ] **Install gh act**: `gh extension install https://github.com/nektos/gh-act`
- [ ] **List workflows**: `gh act -l` - review jobs
- [ ] **Test locally**: `gh act` - **MUST SUCCEED**
- [ ] **Test specific jobs**: `gh act -j <job-name>` if needed
- [ ] **Test pull_request**: `gh act pull_request` if applicable
- [ ] **Test with secrets**: Add required secrets for testing
- [ ] **Validate syntax**: Run `actionlint` if available
- [ ] **Verify outputs**: Check that artifacts and outputs are correct
- [ ] **Fix all issues** found during local testing
- [ ] **Re-test** until `gh act` succeeds completely
- [ ] Push to repository
- [ ] Monitor first workflow run in GitHub
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Simplicity**: Keep workflows simple and easy to understand
3. **DRY**: Use reusable workflows and composite actions to avoid duplication
4. **Consistency**: Follow established patterns within the repository
5. **Fail-Fast**: Configure workflows to fail immediately on errors with appropriate timeouts
6. **Latest Versions**: Use latest stable versions of actions (e.g., `@v4`)
7. **Encapsulation**: Limit what workflows expose with proper permissions
8. **Early Returns**: Use path filters and conditions to skip unnecessary work
9. **Monorepo Organization**: For monorepos, create separate workflows per subproject with path filters
10. **Test Locally First**: Always use `gh act` to test workflows before pushing
11. **Security**: Never expose secrets, use least privilege permissions
12. **Efficiency**: Use caching and parallel jobs where appropriate
13. **Reliability**: Set timeouts, use error handling, test failure scenarios
14. **Maintainability**: Use clear names, add comments, create reusable workflows
15. **Iterative Testing**: Run `gh act` repeatedly until all issues are resolved

## Version History

### v2.0.0 (2025-11-16)
**Major Update: Industry Best Practices Integration**

Added comprehensive industry best practices based on 2025 research:

- **Core Principles Section**: Added 7 fundamental principles (Simplicity, DRY, Consistency, Fail-Fast, Latest Versions, Encapsulation, Early Returns)
- **Enhanced Security Section**:
  - GITHUB_TOKEN least privilege permissions with examples
  - Pinned action versions for supply-chain security
  - OIDC authentication for cloud providers
  - Comprehensive secrets management guidelines
- **Reusable Workflows**: Updated with Nov 2025 limits (10 nested, 50 total workflows)
  - Added input/output examples
  - Best practices for environment variable passing
  - Semantic versioning guidance
- **Advanced Monorepo Patterns**:
  - Integration of `dorny/paths-filter` action for sophisticated change detection
  - Matrix strategy examples for parallel testing
  - 300-file diff limitation documentation
  - Performance optimization through smart filtering
- **Updated Examples**: All examples now use latest action versions (@v4, @v5)

### v1.0.0 (Initial Version)
- Basic GitHub Actions workflow development guidance
- Local testing with `gh act`
- Monorepo workflow organization
- MCP server setup instructions
- Security best practices
- Reusable workflows and composite actions
