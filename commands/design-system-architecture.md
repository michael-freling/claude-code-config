---
description: Design system architecture including API schema, database schema, and security
argument-hint: "<feature-name> <system to design>"
allowed-tools: ["Read", "Glob", "Grep", "Task", "Write", "LSP"]
---

# Design System Architecture

$ARGUMENTS

## Workflow

### 1. Gather Context (Software Architect Agent)

Use the software-architect agent to:
- Review product specification (if exists at `docs/<feature-name>/product-spec.md`)
- Analyze existing architecture
- Identify current patterns and constraints

### 2. Design API (Software Architect Agent)

Use the software-architect agent to:
- Design endpoints (REST/gRPC)
- Define request/response schemas
- Document authentication requirements
- Plan error handling

### 3. Design Database Schema (Software Architect Agent)

Use the software-architect agent to:
- Create entity-relationship diagrams
- Define table structures
- Plan indexes
- Document migration approach

### 4. Design Security Architecture (Software Architect Agent)

Use the software-architect agent to:
- Design authentication flow
- Define authorization model
- Plan data protection

### 5. Review Design (Architecture Reviewer Agent)

Use the architecture-reviewer agent to:
- Review scalability
- Check security considerations
- Validate API design consistency
- Review database normalization

### 6. Confirm with User

**IMPORTANT**: Present the architecture design and review feedback to the user. Confirm the design is acceptable before finalizing. Do not finalize until you get approval from the user.

### 7. Write Output

Write the architecture design to: `docs/<feature-name>/architecture.md`

## Output Structure

The output file should include:
- Executive Summary
- Requirements Analysis
- Current State Assessment
- API Design (endpoints, schemas, authentication)
- Database Schema (ERD, tables, indexes)
- Security Architecture
- Trade-offs and Decisions
- Open Questions

## Guidelines

- NO implementation code
- Focus on design decisions and trade-offs
- Document alternatives considered
- Address security in every design
- Be explicit about assumptions
