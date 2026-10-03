# CI See — System Architecture

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/architecture.md`  
**Last updated:** 2026-10-03

## 1. Purpose

CI See is a local-first CLI that runs a repository's existing GitHub Actions locally and gates normal Git operations using exact-state CI results.

This document defines MVP system boundaries and major technical decisions. Detailed behavior belongs in focused specs.

## 2. Architecture Goals

- Keep GitHub Actions as the single CI definition.
- Reuse `nektos/act`; do not build a workflow engine.
- Integrate with normal `git commit` and `git push`.
- Accept PASS only for the exact validated state.
- Fail safely on missing, stale, interrupted, or corrupt state.
- Keep the MVP small, testable, and easy for humans and AI agents to understand.

## 3. MVP Decisions

| Area | Decision |
| --- | --- |
| Language | TypeScript / Node.js |
| Application | CLI-first; no long-lived daemon or desktop app. A short-lived human approval surface may be invoked when required. |
| Execution engine | User-installed `act`, discovered from `PATH` |
| Container runtime | Docker-compatible runtime used by `act` |
| Git integration | Repository Git hooks |
| Persistence | Simple repository-local files under Git metadata |
| Primary development | Windows 11 + WSL2 |
| Execution focus | Linux-container GitHub Actions workloads |

CI See owns orchestration, state validation, Git gating, diagnostics, and developer-facing decisions. `act` owns GitHub Actions execution.

## 4. System Context

```text
Developer / AI agent
        |
        v
    Git / ci-see
        |
        v
    CI See CLI
   /     |      \
 Git   Local    act
hooks  state     |
                 v
               Docker
                 |
                 v
       workflow containers
```

CI See invokes `git` and `act` as child processes. It does not embed or reimplement either tool.

## 5. Components

### CLI

Provides the user-facing commands, initially:

```text
ci-see
ci-see init
ci-see status
```

### Core

Owns product rules independent of Git, `act`, files, or terminal presentation:

- validation states;
- exact-state matching;
- gate decisions;
- run lifecycle;
- approval boundary.

### Adapters

- **Git adapter:** repository discovery, state identity, hook integration, commit/push context.
- **`act` adapter:** binary detection, safe process execution, result mapping.
- **State adapter:** atomic local persistence and repository/worktree isolation.

### Git hooks

`ci-see init` installs thin hooks that call the CLI. Product rules must remain in the core, not in hook scripts.

## 6. Local State

Store private state under Git metadata resolved through Git, not by assuming a literal `.git/` directory.

Conceptually:

```text
<git-dir>/ci-see/
  state.json
  runs/
```

Requirements:

- atomic writes;
- repository/worktree isolation;
- conflicting writes serialized with a local lock;
- interrupted execution can never become PASS.

Exact fingerprint fields belong in the Git gating spec.

## 7. Runtime Flows

### Local run

```text
ci-see
→ inspect repository/workflows
→ determine validation identity
→ record RUNNING
→ execute act
→ persist PASS / FAIL / INCOMPLETE
→ return a clear exit result
```

### Commit gate

```text
git commit
→ hook invokes CI See
→ calculate staged-state identity
→ matching PASS: continue
→ otherwise: block and explain
```

### Push gate

```text
git push
→ hook invokes CI See
→ identify pushed commit state
→ matching acceptable result: continue
→ otherwise: block
```

## 8. Approval Boundary

Approval remains a separate core boundary.

The MVP does not require a long-lived daemon or desktop application, but the CLI may invoke a short-lived human approval surface when a developer decision is required. The exact mechanism belongs in the approval UX spec.

CI See must not expose a trivial machine-only `--force` or equivalent bypass that an AI agent could silently use. Non-interactive Git operations must fail safely if a valid developer approval cannot be obtained.

This keeps the core CLI architecture small while preserving the PRD requirement for developer-controlled bypass across terminal, IDE, and AI-agent initiated Git operations.

## 9. Reliability and Security

- Unknown, stale, corrupt, or interrupted state is not PASS.
- Spawn `git` and `act` with argument arrays; avoid shell-string construction.
- Validate external/process/file data at boundaries.
- Do not persist secrets in result metadata.
- Logs must explain why a gate passed, failed, or became stale.
- Local `act` success is evidence, not proof of exact GitHub-hosted runner parity.

## 10. Platform Strategy

Primary development is Windows 11 + WSL2.

Architecture must avoid WSL-specific assumptions. Platform-specific behavior stays behind adapters so native Windows, Linux, and macOS can be validated independently.

The MVP focuses on Linux-container execution through `act`; full Windows/macOS hosted-runner emulation is not required.

## 11. Code Boundaries

```text
src/
  cli/
  core/
  adapters/
    git/
    act/
    state/
  hooks/
```

Rules:

- `core` must not import infrastructure adapters.
- Git/`act`/filesystem details stay in adapters.
- CLI and hooks call core use cases instead of duplicating rules.
- TypeScript strict mode is required.
- Boundary data should be runtime-validated.
- Tests must enforce exact-state PASS and fail-closed behavior.

## 12. Deferred Decisions

Focused specs or ADRs will define:

- event selection and generated event payloads;
- exact validation fingerprint inputs;
- developer approval/bypass UX;
- secrets, artifacts, and cache behavior;
- runner-image policy;
- duplicate GitHub-hosted execution control;
- native Windows/macOS runner behavior.

## 13. Next Specifications

1. Local CI execution specification.
2. Git gating specification.
3. Developer approval UX specification.
4. MVP implementation plan.

Implementation should begin only when the relevant behavior is specified enough to test deterministically.
