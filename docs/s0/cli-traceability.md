# Revised CLI — Requirement-to-slice and acceptance map

**Status:** Planned evidence only; no runtime acceptance
**Authority:** [PRD](../prd.md), [execution acceptance](../specs/local-ci-execution/results-and-acceptance.md#14-required-acceptance-tests),
[revised slices](../implementation-plan/vertical-slices.md#4-dependency-ordered-vertical-slices)

## Product requirement dispositions

| Requirement | Revised implementation | Acceptance |
| --- | --- | --- |
| FR-1 | S1/S2/S3: optional setup without hooks | L-04, L-12 |
| FR-2 | S2: discovery/listing/selection | L-01, L-02 |
| FR-3 | S3/S4: fresh act/Docker adapter and verified cache reuse | L-04–L-06, L-10, L-12 |
| FR-4 | S2/S4: local events and explicit missing inputs | L-02, L-03 |
| FR-5 | S2/S5: tested-input/run identity, durable reports/JSON | L-03, L-06, L-08, L-10, L-12 |
| FR-6 | S2/S5: report accuracy; gating/PASS reuse removed | L-03, L-06, L-09, L-12 |
| FR-7 | S1/S4: ordinary Git remains independent | L-04, L-12 |
| FR-8 | REMOVED: no approval UI | Negative scope check L-12 |
| FR-9 | REMOVED: no bypass/audit provider | Negative scope check L-12 |
| FR-10 | DEFERRED: no hosted suppression in MVP | Negative scope check L-12 |
| FR-11 | S4/S5: safe diagnostics/progress/scope | L-05, L-09, L-12 |
| FR-12 | S3/S4/S5: explicit secrets/redaction | L-07 |
| FR-13 | S9/S10: supported host matrix and Linux jobs | L-09, L-11 |

## Local execution acceptance mapping

L-n identifiers refer to the twelve current numbered acceptance clauses.

| ID | Slices | Required evidence |
| --- | --- | --- |
| L-01 | S2 | Discovery, YAML/path errors, listing/selection, no empty PASS |
| L-02 | S2/S4 | Deterministic selection/defaults against owner-approved CLI-01; explicit events, dependency/matrix scope and exclusions; historical preference is not an oracle |
| L-03 | S2/S4/S5 | Truthful payloads/current dirty inputs and required-field errors |
| L-04 | S1/S3 | Dependency/setup failure, informational commands, no Git integration |
| L-05 | S4/S5 | Failed/mixed/uncertain execution and complete result aggregation |
| L-06 | S2/S4/S5 | Crash/cancel/mutation/write/cleanup failures and run ownership |
| L-07 | S3/S4/S5 | Secret canaries, explicit input, transport and redaction |
| L-08 | S2/S4/S5 | Parallel runs/worktrees, atomic reports and no cross-association |
| L-09 | S2/S3/S4/S5 | Unsupported/ambient conflicts, identity ambiguity, honest partial scope |
| L-10 | S3/S4/S5 | Cold/warm/cache tamper/ref change/download interruption and provenance |
| L-11 | S9/S10 | Real supported native/WSL2 host acceptance and prebuilt installation |
| L-12 | S1–S5/S10 | Go argv, exit/JSON consistency, fresh execution and removed-feature absence |

## Open decisions and historical evidence

CLI-01/SRC-01/DEP-01/SEC-01/PERF-01 are in the
[current decision register](../implementation-plan/authority-and-gates.md#current-cli-decisions).
Only their dependent implementation is blocked; unrelated documented work can proceed.
D4 governs release support. D1 is retired, D2 deferred, and D3 historical.

## PRD success measurements

[PRD §11](../prd/success-and-decisions.md#11-success-criteria) requires setup-step
count, orchestration overhead, cold/warm duration, disk use and downloaded bytes.
S9 owns the documented representative-workflow baseline and post-baseline budget
approval; S10 owns measured verification against every approved budget before release.
Methods, units, evidence and missing-data behavior are in the
[measurement protocol](../implementation-plan/verification-and-risks.md#success-measurement-and-budget-gate).
PERF-01 remains open; no measurements or approved numeric targets exist in this PR.

Old [G/A mapping](traceability.md) and [fixtures](fixtures.md) are historical.
Useful failure/privacy/concurrency test ideas may be adapted with a new current
oracle; old staged-gate and human-provider assertions are not carried forward.
Existing foundation CI is evidence for existing code, not for these planned fixtures.
