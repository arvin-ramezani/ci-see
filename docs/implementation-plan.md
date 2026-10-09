# CI See — Go MVP Implementation Plan

**Status:** Draft — requires independent plan review and owner approval  
**Product:** CI See  
**Canonical path:** `docs/implementation-plan.md`  
**Planning issue:** #7  
**Last updated:** 2026-10-09

## 1. Objective and Authority

Deliver the first usable, local-first Go CLI that runs existing GitHub Actions via user-installed `act`, automatically gates normal Git commits/pushes on exact-state CI results, and allows only a demonstrably human-approved, durably audited bypass.

Authority, in order: `docs/prd.md` (scope and outcomes), `docs/architecture.md` (Go/system boundaries), `docs/specs/local-ci-execution.md`, `docs/specs/git-gating.md`, and `docs/specs/developer-approval-ux.md` (testable behavior). This plan schedules work; it does **not** override or silently relax a canonical contract. Resolve contradictions through reviewed spec/ADR changes before code.

**Non-goals:** rebuilding GitHub Actions, embedding `act`, promising hosted-runner parity, a daemon/SaaS, or claiming native Git hooks are tamper-proof. Git `--no-verify` and full host compromise remain outside the enforcement boundary.

## 2. Execution Rules

- One narrowly scoped implementation issue and PR per slice. Branch from verified current `main`; avoid concurrent overlapping changes.
- Before each slice, identify the authoritative acceptance clauses and missing decisions. Request owner/ADR decision rather than inventing security or product behavior.
- First obtain independent review **PASS** on this documentation-only plan's exact PR SHA. No application code in planning PR.
- Each implementation PR: tests + deterministic verification, independent code review, security review for trust boundaries, and owner approval before merge. A reviewer cannot self-certify an unverified gate.
- Gates are evidence-driven: report exact head SHA, checks and limitations. Passing unit tests alone never means release-ready.
- Prefer Go standard library, idiomatic `cmd/ci-see` and `internal/{cli,core,adapters/{git,act,state,approval},hooks}`. Keep `core` independent of platform/process adapters; never use shell-string execution.
- Run `gofmt`, `go test ./...`, `go vet ./...`, focused integration tests, and race checks where supported. Establish coverage reporting/threshold in the test plan before enforcing a numeric target; never claim nonexistent coverage.
- At each slice boundary, earlier gates continue to fail closed. No partial implementation may allow an unknown state or missing dependency to become PASS.

## 3. Decisions and Feasibility Gates

| Gate | Decision / required evidence | Deadline |
| --- | --- | --- |
| D1 — Human approval trust | Threat-model Windows 11 + WSL2 and supported native providers; authenticate approval response, bind to one request/state, prevent caller-controlled flags/env/stdin/files/stdout/exit-code spoofing, test UI automation/IPC attacks. If human-only boundary is unproven, disable Continue anyway and **do not claim complete MVP**. | Prototype/ADR early; before S7 |
| D2 — Remote duplicate CI | PRD FR-10 requires a safe, explicit mechanism for avoiding unnecessary duplicate GitHub-hosted work without invisibly weakening remote validation. Select/review mechanism, permissions, failure fallback and acceptance tests in a separate spec/ADR; do not silently drop this requirement. | Before S8 |
| D3 — Snapshot/event compatibility | Decide isolated execution materialization for index tree vs commit tip, worktrees/submodules, generated PR/push event context, and honest unsupported cases. No unstaged/untracked leakage or guessed required event inputs. | Before S2/S4 |
| D4 — Platform/release matrix | Confirm initial supported `GOOS/GOARCH`, executable naming, minimum Git/`act`/Docker conditions, process-tree cancellation strategy, and clean-host acceptance targets. Windows native and WSL2 are separate hosts. | Before S9; relevant platform adapters earlier |

Feasibility investigations may run in a **separate, explicitly approved spike issue** before feature work, but must not be merged as unreviewed production bypass behavior. An unresolved gate blocks its dependent slice/release, not unrelated earlier slices.

## 4. Dependency-Ordered Vertical Slices

| Slice | Deliverable and boundaries | Acceptance evidence |
| --- | --- | --- |
| **S0 — Contracts & test harness** | Finalize blocking D3 assumptions; create a deterministic fixture plan for synthetic Git repositories, fake `act`, temporary Git metadata/worktrees, process failure and headless conditions. Document all undecided scope explicitly. | Test matrix maps PRD FR-1–FR-13 and existing spec acceptance clauses; no unsupported decision treated as resolved. |
| **S1 — Go CLI foundation** | Go module, pinned build toolchain for implementation, `cmd/ci-see`, typed core outcomes, `ci-see`/`init`/`status` skeleton, platform/process adapter interfaces, structured errors/exit meanings. No Node.js application runtime. | Build/tests/vet on chosen hosts; unknown commands/dependency errors never falsely PASS; packages obey dependency boundary. |
| **S2 — Repository, discovery & isolated snapshots** | Resolve Git/worktree paths via Git, discover YAML workflows, select supported events deterministically, construct temporary `push`/`pull_request`/`workflow_dispatch` context and **exact index-tree or commit-tip execution snapshots** as appropriate. | Tests for staged vs unstaged/untracked files, detached HEAD, worktrees, missing branch/base context, ambiguous/required inputs, unsupported workflows, cleanup and path safety. |
| **S3 — Local state & fingerprints** | Git-metadata state adapter, run IDs, RUNNING/PASS/FAIL/INCOMPLETE, stable schema/version, exact fingerprint (tree, repo/worktree, workflows, events, config, engine/image, variables and non-plaintext secret digest), atomic durable writes, serialization/lock recovery. | Fault injection for crashes, disk full, partial/corrupt files, concurrency and staleness; no cross-worktree reuse, false PASS, or plaintext secrets. |
| **S4 — Local execution via `act`** | Dependency preflight; Go `os/exec` argument arrays; bounded context/process-tree cancellation; execute selected workflows against the isolated snapshot; normalize per-job and overall results; sanitize logs, temp secrets/vars and generated payloads. | Fake-`act` then real-`act` fixture tests: exit/nonzero/crash/cancel/timeouts, version detection, missing Docker/`act`, unsupported cases, multiple workflows, cleanup and no leaked secrets. |
| **S5 — Git hooks and automatic gating** | Idempotent `ci-see init` (commit/push/both; default both); preserve existing hooks/`core.hooksPath` arrangements; pre-commit exact index gate; safe post-commit result binding; pre-push per-new-tip multi-ref gate; execute local CI automatically when no reusable PASS. | Real temp-repo tests for all modes, mode changes, no overwrite, changed staged tree, equivalent commit reuse, multi-ref push, deletion ref, stale PASS, CI crash, missing CLI, and stable allow/block exit behavior. |
| **S6 — Approval-provider feasibility** | Security/UX spike per D1: platform-specific trusted approval source and reply channel, WSL2-to-Windows route, native Windows/Linux/macOS approach; state capability limits and adversarial test strategy. Separate from allowing bypass in code. | Independent security review records what the provider **can/cannot** resist. Unverified provider means no Continue anyway, with a fail-closed UX and MVP release blocker. |
| **S7 — Developer decision and durable bypass audit** | Implement only reviewed provider(s). Show Cancel/Continue with operation, repo, state and reason. Authenticate a request-specific, expiring, single-use response; recheck state; persist, sync and verify audit **before** ALLOW_BYPASS. Headless/closed/crashed/expired = block. | Tests for caller-controlled input spoofing, replay/races, UI automation limits, cancellation, disk full/denied audit, changed fingerprint while UI open, noninteractive Git/IDE/agent, and response semantics. No reusable PASS or token. |
| **S8 — Hosted duplicate-execution control** | Implement the reviewed D2 mechanism only. Preserve visible hosted security checks; explicit opt-in/config where required; safe fallback if GitHub API/config/permissions do not permit suppression. | Representative GitHub workflow/repository tests prove no accidental skip of required remote checks, no redundant work in eligible cases, and no silent bypass on failure. |
| **S9 — Release, install and operability** | Publish signed or checksum-verifiable (policy to decide) per-OS/architecture native binaries; pin Go/tool dependencies; document installation, `PATH`, Git-hook lifecycle, upgrade and safe uninstall/rollback. No requirement for Go/Node/npm on hosts using prebuilt CLI. | Fresh-host smoke tests for supported OS/architectures, especially native Windows Git vs Git in WSL2; checksums, reproducible version metadata, missing/old `act`, Docker unavailable, hook restoration and uninstall. |
| **S10 — Full MVP acceptance** | Exercise actual `change → Git commit → local act CI → exact PASS/non-PASS → human decision → push → PR` flow, supported platform matrix, documented parity/security limitations and diagnostic UX. | Entire acceptance matrix PASS on real Git plus `act`/Docker, verified provider/security tests, hosted duplicate strategy, release artifacts and documented recovery. Then independent release/security review + explicit owner go/no-go. |

**Sequence:** S0 → S1 → S2 → S3 → S4 → S5 → S7 → S8 → S9 → S10. S6 security feasibility begins alongside S0/S1 and must PASS its intended security boundary before S7. D2 analysis can begin in parallel but blocks S8 and MVP acceptance.

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
