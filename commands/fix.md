---
description: Fix a bug by reproducing the error, understanding root cause, and planning fixes
argument-hint: "describe the error and what fix is needed"
allowed-tools: ["*"]
---

# Fix

$ARGUMENTS

## Workflow

### 1. Analyze Existing Codebase (Architect Agent)

Analyze the codebase to understand where the error happens:
- Identify affected components
- Understand code flow leading to the error
- Map out relevant parts of the codebase

### 2. Reproduce Error and Find Root Cause (Software Engineer Agent)

A software engineer with the appropriate tech stack must:
- Reproduce the error to confirm the issue
- Understand the root cause of the error
- Document findings and analysis

### 3. Plan Changes (Architect Agent)

Based on the analysis and root cause:
- Plan the necessary fixes
- Identify opportunities for parallel changes
- Ensure the fix addresses the root cause, not just symptoms

### 4. Confirm Plan with User

**IMPORTANT**: Present the analysis, root cause findings, and planned fixes to the user. Wait for user approval before proceeding with any implementation.

## Guidelines

- Always reproduce the error first before attempting fixes
- Understand root cause before fixing
- Follow general coding guidelines (DRY, fail-fast, simplicity)
- Adhere to project-specific guidelines from `.claude/docs/guideline.md`
