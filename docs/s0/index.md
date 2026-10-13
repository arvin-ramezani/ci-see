# S0 — Current CLI traceability and historical feasibility routing

**Status:** Revised planning index; not implementation evidence
**Authority:** [PRD](../prd.md), [architecture](../architecture.md),
[execution spec](../specs/local-ci-execution.md), [plan](../implementation-plan.md)

## Read by task

| Task | Read |
| --- | --- |
| Implement revised CLI requirements and acceptance | [Current traceability](cli-traceability.md) |
| Plan focused tests and package boundaries | [Testing and interfaces](testing-and-interfaces.md) |
| Review requirement retirement/preservation | [Alignment record](../migration/cli-contract-alignment.md) |
| Inspect previous requirement mapping | [Historical matrix](traceability.md) |
| Inspect previous staged/push design | [Historical decisions](snapshots-and-decisions.md) |
| Inspect old fixture descriptions | [Historical fixture catalog](fixtures.md) |
| Inspect recorded restricted results | [D3 feasibility record](d3-feasibility-record.md) |
| Inspect historical D3 procedures/tooling | [Phase A plan](d3-phase-a-verification-plan.md), [runtime catalogue](d3-phase-a-runtime-catalogue.md), [Git fixtures](d3-phase-a-git-fixtures.md), [offline preflight](d3-phase-b-preflight.md) |

## Contract preservation and interpretation

Current implementation follows revised requirements, not historical gate/provider
clauses. G-01–G-20 and A-01–A-23 remain identifiable historical obligations;
they are not revised MVP acceptance requirements.
FR-8/FR-9 and gates are removed; FR-10 is deferred.

Keep historical result labels, evidence JSON and hashes unchanged. D3 remains
PARTIAL and unverified records remain unverified. Nothing here proves runtime
networking, credentials, event fidelity or the new CLI works.

## S0 completion and handoff

Validate links/anchors, requirement IDs, current mappings, retirement status and
preservation evidence. Independent review and owner-approved merge remain pending.
Next implement the scoped foundation cleanup and resolve the relevant new contract
decisions before fresh runner work. Fake tools may be used without live act.

## S0 evidence boundary

This is a planning map, not code or passing fixture evidence. Existing PR #15 is
CLI-foundation evidence only; PR #18 is offline D3 tooling; PR #19 is still unmerged.
This PR does not close issues/PRs or authorize live tests.
