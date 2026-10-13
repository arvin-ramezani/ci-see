# Local CI execution — Outcomes, dependencies and acceptance

**Status:** Draft aligned with merged PRD #20
**Canonical authority:** [Local CI execution](../local-ci-execution.md)

## 11. Result Model

- **PASS:** nonempty selected validation completed successfully with supported
  required behavior, identified inputs, verified artifacts and complete reporting.
- **FAIL:** selected validation completed reliably and a required job/step failed.
- **INCOMPLETE:** required work could not be completed/reliably classified, including
  unsupported required behavior, unresolved input, dependency/runtime error,
  interrupted execution, uncertain cleanup or failed final report persistence.

If an incomplete required job coexists with a failed job, overall INCOMPLETE retains
the failure details; incomplete validation cannot be masked by aggregate exit zero.
An abandoned RUNNING record is incomplete. NOT_RUN means no execution; no stale
gate state, allow/block outcome or reusable-PASS authorization is required.

Record run/repository/worktree and actual source identity, dirty state, selected
workflows/jobs/matrices/events, effective non-secret inputs/configuration, tool/image/
dependency identities, timestamps, outcomes, exclusions and limitations.
Required selection conditions legitimately skipped by workflow semantics are visible;
zero executed required validation cannot silently become PASS.

**CLI-01:** numeric exits and versioned JSON/report schema require a focused contract.
Exit zero is limited to successful informational/setup commands or proven completed
PASS; FAIL/INCOMPLETE/invalid usage must be nonzero. Current S1 exits remain unchanged.

## 12. Artifacts and Cache

Use `act`, Docker and workflow package-manager cache facilities; no new cache engine
or artifact service. Local execution must not depend on cache presence for correctness.
Expose artifact availability/limitations, but do not claim unavailable hosted services work.

## 13. Compatibility Boundary

Initially execute Linux-container-compatible jobs on tested host platforms.
Keep existing YAML; configuration can select supported images without claiming
unsupported native runner labels are equivalent. No automatic Windows/macOS remapping.
Unsupported required features yield non-PASS, even when `act` ignores them.
Publish a tested engine/feature/runner matrix; local PASS is not hosted parity.

## 14. Required Acceptance Tests

Required acceptance coverage; these are obligations, not existing passing fixtures:

1. Discover/list/select compatible workflow/job files; invalid YAML, unsafe paths,
   ambiguous jobs, unsupported selections and empty execution cannot PASS.
2. Deterministic owner-approved defaults/selection and explicit events, visible exclusions,
   supported dependencies/matrices and no silent event fallback.
   CLI-01 is unresolved; the historical preference is PROPOSED, not a test oracle.
3. Derive truthful local payloads without normal manual JSON; missing/invalid/conflicting
   inputs and dirty/HEAD mismatch are non-PASS or accurately represented.
4. Missing/unusable `act`, Docker or repository prerequisites cannot PASS;
   setup/help remains usable and init never changes hooks/Git configuration.
5. Failed jobs map to FAIL; mixed failure/incomplete work remains accurately reported;
   fake zero-exit or unavailable engine outcome evidence cannot fabricate PASS.
6. Crash/cancel/timeout/input mutation/uncertain cleanup/final-write failure never
   leave PASS; only run-owned resources may be cleaned.
7. Explicit secrets/variables are safely transported and redacted, never implicitly
   inherited or exposed through argv, source snapshots, logs or reports.
8. Concurrent repositories/worktrees/runs cannot overwrite or cross-associate
   source, payloads, reports or artifacts; report writes are atomic.
9. Unsupported required behavior, ambient config conflicts and identity uncertainty
   are visible non-PASS; partial-selection PASS never means hosted/full-repository PASS.
10. Record engine/config/source/action/image identities; prove cold acquisition,
    warm verified reuse, changed references, corrupt caches and interrupted/denied
    downloads. No unverified artifact execution or custom cache service.
11. Native Windows, Linux and macOS prebuilt CLI runs need no host Go/Node runtime;
    WSL2 is tested separately, with actual Docker/act and permission/cancellation evidence.
12. Go orchestration uses argv, stable tested exits/JSON, no hook/approval/background
    service and no automatic hosted suppression; each invocation requests a new run.

## 15. Deferred

Hosted suppression, additional OS/architectures, native Windows/macOS jobs,
alternative engines and advanced artifacts/cache administration are deferred.
Automatic Git gates and approval/bypass are removed, not deferred MVP dependencies.

## Dependency acquisition and reuse (2026-10-12)

Retained owner policy, aligned with PRD FR-3:
- Download missing approved artifacts when needed; normal use is not permanently offline.
- Resolve action refs to trusted full commits, verify cached Git objects/materialized
  content against those identities; directory names alone do not prove integrity.
- Resolve image refs to trusted platform-specific OCI digests and verify actual
  execution/reuse matches them. Mutable tags or local image IDs alone are insufficient.
- Treat changed identities as new dependencies; never silently substitute versions.
- Respect workflow package-manager lockfile/integrity and configured-cache boundaries.
- The adapter must establish trusted expected versus actual identities before job
  execution and prevent engine substitution/implicit pulls from bypassing verification.
  If the tested engine cannot expose/enforce this handoff, INCOMPLETE.
- Acquisition credentials/network policy are separate from runtime credentials/network.
  Workflow steps can need networking; prefetching does not make every workflow offline.
- Record non-secret provenance/outcomes; denied networking, corruption, interruptions
  or unverifiable dependencies are INCOMPLETE, never PASS.
- Offline/cache-only is optional and must fail on missing content. `act` offline mode
  can fetch missing artifacts; it is not a security boundary.

The implementation contract for configuration/trusted sources and verified handoff
belongs to **DEP-01/SEC-01**, resolved before that slice launches workflows.
References checked 2026-10-14: [act usage](https://nektosact.com/usage/index.html),
[limitations](https://nektosact.com/not_supported.html),
[GitHub workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax).
