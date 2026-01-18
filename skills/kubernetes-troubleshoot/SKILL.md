---
name: kubernetes-troubleshoot
description: Diagnose Kubernetes issues including pod failures, ArgoCD sync problems, and controller errors. Use when pods crash, ArgoCD applications show sync issues, or controllers fail to reconcile resources.
---

# Kubernetes Troubleshooting

Follow the general diagnostic workflow from the `troubleshoot` skill. This skill provides K8s-specific commands and patterns.

For controller-specific notes, see:
- `references/prometheus-operator.md`
- `references/cert-manager.md`
- `references/external-secrets.md`

## Getting Logs

### Pod Logs

```bash
# Get pod logs
kubectl logs -n <namespace> <pod-name> --tail=100

# Get logs from previous crashed container
kubectl logs -n <namespace> <pod-name> --previous

# Get logs from specific container in multi-container pod
kubectl logs -n <namespace> <pod-name> -c <container-name>

# Follow logs
kubectl logs -n <namespace> <pod-name> -f
```

### Controller/Operator Logs

```bash
# Find controller pod
kubectl get pods -n <namespace> | grep -E "(operator|controller|manager)"

# Get controller logs
kubectl logs -n <namespace> deploy/<controller-deployment> --tail=100
```

## Checking Resource State

```bash
# Get resource status
kubectl get <kind> <name> -n <namespace>

# Describe for events and conditions
kubectl describe <kind> <name> -n <namespace>

# Get labels
kubectl get <kind> <name> -n <namespace> -o jsonpath='{.metadata.labels}' | jq .

# Get owner references
kubectl get <kind> <name> -n <namespace> -o jsonpath='{.metadata.ownerReferences}' | jq .

# Get CR status
kubectl get <kind> <name> -n <namespace> -o jsonpath='{.status}' | jq .

# Get conditions
kubectl get <kind> <name> -n <namespace> -o jsonpath='{.status.conditions}' | jq .
```

## Common Error Patterns

| Error Pattern | Likely Cause | How to Verify |
|---------------|--------------|---------------|
| `ImagePullBackOff` | Wrong image name/tag or registry auth | Check image name, registry credentials |
| `CrashLoopBackOff` | Application crash on startup | Check pod logs with `--previous` |
| `Pending` | No schedulable node | Check node resources, taints, affinity |
| `CreateContainerConfigError` | Missing ConfigMap/Secret | Check referenced configs exist |
| `OOMKilled` | Out of memory | Check memory limits vs actual usage |
| `creating X failed: X already exists` | Controller cache stale | Check resource labels, consider restart |

---

## ArgoCD Sync Issues

### Check Application Status

```bash
kubectl get applications -n <argocd-namespace> -o custom-columns='NAME:.metadata.name,SYNC:.status.sync.status,HEALTH:.status.health.status,PHASE:.status.operationState.phase'
```

### Check Actual Diff (CRITICAL)

Never assume the cause of OutOfSync. Always check the actual diff:

```bash
kubectl get application <app-name> -n <argocd-namespace> -o jsonpath='{.status.resources[?(@.status=="OutOfSync")]}'
```

Or use ArgoCD UI to see the detailed diff.

### Common Causes and Fixes

| Cause | How to Identify | Fix |
|-------|-----------------|-----|
| Kubernetes default values | Diff shows fields like `group: ""`, `conversionStrategy: Default` | Add defaults to manifest |
| Runtime-managed fields | Diff shows cloud provider annotations, finalizers | Use `ignoreDifferences` in Application |
| Manifest error | Sync fails with validation error | Fix the manifest syntax/schema |
| Hook stuck | Phase shows "Running" | Check Job/Pod status, may need to delete |

### Fixing OutOfSync

#### Kubernetes Default Values

Add defaults explicitly to manifests:

```yaml
# Example: Gateway certificateRefs
certificateRefs:
  - name: gateway-tls
    kind: Secret
    group: ""  # Add this default explicitly

# Example: ExternalSecret
remoteRef:
  key: secret-name
  conversionStrategy: Default
  decodingStrategy: None
  metadataPolicy: None
```

#### Runtime-Managed Fields

For fields managed at runtime, use `ignoreDifferences`:

```yaml
# In ArgoCD Application spec
ignoreDifferences:
  - group: gateway.networking.k8s.io
    kind: Gateway
    jqPathExpressions:
      - .metadata.annotations | keys | map(select(startswith("networking.gke.io")))
      - .metadata.finalizers
```

### Verify Fix

```bash
kubectl get application <app-name> -n <argocd-namespace> -o jsonpath='Sync: {.status.sync.status}, Health: {.status.health.status}, Phase: {.status.operationState.phase}'
```

---

## Controller Troubleshooting

### Cache/Informer Issues

If controller cache doesn't see existing resources:

1. Check if resource has expected labels (controller uses label selectors)
2. Restart controller to refresh cache: `kubectl rollout restart deploy/<controller>`
3. If label mismatch: delete resource and let controller recreate

### Resource Conflicts

If resource exists but shouldn't:

1. Check who created it (ownerReferences, annotations)
2. Delete the conflicting resource
3. Let the correct controller recreate it

### RBAC Issues

If permission denied:

1. Check controller's ServiceAccount
2. Verify ClusterRole/Role bindings
3. Add missing permissions
