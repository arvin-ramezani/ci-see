# CLI contract alignment — Preservation and status record

**Date:** 2026-10-14
**Baseline:** `63ef75bf5a7c475675e0602a79df0294206d529a` (merged PRD #20)
**Status:** Documentation alignment only; independent review pending

## Authorization and boundaries

The owner requested architecture/spec/plan alignment after approving explicit
CLI scope in PR #20. No application code, tests, workflows, host configuration,
D3 evidence bytes or GitHub issue/PR status is modified.
The earlier implementation audit found CLI foundation and offline D3 validators,
not an implemented workflow runner. Runtime source remains unchanged.

## Semantic disposition map

| Baseline source/sections | Destination and disposition |
| --- | --- |
| Architecture §1–6 | Same child path: explicit CLI/act/Docker, current tested source and per-run reports replace hooks/gate state |
| Architecture §7–10 | Same path: independent Git, retired approval component; retain argv, truthful outcomes, credentials/cleanup and host distinctions |
| Architecture §11–13 | Same path: fresh workflow/engine/report responsibilities; retain consumer seams, boundary validation, Go and unresolved decisions |
| Execution §1–5 | Same path: discovery/selection retained; historical event preference is PROPOSED pending owner-approved CLI-01, not a default contract; classify missing/unsupported input |
| Execution §6–10 | Same path: current-file identity replaces index/push target; retain truthful payloads, concurrency, private transport, redaction and cancellation |
| Execution §11–15 and dated dependency policy | Same path: retain twelve acceptance IDs, verified download/reuse and no hosted parity; revise detailed oracles for explicit runs/reports |
| Plan §1–3 | Same path: D1 retired, D2 deferred, D3 historical, D4 retained; CLI-01/SRC-01/DEP-01/SEC-01 record unresolved dependent decisions |
| Plan §4 | Same path: S0/S1 retain foundation; S2–S5 fresh discovery/dependency/execution/report work; S6/S7 retired, S8 deferred, S9/S10 revised release/acceptance |
| Plan §5–7 | Same path: focused verification, independent/owner gates and host acceptance; retain PRD §11 measurement obligations through S9/S10/PERF-01; remove gate/provider/remote suppression prerequisites |
| Git gating/approval indexes and six children | Historical banner only; original bodies, IDs, headings and acceptance numbers unchanged |
| S0 historical matrix/fixtures/snapshot decisions/D3 procedures | Historical banner only; original bodies remain inspectable, not revised product obligations |
| S0 index/testing interfaces | Current routing/testing seams; preserve old headings for known links; replace gate/provider advice |
| New current traceability | `docs/s0/cli-traceability.md`: FR dispositions and L-01–L-12 current oracles/slices |
| README/AGENTS/docs index/PRD handoff | Route current contracts, historical records and pending cleanup correctly |
| D3 evidence report/manifest | Unchanged bytes/hashes and PARTIAL/UNVERIFIED labels |

The older contract wording remains recoverable at the baseline. This is an
authorized semantic revision, not a claim of unchanged historical meaning.
Deprecated paths and legacy fragment routes remain valid.

## Size decisions and checks

KEEP existing active topic boundaries: each is an independently useful retrieval unit.
KEEP the two long historical D3 catalogue/fixture documents as evidence references:
their cohesive recipes and old cleanup conditions are not mechanically split.
A historical banner does not require redesigning their code examples.
REORGANIZE authority through current versus historical routing; no file moves.

Required checks: repository doc validator (paths/fragments/IDs/legacy routes),
context scanner (reachability/size), diff whitespace, original-versus-result
semantic comparison, preserved historical-body checks and evidence-byte checks.
These cannot prove runtime parity, isolation, platform support or independent review.

Local verification: both documentation scanners and diff whitespace checks PASS;
46 indexed documentation files, no broken paths/fragments or unindexed files.
All 16 banner-only historical bodies match baseline exactly; executable code,
tooling/workflows and evidence JSON are byte-unchanged. All changed active prose
is at most 150 lines. Runtime/unit/E2E checks were not executed for this docs-only diff.

## Implementation handoff

Next is a separate scoped cleanup/command-contract PR, followed by fresh features
under the revised sequence. Do not adapt or ship old gate models as runner behavior.
Do not close PR #19, run act/Docker, modify firewall settings or claim production
acceptance as a side effect of documentation alignment.
