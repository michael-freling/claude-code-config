# GitHub Actions Workflow Guidelines

## Reduce CI Log Noise

When creating or modifying GitHub Actions workflows, use quiet flags to minimize verbose output while preserving error visibility. This reduces token consumption when AI agents analyze CI logs.

### Required Quiet Flags

**Docker:** See `docker-commands.md` for Docker and Docker Compose quiet flags.

**Package Managers:**
- `pnpm install --frozen-lockfile --silent` - Suppress installation output
- `npm install --silent` - Suppress installation output

**System Commands:**
- `apt-get update > /dev/null` - Redirect stdout to suppress output
- `apt-get install -y package > /dev/null` - Redirect stdout to suppress output
- `curl -sSL` - Silent mode with show errors (not `-SL`)

Note: `-qq` flag for apt-get still shows "Reading database..." output. Use `> /dev/null` to truly silence it while preserving stderr for errors.

### Preserve Error Visibility

These flags suppress success output but preserve error messages. Do not redirect stderr to `/dev/null` unless explicitly handling errors.

### Example

```yaml
- name: Install dependencies
  run: |
    sudo apt-get update > /dev/null
    sudo apt-get install -y make > /dev/null
    sudo curl -sSL "https://example.com/file" -o /tmp/file
```

For Docker examples, see `docker-commands.md`.
