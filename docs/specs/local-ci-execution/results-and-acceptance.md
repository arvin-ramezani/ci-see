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

## Dependency acquisition and reuse (2026-10-12)

**Owner decision:** Missing workflow dependencies should be downloaded when needed, then reused from a local cache. This clarifies §3 and §12 without requiring CI See to build its own dependency/cache engine.

- Reuse previously fetched actions, runner images and other compatible dependencies only when their identity and integrity can be established. Stale, changed, untrusted or corrupt cache entries are not executable.
- On a cache miss, acquire the dependency from an approved source with controlled network access; verify the expected commit, digest or other approved provenance/integrity policy **before execution**, then reuse the cached copy. A cache miss alone is not a CI failure.
- **Actions:** For an immutable `uses: owner/repo@<full-commit-SHA>` reference, verify the cached Git object/tree and materialized files match that commit and a trusted origin. For tag/branch references, resolve the reference at an approved source to a full commit SHA; record the resolved identity and verify it before reuse. If the reference later resolves differently, treat it as a new dependency and invalidate affected prior results. Do not trust a matching cache directory name alone.
- **Images:** Resolve image references at an approved registry to an immutable OCI manifest digest (for the selected platform), then verify the local image corresponds to that trusted digest before use or reuse. A mutable tag, local image name, or image ID alone does not establish trusted provenance; a changed digest invalidates dependent results.
- **Enforcement boundary:** Use the existing `act`/Git/Docker download and cache facilities, not a new CI See cache engine. Before executing untrusted workflow code, the Go adapter must validate authoritative expected identities against the actual cached artifacts; if the engine cannot expose evidence or enforce the selected verified artifacts, stop as INCOMPLETE. Do not allow an implicit `act`/Docker pull or cache substitution to bypass this check.
- **Evidence and failures:** Record non-secret source, resolved commit/digest, platform and verification outcome in the execution fingerprint; secrets and tokens must not be persisted. Reject cache poisoning, mismatched refs or digests, unsupported identity formats, unverifiable downloaded content and missing authorization with INCOMPLETE.
- If acquisition, trust verification, or required network access fails, report **INCOMPLETE**, never PASS; do not silently substitute a different version.
- Treat dependency acquisition and execution as distinct security boundaries. Restricted offline/D3 verification remains offline, while normal CI See operation may download missing dependencies under a separately reviewed network policy. Some workflow steps need runtime networking; do not promise that prefetching alone makes all workflows offline-compatible.
- Preserve the requirement that users install `act` and Docker or an equivalent runtime; this decision does not authorize automatic installation of host tools or create an independent CI See cache service.
- Before real-`act` acceptance, test cold-cache acquisition, warm-cache reuse, pinned-identity changes, invalid sources, corrupted cache entries, interrupted downloads, missing network access and non-PASS failures. Verify that credentials are not implicitly exposed to download subprocesses.

This is a behavior decision, **not evidence that dependency acquisition or network security has been implemented**.
