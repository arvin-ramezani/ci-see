> [Back to Implementation planning index](../implementation-plan.md)

<!-- BEGIN ORIGINAL CONTRACT -->
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

<!-- END ORIGINAL CONTRACT -->
