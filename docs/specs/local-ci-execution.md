# CI See — Local CI Execution Specification

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/specs/local-ci-execution.md`  
**Depends on:** `docs/prd.md`, `docs/architecture.md`  
**Last updated:** 2026-10-03

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

- Node.js supported by CI See;
- `act`, available on `PATH`;
- Docker or another `act`-compatible container runtime.

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

## 6. Event Payloads

Users must not need to create event JSON for normal usage.

CI See generates the temporary event payload needed by `act` from local Git/repository context.

For `pull_request`, the generated context should include at minimum:

- current/head branch;
- base/default branch when known;
- head commit SHA when available;
- base commit SHA when available;
- repository identity.

For `push`, it should include at minimum:

- current ref;
- current commit SHA;
- repository identity.

Generated payloads are temporary execution inputs and must not be committed.

If required event context cannot be determined safely, the affected workflow is INCOMPLETE rather than guessed.

## 7. Execution Flow

```text
ci-see
→ verify repository
→ discover workflows
→ verify act/runtime
→ resolve event for each eligible workflow
→ generate temporary event context
→ record RUNNING
→ invoke act
→ capture result/logs
→ persist normalized result
→ remove temporary sensitive files
→ return final status
```

CI See invokes `act` as a child process using argument arrays, not shell-built command strings.

## 8. Execution Isolation

Each run is bound to:

- one repository/worktree;
- one exact local validation state;
- the selected workflows/events;
- the detected `act` version;
- relevant CI See execution configuration.

Concurrent runs must not overwrite each other's temporary files or result records.

Detailed validation fingerprint rules belong in the Git gating specification.

## 9. Secrets and Variables

CI See must not place secret values directly in generated shell command strings.

MVP sources:

- environment variables;
- optional repository-local private secret file under CI See Git metadata;
- optional repository-local private variable file under CI See Git metadata.

Conceptually:

```text
<git-dir>/ci-see/
  secrets.env
  vars.env
```

These files are outside tracked repository content.

CI See passes them to `act` through supported secret/variable mechanisms.

Rules:

- secret values must not be written to normal CI See logs;
- CI See must not copy secrets into result metadata;
- missing required secrets must produce a clear non-PASS result;
- file permissions should be restricted where the host platform supports it.

## 10. Logs

Normal output should show:

- workflow/job currently running;
- PASS/FAIL/SKIPPED outcome;
- engine/runtime failures;
- unsupported workflow/event reasons;
- final normalized result.

Detailed `act` output remains available for diagnosis.

CI See should preserve enough local run metadata to explain later why a result passed, failed, or was incomplete.

## 11. Result Model

A completed local execution produces one overall result:

- **PASS** — every selected executable workflow completed successfully.
- **FAIL** — execution completed and at least one selected workflow/job failed.
- **INCOMPLETE** — validation could not complete reliably, including interruption, dependency/runtime failure, unresolved required event input, or corrupt execution state.

`NOT_RUN` and `STALE` are state/gating concepts handled outside the active execution lifecycle.

A process crash or cancellation must never leave PASS.

## 12. Artifacts and Cache

For the MVP:

- CI See does not implement its own artifact service;
- CI See does not implement its own dependency/cache engine;
- behavior provided directly by `act` may be used;
- CI See must not depend on cache presence for correctness.

Additional artifact/cache product behavior is deferred until a concrete user requirement exists.

## 13. Compatibility Boundary

MVP execution targets Linux-container-compatible GitHub Actions workloads.

CI See does not promise exact GitHub-hosted runner parity.

A local PASS means:

> the selected workflows completed successfully under the recorded local CI See + `act` execution context.

It does not mean GitHub-hosted execution is guaranteed to produce the same result.

Unsupported features must be surfaced rather than silently converted into PASS.

## 14. Required Acceptance Tests

Before this spec is considered implemented, tests must prove:

1. workflow files are discovered automatically;
2. supported events are selected deterministically;
3. users do not need to create event JSON manually;
4. missing `act` cannot produce PASS;
5. failed `act` execution cannot produce PASS;
6. interrupted execution cannot produce PASS;
7. secrets are not included in normal result metadata;
8. concurrent runs cannot corrupt one another's state;
9. unsupported execution is visible and non-PASS when it prevents required validation;
10. recorded results include the detected `act` version.

## 15. Deferred

The following remain outside this spec:

- exact Git validation fingerprint;
- commit and push gate rules;
- developer approval/bypass UX;
- remote GitHub-hosted duplicate-execution control;
- native Windows/macOS hosted-runner emulation;
- advanced artifact/cache management.
