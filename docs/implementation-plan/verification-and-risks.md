# Implementation planning — Test matrix, review/merge gates, risks

**Status:** Draft (inherited)  
**Canonical authority:** [Implementation planning](../implementation-plan.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 5. Required Test Matrix

| Category | Minimum fixtures and assertions |
| --- | --- |
| Core/invariants | Typed states; matching PASS only; unknown, missing, stale, corrupt, interrupted, unsupported-required validation never PASS; stable exit codes. |
| Git identity | Index-only staged snapshots; worktree and detached cases; post-commit tree equivalence; multi-ref pre-push; relevant workflow/event/config/secret changes invalidate PASS. |
| Files and processes | Concurrent runs, lock contention/stale lock, partial/crashed writes, temp cleanup, spawned-child termination, bounded headless timeouts, no secret logs/metadata. |
| Approval security | Noninteractive agent/IDE Git; forged/replayed responses, changed state, agent-controlled environment/stdin/files, UI automation/IPC exposures, missing provider, mandatory persisted audit before ALLOW_BYPASS. |
| Hosts and releases | Windows native vs WSL2 separation, Linux/macOS supported targets, Git hook shell/path differences, binaries with Go/Node/npm absent, required external `act` and container runtime. |
| Remote CI | FR-10 approved design's exact-state and no-weakened-required-checks tests plus graceful fallback when duplicate suppression is not possible. |

For each slice, check the **full existing** acceptance clauses in its owning canonical spec; the matrix above is an index, not a replacement. Treat both mocked tests and real-tool integration as necessary, with clear environment prerequisites and no invented green checks.

## 6. Review, Merge and Release Gates

1. **Plan gate:** independent architecture/security/contract plan reviewer audits this exact PR head against the five canonical docs; returns PASS or CHANGES REQUIRED. Owner approves merge. This plan PR adds no application code.
2. **Slice planning gate:** create a scoped issue (authority links, inputs, out-of-scope, risks, tests, expected files). Independently review significant architecture/security/UX plans before implementation.
3. **Implementation gate:** Codex/executor implements on a fresh issue branch; `gofmt`, `go test ./...`, `go vet ./...`, targeted integration, changed-file coverage once an agreed threshold exists, and diff checks. Record exact SHA and environment.
4. **Independent gates:** code review for every PR; security review on Git/process/IPC/secrets/audit changes; UX validation on approval behavior. Fix findings and re-review the new exact head. **Never merge without explicit owner authorization.**
5. **MVP gate:** S0–S10 accepted; D1–D4 resolved; cross-platform/real-tool E2E and release install/uninstall PASS; FR-10 explicitly covered; residual threats documented. `PRODUCTION_READY = NO` until this gate is met.

## 7. Known Risks, Dependencies and Next Action

- Local `act` PASS is not hosted parity; unsupported GitHub Actions capabilities must remain visible.
- Local Git hooks are bypassable outside CI See; do not promise central enforcement.
- Same-user desktop automation may defeat naive approval windows; **never** call an ordinary popup proof of human approval.
- Runner images, Docker availability, filesystem locks, and Windows/WSL2 process behavior can change tests and must be validated on real hosts.
- CI See should not add a JavaScript runtime merely to implement its own UI/orchestration; independently executed workflow actions may use Node.js inside their execution environment.
- FR-10 duplicate-hosted-CI behavior is unresolved and cannot be marked complete until D2 is accepted.

**Next:** independent review this plan at its exact PR SHA. After owner-approved merge, open a separate S0 execution issue; no implementation or bypass enablement in this planning PR.
<!-- END ORIGINAL CONTRACT -->
