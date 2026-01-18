---
name: github-actions-troubleshoot
description: Troubleshoot GitHub Actions CI errors. Use when GitHub Actions pipelines fail. Provides GitHub-specific commands and CI-specific patterns.
allowed-tools: [Read, Write, Edit, Glob, Grep, Bash]
---

# GitHub Actions Troubleshooting

Follow the general diagnostic workflow from the `troubleshoot` skill. This skill provides GitHub Actions-specific commands and patterns.

## GitHub Actions Commands

### View CI Status

```bash
# View recent CI runs on main branch
gh run list --branch main --limit 5

# View recent CI runs on current branch
gh run list --limit 5

# View details of a specific run
gh run view <run-id>

# View logs of a failed job
gh run view <run-id> --log-failed

# Watch CI run progress
gh run watch

# Re-run failed jobs
gh run rerun <run-id> --failed
```

## Check Main Branch First

**Determine if this is a new error or pre-existing issue before fixing.**

1. Check if main branch CI has the same error using `gh run list --branch main`
2. Compare with your branch:
   - Same error on main: Pre-existing issue, consider if fixing is in scope
   - Main passes: Your changes introduced the error

## Reproduce Locally

**Always reproduce before fixing.**

Match the CI environment:
- Use the same Node/Go/Python version
- Install the same dependencies
- Set similar environment variables

### Using Docker/Act to Match CI

```bash
# If CI uses a specific image
docker run -v $(pwd):/app -w /app <ci-image> <failing-command>

# For GitHub Actions, use act
gh act -l  # List available jobs
gh act -j <job-name>  # Run specific job
```

### Check Versions

```bash
node --version && npm --version
go version
python --version
```

## CI Error Categories

| Category | Examples | Common Causes |
|----------|----------|---------------|
| Build errors | Compilation failures | Missing dependencies, type errors |
| Lint errors | Style violations | Code formatting, static analysis |
| Test failures | Unit/integration/E2E | Logic errors, flaky tests, environment |
| Config errors | CI config issues | YAML syntax, missing secrets |
| Resource errors | Timeouts, OOM | Memory limits, slow tests |

## Common Error Patterns

### Test Failures

```bash
# Run specific failing test
npm test -- --testNamePattern="test name"
go test -v -run TestName ./package/...
pytest -v -k "test_name"
```

### Linting Errors

```bash
# JavaScript/TypeScript
npm run lint
npx eslint --fix .

# Go
golangci-lint run
gofmt -w .

# Python
flake8 .
black .
```

### Dependency Issues

```bash
# Clear caches and reinstall
rm -rf node_modules && npm ci
go clean -modcache && go mod download

# Check for version mismatches
npm outdated
go list -m -u all
```

## Iterate Until CI Passes

After fixing locally, push and monitor:

```bash
git add -A && git commit -m "fix: resolve CI error" && git push
gh run watch
```

If CI still fails, repeat the diagnostic workflow with the new error.
