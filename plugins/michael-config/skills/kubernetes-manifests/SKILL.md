---
name: kubernetes-manifests
description: Kubernetes manifest best practices. Use when creating or modifying K8s manifests, Kustomize overlays, or Helm charts.
---

# Kubernetes Manifests

## Validation

- Validate before committing: `kubectl apply --dry-run=client -f <manifest>`
- For Kustomize: `kubectl kustomize <overlay-path> | kubectl apply --dry-run=client -f -`

## Kustomize

When using Kustomize:
- Create namespace as a Kubernetes resource (`namespace.yaml`)
- Set namespace in `kustomization.yaml`, not in individual resources
- Use base for shared configs, overlays for environment-specific differences only

## Helm Charts

When using Helm charts, search the web for the latest stable version.

## ArgoCD

When using ArgoCD:
- If OutOfSync is caused by K8s default values, add them explicitly to manifest
- Only use `ignoreDifferences` for runtime-managed fields (e.g., GKE annotations)
