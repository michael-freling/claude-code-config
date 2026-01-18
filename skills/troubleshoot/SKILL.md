---
name: troubleshoot
description: General diagnostic principles for troubleshooting any system. Use when diagnosing errors, failures, or unexpected behavior in applications, infrastructure, or CI/CD. CRITICAL - Always identify the exact error before suggesting any fix.
---

# Troubleshooting Skill

## Critical Rule

**NEVER suggest a fix or root cause without completing the diagnostic workflow.**

Do not speculate. Do not assume. Read the logs, identify the exact error, and analyze what it means before proposing any solution.

## Diagnostic Workflow

### Step 1: Get Logs

Retrieve relevant logs from the failing component. Look for:
- Error messages
- Stack traces
- Warning messages preceding the error
- Timestamps to correlate events

### Step 2: Identify the Exact Error

Find and **quote the exact error message**. Look for:
- `error`, `failed`, `unable to`
- Exception names and messages
- Exit codes

**You MUST quote the exact error message before proceeding.**

### Step 3: Analyze the Error Semantically

Break down what each part of the error message means:

1. **What operation failed?** (request, sync, create, update, delete, connect)
2. **What resource/component was involved?** (file, service, endpoint, object)
3. **What was the specific failure reason?** (not found, permission denied, timeout, invalid)
4. **Is there additional context?** (field path, line number, stack trace)

### Step 4: Check Current State

Verify the actual state of the system:
- Does the resource/file exist?
- Is the service running?
- Are permissions correct?
- What are the current values/configuration?

### Step 5: Identify Root Cause from Evidence

Based on the error message AND current state, identify the root cause. Common patterns:

| Error Pattern | Likely Cause |
|---------------|--------------|
| `not found` | Resource missing or wrong path |
| `permission denied` / `forbidden` | Access control issue |
| `connection refused` | Service not running or wrong endpoint |
| `timeout` | Network issue or service overloaded |
| `already exists` | Duplicate resource or stale cache |
| `invalid` | Schema/format validation failed |

### Step 6: Apply Targeted Fix

Only after identifying the root cause with evidence, apply the fix:
- Address the specific cause identified
- Make minimal changes needed
- Don't fix unrelated issues simultaneously

### Step 7: Verify Fix

After applying the fix:
- Check logs for errors
- Verify the operation now succeeds
- Confirm the system is in expected state

## Anti-Patterns to Avoid

1. **DO NOT** suggest "restart" as the first step
2. **DO NOT** assume the cause without reading the error message
3. **DO NOT** suggest deleting/recreating without understanding why it failed
4. **DO NOT** speculate about race conditions or edge cases without evidence
5. **DO NOT** skip verifying the fix worked
6. **DO NOT** fix multiple unrelated issues at once
