# Implementation planning — Revised authority and decisions

**Status:** Draft aligned with merged PRD #20; review pending
**Canonical authority:** [Implementation planning](../implementation-plan.md)

## 1. Objective and Authority

Deliver explicitly invoked Go local CI through act/Docker with accurate reports.
Authority: [PRD](../prd.md) → [architecture](../architecture.md) →
[execution contract](../specs/local-ci-execution.md) → this sequencing plan.
Retired gating/approval docs and historical S0/D3 procedures are not product authority.

## 2. Execution Rules

- Separate scoped implementation issue/feature branch and draft PR per slice;
  verify current main and preserve unrelated work.
- Reuse Go entrypoint/help/version; build fresh execution features.
- Resolve the relevant decision before implementing its dependent behavior.
- Run related deterministic verification; report missing evidence explicitly.
- Independent reviews and explicit owner merge authorization remain separate gates.
- No shell strings, global host changes, real credentials or live act in this docs task.
- No false PASS, implied platform acceptance or automatic production-readiness claim.

## 3. Decisions and Feasibility Gates

| Gate | Revised disposition | Effect |
| --- | --- | --- |
| D1 — Human approval trust | RETIRED with FR-8/FR-9 | Not an MVP dependency |
| D2 — Remote duplicate CI | DEFERRED with FR-10 | Not an MVP release dependency |
| D3 — Snapshot/event compatibility | HISTORICAL staged/push acceptance; PARTIAL evidence preserved | No staged-gate/PASS-reuse requirement; applicable current source/runtime safety is covered by SRC-01/SEC-01 |
| D4 — Platform/release matrix | OPEN | Confirm real supported hosts/architectures, versions, permissions, cancellation and installation before release |

### Current CLI decisions

| ID | Required contract | Blocking scope |
| --- | --- | --- |
| CLI-01 | Owner-approved default event/workflow/job selection, CLI flags/subcommands, input precedence, numeric exits, JSON/report schema and stdout/stderr | Dependent S2 selection behavior and new CLI/report interface; historical event preference is PROPOSED only |
| SRC-01 | Tested current-file capture, ignored/untracked policy, links/submodules, Git metadata and change detection | Source execution |
| DEP-01 | Trusted sources/ref resolution and enforcing verified act/Docker artifacts without engine substitution | Dependency execution handoff |
| SEC-01 | Effective config/env/credentials, acquisition/runtime networks, mounts/privilege, bounded processes and run-owned cleanup | Real workflow integration |
| PERF-01 | Representative-workflow measurement protocol, baseline evidence, post-baseline owner-approved numeric budgets and comparison rules for every PRD §11 metric | S9/S10 acceptance and release; no invented numbers or live measurements in this docs PR |

These are unresolved implementation contracts, not permissions to invent behavior.
An open decision blocks dependent work, not the existing foundation or documentation.

Normal product architecture may download missing dependencies and allow explicit
workflow networking under a reviewed policy. It does not require D3's signed VM
authorization machinery. This does not authorize running historical D3 recipes
or real act tests: those still need separately reviewed scope and explicit approval.

## D3 progress and S2 development boundary — 2026-10-12

**Historical disposition:** the earlier restricted S2 permission and staged/push
acceptance obligations are superseded for product sequencing by the revised slices.
D3 evidence remains PARTIAL; UNVERIFIED results do not become PASS.
Original wording is preserved in Git at baseline
`63ef75bf5a7c475675e0602a79df0294206d529a`.

## Owner-authorized PR #18 D3 preparation exception — 2026-10-13

**Historical disposition:** PR #18 authorized only offline D3 preparation.
It did not enable runtime execution. The existing validator code remains untouched;
it is separate test tooling, not normal CLI runtime authorization.
PR #19 stays unmerged pending a separate owner decision. This alignment neither
merges/closes it nor extends its approval scope.
