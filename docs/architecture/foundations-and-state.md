> [Back to System architecture index](../architecture.md)

<!-- BEGIN ORIGINAL CONTRACT -->
# CI See — System Architecture

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/architecture.md`  
**Last updated:** 2026-10-09

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
| Language | Go, using Go modules; compiled CLI (no Node.js runtime) |
| Application | CLI-first Go executable; no long-lived daemon or desktop app. A short-lived human approval surface may be invoked when required. |
| Execution engine | User-installed `act`, discovered from `PATH` |
| Container runtime | Docker-compatible runtime used by `act` |
| Git integration | Repository Git hooks |
| Persistence | Simple repository-local files under Git metadata |
| Primary development | Windows 11 + WSL2 |
| Execution focus | Linux-container GitHub Actions workloads |
| Packaging | Native binaries per supported host OS/architecture; Go toolchain needed only to build from source |

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

The Go executable invokes `git` and `act` as child processes using argument arrays (Go `os/exec`, with cancellation where needed). It does not embed or reimplement either tool. Existing `act`, Docker, and Git are external dependencies; Node.js is not.

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

<!-- END ORIGINAL CONTRACT -->
