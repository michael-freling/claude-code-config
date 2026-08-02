---
name: validate-plugin
description: Validates this repo's marketplace and plugin manifests plus agent, command, and skill frontmatter. Use before opening a PR that touches .claude-plugin/, plugins/, or .claude/skills/.
allowed-tools: Bash
---

# Validate Plugin Skill

Catch broken plugin metadata locally, before it ships to anyone installing
`michael-config@michael-freling`. Run every command from the repo root.

The `Makefile` has no validation target (only `install`, `update`,
`uninstall`, `gemini-install`), so these checks run directly.

## 1. Manifests and frontmatter

`claude plugin validate` parses each JSON manifest and the YAML frontmatter of
every agent, command, and skill inside a plugin. It only prints files with
findings.

```bash
claude plugin validate .
claude plugin validate ./plugins/michael-config
```

Warnings (missing `description:`, unrecognized fields) exit 0. Use `--strict`
to fail on them:

```bash
claude plugin validate ./plugins/michael-config --strict
```

The most common real failure is an agent `description:` written as a plain
scalar that then spans unindented lines — YAML cannot parse it and the whole
frontmatter is silently dropped at runtime. Use a block scalar instead:

```yaml
description: |
  First line.

  More lines, all indented.
```

## 2. Every JSON file parses

Covers `settings.json` and anything else `claude plugin validate` does not
read. `git ls-files -co --exclude-standard` skips ignored paths such as
`.claude/worktrees/`.

```bash
git ls-files -co --exclude-standard '*.json' | xargs -r -n1 jq empty \
  && echo "all JSON parses"
```

## 3. Marketplace sources resolve

`claude plugin validate` does **not** check that a `source` path exists, so
check it here.

```bash
for src in $(jq -r '.plugins[].source' .claude-plugin/marketplace.json); do
  [ -f "$src/.claude-plugin/plugin.json" ] \
    && echo "ok   $src" \
    || echo "FAIL $src/.claude-plugin/plugin.json is missing"
done
```

## 4. Skill frontmatter and directory names

`claude plugin validate` does not compare a skill's `name:` to its directory,
and it cannot see project-local skills under `.claude/skills/` at all (they are
not part of the plugin). This covers both.

```bash
for f in plugins/*/skills/*/SKILL.md .claude/skills/*/SKILL.md; do
  [ -f "$f" ] || continue
  dir=$(basename "$(dirname "$f")")
  if [ "$(head -n1 "$f")" != "---" ] || [ "$(grep -c '^---$' "$f")" -lt 2 ]; then
    echo "FAIL $f: frontmatter must start on line 1 with --- and be closed"
    continue
  fi
  fm=$(sed -n '2,/^---$/p' "$f" | sed '$d')
  printf '%s\n' "$fm" | grep -q '^description:' \
    || echo "FAIL $f: frontmatter missing description:"
  [ "$(printf '%s\n' "$fm" | sed -n 's/^name: *//p')" = "$dir" ] \
    || echo "FAIL $f: name: must match directory '$dir'"
done
```

## Notes

- Agents need `name:` and `description:`; the `name:` must match the filename
  stem. Commands need only `description:` — their name comes from the filename.
- Skills under `.claude/skills/` are repo-local. Only the `source` directory
  named in `.claude-plugin/marketplace.json` is distributed to installers, so
  nothing there reaches plugin users.
- Report every failure found. Do not fix unrelated pre-existing failures unless
  asked.
