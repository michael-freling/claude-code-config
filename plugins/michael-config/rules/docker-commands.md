---
paths: .github/workflows/*.{yml,yaml}, **/*.sh, Makefile, *.mk
---

# Docker Command Best Practices

When running Docker commands in GitHub Actions workflows, bash scripts, or Makefiles, follow these recommendations:

## Why Use Quiet Flags

Docker commands produce verbose output with progress bars, layer information, and download status. When AI agents execute these commands, the output consumes significant context window space without providing actionable information. Using quiet flags keeps output minimal and preserves context for meaningful content.

## Docker Pull

Use `--quiet` or `-q` flag to reduce noisy stdout:

```bash
docker pull --quiet <image>
```

## Docker Build

Use `docker buildx build` instead of `docker build` for improved features:

```bash
docker buildx build --quiet -t <image> .
```

Benefits of buildx:
- Multi-platform builds
- Better caching
- BuildKit features by default

## Docker Compose

Use `--quiet-pull` flag when pulling images:

```bash
docker compose pull --quiet
docker compose up -d --quiet-pull
```

## Summary of Quiet Flags

| Command | Flag |
|---------|------|
| `docker pull` | `--quiet` or `-q` |
| `docker build` / `docker buildx build` | `--quiet` or `-q` |
| `docker compose pull` | `--quiet` or `-q` |
| `docker compose up` | `--quiet-pull` |
