# Local CI execution specification — Results, artifacts, compatibility, acceptance 1–12

**Status:** Draft (inherited)  
**Canonical authority:** [Local CI execution specification](../local-ci-execution.md)

<!-- BEGIN ORIGINAL CONTRACT -->
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
10. recorded results include the detected `act` version;
11. a prebuilt CI See binary can run without Go or Node.js installed on the host;
12. CI See's own orchestration runs through Go and does not depend on npm or a JavaScript runtime.

## 15. Deferred

The following remain outside this spec:

- exact Git validation fingerprint;
- commit and push gate rules;
- developer approval/bypass UX;
- remote GitHub-hosted duplicate-execution control;
- native Windows/macOS hosted-runner emulation;
- advanced artifact/cache management.
<!-- END ORIGINAL CONTRACT -->
