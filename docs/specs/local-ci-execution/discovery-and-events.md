> [Back to Local CI execution specification index](../local-ci-execution.md)

<!-- BEGIN ORIGINAL CONTRACT -->
# CI See — Local CI Execution Specification

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/specs/local-ci-execution.md`  
**Depends on:** `docs/prd.md`, `docs/architecture.md`  
**Last updated:** 2026-10-09

## 1. Purpose

Define how CI See discovers and runs GitHub Actions locally for the MVP.

CI See orchestrates execution. `nektos/act` remains the workflow execution engine.

## 2. Scope

This spec owns:

- workflow discovery;
- supported local events;
- generated event context;
- `act` and container-runtime preflight;
- execution lifecycle;
- local secrets and variables;
- logs and normalized results;
- local parity boundaries.

This spec does not own Git commit/push gating or developer bypass behavior.

## 3. Dependencies

The user installs:

- the **prebuilt Go `ci-see` executable** for their host OS/architecture (or builds it from source);
- `act`, available on `PATH`;
- Docker or another `act`-compatible container runtime.

Users of a prebuilt CI See executable do **not** need to install Go, Node.js, or npm. The Go toolchain is needed only when building CI See from source. CI See must not invoke Node.js as part of its own orchestration; GitHub Actions workflows executed by `act` may independently require Node.js or other runtimes inside their jobs.

CI See must detect missing or unusable dependencies before starting a run and return INCOMPLETE, never PASS.

CI See records the detected `act` version with each run.

## 4. Workflow Discovery

CI See discovers workflow files under:

```text
.github/workflows/**/*.yml
.github/workflows/**/*.yaml
```

Users do not register workflows manually.

A workflow is locally eligible when it exposes at least one MVP-supported event.

Unsupported-only workflows are reported as unsupported and are not silently treated as executed.

## 5. MVP Event Policy

Supported events:

1. `pull_request`
2. `push`
3. `workflow_dispatch` where no unresolved required input prevents execution

For normal `ci-see` execution, each eligible workflow is run once.

Event preference for a workflow:

```text
pull_request
→ otherwise push
→ otherwise workflow_dispatch
```

This preference favors the event most representative of pre-merge validation while avoiding duplicate execution of a workflow that supports multiple events.

If a `workflow_dispatch` workflow requires inputs that CI See cannot derive, it is reported as not executable for that run.

<!-- END ORIGINAL CONTRACT -->
